package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"net/http"
)

type APIError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func APIKeyMiddleware(APIKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")

			expectedHash := sha256.Sum256([]byte(APIKey))
			inputHash := sha256.Sum256([]byte(apiKey))

			match := subtle.ConstantTimeCompare(expectedHash[:], inputHash[:])

			if apiKey == "" || match != 1 {
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
