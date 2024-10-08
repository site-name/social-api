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

// CustomProductAttributeValue is an object representing the database table.
type CustomProductAttributeValue struct {
	ID          string `boil:"id" json:"id" toml:"id" yaml:"id"`
	Value       string `boil:"value" json:"value" toml:"value" yaml:"value"`
	AttributeID string `boil:"attribute_id" json:"attribute_id" toml:"attribute_id" yaml:"attribute_id"`

	R *customProductAttributeValueR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L customProductAttributeValueL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CustomProductAttributeValueColumns = struct {
	ID          string
	Value       string
	AttributeID string
}{
	ID:          "id",
	Value:       "value",
	AttributeID: "attribute_id",
}

var CustomProductAttributeValueTableColumns = struct {
	ID          string
	Value       string
	AttributeID string
}{
	ID:          "custom_product_attribute_values.id",
	Value:       "custom_product_attribute_values.value",
	AttributeID: "custom_product_attribute_values.attribute_id",
}

// Generated where

var CustomProductAttributeValueWhere = struct {
	ID          whereHelperstring
	Value       whereHelperstring
	AttributeID whereHelperstring
}{
	ID:          whereHelperstring{field: "\"custom_product_attribute_values\".\"id\""},
	Value:       whereHelperstring{field: "\"custom_product_attribute_values\".\"value\""},
	AttributeID: whereHelperstring{field: "\"custom_product_attribute_values\".\"attribute_id\""},
}

// CustomProductAttributeValueRels is where relationship names are stored.
var CustomProductAttributeValueRels = struct {
	Attribute                                           string
	AttributeValueAssignedProductVariantAttributeValues string
}{
	Attribute: "Attribute",
	AttributeValueAssignedProductVariantAttributeValues: "AttributeValueAssignedProductVariantAttributeValues",
}

// customProductAttributeValueR is where relationships are stored.
type customProductAttributeValueR struct {
	Attribute                                           *CustomProductAttribute                   `boil:"Attribute" json:"Attribute" toml:"Attribute" yaml:"Attribute"`
	AttributeValueAssignedProductVariantAttributeValues AssignedProductVariantAttributeValueSlice `boil:"AttributeValueAssignedProductVariantAttributeValues" json:"AttributeValueAssignedProductVariantAttributeValues" toml:"AttributeValueAssignedProductVariantAttributeValues" yaml:"AttributeValueAssignedProductVariantAttributeValues"`
}

// NewStruct creates a new relationship struct
func (*customProductAttributeValueR) NewStruct() *customProductAttributeValueR {
	return &customProductAttributeValueR{}
}

func (r *customProductAttributeValueR) GetAttribute() *CustomProductAttribute {
	if r == nil {
		return nil
	}
	return r.Attribute
}

func (r *customProductAttributeValueR) GetAttributeValueAssignedProductVariantAttributeValues() AssignedProductVariantAttributeValueSlice {
	if r == nil {
		return nil
	}
	return r.AttributeValueAssignedProductVariantAttributeValues
}

// customProductAttributeValueL is where Load methods for each relationship are stored.
type customProductAttributeValueL struct{}

var (
	customProductAttributeValueAllColumns            = []string{"id", "value", "attribute_id"}
	customProductAttributeValueColumnsWithoutDefault = []string{"id", "value", "attribute_id"}
	customProductAttributeValueColumnsWithDefault    = []string{}
	customProductAttributeValuePrimaryKeyColumns     = []string{"id"}
	customProductAttributeValueGeneratedColumns      = []string{}
)

type (
	// CustomProductAttributeValueSlice is an alias for a slice of pointers to CustomProductAttributeValue.
	// This should almost always be used instead of []CustomProductAttributeValue.
	CustomProductAttributeValueSlice []*CustomProductAttributeValue

	customProductAttributeValueQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	customProductAttributeValueType                 = reflect.TypeOf(&CustomProductAttributeValue{})
	customProductAttributeValueMapping              = queries.MakeStructMapping(customProductAttributeValueType)
	customProductAttributeValuePrimaryKeyMapping, _ = queries.BindMapping(customProductAttributeValueType, customProductAttributeValueMapping, customProductAttributeValuePrimaryKeyColumns)
	customProductAttributeValueInsertCacheMut       sync.RWMutex
	customProductAttributeValueInsertCache          = make(map[string]insertCache)
	customProductAttributeValueUpdateCacheMut       sync.RWMutex
	customProductAttributeValueUpdateCache          = make(map[string]updateCache)
	customProductAttributeValueUpsertCacheMut       sync.RWMutex
	customProductAttributeValueUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single customProductAttributeValue record from the query.
func (q customProductAttributeValueQuery) One(exec boil.Executor) (*CustomProductAttributeValue, error) {
	o := &CustomProductAttributeValue{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for custom_product_attribute_values")
	}

	return o, nil
}

// All returns all CustomProductAttributeValue records from the query.
func (q customProductAttributeValueQuery) All(exec boil.Executor) (CustomProductAttributeValueSlice, error) {
	var o []*CustomProductAttributeValue

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to CustomProductAttributeValue slice")
	}

	return o, nil
}

// Count returns the count of all CustomProductAttributeValue records in the query.
func (q customProductAttributeValueQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count custom_product_attribute_values rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q customProductAttributeValueQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if custom_product_attribute_values exists")
	}

	return count > 0, nil
}

