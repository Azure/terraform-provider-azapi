package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BicepTypesAzCommitSHA      string `yaml:"bicep-types-az-commit-sha"`
	AzureRestAPISpecsCommitSHA string `yaml:"azure-rest-api-specs-commit-sha"`
}

func ReadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.BicepTypesAzCommitSHA == "" {
		return nil, fmt.Errorf("commit SHA not found in %s", configPath)
	}
	if cfg.AzureRestAPISpecsCommitSHA == "" {
		return nil, fmt.Errorf("azure-rest-api-specs commit SHA not found in %s", configPath)
	}

	return &cfg, nil
}
