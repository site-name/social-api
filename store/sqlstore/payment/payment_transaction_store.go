package payment

import (
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
		table := db.AddTableWithName(payment.PaymentTransaction{}, "Transactions").SetKeys(false, "Id")
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

}
