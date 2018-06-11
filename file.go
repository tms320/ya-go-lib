package yagolib

import (
	"os"
	"path/filepath"
)

func IsFileExists(path string) bool {
	// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	path, _ = filepath.Abs(path)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
