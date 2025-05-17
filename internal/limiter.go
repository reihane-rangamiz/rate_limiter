package internal

import (
	"strings"
	"sync"
	"time"
)

var clients map[string]*RateLimiterStore
var rateLimiter *RateLimiterConfig
var clientsMu sync.Mutex

// var ClientsLimit map[string]*RateLimiterConfig

type RateLimiterConfig struct {
	Rate             int
	Interval         time.Duration
	IncludedUserIPs  []string
	IncludedUserIDs  []string
	IncludedIPRanges []*IPRange
	ExcludedIPRange  []*IPRange
	ExcludedUserIPs  []string
	ExcludedUserIDs  []string
}

type IPRange struct {
	StartRange string
	EndRange   string
}

// stores user request track and the global set limiter config
type RateLimiterStore struct {
	mu       sync.Mutex
	configs  RateLimiterConfig
	requests []time.Time
}

func NewUserRateLimiter(ip string, config RateLimiterConfig) *RateLimiterStore {
	clients = make(map[string]*RateLimiterStore)
	return &RateLimiterStore{
		configs:  config,
		requests: make([]time.Time, 0),
	}
}

// return a new rate limiter
func NewRateLimiter(configs *RateLimiterConfig) RateLimiterConfig {
	var limiter = &RateLimiterConfig{
		Rate:             configs.Rate,
		Interval:         configs.Interval,
		IncludedUserIPs:  configs.IncludedUserIPs,
		IncludedUserIDs:  configs.IncludedUserIDs,
		IncludedIPRanges: configs.IncludedIPRanges,
		ExcludedIPRange:  configs.ExcludedIPRange,
		ExcludedUserIPs:  configs.ExcludedUserIPs,
		ExcludedUserIDs:  configs.ExcludedUserIDs,
	}
	rateLimiter = limiter
	return *rateLimiter
}

// userUniqID can be ip / id / or combanition
func (s RateLimiterConfig) Allow(userUniqID string) bool {
	if strings.TrimSpace(userUniqID) == "" {
		return false
	}

	for _, i := range s.ExcludedUserIPs {
		if i == userUniqID {
			return true
		}
	}

	for _, i := range s.ExcludedUserIDs {
		if i == userUniqID {
			return true
		}
	}

	for _, i := range s.IncludedUserIDs {
		if i != userUniqID {
			continue
		}
		return checkAndUpdateLimiter(userUniqID, s)
	}

	for _, i := range s.IncludedUserIPs {
		if i != userUniqID {
			continue
		}
		return checkAndUpdateLimiter(userUniqID, s)
	}

	return true
}

func checkAndUpdateLimiter(userUniqID string, config RateLimiterConfig) bool {
	clientsMu.Lock()
	limiter, exists := clients[userUniqID]
	if !exists {
		limiter = NewUserRateLimiter(userUniqID, config)
		clients[userUniqID] = limiter
	}
	clientsMu.Unlock()

	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	now := time.Now()
	limitTime := now.Add(-config.Interval)

	reqCount := 0
	for _, r := range clients[userUniqID].requests {
		if r.After(limitTime) {
			reqCount++
		}
	}

	if reqCount >= config.Rate {
		return false
	}
	clients[userUniqID].requests = append(clients[userUniqID].requests, now)
	return true
}

// func MustParseCIDR(cidr string) *net.IPNet {
// 	_, ipnet, err := net.ParseCIDR(cidr)
// 	if err != nil {
// 		log.Fatalf("invalid CIDR %q: %v", cidr, err)
// 	}
// 	return ipnet
// }

// multiple limiters
// func NewRateLimiterStore(configs []*RateLimiterConfig) *RateLimiterStore {
// 	return &RateLimiterStore{
// 		configs: configs,
// 	}
// }

// func RateLimiterMiddleware(r rate.Limit, b int) func() {
// Create a map to track rate limiters for different clients (IPs)
// 	var clients = make(map[string]*rate.Limiter)
// 	var mu sync.Mutex

// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			// Extract client IP (you can adjust this based on your needs)
// 			ip := c.RealIP()

// 			// Lock the map while accessing/modifying it
// 			mu.Lock()
// 			limiter, exists := clients[ip]
// 			if !exists {
// 				// If no limiter exists for the client, create a new one
// 				limiter = rate.NewLimiter(r, b)
// 				clients[ip] = limiter
// 			}
// 			mu.Unlock()

// 			// Check if the request is allowed based on the rate limit
// 			if !limiter.Allow() {
// 				return c.JSON(http.StatusTooManyRequests, map[string]string{
// 					"error": "too many requests",
// 				})
// 			}

// 			// Proceed to the next handler
// 			return next(c)
// 		}
// 	}
// }

// func NewCustomLimiter(config *RateLimiterConfig) *RateLimiterConfig {
// 	return &RateLimiterConfig{
// 		Limiter:         config.Limiter,
// 		UserIP:          config.UserIP,
// 		IPRanges:        config.IPRanges,
// 		ExcludedUserIPs: config.ExcludedUserIPs,
// 	}
// }
// func retrieveLimiters(ip string) []*rate.Limiter {
// 	var limiters []*rate.Limiter
// 	parsedIP := net.ParseIP(ip)

// 	for _, conf := range RateLimiterStore {
// 		for _, ex := range conf.ExcludedUserIps {
// 			if ip == ex {
// 				return nil
// 			}
// 		}

// 		if conf.UserIp != "" && conf.UserIp == ip {
// 			limiters = append(limiters, conf.Limiter)
// 			continue
// 		}

// 		for _, cidr := range conf.IPRanges {
// 			_, subnet, err := net.ParseCIDR(cidr)
// 			if err == nil && subnet.Contains(parsedIP) {
// 				limiters = append(limiters, conf.Limiter)
// 				break
// 			}
// 		}
// 	}

// 	return limiters
// }

// func RateLimitMiddleware(c *gin.Context) {
// 	ip := c.ClientIP()
// 	limiters := retrieveLimiters(ip)

// 	for _, limiter := range limiters {
// 		if !limiter.Allow() {
// 			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
// 				"error": "Rate limit exceeded",
// 			})
// 			return
// 		}
// 	}

// 	c.Next()
// }

// for _, config := range s.configs {
// 	// Excluded IPs
// 	if _, excluded := config.ExcludedUserIPs[ip]; excluded {
// 		return true
// 	}

// 	// Match exact IP
// 	if config.UserIP == ip {
// 		if !config.Limiter.Allow() {
// 			return false
// 		}
// 		continue
// 	}

// 	// Match IP ranges
// 	for _, cidr := range config.IPRanges {
// 		if cidr.Contains(parsedIP) {
// 			if !config.Limiter.Allow() {
// 				return false
// 			}
// 			break
// 		}
// 	}
// }
