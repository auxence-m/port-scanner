package scan_test

import (
	"net"
	"pScan/scan"
	"strconv"
	"testing"
)

func TestState_String(t *testing.T) {
	portState := scan.PortState{}

	if portState.Open.String() != "closed" {
		t.Errorf("Expected %q, got %q instead\n", "closed", portState.Open.String())
	}

	portState.Open = true

	if portState.Open.String() != "open" {
		t.Errorf("Expected %q, got %q instead\n", "Open", portState.Open.String())
	}
}

func TestRunHost_Found(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}

	host := "localhost"
	hosts := &scan.HostsList{}
	if err := hosts.Add(host); err != nil {
		t.Fatal(err)
	}

	var ports []int

	// Init ports, 1 open, 1 closed
	for _, tc := range testCases {
		// Find an open port on localhost
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)
		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}

	result := scan.Run(hosts, ports, "tcp")

	// There should be only one element in the Results slice returned
	// Verify results for HostFound test
	if len(result) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(result))
	}

	if result[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, result[0].Host)
	}

	if result[0].NotFound {
		t.Errorf("Expected host %q to be found\n", host)
	}

	if len(result[0].PortStates) != 2 {
		t.Fatalf("Expected 2 port states, got %d instead\n", len(result[0].PortStates))
	}

	for i, tc := range testCases {
		if result[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expected port %d, got %d instead\n", ports[i], result[0].PortStates[i].Port)
		}

		if result[0].PortStates[i].Open.String() != tc.expectState {
			t.Errorf("Expected port %d to be %s\n", ports[i], tc.expectState)
		}
	}
}

func TestRunHost_NotFound(t *testing.T) {
	host := "389.389.389.389"
	hosts := &scan.HostsList{}
	if err := hosts.Add(host); err != nil {
		t.Fatal(err)
	}

	result := scan.Run(hosts, []int{}, "tcp")

	// Verify results for HostNotFound test
	if len(result) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(result))
	}

	if result[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, result[0].Host)
	}

	if !result[0].NotFound {
		t.Errorf("Expected host %q NOT to be found\n", host)
	}

	if len(result[0].PortStates) != 0 {
		t.Fatalf("Expected 0 port states, got %d instead\n", len(result[0].PortStates))
	}
}
