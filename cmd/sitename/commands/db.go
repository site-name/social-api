package commands

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/config"
	"github.com/spf13/cobra"
)

var DbCmd = &cobra.Command{
	Use:   "db",
	Short: "Commands related to the database",
}

var InitDbCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database",
	Long: `Initialize the database for a given DSN, executing the migrations and loading the custom defaults if any.

This command should be run using a database configuration DSN.`,
	Example: `  # you can use the config flag to pass the DSN
  $ mattermost db init --config postgres://localhost/mattermost

  # or you can use the SN_CONFIG environment variable
  $ SN_CONFIG=postgres://localhost/mattermost mattermost db init

  # and you can set a custom defaults file to be loaded into the database
  $ SN_CUSTOM_DEFAULTS_PATH=custom.json SN_CONFIG=postgres://localhost/mattermost mattermost db init`,
	Args: cobra.NoArgs,
	RunE: initDbCmdF,
}

func init() {
	DbCmd.AddCommand(
		InitDbCmd
	)

	RootCmd.AddCommand(
		DbCmd,
	)
}

func initDbCmdF(command *cobra.Command, _ []string) error {
	dsn := getConfigDSN(command, config.GetEnvironment())

	if !config.IsDatabaseDSN(dsn) {
		return errors.New("this command should be run using a database configuration DSN")
	}

	customDefaults, err := loadCustomDefaults()
	if err != nil {
		return errors.Wrap(err, "error loading custom configuration defaults")
	}
}
