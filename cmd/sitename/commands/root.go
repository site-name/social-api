package commands

import (
	"github.com/spf13/cobra"
)

type Command cobra.Command

var RootCmd = &cobra.Command{
	Use:   "sitename",
	Short: "Open source selling platform",
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "", "Configuration file to use.")
	RootCmd.PersistentFlags().Bool("disableconfigwatch", false, "When set config.json will not be loaded from disk when the file is changed.")
}

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}
