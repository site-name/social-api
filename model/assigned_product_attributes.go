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

// AssignedProductAttribute is an object representing the database table.
type AssignedProductAttribute struct {
	ID           string `boil:"id" json:"id" toml:"id" yaml:"id"`
	ProductID    string `boil:"product_id" json:"product_id" toml:"product_id" yaml:"product_id"`
	AssignmentID string `boil:"assignment_id" json:"assignment_id" toml:"assignment_id" yaml:"assignment_id"`

	R *assignedProductAttributeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L assignedProductAttributeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AssignedProductAttributeColumns = struct {
	ID           string
	ProductID    string
	AssignmentID string
}{
	ID:           "id",
	ProductID:    "product_id",
	AssignmentID: "assignment_id",
}

var AssignedProductAttributeTableColumns = struct {
	ID           string
	ProductID    string
	AssignmentID string
}{
	ID:           "assigned_product_attributes.id",
	ProductID:    "assigned_product_attributes.product_id",
	AssignmentID: "assigned_product_attributes.assignment_id",
}

// Generated where

var AssignedProductAttributeWhere = struct {
	ID           whereHelperstring
	ProductID    whereHelperstring
	AssignmentID whereHelperstring
}{
	ID:           whereHelperstring{field: "\"assigned_product_attributes\".\"id\""},
	ProductID:    whereHelperstring{field: "\"assigned_product_attributes\".\"product_id\""},
	AssignmentID: whereHelperstring{field: "\"assigned_product_attributes\".\"assignment_id\""},
}

// AssignedProductAttributeRels is where relationship names are stored.
var AssignedProductAttributeRels = struct {
	Assignment                               string
	Product                                  string
	AssignmentAssignedProductAttributeValues string
}{
	Assignment:                               "Assignment",
	Product:                                  "Product",
	AssignmentAssignedProductAttributeValues: "AssignmentAssignedProductAttributeValues",
}

// assignedProductAttributeR is where relationships are stored.
type assignedProductAttributeR struct {
	Assignment                               *CategoryAttribute                 `boil:"Assignment" json:"Assignment" toml:"Assignment" yaml:"Assignment"`
	Product                                  *Product                           `boil:"Product" json:"Product" toml:"Product" yaml:"Product"`
	AssignmentAssignedProductAttributeValues AssignedProductAttributeValueSlice `boil:"AssignmentAssignedProductAttributeValues" json:"AssignmentAssignedProductAttributeValues" toml:"AssignmentAssignedProductAttributeValues" yaml:"AssignmentAssignedProductAttributeValues"`
}

// NewStruct creates a new relationship struct
func (*assignedProductAttributeR) NewStruct() *assignedProductAttributeR {
	return &assignedProductAttributeR{}
}

func (r *assignedProductAttributeR) GetAssignment() *CategoryAttribute {
	if r == nil {
		return nil
	}
	return r.Assignment
}

func (r *assignedProductAttributeR) GetProduct() *Product {
	if r == nil {
		return nil
	}
	return r.Product
}

func (r *assignedProductAttributeR) GetAssignmentAssignedProductAttributeValues() AssignedProductAttributeValueSlice {
	if r == nil {
		return nil
	}
	return r.AssignmentAssignedProductAttributeValues
}

// assignedProductAttributeL is where Load methods for each relationship are stored.
type assignedProductAttributeL struct{}

var (
	assignedProductAttributeAllColumns            = []string{"id", "product_id", "assignment_id"}
	assignedProductAttributeColumnsWithoutDefault = []string{"id", "product_id", "assignment_id"}
	assignedProductAttributeColumnsWithDefault    = []string{}
	assignedProductAttributePrimaryKeyColumns     = []string{"id"}
	assignedProductAttributeGeneratedColumns      = []string{}
)

type (
	// AssignedProductAttributeSlice is an alias for a slice of pointers to AssignedProductAttribute.
	// This should almost always be used instead of []AssignedProductAttribute.
	AssignedProductAttributeSlice []*AssignedProductAttribute

	assignedProductAttributeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	assignedProductAttributeType                 = reflect.TypeOf(&AssignedProductAttribute{})
	assignedProductAttributeMapping              = queries.MakeStructMapping(assignedProductAttributeType)
	assignedProductAttributePrimaryKeyMapping, _ = queries.BindMapping(assignedProductAttributeType, assignedProductAttributeMapping, assignedProductAttributePrimaryKeyColumns)
	assignedProductAttributeInsertCacheMut       sync.RWMutex
	assignedProductAttributeInsertCache          = make(map[string]insertCache)
	assignedProductAttributeUpdateCacheMut       sync.RWMutex
	assignedProductAttributeUpdateCache          = make(map[string]updateCache)
	assignedProductAttributeUpsertCacheMut       sync.RWMutex
	assignedProductAttributeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single assignedProductAttribute record from the query.
func (q assignedProductAttributeQuery) One(exec boil.Executor) (*AssignedProductAttribute, error) {
	o := &AssignedProductAttribute{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for assigned_product_attributes")
	}

	return o, nil
}

// All returns all AssignedProductAttribute records from the query.
func (q assignedProductAttributeQuery) All(exec boil.Executor) (AssignedProductAttributeSlice, error) {
	var o []*AssignedProductAttribute

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to AssignedProductAttribute slice")
	}

	return o, nil
}

// Count returns the count of all AssignedProductAttribute records in the query.
func (q assignedProductAttributeQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count assigned_product_attributes rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q assignedProductAttributeQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if assigned_product_attributes exists")
	}

	return count > 0, nil
}

