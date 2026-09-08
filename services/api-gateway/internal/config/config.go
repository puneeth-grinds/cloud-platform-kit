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
