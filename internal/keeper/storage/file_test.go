package storage

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"gophkeeper/internal/keeper/entities"
	"gophkeeper/pkg/models"

	"github.com/stretchr/testify/assert"
)

type MockEncrypter struct{}

func (m *MockEncrypter) Encrypt(data []byte, password string) ([]byte, error) {
	return data, nil // No-op encryption for testing
}

func (m *MockEncrypter) Decrypt(data []byte, password string) ([]byte, error) {
	return data, nil // No-op decryption for testing
}

func TestFileStorage(t *testing.T) {
	tempFile, err := os.CreateTemp("", "filestorage_test.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Initialize file with valid JSON structure
	initialData := map[uint64]models.Secret{}
	initialJSON, err := json.Marshal(initialData)
	assert.NoError(t, err)
	_, err = tempFile.Write(initialJSON)
	assert.NoError(t, err)
	tempFile.Close()

	encrypter := &MockEncrypter{}
	password := "testpassword"

	store, err := NewFileStorage(tempFile.Name(), password, encrypter)
	assert.NoError(t, err)
	t.Run("Create and Get Secret", func(t *testing.T) {
		secret := &models.Secret{
			ID:         1,
			Title:      "Test Secret",
			Metadata:   "metadata",
			SecretType: "credential",
			Payload:    []byte("payload"),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := store.Create(context.Background(), secret)
		assert.NoError(t, err)

		storedSecret, err := store.Get(context.Background(), secret.ID)
		assert.NoError(t, err)
		assert.Equal(t, secret.Title, storedSecret.Title)
		assert.Equal(t, secret.Metadata, storedSecret.Metadata)
	})

	t.Run("Update Secret", func(t *testing.T) {
		updatedSecret := &models.Secret{
			ID:         1,
			Title:      "Updated Secret",
			Metadata:   "updated metadata",
			SecretType: "credential",
			Payload:    []byte("updated payload"),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := store.Update(context.Background(), updatedSecret)
		assert.NoError(t, err)

		storedSecret, err := store.Get(context.Background(), updatedSecret.ID)
		assert.NoError(t, err)
		assert.Equal(t, updatedSecret.Title, storedSecret.Title)
		assert.Equal(t, updatedSecret.Metadata, storedSecret.Metadata)
	})

	t.Run("Delete Secret", func(t *testing.T) {
		err := store.Delete(context.Background(), 1)
		assert.NoError(t, err)

		_, err = store.Get(context.Background(), 1)
		assert.ErrorIs(t, err, entities.ErrSecretNotFound)
	})

	t.Run("Get All Secrets", func(t *testing.T) {
		secrets := []*models.Secret{
			{ID: 1, Title: "Secret 1"},
			{ID: 2, Title: "Secret 2"},
		}

		for _, secret := range secrets {
			err := store.Create(context.Background(), secret)
			assert.NoError(t, err)
		}

		allSecrets, err := store.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Len(t, allSecrets, len(secrets))
	})

	t.Run("Dump and Load", func(t *testing.T) {
		secret := &models.Secret{
			ID:         3,
			Title:      "Test Dump",
			Metadata:   "dump metadata",
			SecretType: "credential",
			Payload:    []byte("dump payload"),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := store.Create(context.Background(), secret)
		assert.NoError(t, err)

		store.Close(context.Background())

		newStore, err := NewFileStorage(tempFile.Name(), password, encrypter)
		assert.NoError(t, err)

		loadedSecret, err := newStore.Get(context.Background(), secret.ID)
		assert.NoError(t, err)
		assert.Equal(t, secret.Title, loadedSecret.Title)
	})
}
