package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCollectionTranslationStore struct {
	store.Store
}

func NewSqlCollectionTranslationStore(s store.Store) store.CollectionTranslationStore {
	cts := &SqlCollectionTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.CollectionTranslation{}, "CollectionTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CollectionID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.COLLECTION_NAME_MAX_LENGTH).SetUnique(true)

		s.CommonSeoMaxLength(table)
		table.SetUniqueTogether("LanguageCode", "CollectionID")
	}
	return cts
}

func (ps *SqlCollectionTranslationStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_collection_translations_name", "CollectionTranslations", "Name")
	ps.CreateIndexIfNotExists("idx_collections_translations_name_lower_textpattern", "CollectionTranslations", "lower(Name) text_pattern_ops")
}
