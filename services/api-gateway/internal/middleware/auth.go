package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type APIError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// APIKeyMiddleware blocks requests that do not provide the configured API key
// in the X-API-Key header.
func APIKeyMiddleware(APIKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")

			// Hash both values before comparing so the comparison does not leak
			// timing information about the expected key.
			expectedHash := sha256.Sum256([]byte(APIKey))
			inputHash := sha256.Sum256([]byte(apiKey))
			match := subtle.ConstantTimeCompare(expectedHash[:], inputHash[:])
			
			span := trace.SpanFromContext(r.Context())
			span.SetAttributes(attribute.Bool("auth.api_key.valid", apiKey != "" && match == 1))

			if apiKey == "" || match != 1 {
				span.AddEvent("api_key.authentication_failed")
				span.SetStatus(codes.Error, "unauthorized")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				apiError := APIError{
					Error: "Invalid or missing API Key",
					Code:  http.StatusUnauthorized,
				}
				json.NewEncoder(w).Encode(apiError)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
