package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductTypeStore struct {
	*SqlStore
}

func newSqlProductTypeStore(s *SqlStore) store.ProductTypeStore {
	pts := &SqlProductTypeStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductType{}, "ProductTypes").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_TYPE_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(product_and_discount.PRODUCT_TYPE_SLUG_MAX_LENGTH)
	}
	return pts
}

func (ps *SqlProductTypeStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_product_types_name", "ProductTypes", "Name")
	ps.CreateIndexIfNotExists("idx_product_types_name_lower_textpattern", "ProductTypes", "lower(Name) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_product_types_slug", "ProductTypes", "Slug")
}
