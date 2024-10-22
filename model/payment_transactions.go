// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// PaymentTransaction is an object representing the database table.
type PaymentTransaction struct {
	ID                 string                 `boil:"id" json:"id" toml:"id" yaml:"id"`
	CreatedAt          int64                  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	PaymentID          string                 `boil:"payment_id" json:"payment_id" toml:"payment_id" yaml:"payment_id"`
	Token              string                 `boil:"token" json:"token" toml:"token" yaml:"token"`
	Kind               TransactionKind        `boil:"kind" json:"kind" toml:"kind" yaml:"kind"`
	IsSuccess          bool                   `boil:"is_success" json:"is_success" toml:"is_success" yaml:"is_success"`
	ActionRequired     bool                   `boil:"action_required" json:"action_required" toml:"action_required" yaml:"action_required"`
	ActionRequiredData model_types.JSONString `boil:"action_required_data" json:"action_required_data" toml:"action_required_data" yaml:"action_required_data"`
	Currency           Currency               `boil:"currency" json:"currency" toml:"currency" yaml:"currency"`
	Amount             decimal.Decimal        `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Error              model_types.NullString `boil:"error" json:"error,omitempty" toml:"error" yaml:"error,omitempty"`
	CustomerID         model_types.NullString `boil:"customer_id" json:"customer_id,omitempty" toml:"customer_id" yaml:"customer_id,omitempty"`
	GatewayResponse    model_types.JSONString `boil:"gateway_response" json:"gateway_response" toml:"gateway_response" yaml:"gateway_response"`
	AlreadyProcessed   bool                   `boil:"already_processed" json:"already_processed" toml:"already_processed" yaml:"already_processed"`

	R *paymentTransactionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L paymentTransactionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PaymentTransactionColumns = struct {
	ID                 string
	CreatedAt          string
	PaymentID          string
	Token              string
	Kind               string
	IsSuccess          string
	ActionRequired     string
	ActionRequiredData string
	Currency           string
	Amount             string
	Error              string
	CustomerID         string
	GatewayResponse    string
	AlreadyProcessed   string
}{
	ID:                 "id",
	CreatedAt:          "created_at",
	PaymentID:          "payment_id",
	Token:              "token",
	Kind:               "kind",
	IsSuccess:          "is_success",
	ActionRequired:     "action_required",
	ActionRequiredData: "action_required_data",
	Currency:           "currency",
	Amount:             "amount",
	Error:              "error",
	CustomerID:         "customer_id",
	GatewayResponse:    "gateway_response",
	AlreadyProcessed:   "already_processed",
}

var PaymentTransactionTableColumns = struct {
	ID                 string
	CreatedAt          string
	PaymentID          string
	Token              string
	Kind               string
	IsSuccess          string
	ActionRequired     string
	ActionRequiredData string
	Currency           string
	Amount             string
	Error              string
	CustomerID         string
	GatewayResponse    string
	AlreadyProcessed   string
}{
	ID:                 "payment_transactions.id",
	CreatedAt:          "payment_transactions.created_at",
	PaymentID:          "payment_transactions.payment_id",
	Token:              "payment_transactions.token",
	Kind:               "payment_transactions.kind",
	IsSuccess:          "payment_transactions.is_success",
	ActionRequired:     "payment_transactions.action_required",
	ActionRequiredData: "payment_transactions.action_required_data",
	Currency:           "payment_transactions.currency",
	Amount:             "payment_transactions.amount",
	Error:              "payment_transactions.error",
	CustomerID:         "payment_transactions.customer_id",
	GatewayResponse:    "payment_transactions.gateway_response",
	AlreadyProcessed:   "payment_transactions.already_processed",
}

// Generated where

type whereHelperTransactionKind struct{ field string }

func (w whereHelperTransactionKind) EQ(x TransactionKind) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelperTransactionKind) NEQ(x TransactionKind) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperTransactionKind) LT(x TransactionKind) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelperTransactionKind) LTE(x TransactionKind) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperTransactionKind) GT(x TransactionKind) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelperTransactionKind) GTE(x TransactionKind) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelperTransactionKind) IN(slice []TransactionKind) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperTransactionKind) NIN(slice []TransactionKind) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var PaymentTransactionWhere = struct {
	ID                 whereHelperstring
	CreatedAt          whereHelperint64
	PaymentID          whereHelperstring
	Token              whereHelperstring
	Kind               whereHelperTransactionKind
	IsSuccess          whereHelperbool
	ActionRequired     whereHelperbool
	ActionRequiredData whereHelpermodel_types_JSONString
	Currency           whereHelperCurrency
	Amount             whereHelperdecimal_Decimal
	Error              whereHelpermodel_types_NullString
	CustomerID         whereHelpermodel_types_NullString
	GatewayResponse    whereHelpermodel_types_JSONString
	AlreadyProcessed   whereHelperbool
}{
	ID:                 whereHelperstring{field: "\"payment_transactions\".\"id\""},
	CreatedAt:          whereHelperint64{field: "\"payment_transactions\".\"created_at\""},
	PaymentID:          whereHelperstring{field: "\"payment_transactions\".\"payment_id\""},
	Token:              whereHelperstring{field: "\"payment_transactions\".\"token\""},
	Kind:               whereHelperTransactionKind{field: "\"payment_transactions\".\"kind\""},
	IsSuccess:          whereHelperbool{field: "\"payment_transactions\".\"is_success\""},
	ActionRequired:     whereHelperbool{field: "\"payment_transactions\".\"action_required\""},
	ActionRequiredData: whereHelpermodel_types_JSONString{field: "\"payment_transactions\".\"action_required_data\""},
	Currency:           whereHelperCurrency{field: "\"payment_transactions\".\"currency\""},
	Amount:             whereHelperdecimal_Decimal{field: "\"payment_transactions\".\"amount\""},
	Error:              whereHelpermodel_types_NullString{field: "\"payment_transactions\".\"error\""},
	CustomerID:         whereHelpermodel_types_NullString{field: "\"payment_transactions\".\"customer_id\""},
	GatewayResponse:    whereHelpermodel_types_JSONString{field: "\"payment_transactions\".\"gateway_response\""},
	AlreadyProcessed:   whereHelperbool{field: "\"payment_transactions\".\"already_processed\""},
}

// PaymentTransactionRels is where relationship names are stored.
var PaymentTransactionRels = struct {
	Payment string
}{
	Payment: "Payment",
}

// paymentTransactionR is where relationships are stored.
type paymentTransactionR struct {
	Payment *Payment `boil:"Payment" json:"Payment" toml:"Payment" yaml:"Payment"`
}

// NewStruct creates a new relationship struct
func (*paymentTransactionR) NewStruct() *paymentTransactionR {
	return &paymentTransactionR{}
}

func (r *paymentTransactionR) GetPayment() *Payment {
	if r == nil {
		return nil
	}
	return r.Payment
}

// paymentTransactionL is where Load methods for each relationship are stored.
type paymentTransactionL struct{}

var (
	paymentTransactionAllColumns            = []string{"id", "created_at", "payment_id", "token", "kind", "is_success", "action_required", "action_required_data", "currency", "amount", "error", "customer_id", "gateway_response", "already_processed"}
	paymentTransactionColumnsWithoutDefault = []string{"id", "created_at", "payment_id", "token", "kind", "is_success", "action_required", "action_required_data", "currency", "gateway_response", "already_processed"}
	paymentTransactionColumnsWithDefault    = []string{"amount", "error", "customer_id"}
	paymentTransactionPrimaryKeyColumns     = []string{"id"}
	paymentTransactionGeneratedColumns      = []string{}
)

type (
	// PaymentTransactionSlice is an alias for a slice of pointers to PaymentTransaction.
	// This should almost always be used instead of []PaymentTransaction.
	PaymentTransactionSlice []*PaymentTransaction

	paymentTransactionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	paymentTransactionType                 = reflect.TypeOf(&PaymentTransaction{})
	paymentTransactionMapping              = queries.MakeStructMapping(paymentTransactionType)
	paymentTransactionPrimaryKeyMapping, _ = queries.BindMapping(paymentTransactionType, paymentTransactionMapping, paymentTransactionPrimaryKeyColumns)
	paymentTransactionInsertCacheMut       sync.RWMutex
	paymentTransactionInsertCache          = make(map[string]insertCache)
	paymentTransactionUpdateCacheMut       sync.RWMutex
	paymentTransactionUpdateCache          = make(map[string]updateCache)
	paymentTransactionUpsertCacheMut       sync.RWMutex
	paymentTransactionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single paymentTransaction record from the query.
func (q paymentTransactionQuery) One(exec boil.Executor) (*PaymentTransaction, error) {
	o := &PaymentTransaction{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for payment_transactions")
	}

	return o, nil
}

// All returns all PaymentTransaction records from the query.
func (q paymentTransactionQuery) All(exec boil.Executor) (PaymentTransactionSlice, error) {
	var o []*PaymentTransaction

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to PaymentTransaction slice")
	}

	return o, nil
}

// Count returns the count of all PaymentTransaction records in the query.
func (q paymentTransactionQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count payment_transactions rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q paymentTransactionQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if payment_transactions exists")
	}

	return count > 0, nil
}

// Payment pointed to by the foreign key.
func (o *PaymentTransaction) Payment(mods ...qm.QueryMod) paymentQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.PaymentID),
	}

	queryMods = append(queryMods, mods...)

	return Payments(queryMods...)
}

// LoadPayment allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (paymentTransactionL) LoadPayment(e boil.Executor, singular bool, maybePaymentTransaction interface{}, mods queries.Applicator) error {
	var slice []*PaymentTransaction
	var object *PaymentTransaction

	if singular {
		var ok bool
		object, ok = maybePaymentTransaction.(*PaymentTransaction)
		if !ok {
			object = new(PaymentTransaction)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePaymentTransaction)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePaymentTransaction))
			}
		}
	} else {
		s, ok := maybePaymentTransaction.(*[]*PaymentTransaction)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePaymentTransaction)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePaymentTransaction))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &paymentTransactionR{}
		}
		args[object.PaymentID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &paymentTransactionR{}
			}

			args[obj.PaymentID] = struct{}{}

		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`payments`),
		qm.WhereIn(`payments.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Payment")
	}

	var resultSlice []*Payment
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Payment")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for payments")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for payments")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Payment = foreign
		if foreign.R == nil {
			foreign.R = &paymentR{}
		}
		foreign.R.PaymentTransactions = append(foreign.R.PaymentTransactions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.PaymentID == foreign.ID {
				local.R.Payment = foreign
				if foreign.R == nil {
					foreign.R = &paymentR{}
				}
				foreign.R.PaymentTransactions = append(foreign.R.PaymentTransactions, local)
				break
			}
		}
	}

	return nil
}

