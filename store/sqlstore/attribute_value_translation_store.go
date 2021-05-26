package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeValueTranslationStore struct {
	*SqlStore
}

func newSqlAttributeValueTranslationStore(s *SqlStore) store.AttributeValueTranslationStore {
	as := &SqlAttributeValueTranslationStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeValueTranslation{}, "AttributeValueTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AttributeValueID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(attribute.ATTRIBUTE_VALUE_TRANSLATION_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "AttributeValueID")
	}
	return as
}

func (as *SqlAttributeValueTranslationStore) createIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_attribute_value_translations_name", "AttributeValueTranslations", "Name")
	as.CreateIndexIfNotExists("idx_attribute_value_translations_name_lower_textpattern", "AttributeValueTranslations", "lower(Name) text_pattern_ops")
}
