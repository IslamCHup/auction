package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenBucket struct {
	capacity   int
	tokens     float64
	refillRate float64
	lastRefill time.Time
	mutex      sync.Mutex
}

func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     float64(capacity),
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

var UserBucket sync.Map
var BidBucket sync.Map

func getOrCreateBucket(store *sync.Map, userID uint64, capacity int, refillRate float64) *TokenBucket {
	actual, _ := store.LoadOrStore(userID, NewTokenBucket(capacity, refillRate))
	return actual.(*TokenBucket)
}

func UserRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uidAny, exist := c.Get("user_id")
		if !exist {
			c.Next()
			return
		}
		uid := uidAny.(uint64)
		bucket := getOrCreateBucket(&UserBucket, uid, 100, 100.0/60.0)

		if !bucket.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}
		c.Next()
	}
}

func BidRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasSuffix(c.FullPath(), "/bids") && !strings.Contains(c.FullPath(), "/bids/") {
			c.Next()
			return
		}
		uidAny, exist := c.Get("user_id")
		if !exist {
			c.Next()
			return
		}
		uid := uidAny.(uint64)
		bucket := getOrCreateBucket(&BidBucket, uid, 10, 10.0/60.0)

		if !bucket.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}
		c.Next()
	}
}
