package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductVariantTranslationStore struct {
	store.Store
}

func NewSqlProductVariantTranslationStore(s store.Store) store.ProductVariantTranslationStore {
	pvts := &SqlProductVariantTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductVariantTranslation{}, "ProductVariantTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_VARIANT_NAME_MAX_LENGTH).SetUnique(true)

		table.SetUniqueTogether("LanguageCode", "ProductVariantID")
	}
	return pvts
}

func (ps *SqlProductVariantTranslationStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_product_variant_translations_name", "ProductVariantTranslations", "Name")
	ps.CreateIndexIfNotExists("idx_product_variant_translations_name_lower_textpattern", "ProductVariantTranslations", "lower(Name) text_pattern_ops")

}