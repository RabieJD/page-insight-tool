package helper

import "net"

// IsPrivateIP checks if an IP address is private
func IsPrivateIP(ip net.IP) bool {
	// Handle IPv6 addresses
	if ip.To4() == nil {
		// For IPv6, check for private ranges
		// ::1 is localhost
		if ip.Equal(net.ParseIP("::1")) {
			return true
		}

		// Ensure we have enough bytes before accessing
		if len(ip) < 2 {
			return false
		}

		// fc00::/7 is unique local address
		if ip[0] == 0xfc || ip[0] == 0xfd {
			return true
		}
		// fe80::/10 is link-local
		if ip[0] == 0xfe && (ip[1]&0xc0) == 0x80 {
			return true
		}
		return false
	}

	// Check for private IPv4 ranges
	privateRanges := []struct {
		start, end net.IP
	}{
		{net.ParseIP("10.0.0.0"), net.ParseIP("10.255.255.255")},
		{net.ParseIP("172.16.0.0"), net.ParseIP("172.31.255.255")},
		{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")},
		{net.ParseIP("127.0.0.0"), net.ParseIP("127.255.255.255")},
		{net.ParseIP("169.254.0.0"), net.ParseIP("169.254.255.255")},
		{net.ParseIP("0.0.0.0"), net.ParseIP("0.255.255.255")},
	}

	for _, r := range privateRanges {
		if InRange(ip, r.start, r.end) {
			return true
		}
	}

	return false
}

// InRange checks if an IP is within a range
func InRange(ip, start, end net.IP) bool {
	// Convert all IPs to IPv4 for comparison
	ip4 := ip.To4()
	start4 := start.To4()
	end4 := end.To4()

	// If any of the IPs are not IPv4, we can't compare them
	if ip4 == nil || start4 == nil || end4 == nil {
		return false
	}

	return Bytes2Int(ip4) >= Bytes2Int(start4) && Bytes2Int(ip4) <= Bytes2Int(end4)
}

// Bytes2Int converts IP bytes to integer for comparison
func Bytes2Int(ip net.IP) uint32 {
	// Convert to IPv4 if it's an IPv6 address
	ip = ip.To4()
	if ip == nil {
		// If it's not an IPv4 address, return 0 to avoid panic
		return 0
	}

	// Ensure we have exactly 4 bytes before accessing them
	if len(ip) != 4 {
		return 0
	}

	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}
