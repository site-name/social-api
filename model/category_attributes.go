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

// CategoryAttribute is an object representing the database table.
type CategoryAttribute struct {
	ID          string              `boil:"id" json:"id" toml:"id" yaml:"id"`
	AttributeID string              `boil:"attribute_id" json:"attribute_id" toml:"attribute_id" yaml:"attribute_id"`
	CategoryID  string              `boil:"category_id" json:"category_id" toml:"category_id" yaml:"category_id"`
	SortOrder   model_types.NullInt `boil:"sort_order" json:"sort_order,omitempty" toml:"sort_order" yaml:"sort_order,omitempty"`

	R *categoryAttributeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L categoryAttributeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CategoryAttributeColumns = struct {
	ID          string
	AttributeID string
	CategoryID  string
	SortOrder   string
}{
	ID:          "id",
	AttributeID: "attribute_id",
	CategoryID:  "category_id",
	SortOrder:   "sort_order",
}

var CategoryAttributeTableColumns = struct {
	ID          string
	AttributeID string
	CategoryID  string
	SortOrder   string
}{
	ID:          "category_attributes.id",
	AttributeID: "category_attributes.attribute_id",
	CategoryID:  "category_attributes.category_id",
	SortOrder:   "category_attributes.sort_order",
}

// Generated where

var CategoryAttributeWhere = struct {
	ID          whereHelperstring
	AttributeID whereHelperstring
	CategoryID  whereHelperstring
	SortOrder   whereHelpermodel_types_NullInt
}{
	ID:          whereHelperstring{field: "\"category_attributes\".\"id\""},
	AttributeID: whereHelperstring{field: "\"category_attributes\".\"attribute_id\""},
	CategoryID:  whereHelperstring{field: "\"category_attributes\".\"category_id\""},
	SortOrder:   whereHelpermodel_types_NullInt{field: "\"category_attributes\".\"sort_order\""},
}

// CategoryAttributeRels is where relationship names are stored.
var CategoryAttributeRels = struct {
	Attribute                           string
	Category                            string
	AssignmentAssignedProductAttributes string
}{
	Attribute:                           "Attribute",
	Category:                            "Category",
	AssignmentAssignedProductAttributes: "AssignmentAssignedProductAttributes",
}

// categoryAttributeR is where relationships are stored.
type categoryAttributeR struct {
	Attribute                           *Attribute                    `boil:"Attribute" json:"Attribute" toml:"Attribute" yaml:"Attribute"`
	Category                            *Category                     `boil:"Category" json:"Category" toml:"Category" yaml:"Category"`
	AssignmentAssignedProductAttributes AssignedProductAttributeSlice `boil:"AssignmentAssignedProductAttributes" json:"AssignmentAssignedProductAttributes" toml:"AssignmentAssignedProductAttributes" yaml:"AssignmentAssignedProductAttributes"`
}

// NewStruct creates a new relationship struct
func (*categoryAttributeR) NewStruct() *categoryAttributeR {
	return &categoryAttributeR{}
}

func (r *categoryAttributeR) GetAttribute() *Attribute {
	if r == nil {
		return nil
	}
	return r.Attribute
}

func (r *categoryAttributeR) GetCategory() *Category {
	if r == nil {
		return nil
	}
	return r.Category
}

func (r *categoryAttributeR) GetAssignmentAssignedProductAttributes() AssignedProductAttributeSlice {
	if r == nil {
		return nil
	}
	return r.AssignmentAssignedProductAttributes
}

// categoryAttributeL is where Load methods for each relationship are stored.
type categoryAttributeL struct{}

var (
	categoryAttributeAllColumns            = []string{"id", "attribute_id", "category_id", "sort_order"}
	categoryAttributeColumnsWithoutDefault = []string{"id", "attribute_id", "category_id"}
	categoryAttributeColumnsWithDefault    = []string{"sort_order"}
	categoryAttributePrimaryKeyColumns     = []string{"id"}
	categoryAttributeGeneratedColumns      = []string{}
)

type (
	// CategoryAttributeSlice is an alias for a slice of pointers to CategoryAttribute.
	// This should almost always be used instead of []CategoryAttribute.
	CategoryAttributeSlice []*CategoryAttribute

	categoryAttributeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	categoryAttributeType                 = reflect.TypeOf(&CategoryAttribute{})
	categoryAttributeMapping              = queries.MakeStructMapping(categoryAttributeType)
	categoryAttributePrimaryKeyMapping, _ = queries.BindMapping(categoryAttributeType, categoryAttributeMapping, categoryAttributePrimaryKeyColumns)
	categoryAttributeInsertCacheMut       sync.RWMutex
	categoryAttributeInsertCache          = make(map[string]insertCache)
	categoryAttributeUpdateCacheMut       sync.RWMutex
	categoryAttributeUpdateCache          = make(map[string]updateCache)
	categoryAttributeUpsertCacheMut       sync.RWMutex
	categoryAttributeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single categoryAttribute record from the query.
func (q categoryAttributeQuery) One(exec boil.Executor) (*CategoryAttribute, error) {
	o := &CategoryAttribute{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for category_attributes")
	}

	return o, nil
}

// All returns all CategoryAttribute records from the query.
func (q categoryAttributeQuery) All(exec boil.Executor) (CategoryAttributeSlice, error) {
	var o []*CategoryAttribute

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to CategoryAttribute slice")
	}

	return o, nil
}

// Count returns the count of all CategoryAttribute records in the query.
func (q categoryAttributeQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count category_attributes rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q categoryAttributeQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if category_attributes exists")
	}

	return count > 0, nil
}

