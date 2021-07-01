package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

const (
	CategoryTableName = "Categories"
)

type SqlCategoryStore struct {
	store.Store
}

func NewSqlCategoryStore(s store.Store) store.CategoryStore {
	cs := &SqlCategoryStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Category{}, CategoryTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ParentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.CATEGORY_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(product_and_discount.CATEGORY_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("BackgroundImage").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("BackgroundImageAlt").SetMaxSize(product_and_discount.COLLECTION_BACKGROUND_ALT_MAX_LENGTH)

		s.CommonSeoMaxLength(table)
	}
	return cs
}

func (ps *SqlCategoryStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_categories_name", CategoryTableName, "Name")
	ps.CreateIndexIfNotExists("idx_categories_slug", CategoryTableName, "Slug")
	ps.CreateIndexIfNotExists("idx_categories_name_lower_textpattern", CategoryTableName, "lower(Name) text_pattern_ops")
}
