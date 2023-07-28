package warehouse

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlWareHouseStore struct {
	store.Store
}

func NewSqlWarehouseStore(s store.Store) store.WarehouseStore {
	return &SqlWareHouseStore{s}
}

func (ws *SqlWareHouseStore) ScanFields(wareHouse *model.WareHouse) []interface{} {
	return []interface{}{
		&wareHouse.Id,
		&wareHouse.Name,
		&wareHouse.Slug,
		&wareHouse.AddressID,
		&wareHouse.Email,
		&wareHouse.ClickAndCollectOption,
		&wareHouse.IsPrivate,
		&wareHouse.Metadata,
		&wareHouse.PrivateMetadata,
	}
}

func (ws *SqlWareHouseStore) Save(wh *model.WareHouse) (*model.WareHouse, error) {
	if err := ws.GetMaster().Create(wh).Error; err != nil {
		if ws.IsUniqueConstraintError(err, []string{"slug", "warehouses_slug_key", "idx_warehouses_slug_unique"}) {
			return nil, store.NewErrInvalidInput("Warehouses", "Slug", wh.Slug)
		}
		return nil, errors.Wrap(err, "failed to save Warehouse")
	}

	return wh, nil
}

func (ws *SqlWareHouseStore) Update(warehouse *model.WareHouse) (*model.WareHouse, error) {
	err := ws.GetMaster().Model(warehouse).Updates(warehouse).Error
	if err != nil {
		if ws.IsUniqueConstraintError(err, []string{"Slug", "warehouses_slug_key", "idx_warehouses_slug_unique"}) {
			return nil, store.NewErrInvalidInput("Warehouses", "Slug", warehouse.Slug)
		}
		return nil, errors.Wrap(err, "failed to update warehouse")
	}

	return warehouse, nil
}

// NOTE: if option is nil, all warehouses query is returned.
func (ws *SqlWareHouseStore) commonQueryBuilder(option *model.WarehouseFilterOption) squirrel.SelectBuilder {
	selectFields := []string{model.WarehouseTableName + ".*"}
	if option.SelectRelatedAddress {
		selectFields = append(selectFields, model.AddressTableName+".*")
	}

	query := ws.GetQueryBuilder().
		Select(selectFields...).
		From(model.WarehouseTableName).Where(option.Conditions)

	for _, opt := range []squirrel.Sqlizer{
		option.ShippingZonesCountries,
		option.ShippingZonesId,
	} {
		query = query.Where(opt)
	}

	if option.ShippingZonesCountries != nil || option.ShippingZonesId != nil {
		query = query.
			InnerJoin(model.WarehouseShippingZoneTableName + " ON Warehouses.Id = WarehouseShippingZones.WarehouseID").
			InnerJoin(model.ShippingZoneTableName + " ON WarehouseShippingZones.ShippingZoneID = ShippingZones.Id")
	}
	if option.SelectRelatedAddress || option.Search != "" {
		query = query.InnerJoin(model.AddressTableName + " ON (Addresses.Id = Warehouses.AddressID)")

		if option.Search != "" {
			expr := "%" + option.Search + "%"

			query = query.Where(squirrel.Or{
				squirrel.ILike{model.WarehouseTableName + ".Name": expr},
				squirrel.ILike{model.WarehouseTableName + ".Email": expr},

				squirrel.ILike{model.AddressTableName + ".CompanyName": expr},
				squirrel.ILike{model.AddressTableName + ".StreetAddress1": expr},
				squirrel.ILike{model.AddressTableName + ".StreetAddress2": expr},
				squirrel.ILike{model.AddressTableName + ".City": expr},
				squirrel.ILike{model.AddressTableName + ".PostalCode": expr},
				squirrel.ILike{model.AddressTableName + ".Phone": expr},
			})
		}
	}

	if option.Distinct {
		query = query.Distinct()
	}

	return query
}

// GetByOption finds and returns a warehouse filtered given option
func (ws *SqlWareHouseStore) GetByOption(option *model.WarehouseFilterOption) (*model.WareHouse, error) {
	query, args, err := ws.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}
	var (
		res        model.WareHouse
		address    model.Address
		scanFields = ws.ScanFields(&res)
	)
	if option.SelectRelatedAddress {
		scanFields = append(scanFields, ws.Address().ScanFields(&address)...)
	}

	err = ws.GetReplica().Raw(query, args...).Row().Scan(scanFields...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.NewErrNotFound(model.WarehouseTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find warehouse with given option")
	}

	if option.SelectRelatedAddress {
		res.SetAddress(&address)
	}

	// check if we need to prefetch shipping zones:
	// 1) prefetching shipping zones is required
	// 2) returning warehouse is valid
	if option.PrefetchShippingZones {
		queryString, args, err := ws.GetQueryBuilder().
			Select(model.ShippingZoneTableName+".*").
			From(model.ShippingZoneTableName).
			InnerJoin(model.WarehouseShippingZoneTableName+" ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID").
			Where("WarehouseShippingZones.WarehouseID = ?", res.Id).
			ToSql()

		if err != nil {
			return nil, errors.Wrap(err, "GetByOption_Warehouse_ToSql")
		}

		var shippingZones model.ShippingZones
		err = ws.GetReplica().Raw(queryString, args...).Scan(&shippingZones).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to find shipping zones by warehouse ids")
		}

		res.ShippingZones = shippingZones
	}

	return &res, nil
}

