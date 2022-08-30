package commands

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"

	"github.com/sitename/sitename/modules/config"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store/sqlstore"
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
  $ sitename db init --config modules/config/config.json 
  OR 
  $ ` + CustomDefaultsEnvVar + `=<path_to_config> sitename db init`,
	Args: cobra.NoArgs,
	RunE: initDbCmdF,
}

var DBVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the recent applied version number",
	RunE:  dbVersionCmdF,
}

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database if there are any unapplied migrations",
	Long:  "Run the missing migrations from the migrations table.",
	RunE:  migrateCmdF,
}

func init() {
	InitDbCmd.Flags().StringP("config", "c", "modules/config/config.json", "path to config.json file.")

	DbCmd.AddCommand(
		InitDbCmd,
		DBVersionCmd,
		MigrateCmd,
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

func dbVersionCmdF(command *cobra.Command, args []string) error {
	cfgDSN := getConfigDSN(command, config.GetEnvironment())
	cfgStore, err := config.NewStoreFromDSN(cfgDSN, true, nil, true)
	if err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}
	config := cfgStore.Get()

	store := sqlstore.New(config.SqlSettings, nil)
	defer store.Close()

	allFlags, _ := command.Flags().GetBool("all")
	if allFlags {
		applied, err2 := store.GetAppliedMigrations()
		if err2 != nil {
			return errors.Wrap(err2, "failed to get applied migrations")
		}
		for _, migration := range applied {
			CommandPrettyPrintln(fmt.Sprintf("Version: %d, Name: %s", migration.Version, migration.Name))
		}
		return nil
	}

	v, err := store.GetDBSchemaVersion()
	if err != nil {
		return errors.Wrap(err, "failed to get schema version")
	}

	CommandPrettyPrintln("Current database schema version is: " + strconv.Itoa(v))

	return nil
}

func migrateCmdF(command *cobra.Command, args []string) error {
	cfgDSN := getConfigDSN(command, config.GetEnvironment())
	cfgStore, err := config.NewStoreFromDSN(cfgDSN, true, nil, true)
	if err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}
	config := cfgStore.Get()

	store := sqlstore.New(config.SqlSettings, nil)
	defer store.Close()

	CommandPrettyPrintln("Database successfully migrated")

	return nil
}
