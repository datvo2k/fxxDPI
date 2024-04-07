package dns

import (
	"errors"
	"net"
	"strconv"
)

func LookupIP(hostname string, portStr string) (*net.TCPAddr, error) {
	// Convert portStr to an integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	// Parse the hostname as an IP address
	ip := net.ParseIP(hostname)
	if ip == nil {
		return nil, errors.New("no suitable IP address found")
	}

	// Check if the IP address is private
	if isPrivateIP(ip) {
		return nil, errors.New("private IPs are not allowed")
	}

	// Create a new TCPAddr with the IP address and port number
	addr := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}

	return addr, nil
}
