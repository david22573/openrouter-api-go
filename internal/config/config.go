package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
}

var ErrNoAPIKey = fmt.Errorf("OPENROUTER_API_KEY not set in config or environment")

func LoadConfig() (*Config, error) {
	v := viper.New()

	// --- Defaults ---
	v.SetDefault("model", "openrouter/llama3.1")

	// --- Read YAML config ---
	v.SetConfigName("config") // config.yml
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.config/myapp")
	_ = v.ReadInConfig() // optional; ignore error if no file

	// --- Environment variable override ---
	v.AllowEmptyEnv(false)
	v.BindEnv("api_key", "OPENROUTER_API_KEY")

	// --- Unmarshal into struct ---
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// --- Validation ---
	if cfg.APIKey == "" {
		return nil, ErrNoAPIKey
	}

	return &cfg, nil
}