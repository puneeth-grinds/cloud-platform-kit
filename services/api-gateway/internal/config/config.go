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

func Load()(Config, error) {
	cfg := Config(
		Port: getEnv("PORT", "8080")
		LogLevel: getEnv("Log_Level", "info")
		ScannerURL: os.Getenv("Scanner_URL")
		APIKey: os.Getenv("APIKEY")
	)
}