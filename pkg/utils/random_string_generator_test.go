package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	randomString5 := GenerateRandomString(5)
	assert.NotEmpty(t, randomString5, "Random string should not be empty")
	assert.Len(t, randomString5, 5, "Random string should have 5 characters")

	randomString8 := GenerateRandomString(8)
	assert.NotEmpty(t, randomString8)
	assert.Len(t, randomString8, 8, "Random string should have 8 characters")

	randomString15 := GenerateRandomString(15)
	assert.NotEmpty(t, randomString15)
	assert.Len(t, randomString15, 15, "Random string should have 15 characters")
}
