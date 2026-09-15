package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
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

		now := time.Now()

		// Load the existing counter for this API key, or create the first counter
		// when this key is seen for the first time.
		entry := &rateLimitingEntry{
			Count:       1,
			WindowStart: now,
		}
		actual, loaded := rl.requests.LoadOrStore(apiKey, entry)
		entry = actual.(*rateLimitingEntry)
		entry.mu.Lock()

		// Existing keys either start a new window or increment the current one.
		if loaded {
			if now.Sub(entry.WindowStart) >= rl.per {
				entry.Count = 1
				entry.WindowStart = now
			} else {
				entry.Count++
			}
		}

		if entry.Count > rl.rate {
			w.Header().Set("Content-Type", "application/json")

			// Retry-After tells the client how many seconds remain before the
			// current rate-limit window resets.
			windowEndsAt := entry.WindowStart.Add(rl.per)
			retryAfter := windowEndsAt.Sub(now)
			retryAfterSec := retryAfter.Seconds()

			if retryAfterSec < 1 {
				retryAfterSec = 1
			}
			entry.mu.Unlock()

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
