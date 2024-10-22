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

// Session is an object representing the database table.
type Session struct {
	ID             string                 `boil:"id" json:"id" toml:"id" yaml:"id"`
	Token          string                 `boil:"token" json:"token" toml:"token" yaml:"token"`
	CreatedAt      int64                  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	ExpiresAt      int64                  `boil:"expires_at" json:"expires_at" toml:"expires_at" yaml:"expires_at"`
	LastActivityAt int64                  `boil:"last_activity_at" json:"last_activity_at" toml:"last_activity_at" yaml:"last_activity_at"`
	UserID         string                 `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	DeviceID       string                 `boil:"device_id" json:"device_id" toml:"device_id" yaml:"device_id"`
	Roles          string                 `boil:"roles" json:"roles" toml:"roles" yaml:"roles"`
	IsOauth        bool                   `boil:"is_oauth" json:"is_oauth" toml:"is_oauth" yaml:"is_oauth"`
	ExpiredNotify  bool                   `boil:"expired_notify" json:"expired_notify" toml:"expired_notify" yaml:"expired_notify"`
	Props          model_types.JSONString `boil:"props" json:"props,omitempty" toml:"props" yaml:"props,omitempty"`
	Local          bool                   `boil:"local" json:"local" toml:"local" yaml:"local"`

	R *sessionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L sessionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SessionColumns = struct {
	ID             string
	Token          string
	CreatedAt      string
	ExpiresAt      string
	LastActivityAt string
	UserID         string
	DeviceID       string
	Roles          string
	IsOauth        string
	ExpiredNotify  string
	Props          string
	Local          string
}{
	ID:             "id",
	Token:          "token",
	CreatedAt:      "created_at",
	ExpiresAt:      "expires_at",
	LastActivityAt: "last_activity_at",
	UserID:         "user_id",
	DeviceID:       "device_id",
	Roles:          "roles",
	IsOauth:        "is_oauth",
	ExpiredNotify:  "expired_notify",
	Props:          "props",
	Local:          "local",
}

var SessionTableColumns = struct {
	ID             string
	Token          string
	CreatedAt      string
	ExpiresAt      string
	LastActivityAt string
	UserID         string
	DeviceID       string
	Roles          string
	IsOauth        string
	ExpiredNotify  string
	Props          string
	Local          string
}{
	ID:             "sessions.id",
	Token:          "sessions.token",
	CreatedAt:      "sessions.created_at",
	ExpiresAt:      "sessions.expires_at",
	LastActivityAt: "sessions.last_activity_at",
	UserID:         "sessions.user_id",
	DeviceID:       "sessions.device_id",
	Roles:          "sessions.roles",
	IsOauth:        "sessions.is_oauth",
	ExpiredNotify:  "sessions.expired_notify",
	Props:          "sessions.props",
	Local:          "sessions.local",
}

// Generated where

var SessionWhere = struct {
	ID             whereHelperstring
	Token          whereHelperstring
	CreatedAt      whereHelperint64
	ExpiresAt      whereHelperint64
	LastActivityAt whereHelperint64
	UserID         whereHelperstring
	DeviceID       whereHelperstring
	Roles          whereHelperstring
	IsOauth        whereHelperbool
	ExpiredNotify  whereHelperbool
	Props          whereHelpermodel_types_JSONString
	Local          whereHelperbool
}{
	ID:             whereHelperstring{field: "\"sessions\".\"id\""},
	Token:          whereHelperstring{field: "\"sessions\".\"token\""},
	CreatedAt:      whereHelperint64{field: "\"sessions\".\"created_at\""},
	ExpiresAt:      whereHelperint64{field: "\"sessions\".\"expires_at\""},
	LastActivityAt: whereHelperint64{field: "\"sessions\".\"last_activity_at\""},
	UserID:         whereHelperstring{field: "\"sessions\".\"user_id\""},
	DeviceID:       whereHelperstring{field: "\"sessions\".\"device_id\""},
	Roles:          whereHelperstring{field: "\"sessions\".\"roles\""},
	IsOauth:        whereHelperbool{field: "\"sessions\".\"is_oauth\""},
	ExpiredNotify:  whereHelperbool{field: "\"sessions\".\"expired_notify\""},
	Props:          whereHelpermodel_types_JSONString{field: "\"sessions\".\"props\""},
	Local:          whereHelperbool{field: "\"sessions\".\"local\""},
}

// SessionRels is where relationship names are stored.
var SessionRels = struct {
}{}

// sessionR is where relationships are stored.
type sessionR struct {
}

// NewStruct creates a new relationship struct
func (*sessionR) NewStruct() *sessionR {
	return &sessionR{}
}

// sessionL is where Load methods for each relationship are stored.
type sessionL struct{}

var (
	sessionAllColumns            = []string{"id", "token", "created_at", "expires_at", "last_activity_at", "user_id", "device_id", "roles", "is_oauth", "expired_notify", "props", "local"}
	sessionColumnsWithoutDefault = []string{"id", "token", "created_at", "expires_at", "last_activity_at", "user_id", "device_id", "roles", "is_oauth", "expired_notify", "local"}
	sessionColumnsWithDefault    = []string{"props"}
	sessionPrimaryKeyColumns     = []string{"id"}
	sessionGeneratedColumns      = []string{}
)

type (
	// SessionSlice is an alias for a slice of pointers to Session.
	// This should almost always be used instead of []Session.
	SessionSlice []*Session

	sessionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	sessionType                 = reflect.TypeOf(&Session{})
	sessionMapping              = queries.MakeStructMapping(sessionType)
	sessionPrimaryKeyMapping, _ = queries.BindMapping(sessionType, sessionMapping, sessionPrimaryKeyColumns)
	sessionInsertCacheMut       sync.RWMutex
	sessionInsertCache          = make(map[string]insertCache)
	sessionUpdateCacheMut       sync.RWMutex
	sessionUpdateCache          = make(map[string]updateCache)
	sessionUpsertCacheMut       sync.RWMutex
	sessionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single session record from the query.
func (q sessionQuery) One(exec boil.Executor) (*Session, error) {
	o := &Session{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for sessions")
	}

	return o, nil
}

// All returns all Session records from the query.
func (q sessionQuery) All(exec boil.Executor) (SessionSlice, error) {
	var o []*Session

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to Session slice")
	}

	return o, nil
}

// Count returns the count of all Session records in the query.
func (q sessionQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count sessions rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q sessionQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if sessions exists")
	}

	return count > 0, nil
}

// Sessions retrieves all the records using an executor.
func Sessions(mods ...qm.QueryMod) sessionQuery {
	mods = append(mods, qm.From("\"sessions\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"sessions\".*"})
	}

	return sessionQuery{q}
}

// FindSession retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSession(exec boil.Executor, iD string, selectCols ...string) (*Session, error) {
	sessionObj := &Session{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"sessions\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, sessionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from sessions")
	}

	return sessionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Session) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no sessions provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(sessionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	sessionInsertCacheMut.RLock()
	cache, cached := sessionInsertCache[key]
	sessionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			sessionAllColumns,
			sessionColumnsWithDefault,
			sessionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(sessionType, sessionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(sessionType, sessionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"sessions\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"sessions\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into sessions")
	}

	if !cached {
		sessionInsertCacheMut.Lock()
		sessionInsertCache[key] = cache
		sessionInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Session.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Session) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	sessionUpdateCacheMut.RLock()
	cache, cached := sessionUpdateCache[key]
	sessionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			sessionAllColumns,
			sessionPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update sessions, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"sessions\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, sessionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(sessionType, sessionMapping, append(wl, sessionPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update sessions row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for sessions")
	}

	if !cached {
		sessionUpdateCacheMut.Lock()
		sessionUpdateCache[key] = cache
		sessionUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q sessionQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for sessions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for sessions")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SessionSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sessionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"sessions\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, sessionPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in session slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all session")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Session) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no sessions provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(sessionColumnsWithDefault, o)

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

	sessionUpsertCacheMut.RLock()
	cache, cached := sessionUpsertCache[key]
	sessionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			sessionAllColumns,
			sessionColumnsWithDefault,
			sessionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			sessionAllColumns,
			sessionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert sessions, could not build update column list")
		}

		ret := strmangle.SetComplement(sessionAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(sessionPrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert sessions, could not build conflict column list")
			}

			conflict = make([]string, len(sessionPrimaryKeyColumns))
			copy(conflict, sessionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"sessions\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(sessionType, sessionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(sessionType, sessionMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert sessions")
	}

	if !cached {
		sessionUpsertCacheMut.Lock()
		sessionUpsertCache[key] = cache
		sessionUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Session record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Session) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no Session provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sessionPrimaryKeyMapping)
	sql := "DELETE FROM \"sessions\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from sessions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for sessions")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q sessionQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no sessionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from sessions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for sessions")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SessionSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sessionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"sessions\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sessionPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from session slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for sessions")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Session) Reload(exec boil.Executor) error {
	ret, err := FindSession(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SessionSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SessionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sessionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"sessions\".* FROM \"sessions\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sessionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in SessionSlice")
	}

	*o = slice

	return nil
}

// SessionExists checks if the Session row exists.
func SessionExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"sessions\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if sessions exists")
	}

	return exists, nil
}

// Exists checks if the Session row exists.
func (o *Session) Exists(exec boil.Executor) (bool, error) {
	return SessionExists(exec, o.ID)
}
