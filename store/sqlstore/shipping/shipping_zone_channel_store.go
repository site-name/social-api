package shipping

import (
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingZoneChannelStore struct {
	store.Store
}

func NewSqlShippingZoneChannelStore(s store.Store) store.ShippingZoneChannelStore {
	ss := &SqlShippingZoneChannelStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingZoneChannel{}, store.ShippingZoneChannelTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingZoneID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ShippingZoneID", "ChannelID")
	}

	return ss
}

func (ss *SqlShippingZoneChannelStore) CreateIndexesIfNotExists() {
	ss.CreateForeignKeyIfNotExists(store.ShippingZoneChannelTableName, "ShippingZoneID", store.ShippingZoneTableName, "Id", false)
	ss.CreateForeignKeyIfNotExists(store.ShippingZoneChannelTableName, "ChannelID", store.ChannelTableName, "Id", false)
}