// Attribute pointed to by the foreign key.
func (o *CategoryAttribute) Attribute(mods ...qm.QueryMod) attributeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.AttributeID),
	}

	queryMods = append(queryMods, mods...)

	return Attributes(queryMods...)
}

// Category pointed to by the foreign key.
func (o *CategoryAttribute) Category(mods ...qm.QueryMod) categoryQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.CategoryID),
	}

	queryMods = append(queryMods, mods...)

	return Categories(queryMods...)
}

// AssignmentAssignedProductAttributes retrieves all the assigned_product_attribute's AssignedProductAttributes with an executor via assignment_id column.
func (o *CategoryAttribute) AssignmentAssignedProductAttributes(mods ...qm.QueryMod) assignedProductAttributeQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"assigned_product_attributes\".\"assignment_id\"=?", o.ID),
	)

	return AssignedProductAttributes(queryMods...)
}

// LoadAttribute allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (categoryAttributeL) LoadAttribute(e boil.Executor, singular bool, maybeCategoryAttribute interface{}, mods queries.Applicator) error {
	var slice []*CategoryAttribute
	var object *CategoryAttribute

	if singular {
		var ok bool
		object, ok = maybeCategoryAttribute.(*CategoryAttribute)
		if !ok {
			object = new(CategoryAttribute)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCategoryAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCategoryAttribute))
			}
		}
	} else {
		s, ok := maybeCategoryAttribute.(*[]*CategoryAttribute)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCategoryAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCategoryAttribute))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &categoryAttributeR{}
		}
		args[object.AttributeID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &categoryAttributeR{}
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
		qm.From(`attributes`),
		qm.WhereIn(`attributes.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Attribute")
	}

	var resultSlice []*Attribute
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Attribute")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for attributes")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for attributes")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Attribute = foreign
		if foreign.R == nil {
			foreign.R = &attributeR{}
		}
		foreign.R.CategoryAttributes = append(foreign.R.CategoryAttributes, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.AttributeID == foreign.ID {
				local.R.Attribute = foreign
				if foreign.R == nil {
					foreign.R = &attributeR{}
				}
				foreign.R.CategoryAttributes = append(foreign.R.CategoryAttributes, local)
				break
			}
		}
	}

	return nil
}

// LoadCategory allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (categoryAttributeL) LoadCategory(e boil.Executor, singular bool, maybeCategoryAttribute interface{}, mods queries.Applicator) error {
	var slice []*CategoryAttribute
	var object *CategoryAttribute

	if singular {
		var ok bool
		object, ok = maybeCategoryAttribute.(*CategoryAttribute)
		if !ok {
			object = new(CategoryAttribute)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCategoryAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCategoryAttribute))
			}
		}
	} else {
		s, ok := maybeCategoryAttribute.(*[]*CategoryAttribute)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCategoryAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCategoryAttribute))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &categoryAttributeR{}
		}
		args[object.CategoryID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &categoryAttributeR{}
			}

			args[obj.CategoryID] = struct{}{}

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
		qm.From(`categories`),
		qm.WhereIn(`categories.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
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
		foreign.R.CategoryAttributes = append(foreign.R.CategoryAttributes, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CategoryID == foreign.ID {
				local.R.Category = foreign
				if foreign.R == nil {
					foreign.R = &categoryR{}
				}
				foreign.R.CategoryAttributes = append(foreign.R.CategoryAttributes, local)
				break
			}
		}
	}

	return nil
}

