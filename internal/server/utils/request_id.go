package utils

import (
	uuid "github.com/satori/go.uuid"
)

func GenerateRequestID() string {
	return uuid.NewV4().String()
}
