package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductVariantStore struct {
	*SqlStore
}

func newSqlProductVariantStore(s *SqlStore) store.ProductVariantStore {
	pvs := &SqlProductVariantStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductVariant{}, "ProductVariants").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Sku").SetMaxSize(product_and_discount.PRODUCT_VARIANT_SKU_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_VARIANT_NAME_MAX_LENGTH)

	}
	return pvs
}

func (ps *SqlProductVariantStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_product_variants_name", "ProductVariants", "Name")
	ps.CreateIndexIfNotExists("idx_product_variants_name_lower_textpattern", "ProductVariants", "lower(Name) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_product_variants_sku", "ProductVariants", "Sku")
}
