package cmd

import (
	"fmt"
	"io"
	"os"
	"pScan/scan"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:          "delete <host1>...<hostN>",
	Aliases:      []string{"d"},
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	Short:        "Delete hosts(s) from list",
	RunE:         deleteRun,
}

func deleteRun(cmd *cobra.Command, args []string) error {
	hostFile, err := cmd.Flags().GetString("hosts-file")
	if err != nil {
		return err
	}

	return deleteAction(os.Stdout, hostFile, args)
}

func deleteAction(out io.Writer, file string, args []string) error {
	hosts := &scan.HostsList{}

	if err := hosts.Load(file); err != nil {
		return err
	}

	for _, host := range args {
		if err := hosts.Remove(host); err != nil {
			return err
		}

		if _, err := fmt.Fprintln(out, "Deleted host:", host); err != nil {
			return err
		}
	}

	return hosts.Save(file)
}

func init() {
	hostsCmd.AddCommand(deleteCmd)
}
