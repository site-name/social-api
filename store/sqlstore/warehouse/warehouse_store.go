package warehouse

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlWareHouseStore struct {
	store.Store
}

func NewSqlWareHouseStore(s store.Store) store.WarehouseStore {
	ws := &SqlWareHouseStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.WareHouse{}, "WareHouses").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(warehouse.WAREHOUSE_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(warehouse.WAREHOUSE_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("Email").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
	}
	return ws
}

func (ws *SqlWareHouseStore) CreateIndexesIfNotExists() {
	ws.CreateIndexIfNotExists("idx_warehouses_name", "WareHouses", "Name")
	ws.CreateIndexIfNotExists("idx_warehouses_name_lower_textpattern", "WareHouses", "lower(Name) text_pattern_ops")
	ws.CreateIndexIfNotExists("idx_warehouses_slug", "WareHouses", "Slug")
	ws.CreateIndexIfNotExists("idx_warehouses_email", "WareHouses", "Email")
	ws.CreateIndexIfNotExists("idx_warehouses_email_lower_textpattern", "WareHouses", "lower(Email) text_pattern_ops")
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
	inter, err := ws.GetMaster().Get(warehouse.WareHouse{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Warehouse", id)
		}
		return nil, errors.Wrapf(err, "failed to get warehouse with Id=%s", id)
	}

	return inter.(*warehouse.WareHouse), nil
}

func (wh *SqlWareHouseStore) GetWarehousesHeaders(ids []string) ([]string, error) {
	var headers []string
	_, err := wh.GetReplica().Select(
		&headers,
		`SELECT
			CONCAT(wh.Slug, ' (warehouse quantity)') AS header
		FROM 
			WareHouses AS wh
		WHERE 
			wh.Id IN :IDS
		ORDER BY wh.Slug`,
		map[string]interface{}{"IDS": ids},
	)

	if err != nil {
		return nil, err
	}

	return headers, nil
}
