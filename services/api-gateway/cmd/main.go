package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/handler"
	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/middleware"
	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/proxy"
)

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
		"ratelimitrpm", cfg.RateLimitRPM,
	)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPM)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.HealthHandler)

	scannerProxy := proxy.NewScannerProxy(cfg.ScannerURL, logger)

	scanHandler := handler.NewScanHandler(scannerProxy)

	rateLimitedScanHandler := rateLimiter.Middleware(scanHandler)

	protectedScanHandler := middleware.APIKeyMiddleware(cfg.APIKey)(rateLimitedScanHandler)

	mux.Handle("POST /scan", protectedScanHandler)

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
		logger.Error("server shutdown failed", "error", err)
	} else {
		logger.Info("server shutdown complete")
	}

}
