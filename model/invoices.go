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
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Invoice is an object representing the database table.
type Invoice struct {
	ID              string                 `boil:"id" json:"id" toml:"id" yaml:"id"`
	OrderID         model_types.NullString `boil:"order_id" json:"order_id,omitempty" toml:"order_id" yaml:"order_id,omitempty"`
	Number          string                 `boil:"number" json:"number" toml:"number" yaml:"number"`
	CreatedAt       int64                  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	ExternalURL     string                 `boil:"external_url" json:"external_url" toml:"external_url" yaml:"external_url"`
	Status          string                 `boil:"status" json:"status" toml:"status" yaml:"status"`
	Message         string                 `boil:"message" json:"message" toml:"message" yaml:"message"`
	UpdatedAt       int64                  `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	InvoiceFile     model_types.NullString `boil:"invoice_file" json:"invoice_file,omitempty" toml:"invoice_file" yaml:"invoice_file,omitempty"`
	Metadata        model_types.JSONString `boil:"metadata" json:"metadata,omitempty" toml:"metadata" yaml:"metadata,omitempty"`
	PrivateMetadata model_types.JSONString `boil:"private_metadata" json:"private_metadata,omitempty" toml:"private_metadata" yaml:"private_metadata,omitempty"`

	R *invoiceR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L invoiceL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var InvoiceColumns = struct {
	ID              string
	OrderID         string
	Number          string
	CreatedAt       string
	ExternalURL     string
	Status          string
	Message         string
	UpdatedAt       string
	InvoiceFile     string
	Metadata        string
	PrivateMetadata string
}{
	ID:              "id",
	OrderID:         "order_id",
	Number:          "number",
	CreatedAt:       "created_at",
	ExternalURL:     "external_url",
	Status:          "status",
	Message:         "message",
	UpdatedAt:       "updated_at",
	InvoiceFile:     "invoice_file",
	Metadata:        "metadata",
	PrivateMetadata: "private_metadata",
}

var InvoiceTableColumns = struct {
	ID              string
	OrderID         string
	Number          string
	CreatedAt       string
	ExternalURL     string
	Status          string
	Message         string
	UpdatedAt       string
	InvoiceFile     string
	Metadata        string
	PrivateMetadata string
}{
	ID:              "invoices.id",
	OrderID:         "invoices.order_id",
	Number:          "invoices.number",
	CreatedAt:       "invoices.created_at",
	ExternalURL:     "invoices.external_url",
	Status:          "invoices.status",
	Message:         "invoices.message",
	UpdatedAt:       "invoices.updated_at",
	InvoiceFile:     "invoices.invoice_file",
	Metadata:        "invoices.metadata",
	PrivateMetadata: "invoices.private_metadata",
}

// Generated where

var InvoiceWhere = struct {
	ID              whereHelperstring
	OrderID         whereHelpermodel_types_NullString
	Number          whereHelperstring
	CreatedAt       whereHelperint64
	ExternalURL     whereHelperstring
	Status          whereHelperstring
	Message         whereHelperstring
	UpdatedAt       whereHelperint64
	InvoiceFile     whereHelpermodel_types_NullString
	Metadata        whereHelpermodel_types_JSONString
	PrivateMetadata whereHelpermodel_types_JSONString
}{
	ID:              whereHelperstring{field: "\"invoices\".\"id\""},
	OrderID:         whereHelpermodel_types_NullString{field: "\"invoices\".\"order_id\""},
	Number:          whereHelperstring{field: "\"invoices\".\"number\""},
	CreatedAt:       whereHelperint64{field: "\"invoices\".\"created_at\""},
	ExternalURL:     whereHelperstring{field: "\"invoices\".\"external_url\""},
	Status:          whereHelperstring{field: "\"invoices\".\"status\""},
	Message:         whereHelperstring{field: "\"invoices\".\"message\""},
	UpdatedAt:       whereHelperint64{field: "\"invoices\".\"updated_at\""},
	InvoiceFile:     whereHelpermodel_types_NullString{field: "\"invoices\".\"invoice_file\""},
	Metadata:        whereHelpermodel_types_JSONString{field: "\"invoices\".\"metadata\""},
	PrivateMetadata: whereHelpermodel_types_JSONString{field: "\"invoices\".\"private_metadata\""},
}

// InvoiceRels is where relationship names are stored.
var InvoiceRels = struct {
	Order         string
	InvoiceEvents string
}{
	Order:         "Order",
	InvoiceEvents: "InvoiceEvents",
}

// invoiceR is where relationships are stored.
type invoiceR struct {
	Order         *Order            `boil:"Order" json:"Order" toml:"Order" yaml:"Order"`
	InvoiceEvents InvoiceEventSlice `boil:"InvoiceEvents" json:"InvoiceEvents" toml:"InvoiceEvents" yaml:"InvoiceEvents"`
}

// NewStruct creates a new relationship struct
func (*invoiceR) NewStruct() *invoiceR {
	return &invoiceR{}
}

func (r *invoiceR) GetOrder() *Order {
	if r == nil {
		return nil
	}
	return r.Order
}

func (r *invoiceR) GetInvoiceEvents() InvoiceEventSlice {
	if r == nil {
		return nil
	}
	return r.InvoiceEvents
}

// invoiceL is where Load methods for each relationship are stored.
type invoiceL struct{}

var (
	invoiceAllColumns            = []string{"id", "order_id", "number", "created_at", "external_url", "status", "message", "updated_at", "invoice_file", "metadata", "private_metadata"}
	invoiceColumnsWithoutDefault = []string{"id", "number", "created_at", "external_url", "status", "message", "updated_at"}
	invoiceColumnsWithDefault    = []string{"order_id", "invoice_file", "metadata", "private_metadata"}
	invoicePrimaryKeyColumns     = []string{"id"}
	invoiceGeneratedColumns      = []string{}
)

type (
	// InvoiceSlice is an alias for a slice of pointers to Invoice.
	// This should almost always be used instead of []Invoice.
	InvoiceSlice []*Invoice

	invoiceQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	invoiceType                 = reflect.TypeOf(&Invoice{})
	invoiceMapping              = queries.MakeStructMapping(invoiceType)
	invoicePrimaryKeyMapping, _ = queries.BindMapping(invoiceType, invoiceMapping, invoicePrimaryKeyColumns)
	invoiceInsertCacheMut       sync.RWMutex
	invoiceInsertCache          = make(map[string]insertCache)
	invoiceUpdateCacheMut       sync.RWMutex
	invoiceUpdateCache          = make(map[string]updateCache)
	invoiceUpsertCacheMut       sync.RWMutex
	invoiceUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single invoice record from the query.
func (q invoiceQuery) One(exec boil.Executor) (*Invoice, error) {
	o := &Invoice{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for invoices")
	}

	return o, nil
}

// All returns all Invoice records from the query.
func (q invoiceQuery) All(exec boil.Executor) (InvoiceSlice, error) {
	var o []*Invoice

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to Invoice slice")
	}

	return o, nil
}

// Count returns the count of all Invoice records in the query.
func (q invoiceQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count invoices rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q invoiceQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if invoices exists")
	}

	return count > 0, nil
}

// Order pointed to by the foreign key.
func (o *Invoice) Order(mods ...qm.QueryMod) orderQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.OrderID),
	}

	queryMods = append(queryMods, mods...)

	return Orders(queryMods...)
}

// InvoiceEvents retrieves all the invoice_event's InvoiceEvents with an executor.
func (o *Invoice) InvoiceEvents(mods ...qm.QueryMod) invoiceEventQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"invoice_events\".\"invoice_id\"=?", o.ID),
	)

	return InvoiceEvents(queryMods...)
}

// LoadOrder allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (invoiceL) LoadOrder(e boil.Executor, singular bool, maybeInvoice interface{}, mods queries.Applicator) error {
	var slice []*Invoice
	var object *Invoice

	if singular {
		var ok bool
		object, ok = maybeInvoice.(*Invoice)
		if !ok {
			object = new(Invoice)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeInvoice)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeInvoice))
			}
		}
	} else {
		s, ok := maybeInvoice.(*[]*Invoice)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeInvoice)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeInvoice))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &invoiceR{}
		}
		if !queries.IsNil(object.OrderID) {
			args[object.OrderID] = struct{}{}
		}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &invoiceR{}
			}

			if !queries.IsNil(obj.OrderID) {
				args[obj.OrderID] = struct{}{}
			}

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
		qm.From(`orders`),
		qm.WhereIn(`orders.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Order")
	}

	var resultSlice []*Order
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Order")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for orders")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for orders")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Order = foreign
		if foreign.R == nil {
			foreign.R = &orderR{}
		}
		foreign.R.Invoices = append(foreign.R.Invoices, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.OrderID, foreign.ID) {
				local.R.Order = foreign
				if foreign.R == nil {
					foreign.R = &orderR{}
				}
				foreign.R.Invoices = append(foreign.R.Invoices, local)
				break
			}
		}
	}

	return nil
}

