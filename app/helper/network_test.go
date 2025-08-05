package helper

import (
	"net"
	"testing"
)

func TestIsPrivateIP_PrivateRanges(t *testing.T) {
	privateIPs := []string{
		"10.0.0.1",
		"172.16.0.1",
		"192.168.1.1",
		"127.0.0.1",
		"169.254.1.1",
	}

	for _, ipStr := range privateIPs {
		ip := net.ParseIP(ipStr)
		if !IsPrivateIP(ip) {
			t.Errorf("expected %s to be private", ipStr)
		}
	}
}

func TestIsPrivateIP_PublicRanges(t *testing.T) {
	publicIPs := []string{
		"8.8.8.8",
		"1.1.1.1",
		"208.67.222.222",
	}

	for _, ipStr := range publicIPs {
		ip := net.ParseIP(ipStr)
		if IsPrivateIP(ip) {
			t.Errorf("expected %s to be public", ipStr)
		}
	}
}
