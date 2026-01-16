package settings

import (
	"encoding/json"
	"io"

	"github.com/neo4j/cli/common/clierr"
)

type AuraSettings struct {
	DefaultSetting string         `json:"default-setting"`
	Settings       []*AuraSetting `json:"settings"`
	onUpdate       func()
}

func (c *AuraSettings) List() []*AuraSetting {
	return c.Settings
}

func (c *AuraSettings) Print(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(c.Settings); err != nil {
		return err
	}

	return nil
}

func (c *AuraSettings) Add(name string, organizationId string, projectId string) error {
	auraSettings := c.Settings
	for _, setting := range auraSettings {
		if setting.Name == name {
			return clierr.NewUsageError("already have setting with name %s", name)
		}
	}

	c.Settings = append(c.Settings, &AuraSetting{Name: name, OrganizationId: organizationId, ProjectId: projectId})
	if len(c.Settings) == 1 {
		c.SetDefault(name)
	}
	c.onUpdate()
	return nil
}

func (c *AuraSettings) Remove(name string) error {
	var indexToRemove = -1

	for i, setting := range c.Settings {
		if setting.Name == name {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		return clierr.NewUsageError("could not find setting with name %s to remove", name)
	}

	if c.DefaultSetting == name {
		c.DefaultSetting = ""
	}

	c.Settings = append(c.Settings[:indexToRemove], c.Settings[indexToRemove+1:]...)
	c.onUpdate()
	return nil
}

func (c *AuraSettings) SetDefault(name string) error {
	if !c.settingExists(name) {
		return clierr.NewUsageError("could not find setting with name %s", name)
	}

	c.DefaultSetting = name
	c.onUpdate()
	return nil
}

func (c *AuraSettings) GetDefault() (*AuraSetting, error) {
	if c.DefaultSetting == "" {
		return nil, clierr.NewUsageError("default setting not set, use the `setting add` subcommand to add a new setting")
	}
	return c.Get(c.DefaultSetting)
}

func (c *AuraSettings) Get(name string) (*AuraSetting, error) {
	for _, setting := range c.Settings {
		if setting.Name == name {
			return setting, nil
		}
	}
	return nil, clierr.NewUsageError("could not find credential with name %s", name)
}

func (c *AuraSettings) settingExists(name string) bool {
	for _, setting := range c.Settings {
		if setting.Name == name {
			return true
		}
	}
	return false
}

type AuraSetting struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organization-id"`
	ProjectId      string `json:"project-id"`
}
