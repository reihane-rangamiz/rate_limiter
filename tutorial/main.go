package tutorial

import (
	"fmt"
	"rate_limiter/tutorial/examples"
	"time"
)


// examples.FixedWindowReal() output:
// Request 1: ✅ allowed  Request 2: ✅ allowed  Request 3: ✅ allowed  Request 4: ❌ denied Request 5: ❌ denied Waiting for next window...  Request after reset: ✅ allowed
func main() {
	fmt.Print("🪟Starting Fixed Window Basic Test...🪟")
	examples.FixedWindowBasicTest()

	fmt.Print("Resetting Window...⏰")
	time.Sleep(2 * time.Second)

	fmt.Print("😰Starting Fixed Window Stress Test...🪟😰")
	examples.FixedWindowStressTest()

	fmt.Print("Resetting Window...⏰")
	time.Sleep(2 * time.Second)

}
