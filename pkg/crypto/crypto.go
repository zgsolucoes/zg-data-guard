package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type CryptoHelper struct {
	Key []byte
}

func NewCryptoHelper(key string) *CryptoHelper {
	return &CryptoHelper{
		Key: []byte(key),
	}
}

func (helper *CryptoHelper) Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(helper.Key)
	if err != nil {
		return "", err
	}

	plainTextBytes := []byte(plainText)
	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextBytes)

	return hex.EncodeToString(cipherText), nil
}

func (helper *CryptoHelper) Decrypt(cipherTextHex string) (string, error) {
	block, err := aes.NewCipher(helper.Key)
	if err != nil {
		return "", err
	}

	cipherText, _ := hex.DecodeString(cipherTextHex)
	if len(cipherText) < aes.BlockSize {
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