// Attribute pointed to by the foreign key.
func (o *CustomProductAttributeValue) Attribute(mods ...qm.QueryMod) customProductAttributeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.AttributeID),
	}

	queryMods = append(queryMods, mods...)

	return CustomProductAttributes(queryMods...)
}

// AttributeValueAssignedProductVariantAttributeValues retrieves all the assigned_product_variant_attribute_value's AssignedProductVariantAttributeValues with an executor via attribute_value_id column.
func (o *CustomProductAttributeValue) AttributeValueAssignedProductVariantAttributeValues(mods ...qm.QueryMod) assignedProductVariantAttributeValueQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"assigned_product_variant_attribute_values\".\"attribute_value_id\"=?", o.ID),
	)

	return AssignedProductVariantAttributeValues(queryMods...)
}

// LoadAttribute allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (customProductAttributeValueL) LoadAttribute(e boil.Executor, singular bool, maybeCustomProductAttributeValue interface{}, mods queries.Applicator) error {
	var slice []*CustomProductAttributeValue
	var object *CustomProductAttributeValue

	if singular {
		var ok bool
		object, ok = maybeCustomProductAttributeValue.(*CustomProductAttributeValue)
		if !ok {
			object = new(CustomProductAttributeValue)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCustomProductAttributeValue)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCustomProductAttributeValue))
			}
		}
	} else {
		s, ok := maybeCustomProductAttributeValue.(*[]*CustomProductAttributeValue)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCustomProductAttributeValue)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCustomProductAttributeValue))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &customProductAttributeValueR{}
		}
		args[object.AttributeID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &customProductAttributeValueR{}
			}

			args[obj.AttributeID] = struct{}{}

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
		qm.From(`custom_product_attributes`),
		qm.WhereIn(`custom_product_attributes.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load CustomProductAttribute")
	}

	var resultSlice []*CustomProductAttribute
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice CustomProductAttribute")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for custom_product_attributes")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for custom_product_attributes")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Attribute = foreign
		if foreign.R == nil {
			foreign.R = &customProductAttributeR{}
		}
		foreign.R.AttributeCustomProductAttributeValues = append(foreign.R.AttributeCustomProductAttributeValues, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.AttributeID == foreign.ID {
				local.R.Attribute = foreign
				if foreign.R == nil {
					foreign.R = &customProductAttributeR{}
				}
				foreign.R.AttributeCustomProductAttributeValues = append(foreign.R.AttributeCustomProductAttributeValues, local)
				break
			}
		}
	}

	return nil
}

// LoadAttributeValueAssignedProductVariantAttributeValues allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (customProductAttributeValueL) LoadAttributeValueAssignedProductVariantAttributeValues(e boil.Executor, singular bool, maybeCustomProductAttributeValue interface{}, mods queries.Applicator) error {
	var slice []*CustomProductAttributeValue
	var object *CustomProductAttributeValue

	if singular {
		var ok bool
		object, ok = maybeCustomProductAttributeValue.(*CustomProductAttributeValue)
		if !ok {
			object = new(CustomProductAttributeValue)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCustomProductAttributeValue)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCustomProductAttributeValue))
			}
		}
	} else {
		s, ok := maybeCustomProductAttributeValue.(*[]*CustomProductAttributeValue)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCustomProductAttributeValue)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCustomProductAttributeValue))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &customProductAttributeValueR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &customProductAttributeValueR{}
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
		qm.From(`assigned_product_variant_attribute_values`),
		qm.WhereIn(`assigned_product_variant_attribute_values.attribute_value_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load assigned_product_variant_attribute_values")
	}

	var resultSlice []*AssignedProductVariantAttributeValue
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice assigned_product_variant_attribute_values")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on assigned_product_variant_attribute_values")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for assigned_product_variant_attribute_values")
	}

	if singular {
		object.R.AttributeValueAssignedProductVariantAttributeValues = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &assignedProductVariantAttributeValueR{}
			}
			foreign.R.AttributeValue = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.AttributeValueID {
				local.R.AttributeValueAssignedProductVariantAttributeValues = append(local.R.AttributeValueAssignedProductVariantAttributeValues, foreign)
				if foreign.R == nil {
					foreign.R = &assignedProductVariantAttributeValueR{}
				}
				foreign.R.AttributeValue = local
				break
			}
		}
	}

	return nil
}

