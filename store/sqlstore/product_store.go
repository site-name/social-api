package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductStore struct {
	*SqlStore
}

func newSqlProductStore(s *SqlStore) store.ProductStore {
	ps := &SqlProductStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Product{}, "Products").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductTypeID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("DefaultVariantID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("CategoryID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(product_and_discount.PRODUCT_SLUG_MAX_LENGTH).SetUnique(true)

		s.commonSeoMaxLength(table)
	}
	return ps
}

func (ps *SqlProductStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_products_name", "Products", "Name")
	ps.CreateIndexIfNotExists("idx_products_slug", "Products", "Slug")
	ps.CreateIndexIfNotExists("idx_products_name_lower_textpattern", "Products", "lower(Name) text_pattern_ops")

	ps.CommonMetaDataIndex("Products")
}
