package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductVariantChannelListingStore struct {
	*SqlStore
}

func newSqlProductVariantChannelListingStore(s *SqlStore) store.ProductVariantChannelListingStore {
	pvcls := &SqlProductVariantChannelListingStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductVariantChannelListing{}, "ProductVariantChannelListings").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("ChannelID").SetMaxSize(UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("VariantID", "ChannelID")
	}
	return pvcls
}

func (ps *SqlProductVariantChannelListingStore) createIndexesIfNotExists() {

}
