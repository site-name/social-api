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

	annotations := model_helper.AnnotationAggregator{}
	if conds.AnnotateHasOrders {
		annotations[model_helper.ChannelAnnotationKeys.HasOrders] = fmt.Sprintf(`EXISTS ( SELECT (1) AS "a" FROM %s WHERE %s = %s LIMIT 1 )`, model.TableNames.Orders, model.OrderColumns.ChannelID, model.ChannelTableColumns.ID)
	}
	res = append(res, annotations)

	return res
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

func (s *SqlChannelStore) FilterByOptions(conds model_helper.ChannelFilterOptions) (model.ChannelSlice, error) {
	mods := commonQueryBuilder(conds)
	return model.Channels(mods...).All(s.GetReplica())
}