// SetAttribute of the customProductAttributeValue to the related item.
// Sets o.R.Attribute to related.
// Adds o to related.R.AttributeCustomProductAttributeValues.
func (o *CustomProductAttributeValue) SetAttribute(exec boil.Executor, insert bool, related *CustomProductAttribute) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"custom_product_attribute_values\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"attribute_id"}),
		strmangle.WhereClause("\"", "\"", 2, customProductAttributeValuePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.AttributeID = related.ID
	if o.R == nil {
		o.R = &customProductAttributeValueR{
			Attribute: related,
		}
	} else {
		o.R.Attribute = related
	}

	if related.R == nil {
		related.R = &customProductAttributeR{
			AttributeCustomProductAttributeValues: CustomProductAttributeValueSlice{o},
		}
	} else {
		related.R.AttributeCustomProductAttributeValues = append(related.R.AttributeCustomProductAttributeValues, o)
	}

	return nil
}

// AddAttributeValueAssignedProductVariantAttributeValues adds the given related objects to the existing relationships
// of the custom_product_attribute_value, optionally inserting them as new records.
// Appends related to o.R.AttributeValueAssignedProductVariantAttributeValues.
// Sets related.R.AttributeValue appropriately.
func (o *CustomProductAttributeValue) AddAttributeValueAssignedProductVariantAttributeValues(exec boil.Executor, insert bool, related ...*AssignedProductVariantAttributeValue) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.AttributeValueID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"assigned_product_variant_attribute_values\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"attribute_value_id"}),
				strmangle.WhereClause("\"", "\"", 2, assignedProductVariantAttributeValuePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.AttributeValueID = o.ID
		}
	}

	if o.R == nil {
		o.R = &customProductAttributeValueR{
			AttributeValueAssignedProductVariantAttributeValues: related,
		}
	} else {
		o.R.AttributeValueAssignedProductVariantAttributeValues = append(o.R.AttributeValueAssignedProductVariantAttributeValues, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &assignedProductVariantAttributeValueR{
				AttributeValue: o,
			}
		} else {
			rel.R.AttributeValue = o
		}
	}
	return nil
}

// CustomProductAttributeValues retrieves all the records using an executor.
func CustomProductAttributeValues(mods ...qm.QueryMod) customProductAttributeValueQuery {
	mods = append(mods, qm.From("\"custom_product_attribute_values\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"custom_product_attribute_values\".*"})
	}

	return customProductAttributeValueQuery{q}
}

