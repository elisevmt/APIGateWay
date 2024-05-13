package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"fmt"
)

type Config struct {
	System struct {
		MaxGoroutines string `yaml:"max_goroutines"`
	} `yaml:"system"`
	Rabbit struct {
		LogPublisher struct {
			Url       string `yaml:"url"`
			QueueName string `yaml:"queueName"`
		} `yaml:"logPublisher"`
	} `yaml:"Rabbit"`
	Server struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Version string `yaml:"version"`
	} `yaml:"server"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
		Sslmode  string `yaml:"sslmode"`
	} `yaml:"postgres"`
}

func LoadConfig(path string) (*Config, error) {
	// #nosec
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file not found %s", path)
	}
	yamlDecoder := yaml.NewDecoder(file)

	cfg := &Config{}
	err = yamlDecoder.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("%v", err.Error())
	}

	return cfg, nil
}

