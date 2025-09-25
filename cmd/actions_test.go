package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"pScan/scan"
	"strings"
	"testing"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	// Create temp file
	tempFile, err := os.CreateTemp("", "pScan")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	// Initialize list if needed
	if initList {
		hostsList := &scan.HostsList{}

		for _, host := range hosts {
			if err := hostsList.Add(host); err != nil {
				t.Fatal(err)
			}
		}

		if err := hostsList.Save(tempFile.Name()); err != nil {
			t.Fatal(err)
		}
	}

	// Return temp file name and cleanup function
	return tempFile.Name(), func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Fatal(err)
		}
	}
}

func TestHostActions(t *testing.T) {
	// Define hosts for actions test
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	// Test cases for Action test
	testCases := []struct {
		name           string
		args           []string
		expectedOut    string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOut:    "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "ListAction",
			expectedOut:    "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "DeleteAction",
			args:           hosts,
			expectedOut:    "Deleted host: host1\nDeleted host: host2\nDeleted host: host3\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup Action test
			tempFile, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()

			// Define var to capture Action output
			var out bytes.Buffer

			// Execute Action and capture output
			if err := tc.actionFunction(&out, tempFile, tc.args); err != nil {
				t.Fatalf("Expected no error, got %q\n", err)
			}

			// Test Actions output
			if out.String() != tc.expectedOut {
				t.Errorf("Expected output %q, got %q\n", tc.expectedOut, out.String())
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	// Define hosts for actions test
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	// Setup integration test
	tempFile, cleanup := setup(t, hosts, false)
	defer cleanup()

	deleteHost := "host2"
	hostsEnd := []string{
		"host1",
		"host3",
	}

	// Define var to capture output
	var out bytes.Buffer

	// Define expected output for all actions
	expectedOut := ""

	// First, loop through the hosts slice to create the output for the add operation
	// Join the items of the hosts slice with a newline character \n as the output of the list operation
	// Add the output of the delete operation
	// Repeat the list output
	for _, host := range hosts {
		expectedOut += fmt.Sprintf("Added host: %s\n", host)
	}
	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintf("Deleted host: %s\n", deleteHost)
	expectedOut += strings.Join(hostsEnd, "\n")
	expectedOut += fmt.Sprintln()

	// Add hosts to the list
	if err := addAction(&out, tempFile, hosts); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// List hosts
	if err := listAction(&out, tempFile, nil); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// Delete host2
	if err := deleteAction(&out, tempFile, []string{deleteHost}); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// List hosts after delete
	if err := listAction(&out, tempFile, nil); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// Test integration output
	if out.String() != expectedOut {
		t.Errorf("Expected output %q, got %q\n", expectedOut, out.String())
	}
}
