package config

type LoggerConf struct {
	Level  string `yaml:"level"`
	Source bool   `yaml:"source"`
}

type StorageConf struct {
	Type string        `yaml:"type"`
	DB   DBStorageConf `yaml:"db"`
}

type DBStorageConf struct {
	Dsn string `yaml:"dsn"`
}

type QueueConf struct {
	Type string `yaml:"type"`
	RMQ  RMQ    `yaml:"rmq"`
}

type RMQ struct {
	Dsn   string `yaml:"dsn"`
	QName string `yaml:"qname"`
}
