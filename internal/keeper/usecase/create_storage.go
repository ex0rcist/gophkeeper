package usecase

import (
	"context"
	"gophkeeper/internal/keeper/storage"
)

type createLocalStoreUseCase struct {
}

func NewCreateStorageUsecase() *createLocalStoreUseCase {
	return &createLocalStoreUseCase{}
}

func (uc createLocalStoreUseCase) Call(path string) (*storage.FileStorage, error) {
	var err error

	ctx := context.Background()

	if path == "" {
		panic("empty path") // todo
	}

	fs, err := storage.NewFileStorage(path)
	defer fs.Close(ctx)

	return fs, err
}
