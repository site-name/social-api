package channel

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
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

func (cs *SqlChannelStore) ModelFields() []string {
	return []string{
		"Channels.Id",
		"Channels.Name",
		"Channels.IsActive",
		"Channels.Slug",
		"Channels.Currency",
	}
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

func (cs *SqlChannelStore) Get(id string) (*channel.Channel, error) {
	var channel channel.Channel
	err := cs.GetReplica().SelectOne(&channel, "SELECT * FROM "+store.ChannelTableName+" WHERE Id = :id", map[string]interface{}{"id": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ChannelTableName, id)
		}
		return nil, errors.Wrapf(err, "Failed to get Channel with ChannelID=%s", id)
	}

	return &channel, nil
}

func (cs *SqlChannelStore) GetRandomActiveChannel() (*channel.Channel, error) {
	var channels = []*channel.Channel{}
	_, err := cs.GetReplica().Select(&channels, "SELECT * FROM "+store.ChannelTableName+" WHERE IsActive")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Channel with Active=true")
	}

	return channels[0], nil
}

func (cs *SqlChannelStore) commonQueryBuilder(option *channel.ChannelFilterOption) (string, []interface{}, error) {
	query := cs.GetQueryBuilder().
		Select(cs.ModelFields()...).
		From(store.ChannelTableName).
		OrderBy(store.TableOrderingMap[store.ChannelTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.Name != nil {
		query = query.Where(option.Name.ToSquirrel("Name"))
	}
	if option.IsActive != nil {
		query = query.Where(squirrel.Eq{"IsActive": *option.IsActive})
	}
	if option.Slug != nil {
		query = query.Where(option.Slug.ToSquirrel("Slug"))
	}
	if option.Currency != nil {
		query = query.Where(option.Currency.ToSquirrel("Currency"))
	}

	return query.ToSql()
}

// GetbyOption finds and returns 1 channel filtered using given options
func (cs *SqlChannelStore) GetbyOption(option *channel.ChannelFilterOption) (*channel.Channel, error) {

	queryString, args, err := cs.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}

	var res channel.Channel
	err = cs.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ChannelTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find channel by given options")
	}

	return &res, nil
}

// FilterByOption returns a list of channels with given option
func (cs *SqlChannelStore) FilterByOption(option *channel.ChannelFilterOption) ([]*channel.Channel, error) {

	queryString, args, err := cs.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*channel.Channel
	_, err = cs.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find channels with given option")
	}

	return res, nil
}
