package credentials

import (
	"encoding/json"
	"io"
	"time"

	"github.com/neo4j/cli/common/clierr"
)

type AuraCredentials struct {
	DefaultCredential string            `json:"default-credential"`
	Credentials       []*AuraCredential `json:"credentials"`
	onUpdate          func()
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
	auraCredentials := c.Credentials
	for _, credential := range auraCredentials {
		if credential.Name == name {
			return clierr.NewUsageError("already have credential with name %s", name)
		}
	}

	c.Credentials = append(c.Credentials, &AuraCredential{Name: name, ClientId: clientId, ClientSecret: clientSecret})
	if len(c.Credentials) == 1 {
		c.SetDefault(name)
	}
	c.onUpdate()
	return nil
}

func (c *AuraCredentials) Remove(name string) error {
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
	c.onUpdate()
	return nil
}

func (c *AuraCredentials) SetDefault(name string) error {
	if !c.credentialExists(name) {
		return clierr.NewUsageError("could not find credential with name %s", name)
	}

	c.DefaultCredential = name
	c.onUpdate()
	return nil
}

func (c *AuraCredentials) GetDefault() (*AuraCredential, error) {
	if c.DefaultCredential == "" {
		return nil, clierr.NewUsageError("default credential not set")
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
	credential, err := c.Get(cred.Name)
	if err != nil {
		panic(err)
	}
	const expireToleranceSeconds = 60

	now := time.Now().UnixMilli()

	credential.TokenExpiry = now + (expiresInSeconds-expireToleranceSeconds)*1000
	credential.AccessToken = accessToken
	c.onUpdate()
	return credential
}

func (c *AuraCredentials) ClearAccessToken(cred *AuraCredential) (*AuraCredential, error) {
	credential, err := c.Get(cred.Name)
	if err != nil {
		return nil, err
	}

	credential.TokenExpiry = 0
	credential.AccessToken = ""
	c.onUpdate()
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
	Name         string `json:"name"`
	ClientId     string `json:"client-id"`
	ClientSecret string `json:"client-secret"`
	AccessToken  string `json:"access-token"`
	TokenExpiry  int64  `json:"token-expiry"`
}

func (credential *AuraCredential) HasValidAccessToken() bool {
	now := time.Now().UnixMilli()

	if credential.AccessToken == "" {
		return false
	}

	if now >= credential.TokenExpiry {
		return false
	}

	return true
}
