package shop

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlShopStaffStore struct {
	store.Store
}

func NewSqlShopStaffStore(s store.Store) store.ShopStaffStore {
	return &SqlShopStaffStore{s}
}

func (sss *SqlShopStaffStore) Upsert(shopStaff model.ShopStaff) (*model.ShopStaff, error) {
	isSaving := shopStaff.ID != ""
	if isSaving {
		model_helper.ShopStaffPreSave(&shopStaff)
	} else {
		model_helper.ShopStaffCommonPre(&shopStaff)
	}

	if err := model_helper.ShopStaffIsValid(shopStaff); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = shopStaff.Insert(sss.GetMaster(), boil.Infer())
	} else {
		_, err = shopStaff.Update(sss.GetMaster(), boil.Blacklist(model.ShopStaffColumns.CreatedAt))
	}

	if err != nil {
		if sss.IsUniqueConstraintError(err, []string{"shop_staff_staff_id_unique_idx"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.ShopStaffs, model.ShopStaffColumns.StaffID, shopStaff.StaffID)
		}
		return nil, err
	}

	return &shopStaff, nil
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
	return model.ShopStaffs(options.Conditions...).All(s.GetReplica())
}

func (s *SqlShopStaffStore) GetByOptions(options model_helper.ShopStaffFilterOptions) (*model.ShopStaff, error) {
	record, err := model.ShopStaffs(options.Conditions...).One(s.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ShopStaffs, "options")
		}
		return nil, err
	}

	return record, nil
}
