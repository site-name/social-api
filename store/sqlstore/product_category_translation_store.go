package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCategoryTranslationStore struct {
	*SqlStore
}

func newSqlCategoryTranslationStore(s *SqlStore) store.CategoryTranslationStore {
	cts := &SqlCategoryTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.CategoryTranslation{}, "CategoryTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("CategoryID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.CATEGORY_NAME_MAX_LENGTH)

		s.commonSeoMaxLength(table)
	}
	return cts
}

func (ps *SqlCategoryTranslationStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_category_translations_name", "CategoryTranslations", "Name")
	ps.CreateIndexIfNotExists("idx_category_translations_name_lower_textpattern", "CategoryTranslations", "lower(Name) text_pattern_ops")
}
