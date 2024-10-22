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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// PluginKeyValueStore is an object representing the database table.
type PluginKeyValueStore struct {
	PluginID string                `boil:"plugin_id" json:"plugin_id" toml:"plugin_id" yaml:"plugin_id"`
	PKey     string                `boil:"p_key" json:"p_key" toml:"p_key" yaml:"p_key"`
	PValue   null.Bytes            `boil:"p_value" json:"p_value,omitempty" toml:"p_value" yaml:"p_value,omitempty"`
	ExpireAt model_types.NullInt64 `boil:"expire_at" json:"expire_at,omitempty" toml:"expire_at" yaml:"expire_at,omitempty"`

	R *pluginKeyValueStoreR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L pluginKeyValueStoreL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PluginKeyValueStoreColumns = struct {
	PluginID string
	PKey     string
	PValue   string
	ExpireAt string
}{
	PluginID: "plugin_id",
	PKey:     "p_key",
	PValue:   "p_value",
	ExpireAt: "expire_at",
}

var PluginKeyValueStoreTableColumns = struct {
	PluginID string
	PKey     string
	PValue   string
	ExpireAt string
}{
	PluginID: "plugin_key_value_store.plugin_id",
	PKey:     "plugin_key_value_store.p_key",
	PValue:   "plugin_key_value_store.p_value",
	ExpireAt: "plugin_key_value_store.expire_at",
}

// Generated where

var PluginKeyValueStoreWhere = struct {
	PluginID whereHelperstring
	PKey     whereHelperstring
	PValue   whereHelpernull_Bytes
	ExpireAt whereHelpermodel_types_NullInt64
}{
	PluginID: whereHelperstring{field: "\"plugin_key_value_store\".\"plugin_id\""},
	PKey:     whereHelperstring{field: "\"plugin_key_value_store\".\"p_key\""},
	PValue:   whereHelpernull_Bytes{field: "\"plugin_key_value_store\".\"p_value\""},
	ExpireAt: whereHelpermodel_types_NullInt64{field: "\"plugin_key_value_store\".\"expire_at\""},
}

// PluginKeyValueStoreRels is where relationship names are stored.
var PluginKeyValueStoreRels = struct {
}{}

// pluginKeyValueStoreR is where relationships are stored.
type pluginKeyValueStoreR struct {
}

// NewStruct creates a new relationship struct
func (*pluginKeyValueStoreR) NewStruct() *pluginKeyValueStoreR {
	return &pluginKeyValueStoreR{}
}

// pluginKeyValueStoreL is where Load methods for each relationship are stored.
type pluginKeyValueStoreL struct{}

var (
	pluginKeyValueStoreAllColumns            = []string{"plugin_id", "p_key", "p_value", "expire_at"}
	pluginKeyValueStoreColumnsWithoutDefault = []string{"plugin_id", "p_key"}
	pluginKeyValueStoreColumnsWithDefault    = []string{"p_value", "expire_at"}
	pluginKeyValueStorePrimaryKeyColumns     = []string{"plugin_id"}
	pluginKeyValueStoreGeneratedColumns      = []string{}
)

type (
	// PluginKeyValueStoreSlice is an alias for a slice of pointers to PluginKeyValueStore.
	// This should almost always be used instead of []PluginKeyValueStore.
	PluginKeyValueStoreSlice []*PluginKeyValueStore

	pluginKeyValueStoreQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	pluginKeyValueStoreType                 = reflect.TypeOf(&PluginKeyValueStore{})
	pluginKeyValueStoreMapping              = queries.MakeStructMapping(pluginKeyValueStoreType)
	pluginKeyValueStorePrimaryKeyMapping, _ = queries.BindMapping(pluginKeyValueStoreType, pluginKeyValueStoreMapping, pluginKeyValueStorePrimaryKeyColumns)
	pluginKeyValueStoreInsertCacheMut       sync.RWMutex
	pluginKeyValueStoreInsertCache          = make(map[string]insertCache)
	pluginKeyValueStoreUpdateCacheMut       sync.RWMutex
	pluginKeyValueStoreUpdateCache          = make(map[string]updateCache)
	pluginKeyValueStoreUpsertCacheMut       sync.RWMutex
	pluginKeyValueStoreUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single pluginKeyValueStore record from the query.
func (q pluginKeyValueStoreQuery) One(exec boil.Executor) (*PluginKeyValueStore, error) {
	o := &PluginKeyValueStore{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for plugin_key_value_store")
	}

	return o, nil
}

// All returns all PluginKeyValueStore records from the query.
func (q pluginKeyValueStoreQuery) All(exec boil.Executor) (PluginKeyValueStoreSlice, error) {
	var o []*PluginKeyValueStore

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to PluginKeyValueStore slice")
	}

	return o, nil
}

// Count returns the count of all PluginKeyValueStore records in the query.
func (q pluginKeyValueStoreQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count plugin_key_value_store rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q pluginKeyValueStoreQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if plugin_key_value_store exists")
	}

	return count > 0, nil
}

// PluginKeyValueStores retrieves all the records using an executor.
func PluginKeyValueStores(mods ...qm.QueryMod) pluginKeyValueStoreQuery {
	mods = append(mods, qm.From("\"plugin_key_value_store\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"plugin_key_value_store\".*"})
	}

	return pluginKeyValueStoreQuery{q}
}

