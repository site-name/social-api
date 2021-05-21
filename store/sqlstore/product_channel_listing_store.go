package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductChannelListingStore struct {
	*SqlStore
}

func newSqlProductChannelListingStore(s *SqlStore) store.ProductChannelListingStore {
	pcls := &SqlProductChannelListingStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductChannelListing{}, "ProductChannelListings").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

	}
	return pcls
}

func (ps *SqlProductChannelListingStore) createIndexesIfNotExists() {

}
