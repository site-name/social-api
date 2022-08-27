package warehouse

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlWareHouseStore struct {
	store.Store
}

func NewSqlWarehouseStore(s store.Store) store.WarehouseStore {
	return &SqlWareHouseStore{s}
}

func (ws *SqlWareHouseStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"Name",
		"Slug",
		"AddressID",
		"Email",
		"ClickAndCollectOption",
		"IsPrivate",
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

func (ws *SqlWareHouseStore) ScanFields(wareHouse warehouse.WareHouse) []interface{} {
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

func (ws *SqlWareHouseStore) Save(wh *warehouse.WareHouse) (*warehouse.WareHouse, error) {
	wh.PreSave()
	if err := wh.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.WarehouseTableName + "(" + ws.ModelFields("").Join(",") + ") VALUES (" + ws.ModelFields(":").Join(",") + ")"
	if _, err := ws.GetMasterX().NamedExec(query, wh); err != nil {
		if ws.IsUniqueConstraintError(err, []string{"Slug", "warehouses_slug_key", "idx_warehouses_slug_unique"}) {
			return nil, store.NewErrInvalidInput("Warehouses", "Slug", wh.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save Warehouse with Id=%s", wh.Id)
	}

	return wh, nil
}

func (ws *SqlWareHouseStore) Get(id string) (*warehouse.WareHouse, error) {
	var res warehouse.WareHouse
	err := ws.GetReplicaX().Get(
		&res,
		"SELECT * FROM "+store.WarehouseTableName+" WHERE Id = ?",
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Warehouse", id)
		}
		return nil, errors.Wrapf(err, "failed to get warehouse with Id=%s", id)
	}

	return &res, nil
}

func (ws *SqlWareHouseStore) commonQueryBuilder(option *warehouse.WarehouseFilterOption) squirrel.SelectBuilder {
	if option == nil {
		return ws.GetQueryBuilder().
			Select(ws.ModelFields(store.WarehouseTableName + ".")...).
			From(store.WarehouseTableName).
			OrderBy(store.TableOrderingMap[store.WarehouseTableName])
	}

	selectFields := ws.ModelFields(store.WarehouseTableName + ".")
	if option.SelectRelatedAddress {
		selectFields = append(selectFields, ws.Address().ModelFields(store.AddressTableName+".")...)
	}

	query := ws.GetQueryBuilder().
		Select(selectFields...).
		From(store.WarehouseTableName).
		OrderBy(store.TableOrderingMap[store.WarehouseTableName])

	// parse option
	if option.Distinct {
		query = query.Distinct()
	}
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.Name != nil {
		query = query.Where(option.Name)
	}
	if option.Slug != nil {
		query = query.Where(option.Slug)
	}
	if option.AddressID != nil {
		query = query.Where(option.AddressID)
	}
	if option.Email != nil {
		query = query.Where(option.Email)
	}
	if option.ShippingZonesCountries != nil || option.ShippingZonesId != nil {
		query = query.
			InnerJoin(store.WarehouseShippingZoneTableName + " ON Warehouses.Id = WarehouseShippingZones.WarehouseID").
			InnerJoin(store.ShippingZoneTableName + " ON WarehouseShippingZones.ShippingZoneID = ShippingZones.Id")

		if option.ShippingZonesCountries != nil {
			query = query.Where(option.ShippingZonesCountries)
		}
		if option.ShippingZonesId != nil {
			query = query.Where(option.ShippingZonesId)
		}
	}
	if option.SelectRelatedAddress {
		query.InnerJoin(store.AddressTableName + " ON (Addresses.Id = Warehouses.AddressID)")
	}

	return query
}

// GetByOption finds and returns a warehouse filtered given option
func (ws *SqlWareHouseStore) GetByOption(option *warehouse.WarehouseFilterOption) (*warehouse.WareHouse, error) {
	query, args, err := ws.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}
	var (
		res        warehouse.WareHouse
		address    account.Address
		scanFields = ws.ScanFields(res)
	)
	if option.SelectRelatedAddress {
		scanFields = append(scanFields, ws.Address().ScanFields(address)...)
	}

	err = ws.GetReplicaX().QueryRowX(query, args...).Scan(scanFields...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WarehouseTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find warehouse with given option")
	}

	if option.SelectRelatedAddress {
		res.Address = address.DeepCopy()
	}

	// check if we need to prefetch shipping zones:
	// 1) prefetching shipping zones is required
	// 2) returning warehouse is valid
	if option.PrefetchShippingZones && model.IsValidId(res.Id) {
		queryString, args, err := ws.GetQueryBuilder().
			Select(ws.ShippingZone().ModelFields(store.ShippingZoneTableName+".")...).
			From(store.ShippingZoneTableName).
			InnerJoin(store.WarehouseShippingZoneTableName+" ON (ShippingZones.Id = WarehouseShippingZones.ShippingZoneID)").
			Where("WarehouseShippingZones.WarehouseID = ?", res.Id).
			ToSql()

		if err != nil {
			return nil, errors.Wrap(err, "GetByOption_ToSql")
		}

		rows, err := ws.GetReplicaX().QueryX(queryString, args...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find shipping zones related to returning warehouse")
		}
		var (
			shippingZone shipping.ShippingZone
			scanFields   = ws.ShippingZone().ScanFields(shippingZone)
		)

		for rows.Next() {
			err = rows.Scan(scanFields...)
			if err != nil {
				return nil, errors.Wrap(err, "failed to scan a row of shipping zone")
			}

			res.ShippingZones = append(res.ShippingZones, shippingZone.DeepCopy())
		}

		if err = rows.Close(); err != nil {
			return nil, errors.Wrap(err, "failed to close rows of shipping zones")
		}
	}

	return &res, nil
}

// FilterByOprion returns a slice of warehouses with given option
func (wh *SqlWareHouseStore) FilterByOprion(option *warehouse.WarehouseFilterOption) ([]*warehouse.WareHouse, error) {
	query, args, err := wh.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}
	rows, err := wh.GetReplicaX().QueryX(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find warehouses with given option")
	}

	var (
		returningWarehouses warehouse.Warehouses
		warehousesMap       = map[string]*warehouse.WareHouse{} // keys are warehouse IDs
		wareHouse           warehouse.WareHouse
		address             account.Address
		scanFields          = wh.ScanFields(wareHouse)
	)
	if option.SelectRelatedAddress {
		scanFields = append(scanFields, wh.Address().ScanFields(address)...)
	}
	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of warehouse and address")
		}

		copiedWarehouse := wareHouse.DeepCopy()

		if option.SelectRelatedAddress {
			copiedWarehouse.Address = address.DeepCopy()
		}

		returningWarehouses = append(returningWarehouses, copiedWarehouse)
		warehousesMap[wareHouse.Id] = copiedWarehouse
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed closing rows of warehouses and addresses")
	}

	// check if we need prefetch related shipping zones:
	if option.PrefetchShippingZones && len(returningWarehouses) > 0 {
		query, args, err = wh.GetQueryBuilder().
			Select(wh.ShippingZone().ModelFields(store.ShippingZoneTableName+".")...).
			Column(squirrel.Alias(squirrel.Expr("WarehouseShippingZones.WarehouseID"), "PrefetchRelatedWarehouseID")). // <- this column selection helps determine which shipping zone is related to which warehouse
			From(store.ShippingZoneTableName).
			InnerJoin(store.WarehouseShippingZoneTableName+" ON (ShippingZones.Id = WarehouseShippingZones.ShippingZoneID)").
			Where("WarehouseShippingZones.WarehouseID IN ?", returningWarehouses.IDs()).
			ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "FilerByOption_ToSql")
		}

		rows, err := wh.GetReplicaX().QueryX(query, args...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find shipping zones of warehouses")
		}
		var (
			shippingZone shipping.ShippingZone
			warehouseID  string
			scanFields   = append(wh.ShippingZone().ScanFields(shippingZone), &warehouseID)
		)

		for rows.Next() {
			err = rows.Scan(scanFields...)
			if err != nil {
				return nil, errors.Wrap(err, "failed to scan a row of shipping zone and warehouse id")
			}

			if warehousesMap[warehouseID] != nil {
				warehousesMap[warehouseID].ShippingZones = append(warehousesMap[warehouseID].ShippingZones, shippingZone.DeepCopy())
			}
		}

		if err = rows.Close(); err != nil {
			return nil, errors.Wrap(err, "failed closing rows of shipping zones")
		}
	}

	return returningWarehouses, nil
}

