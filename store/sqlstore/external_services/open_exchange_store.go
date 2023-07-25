package external_services

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlOpenExchangeRateStore struct {
	store.Store
}

func NewSqlOpenExchangeRateStore(s store.Store) store.OpenExchangeRateStore {
	return &SqlOpenExchangeRateStore{s}
}

// BulkUpsert performs bulk update/insert to given exchange rates
func (os *SqlOpenExchangeRateStore) BulkUpsert(rates []*model.OpenExchangeRate) ([]*model.OpenExchangeRate, error) {
	for _, rate := range rates {
		err := os.GetMaster().Save(rate).Error
		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert exchange rate with ToCurrency=%s", rate.ToCurrency)
		}
	}

	return rates, nil
}

// GetAll returns all exchange currency rates
func (os *SqlOpenExchangeRateStore) GetAll() ([]*model.OpenExchangeRate, error) {
	var res []*model.OpenExchangeRate
	err := os.GetReplica().Order("ToCurrency ASC").Find(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all exchange rates")
	}

	return res, nil
}
