package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductTranslationStore struct {
	*SqlStore
}

func newSqlProductTranslationStore(s *SqlStore) store.ProductTranslationStore {
	pts := &SqlProductTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductTranslation{}, "ProductTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "ProductID")
		s.commonSeoMaxLength(table)
	}
	return pts
}

func (ps *SqlProductTranslationStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_product_translations_name", "ProductTranslations", "Name")
	ps.CreateIndexIfNotExists("idx_product_translations_name_lower_textpattern", "ProductTranslations", "lower(Name) text_pattern_ops")
}
