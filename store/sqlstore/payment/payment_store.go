package payment

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlPaymentStore struct {
	store.Store
}

func NewSqlPaymentStore(s store.Store) store.PaymentStore {
	return &SqlPaymentStore{s}
}

func (ps *SqlPaymentStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"GateWay",
		"IsActive",
		"ToConfirm",
		"CreateAt",
		"UpdateAt",
		"ChargeStatus",
		"Token",
		"Total",
		"CapturedAmount",
		"Currency",
		"CheckoutID",
		"OrderID",
		"BillingEmail",
		"BillingFirstName",
		"BillingLastName",
		"BillingCompanyName",
		"BillingAddress1",
		"BillingAddress2",
		"BillingCity",
		"BillingCityArea",
		"BillingPostalCode",
		"BillingCountryCode",
		"BillingCountryArea",
		"CcFirstDigits",
		"CcLastDigits",
		"CcBrand",
		"CcExpMonth",
		"CcExpYear",
		"PaymentMethodType",
		"CustomerIpAddress",
		"ExtraData",
		"ReturnUrl",
		"PspReference",
		"StorePaymentMethod",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
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
func (ps *SqlPaymentStore) Save(transaction store_iface.SqlxTxExecutor, payment *model.Payment) (*model.Payment, error) {
	var executor store_iface.SqlxExecutor = ps.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	payment.PreSave()
	if err := payment.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.PaymentTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
	if _, err := executor.NamedExec(query, payment); err != nil {
		return nil, errors.Wrapf(err, "failed to insert new payment with id=%s", payment.Id)
	}

	return payment, nil
}

// Update updates given payment and returns the updated value
func (ps *SqlPaymentStore) Update(transaction store_iface.SqlxTxExecutor, payment *model.Payment) (*model.Payment, error) {
	var executor store_iface.SqlxExecutor = ps.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	payment.PreUpdate()
	if err := payment.IsValid(); err != nil {
		return nil, err
	}

	query := "UPDATE " + store.PaymentTableName + " SET " + ps.
		ModelFields("").
		Map(func(_ int, s string) string {
			return s + "=:" + s
		}).
		Join(",") + " WHERE Id=:Id"

	result, err := executor.NamedExec(query, payment)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update payment with PaymentId=%s", payment.Id)
	}
	if numUpdated, _ := result.RowsAffected(); numUpdated > 1 {
		return nil, errors.Errorf("more than one payment updated: %d", numUpdated)
	}

	return payment, nil
}

// Get finds and returns the payment with given id
func (ps *SqlPaymentStore) Get(transaction store_iface.SqlxTxExecutor, id string, lockForUpdate bool) (*model.Payment, error) {
	var selector store_iface.SqlxExecutor = ps.GetReplicaX()
	if transaction != nil {
		selector = transaction
	}

	var (
		res          model.Payment
		forUpdateSql string
	)
	if lockForUpdate {
		forUpdateSql = " FOR UPDATE"
	}

	err := selector.Get(
		&res,
		"SELECT * FROM "+store.PaymentTableName+" WHERE Id = ?"+forUpdateSql,
		id,
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
	_, err := ps.GetMasterX().Exec("UPDATE "+store.PaymentTableName+" SET IsActive = false WHERE CheckoutID = ? AND IsActive = true", checkoutID)
	if err != nil {
		return errors.Wrapf(err, "failed to deactivate payments that are active and belong to checkout with id=%s", checkoutID)
	}

	return nil
}

// FilterByOption finds and returns a list of payments that satisfy given option
func (ps *SqlPaymentStore) FilterByOption(option *model.PaymentFilterOption) ([]*model.Payment, error) {
	query := ps.GetQueryBuilder().
		Select(ps.ModelFields(store.PaymentTableName + ".")...).
		From(store.PaymentTableName).
		OrderBy(store.TableOrderingMap[store.PaymentTableName])

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.OrderID != nil {
		query = query.Where(option.OrderID)
	}
	if option.CheckoutID != nil {
		query = query.Where(option.CheckoutID)
	}
	if option.IsActive != nil {
		query = query.Where(squirrel.Eq{"Payments.IsActive": *option.IsActive})
	}

	var joinedTransactionTable bool

	if option.TransactionsKind != nil {
		query = query.
			InnerJoin(store.TransactionTableName + " ON (Transactions.PaymentID = Payments.Id)").
			Where(option.TransactionsKind)

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

	var payments []*model.Payment
	err = ps.GetReplicaX().Select(&payments, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to finds payments with given option")
	}

	return payments, nil
}

// UpdatePaymentsOfCheckout updates payments of given checkout
func (ps *SqlPaymentStore) UpdatePaymentsOfCheckout(transaction store_iface.SqlxTxExecutor, checkoutToken string, option *model.PaymentPatch) error {
	var executor store_iface.SqlxExecutor = ps.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	query := ps.GetQueryBuilder().Update(store.PaymentTableName).Where("CheckoutID = ?", checkoutToken)

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

	_, err = executor.Exec(queryString, args...)
	if err != nil {
		return errors.Wrap(err, "failed to update payments of given checkout and options")
	}

	return nil
}
