package shipping

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingZoneStore struct {
	store.Store
}

func NewSqlShippingZoneStore(s store.Store) store.ShippingZoneStore {
	return &SqlShippingZoneStore{s}
}

func (s *SqlShippingZoneStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
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

func (s *SqlShippingZoneStore) ScanFields(shippingZone shipping.ShippingZone) []interface{} {
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
func (s *SqlShippingZoneStore) Upsert(shippingZone *shipping.ShippingZone) (*shipping.ShippingZone, error) {
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
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + store.ShippingZoneTableName + "(" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
		_, err = s.GetMasterX().NamedExec(query, shippingZone)

	} else {
		query := "UPDATE " + store.ShippingZoneTableName + " SET " + s.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result

		result, err = s.GetMasterX().NamedExec(query, shippingZone)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert shipping zone with id=%s", shippingZone.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple shipping zones were updated: %d instead of 1", numUpdated)
	}

	return shippingZone, nil
}

// Get finds 1 shipping zone for given shippingZoneID
func (s *SqlShippingZoneStore) Get(shippingZoneID string) (*shipping.ShippingZone, error) {
	var res shipping.ShippingZone
	err := s.GetReplicaX().Get(&res, "SELECT * FROM "+store.ShippingZoneTableName+" WHERE Id = ?", shippingZoneID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShippingZoneTableName, shippingZoneID)
		}
		return nil, errors.Wrapf(err, "failed to find shipping zone with id=%s", shippingZoneID)
	}

	return &res, nil
}

// FilterByOption finds a list of shipping zones based on given option
func (s *SqlShippingZoneStore) FilterByOption(option *shipping.ShippingZoneFilterOption) ([]*shipping.ShippingZone, error) {
	selectFields := s.ModelFields(store.ShippingZoneTableName + ".")
	if option.SelectRelatedThroughData {
		selectFields = append(selectFields, "WarehouseShippingZones.WarehouseID")
	}

	query := s.GetQueryBuilder().
		Select(selectFields...).
		From(store.ShippingZoneTableName).
		OrderBy(store.TableOrderingMap[store.ShippingZoneTableName])

	// check option id
	if option != nil && option.Id != nil {
		query = query.Where(option.Id)
	}
	if option != nil && option.DefaultValue != nil {
		query = query.Where(squirrel.Eq{"ShippingZones.Default": *option.DefaultValue})
	}
	if option.WarehouseID != nil {
		query = query.
			InnerJoin(store.WarehouseShippingZoneTableName + " ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID").
			Where(option.WarehouseID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}
	rows, err := s.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping zones with given options")
	}
	var (
		shippingZone           shipping.ShippingZone
		returningShippingZones shipping.ShippingZones
		warehouseID            string
		shippingZonesMap       = map[string]*shipping.ShippingZone{} // shippingZonesMap is a map with keys are shipping zones's ids
		scanFields             = s.ScanFields(shippingZone)
	)

	if option.SelectRelatedThroughData {
		scanFields = append(scanFields, &warehouseID)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to to scan a row contains shipping zones")
		}

		copiedShippingZone := shippingZone.DeepCopy()

		if _, exist := shippingZonesMap[copiedShippingZone.Id]; !exist {
			returningShippingZones = append(returningShippingZones, copiedShippingZone)
			shippingZonesMap[copiedShippingZone.Id] = copiedShippingZone
		}
		duplicateWarehouseID := warehouseID
		shippingZonesMap[copiedShippingZone.Id].RelativeWarehouseIDs = append(shippingZonesMap[copiedShippingZone.Id].RelativeWarehouseIDs, duplicateWarehouseID)
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows of shipping zones")
	}

	return returningShippingZones, nil
}
