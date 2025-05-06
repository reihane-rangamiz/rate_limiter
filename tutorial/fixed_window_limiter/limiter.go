package fixedwindowlimiter

import (
	"sync"
	"time"
)

type RateLimiter interface {
	Allow() bool
}

type FixedWindowRateLimiter struct {
	Rate     int           // maximum request
	Interval time.Duration // the determined duration
	Mu       sync.Mutex
	Requests []time.Time // recent request tracks
}

func NewFixedWindowRateLimiter(rate int, interval time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		Rate:     rate,
		Interval: interval,
		Requests: make([]time.Time, 0),
	}
}

// Allow() Tracks timestamps of requests.
// On every Allow():
// Removes old timestamps (outside the current window).
// Checks how many remain.(based on the rate)
// If less than rate, allow and add new timestamp.
// Otherwise, deny.
// Uses sync.Mutex to protect data in concurrent access.
func (f *FixedWindowRateLimiter) Allow() bool {
	f.Mu.Lock() // to ensure safe concurrent access to the request
	defer f.Mu.Unlock()

	now := time.Now()
	limitTime := now.Add(-f.Interval) // check the start of the duration which we should consider(calc limit window)

	i := 0
	for _, t := range f.Requests {
		if t.After(limitTime) {
			break
		}
		i++
	}
	f.Requests = f.Requests[i:] // remove the checked request time and keep the rest to examine

	if len(f.Requests) < f.Rate {
		f.Requests = append(f.Requests, now) // we should add the new request to the request list and keep the track of them
		return true
	}

	return false
}

// ----------------------------------------------------------------------

// FIXED WINDOW WITH BURST

// Rate     int          allowed requests per window
// Burst    int           max burst capacity
// Interval time.Duration  time window
// Mu       sync.Mutex
// Requests []time.Time    tracked request times
type FixedWindowWithBurstRateLimiter struct {
	Rate     int          
	Burst    int           
	Interval time.Duration 
	Mu       sync.Mutex
	Requests []time.Time 
}


func NewFixedWindowWithBurstRateLimiter(rate, burst int, interval time.Duration) *FixedWindowWithBurstRateLimiter {
	return &FixedWindowWithBurstRateLimiter{
		Rate:     rate,
		Burst:    burst,
		Interval: interval,
		Requests: make([]time.Time, 0),
	}
}


func (f *FixedWindowWithBurstRateLimiter) AllowBurst() bool {
	f.Mu.Lock()
	defer f.Mu.Unlock()

	now := time.Now()
	limitTime := now.Add(-f.Interval)

	i := 0
	for _, t := range f.Requests {
		if t.After(limitTime) {
			break
		}
		i++
	}
	f.Requests = f.Requests[i:]

	
	if len(f.Requests) < f.Burst {
		f.Requests = append(f.Requests, now)
		return true
	}

	return false
}
