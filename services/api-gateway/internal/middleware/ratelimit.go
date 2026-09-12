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

		if now.Sub(entry.WindowStart) >= rl.per {
			entry.Count = 1
			entry.WindowStart = now
		}

		if loaded == true {
			entry.Count++
		} else {
			entry.Count = 1
		}

		if entry.Count > rl.rate {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		} else {
			next.ServeHTTP(w, r)
		}

	})

}
