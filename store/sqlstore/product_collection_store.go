package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCollectionStore struct {
	*SqlStore
}

func newSqlCollectionStore(s *SqlStore) store.CollectionStore {
	cs := &SqlCollectionStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Collection{}, "Collections").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.COLLECTION_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(product_and_discount.COLLECTION_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("BackgroundImage").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("BackgroundImageAlt").SetMaxSize(product_and_discount.COLLECTION_BACKGROUND_ALT_MAX_LENGTH)

		s.commonSeoMaxLength(table)
	}
	return cs
}

func (ps *SqlCollectionStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_collections_name", "Collections", "Name")
	ps.CreateIndexIfNotExists("idx_collections_slug", "Collections", "Slug")
	ps.CreateIndexIfNotExists("idx_collections_name_lower_textpattern", "Collections", "lower(Name) text_pattern_ops")
}
