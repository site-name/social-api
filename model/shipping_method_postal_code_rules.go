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

// ShippingMethodPostalCodeRule is an object representing the database table.
type ShippingMethodPostalCodeRule struct {
	ID               string `boil:"id" json:"id" toml:"id" yaml:"id"`
	ShippingMethodID string `boil:"shipping_method_id" json:"shipping_method_id" toml:"shipping_method_id" yaml:"shipping_method_id"`
	Start            string `boil:"start" json:"start" toml:"start" yaml:"start"`
	End              string `boil:"end" json:"end" toml:"end" yaml:"end"`
	InclusionType    string `boil:"inclusion_type" json:"inclusion_type" toml:"inclusion_type" yaml:"inclusion_type"`

	R *shippingMethodPostalCodeRuleR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L shippingMethodPostalCodeRuleL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ShippingMethodPostalCodeRuleColumns = struct {
	ID               string
	ShippingMethodID string
	Start            string
	End              string
	InclusionType    string
}{
	ID:               "id",
	ShippingMethodID: "shipping_method_id",
	Start:            "start",
	End:              "end",
	InclusionType:    "inclusion_type",
}

var ShippingMethodPostalCodeRuleTableColumns = struct {
	ID               string
	ShippingMethodID string
	Start            string
	End              string
	InclusionType    string
}{
	ID:               "shipping_method_postal_code_rules.id",
	ShippingMethodID: "shipping_method_postal_code_rules.shipping_method_id",
	Start:            "shipping_method_postal_code_rules.start",
	End:              "shipping_method_postal_code_rules.end",
	InclusionType:    "shipping_method_postal_code_rules.inclusion_type",
}

// Generated where

var ShippingMethodPostalCodeRuleWhere = struct {
	ID               whereHelperstring
	ShippingMethodID whereHelperstring
	Start            whereHelperstring
	End              whereHelperstring
	InclusionType    whereHelperstring
}{
	ID:               whereHelperstring{field: "\"shipping_method_postal_code_rules\".\"id\""},
	ShippingMethodID: whereHelperstring{field: "\"shipping_method_postal_code_rules\".\"shipping_method_id\""},
	Start:            whereHelperstring{field: "\"shipping_method_postal_code_rules\".\"start\""},
	End:              whereHelperstring{field: "\"shipping_method_postal_code_rules\".\"end\""},
	InclusionType:    whereHelperstring{field: "\"shipping_method_postal_code_rules\".\"inclusion_type\""},
}

// ShippingMethodPostalCodeRuleRels is where relationship names are stored.
var ShippingMethodPostalCodeRuleRels = struct {
	ShippingMethod string
}{
	ShippingMethod: "ShippingMethod",
}

// shippingMethodPostalCodeRuleR is where relationships are stored.
type shippingMethodPostalCodeRuleR struct {
	ShippingMethod *ShippingMethod `boil:"ShippingMethod" json:"ShippingMethod" toml:"ShippingMethod" yaml:"ShippingMethod"`
}

// NewStruct creates a new relationship struct
func (*shippingMethodPostalCodeRuleR) NewStruct() *shippingMethodPostalCodeRuleR {
	return &shippingMethodPostalCodeRuleR{}
}

func (r *shippingMethodPostalCodeRuleR) GetShippingMethod() *ShippingMethod {
	if r == nil {
		return nil
	}
	return r.ShippingMethod
}

// shippingMethodPostalCodeRuleL is where Load methods for each relationship are stored.
type shippingMethodPostalCodeRuleL struct{}

var (
	shippingMethodPostalCodeRuleAllColumns            = []string{"id", "shipping_method_id", "start", "end", "inclusion_type"}
	shippingMethodPostalCodeRuleColumnsWithoutDefault = []string{"shipping_method_id", "start", "end", "inclusion_type"}
	shippingMethodPostalCodeRuleColumnsWithDefault    = []string{"id"}
	shippingMethodPostalCodeRulePrimaryKeyColumns     = []string{"id"}
	shippingMethodPostalCodeRuleGeneratedColumns      = []string{}
)

type (
	// ShippingMethodPostalCodeRuleSlice is an alias for a slice of pointers to ShippingMethodPostalCodeRule.
	// This should almost always be used instead of []ShippingMethodPostalCodeRule.
	ShippingMethodPostalCodeRuleSlice []*ShippingMethodPostalCodeRule

	shippingMethodPostalCodeRuleQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	shippingMethodPostalCodeRuleType                 = reflect.TypeOf(&ShippingMethodPostalCodeRule{})
	shippingMethodPostalCodeRuleMapping              = queries.MakeStructMapping(shippingMethodPostalCodeRuleType)
	shippingMethodPostalCodeRulePrimaryKeyMapping, _ = queries.BindMapping(shippingMethodPostalCodeRuleType, shippingMethodPostalCodeRuleMapping, shippingMethodPostalCodeRulePrimaryKeyColumns)
	shippingMethodPostalCodeRuleInsertCacheMut       sync.RWMutex
	shippingMethodPostalCodeRuleInsertCache          = make(map[string]insertCache)
	shippingMethodPostalCodeRuleUpdateCacheMut       sync.RWMutex
	shippingMethodPostalCodeRuleUpdateCache          = make(map[string]updateCache)
	shippingMethodPostalCodeRuleUpsertCacheMut       sync.RWMutex
	shippingMethodPostalCodeRuleUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single shippingMethodPostalCodeRule record from the query.
func (q shippingMethodPostalCodeRuleQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ShippingMethodPostalCodeRule, error) {
	o := &ShippingMethodPostalCodeRule{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for shipping_method_postal_code_rules")
	}

	return o, nil
}

// All returns all ShippingMethodPostalCodeRule records from the query.
func (q shippingMethodPostalCodeRuleQuery) All(ctx context.Context, exec boil.ContextExecutor) (ShippingMethodPostalCodeRuleSlice, error) {
	var o []*ShippingMethodPostalCodeRule

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to ShippingMethodPostalCodeRule slice")
	}

	return o, nil
}

// Count returns the count of all ShippingMethodPostalCodeRule records in the query.
func (q shippingMethodPostalCodeRuleQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count shipping_method_postal_code_rules rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q shippingMethodPostalCodeRuleQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if shipping_method_postal_code_rules exists")
	}

	return count > 0, nil
}

