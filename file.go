package yagolib

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func IsFileExists(path string) bool {
	// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	if path, err := NormalizePath(path); err == nil {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return true
		}
	}
	return false
}

func NormalizePath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err == nil {
			path = filepath.Join(usr.HomeDir, path[1:])
		} else {
			return path, err
		}
	}
	return filepath.Abs(path)
}
