package checkout

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

type SqlCheckoutLineStore struct {
	store.Store
}

func NewSqlCheckoutLineStore(sqlStore store.Store) store.CheckoutLineStore {
	cls := &SqlCheckoutLineStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(checkout.CheckoutLine{}, store.CheckoutLineTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CheckoutID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH)
	}
	return cls
}

func (cls *SqlCheckoutLineStore) CreateIndexesIfNotExists() {
	cls.CreateIndexIfNotExists("idx_checkoutlines_checkout_id", store.CheckoutLineTableName, "CheckoutID")
	cls.CreateIndexIfNotExists("idx_checkoutlines_variant_id", store.CheckoutLineTableName, "VariantID")

	// foreign keys:
	cls.CreateForeignKeyIfNotExists(store.CheckoutLineTableName, "CheckoutID", store.CheckoutTableName, "Id", true)
	cls.CreateForeignKeyIfNotExists(store.CheckoutLineTableName, "VariantID", store.ProductVariantTableName, "Id", true)
}

func (cls *SqlCheckoutLineStore) Save(cl *checkout.CheckoutLine) (*checkout.CheckoutLine, error) {
	cl.PreSave()
	if err := cl.IsValid(); err != nil {
		return nil, err
	}

	if err := cls.GetMaster().Insert(cl); err != nil {
		return nil, errors.Wrapf(err, "failed to save checkout line with id=%s", cl.Id)
	}

	return cl, nil
}

func (cls *SqlCheckoutLineStore) Get(id string) (*checkout.CheckoutLine, error) {
	res, err := cls.GetReplica().Get(checkout.CheckoutLine{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CheckoutLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to to find checkout line with id=%s", id)
	}

	return res.(*checkout.CheckoutLine), nil
}

func (cls *SqlCheckoutLineStore) CheckoutLinesByCheckoutID(checkoutID string) ([]*checkout.CheckoutLine, error) {
	var res []*checkout.CheckoutLine
	_, err := cls.GetReplica().Select(&res, "SELECT * FROM "+store.CheckoutLineTableName+" WHERE CheckoutID = :CheckoutID", map[string]interface{}{"CheckoutID": checkoutID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CheckoutLineTableName, "checkoutID="+checkoutID)
		}
		return nil, errors.Wrapf(err, "failed to get checkout lines belong to checkout with id=%s", checkoutID)
	}

	return res, nil
}