package payment

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlPaymentTransactionStore struct {
	store.Store
}

func NewSqlPaymentTransactionStore(s store.Store) store.PaymentTransactionStore {
	return &SqlPaymentTransactionStore{s}
}

func (ps *SqlPaymentTransactionStore) Upsert(transaction boil.ContextTransactor, paymentTransaction model.PaymentTransaction) (*model.PaymentTransaction, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}
	isSaving := paymentTransaction.ID == ""

	if isSaving {
		model_helper.PaymentTransactionPreSave(&paymentTransaction)
	} else {
		model_helper.PaymentTransactionCommonPre(&paymentTransaction)
	}

	if err := model_helper.PaymentTransactionIsValid(paymentTransaction); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = paymentTransaction.Insert(transaction, boil.Infer())
	} else {
		_, err = paymentTransaction.Update(transaction, boil.Infer())
	}

	if err != nil {
		return nil, err
	}

	return &paymentTransaction, nil
}

func (ps *SqlPaymentTransactionStore) Get(id string) (*model.PaymentTransaction, error) {
	tran, err := model.FindPaymentTransaction(ps.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.PaymentTransactions, id)
		}
		return nil, err
	}

	return tran, nil
}

func (ps *SqlPaymentTransactionStore) FilterByOption(option model_helper.PaymentTransactionFilterOpts) ([]*model.PaymentTransaction, error) {
	return model.PaymentTransactions(option.Conditions...).All(ps.GetReplica())
}
