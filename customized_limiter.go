package main

import (
	"fmt"
	"net"

	"golang.org/x/time/rate"
)

var RateLimit *CustomLimiter

type CustomLimiter struct {
	*rate.Limiter
	UserIp          string
	IPRange         []IPRange
	ExcludedUserIps []string
}

type IPRange struct {
	StartRange string
	EndRange   string
}

func NewCustomLimiter(rateLimit rate.Limit, burst int, userIp string, rangeIp []IPRange, excludedUserIps []string) *CustomLimiter {
	limiter := rate.NewLimiter(rateLimit, burst)
	return &CustomLimiter{
		Limiter:         limiter,
		UserIp:          userIp,
		IPRange:         rangeIp,
		ExcludedUserIps: excludedUserIps,
	}
}

func CreateNewRateLimiter(rateLimit rate.Limit, burst int, name string, userIp string, ipRagne []IPRange, excludedUserIps []string) {
	rateLimiter := NewCustomLimiter(rateLimit, burst, name, userIp, ipRagne, excludedUserIps)
	var rangeIPLimits = make(map[IPRange]*CustomLimiter, 0)

	// TODO
	// var userLimits = make(map[string]*CustomLimiter, 0)
	// var excludedUserIpLimit = make(map[string]*CustomLimiter, 0)

	isInIpRange := CheckIsInIPRange(ipRagne)

}

func CheckIsInIPRange(ip string) bool {
	userIP := net.ParseIP(ip)
	if userIP.To4() == nil {
		fmt.Printf("%v is not an IPv4 address\n", userIP)
		return false
	}

	// for key, value := range
	// if bytes.Compare(userIP, ip1) >= 0 && bytes.Compare(userIP, ip2) <= 0 {
	// 	fmt.Printf("%v is between %v and %v\n", userIP, ip1, ip2)
	// 	return true
	// }
	// fmt.Printf("%v is NOT between %v and %v\n", userIP, ip1, ip2)
	return false
}
