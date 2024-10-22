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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// VoucherTranslation is an object representing the database table.
type VoucherTranslation struct {
	ID           string       `boil:"id" json:"id" toml:"id" yaml:"id"`
	LanguageCode LanguageCode `boil:"language_code" json:"language_code" toml:"language_code" yaml:"language_code"`
	Name         string       `boil:"name" json:"name" toml:"name" yaml:"name"`
	VoucherID    string       `boil:"voucher_id" json:"voucher_id" toml:"voucher_id" yaml:"voucher_id"`
	CreatedAt    int64        `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *voucherTranslationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L voucherTranslationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var VoucherTranslationColumns = struct {
	ID           string
	LanguageCode string
	Name         string
	VoucherID    string
	CreatedAt    string
}{
	ID:           "id",
	LanguageCode: "language_code",
	Name:         "name",
	VoucherID:    "voucher_id",
	CreatedAt:    "created_at",
}

var VoucherTranslationTableColumns = struct {
	ID           string
	LanguageCode string
	Name         string
	VoucherID    string
	CreatedAt    string
}{
	ID:           "voucher_translations.id",
	LanguageCode: "voucher_translations.language_code",
	Name:         "voucher_translations.name",
	VoucherID:    "voucher_translations.voucher_id",
	CreatedAt:    "voucher_translations.created_at",
}

// Generated where

var VoucherTranslationWhere = struct {
	ID           whereHelperstring
	LanguageCode whereHelperLanguageCode
	Name         whereHelperstring
	VoucherID    whereHelperstring
	CreatedAt    whereHelperint64
}{
	ID:           whereHelperstring{field: "\"voucher_translations\".\"id\""},
	LanguageCode: whereHelperLanguageCode{field: "\"voucher_translations\".\"language_code\""},
	Name:         whereHelperstring{field: "\"voucher_translations\".\"name\""},
	VoucherID:    whereHelperstring{field: "\"voucher_translations\".\"voucher_id\""},
	CreatedAt:    whereHelperint64{field: "\"voucher_translations\".\"created_at\""},
}

// VoucherTranslationRels is where relationship names are stored.
var VoucherTranslationRels = struct {
	Voucher string
}{
	Voucher: "Voucher",
}

// voucherTranslationR is where relationships are stored.
type voucherTranslationR struct {
	Voucher *Voucher `boil:"Voucher" json:"Voucher" toml:"Voucher" yaml:"Voucher"`
}

// NewStruct creates a new relationship struct
func (*voucherTranslationR) NewStruct() *voucherTranslationR {
	return &voucherTranslationR{}
}

func (r *voucherTranslationR) GetVoucher() *Voucher {
	if r == nil {
		return nil
	}
	return r.Voucher
}

// voucherTranslationL is where Load methods for each relationship are stored.
type voucherTranslationL struct{}

var (
	voucherTranslationAllColumns            = []string{"id", "language_code", "name", "voucher_id", "created_at"}
	voucherTranslationColumnsWithoutDefault = []string{"id", "language_code", "name", "voucher_id", "created_at"}
	voucherTranslationColumnsWithDefault    = []string{}
	voucherTranslationPrimaryKeyColumns     = []string{"id"}
	voucherTranslationGeneratedColumns      = []string{}
)

type (
	// VoucherTranslationSlice is an alias for a slice of pointers to VoucherTranslation.
	// This should almost always be used instead of []VoucherTranslation.
	VoucherTranslationSlice []*VoucherTranslation

	voucherTranslationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	voucherTranslationType                 = reflect.TypeOf(&VoucherTranslation{})
	voucherTranslationMapping              = queries.MakeStructMapping(voucherTranslationType)
	voucherTranslationPrimaryKeyMapping, _ = queries.BindMapping(voucherTranslationType, voucherTranslationMapping, voucherTranslationPrimaryKeyColumns)
	voucherTranslationInsertCacheMut       sync.RWMutex
	voucherTranslationInsertCache          = make(map[string]insertCache)
	voucherTranslationUpdateCacheMut       sync.RWMutex
	voucherTranslationUpdateCache          = make(map[string]updateCache)
	voucherTranslationUpsertCacheMut       sync.RWMutex
	voucherTranslationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single voucherTranslation record from the query.
func (q voucherTranslationQuery) One(exec boil.Executor) (*VoucherTranslation, error) {
	o := &VoucherTranslation{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for voucher_translations")
	}

	return o, nil
}

// All returns all VoucherTranslation records from the query.
func (q voucherTranslationQuery) All(exec boil.Executor) (VoucherTranslationSlice, error) {
	var o []*VoucherTranslation

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to VoucherTranslation slice")
	}

	return o, nil
}

// Count returns the count of all VoucherTranslation records in the query.
func (q voucherTranslationQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count voucher_translations rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q voucherTranslationQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if voucher_translations exists")
	}

	return count > 0, nil
}

// Voucher pointed to by the foreign key.
func (o *VoucherTranslation) Voucher(mods ...qm.QueryMod) voucherQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.VoucherID),
	}

	queryMods = append(queryMods, mods...)

	return Vouchers(queryMods...)
}

// LoadVoucher allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (voucherTranslationL) LoadVoucher(e boil.Executor, singular bool, maybeVoucherTranslation interface{}, mods queries.Applicator) error {
	var slice []*VoucherTranslation
	var object *VoucherTranslation

	if singular {
		var ok bool
		object, ok = maybeVoucherTranslation.(*VoucherTranslation)
		if !ok {
			object = new(VoucherTranslation)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeVoucherTranslation)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeVoucherTranslation))
			}
		}
	} else {
		s, ok := maybeVoucherTranslation.(*[]*VoucherTranslation)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeVoucherTranslation)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeVoucherTranslation))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &voucherTranslationR{}
		}
		args[object.VoucherID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &voucherTranslationR{}
			}

			args[obj.VoucherID] = struct{}{}

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
		qm.From(`vouchers`),
		qm.WhereIn(`vouchers.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Voucher")
	}

	var resultSlice []*Voucher
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Voucher")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for vouchers")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for vouchers")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Voucher = foreign
		if foreign.R == nil {
			foreign.R = &voucherR{}
		}
		foreign.R.VoucherTranslations = append(foreign.R.VoucherTranslations, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.VoucherID == foreign.ID {
				local.R.Voucher = foreign
				if foreign.R == nil {
					foreign.R = &voucherR{}
				}
				foreign.R.VoucherTranslations = append(foreign.R.VoucherTranslations, local)
				break
			}
		}
	}

	return nil
}

