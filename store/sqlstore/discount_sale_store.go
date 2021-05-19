package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDiscountSaleStore struct {
	*SqlStore
}

func newSqlDiscountSaleStore(sqlStore *SqlStore) store.DiscountSaleStore {
	ss := &SqlDiscountSaleStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Sale{}, "Sales").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.SALE_NAME_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(10).SetDefaultConstraint(model.NewString(product_and_discount.FIXED))

	}
	return ss
}

func (ss *SqlDiscountSaleStore) createIndexesIfNotExists() {
	ss.CreateIndexIfNotExists("idx_sales_name", "Sales", "Name")
	ss.CreateIndexIfNotExists("idx_sales_type", "Sales", "Type")
}
