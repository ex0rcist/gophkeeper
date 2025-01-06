package entities

import (
	"os"
	"testing"
)

func TestFilePathSet_ValidPath(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		rErr := os.Remove(tempFile.Name())
		if rErr != nil {
			t.Error(rErr)
		}
	}()

	var fp FilePath
	if err := fp.Set(tempFile.Name()); err != nil {
		t.Errorf("expected Set to succeed, got error: %v", err)
	}

	if fp.String() != tempFile.Name() {
		t.Errorf("expected FilePath to be %s, got %s", tempFile.Name(), fp.String())
	}
}

func TestFilePathSet_InvalidPath(t *testing.T) {
	var fp FilePath
	invalidPath := "/nonexistent/path/to/file"

	err := fp.Set(invalidPath)
	if err == nil {
		t.Errorf("expected error for invalid path, got nil")
	}

	expectedErrMsg := "invalid file path"
	if err != nil && err.Error()[:len(expectedErrMsg)] != expectedErrMsg {
		t.Errorf("expected error message to start with %q, got %q", expectedErrMsg, err.Error())
	}
}

func TestFilePathString(t *testing.T) {
	fp := FilePath("/path/to/file")

	if fp.String() != "/path/to/file" {
		t.Errorf("expected String to return '/path/to/file', got %s", fp.String())
	}
}

func TestFilePathType(t *testing.T) {
	fp := FilePath("/path/to/file")

	if fp.Type() != "string" {
		t.Errorf("expected Type to return 'string', got %s", fp.Type())
	}
}
