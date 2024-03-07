package shop

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type sqlVatStore struct {
	store.Store
}

func NewSqlVatStore(s store.Store) store.VatStore {
	return &sqlVatStore{s}
}

func (s *sqlVatStore) Upsert(tx boil.ContextTransactor, vats model.VatSlice) (model.VatSlice, error) {
	if tx == nil {
		tx = s.GetMaster()
	}

	for _, vat := range vats {
		var err error
		if vat.ID == "" {
			err = vat.Insert(tx, boil.Infer())
		} else {
			_, err = vat.Update(tx, boil.Infer())
		}

		if err != nil {
			return nil, errors.Wrap(err, "failed to upsert a vat")
		}
	}

	return vats, nil
}

func (s *sqlVatStore) FilterByOptions(options model_helper.VatFilterOptions) (model.VatSlice, error) {
	return model.Vats(options.Conditions...).All(s.GetReplica())
}
