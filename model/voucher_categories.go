// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package model

import (
	"context"
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

// VoucherCategory is an object representing the database table.
type VoucherCategory struct {
	ID         string `boil:"id" json:"id" toml:"id" yaml:"id"`
	VoucherID  string `boil:"voucher_id" json:"voucher_id" toml:"voucher_id" yaml:"voucher_id"`
	CategoryID string `boil:"category_id" json:"category_id" toml:"category_id" yaml:"category_id"`
	CreatedAt  int64  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *voucherCategoryR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L voucherCategoryL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var VoucherCategoryColumns = struct {
	ID         string
	VoucherID  string
	CategoryID string
	CreatedAt  string
}{
	ID:         "id",
	VoucherID:  "voucher_id",
	CategoryID: "category_id",
	CreatedAt:  "created_at",
}

var VoucherCategoryTableColumns = struct {
	ID         string
	VoucherID  string
	CategoryID string
	CreatedAt  string
}{
	ID:         "voucher_categories.id",
	VoucherID:  "voucher_categories.voucher_id",
	CategoryID: "voucher_categories.category_id",
	CreatedAt:  "voucher_categories.created_at",
}

// Generated where

var VoucherCategoryWhere = struct {
	ID         whereHelperstring
	VoucherID  whereHelperstring
	CategoryID whereHelperstring
	CreatedAt  whereHelperint64
}{
	ID:         whereHelperstring{field: "\"voucher_categories\".\"id\""},
	VoucherID:  whereHelperstring{field: "\"voucher_categories\".\"voucher_id\""},
	CategoryID: whereHelperstring{field: "\"voucher_categories\".\"category_id\""},
	CreatedAt:  whereHelperint64{field: "\"voucher_categories\".\"created_at\""},
}

// VoucherCategoryRels is where relationship names are stored.
var VoucherCategoryRels = struct {
	Category string
	Voucher  string
}{
	Category: "Category",
	Voucher:  "Voucher",
}

// voucherCategoryR is where relationships are stored.
type voucherCategoryR struct {
	Category *Category `boil:"Category" json:"Category" toml:"Category" yaml:"Category"`
	Voucher  *Voucher  `boil:"Voucher" json:"Voucher" toml:"Voucher" yaml:"Voucher"`
}

// NewStruct creates a new relationship struct
func (*voucherCategoryR) NewStruct() *voucherCategoryR {
	return &voucherCategoryR{}
}

func (r *voucherCategoryR) GetCategory() *Category {
	if r == nil {
		return nil
	}
	return r.Category
}

func (r *voucherCategoryR) GetVoucher() *Voucher {
	if r == nil {
		return nil
	}
	return r.Voucher
}

// voucherCategoryL is where Load methods for each relationship are stored.
type voucherCategoryL struct{}

var (
	voucherCategoryAllColumns            = []string{"id", "voucher_id", "category_id", "created_at"}
	voucherCategoryColumnsWithoutDefault = []string{"voucher_id", "category_id", "created_at"}
	voucherCategoryColumnsWithDefault    = []string{"id"}
	voucherCategoryPrimaryKeyColumns     = []string{"id"}
	voucherCategoryGeneratedColumns      = []string{}
)

type (
	// VoucherCategorySlice is an alias for a slice of pointers to VoucherCategory.
	// This should almost always be used instead of []VoucherCategory.
	VoucherCategorySlice []*VoucherCategory

	voucherCategoryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	voucherCategoryType                 = reflect.TypeOf(&VoucherCategory{})
	voucherCategoryMapping              = queries.MakeStructMapping(voucherCategoryType)
	voucherCategoryPrimaryKeyMapping, _ = queries.BindMapping(voucherCategoryType, voucherCategoryMapping, voucherCategoryPrimaryKeyColumns)
	voucherCategoryInsertCacheMut       sync.RWMutex
	voucherCategoryInsertCache          = make(map[string]insertCache)
	voucherCategoryUpdateCacheMut       sync.RWMutex
	voucherCategoryUpdateCache          = make(map[string]updateCache)
	voucherCategoryUpsertCacheMut       sync.RWMutex
	voucherCategoryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single voucherCategory record from the query.
func (q voucherCategoryQuery) One(ctx context.Context, exec boil.ContextExecutor) (*VoucherCategory, error) {
	o := &VoucherCategory{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for voucher_categories")
	}

	return o, nil
}

// All returns all VoucherCategory records from the query.
func (q voucherCategoryQuery) All(ctx context.Context, exec boil.ContextExecutor) (VoucherCategorySlice, error) {
	var o []*VoucherCategory

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to VoucherCategory slice")
	}

	return o, nil
}

// Count returns the count of all VoucherCategory records in the query.
func (q voucherCategoryQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count voucher_categories rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q voucherCategoryQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if voucher_categories exists")
	}

	return count > 0, nil
}

// Category pointed to by the foreign key.
func (o *VoucherCategory) Category(mods ...qm.QueryMod) categoryQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.CategoryID),
	}

	queryMods = append(queryMods, mods...)

	return Categories(queryMods...)
}

