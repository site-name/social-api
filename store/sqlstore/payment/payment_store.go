package payment

import (
	"fmt"

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

// CancelActivePaymentsOfCheckout inactivate all payments that belong to given checkout and in active status
func (ps *SqlPaymentStore) CancelActivePaymentsOfCheckout(checkoutID string) error {
	err := ps.GetMaster().Raw("UPDATE "+model.PaymentTableName+" SET IsActive = false WHERE CheckoutID = ? AND IsActive = true", checkoutID).Error
	if err != nil {
		return errors.Wrapf(err, "failed to deactivate payments that are active and belong to checkout with id=%s", checkoutID)
	}

	return nil
}

// FilterByOption finds and returns a list of payments that satisfy given option
func (ps *SqlPaymentStore) FilterByOption(option *model.PaymentFilterOption) (int64, []*model.Payment, error) {
	query := ps.GetQueryBuilder().
		Select(model.PaymentTableName + ".*").
		From(model.PaymentTableName).
		Where(option.Conditions)

	if option.LockForUpdate && option.DbTransaction != nil {
		query = query.Suffix("FOR UPDATE")
	}

	if option.RelatedTransactionConditions != nil {
		query = query.
			InnerJoin(
				fmt.Sprintf(
					"%[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.TransactionTableName,       //
					model.PaymentTableName,           // 2
					model.TransactionColumnPaymentID, // 3
					model.PaymentColumnId,            // 4
				),
			).
			Where(option.RelatedTransactionConditions)
	}

	// count if needed
	var totalCount int64
	if option.CountTotal {
		countQuery, args, err := ps.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "FilterByOption_CountTotal_ToSql")
		}

		err = ps.GetReplica().Raw(countQuery, args...).Scan(&totalCount).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total payments by options")
		}
	}

	// paginate if needed
	option.PaginationValues.AddPaginationToSelectBuilderIfNeeded(&query)

	queryString, args, err := query.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var payments []*model.Payment
	err = ps.GetReplica().Raw(queryString, args...).Scan(&payments).Error
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to finds payments with given option")
	}

	return totalCount, payments, nil
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
	query := fmt.Sprintf(
		`SELECT
			COUNT(*)
		FROM
			%[1]s
		INNER JOIN 
			%[2]s ON %[1]s.%[3]s = %[2]s.%[4]s
		INNER JOIN
			%[5]s ON %[5]s.%[6]s = %[1]s.%[7]s
		WHERE (
			%[2]s.%[8]s = $1
			OR %[5]s.%[9]s = $1
		)
		AND (
			%[1]s.%[10]s = $2
			OR %[1]s.%[11]s = $2
		)`,

		model.PaymentTableName,        // 1
		model.OrderTableName,          // 2
		model.PaymentColumnOrderID,    // 3
		model.OrderColumnId,           // 4
		model.CheckoutTableName,       // 5
		model.CheckoutColumnToken,     // 6
		model.PaymentColumnCheckoutID, // 7
		model.OrderColumnUserId,       // 8
		model.CheckoutColumnUserID,    // 9
		model.PaymentColumnId,         // 10
		model.PaymentColumnToken,      // 11
	)

	var paymentCount int64
	err := ps.GetReplica().Raw(query, userID, paymentID).Scan(&paymentCount).Error
	if err != nil {
		return false, errors.Wrap(err, "failed to find payments belong to given user")
	}

	return paymentCount > 0, nil
}
