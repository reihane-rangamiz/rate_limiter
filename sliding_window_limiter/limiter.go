package slidingwindowlimiter

import (
	"sync"
	"time"
)

// Rate     int           	 maximum request count
// Interval time.Duration 	 the time window for rate limiting
// Window   time.Duration 	 the duration of each time slot
// Requests []time.Time 	 keeps track of timestamps of requests
type SlidingWindowRateLimiter struct {
	Rate     int           
	Interval time.Duration 
	Window   time.Duration 
	Mu       sync.Mutex
	Requests []time.Time 
}


func NewSlidingWindowRateLimiter(rate int, interval, window time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		Rate:     rate,
		Interval: interval,
		Window:   window,
		Requests: make([]time.Time, 0),
	}
}

// Allow checks if the request is allowed within the sliding window
func (s *SlidingWindowRateLimiter) AllowSliding() bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	now := time.Now()
	limitTime := now.Add(-s.Interval) 

	
	i := 0
	for _, t := range s.Requests {
		if t.After(limitTime) {
			break
		}
		i++
	}
	s.Requests = s.Requests[i:]

	if len(s.Requests) < s.Rate {
		s.Requests = append(s.Requests, now)
		return true
	}
	return false
}
