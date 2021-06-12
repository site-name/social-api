package discount

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherChannelListingStore struct {
	store.Store
}

func NewSqlVoucherChannelListingStore(sqlStore store.Store) store.VoucherChannelListingStore {
	vcls := &SqlVoucherChannelListingStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherChannelListing{}, "VoucherChannelListings").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("VoucherID", "ChannelID")
	}

	return vcls
}

func (vcls *SqlVoucherChannelListingStore) CreateIndexesIfNotExists() {

}
