package scan

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single TCP port
type PortState struct {
	Port int
	Open state
}

type state bool

// String converts the boolean value of state to a human-readable string
func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

// scanPortTCP performs a port scan on a single TCP port
func scanPortTCP(host string, port int) PortState {
	portState := PortState{
		Port: port,
	}

	// define the network address based on the host and port
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	// connexion attempt to the network address
	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return portState
	}

	err = scanConn.Close()
	if err != nil {
		return PortState{}
	}

	portState.Open = true
	return portState
}

// scanPortUDP performs a port scan on a single UDP port
func scanPortUDP(host string, port int) PortState {
	portState := PortState{
		Port: port,
	}

	// define the network address based on the host and port
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	// define the network address based on the host and port
	scanConn, err := net.DialTimeout("udp", address, 1*time.Second)
	if err != nil {
		return portState
	}

	err = scanConn.Close()
	if err != nil {
		return portState
	}

	portState.Open = true
	return portState
}

// Results represents the scan results for a single host
type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

// Run performs a port scan on the hosts list
func Run(hosts *HostsList, ports []int, protocol string) []Results {
	results := make([]Results, 0, len(hosts.Hosts))

	// for every host
	for _, host := range hosts.Hosts {
		res := Results{
			Host: host,
		}

		// resolves hostname into a valid IP address
		if _, err := net.LookupHost(host); err != nil {
			res.NotFound = true
			results = append(results, res)
			continue // process next item in the loop
		}

		// scan ports for the given IP address
		if protocol == "tcp" {
			for _, port := range ports {
				res.PortStates = append(res.PortStates, scanPortTCP(host, port))
			}
		}

		if protocol == "udp" {
			for _, port := range ports {
				res.PortStates = append(res.PortStates, scanPortUDP(host, port))
			}
		}

		results = append(results, res)
	}

	return results
}
