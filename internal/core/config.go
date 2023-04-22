package core

import (
	"fmt"

	"github.com/twelvelabs/termite/conf"
)

type Config struct {
	ConfigPath string
	Debug      bool   `yaml:"debug" env:"DEBUG"`
	LogLevel   string `yaml:"log_level" env:"LOG_LEVEL" default:"warn" validate:"oneof=debug info warn error fatal"`
}

// NewTestConfig returns a new Config for unit tests
// populated with default values.
func NewTestConfig() (*Config, error) {
	return NewConfigFromPath("")
}

// NewConfigFromPath returns a new config for the file at path.
// If path is empty, uses `.envctl.yaml`.
func NewConfigFromPath(path string) (*Config, error) {
	config, err := conf.NewLoader(&Config{}, path).Load()
	if err != nil {
		return nil, fmt.Errorf("config load: %w", err)
	}
	config.ConfigPath = path

	return config, nil
}
