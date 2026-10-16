package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIKeyMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     string
		wantStatus int
		wantNext   bool
	}{
		{name: "missing key", wantStatus: http.StatusUnauthorized},
		{name: "wrong key", apiKey: "wrong-key", wantStatus: http.StatusUnauthorized},
		{name: "valid key", apiKey: "correct-key", wantStatus: http.StatusNoContent, wantNext: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusNoContent)
			})

			handler := APIKeyMiddleware("correct-key")(next)
			request := httptest.NewRequest(http.MethodGet, "/scan", nil)
			if test.apiKey != "" {
				request.Header.Set("X-API-Key", test.apiKey)
			}
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, test.wantStatus)
			}
			if nextCalled != test.wantNext {
				t.Fatalf("nextCalled = %t, want %t", nextCalled, test.wantNext)
			}
		})
	}
}
