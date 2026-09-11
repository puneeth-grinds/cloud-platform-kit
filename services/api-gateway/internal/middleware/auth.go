package middleware

import (
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

			if apiKey == "" || apiKey != APIKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				apiError := APIError{
					Error: "Invalid or missing API Key",
					Code:  http.StatusUnauthorized,
				}
				json.NewEncoder(w).Encode(apiError)
			}

		})
	}
}
