package middleware

import (
	"net/http"
)

type APIError struct {
	Error string `json:"error`
	Code  int `json:"code"`
}


func APIKeyMiddleware(APIKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			w.Header().Set("Content-Type", "application/json")

			apierror := APIError{
				Error: http.Error(http.StatusUnauthorized),
				Code: http.StatusUnauthorized,
			}

			if apiKey == "" || apiKey != APIKey {
				return 
			}
			next.ServeHTTP(w, r)
		})
	}
}
