package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductTranslationStore struct {
	store.Store
}

func NewSqlProductTranslationStore(s store.Store) store.ProductTranslationStore {
	pts := &SqlProductTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductTranslation{}, "ProductTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_NAME_MAX_LENGTH).SetUnique(true)

		table.SetUniqueTogether("LanguageCode", "ProductID")
		s.CommonSeoMaxLength(table)
	}
	return pts
}

func (ps *SqlProductTranslationStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_product_translations_name", "ProductTranslations", "Name")
	ps.CreateIndexIfNotExists("idx_product_translations_name_lower_textpattern", "ProductTranslations", "lower(Name) text_pattern_ops")
}
