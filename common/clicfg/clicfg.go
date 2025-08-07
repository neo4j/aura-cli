package clicfg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"slices"

	"github.com/neo4j/cli/common/clicfg/credentials"
	"github.com/neo4j/cli/common/clicfg/fileutils"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

var ConfigPrefix string

const (
	DefaultAuraBaseUrl     = "https://api.neo4j.io"
	DefaultAuraBetaPathV1  = "v1beta5"
	DefaultAuraBetaPathV2  = "v2beta1"
	DefaultAuraAuthUrl     = "https://api.neo4j.io/oauth/token"
	DefaultAuraBetaEnabled = false
)

var ValidOutputValues = [3]string{"default", "json", "table"}

type Config struct {
	Version     string
	Aura        *AuraConfig
	Credentials *credentials.Credentials
}

func NewConfig(fs afero.Fs, version string) *Config {
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
				panic(err)
			}
			if err = Viper.SafeWriteConfig(); err != nil {
				panic(err)
			}
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}

	credentials := credentials.NewCredentials(fs, ConfigPrefix)

	return &Config{
		Version: version,
		Aura: &AuraConfig{
			fs:    fs,
			viper: Viper, pollingOverride: PollingConfig{
				MaxRetries: 60,
				Interval:   20,
			},
			ValidConfigKeys: []string{"auth-url", "base-url", "default-tenant", "output", "beta-enabled"},
		},
		Credentials: credentials,
	}
}

func bindEnvironmentVariables(Viper *viper.Viper) {
	Viper.BindEnv("aura.base-url", "AURA_BASE_URL")
	Viper.BindEnv("aura.auth-url", "AURA_AUTH_URL")
}

func setDefaultValues(Viper *viper.Viper) {
	Viper.SetDefault("aura.base-url", DefaultAuraBaseUrl)
	Viper.SetDefault("aura.auth-url", DefaultAuraAuthUrl)
	Viper.SetDefault("aura.output", "default")
	Viper.SetDefault("aura.beta-enabled", DefaultAuraBetaEnabled)
}

type AuraConfig struct {
	viper           *viper.Viper
	fs              afero.Fs
	pollingOverride PollingConfig
	ValidConfigKeys []string
}

type PollingConfig struct {
	Interval   int
	MaxRetries int
}

func (config *AuraConfig) IsValidConfigKey(key string) bool {
	return slices.Contains(config.ValidConfigKeys, key)
}

func (config *AuraConfig) Get(key string) interface{} {
	return config.viper.Get(fmt.Sprintf("aura.%s", key))
}

func (config *AuraConfig) Set(key string, value string) {
	filename := config.viper.ConfigFileUsed()
	data := fileutils.ReadFileSafe(config.fs, filename)

	updateConfig, err := sjson.Set(string(data), fmt.Sprintf("aura.%s", key), value)
	if err != nil {
		panic(err)
	}

	baseUrlToUpdate := ""
	if key == "base-url" {
		baseUrlToUpdate = value
	}
	log.Printf("updating aura base url: %s", updateConfig)
	updatedAuraBaseUrl := config.auraBaseUrlOnConfigChange(baseUrlToUpdate)
	log.Printf("updated aura base url: %s", updatedAuraBaseUrl)
	if updatedAuraBaseUrl != "" {
		intermediateUpdateConfig, err := sjson.Set(string(updateConfig), "aura.base-url", updatedAuraBaseUrl)
		if err != nil {
			panic(err)
		}
		updateConfig = intermediateUpdateConfig
		log.Printf("updated aura config: %s", updateConfig)
	}

	fileutils.WriteFile(config.fs, filename, []byte(updateConfig))
}

func (config *AuraConfig) Print(cmd *cobra.Command) {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(config.viper.Get("aura")); err != nil {
		panic(err)
	}
}

func (config *AuraConfig) BaseUrl() string {
	originalUrl := config.viper.Get("aura.base-url")
	return removePathParametersFromUrl(originalUrl.(string))
}

func removePathParametersFromUrl(originalUrl string) string {
	parsedUrl, err := url.Parse(originalUrl)
	if err != nil {
		panic(err)
	}
	log.Printf("aura.base-url: %s", parsedUrl.Host)
	log.Printf("aura.auth-url: %s", parsedUrl.Scheme)
	return fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)
}

func (config *AuraConfig) BetaPathV1() string {
	return DefaultAuraBetaPathV1
}

func (config *AuraConfig) BetaPathV2() string {
	return DefaultAuraBetaPathV2
}

func (config *AuraConfig) BindBaseUrl(flag *pflag.Flag) {
	if err := config.viper.BindPFlag("aura.base-url", flag); err != nil {
		panic(err)
	}
}

func (config *AuraConfig) AuthUrl() string {
	return config.viper.GetString("aura.auth-url")
}

func (config *AuraConfig) BindAuthUrl(flag *pflag.Flag) {
	if err := config.viper.BindPFlag("aura.auth-url", flag); err != nil {
		panic(err)
	}
}

func (config *AuraConfig) Output() string {
	return config.viper.GetString("aura.output")
}

func (config *AuraConfig) BindOutput(flag *pflag.Flag) {
	if err := config.viper.BindPFlag("aura.output", flag); err != nil {
		panic(err)
	}
}

func (config *AuraConfig) AuraBetaEnabled() bool {
	return config.viper.GetBool("aura.beta-enabled")
}

func (config *AuraConfig) DefaultTenant() string {
	return config.viper.GetString("aura.default-tenant")
}

func (config *AuraConfig) Fs() afero.Fs {
	return config.fs
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

func (config *AuraConfig) auraBaseUrlOnConfigChange(url string) string {
	if url == "" {
		return DefaultAuraBaseUrl
	}
	return removePathParametersFromUrl(url)
}
