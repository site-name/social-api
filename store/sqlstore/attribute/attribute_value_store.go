package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeValueStore struct {
	store.Store
}

func NewSqlAttributeValueStore(s store.Store) store.AttributeValueStore {
	as := &SqlAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeValue{}, "AttributeValues").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(attribute.ATTRIBUTE_VALUE_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(attribute.ATTRIBUTE_VALUE_SLUG_MAX_LENGTH)
		table.ColMap("Value").SetMaxSize(attribute.ATTRIBUTE_VALUE_VALUE_MAX_LENGTH)
		table.ColMap("ContentType").SetMaxSize(attribute.ATTRIBUTE_VALUE_CONTENT_TYPE_MAX_LENGTH)
		table.ColMap("FileUrl").SetMaxSize(model.URL_LINK_MAX_LENGTH)

		table.SetUniqueTogether("Slug", "AttributeID")
	}
	return as
}

func (as *SqlAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_attributevalues_name", "AttributeValues", "Name")
	as.CreateIndexIfNotExists("idx_attributevalues_name_lower_textpattern", "AttributeValues", "lower(Name) text_pattern_ops")
	as.CreateIndexIfNotExists("idx_attributevalues_slug", "AttributeValues", "Slug")
	as.CreateIndexIfNotExists("idx_attributevalues_value", "AttributeValues", "Value")

	as.CreateForeignKeyIfNotExists("AttributeValues", "AttributeID", "Attributes", "Id", true)
}
