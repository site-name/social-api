package shop

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
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
		"StaffID",
		"CreateAt",
		"EndAt",
		"SalaryPeriod",
		"Salary",
		"SalaryCurrency",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (s *SqlShopStaffStore) ScanFields(rel *model.ShopStaff) []interface{} {
	return []interface{}{
		&rel.Id,
		&rel.StaffID,
		&rel.CreateAt,
		&rel.EndAt,
		&rel.SalaryPeriod,
		&rel.Salary,
		&rel.SalaryCurrency,
	}
}

// Save inserts given shopStaff into database then returns it with an error
func (sss *SqlShopStaffStore) Save(shopStaff *model.ShopStaff) (*model.ShopStaff, error) {
	shopStaff.PreSave()
	if err := shopStaff.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.ShopStaffTableName + "(" + sss.ModelFields("").Join(",") + ") VALUES (" + sss.ModelFields(":").Join(",") + ")"
	if _, err := sss.GetMasterX().NamedExec(query, shopStaff); err != nil {
		return nil, errors.Wrapf(err, "failed to save shop-staff relation with id=%s", shopStaff.Id)
	}

	return shopStaff, nil
}

// Get finds a shop staff with given id then returns it with an error
func (sss *SqlShopStaffStore) Get(shopStaffID string) (*model.ShopStaff, error) {
	var res model.ShopStaff
	err := sss.GetReplicaX().Get(&res, "SELECT * FROM "+model.ShopStaffTableName+" WHERE Id = ?", shopStaffID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ShopStaffTableName, shopStaffID)
		}
		return nil, errors.Wrapf(err, "failed to finds shop staff relation with id=%s", shopStaffID)
	}

	return &res, nil
}

func (s *SqlShopStaffStore) commonQueryBuilder(options *model.ShopStaffFilterOptions) squirrel.SelectBuilder {
	selectFields := s.ModelFields(model.ShopStaffTableName + ".")
	if options.SelectRelatedStaff {
		selectFields = append(selectFields, s.User().ModelFields(model.UserTableName+".")...)
	}

	query := s.GetQueryBuilder().Select(selectFields...).From(model.ShopStaffTableName)

	if options.StaffID != nil {
		query = query.Where(options.StaffID)
	}
	if options.CreateAt != nil {
		query = query.Where(options.CreateAt)
	}
	if options.EndAt != nil {
		query = query.Where(options.EndAt)
	}
	if options.SelectRelatedStaff {
		query = query.InnerJoin(model.UserTableName + " ON Users.Id = ShopStaffs.StaffID")
	}

	return query
}

func (s *SqlShopStaffStore) FilterByOptions(options *model.ShopStaffFilterOptions) ([]*model.ShopStaff, error) {
	queryString, args, err := s.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	rows, err := s.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shop staff relations with given opsitons")
	}
	defer rows.Close()

	var res []*model.ShopStaff

	for rows.Next() {
		var relation model.ShopStaff
		var staff model.User
		var scanFields = s.ScanFields(&relation)
		if options.SelectRelatedStaff {
			scanFields = append(scanFields, s.User().ScanFields(&staff)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of shop-staff relation")
		}

		if options.SelectRelatedStaff {
			relation.SetStaff(&staff)
		}
		res = append(res, &relation)
	}
	return res, nil
}

func (s *SqlShopStaffStore) GetByOptions(options *model.ShopStaffFilterOptions) (*model.ShopStaff, error) {
	queryString, args, err := s.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var relation model.ShopStaff
	var staff model.User
	var scanFields = s.ScanFields(&relation)

	if options.SelectRelatedStaff {
		scanFields = append(scanFields, s.User().ScanFields(&staff))
	}

	err = s.GetReplicaX().QueryRowX(queryString, args...).Scan(scanFields...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("ShopStaffs", "options")
		}
		return nil, errors.Wrap(err, "failed to scan shop-staff relation with given options")
	}

	if options.SelectRelatedStaff {
		relation.SetStaff(&staff)
	}

	return &relation, nil
}
