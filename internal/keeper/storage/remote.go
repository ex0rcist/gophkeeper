package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/keeper/api"
	"gophkeeper/internal/keeper/crypto"
	"gophkeeper/pkg/models"
)

var _ Storage = (*RemoteStorage)(nil)

// Remote storage
type RemoteStorage struct {
	client    api.IApiClient
	encrypter crypto.Encrypter
	password  string // passw to encrypt payload
}

func NewRemoteStorage(client api.IApiClient, encrypter crypto.Encrypter) (*RemoteStorage, error) {
	if encrypter == nil {
		encrypter = crypto.NewKeeperEncrypter()
	}

	store := &RemoteStorage{
		client:    client,
		encrypter: encrypter,
		password:  client.GetPassword(),
	}

	return store, nil
}

func (store *RemoteStorage) Get(_ context.Context, id uint64) (*models.Secret, error) {
	secret, err := store.client.LoadSecret(context.Background(), id)
	if err != nil {
		return nil, err
	}

	err = store.decryptPayload(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (store *RemoteStorage) GetAll(_ context.Context) ([]*models.Secret, error) {
	secrets, err := store.client.LoadSecrets(context.Background())
	if err != nil {
		return nil, err
	}

	for _, s := range secrets {
		err = store.decryptPayload(s)
		if err != nil {
			return nil, err
		}
	}

	return secrets, nil
}

func (store *RemoteStorage) Create(ctx context.Context, secret *models.Secret) (err error) {
	err = store.encryptPayload(secret)
	if err != nil {
		return
	}

	err = store.client.SaveSecret(context.Background(), secret)
	return err
}

func (store *RemoteStorage) Update(ctx context.Context, secret *models.Secret) (err error) {
	err = store.encryptPayload(secret)
	if err != nil {
		return
	}

	err = store.client.SaveSecret(context.Background(), secret)
	return err
}

func (store *RemoteStorage) Delete(ctx context.Context, id uint64) (err error) {
	err = store.client.DeleteSecret(context.Background(), id)
	return err
}

func (store *RemoteStorage) encryptPayload(secret *models.Secret) (err error) {
	// Marshal
	data, err := marshalSecret(secret)
	if err != nil {
		return fmt.Errorf("encryptPayload(): error serializing data: %w", err)
	}

	// Encrypt data
	encryptedData, err := store.encrypter.Encrypt(data, store.password)
	if err != nil {
		return fmt.Errorf("encryptPayload(): error encrypting Data: %w", err)
	} else {
		secret.Payload = encryptedData
	}

	return err
}

func (store *RemoteStorage) decryptPayload(secret *models.Secret) (err error) {
	// Decrypt data
	decryptedData, err := store.encrypter.Decrypt(secret.Payload, store.password)
	if err != nil {
		return fmt.Errorf("decryptPayload: failed to decrypt data: %w", err)

	}

	// Unmarshal
	err = unmarshalSecret(secret, decryptedData)
	if err != nil {
		return fmt.Errorf("decryptPayload: failed to unmarshal data: %w", err)
	}

	return nil
}

func marshalSecret(secret *models.Secret) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	switch models.SecretType(secret.SecretType) {
	case models.CredSecret:
		data, err = json.Marshal(secret.Creds)
	case models.TextSecret:
		data, err = json.Marshal(secret.Text)
	case models.CardSecret:
		data, err = json.Marshal(secret.Card)
	case models.BlobSecret:
		data, err = json.Marshal(secret.Blob)
	}

	return data, err
}

func unmarshalSecret(secret *models.Secret, data []byte) error {
	var err error

	switch models.SecretType(secret.SecretType) {
	case models.CredSecret:
		err = json.Unmarshal(data, &secret.Creds)
	case models.TextSecret:
		err = json.Unmarshal(data, &secret.Text)
	case models.CardSecret:
		err = json.Unmarshal(data, &secret.Card)
	case models.BlobSecret:
		err = json.Unmarshal(data, &secret.Blob)
	}

	return err
}

func (store *RemoteStorage) String() string {
	return "remote storage"
}

func (store *RemoteStorage) Close(_ context.Context) error {
	return nil
}
