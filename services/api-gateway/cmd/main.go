package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type ScanResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// statusResponseWriter wraps the real response writer so middleware can record
// the status code written by the handler.
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

			// Default to 200 because Go sends that status when a handler writes
			// a body without explicitly calling WriteHeader.
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

// healthHandler is used by local checks and the load balancer to confirm that
// the api-gateway process is running.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status:  "ok",
		Service: "api-gateway",
	}

	json.NewEncoder(w).Encode(response)

}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := ScanResponse{
		Status:  "accepted",
		Service: "api-gateway",
	}
	json.NewEncoder(w).Encode(response)
}

// parseLogLevel converts the LOG_LEVEL string from config into slog's typed
// log level value.
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
	// create context that listens for the SIGNINT signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// Load config before starting the server so missing required environment
	// variables fail fast.
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Use JSON logs so ECS and CloudWatch receive structured fields.
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
	mux.HandleFunc("GET /scan", scanHandler )

	wrappedMux := loggingMiddleware(logger)(mux)

	server := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrappedMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info(
			"server starting",
			"addr", server.Addr,
			"service", "api-gateway",
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	logger.Info("server is shutting down gracefully")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("error shutdown failed", "error", err)
	} else {
		logger.Info("server shutdown complete")
	}

}
