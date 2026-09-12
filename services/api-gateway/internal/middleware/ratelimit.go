package middleware

import (
	"time"
)

type rateLimitingEntry struct {
	Count       int
	WindowStart time.Time
}
