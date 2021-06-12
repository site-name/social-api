package shipping

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodTranslationStore struct {
	store.Store
}

func NewSqlShippingMethodTranslationStore(s store.Store) store.ShippingMethodTranslationStore {
	smls := &SqlShippingMethodTranslationStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethodTranslation{}, "ShippingMethodTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(shipping.SHIPPING_METHOD_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "ShippingMethodID")
	}
	return smls
}

func (s *SqlShippingMethodTranslationStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_shipping_method_translations_name", "ShippingMethodTranslations", "Name")
	s.CreateIndexIfNotExists("idx_shipping_method_translations_name_lower_textpattern", "ShippingMethodTranslations", "lower(Name) text_pattern_ops")
}
