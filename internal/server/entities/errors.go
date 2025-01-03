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

	// ErrWrongSecretType = errors.New("invalid secret type")

	// ErrNumberInvalid  = errors.New("card number is invalid")
	// ErrNoSubscribers   = errors.New("no clients subscribed")
)

func ErrorUserAlreadyExists(login string) error {
	return fmt.Errorf("%w (%s)", ErrUserAlreadyExists, login)
}

func ErrorSecretNotFound(secretID uint64) error {
	return fmt.Errorf("%w (id=%d)", ErrSecretNotFound, secretID)
}

// type ErrUserExists struct {
// 	Login string
// }

// func (e *ErrUserExists) Error() string {
// 	return fmt.Sprintf("user with login '%s' has already registered", e.Login)
// }

// func (e *ErrUserExists) Is(tgt error) bool {
// 	target, ok := tgt.(*ErrUserExists)
// 	if !ok {
// 		return false
// 	}
// 	return e.Login == target.Login
// }
