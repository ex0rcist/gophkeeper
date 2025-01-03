package entities

import (
	"fmt"
	"os"
)

// FilePath is a path to a file on local filesystem.
type FilePath string

// Set validates that path exists and assigns it to FilePath.
// Required by pflags interface.
func (p *FilePath) Set(src string) error {
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	*p = FilePath(src)

	return nil
}

// Returns string representation of stored path.
// Required by pflags interface.
func (p FilePath) String() string {
	return string(p)
}

// Required by pflags interface.
func (p FilePath) Type() string {
	return "string"
}
