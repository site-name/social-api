package shop

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
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

func (s *SqlShopStaffStore) ScanFields(rel *model.ShopStaffRelation) []interface{} {
	return []interface{}{
		&rel.Id,
		&rel.ShopID,
		&rel.StaffID,
		&rel.CreateAt,
		&rel.EndAt,
	}
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

func (s *SqlShopStaffStore) commonQueryBuilder(options *model.ShopStaffRelationFilterOptions) squirrel.SelectBuilder {
	selectFields := s.ModelFields(store.ShopStaffTableName + ".")
	if options.SelectRelatedShop {
		selectFields = append(selectFields, s.Shop().ModelFields(store.ShopTableName+".")...)
	}
	if options.SelectRelatedStaff {
		selectFields = append(selectFields, s.User().ModelFields(store.UserTableName+".")...)
	}

	query := s.GetQueryBuilder().Select(selectFields...).From(store.ShopStaffTableName)

	if options.ShopID != nil {
		query = query.Where(options.ShopID)
	}
	if options.StaffID != nil {
		query = query.Where(options.StaffID)
	}
	if options.CreateAt != nil {
		query = query.Where(options.CreateAt)
	}
	if options.EndAt != nil {
		query = query.Where(options.EndAt)
	}
	if options.SelectRelatedShop {
		query = query.InnerJoin(store.ShopTableName + " ON Shops.Id = ShopStaffs.ShopID")
	}
	if options.SelectRelatedStaff {
		query = query.InnerJoin(store.UserTableName + " ON Users.Id = ShopStaffs.StaffID")
	}

	return query
}

func (s *SqlShopStaffStore) FilterByOptions(options *model.ShopStaffRelationFilterOptions) ([]*model.ShopStaffRelation, error) {
	queryString, args, err := s.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	rows, err := s.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shop staff relations with given opsitons")
	}
	defer rows.Close()

	var res []*model.ShopStaffRelation
	var relation model.ShopStaffRelation
	var shop model.Shop
	var staff model.User
	var scanFields = s.ScanFields(&relation)
	if options.SelectRelatedShop {
		scanFields = append(scanFields, s.Shop().ScanFields(&shop)...)
	}
	if options.SelectRelatedStaff {
		scanFields = append(scanFields, s.User().ScanFields(&staff)...)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of shop-staff relation")
		}

		if options.SelectRelatedShop {
			relation.SetShop(&shop)
		}
		if options.SelectRelatedStaff {
			relation.SetStaff(&staff)
		}
		res = append(res, relation.DeepCopy())
	}
	return res, nil
}

func (s *SqlShopStaffStore) GetByOptions(options *model.ShopStaffRelationFilterOptions) (*model.ShopStaffRelation, error) {
	queryString, args, err := s.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var relation model.ShopStaffRelation
	var shop model.Shop
	var staff model.User
	var scanFields = s.ScanFields(&relation)

	if options.SelectRelatedShop {
		scanFields = append(scanFields, s.Shop().ScanFields(&shop)...)
	}
	if options.SelectRelatedStaff {
		scanFields = append(scanFields, s.User().ScanFields(&staff))
	}

	err = s.GetReplicaX().QueryRowX(queryString, args...).Scan(scanFields...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("ShopStaffs", "options")
		}
		return nil, errors.Wrap(err, "failed to scan shop-staff relation with given options")
	}
	if options.SelectRelatedShop {
		relation.SetShop(&shop)
	}
	if options.SelectRelatedStaff {
		relation.SetStaff(&staff)
	}

	return &relation, nil
}
