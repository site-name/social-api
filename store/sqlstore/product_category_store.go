package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCategoryStore struct {
	*SqlStore
}

func newSqlCategoryStore(s *SqlStore) store.CategoryStore {
	cs := &SqlCategoryStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Category{}, "Categories").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ParentID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.CATEGORY_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(product_and_discount.CATEGORY_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("BackgroundImage").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("BackgroundImageAlt").SetMaxSize(product_and_discount.COLLECTION_BACKGROUND_ALT_MAX_LENGTH)

		s.commonSeoMaxLength(table)
	}
	return cs
}

func (ps *SqlCategoryStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_categories_name", "Categories", "Name")
	ps.CreateIndexIfNotExists("idx_categories_slug", "Categories", "Slug")
	ps.CreateIndexIfNotExists("idx_categories_name_lower_textpattern", "Users", "lower(Name) text_pattern_ops")
}
