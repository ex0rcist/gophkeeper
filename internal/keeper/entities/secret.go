package entities

import (
	"strings"
)

// Entity to hide sensitive values which shouldn't leak to logs.
type Secret string

// Set value.
func (s *Secret) Set(src string) error {
	*s = Secret(src)

	return nil
}

// Return a string for correct type conversion.
func (s Secret) Type() string {
	return "string"
}

// Stringer.
func (s Secret) String() string {
	if len(s) <= 2 {
		return string(s)
	}

	masked := strings.Repeat("*", len(s)-2)
	return string(s[0]) + masked + string(s[len(s)-1])
}
