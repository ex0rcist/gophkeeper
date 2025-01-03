package entities

import (
	"testing"
)

func TestSecretSet(t *testing.T) {
	var s Secret
	err := s.Set("mysecrettokenvalueforauthentication")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if s != "mysecrettokenvalueforauthentication" {
		t.Errorf("expected Secret to be set to 'mysecrettokenvalueforauthentication', got %s", s)
	}
}

func TestSecretSet_WarnOnShortSecret(t *testing.T) {
	var s Secret
	err := s.Set("short")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if s != "short" {
		t.Errorf("expected Secret to be set to 'short', got %s", s)
	}
}

func TestSecretType(t *testing.T) {
	var s Secret
	if s.Type() != "string" {
		t.Errorf("expected Type to return 'string', got %s", s.Type())
	}
}

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
		var s Secret
		err := s.Set(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		result := s.String()

		if result != tt.expected {
			t.Errorf("expected masked output to be '%s', got '%s' for input '%s'", tt.expected, result, tt.input)
		}
	}
}
