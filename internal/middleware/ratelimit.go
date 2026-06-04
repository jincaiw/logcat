package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiterEntry struct {
	count    int
	expireAt time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	entries  map[string]*rateLimiterEntry
	max      int
	window   time.Duration
	stopCh   chan struct{}
	stopOnce sync.Once
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*rateLimiterEntry),
		max:     maxRequests,
		window:  window,
		stopCh:  make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Stop() {
	rl.stopOnce.Do(func() {
		close(rl.stopCh)
	})
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for k, v := range rl.entries {
				if now.After(v.expireAt) {
					delete(rl.entries, k)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopCh:
			return
		}
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	entry, ok := rl.entries[key]
	if !ok || now.After(entry.expireAt) {
		rl.entries[key] = &rateLimiterEntry{
			count:    1,
			expireAt: now.Add(rl.window),
		}
		return true
	}
	if entry.count >= rl.max {
		return false
	}
	entry.count++
	return true
}

var (
	defaultAPILimiter    *RateLimiter
	defaultLoginLimiter  *RateLimiter
	defaultInitLimiter   *RateLimiter
)

func init() {
	defaultAPILimiter = NewRateLimiter(300, time.Minute)
	defaultLoginLimiter = NewRateLimiter(10, time.Minute)
	defaultInitLimiter = NewRateLimiter(5, time.Minute)
}

func StopRateLimiters() {
	if defaultAPILimiter != nil {
		defaultAPILimiter.Stop()
	}
	if defaultLoginLimiter != nil {
		defaultLoginLimiter.Stop()
	}
	if defaultInitLimiter != nil {
		defaultInitLimiter.Stop()
	}
}

func APIRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !defaultAPILimiter.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "too many requests, please try again later",
			})
			return
		}
		c.Next()
	}
}

func LoginRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !defaultLoginLimiter.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "too many login attempts, please try again later",
			})
			return
		}
		c.Next()
	}
}

func InitRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !defaultInitLimiter.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "too many attempts, please try again later",
			})
			return
		}
		c.Next()
	}
}

func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(maxRequests, window)
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !limiter.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "too many requests, please try again later",
			})
			return
		}
		c.Next()
	}
}
