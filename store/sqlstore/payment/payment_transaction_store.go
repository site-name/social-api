package payment

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

type SqlPaymentTransactionStore struct {
	store.Store
}

const transactionTableName = "Transactions"

func NewSqlPaymentTransactionStore(s store.Store) store.PaymentTransactionStore {
	ps := &SqlPaymentTransactionStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(payment.PaymentTransaction{}, transactionTableName).SetKeys(false, "Id")
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

func (ps *SqlPaymentTransactionStore) CreateIndexesIfNotExists() {}

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

func (ps *SqlPaymentTransactionStore) Get(id string) (*payment.PaymentTransaction, error) {
	transacResult, err := ps.GetReplica().Get(payment.PaymentTransaction{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(transactionTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find payment transaction withh id=%s", id)
	}

	return transacResult.(*payment.PaymentTransaction), nil
}

func (ps *SqlPaymentTransactionStore) GetAllByPaymentID(paymentID string) ([]*payment.PaymentTransaction, error) {
	var transactions []*payment.PaymentTransaction

	if _, err := ps.GetReplica().Select(
		&transactions,
		"SELECT * FROM "+transactionTableName+" WHERE PaymentID = :paymentID",
		map[string]interface{}{"paymentID": paymentID},
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(transactionTableName, "paymentID="+paymentID)
		}
		return nil, errors.Wrapf(err, "failed to find transactions belong to payment with id=%s", paymentID)
	}

	return transactions, nil
}
