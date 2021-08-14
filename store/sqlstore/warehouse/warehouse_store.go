package warehouse

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
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
		"Warehouses.Metadata",
		"Warehouses.PrivateMetadata",
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

// FilterByOprion returns a slice of warehouses with given option
func (wh *SqlWareHouseStore) FilterByOprion(option *warehouse.WarehouseFilterOption) ([]*warehouse.WareHouse, error) {
	query := wh.GetQueryBuilder().
		Select(wh.ModelFields()...).
		Distinct().
		From(store.WarehouseTableName).
		OrderBy(store.TableOrderingMap[store.WarehouseTableName])

	// parse option
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
			InnerJoin(store.WarehouseShippingZoneTableName + " ON (Warehouses.Id = WarehouseShippingZones.WarehouseID)").
			InnerJoin(store.ShippingZoneTableName + " ON (WarehouseShippingZones.ShippingZoneID = ShippingZones.Id)").
			Where(option.ShippingZonesCountries.ToSquirrel("ShippingZones.Countries"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*warehouse.WareHouse
	_, err = wh.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find warehouses with given option")
	}

	return res, nil
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
