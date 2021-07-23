package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeValueTranslationStore struct {
	store.Store
}

func NewSqlAttributeValueTranslationStore(s store.Store) store.AttributeValueTranslationStore {
	as := &SqlAttributeValueTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeValueTranslation{}, store.AttributeValueTranslationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AttributeValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(attribute.ATTRIBUTE_VALUE_TRANSLATION_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "AttributeValueID")
	}
	return as
}

func (as *SqlAttributeValueTranslationStore) CreateIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_attribute_value_translations_name", store.AttributeValueTranslationTableName, "Name")
	as.CreateIndexIfNotExists("idx_attribute_value_translations_name_lower_textpattern", store.AttributeValueTranslationTableName, "lower(Name) text_pattern_ops")
}
