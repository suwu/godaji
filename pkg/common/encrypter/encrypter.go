package encrypter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"mitaitech.com/oa/pkg/common/config"
)

func PKCS5UnPadding(originalBytes []byte) []byte {
	length := len(originalBytes)
	unpadding := int(originalBytes[length-1])
	return originalBytes[:(length - unpadding)]
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func AesEncrypt(originalBytes []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	originalBytes = PKCS5Padding(originalBytes, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	cipherBytes := make([]byte, len(originalBytes))
	blockMode.CryptBlocks(cipherBytes, originalBytes)
	return cipherBytes, nil
}

func AesDecrypt(cipherBytes []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originalBytes := make([]byte, len(cipherBytes))
	blockMode.CryptBlocks(originalBytes, cipherBytes)
	originalBytes = PKCS5UnPadding(originalBytes)
	return originalBytes, nil
}

func Encrypt(originalText string) (string, error) {
	cfg := config.GetConfig()
	key := []byte(cfg.EncryptKey)
	cipherBytes, err := AesEncrypt([]byte(originalText), key)
	if err != nil {
		return "", err
	}
	base64str := base64.StdEncoding.EncodeToString(cipherBytes)
	return base64str, nil
}

func Decrypt(cipherText string) (string, error) {
	cfg := config.GetConfig()
	key := []byte(cfg.EncryptKey)
	cipherBytes, _ := base64.StdEncoding.DecodeString(cipherText)
	cipherBytes, err := AesDecrypt(cipherBytes, key)
	if err != nil {
		return "", err
	}
	return string(cipherBytes), nil
}
