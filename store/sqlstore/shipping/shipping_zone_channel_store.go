package shipping

import (
	"github.com/sitename/sitename/store"
)

type SqlShippingZoneChannelStore struct {
	store.Store
}

func NewSqlShippingZoneChannelStore(s store.Store) store.ShippingZoneChannelStore {
	return &SqlShippingZoneChannelStore{s}
}
