package cmd

import (
	"fmt"
	"io"
	"os"
	"pScan/scan"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <host1>...<hostN>",
	Aliases:      []string{"a"},
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	Short:        "Add new host(s) to list",
	RunE:         addRun,
	Example: `  
# To add a new host
pScan hosts add 192.168.0.199

# To add multiple hosts
pScan hosts add 192.168.0.199 192.168.0.56

# To add host(s) to a specific host file
pScan hosts add 192.168.0.199 192.168.0.56 --host-file file.hosts
`,
}

func addRun(cmd *cobra.Command, args []string) error {
	hostsFile := viper.GetString("hosts-file")

	return addAction(os.Stdout, hostsFile, args)
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
