package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiterLimitsEachAPIKeyIndependently(t *testing.T) {
	limiter := NewRateLimiter(1)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler := limiter.Middleware(next)

	firstRequest := httptest.NewRequest(http.MethodPost, "/scan", nil)
	firstRequest.Header.Set("X-API-Key", "first-key")
	firstResponse := httptest.NewRecorder()
	handler.ServeHTTP(firstResponse, firstRequest)
	if firstResponse.Code != http.StatusNoContent {
		t.Fatalf("first request status = %d, want %d", firstResponse.Code, http.StatusNoContent)
	}

	limitedRequest := httptest.NewRequest(http.MethodPost, "/scan", nil)
	limitedRequest.Header.Set("X-API-Key", "first-key")
	limitedResponse := httptest.NewRecorder()
	handler.ServeHTTP(limitedResponse, limitedRequest)
	if limitedResponse.Code != http.StatusTooManyRequests {
		t.Fatalf("limited request status = %d, want %d", limitedResponse.Code, http.StatusTooManyRequests)
	}
	if limitedResponse.Header().Get("Retry-After") == "" {
		t.Fatal("Retry-After header is empty")
	}

	secondRequest := httptest.NewRequest(http.MethodPost, "/scan", nil)
	secondRequest.Header.Set("X-API-Key", "second-key")
	secondResponse := httptest.NewRecorder()
	handler.ServeHTTP(secondResponse, secondRequest)
	if secondResponse.Code != http.StatusNoContent {
		t.Fatalf("second key status = %d, want %d", secondResponse.Code, http.StatusNoContent)
	}
}
