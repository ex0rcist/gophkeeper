package entities

import (
	"testing"
)

func TestSecretString_Masking(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"supersecretpassword", "s*****************d"},
		{"short", "s***t"},
		{"ab", "ab"},
		{"a", "a"},
		{"", ""},
	}

	for _, tt := range tests {
		result := SecretString(tt.input).String()

		if result != tt.expected {
			t.Errorf("expected masked output to be '%s', got '%s' for input '%s'", tt.expected, result, tt.input)
		}
	}
}
