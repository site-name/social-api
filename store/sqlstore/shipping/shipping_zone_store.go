package shipping

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlShippingZoneStore struct {
	store.Store
}

func NewSqlShippingZoneStore(s store.Store) store.ShippingZoneStore {
	return &SqlShippingZoneStore{s}
}

func (s *SqlShippingZoneStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Name",
		"Countries",
		"Default",
		"Description",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (s *SqlShippingZoneStore) ScanFields(shippingZone *model.ShippingZone) []interface{} {
	return []interface{}{
		&shippingZone.Id,
		&shippingZone.Name,
		&shippingZone.Countries,
		&shippingZone.Default,
		&shippingZone.Description,
		&shippingZone.Metadata,
		&shippingZone.PrivateMetadata,
	}
}

// Upsert depends on given shipping zone's Id to decide update or insert the zone
func (s *SqlShippingZoneStore) Upsert(tran *gorm.DB, shippingZone *model.ShippingZone) (*model.ShippingZone, error) {
	var isSaving bool

	if shippingZone.Id == "" {
		isSaving = true
		shippingZone.PreSave()
	} else {
		shippingZone.PreUpdate()
	}

	if err := shippingZone.IsValid(); err != nil {
		return nil, err
	}

	var (
		err    error
		result sql.Result
		runner = s.GetMaster()
	)
	if tran != nil {
		runner = tran
	}

	if isSaving {
		query := "INSERT INTO " + model.ShippingZoneTableName + "(" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
		result, err = runner.NamedExec(query, shippingZone)

	} else {
		query := "UPDATE " + model.ShippingZoneTableName + " SET " + s.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		result, err = runner.NamedExec(query, shippingZone)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert shipping zone with id=%s", shippingZone.Id)
	}
	numUpserted, _ := result.RowsAffected()

	if numUpserted > 1 {
		return nil, errors.Errorf("%d shipping zone(s) upserted instead of 1", numUpserted)
	}

	return shippingZone, nil
}

// Get finds 1 shipping zone for given shippingZoneID
func (s *SqlShippingZoneStore) Get(shippingZoneID string) (*model.ShippingZone, error) {
	var res model.ShippingZone
	err := s.GetReplica().Get(&res, "SELECT * FROM "+model.ShippingZoneTableName+" WHERE Id = ?", shippingZoneID)
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
	selectFields := s.ModelFields(model.ShippingZoneTableName + ".")
	if option.SelectRelatedWarehouseIDs {
		selectFields = append(selectFields, "WarehouseShippingZones.WarehouseID")
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
		if opt != nil {
			query = query.Where(opt)
		}
	}

	if option.WarehouseID != nil || option.SelectRelatedWarehouseIDs {
		query = query.InnerJoin(model.WarehouseShippingZoneTableName + " ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID")
	}
	if option.ChannelID != nil {
		query = query.InnerJoin(model.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}
	rows, err := s.GetReplica().Query(queryString, args...)
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
			warehouseID  string
			scanFields   = s.ScanFields(&shippingZone)
		)
		if option.SelectRelatedWarehouseIDs {
			scanFields = append(scanFields, &warehouseID)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to to scan a row contains shipping zones")
		}

		if _, met := shippingZonesMap[shippingZone.Id]; !met {
			returningShippingZones = append(returningShippingZones, &shippingZone)
			shippingZonesMap[shippingZone.Id] = &shippingZone
		}
		if option.SelectRelatedWarehouseIDs {
			shippingZonesMap[shippingZone.Id].RelativeWarehouseIDs = append(shippingZonesMap[shippingZone.Id].RelativeWarehouseIDs, warehouseID)
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
		if opt != nil {
			query = query.Where(opt)
		}
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
	err = s.GetReplica().Select(&res, queryStr, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of shipping zone by given options")
	}

	return res, nil
}

func (s *SqlShippingZoneStore) Delete(transaction *gorm.DB, conditions *model.ShippingZoneFilterOption) (int64, error) {
	if conditions.Conditions == nil {
		return 0, errors.New("please provide conditions to delete")
	}

	query, args, err := s.GetQueryBuilder().Delete(model.ShippingZoneTableName).Where(conditions.Conditions).ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Delete_ToSql")
	}

	runner := s.GetMaster()
	if transaction != nil {
		runner = transaction
	}

	result, err := runner.Exec(query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete shipping zones by given options")
	}

	numDeleted, _ := result.RowsAffected()

	return numDeleted, nil
}
