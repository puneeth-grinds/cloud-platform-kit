package proxy

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
)

type ScannerProxy struct {
	baseURL           string
	client            *http.Client
	logger            *slog.Logger
	scannerErrorCount metric.Int64Counter
}

// NewScannerProxy stores the scanner base URL and creates an HTTP client with a
// timeout so api-gateway does not wait forever on the downstream service.
func NewScannerProxy(scannerURL string, logger *slog.Logger) *ScannerProxy {
	httpClient := http.Client{Timeout: 30 * time.Second}

	proxy := ScannerProxy{
		baseURL: scannerURL,
		client:  &httpClient,
		logger:  logger,
	}
	return &proxy
}

// Forward sends the caller's scan request body to the vulnerability-scanner and
// returns the scanner's status code and response body.
func (p *ScannerProxy) Forward(ctx context.Context, body io.Reader, contentType string) (int, []byte, error) {
	ctx, span := otel.Tracer("api-gateway/proxy").Start(ctx, "proxy.scan")
	defer span.End()

	start := time.Now()

	targetURL := strings.TrimRight(p.baseURL, "/") + "/scan"
	span.SetAttributes(attribute.String("scanner.url", targetURL))

	// Reuse the caller's context so cancelled client requests also cancel the
	// downstream scanner request.
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		targetURL,
		body,
	)
	if err != nil {
		return 0, nil, err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	if req.ContentLength >= 0 {
		span.SetAttributes(attribute.Int64("http.request.body.size", req.ContentLength))
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := p.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "scanner request failed")

		p.logger.ErrorContext(
			ctx,
			"scanner request failed",
			"error", err,
			"duration", time.Since(start),
		)
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, nil, context.DeadlineExceeded
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return 0, nil, context.DeadlineExceeded
		}
		return 0, nil, err
	}
	span.SetAttributes(attribute.Int("http.response.status_code", resp.StatusCode))

	defer resp.Body.Close()

	duration := time.Since(start)
	p.logger.InfoContext(
		ctx,
		"scanner request complete",
		"status", resp.StatusCode,
		"duration", duration,
	)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, responseBody, nil
}
