package scan_test

import (
	"errors"
	"os"
	"pScan/scan"
	"testing"
)

func TestHostsList_Add(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hosts := scan.HostsList{}

			// Initialize list
			if err := hosts.Add("host1"); err != nil {
				t.Fatal(err)
			}

			err := hosts.Add(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("Expected error, got nil instead\n")
				}

				if !errors.Is(err, tc.expectErr) {
					t.Errorf("Expected error %q, got %q instead\n", tc.expectErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %q instead\n", err)
			}

			if len(hosts.Hosts) != tc.expectLen {
				t.Errorf("Expected list length %d, got %d instead\n", tc.expectLen, len(hosts.Hosts))
			}

			if hosts.Hosts[1] != tc.host {
				t.Errorf("Expected host name %q as index 1, got %q instead\n", tc.host, hosts.Hosts[1])
			}
		})
	}
}

func TestHostsList_Remove(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"RemoveExisting", "host1", 1, nil},
		{"RemoveNotFound", "host3", 1, scan.ErrNotExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hosts := scan.HostsList{}

			// Initialize list
			for _, host := range []string{"host1", "host2"} {
				if err := hosts.Add(host); err != nil {
					t.Fatal(err)
				}
			}

			err := hosts.Remove(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("Expected error, got nil instead\n")
				}

				if !errors.Is(err, tc.expectErr) {
					t.Errorf("Expected error %q, got %q instead\n", tc.expectErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %q instead\n", err)
			}

			if len(hosts.Hosts) != tc.expectLen {
				t.Errorf("Expected list length %d, got %d instead\n", tc.expectLen, len(hosts.Hosts))
			}

			if hosts.Hosts[0] == tc.host {
				t.Errorf("Host name %q should not be in the list\n", tc.host)
			}
		})
	}
}

func TestHostsList_SaveLoad(t *testing.T) {
	hosts1 := scan.HostsList{}
	hosts2 := scan.HostsList{}

	hostName := "host1"
	err := hosts1.Add(hostName)
	if err != nil {
		return
	}

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tempFile.Name())

	if err := hosts1.Save(tempFile.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}

	if err := hosts2.Load(tempFile.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}

	if hosts1.Hosts[0] != hosts2.Hosts[0] {
		t.Errorf("Host %q should match %q host.", hosts1.Hosts[0], hosts2.Hosts[0])
	}
}

func TestHostsList_LoadNoFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	tempFile.Close()

	if err := os.Remove(tempFile.Name()); err != nil {
		t.Fatalf("Error deleting temp file: %s", err)
	}

	hosts := scan.HostsList{}

	if err := hosts.Load(tempFile.Name()); err != nil {
		t.Errorf("Expected no error, got %q instead\n", err)
	}
}