// FilterByOprion returns a slice of warehouses with given option
func (wh *SqlWareHouseStore) FilterByOprion(option *model.WarehouseFilterOption) ([]*model.WareHouse, error) {
	query, args, err := wh.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}
	rows, err := wh.GetReplica().Raw(query, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find warehouses with given option")
	}
	defer rows.Close()

	var returningWarehouses model.Warehouses

	for rows.Next() {
		var (
			wareHouse  model.WareHouse
			address    model.Address
			scanFields = wh.ScanFields(&wareHouse)
		)
		if option.SelectRelatedAddress {
			scanFields = append(scanFields, wh.Address().ScanFields(&address)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of warehouse and address")
		}

		if option.SelectRelatedAddress {
			wareHouse.SetAddress(&address)
		}
		returningWarehouses = append(returningWarehouses, &wareHouse)
	}

	// check if we need prefetch related shipping zones:
	if option.PrefetchShippingZones && len(returningWarehouses) > 0 {
		query, args, err = wh.GetQueryBuilder().
			Select(model.ShippingZoneTableName + ".*").
			Column("WarehouseShippingZones.WarehouseID AS PrefetchRelatedWarehouseID"). // <- this column selection helps determine which shipping zone is related to which warehouse
			From(model.ShippingZoneTableName).
			InnerJoin(model.WarehouseShippingZoneTableName + " ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID").
			Where(squirrel.Eq{"PrefetchRelatedWarehouseID": returningWarehouses.IDs()}).
			ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "FilerByOption_Prefetch_ToSql")
		}

		rows, err := wh.GetReplica().Raw(query, args...).Rows()
		if err != nil {
			return nil, errors.Wrap(err, "failed to find shipping zones of warehouses")
		}
		defer rows.Close()
		var warehousesMap = lo.SliceToMap(returningWarehouses, func(w *model.WareHouse) (string, *model.WareHouse) { return w.Id, w })

		for rows.Next() {
			var (
				shippingZone       model.ShippingZone
				relatedWarehouseID string
				scanFields         = append(wh.ShippingZone().ScanFields(&shippingZone), &relatedWarehouseID)
			)

			err = rows.Scan(scanFields...)
			if err != nil {
				return nil, errors.Wrap(err, "failed to scan a row of shipping zone and warehouse id")
			}

			if warehousesMap[relatedWarehouseID] != nil {
				warehousesMap[relatedWarehouseID].ShippingZones = append(warehousesMap[relatedWarehouseID].ShippingZones, &shippingZone)
			}
		}
	}

	return returningWarehouses, nil
}

// WarehouseByStockID returns 1 warehouse by given stock id
func (ws *SqlWareHouseStore) WarehouseByStockID(stockID string) (*model.WareHouse, error) {
	var res model.WareHouse
	err := ws.GetReplica().Raw(
		`SELECT `+model.WarehouseTableName+".*"+`
		FROM `+model.WarehouseTableName+`
		INNER JOIN `+model.StockTableName+` ON Stocks.WarehouseID = Warehouses.Id
		WHERE Stocks.Id = ?`,
		stockID,
	).
		Scan(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.WarehouseTableName, "StockID="+stockID)
		}
		return nil, errors.Wrapf(err, "failed to find warehouse with StockID=%s", stockID)
	}

	return &res, nil
}

func (ws *SqlWareHouseStore) ApplicableForClickAndCollectNoQuantityCheck(checkoutLines model.CheckoutLines, country model.CountryCode) (model.Warehouses, error) {
	stocks, err := ws.Stock().FilterByOption(&model.StockFilterOption{
		SelectRelatedProductVariant: true,
		Conditions:                  squirrel.Eq{model.StockTableName + ".ProductVariantID": checkoutLines.VariantIDs()},
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks")
	}

	return ws.forCountryLinesAndStocks(checkoutLines, stocks, country)
}

func (w *SqlWareHouseStore) Delete(transaction *gorm.DB, ids ...string) error {
	if transaction == nil {
		transaction = w.GetMaster()
	}

	err := transaction.Raw("DELETE FROM "+model.WarehouseTableName+" WHERE Id IN ?", ids).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete warehouse(s) by given ids")
	}

	return nil
}

func (ws *SqlWareHouseStore) ApplicableForClickAndCollectCheckoutLines(checkoutLines model.CheckoutLines, country model.CountryCode) (model.Warehouses, error) {
	panic("not implemented")
}

func (s *SqlWareHouseStore) ApplicableForClickAndCollectOrderLines(orderLines model.OrderLines, country model.CountryCode) (model.Warehouses, error) {
	panic("not implemented")
}

func (ws *SqlWareHouseStore) forCountryLinesAndStocks(checkoutLines model.CheckoutLines, stocks model.Stocks, country model.CountryCode) (model.Warehouses, error) {
	panic("not implemented")
}
