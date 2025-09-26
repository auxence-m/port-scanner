package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var hostFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pScan",
	Version: "1.0",
	Short:   "A fast TCP port scanner",
	Long: `pScan - short for Port Scanner - executes TCP port scan
on a list of hosts.

pScan allows you to add, list, and delete hosts from the list.

pScan executes a port scan on specified TCP ports. You can customize the
target ports using a command line flag.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&hostFile, "hosts-file", "f", "pScan.hosts", "pScan hosts file")
	rootCmd.SetVersionTemplate(`{{printf "%s : %s - version %s\n" .Name .Short .Version}}`)
}
