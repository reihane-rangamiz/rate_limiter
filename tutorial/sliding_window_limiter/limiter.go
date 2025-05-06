package slidingwindowlimiter

import (
	"sync"
	"time"
)

type SlidingWindowRateLimiter struct {
	mu         sync.Mutex
	maxReq     int           // max requests allowed
	windowSize time.Duration // e.g. 1 second
	timestamps []time.Time   // when requests happened
}

func NewSlidingWindowRateLimiter(maxReq int, window time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		maxReq:     maxReq,
		windowSize: window,
		timestamps: make([]time.Time, 0),
	}
}

func (r *SlidingWindowRateLimiter) AllowSliding() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.windowSize)

	i := 0
	for _, ts := range r.timestamps {
		if ts.After(windowStart) {
			break
		}
		i++
	}
	r.timestamps = r.timestamps[i:] // drop old ones

	if len(r.timestamps) < r.maxReq {
		r.timestamps = append(r.timestamps, now)
		return true
	}
	return false
}
