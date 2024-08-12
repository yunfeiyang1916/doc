package https

import (
	"bytes"
	"crypto/cipher"
	"crypto/subtle"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// 返回tls服务端连接
func Server(conn net.Conn, config *Config) *Conn {
	return &Conn{conn: conn, config: config}
}

// 返回tls客户端连接
func Client(conn net.Conn, config *Config) *Conn {
	return &Conn{conn: conn, config: config, isClient: true}
}

// tls连接，实现了net.Conn接口
type Conn struct {
	conn net.Conn
	// 是否为客户端
	isClient bool

	// handshakeMutex < in.Mutex, out.Mutex, errMutex
	handshakeMutex sync.Mutex
	// 握手导致的错误
	handshakeErr error
	// tls version
	vers uint16
	// 版本已协商
	haveVers bool
	config   *Config
	// 是否握手完成
	handshakeComplete bool
	// 密码套件
	cipherSuite uint16
	// 用于存储我们构建的证书链，而不是服务器提供的证书链
	verifiedChains [][]*x509.Certificate
	// 用于存储在握手过程中发送的第一个Finished哈希值。这个值被用作“tls-unique”通道绑定值。
	firstFinished [12]byte

	// 表示发送或者接收方向的连接
	in, out halfConn
	// 原始输入数据，直接从网络接收
	rawInput *block
	// 应用程序数据，等待被读取
	input *block
	// 握手数据，等待被读取
	hand bytes.Buffer

	// 最低位表示Close方法是否被调用过，其余位表示Conn.Write方法中运行的goroutine数量。
	activeCall int32

	tmp [16]byte
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated with the connection.
// A zero value for t means Read and Write will not time out.
// After a Write has timed out, the TLS state is corrupt and all future writes will return the same error.
func (c *Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadDeadline sets the read deadline on the underlying connection.
// A zero value for t means Read will not time out.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the write deadline on the underlying connection.
// A zero value for t means Write will not time out.
// After a Write has timed out, the TLS state is corrupt and all future writes will return the same error.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// sendAlert sends a TLS alert message.
// c.out.Mutex <= L.
func (c *Conn) sendAlertLocked(err alert) error {
	switch err {
	case alertNoRenegotiation, alertCloseNotify:
		c.tmp[0] = alertLevelWarning
	default:
		c.tmp[0] = alertLevelError
	}
	c.tmp[1] = byte(err)
	c.writeRecord(recordTypeAlert, c.tmp[0:2])
	// 如果错误类型不是closeNotify，则将错误信息设置为local error并返回
	if err != alertCloseNotify {
		return c.out.setErrorLocked(&net.OpError{Op: "local error", Err: err})
	}
	return nil
}

// 发送tls警告信息
func (c *Conn) sendAlert(err alert) error {
	c.out.Lock()
	defer c.out.Unlock()
	return c.sendAlertLocked(err)
}

func (c *Conn) writeRecord(typ recordType, data []byte) (n int, err error) {
	b := c.out.newBlock()
	for len(data) > 0 {
		m := len(data)
		if m > maxPlaintext {
			m = maxPlaintext
		}
		explicitIVLen := 0
		explicitIVIsSeq := false

		if cbc, ok := c.out.cipher.(cbcMode); ok {
			explicitIVLen = cbc.BlockSize()
		}
		if explicitIVLen == 0 {
			if _, ok := c.out.cipher.(cipher.AEAD); ok {
				explicitIVLen = 8
				// The AES-GCM construction in TLS has an
				// explicit nonce so that the nonce can be
				// random. However, the nonce is only 8 bytes
				// which is too small for a secure, random
				// nonce. Therefore we use the sequence number
				// as the nonce.
				explicitIVIsSeq = true
			}
		}
		b.resize(recordHeaderLen + explicitIVLen + m)
		b.data[0] = byte(typ)
		vers := c.vers
		if vers == 0 {
			// Some TLS servers fail if the record version is
			// greater than TLS 1.0 for the initial ClientHello.
			vers = VersionTLS12
		}
		b.data[1] = byte(vers >> 8)
		b.data[2] = byte(vers)
		b.data[3] = byte(m >> 8)
		b.data[4] = byte(m)
		if explicitIVLen > 0 {
			explicitIV := b.data[recordHeaderLen : recordHeaderLen+explicitIVLen]
			if explicitIVIsSeq {
				copy(explicitIV, c.out.seq[:])
			} else {
				if _, err = io.ReadFull(c.config.rand(), explicitIV); err != nil {
					break
				}
			}
		}
		copy(b.data[recordHeaderLen+explicitIVLen:], data)
		c.out.encrypt(b, explicitIVLen)
		_, err = c.conn.Write(b.data)
		if err != nil {
			break
		}
		n += m
		data = data[m:]
	}
	c.out.freeBlock(b)
	// 如果写入的是ChangeCipherSpec记录类型，则在写入后更新加密状态。
	if typ == recordTypeChangeCipherSpec {
		err = c.out.changeCipherSpec()
		if err != nil {
			// Cannot call sendAlert directly,
			// because we already hold c.out.Mutex.
			c.tmp[0] = alertLevelError
			c.tmp[1] = byte(err.(alert))
			c.writeRecord(recordTypeAlert, c.tmp[0:2])
			return n, c.out.setErrorLocked(&net.OpError{Op: "local error", Err: err})
		}
	}
	return
}

// 用于表示TLS记录头无效时产生的错误
type RecordHeaderError struct {
	// 是一个人类可读的字符串，用于描述错误的详细信息
	Msg string
	// 是一个长度为5字节的数组，用于存储触发错误的TLS记录头的前五个字节
	RecordHeader [5]byte
}

func (e RecordHeaderError) Error() string { return "tls: " + e.Msg }

func (c *Conn) newRecordHeaderError(msg string) (err RecordHeaderError) {
	err.Msg = msg
	copy(err.RecordHeader[:], c.rawInput.data)
	return err
}

// 该函数用于从连接中读取下一个记录，并更新记录状态。
// c.in.Mutex <= L; c.input == nil.
func (c *Conn) readRecord(want recordType) error {
	// Caller must be in sync with connection:
	// handshake data if handshake not yet completed,
	// else application data.  (We don't support renegotiation.)
	switch want {
	default:
		c.sendAlert(alertInternalError)
		return c.in.setErrorLocked(errors.New("tls: unknown record type requested"))
	case recordTypeHandshake, recordTypeChangeCipherSpec:
		if c.handshakeComplete {
			c.sendAlert(alertInternalError)
			return c.in.setErrorLocked(errors.New("tls: handshake or ChangeCipherSpec requested after handshake complete"))
		}
	case recordTypeApplicationData:
		if !c.handshakeComplete {
			c.sendAlert(alertInternalError)
			return c.in.setErrorLocked(errors.New("tls: application data record requested before handshake complete"))
		}
	}

Again:
	if c.rawInput == nil {
		c.rawInput = c.in.newBlock()
	}
	b := c.rawInput

	// Read header, payload.
	if err := b.readFromUntil(c.conn, recordHeaderLen); err != nil {
		// RFC suggests that EOF without an alertCloseNotify is
		// an error, but popular web sites seem to do this,
		// so we can't make it an error.
		// if err == io.EOF {
		// 	err = io.ErrUnexpectedEOF
		// }
		if e, ok := err.(net.Error); !ok || !e.Temporary() {
			c.in.setErrorLocked(err)
		}
		return err
	}
	typ := recordType(b.data[0])

	// No valid TLS record has a type of 0x80, however SSLv2 handshakes
	// start with a uint16 length where the MSB is set and the first record
	// is always < 256 bytes long. Therefore typ == 0x80 strongly suggests
	// an SSLv2 client.
	if want == recordTypeHandshake && typ == 0x80 {
		c.sendAlert(alertProtocolVersion)
		return c.in.setErrorLocked(c.newRecordHeaderError("unsupported SSLv2 handshake received"))
	}

	vers := uint16(b.data[1])<<8 | uint16(b.data[2])
	n := int(b.data[3])<<8 | int(b.data[4])
	if c.haveVers && vers != c.vers {
		c.sendAlert(alertProtocolVersion)
		msg := fmt.Sprintf("received record with version %x when expecting version %x", vers, c.vers)
		return c.in.setErrorLocked(c.newRecordHeaderError(msg))
	}
	if n > maxCiphertext {
		c.sendAlert(alertRecordOverflow)
		msg := fmt.Sprintf("oversized record received with length %d", n)
		return c.in.setErrorLocked(c.newRecordHeaderError(msg))
	}
	if !c.haveVers {
		// First message, be extra suspicious: this might not be a TLS
		// client. Bail out before reading a full 'body', if possible.
		// The current max version is 3.3 so if the version is >= 16.0,
		// it's probably not real.
		if (typ != recordTypeAlert && typ != want) || vers >= 0x1000 {
			c.sendAlert(alertUnexpectedMessage)
			return c.in.setErrorLocked(c.newRecordHeaderError("first record does not look like a TLS handshake"))
		}
	}

	if err := b.readFromUntil(c.conn, recordHeaderLen+n); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		if e, ok := err.(net.Error); !ok || !e.Temporary() {
			c.in.setErrorLocked(err)
		}
		return err
	}

	// Process message.
	b, c.rawInput = c.in.splitBlock(b, recordHeaderLen+n)

	// mac validation from record data
	ok, off, err := c.in.decrypt(b)
	if !ok {
		c.in.setErrorLocked(c.sendAlert(err))
	}
	b.off = off
	data := b.data[b.off:]
	if len(data) > maxPlaintext {
		err := c.sendAlert(alertRecordOverflow)
		c.in.freeBlock(b)
		return c.in.setErrorLocked(err)
	}

	switch typ {
	default:
		c.in.setErrorLocked(c.sendAlert(alertUnexpectedMessage))

	case recordTypeAlert:
		if len(data) != 2 {
			c.in.setErrorLocked(c.sendAlert(alertUnexpectedMessage))
			break
		}
		if alert(data[1]) == alertCloseNotify {
			c.in.setErrorLocked(io.EOF)
			break
		}
		switch data[0] {
		case alertLevelWarning:
			// drop on the floor
			c.in.freeBlock(b)
			goto Again
		case alertLevelError:
			c.in.setErrorLocked(&net.OpError{Op: "remote error", Err: alert(data[1])})
		default:
			c.in.setErrorLocked(c.sendAlert(alertUnexpectedMessage))
		}

	case recordTypeChangeCipherSpec:
		if typ != want || len(data) != 1 || data[0] != 1 {
			c.in.setErrorLocked(c.sendAlert(alertUnexpectedMessage))
			break
		}
		err := c.in.changeCipherSpec()
		if err != nil {
			c.in.setErrorLocked(c.sendAlert(err.(alert)))
		}

	case recordTypeApplicationData:
		if typ != want {
			c.in.setErrorLocked(c.sendAlert(alertUnexpectedMessage))
			break
		}
		c.input = b
		b = nil

	case recordTypeHandshake:
		// TODO(rsc): Should at least pick off connection close.
		if typ != want {
			return c.in.setErrorLocked(c.sendAlert(alertNoRenegotiation))
		}
		c.hand.Write(data)
	}

	if b != nil {
		c.in.freeBlock(b)
	}
	return c.in.err
}

// 读取下一个握手消息
// c.in.Mutex < L; c.out.Mutex < L.
func (c *Conn) readHandshake() (interface{}, error) {
	// 先检查c.hand的长度是否小于4，如果是，则会通过调用readRecord函数读取记录，直到长度大于等于4
	for c.hand.Len() < 4 {
		if err := c.in.err; err != nil {
			return nil, err
		}
		if err := c.readRecord(recordTypeHandshake); err != nil {
			return nil, err
		}
	}

	data := c.hand.Bytes()
	n := int(data[1])<<16 | int(data[2])<<8 | int(data[3])
	if n > maxHandshake {
		return nil, c.in.setErrorLocked(c.sendAlert(alertInternalError))
	}
	for c.hand.Len() < 4+n {
		if err := c.in.err; err != nil {
			return nil, err
		}
		if err := c.readRecord(recordTypeHandshake); err != nil {
			return nil, err
		}
	}
	data = c.hand.Next(4 + n)
	var m handshakeMessage
	switch data[0] {
	case typeClientHello:
		m = new(clientHelloMsg)
	case typeServerHello:
		m = new(serverHelloMsg)
	case typeCertificate:
		m = new(certificateMsg)
	case typeServerHelloDone:
		m = new(serverHelloDoneMsg)
	case typeClientKeyExchange:
		m = new(clientKeyExchangeMsg)
	case typeFinished:
		m = new(finishedMsg)
	default:
		return nil, c.in.setErrorLocked(c.sendAlert(alertUnexpectedMessage))
	}

	// The handshake message unmarshallers
	// expect to be able to keep references to data,
	// so pass in a fresh copy that won't be overwritten.
	data = append([]byte(nil), data...)

	if !m.unmarshal(data) {
		return nil, c.in.setErrorLocked(c.sendAlert(alertUnexpectedMessage))
	}
	return m, nil
}

var errClosed = errors.New("crypto/tls: use of closed connection")

// Read can be made to time out and return a net.Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(b []byte) (n int, err error) {
	if err = c.Handshake(); err != nil {
		return
	}
	if len(b) == 0 {
		// Put this after Handshake, in case people were calling
		// Read(nil) for the side effect of the Handshake.
		return
	}

	c.in.Lock()
	defer c.in.Unlock()

	// Some OpenSSL servers send empty records in order to randomize the
	// CBC IV. So this loop ignores a limited number of empty records.
	const maxConsecutiveEmptyRecords = 100
	for emptyRecordCount := 0; emptyRecordCount <= maxConsecutiveEmptyRecords; emptyRecordCount++ {
		for c.input == nil && c.in.err == nil {
			if err := c.readRecord(recordTypeApplicationData); err != nil {
				// Soft error, like EAGAIN
				return 0, err
			}
		}
		if err := c.in.err; err != nil {
			return 0, err
		}

		n, err = c.input.Read(b)
		if c.input.off >= len(c.input.data) {
			c.in.freeBlock(c.input)
			c.input = nil
		}

		// If a close-notify alert is waiting, read it so that
		// we can return (n, EOF) instead of (n, nil), to signal
		// to the HTTP response reading goroutine that the
		// connection is now closed. This eliminates a race
		// where the HTTP response reading goroutine would
		// otherwise not observe the EOF until its next read,
		// by which time a client goroutine might have already
		// tried to reuse the HTTP connection for a new
		// request.
		// See https://codereview.appspot.com/76400046
		// and https://golang.org/issue/3514
		if ri := c.rawInput; ri != nil &&
			n != 0 && err == nil &&
			c.input == nil && len(ri.data) > 0 && recordType(ri.data[0]) == recordTypeAlert {
			if recErr := c.readRecord(recordTypeApplicationData); recErr != nil {
				err = recErr // will be io.EOF on closeNotify
			}
		}

		if n != 0 || err != nil {
			return n, err
		}
	}

	return 0, io.ErrNoProgress
}

// Close closes the connection.
func (c *Conn) Close() error {
	// Interlock with Conn.Write above.
	var x int32
	for {
		x = atomic.LoadInt32(&c.activeCall)
		if x&1 != 0 {
			return errClosed
		}
		if atomic.CompareAndSwapInt32(&c.activeCall, x, x|1) {
			break
		}
	}
	if x != 0 {
		// io.Writer and io.Closer should not be used concurrently.
		// If Close is called while a Write is currently in-flight,
		// interpret that as a sign that this Close is really just
		// being used to break the Write and/or clean up resources and
		// avoid sending the alertCloseNotify, which may block
		// waiting on handshakeMutex or the c.out mutex.
		return c.conn.Close()
	}

	var alertErr error

	c.handshakeMutex.Lock()
	defer c.handshakeMutex.Unlock()
	if c.handshakeComplete {
		alertErr = c.sendAlert(alertCloseNotify)
	}

	if err := c.conn.Close(); err != nil {
		return err
	}
	return alertErr
}

// 向连接写入数据
func (c *Conn) Write(b []byte) (int, error) {
	for {
		// 首先通过原子操作检查连接是否已关闭，如果已关闭则返回错误
		x := atomic.LoadInt32(&c.activeCall)
		if x&1 != 0 {
			return 0, errClosed
		}
		if atomic.CompareAndSwapInt32(&c.activeCall, x, x+2) {
			defer atomic.AddInt32(&c.activeCall, -2)
			break
		}
	}
	// 进行握手
	if err := c.Handshake(); err != nil {
		return 0, err
	}

	c.out.Lock()
	defer c.out.Unlock()

	if err := c.out.err; err != nil {
		return 0, err
	}

	if !c.handshakeComplete {
		return 0, alertInternalError
	}

	n, err := c.writeRecord(recordTypeApplicationData, b)
	return n, c.out.setErrorLocked(err)
}

// 用于执行客户端或服务端握手协议。如果握手协议尚未执行，则会执行相应的握手流程。
// 大多数情况下，无需显式调用Handshake方法，因为首次读写操作会自动调用它。
// 该方法首先通过互斥锁确保并发安全，然后检查握手是否已完成或存在错误，如果是，则直接返回。
func (c *Conn) Handshake() error {
	c.handshakeMutex.Lock()
	defer c.handshakeMutex.Unlock()
	if err := c.handshakeErr; err != nil {
		return err
	}
	if c.handshakeComplete {
		return nil
	}

	if c.isClient {
		c.handshakeErr = c.clientHandshake()
	} else {
		c.handshakeErr = c.serverHandshake()
	}
	return c.handshakeErr
}

// 表示发送或者接收方向的连接
type halfConn struct {
	sync.Mutex

	err error
	// 协议版本
	version uint16
	// 加密算法
	cipher interface{}
	mac    macFunction
	// 用于表示64位的序列号
	seq [8]byte
	// 用于表示一系列可用的块
	bfree *block
	// 用于避免分配内存；作为接口方法参数时，可以避免对象逃逸
	additionalData [13]byte

	// 下一个加密状态
	nextCipher interface{} // next encryption state
	// 下一个MAC算法
	nextMac macFunction // next MAC algorithm

	// 用于在计算MAC时避免为每个MAC分配新的缓冲区
	inDigestBuf, outDigestBuf []byte // used to save allocating a new buffer for each MAC.
}

func (hc *halfConn) setErrorLocked(err error) error {
	hc.err = err
	return err
}

func (hc *halfConn) error() error {
	hc.Lock()
	err := hc.err
	hc.Unlock()
	return err
}

// 用于分配一个新的block，如果可能的话，从hc的空闲列表中分配
func (hc *halfConn) newBlock() *block {
	b := hc.bfree
	if b == nil {
		return new(block)
	}
	hc.bfree = b.link
	b.link = nil
	b.resize(0)
	return b
}

// 将一个block对象放回到hc的空闲列表中。
// 函数通过将传入的b对象的link字段设置为hc.bfree，并将hc.bfree设置为b来实现。由于每侧只有一个或两个块在空闲列表上，因此不需要担心列表的修剪等问题。
func (hc *halfConn) freeBlock(b *block) {
	b.link = hc.bfree
	hc.bfree = b
}

// 该函数用于将一个block在第n个字节之后进行分割，返回一个包含前n个字节的block和一个包含剩余字节的block。
// 如果剩余字节为空，则返回nil。
func (hc *halfConn) splitBlock(b *block, n int) (*block, *block) {
	if len(b.data) <= n {
		return b, nil
	}
	bb := hc.newBlock()
	bb.resize(len(b.data) - n)
	copy(bb.data, b.data[n:])
	b.data = b.data[0:n]
	return b, bb
}

// padToBlockSize calculates the needed padding block, if any, for a payload.
// On exit, prefix aliases payload and extends to the end of the last full
// block of payload. finalBlock is a fresh slice which contains the contents of
// any suffix of payload as well as the needed padding to make finalBlock a
// full block.
// 用于计算给定负载需要的填充块。
// 函数接受两个参数，payload为要填充的数据，blockSize为块大小。函数返回两个切片，prefix为不包含填充的数据部分，finalBlock为包含填充的数据部分。
func padToBlockSize(payload []byte, blockSize int) (prefix, finalBlock []byte) {
	overrun := len(payload) % blockSize
	paddingLen := blockSize - overrun
	prefix = payload[:len(payload)-overrun]
	finalBlock = make([]byte, blockSize)
	copy(finalBlock, payload[len(payload)-overrun:])
	for i := overrun; i < blockSize; i++ {
		finalBlock[i] = byte(paddingLen - 1)
	}
	return
}

// 该函数用于移除给定负载（payload）中的填充（padding），并以常量时间返回未填充的切片和一个字节，该字节在填充有效时等于255，在填充无效时等于0。
// 函数遵循RFC 2246中6.2.3.2节的规定。
func removePadding(payload []byte) ([]byte, byte) {
	if len(payload) < 1 {
		return payload, 0
	}

	paddingLen := payload[len(payload)-1]
	t := uint(len(payload)-1) - uint(paddingLen)
	// if len(payload) >= (paddingLen - 1) then the MSB of t is zero
	good := byte(int32(^t) >> 31)

	toCheck := 255 // the maximum possible padding length
	// The length of the padded data is public, so we can use an if here
	if toCheck+1 > len(payload) {
		toCheck = len(payload) - 1
	}

	for i := 0; i < toCheck; i++ {
		t := uint(paddingLen) - uint(i)
		// if i <= paddingLen then the MSB of t is zero
		mask := byte(int32(^t) >> 31)
		b := payload[len(payload)-1-i]
		good &^= mask&paddingLen ^ mask&b
	}

	// We AND together the bits of good and replicate the result across
	// all the bits.
	good &= good << 4
	good &= good << 2
	good &= good << 1
	good = uint8(int8(good) >> 7)

	toRemove := good&paddingLen + 1
	return payload[:len(payload)-int(toRemove)], good
}

// 该函数用于将a向上舍入到b的倍数
func roundUp(a, b int) int {
	return a + (b-a%b)%b
}

// prepareCipherSpec sets the encryption and MAC states
// that a subsequent changeCipherSpec will use.
func (hc *halfConn) prepareCipherSpec(version uint16, cipher interface{}, mac macFunction) {
	hc.version = version
	hc.nextCipher = cipher
	hc.nextMac = mac
}

// 用于切换加密和MAC状态到之前通过prepareCipherSpec传递的参数
func (hc *halfConn) changeCipherSpec() error {
	if hc.nextCipher == nil {
		return alertInternalError
	}
	hc.cipher = hc.nextCipher
	hc.mac = hc.nextMac
	hc.nextCipher = nil
	hc.nextMac = nil
	for i := range hc.seq {
		hc.seq[i] = 0
	}
	return nil
}

// 用于递增TLS连接的序列号。
// 它从最低位开始递增，如果递增后不为0，则返回。如果所有位都递增到0，则抛出一个panic，表示序列号即将溢出，需要重新协商。
func (hc *halfConn) incSeq() {
	for i := 7; i >= 0; i-- {
		hc.seq[i]++
		if hc.seq[i] != 0 {
			return
		}
	}

	// Not allowed to let sequence number wrap.
	// Instead, must renegotiate before it does.
	// Not likely enough to bother.
	panic("TLS: sequence number wraparound")
}

// 用于对数据进行加密和MAC认证
func (hc *halfConn) encrypt(b *block, explicitIVLen int) (bool, alert) {
	// 如果hc.mac不为空，则计算MAC值，并将其注入到记录中
	if hc.mac != nil {
		mac := hc.mac.MAC(hc.outDigestBuf, hc.seq[0:], b.data[:recordHeaderLen], b.data[recordHeaderLen+explicitIVLen:])

		n := len(b.data)
		b.resize(n + len(mac))
		copy(b.data[n:], mac)
		hc.outDigestBuf = mac
	}

	payload := b.data[recordHeaderLen:]

	// 根据加密算法类型进行加密。目前只支持CBC模式的加密算法
	if hc.cipher != nil {
		switch c := hc.cipher.(type) {
		case cbcMode:
			blockSize := c.BlockSize()
			if explicitIVLen > 0 {
				c.SetIV(payload[:explicitIVLen])
				payload = payload[explicitIVLen:]
			}
			prefix, finalBlock := padToBlockSize(payload, blockSize)
			b.resize(recordHeaderLen + explicitIVLen + len(prefix) + len(finalBlock))
			c.CryptBlocks(b.data[recordHeaderLen+explicitIVLen:], prefix)
			c.CryptBlocks(b.data[recordHeaderLen+explicitIVLen+len(prefix):], finalBlock)
		default:
			panic("unknown cipher type")
		}
	}

	// update length to include MAC and any block padding needed.
	n := len(b.data) - recordHeaderLen
	b.data[3] = byte(n >> 8)
	b.data[4] = byte(n)
	hc.incSeq()

	return true, 0
}

// 解密数据块b，并返回是否解密成功的布尔值、跳过记录开始部分的字节数以获取应用程序负载、以及可选的警告值
func (hc *halfConn) decrypt(b *block) (ok bool, prefixLen int, alertValue alert) {
	// 提取负载
	payload := b.data[recordHeaderLen:]

	macSize := 0
	if hc.mac != nil {
		macSize = hc.mac.Size()
	}

	paddingGood := byte(255)
	explicitIVLen := 0

	// decrypt
	if hc.cipher != nil {
		switch c := hc.cipher.(type) {
		case cbcMode:
			blockSize := c.BlockSize()
			explicitIVLen = blockSize

			if len(payload)%blockSize != 0 || len(payload) < roundUp(explicitIVLen+macSize+1, blockSize) {
				return false, 0, alertBadRecordMAC
			}

			if explicitIVLen > 0 {
				c.SetIV(payload[:explicitIVLen])
				payload = payload[explicitIVLen:]
			}
			c.CryptBlocks(payload, payload)
			payload, paddingGood = removePadding(payload)
			b.resize(recordHeaderLen + explicitIVLen + len(payload))
		default:
			panic("unknown cipher type")
		}
	}

	// check, strip mac
	if hc.mac != nil {
		if len(payload) < macSize {
			return false, 0, alertBadRecordMAC
		}

		// strip mac off payload, b.data
		n := len(payload) - macSize
		b.data[3] = byte(n >> 8)
		b.data[4] = byte(n)
		b.resize(recordHeaderLen + explicitIVLen + n)
		remoteMAC := payload[n:]
		localMAC := hc.mac.MAC(hc.inDigestBuf, hc.seq[0:], b.data[:recordHeaderLen], payload[:n])

		if subtle.ConstantTimeCompare(localMAC, remoteMAC) != 1 || paddingGood != 255 {
			return false, 0, alertBadRecordMAC
		}
		hc.inDigestBuf = localMAC
	}
	hc.incSeq()

	return true, recordHeaderLen + explicitIVLen, 0
}

// 一个简单的数据缓存
type block struct {
	data []byte
	// index for Read
	off  int
	link *block
}

// 该方法从io.Reader接口类型的r中读取数据，并将其存储到b.data切片中，直到b.data切片中至少包含n个字节或出现错误。
func (b *block) readFromUntil(r io.Reader, n int) error {
	// quick case
	if len(b.data) >= n {
		return nil
	}

	// read until have enough.
	b.reserve(n)
	for {
		m, err := r.Read(b.data[len(b.data):cap(b.data)])
		b.data = b.data[0 : len(b.data)+m]
		if len(b.data) >= n {
			// TODO(bradfitz,agl): slightly suspicious
			// that we're throwing away r.Read's err here.
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// 从block的数据中读取内容到切片p，并更新block的偏移量off
func (b *block) Read(p []byte) (n int, err error) {
	n = copy(p, b.data[b.off:])
	b.off += n
	return
}

// 该函数用于调整block的数据切片大小，确保其容量至少为n字节。如果当前容量小于n，则调用reserve方法增加容量。最后将数据切片的长度调整为n
func (b *block) resize(n int) {
	if n > cap(b.data) {
		b.reserve(n)
	}
	b.data = b.data[0:n]
}

// 通过动态扩容的方式，确保了block中的data切片容量至少为n字节，满足对数据容量扩展的需求
func (b *block) reserve(n int) {
	if cap(b.data) >= n {
		return
	}
	m := cap(b.data)
	if m == 0 {
		m = 1024
	}
	for m < n {
		m *= 2
	}
	data := make([]byte, len(b.data), m)
	copy(data, b.data)
	b.data = data
}

// cbc模式加密，密码块链接（Cipher Block Chaining）模式的块密码
type cbcMode interface {
	cipher.BlockMode
	SetIV([]byte)
}
