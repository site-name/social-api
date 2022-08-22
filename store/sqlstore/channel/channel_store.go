package channel

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/store"
)

type SqlChannelStore struct {
	store.Store
}

func NewSqlChannelStore(sqlStore store.Store) store.ChannelStore {
	return &SqlChannelStore{sqlStore}
}

func (cs *SqlChannelStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"ShopID",
		"Name",
		"IsActive",
		"Slug",
		"Currency",
		"DefaultCountry",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (cs *SqlChannelStore) ScanFields(ch channel.Channel) []interface{} {
	return []interface{}{
		&ch.Id,
		&ch.ShopID,
		&ch.Name,
		&ch.IsActive,
		&ch.Slug,
		&ch.Currency,
		&ch.DefaultCountry,
	}
}

func (cs *SqlChannelStore) Save(ch *channel.Channel) (*channel.Channel, error) {
	ch.PreSave()
	if err := ch.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.ChannelTableName + "(" + cs.ModelFields("").Join(",") + ") VALUES (" + cs.ModelFields(":").Join(",") + ")"
	if _, err := cs.GetMasterX().NamedExec(query, ch); err != nil {
		if cs.IsUniqueConstraintError(err, []string{"Slug", "channels_slug_key", "idx_channels_slug_unique"}) {
			return nil, store.NewErrInvalidInput("Channel", "Slug", ch.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save channel with id=%s", ch.Id)
	}

	return ch, nil
}

func (cs *SqlChannelStore) Get(id string) (*channel.Channel, error) {
	var channel channel.Channel

	err := cs.GetReplicaX().Get(&channel, "SELECT * FROM "+store.ChannelTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ChannelTableName, id)
		}
		return nil, errors.Wrapf(err, "Failed to get Channel with ChannelID=%s", id)
	}

	return &channel, nil
}

func (cs *SqlChannelStore) commonQueryBuilder(option *channel.ChannelFilterOption) (string, []interface{}, error) {
	query := cs.GetQueryBuilder().
		Select(cs.ModelFields("")...).
		From(store.ChannelTableName).
		OrderBy(store.TableOrderingMap[store.ChannelTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.ShopID != nil {
		query = query.Where(option.ShopID)
	}
	if option.Name != nil {
		query = query.Where(option.Name)
	}
	if option.IsActive != nil {
		query = query.Where(squirrel.Eq{"IsActive": *option.IsActive})
	}
	if option.Slug != nil {
		query = query.Where(option.Slug)
	}
	if option.Currency != nil {
		query = query.Where(option.Currency)
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
	err = cs.GetReplicaX().Get(&res, queryString, args...)
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
	err = cs.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find channels with given option")
	}

	return res, nil
}
