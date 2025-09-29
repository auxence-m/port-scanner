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

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:          "scan",
	Short:        "Run a port scan on the hosts",
	SilenceUsage: true,
	RunE:         scanRun,
}

func scanRun(cmd *cobra.Command, args []string) error {
	hostsFile := viper.GetString("hosts-file")

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

	return scanAction(os.Stdout, hostsFile, scannedPorts)
}

func scanAction(out io.Writer, file string, ports []int) error {
	hosts := &scan.HostsList{}

	if err := hosts.Load(file); err != nil {
		return err
	}

	results := scan.Run(hosts, ports)

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
	scanCmd.Flags().IntSliceVarP(&scannedPorts, "ports", "p", []int{22, 80, 443}, "ports to scan")
	scanCmd.Flags().StringVarP(&portRange, "range", "r", "", "port range to scan")
}
