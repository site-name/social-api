package sqlstore

import (
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherChannelListingStore struct {
	*SqlStore
}

func newSqlVoucherChannelListingStore(sqlStore *SqlStore) store.VoucherChannelListingStore {
	vcls := &SqlVoucherChannelListingStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherChannelListing{}, "VoucherChannelListings").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(UUID_MAX_LENGTH).SetNotNull(true)

	}

	return vcls
}

func (vcls *SqlVoucherChannelListingStore) createIndexesIfNotExists() {

}
