package entities

import (
	"errors"
)

var (
	ErrBadCredentials = errors.New("bad auth credentials")

	ErrUnexpected       = errors.New("unexpected error")
	ErrBadAddressFormat = errors.New("bad net address format")
	ErrSecretNotFound   = errors.New("secret not found in storage")
	ErrBadFileStorePath = errors.New("file at store path was not found")

	ErrBadPassword   = errors.New("incorrect password")
	ErrBadEncryption = errors.New("failed to decrypt file")

	// ErrWrongSecretType = errors.New("invalid secret type")

	// ErrNumberInvalid  = errors.New("card number is invalid")
	// ErrNoSubscribers   = errors.New("no clients subscribed")
)

// func ErrorUserAlreadyExists(login string) error {
// 	return fmt.Errorf("%w (%s)", ErrUserAlreadyExists, login)
// }

// func ErrorSecretNotFound(secretID int) error {
// 	return fmt.Errorf("%w (id=%d)", ErrSecretNotFound, secretID)
// }
