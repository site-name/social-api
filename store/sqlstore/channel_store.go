package sqlstore

import (
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/store"
)

type SqlChannelStore struct {
	*SqlStore
}

func newSqlChannelStore(sqlStore *SqlStore) store.ChannelStore {
	cs := &SqlChannelStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(channel.Channel{}, "Channels").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(channel.CHANNEL_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(channel.CHANNEL_SLUG_MAX_LENGTH)
	}

	return cs
}

func (cs *SqlChannelStore) createIndexesIfNotExists() {
	cs.CreateIndexIfNotExists("idx_channels_name", "Channels", "Name")
	cs.CreateIndexIfNotExists("idx_channels_slug", "Channels", "Slug")
	cs.CreateIndexIfNotExists("idx_channels_isactive", "Channels", "IsActive")
	cs.CreateIndexIfNotExists("idx_channels_currency", "Channels", "Currency")

	cs.CreateIndexIfNotExists("idx_channels_name_lower_textpattern", "Channels", "lower(Name) text_pattern_ops")
}
