package payment

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

type SqlPaymentTransactionStore struct {
	store.Store
}

func NewSqlPaymentTransactionStore(s store.Store) store.PaymentTransactionStore {
	ps := &SqlPaymentTransactionStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(payment.PaymentTransaction{}, store.TransactionTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("PaymentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(payment.MAX_LENGTH_PAYMENT_TOKEN)
		table.ColMap("Kind").SetMaxSize(payment.TRANSACTION_KIND_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)
		table.ColMap("Error").SetMaxSize(payment.TRANSACTION_ERROR_MAX_LENGTH)
		table.ColMap("CustomerID").SetMaxSize(payment.TRANSACTION_CUSTOMER_ID_MAX_LENGTH)
	}
	return ps
}

func (ps *SqlPaymentTransactionStore) CreateIndexesIfNotExists() {
	ps.CreateForeignKeyIfNotExists(store.TransactionTableName, "PaymentID", store.PaymentTableName, "Id", false)
}

// Save insert given transaction into database then returns it
func (ps *SqlPaymentTransactionStore) Save(transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, error) {
	transaction.PreSave()
	if err := transaction.IsValid(); err != nil {
		return nil, err
	}

	if err := ps.GetMaster().Insert(transaction); err != nil {
		return nil, errors.Wrapf(err, "failed to save payment transaction with id=%s", transaction.Id)
	}

	return transaction, nil
}

// Update updates given transaction then return it
func (ps *SqlPaymentTransactionStore) Update(transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, error) {
	transaction.PreUpdate()
	if err := transaction.IsValid(); err != nil {
		return nil, err
	}

	oldTran, err := ps.Get(transaction.Id)
	if err != nil {
		return nil, err
	}

	transaction.CreateAt = oldTran.CreateAt
	transaction.PaymentID = oldTran.PaymentID

	numUpdates, err := ps.GetMaster().Update(transaction)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update transaction with id=%s", transaction.Id)
	}

	if numUpdates > 1 {
		return nil, errors.Errorf("multiple transactions updated: %d", numUpdates)
	}

	return transaction, nil
}

func (ps *SqlPaymentTransactionStore) Get(id string) (*payment.PaymentTransaction, error) {
	var res payment.PaymentTransaction
	err := ps.GetReplica().SelectOne(&res, "SELECT * FROM "+store.TransactionTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.TransactionTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find payment transaction withh id=%s", id)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of transactions with given option
func (ps *SqlPaymentTransactionStore) FilterByOption(option *payment.PaymentTransactionFilterOpts) ([]*payment.PaymentTransaction, error) {
	query := ps.GetQueryBuilder().
		Select("*").
		From(store.TransactionTableName).
		OrderBy(store.TableOrderingMap[store.TransactionTableName])

	// parse options:
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}
	if option.PaymentID != nil {
		query = query.Where(option.PaymentID.ToSquirrel("PaymentID"))
	}
	if option.Kind != nil {
		query = query.Where(option.Kind.ToSquirrel("Kind"))
	}
	if option.ActionRequired != nil {
		query = query.Where(squirrel.Eq{"ActionRequired": *option.ActionRequired})
	}
	if option.IsSuccess != nil {
		query = query.Where(squirrel.Eq{"IsSuccess": *option.IsSuccess})
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*payment.PaymentTransaction
	_, err = ps.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find payment transactions based on given option")
	}

	return res, nil
}
