package util

import (
	"os"
)

const filesMode os.FileMode = 0644

func WriteFile(target string, content []byte) error {
	if err := os.WriteFile(target, content, filesMode); err != nil {
		return err
	}

	return nil
}
