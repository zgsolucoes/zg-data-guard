package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const invalidKey = "testkey"
const valid32lenghtKey = "my32lengthsupersecretkeyyyy11111"

func TestNewCryptoHelper(t *testing.T) {
	helper := NewCryptoHelper(valid32lenghtKey)

	assert.Equal(t, []byte(valid32lenghtKey), helper.Key, "Key should be correctly set")
}

func TestGivenInvalidKey_WhenEncrypt_ThenShouldReturnError(t *testing.T) {
	helper := NewCryptoHelper(invalidKey)
	plainText := "Hello, World!"

	// Test encryption
	cipherText, err := helper.Encrypt(plainText)
	assert.Error(t, err, "Encryption should return an error")
	assert.Empty(t, cipherText, "Cipher text should be empty")
}

func TestGivenValidInput_WhenEncrypt_ThenShouldEncrypt(t *testing.T) {
	helper := NewCryptoHelper(valid32lenghtKey)
	plainText := "Hello, World!"

	// Test encryption
	cipherText, err := helper.Encrypt(plainText)
	assert.NoError(t, err, "Encryption should not return an error")
	assert.NotEqual(t, plainText, cipherText, "Cipher text should not be the same as the plain text")
}

func TestGivenInvalidKey_WhenDecrypt_ThenShouldReturnError(t *testing.T) {
	helper := NewCryptoHelper(invalidKey)
	plainText := "Hello, World!"

	// Test encryption
	decryptedText, err := helper.Decrypt(plainText)
	assert.Error(t, err, "Decryption should return an error")
	assert.Empty(t, decryptedText, "Cipher text should be empty")
}

func TestGivenValidKey_WhenDecrypt_ThenShouldDecrypt(t *testing.T) {
	helper := NewCryptoHelper(valid32lenghtKey)
	plainText := "Hello, World!"

	// Test encryption
	cipherText, err := helper.Encrypt(plainText)

	// Test decryption
	decryptedText, err := helper.Decrypt(cipherText)
	assert.NoError(t, err, "Decryption should not return an error")
	assert.Equal(t, plainText, decryptedText, "Decrypted text should be the same as the original plain text")
}
