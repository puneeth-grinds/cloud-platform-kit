package telemetry

import (
	"context"
	"fmt"
	"strings"

	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
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
	
	headersMap := parseHeaders(cfg.OTELExporterOTLPHeaders)

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpointURL(cfg.OTELExporterOTLPEndpoint),
		otlptracegrpc.WithHeaders(headersMap),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otlp exporter %w", err)
	}

	shutdown := func(shutdownCtx context.Context) error {
		return nil
	}

	return shutdown, nil
}

func parseHeaders(rawHeaders string) map[string]string {
	headers := make(map[string]string)
	if strings.TrimSpace(rawHeaders) == "" {
		return headers
	}

	// Using SplitN with 2 elements prevents tokens containing "=" from fragmenting.
	parts := strings.SplitN(rawHeaders, "=", 2)
	if len(parts) == 2 {
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key != "" {
			headers[key] = value
		}
	}

	return headers
}