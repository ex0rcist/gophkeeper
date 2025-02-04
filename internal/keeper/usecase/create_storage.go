package usecase

import (
	"fmt"
	"gophkeeper/internal/keeper/crypto"
	"gophkeeper/internal/keeper/storage"
)

type CreateLocalStoreUseCase struct {
}

func NewCreateStorageUsecase() *CreateLocalStoreUseCase {
	return &CreateLocalStoreUseCase{}
}

func (uc CreateLocalStoreUseCase) Call(path string, password string, encrypter crypto.Encrypter) (*storage.FileStorage, error) {
	if path == "" {
		return nil, fmt.Errorf("no path provided")
	}

	fs, err := storage.NewFileStorage(path, password, encrypter)
	return fs, err
}
