package checkout

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore/account"
	"github.com/sitename/sitename/store/sqlstore/channel"
	"github.com/sitename/sitename/store/sqlstore/shipping"
)

type SqlCheckoutStore struct {
	store.Store
}

const (
	CheckoutTableName = "Checkouts"
)

func NewSqlCheckoutStore(sqlStore store.Store) store.CheckoutStore {
	cs := &SqlCheckoutStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(checkout.Checkout{}, CheckoutTableName).SetKeys(false, "Id")
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
	cs.CreateIndexIfNotExists("idx_checkouts_userid", CheckoutTableName, "UserID")
	cs.CreateIndexIfNotExists("idx_checkouts_token", CheckoutTableName, "Token")
	cs.CreateIndexIfNotExists("idx_checkouts_channelid", CheckoutTableName, "ChannelID")
	cs.CreateIndexIfNotExists("idx_checkouts_billing_address_id", CheckoutTableName, "BillingAddressID")
	cs.CreateIndexIfNotExists("idx_checkouts_shipping_address_id", CheckoutTableName, "ShippingAddressID")
	cs.CreateIndexIfNotExists("idx_checkouts_shipping_method_id", CheckoutTableName, "ShippingMethodID")

	cs.CreateForeignKeyIfNotExists(CheckoutTableName, "UserID", account.UserTableName, "Id", true)
	cs.CreateForeignKeyIfNotExists(CheckoutTableName, "ChannelID", channel.ChannelTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(CheckoutTableName, "BillingAddressID", account.AddressTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(CheckoutTableName, "ShippingAddressID", account.AddressTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(CheckoutTableName, "ShippingMethodID", shipping.ShippingMethodTableName, "Id", false)
}

func (cs *SqlCheckoutStore) Save(checkout *checkout.Checkout) (*checkout.Checkout, error) {
	checkout.PreSave()
	if err := checkout.IsValid(); err != nil {
		return nil, err
	}
	if err := cs.GetMaster().Insert(checkout); err != nil {
		return nil, err
	}
	return checkout, nil
}

func (cs *SqlCheckoutStore) Get(id string) (*checkout.Checkout, error) {
	iface, err := cs.GetReplica().Get(checkout.Checkout{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(CheckoutTableName, id)
		}
		return nil, err
	}

	return iface.(*checkout.Checkout), nil
}
