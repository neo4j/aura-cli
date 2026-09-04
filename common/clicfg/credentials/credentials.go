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

type credentialsFileOnDisk struct {
	Aura auraCredentialsOnDisk `json:"aura"`
}

type auraCredentialOnDisk struct {
	AuraCredential
	ClientSecret string `json:"client-secret"`
	AccessToken  string `json:"access-token"`
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
		},
	}
	if fileHasData {
		var onDisk credentialsFileOnDisk
		if err := json.Unmarshal(data, &onDisk); err != nil {
			panic(err)
		}
		credentials.Aura = onDisk.Aura.toAuraCredentials()
	}

	c.Aura = credentials.Aura
	c.Aura.refresh = c.refreshAura
	c.Aura.persist = c.save

	if !fileHasData {
		c.save()
	}
}

// refreshAura re-reads the credentials file from disk into the existing c.Aura, discarding
// whatever was loaded or mutated in memory before this call. It updates the struct in place
// rather than replacing it, so the refresh/persist closures wired up in load stay intact.
func (c *Credentials) refreshAura() error {
	data := fileutils.ReadFileSafe(c.fs, c.filePath)

	if len(data) == 0 {
		c.Aura.Credentials = []*AuraCredential{}
		c.Aura.DefaultCredential = ""
		return nil
	}

	var onDisk credentialsFileOnDisk
	if err := json.Unmarshal(data, &onDisk); err != nil {
		return err
	}

	fresh := onDisk.Aura.toAuraCredentials()
	c.Aura.Credentials = fresh.Credentials
	c.Aura.DefaultCredential = fresh.DefaultCredential
	return nil
}

func (c *Credentials) save() error {
	onDisk := credentialsFileOnDisk{
		Aura: c.Aura.toOnDisk(),
	}

	data, err := json.Marshal(onDisk)
	if err != nil {
		panic(err)
	}

	fileutils.WriteFile(c.fs, c.filePath, data)
	return nil
}
