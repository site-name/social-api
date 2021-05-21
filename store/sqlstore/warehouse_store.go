package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlWareHouseStore struct {
	*SqlStore
}

func newSqlWareHouseStore(s *SqlStore) store.WarehouseStore {
	ws := &SqlWareHouseStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.WareHouse{}, "WareHouses").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AddressID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(warehouse.WAREHOUSE_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(warehouse.WAREHOUSE_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("CompanyName").SetMaxSize(warehouse.WAREHOUSE_COMPANY_NAME_MAX_LENGTH)
		table.ColMap("Email").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
	}
	return ws
}

func (ws *SqlWareHouseStore) createIndexesIfNotExists() {
	ws.CreateIndexIfNotExists("idx_warehouses_name", "WareHouses", "Name")
	ws.CreateIndexIfNotExists("idx_warehouses_name_lower_textpattern", "WareHouses", "lower(Name) text_pattern_ops")
	ws.CreateIndexIfNotExists("idx_warehouses_slug", "WareHouses", "Slug")
	ws.CreateIndexIfNotExists("idx_warehouses_email", "WareHouses", "Email")
	ws.CreateIndexIfNotExists("idx_warehouses_email_lower_textpattern", "WareHouses", "lower(Email) text_pattern_ops")
}
