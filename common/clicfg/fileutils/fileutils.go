package fileutils

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

/* Reads a file, if it doesn't exist returns empty []byte */
func ReadFileSafe(fs afero.Fs, path string) []byte {
	exists := FileExists(fs, path)

	if exists {
		data, err := afero.ReadFile(fs, path)
		if err != nil {
			panic(err)
		}
		return data
	} else {
		return []byte{}
	}
}

func ReadOrCreateFile(fs afero.Fs, path string) []byte {
	exists := FileExists(fs, path)

	if exists {
		data, err := afero.ReadFile(fs, path)
		if err != nil {
			panic(err)
		}
		return data
	} else {
		createFile(fs, path)
		return []byte{}
	}
}

func WriteFile(fs afero.Fs, path string, data []byte) {
	err := afero.WriteFile(fs, path, data, 0600)
	if err != nil {
		panic(err)
	}
}

func FileExists(fs afero.Fs, path string) bool {
	if _, err := fs.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		panic(err)
	}
}

func createFile(fs afero.Fs, path string) {
	if err := fs.MkdirAll(filepath.Dir(path), 0755); err != nil {
		panic(err)
	}

	_, err := fs.Create(path)
	if err != nil {
		panic(err)
	}

	if err = fs.Chmod(path, 0600); err != nil {
		panic(err)
	}
}
