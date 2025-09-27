package cmd

import (
	"fmt"
	"io"
	"os"
	"pScan/scan"

	"github.com/spf13/cobra"
)

var scannedPorts []int

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE:  scanRun,
}

func scanRun(cmd *cobra.Command, args []string) error {
	return scanAction(os.Stdout, hostFile, scannedPorts)
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
}
