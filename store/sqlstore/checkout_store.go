package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

type SqlCheckoutStore struct {
	*SqlStore
}

func newSqlCheckoutStore(sqlStore *SqlStore) store.CheckoutStore {
	cs := &SqlCheckoutStore{
		SqlStore: sqlStore,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(checkout.Checkout{}, "Checkouts").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("UserID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("BillingAddressID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ShippingAddressID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("DiscountName").SetMaxSize(checkout.CHECKOUT_DISCOUNT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedDiscountName").SetMaxSize(checkout.CHECKOUT_TRANSLATED_DISCOUNT_NAME_MAX_LENGTH)
		table.ColMap("VoucherCode").SetMaxSize(checkout.CHECKOUT_VOUCHER_CODE_MAX_LENGTH)
		table.ColMap("TrackingCode").SetMaxSize(checkout.CHECKOUT_TRACKING_CODE_MAX_LENGTH)
		table.ColMap("Country").SetMaxSize(model.SingleCountryMaxLength)
	}

	return cs
}
