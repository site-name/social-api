package sqlstore

import (
	"github.com/sitename/sitename/model/app"
	"github.com/sitename/sitename/store"
)

type SqlAppTokenStore struct {
	*SqlStore
}

func newSqlAppTokenStore(sqlStore *SqlStore) store.AppTokenStore {
	as := &SqlAppTokenStore{
		SqlStore: sqlStore,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(app.AppToken{}, "AppTokens").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AppId").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(app.APP_TOKEN_NAME_MAX_LENGTH)
		table.ColMap("AuthToken").SetMaxSize(app.APP_TOKEN_AUTH_TOKEN_MAX_LENGTH).SetUnique(true)
	}

	return as
}

func (as *SqlAppTokenStore) createIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_app_tokens_name", "AppTokens", "Name")
	as.CreateIndexIfNotExists("idx_app_tokens_name_lower_textpattern", "AppTokens", "lower(Name) text_pattern_ops")
}

func (as *SqlAppTokenStore) Save(appToken *app.AppToken) (*app.AppToken, error) {
	panic("not implemented")
}
