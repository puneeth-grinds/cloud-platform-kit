package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Port         string
	LogLevel     string
	ScannerURL   string
	APIKey       string
	RateLimitRPM int
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
	rateLimitRPMValue := getEnv("RATE_LIMIT_RPM", "60")
	rateLimitRPMInt, err := strconv.Atoi(rateLimitRPMValue)
	if err != nil {
		return Config{}, err
	}

	if rateLimitRPMInt <= 0 {
		return Config{}, errors.New("RATE_LIMIT_RPM must be greater than 0")
	}
	cfg := Config{
		Port:       getEnv("PORT", "8080"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		ScannerURL: os.Getenv("SCANNER_URL"),
		APIKey:     os.Getenv("API_KEY"),
	}
	rateLimitRPMInt, err := strconv.Atoi(cfg.RateLimitRPM)
	if err != nil {
		panic(err)
	}

	if rateLimitRPMInt <= 0 {
		return Config{}, errors.New("RateLimitRPM value should be greater than 0")
	}

	if cfg.ScannerURL == "" {
		return Config{}, errors.New("SCANNER_URL is required")
	}
	if cfg.APIKey == "" {
		return Config{}, errors.New("API_KEY is required")
	}
	return cfg, nil
}
