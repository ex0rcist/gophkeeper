package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	password := "password"
	plaintext := []byte{0x55, 0x44, 0x33, 0x22}

	encrypted, err := Encrypt([]byte(password), plaintext)

	assert.NoError(t, err)

	decrypted, err := Decrypt(password, []byte(encrypted))

	assert.NoError(t, err)

	assert.Equal(t, plaintext, decrypted)
}
