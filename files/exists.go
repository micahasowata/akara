package files

import (
	"errors"
	"os"
)

// IsFile ensure that location contains a file that actually exists
func IsFile(location string) bool {
	info, err := os.Stat(location)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}

	return !info.IsDir()
}

// IsDir ensure that location contains a directory that actually exists
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}

	return info.IsDir()
}
