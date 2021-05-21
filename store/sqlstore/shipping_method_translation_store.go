package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodTranslationStore struct {
	*SqlStore
}

func newSqlShippingMethodTranslationStore(s *SqlStore) store.ShippingMethodTranslationStore {
	smls := &SqlShippingMethodTranslationStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethodTranslation{}, "ShippingMethodTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(shipping.SHIPPING_METHOD_NAME_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "ShippingMethodID")
	}
	return smls
}

func (s *SqlShippingMethodTranslationStore) createIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_shipping_method_translations_name", "ShippingMethodTranslations", "Name")
	s.CreateIndexIfNotExists("idx_shipping_method_translations_name_lower_textpattern", "ShippingMethodTranslations", "lower(Name) text_pattern_ops")
}
