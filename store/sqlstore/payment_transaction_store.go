package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

type SqlPaymentTransactionStore struct {
	*SqlStore
}

func newSqlPaymentTransactionStore(s *SqlStore) store.PaymentTransactionStore {
	ps := &SqlPaymentTransactionStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(payment.PaymentTransaction{}, "Transactions").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("PaymentID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(payment.MAX_LENGTH_PAYMENT_TOKEN)
		table.ColMap("Kind").SetMaxSize(payment.TRANSACTION_KIND_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)
		table.ColMap("Error").SetMaxSize(payment.TRANSACTION_ERROR_MAX_LENGTH)
		table.ColMap("CustomerID").SetMaxSize(payment.TRANSACTION_CUSTOMER_ID_MAX_LENGTH)
	}
	return ps
}

func (ps *SqlPaymentTransactionStore) createIndexesIfNotExists() {
	// ps.CreateIndexIfNotExists("idx_pages_title", "Transactions", "Title")
	// ps.CreateIndexIfNotExists("idx_pages_slug", "Transactions", "Slug")

	// ps.CreateIndexIfNotExists("idx_pages_title_lower_textpattern", "Transactions", "lower(Title) text_pattern_ops")
}