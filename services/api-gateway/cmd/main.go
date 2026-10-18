package main

import (
	"context"
	"log"
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
	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/telemetry"
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
	// Create a root context that is cancelled when the process receives
	// SIGINT or SIGTERM. This is what triggers graceful shutdown.
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
		"version", cfg.Version,
	)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPM)

	if cfg.OTELExporterOTLPEndpoint == "" {
		logger.Warn("Telementry is disabled: OLTP endpoint is not configured")
	} else {
		shutdown, err := telemetry.Init(ctx, cfg)
		if err != nil {
			log.Fatal("Failed to initialize telemetry","error",err)
		}
	}

	mux := http.NewServeMux()

	// Health is public because load balancers and local checks must be able to
	// call it without an API key.
	mux.HandleFunc("GET /health", handler.NewHealthHandler(cfg.Version))

	scannerProxy := proxy.NewScannerProxy(cfg.ScannerURL, logger)

	// /scan flow: auth checks the API key, rate limit checks request volume,
	// then the handler forwards the request to the vulnerability-scanner.
	scanHandler := handler.NewScanHandler(scannerProxy)
	rateLimitedScanHandler := rateLimiter.Middleware(scanHandler)
	protectedScanHandler := middleware.APIKeyMiddleware(cfg.APIKey)(rateLimitedScanHandler)
	mux.Handle("POST /scan", protectedScanHandler)

	// Logging wraps the whole mux so every route gets one structured request log.
	wrappedMux := middleware.LoggingMiddleware(logger)(mux)

	server := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrappedMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 45 * time.Second,
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

	if err := shutdown(shutdownCtx); err != nil {
		logger.Error("telemetry shutdown failed", "error", err)
	} else {
		logger.Info("telemetry shutdown complete")
	}
}
