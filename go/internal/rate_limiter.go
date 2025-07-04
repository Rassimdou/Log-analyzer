package main

import (
	"net/http"
	"sync"
	"time"
	
	"github.com/gin-gonic/gin"
)

// RateLimiter stores request counters
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.Mutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new instance
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Middleware returns the Gin middleware handler
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		rl.mutex.Lock()
		defer rl.mutex.Unlock()
		
		// Clean up old requests
		now := time.Now()
		var validRequests []time.Time
		for _, t := range rl.requests[ip] {
			if now.Sub(t) <= rl.window {
				validRequests = append(validRequests, t)
			}
		}
		
		// Check if over limit
		if len(validRequests) >= rl.limit {
			c.Header("Retry-After", rl.window.String())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
				"limit": rl.limit,
				"window": rl.window.String(),
			})
			return
		}
		
		// Add current request
		rl.requests[ip] = append(validRequests, now)
		c.Next()
	}
}