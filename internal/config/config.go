package config

import (
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

// Config contains linter settings
type Config struct {
	// Patterns for ignoring files and directories
	Ignore []string `yaml:"ignore"`
}

// defaultConfig returns the default configuration
func defaultConfig() *Config {
	return &Config{
		Ignore: []string{
			"**/*_test.go",
			"test/**",
			"**/*_mock.go",
			"**/mock/**",
			"**/mocks/**",
		},
	}
}

// findConfigFile searches for a configuration file in standard locations
func findConfigFile() string {
	candidates := []string{
		".unused-interface-methods.yml",
		"unused-interface-methods.yml",
		".config/unused-interface-methods.yml",
		".unused-interface-methods.yaml",
		"unused-interface-methods.yaml",
		".config/unused-interface-methods.yaml",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}

// matchPattern checks if a file matches the pattern
func (c *Config) matchPattern(pattern, filePath string) bool {
	// Normalize path
	filePath = filepath.Clean(filePath)

	// Use doublestar for pattern matching
	matched, _ := doublestar.Match(pattern, filePath)
	return matched
}

// ShouldIgnore checks if a file or directory should be ignored
func (c *Config) ShouldIgnore(filePath string) bool {
	// Normalize path
	filePath = filepath.Clean(filePath)

	for _, pattern := range c.Ignore {
		if c.matchPattern(pattern, filePath) {
			return true
		}
	}

	return false
}

// LoadConfig loads configuration from a file or returns default configuration
func LoadConfig(configPath string) (*Config, error) {
	// If path is not specified, look in standard locations
	if configPath == "" {
		configPath = findConfigFile()
	}

	// If file is not found, use default configuration
	if configPath == "" {
		return defaultConfig(), nil
	}

	// Check if file exists
	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := defaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
