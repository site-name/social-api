package channel

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlChannelStore struct {
	store.Store
}

func NewSqlChannelStore(sqlStore store.Store) store.ChannelStore {
	return &SqlChannelStore{sqlStore}
}

func (cs *SqlChannelStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
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
		&ch.Name,
		&ch.IsActive,
		&ch.Slug,
		&ch.Currency,
		&ch.DefaultCountry,
	}
}

func (s *SqlChannelStore) Upsert(transaction store_iface.SqlxTxExecutor, channel *model.Channel) (*model.Channel, error) {
	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	isSaving := false
	if !model.IsValidId(channel.Id) {
		channel.Id = ""
		isSaving = true
		channel.PreSave()
	} else {
		channel.PreUpdate()
	}

	if appErr := channel.IsValid(); appErr != nil {
		return nil, appErr
	}

	var (
		err    error
		result sql.Result
	)
	if isSaving {
		query := "INSERT INTO " + store.ChannelTableName + "(" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
		result, err = runner.NamedExec(query, channel)
	} else {
		query := "UPDATE " + store.ChannelTableName + " SET " + s.ModelFields(":").Join(",") + " WHERE Id=:Id"
		result, err = runner.NamedExec(query, channel)
	}

	if err != nil {
		if s.IsUniqueConstraintError(err, []string{"Slug", "channels_slug_key", "idx_channels_slug_unique"}) {
			return nil, store.NewErrInvalidInput(store.ChannelTableName, "Slug", channel.Slug)
		}
		return nil, errors.Wrap(err, "failed to upsert channel")
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected != 1 {
		return nil, errors.Errorf("%d rows affected instead of 1", rowsAffected)
	}

	return channel, nil
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
	selectFields := cs.ModelFields(store.ChannelTableName + ".")
	if option.AnnotateHasOrders {
		selectFields = append(selectFields, `EXISTS ( SELECT (1) AS "a" FROM Orders WHERE Orders.ChannelID = Channels.Id LIMIT 1 ) AS HasOrders`)
	}

	query := cs.GetQueryBuilder().
		Select(selectFields...).
		From(store.ChannelTableName)

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.Name != nil {
		query = query.Where(option.Name)
	}
	if option.IsActive != nil {
		query = query.Where(squirrel.Eq{"Channels.IsActive": *option.IsActive})
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
	if option.ShippingZoneChannels_ShippingZoneID != nil {
		query = query.InnerJoin(store.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ChannelID = Channels.Id").
			Where(option.ShippingZoneChannels_ShippingZoneID)
	}

	return query.ToSql()
}

// GetbyOption finds and returns 1 channel filtered using given options
func (cs *SqlChannelStore) GetbyOption(option *model.ChannelFilterOption) (*model.Channel, error) {
	queryString, args, err := cs.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}

	var (
		res        model.Channel
		hasOrder   bool
		row        = cs.GetReplicaX().QueryRowX(queryString, args...)
		scanFields = cs.ScanFields(&res)
	)

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

	var (
		res        model.Channels
		hasOrder   bool
		channel    model.Channel
		scanFields = cs.ScanFields(&channel)
	)
	if option.AnnotateHasOrders {
		scanFields = append(scanFields, &hasOrder)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan channel row")
		}

		if option.AnnotateHasOrders {
			channel.SetHasOrders(hasOrder)
		}
		res = append(res, channel.DeepCopy())
	}

	return res, nil
}
