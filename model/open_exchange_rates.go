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

// OpenExchangeRate is an object representing the database table.
type OpenExchangeRate struct {
	ID         string                  `boil:"id" json:"id" toml:"id" yaml:"id"`
	ToCurrency Currency                `boil:"to_currency" json:"to_currency" toml:"to_currency" yaml:"to_currency"`
	Rate       model_types.NullDecimal `boil:"rate" json:"rate,omitempty" toml:"rate" yaml:"rate,omitempty"`
	CreatedAt  int64                   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *openExchangeRateR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L openExchangeRateL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var OpenExchangeRateColumns = struct {
	ID         string
	ToCurrency string
	Rate       string
	CreatedAt  string
}{
	ID:         "id",
	ToCurrency: "to_currency",
	Rate:       "rate",
	CreatedAt:  "created_at",
}

var OpenExchangeRateTableColumns = struct {
	ID         string
	ToCurrency string
	Rate       string
	CreatedAt  string
}{
	ID:         "open_exchange_rates.id",
	ToCurrency: "open_exchange_rates.to_currency",
	Rate:       "open_exchange_rates.rate",
	CreatedAt:  "open_exchange_rates.created_at",
}

// Generated where

var OpenExchangeRateWhere = struct {
	ID         whereHelperstring
	ToCurrency whereHelperCurrency
	Rate       whereHelpermodel_types_NullDecimal
	CreatedAt  whereHelperint64
}{
	ID:         whereHelperstring{field: "\"open_exchange_rates\".\"id\""},
	ToCurrency: whereHelperCurrency{field: "\"open_exchange_rates\".\"to_currency\""},
	Rate:       whereHelpermodel_types_NullDecimal{field: "\"open_exchange_rates\".\"rate\""},
	CreatedAt:  whereHelperint64{field: "\"open_exchange_rates\".\"created_at\""},
}

// OpenExchangeRateRels is where relationship names are stored.
var OpenExchangeRateRels = struct {
}{}

// openExchangeRateR is where relationships are stored.
type openExchangeRateR struct {
}

// NewStruct creates a new relationship struct
func (*openExchangeRateR) NewStruct() *openExchangeRateR {
	return &openExchangeRateR{}
}

// openExchangeRateL is where Load methods for each relationship are stored.
type openExchangeRateL struct{}

var (
	openExchangeRateAllColumns            = []string{"id", "to_currency", "rate", "created_at"}
	openExchangeRateColumnsWithoutDefault = []string{"id", "to_currency", "created_at"}
	openExchangeRateColumnsWithDefault    = []string{"rate"}
	openExchangeRatePrimaryKeyColumns     = []string{"id"}
	openExchangeRateGeneratedColumns      = []string{}
)

type (
	// OpenExchangeRateSlice is an alias for a slice of pointers to OpenExchangeRate.
	// This should almost always be used instead of []OpenExchangeRate.
	OpenExchangeRateSlice []*OpenExchangeRate

	openExchangeRateQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	openExchangeRateType                 = reflect.TypeOf(&OpenExchangeRate{})
	openExchangeRateMapping              = queries.MakeStructMapping(openExchangeRateType)
	openExchangeRatePrimaryKeyMapping, _ = queries.BindMapping(openExchangeRateType, openExchangeRateMapping, openExchangeRatePrimaryKeyColumns)
	openExchangeRateInsertCacheMut       sync.RWMutex
	openExchangeRateInsertCache          = make(map[string]insertCache)
	openExchangeRateUpdateCacheMut       sync.RWMutex
	openExchangeRateUpdateCache          = make(map[string]updateCache)
	openExchangeRateUpsertCacheMut       sync.RWMutex
	openExchangeRateUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single openExchangeRate record from the query.
func (q openExchangeRateQuery) One(exec boil.Executor) (*OpenExchangeRate, error) {
	o := &OpenExchangeRate{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for open_exchange_rates")
	}

	return o, nil
}

// All returns all OpenExchangeRate records from the query.
func (q openExchangeRateQuery) All(exec boil.Executor) (OpenExchangeRateSlice, error) {
	var o []*OpenExchangeRate

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to OpenExchangeRate slice")
	}

	return o, nil
}

// Count returns the count of all OpenExchangeRate records in the query.
func (q openExchangeRateQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count open_exchange_rates rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q openExchangeRateQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if open_exchange_rates exists")
	}

	return count > 0, nil
}

