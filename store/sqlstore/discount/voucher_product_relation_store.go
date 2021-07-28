package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherProductStore struct {
	store.Store
}

var (
	VoucherProductDuplicateList = []string{
		"VoucherID", "ProductID", "voucherproducts_voucherid_productid_key",
	}
)

func NewSqlVoucherProductStore(s store.Store) store.VoucherProductStore {
	vps := &SqlVoucherProductStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherProduct{}, store.VoucherProductTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("VoucherID", "ProductID")
	}

	return vps
}

func (vps *SqlVoucherProductStore) CreateIndexesIfNotExists() {
	vps.CreateForeignKeyIfNotExists(store.VoucherProductTableName, "VoucherID", store.VoucherTableName, "Id", true)
	vps.CreateForeignKeyIfNotExists(store.VoucherProductTableName, "ProductID", store.ProductTableName, "Id", true)
}

// Upsert saves or updates given voucher product then returns it with an error
func (vps *SqlVoucherProductStore) Upsert(voucherProduct *product_and_discount.VoucherProduct) (*product_and_discount.VoucherProduct, error) {
	var saving bool
	if voucherProduct.Id == "" {
		voucherProduct.PreSave()
		saving = true
	}
	if err := voucherProduct.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if saving {
		err = vps.GetMaster().Insert(voucherProduct)
	} else {
		numUpdated, err = vps.GetMaster().Update(voucherProduct)
	}

	if err != nil {
		if vps.IsUniqueConstraintError(err, VoucherProductDuplicateList) {
			return nil, store.NewErrInvalidInput(store.VoucherProductTableName, "VoucherID/ProductID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher-product relation with id=%s", voucherProduct.Id)
	} else if numUpdated > 1 {
		return nil, errors.Errorf("multiple voucher-product relations updated: %d instead of 1", numUpdated)
	}

	return voucherProduct, nil
}

// Get finds a voucher product with given id, then returns it with an error
func (vps *SqlVoucherProductStore) Get(voucherProductID string) (*product_and_discount.VoucherProduct, error) {
	result, err := vps.GetReplica().Get(product_and_discount.VoucherProduct{}, voucherProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherProductTableName, voucherProductID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher-product relation with id=%s", voucherProductID)
	}

	return result.(*product_and_discount.VoucherProduct), nil
}

// ProductsByVoucherID finds all products that have relationships with given voucher
func (vps *SqlVoucherProductStore) ProductsByVoucherID(voucherID string) ([]*product_and_discount.Product, error) {
	var products []*product_and_discount.Product
	_, err := vps.GetReplica().Select(
		&products,
		`SELECT * FROM `+store.ProductTableName+` WHERE Id IN (
			SELECT ProductID from `+store.VoucherProductTableName+` WHERE (
				VoucherID = :VoucherID
			)
		)`,
		map[string]interface{}{
			"VoucherID": voucherID,
		},
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return []*product_and_discount.Product{}, store.NewErrNotFound(store.ProductTableName, "voucherID="+voucherID)
		}
		return nil, errors.Wrapf(err, "failed to find products with relation to voucher with voucherId=%s", voucherID)
	}

	return products, nil
}
