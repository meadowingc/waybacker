package site

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessKey string   `yaml:"access_key"`
	SecretKey string   `yaml:"secret_key"`
	SleepDays int      `yaml:"sleep_days"`
	URLs      []string `yaml:"urls"`
}

func ReadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