// OpenExchangeRates retrieves all the records using an executor.
func OpenExchangeRates(mods ...qm.QueryMod) openExchangeRateQuery {
	mods = append(mods, qm.From("\"open_exchange_rates\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"open_exchange_rates\".*"})
	}

	return openExchangeRateQuery{q}
}

// FindOpenExchangeRate retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindOpenExchangeRate(exec boil.Executor, iD string, selectCols ...string) (*OpenExchangeRate, error) {
	openExchangeRateObj := &OpenExchangeRate{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"open_exchange_rates\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, openExchangeRateObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from open_exchange_rates")
	}

	return openExchangeRateObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *OpenExchangeRate) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no open_exchange_rates provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(openExchangeRateColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	openExchangeRateInsertCacheMut.RLock()
	cache, cached := openExchangeRateInsertCache[key]
	openExchangeRateInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			openExchangeRateAllColumns,
			openExchangeRateColumnsWithDefault,
			openExchangeRateColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(openExchangeRateType, openExchangeRateMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(openExchangeRateType, openExchangeRateMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"open_exchange_rates\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"open_exchange_rates\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into open_exchange_rates")
	}

	if !cached {
		openExchangeRateInsertCacheMut.Lock()
		openExchangeRateInsertCache[key] = cache
		openExchangeRateInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the OpenExchangeRate.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *OpenExchangeRate) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	openExchangeRateUpdateCacheMut.RLock()
	cache, cached := openExchangeRateUpdateCache[key]
	openExchangeRateUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			openExchangeRateAllColumns,
			openExchangeRatePrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update open_exchange_rates, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"open_exchange_rates\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, openExchangeRatePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(openExchangeRateType, openExchangeRateMapping, append(wl, openExchangeRatePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update open_exchange_rates row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for open_exchange_rates")
	}

	if !cached {
		openExchangeRateUpdateCacheMut.Lock()
		openExchangeRateUpdateCache[key] = cache
		openExchangeRateUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q openExchangeRateQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for open_exchange_rates")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for open_exchange_rates")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o OpenExchangeRateSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), openExchangeRatePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"open_exchange_rates\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, openExchangeRatePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in openExchangeRate slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all openExchangeRate")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *OpenExchangeRate) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no open_exchange_rates provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(openExchangeRateColumnsWithDefault, o)

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

	openExchangeRateUpsertCacheMut.RLock()
	cache, cached := openExchangeRateUpsertCache[key]
	openExchangeRateUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			openExchangeRateAllColumns,
			openExchangeRateColumnsWithDefault,
			openExchangeRateColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			openExchangeRateAllColumns,
			openExchangeRatePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert open_exchange_rates, could not build update column list")
		}

		ret := strmangle.SetComplement(openExchangeRateAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(openExchangeRatePrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert open_exchange_rates, could not build conflict column list")
			}

			conflict = make([]string, len(openExchangeRatePrimaryKeyColumns))
			copy(conflict, openExchangeRatePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"open_exchange_rates\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(openExchangeRateType, openExchangeRateMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(openExchangeRateType, openExchangeRateMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert open_exchange_rates")
	}

	if !cached {
		openExchangeRateUpsertCacheMut.Lock()
		openExchangeRateUpsertCache[key] = cache
		openExchangeRateUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single OpenExchangeRate record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *OpenExchangeRate) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no OpenExchangeRate provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), openExchangeRatePrimaryKeyMapping)
	sql := "DELETE FROM \"open_exchange_rates\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from open_exchange_rates")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for open_exchange_rates")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q openExchangeRateQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no openExchangeRateQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from open_exchange_rates")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for open_exchange_rates")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o OpenExchangeRateSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), openExchangeRatePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"open_exchange_rates\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, openExchangeRatePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from openExchangeRate slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for open_exchange_rates")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *OpenExchangeRate) Reload(exec boil.Executor) error {
	ret, err := FindOpenExchangeRate(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *OpenExchangeRateSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := OpenExchangeRateSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), openExchangeRatePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"open_exchange_rates\".* FROM \"open_exchange_rates\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, openExchangeRatePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in OpenExchangeRateSlice")
	}

	*o = slice

	return nil
}

// OpenExchangeRateExists checks if the OpenExchangeRate row exists.
func OpenExchangeRateExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"open_exchange_rates\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if open_exchange_rates exists")
	}

	return exists, nil
}

// Exists checks if the OpenExchangeRate row exists.
func (o *OpenExchangeRate) Exists(exec boil.Executor) (bool, error) {
	return OpenExchangeRateExists(exec, o.ID)
}
