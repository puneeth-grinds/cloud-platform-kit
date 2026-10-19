package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

func NewRateLimiter(rate int) *RateLimiter {
	return &RateLimiter{
		rate: rate,
		per:  time.Minute,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		allowed, remaining, retryAfterSec := rl.check(apiKey, time.Now())

		trace.SpanFromContext(r.Context()).SetAttributes(
			attribute.Int("rate_limit.remaining", remaining),
		)

		if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", strconv.Itoa(retryAfterSec))
			w.WriteHeader(http.StatusTooManyRequests)

			apiError := APIError{
				Error: "Too Many Requests",
				Code:  http.StatusTooManyRequests,
			}
			_ = json.NewEncoder(w).Encode(apiError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) check(
	apiKey string,
	now time.Time,
) (allowed bool, remaining int, retryAfterSec int) {
	entry := &rateLimitingEntry{Count: 1, WindowStart: now}
	actual, loaded := rl.requests.LoadOrStore(apiKey, entry)
	entry = actual.(*rateLimitingEntry)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if loaded {
		if now.Sub(entry.WindowStart) >= rl.per {
			entry.Count = 1
			entry.WindowStart = now
		} else {
			entry.Count++
		}
	}

	remaining = rl.rate - entry.Count
	if remaining < 0 {
		remaining = 0
	}

	allowed = entry.Count <= rl.rate
	if !allowed {
		retryAfter := entry.WindowStart.Add(rl.per).Sub(now).Seconds()
		retryAfterSec = int(retryAfter)
		if retryAfterSec < 1 {
			retryAfterSec = 1
		}
	}

	return allowed, remaining, retryAfterSec
}