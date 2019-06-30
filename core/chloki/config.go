package chloki

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	ClickhouseMaxConnections = 100
)

type ClickhouseConfig struct {
	Server         string
	Login          string
	Password       string
	Database       string
	LogTableName   string
	MaxConnections int
}

type Config struct {
	Clickhouse ClickhouseConfig
	ListenAddr string
	UseTLS     bool   `toml:"use_tls"`
	CertFile   string `toml:"cert_file"`
	KeyFile    string `toml:"key_file"`
}

// Construct config struct with default values
func NewConfig() *Config {
	return &Config{
		Clickhouse: ClickhouseConfig{MaxConnections: ClickhouseMaxConnections, Server: "127.0.0.1:8123"},
		ListenAddr: "0.0.0.0:18123",
	}
}

// parse config file in toml format
func ParseConfig(fileName string) *Config {
	_, err := os.Stat(fileName)
	if err != nil {
		log.Fatal("Config file is missing: ", fileName)
	}

	config := NewConfig()
	_, err = toml.DecodeFile(fileName, config)
	if err != nil {
		log.Fatal(err)
	}

	if config.UseTLS {
		if _, err := os.Stat(config.CertFile); err != nil {
			log.Fatal("Cert file is missing: ", config.CertFile)
		}
		if _, err := os.Stat(config.KeyFile); err != nil {
			log.Fatal("Key file is missing: ", config.KeyFile)
		}
	}

	log.Debug("Config: ", config)

	return config
}
