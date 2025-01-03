package entities_test

import (
	"gophkeeper/internal/server/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressType(t *testing.T) {
	addr := entities.Address("0.0.0.0:8080")

	require.Equal(t, "string", addr.Type())
}

func TestAddressSet(t *testing.T) {
	tt := []struct {
		name string
		src  string
		want bool
	}{
		{name: "correct interface", src: "0.0.0.0:8080", want: true},
		{name: "with ip", src: "42.42.42.42:4242", want: true},
		{name: "localhost and port", src: "localhost:8080", want: true},
		{name: "no colon", src: "localhost8080", want: false},
		{name: "no port", src: "localhost", want: false},
		{name: "invalid format", src: "localhost/42", want: false},
		{name: "invalid port ", src: "localhost:abc", want: false},
		{name: "no host", src: "8080", want: false},
		{name: "stolen", src: "", want: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			addr := new(entities.Address)

			err := addr.Set(tc.src)

			if !tc.want {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.src, addr.String())
			}
		})
	}
}
