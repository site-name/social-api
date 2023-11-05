package discount

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlVoucherTranslationStore struct {
	store.Store
}

func NewSqlVoucherTranslationStore(sqlStore store.Store) store.VoucherTranslationStore {
	return &SqlVoucherTranslationStore{sqlStore}
}

// Save inserts given translation into database and returns it
func (vts *SqlVoucherTranslationStore) Save(translation *model.VoucherTranslation) (*model.VoucherTranslation, error) {
	err := vts.GetMaster().Create(translation).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save voucher translation with id=%s", translation.Id)
	}

	return translation, nil
}

// Get finds and returns a voucher translation with given id
func (vts *SqlVoucherTranslationStore) Get(id string) (*model.VoucherTranslation, error) {
	var res model.VoucherTranslation
	err := vts.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.VoucherTranslationTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find voucher translation with id=%s", id)
	}

	return &res, nil
}

// FilterByOption returns a list of voucher translations filtered using given options
func (vts *SqlVoucherTranslationStore) FilterByOption(option *model.VoucherTranslationFilterOption) ([]*model.VoucherTranslation, error) {
	args, err := store.BuildSqlizer(option.Conditions, "VoucherTranslation_FilterByOption")
	if err != nil {
		return nil, err
	}

	var res []*model.VoucherTranslation
	err = vts.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher translations with given options")
	}

	return res, nil
}

// GetByOption finds and returns 1 voucher translation by given options
func (vts *SqlVoucherTranslationStore) GetByOption(option *model.VoucherTranslationFilterOption) (*model.VoucherTranslation, error) {
	args, err := store.BuildSqlizer(option.Conditions, "VoucherTranslation_GetByOptions")
	if err != nil {
		return nil, err
	}

	var res model.VoucherTranslation
	err = vts.GetReplica().First(&res, args...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.VoucherTranslationTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find a voucher translation by given option")
	}

	return &res, nil
}
