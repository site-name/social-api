package external_services

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/external_services"
	"github.com/sitename/sitename/store"
)

type SqlOpenExchangeRateStore struct {
	store.Store
}

func NewSqlOpenExchangeRateStore(s store.Store) store.OpenExchangeRateStore {
	os := &SqlOpenExchangeRateStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(external_services.OpenExchangeRate{}, store.OpenExchangeRateTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ToCurrency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH).SetUnique(true)
	}

	return os
}

func (os *SqlOpenExchangeRateStore) CreateIndexesIfNotExists() {
	os.CreateIndexIfNotExists("idx_openexchange_to_currency", store.OpenExchangeRateTableName, "ToCurrency")
}

// BulkUpsert performs bulk update/insert to given exchange rates
func (os *SqlOpenExchangeRateStore) BulkUpsert(rates []*external_services.OpenExchangeRate) ([]*external_services.OpenExchangeRate, error) {

	transaction, err := os.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	defer store.FinalizeTransaction(transaction)

	var isSaving bool
	for _, rate := range rates {
		if rate.Id == "" {
			rate.PreSave()
			isSaving = true
		} else {
			rate.PreUpdate()
		}

		if err := rate.IsValid(); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = transaction.Insert(rate)
		} else {
			_, err = transaction.Update(rate)
		}
		if err != nil {
			if os.IsUniqueConstraintError(err, []string{"ToCurrency", "openexchangerates_tocurrency_key"}) {
				return nil, store.NewErrInvalidInput(store.OpenExchangeRateTableName, "ToCurrency", rate.ToCurrency)
			}
			return nil, errors.Wrapf(err, "failed to upsert exchange rate with id=%s", rate.Id)
		}
	}

	if err = transaction.Commit(); err != nil {
		return nil, errors.Wrap(err, "transaction_commit")
	}

	return rates, nil
}

// GetAll returns all exchange currency rates
func (os *SqlOpenExchangeRateStore) GetAll() ([]*external_services.OpenExchangeRate, error) {
	var res []*external_services.OpenExchangeRate
	if _, err := os.GetReplica().Select(
		&res,
		"SELECT * FROM "+store.OpenExchangeRateTableName+" ORDER BY :OrderBy",
		map[string]interface{}{
			"OrderBy": store.TableOrderingMap[store.OpenExchangeRateTableName],
		},
	); err != nil {
		return nil, errors.Wrap(err, "failed to get all exchange rates")
	}

	return res, nil
}
