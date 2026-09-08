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
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func Load() (Config, error) {
	cfg := Config{
		Port:       getEnv("PORT", "8080"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
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
