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

func NewScanHandler(scannerProxy *proxy.ScannerProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scanHandler(w, r, scannerProxy)
	})
}

func scanHandler(w http.ResponseWriter, r *http.Request, scannerProxy *proxy.ScannerProxy) {
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
