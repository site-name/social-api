package channel

import (
	"database/sql"
	"fmt"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlChannelStore struct {
	store.Store
}

func NewSqlChannelStore(sqlStore store.Store) store.ChannelStore {
	return &SqlChannelStore{sqlStore}
}

func (s *SqlChannelStore) DeleteChannels(tx boil.ContextTransactor, ids []string) error {
	if tx == nil {
		tx = s.GetMaster()
	}

	_, err := model.Channels(model.ChannelWhere.ID.IN(ids)).DeleteAll(tx)
	return err
}

func (s *SqlChannelStore) Get(id string) (*model.Channel, error) {
	channel, err := model.FindChannel(s.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Channels, id)
		}
		return nil, err
	}
	return channel, nil
}

func (s *SqlChannelStore) Find(conds model_helper.ChannelFilterOptions) (model.ChannelSlice, error) {
	mods := commonQueryBuilder(conds)
	return model.Channels(mods...).All(s.GetReplica())
}

func (s *SqlChannelStore) GetByOptions(conds model_helper.ChannelFilterOptions) (*model.Channel, error) {
	mods := commonQueryBuilder(conds)
	channel, err := model.Channels(mods...).One(s.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Channels, "options")
		}
		return nil, err
	}
	return channel, nil
}

func (s *SqlChannelStore) Upsert(tx boil.ContextTransactor, channel model.Channel) (*model.Channel, error) {
	isSaving := channel.ID == ""
	if isSaving {
		model_helper.ChannelPreSave(&channel)
	} else {
		model_helper.ChannelCommonPre(&channel)
	}

	if err := model_helper.ChannelIsValid(channel); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = channel.Insert(tx, boil.Infer())
	} else {
		_, err = channel.Update(tx, boil.Infer())
	}

	if err != nil {
		if s.IsUniqueConstraintError(err, []string{model.ChannelColumns.Slug, "channels_slug_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Channels, model.ChannelColumns.Slug, channel.Slug)
		}
		return nil, err
	}

	return &channel, nil
}

func commonQueryBuilder(conds model_helper.ChannelFilterOptions) []qm.QueryMod {
	res := []qm.QueryMod{}

	for _, cond := range conds.Conditions {
		if cond != nil {
			res = append(res, cond)
		}
	}
	if conds.ShippingZoneID != nil {
		res = append(
			res,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZoneChannels, model.ShippingZoneChannelTableColumns.ChannelID, model.ChannelTableColumns.ID)),
			conds.ShippingZoneID,
		)
	}
	if conds.VoucherID != nil {
		res = append(res,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.VoucherChannelListings, model.VoucherChannelListingColumns.ChannelID, model.ChannelTableColumns.ID)),
			conds.VoucherID,
		)
	}

	return res
}

// // FilterByOption returns a list of channels with given option
// func (cs *SqlChannelStore) Find(option model_helper.ChannelFilterOption) (model.ChannelSlice, error) {
// 	selectFields := []string{model.ChannelTableName + ".*"}
// 	if option.AnnotateHasOrders {
// 		selectFields = append(selectFields, `EXISTS ( SELECT (1) AS "a" FROM Orders WHERE Orders.ChannelID = Channels.Id LIMIT 1 ) AS HasOrders`)
// 	}

// 	query := cs.GetQueryBuilder().
// 		Select(selectFields...).
// 		From(model.ChannelTableName).
// 		Where(option.Conditions)

// 	// parse options
// 	if option.ShippingZoneChannels_ShippingZoneID != nil {
// 		query = query.
// 			InnerJoin(model.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ChannelID = Channels.Id").
// 			Where(option.ShippingZoneChannels_ShippingZoneID)
// 	}
// 	if option.VoucherChannelListing_VoucherID != nil {
// 		query = query.
// 			InnerJoin(model.VoucherChannelListingTableName + " ON VoucherChannelListings.ChannelID = Channels.Id").
// 			Where(option.VoucherChannelListing_VoucherID)
// 	}
// 	if option.Limit > 0 {
// 		query = query.Limit(uint64(option.Limit))
// 	}

// 	queryString, args, err := query.ToSql()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "FilterByOption_ToSql")
// 	}

// 	var res model.Channels

// 	rows, err := cs.GetReplica().Raw(queryString, args...).Rows()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to find channels with given option")
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var (
// 			hasOrder   bool
// 			channel    model.Channel
// 			scanFields = cs.ScanFields(&channel)
// 		)
// 		if option.AnnotateHasOrders {
// 			scanFields = append(scanFields, &hasOrder)
// 		}

// 		err = rows.Scan(scanFields...)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "failed to scan channel row")
// 		}

// 		if option.AnnotateHasOrders {
// 			channel.SetHasOrders(hasOrder)
// 		}
// 		res = append(res, &channel)
// 	}

// 	return res, nil
// }