// Voucher pointed to by the foreign key.
func (o *VoucherCategory) Voucher(mods ...qm.QueryMod) voucherQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.VoucherID),
	}

	queryMods = append(queryMods, mods...)

	return Vouchers(queryMods...)
}

// LoadCategory allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (voucherCategoryL) LoadCategory(ctx context.Context, e boil.ContextExecutor, singular bool, maybeVoucherCategory interface{}, mods queries.Applicator) error {
	var slice []*VoucherCategory
	var object *VoucherCategory

	if singular {
		var ok bool
		object, ok = maybeVoucherCategory.(*VoucherCategory)
		if !ok {
			object = new(VoucherCategory)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeVoucherCategory)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeVoucherCategory))
			}
		}
	} else {
		s, ok := maybeVoucherCategory.(*[]*VoucherCategory)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeVoucherCategory)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeVoucherCategory))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &voucherCategoryR{}
		}
		args = append(args, object.CategoryID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &voucherCategoryR{}
			}

			for _, a := range args {
				if a == obj.CategoryID {
					continue Outer
				}
			}

			args = append(args, obj.CategoryID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`categories`),
		qm.WhereIn(`categories.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Category")
	}

	var resultSlice []*Category
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Category")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for categories")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for categories")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Category = foreign
		if foreign.R == nil {
			foreign.R = &categoryR{}
		}
		foreign.R.VoucherCategories = append(foreign.R.VoucherCategories, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CategoryID == foreign.ID {
				local.R.Category = foreign
				if foreign.R == nil {
					foreign.R = &categoryR{}
				}
				foreign.R.VoucherCategories = append(foreign.R.VoucherCategories, local)
				break
			}
		}
	}

	return nil
}

// LoadVoucher allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (voucherCategoryL) LoadVoucher(ctx context.Context, e boil.ContextExecutor, singular bool, maybeVoucherCategory interface{}, mods queries.Applicator) error {
	var slice []*VoucherCategory
	var object *VoucherCategory

	if singular {
		var ok bool
		object, ok = maybeVoucherCategory.(*VoucherCategory)
		if !ok {
			object = new(VoucherCategory)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeVoucherCategory)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeVoucherCategory))
			}
		}
	} else {
		s, ok := maybeVoucherCategory.(*[]*VoucherCategory)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeVoucherCategory)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeVoucherCategory))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &voucherCategoryR{}
		}
		args = append(args, object.VoucherID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &voucherCategoryR{}
			}

			for _, a := range args {
				if a == obj.VoucherID {
					continue Outer
				}
			}

			args = append(args, obj.VoucherID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`vouchers`),
		qm.WhereIn(`vouchers.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
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
		foreign.R.VoucherCategories = append(foreign.R.VoucherCategories, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.VoucherID == foreign.ID {
				local.R.Voucher = foreign
				if foreign.R == nil {
					foreign.R = &voucherR{}
				}
				foreign.R.VoucherCategories = append(foreign.R.VoucherCategories, local)
				break
			}
		}
	}

	return nil
}

// SetCategory of the voucherCategory to the related item.
// Sets o.R.Category to related.
// Adds o to related.R.VoucherCategories.
func (o *VoucherCategory) SetCategory(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Category) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"voucher_categories\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"category_id"}),
		strmangle.WhereClause("\"", "\"", 2, voucherCategoryPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CategoryID = related.ID
	if o.R == nil {
		o.R = &voucherCategoryR{
			Category: related,
		}
	} else {
		o.R.Category = related
	}

	if related.R == nil {
		related.R = &categoryR{
			VoucherCategories: VoucherCategorySlice{o},
		}
	} else {
		related.R.VoucherCategories = append(related.R.VoucherCategories, o)
	}

	return nil
}

