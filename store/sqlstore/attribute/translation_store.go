package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeTranslationStore struct {
	store.Store
}

func NewSqlAttributeTranslationStore(s store.Store) store.AttributeTranslationStore {
	as := &SqlAttributeTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeTranslation{}, store.AttributeTranslationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(attribute.ATTRIBUTE_TRANSLATION_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "AttributeID")
	}
	return as
}

func (as *SqlAttributeTranslationStore) CreateIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_attributetranslations_name", store.AttributeTranslationTableName, "Name")
	as.CreateIndexIfNotExists("idx_attributetranslations_name_lower_textpattern", store.AttributeTranslationTableName, "lower(Name) text_pattern_ops")
}
