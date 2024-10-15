package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const allowedChars = "abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-!*@"

func GenerateRandomString(length int) string {
	var stringBuilder strings.Builder
	stringBuilder.Grow(length)

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(allowedChars))))
		// Pick and append the character at that index
		stringBuilder.WriteByte(allowedChars[randomIndex.Int64()])
	}
	return stringBuilder.String()
}
