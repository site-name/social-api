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

func NewSqlPaymentStore(s store.Store) store.PaymentStore {
	ps := &SqlPaymentStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(payment.Payment{}, store.PaymentTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CheckoutID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("GateWay").SetMaxSize(payment.MAX_LENGTH_PAYMENT_GATEWAY)
		table.ColMap("ChargeStatus").SetMaxSize(payment.MAX_LENGTH_PAYMENT_CHARGE_STATUS)
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
	ps.CreateIndexIfNotExists("idx_payments_order_id", store.PaymentTableName, "OrderID")
	ps.CreateIndexIfNotExists("idx_payments_is_active", store.PaymentTableName, "IsActive")
	ps.CreateIndexIfNotExists("idx_payments_charge_status", store.PaymentTableName, "ChargeStatus")
	ps.CreateIndexIfNotExists("idx_payments_psp_reference", store.PaymentTableName, "PspReference")

	ps.CreateForeignKeyIfNotExists(store.PaymentTableName, "OrderID", "Orders", "Id", false)
	ps.CreateForeignKeyIfNotExists(store.PaymentTableName, "CheckoutID", "Checkouts", "Id", false)
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
	err := ps.GetReplica().SelectOne(&payment, "SELECT * FROM "+store.PaymentTableName+" WHERE Id = :id", map[string]interface{}{"id": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.PaymentTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find payment with id=%s", id)
	}

	return &payment, nil
}

func (ps *SqlPaymentStore) GetPaymentsByOrderID(orderID string) ([]*payment.Payment, error) {
	var payments []*payment.Payment
	_, err := ps.GetReplica().Select(&payments, "SELECT * FROM "+store.PaymentTableName+" WHERE OrderID = :orderID", map[string]interface{}{"orderID": orderID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.PaymentTableName, "orderID="+orderID)
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
	_, err = ps.GetReplica().Select(
		&payments,
		query,
		map[string]interface{}{
			"orderId":        opts.OrderID,
			"isActive":       opts.IsActive,
			"kind":           opts.Kind,
			"actionRequired": opts.ActionRequired,
			"isSuccess":      opts.IsSuccess,
		},
	)
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

func (ps *SqlPaymentStore) GetPaymentsByCheckoutID(checkoutID string) ([]*payment.Payment, error) {
	var payments []*payment.Payment
	_, err := ps.GetReplica().Select(&payments, "SELECT * FROM "+store.PaymentTableName+" WHERE CheckoutID = :CheckoutID", map[string]interface{}{"CheckoutID": checkoutID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.PaymentTableName, "checkoutID="+checkoutID)
		}
		return nil, errors.Wrapf(err, "failed to find payments for checkout with id=%s", checkoutID)
	}

	return payments, nil
}