// FindPluginKeyValueStore retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPluginKeyValueStore(exec boil.Executor, pluginID string, selectCols ...string) (*PluginKeyValueStore, error) {
	pluginKeyValueStoreObj := &PluginKeyValueStore{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"plugin_key_value_store\" where \"plugin_id\"=$1", sel,
	)

	q := queries.Raw(query, pluginID)

	err := q.Bind(nil, exec, pluginKeyValueStoreObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from plugin_key_value_store")
	}

	return pluginKeyValueStoreObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PluginKeyValueStore) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no plugin_key_value_store provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(pluginKeyValueStoreColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	pluginKeyValueStoreInsertCacheMut.RLock()
	cache, cached := pluginKeyValueStoreInsertCache[key]
	pluginKeyValueStoreInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			pluginKeyValueStoreAllColumns,
			pluginKeyValueStoreColumnsWithDefault,
			pluginKeyValueStoreColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(pluginKeyValueStoreType, pluginKeyValueStoreMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(pluginKeyValueStoreType, pluginKeyValueStoreMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"plugin_key_value_store\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"plugin_key_value_store\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into plugin_key_value_store")
	}

	if !cached {
		pluginKeyValueStoreInsertCacheMut.Lock()
		pluginKeyValueStoreInsertCache[key] = cache
		pluginKeyValueStoreInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the PluginKeyValueStore.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PluginKeyValueStore) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	pluginKeyValueStoreUpdateCacheMut.RLock()
	cache, cached := pluginKeyValueStoreUpdateCache[key]
	pluginKeyValueStoreUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			pluginKeyValueStoreAllColumns,
			pluginKeyValueStorePrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update plugin_key_value_store, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"plugin_key_value_store\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, pluginKeyValueStorePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(pluginKeyValueStoreType, pluginKeyValueStoreMapping, append(wl, pluginKeyValueStorePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update plugin_key_value_store row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for plugin_key_value_store")
	}

	if !cached {
		pluginKeyValueStoreUpdateCacheMut.Lock()
		pluginKeyValueStoreUpdateCache[key] = cache
		pluginKeyValueStoreUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q pluginKeyValueStoreQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for plugin_key_value_store")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for plugin_key_value_store")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PluginKeyValueStoreSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pluginKeyValueStorePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"plugin_key_value_store\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, pluginKeyValueStorePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in pluginKeyValueStore slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all pluginKeyValueStore")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PluginKeyValueStore) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no plugin_key_value_store provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(pluginKeyValueStoreColumnsWithDefault, o)

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

	pluginKeyValueStoreUpsertCacheMut.RLock()
	cache, cached := pluginKeyValueStoreUpsertCache[key]
	pluginKeyValueStoreUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			pluginKeyValueStoreAllColumns,
			pluginKeyValueStoreColumnsWithDefault,
			pluginKeyValueStoreColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			pluginKeyValueStoreAllColumns,
			pluginKeyValueStorePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert plugin_key_value_store, could not build update column list")
		}

		ret := strmangle.SetComplement(pluginKeyValueStoreAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(pluginKeyValueStorePrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert plugin_key_value_store, could not build conflict column list")
			}

			conflict = make([]string, len(pluginKeyValueStorePrimaryKeyColumns))
			copy(conflict, pluginKeyValueStorePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"plugin_key_value_store\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(pluginKeyValueStoreType, pluginKeyValueStoreMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(pluginKeyValueStoreType, pluginKeyValueStoreMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert plugin_key_value_store")
	}

	if !cached {
		pluginKeyValueStoreUpsertCacheMut.Lock()
		pluginKeyValueStoreUpsertCache[key] = cache
		pluginKeyValueStoreUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single PluginKeyValueStore record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PluginKeyValueStore) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no PluginKeyValueStore provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), pluginKeyValueStorePrimaryKeyMapping)
	sql := "DELETE FROM \"plugin_key_value_store\" WHERE \"plugin_id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from plugin_key_value_store")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for plugin_key_value_store")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q pluginKeyValueStoreQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no pluginKeyValueStoreQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from plugin_key_value_store")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for plugin_key_value_store")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PluginKeyValueStoreSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pluginKeyValueStorePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"plugin_key_value_store\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, pluginKeyValueStorePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from pluginKeyValueStore slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for plugin_key_value_store")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PluginKeyValueStore) Reload(exec boil.Executor) error {
	ret, err := FindPluginKeyValueStore(exec, o.PluginID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PluginKeyValueStoreSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PluginKeyValueStoreSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pluginKeyValueStorePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"plugin_key_value_store\".* FROM \"plugin_key_value_store\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, pluginKeyValueStorePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in PluginKeyValueStoreSlice")
	}

	*o = slice

	return nil
}

// PluginKeyValueStoreExists checks if the PluginKeyValueStore row exists.
func PluginKeyValueStoreExists(exec boil.Executor, pluginID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"plugin_key_value_store\" where \"plugin_id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, pluginID)
	}
	row := exec.QueryRow(sql, pluginID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if plugin_key_value_store exists")
	}

	return exists, nil
}

// Exists checks if the PluginKeyValueStore row exists.
func (o *PluginKeyValueStore) Exists(exec boil.Executor) (bool, error) {
	return PluginKeyValueStoreExists(exec, o.PluginID)
}
