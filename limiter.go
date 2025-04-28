package ratelimiter


import (
    "time"
)

type RateLimiter struct {
    rate       int           // tokens per interval
    interval   time.Duration // refill interval
    tokens     chan struct{} // token bucket
}
