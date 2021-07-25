package checkout

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore/shipping"
)

type SqlCheckoutStore struct {
	store.Store
}

func NewSqlCheckoutStore(sqlStore store.Store) store.CheckoutStore {
	cs := &SqlCheckoutStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(checkout.Checkout{}, store.CheckoutTableName).SetKeys(false, "Token")
		table.ColMap("Token").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
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
	cs.CreateIndexIfNotExists("idx_checkouts_userid", store.CheckoutTableName, "UserID")
	cs.CreateIndexIfNotExists("idx_checkouts_token", store.CheckoutTableName, "Token")
	cs.CreateIndexIfNotExists("idx_checkouts_channelid", store.CheckoutTableName, "ChannelID")
	cs.CreateIndexIfNotExists("idx_checkouts_billing_address_id", store.CheckoutTableName, "BillingAddressID")
	cs.CreateIndexIfNotExists("idx_checkouts_shipping_address_id", store.CheckoutTableName, "ShippingAddressID")
	cs.CreateIndexIfNotExists("idx_checkouts_shipping_method_id", store.CheckoutTableName, "ShippingMethodID")

	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "UserID", store.UserTableName, "Id", true)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "ChannelID", store.ChannelTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "BillingAddressID", store.AddressTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "ShippingAddressID", store.AddressTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "ShippingMethodID", shipping.ShippingMethodTableName, "Id", false)
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
			return nil, store.NewErrNotFound(store.CheckoutTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find checkout with id=%s", id)
	}

	return iface.(*checkout.Checkout), nil
}

func (cs *SqlCheckoutStore) Update(ckout *checkout.Checkout) (*checkout.Checkout, error) {
	ckout.PreUpdate()
	if err := ckout.IsValid(); err != nil {
		return nil, err
	}

	// try finding if checkeout exist:
	result, err := cs.GetReplica().Get(checkout.Checkout{}, ckout.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CheckoutTableName, "Token="+ckout.Token)
		}
		return nil, errors.Wrapf(err, "failed to find a checkout with token=%s", ckout.Token)
	}
	// Found a checkout. now set fields that cannot be modified:
	oldCheckout := result.(*checkout.Checkout)
	ckout.BillingAddressID = oldCheckout.BillingAddressID
	ckout.ShippingAddressID = oldCheckout.ShippingAddressID

	numUpdate, err := cs.GetMaster().Update(ckout)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update checkout with token=%s", ckout.Token)
	}
	if numUpdate > 1 {
		return nil, errors.New("multiple checkout was updated instead of one")
	}

	return ckout, nil
}

func (cs *SqlCheckoutStore) CheckoutsByUserID(userID string, channelActive bool) ([]*checkout.Checkout, error) {
	var checkouts []*checkout.Checkout

	query := `SELECT * FROM ` + store.CheckoutTableName + ` AS Ck 
	INNER JOIN ` + store.ChannelTableName + ` AS Cn ON (
		Cn.Id = Ck.ChannelID
	)`
	condition := `Ck.UserID = :UserID`

	if channelActive {
		condition += ` AND Cn.IsActive`
	} else {
		condition += ` AND NOT Cn.IsActive`
	}

	query += `WHERE (` + condition + `)`

	_, err := cs.GetReplica().Select(&checkouts, query, map[string]interface{}{"UserID": userID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CheckoutTableName, "userID="+userID)
		}
		return nil, errors.Wrapf(err, "failed to find checkouts for user with Id=%s", userID)
	}

	return checkouts, nil
}
