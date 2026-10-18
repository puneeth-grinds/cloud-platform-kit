package telemetry

import (
	"context"
	"fmt"
	"strings"

	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	shutdown := func(shutdownCtx context.Context) error {
		return nil
	}

	return tp.Shutdown, nil
}

func parseHeaders(rawHeaders string) map[string]string {
	headers := make(map[string]string)
	if strings.TrimSpace(rawHeaders) == "" {
		return headers
	}

	// First, split the entire configuration string by commas to get individual header strings
	pairs := strings.Split(rawHeaders, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		// Use SplitN with 2 elements on each pair so tokens containing "=" remain uncorrupted [1]
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				headers[key] = value
			}
		}
	}

	return headers
}
