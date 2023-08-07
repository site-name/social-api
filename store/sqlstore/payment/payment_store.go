package payment

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlPaymentStore struct {
	store.Store
}

func NewSqlPaymentStore(s store.Store) store.PaymentStore {
	return &SqlPaymentStore{s}
}

func (ps *SqlPaymentStore) ScanFields(payMent *model.Payment) []interface{} {
	return []interface{}{
		&payMent.Id,
		&payMent.GateWay,
		&payMent.IsActive,
		&payMent.ToConfirm,
		&payMent.CreateAt,
		&payMent.UpdateAt,
		&payMent.ChargeStatus,
		&payMent.Token,
		&payMent.Total,
		&payMent.CapturedAmount,
		&payMent.Currency,
		&payMent.CheckoutID,
		&payMent.OrderID,
		&payMent.BillingEmail,
		&payMent.BillingFirstName,
		&payMent.BillingLastName,
		&payMent.BillingCompanyName,
		&payMent.BillingAddress1,
		&payMent.BillingAddress2,
		&payMent.BillingCity,
		&payMent.BillingCityArea,
		&payMent.BillingPostalCode,
		&payMent.BillingCountryCode,
		&payMent.BillingCountryArea,
		&payMent.CcFirstDigits,
		&payMent.CcLastDigits,
		&payMent.CcBrand,
		&payMent.CcExpMonth,
		&payMent.CcExpYear,
		&payMent.PaymentMethodType,
		&payMent.CustomerIpAddress,
		&payMent.ExtraData,
		&payMent.ReturnUrl,
		&payMent.PspReference,
		&payMent.StorePaymentMethod,
		&payMent.Metadata,
		&payMent.PrivateMetadata,
	}
}

// Save inserts given payment into database then returns it
func (ps *SqlPaymentStore) Save(transaction *gorm.DB, payment *model.Payment) (*model.Payment, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	if err := transaction.Create(payment).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to insert new payment with id=%s", payment.Id)
	}

	return payment, nil
}

// Update updates given payment and returns the updated value
func (ps *SqlPaymentStore) Update(transaction *gorm.DB, payment *model.Payment) (*model.Payment, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	payment.CreateAt = 0 // prevent update

	err := transaction.Model(payment).Updates(payment).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update payment with PaymentId=%s", payment.Id)
	}

	return payment, nil
}

// Get finds and returns the payment with given id
func (ps *SqlPaymentStore) Get(transaction *gorm.DB, id string, lockForUpdate bool) (*model.Payment, error) {
	var (
		res          model.Payment
		forUpdateSql string
	)
	if lockForUpdate && transaction != nil {
		forUpdateSql = " FOR UPDATE"
	}

	if transaction == nil {
		transaction = ps.GetReplica()
	}
	err := transaction.Raw("SELECT * FROM "+model.PaymentTableName+" WHERE Id = ?"+forUpdateSql, id).Scan(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.PaymentTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find payment with id=%s", id)
	}

	return &res, nil
}

// CancelActivePaymentsOfCheckout inactivate all payments that belong to given checkout and in active status
func (ps *SqlPaymentStore) CancelActivePaymentsOfCheckout(checkoutID string) error {
	err := ps.GetMaster().Raw("UPDATE "+model.PaymentTableName+" SET IsActive = false WHERE CheckoutID = ? AND IsActive = true", checkoutID).Error
	if err != nil {
		return errors.Wrapf(err, "failed to deactivate payments that are active and belong to checkout with id=%s", checkoutID)
	}

	return nil
}

// FilterByOption finds and returns a list of payments that satisfy given option
func (ps *SqlPaymentStore) FilterByOption(option *model.PaymentFilterOption) ([]*model.Payment, error) {
	query := ps.GetQueryBuilder().
		Select(model.PaymentTableName + ".*").
		From(model.PaymentTableName).
		Where(option.Conditions)

	if option.TransactionsKind != nil ||
		option.TransactionsActionRequired != nil ||
		option.TransactionsIsSuccess != nil {
		andConds := squirrel.And{
			option.TransactionsKind,
			option.TransactionsActionRequired,
			option.TransactionsIsSuccess,
		}

		query = query.
			InnerJoin(model.TransactionTableName + " ON (Transactions.PaymentID = Payments.Id)").
			Where(andConds)

	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var payments []*model.Payment
	err = ps.GetReplica().Raw(queryString, args...).Scan(&payments).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to finds payments with given option")
	}

	return payments, nil
}

// UpdatePaymentsOfCheckout updates payments of given checkout
func (ps *SqlPaymentStore) UpdatePaymentsOfCheckout(transaction *gorm.DB, checkoutToken string, option *model.PaymentPatch) error {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	query := ps.GetQueryBuilder().Update(model.PaymentTableName).Where("CheckoutID = ?", checkoutToken)

	// parse option
	if model.IsValidEmail(option.BillingEmail) {
		query = query.Set("BillingEmail", option.BillingEmail)
	}
	if model.IsValidId(option.OrderID) {
		query = query.Set("OrderID", option.OrderID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "UpdatePaymentsOfCheckout_ToSql")
	}

	err = transaction.Raw(queryString, args...).Error
	if err != nil {
		return errors.Wrap(err, "failed to update payments of given checkout and options")
	}

	return nil
}

func (ps *SqlPaymentStore) PaymentOwnedByUser(userID, paymentID string) (bool, error) {
	query := `SELECT * FROM ` +
		model.PaymentTableName +
		` P INNER JOIN ` +
		model.OrderTableName +
		` O ON O.Id = P.OrderID INNER JOIN ` +
		model.CheckoutTableName +
		` C ON C.Id = P.CheckoutID WHERE (O.UserID = $1 OR C.UserID = $2) AND (P.Id = $3 OR P.Token = $4)`

	var payments []*model.Payment
	err := ps.GetReplica().Raw(query, userID, userID, paymentID, paymentID).Scan(&payments).Error
	if err != nil {
		return false, errors.Wrap(err, "failed to find payments belong to given user")
	}

	return len(payments) > 0, nil
}
