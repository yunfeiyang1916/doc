package https

import (
	"crypto"
	"crypto/x509"
	"fmt"
)

const (
	VersionTLS12 = 0x0303
)

const (
	// 最大明文负载长度
	maxPlaintext = 16384 // maximum plaintext payload length
	// 最大密文负载长度
	maxCiphertext = 16384 + 2048 // maximum ciphertext payload length
	// 记录头长度
	recordHeaderLen = 5 // record header length
	// 最大握手长度
	maxHandshake = 65536 // maximum handshake we support (protocol max is 16 MB)
)

// TLS compression types.
const (
	compressionNone uint8 = 0
)

// TLS record types.
type recordType uint8

const (
	recordTypeChangeCipherSpec recordType = 20
	recordTypeAlert            recordType = 21
	recordTypeHandshake        recordType = 22
	recordTypeApplicationData  recordType = 23
)

// TLS handshake message types.
const (
	typeClientHello       uint8 = 1
	typeServerHello       uint8 = 2
	typeCertificate       uint8 = 11
	typeServerHelloDone   uint8 = 14
	typeClientKeyExchange uint8 = 16
	typeFinished          uint8 = 20
)

// A Certificate is a chain of one or more certificates, leaf first.
type Certificate struct {
	Certificate [][]byte
	// PrivateKey contains the private key corresponding to the public key
	// in Leaf. For a server, this must implement crypto.Signer and/or
	// crypto.Decrypter, with an RSA.
	PrivateKey crypto.PrivateKey
	// Leaf is the parsed form of the leaf certificate, which may be
	// initialized using x509.ParseCertificate to reduce per-handshake
	// processing for TLS clients doing client authentication. If nil, the
	// leaf certificate will be parsed as needed.
	Leaf *x509.Certificate
}

func unexpectedMessageError(wanted, got interface{}) error {
	return fmt.Errorf("tls: received unexpected handshake message of type %T when waiting for %T", got, wanted)
}
