package tokenbucket

import (
	"sync"
	"time"
)


// capacity   int       	max tokens
// tokens     float64   	current tokens
// refillRate float64   	tokens per second
// lastRefill time.Time 	last time we refilled tokens
type TokenBucketLimiter struct {
	capacity   int       
	tokens     float64   
	refillRate float64   
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucketLimiter(rate int, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity:   burst,
		tokens:     float64(burst),
		refillRate: float64(rate),
		lastRefill: time.Now(),
	}
}

func (l *TokenBucketLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastRefill).Seconds()
	l.tokens += elapsed * l.refillRate
	if l.tokens > float64(l.capacity) {
		l.tokens = float64(l.capacity)
	}
	l.lastRefill = now

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}
