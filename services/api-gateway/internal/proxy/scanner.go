package proxy

import (
	"net/http"
)

type ScannerProxy struct {
	baseURL    string 
	httpClient *http.Client
}
