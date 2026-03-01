package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type SenderConfig struct {
	Logger LoggerConf `yaml:"logger"`
	Queue  QueueConf  `yaml:"queue"`
}

func NewSenderConfig(configPath string) (*SenderConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening sender config file: %w", err)
	}
	defer file.Close()

	config := &SenderConfig{}
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("error decoding sender config file: %w", err)
	}
	return config, nil
}
