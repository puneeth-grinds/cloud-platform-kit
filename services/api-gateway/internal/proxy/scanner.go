package proxy

import (
	"net/http"
)

type ScannerProxy struct {
	baseURL string
	client  *http.Client
}
