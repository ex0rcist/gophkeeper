// Provides functions necessary for encryption and decryption
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"gophkeeper/internal/keeper/entities"
	"gophkeeper/internal/keeper/utils"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const saltLen = 8

var _ Encrypter = (*KeeperEncrypter)(nil)

type Encrypter interface {
	Encrypt(data []byte, password string) ([]byte, error)
	Decrypt(encrypted []byte, password string) ([]byte, error)
}

type KeeperEncrypter struct {
	saltLen int
}

func NewKeeperEncrypter() *KeeperEncrypter {
	return &KeeperEncrypter{saltLen: saltLen}
}

func (e KeeperEncrypter) Encrypt(plaintext []byte, password string) ([]byte, error) {
	key, salt, err := e.deriveKey(password, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// creating AES block
	AESBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// creating GCM
	GCM, err := cipher.NewGCM(AESBlock)
	if err != nil {
		return nil, err
	}

	// generating nonce
	nonce, err := utils.GenerateRandom(GCM.NonceSize())
	if err != nil {
		return nil, err
	}

	// encrypt data
	encrypted := GCM.Seal(nonce, nonce, plaintext, nil)

	// store salt alongside encrypted data
	encrypted = append(encrypted, salt...)

	return encrypted, nil
}

func (e KeeperEncrypter) Decrypt(encrypted []byte, password string) ([]byte, error) {
	// extract salt
	saltIdx := len(encrypted) - e.saltLen
	salt := encrypted[saltIdx:]

	encrypted = encrypted[:saltIdx]

	key, _, err := e.deriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// creating AES block
	AESBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// creating GCM
	GCM, err := cipher.NewGCM(AESBlock)
	if err != nil {
		return nil, err
	}

	// extract nonce
	nonce := encrypted[:GCM.NonceSize()]
	encrypted = encrypted[GCM.NonceSize():]

	// decrypt data
	decrypted, err := GCM.Open(nil, nonce, encrypted, nil)
	if err != nil {
		if strings.Contains(err.Error(), "message authentication failed") {
			return nil, entities.ErrBadPassword
		}
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return decrypted, nil
}

func (e KeeperEncrypter) deriveKey(password string, salt []byte) ([]byte, []byte, error) {
	if len(salt) == 0 {
		salt = make([]byte, e.saltLen)
		_, err := rand.Read(salt)
		if err != nil {
			return nil, nil, err
		}
	}
	return pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New), salt, nil
}
