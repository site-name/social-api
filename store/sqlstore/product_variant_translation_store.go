package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductVariantTranslationStore struct {
	*SqlStore
}

func newSqlProductVariantTranslationStore(s *SqlStore) store.ProductVariantTranslationStore {
	pvts := &SqlProductVariantTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductVariantTranslation{}, "ProductVariantTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_VARIANT_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "ProductVariantID")
	}
	return pvts
}

func (ps *SqlProductVariantTranslationStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_product_variant_translations_name", "ProductVariantTranslations", "Name")
	ps.CreateIndexIfNotExists("idx_product_variant_translations_name_lower_textpattern", "ProductVariantTranslations", "lower(Name) text_pattern_ops")

}
