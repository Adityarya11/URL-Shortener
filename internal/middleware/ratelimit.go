package middleware

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RateLimiter is a simple in-memory token-bucket limiter per client (IP or API key).
type RateLimiter struct {
	clients     map[string]*client
	mu          sync.Mutex
	capacity    int64         // max tokens
	refillDur   time.Duration // window duration
	refillRate  float64       // tokens per second
	cleanupTick time.Duration
}

type client struct {
	tokens     float64
	lastRefill time.Time
	lastSeen   time.Time
}

// NewRateLimiterFromEnv constructs limiter using env vars:
// RATE_LIMIT_REQUESTS (int), RATE_LIMIT_WINDOW (duration string, e.g. "1h")
func NewRateLimiterFromEnv() *RateLimiter {
	limitStr := os.Getenv("RATE_LIMIT_REQUESTS")
	windowStr := os.Getenv("RATE_LIMIT_WINDOW")

	limit := int64(100) // default
	if v, err := strconv.ParseInt(limitStr, 10, 64); err == nil && v > 0 {
		limit = v
	}

	window := 1 * time.Hour
	if d, err := time.ParseDuration(windowStr); err == nil {
		window = d
	}

	rl := &RateLimiter{
		clients:     make(map[string]*client),
		capacity:    limit,
		refillDur:   window,
		refillRate:  float64(limit) / window.Seconds(),
		cleanupTick: window, // run cleanup every window
	}

	go rl.cleanupLoop()
	return rl
}

// Middleware wraps an http.Handler and enforces rate limits per client IP.
// It sets headers: X-RateLimit-Limit, X-RateLimit-Remaining, Retry-After (when limited)
func (r *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := r.clientKey(req)

		allowed, remaining, retryAfter := r.allow(key)
		// Set headers for clients
		w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(r.capacity, 10))
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(int64(remaining), 10))

		if !allowed {
			// 429 Too Many Requests
			w.Header().Set("Retry-After", strconv.FormatInt(int64(retryAfter.Seconds()), 10))
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, req)
	})
}

// allow checks and consumes one token for the given key.
// returns: allowed, remainingTokens, retryAfterDuration
func (r *RateLimiter) allow(key string) (bool, int64, time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	c, ok := r.clients[key]
	if !ok {
		c = &client{
			tokens:     float64(r.capacity),
			lastRefill: now,
			lastSeen:   now,
		}
		r.clients[key] = c
	}

	// refill tokens
	elapsed := now.Sub(c.lastRefill).Seconds()
	if elapsed > 0 {
		c.tokens += elapsed * r.refillRate
		if c.tokens > float64(r.capacity) {
			c.tokens = float64(r.capacity)
		}
		c.lastRefill = now
	}

	c.lastSeen = now

	if c.tokens >= 1 {
		c.tokens -= 1
		remaining := int64(c.tokens)
		return true, remaining, 0
	}

	// not enough tokens -> calculate retry after (time needed to acquire 1 token)
	needed := 1.0 - c.tokens
	seconds := needed / r.refillRate
	retry := time.Duration(seconds * float64(time.Second))
	return false, int64(c.tokens), retry
}

// clientKey returns a stable key to rate limit by.
// Prefers X-Forwarded-For (first IP) then RemoteAddr IP.
func (r *RateLimiter) clientKey(req *http.Request) string {
	// If you later support API keys, check Authorization / X-Api-Key here first.
	xff := req.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil && ip != "" {
		return ip
	}

	// fallback
	return req.RemoteAddr
}

// cleanupLoop periodically removes clients not seen for > 2 * window
func (r *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(r.cleanupTick)
	defer ticker.Stop()
	for range ticker.C {
		threshold := time.Now().Add(-2 * r.refillDura())
		r.mu.Lock()
		for k, c := range r.clients {
			if c.lastSeen.Before(threshold) {
				delete(r.clients, k)
			}
		}
		r.mu.Unlock()
	}
}

func (r *RateLimiter) refillDura() time.Duration {
	// approximate refill duration equal to configured window
	return r.refillDur
}
