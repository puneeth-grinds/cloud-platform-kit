package middleware

import {
	"time"
}
type rateLimitingEntry struct {
	Count       int `json:"count"`
	WindowStart time.Time `json:"windowStart"`
}
