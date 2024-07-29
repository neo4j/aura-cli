package clicfg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var configPrefix string

const DefaultAuraBaseUrl = "https://api.neo4j.io/v1"
const DefaultAuraAuthUrl = "https://api.neo4j.io/oauth/token"

type Config struct {
	viper *viper.Viper
	out   *bufio.Writer
	Aura  AuraConfig `mapstructure:"aura" json:"aura"`
}

func (config *Config) Get(key string) (interface{}, error) {
	val := config.viper.Get(key)

	if val == nil {
		return nil, fmt.Errorf("could not find config value with key %s", key)
	}

	return val, nil
}

func (config *Config) Set(key string, value interface{}) error {
	config.viper.Set(key, value)

	err := config.viper.Unmarshal(config)
	if err != nil {
		return err
	}

	return nil
}

func (config *Config) Write() error {
	content, err := json.Marshal(config)
	if err != nil {
		return err
	}

	n, err := config.out.Write(content)
	if err != nil {
		return err
	}

	if err = config.out.Flush(); err != nil {
		return err
	}

	fmt.Printf("wrote %d bytes to config file\n", n)

	return nil
}

func (config *Config) BindPFlag(key string, flag *pflag.Flag) error {
	return config.viper.BindPFlag(key, flag)
}

type AuraConfig struct {
	BaseUrl           string           `mapstructure:"base-url" json:"base-url"`
	AuthUrl           string           `mapstructure:"auth-url" json:"auth-url"`
	Output            string           `mapstructure:"output" json:"output"`
	DefaultTenant     string           `mapstructure:"default-tenant" json:"default-tenant,omitempty"`
	DefaultCredential string           `mapstructure:"default-credential" json:"default-credential,omitempty"`
	Credentials       []AuraCredential `mapstructure:"credentials" json:"credentials"`
}

func (auraConfig *AuraConfig) AddCredential(name string, clientId string, clientSecret string) error {
	for _, credential := range auraConfig.Credentials {
		if credential.Name == name {
			return fmt.Errorf("already have credential with name %s", name)
		}
	}

	auraConfig.Credentials = append(auraConfig.Credentials, AuraCredential{Name: name, ClientId: clientId, ClientSecret: clientSecret})

	return nil
}

func (auraConfig *AuraConfig) RemoveCredential(name string) error {
	var indexToRemove = -1

	for i, credential := range auraConfig.Credentials {
		if credential.Name == name {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		return fmt.Errorf("could not find credential with name %s to remove", name)
	}

	if auraConfig.DefaultCredential == name {
		auraConfig.DefaultCredential = ""
	}

	auraConfig.Credentials = append(auraConfig.Credentials[:indexToRemove], auraConfig.Credentials[indexToRemove+1:]...)

	return nil
}

func (auraConfig *AuraConfig) SetDefaultCredential(name string) error {
	var credentialExists = false

	for _, credential := range auraConfig.Credentials {
		if credential.Name == name {
			credentialExists = true
			break
		}
	}

	if !credentialExists {
		return fmt.Errorf("could not find credential with name %s", name)
	}

	auraConfig.DefaultCredential = name

	return nil
}

func (auraConfig *AuraConfig) Print(cmd *cobra.Command) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(auraConfig); err != nil {
		return err
	}

	return nil
}

type AuraCredential struct {
	Name         string `mapstructure:"name" json:"name"`
	ClientId     string `mapstructure:"client-id" json:"client-id"`
	ClientSecret string `mapstructure:"client-secret" json:"client-secret"`
}

func NewConfig() (*Config, error) {
	configPath := filepath.Join(configPrefix, "neo4j", "cli", "config.json")

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return &Config{}, err
	}
	defer f.Close()

	fi, err := os.Stat(configPath)

	if err != nil {
		return &Config{}, err
	}

	out := bufio.NewWriter(f)

	if fi.Size() == 0 {
		if _, err := out.Write([]byte("{}")); err != nil {
			return nil, err
		}
		if err := out.Flush(); err != nil {
			return nil, err
		}
	}

	in, err := os.Open(configPath)
	if err != nil {
		return &Config{}, err
	}
	defer in.Close()

	return NewConfigFrom(in, out)
}

func NewConfigFrom(in io.Reader, out *bufio.Writer) (*Config, error) {
	Viper := viper.New()

	Viper.SetConfigType("json")

	bindEnvironmentVariables(Viper)
	setDefaultValues(Viper)

	if err := Viper.ReadConfig(in); err != nil {
		return &Config{}, err
	}

	var config Config
	err := Viper.Unmarshal(&config)

	if err != nil {
		return &Config{}, err
	}

	config.viper = Viper
	config.out = out

	return &config, nil
}

func bindEnvironmentVariables(Viper *viper.Viper) {
	Viper.BindEnv("aura.base-url", "AURA_BASE_URL")
	Viper.BindEnv("aura.auth-url", "AURA_AUTH_URL")
}

func setDefaultValues(Viper *viper.Viper) {
	Viper.SetDefault("aura.base-url", DefaultAuraBaseUrl)
	Viper.SetDefault("aura.auth-url", DefaultAuraAuthUrl)
	Viper.SetDefault("aura.output", "json")
	Viper.SetDefault("aura.credentials", []AuraCredential{})
}
