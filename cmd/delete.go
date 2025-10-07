package cmd

import (
	"fmt"
	"io"
	"os"
	"pScan/scan"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:          "delete <host1>...<hostN>",
	Aliases:      []string{"d"},
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	Short:        "Delete hosts(s) from list",
	RunE:         deleteRun,
	Example: `  
# To delete a new host
pScan hosts delete 192.168.0.199

# To delete multiple hosts
pScan hosts delete 192.168.0.199 192.168.0.56

# To delete host(s) from a specific host file
pScan hosts delete 192.168.0.199 192.168.0.56 --host-file file.hosts
`,
}

func deleteRun(cmd *cobra.Command, args []string) error {
	hostsFile := viper.GetString("hosts-file")

	return deleteAction(os.Stdout, hostsFile, args)
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
