package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int
	}
}

func GetConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("No configuration file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config

	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("Can't parse yaml: %w", err)
	}

	return &config, nil
}
