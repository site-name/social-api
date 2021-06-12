package shipping

import (
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodStore struct {
	store.Store
}

func NewSqlShippingMethodStore(s store.Store) store.ShippingMethodStore {
	smls := &SqlShippingMethodStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethod{}, "ShippingMethods").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingZoneID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(shipping.SHIPPING_METHOD_NAME_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(shipping.SHIPPING_METHOD_TYPE_MAX_LENGTH)
	}
	return smls
}

func (s *SqlShippingMethodStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_shipping_methods_name", "ShippingMethods", "Name")
	s.CreateIndexIfNotExists("idx_shipping_methods_name_lower_textpattern", "ShippingMethods", "lower(Name) text_pattern_ops")

}
