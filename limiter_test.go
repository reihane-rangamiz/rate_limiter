package ratelimiter

import (
	"sync"
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

// in this test we check the process with logs + small delay between request (check debug console to see the behavior)
// use this inside your terminal to see the logs because go hides the logs inside the test files: "go test -v -run ^TestFixedWindowRateLimiterWithLogging$"
// result should look like:
// === RUN   TestFixedWindowRateLimiterWithLogging
//
//	limiter_test.go:39: Request 1: allowed = true
//	limiter_test.go:39: Request 2: allowed = true
//	limiter_test.go:39: Request 3: allowed = true
//	limiter_test.go:39: Request 4: allowed = false
//	limiter_test.go:39: Request 5: allowed = false
//
// --- PASS: TestFixedWindowRateLimiterWithLogging (0.52s)
// PASS
// ok      rate_limiter    1.320s
func TestFixedWindowRateLimiterWithLogging(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(3, time.Minute) // changed this from second to minute so you can debug line by line as well, otherwise return it to second

	for i := 1; i <= 5; i++ {
		allowed := limiter.Allow()
		t.Logf("Request %d: allowed = %v", i, allowed)
		time.Sleep(100 * time.Millisecond) // space out slightly
	}
}

// in this test we check the concurrency pattern behaivior on our limiter
// this func simulates real-world parallel traffic.
// The limiter must: Block after 3 requests - Never panic -Never allow more than 3
// go test -v -run ^TestFixedWindowRateLimiterWithConcurrency$
func TestFixedWindowRateLimiterWithConcurrency(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(3, time.Second)
	var wg sync.WaitGroup
	successCount := 0
	mu := sync.Mutex{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if limiter.Allow() {
				mu.Lock()
				successCount++
				mu.Unlock()
				t.Logf("Goroutine %d: allowed ✅", id)
			} else {
				t.Logf("Goroutine %d: denied ❌", id)
			}
		}(i)
	}

	wg.Wait()

	if successCount > 3 {
		t.Errorf("Too many allowed requests: %d", successCount)
	}
}

// this func is for stress testing / benchmarking :
// Runs .Allow() millions of times.
// b.N is set automatically by Go’s benchmarking engine.
// Used to check performance (how fast is the limiter?).
// Shows how the limiter performs when called at scale, like in real traffic.
// This limiter is shared across all iterations of the benchmark (not reset every time).
// command : go test -bench=.
func BenchmarkLimiterAllow(b *testing.B) {
	limiter := NewFixedWindowRateLimiter(3, time.Second)

	for i := 0; i < b.N; i++ {
		limiter.Allow()

	}
}

// benchmark with log
func BenchmarkLimiterAllowlogs(b *testing.B) {
	limiter := NewFixedWindowRateLimiter(3, time.Second)

	allowedCount := 0
	deniedCount := 0

	for i := 0; i < b.N; i++ {
		if limiter.Allow() {
			allowedCount++
		} else {
			deniedCount++
		}
	}

	b.Logf("Allowed: %d, Denied: %d", allowedCount, deniedCount)
}
