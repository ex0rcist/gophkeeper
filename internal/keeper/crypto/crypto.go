package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const saltSize = 16
const keySize = 32
const iterations = 100_000

// DeriveKey generates a key from a password and a salt using PBKDF2.
func deriveKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, iterations, keySize, sha256.New)
}

// Encrypt encrypts data using AES-GCM with a key derived from the password.
func Encrypt(data, password []byte) (string, error) {
	// Generate a random salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive the encryption key
	key := deriveKey(password, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	ciphertext := aesGCM.Seal(nil, nonce, data, nil)

	// Combine salt, nonce, and ciphertext
	result := append(salt, nonce...)
	result = append(result, ciphertext...)

	// Return as Base64-encoded string
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts data using AES-GCM with a key derived from the password.
func Decrypt(encodedCiphertext string, password []byte) ([]byte, error) {
	// Decode the Base64 string
	data, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Base64: %w", err)
	}

	// Extract salt, nonce, and ciphertext
	if len(data) < saltSize {
		return nil, errors.New("invalid ciphertext")
	}
	salt := data[:saltSize]
	data = data[saltSize:]

	// Derive the decryption key
	key := deriveKey(password, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce and actual ciphertext
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("invalid ciphertext")
	}
	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}
