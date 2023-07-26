package shop

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlShopStaffStore struct {
	store.Store
}

func NewSqlShopStaffStore(s store.Store) store.ShopStaffStore {
	return &SqlShopStaffStore{s}
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
	if err := sss.GetMaster().Create(shopStaff).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to save shop-staff relation with id=%s", shopStaff.Id)
	}

	return shopStaff, nil
}

// Get finds a shop staff with given id then returns it with an error
func (sss *SqlShopStaffStore) Get(shopStaffID string) (*model.ShopStaff, error) {
	var res model.ShopStaff
	err := sss.GetReplica().First(&res, "Id = ?", shopStaffID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ShopStaffTableName, shopStaffID)
		}
		return nil, errors.Wrapf(err, "failed to finds shop staff relation with id=%s", shopStaffID)
	}

	return &res, nil
}

func (s *SqlShopStaffStore) commonQueryBuilder(options *model.ShopStaffFilterOptions) squirrel.SelectBuilder {
	selectFields := []string{model.ShopStaffTableName + ".*"}
	if options.SelectRelatedStaff {
		selectFields = append(selectFields, model.UserTableName+".*")
	}

	query := s.GetQueryBuilder().Select(selectFields...).
		From(model.ShopStaffTableName).
		Where(options.Conditions)
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

	rows, err := s.GetReplica().Raw(queryString, args...).Rows()
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

	err = s.GetReplica().Raw(queryString, args...).Row().Scan(scanFields...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.NewErrNotFound("ShopStaffs", "options")
		}
		return nil, errors.Wrap(err, "failed to scan shop-staff relation with given options")
	}

	if options.SelectRelatedStaff {
		relation.SetStaff(&staff)
	}

	return &relation, nil
}
