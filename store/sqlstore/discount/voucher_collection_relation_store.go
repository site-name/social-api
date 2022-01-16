package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherCollectionStore struct {
	store.Store
}

func NewSqlVoucherCollectionStore(s store.Store) store.VoucherCollectionStore {
	vcs := &SqlVoucherCollectionStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherCollection{}, store.VoucherCollectionTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CollectionID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("VoucherID", "CollectionID")
	}

	return vcs
}

func (vcs *SqlVoucherCollectionStore) TableName(withField string) string {
	name := "VoucherCollections"
	if withField != "" {
		name += "." + withField
	}

	return name
}

func (vcs *SqlVoucherCollectionStore) CreateIndexesIfNotExists() {
	vcs.CreateForeignKeyIfNotExists(store.VoucherCollectionTableName, "VoucherID", store.VoucherTableName, "Id", true)
	vcs.CreateForeignKeyIfNotExists(store.VoucherCollectionTableName, "CollectionID", store.ProductCollectionTableName, "Id", true)
}

// Upsert saves or updates given voucher collection then returns it with an error
func (vcs *SqlVoucherCollectionStore) Upsert(voucherCollection *product_and_discount.VoucherCollection) (*product_and_discount.VoucherCollection, error) {
	var saving bool
	if voucherCollection.Id == "" {
		voucherCollection.PreSave()
		saving = true
	}
	if err := voucherCollection.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if saving {
		err = vcs.GetMaster().Insert(voucherCollection)
	} else {
		numUpdated, err = vcs.GetMaster().Update(voucherCollection)
	}

	if err != nil {
		if vcs.IsUniqueConstraintError(err, []string{"VoucherID", "CollectionID", "vouchercollections_voucherid_collectionid_key"}) {
			return nil, store.NewErrInvalidInput(store.VoucherCollectionTableName, "VoucherID/CollectionID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher-collection relation with id=%s", voucherCollection.Id)
	} else if numUpdated > 1 {
		return nil, errors.Errorf("multiple voucher-collection relations updated: %d instead of 1", numUpdated)
	}

	return voucherCollection, nil
}

// Get finds a voucher collection with given id, then returns it with an error
func (vcs *SqlVoucherCollectionStore) Get(voucherCollectionID string) (*product_and_discount.VoucherCollection, error) {
	var res product_and_discount.VoucherCollection
	err := vcs.GetReplica().SelectOne(&res, "SELECT * FROM "+store.VoucherCollectionTableName+" WHERE Id = :ID", map[string]interface{}{"ID": voucherCollectionID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherCollectionTableName, voucherCollectionID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher-collection relation with id=%s", voucherCollectionID)
	}

	return &res, nil
}
