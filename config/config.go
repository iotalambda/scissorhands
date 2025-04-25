package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AzureSpeechServiceKey string `json:"azureSpeechServiceKey"`
	OpenAIApiKey          string `json:"openAIApiKey"`
}

var Global *Config

func InitConfig() error {
	v := viper.New()
	v.SetConfigFile("config.json")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("reading the config file failed: %v", err)
	}

	if err := v.Unmarshal(&Global); err != nil {
		return fmt.Errorf("unmarshalling the config file failed: %v", err)
	}

	return nil
}
