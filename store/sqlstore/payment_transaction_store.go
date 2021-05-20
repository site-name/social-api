package sqlstore

import (
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

type SqlPaymentTransactionStore struct {
	*SqlStore
}

func newSqlPaymentTransactionStore(s *SqlStore) store.PaymentTransactionStore {
	ps := &SqlPaymentTransactionStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(payment.Payment{}, "Transactions").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("PageTypeID").SetMaxSize(UUID_MAX_LENGTH)
		// table.ColMap("Title").SetMaxSize(payment.PAGE_TITLE_MAX_LENGTH)
		// table.ColMap("Slug").SetMaxSize(payment.PAGE_SLUG_MAX_LENGTH)
	}
	return ps
}

func (ps *SqlPaymentTransactionStore) createIndexesIfNotExists() {
	// ps.CreateIndexIfNotExists("idx_pages_title", "Transactions", "Title")
	// ps.CreateIndexIfNotExists("idx_pages_slug", "Transactions", "Slug")

	// ps.CreateIndexIfNotExists("idx_pages_title_lower_textpattern", "Transactions", "lower(Title) text_pattern_ops")
}
