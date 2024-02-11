package shipping

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlShippingZoneStore struct {
	store.Store
}

func NewSqlShippingZoneStore(s store.Store) store.ShippingZoneStore {
	return &SqlShippingZoneStore{s}
}

// Upsert depends on given shipping zone's Id to decide update or insert the zone
func (s *SqlShippingZoneStore) Upsert(tran *gorm.DB, shippingZone *model.ShippingZone) (*model.ShippingZone, error) {
	if tran == nil {
		tran = s.GetMaster()
	}

	err := tran.Save(shippingZone).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert shipping zone")
	}
	return shippingZone, nil
}

// Get finds 1 shipping zone for given shippingZoneID
func (s *SqlShippingZoneStore) Get(shippingZoneID string) (*model.ShippingZone, error) {
	var res model.ShippingZone
	err := s.GetReplica().First(&res, "Id = ?", shippingZoneID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ShippingZoneTableName, shippingZoneID)
		}
		return nil, errors.Wrapf(err, "failed to find shipping zone with id=%s", shippingZoneID)
	}

	return &res, nil
}

// FilterByOption finds a list of shipping zones based on given option
func (s *SqlShippingZoneStore) FilterByOption(option *model.ShippingZoneFilterOption) ([]*model.ShippingZone, error) {
	selectFields := []string{model.ShippingZoneTableName + ".*"}
	if option.SelectRelatedWarehouses {
		selectFields = append(selectFields, model.WarehouseTableName+".*")
	}

	query := s.GetQueryBuilder().
		Select(selectFields...).
		From(model.ShippingZoneTableName)

	// parse options
	for _, opt := range []squirrel.Sqlizer{
		option.Conditions,
		option.WarehouseID,
		option.ChannelID,
	} {
		query = query.Where(opt)
	}

	if option.WarehouseID != nil || option.SelectRelatedWarehouses {
		query = query.InnerJoin(model.WarehouseShippingZoneTableName + " ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID")
	}
	if option.ChannelID != nil {
		query = query.InnerJoin(model.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}
	rows, err := s.GetReplica().Raw(queryString, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping zones with given options")
	}
	defer rows.Close()

	var (
		returningShippingZones model.ShippingZones
		shippingZonesMap       = map[string]*model.ShippingZone{} // keys are shipping zones' ids
	)

	for rows.Next() {
		var (
			shippingZone model.ShippingZone
			warehouse    model.WareHouse
			scanFields   = s.ScanFields(&shippingZone)
		)
		if option.SelectRelatedWarehouses {
			scanFields = append(scanFields, s.Warehouse().ScanFields(&warehouse)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to to scan a row contains shipping zones")
		}

		if _, met := shippingZonesMap[shippingZone.Id]; !met {
			returningShippingZones = append(returningShippingZones, &shippingZone)
			shippingZonesMap[shippingZone.Id] = &shippingZone
		}
		if option.SelectRelatedWarehouses {
			shippingZonesMap[shippingZone.Id].Warehouses = append(shippingZonesMap[shippingZone.Id].Warehouses, &warehouse)
		}
	}

	return returningShippingZones, nil
}

func (s *SqlShippingZoneStore) CountByOptions(options *model.ShippingZoneFilterOption) (int64, error) {
	query := s.GetQueryBuilder().Select("COUNT( DISTINCT ShippingZones.Id)").From(model.ShippingZoneTableName)

	// parse options
	for _, opt := range []squirrel.Sqlizer{
		options.Conditions,
		options.WarehouseID,
		options.ChannelID,
	} {
		query = query.Where(opt)
	}

	if options.WarehouseID != nil {
		query = query.InnerJoin(model.WarehouseShippingZoneTableName + " ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID")
	}
	if options.ChannelID != nil {
		query = query.InnerJoin(model.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id")
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "CountByOptions_ToSql")
	}

	var res int64
	err = s.GetReplica().Raw(queryStr, args...).Scan(&res).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of shipping zone by given options")
	}

	return res, nil
}

func (s *SqlShippingZoneStore) Delete(transaction *gorm.DB, conditions *model.ShippingZoneFilterOption) (int64, error) {
	query, args, err := s.GetQueryBuilder().Delete(model.ShippingZoneTableName).Where(conditions.Conditions).ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Delete_ToSql")
	}

	if transaction == nil {
		transaction = s.GetMaster()
	}

	result := transaction.Raw(query, args...)
	if result.Error != nil {
		return 0, errors.Wrap(result.Error, "failed to delete shipping zones by given options")
	}

	return result.RowsAffected, nil
}

func (s *SqlShippingZoneStore) ToggleRelations(transaction *gorm.DB, zones model.ShippingZones, warehouseIds, channelIds []string, delete bool) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	var relationsMap = map[string]any{
		"Channels":   lo.Map(channelIds, func(id string, _ int) *model.Channel { return &model.Channel{Id: id} }),
		"Warehouses": lo.Map(warehouseIds, func(id string, _ int) *model.WareHouse { return &model.WareHouse{Id: id} }),
	}

	for _, shippingZone := range zones {
		if shippingZone != nil {
			for assoName, relations := range relationsMap {
				association := transaction.Model(shippingZone).Association(assoName)
				var err error
				switch {
				case delete:
					err = association.Delete(relations)
				default:
					err = association.Append(relations)
				}
				if err != nil {
					return errors.Wrap(err, "failed to toggle "+strings.ToLower(assoName)+" to shipping zone with id = "+shippingZone.Id)
				}
			}
		}
	}

	return nil
}
