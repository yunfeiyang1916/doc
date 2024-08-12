package https

import (
	"crypto/rsa"
	"crypto/subtle"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
)

// 执行客户端握手
func (c *Conn) clientHandshake() error {
	if c.config == nil {
		c.config = defaultConfig()
	}
	// 发送客户端Hello消息
	hello, err := c.sendClientHelloMsg()
	if err != nil {
		return err
	}
	// 等待服务端Hello消息
	serverHello, err := c.waitServerHelloMsg()
	if err != nil {
		return err
	}

	c.vers = VersionTLS12
	// 版本已经协商
	c.haveVers = true

	// 根据客户端和服务器协商的密码套件选择一个双方都支持的密码套件
	suite := mutualCipherSuite(hello.cipherSuites, serverHello.cipherSuite)
	if suite == nil {
		c.sendAlert(alertHandshakeFailure)
		return errors.New("tls: server chose an unconfigured cipher suite")
	}
	// 创建握手状态
	hs := &clientHandshakeState{
		c:            c,
		serverHello:  serverHello,
		hello:        hello,
		suite:        suite,
		finishedHash: newFinishedHash(c.vers),
	}

	hs.finishedHash.Write(hs.hello.marshal())
	hs.finishedHash.Write(hs.serverHello.marshal())
	// 执行完整握手过程
	if err := hs.doFullHandshake(); err != nil {
		return err
	}
	// 建立会话密钥
	if err := hs.establishKeys(); err != nil {
		return err
	}
	// 发送完成消息
	if err := hs.sendFinished(c.firstFinished[:]); err != nil {
		return err
	}
	// 等待完成消息
	if err := hs.readFinished(nil); err != nil {
		return err
	}

	c.handshakeComplete = true
	c.cipherSuite = suite.id
	return nil
}

// 发送客户端Hello消息,包括客户端随机数，自己的TLS版本号，以及自己支持的加密套件及支持的压缩方式
func (c *Conn) sendClientHelloMsg() (*clientHelloMsg, error) {
	hello := &clientHelloMsg{
		vers:               VersionTLS12,
		compressionMethods: []uint8{compressionNone},
		random:             make([]byte, 32),
	}

	possibleCipherSuites := c.config.cipherSuites()
	hello.cipherSuites = make([]uint16, 0, len(possibleCipherSuites))

NextCipherSuite:
	for _, suiteId := range possibleCipherSuites {
		for _, suite := range cipherSuites {
			if suite.id != suiteId {
				continue
			}
			hello.cipherSuites = append(hello.cipherSuites, suiteId)
			continue NextCipherSuite
		}
	}

	_, err := io.ReadFull(c.config.rand(), hello.random)
	if err != nil {
		c.sendAlert(alertInternalError)
		return nil, errors.New("tls: short read from Rand: " + err.Error())
	}

	_, err = c.writeRecord(recordTypeHandshake, hello.marshal())
	return hello, err
}

// 等待并读取服务器的Hello消息
func (c *Conn) waitServerHelloMsg() (*serverHelloMsg, error) {
	msg, err := c.readHandshake()
	if err != nil {
		return nil, err
	}
	serverHello, ok := msg.(*serverHelloMsg)
	if !ok {
		c.sendAlert(alertUnexpectedMessage)
		return nil, unexpectedMessageError(serverHello, msg)
	}
	return serverHello, nil
}

// 用于保存客户端握手状态
type clientHandshakeState struct {
	c            *Conn
	serverHello  *serverHelloMsg
	hello        *clientHelloMsg
	suite        *cipherSuite
	finishedHash finishedHash
	masterSecret []byte
}

// 进行完整握手
func (hs *clientHandshakeState) doFullHandshake() error {
	c := hs.c
	// 等待并处理服务端的证书消息
	certMsg, err := hs.waitCertificateMsg()
	hs.finishedHash.Write(certMsg.marshal())
	// 验证服务端证书
	certs, err := hs.verifyServer(certMsg.certificates)
	if err != nil {
		c.sendAlert(alertInternalError)
		return err
	}
	// 只支持RSA公钥
	switch certs[0].PublicKey.(type) {
	case *rsa.PublicKey:
		break
	default:
		c.sendAlert(alertUnsupportedCertificate)
		return fmt.Errorf("tls: server's certificate contains an unsupported type of public key: %T", certs[0].PublicKey)
	}
	// 等待并处理服务端的HelloDone消息
	shd, err := hs.waitServerHelloDoneMsg()
	if err != nil {
		c.sendAlert(alertInternalError)
		return err
	}
	hs.finishedHash.Write(shd.marshal())
	// 发送客户端密钥交换消息，并获取预主密钥
	preMasterSecret, err := hs.sendClientKeyExchangeMsg(certs[0])
	if err != nil {
		return err
	}
	// 根据预主密钥计算主密钥，并清空finishedHash的握手缓冲区
	hs.masterSecret = masterFromPreMasterSecret(c.vers, preMasterSecret, hs.hello.random, hs.serverHello.random)
	hs.finishedHash.discardHandshakeBuffer()

	return nil
}

