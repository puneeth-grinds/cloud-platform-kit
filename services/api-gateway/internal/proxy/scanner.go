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
)

type ScannerProxy struct {
	baseURL string
	client  *http.Client
	logger *slog.Logger
}

func NewScannerProxy(scannerURL string, logger *slog.Logger) *ScannerProxy {
	httpClient := http.Client{Timeout: 30 * time.Second}

	proxy := ScannerProxy{
		baseURL: scannerURL,
		client:  &httpClient,
		logger: logger,
	}
	return &proxy
}

func (p *ScannerProxy, ) Forward(ctx context.Context, body io.Reader, contentType string) (int, []byte, error) {
	trimmedBaseURL := strings.TrimRight(p.baseURL, "/")
	start := time.Now()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		trimmedBaseURL+"/scan",
		body,
	)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := p.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, nil, context.DeadlineExceeded
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return 0, nil, context.DeadlineExceeded
		}
		return 0, nil, err
	}
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
