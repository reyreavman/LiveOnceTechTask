package ratelimiter

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     int
	window   time.Duration
	stopChan chan struct{}
}

type tokenBucket struct {
	tokens    int
	lastReset time.Time
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*tokenBucket),
		rate:     rate,
		window:   window,
		stopChan: make(chan struct{}),
	}

	go rl.cleanupRoutine()
	return rl
}

func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[clientID]
	if !exists {
		rl.buckets[clientID] = &tokenBucket{
			tokens:    rl.rate - 1,
			lastReset: time.Now(),
		}
		return true
	}

	elapsed := time.Since(bucket.lastReset)
	if elapsed >= rl.window {
		bucket.tokens = rl.rate - 1
		bucket.lastReset = time.Now()
		return true
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}
	return false
}

func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanupExpired()
		case <-rl.stopChan:
			return
		}
	}
}

func (rl *RateLimiter) cleanupExpired() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	expiryTime := time.Now().Add(-2 * rl.window)
	for id, bucket := range rl.buckets {
		if bucket.lastReset.Before(expiryTime) {
			delete(rl.buckets, id)
		}
	}
}

func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}