// SetVoucher of the voucherTranslation to the related item.
// Sets o.R.Voucher to related.
// Adds o to related.R.VoucherTranslations.
func (o *VoucherTranslation) SetVoucher(exec boil.Executor, insert bool, related *Voucher) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"voucher_translations\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"voucher_id"}),
		strmangle.WhereClause("\"", "\"", 2, voucherTranslationPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.VoucherID = related.ID
	if o.R == nil {
		o.R = &voucherTranslationR{
			Voucher: related,
		}
	} else {
		o.R.Voucher = related
	}

	if related.R == nil {
		related.R = &voucherR{
			VoucherTranslations: VoucherTranslationSlice{o},
		}
	} else {
		related.R.VoucherTranslations = append(related.R.VoucherTranslations, o)
	}

	return nil
}

// VoucherTranslations retrieves all the records using an executor.
func VoucherTranslations(mods ...qm.QueryMod) voucherTranslationQuery {
	mods = append(mods, qm.From("\"voucher_translations\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"voucher_translations\".*"})
	}

	return voucherTranslationQuery{q}
}

// FindVoucherTranslation retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindVoucherTranslation(exec boil.Executor, iD string, selectCols ...string) (*VoucherTranslation, error) {
	voucherTranslationObj := &VoucherTranslation{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"voucher_translations\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, voucherTranslationObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from voucher_translations")
	}

	return voucherTranslationObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *VoucherTranslation) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no voucher_translations provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(voucherTranslationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	voucherTranslationInsertCacheMut.RLock()
	cache, cached := voucherTranslationInsertCache[key]
	voucherTranslationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			voucherTranslationAllColumns,
			voucherTranslationColumnsWithDefault,
			voucherTranslationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(voucherTranslationType, voucherTranslationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(voucherTranslationType, voucherTranslationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"voucher_translations\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"voucher_translations\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into voucher_translations")
	}

	if !cached {
		voucherTranslationInsertCacheMut.Lock()
		voucherTranslationInsertCache[key] = cache
		voucherTranslationInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the VoucherTranslation.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *VoucherTranslation) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	voucherTranslationUpdateCacheMut.RLock()
	cache, cached := voucherTranslationUpdateCache[key]
	voucherTranslationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			voucherTranslationAllColumns,
			voucherTranslationPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update voucher_translations, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"voucher_translations\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, voucherTranslationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(voucherTranslationType, voucherTranslationMapping, append(wl, voucherTranslationPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update voucher_translations row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for voucher_translations")
	}

	if !cached {
		voucherTranslationUpdateCacheMut.Lock()
		voucherTranslationUpdateCache[key] = cache
		voucherTranslationUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q voucherTranslationQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for voucher_translations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for voucher_translations")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o VoucherTranslationSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), voucherTranslationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"voucher_translations\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, voucherTranslationPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in voucherTranslation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all voucherTranslation")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *VoucherTranslation) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no voucher_translations provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(voucherTranslationColumnsWithDefault, o)

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

	voucherTranslationUpsertCacheMut.RLock()
	cache, cached := voucherTranslationUpsertCache[key]
	voucherTranslationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			voucherTranslationAllColumns,
			voucherTranslationColumnsWithDefault,
			voucherTranslationColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			voucherTranslationAllColumns,
			voucherTranslationPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert voucher_translations, could not build update column list")
		}

		ret := strmangle.SetComplement(voucherTranslationAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(voucherTranslationPrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert voucher_translations, could not build conflict column list")
			}

			conflict = make([]string, len(voucherTranslationPrimaryKeyColumns))
			copy(conflict, voucherTranslationPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"voucher_translations\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(voucherTranslationType, voucherTranslationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(voucherTranslationType, voucherTranslationMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert voucher_translations")
	}

	if !cached {
		voucherTranslationUpsertCacheMut.Lock()
		voucherTranslationUpsertCache[key] = cache
		voucherTranslationUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single VoucherTranslation record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *VoucherTranslation) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no VoucherTranslation provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), voucherTranslationPrimaryKeyMapping)
	sql := "DELETE FROM \"voucher_translations\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from voucher_translations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for voucher_translations")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q voucherTranslationQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no voucherTranslationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from voucher_translations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for voucher_translations")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o VoucherTranslationSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), voucherTranslationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"voucher_translations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, voucherTranslationPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from voucherTranslation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for voucher_translations")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *VoucherTranslation) Reload(exec boil.Executor) error {
	ret, err := FindVoucherTranslation(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *VoucherTranslationSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := VoucherTranslationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), voucherTranslationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"voucher_translations\".* FROM \"voucher_translations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, voucherTranslationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in VoucherTranslationSlice")
	}

	*o = slice

	return nil
}

// VoucherTranslationExists checks if the VoucherTranslation row exists.
func VoucherTranslationExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"voucher_translations\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if voucher_translations exists")
	}

	return exists, nil
}

// Exists checks if the VoucherTranslation row exists.
func (o *VoucherTranslation) Exists(exec boil.Executor) (bool, error) {
	return VoucherTranslationExists(exec, o.ID)
}
