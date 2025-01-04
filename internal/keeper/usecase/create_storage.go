package usecase

import (
	"gophkeeper/internal/keeper/crypto"
	"gophkeeper/internal/keeper/storage"
)

type CreateLocalStoreUseCase struct {
}

func NewCreateStorageUsecase() *CreateLocalStoreUseCase {
	return &CreateLocalStoreUseCase{}
}

func (uc CreateLocalStoreUseCase) Call(path string, password string, encrypter crypto.Encrypter) (*storage.FileStorage, error) {
	var err error

	if path == "" {
		panic("empty path") // todo
	}

	fs, err := storage.NewFileStorage(path, password, encrypter)
	return fs, err
}
