package payment

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
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

	ps.CreateForeignKeyIfNotExists(store.PaymentTableName, "OrderID", store.OrderTableName, "Id", false)
	ps.CreateForeignKeyIfNotExists(store.PaymentTableName, "CheckoutID", store.CheckoutTableName, "Token", false)
}

func (ps *SqlPaymentStore) ModelFields() []string {
	return []string{
		"Payments.Id",
		"Payments.GateWay",
		"Payments.IsActive",
		"Payments.ToConfirm",
		"Payments.CreateAt",
		"Payments.UpdateAt",
		"Payments.ChargeStatus",
		"Payments.Token",
		"Payments.Total",
		"Payments.CapturedAmount",
		"Payments.Currency",
		"Payments.CheckoutID",
		"Payments.OrderID",
		"Payments.BillingEmail",
		"Payments.BillingFirstName",
		"Payments.BillingLastName",
		"Payments.BillingCompanyName",
		"Payments.BillingAddress1",
		"Payments.BillingAddress2",
		"Payments.BillingCity",
		"Payments.BillingCityArea",
		"Payments.BillingPostalCode",
		"Payments.BillingCountryCode",
		"Payments.BillingCountryArea",
		"Payments.CcFirstDigits",
		"Payments.CcLastDigits",
		"Payments.CcBrand",
		"Payments.CcExpMonth",
		"Payments.CcExpYear",
		"Payments.PaymentMethodType",
		"Payments.CustomerIpAddress",
		"Payments.ExtraData",
		"Payments.ReturnUrl",
		"Payments.PspReference",
	}
}

// Save inserts given payment into database then returns it
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

// Update updates given payment and returns the updated value
func (ps *SqlPaymentStore) Update(payment *payment.Payment) (*payment.Payment, error) {
	payment.PreUpdate()
	if err := payment.IsValid(); err != nil {
		return nil, err
	}

	oldPayment, err := ps.Get(payment.Id, false)
	if err != nil {
		return nil, err
	}

	payment.CreateAt = oldPayment.CreateAt
	payment.OrderID = oldPayment.OrderID
	payment.CheckoutID = oldPayment.CheckoutID

	numUpdated, err := ps.GetMaster().Update(payment)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update payment with PaymentId=%s", payment.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("more than one payment updated: %d", numUpdated)
	}

	return payment, nil
}

// Get finds and returns the payment with given id
func (ps *SqlPaymentStore) Get(id string, lockForUpdate bool) (*payment.Payment, error) {
	var res payment.Payment
	var forUpdateSql string
	if lockForUpdate {
		forUpdateSql = " FOR UPDATE"
	}
	err := ps.GetReplica().
		SelectOne(
			&res,
			"SELECT * FROM "+store.PaymentTableName+" WHERE Id :ID"+forUpdateSql,
			map[string]interface{}{"ID": id},
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.PaymentTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find payment with id=%s", id)
	}

	return &res, nil
}

// CancelActivePaymentsOfCheckout inactivate all payments that belong to given checkout and in active status
func (ps *SqlPaymentStore) CancelActivePaymentsOfCheckout(checkoutID string) error {
	_, err := ps.GetQueryBuilder().
		Update(store.PaymentTableName).
		Set("IsActive", false).
		Where(squirrel.And{
			squirrel.Eq{"CheckoutID": checkoutID},
			squirrel.Eq{"IsActive": true},
		}).
		RunWith(ps.GetMaster()).
		Exec()

	if err != nil {
		return errors.Wrapf(err, "failed to deactivate payments that are active and belong to checkout with id=%s", checkoutID)
	}

	return nil
}

// FilterByOption finds and returns a list of payments that satisfy given option
func (ps *SqlPaymentStore) FilterByOption(option *payment.PaymentFilterOption) ([]*payment.Payment, error) {
	query := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		Distinct().
		From(store.PaymentTableName).
		OrderBy(store.TableOrderingMap[store.PaymentTableName])

	var joinedTransactionTable bool

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Payments.Id"))
	}
	if model.IsValidId(option.OrderID) {
		query = query.Where(squirrel.Eq{"Payments.OrderID": option.OrderID})
	}
	if model.IsValidId(option.CheckoutToken) {
		query = query.Where(squirrel.Eq{"Payments.CheckoutID": option.CheckoutToken})
	}
	if option.IsActive != nil {
		query = query.Where(squirrel.Eq{"Payments.IsActive": *option.IsActive})
	}
	if option.TransactionsKind != nil {
		query = query.
			LeftJoin(store.TransactionTableName + " ON (Transactions.PaymentID = Payments.Id)").
			Where(option.TransactionsKind.ToSquirrel("Transactions.Kind"))

		// let later checks know that this query has already joined transaction table
		joinedTransactionTable = true
	}
	if option.TransactionsActionRequired != nil {
		// check if already joined table transactions
		if !joinedTransactionTable {
			query = query.
				LeftJoin(store.TransactionTableName + " ON (Transactions.PaymentID = Payments.Id)")
		}
		query = query.Where(squirrel.Eq{"Transactions.ActionRequired": *option.TransactionsActionRequired})
	}
	if option.TransactionsIsSuccess != nil {
		// check if already joined table transactions
		if !joinedTransactionTable {
			query = query.
				LeftJoin(store.TransactionTableName + " ON (Transactions.PaymentID = Payments.Id)")
		}
		query = query.Where(squirrel.Eq{"Transactions.IsSuccess": *option.TransactionsIsSuccess})
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var payments []*payment.Payment
	_, err = ps.GetReplica().Select(&payments, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to finds payments with given option")
	}

	return payments, nil
}
