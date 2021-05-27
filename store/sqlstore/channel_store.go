package sqlstore

import (
	"github.com/pkg/errors"
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
		table.ColMap("Slug").SetMaxSize(channel.CHANNEL_SLUG_MAX_LENGTH).SetUnique(true)
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

func (cs *SqlChannelStore) Save(ch *channel.Channel) (*channel.Channel, error) {
	ch.PreSave()
	if err := ch.IsValid(); err != nil {
		return nil, err
	}

	if err := cs.GetMaster().Insert(ch); err != nil {
		if IsUniqueConstraintError(err, []string{"Slug", "channels_slug_key", "idx_channels_slug_unique"}) {
			return nil, store.NewErrInvalidInput("Channel", "Slug", ch.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save channel with id=%s", ch.Id)
	}

	return ch, nil
}

func (cs *SqlChannelStore) GetChannelsByIdsAndOrder(ids []string, order string) ([]*channel.Channel, error) {
	var channels []*channel.Channel
	_, err := cs.GetReplica().Select(
		&channels,
		`SELECT 
			*
		FROM
			Channels
		WHERE
			Id IN :IDS
		ORDER BY :Order`,
		map[string]interface{}{"IDS": ids, "Order": order},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find channels by ids")
	}

	return channels, nil
}
