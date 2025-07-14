package core

import (
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/twelvelabs/termite/conf"

	"github.com/twelvelabs/envctl/internal/dotenv"
)

//go:embed config.default.yaml
var ConfigContentDefault []byte

const (
	ConfigPathDefault = ".envctl.yaml"
	ConfigPathEnv     = "ENVCTL_CONFIG"
)

type Config struct {
	ConfigPath string
	Color      bool   `yaml:"color" env:"ENVCTL_COLOR" default:"true"`
	Debug      bool   `yaml:"debug" env:"ENVCTL_DEBUG"`
	Prompt     bool   `yaml:"prompt" env:"ENVCTL_PROMPT" default:"true"`
	LogLevel   string `yaml:"log_level" env:"ENVCTL_LOG_LEVEL" default:"warn" validate:"oneof=debug info warn error fatal"` //nolint: lll

	Version      string        `yaml:"version"`
	DotEnv       DotEnvConfig  `yaml:"dotenv"`
	Environments []Environment `yaml:"environments"`
}

type DotEnvConfig struct {
	Enabled     bool               `yaml:"enabled" env:"ENVCTL_DOTENV_ENABLED"`
	BasePath    string             `yaml:"base_path" env:"ENVCTL_DOTENV_BASE_PATH"`
	QuoteStyle  dotenv.QuoteStyle  `yaml:"quote_style" env:"ENVCTL_DOTENV_QUOTE_STYLE" default:"double" validate:"oneof=none single double"` //nolint: lll
	EscapeStyle dotenv.EscapeStyle `yaml:"escape_style" env:"ENVCTL_DOTENV_ESCAPE_STYLE" default:"default" validate:"oneof=default compose"` //nolint: lll
}

func (c *Config) EnvironmentNames() []string {
	names := []string{}
	for _, env := range c.Environments {
		names = append(names, env.Name)
	}
	return names
}

// NewTestConfig returns a new Config for unit tests
// populated with default values.
func NewTestConfig() (*Config, error) {
	return NewConfigFromPath("")
}

// NewConfigFromPath returns a new config for the file at path.
func NewConfigFromPath(path string) (*Config, error) {
	config, err := conf.NewLoader(&Config{}, path).Load()
	if err != nil {
		return nil, fmt.Errorf("config load: %w", err)
	}
	config.ConfigPath = path
	return config, nil
}

// ConfigPath resolves and returns the config path.
// Lookup order:
//   - Flag
//   - Environment variable
//   - Default path name
func ConfigPath(args []string) (string, error) {
	path := ConfigPathDefault
	if p := os.Getenv(ConfigPathEnv); p != "" {
		path = p
	}

	// Create a minimal, duplicate flag set to determine just the config path
	// (the remaining flags are defined on the cobra.Command flag set).
	// Using two different sets because Cobra doesn't parse flags until _after_
	// we have instantiated the app (and thus the Config).
	fs := pflag.NewFlagSet("config-args", pflag.ContinueOnError)
	fs.StringVarP(&path, "config", "c", path, "")
	// Ignore all the flags used by the main Cobra flagset.
	fs.ParseErrorsWhitelist.UnknownFlags = true
	// Suppress the default usage shown when the `--help` flag is present
	// (otherwise we end up w/ a duplicate of what Cobra shows).
	fs.Usage = func() {}

	err := fs.Parse(args)
	if err != nil && !errors.Is(err, pflag.ErrHelp) {
		return "", fmt.Errorf("unable to parse config flag: %w", err)
	}

	return path, nil
}
