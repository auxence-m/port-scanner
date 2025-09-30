package cmd

import (
	"fmt"
	"io"
	"os"
	"pScan/scan"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var scannedPorts []int
var portRange string
var udp bool

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:          "scan",
	Short:        "Run a port scan on the hosts",
	SilenceUsage: true,
	RunE:         scanRun,
}

func scanRun(cmd *cobra.Command, args []string) error {
	hostsFile := viper.GetString("hosts-file")

	// When performing a UDP port scan, change default ports to well known UDP ports
	if udp {
		scannedPorts = []int{53, 67, 68, 123, 135, 161}
	}

	// Verifying provided port range format
	if portRange != "" {
		rangeStr := strings.Split(portRange, "-")
		if len(rangeStr) != 2 {
			return fmt.Errorf("parsing \"%s\": invalid port range format", portRange)
		}

		start, err := strconv.Atoi(rangeStr[0])
		if err != nil {
			return err
		}
		end, err := strconv.Atoi(rangeStr[1])
		if err != nil {
			return err
		}

		for i := start; i <= end; i++ {
			scannedPorts = append(scannedPorts, i)
		}
	}

	for _, port := range scannedPorts {
		if port < 1 || port > 65535 {
			return fmt.Errorf("port %d is out of range [1-65535]", port)
		}
	}

	return scanAction(os.Stdout, hostsFile, scannedPorts)
}

func scanAction(out io.Writer, file string, ports []int) error {
	hosts := &scan.HostsList{}
	results := make([]scan.Results, 0, len(hosts.Hosts))

	if err := hosts.Load(file); err != nil {
		return err
	}

	if udp {
		results = scan.Run(hosts, ports, "udp")
	} else {
		results = scan.Run(hosts, ports, "tcp")
	}

	return printResults(out, results)
}

func printResults(out io.Writer, results []scan.Results) error {
	message := ""

	for _, res := range results {
		message += fmt.Sprintf("%s:", res.Host)

		if res.NotFound {
			message += fmt.Sprintf(" Host not found\n\n")
			continue
		}

		message += fmt.Sprintln()

		for _, port := range res.PortStates {
			message += fmt.Sprintf("\t%d: %s\n", port.Port, port.Open)
		}

		message += fmt.Sprintln()
	}

	_, err := fmt.Fprint(out, message)
	return err
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceVarP(&scannedPorts, "ports", "p", []int{21, 22, 25, 80, 443}, "ports to scan")
	scanCmd.Flags().StringVarP(&portRange, "range", "r", "", "port range to scan")
	scanCmd.Flags().BoolVar(&udp, "udp", false, "enable UDP port scans")
}
