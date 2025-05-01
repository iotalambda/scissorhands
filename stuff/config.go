package stuff

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AzureSpeechServiceKey string `json:"azureSpeechServiceKey"`
	OpenAIApiKey          string `json:"openAIApiKey"`
}

var GlobalConfig *Config

func InitConfig() error {
	v := viper.New()
	v.SetConfigFile("config.json")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read the config file: %v", err)
	}

	if err := v.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("unmarshal the config file: %v", err)
	}

	return nil
}
