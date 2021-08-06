package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCollectionStore struct {
	store.Store
}

func NewSqlCollectionStore(s store.Store) store.CollectionStore {
	cs := &SqlCollectionStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Collection{}, store.ProductCollectionTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.COLLECTION_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(product_and_discount.COLLECTION_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("BackgroundImage").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("BackgroundImageAlt").SetMaxSize(product_and_discount.COLLECTION_BACKGROUND_ALT_MAX_LENGTH)

		s.CommonSeoMaxLength(table)
	}
	return cs
}

func (ps *SqlCollectionStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_collections_name", store.ProductCollectionTableName, "Name")
	ps.CreateIndexIfNotExists("idx_collections_slug", store.ProductCollectionTableName, "Slug")
	ps.CreateIndexIfNotExists("idx_collections_name_lower_textpattern", store.ProductCollectionTableName, "lower(Name) text_pattern_ops")
}

func (ps *SqlCollectionStore) ModelFields() []string {
	return []string{
		"Collections.Id",
		"Collections.Name",
		"Collections.Slug",
		"Collections.BackgroundImage",
		"Collections.BackgroundImageAlt",
		"Collections.Description",
		"Collections.Metadata",
		"Collections.PrivateMetadata",
		"Collections.SeoTitle",
		"Collections.SeoDescription",
	}
}
