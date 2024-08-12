package https

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/x509"
	"hash"
	"sync"
)

/*浏览器与服务器在使用 TLS 建立连接的时候实际上就是选了一组加密算法实现安全通信，这些算法组合叫做 “密码套件（cipher suite）”。

套件命名很有规律，比如“ECDHE-RSA-AES256-GCM-SHA384”。按照 密钥交换算法 + 签名算法 + 对称加密算法 + 摘要算法”组成的.

所以这个套件的意思就是：使用 ECDHE 算法进行密钥交换，使用 RSA 签名和身份验证，
握手后使用 AES 对称加密，密钥长度 256 位，分组模式 GCM，消息认证和随机数生成使用摘要算法 SHA384。
*/

var (
	once                   sync.Once
	varDefaultCipherSuites []uint16
)

func defaultCipherSuites() []uint16 {
	once.Do(initDefaultCipherSuites)
	return varDefaultCipherSuites
}

func initDefaultCipherSuites() {
	varDefaultCipherSuites = make([]uint16, 0, len(cipherSuites))
	for _, suite := range cipherSuites {
		varDefaultCipherSuites = append(varDefaultCipherSuites, suite.id)
	}
}

// 用于实现TLS密钥协议的客户端和服务器端，它们负责处理密钥交换消息。
type keyAgreement interface {
	// 处理客户端密钥交换消息，根据传入的参数生成密钥并返回
	processClientKeyExchange(*Config, *Certificate, *clientKeyExchangeMsg, uint16) ([]byte, error)
	// 生成客户端密钥交换消息，根据传入的参数生成密钥及相关消息并返回
	generateClientKeyExchange(*Config, *clientHelloMsg, *x509.Certificate) ([]byte, *clientKeyExchangeMsg, error)
}

// 表示一个特定的密钥协商、加密和MAC函数的组合。所有的cipherSuite目前都假设使用RSA密钥协商
type cipherSuite struct {
	id uint16
	// 表示每个组件所需的密钥材料的长度，以字节为单位
	keyLen int
	// 表示MAC的长度，以字节为单位
	macLen int
	// 表示初始化向量的长度，以字节为单位
	ivLen int
	// 根据版本号返回一个keyAgreement对象，用于密钥协商
	ka func(version uint16) keyAgreement
	//根据密钥、初始化向量和是否为读操作返回一个加密算法对象
	cipher func(key, iv []byte, isRead bool) interface{}
	// 根据版本号和MAC密钥返回一个macFunction对象，用于计算MAC值
	mac func(version uint16, macKey []byte) macFunction
}

var cipherSuites = []*cipherSuite{
	{TLS_RSA_WITH_AES_128_CBC_SHA, 16, 20, 16, rsaKA, cipherAES, macSHA1},
	{TLS_RSA_WITH_AES_256_CBC_SHA, 32, 20, 16, rsaKA, cipherAES, macSHA1},
	{TLS_RSA_WITH_3DES_EDE_CBC_SHA, 24, 20, 8, rsaKA, cipher3DES, macSHA1},
}

func cipher3DES(key, iv []byte, isRead bool) interface{} {
	block, _ := des.NewTripleDESCipher(key)
	if isRead {
		return cipher.NewCBCDecrypter(block, iv)
	}
	return cipher.NewCBCEncrypter(block, iv)
}

func cipherAES(key, iv []byte, isRead bool) interface{} {
	block, _ := aes.NewCipher(key)
	if isRead {
		return cipher.NewCBCDecrypter(block, iv)
	}
	return cipher.NewCBCEncrypter(block, iv)
}

// macSHA1 returns a macFunction for the given protocol version.
func macSHA1(version uint16, key []byte) macFunction {
	return tls10MAC{hmac.New(sha1.New, key)}
}

type macFunction interface {
	Size() int
	MAC(digestBuf, seq, header, data []byte) []byte
}

// tls10MAC implements the TLS 1.0 MAC function. RFC 2246, section 6.2.3.
type tls10MAC struct {
	h hash.Hash
}

func (s tls10MAC) Size() int {
	return s.h.Size()
}

func (s tls10MAC) MAC(digestBuf, seq, header, data []byte) []byte {
	s.h.Reset()
	s.h.Write(seq)
	s.h.Write(header)
	s.h.Write(data)
	return s.h.Sum(digestBuf[:0])
}

func rsaKA(version uint16) keyAgreement {
	return rsaKeyAgreement{}
}

// 该函数用于从给定的支持加密套件列表中，根据对端请求的ID返回一个加密套件
func mutualCipherSuite(have []uint16, want uint16) *cipherSuite {
	for _, id := range have {
		if id == want {
			for _, suite := range cipherSuites {
				if suite.id == want {
					return suite
				}
			}
			return nil
		}
	}
	return nil
}

// A list of the possible cipher suite ids. Taken from
// http://www.iana.org/assignments/tls-parameters/tls-parameters.xml
const (
	TLS_RSA_WITH_3DES_EDE_CBC_SHA uint16 = 0x000a
	TLS_RSA_WITH_AES_128_CBC_SHA  uint16 = 0x002f
	TLS_RSA_WITH_AES_256_CBC_SHA  uint16 = 0x0035
)
