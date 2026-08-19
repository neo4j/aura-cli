// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package credentials

import (
	"encoding/json"
	"path/filepath"

	"github.com/neo4j/cli/common/clicfg/fileutils"
	"github.com/spf13/afero"
)

type auraCredentialsOnDisk struct {
	DefaultCredential string                 `json:"default-credential"`
	Credentials       []auraCredentialOnDisk `json:"credentials"`
}

type auraCredentialOnDisk struct {
	AuraCredential
	ClientSecret string `json:"client-secret"`
}

type CredentialsFile struct {
	Aura *AuraCredentials `json:"aura"`
}

type Credentials struct {
	fs       afero.Fs
	Aura     *AuraCredentials
	filePath string
}

func NewCredentials(fs afero.Fs, configPrefix string) *Credentials {
	configPath := filepath.Join(configPrefix, "neo4j", "cli", "credentials.json")
	c := Credentials{
		fs:       fs,
		filePath: configPath,
	}
	c.load()
	return &c
}

func (c *Credentials) load() {
	data := fileutils.ReadFileSafe(c.fs, c.filePath)
	fileHasData := len(data) != 0

	var credentials CredentialsFile = CredentialsFile{
		Aura: &AuraCredentials{
			Credentials: []*AuraCredential{},
			onUpdate:    c.save,
		},
	}
	if fileHasData {
		var onDisk struct {
			Aura auraCredentialsOnDisk `json:"aura"`
		}
		if err := json.Unmarshal(data, &onDisk); err != nil {
			panic(err)
		}
		credentials.Aura = onDisk.Aura.toAuraCredentials(c.save)
	}

	c.Aura = credentials.Aura

	if !fileHasData {
		c.save()
	}
}

func (c *Credentials) save() {
	onDisk := struct {
		Aura auraCredentialsOnDisk `json:"aura"`
	}{
		Aura: c.Aura.toOnDisk(),
	}

	data, err := json.Marshal(onDisk)
	if err != nil {
		panic(err)
	}

	fileutils.WriteFile(c.fs, c.filePath, data)
}
