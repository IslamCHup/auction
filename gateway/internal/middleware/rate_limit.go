package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenBucket struct{
	capacity int
	tokens float64
	refillRate float64
	lastRefill time.Time
	mutex sync.Mutex
}

func NewTokenBucket(capacity int, refillRate float64) *TokenBucket{
	return &TokenBucket{
		capacity: capacity,
		tokens: float64(capacity),
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	tb.lastRefill = now

	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}

	return false
}

func RateLimitMiddleware() gin.HandlerFunc {
	tb := NewTokenBucket(2000, 100) // 100 запросов в секунду, после 2000 запросов 

	return func(c *gin.Context) {
		if !tb.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
