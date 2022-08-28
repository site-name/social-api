package commands

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"

	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store/sqlstore"
	"github.com/spf13/cobra"
)

var DbCmd = &cobra.Command{
	Use:   "db",
	Short: "Commands related to the database",
}

var PopulateDbCmd = &cobra.Command{
	Use:     "populate",
	Short:   "Popularize database with fake data",
	RunE:    populateDbCmdF,
	Example: ` sitename db populate --type [channel, all]`,
}

var InitDbCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database",
	Long: `Initialize the database for a given DSN, executing the migrations and loading the custom defaults if any.

This command should be run using a database configuration DSN.`,
	Example: `  # you can use the config flag to pass the DSN
  $ sitename db init --config modules/config/config.json 
  OR 
  $ ` + CustomDefaultsEnvVar + `=<path_to_config> sitename db init`,
	Args: cobra.NoArgs,
	RunE: initDbCmdF,
}

func init() {
	PopulateDbCmd.Flags().StringP("type", "t", "all", "specify which table to populate")
	PopulateDbCmd.Flags().IntP("amount", "a", 5, "specify which table to populate")

	InitDbCmd.Flags().StringP("config", "c", "modules/config/config.json", "path to config.json file.")

	DbCmd.AddCommand(
		InitDbCmd,
		PopulateDbCmd,
	)

	RootCmd.AddCommand(
		DbCmd,
	)
}

func populateAll(a *app.App, cmd *cobra.Command, args []string, amount int) error {
	err := populateChannel(a, cmd, args, amount)
	if err != nil {
		return errors.Wrap(err, "populate all failed to populate channel table")
	}

	return nil
}

func populateChannel(a *app.App, cmd *cobra.Command, args []string, amount int) error {
	if amount == 0 {
		return errors.New("amount must be greater than 0")
	}
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < amount; i++ {
		active := i%2 == 0
		channel := &channel.Channel{
			Name:     fmt.Sprintf("This is first channel #%d name", i+1),
			IsActive: active,
			Currency: Currencies[rand.Intn(len(Currencies))],
		}
		_, err := a.Srv().Store.Channel().Save(channel)
		if err != nil {
			return errors.Wrap(err, "could not create channel")
		}
	}
	return nil
}

func populateDbCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Srv().Shutdown()

	populateType, err := command.Flags().GetString("type")
	if err != nil {
		return err
	}
	populateType = strings.TrimSpace(strings.ToLower(populateType))

	amount, err := command.Flags().GetInt("amount")
	if err != nil {
		return err
	}

	switch populateType {
	case "all":
		return populateAll(a, command, args, amount)
	case "channel":
		return populateChannel(a, command, args, amount)
	}

	return nil
}

func initDbCmdF(command *cobra.Command, _ []string) error {
	dsn, err := command.Flags().GetString("config")
	if err != nil {
		slog.Error("Error getting config path. Trying get config path from environment var: "+CustomDefaultsEnvVar, slog.Err(err))
		dsn = os.Getenv(CustomDefaultsEnvVar)
		if dsn == "" {
			return errors.New("cannot get path to config.json file.")
		}
	}

	file, err := os.Open(dsn)
	if err != nil {
		return fmt.Errorf("unable to open custom defaults file at %q: %w", dsn, err)
	}
	defer file.Close()

	var config *model.Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return fmt.Errorf("unable to decode custom defaults configuration: %w", err)
	}

	sqlStore := sqlstore.New(config.SqlSettings, nil)
	defer sqlStore.Close()

	fmt.Println("Database store correctly initialized")

	return nil
}
