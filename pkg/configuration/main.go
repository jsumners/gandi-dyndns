package configuration

import (
	"fmt"
	"os"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

type Record struct {
	Zone string `mapstructure:"zone"`
	Type string `mapstructure:"type"`
	Name string `mapstructure:"name"`
	TTL int `mapstructure:"ttl" default:"300"`
}

type Configuration struct {
	GandiV5ApiKey string `mapstructure:"gandi_v5_api_key"`
	Records []Record `mapstructure:"records"`
}

func New() *Configuration {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix("GDYNDNS")
	v.AddConfigPath("/etc/gdyndns")
	v.AddConfigPath("$HOME/.gdyndns")

	configFile := os.Getenv("GDYNDNS_CONFIG_FILE")
	if configFile != "" {
		v.SetConfigFile(configFile)
	}

	viperError := v.ReadInConfig()
	if viperError != nil {
		_, isFileNotFoundError := viperError.(viper.ConfigFileNotFoundError)
		if isFileNotFoundError && configFile != "" {
			panic(fmt.Errorf("fatal error loading configuration (%s): %w", configFile, viperError))
		} else if !isFileNotFoundError {
			panic(fmt.Errorf("fatal error loading configuration: %w", viperError))
		}
	}

	result := &Configuration{}
	err := v.Unmarshal(result)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshaling configuration: %w", err))
	}
	defaults.SetDefaults(result)

	return result
}
