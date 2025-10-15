package cmd

import (
	"fmt"
	"io"
	"os"
	"pScan/scan"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var scannedPorts []int
var portRange string
var udp bool
var open bool

// Top 20 (most commonly opened) TCP ports
var tcpPorts = []int{80, 23, 443, 21, 22, 25, 3389, 110, 445, 139, 143, 53, 135, 3306, 8080, 1723, 111, 995, 993, 5900}

// Top 20 (most commonly opened) UDP ports
var udpPorts = []int{631, 161, 137, 123, 138, 1434, 445, 135, 67, 53, 139, 500, 68, 520, 1900, 4500, 514, 49152, 162, 69}

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:          "scan",
	Short:        "Run a port scan on the hosts",
	SilenceUsage: true,
	RunE:         scanRun,
	Example: `
# To scan one or multiple TCP ports
pScan scan --ports 80,135,445,139,50477,54672,59869

# To scan TCP ports within a specific range
pScan scan --range 59860-59890

# To scan one or multiple UDP ports
pScan scan --udp --ports 53,67,163,56448,57674

# To scan UDP ports within a specific range
pScan scan --udp --range 59860-59890

# To show only open port
pScan scan --ports 80,135,445,139,50477,54672,59869 --open

# To combine multiple scan options
pScan scan --ports 80,135,445,139,50477,54672,59869 --range 59860-59890 --open
`,
}

func scanRun(cmd *cobra.Command, args []string) error {
	hostsFile := viper.GetString("hosts-file")

	// When performing a UDP port scan, change default ports to well known UDP ports
	if udp {
		scannedPorts = []int{53, 67, 68, 123, 135}
	}

	// Validates the provided port range format
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

		if start > end {
			return fmt.Errorf("invalid port range [%d-%d]", start, end)
		}

		for i := start; i <= end; i++ {
			scannedPorts = append(scannedPorts, i)
		}
	}

	// Validates the provided port numbers are within the proper range for TCP ports from 1 to 65535
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
	protocol := "tcp"

	if err := hosts.Load(file); err != nil {
		return err
	}

	if udp {
		results = scan.Run(hosts, ports, "udp")
		protocol = "udp"
	} else {
		results = scan.Run(hosts, ports, "tcp")
	}

	return printResults(out, results, protocol)
}

func printResults(out io.Writer, results []scan.Results, protocol string) error {
	message := ""
	w := tabwriter.NewWriter(out, 0, 0, 5, ' ', 0)

	for _, res := range results {
		message += fmt.Sprintf("Scan report for %s:", res.Host)

		if res.NotFound {
			message += fmt.Sprintf(" Host Not Found\n\n")
			continue
		}

		message += fmt.Sprintln()
		message += fmt.Sprintln("PORT\tSTATE")

		for _, port := range res.PortStates {
			if open {
				if port.Open {
					message += fmt.Sprintf("%d/%s\t%s\n", port.Port, protocol, port.Open)
				}
			} else {
				message += fmt.Sprintf("%d/%s\t%s\n", port.Port, protocol, port.Open)
			}
		}

		message += fmt.Sprintln()
	}

	_, err := fmt.Fprintln(w, message)
	return err
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceVarP(&scannedPorts, "ports", "p", []int{21, 22, 25, 80, 443}, "ports to scan")
	scanCmd.Flags().StringVarP(&portRange, "range", "r", "", "port range to scan")
	scanCmd.Flags().BoolVar(&udp, "udp", false, "enable UDP port scans")
	scanCmd.Flags().BoolVar(&open, "open", false, "show only open ports")
}
