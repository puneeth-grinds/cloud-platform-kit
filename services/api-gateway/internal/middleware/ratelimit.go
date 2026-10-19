package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/attributes"
)

type rateLimitingEntry struct {
	Count       int
	WindowStart time.Time
	mu          sync.Mutex
}

type RateLimiter struct {
	requests sync.Map
	rate     int
	per      time.Duration
}

// NewRateLimiter creates a fixed-window limiter using requests per minute.
func NewRateLimiter(rate int) *RateLimiter {
	return &RateLimiter{
		rate: rate,
		per:  time.Minute,
	}
}

// Middleware tracks request counts per API key and returns 429 when the key has
// used all requests in the current one-minute window.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		allowed, remaining := rl.check(apiKey, time.Now())
		trace.SpanFromContext(r.Context()).SetAttributes(
			attributes.Int("rate_limit.remaining", remaining),
		)

			w.Header().Set("Retry-After", strconv.Itoa(int(retryAfterSec)))
			w.WriteHeader(http.StatusTooManyRequests)

			apiError := APIError{
				Error: "Too Many Requests",
				Code:  http.StatusTooManyRequests,
			}
			json.NewEncoder(w).Encode(apiError)
			return
		}

		entry.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) check(apiKey string, now time.Time)(allowed bool, remaining int){
	entry := &rateLimitingEntry{Count: 1, WindowStart: now}
	actual, loaded := rl.requests.LoadOrStore(apiKey, entry)
	entry = actual.(*rateLimitingEntry)

	if loaded{
		if now.Sub(entry.WindowStart) >= rl.per{
			entry.Count = 1
			entry.WindowStart = now
		} else {
			entry.Count++
		}
	}
	remaining = rl.rate - entry.Count
	if remaining< 0 {
		remaining = 0
	}
	return entry.Count <= rl.rate, remaining
}