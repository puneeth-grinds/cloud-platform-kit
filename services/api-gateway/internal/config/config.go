package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	// Port is the HTTP port the api-gateway listens on.
	Port string
	// LogLevel controls how much detail slog writes.
	LogLevel string
	// ScannerURL is the base URL for the vulnerability-scanner service.
	ScannerURL string
	// APIKey is the shared secret expected in the X-API-Key header.
	APIKey string
	// RateLimitRPM is the max requests per minute allowed per API key.
	RateLimitRPM int
	// Version identifies the deployed application revision.
	Version string

	OTELExporterOTLPEndpoint string `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"`

	OTELExporterOTLPHeaders string `envconfig:"OTEL_EXPORTER_OTLP_HEADERS"`

	ServiceName string `envconfig:"SERVICE_NAME" default:"api-gateway"`

	ServiceVersion string `envconfig:"SERVICE_VERSION" default:"0.1.0"`

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
	// Environment variables are strings, so convert RATE_LIMIT_RPM before
	// storing it in Config.
	rateLimitRPMValue := getEnv("RATE_LIMIT_RPM", "60")
	rateLimitRPMInt, err := strconv.Atoi(rateLimitRPMValue)
	if err != nil {
		return Config{}, err
	}

	if rateLimitRPMInt <= 0 {
		return Config{}, errors.New("RATE_LIMIT_RPM must be greater than 0")
	}

	// Optional values use safe defaults. Required values use os.Getenv and are
	// validated below.
	cfg := Config{
		Port:         getEnv("PORT", "8080"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		ScannerURL:   os.Getenv("SCANNER_URL"),
		APIKey:       os.Getenv("API_KEY"),
		RateLimitRPM: rateLimitRPMInt,
		Version:      getEnv("APP_VERSION", "dev"),
		OTELExporterOTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OTELExporterOTLPHeaders: os.Getenv("OTEL_EXPORTER_OTLP_HEADERS"),
		ServiceName: getEnv("SERVICE_NAME","api-gateway"),
		ServiceVersion: getEnv("SERVICE_VERSION","0.1.0"),
	}

	if cfg.ScannerURL == "" {
		return Config{}, errors.New("SCANNER_URL is required")
	}
	if cfg.APIKey == "" {
		return Config{}, errors.New("API_KEY is required")
	}

	return cfg, nil
}
