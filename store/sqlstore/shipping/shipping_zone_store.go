package shipping

import (
	"database/sql"
	"fmt"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlShippingZoneStore struct {
	store.Store
}

func NewSqlShippingZoneStore(s store.Store) store.ShippingZoneStore {
	return &SqlShippingZoneStore{s}
}

func (s *SqlShippingZoneStore) Upsert(tran boil.ContextTransactor, shippingZone model.ShippingZone) (*model.ShippingZone, error) {
	if tran == nil {
		tran = s.GetMaster()
	}

	isSaving := shippingZone.ID == ""
	if isSaving {
		model_helper.ShippingZonePreSave(&shippingZone)
	} else {
		model_helper.ShippingZoneCommonPre(&shippingZone)
	}

	if err := model_helper.ShippingZoneIsValid(shippingZone); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = shippingZone.Insert(tran, boil.Infer())
	} else {
		_, err = shippingZone.Update(tran, boil.Blacklist(model.ShippingZoneColumns.CreatedAt))
	}

	if err != nil {
		return nil, err
	}

	return &shippingZone, nil
}

func (s *SqlShippingZoneStore) Get(shippingZoneID string) (*model.ShippingZone, error) {
	zone, err := model.FindShippingZone(s.GetReplica(), shippingZoneID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ShippingZones, shippingZoneID)
		}
		return nil, err
	}

	return zone, nil
}

func (s *SqlShippingZoneStore) commonQueryConditionBuilder(option model_helper.ShippingZoneFilterOption) []qm.QueryMod {
	result := option.Conditions
	if option.WarehouseID != nil {
		result = append(
			result,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.WarehouseShippingZones, model.WarehouseShippingZoneTableColumns.ShippingZoneID, model.ShippingZoneTableColumns.ID)),
			option.WarehouseID,
		)
	}
	if option.ChannelID != nil {
		result = append(
			result,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZoneChannels, model.ShippingZoneChannelTableColumns.ShippingZoneID, model.ShippingZoneTableColumns.ID)),
			option.ChannelID,
		)
	}

	return result
}

func (s *SqlShippingZoneStore) FilterByOption(option model_helper.ShippingZoneFilterOption) (model.ShippingZoneSlice, error) {
	conditions := s.commonQueryConditionBuilder(option)
	return model.ShippingZones(conditions...).All(s.GetReplica())
}

func (s *SqlShippingZoneStore) CountByOptions(options model_helper.ShippingZoneFilterOption) (int64, error) {
	conditions := s.commonQueryConditionBuilder(options)
	return model.ShippingZones(conditions...).Count(s.GetReplica())
}

func (s *SqlShippingZoneStore) Delete(transaction boil.ContextTransactor, ids []string) (int64, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	res, err := model.ShippingZones(model.ShippingZoneWhere.ID.IN(ids)).DeleteAll(transaction)
	if err != nil {
		return 0, err
	}

	return res, nil
}
