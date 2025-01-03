package entities

import (
	"regexp"
	"strings"
)

// Entity to hide sensitive strings which shouldn't leak to logs.
type SecretString string

// Stringer.
func (s SecretString) String() string {
	if len(s) <= 2 {
		return string(s)
	}

	masked := strings.Repeat("*", len(s)-2)
	return string(s[0]) + masked + string(s[len(s)-1])
}

// Entity to hide sensitive URI values (e.g. login:password) which shouldn't leak to logs.
type SecretConnURI string

var _URISecrets = regexp.MustCompile(`(://).*:.*(@)`)

func (u SecretConnURI) String() string {
	return string(_URISecrets.ReplaceAll([]byte(u), []byte("$1*****:*****$2")))
}
