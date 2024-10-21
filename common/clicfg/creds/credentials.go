package creds

import (
	"encoding/json"
	"path/filepath"

	"github.com/neo4j/cli/common/clicfg/fileutils"
	"github.com/spf13/afero"
)

type CredentialsFile struct {
	Aura AuraCredentials `json:"aura"`
}

type Credentials struct {
	fs afero.Fs

	Aura     *AuraCredentials
	filePath string
}

func NewCredentials(fs afero.Fs, configPrefix string) (*Credentials, error) {
	configPath := filepath.Join(configPrefix, "neo4j", "cli", "credentials.json")
	c := Credentials{
		fs:       fs,
		filePath: configPath,
	}
	c.load()
	return &c, nil
}

func (c *Credentials) load() error {
	var data, err = fileutils.ReadFileSafe(c.fs, c.filePath)
	fileHasData := len(data) != 0
	if err != nil {
		return err
	}

	var credentials CredentialsFile = CredentialsFile{
		Aura: AuraCredentials{
			Credentials: []AuraCredential{},
			onSave:      c.save,
		},
	}
	if fileHasData {
		if err := json.Unmarshal(data, &credentials); err != nil {
			return err
		}
	}

	c.Aura = &credentials.Aura

	if !fileHasData {
		return c.save()
	}

	return nil
}

func (c *Credentials) save() error {
	data, err := json.Marshal(CredentialsFile{
		Aura: *c.Aura,
	})
	if err != nil {
		return err
	}

	return fileutils.WriteFile(c.fs, c.filePath, data)

}
