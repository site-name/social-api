package payment

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
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
		table.ColMap("StorePaymentMethod").SetMaxSize(payment.PAYMENT_STORE_PAYMENT_METHOD_MAX_LENGTH)
	}
	return ps
}

func (ps *SqlPaymentStore) CreateIndexesIfNotExists() {
	ps.CommonMetaDataIndex(store.PaymentTableName)

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
		"Payments.StorePaymentMethod",
		"Payments.Metadata",
		"Payments.PrivateMetadata",
	}
}

func (ps *SqlPaymentStore) ScanFields(payMent payment.Payment) []interface{} {
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
func (ps *SqlPaymentStore) Save(transaction *gorp.Transaction, payment *payment.Payment) (*payment.Payment, error) {

	var upsertor store.Upsertor = ps.GetMaster()
	if transaction != nil {
		upsertor = transaction
	}

	payment.PreSave()
	if err := payment.IsValid(); err != nil {
		return nil, err
	}

	if err := upsertor.Insert(payment); err != nil {
		return nil, errors.Wrapf(err, "failed to insert new payment with id=%s", payment.Id)
	}

	return payment, nil
}

// Update updates given payment and returns the updated value
func (ps *SqlPaymentStore) Update(transaction *gorp.Transaction, payment *payment.Payment) (*payment.Payment, error) {
	var upsertor store.Upsertor = ps.GetMaster()
	if transaction != nil {
		upsertor = transaction
	}

	payment.PreUpdate()
	if err := payment.IsValid(); err != nil {
		return nil, err
	}

	oldPayment, err := ps.Get(transaction, payment.Id, false)
	if err != nil {
		return nil, err
	}

	payment.CreateAt = oldPayment.CreateAt
	payment.OrderID = oldPayment.OrderID
	payment.CheckoutID = oldPayment.CheckoutID

	numUpdated, err := upsertor.Update(payment)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update payment with PaymentId=%s", payment.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("more than one payment updated: %d", numUpdated)
	}

	return payment, nil
}

// Get finds and returns the payment with given id
func (ps *SqlPaymentStore) Get(transaction *gorp.Transaction, id string, lockForUpdate bool) (*payment.Payment, error) {
	var selector store.Selector = ps.GetReplica()
	if transaction != nil {
		selector = transaction
	}

	var (
		res          payment.Payment
		forUpdateSql string
	)
	if lockForUpdate {
		forUpdateSql = " FOR UPDATE"
	}

	err := selector.SelectOne(
		&res,
		"SELECT * FROM "+store.PaymentTableName+" WHERE Id = :ID"+forUpdateSql,
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
		From(store.PaymentTableName).
		OrderBy(store.TableOrderingMap[store.PaymentTableName])

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

	var joinedTransactionTable bool

	if option.TransactionsKind != nil {
		query = query.
			InnerJoin(store.TransactionTableName + " ON (Transactions.PaymentID = Payments.Id)").
			Where(option.TransactionsKind.ToSquirrel("Transactions.Kind"))

		joinedTransactionTable = true // indicate that we have joined transaction table
	}
	if option.TransactionsActionRequired != nil {
		if !joinedTransactionTable {
			query = query.InnerJoin(store.TransactionTableName + " ON (Transactions.PaymentID = Payments.Id)")
		}
		query = query.Where(squirrel.Eq{"Transactions.ActionRequired": *option.TransactionsActionRequired})
	}
	if option.TransactionsIsSuccess != nil {
		if !joinedTransactionTable {
			query = query.InnerJoin(store.TransactionTableName + " ON (Transactions.PaymentID = Payments.Id)")
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

// UpdatePaymentsOfCheckout updates payments of given checkout
func (ps *SqlPaymentStore) UpdatePaymentsOfCheckout(transaction *gorp.Transaction, checkoutToken string, option *payment.PaymentPatch) error {
	var executor squirrel.BaseRunner = ps.GetMaster()
	if transaction != nil {
		executor = transaction
	}

	query := ps.GetQueryBuilder().Update(store.PaymentTableName).Where(squirrel.Expr("CheckoutID = ?", checkoutToken))

	// parse option
	if model.IsValidEmail(option.BillingEmail) {
		query = query.Set("BillingEmail", option.BillingEmail)
	}
	if model.IsValidId(option.OrderID) {
		query = query.Set("OrderID", option.OrderID)
	}

	_, err := query.RunWith(executor).Exec()
	if err != nil {
		return errors.Wrap(err, "failed to update payments of given checkout and options")
	}

	return nil
}
