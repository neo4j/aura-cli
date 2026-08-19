// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package credentials

import (
	"encoding/json"
	"path/filepath"

	"github.com/neo4j/cli/common/clicfg/fileutils"
	"github.com/spf13/afero"
)

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
		if err := json.Unmarshal(data, &credentials); err != nil {
			panic(err)
		}
	}

	c.Aura = credentials.Aura

	if !fileHasData {
		c.save()
	}
}

func (c *Credentials) save() {
	// Build a structure with secrets revealed for on-disk storage.
	// The stored file must contain the real, unredacted secrets,
	// while display/print paths mask them via Secret.MarshalJSON().
	data, err := json.Marshal(map[string]interface{}{
		"aura": map[string]interface{}{
			"default-credential": c.Aura.DefaultCredential,
			"credentials": func() []map[string]interface{} {
				result := make([]map[string]interface{}, len(c.Aura.Credentials))
				for i, cred := range c.Aura.Credentials {
					result[i] = map[string]interface{}{
						"name":              cred.Name,
						"client-id":         cred.ClientId,
						"client-secret":     cred.ClientSecret.Reveal(),
						"access-token":      cred.AccessToken,
						"token-expiry":      cred.TokenExpiry,
					}
				}
				return result
			}(),
		},
	})
	if err != nil {
		panic(err)
	}

	fileutils.WriteFile(c.fs, c.filePath, data)
}