// WarehouseByStockID returns 1 warehouse by given stock id
func (ws *SqlWareHouseStore) WarehouseByStockID(stockID string) (*warehouse.WareHouse, error) {
	var res warehouse.WareHouse
	err := ws.GetReplicaX().Get(
		&res,
		`SELECT `+ws.ModelFields(store.WarehouseTableName+".").Join(",")+`
		FROM `+store.StockTableName+`
		INNER JOIN `+store.WarehouseTableName+` ON (
			Stocks.WarehouseID = Warehouses.Id
		)
		WHERE Stocks.Id = ?`,
		stockID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WarehouseTableName, "StockID="+stockID)
		}
		return nil, errors.Wrapf(err, "failed to find warehouse with StockID=%s", stockID)
	}

	return &res, nil
}

func (ws *SqlWareHouseStore) ApplicableForClickAndCollectNoQuantityCheck(checkoutLines checkout.CheckoutLines, country string) (warehouse.Warehouses, error) {
	stocks, err := ws.Stock().FilterByOption(nil, &warehouse.StockFilterOption{
		ProductVariantID:            squirrel.Eq{store.StockTableName + ".ProductVariantID": checkoutLines.VariantIDs()},
		SelectRelatedProductVariant: true,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks")
	}

	return ws.forCountryLinesAndStocks(checkoutLines, stocks, country)
}

func (ws *SqlWareHouseStore) ApplicableForClickAndCollect(checkoutLines checkout.CheckoutLines, country string) (warehouse.Warehouses, error) {
	panic("not implemented")
}

func (ws *SqlWareHouseStore) forCountryLinesAndStocks(checkoutLines checkout.CheckoutLines, stocks warehouse.Stocks, country string) (warehouse.Warehouses, error) {
	panic("not implemented")
}
