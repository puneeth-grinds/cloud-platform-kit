package middleware

import (
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

func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
	return &RateLimiter{
		rate: rate,
		per:  per,
	}
}
