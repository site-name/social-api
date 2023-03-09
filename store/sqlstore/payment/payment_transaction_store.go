package payment

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlPaymentTransactionStore struct {
	store.Store
}

func NewSqlPaymentTransactionStore(s store.Store) store.PaymentTransactionStore {
	return &SqlPaymentTransactionStore{s}
}

func (s *SqlPaymentTransactionStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"CreateAt",
		"PaymentID",
		"Token",
		"Kind",
		"IsSuccess",
		"ActionRequired",
		"ActionRequiredData",
		"Currency",
		"Amount",
		"Error",
		"CustomerID",
		"GatewayResponse",
		"AlreadyProcessed",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save insert given transaction into database then returns it
func (ps *SqlPaymentTransactionStore) Save(transaction store_iface.SqlxTxExecutor, paymentTransaction *model.PaymentTransaction) (*model.PaymentTransaction, error) {
	var executor store_iface.SqlxExecutor = ps.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	paymentTransaction.PreSave()
	if err := paymentTransaction.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.TransactionTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
	if _, err := executor.NamedExec(query, paymentTransaction); err != nil {
		return nil, errors.Wrapf(err, "failed to save payment paymentTransaction with id=%s", paymentTransaction.Id)
	}

	return paymentTransaction, nil
}

// Update updates given transaction then return it
func (ps *SqlPaymentTransactionStore) Update(transaction *model.PaymentTransaction) (*model.PaymentTransaction, error) {
	transaction.PreUpdate()
	if err := transaction.IsValid(); err != nil {
		return nil, err
	}

	query := "UPDATE " + store.TransactionTableName + " SET " + ps.
		ModelFields("").
		Map(func(_ int, s string) string {
			return s + "=:" + s
		}).
		Join(",") + " WHERE Id=:Id"

	result, err := ps.GetMasterX().NamedExec(query, transaction)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update transaction with id=%s", transaction.Id)
	}

	if numUpdates, _ := result.RowsAffected(); numUpdates > 1 {
		return nil, errors.Errorf("multiple transactions updated: %d", numUpdates)
	}

	return transaction, nil
}

func (ps *SqlPaymentTransactionStore) Get(id string) (*model.PaymentTransaction, error) {
	var res model.PaymentTransaction
	err := ps.GetReplicaX().Get(&res, "SELECT * FROM "+store.TransactionTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.TransactionTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find payment transaction withh id=%s", id)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of transactions with given option
func (ps *SqlPaymentTransactionStore) FilterByOption(option *model.PaymentTransactionFilterOpts) ([]*model.PaymentTransaction, error) {
	query := ps.GetQueryBuilder().
		Select("*").
		From(store.TransactionTableName).
		OrderBy(store.TableOrderingMap[store.TransactionTableName])

	// parse options:
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.PaymentID != nil {
		query = query.Where(option.PaymentID)
	}
	if option.Kind != nil {
		query = query.Where(option.Kind)
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

	var res []*model.PaymentTransaction
	err = ps.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find payment transactions based on given option")
	}

	return res, nil
}