// ShippingMethod pointed to by the foreign key.
func (o *ShippingMethodPostalCodeRule) ShippingMethod(mods ...qm.QueryMod) shippingMethodQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ShippingMethodID),
	}

	queryMods = append(queryMods, mods...)

	return ShippingMethods(queryMods...)
}

// LoadShippingMethod allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (shippingMethodPostalCodeRuleL) LoadShippingMethod(ctx context.Context, e boil.ContextExecutor, singular bool, maybeShippingMethodPostalCodeRule interface{}, mods queries.Applicator) error {
	var slice []*ShippingMethodPostalCodeRule
	var object *ShippingMethodPostalCodeRule

	if singular {
		var ok bool
		object, ok = maybeShippingMethodPostalCodeRule.(*ShippingMethodPostalCodeRule)
		if !ok {
			object = new(ShippingMethodPostalCodeRule)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeShippingMethodPostalCodeRule)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeShippingMethodPostalCodeRule))
			}
		}
	} else {
		s, ok := maybeShippingMethodPostalCodeRule.(*[]*ShippingMethodPostalCodeRule)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeShippingMethodPostalCodeRule)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeShippingMethodPostalCodeRule))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &shippingMethodPostalCodeRuleR{}
		}
		args = append(args, object.ShippingMethodID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &shippingMethodPostalCodeRuleR{}
			}

			for _, a := range args {
				if a == obj.ShippingMethodID {
					continue Outer
				}
			}

			args = append(args, obj.ShippingMethodID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`shipping_methods`),
		qm.WhereIn(`shipping_methods.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ShippingMethod")
	}

	var resultSlice []*ShippingMethod
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ShippingMethod")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for shipping_methods")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for shipping_methods")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.ShippingMethod = foreign
		if foreign.R == nil {
			foreign.R = &shippingMethodR{}
		}
		foreign.R.ShippingMethodPostalCodeRules = append(foreign.R.ShippingMethodPostalCodeRules, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ShippingMethodID == foreign.ID {
				local.R.ShippingMethod = foreign
				if foreign.R == nil {
					foreign.R = &shippingMethodR{}
				}
				foreign.R.ShippingMethodPostalCodeRules = append(foreign.R.ShippingMethodPostalCodeRules, local)
				break
			}
		}
	}

	return nil
}

// SetShippingMethod of the shippingMethodPostalCodeRule to the related item.
// Sets o.R.ShippingMethod to related.
// Adds o to related.R.ShippingMethodPostalCodeRules.
func (o *ShippingMethodPostalCodeRule) SetShippingMethod(ctx context.Context, exec boil.ContextExecutor, insert bool, related *ShippingMethod) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"shipping_method_postal_code_rules\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"shipping_method_id"}),
		strmangle.WhereClause("\"", "\"", 2, shippingMethodPostalCodeRulePrimaryKeyColumns),
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

	o.ShippingMethodID = related.ID
	if o.R == nil {
		o.R = &shippingMethodPostalCodeRuleR{
			ShippingMethod: related,
		}
	} else {
		o.R.ShippingMethod = related
	}

	if related.R == nil {
		related.R = &shippingMethodR{
			ShippingMethodPostalCodeRules: ShippingMethodPostalCodeRuleSlice{o},
		}
	} else {
		related.R.ShippingMethodPostalCodeRules = append(related.R.ShippingMethodPostalCodeRules, o)
	}

	return nil
}

// ShippingMethodPostalCodeRules retrieves all the records using an executor.
func ShippingMethodPostalCodeRules(mods ...qm.QueryMod) shippingMethodPostalCodeRuleQuery {
	mods = append(mods, qm.From("\"shipping_method_postal_code_rules\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"shipping_method_postal_code_rules\".*"})
	}

	return shippingMethodPostalCodeRuleQuery{q}
}

// FindShippingMethodPostalCodeRule retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindShippingMethodPostalCodeRule(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*ShippingMethodPostalCodeRule, error) {
	shippingMethodPostalCodeRuleObj := &ShippingMethodPostalCodeRule{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"shipping_method_postal_code_rules\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, shippingMethodPostalCodeRuleObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from shipping_method_postal_code_rules")
	}

	return shippingMethodPostalCodeRuleObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ShippingMethodPostalCodeRule) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no shipping_method_postal_code_rules provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(shippingMethodPostalCodeRuleColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	shippingMethodPostalCodeRuleInsertCacheMut.RLock()
	cache, cached := shippingMethodPostalCodeRuleInsertCache[key]
	shippingMethodPostalCodeRuleInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			shippingMethodPostalCodeRuleAllColumns,
			shippingMethodPostalCodeRuleColumnsWithDefault,
			shippingMethodPostalCodeRuleColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(shippingMethodPostalCodeRuleType, shippingMethodPostalCodeRuleMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(shippingMethodPostalCodeRuleType, shippingMethodPostalCodeRuleMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"shipping_method_postal_code_rules\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"shipping_method_postal_code_rules\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into shipping_method_postal_code_rules")
	}

	if !cached {
		shippingMethodPostalCodeRuleInsertCacheMut.Lock()
		shippingMethodPostalCodeRuleInsertCache[key] = cache
		shippingMethodPostalCodeRuleInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the ShippingMethodPostalCodeRule.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ShippingMethodPostalCodeRule) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	shippingMethodPostalCodeRuleUpdateCacheMut.RLock()
	cache, cached := shippingMethodPostalCodeRuleUpdateCache[key]
	shippingMethodPostalCodeRuleUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			shippingMethodPostalCodeRuleAllColumns,
			shippingMethodPostalCodeRulePrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update shipping_method_postal_code_rules, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"shipping_method_postal_code_rules\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, shippingMethodPostalCodeRulePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(shippingMethodPostalCodeRuleType, shippingMethodPostalCodeRuleMapping, append(wl, shippingMethodPostalCodeRulePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update shipping_method_postal_code_rules row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for shipping_method_postal_code_rules")
	}

	if !cached {
		shippingMethodPostalCodeRuleUpdateCacheMut.Lock()
		shippingMethodPostalCodeRuleUpdateCache[key] = cache
		shippingMethodPostalCodeRuleUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q shippingMethodPostalCodeRuleQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for shipping_method_postal_code_rules")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for shipping_method_postal_code_rules")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ShippingMethodPostalCodeRuleSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), shippingMethodPostalCodeRulePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"shipping_method_postal_code_rules\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, shippingMethodPostalCodeRulePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in shippingMethodPostalCodeRule slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all shippingMethodPostalCodeRule")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ShippingMethodPostalCodeRule) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no shipping_method_postal_code_rules provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(shippingMethodPostalCodeRuleColumnsWithDefault, o)

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

	shippingMethodPostalCodeRuleUpsertCacheMut.RLock()
	cache, cached := shippingMethodPostalCodeRuleUpsertCache[key]
	shippingMethodPostalCodeRuleUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			shippingMethodPostalCodeRuleAllColumns,
			shippingMethodPostalCodeRuleColumnsWithDefault,
			shippingMethodPostalCodeRuleColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			shippingMethodPostalCodeRuleAllColumns,
			shippingMethodPostalCodeRulePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert shipping_method_postal_code_rules, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(shippingMethodPostalCodeRulePrimaryKeyColumns))
			copy(conflict, shippingMethodPostalCodeRulePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"shipping_method_postal_code_rules\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(shippingMethodPostalCodeRuleType, shippingMethodPostalCodeRuleMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(shippingMethodPostalCodeRuleType, shippingMethodPostalCodeRuleMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert shipping_method_postal_code_rules")
	}

	if !cached {
		shippingMethodPostalCodeRuleUpsertCacheMut.Lock()
		shippingMethodPostalCodeRuleUpsertCache[key] = cache
		shippingMethodPostalCodeRuleUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single ShippingMethodPostalCodeRule record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ShippingMethodPostalCodeRule) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no ShippingMethodPostalCodeRule provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), shippingMethodPostalCodeRulePrimaryKeyMapping)
	sql := "DELETE FROM \"shipping_method_postal_code_rules\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from shipping_method_postal_code_rules")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for shipping_method_postal_code_rules")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q shippingMethodPostalCodeRuleQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no shippingMethodPostalCodeRuleQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from shipping_method_postal_code_rules")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for shipping_method_postal_code_rules")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ShippingMethodPostalCodeRuleSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), shippingMethodPostalCodeRulePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"shipping_method_postal_code_rules\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, shippingMethodPostalCodeRulePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from shippingMethodPostalCodeRule slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for shipping_method_postal_code_rules")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ShippingMethodPostalCodeRule) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindShippingMethodPostalCodeRule(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ShippingMethodPostalCodeRuleSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ShippingMethodPostalCodeRuleSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), shippingMethodPostalCodeRulePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"shipping_method_postal_code_rules\".* FROM \"shipping_method_postal_code_rules\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, shippingMethodPostalCodeRulePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in ShippingMethodPostalCodeRuleSlice")
	}

	*o = slice

	return nil
}

// ShippingMethodPostalCodeRuleExists checks if the ShippingMethodPostalCodeRule row exists.
func ShippingMethodPostalCodeRuleExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"shipping_method_postal_code_rules\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if shipping_method_postal_code_rules exists")
	}

	return exists, nil
}

// Exists checks if the ShippingMethodPostalCodeRule row exists.
func (o *ShippingMethodPostalCodeRule) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return ShippingMethodPostalCodeRuleExists(ctx, exec, o.ID)
}