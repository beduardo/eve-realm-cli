package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// HostConfig holds runtime configuration values for the CLI host binary.
type HostConfig struct {
	MCPServerAddr string `yaml:"mcp_server_addr"`
}

// LoadHostConfig reads a YAML file at path and returns a populated HostConfig.
// If the file does not exist, a zero-value HostConfig and nil error are returned.
func LoadHostConfig(path string) (HostConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return HostConfig{}, nil
		}
		return HostConfig{}, err
	}

	var cfg HostConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return HostConfig{}, err
	}
	return cfg, nil
}

// Resolve returns the value of the environment variable named envKey when it is
// non-empty. Otherwise it returns yamlValue.
func Resolve(envKey, yamlValue string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return yamlValue
}
