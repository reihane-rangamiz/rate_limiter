package slidingwindowlimiter

import (
	"sync"
	"testing"
	"time"
)

func TestSlidingWindowRateLimiter(t *testing.T) {
	limiter := NewSlidingWindowRateLimiter(3, time.Second, time.Second)

	for i := 0; i < 3; i++ {
		if !limiter.AllowSliding() {
			t.Errorf("Request %d should have been allowed", i+1)
		}
	}

	if limiter.AllowSliding() {
		t.Errorf("4th request should have been denied")
	}

	time.Sleep(500 * time.Millisecond)

	if !limiter.AllowSliding() {
		t.Errorf("Request after window reset should be allowed")
	}
}

func TestSlidingWindowRateLimiterWithLogs(t *testing.T) {
	limiter := NewSlidingWindowRateLimiter(3, time.Second, time.Second)
	allowed := 0
	denied := 0

	for i := 0; i < 10; i++ {
		if limiter.AllowSliding() {
			t.Logf("Request %d: ✅ allowed", i)
			allowed++
		} else {
			t.Logf("Request %d: ❌ denied", i)
			denied++
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Logf("Total allowed: %d, denied: %d", allowed, denied)
}

func TestSlidingWindowLimiterConcurrency(t *testing.T) {
	limiter := NewSlidingWindowRateLimiter(3, time.Second, time.Second)
	var wg sync.WaitGroup
	allowedCount := 0
	deniedCount := 0
	mu := sync.Mutex{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if limiter.AllowSliding() {
				mu.Lock()
				allowedCount++
				t.Logf("Goroutine %d: ✅ allowed", id)
				mu.Unlock()
			} else {
				mu.Lock()
				deniedCount++
				t.Logf("Goroutine %d: ❌ denied", id)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	t.Logf("Total allowed: %d, denied: %d", allowedCount, deniedCount)
	if allowedCount > 3 {
		t.Errorf("Allowed more than rate limit")
	}
}

func BenchmarkSlidingWindowLimiter(b *testing.B) {
	limiter := NewSlidingWindowRateLimiter(10, time.Second, time.Second)
	allowed := 0
	denied := 0

	for i := 0; i < b.N; i++ {
		if limiter.AllowSliding() {
			allowed++
		} else {
			denied++
		}
	}

	b.Logf("Benchmark done: Allowed = %d, Denied = %d", allowed, denied)
}

