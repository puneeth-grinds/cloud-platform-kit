package config

import (
	"errors"
	"os"
)

type Config struct {
	Port       string
	LogLevel   string
	ScannerURL string
	APIKey     string
	RateLimitRPM string
}

// getEnv reads optional environment variables that have safe defaults.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Load builds the application config from environment variables and fails fast
// when required values are missing.
func Load() (Config, error) {
	cfg := Config{
		Port:       getEnv("PORT", "8080"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		RateLimitRPM: getEnv("RATE_LIMIT_RPM","60"),
		ScannerURL: os.Getenv("SCANNER_URL"),
		APIKey:     os.Getenv("API_KEY"),

	}
	if cfg.ScannerURL == "" {
		return Config{}, errors.New("SCANNER_URL is required")
	}
	if cfg.APIKey == "" {
		return Config{}, errors.New("API_KEY is required")
	}
	return cfg, nil
}
