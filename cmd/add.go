package cmd

import (
	"fmt"
	"io"
	"os"
	"pScan/scan"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <host1>...<hostN>",
	Aliases:      []string{"a"},
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	Short:        "Add new host(s) to list",
	RunE:         addRun,
}

func addRun(cmd *cobra.Command, args []string) error {
	hostFile, err := cmd.Flags().GetString("hosts-file")
	if err != nil {
		return err
	}

	return addAction(os.Stdout, hostFile, args)
}

func addAction(out io.Writer, file string, args []string) error {
	hosts := &scan.HostsList{}

	if err := hosts.Load(file); err != nil {
		return err
	}

	for _, host := range args {
		if err := hosts.Add(host); err != nil {
			return err
		}

		if _, err := fmt.Fprintln(out, "Added host:", host); err != nil {
			return err
		}
	}

	return hosts.Save(file)
}

func init() {
	hostsCmd.AddCommand(addCmd)
}
