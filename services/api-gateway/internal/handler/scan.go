package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/proxy"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// NewScanHandler creates the /scan HTTP handler and injects the scanner proxy
// dependency it needs to forward requests.
func NewScanHandler(scannerProxy *proxy.ScannerProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scanHandler(w, r, scannerProxy)
	})
}

func scanHandler(w http.ResponseWriter, r *http.Request, scannerProxy *proxy.ScannerProxy) {
	// Only JSON scan requests are accepted because the scanner expects JSON.
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		errorResponse := ErrorResponse{
			Error: "content type must be application/json",
			Code:  http.StatusUnsupportedMediaType,
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Forward the original request body to the vulnerability-scanner and return
	// the scanner response back to the caller.
	statusCode, respBytes, err := scannerProxy.Forward(r.Context(), r.Body, contentType)
	w.Header().Set("Content-Type", "application/json")

	if errors.Is(err, context.DeadlineExceeded) {
		errorResponse := ErrorResponse{
			Error: "Gateway Timeout",
			Code:  http.StatusGatewayTimeout,
		}
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err != nil {
		errorResponse := ErrorResponse{
			Error: "Bad Gateway",
			Code:  http.StatusBadGateway,
		}
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(respBytes)
}
