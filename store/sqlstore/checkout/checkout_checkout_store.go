package checkout

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

type SqlCheckoutStore struct {
	store.Store
}

func NewSqlCheckoutStore(sqlStore store.Store) store.CheckoutStore {
	cs := &SqlCheckoutStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(checkout.Checkout{}, "Checkouts").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("BillingAddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingAddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("DiscountName").SetMaxSize(checkout.CHECKOUT_DISCOUNT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedDiscountName").SetMaxSize(checkout.CHECKOUT_TRANSLATED_DISCOUNT_NAME_MAX_LENGTH)
		table.ColMap("VoucherCode").SetMaxSize(checkout.CHECKOUT_VOUCHER_CODE_MAX_LENGTH)
		table.ColMap("TrackingCode").SetMaxSize(checkout.CHECKOUT_TRACKING_CODE_MAX_LENGTH)
		table.ColMap("Country").SetMaxSize(model.SINGLE_COUNTRY_CODE_MAX_LENGTH)
	}

	return cs
}

func (cs *SqlCheckoutStore) CreateIndexesIfNotExists() {
	cs.CreateIndexIfNotExists("idx_checkouts_userid", "Checkouts", "UserID")
	cs.CreateIndexIfNotExists("idx_checkouts_token", "Checkouts", "Token")
	cs.CreateIndexIfNotExists("idx_checkouts_channelid", "Checkouts", "ChannelID")
	cs.CreateIndexIfNotExists("idx_checkouts_billing_address_id", "Checkouts", "BillingAddressID")
	cs.CreateIndexIfNotExists("idx_checkouts_shipping_address_id", "Checkouts", "ShippingAddressID")
	cs.CreateIndexIfNotExists("idx_checkouts_shipping_method_id", "Checkouts", "ShippingMethodID")

	cs.CommonMetaDataIndex("Checkouts")
}
