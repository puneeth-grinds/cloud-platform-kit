package middleware

type rateLimitingEntry struct {
	Count int `json:"count"`
	WindowStart int `json:"windowStart"`
}