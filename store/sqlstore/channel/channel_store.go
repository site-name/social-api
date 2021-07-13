package channel

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/store"
)

type SqlChannelStore struct {
	store.Store
}

func NewSqlChannelStore(sqlStore store.Store) store.ChannelStore {
	cs := &SqlChannelStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(channel.Channel{}, store.ChannelTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(channel.CHANNEL_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(channel.CHANNEL_SLUG_MAX_LENGTH).SetUnique(true)
	}

	return cs
}

func (cs *SqlChannelStore) CreateIndexesIfNotExists() {
	cs.CreateIndexIfNotExists("idx_channels_name", store.ChannelTableName, "Name")
	cs.CreateIndexIfNotExists("idx_channels_slug", store.ChannelTableName, "Slug")
	cs.CreateIndexIfNotExists("idx_channels_isactive", store.ChannelTableName, "IsActive")
	cs.CreateIndexIfNotExists("idx_channels_currency", store.ChannelTableName, "Currency")

	cs.CreateIndexIfNotExists("idx_channels_name_lower_textpattern", store.ChannelTableName, "lower(Name) text_pattern_ops")
}

func (cs *SqlChannelStore) Save(ch *channel.Channel) (*channel.Channel, error) {
	ch.PreSave()
	if err := ch.IsValid(); err != nil {
		return nil, err
	}

	if err := cs.GetMaster().Insert(ch); err != nil {
		if cs.IsUniqueConstraintError(err, []string{"Slug", "channels_slug_key", "idx_channels_slug_unique"}) {
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
		`SELECT * FROM `+store.ChannelTableName+` WHERE Id IN :IDS ORDER BY :Order`,
		map[string]interface{}{"IDS": ids, "Order": order},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find channels by ids")
	}

	return channels, nil
}

func (cs *SqlChannelStore) Get(id string) (*channel.Channel, error) {
	var channel *channel.Channel
	err := cs.GetReplica().SelectOne(&channel, "SELECT * FROM "+store.ChannelTableName+" WHERE Id = :id", map[string]interface{}{"id": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ChannelTableName, id)
		}
		return nil, errors.Wrapf(err, "Failed to get Channel with ChannelID=%s", id)
	}

	return channel, nil
}

func (cs *SqlChannelStore) GetBySlug(slug string) (*channel.Channel, error) {
	var channel *channel.Channel
	err := cs.GetReplica().SelectOne(&channel, "SELECT * FROM "+store.ChannelTableName+" WHERE Slug = :slug", map[string]interface{}{"slug": slug})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ChannelTableName, "slug="+slug)
		}
		return nil, errors.Wrapf(err, "Failed to get Channel with slug=%s", slug)
	}

	return channel, nil
}

func (cs *SqlChannelStore) GetRandomActiveChannel() (*channel.Channel, error) {
	var channels = []*channel.Channel{}
	_, err := cs.GetReplica().Select(&channels, "SELECT * FROM "+store.ChannelTableName+" WHERE IsActive")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ChannelTableName, "")
		}
		return nil, errors.Wrap(err, "failed to get Channel with Active=true")
	}

	return channels[0], nil
}
