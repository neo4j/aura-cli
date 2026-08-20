// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package credentials

import (
	"encoding/json"
	"io"
	"time"

	"github.com/neo4j/cli/common/clierr"
	"github.com/neo4j/cli/common/redact"
)

type AuraCredentials struct {
	DefaultCredential string            `json:"default-credential"`
	Credentials       []*AuraCredential `json:"credentials"`
}

func (c *AuraCredentials) List() []*AuraCredential {
	return c.Credentials
}

func (config *AuraCredentials) Print(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(config.Credentials); err != nil {
		return err
	}

	return nil
}

func (c *AuraCredentials) GetDefault() (*AuraCredential, error) {
	if c.DefaultCredential == "" {
		return nil, clierr.NewUsageError("default credential not set, please follow the instructions at https://neo4j.com/docs/aura/classic/platform/api/authentication/#_creating_credentials and use the `credential add` subcommand to add the created credentials")
	}
	return c.Get(c.DefaultCredential)
}

func (c *AuraCredentials) Get(name string) (*AuraCredential, error) {
	for _, credential := range c.Credentials {
		if credential.Name == name {
			return credential, nil
		}
	}
	return nil, clierr.NewUsageError("could not find credential with name %s", name)
}

type AuraCredential struct {
	Name         string        `json:"name"`
	ClientId     string        `json:"client-id"`
	ClientSecret redact.Secret `json:"client-secret"`
	AccessToken  redact.Secret `json:"access-token"`
	TokenExpiry  int64         `json:"token-expiry"`
}

func (credential *AuraCredential) HasValidAccessToken() bool {
	now := time.Now().UnixMilli()

	if credential.AccessToken.Reveal() == "" {
		return false
	}

	if now >= credential.TokenExpiry {
		return false
	}

	return true
}

func (c *AuraCredentials) toOnDisk() auraCredentialsOnDisk {
	result := auraCredentialsOnDisk{
		DefaultCredential: c.DefaultCredential,
		Credentials:       make([]auraCredentialOnDisk, len(c.Credentials)),
	}
	for i, cred := range c.Credentials {
		result.Credentials[i] = auraCredentialOnDisk{
			AuraCredential: *cred,
			ClientSecret:   cred.ClientSecret.Reveal(),
			AccessToken:    cred.AccessToken.Reveal(),
		}
	}
	return result
}

func (od auraCredentialsOnDisk) toAuraCredentials() *AuraCredentials {
	result := &AuraCredentials{
		DefaultCredential: od.DefaultCredential,
		Credentials:       make([]*AuraCredential, len(od.Credentials)),
	}
	for i, cred := range od.Credentials {
		newCred := cred.AuraCredential
		newCred.ClientSecret = redact.NewSecret(cred.ClientSecret)
		newCred.AccessToken = redact.NewSecret(cred.AccessToken)
		result.Credentials[i] = &newCred
	}
	return result
}
