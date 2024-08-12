package https

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func ListenAndServer(addr, certFile, keyFile string, handler http.Handler) error {
	cert, err := LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	config := &Config{
		Certificates: []Certificate{cert},
	}
	lis, err := Listen("tcp", addr, config)
	if err != nil {
		return err
	}
	defer lis.Close()
	srv := &http.Server{Handler: handler}
	return srv.Serve(lis)
}

// 创建一个TLS监听器
func Listen(network, laddr string, config *Config) (net.Listener, error) {
	if config == nil || len(config.Certificates) == 0 {
		return nil, errors.New("tls: neither Certificates nor GetCertificate set in Config")
	}
	l, err := net.Listen(network, laddr)
	if err != nil {
		return nil, err
	}
	return NewListener(l, config), nil
}

// 用于tls连接的网络监听器，实现了net.Listener接口
type listener struct {
	net.Listener
	config *Config
}

// 接收tcp连接，并返回经过TLS处理的安全连接
func (l *listener) Accept() (c net.Conn, err error) {
	c, err = l.Listener.Accept()
	if err != nil {
		return
	}
	c = Server(c, l.config)
	return
}

// 创建tls Listener，它从内部Listener接受连接，并将每个连接用Server包装
func NewListener(inner net.Listener, config *Config) net.Listener {
	l := new(listener)
	l.Listener = inner
	l.config = config
	return l
}

// 该函数用于从一对文件中读取和解析公钥/私钥对。这些文件必须包含PEM编码的数据。
// 成功返回后，Certificate.Leaf将为nil，因为证书的解析形式不会被保留。
// 函数首先从certFile中读取公钥的PEM编码数据，然后从keyFile中读取私钥的PEM编码数据。
// 最后，它会调用X509KeyPair函数来解析并返回公钥/私钥对。如果读取文件或解析密钥对时发生错误，函数将返回一个错误。
func LoadX509KeyPair(certFile, keyFile string) (Certificate, error) {
	certPEMBlock, err := ioutil.ReadFile(certFile)
	if err != nil {
		return Certificate{}, err
	}
	keyPEMBlock, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return Certificate{}, err
	}
	return X509KeyPair(certPEMBlock, keyPEMBlock)
}

// X509KeyPair parses a public/private key pair from a pair of
// PEM encoded data. On successful return, Certificate.Leaf will be nil because
// the parsed form of the certificate is not retained.
func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (Certificate, error) {
	fail := func(err error) (Certificate, error) { return Certificate{}, err }

	var cert Certificate
	var skippedBlockTypes []string
	for {
		var certDERBlock *pem.Block
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		} else {
			skippedBlockTypes = append(skippedBlockTypes, certDERBlock.Type)
		}
	}

	if len(cert.Certificate) == 0 {
		if len(skippedBlockTypes) == 0 {
			return fail(errors.New("crypto/tls: failed to find any PEM data in certificate input"))
		} else if len(skippedBlockTypes) == 1 && strings.HasSuffix(skippedBlockTypes[0], "PRIVATE KEY") {
			return fail(errors.New("crypto/tls: failed to find certificate PEM data in certificate input, but did find a private key; PEM inputs may have been switched"))
		} else {
			return fail(fmt.Errorf("crypto/tls: failed to find \"CERTIFICATE\" PEM block in certificate input after skipping PEM blocks of the following types: %v", skippedBlockTypes))
		}
	}

	skippedBlockTypes = skippedBlockTypes[:0]
	var keyDERBlock *pem.Block
	for {
		keyDERBlock, keyPEMBlock = pem.Decode(keyPEMBlock)
		if keyDERBlock == nil {
			if len(skippedBlockTypes) == 0 {
				return fail(errors.New("crypto/tls: failed to find any PEM data in key input"))
			} else if len(skippedBlockTypes) == 1 && skippedBlockTypes[0] == "CERTIFICATE" {
				return fail(errors.New("crypto/tls: found a certificate rather than a key in the PEM for the private key"))
			} else {
				return fail(fmt.Errorf("crypto/tls: failed to find PEM block with type ending in \"PRIVATE KEY\" in key input after skipping PEM blocks of the following types: %v", skippedBlockTypes))
			}
		}
		if keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY") {
			break
		}
		skippedBlockTypes = append(skippedBlockTypes, keyDERBlock.Type)
	}

	var err error
	cert.PrivateKey, err = parsePrivateKey(keyDERBlock.Bytes)
	if err != nil {
		return fail(err)
	}

	// We don't need to parse the public key for TLS, but we so do anyway
	// to check that it looks sane and matches the private key.
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fail(err)
	}

	switch pub := x509Cert.PublicKey.(type) {
	case *rsa.PublicKey:
		priv, ok := cert.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return fail(errors.New("crypto/tls: private key type does not match public key type"))
		}
		if pub.N.Cmp(priv.N) != 0 {
			return fail(errors.New("crypto/tls: private key does not match public key"))
		}
	default:
		return fail(errors.New("crypto/tls: unknown public key algorithm"))
	}

	return cert, nil
}

// Attempt to parse the given private key DER block. OpenSSL 0.9.8 generates
// PKCS#1 private keys by default, while OpenSSL 1.0.0 generates PKCS#8 keys.
// OpenSSL ecparam generates SEC1 EC private keys for ECDSA. We try all three.
func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("crypto/tls: found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}

	return nil, errors.New("crypto/tls: failed to parse private key")
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "tls: DialWithDialer timed out" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

// DialWithDialer connects to the given network address using dialer.Dial and
// then initiates a TLS handshake, returning the resulting TLS connection. Any
// timeout or deadline given in the dialer apply to connection and TLS
// handshake as a whole.
//
// DialWithDialer interprets a nil configuration as equivalent to the zero
// configuration; see the documentation of Config for the defaults.
func DialWithDialer(dialer *net.Dialer, network, addr string, config *Config) (*Conn, error) {
	// We want the Timeout and Deadline values from dialer to cover the
	// whole process: TCP connection and TLS handshake. This means that we
	// also need to start our own timers now.
	timeout := dialer.Timeout

	if !dialer.Deadline.IsZero() {
		deadlineTimeout := dialer.Deadline.Sub(time.Now())
		if timeout == 0 || deadlineTimeout < timeout {
			timeout = deadlineTimeout
		}
	}

	var errChannel chan error

	if timeout != 0 {
		errChannel = make(chan error, 2)
		time.AfterFunc(timeout, func() {
			errChannel <- timeoutError{}
		})
	}

	rawConn, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = defaultConfig()
	}

	conn := Client(rawConn, config)

	if timeout == 0 {
		err = conn.Handshake()
	} else {
		go func() {
			errChannel <- conn.Handshake()
		}()

		err = <-errChannel
	}

	if err != nil {
		rawConn.Close()
		return nil, err
	}

	return conn, nil
}

// Dial connects to the given network address using net.Dial
// and then initiates a TLS handshake, returning the resulting
// TLS connection.
// Dial interprets a nil configuration as equivalent to
// the zero configuration; see the documentation of Config
// for the defaults.
func Dial(network, addr string, config *Config) (*Conn, error) {
	return DialWithDialer(new(net.Dialer), network, addr, config)
}
