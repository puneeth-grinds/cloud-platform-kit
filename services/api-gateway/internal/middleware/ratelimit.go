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

		val, found := rl.requests.Load(apiKey)
		if ! found{
			entry := &rateLimitingEntry{
				Count: 0,
				WindowStart: time.Now(),
			}
			actual, _ := rl.requests.LoadOrStore(apiKey, entry)
			entry = actual.(*rateLimitingEntry)
		}

		next.ServeHTTP(w, r)

	})

}
