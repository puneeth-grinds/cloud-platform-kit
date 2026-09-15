package proxy

import (
	"context"
	"errors"
	"io"
	"net/http"
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
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.baseURL+"/scan",
		body,
	)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := p.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded){
			return 0, nil, context.DeadlineExceeded
		}
		return 0, nil, err
	}
	defer resp.Body.Close()
	
}
