package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	encrypter := NewKeeperEncrypter()

	password := "password"
	plaintext := []byte{0x55, 0x44, 0x33, 0x22}

	encrypted, err := encrypter.Encrypt(plaintext, password)
	assert.NoError(t, err)

	decrypted, err := encrypter.Decrypt(encrypted, password)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}
