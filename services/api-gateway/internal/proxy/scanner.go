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

func ForwardRequest(ctx context)