// 等待服务端的证书消息
func (hs *clientHandshakeState) waitCertificateMsg() (*certificateMsg, error) {
	c := hs.c
	msg, err := c.readHandshake()
	if err != nil {
		return nil, err
	}
	certMsg, ok := msg.(*certificateMsg)
	if !ok || len(certMsg.certificates) == 0 {
		c.sendAlert(alertUnexpectedMessage)
		return nil, unexpectedMessageError(certMsg, msg)
	}
	return certMsg, nil
}

// 等待服务端的HelloDone消息
func (hs *clientHandshakeState) waitServerHelloDoneMsg() (*serverHelloDoneMsg, error) {
	c := hs.c
	msg, err := c.readHandshake()
	if err != nil {
		return nil, err
	}
	shd, ok := msg.(*serverHelloDoneMsg)
	if !ok {
		c.sendAlert(alertUnexpectedMessage)
		return nil, unexpectedMessageError(shd, msg)
	}
	return shd, nil
}

// 该函数用于验证服务端的证书
func (hs *clientHandshakeState) verifyServer(certificates [][]byte) ([]*x509.Certificate, error) {
	c := hs.c

	certs := make([]*x509.Certificate, len(certificates))
	for i, asn1Data := range certificates {
		cert, err := x509.ParseCertificate(asn1Data)
		if err != nil {
			c.sendAlert(alertBadCertificate)
			return nil, errors.New("tls: failed to parse certificate from server: " + err.Error())
		}
		certs[i] = cert
	}
	// 使用x509.VerifyOptions选项验证证书链
	opts := x509.VerifyOptions{
		Roots:         c.config.RootCAs,
		CurrentTime:   c.config.time(),
		DNSName:       c.config.ServerName,
		Intermediates: x509.NewCertPool(),
	}

	for i, cert := range certs {
		if i == 0 {
			continue
		}
		opts.Intermediates.AddCert(cert)
	}
	var err error
	c.verifiedChains, err = certs[0].Verify(opts)
	if err != nil {
		// NOTE: just left log, then proceed
		fmt.Printf("failed to verify server. %s\n", err)
	}
	return certs, nil
}

// 发送客户端密钥交换消息，并获取预主密钥
func (hs *clientHandshakeState) sendClientKeyExchangeMsg(cert *x509.Certificate) ([]byte, error) {
	c := hs.c

	preMasterSecret, ckx, err := hs.suite.ka(c.vers).generateClientKeyExchange(c.config, hs.hello, cert)
	if err != nil {
		c.sendAlert(alertInternalError)
		return nil, err
	}
	if ckx == nil {
		c.sendAlert(alertInternalError)
		return nil, errors.New("tls: unexpected ServerKeyExchange")
	}
	hs.finishedHash.Write(ckx.marshal())
	// send clientKeyExchangeMsg message
	c.writeRecord(recordTypeHandshake, ckx.marshal())
	return preMasterSecret, nil
}

// 该函数用于在TLS握手过程中建立客户端和服务器的加密密钥和认证密钥
func (hs *clientHandshakeState) establishKeys() error {
	c := hs.c

	clientMAC, serverMAC, clientKey, serverKey, clientIV, serverIV :=
		keysFromMasterSecret(c.vers, hs.masterSecret, hs.hello.random, hs.serverHello.random, hs.suite.macLen, hs.suite.keyLen, hs.suite.ivLen)
	var clientCipher, serverCipher interface{}
	var clientHash, serverHash macFunction

	clientCipher = hs.suite.cipher(clientKey, clientIV, false /* not for reading */)
	clientHash = hs.suite.mac(c.vers, clientMAC)
	serverCipher = hs.suite.cipher(serverKey, serverIV, true /* for reading */)
	serverHash = hs.suite.mac(c.vers, serverMAC)

	c.in.prepareCipherSpec(c.vers, serverCipher, serverHash)
	c.out.prepareCipherSpec(c.vers, clientCipher, clientHash)
	return nil
}

// 该函数用于读取服务端的"Finished"消息，并验证其正确性。
func (hs *clientHandshakeState) readFinished(out []byte) error {
	c := hs.c

	c.readRecord(recordTypeChangeCipherSpec)
	if err := c.in.error(); err != nil {
		return err
	}

	msg, err := c.readHandshake()
	if err != nil {
		return err
	}
	serverFinished, ok := msg.(*finishedMsg)
	if !ok {
		c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(serverFinished, msg)
	}

	verify := hs.finishedHash.serverSum(hs.masterSecret)
	if len(verify) != len(serverFinished.verifyData) ||
		subtle.ConstantTimeCompare(verify, serverFinished.verifyData) != 1 {
		c.sendAlert(alertHandshakeFailure)
		return errors.New("tls: server's Finished message was incorrect")
	}
	hs.finishedHash.Write(serverFinished.marshal())
	copy(out, verify)
	return nil
}

// 向服务端发送Finished消息
func (hs *clientHandshakeState) sendFinished(out []byte) error {
	c := hs.c

	c.writeRecord(recordTypeChangeCipherSpec, []byte{1})

	finished := new(finishedMsg)
	finished.verifyData = hs.finishedHash.clientSum(hs.masterSecret)
	hs.finishedHash.Write(finished.marshal())
	c.writeRecord(recordTypeHandshake, finished.marshal())
	copy(out, finished.verifyData)
	return nil
}
