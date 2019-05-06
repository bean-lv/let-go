package cryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

type AES interface {
	Encrypt(origin string) (string, error)
	EncryptData(origin []byte) ([]byte, error)
	Decrypt(crypted string) (string, error)
	DecryptData(crypted []byte) ([]byte, error)
}

type xAes struct {
	key string
}

var XAES AES

func InitAES(key string) {
	XAES = &xAes{key: key}
}

func (x *xAes) Encrypt(origin string) (string, error) {
	crypted, err := x.EncryptData([]byte(origin))
	return string(crypted), err
}

func (x *xAes) EncryptData(origin []byte) ([]byte, error) {
	key := []byte(x.key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])

	// 填充原始数据
	origin = PKCS5Padding(origin, blockSize)

	// 加密
	crypted := make([]byte, len(origin))
	blockMode.CryptBlocks(crypted, origin)

	return crypted, nil
}

func (x *xAes) Decrypt(crypted string) (string, error) {
	origin, err := x.DecryptData([]byte(crypted))
	return string(origin), err
}

func (x *xAes) DecryptData(crypted []byte) ([]byte, error) {
	key := []byte(x.key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])

	// 解密
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	// 还原原始数据
	origData = PKCS5Unpadding(origData)

	return origData, nil
}

// 填充模式，将待加密数据填充为blockSize的整数倍
// 取模后如果不足blockSize，则补差值
// 取模后如果正好是blockSize，则再补blockSize
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
