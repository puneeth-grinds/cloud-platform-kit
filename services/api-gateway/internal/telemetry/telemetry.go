package telemetry

import (
	"context"

	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
)

func Init(ctx context.Context, cfg config.Config) (func(context.Context) error, error) {
	shutdown := func(shutdownCtx context.Context) error {
		return nil
	}
	if cfg.OTELExporterOTLPEndpoint == "" {
		return func(ctx context.Context) error { return nil }, nil
	}
	return shutdown, nil
}
