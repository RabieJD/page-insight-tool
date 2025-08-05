package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

// Environment represents environment-specific configuration
type Environment struct {
	Port string `yaml:"Port"`
}

// Config represents the application configuration
type Config struct {
	ServerAddress string
}

var envs map[string]Environment

// LoadConfig loads configuration from YAML file or environment variables
func LoadConfig(configFile string) *Config {
	cfg := &Config{
		ServerAddress: ":8080",
	}

	// Load from YAML file if provided
	if configFile != "" {
		if yamlConfig := loadYamlConfig(configFile); yamlConfig != nil {
			cfg = yamlConfig
		}
	}

	// Override with environment variables if available
	if port := os.Getenv("PORT"); port != "" {
		cfg.ServerAddress = ":" + port
	}

	return cfg
}

// loadYamlConfig loads configuration from YAML file
func loadYamlConfig(filename string) *Config {
	if filename == "" {
		filename, _ = filepath.Abs("./app/config/page-insight-tool.yaml")
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	err = yaml.Unmarshal(yamlFile, &envs)
	if err != nil {
		return nil
	}

	// Default to Local environment
	env := "Local"
	// Override with environment if available
	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		env = appEnv
	}
	if envs[env].Port != "" {
		return &Config{
			ServerAddress: ":" + envs[env].Port,
		}
	}
	return nil
}
