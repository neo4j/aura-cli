package settings

import (
	"encoding/json"
	"path/filepath"

	"github.com/neo4j/cli/common/clicfg/fileutils"
	"github.com/spf13/afero"
)

type SettingsFile struct {
	Aura *AuraSettings `json:"aura"`
}

type Settings struct {
	fs       afero.Fs
	Aura     *AuraSettings
	filePath string
}

func NewCredentials(fs afero.Fs, configPrefix string) *Settings {
	configPath := filepath.Join(configPrefix, "neo4j", "cli", "settings.json")
	c := Settings{
		fs:       fs,
		filePath: configPath,
	}
	c.load()
	return &c
}

func (c *Settings) load() {
	data := fileutils.ReadFileSafe(c.fs, c.filePath)
	fileHasData := len(data) != 0

	var credentials SettingsFile = SettingsFile{
		Aura: &AuraSettings{
			Settings: []*AuraSetting{},
			onUpdate: c.save,
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

func (c *Settings) save() {
	data, err := json.Marshal(SettingsFile{
		Aura: c.Aura,
	})
	if err != nil {
		panic(err)
	}

	fileutils.WriteFile(c.fs, c.filePath, data)
}
