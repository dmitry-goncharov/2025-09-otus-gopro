package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
	Server  ServerConf  `yaml:"server"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type StorageConf struct {
	Type string        `yaml:"type"`
	DB   DBStorageConf `yaml:"db"`
}

type DBStorageConf struct {
	Dsn string `yaml:"dsn"`
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
