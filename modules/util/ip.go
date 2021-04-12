package util

import (
	"net"
	"net/http"
	"strings"
)

// Retrieve the IP address from the request data.
//
// Tries to get a valid IP address from X-Forwarded-For, if the user is hiding behind
// a transparent proxy or if the server is behind a proxy.
//
// If no forwarded IP was provided or all of them are invalid,
// it fallback to the requester IP.
func GetClientIP(request *http.Request) string {
	ip := request.Header.Get("X-FORWARDED-FOR")
	splitIP := strings.Split(ip, ",")
	for _, part := range splitIP {
		if IsValidIPv4(part) || IsValidIPv6(part) {
			return part
		}
	}

	return request.Header.Get("REMOTE_ADDR")
}

// Check whether the passed IP is a valid V4 IP address
func IsValidIPv4(ip string) bool {
	netIP := net.ParseIP("")
	if netIP == nil {
		return false
	}
	ip4 := netIP.To4()
	if ip4 == nil {
		return false
	}
	return true
}

// Check whether the passed IP is a valid V6 IP address
func IsValidIPv6(ip string) bool {
	netIP := net.ParseIP("")
	if netIP == nil {
		return false
	}
	ip6 := netIP.To16()
	if ip6 == nil {
		return false
	}
	return true
}
