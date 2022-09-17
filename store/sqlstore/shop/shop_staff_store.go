package shop

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlShopStaffStore struct {
	store.Store
}

func NewSqlShopStaffStore(s store.Store) store.ShopStaffStore {
	return &SqlShopStaffStore{s}
}

func (s *SqlShopStaffStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"ShopID",
		"StaffID",
		"CreateAt",
		"EndAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given shopStaff into database then returns it with an error
func (sss *SqlShopStaffStore) Save(shopStaff *model.ShopStaffRelation) (*model.ShopStaffRelation, error) {
	shopStaff.PreSave()
	if err := shopStaff.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.ShopStaffTableName + "(" + sss.ModelFields("").Join(",") + ") VALUES (" + sss.ModelFields(":").Join(",") + ")"
	if _, err := sss.GetMasterX().NamedExec(query, shopStaff); err != nil {
		if sss.IsUniqueConstraintError(err, []string{"ShopID", "StaffID", "shopstaffs_shopid_staffid_key"}) {
			return nil, store.NewErrInvalidInput(store.ShopStaffTableName, "ShopID/StaffID", "unique values")
		}
		return nil, errors.Wrapf(err, "failed to save shop-staff relation with id=%s", shopStaff.Id)
	}

	return shopStaff, nil
}

// Get finds a shop staff with given id then returns it with an error
func (sss *SqlShopStaffStore) Get(shopStaffID string) (*model.ShopStaffRelation, error) {
	var res model.ShopStaffRelation
	err := sss.GetReplicaX().Get(&res, "SELECT * FROM "+store.ShopStaffTableName+" WHERE Id = ?", shopStaffID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShopStaffTableName, shopStaffID)
		}
		return nil, errors.Wrapf(err, "failed to finds shop staff relation with id=%s", shopStaffID)
	}

	return &res, nil
}

// FilterByShopAndStaff finds a relation ship with given shopId and staffId
func (sss *SqlShopStaffStore) FilterByShopAndStaff(shopID string, staffID string) (*model.ShopStaffRelation, error) {
	var result model.ShopStaffRelation
	err := sss.GetReplicaX().Get(
		&result,
		`SELECT * FROM `+store.ShopStaffTableName+`
		WHERE (
			ShopID = ? AND StaffID = ?
		)`,
		shopID,
		staffID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShopStaffTableName, fmt.Sprintf("ShopID=%s, StaffID=%s", shopID, staffID))
		}
		return nil, errors.Wrapf(err, "failed to find shop-staff relation with ShopID=%s, StaffID=%s", shopID, staffID)
	}

	return &result, nil
}
