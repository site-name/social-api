package warehouse

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlWareHouseStore struct {
	store.Store
}

func NewSqlWarehouseStore(s store.Store) store.WarehouseStore {
	ws := &SqlWareHouseStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.WareHouse{}, store.WarehouseTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(warehouse.WAREHOUSE_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(warehouse.WAREHOUSE_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("Email").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
		table.ColMap("ClickAndCollectOption").SetMaxSize(warehouse.WAREHOUSE_CLICK_AND_COLLECT_OPTION_MAX_LENGTH)
	}
	return ws
}

func (ws *SqlWareHouseStore) ModelFields() []string {
	return []string{
		"Warehouses.Id",
		"Warehouses.Name",
		"Warehouses.Slug",
		"Warehouses.AddressID",
		"Warehouses.Email",
		"Warehouses.ClickAndCollectOption",
		"Warehouses.IsPrivate",
		"Warehouses.Metadata",
		"Warehouses.PrivateMetadata",
	}
}

func (ws *SqlWareHouseStore) ScanFields(wh warehouse.WareHouse) []interface{} {
	return []interface{}{
		&wh.Id,
		&wh.Name,
		&wh.Slug,
		&wh.AddressID,
		&wh.Email,
		&wh.ClickAndCollectOption,
		&wh.IsPrivate,
		&wh.Metadata,
		&wh.PrivateMetadata,
	}
}

func (ws *SqlWareHouseStore) CreateIndexesIfNotExists() {
	ws.CreateIndexIfNotExists("idx_warehouses_email", store.WarehouseTableName, "Email")
	ws.CreateIndexIfNotExists("idx_warehouses_email_lower_textpattern", store.WarehouseTableName, "lower(Email) text_pattern_ops")

	ws.CreateForeignKeyIfNotExists(store.WarehouseTableName, "AddressID", store.AddressTableName, "Id", false)
}

func (ws *SqlWareHouseStore) Save(wh *warehouse.WareHouse) (*warehouse.WareHouse, error) {
	wh.PreSave()
	if err := wh.IsValid(); err != nil {
		return nil, err
	}

	if err := ws.GetMaster().Insert(wh); err != nil {
		if ws.IsUniqueConstraintError(err, []string{"Slug", "warehouses_slug_key", "idx_warehouses_slug_unique"}) {
			return nil, store.NewErrInvalidInput("Warehouses", "Slug", wh.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save Warehouse with Id=%s", wh.Id)
	}

	return wh, nil
}

func (ws *SqlWareHouseStore) Get(id string) (*warehouse.WareHouse, error) {
	var res warehouse.WareHouse
	err := ws.GetMaster().SelectOne(
		&res,
		"SELECT * FROM "+store.WarehouseTableName+" WHERE Id = :ID",
		map[string]interface{}{
			"ID": id,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Warehouse", id)
		}
		return nil, errors.Wrapf(err, "failed to get warehouse with Id=%s", id)
	}

	return &res, nil
}

func (wh *SqlWareHouseStore) commonQueryBuilder(option *warehouse.WarehouseFilterOption) squirrel.SelectBuilder {
	if option == nil {
		return wh.GetQueryBuilder().
			Select(wh.ModelFields()...).
			From(store.WarehouseTableName).
			OrderBy(store.TableOrderingMap[store.WarehouseTableName])
	}

	selectFields := wh.ModelFields()
	if option.SelectRelatedAddress {
		selectFields = append(selectFields, wh.Address().ModelFields()...)
	}

	query := wh.GetQueryBuilder().
		Select(selectFields...).
		From(store.WarehouseTableName + " WH").
		OrderBy(store.TableOrderingMap[store.WarehouseTableName])

	// parse option
	if option.Distinct {
		query = query.Distinct()
	}
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.Name != nil {
		query = query.Where(option.Name.ToSquirrel("Name"))
	}
	if option.Slug != nil {
		query = query.Where(option.Slug.ToSquirrel("Slug"))
	}
	if option.AddressID != nil {
		query = query.Where(option.AddressID.ToSquirrel("AddressID"))
	}
	if option.Email != nil {
		query = query.Where(option.Email.ToSquirrel("Email"))
	}
	if option.ShippingZonesCountries != nil {
		query = query.
			InnerJoin(store.WarehouseShippingZoneTableName + " AS WhSz ON (WH.Id = WhSz.WarehouseID)").
			InnerJoin(store.ShippingZoneTableName + " AS SZ ON (WhSz.ShippingZoneID = SZ.Id)").
			Where(option.ShippingZonesCountries.ToSquirrel("SZ.Countries"))
	}
	// check if we need to join address table:
	if option.SelectRelatedAddress {
		query.InnerJoin(store.AddressTableName + " AS AD ON (AD.Id = WH.AddressID)")
	}

	return query
}

// GetByOption finds and returns a warehouse filtered given option
func (wh *SqlWareHouseStore) GetByOption(option *warehouse.WarehouseFilterOption) (*warehouse.WareHouse, error) {
	rowScanner := wh.commonQueryBuilder(option).RunWith(wh.GetReplica()).QueryRow()

	var (
		res     warehouse.WareHouse
		address account.Address
	)
	scanFields := wh.ScanFields(res)
	if option.SelectRelatedAddress {
		scanFields = append(scanFields, wh.Address().ScanFields(address)...)
	}

	err := rowScanner.Scan(scanFields...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WarehouseTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find warehouse with given option")
	}

	if option.SelectRelatedAddress {
		res.Address = &address
	}

	// check if we need to prefetch shipping zones
	if option.PrefetchShippingZones && model.IsValidId(res.Id) {
		rows, err := wh.GetQueryBuilder().
			Select(wh.ShippingZone().ModelFields()...).
			From(store.ShippingZoneTableName+" SZ").
			InnerJoin(store.WarehouseShippingZoneTableName+" WHSZ ON (SZ.Id = WHSZ.ShippingZoneID)").
			Where("WHSZ.WarehouseID = ?", res.Id).
			RunWith(wh.GetReplica()).Query()

		if err != nil {
			return nil, errors.Wrap(err, "failed to find shipping zones related to returning warehouse")
		}
		var shippingZone shipping.ShippingZone
		scanFields := wh.ShippingZone().ScanFields(shippingZone)

		for rows.Next() {
			err = rows.Scan(scanFields...)
			if err != nil {
				return nil, errors.Wrap(err, "failed to scan a row of shipping zone")
			}

			res.ShippingZones = append(res.ShippingZones, &shippingZone)
		}

		if err = rows.Close(); err != nil {
			return nil, errors.Wrap(err, "failed to close rows of shipping zones")
		}
	}

	return &res, nil
}

// FilterByOprion returns a slice of warehouses with given option
func (wh *SqlWareHouseStore) FilterByOprion(option *warehouse.WarehouseFilterOption) ([]*warehouse.WareHouse, error) {

	rows, err := wh.commonQueryBuilder(option).RunWith(wh.GetReplica()).Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find warehouses with given option")
	}

	var (
		returningWarehouses warehouse.Warehouses
		warehousesMap       = map[string]*warehouse.WareHouse{} // keys are warehouse IDs
		wareHouse           warehouse.WareHouse
		address             account.Address
		scanItems           = wh.ScanFields(wareHouse)
	)
	if option.SelectRelatedAddress {
		scanItems = append(scanItems, wh.Address().ScanFields(address)...)
	}
	for rows.Next() {
		err = rows.Scan(scanItems...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of warehouse and address")
		}

		if option.SelectRelatedAddress {
			wareHouse.Address = &address
		}
		returningWarehouses = append(returningWarehouses, &wareHouse)
		warehousesMap[wareHouse.Id] = &wareHouse
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed closing rows of warehouses and addresses")
	}

	// check if we need prefetch related shipping zones:
	if option.PrefetchShippingZones && len(returningWarehouses) > 0 {
		rows, err = wh.GetQueryBuilder().
			Select(wh.ShippingZone().ModelFields()...).
			Column(squirrel.Alias(squirrel.Expr("WHSZ.WarehouseID"), "PrefetchRelatedWarehouseID")). // <- this column selection helps determine which shipping zone is related to which warehouse
			From(store.ShippingZoneTableName+" SZ").
			InnerJoin(store.WarehouseShippingZoneTableName+" WHSZ ON (SZ.Id = WHSZ.ShippingZoneID)").
			Where("WHSZ.WarehouseID IN ?", returningWarehouses.IDs()).
			RunWith(wh.GetReplica()).Query()

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
				warehousesMap[warehouseID].ShippingZones = append(warehousesMap[warehouseID].ShippingZones, &shippingZone)
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
	err := ws.GetReplica().SelectOne(
		&res,
		`SELECT `+strings.Join(ws.ModelFields(), ", ")+`
		FROM `+store.StockTableName+`
		INNER JOIN `+store.WarehouseTableName+` ON (
			Stocks.WarehouseID = Warehouses.Id
		)
		WHERE Stocks.Id = :StockID`,
		map[string]interface{}{
			"StockID": stockID,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.WarehouseTableName, "StockID="+stockID)
		}
		return nil, errors.Wrapf(err, "failed to find warehouse with StockID=%s", stockID)
	}

	return &res, nil
}