// SetVoucher of the voucherCategory to the related item.
// Sets o.R.Voucher to related.
// Adds o to related.R.VoucherCategories.
func (o *VoucherCategory) SetVoucher(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Voucher) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"voucher_categories\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"voucher_id"}),
		strmangle.WhereClause("\"", "\"", 2, voucherCategoryPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.VoucherID = related.ID
	if o.R == nil {
		o.R = &voucherCategoryR{
			Voucher: related,
		}
	} else {
		o.R.Voucher = related
	}

	if related.R == nil {
		related.R = &voucherR{
			VoucherCategories: VoucherCategorySlice{o},
		}
	} else {
		related.R.VoucherCategories = append(related.R.VoucherCategories, o)
	}

	return nil
}

// VoucherCategories retrieves all the records using an executor.
func VoucherCategories(mods ...qm.QueryMod) voucherCategoryQuery {
	mods = append(mods, qm.From("\"voucher_categories\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"voucher_categories\".*"})
	}

	return voucherCategoryQuery{q}
}

// FindVoucherCategory retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindVoucherCategory(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*VoucherCategory, error) {
	voucherCategoryObj := &VoucherCategory{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"voucher_categories\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, voucherCategoryObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from voucher_categories")
	}

	return voucherCategoryObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *VoucherCategory) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no voucher_categories provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(voucherCategoryColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	voucherCategoryInsertCacheMut.RLock()
	cache, cached := voucherCategoryInsertCache[key]
	voucherCategoryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			voucherCategoryAllColumns,
			voucherCategoryColumnsWithDefault,
			voucherCategoryColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(voucherCategoryType, voucherCategoryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(voucherCategoryType, voucherCategoryMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"voucher_categories\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"voucher_categories\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "model: unable to insert into voucher_categories")
	}

	if !cached {
		voucherCategoryInsertCacheMut.Lock()
		voucherCategoryInsertCache[key] = cache
		voucherCategoryInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the VoucherCategory.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *VoucherCategory) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	voucherCategoryUpdateCacheMut.RLock()
	cache, cached := voucherCategoryUpdateCache[key]
	voucherCategoryUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			voucherCategoryAllColumns,
			voucherCategoryPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update voucher_categories, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"voucher_categories\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, voucherCategoryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(voucherCategoryType, voucherCategoryMapping, append(wl, voucherCategoryPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update voucher_categories row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for voucher_categories")
	}

	if !cached {
		voucherCategoryUpdateCacheMut.Lock()
		voucherCategoryUpdateCache[key] = cache
		voucherCategoryUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q voucherCategoryQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for voucher_categories")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for voucher_categories")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o VoucherCategorySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), voucherCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"voucher_categories\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, voucherCategoryPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in voucherCategory slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all voucherCategory")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *VoucherCategory) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no voucher_categories provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(voucherCategoryColumnsWithDefault, o)

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

	voucherCategoryUpsertCacheMut.RLock()
	cache, cached := voucherCategoryUpsertCache[key]
	voucherCategoryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			voucherCategoryAllColumns,
			voucherCategoryColumnsWithDefault,
			voucherCategoryColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			voucherCategoryAllColumns,
			voucherCategoryPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert voucher_categories, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(voucherCategoryPrimaryKeyColumns))
			copy(conflict, voucherCategoryPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"voucher_categories\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(voucherCategoryType, voucherCategoryMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(voucherCategoryType, voucherCategoryMapping, ret)
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

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "model: unable to upsert voucher_categories")
	}

	if !cached {
		voucherCategoryUpsertCacheMut.Lock()
		voucherCategoryUpsertCache[key] = cache
		voucherCategoryUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single VoucherCategory record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *VoucherCategory) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no VoucherCategory provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), voucherCategoryPrimaryKeyMapping)
	sql := "DELETE FROM \"voucher_categories\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from voucher_categories")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for voucher_categories")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q voucherCategoryQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no voucherCategoryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from voucher_categories")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for voucher_categories")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o VoucherCategorySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), voucherCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"voucher_categories\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, voucherCategoryPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from voucherCategory slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for voucher_categories")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *VoucherCategory) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindVoucherCategory(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *VoucherCategorySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := VoucherCategorySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), voucherCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"voucher_categories\".* FROM \"voucher_categories\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, voucherCategoryPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in VoucherCategorySlice")
	}

	*o = slice

	return nil
}

// VoucherCategoryExists checks if the VoucherCategory row exists.
func VoucherCategoryExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"voucher_categories\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if voucher_categories exists")
	}

	return exists, nil
}

// Exists checks if the VoucherCategory row exists.
func (o *VoucherCategory) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return VoucherCategoryExists(ctx, exec, o.ID)
}