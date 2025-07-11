package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Watch struct {
		Resources  []string `mapstructure:"resources"`
		Namespaces []string `mapstructure:"namespaces"`
	} `mapstructure:"watch"`

	Notifier struct {
		Type string `mapstructure:"type"`
	} `mapstructure:"notifier"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
