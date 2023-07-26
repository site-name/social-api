package shop

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type sqlVatStore struct {
	store.Store
}

func NewSqlVatStore(s store.Store) store.VatStore {
	return &sqlVatStore{s}
}

func (s *sqlVatStore) Upsert(transaction *gorm.DB, vats []*model.Vat) ([]*model.Vat, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	for _, vat := range vats {
		var err error
		if vat.Id == "" {
			err = transaction.Create(vat).Error
		} else {
			err = transaction.Model(vat).Updates(vat).Error
		}

		if err != nil {
			return nil, errors.Wrap(err, "failed to upsert a vat")
		}
	}

	return vats, nil
}

func (s *sqlVatStore) FilterByOptions(options *model.VatFilterOptions) ([]*model.Vat, error) {
	var res []*model.Vat
	err := s.GetReplica().Find(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find vat objects by options")
	}

	return res, nil
}
