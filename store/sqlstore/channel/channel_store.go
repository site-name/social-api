package channel

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlChannelStore struct {
	store.Store
}

func NewSqlChannelStore(sqlStore store.Store) store.ChannelStore {
	return &SqlChannelStore{sqlStore}
}

func (cs *SqlChannelStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
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

func (cs *SqlChannelStore) ScanFields(ch *model.Channel) []interface{} {
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

func (cs *SqlChannelStore) Save(ch *model.Channel) (*model.Channel, error) {
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

func (cs *SqlChannelStore) Get(id string) (*model.Channel, error) {
	var channel model.Channel

	err := cs.GetReplicaX().Get(&channel, "SELECT * FROM "+store.ChannelTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ChannelTableName, id)
		}
		return nil, errors.Wrapf(err, "Failed to get Channel with ChannelID=%s", id)
	}

	return &channel, nil
}

func (cs *SqlChannelStore) commonQueryBuilder(option *model.ChannelFilterOption) (string, []interface{}, error) {
	query := cs.GetQueryBuilder().
		Select(cs.ModelFields("")...).
		From(store.ChannelTableName)

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
	if option.Extra != nil {
		query = query.Where(option.Extra)
	}
	if option.AnnotateHasOrders {
		query = query.Column(`EXISTS ( SELECT (1) AS "a" FROM Orders WHERE Orders.ChannelID = Channels.Id LIMIT 1 ) AS HasOrders`)
	}

	return query.ToSql()
}

// GetbyOption finds and returns 1 channel filtered using given options
func (cs *SqlChannelStore) GetbyOption(option *model.ChannelFilterOption) (*model.Channel, error) {
	queryString, args, err := cs.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}

	var res model.Channel
	var hasOrder bool

	row := cs.GetReplicaX().QueryRowX(queryString, args...)

	scanFields := cs.ScanFields(&res)
	if option.AnnotateHasOrders {
		scanFields = append(scanFields, &hasOrder)
	}
	err = row.Scan(scanFields...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ChannelTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find channel by given options")
	}

	if option.AnnotateHasOrders {
		res.SetHasOrders(hasOrder)
	}

	return &res, nil
}

// FilterByOption returns a list of channels with given option
func (cs *SqlChannelStore) FilterByOption(option *model.ChannelFilterOption) ([]*model.Channel, error) {
	queryString, args, err := cs.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	rows, err := cs.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find channels with given option")
	}
	defer rows.Close()

	var res model.Channels
	var hasOrder bool
	var channel model.Channel
	var scanFields = cs.ScanFields(&channel)

	if option.AnnotateHasOrders {
		scanFields = append(scanFields, &hasOrder)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan channel row")
		}

		channel.SetHasOrders(hasOrder)
		res = append(res, channel.DeepCopy())
	}

	return res, nil
}
