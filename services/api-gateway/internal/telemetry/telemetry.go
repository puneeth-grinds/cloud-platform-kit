package telemetry

import (
	"context"
	"fmt"

	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

func Init(ctx context.Context, cfg config.Config) (func(context.Context) error, error) {
	if cfg.OTELExporterOTLPEndpoint == "" {
		return func(ctx context.Context) error { return nil }, nil
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
		),
		resource.WithTelemetrySDK(),

	)
	if err != nil {
		return nil, fmt.Errorf("failed to created OTEL resource %w", err)
	}

	shutdown := func(shutdownCtx context.Context) error {
		return nil
	}

	return shutdown, nil
}
