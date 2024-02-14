package external_services

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlOpenExchangeRateStore struct {
	store.Store
}

func NewSqlOpenExchangeRateStore(s store.Store) store.OpenExchangeRateStore {
	return &SqlOpenExchangeRateStore{s}
}

// BulkUpsert performs bulk update/insert to given exchange rates
func (os *SqlOpenExchangeRateStore) BulkUpsert(rates model.OpenExchangeRateSlice) (model.OpenExchangeRateSlice, error) {
	for _, rate := range rates {
		if rate == nil {
			continue
		}

		isSaving := false
		if rate.ID == "" {
			isSaving = true
		}
		model_helper.OpenExchangeRateCommonPre(rate)

		if err := model_helper.OpenExchangeRateIsValid(*rate); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = rate.Insert(os.GetMaster(), boil.Infer())
		} else {
			_, err = rate.Update(os.GetMaster(), boil.Blacklist(model.OpenExchangeRateColumns.CreatedAt))
		}

		if err != nil {
			if os.IsUniqueConstraintError(err, []string{model.OpenExchangeRateColumns.ToCurrency, "open_exchange_rates_to_currency_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.OpenExchangeRates, model.OpenExchangeRateColumns.ToCurrency, "unique")
			}
			return nil, err
		}
	}

	return rates, nil
}

// GetAll returns all exchange currency rates
func (os *SqlOpenExchangeRateStore) GetAll() (model.OpenExchangeRateSlice, error) {
	return model.OpenExchangeRates().All(os.GetReplica())
}
