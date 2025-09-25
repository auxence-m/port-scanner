package cmd

import (
	"fmt"
	"io"
	"os"
	"pScan/scan"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List hosts in hosts list",
	RunE:    listRun,
}

func listRun(cmd *cobra.Command, args []string) error {
	hostFile, err := cmd.Flags().GetString("hosts-file")
	if err != nil {
		return err
	}

	return listAction(os.Stdout, hostFile, args)
}

func listAction(out io.Writer, file string, args []string) error {
	hosts := &scan.HostsList{}

	if err := hosts.Load(file); err != nil {
		return err
	}

	for _, host := range hosts.Hosts {
		if _, err := fmt.Fprintln(out, host); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	hostsCmd.AddCommand(listCmd)
}