// SetPayment of the paymentTransaction to the related item.
// Sets o.R.Payment to related.
// Adds o to related.R.PaymentTransactions.
func (o *PaymentTransaction) SetPayment(exec boil.Executor, insert bool, related *Payment) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"payment_transactions\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"payment_id"}),
		strmangle.WhereClause("\"", "\"", 2, paymentTransactionPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.PaymentID = related.ID
	if o.R == nil {
		o.R = &paymentTransactionR{
			Payment: related,
		}
	} else {
		o.R.Payment = related
	}

	if related.R == nil {
		related.R = &paymentR{
			PaymentTransactions: PaymentTransactionSlice{o},
		}
	} else {
		related.R.PaymentTransactions = append(related.R.PaymentTransactions, o)
	}

	return nil
}

// PaymentTransactions retrieves all the records using an executor.
func PaymentTransactions(mods ...qm.QueryMod) paymentTransactionQuery {
	mods = append(mods, qm.From("\"payment_transactions\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"payment_transactions\".*"})
	}

	return paymentTransactionQuery{q}
}

// FindPaymentTransaction retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPaymentTransaction(exec boil.Executor, iD string, selectCols ...string) (*PaymentTransaction, error) {
	paymentTransactionObj := &PaymentTransaction{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"payment_transactions\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, paymentTransactionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from payment_transactions")
	}

	return paymentTransactionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PaymentTransaction) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no payment_transactions provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(paymentTransactionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	paymentTransactionInsertCacheMut.RLock()
	cache, cached := paymentTransactionInsertCache[key]
	paymentTransactionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			paymentTransactionAllColumns,
			paymentTransactionColumnsWithDefault,
			paymentTransactionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(paymentTransactionType, paymentTransactionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(paymentTransactionType, paymentTransactionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"payment_transactions\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"payment_transactions\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "model: unable to insert into payment_transactions")
	}

	if !cached {
		paymentTransactionInsertCacheMut.Lock()
		paymentTransactionInsertCache[key] = cache
		paymentTransactionInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the PaymentTransaction.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PaymentTransaction) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	paymentTransactionUpdateCacheMut.RLock()
	cache, cached := paymentTransactionUpdateCache[key]
	paymentTransactionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			paymentTransactionAllColumns,
			paymentTransactionPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update payment_transactions, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"payment_transactions\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, paymentTransactionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(paymentTransactionType, paymentTransactionMapping, append(wl, paymentTransactionPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	var result sql.Result
	result, err = exec.Exec(cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update payment_transactions row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for payment_transactions")
	}

	if !cached {
		paymentTransactionUpdateCacheMut.Lock()
		paymentTransactionUpdateCache[key] = cache
		paymentTransactionUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q paymentTransactionQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for payment_transactions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for payment_transactions")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PaymentTransactionSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("model: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), paymentTransactionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"payment_transactions\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, paymentTransactionPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in paymentTransaction slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all paymentTransaction")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PaymentTransaction) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no payment_transactions provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(paymentTransactionColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	paymentTransactionUpsertCacheMut.RLock()
	cache, cached := paymentTransactionUpsertCache[key]
	paymentTransactionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			paymentTransactionAllColumns,
			paymentTransactionColumnsWithDefault,
			paymentTransactionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			paymentTransactionAllColumns,
			paymentTransactionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert payment_transactions, could not build update column list")
		}

		ret := strmangle.SetComplement(paymentTransactionAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(paymentTransactionPrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert payment_transactions, could not build conflict column list")
			}

			conflict = make([]string, len(paymentTransactionPrimaryKeyColumns))
			copy(conflict, paymentTransactionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"payment_transactions\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(paymentTransactionType, paymentTransactionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(paymentTransactionType, paymentTransactionMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "model: unable to upsert payment_transactions")
	}

	if !cached {
		paymentTransactionUpsertCacheMut.Lock()
		paymentTransactionUpsertCache[key] = cache
		paymentTransactionUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single PaymentTransaction record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PaymentTransaction) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no PaymentTransaction provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), paymentTransactionPrimaryKeyMapping)
	sql := "DELETE FROM \"payment_transactions\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from payment_transactions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for payment_transactions")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q paymentTransactionQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no paymentTransactionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from payment_transactions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for payment_transactions")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PaymentTransactionSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), paymentTransactionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"payment_transactions\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, paymentTransactionPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from paymentTransaction slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for payment_transactions")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PaymentTransaction) Reload(exec boil.Executor) error {
	ret, err := FindPaymentTransaction(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PaymentTransactionSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PaymentTransactionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), paymentTransactionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"payment_transactions\".* FROM \"payment_transactions\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, paymentTransactionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in PaymentTransactionSlice")
	}

	*o = slice

	return nil
}

// PaymentTransactionExists checks if the PaymentTransaction row exists.
func PaymentTransactionExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"payment_transactions\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if payment_transactions exists")
	}

	return exists, nil
}

// Exists checks if the PaymentTransaction row exists.
func (o *PaymentTransaction) Exists(exec boil.Executor) (bool, error) {
	return PaymentTransactionExists(exec, o.ID)
}
