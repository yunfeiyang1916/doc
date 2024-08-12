package cert

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// 表示用于创建证书唯一序列号时所使用的上限值。该变量可以容纳最大长度为20字节（即2^(8*20)-1）的无符号整数作为证书的唯一序列号。
// 具体实现上，它是通过将字节255重复20次，并将这20个字节解析为一个整数值来初始化一个*big.Int类型的变量。这样确保生成的序列号始终在一个预定义的有效范围内。
var MaxSerialNumber = big.NewInt(0).SetBytes(bytes.Repeat([]byte{255}, 20))

// tls证书配置
type CertConfig struct {
	// ca根证书
	ca *x509.Certificate
	// ca根证书的私钥
	caPriv crypto.PrivateKey
	// 需要生成的证书的私钥
	priv *rsa.PrivateKey
	// 密钥id
	keyID []byte
}

// 使用CA根证书和私钥生成动态证书
func NewCertConfig(ca *x509.Certificate, caPriv crypto.PrivateKey) (*CertConfig, error) {
	// 生成一个2048位的RSA公私钥对
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	// 公钥
	pub := priv.PublicKey
	// 计算公钥的Subject Key Identifier：根据RFC 3280中关于X.509证书规范的要求，为生成的公钥计算其Subject Key Identifier（SKID）。
	// 这通过先序列化公钥，然后使用SHA-1哈希算法生成哈希值来实现
	pkixPubKey, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	h := sha1.New()
	h.Write(pkixPubKey)
	keyID := h.Sum(nil)

	return &CertConfig{
		ca:     ca,
		caPriv: caPriv,
		priv:   priv,
		keyID:  keyID,
	}, nil
}

// 生成TLS配置。该TLS配置会根据客户端在TLS ClientHello中提供的SNI扩展动态生成证书
func (c *CertConfig) TLSConfig() *tls.Config {
	return &tls.Config{
		GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			if clientHello.ServerName == "" {
				return nil, errors.New("missing server name (SNI)")
			}

			return c.cert(clientHello.ServerName)
		},
		MinVersion: tls.VersionTLS12,
		NextProtos: []string{"http/1.1"},
	}
}

// 生成子证书
func (c *CertConfig) cert(hostname string) (*tls.Certificate, error) {
	// Remove the port if it exists.
	host, _, err := net.SplitHostPort(hostname)
	if err == nil {
		hostname = host
	}

	serial, err := rand.Int(rand.Reader, MaxSerialNumber)
	if err != nil {
		return nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   hostname,
			Organization: []string{"yunfeiyang"},
		},
		SubjectKeyId:          c.keyID,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
	}
	if ip := net.ParseIP(hostname); ip != nil {
		tmpl.IPAddresses = []net.IP{ip}
	} else {
		tmpl.DNSNames = []string{hostname}
	}
	raw, err := x509.CreateCertificate(rand.Reader, tmpl, c.ca, c.priv.Public(), c.caPriv)
	if err != nil {
		return nil, err
	}

	// Parse certificate bytes so that we have a leaf certificate.
	x509c, err := x509.ParseCertificate(raw)
	if err != nil {
		return nil, err
	}

	return &tls.Certificate{
		Certificate: [][]byte{raw, c.ca.Raw},
		PrivateKey:  c.priv,
		Leaf:        x509c,
	}, nil
}

// 加载或创建证书，如果证书或密钥文件不存在，则创建新的密钥对并保存到磁盘
func LoadOrCreateCA(caKeyFile, caCertFile string) (*x509.Certificate, *rsa.PrivateKey, error) {
	tlsCA, err := tls.LoadX509KeyPair(caCertFile, caKeyFile)
	if err == nil {
		caCert, err := x509.ParseCertificate(tlsCA.Certificate[0])
		if err != nil {
			return nil, nil, fmt.Errorf("proxy: could not parse CA: %w", err)
		}

		caKey, ok := tlsCA.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, nil, errors.New("proxy: private key is not RSA")
		}

		return caCert, caKey, nil
	}
	if !os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("proxy: could not load CA key pair: %w", err)
	}

	// 创建目录
	keyDir, _ := filepath.Split(caKeyFile)
	if keyDir != "" {
		if _, err := os.Stat(keyDir); os.IsNotExist(err) {
			if err := os.MkdirAll(keyDir, 0o755); err != nil {
				return nil, nil, fmt.Errorf("proxy: could not create directory for CA key: %w", err)
			}
		}
	}

	keyDir, _ = filepath.Split(caCertFile)
	if keyDir != "" {
		if _, err := os.Stat("keyDir"); os.IsNotExist(err) {
			if err := os.MkdirAll(keyDir, 0o755); err != nil {
				return nil, nil, fmt.Errorf("proxy: could not create directory for CA cert: %w", err)
			}
		}
	}
	caCert, caKey, err := NewCA("yunfeiyang-proxy-ca", "yunfeiyang CA", 365*24*time.Hour)
	if err != nil {
		return nil, nil, fmt.Errorf("proxy: could not generate new CA keypair: %w", err)
	}
	// Open CA certificate and key files for writing.
	certOut, err := os.Create(caCertFile)
	if err != nil {
		return nil, nil, fmt.Errorf("proxy: could not open cert file for writing: %w", err)
	}

	keyOut, err := os.OpenFile(caKeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, nil, fmt.Errorf("proxy: could not open key file for writing: %w", err)
	}
	// Write PEM blocks to CA certificate and key files.
	if err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: caCert.Raw}); err != nil {
		return nil, nil, fmt.Errorf("proxy: could not write CA certificate to disk: %w", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("proxy: could not convert private key to DER format: %w", err)
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return nil, nil, fmt.Errorf("proxy: could not write CA key to disk: %w", err)
	}

	return caCert, caKey, nil
}

// 该函数用于创建一个新的CA证书和关联的私钥。函数参数包括证书名称、组织名称和证书有效期。
// 函数内部通过生成2048位的RSA私钥来创建证书，并设置证书的主题、序列号、有效期等信息。最后返回创建的x509证书和私钥
func NewCA(name, organization string, validity time.Duration) (*x509.Certificate, *rsa.PrivateKey, error) {
	// 生成一个2048位的RSA公私钥对
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	// 公钥
	pub := priv.Public()
	// 计算公钥的Subject Key Identifier：根据RFC 3280中关于X.509证书规范的要求，为生成的公钥计算其Subject Key Identifier（SKID）。
	// 这通过先序列化公钥，然后使用SHA-1哈希算法生成哈希值来实现
	pkixPubKey, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, nil, err
	}

	h := sha1.New()
	h.Write(pkixPubKey)
	keyID := h.Sum(nil)
	// serial multiple times.
	serial, err := rand.Int(rand.Reader, MaxSerialNumber)
	if err != nil {
		return nil, nil, err
	}
	tmpl := &x509.Certificate{
		// 序列号
		SerialNumber: serial,
		// 主体
		Subject: pkix.Name{
			CommonName:   name,
			Organization: []string{organization},
		},
		SubjectKeyId:          keyID,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(validity),
		DNSNames:              []string{name},
		IsCA:                  true,
	}
	// 创建自签证书
	raw, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, pub, priv)
	if err != nil {
		return nil, nil, err
	}

	// Parse certificate bytes so that we have a leaf certificate.
	x509c, err := x509.ParseCertificate(raw)
	if err != nil {
		return nil, nil, err
	}

	return x509c, priv, nil
}
