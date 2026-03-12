package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Logger     LoggerConf  `yaml:"logger"`
	Storage    StorageConf `yaml:"storage"`
	HTTPServer ServerConf  `yaml:"httpserver"`
	GRPCServer ServerConf  `yaml:"grpcserver"`
}

type ServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func NewConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	config := &Config{}
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}
	return config, nil
}
