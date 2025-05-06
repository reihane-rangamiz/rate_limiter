package slidingwindowlimiter

import (
	fixedwindowlimiter "rate_limiter/tutorial/fixed_window_limiter"
	"testing"
	"time"
)

// go test -v -run ^TestFixedVsSlidingWindowComparison$
func TestFixedVsSlidingWindowComparison(t *testing.T) {
	fixed := fixedwindowlimiter.NewFixedWindowRateLimiter(5, time.Second)
	sliding := NewSlidingWindowRateLimiter(5, time.Second) // real implementation below

	var fixedAllowed, slidingAllowed int

	// Send 5 requests quickly — both should allow them
	for i := 0; i < 5; i++ {
		if fixed.Allow() {
			fixedAllowed++
		}
		if sliding.AllowSliding() {
			slidingAllowed++
		}
	}

	// Wait just enough to cross the fixed window boundary
	time.Sleep(900 * time.Millisecond)

	// Send 5 more requests — fixed should reset, sliding should still consider old ones
	for i := 0; i < 5; i++ {
		if fixed.Allow() {
			fixedAllowed++
		}
		if sliding.AllowSliding() {
			slidingAllowed++
		}
	}

	t.Logf("✅ FixedWindow allowed: %d", fixedAllowed)     // Expect 10
	t.Logf("🌀 SlidingWindow allowed: %d", slidingAllowed) // Expect 5–7 depending on timestamps
}
