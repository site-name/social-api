package shop

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlShopStaffStore struct {
	store.Store
}

func NewSqlShopStaffStore(s store.Store) store.ShopStaffStore {
	return &SqlShopStaffStore{s}
}

func (s *SqlShopStaffStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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

func (s *SqlShopStaffStore) FilterByOptions(options *model.ShopStaffRelationFilterOptions) ([]*model.ShopStaffRelation, error) {
	query := s.GetQueryBuilder().Select("*").From(store.ShopStaffTableName)

	if options.ShopID != nil {
		query = query.Where(options.ShopID)
	}
	if options.StaffID != nil {
		query = query.Where(options.StaffID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.ShopStaffRelation
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shop staff relations with given opsitons")
	}
	return res, nil
}
