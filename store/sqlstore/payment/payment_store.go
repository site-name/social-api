package payment

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

type SqlPaymentStore struct {
	store.Store
}

const (
	paymentTableName = "Payments"
)

func NewSqlPaymentStore(s store.Store) store.PaymentStore {
	ps := &SqlPaymentStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(payment.Payment{}, paymentTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CheckoutID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("GateWay").SetMaxSize(payment.MAX_LENGTH_PAYMENT_GATEWAY)
		table.ColMap("ChargeStatus").SetMaxSize(payment.MAX_LENGTH_PAYMENT_CHARGE_STATUS).
			SetDefaultConstraint(model.NewString(payment.NOT_CHARGED))
		table.ColMap("Token").SetMaxSize(payment.MAX_LENGTH_PAYMENT_TOKEN)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)
		table.ColMap("BillingEmail").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
		table.ColMap("BillingFirstName").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingLastName").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingCompanyName").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingAddress1").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingAddress2").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingCity").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingCityArea").SetMaxSize(account.ADDRESS_CITY_AREA_MAX_LENGTH)
		table.ColMap("BillingPostalCode").SetMaxSize(account.ADDRESS_POSTAL_CODE_MAX_LENGTH)
		table.ColMap("BillingCountryCode").SetMaxSize(model.SINGLE_COUNTRY_CODE_MAX_LENGTH)
		table.ColMap("BillingCountryArea").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)

		table.ColMap("CcFirstDigits").SetMaxSize(payment.MAX_LENGTH_CC_FIRST_DIGITS)
		table.ColMap("CcLastDigits").SetMaxSize(payment.MAX_LENGTH_CC_LAST_DIGITS)
		table.ColMap("CcBrand").SetMaxSize(payment.MAX_LENGTH_CC_BRAND)
		table.ColMap("PaymentMethodType").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("CustomerIpAddress").SetMaxSize(model.IP_ADDRESS_MAX_LENGTH)
		table.ColMap("ReturnUrl").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("PspReference").SetMaxSize(payment.PAYMENT_PSP_REFERENCE_MAX_LENGTH)
	}
	return ps
}

func (ps *SqlPaymentStore) CreateIndexesIfNotExists() {
	// NOTE: need more investigation in the future
	ps.CreateIndexIfNotExists("idx_payments_billing_email", paymentTableName, "BillingEmail")
	ps.CreateIndexIfNotExists("idx_payments_billing_first_name", paymentTableName, "BillingFirstName")
	ps.CreateIndexIfNotExists("idx_payments_billing_last_name", paymentTableName, "BillingLastName")
	ps.CreateIndexIfNotExists("idx_payments_billing_company_name", paymentTableName, "BillingCompanyName")
	ps.CreateIndexIfNotExists("idx_payments_billing_address_1", paymentTableName, "BillingAddress1")
	ps.CreateIndexIfNotExists("idx_payments_billing_city", paymentTableName, "BillingCity")
	ps.CreateIndexIfNotExists("idx_payments_billing_city_area", paymentTableName, "BillingCityArea")

	ps.CreateIndexIfNotExists("idx_payments_billing_email_lower_textpattern", paymentTableName, "lower(BillingEmail) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_payments_billing_first_name_lower_textpattern", paymentTableName, "lower(BillingFirstName) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_payments_billing_last_name_lower_textpattern", paymentTableName, "lower(BillingLastName) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_payments_billing_city_area_lower_textpattern", paymentTableName, "lower(BillingCityArea) text_pattern_ops")

	ps.CreateIndexIfNotExists("idx_payments_psp_reference", paymentTableName, "PspReference")
}

func (ps *SqlPaymentStore) Save(payment *payment.Payment) (*payment.Payment, error) {
	payment.PreSave()
	if err := payment.IsValid(); err != nil {
		return nil, err
	}
	if err := ps.GetMaster().Insert(payment); err != nil {
		return nil, errors.Wrapf(err, "failed to insert new payment with id=%s", payment.Id)
	}

	return payment, nil
}

func (ps *SqlPaymentStore) Get(id string) (*payment.Payment, error) {
	var payment payment.Payment
	err := ps.GetReplica().SelectOne(&payment, "SELECT * FROM "+paymentTableName+" WHERE Id = :id", map[string]interface{}{"id": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(paymentTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find payment with id=%s", id)
	}

	return &payment, nil
}

func (ps *SqlPaymentStore) GetPaymentsByOrderID(orderID string) ([]*payment.Payment, error) {
	var payments []*payment.Payment
	_, err := ps.GetReplica().Select(&payments, "SELECT * FROM "+paymentTableName+" WHERE OrderID = :orderID", map[string]interface{}{"orderID": orderID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(paymentTableName, "orderID="+orderID)
		}
		return nil, errors.Wrapf(err, "failed to find payments belong to order with Id=%s", orderID)
	}

	return payments, nil
}

func (ps *SqlPaymentStore) PaymentExistWithOptions(opts *payment.PaymentFilterOpts) (paymentExist bool, err error) {
	query := `SELECT *
		FROM
			Payments AS P
		INNER JOIN
			Transactions AS T
		ON (
			P.Id = T.PaymentID
		)
		WHERE (
			P.OrderId = :orderId
			AND P.IsActive = :isActive
			AND T.Kind = :kind
			AND T.ActionRequired = :actionRequired
			AND T.IsSuccess = :isSuccess
		)`
	var payments []*payment.Payment
	_, err = ps.GetReplica().Select(&payments, query, map[string]interface{}{
		"orderId":        opts.OrderID,
		"isActive":       opts.IsActive,
		"kind":           opts.Kind,
		"actionRequired": opts.ActionRequired,
		"isSuccess":      opts.IsSuccess,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// not found means does not exist
			return false, nil
		}
		// other errors mean system error
		return false, errors.Wrap(err, "failed to find transactions with given options")
	}

	if len(payments) == 0 {
		return false, nil
	}
	return true, nil
}
