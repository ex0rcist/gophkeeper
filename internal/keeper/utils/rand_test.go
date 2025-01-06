package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRand(t *testing.T) {
	size := 20
	rand, err := GenerateRandom(size)

	assert.NoError(t, err)
	assert.Len(t, rand, size)
}