// Assignment pointed to by the foreign key.
func (o *AssignedProductAttribute) Assignment(mods ...qm.QueryMod) categoryAttributeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.AssignmentID),
	}

	queryMods = append(queryMods, mods...)

	return CategoryAttributes(queryMods...)
}

// Product pointed to by the foreign key.
func (o *AssignedProductAttribute) Product(mods ...qm.QueryMod) productQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ProductID),
	}

	queryMods = append(queryMods, mods...)

	return Products(queryMods...)
}

// AssignmentAssignedProductAttributeValues retrieves all the assigned_product_attribute_value's AssignedProductAttributeValues with an executor via assignment_id column.
func (o *AssignedProductAttribute) AssignmentAssignedProductAttributeValues(mods ...qm.QueryMod) assignedProductAttributeValueQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"assigned_product_attribute_values\".\"assignment_id\"=?", o.ID),
	)

	return AssignedProductAttributeValues(queryMods...)
}

// LoadAssignment allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (assignedProductAttributeL) LoadAssignment(e boil.Executor, singular bool, maybeAssignedProductAttribute interface{}, mods queries.Applicator) error {
	var slice []*AssignedProductAttribute
	var object *AssignedProductAttribute

	if singular {
		var ok bool
		object, ok = maybeAssignedProductAttribute.(*AssignedProductAttribute)
		if !ok {
			object = new(AssignedProductAttribute)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAssignedProductAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAssignedProductAttribute))
			}
		}
	} else {
		s, ok := maybeAssignedProductAttribute.(*[]*AssignedProductAttribute)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAssignedProductAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAssignedProductAttribute))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &assignedProductAttributeR{}
		}
		args[object.AssignmentID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &assignedProductAttributeR{}
			}

			args[obj.AssignmentID] = struct{}{}

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
		qm.From(`category_attributes`),
		qm.WhereIn(`category_attributes.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load CategoryAttribute")
	}

	var resultSlice []*CategoryAttribute
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice CategoryAttribute")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for category_attributes")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for category_attributes")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Assignment = foreign
		if foreign.R == nil {
			foreign.R = &categoryAttributeR{}
		}
		foreign.R.AssignmentAssignedProductAttributes = append(foreign.R.AssignmentAssignedProductAttributes, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.AssignmentID == foreign.ID {
				local.R.Assignment = foreign
				if foreign.R == nil {
					foreign.R = &categoryAttributeR{}
				}
				foreign.R.AssignmentAssignedProductAttributes = append(foreign.R.AssignmentAssignedProductAttributes, local)
				break
			}
		}
	}

	return nil
}

// LoadProduct allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (assignedProductAttributeL) LoadProduct(e boil.Executor, singular bool, maybeAssignedProductAttribute interface{}, mods queries.Applicator) error {
	var slice []*AssignedProductAttribute
	var object *AssignedProductAttribute

	if singular {
		var ok bool
		object, ok = maybeAssignedProductAttribute.(*AssignedProductAttribute)
		if !ok {
			object = new(AssignedProductAttribute)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAssignedProductAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAssignedProductAttribute))
			}
		}
	} else {
		s, ok := maybeAssignedProductAttribute.(*[]*AssignedProductAttribute)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAssignedProductAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAssignedProductAttribute))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &assignedProductAttributeR{}
		}
		args[object.ProductID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &assignedProductAttributeR{}
			}

			args[obj.ProductID] = struct{}{}

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
		qm.From(`products`),
		qm.WhereIn(`products.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Product")
	}

	var resultSlice []*Product
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Product")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for products")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for products")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Product = foreign
		if foreign.R == nil {
			foreign.R = &productR{}
		}
		foreign.R.AssignedProductAttributes = append(foreign.R.AssignedProductAttributes, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ProductID == foreign.ID {
				local.R.Product = foreign
				if foreign.R == nil {
					foreign.R = &productR{}
				}
				foreign.R.AssignedProductAttributes = append(foreign.R.AssignedProductAttributes, local)
				break
			}
		}
	}

	return nil
}

