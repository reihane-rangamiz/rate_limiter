package ratelimiter

import (
	"testing"
	"time"
)

func TestFixedWindowRateLimiter(t *testing.T) {
	// sec:= time.Duration(time.Second)*10
	newLimiter := NewFixedWindowRateLimiter(3, time.Second)

	for i := 0; i < 3; i++ {
		if !newLimiter.Allow() {
			t.Errorf("request must have been allowed since we didnt hit the rate")
		}
	}

	// this is more than rate, If it returns true, it means the limiter failed to block the extra request (→ That's a bug). hence we wait a second and then check allow 
	if newLimiter.Allow() { //negative test 
		// Bad! We should NOT be allowed to send the 4th request in the same window
		t.Errorf("request should have been blocked since we reached the rate limit")
	}

	time.Sleep(time.Second)

	if !newLimiter.Allow() {
		// now we reset the window and first 3 requests supposed to be allowed 
		t.Errorf("Request after window reset should be allowed")
	}
}
