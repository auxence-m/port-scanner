package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrExists    = errors.New("host already in the list")
	ErrNotExists = errors.New("host not in the list")
)

// HostsList represents a list of hosts to run port scan
type HostsList struct {
	Hosts []string
}

// search searches for hosts in the list
func (hosts *HostsList) search(host string) (bool, int) {
	sort.Strings(hosts.Hosts)

	i := sort.SearchStrings(hosts.Hosts, host)
	if i < len(hosts.Hosts) && hosts.Hosts[i] == host {
		return true, i
	}

	return false, -1
}

func (hosts *HostsList) Add(host string) error {
	found, _ := hosts.search(host)
	if found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}

	hosts.Hosts = append(hosts.Hosts, host)
	return nil
}

// Remove deletes a host from the list
func (hosts *HostsList) Remove(host string) error {
	found, i := hosts.search(host)
	if found {
		hosts.Hosts = append(hosts.Hosts[:i], hosts.Hosts[i+1:]...)
		return nil
	}

	return fmt.Errorf("%w: %s", ErrNotExists, host)
}

// Load obtains hosts from a hosts file
func (hosts *HostsList) Load(hostsFile string) error {
	file, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hosts.Hosts = append(hosts.Hosts, scanner.Text())
	}

	return nil
}

// Save saves hosts to a hosts file
func (hosts *HostsList) Save(hostsFile string) error {
	output := ""

	for _, host := range hosts.Hosts {
		output += fmt.Sprintln(host)
	}

	return os.WriteFile(hostsFile, []byte(output), 0644)
}
