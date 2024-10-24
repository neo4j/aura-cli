package fileutils

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

/* Reads a file, if it doesn't exist returns empty []byte */
func ReadFileSafe(fs afero.Fs, path string) ([]byte, error) {
	exists, err := FileExists(fs, path)

	if err != nil {
		return nil, err
	}
	if exists {
		return afero.ReadFile(fs, path)
	} else {
		return []byte{}, nil
	}
}

func ReadOrCreateFile(fs afero.Fs, path string) ([]byte, error) {
	exists, err := FileExists(fs, path)

	if err != nil {
		return nil, err
	}
	if exists {
		return afero.ReadFile(fs, path)
	} else {
		return []byte{}, createFile(fs, path)
	}
}

func WriteFile(fs afero.Fs, path string, data []byte) error {
	return afero.WriteFile(fs, path, data, 0600)
}

func FileExists(fs afero.Fs, path string) (bool, error) {
	if _, err := fs.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, err
	}
}

func createFile(fs afero.Fs, path string) error {
	if err := fs.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	_, err := fs.Create(path)
	if err != nil {
		return err
	}

	if err = fs.Chmod(path, 0600); err != nil {
		return err
	}

	return nil
}
