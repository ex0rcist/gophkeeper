package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/keeper/crypto"
	"gophkeeper/internal/keeper/entities"
	"gophkeeper/pkg/models"
	"io"
	"log"
	"os"
	"slices"
	"sync"
)

var _ Storage = (*FileStorage)(nil)

// File-backed storage
type FileStorage struct {
	sync.RWMutex

	file *os.File
	Data map[uint64]models.Secret `json:"secrets"`

	encrypter crypto.Encrypter
	password  string
}

func NewFileStorage(path string, password string, encrypter crypto.Encrypter) (*FileStorage, error) {
	store := &FileStorage{
		Data:      make(map[uint64]models.Secret),
		encrypter: encrypter,
		password:  password,
	}

	err := store.openOrCreateFile(path)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (store *FileStorage) Get(_ context.Context, id uint64) (*models.Secret, error) {
	store.Lock()
	defer store.Unlock()

	secret, ok := store.Data[id]
	secret.ID = id
	if !ok {
		return nil, entities.ErrSecretNotFound
	}

	return &secret, nil
}

func (store *FileStorage) GetAll(_ context.Context) ([]*models.Secret, error) {
	store.Lock()
	defer store.Unlock()

	arr := make([]*models.Secret, len(store.Data))

	i := 0
	for id, secret := range store.Data {
		secret.ID = id
		arr[i] = &secret
		i++
	}

	return arr, nil
}

func (store *FileStorage) Create(ctx context.Context, secret *models.Secret) error {
	store.Lock()
	defer store.Unlock()

	id := store.nextID()
	store.Data[id] = *secret

	return store.dump()
}

func (store *FileStorage) Update(ctx context.Context, secret *models.Secret) error {
	store.Lock()
	defer store.Unlock()

	id := secret.ID
	store.Data[id] = *secret

	return store.dump()
}

func (store *FileStorage) Delete(ctx context.Context, id uint64) error {
	store.Lock()
	defer store.Unlock()

	delete(store.Data, id)

	return store.dump()
}

func (store *FileStorage) String() string {
	return store.file.Name()
}

func (store *FileStorage) Close(_ context.Context) error {
	var err error

	err = store.dump()
	if err != nil {
		return err
	}

	defer func() {
		if err := store.file.Close(); err != nil {
			log.Fatal("dump(): failed to close file: %w")
		}
	}()

	return err
}

func (store *FileStorage) openOrCreateFile(path string) error {
	store.Lock()
	defer store.Unlock()

	var (
		existedFile bool
		err         error
	)

	// Check if file already existed
	if _, err = os.Stat(path); err == nil {
		existedFile = true
	}

	// Open file
	store.file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("openOrCreateFile(): failed to open file: %w", err)
	}

	if !existedFile {
		if err := store.dump(); err != nil {
			return fmt.Errorf("openOrCreateFile -> dump(): failed to write initial data: %w", err)
		}
	} else {
		// Read
		encryptedData, err := io.ReadAll(store.file)
		if err != nil {
			return fmt.Errorf("openOrCreateFile -> ReadAll(): failed to read file: %w", err)
		}

		// Decrypt
		decryptedData, err := store.DecryptWithRecover(encryptedData, store.password)
		if err != nil {
			switch err {
			case entities.ErrBadPassword, entities.ErrBadEncryption:
				return err
			default:
				return fmt.Errorf("openOrCreateFile -> Decrypt(): failed to decrypt data: %w", err)
			}
		}

		// Unmarshal
		if err := json.Unmarshal(decryptedData, &store.Data); err != nil {
			return fmt.Errorf("openOrCreateFile -> Unmarshal(): failed to decode Data: %w", err)
		}

		// Reset pointer
		if _, err := store.file.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("restore(): failed to reset file pointer: %w", err)
		}
	}

	return nil
}

// Attempt to decode non-encoded file may cause panic
func (store *FileStorage) DecryptWithRecover(data []byte, password string) (res []byte, err error) {
	defer func() { // defer can replace named return values
		if r := recover(); r != nil {
			err = entities.ErrBadEncryption
		}
	}()

	return store.encrypter.Decrypt(data, password)
}

// Dump storage to file
func (store *FileStorage) dump() (err error) {
	// Serialize data
	data, err := json.Marshal(store.Data)
	if err != nil {
		return fmt.Errorf("dump(): error serializing Data: %w", err)
	}

	// Clear any existing data
	if err := store.file.Truncate(0); err != nil {
		return fmt.Errorf("dump(): failed to truncate file: %w", err)
	}

	// Encrypt data
	encryptedData, err := store.encrypter.Encrypt(data, store.password)
	if err != nil {
		return fmt.Errorf("dump(): error encrypting Data: %w", err)
	}

	// Dump to file
	if _, err := store.file.Write(encryptedData); err != nil {
		return fmt.Errorf("dump(): error writing encrypted Data to file: %w", err)
	}

	// Reset pointer
	if _, err := store.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("dump(): failed to reset file pointer: %w", err)
	}

	return nil
}

func (store *FileStorage) nextID() uint64 {
	if len(store.Data) == 0 {
		return 1
	}

	ids := make([]uint64, len(store.Data))
	for k := range store.Data {
		ids = append(ids, k)
	}

	return slices.Max(ids) + 1
}
