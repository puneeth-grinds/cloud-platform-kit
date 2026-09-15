package proxy

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type ScannerProxy struct {
	baseURL string
	client  *http.Client
}

func NewScannerProxy(scannerURL string) *ScannerProxy {
	httpClient := http.Client{Timeout: 30 * time.Second}

	proxy := ScannerProxy{
		baseURL: scannerURL,
		client:  &httpClient,
	}
	return &proxy
}

func (p *ScannerProxy) Forward(ctx context.Context, body io.Reader, contentType string) (int, []byte, error) {
	trimmedBaseURL := strings.TrimRight(p.baseURL, "/")
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
		return 0, nil, err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, responseBody, nil

}
