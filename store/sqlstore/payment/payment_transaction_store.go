package payment

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlPaymentTransactionStore struct {
	store.Store
}

func NewSqlPaymentTransactionStore(s store.Store) store.PaymentTransactionStore {
	return &SqlPaymentTransactionStore{s}
}

// Save insert given transaction into database then returns it
func (ps *SqlPaymentTransactionStore) Save(transaction *gorm.DB, paymentTransaction *model.PaymentTransaction) (*model.PaymentTransaction, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}
	if err := transaction.Create(paymentTransaction).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to save payment paymentTransaction with id=%s", paymentTransaction.Id)
	}

	return paymentTransaction, nil
}

// Update updates given transaction then return it
func (ps *SqlPaymentTransactionStore) Update(transaction *model.PaymentTransaction) (*model.PaymentTransaction, error) {
	transaction.CreateAt = 0 // prevent update this field
	err := ps.GetMaster().Model(transaction).Updates(transaction).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update transaction with id=%s", transaction.Id)
	}

	return transaction, nil
}

func (ps *SqlPaymentTransactionStore) Get(id string) (*model.PaymentTransaction, error) {
	var res model.PaymentTransaction
	err := ps.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.TransactionTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find payment transaction withh id=%s", id)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of transactions with given option
func (ps *SqlPaymentTransactionStore) FilterByOption(option *model.PaymentTransactionFilterOpts) ([]*model.PaymentTransaction, error) {
	var res []*model.PaymentTransaction
	err := ps.GetReplica().Find(&res, store.BuildSqlizer(option.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find payment transactions based on given option")
	}

	return res, nil
}
