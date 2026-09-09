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

func (sw *statusResponseWriter) WriteHeader(statusCode int) {
	sw.statusCode = statusCode
	sw.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrappedWriter := &statusResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(wrappedWriter, r)

			logger.Info("HTTP Request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrappedWriter.statusCode),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
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

	wrappedmux := loggingMiddleware(logger)(mux)

	mux.HandleFunc("GET /health", healthHandler)
	server := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrappedmux,
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
