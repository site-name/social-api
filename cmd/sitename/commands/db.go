package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"

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

var PopulateCmd = &cobra.Command{
	Use:   "populate",
	Short: "Populate categories",
	RunE:  populateCategoriesCmdF,
}

func init() {
	InitDbCmd.Flags().StringP("config", "c", "modules/config/config.json", "path to config.json file.")
	PopulateCmd.Flags().StringP("type", "t", "category", "which type of table to populate")

	DbCmd.AddCommand(
		InitDbCmd,
		DBVersionCmd,
		MigrateCmd,
		PopulateCmd,
	)
	RootCmd.AddCommand(
		DbCmd,
	)
}

// func populateAll(a *app.App, cmd *cobra.Command, args []string, amount int) error {
// 	err := populateChannel(a, cmd, args, amount)
// 	if err != nil {
// 		return errors.Wrap(err, "populate all failed to populate channel table")
// 	}

// 	return nil
// }

// func populateChannel(a *app.App, cmd *cobra.Command, args []string, amount int) error {
// 	if amount == 0 {
// 		return errors.New("amount must be greater than 0")
// 	}
// 	rd := rand.New(rand.NewSource(time.Now().UnixNano()))

// 	for i := 0; i < amount; i++ {
// 		active := i%2 == 0
// 		channel := &model.Channel{
// 			Name:     fmt.Sprintf("This is first channel #%d name", i+1),
// 			IsActive: active,
// 			Currency: Currencies[rd.Intn(len(Currencies))],
// 		}
// 		_, err := a.Srv().Store.Channel().Save(channel)
// 		if err != nil {
// 			return errors.Wrap(err, "could not create channel")
// 		}
// 	}
// 	return nil
// }

type categoryPath struct {
	CategoryID     int    `json:"category_id"`
	CategoryName   string `json:"category_name"`
	CategoryNameEn string `json:"category_name_en"`
}

type rawCategory struct {
	CategoryID     int            `json:"category_id"`
	CategoryName   string         `json:"category_name"`
	CategoryNameEn string         `json:"category_name_en"`
	Toggle         bool           `json:"toggle"`
	Images         []string       `json:"images"`
	Path           []categoryPath `json:"path"`
}

func populateCategoriesCmdF(command *cobra.Command, args []string) error {
	cfgDSN := getConfigDSN(command, config.GetEnvironment())
	cfgStore, err := config.NewStoreFromDSN(cfgDSN, true, nil, true)
	if err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}
	config := cfgStore.Get()

	sqlStore := sqlstore.New(config.SqlSettings, nil)
	defer sqlStore.Close()

	// check if we already populated categories for the first time
	_, err = sqlStore.System().GetByName(model.PopulateCategoriesForTheFirstTimeKey)
	if err == nil {
		// means populated, return now
		slog.Info("categories already populated. Returning now.")
		return nil
	}

	currentDir, _ := os.Getwd()
	filePath := path.Join(currentDir, "model/raw_data/raw_categories.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var rawCategories []*rawCategory
	err = json.Unmarshal(data, &rawCategories)
	if err != nil {
		return err
	}

	categories := model.Categories{}
	meetMap := map[string]*model.Category{}

	for _, cate := range rawCategories {
		if cate.CategoryNameEn == "" {
			slog.Debug("a category is not translated", slog.Int("id", cate.CategoryID))
			continue
		}

		var (
			named     = "Category"
			slugg     = ""
			parentKey string
		)
		for pathIdx, path := range cate.Path {
			slugg += " " + path.CategoryNameEn
			if pathIdx > 0 {
				parentKey = named
			}
			named += path.CategoryNameEn

			if _, met := meetMap[named]; !met {
				desired := &model.Category{
					Slug:  slug.Make(slugg),
					Name:  path.CategoryNameEn,
					Level: uint8(pathIdx),
					NameTranslation: model.StringMAP{
						"vi": path.CategoryName,
					},
				}
				if pathIdx == len(cate.Path)-1 {
					desired.Images = strings.Join(cate.Images, " ")
				}
				if pathIdx > 0 {
					desired.ParentID = &meetMap[parentKey].Id
				}

				categories = append(categories, desired)
				meetMap[named] = desired
			}
		}
	}

	slog.Info("Populating categories for the first time...")

	_, err = sqlStore.GetMaster().Exec("DELETE FROM " + model.CategoryTableName)
	if err != nil {
		slog.Error("failed to delete categories from db")
		return err
	}

	for _, cate := range categories {
		_, err = sqlStore.Category().Upsert(cate)
		if err != nil {
			slog.Error("failed to insert category", slog.String("name", cate.Name))
			return err
		}
	}

	slog.Info("Successfully populated categories.")

	// indicate populated
	return sqlStore.System().Save(&model.System{
		Name:  model.PopulateCategoriesForTheFirstTimeKey,
		Value: "true",
	})
}

// func populateDbCmdF(command *cobra.Command, args []string) error {
// 	a, err := InitDBCommandContextCobra(command)
// 	if err != nil {
// 		return err
// 	}
// 	defer a.Srv().Shutdown()

// 	populateType, err := command.Flags().GetString("type")
// 	if err != nil {
// 		return err
// 	}
// 	populateType = strings.TrimSpace(strings.ToLower(populateType))

// 	amount, err := command.Flags().GetInt("amount")
// 	if err != nil {
// 		return err
// 	}

// 	switch populateType {
// 	case "all":
// 		return populateAll(a, command, args, amount)
// 	case "channel":
// 		return populateChannel(a, command, args, amount)
// 	case "category":
// 		return populateCategoriesCmdF(command, args)
// 	}

// 	return nil
// }

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
