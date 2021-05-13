package sqlstore

import (
	"github.com/sitename/sitename/model/app"
	"github.com/sitename/sitename/store"
)

type AppSqlStore struct {
	*SqlStore
}

func newAppSqlStore(sqlStore *SqlStore) store.App {
	as := &AppSqlStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(app.App{}, "Apps").SetKeys(false, "id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(app.APP_NAME_MAX_LENGTH)
		table.ColMap("Identifier").SetMaxSize(app.APP_IDENTIFIER_MAX_LENGTH)
		table.ColMap("Version").SetMaxSize(app.APP_VERSION_MAX_LENGTH)
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
	}

	return as
}

func (as *AppSqlStore) createIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_apps_name", "Apps", "Name")
	as.CreateIndexIfNotExists("idx_apps_identifier", "Apps", "Identifier")
	as.CreateIndexIfNotExists("idx_apps_name_lower_textpattern", "Apps", "lower(Name) text_pattern_ops")
	as.CreateIndexIfNotExists("idx_apps_identifier_lower_textpattern", "Apps", "lower(Identifier) text_pattern_ops")
}

func (as *AppSqlStore) Save(app *app.App) (*app.App, error) {
	panic("not implemented") // NOTE: fixme
}
