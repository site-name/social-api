package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherCategoryStore struct {
	store.Store
}

var (
	VoucherCategoryUniqueList = []string{
		"VoucherID", "CategoryID", "vouchercategories_voucherid_categoryid_key",
	}
)

func NewSqlVoucherCategoryStore(s store.Store) store.VoucherCategoryStore {
	vcs := &SqlVoucherCategoryStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherCategory{}, store.VoucherCategoryTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CategoryID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("VoucherID", "CategoryID")
	}

	return vcs
}

func (vcs *SqlVoucherCategoryStore) CreateIndexesIfNotExists() {
	vcs.CreateForeignKeyIfNotExists(store.VoucherCategoryTableName, "VoucherID", store.VoucherTableName, "Id", true)
	vcs.CreateForeignKeyIfNotExists(store.VoucherCategoryTableName, "CategoryID", store.ProductCategoryTableName, "Id", true)
}

// Upsert saves or updates given voucher category then returns it with an error
func (vcs *SqlVoucherCategoryStore) Upsert(voucherCategory *product_and_discount.VoucherCategory) (*product_and_discount.VoucherCategory, error) {
	var saving bool
	if voucherCategory.Id == "" {
		voucherCategory.PreSave()
		saving = true
	}
	if err := voucherCategory.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if saving {
		err = vcs.GetMaster().Insert(voucherCategory)
	} else {
		numUpdated, err = vcs.GetMaster().Update(voucherCategory)
	}

	if err != nil {
		if vcs.IsUniqueConstraintError(err, VoucherCategoryUniqueList) {
			return nil, store.NewErrInvalidInput(store.VoucherCategoryTableName, "VoucherID/CategoryID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher-category relation with id=%s", voucherCategory.Id)
	} else if numUpdated > 1 {
		return nil, errors.Errorf("multiple voucher-category relations updated: %d instead of 1", numUpdated)
	}

	return voucherCategory, nil
}

// Get finds a voucher category with given id, then returns it with an error
func (vcs *SqlVoucherCategoryStore) Get(voucherCategoryID string) (*product_and_discount.VoucherCategory, error) {
	result, err := vcs.GetReplica().Get(product_and_discount.VoucherCategory{}, voucherCategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherCategoryTableName, voucherCategoryID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher-category relation with id=%s", voucherCategoryID)
	}

	return result.(*product_and_discount.VoucherCategory), nil
}

// ProductCategoriesByVoucherID finds a list of product categories that have relationships with given voucher
func (vcs *SqlVoucherCategoryStore) ProductCategoriesByVoucherID(voucherID string) ([]*product_and_discount.Category, error) {
	var categories []*product_and_discount.Category
	_, err := vcs.GetReplica().Select(
		&categories,
		`SELECT * FROM `+store.ProductCategoryTableName+` WHERE Id IN (
			SELECT CategoryID from `+store.VoucherCategoryTableName+` WHERE (
				VoucherID = :VoucherID
			)
		)`,
		map[string]interface{}{
			"VoucherID": voucherID,
		},
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return []*product_and_discount.Category{}, store.NewErrNotFound(store.ProductCategoryTableName, "voucherID="+voucherID)
		}
		return nil, errors.Wrapf(err, "failed to find categories with relation to voucher with voucherId=%s", voucherID)
	}

	return categories, nil
}
