package external_services

import (
	"database/sql"

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
	transaction, err := os.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	defer store.FinalizeTransaction(transaction)

	for _, rate := range rates {
		var (
			oldRate  model.OpenExchangeRate
			isSaving bool
		)
		// try lookup:
		err := transaction.QueryRow(
			"SELECT * FROM "+model.OpenExchangeRateTableName+" WHERE ToCurrency = $1 FOR UPDATE",
			rate.ToCurrency,
		).
			Scan(&oldRate.Id, &oldRate.ToCurrency, &oldRate.Rate)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, errors.Wrapf(err, "failed to find exchange rate with ToCurrency=%s", rate.ToCurrency)
			}
			isSaving = true
		}

		if isSaving {
			rate.Id = ""
			rate.PreSave()
		} else {
			rate.PreUpdate()
		}

		if err := rate.IsValid(); err != nil {
			return nil, err
		}

		if isSaving {
			_, err = transaction.NamedExec("INSERT INTO "+model.OpenExchangeRateTableName+"(Id, ToCurrency, Rate) VALUES (:Id, :ToCurrency, :Rate)", rate)
		} else {
			// check if rates are different then update
			if !rate.Rate.Equal(*oldRate.Rate) {
				rate.Id = oldRate.Id
				_, err = transaction.NamedExec("UPDATE "+model.OpenExchangeRateTableName+" SET Rate=:Rate WHERE Id=:Id", rate)
			}
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert exchange rate with ToCurrency=%s", rate.ToCurrency)
		}
	}

	if err = transaction.Commit(); err != nil {
		return nil, errors.Wrap(err, "transaction_commit")
	}

	return rates, nil
}

// GetAll returns all exchange currency rates
func (os *SqlOpenExchangeRateStore) GetAll() ([]*model.OpenExchangeRate, error) {
	var res []*model.OpenExchangeRate
	err := os.GetReplica().Select(
		&res,
		"SELECT * FROM "+model.OpenExchangeRateTableName+" ORDER BY ToCurrency ASC",
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get all exchange rates")
	}

	return res, nil
}
