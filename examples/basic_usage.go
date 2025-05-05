package examples

import (
	"fmt"
	ratelimiter "rate_limiter/fixed_window_limiter"
	"sync"
	"time"
)

func FixedWindowBasicTest() {
	limiter := ratelimiter.NewFixedWindowRateLimiter(3, time.Second)

	for i := 1; i <= 5; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: ✅ allowed\n", i)
		} else {
			fmt.Printf("Request %d: ❌ denied\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Waiting for next window...")
	time.Sleep(1 * time.Second)

	if limiter.Allow() {
		fmt.Println("Request after reset: ✅ allowed")
	} else {
		fmt.Println("Request after reset: ❌ denied")
	}
}

func FixedWindowStressTest() {
	limiter := ratelimiter.NewFixedWindowRateLimiter(5, time.Second)

	var wg sync.WaitGroup
	totalRequests := 20

	var allowedCount, deniedCount int64
	var mu sync.Mutex

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			if limiter.Allow() {
				mu.Lock()
				allowedCount++
				mu.Unlock()
				fmt.Printf("Goroutine %d: ✅ allowed\n", id)
			} else {
				mu.Lock()
				deniedCount++
				mu.Unlock()
				fmt.Printf("Goroutine %d: ❌ denied\n", id)
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("\n--- Summary ---")
	fmt.Printf("Allowed: %d\n", allowedCount)
	fmt.Printf("Denied:  %d\n", deniedCount)
}