// FindCustomProductAttributeValue retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCustomProductAttributeValue(exec boil.Executor, iD string, selectCols ...string) (*CustomProductAttributeValue, error) {
	customProductAttributeValueObj := &CustomProductAttributeValue{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"custom_product_attribute_values\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, customProductAttributeValueObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from custom_product_attribute_values")
	}

	return customProductAttributeValueObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *CustomProductAttributeValue) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no custom_product_attribute_values provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(customProductAttributeValueColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	customProductAttributeValueInsertCacheMut.RLock()
	cache, cached := customProductAttributeValueInsertCache[key]
	customProductAttributeValueInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			customProductAttributeValueAllColumns,
			customProductAttributeValueColumnsWithDefault,
			customProductAttributeValueColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(customProductAttributeValueType, customProductAttributeValueMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(customProductAttributeValueType, customProductAttributeValueMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"custom_product_attribute_values\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"custom_product_attribute_values\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into custom_product_attribute_values")
	}

	if !cached {
		customProductAttributeValueInsertCacheMut.Lock()
		customProductAttributeValueInsertCache[key] = cache
		customProductAttributeValueInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the CustomProductAttributeValue.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *CustomProductAttributeValue) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	customProductAttributeValueUpdateCacheMut.RLock()
	cache, cached := customProductAttributeValueUpdateCache[key]
	customProductAttributeValueUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			customProductAttributeValueAllColumns,
			customProductAttributeValuePrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update custom_product_attribute_values, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"custom_product_attribute_values\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, customProductAttributeValuePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(customProductAttributeValueType, customProductAttributeValueMapping, append(wl, customProductAttributeValuePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update custom_product_attribute_values row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for custom_product_attribute_values")
	}

	if !cached {
		customProductAttributeValueUpdateCacheMut.Lock()
		customProductAttributeValueUpdateCache[key] = cache
		customProductAttributeValueUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q customProductAttributeValueQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for custom_product_attribute_values")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for custom_product_attribute_values")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CustomProductAttributeValueSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), customProductAttributeValuePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"custom_product_attribute_values\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, customProductAttributeValuePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in customProductAttributeValue slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all customProductAttributeValue")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *CustomProductAttributeValue) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no custom_product_attribute_values provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(customProductAttributeValueColumnsWithDefault, o)

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

	customProductAttributeValueUpsertCacheMut.RLock()
	cache, cached := customProductAttributeValueUpsertCache[key]
	customProductAttributeValueUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			customProductAttributeValueAllColumns,
			customProductAttributeValueColumnsWithDefault,
			customProductAttributeValueColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			customProductAttributeValueAllColumns,
			customProductAttributeValuePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert custom_product_attribute_values, could not build update column list")
		}

		ret := strmangle.SetComplement(customProductAttributeValueAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(customProductAttributeValuePrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert custom_product_attribute_values, could not build conflict column list")
			}

			conflict = make([]string, len(customProductAttributeValuePrimaryKeyColumns))
			copy(conflict, customProductAttributeValuePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"custom_product_attribute_values\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(customProductAttributeValueType, customProductAttributeValueMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(customProductAttributeValueType, customProductAttributeValueMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert custom_product_attribute_values")
	}

	if !cached {
		customProductAttributeValueUpsertCacheMut.Lock()
		customProductAttributeValueUpsertCache[key] = cache
		customProductAttributeValueUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single CustomProductAttributeValue record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *CustomProductAttributeValue) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no CustomProductAttributeValue provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), customProductAttributeValuePrimaryKeyMapping)
	sql := "DELETE FROM \"custom_product_attribute_values\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from custom_product_attribute_values")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for custom_product_attribute_values")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q customProductAttributeValueQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no customProductAttributeValueQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from custom_product_attribute_values")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for custom_product_attribute_values")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CustomProductAttributeValueSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), customProductAttributeValuePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"custom_product_attribute_values\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, customProductAttributeValuePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from customProductAttributeValue slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for custom_product_attribute_values")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *CustomProductAttributeValue) Reload(exec boil.Executor) error {
	ret, err := FindCustomProductAttributeValue(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CustomProductAttributeValueSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CustomProductAttributeValueSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), customProductAttributeValuePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"custom_product_attribute_values\".* FROM \"custom_product_attribute_values\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, customProductAttributeValuePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in CustomProductAttributeValueSlice")
	}

	*o = slice

	return nil
}

// CustomProductAttributeValueExists checks if the CustomProductAttributeValue row exists.
func CustomProductAttributeValueExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"custom_product_attribute_values\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if custom_product_attribute_values exists")
	}

	return exists, nil
}

// Exists checks if the CustomProductAttributeValue row exists.
func (o *CustomProductAttributeValue) Exists(exec boil.Executor) (bool, error) {
	return CustomProductAttributeValueExists(exec, o.ID)
}
