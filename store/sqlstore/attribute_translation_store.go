package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeTranslationStore struct {
	*SqlStore
}

func newSqlAttributeTranslationStore(s *SqlStore) store.AttributeTranslationStore {
	as := &SqlAttributeTranslationStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeTranslation{}, "AttributeTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(attribute.ATTRIBUTE_TRANSLATION_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "AttributeID")
	}
	return as
}

func (as *SqlAttributeTranslationStore) createIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_attributetranslations_name", "AttributeTranslations", "Name")
	as.CreateIndexIfNotExists("idx_attributetranslations_name_lower_textpattern", "AttributeTranslations", "lower(Name) text_pattern_ops")
}
