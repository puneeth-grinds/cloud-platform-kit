package config

import (
	"errors"
	"os"

	"golang.org/x/tools/go/cfg"
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

func Load ()(Config, error){
	cfg := {
		Port: getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		ScannerURL: os.Getenv("SCANNER_URL"),
		APIKey: os.Getenv("API_KEY")
	}
}
