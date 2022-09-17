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
	transaction, err := os.GetMasterX().Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	defer store.FinalizeTransaction(transaction)

	var (
		oldRate    model.OpenExchangeRate
		numUpdated int64
	)

	for _, rate := range rates {
		isSaving := false
		// try lookup:
		err := transaction.Get(
			&oldRate,
			"SELECT * FROM "+store.OpenExchangeRateTableName+" WHERE ToCurrency = ? FOR UPDATE",
			rate.ToCurrency,
		)
		if err != nil {
			if err == sql.ErrNoRows { // does not exist
				isSaving = true
			} else {
				return nil, errors.Wrapf(err, "failed to find exchange rate with ToCurrency=%s", rate.ToCurrency)
			}
		}

		if isSaving {
			rate.PreSave()
		} else {
			rate.PreUpdate()
		}

		if err := rate.IsValid(); err != nil {
			return nil, err
		}

		if isSaving {
			query := "INSERT INTO " + store.OpenExchangeRateTableName + "(Id, ToCurrency, Rate) VALUES (:Id, :ToCurrency, :Rate)"
			_, err = transaction.NamedExec(query, rate)

		} else {
			// check if rates are different then update
			if !rate.Rate.Equal(*oldRate.Rate) {
				rate.Id = oldRate.Id

				query := "UPDATE " + store.OpenExchangeRateTableName + " SET Id=:Id, ToCurrency=:ToCurrency, Rate=:Rate WHERE Id=:Id"
				var result sql.Result
				result, err = transaction.NamedExec(query, rate)
				if err == nil && result != nil {
					numUpdated, _ = result.RowsAffected()
				}
			}
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert exchange rate with ToCurrency=%s", rate.ToCurrency)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("multiple exchange rates were updated: %d instead of 1", numUpdated)
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
	err := os.GetReplicaX().Select(
		&res,
		"SELECT * FROM "+store.OpenExchangeRateTableName+" ORDER BY ?",
		store.TableOrderingMap[store.OpenExchangeRateTableName],
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get all exchange rates")
	}

	return res, nil
}
