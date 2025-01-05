package entities

import (
	"errors"
)

var (
	ErrUnexpected        = errors.New("unexpected error")
	ErrBadAddressFormat  = errors.New("bad net address format")
	ErrSecretNotFound    = errors.New("secret not found in storage")
	ErrBadFileStorePath  = errors.New("file at store path was not found")
	ErrBadPassword       = errors.New("incorrect password")
	ErrBadEncryption     = errors.New("failed to decrypt file")
	ErrServerUnavailable = errors.New("server unavailable")
	ErrUnauthenticated   = errors.New("failed to authenticate")
	ErrAlreadyExist      = errors.New("user already exists")
	// ErrNoSubscribers   = errors.New("no clients subscribed")
)
