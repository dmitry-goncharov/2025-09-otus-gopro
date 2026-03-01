package config

import (
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type SchedulerConfig struct {
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
	Queue   QueueConf   `yaml:"queue"`
	Cleaner CleanerConf `yaml:"cleaner"`
	Scanner ScannerConf `yaml:"scanner"`
}

type CleanerConf struct {
	TickInterval time.Duration `yaml:"tickinterval"`
	OutDateTime  time.Duration `yaml:"outdatetime"`
}

type ScannerConf struct {
	TickInterval         time.Duration `yaml:"tickinterval"`
	NotificationInterval time.Duration `yaml:"notificationinterval"`
}

func NewSchedulerConfig(configPath string) (*SchedulerConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening scheduler config file: %w", err)
	}
	defer file.Close()

	config := &SchedulerConfig{}
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("error decoding scheduler config file: %w", err)
	}
	return config, nil
}
