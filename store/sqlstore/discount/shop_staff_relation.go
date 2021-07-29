package discount

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/store"
)

type SqlShopStaffStore struct {
	store.Store
}

func NewSqlShopStaffStore(s store.Store) store.ShopStaffStore {
	sss := &SqlShopStaffStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shop.ShopStaffRelation{}, store.ShopStaffTableName).SetKeys(false, "id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShopID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("StaffID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ShopID", "StaffID")
	}

	return sss
}

func (sss *SqlShopStaffStore) CreateIndexesIfNotExists() {
	sss.CreateForeignKeyIfNotExists(store.ShopStaffTableName, "ShopID", store.ShopTableName, "Id", false)
	sss.CreateForeignKeyIfNotExists(store.ShopStaffTableName, "StaffID", store.UserTableName, "Id", false)
}

// Save inserts given shopStaff into database then returns it with an error
func (sss *SqlShopStaffStore) Save(shopStaff *shop.ShopStaffRelation) (*shop.ShopStaffRelation, error) {
	shopStaff.PreSave()
	if err := shopStaff.IsValid(); err != nil {
		return nil, err
	}

	if err := sss.GetMaster().Insert(shopStaff); err != nil {
		if sss.IsUniqueConstraintError(err, []string{"ShopID", "StaffID", "shopstaffs_shopid_staffid_key"}) {
			return nil, store.NewErrInvalidInput(store.ShopStaffTableName, "ShopID/StaffID", "unique values")
		}
		return nil, errors.Wrapf(err, "failed to save shop-staff relation with id=%s", shopStaff.Id)
	}

	return shopStaff, nil
}

// Get finds a shop staff with given id then returns it with an error
func (sss *SqlShopStaffStore) Get(shopStaffID string) (*shop.ShopStaffRelation, error) {
	result, err := sss.GetReplica().Get(shop.ShopStaffRelation{}, shopStaffID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShopStaffTableName, shopStaffID)
		}
		return nil, errors.Wrapf(err, "failed to finds shop staff relation with id=%s", shopStaffID)
	}

	return result.(*shop.ShopStaffRelation), nil
}

// FilterByShopAndStaff finds a relation ship with given shopId and staffId
func (sss *SqlShopStaffStore) FilterByShopAndStaff(shopID string, staffID string) (*shop.ShopStaffRelation, error) {
	var result *shop.ShopStaffRelation
	err := sss.GetReplica().SelectOne(
		&result,
		`SELECT * FROM `+store.ShopStaffTableName+`
		WHERE (
			ShopID = :ShopID AND StaffID = :StaffID
		)`,
		map[string]interface{}{
			"ShopID":  shopID,
			"StaffID": staffID,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShopStaffTableName, fmt.Sprintf("ShopID=%s, StaffID=%s", shopID, staffID))
		}
		return nil, errors.Wrapf(err, "failed to find shop-staff relation with ShopID=%s, StaffID=%s", shopID, staffID)
	}

	return result, nil
}
