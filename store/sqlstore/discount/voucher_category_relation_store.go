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

func NewSqlVoucherCategoryStore(s store.Store) store.VoucherCategoryStore {
	vcs := &SqlVoucherCategoryStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherCategory{}, vcs.TableName("")).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CategoryID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("VoucherID", "CategoryID")
	}

	return vcs
}

func (vcs *SqlVoucherCategoryStore) TableName(withField string) string {
	name := "VoucherCategories"
	if withField != "" {
		name += "." + withField
	}

	return name
}

func (vcs *SqlVoucherCategoryStore) CreateIndexesIfNotExists() {
	vcs.CreateForeignKeyIfNotExists(vcs.TableName(""), "VoucherID", store.VoucherTableName, "Id", true)
	vcs.CreateForeignKeyIfNotExists(vcs.TableName(""), "CategoryID", store.ProductCategoryTableName, "Id", true)
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
		if vcs.IsUniqueConstraintError(err, []string{"VoucherID", "CategoryID", "vouchercategories_voucherid_categoryid_key"}) {
			return nil, store.NewErrInvalidInput(vcs.TableName(""), "VoucherID/CategoryID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher-category relation with id=%s", voucherCategory.Id)
	} else if numUpdated > 1 {
		return nil, errors.Errorf("multiple voucher-category relations updated: %d instead of 1", numUpdated)
	}

	return voucherCategory, nil
}

// Get finds a voucher category with given id, then returns it with an error
func (vcs *SqlVoucherCategoryStore) Get(voucherCategoryID string) (*product_and_discount.VoucherCategory, error) {
	var res product_and_discount.VoucherCategory
	err := vcs.GetReplica().SelectOne(&res, "SELECT * FROM "+vcs.TableName("")+" WHERE Id = :ID", map[string]interface{}{"ID": voucherCategoryID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(vcs.TableName(""), voucherCategoryID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher-category relation with id=%s", voucherCategoryID)
	}

	return &res, nil
}
