package main

import (
	"log"
	"net"
	"sync"

	"golang.org/x/time/rate"
)

// var RateLimiterStore []*RateLimiterConfig

type RateLimiterConfig struct {
	Limiter         *rate.Limiter
	IncludedUserIPs []string
	IncludedUserIDs []string
	IPRanges        []*net.IPNet
	ExcludedUserIPs map[string]struct{}
	ExcludedUserIDs map[string]struct{}
}

type RateLimiterStore struct {
	mu sync.RWMutex
	// clients map[string][]*rate.Limiter
	configs []*RateLimiterConfig
}

func NewRateLimiterStore(configs []*RateLimiterConfig) *RateLimiterStore {
	return &RateLimiterStore{
		configs: configs,
	}
}

func (s *RateLimiterStore) Allow(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, config := range s.configs {
		// Excluded IPs
		if _, excluded := config.ExcludedUserIPs[ip]; excluded {
			return true
		}

		// Match exact IP
		if config.UserIP == ip {
			if !config.Limiter.Allow() {
				return false
			}
			continue
		}

		// Match IP ranges
		for _, cidr := range config.IPRanges {
			if cidr.Contains(parsedIP) {
				if !config.Limiter.Allow() {
					return false
				}
				break
			}
		}
	}

	return true
}

func MustParseCIDR(cidr string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatalf("invalid CIDR %q: %v", cidr, err)
	}
	return ipnet
}

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
