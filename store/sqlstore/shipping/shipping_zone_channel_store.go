package shipping

import (
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore/channel"
)

const (
	ShippingZoneChannelTableName = "ShippingZoneChannels"
)

type SqlShippingZoneChannelStore struct {
	store.Store
}

func NewSqlShippingZoneChannelStore(s store.Store) store.ShippingZoneChannelStore {
	ss := &SqlShippingZoneChannelStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingZoneChannel{}, ShippingZoneChannelTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingZoneID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ShippingZoneID", "ChannelID")
	}

	return ss
}

func (ss *SqlShippingZoneChannelStore) CreateIndexesIfNotExists() {
	ss.CreateForeignKeyIfNotExists(ShippingZoneChannelTableName, "ShippingZoneID", ShippingZoneTableName, "Id", false)
	ss.CreateForeignKeyIfNotExists(ShippingZoneChannelTableName, "ChannelID", channel.ChannelTableName, "Id", false)
}
