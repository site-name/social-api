package warehouse

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlWareHouseStore struct {
	store.Store
}

func NewSqlWarehouseStore(s store.Store) store.WarehouseStore {
	return &SqlWareHouseStore{s}
}

func (ws *SqlWareHouseStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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

func (ws *SqlWareHouseStore) Get(id string) (*model.WareHouse, error) {
	var res model.WareHouse
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

func (ws *SqlWareHouseStore) commonQueryBuilder(option *model.WarehouseFilterOption) squirrel.SelectBuilder {
	if option == nil {
		return ws.GetQueryBuilder().
			Select(ws.ModelFields(store.WarehouseTableName + ".")...).
			From(store.WarehouseTableName)
	}

	selectFields := ws.ModelFields(store.WarehouseTableName + ".")
	if option.SelectRelatedAddress {
		selectFields = append(selectFields, ws.Address().ModelFields(store.AddressTableName+".")...)
	}

	query := ws.GetQueryBuilder().
		Select(selectFields...).
		From(store.WarehouseTableName)

	for _, opt := range []squirrel.Sqlizer{
		option.Id,
		option.Name,
		option.Slug,
		option.AddressID,
		option.Email,
		option.ShippingZonesCountries,
		option.ShippingZonesId,
		option.IsPrivate,
		option.ClickAndCollectOption,
	} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	if option.ShippingZonesCountries != nil || option.ShippingZonesId != nil {
		query = query.
			InnerJoin(store.WarehouseShippingZoneTableName + " ON Warehouses.Id = WarehouseShippingZones.WarehouseID").
			InnerJoin(store.ShippingZoneTableName + " ON WarehouseShippingZones.ShippingZoneID = ShippingZones.Id")
	}
	if option.SelectRelatedAddress || option.Search != "" {
		query = query.InnerJoin(store.AddressTableName + " ON (Addresses.Id = Warehouses.AddressID)")

		if option.Search != "" {
			expr := "%" + option.Search + "%"

			query = query.Where(squirrel.Or{
				squirrel.ILike{store.WarehouseTableName + ".Name": expr},
				squirrel.ILike{store.WarehouseTableName + ".Email": expr},

				squirrel.ILike{store.AddressTableName + ".CompanyName": expr},
				squirrel.ILike{store.AddressTableName + ".StreetAddress1": expr},
				squirrel.ILike{store.AddressTableName + ".StreetAddress2": expr},
				squirrel.ILike{store.AddressTableName + ".City": expr},
				squirrel.ILike{store.AddressTableName + ".PostalCode": expr},
				squirrel.ILike{store.AddressTableName + ".Phone": expr},
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

	err = ws.GetReplicaX().QueryRowX(query, args...).Scan(scanFields...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WarehouseTableName, "options")
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
			Select(ws.ShippingZone().ModelFields(store.ShippingZoneTableName+".")...).
			From(store.ShippingZoneTableName).
			InnerJoin(store.WarehouseShippingZoneTableName+" ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID").
			Where("WarehouseShippingZones.WarehouseID = ?", res.Id).
			ToSql()

		if err != nil {
			return nil, errors.Wrap(err, "GetByOption_Warehouse_ToSql")
		}

		var shippingZones model.ShippingZones
		err = ws.GetReplicaX().Select(&shippingZones, queryString, args...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find shipping zones by warehouse ids")
		}

		res.SetShippingZones(shippingZones)
	}

	return &res, nil
}

// FilterByOprion returns a slice of warehouses with given option
func (wh *SqlWareHouseStore) FilterByOprion(option *model.WarehouseFilterOption) ([]*model.WareHouse, error) {
	query, args, err := wh.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}
	rows, err := wh.GetReplicaX().QueryX(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find warehouses with given option")
	}

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

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed closing rows of warehouses and addresses")
	}

	// check if we need prefetch related shipping zones:
	if option.PrefetchShippingZones && len(returningWarehouses) > 0 {
		query, args, err = wh.GetQueryBuilder().
			Select(wh.ShippingZone().ModelFields(store.ShippingZoneTableName + ".")...).
			Column("WarehouseShippingZones.WarehouseID AS PrefetchRelatedWarehouseID"). // <- this column selection helps determine which shipping zone is related to which warehouse
			From(store.ShippingZoneTableName).
			InnerJoin(store.WarehouseShippingZoneTableName + " ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID").
			Where(squirrel.Eq{"PrefetchRelatedWarehouseID": returningWarehouses.IDs()}).
			ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "FilerByOption_Prefetch_ToSql")
		}

		rows, err := wh.GetReplicaX().QueryX(query, args...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find shipping zones of warehouses")
		}
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
				warehousesMap[relatedWarehouseID].AppendShippingZone(&shippingZone)
			}
		}

		if err = rows.Close(); err != nil {
			return nil, errors.Wrap(err, "failed closing rows of shipping zones")
		}
	}

	return returningWarehouses, nil
}

// WarehouseByStockID returns 1 warehouse by given stock id
func (ws *SqlWareHouseStore) WarehouseByStockID(stockID string) (*model.WareHouse, error) {
	var res model.WareHouse
	err := ws.GetReplicaX().QueryRowX(
		`SELECT `+ws.ModelFields(store.WarehouseTableName+".").Join(",")+`
		FROM `+store.WarehouseTableName+`
		INNER JOIN `+store.StockTableName+` ON Stocks.WarehouseID = Warehouses.Id
		WHERE Stocks.Id = ?`,
		stockID,
	).
		Scan(ws.ScanFields(&res)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WarehouseTableName, "StockID="+stockID)
		}
		return nil, errors.Wrapf(err, "failed to find warehouse with StockID=%s", stockID)
	}

	return &res, nil
}

func (ws *SqlWareHouseStore) ApplicableForClickAndCollectNoQuantityCheck(checkoutLines model.CheckoutLines, country model.CountryCode) (model.Warehouses, error) {
	stocks, err := ws.Stock().FilterByOption(nil, &model.StockFilterOption{
		ProductVariantID:            squirrel.Eq{store.StockTableName + ".ProductVariantID": checkoutLines.VariantIDs()},
		SelectRelatedProductVariant: true,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks")
	}

	return ws.forCountryLinesAndStocks(checkoutLines, stocks, country)
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
