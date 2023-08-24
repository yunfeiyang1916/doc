package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"fmt"
)

// 计算md5值
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// DES 对称加密
// @param	data	要加密的数据
// @param	key		密钥
// @param	iv		偏移向量
func AesEncrypt(data, key, iv []byte) ([]byte, error) {
	//加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//使用PKCS5Padding填充方式
	data = PKCS5Padding(data, block.BlockSize())
	//加密模式为CBC
	blockModel := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(data))
	//加密块
	blockModel.CryptBlocks(crypted, data)
	return crypted, nil
}

// PKCS5填充方式
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding) //用0去填充
	return append(ciphertext, padtext...)
}

// PKCS5去除填充
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
func main() {
	//原始密码
	password := Md5("123456")
	fmt.Printf("原始密码：%s\n", password)
	//api授权的apiSecret
	apiSecret := "12345678901234567890123456789012"
	//16位密钥
	key := apiSecret[0:16]
	//16位偏移量
	iv := apiSecret[16:]
	//aes加密
	encrypt, _ := AesEncrypt([]byte(password), []byte(key), []byte(iv))
	fmt.Printf("加密后的密文：%d\n", len(encrypt))
}
