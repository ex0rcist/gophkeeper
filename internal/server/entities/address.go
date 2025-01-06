package entities

import (
	"fmt"
	"strconv"
	"strings"
)

// https://github.com/spf13/pflag/blob/master/README.md#usage
// must comply with https://pkg.go.dev/github.com/spf13/pflag@v1.0.5#Value.
type Address string

// Stringer.
func (a Address) String() string {
	return string(a)
}

// Set value.
func (a *Address) Set(src string) error {
	chunks := strings.Split(src, ":")
	if len(chunks) != 2 {
		return fmt.Errorf("set address failed: %w", ErrBadAddressFormat)
	}

	port := chunks[1]

	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("set address failed: %w", err)
	}

	*a = Address(src)

	return nil
}

// Return string for correct type conversion.
func (a Address) Type() string {
	return "string"
}
