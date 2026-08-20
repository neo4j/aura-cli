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
	refresh           func() error
	persist           func() error
}

// refreshAndPersist re-reads the on-disk state into c immediately before mutate runs, so mutate
// always applies its delta on top of the current file rather than a snapshot that may have gone
// stale since load. It only persists if mutate succeeds.
func (c *AuraCredentials) refreshAndPersist(mutate func() error) error {
	if err := c.refresh(); err != nil {
		return err
	}
	if err := mutate(); err != nil {
		return err
	}
	return c.persist()
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

func (c *AuraCredentials) Add(name string, clientId string, clientSecret string) error {
	return c.refreshAndPersist(func() error {
		for _, credential := range c.Credentials {
			if credential.Name == name {
				return clierr.NewUsageError("already have credential with name %s", name)
			}
		}

		c.Credentials = append(c.Credentials, &AuraCredential{
			Name:         name,
			ClientId:     clientId,
			ClientSecret: redact.NewSecret(clientSecret),
		})
		if len(c.Credentials) == 1 {
			c.DefaultCredential = name
		}
		return nil
	})
}

func (c *AuraCredentials) Remove(name string) error {
	return c.refreshAndPersist(func() error {
		var indexToRemove = -1

		for i, credential := range c.Credentials {
			if credential.Name == name {
				indexToRemove = i
				break
			}
		}

		if indexToRemove == -1 {
			return clierr.NewUsageError("could not find credential with name %s to remove", name)
		}

		if c.DefaultCredential == name {
			c.DefaultCredential = ""
		}

		c.Credentials = append(c.Credentials[:indexToRemove], c.Credentials[indexToRemove+1:]...)
		return nil
	})
}

func (c *AuraCredentials) SetDefault(name string) error {
	return c.refreshAndPersist(func() error {
		if !c.credentialExists(name) {
			return clierr.NewUsageError("could not find credential with name %s", name)
		}

		c.DefaultCredential = name
		return nil
	})
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

func (c *AuraCredentials) UpdateAccessToken(cred *AuraCredential, accessToken string, expiresInSeconds int64) *AuraCredential {
	var credential *AuraCredential
	err := c.refreshAndPersist(func() error {
		var err error
		credential, err = c.Get(cred.Name)
		if err != nil {
			return err
		}

		const expireToleranceSeconds = 60
		now := time.Now().UnixMilli()

		credential.TokenExpiry = now + (expiresInSeconds-expireToleranceSeconds)*1000
		credential.AccessToken = redact.NewSecret(accessToken)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return credential
}

func (c *AuraCredentials) ClearAccessToken(cred *AuraCredential) (*AuraCredential, error) {
	var credential *AuraCredential
	err := c.refreshAndPersist(func() error {
		var err error
		credential, err = c.Get(cred.Name)
		if err != nil {
			return err
		}

		credential.TokenExpiry = 0
		credential.AccessToken = redact.NewSecret("")
		return nil
	})
	if err != nil {
		return nil, err
	}
	return credential, nil
}

func (c *AuraCredentials) credentialExists(name string) bool {
	for _, credential := range c.Credentials {
		if credential.Name == name {
			return true
		}
	}
	return false
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
