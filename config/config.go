package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Host           string `mapstructure:"HOST" default:"0.0.0.0"`
	Port           int    `mapstructure:"PORT" default:"443"`
	KeyFile        string `mapstructure:"KEY_FILE"`
	CertFile       string `mapstructure:"CERT_FILE"`
	Token          string `mapstructure:"TOKEN" required:"true"`
	APIVersion     string `mapstructure:"LATEST_API_VERSION" default:"v19.0"`
	PredictUrl     string `mapstructure:"PREDICT_URL" default:"http://127.0.0.1:5000/predict"`
	QianwenKey     string `mapstructure:"QIANWEN_KEY" required:"true"`
	PageID         string `mapstructure:"PAGE_ID" required:"true"`
	PageAccesToken string `mapstructure:"PAGE_ACCESS_TOKEN" required:"true"`
	AppSecret      string `mapstructure:"APP_SECRET"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error reading configuration file: %w", err)
	}

	var appConf AppConfig
	err := viper.Unmarshal(&appConf)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return &appConf, nil
}