// LoadAssignmentAssignedProductAttributeValues allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (assignedProductAttributeL) LoadAssignmentAssignedProductAttributeValues(e boil.Executor, singular bool, maybeAssignedProductAttribute interface{}, mods queries.Applicator) error {
	var slice []*AssignedProductAttribute
	var object *AssignedProductAttribute

	if singular {
		var ok bool
		object, ok = maybeAssignedProductAttribute.(*AssignedProductAttribute)
		if !ok {
			object = new(AssignedProductAttribute)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAssignedProductAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAssignedProductAttribute))
			}
		}
	} else {
		s, ok := maybeAssignedProductAttribute.(*[]*AssignedProductAttribute)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAssignedProductAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAssignedProductAttribute))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &assignedProductAttributeR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &assignedProductAttributeR{}
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
		qm.From(`assigned_product_attribute_values`),
		qm.WhereIn(`assigned_product_attribute_values.assignment_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load assigned_product_attribute_values")
	}

	var resultSlice []*AssignedProductAttributeValue
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice assigned_product_attribute_values")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on assigned_product_attribute_values")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for assigned_product_attribute_values")
	}

	if singular {
		object.R.AssignmentAssignedProductAttributeValues = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &assignedProductAttributeValueR{}
			}
			foreign.R.Assignment = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.AssignmentID {
				local.R.AssignmentAssignedProductAttributeValues = append(local.R.AssignmentAssignedProductAttributeValues, foreign)
				if foreign.R == nil {
					foreign.R = &assignedProductAttributeValueR{}
				}
				foreign.R.Assignment = local
				break
			}
		}
	}

	return nil
}

// SetAssignment of the assignedProductAttribute to the related item.
// Sets o.R.Assignment to related.
// Adds o to related.R.AssignmentAssignedProductAttributes.
func (o *AssignedProductAttribute) SetAssignment(exec boil.Executor, insert bool, related *CategoryAttribute) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"assigned_product_attributes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"assignment_id"}),
		strmangle.WhereClause("\"", "\"", 2, assignedProductAttributePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.AssignmentID = related.ID
	if o.R == nil {
		o.R = &assignedProductAttributeR{
			Assignment: related,
		}
	} else {
		o.R.Assignment = related
	}

	if related.R == nil {
		related.R = &categoryAttributeR{
			AssignmentAssignedProductAttributes: AssignedProductAttributeSlice{o},
		}
	} else {
		related.R.AssignmentAssignedProductAttributes = append(related.R.AssignmentAssignedProductAttributes, o)
	}

	return nil
}

// SetProduct of the assignedProductAttribute to the related item.
// Sets o.R.Product to related.
// Adds o to related.R.AssignedProductAttributes.
func (o *AssignedProductAttribute) SetProduct(exec boil.Executor, insert bool, related *Product) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"assigned_product_attributes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"product_id"}),
		strmangle.WhereClause("\"", "\"", 2, assignedProductAttributePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ProductID = related.ID
	if o.R == nil {
		o.R = &assignedProductAttributeR{
			Product: related,
		}
	} else {
		o.R.Product = related
	}

	if related.R == nil {
		related.R = &productR{
			AssignedProductAttributes: AssignedProductAttributeSlice{o},
		}
	} else {
		related.R.AssignedProductAttributes = append(related.R.AssignedProductAttributes, o)
	}

	return nil
}

// AddAssignmentAssignedProductAttributeValues adds the given related objects to the existing relationships
// of the assigned_product_attribute, optionally inserting them as new records.
// Appends related to o.R.AssignmentAssignedProductAttributeValues.
// Sets related.R.Assignment appropriately.
func (o *AssignedProductAttribute) AddAssignmentAssignedProductAttributeValues(exec boil.Executor, insert bool, related ...*AssignedProductAttributeValue) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.AssignmentID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"assigned_product_attribute_values\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"assignment_id"}),
				strmangle.WhereClause("\"", "\"", 2, assignedProductAttributeValuePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.AssignmentID = o.ID
		}
	}

	if o.R == nil {
		o.R = &assignedProductAttributeR{
			AssignmentAssignedProductAttributeValues: related,
		}
	} else {
		o.R.AssignmentAssignedProductAttributeValues = append(o.R.AssignmentAssignedProductAttributeValues, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &assignedProductAttributeValueR{
				Assignment: o,
			}
		} else {
			rel.R.Assignment = o
		}
	}
	return nil
}

// AssignedProductAttributes retrieves all the records using an executor.
func AssignedProductAttributes(mods ...qm.QueryMod) assignedProductAttributeQuery {
	mods = append(mods, qm.From("\"assigned_product_attributes\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"assigned_product_attributes\".*"})
	}

	return assignedProductAttributeQuery{q}
}

// FindAssignedProductAttribute retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAssignedProductAttribute(exec boil.Executor, iD string, selectCols ...string) (*AssignedProductAttribute, error) {
	assignedProductAttributeObj := &AssignedProductAttribute{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"assigned_product_attributes\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, assignedProductAttributeObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from assigned_product_attributes")
	}

	return assignedProductAttributeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *AssignedProductAttribute) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no assigned_product_attributes provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(assignedProductAttributeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	assignedProductAttributeInsertCacheMut.RLock()
	cache, cached := assignedProductAttributeInsertCache[key]
	assignedProductAttributeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			assignedProductAttributeAllColumns,
			assignedProductAttributeColumnsWithDefault,
			assignedProductAttributeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(assignedProductAttributeType, assignedProductAttributeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(assignedProductAttributeType, assignedProductAttributeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"assigned_product_attributes\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"assigned_product_attributes\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into assigned_product_attributes")
	}

	if !cached {
		assignedProductAttributeInsertCacheMut.Lock()
		assignedProductAttributeInsertCache[key] = cache
		assignedProductAttributeInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the AssignedProductAttribute.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *AssignedProductAttribute) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	assignedProductAttributeUpdateCacheMut.RLock()
	cache, cached := assignedProductAttributeUpdateCache[key]
	assignedProductAttributeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			assignedProductAttributeAllColumns,
			assignedProductAttributePrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update assigned_product_attributes, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"assigned_product_attributes\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, assignedProductAttributePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(assignedProductAttributeType, assignedProductAttributeMapping, append(wl, assignedProductAttributePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update assigned_product_attributes row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for assigned_product_attributes")
	}

	if !cached {
		assignedProductAttributeUpdateCacheMut.Lock()
		assignedProductAttributeUpdateCache[key] = cache
		assignedProductAttributeUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q assignedProductAttributeQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for assigned_product_attributes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for assigned_product_attributes")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AssignedProductAttributeSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), assignedProductAttributePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"assigned_product_attributes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, assignedProductAttributePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in assignedProductAttribute slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all assignedProductAttribute")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *AssignedProductAttribute) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no assigned_product_attributes provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(assignedProductAttributeColumnsWithDefault, o)

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

	assignedProductAttributeUpsertCacheMut.RLock()
	cache, cached := assignedProductAttributeUpsertCache[key]
	assignedProductAttributeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			assignedProductAttributeAllColumns,
			assignedProductAttributeColumnsWithDefault,
			assignedProductAttributeColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			assignedProductAttributeAllColumns,
			assignedProductAttributePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert assigned_product_attributes, could not build update column list")
		}

		ret := strmangle.SetComplement(assignedProductAttributeAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(assignedProductAttributePrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert assigned_product_attributes, could not build conflict column list")
			}

			conflict = make([]string, len(assignedProductAttributePrimaryKeyColumns))
			copy(conflict, assignedProductAttributePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"assigned_product_attributes\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(assignedProductAttributeType, assignedProductAttributeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(assignedProductAttributeType, assignedProductAttributeMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert assigned_product_attributes")
	}

	if !cached {
		assignedProductAttributeUpsertCacheMut.Lock()
		assignedProductAttributeUpsertCache[key] = cache
		assignedProductAttributeUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single AssignedProductAttribute record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *AssignedProductAttribute) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no AssignedProductAttribute provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), assignedProductAttributePrimaryKeyMapping)
	sql := "DELETE FROM \"assigned_product_attributes\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from assigned_product_attributes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for assigned_product_attributes")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q assignedProductAttributeQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no assignedProductAttributeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from assigned_product_attributes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for assigned_product_attributes")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AssignedProductAttributeSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), assignedProductAttributePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"assigned_product_attributes\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, assignedProductAttributePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from assignedProductAttribute slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for assigned_product_attributes")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *AssignedProductAttribute) Reload(exec boil.Executor) error {
	ret, err := FindAssignedProductAttribute(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AssignedProductAttributeSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AssignedProductAttributeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), assignedProductAttributePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"assigned_product_attributes\".* FROM \"assigned_product_attributes\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, assignedProductAttributePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in AssignedProductAttributeSlice")
	}

	*o = slice

	return nil
}

// AssignedProductAttributeExists checks if the AssignedProductAttribute row exists.
func AssignedProductAttributeExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"assigned_product_attributes\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if assigned_product_attributes exists")
	}

	return exists, nil
}

// Exists checks if the AssignedProductAttribute row exists.
func (o *AssignedProductAttribute) Exists(exec boil.Executor) (bool, error) {
	return AssignedProductAttributeExists(exec, o.ID)
}
