package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// logging to capture the status code
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status:  "ok",
		Service: "api-gateway",
	}

	json.NewEncoder(w).Encode(response)

}
func parseLogLevel(value string) slog.Level {

	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}

}
func main() {
	// Load configs
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// slog logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	logger.Info(
		"config loaded successfully",
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
		"scanner_url", cfg.ScannerURL,
	)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	server := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	logger.Info(
		"server starting",
		"addr", server.Addr,
		"service", "api-gateway",
	)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server failed to start", "error", err)
	}
}
