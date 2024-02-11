package shop

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

type SqlShopStaffStore struct {
	store.Store
}

func NewSqlShopStaffStore(s store.Store) store.ShopStaffStore {
	return &SqlShopStaffStore{s}
}

func (sss *SqlShopStaffStore) Upsert(shopStaff model.ShopStaff) (*model.ShopStaff, error) {
	if err := sss.GetMaster().Create(shopStaff).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to save shop-staff relation with id=%s", shopStaff.ID)
	}

	return shopStaff, nil
}

// Get finds a shop staff with given id then returns it with an error
func (s *SqlShopStaffStore) Get(shopStaffID string) (*model.ShopStaff, error) {
	record, err := model.FindShopStaff(s.GetReplica(), shopStaffID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ShopStaffs, shopStaffID)
		}
		return nil, err
	}

	return record, nil
}

func (s *SqlShopStaffStore) FilterByOptions(options model_helper.ShopStaffFilterOptions) (model.ShopStaffSlice, error) {
	return model.ShopStaffs(options.Conds...).All(s.GetReplica())
}

func (s *SqlShopStaffStore) GetByOptions(options model_helper.ShopStaffFilterOptions) (*model.ShopStaff, error) {
	record, err := model.ShopStaffs(options.Conds...).One(s.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ShopStaffs, "options")
		}
		return nil, err
	}

	return record, nil
}
