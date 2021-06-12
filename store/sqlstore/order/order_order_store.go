package order

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

type SqlOrderStore struct {
	store.Store
}

func NewSqlOrderStore(sqlStore store.Store) store.OrderStore {
	os := &SqlOrderStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(order.Order{}, "Orders").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("BillingAddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingAddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OriginalID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Status").SetMaxSize(order.ORDER_STATUS_MAX_LENGTH)
		table.ColMap("TrackingClientID").SetMaxSize(order.ORDER_TRACKING_CLIENT_ID_MAX_LENGTH)
		table.ColMap("Origin").SetMaxSize(order.ORDER_ORIGIN_MAX_LENGTH)
		table.ColMap("ShippingMethodName").SetMaxSize(order.ORDER_SHIPPING_METHOD_NAME_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(order.ORDER_TOKEN_MAX_LENGTH).SetUnique(true)
		table.ColMap("CheckoutToken").SetMaxSize(order.ORDER_CHECKOUT_TOKEN_MAX_LENGTH)
		table.ColMap("UserEmail").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(35).SetDefaultConstraint(model.NewString(model.DEFAULT_LANGUAGE_CODE))
		table.ColMap("Currency").SetMaxSize(model.URL_LINK_MAX_LENGTH)
	}

	return os
}

func (os *SqlOrderStore) CreateIndexesIfNotExists() {
	os.CommonMetaDataIndex("Orders")
	os.CreateIndexIfNotExists("idx_orders_user_email", "Orders", "UserEmail")
	os.CreateIndexIfNotExists("idx_orders_status", "Orders", "Status")
}