// LoadAssignmentAssignedProductAttributes allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (categoryAttributeL) LoadAssignmentAssignedProductAttributes(e boil.Executor, singular bool, maybeCategoryAttribute interface{}, mods queries.Applicator) error {
	var slice []*CategoryAttribute
	var object *CategoryAttribute

	if singular {
		var ok bool
		object, ok = maybeCategoryAttribute.(*CategoryAttribute)
		if !ok {
			object = new(CategoryAttribute)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCategoryAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCategoryAttribute))
			}
		}
	} else {
		s, ok := maybeCategoryAttribute.(*[]*CategoryAttribute)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCategoryAttribute)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCategoryAttribute))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &categoryAttributeR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &categoryAttributeR{}
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
		qm.From(`assigned_product_attributes`),
		qm.WhereIn(`assigned_product_attributes.assignment_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load assigned_product_attributes")
	}

	var resultSlice []*AssignedProductAttribute
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice assigned_product_attributes")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on assigned_product_attributes")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for assigned_product_attributes")
	}

	if singular {
		object.R.AssignmentAssignedProductAttributes = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &assignedProductAttributeR{}
			}
			foreign.R.Assignment = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.AssignmentID {
				local.R.AssignmentAssignedProductAttributes = append(local.R.AssignmentAssignedProductAttributes, foreign)
				if foreign.R == nil {
					foreign.R = &assignedProductAttributeR{}
				}
				foreign.R.Assignment = local
				break
			}
		}
	}

	return nil
}

// SetAttribute of the categoryAttribute to the related item.
// Sets o.R.Attribute to related.
// Adds o to related.R.CategoryAttributes.
func (o *CategoryAttribute) SetAttribute(exec boil.Executor, insert bool, related *Attribute) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"category_attributes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"attribute_id"}),
		strmangle.WhereClause("\"", "\"", 2, categoryAttributePrimaryKeyColumns),
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
		o.R = &categoryAttributeR{
			Attribute: related,
		}
	} else {
		o.R.Attribute = related
	}

	if related.R == nil {
		related.R = &attributeR{
			CategoryAttributes: CategoryAttributeSlice{o},
		}
	} else {
		related.R.CategoryAttributes = append(related.R.CategoryAttributes, o)
	}

	return nil
}

// SetCategory of the categoryAttribute to the related item.
// Sets o.R.Category to related.
// Adds o to related.R.CategoryAttributes.
func (o *CategoryAttribute) SetCategory(exec boil.Executor, insert bool, related *Category) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"category_attributes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"category_id"}),
		strmangle.WhereClause("\"", "\"", 2, categoryAttributePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CategoryID = related.ID
	if o.R == nil {
		o.R = &categoryAttributeR{
			Category: related,
		}
	} else {
		o.R.Category = related
	}

	if related.R == nil {
		related.R = &categoryR{
			CategoryAttributes: CategoryAttributeSlice{o},
		}
	} else {
		related.R.CategoryAttributes = append(related.R.CategoryAttributes, o)
	}

	return nil
}

// AddAssignmentAssignedProductAttributes adds the given related objects to the existing relationships
// of the category_attribute, optionally inserting them as new records.
// Appends related to o.R.AssignmentAssignedProductAttributes.
// Sets related.R.Assignment appropriately.
func (o *CategoryAttribute) AddAssignmentAssignedProductAttributes(exec boil.Executor, insert bool, related ...*AssignedProductAttribute) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.AssignmentID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"assigned_product_attributes\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"assignment_id"}),
				strmangle.WhereClause("\"", "\"", 2, assignedProductAttributePrimaryKeyColumns),
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
		o.R = &categoryAttributeR{
			AssignmentAssignedProductAttributes: related,
		}
	} else {
		o.R.AssignmentAssignedProductAttributes = append(o.R.AssignmentAssignedProductAttributes, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &assignedProductAttributeR{
				Assignment: o,
			}
		} else {
			rel.R.Assignment = o
		}
	}
	return nil
}

// CategoryAttributes retrieves all the records using an executor.
func CategoryAttributes(mods ...qm.QueryMod) categoryAttributeQuery {
	mods = append(mods, qm.From("\"category_attributes\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"category_attributes\".*"})
	}

	return categoryAttributeQuery{q}
}

// FindCategoryAttribute retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCategoryAttribute(exec boil.Executor, iD string, selectCols ...string) (*CategoryAttribute, error) {
	categoryAttributeObj := &CategoryAttribute{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"category_attributes\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, categoryAttributeObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from category_attributes")
	}

	return categoryAttributeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *CategoryAttribute) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no category_attributes provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(categoryAttributeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	categoryAttributeInsertCacheMut.RLock()
	cache, cached := categoryAttributeInsertCache[key]
	categoryAttributeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			categoryAttributeAllColumns,
			categoryAttributeColumnsWithDefault,
			categoryAttributeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(categoryAttributeType, categoryAttributeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(categoryAttributeType, categoryAttributeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"category_attributes\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"category_attributes\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into category_attributes")
	}

	if !cached {
		categoryAttributeInsertCacheMut.Lock()
		categoryAttributeInsertCache[key] = cache
		categoryAttributeInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the CategoryAttribute.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *CategoryAttribute) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	categoryAttributeUpdateCacheMut.RLock()
	cache, cached := categoryAttributeUpdateCache[key]
	categoryAttributeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			categoryAttributeAllColumns,
			categoryAttributePrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update category_attributes, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"category_attributes\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, categoryAttributePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(categoryAttributeType, categoryAttributeMapping, append(wl, categoryAttributePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update category_attributes row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for category_attributes")
	}

	if !cached {
		categoryAttributeUpdateCacheMut.Lock()
		categoryAttributeUpdateCache[key] = cache
		categoryAttributeUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q categoryAttributeQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for category_attributes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for category_attributes")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CategoryAttributeSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), categoryAttributePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"category_attributes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, categoryAttributePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in categoryAttribute slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all categoryAttribute")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *CategoryAttribute) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no category_attributes provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(categoryAttributeColumnsWithDefault, o)

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

	categoryAttributeUpsertCacheMut.RLock()
	cache, cached := categoryAttributeUpsertCache[key]
	categoryAttributeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			categoryAttributeAllColumns,
			categoryAttributeColumnsWithDefault,
			categoryAttributeColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			categoryAttributeAllColumns,
			categoryAttributePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert category_attributes, could not build update column list")
		}

		ret := strmangle.SetComplement(categoryAttributeAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(categoryAttributePrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert category_attributes, could not build conflict column list")
			}

			conflict = make([]string, len(categoryAttributePrimaryKeyColumns))
			copy(conflict, categoryAttributePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"category_attributes\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(categoryAttributeType, categoryAttributeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(categoryAttributeType, categoryAttributeMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert category_attributes")
	}

	if !cached {
		categoryAttributeUpsertCacheMut.Lock()
		categoryAttributeUpsertCache[key] = cache
		categoryAttributeUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single CategoryAttribute record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *CategoryAttribute) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no CategoryAttribute provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), categoryAttributePrimaryKeyMapping)
	sql := "DELETE FROM \"category_attributes\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from category_attributes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for category_attributes")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q categoryAttributeQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no categoryAttributeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from category_attributes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for category_attributes")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CategoryAttributeSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), categoryAttributePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"category_attributes\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, categoryAttributePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from categoryAttribute slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for category_attributes")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *CategoryAttribute) Reload(exec boil.Executor) error {
	ret, err := FindCategoryAttribute(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CategoryAttributeSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CategoryAttributeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), categoryAttributePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"category_attributes\".* FROM \"category_attributes\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, categoryAttributePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in CategoryAttributeSlice")
	}

	*o = slice

	return nil
}

// CategoryAttributeExists checks if the CategoryAttribute row exists.
func CategoryAttributeExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"category_attributes\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if category_attributes exists")
	}

	return exists, nil
}

// Exists checks if the CategoryAttribute row exists.
func (o *CategoryAttribute) Exists(exec boil.Executor) (bool, error) {
	return CategoryAttributeExists(exec, o.ID)
}