package payment

import (
	"fmt"

	"github.com/mattermost/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlPaymentStore struct {
	store.Store
}

func NewSqlPaymentStore(s store.Store) store.PaymentStore {
	return &SqlPaymentStore{s}
}

func (ps *SqlPaymentStore) Upsert(transaction boil.ContextTransactor, payment model.Payment) (*model.Payment, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	isSaving := payment.ID == ""
	if isSaving {
		model_helper.PaymentPreSave(&payment)
	} else {
		model_helper.PaymentPreUpdate(&payment)
	}

	if err := model_helper.PaymentIsValid(payment); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = payment.Insert(transaction, boil.Infer())
	} else {
		_, err = payment.Update(transaction, boil.Blacklist(model.PaymentColumns.CreatedAt))
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert payment")
	}

	return &payment, nil
}

func (ps *SqlPaymentStore) CancelActivePaymentsOfCheckout(checkoutID string) error {
	_, err := model.Payments(
		model.PaymentWhere.CheckoutID.EQ(model_types.NewNullString(checkoutID)),
		model.PaymentWhere.IsActive.EQ(true),
	).UpdateAll(
		ps.GetMaster(),
		model.M{
			model.PaymentColumns.IsActive:  false,
			model.PaymentColumns.UpdatedAt: model_helper.GetMillis(),
		},
	)
	return err
}

func (ps *SqlPaymentStore) FilterByOption(option model_helper.PaymentFilterOptions) (model.PaymentSlice, error) {
	conds := option.Conditions
	if option.TransactionCondition != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.PaymentTransactions, model.PaymentTransactionTableColumns.PaymentID, model.PaymentTableColumns.ID)),
			option.TransactionCondition,
		)
	}

	return model.Payments(conds...).All(ps.GetReplica())
}

func (ps *SqlPaymentStore) UpdatePaymentsOfCheckout(transaction boil.ContextTransactor, checkoutToken string, option model_helper.PaymentPatch) error {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	updateCols := model.M{
		model.PaymentColumns.UpdatedAt: model_helper.GetMillis(),
	}
	// parse option
	if model_helper.IsValidEmail(option.BillingEmail) {
		updateCols[model.PaymentColumns.BillingEmail] = option.BillingEmail
	}
	if model_helper.IsValidId(option.OrderID) {
		updateCols[model.PaymentColumns.OrderID] = option.OrderID
	}
	_, err := model.Payments(model.PaymentWhere.CheckoutID.EQ(model_types.NewNullString(checkoutToken))).UpdateAll(transaction, updateCols)
	return err
}

func (ps *SqlPaymentStore) PaymentOwnedByUser(userID, paymentID string) (bool, error) {
	return model.Payments(
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Orders, model.PaymentTableColumns.OrderID, model.OrderTableColumns.ID)),
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Checkouts, model.CheckoutTableColumns.Token, model.PaymentTableColumns.CheckoutID)),
		model_helper.Or{
			squirrel.Eq{model.OrderTableColumns.UserID: userID},
			squirrel.Eq{model.CheckoutTableColumns.UserID: userID},
		},
		model_helper.Or{
			squirrel.Eq{model.PaymentTableColumns.ID: paymentID},
			squirrel.Eq{model.PaymentTableColumns.Token: paymentID},
		},
	).Exists(ps.GetReplica())
}
