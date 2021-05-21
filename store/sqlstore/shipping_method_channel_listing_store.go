package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodChannelListingStore struct {
	*SqlStore
}

func newSqlShippingMethodChannelListingStore(s *SqlStore) store.ShippingMethodChannelListingStore {
	smls := &SqlShippingMethodChannelListingStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethodChannelListing{}, "ShippingMethodChannelListings").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("ShippingMethodID", "ChannelID")
	}
	return smls
}

func (s *SqlShippingMethodChannelListingStore) createIndexesIfNotExists() {

}
