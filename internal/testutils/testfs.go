package testutils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/neo4j/cli/pkg/clicfg"
	"github.com/spf13/afero"
)

func GetDefaultTestFs() (afero.Fs, error) {
	return GetTestFs("{}")
}

func GetTestFs(config string) (afero.Fs, error) {
	fs := afero.NewMemMapFs()

	configPath := filepath.Join(clicfg.ConfigPrefix, "neo4j", "cli", "config.json")

	if err := fs.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, err
	}

	f, err := fs.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err = f.Write([]byte(config)); err != nil {
		return nil, err
	}

	return fs, nil
}

func GetTestConfig(fs afero.Fs) (string, error) {
	configPath := filepath.Join(clicfg.ConfigPrefix, "neo4j", "cli", "config.json")

	file, err := fs.Open(configPath)
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	return string(b), nil
}
