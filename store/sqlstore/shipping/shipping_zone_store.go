package shipping

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingZoneStore struct {
	store.Store
}

func NewSqlShippingZoneStore(s store.Store) store.ShippingZoneStore {
	smls := &SqlShippingZoneStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingZone{}, store.ShippingZoneTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(shipping.SHIPPING_ZONE_NAME_MAX_LENGTH)
		table.ColMap("Contries").SetMaxSize(model.MULTIPLE_COUNTRIES_MAX_LENGTH)
	}
	return smls
}

func (s *SqlShippingZoneStore) ModelFields() []string {
	return []string{
		"ShippingZones.Id",
		"ShippingZones.Name",
		"ShippingZones.Contries",
		"ShippingZones.Default",
		"ShippingZones.Description",
		"ShippingZones.Metadata",
		"ShippingZones.PrivateMetadata",
	}
}

func (s *SqlShippingZoneStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_shipping_zone_name", store.ShippingZoneTableName, "Name")
	s.CreateIndexIfNotExists("idx_shipping_zone_name_lower_textpattern", store.ShippingZoneTableName, "lower(Name) text_pattern_ops")
}
