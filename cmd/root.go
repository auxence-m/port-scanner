package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hostsFileFlag string
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pScan",
	Version: "1.0",
	Short:   "A fast TCP port scanner",
	Long: `pScan - short for Port Scanner - executes TCP port scan
on a list of hosts.

pScan allows you to add, list, and delete hosts from the list.

pScan executes a port scan on specified TCP or UDP ports. You can customize the
target ports and the scan output using a command line flag.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&hostsFileFlag, "hosts-file", "f", "pScan.hosts", "pScan hosts file")
	rootCmd.SetVersionTemplate(`{{printf "%s : %s - version %s\n" .Name .Short .Version}}`)

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("PSCAN")

	err := viper.BindPFlag("hosts-file", rootCmd.PersistentFlags().Lookup("hosts-file"))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pScan" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pScan")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Root() *cobra.Command {
	return rootCmd
}
