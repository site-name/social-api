package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductVariantChannelListingStore struct {
	store.Store
}

func NewSqlProductVariantChannelListingStore(s store.Store) store.ProductVariantChannelListingStore {
	pvcls := &SqlProductVariantChannelListingStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductVariantChannelListing{}, "ProductVariantChannelListings").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("VariantID", "ChannelID")
	}
	return pvcls
}

func (ps *SqlProductVariantChannelListingStore) CreateIndexesIfNotExists() {

}
