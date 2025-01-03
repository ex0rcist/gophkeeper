package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/keeper/entities"
	"gophkeeper/pkg/models"
	"io"
	"os"
	"slices"
	"sync"
)

var _ Storage = (*FileStorage)(nil)

// File-backed storage
type FileStorage struct {
	sync.RWMutex

	Data map[uint64]models.Secret `json:"secrets"`

	file     *os.File
	password string
}

func NewFileStorage(path string) (*FileStorage, error) {
	store := &FileStorage{
		Data: make(map[uint64]models.Secret),
	}

	err := store.openOrCreateFile(path)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (store *FileStorage) Get(_ context.Context, id uint64) (models.Secret, error) {
	store.Lock()
	defer store.Unlock()

	secret, ok := store.Data[id]
	secret.ID = id
	if !ok {
		return models.Secret{}, entities.ErrSecretNotFound
	}

	return secret, nil
}

func (store *FileStorage) GetAll(_ context.Context) ([]models.Secret, error) {
	store.Lock()
	defer store.Unlock()

	arr := make([]models.Secret, len(store.Data))

	i := 0
	for id, secret := range store.Data {
		secret.ID = id
		arr[i] = secret
		i++
	}

	return arr, nil
}

func (store *FileStorage) Create(ctx context.Context, secret models.Secret) error {
	store.Lock()
	defer store.Unlock()

	id := store.nextID()
	store.Data[id] = secret

	return store.dump()
}

func (store *FileStorage) Update(ctx context.Context, secret models.Secret) error {
	store.Lock()
	defer store.Unlock()

	id := secret.ID
	store.Data[id] = secret

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
	return store.dump()
}

func (store *FileStorage) openOrCreateFile(path string) error {
	var (
		existedFile bool
		err         error
	)

	store.Lock()
	defer store.Unlock()

	// Проверяем, существует ли файл
	if _, err = os.Stat(path); err == nil {
		existedFile = true
	}

	// Открываем файл для чтения и записи (без очистки содержимого)
	store.file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("restoreOrCreate(): failed to open file: %w", err)
	}

	// Если файл не существовал, создаем начальные данные
	if !existedFile {
		if err := store.dump(); err != nil {
			return fmt.Errorf("dump(): failed to write initial data: %w", err)
		}
	} else {
		// Читаем данные из файла
		decoder := json.NewDecoder(store.file)
		if err := decoder.Decode(&store.Data); err != nil {
			return fmt.Errorf("restore(): failed to decode Data: %w", err)
		}
		// Сбрасываем указатель файла на начало, чтобы избежать проблем с записью
		store.file.Seek(0, io.SeekStart)
	}

	return nil
}

// dump storage to file
func (store *FileStorage) dump() (err error) {
	// store.TryLock()
	// defer store.Unlock()

	// defer func() {
	// 	if closeErr := store.file.Close(); err == nil && closeErr != nil {
	// 		err = fmt.Errorf("dump(): failed to close file: %w", closeErr)
	// 	}
	// }()

	encoder := json.NewEncoder(store.file)
	if err := encoder.Encode(store.Data); err != nil {
		return fmt.Errorf("error encoding Data: %w", err)
	}

	// todo: encryption

	return nil
}

func (store *FileStorage) nextID() uint64 {
	// store.TryRLock()
	// defer store.RUnlock()

	if len(store.Data) == 0 {
		return 1
	}

	ids := make([]uint64, len(store.Data))
	for k := range store.Data {
		ids = append(ids, k)
	}

	return slices.Max(ids) + 1
}