// LoadInvoiceEvents allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (invoiceL) LoadInvoiceEvents(e boil.Executor, singular bool, maybeInvoice interface{}, mods queries.Applicator) error {
	var slice []*Invoice
	var object *Invoice

	if singular {
		var ok bool
		object, ok = maybeInvoice.(*Invoice)
		if !ok {
			object = new(Invoice)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeInvoice)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeInvoice))
			}
		}
	} else {
		s, ok := maybeInvoice.(*[]*Invoice)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeInvoice)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeInvoice))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &invoiceR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &invoiceR{}
			}
			args[obj.ID] = struct{}{}
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
		qm.From(`invoice_events`),
		qm.WhereIn(`invoice_events.invoice_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load invoice_events")
	}

	var resultSlice []*InvoiceEvent
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice invoice_events")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on invoice_events")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for invoice_events")
	}

	if singular {
		object.R.InvoiceEvents = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &invoiceEventR{}
			}
			foreign.R.Invoice = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.InvoiceID) {
				local.R.InvoiceEvents = append(local.R.InvoiceEvents, foreign)
				if foreign.R == nil {
					foreign.R = &invoiceEventR{}
				}
				foreign.R.Invoice = local
				break
			}
		}
	}

	return nil
}

// SetOrder of the invoice to the related item.
// Sets o.R.Order to related.
// Adds o to related.R.Invoices.
func (o *Invoice) SetOrder(exec boil.Executor, insert bool, related *Order) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"invoices\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"order_id"}),
		strmangle.WhereClause("\"", "\"", 2, invoicePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.OrderID, related.ID)
	if o.R == nil {
		o.R = &invoiceR{
			Order: related,
		}
	} else {
		o.R.Order = related
	}

	if related.R == nil {
		related.R = &orderR{
			Invoices: InvoiceSlice{o},
		}
	} else {
		related.R.Invoices = append(related.R.Invoices, o)
	}

	return nil
}

// RemoveOrder relationship.
// Sets o.R.Order to nil.
// Removes o from all passed in related items' relationships struct.
func (o *Invoice) RemoveOrder(exec boil.Executor, related *Order) error {
	var err error

	queries.SetScanner(&o.OrderID, nil)
	if _, err = o.Update(exec, boil.Whitelist("order_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.Order = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.Invoices {
		if queries.Equal(o.OrderID, ri.OrderID) {
			continue
		}

		ln := len(related.R.Invoices)
		if ln > 1 && i < ln-1 {
			related.R.Invoices[i] = related.R.Invoices[ln-1]
		}
		related.R.Invoices = related.R.Invoices[:ln-1]
		break
	}
	return nil
}

// AddInvoiceEvents adds the given related objects to the existing relationships
// of the invoice, optionally inserting them as new records.
// Appends related to o.R.InvoiceEvents.
// Sets related.R.Invoice appropriately.
func (o *Invoice) AddInvoiceEvents(exec boil.Executor, insert bool, related ...*InvoiceEvent) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.InvoiceID, o.ID)
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"invoice_events\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"invoice_id"}),
				strmangle.WhereClause("\"", "\"", 2, invoiceEventPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.InvoiceID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &invoiceR{
			InvoiceEvents: related,
		}
	} else {
		o.R.InvoiceEvents = append(o.R.InvoiceEvents, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &invoiceEventR{
				Invoice: o,
			}
		} else {
			rel.R.Invoice = o
		}
	}
	return nil
}

// SetInvoiceEvents removes all previously related items of the
// invoice replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Invoice's InvoiceEvents accordingly.
// Replaces o.R.InvoiceEvents with related.
// Sets related.R.Invoice's InvoiceEvents accordingly.
func (o *Invoice) SetInvoiceEvents(exec boil.Executor, insert bool, related ...*InvoiceEvent) error {
	query := "update \"invoice_events\" set \"invoice_id\" = null where \"invoice_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.InvoiceEvents {
			queries.SetScanner(&rel.InvoiceID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Invoice = nil
		}
		o.R.InvoiceEvents = nil
	}

	return o.AddInvoiceEvents(exec, insert, related...)
}

// RemoveInvoiceEvents relationships from objects passed in.
// Removes related items from R.InvoiceEvents (uses pointer comparison, removal does not keep order)
// Sets related.R.Invoice.
func (o *Invoice) RemoveInvoiceEvents(exec boil.Executor, related ...*InvoiceEvent) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.InvoiceID, nil)
		if rel.R != nil {
			rel.R.Invoice = nil
		}
		if _, err = rel.Update(exec, boil.Whitelist("invoice_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.InvoiceEvents {
			if rel != ri {
				continue
			}

			ln := len(o.R.InvoiceEvents)
			if ln > 1 && i < ln-1 {
				o.R.InvoiceEvents[i] = o.R.InvoiceEvents[ln-1]
			}
			o.R.InvoiceEvents = o.R.InvoiceEvents[:ln-1]
			break
		}
	}

	return nil
}

// Invoices retrieves all the records using an executor.
func Invoices(mods ...qm.QueryMod) invoiceQuery {
	mods = append(mods, qm.From("\"invoices\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"invoices\".*"})
	}

	return invoiceQuery{q}
}

// FindInvoice retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindInvoice(exec boil.Executor, iD string, selectCols ...string) (*Invoice, error) {
	invoiceObj := &Invoice{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"invoices\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, invoiceObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from invoices")
	}

	return invoiceObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Invoice) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no invoices provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(invoiceColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	invoiceInsertCacheMut.RLock()
	cache, cached := invoiceInsertCache[key]
	invoiceInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			invoiceAllColumns,
			invoiceColumnsWithDefault,
			invoiceColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(invoiceType, invoiceMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(invoiceType, invoiceMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"invoices\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"invoices\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into invoices")
	}

	if !cached {
		invoiceInsertCacheMut.Lock()
		invoiceInsertCache[key] = cache
		invoiceInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Invoice.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Invoice) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	invoiceUpdateCacheMut.RLock()
	cache, cached := invoiceUpdateCache[key]
	invoiceUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			invoiceAllColumns,
			invoicePrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update invoices, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"invoices\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, invoicePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(invoiceType, invoiceMapping, append(wl, invoicePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update invoices row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for invoices")
	}

	if !cached {
		invoiceUpdateCacheMut.Lock()
		invoiceUpdateCache[key] = cache
		invoiceUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q invoiceQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for invoices")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for invoices")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o InvoiceSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), invoicePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"invoices\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, invoicePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in invoice slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all invoice")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Invoice) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no invoices provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(invoiceColumnsWithDefault, o)

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

	invoiceUpsertCacheMut.RLock()
	cache, cached := invoiceUpsertCache[key]
	invoiceUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			invoiceAllColumns,
			invoiceColumnsWithDefault,
			invoiceColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			invoiceAllColumns,
			invoicePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert invoices, could not build update column list")
		}

		ret := strmangle.SetComplement(invoiceAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(invoicePrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert invoices, could not build conflict column list")
			}

			conflict = make([]string, len(invoicePrimaryKeyColumns))
			copy(conflict, invoicePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"invoices\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(invoiceType, invoiceMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(invoiceType, invoiceMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert invoices")
	}

	if !cached {
		invoiceUpsertCacheMut.Lock()
		invoiceUpsertCache[key] = cache
		invoiceUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Invoice record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Invoice) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no Invoice provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), invoicePrimaryKeyMapping)
	sql := "DELETE FROM \"invoices\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from invoices")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for invoices")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q invoiceQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no invoiceQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from invoices")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for invoices")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o InvoiceSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), invoicePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"invoices\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, invoicePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from invoice slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for invoices")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Invoice) Reload(exec boil.Executor) error {
	ret, err := FindInvoice(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *InvoiceSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := InvoiceSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), invoicePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"invoices\".* FROM \"invoices\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, invoicePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in InvoiceSlice")
	}

	*o = slice

	return nil
}

// InvoiceExists checks if the Invoice row exists.
func InvoiceExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"invoices\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if invoices exists")
	}

	return exists, nil
}

// Exists checks if the Invoice row exists.
func (o *Invoice) Exists(exec boil.Executor) (bool, error) {
	return InvoiceExists(exec, o.ID)
}
