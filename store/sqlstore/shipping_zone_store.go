package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingZoneStore struct {
	*SqlStore
}

func newSqlShippingZoneStore(s *SqlStore) store.ShippingZoneStore {
	smls := &SqlShippingZoneStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingZone{}, "ShippingZones").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(shipping.SHIPPING_ZONE_NAME_MAX_LENGTH)
		table.ColMap("Contries").SetMaxSize(model.MULTIPLE_COUNTRIES_MAX_LENGTH)
	}
	return smls
}

func (s *SqlShippingZoneStore) createIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_shipping_zone_name", "ShippingZones", "Name")
	s.CreateIndexIfNotExists("idx_shipping_zone_name_lower_textpattern", "ShippingZones", "lower(Name) text_pattern_ops")
}
