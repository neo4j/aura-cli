package clicfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var ConfigPrefix string

const DefaultAuraBaseUrl = "https://api.neo4j.io/v1"
const DefaultAuraAuthUrl = "https://api.neo4j.io/oauth/token"

func NewConfig(fs afero.Fs, version string) (*Config, error) {
	configPath := filepath.Join(ConfigPrefix, "neo4j", "cli")

	Viper := viper.New()

	Viper.SetFs(fs)
	Viper.SetConfigName("config")
	Viper.SetConfigType("json")
	Viper.AddConfigPath(configPath)
	Viper.SetConfigPermissions(0600)

	bindEnvironmentVariables(Viper)
	setDefaultValues(Viper)

	if err := Viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := fs.MkdirAll(configPath, 0755); err != nil {
				return nil, err
			}
			if err = Viper.SafeWriteConfig(); err != nil {
				return nil, err
			}
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	return &Config{Version: version, Aura: AuraConfig{viper: Viper, pollingOverride: PollingConfig{
		MaxRetries: 60,
		Interval:   20,
	}}}, nil
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

type Config struct {
	Version string
	Aura    AuraConfig
}

type PollingConfig struct {
	Interval   int
	MaxRetries int
}

type AuraConfig struct {
	viper           *viper.Viper
	pollingOverride PollingConfig
}

func (config *AuraConfig) Get(key string) interface{} {
	return config.viper.Get(fmt.Sprintf("aura.%s", key))
}

func (config *AuraConfig) Set(key string, value string) error {
	config.viper.Set(fmt.Sprintf("aura.%s", key), value)
	return config.viper.WriteConfig()
}

func (config *AuraConfig) Print(cmd *cobra.Command) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(config.viper.Get("aura")); err != nil {
		return err
	}

	return nil
}

func (config *AuraConfig) BaseUrl() string {
	return config.viper.GetString("aura.base-url")
}

func (config *AuraConfig) BindBaseUrl(flag *pflag.Flag) error {
	return config.viper.BindPFlag("aura.base-url", flag)
}

func (config *AuraConfig) AuthUrl() string {
	return config.viper.GetString("aura.auth-url")
}

func (config *AuraConfig) BindAuthUrl(flag *pflag.Flag) error {
	return config.viper.BindPFlag("aura.auth-url", flag)
}

func (config *AuraConfig) Output() string {
	return config.viper.GetString("aura.output")
}

func (config *AuraConfig) BindOutput(flag *pflag.Flag) error {
	return config.viper.BindPFlag("aura.output", flag)
}

func (config *AuraConfig) DefaultTenant() string {
	return config.viper.GetString("aura.default-tenant")
}

func (config *AuraConfig) DefaultCredential() (*AuraCredential, error) {
	auraConfig := auraConfig{}

	if err := config.viper.UnmarshalKey("aura", &auraConfig); err != nil {
		return nil, err
	}

	defaultCredential := config.viper.GetString("aura.default-credential")

	if defaultCredential == "" {
		return nil, errors.New("no default credential found")
	}

	for index := range auraConfig.Credentials {
		var credential = &(auraConfig.Credentials[index])
		if credential.Name == defaultCredential {
			credential.viper = config.viper
			return credential, nil
		}
	}

	return nil, fmt.Errorf("could not find credential with name %s", defaultCredential)
}

func (config *AuraConfig) PollingConfig() PollingConfig {
	return config.pollingOverride
}

func (config *AuraConfig) SetPollingConfig(maxRetries int, interval int) {
	config.pollingOverride = PollingConfig{
		MaxRetries: maxRetries,
		Interval:   interval,
	}
}

func (config *AuraConfig) SetDefaultCredential(name string) error {
	auraConfig := auraConfig{}
	config.viper.Sub("aura").Unmarshal(&auraConfig)

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

	config.viper.Set("aura.default-credential", name)
	return config.viper.WriteConfig()
}

type auraConfig struct {
	Credentials []AuraCredential
}

type AuraCredential struct {
	viper        *viper.Viper
	Name         string `mapstructure:"name" json:"name"`
	ClientId     string `mapstructure:"client-id" json:"client-id"`
	ClientSecret string `mapstructure:"client-secret" json:"client-secret"`
	AccessToken  string `mapstructure:"access-token" json:"access-token"`
	TokenExpiry  int64  `mapstructure:"token-expiry" json:"token-expiry"`
}

func (config *AuraConfig) AddCredential(name string, clientId string, clientSecret string) error {
	auraConfig := auraConfig{}
	config.viper.Sub("aura").Unmarshal(&auraConfig)

	for _, credential := range auraConfig.Credentials {
		if credential.Name == name {
			return fmt.Errorf("already have credential with name %s", name)
		}
	}

	auraConfig.Credentials = append(auraConfig.Credentials, AuraCredential{Name: name, ClientId: clientId, ClientSecret: clientSecret})
	config.viper.Set("aura.credentials", auraConfig.Credentials)
	if len(auraConfig.Credentials) == 1 {
		config.viper.Set("aura.default-credential", name)
	}
	return config.viper.WriteConfig()
}

func (config *AuraConfig) RemoveCredential(name string) error {
	auraConfig := auraConfig{}
	config.viper.Sub("aura").Unmarshal(&auraConfig)

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

	if config.viper.Get("aura.default-credential") == name {
		config.viper.Set("aura.default-credential", "")
	}

	auraConfig.Credentials = append(auraConfig.Credentials[:indexToRemove], auraConfig.Credentials[indexToRemove+1:]...)

	config.viper.Set("aura.credentials", auraConfig.Credentials)

	return config.viper.WriteConfig()
}

func (credential *AuraCredential) IsAccessTokenValid() bool {
	now := time.Now().UnixMilli()

	if credential.AccessToken == "" {
		return false
	}

	if now >= credential.TokenExpiry {
		return false
	}

	return true
}

func (credential *AuraCredential) UpdateAccessToken(accessToken string, expiresInSeconds int64) error {
	const expireToleranceSeconds = 60

	now := time.Now().UnixMilli()

	credential.TokenExpiry = now + (expiresInSeconds-expireToleranceSeconds)*1000
	credential.AccessToken = accessToken

	auraConfig := auraConfig{}
	credential.viper.Sub("aura").Unmarshal(&auraConfig)

	for i, c := range auraConfig.Credentials {
		if c.Name == credential.Name {
			auraConfig.Credentials[i].AccessToken = accessToken
			auraConfig.Credentials[i].TokenExpiry = expiresInSeconds
			break
		}
	}

	credential.viper.Set("aura.credentials", auraConfig.Credentials)

	return credential.viper.WriteConfig()
}
