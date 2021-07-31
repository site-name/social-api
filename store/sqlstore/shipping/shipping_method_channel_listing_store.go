package shipping

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodChannelListingStore struct {
	store.Store
}

func NewSqlShippingMethodChannelListingStore(s store.Store) store.ShippingMethodChannelListingStore {
	smls := &SqlShippingMethodChannelListingStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethodChannelListing{}, store.ShippingMethodChannelListingTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("ShippingMethodID", "ChannelID")
	}
	return smls
}

func (s *SqlShippingMethodChannelListingStore) CreateIndexesIfNotExists() {

}
