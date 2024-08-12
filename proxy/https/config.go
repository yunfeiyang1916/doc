package https

import (
	"crypto/rand"
	"crypto/x509"
	"io"
	"sync"
	"time"
)

var emptyConfig Config

func defaultConfig() *Config {
	return &emptyConfig
}

// 用于配置TLS客户端或服务端
type Config struct {
	// 提供非对称加密和nonces and RSA所需的随机数源。
	// 如果Rand为nil，TLS将使用crypto/rand包中的加密随机读取器。该读取器必须适用于多个goroutine并发安全使用
	Rand io.Reader
	// 返回自纪元以来的当前时间，单位为秒。如果Time为nil，TLS将使用time.Now作为当前时间的来源
	Time func() time.Time
	// 包含一个或多个证书链，用于向连接的另一端展示。服务器配置必须至少包含一个证书，或者设置GetCertificate
	Certificates []Certificate
	// 定义了客户端在验证服务器证书时使用的根证书颁发机构集合。如果RootCAs为nil，TLS将使用主机的根CA集合
	RootCAs *x509.CertPool
	// 用于验证返回证书上的主机名
	ServerName string
	// 支持的密码套件列表。如果CipherSuites为nil，TLS将使用实现支持的密码套件列表
	CipherSuites []uint16

	// 保护sessionTicketKeys的互斥锁
	mutex sync.RWMutex
}

func (c *Config) rand() io.Reader {
	r := c.Rand
	if r == nil {
		return rand.Reader
	}
	return r
}

func (c *Config) time() time.Time {
	t := c.Time
	if t == nil {
		t = time.Now
	}
	return t()
}

func (c *Config) cipherSuites() []uint16 {
	s := c.CipherSuites
	if s == nil {
		s = defaultCipherSuites()
	}
	return s
}

// getCertificate returns the best certificate,
// defaulting to the first element of c.Certificates.
func (c *Config) getCertificate() *Certificate {
	// return the first certificate.
	return &c.Certificates[0]
}
