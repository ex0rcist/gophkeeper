package entities

import (
	"errors"
	"fmt"
)

var (
	ErrBadCredentials = errors.New("bad auth credentials")

	ErrStorageUnpingable = errors.New("healthcheck is not supported")
	ErrUnexpected        = errors.New("unexpected error")
	ErrBadAddressFormat  = errors.New("bad net address format")

	ErrSecretNotFound = errors.New("secret not found")
	ErrNoSecrets      = errors.New("no secrets found")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

func ErrorUserAlreadyExists(login string) error {
	return fmt.Errorf("%w (%s)", ErrUserAlreadyExists, login)
}

func ErrorSecretNotFound(secretID uint64) error {
	return fmt.Errorf("%w (id=%d)", ErrSecretNotFound, secretID)
}
