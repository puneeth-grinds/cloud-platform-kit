package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateLimitingEntry struct {
	Count       int
	WindowStart time.Time
}

type RateLimiter struct {
	requests sync.Map
	rate     int           //maximum requests allowed
	per      time.Duration // time window eg: 1 min
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

		now := time.Now()
		entry := &rateLimitingEntry{
			Count:       1,
			WindowStart: now,
		}
		actual, loaded := rl.requests.LoadOrStore(apiKey, entry)
		entry = actual.(*rateLimitingEntry)

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
			w.WriteHeader(http.StatusTooManyRequests)

			apiError := APIError{
				Error: "Too Many Requests",
				Code: http.StatusTooManyRequests,
			}
			
		} else {
			next.ServeHTTP(w, r)
		}

	})

}
