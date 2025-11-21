package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	APIKey string
}

func LoadConfig() (*Config, error) {
	// load .env file if present
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	return &Config{APIKey: apiKey}, nil
}

var ErrNoAPIKey = fmt.Errorf("OPENROUTER_API_KEY not set in environment")
