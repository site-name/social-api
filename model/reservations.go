// Code generated by SQLBoiler 4.17.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// Reservation is an object representing the database table.
type Reservation struct {
	ID               string                `boil:"id" json:"id" toml:"id" yaml:"id"`
	CheckoutLineID   string                `boil:"checkout_line_id" json:"checkout_line_id" toml:"checkout_line_id" yaml:"checkout_line_id"`
	StockID          string                `boil:"stock_id" json:"stock_id" toml:"stock_id" yaml:"stock_id"`
	QuantityReserved int                   `boil:"quantity_reserved" json:"quantity_reserved" toml:"quantity_reserved" yaml:"quantity_reserved"`
	ReservedUntil    model_types.NullInt64 `boil:"reserved_until" json:"reserved_until,omitempty" toml:"reserved_until" yaml:"reserved_until,omitempty"`

	R *reservationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L reservationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ReservationColumns = struct {
	ID               string
	CheckoutLineID   string
	StockID          string
	QuantityReserved string
	ReservedUntil    string
}{
	ID:               "id",
	CheckoutLineID:   "checkout_line_id",
	StockID:          "stock_id",
	QuantityReserved: "quantity_reserved",
	ReservedUntil:    "reserved_until",
}

var ReservationTableColumns = struct {
	ID               string
	CheckoutLineID   string
	StockID          string
	QuantityReserved string
	ReservedUntil    string
}{
	ID:               "reservations.id",
	CheckoutLineID:   "reservations.checkout_line_id",
	StockID:          "reservations.stock_id",
	QuantityReserved: "reservations.quantity_reserved",
	ReservedUntil:    "reservations.reserved_until",
}

// Generated where

var ReservationWhere = struct {
	ID               whereHelperstring
	CheckoutLineID   whereHelperstring
	StockID          whereHelperstring
	QuantityReserved whereHelperint
	ReservedUntil    whereHelpermodel_types_NullInt64
}{
	ID:               whereHelperstring{field: "\"reservations\".\"id\""},
	CheckoutLineID:   whereHelperstring{field: "\"reservations\".\"checkout_line_id\""},
	StockID:          whereHelperstring{field: "\"reservations\".\"stock_id\""},
	QuantityReserved: whereHelperint{field: "\"reservations\".\"quantity_reserved\""},
	ReservedUntil:    whereHelpermodel_types_NullInt64{field: "\"reservations\".\"reserved_until\""},
}

// ReservationRels is where relationship names are stored.
var ReservationRels = struct {
}{}

// reservationR is where relationships are stored.
type reservationR struct {
}

// NewStruct creates a new relationship struct
func (*reservationR) NewStruct() *reservationR {
	return &reservationR{}
}

// reservationL is where Load methods for each relationship are stored.
type reservationL struct{}

var (
	reservationAllColumns            = []string{"id", "checkout_line_id", "stock_id", "quantity_reserved", "reserved_until"}
	reservationColumnsWithoutDefault = []string{"id", "checkout_line_id", "stock_id", "quantity_reserved"}
	reservationColumnsWithDefault    = []string{"reserved_until"}
	reservationPrimaryKeyColumns     = []string{"id"}
	reservationGeneratedColumns      = []string{}
)

type (
	// ReservationSlice is an alias for a slice of pointers to Reservation.
	// This should almost always be used instead of []Reservation.
	ReservationSlice []*Reservation

	reservationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	reservationType                 = reflect.TypeOf(&Reservation{})
	reservationMapping              = queries.MakeStructMapping(reservationType)
	reservationPrimaryKeyMapping, _ = queries.BindMapping(reservationType, reservationMapping, reservationPrimaryKeyColumns)
	reservationInsertCacheMut       sync.RWMutex
	reservationInsertCache          = make(map[string]insertCache)
	reservationUpdateCacheMut       sync.RWMutex
	reservationUpdateCache          = make(map[string]updateCache)
	reservationUpsertCacheMut       sync.RWMutex
	reservationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single reservation record from the query.
func (q reservationQuery) One(exec boil.Executor) (*Reservation, error) {
	o := &Reservation{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for reservations")
	}

	return o, nil
}

// All returns all Reservation records from the query.
func (q reservationQuery) All(exec boil.Executor) (ReservationSlice, error) {
	var o []*Reservation

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to Reservation slice")
	}

	return o, nil
}

// Count returns the count of all Reservation records in the query.
func (q reservationQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count reservations rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q reservationQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if reservations exists")
	}

	return count > 0, nil
}

// Reservations retrieves all the records using an executor.
func Reservations(mods ...qm.QueryMod) reservationQuery {
	mods = append(mods, qm.From("\"reservations\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"reservations\".*"})
	}

	return reservationQuery{q}
}

// FindReservation retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindReservation(exec boil.Executor, iD string, selectCols ...string) (*Reservation, error) {
	reservationObj := &Reservation{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"reservations\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, reservationObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from reservations")
	}

	return reservationObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Reservation) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no reservations provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(reservationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	reservationInsertCacheMut.RLock()
	cache, cached := reservationInsertCache[key]
	reservationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			reservationAllColumns,
			reservationColumnsWithDefault,
			reservationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(reservationType, reservationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(reservationType, reservationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"reservations\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"reservations\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into reservations")
	}

	if !cached {
		reservationInsertCacheMut.Lock()
		reservationInsertCache[key] = cache
		reservationInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Reservation.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Reservation) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	reservationUpdateCacheMut.RLock()
	cache, cached := reservationUpdateCache[key]
	reservationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			reservationAllColumns,
			reservationPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update reservations, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"reservations\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, reservationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(reservationType, reservationMapping, append(wl, reservationPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update reservations row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for reservations")
	}

	if !cached {
		reservationUpdateCacheMut.Lock()
		reservationUpdateCache[key] = cache
		reservationUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q reservationQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for reservations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for reservations")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ReservationSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reservationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"reservations\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, reservationPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in reservation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all reservation")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Reservation) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no reservations provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(reservationColumnsWithDefault, o)

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

	reservationUpsertCacheMut.RLock()
	cache, cached := reservationUpsertCache[key]
	reservationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			reservationAllColumns,
			reservationColumnsWithDefault,
			reservationColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			reservationAllColumns,
			reservationPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert reservations, could not build update column list")
		}

		ret := strmangle.SetComplement(reservationAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(reservationPrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert reservations, could not build conflict column list")
			}

			conflict = make([]string, len(reservationPrimaryKeyColumns))
			copy(conflict, reservationPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"reservations\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(reservationType, reservationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(reservationType, reservationMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert reservations")
	}

	if !cached {
		reservationUpsertCacheMut.Lock()
		reservationUpsertCache[key] = cache
		reservationUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Reservation record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Reservation) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no Reservation provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), reservationPrimaryKeyMapping)
	sql := "DELETE FROM \"reservations\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from reservations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for reservations")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q reservationQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no reservationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from reservations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for reservations")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ReservationSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reservationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"reservations\" WHERE " +
		strmangle.WhereInClause(string(dialect.LQ), string(dialect.RQ), 1, reservationPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from reservation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for reservations")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Reservation) Reload(exec boil.Executor) error {
	ret, err := FindReservation(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ReservationSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ReservationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reservationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"reservations\".* FROM \"reservations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, reservationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in ReservationSlice")
	}

	*o = slice

	return nil
}

// ReservationExists checks if the Reservation row exists.
func ReservationExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"reservations\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if reservations exists")
	}

	return exists, nil
}

// Exists checks if the Reservation row exists.
func (o *Reservation) Exists(exec boil.Executor) (bool, error) {
	return ReservationExists(exec, o.ID)
}