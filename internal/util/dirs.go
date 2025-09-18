package util

import (
	"fmt"
	"os"
)

const (
	dirsMode os.FileMode = 0755
)

// EnsureDir creates dir if not exists.
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, dirsMode); err != nil {
		return fmt.Errorf("dir '%s' does not exist and cannot be created: %w", path, err)
	}

	return nil
}
