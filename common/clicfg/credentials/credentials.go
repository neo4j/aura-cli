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
	fs                 afero.Fs
	Aura               *AuraCredentials
	filePath           string
	initialCredentials map[string]bool
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
		var onDisk credentialsFileOnDisk
		if err := json.Unmarshal(data, &onDisk); err != nil {
			panic(err)
		}
		credentials.Aura = onDisk.Aura.toAuraCredentials(c.save)
	}

	c.Aura = credentials.Aura

	c.initialCredentials = make(map[string]bool)
	for _, cred := range c.Aura.Credentials {
		c.initialCredentials[cred.Name] = true
	}

	if !fileHasData {
		c.save()
	}
}

func (c *Credentials) save() {
	diskData := fileutils.ReadFileSafe(c.fs, c.filePath)
	if len(diskData) > 0 {
		var diskFile credentialsFileOnDisk
		if err := json.Unmarshal(diskData, &diskFile); err == nil {
			c.mergeWithDisk(diskFile.Aura.toAuraCredentials(c.Aura.onUpdate))
		}
	}

	onDisk := credentialsFileOnDisk{
		Aura: c.Aura.toOnDisk(),
	}

	data, err := json.Marshal(onDisk)
	if err != nil {
		panic(err)
	}

	fileutils.WriteFile(c.fs, c.filePath, data)
}

func (c *Credentials) mergeWithDisk(diskAura *AuraCredentials) {
	diskCredentialMap := make(map[string]*AuraCredential)
	for _, cred := range diskAura.Credentials {
		diskCredentialMap[cred.Name] = cred
	}

	merged := make([]*AuraCredential, 0, len(c.Aura.Credentials))
	for _, cred := range c.Aura.Credentials {
		if c.initialCredentials[cred.Name] && diskCredentialMap[cred.Name] == nil {
			continue
		}
		merged = append(merged, cred)
	}
	c.Aura.Credentials = merged

	for _, diskCred := range diskAura.Credentials {
		found := false
		for _, cred := range c.Aura.Credentials {
			if cred.Name == diskCred.Name {
				found = true
				break
			}
		}
		if !found && !c.initialCredentials[diskCred.Name] {
			c.Aura.Credentials = append(c.Aura.Credentials, diskCred)
		}
	}

	credentialMap := make(map[string]bool)
	for _, cred := range c.Aura.Credentials {
		credentialMap[cred.Name] = true
	}

	if c.Aura.DefaultCredential != "" && !credentialMap[c.Aura.DefaultCredential] {
		c.Aura.DefaultCredential = ""
	}
	if c.Aura.DefaultCredential == "" && diskAura.DefaultCredential != "" && credentialMap[diskAura.DefaultCredential] {
		c.Aura.DefaultCredential = diskAura.DefaultCredential
	}
}
