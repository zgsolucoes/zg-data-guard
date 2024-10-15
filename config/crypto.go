package config

import (
	"os"

	"github.com/zgsolucoes/zg-data-guard/pkg/crypto"
)

const defaultAesKey = "my32l3ngthsup3rs3cr3tno0n3kn0ws1"

var cryptoHelper *crypto.CryptoHelper

func GetCryptoHelper() *crypto.CryptoHelper {
	if cryptoHelper == nil {
		initializeCryptography()
	}
	return cryptoHelper
}

func initializeCryptography() {
	aesKey := os.Getenv("AES_PRIVATE_KEY")
	if aesKey == "" {
		aesKey = defaultAesKey
	}
	cryptoHelper = crypto.NewCryptoHelper(aesKey)
}
