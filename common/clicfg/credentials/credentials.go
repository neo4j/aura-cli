// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package credentials

import (
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/neo4j/cli/common/clicfg/fileutils"
	"github.com/neo4j/cli/common/clierr"
	"github.com/neo4j/cli/common/redact"
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
	aura, err := c.readAuraFresh()
	if err != nil {
		panic(err)
	}

	c.Aura = aura

	data := fileutils.ReadFileSafe(c.fs, c.filePath)
	if len(data) == 0 {
		c.writeAura(c.Aura)
	}
}

func (c *Credentials) readAuraFresh() (*AuraCredentials, error) {
	data := fileutils.ReadFileSafe(c.fs, c.filePath)

	if len(data) == 0 {
		return &AuraCredentials{
			Credentials: []*AuraCredential{},
		}, nil
	}

	var onDisk credentialsFileOnDisk
	if err := json.Unmarshal(data, &onDisk); err != nil {
		return nil, err
	}

	return onDisk.Aura.toAuraCredentials(), nil
}

func (c *Credentials) writeAura(aura *AuraCredentials) error {
	onDisk := credentialsFileOnDisk{
		Aura: aura.toOnDisk(),
	}

	data, err := json.Marshal(onDisk)
	if err != nil {
		panic(err)
	}

	fileutils.WriteFile(c.fs, c.filePath, data)
	c.Aura = aura
	return nil
}

func (c *Credentials) Add(name string, clientId string, clientSecret string) error {
	aura, err := c.readAuraFresh()
	if err != nil {
		return err
	}

	for _, credential := range aura.Credentials {
		if credential.Name == name {
			return clierr.NewUsageError("already have credential with name %s", name)
		}
	}

	aura.Credentials = append(aura.Credentials, &AuraCredential{
		Name:         name,
		ClientId:     clientId,
		ClientSecret: redact.NewSecret(clientSecret),
	})
	if len(aura.Credentials) == 1 {
		aura.DefaultCredential = name
	}

	return c.writeAura(aura)
}

func (c *Credentials) Remove(name string) error {
	aura, err := c.readAuraFresh()
	if err != nil {
		return err
	}

	var indexToRemove = -1
	for i, credential := range aura.Credentials {
		if credential.Name == name {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		return clierr.NewUsageError("could not find credential with name %s to remove", name)
	}

	if aura.DefaultCredential == name {
		aura.DefaultCredential = ""
	}

	aura.Credentials = append(aura.Credentials[:indexToRemove], aura.Credentials[indexToRemove+1:]...)

	return c.writeAura(aura)
}

func (c *Credentials) SetDefault(name string) error {
	aura, err := c.readAuraFresh()
	if err != nil {
		return err
	}

	if !c.credentialExists(name, aura) {
		return clierr.NewUsageError("could not find credential with name %s", name)
	}

	aura.DefaultCredential = name

	return c.writeAura(aura)
}

func (c *Credentials) UpdateAccessToken(cred *AuraCredential, accessToken string, expiresInSeconds int64) *AuraCredential {
	aura, err := c.readAuraFresh()
	if err != nil {
		panic(err)
	}

	credential, err := aura.Get(cred.Name)
	if err != nil {
		panic(err)
	}

	const expireToleranceSeconds = 60
	now := time.Now().UnixMilli()

	credential.TokenExpiry = now + (expiresInSeconds-expireToleranceSeconds)*1000
	credential.AccessToken = redact.NewSecret(accessToken)

	err = c.writeAura(aura)
	if err != nil {
		panic(err)
	}

	return credential
}

func (c *Credentials) ClearAccessToken(cred *AuraCredential) (*AuraCredential, error) {
	aura, err := c.readAuraFresh()
	if err != nil {
		return nil, err
	}

	credential, err := aura.Get(cred.Name)
	if err != nil {
		return nil, err
	}

	credential.TokenExpiry = 0
	credential.AccessToken = redact.NewSecret("")

	err = c.writeAura(aura)
	if err != nil {
		return nil, err
	}

	return credential, nil
}

func (c *Credentials) credentialExists(name string, aura *AuraCredentials) bool {
	for _, credential := range aura.Credentials {
		if credential.Name == name {
			return true
		}
	}
	return false
}
