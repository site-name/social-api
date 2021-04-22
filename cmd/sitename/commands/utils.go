package commands

import "github.com/spf13/cobra"

func getConfigDSN(command *cobra.Command, env map[string]string) string {
	configDSN, _ := command.Flags().GetString("config")

	// Config not supplied in flag, check env
	if configDSN == "" {
		configDSN = env["SN_CONFIG"]
	}

	// Config not supplied in env or flag use default
	if configDSN == "" {
		configDSN = "config.json"
	}

	return configDSN
}
