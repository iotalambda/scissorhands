package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	OpenAIApiKey string `json:"openAIApiKey"`
}

var config *Config

func LoadConfigOrErr() {
	v := viper.New()
	v.SetConfigFile("config.json")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("reading the config file failed: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("unmarshalling the config file failed: %v", err)
	}

	config = &cfg
}
