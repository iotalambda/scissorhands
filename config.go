package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	OpenAIApiKey string `json:"openAIApiKey"`
}

var config *Config

func LoadConfig() error {
	v := viper.New()
	v.SetConfigFile("config.json")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("reading the config file failed: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unmarshalling the config file failed: %v", err)
	}

	config = &cfg

	return nil
}
