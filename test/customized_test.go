package test

import (
	"rate_limiter/internal"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	// Create a config that limits to 3 requests per second
	config := &internal.RateLimiterConfig{
		Rate:            3,
		Interval:        time.Minute,
		IncludedUserIDs: []string{"user123"},
		ExcludedUserIPs: []string{"127.0.0.1"},
	}

	internal.NewRateLimiter(config)


	if !config.Allow("127.0.0.1") {
		t.Errorf("Expected excluded IP to be allowed")
	}

	
	for i := 0; i < 3; i++ {
		if !config.Allow("user123") {
			t.Errorf("Expected request #%d to be allowed under rate limit", i+1)
		}
	}

	
	if config.Allow("user123") {
		t.Errorf("Expected rate limit to block this request")
	}

	
	time.Sleep(time.Minute)

	if !config.Allow("user123") {
		t.Errorf("Expected request after interval to be allowed again")
	}
}
