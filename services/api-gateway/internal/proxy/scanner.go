package proxy

import (
	"net/http"
)

type ScannerProxy struct {
	baseURL string
	client  *http.Client
}

func NewProxyService(baseURL ScannerProxy) *ScannerProxy {

}
