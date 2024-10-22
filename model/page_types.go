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

// PageType is an object representing the database table.
type PageType struct {
	ID              string                 `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name            string                 `boil:"name" json:"name" toml:"name" yaml:"name"`
	Slug            string                 `boil:"slug" json:"slug" toml:"slug" yaml:"slug"`
	Metadata        model_types.JSONString `boil:"metadata" json:"metadata,omitempty" toml:"metadata" yaml:"metadata,omitempty"`
	PrivateMetadata model_types.JSONString `boil:"private_metadata" json:"private_metadata,omitempty" toml:"private_metadata" yaml:"private_metadata,omitempty"`

	R *pageTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L pageTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PageTypeColumns = struct {
	ID              string
	Name            string
	Slug            string
	Metadata        string
	PrivateMetadata string
}{
	ID:              "id",
	Name:            "name",
	Slug:            "slug",
	Metadata:        "metadata",
	PrivateMetadata: "private_metadata",
}

var PageTypeTableColumns = struct {
	ID              string
	Name            string
	Slug            string
	Metadata        string
	PrivateMetadata string
}{
	ID:              "page_types.id",
	Name:            "page_types.name",
	Slug:            "page_types.slug",
	Metadata:        "page_types.metadata",
	PrivateMetadata: "page_types.private_metadata",
}

// Generated where

var PageTypeWhere = struct {
	ID              whereHelperstring
	Name            whereHelperstring
	Slug            whereHelperstring
	Metadata        whereHelpermodel_types_JSONString
	PrivateMetadata whereHelpermodel_types_JSONString
}{
	ID:              whereHelperstring{field: "\"page_types\".\"id\""},
	Name:            whereHelperstring{field: "\"page_types\".\"name\""},
	Slug:            whereHelperstring{field: "\"page_types\".\"slug\""},
	Metadata:        whereHelpermodel_types_JSONString{field: "\"page_types\".\"metadata\""},
	PrivateMetadata: whereHelpermodel_types_JSONString{field: "\"page_types\".\"private_metadata\""},
}

// PageTypeRels is where relationship names are stored.
var PageTypeRels = struct {
	AttributePages string
	Pages          string
}{
	AttributePages: "AttributePages",
	Pages:          "Pages",
}

// pageTypeR is where relationships are stored.
type pageTypeR struct {
	AttributePages AttributePageSlice `boil:"AttributePages" json:"AttributePages" toml:"AttributePages" yaml:"AttributePages"`
	Pages          PageSlice          `boil:"Pages" json:"Pages" toml:"Pages" yaml:"Pages"`
}

// NewStruct creates a new relationship struct
func (*pageTypeR) NewStruct() *pageTypeR {
	return &pageTypeR{}
}

func (r *pageTypeR) GetAttributePages() AttributePageSlice {
	if r == nil {
		return nil
	}
	return r.AttributePages
}

func (r *pageTypeR) GetPages() PageSlice {
	if r == nil {
		return nil
	}
	return r.Pages
}

// pageTypeL is where Load methods for each relationship are stored.
type pageTypeL struct{}

var (
	pageTypeAllColumns            = []string{"id", "name", "slug", "metadata", "private_metadata"}
	pageTypeColumnsWithoutDefault = []string{"id", "name", "slug"}
	pageTypeColumnsWithDefault    = []string{"metadata", "private_metadata"}
	pageTypePrimaryKeyColumns     = []string{"id"}
	pageTypeGeneratedColumns      = []string{}
)

type (
	// PageTypeSlice is an alias for a slice of pointers to PageType.
	// This should almost always be used instead of []PageType.
	PageTypeSlice []*PageType

	pageTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	pageTypeType                 = reflect.TypeOf(&PageType{})
	pageTypeMapping              = queries.MakeStructMapping(pageTypeType)
	pageTypePrimaryKeyMapping, _ = queries.BindMapping(pageTypeType, pageTypeMapping, pageTypePrimaryKeyColumns)
	pageTypeInsertCacheMut       sync.RWMutex
	pageTypeInsertCache          = make(map[string]insertCache)
	pageTypeUpdateCacheMut       sync.RWMutex
	pageTypeUpdateCache          = make(map[string]updateCache)
	pageTypeUpsertCacheMut       sync.RWMutex
	pageTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single pageType record from the query.
func (q pageTypeQuery) One(exec boil.Executor) (*PageType, error) {
	o := &PageType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for page_types")
	}

	return o, nil
}

// All returns all PageType records from the query.
func (q pageTypeQuery) All(exec boil.Executor) (PageTypeSlice, error) {
	var o []*PageType

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to PageType slice")
	}

	return o, nil
}

// Count returns the count of all PageType records in the query.
func (q pageTypeQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count page_types rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q pageTypeQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if page_types exists")
	}

	return count > 0, nil
}

// AttributePages retrieves all the attribute_page's AttributePages with an executor.
func (o *PageType) AttributePages(mods ...qm.QueryMod) attributePageQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"attribute_pages\".\"page_type_id\"=?", o.ID),
	)

	return AttributePages(queryMods...)
}

// Pages retrieves all the page's Pages with an executor.
func (o *PageType) Pages(mods ...qm.QueryMod) pageQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"pages\".\"page_type_id\"=?", o.ID),
	)

	return Pages(queryMods...)
}

// LoadAttributePages allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (pageTypeL) LoadAttributePages(e boil.Executor, singular bool, maybePageType interface{}, mods queries.Applicator) error {
	var slice []*PageType
	var object *PageType

	if singular {
		var ok bool
		object, ok = maybePageType.(*PageType)
		if !ok {
			object = new(PageType)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePageType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePageType))
			}
		}
	} else {
		s, ok := maybePageType.(*[]*PageType)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePageType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePageType))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &pageTypeR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &pageTypeR{}
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
		qm.From(`attribute_pages`),
		qm.WhereIn(`attribute_pages.page_type_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load attribute_pages")
	}

	var resultSlice []*AttributePage
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice attribute_pages")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on attribute_pages")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for attribute_pages")
	}

	if singular {
		object.R.AttributePages = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &attributePageR{}
			}
			foreign.R.PageType = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.PageTypeID {
				local.R.AttributePages = append(local.R.AttributePages, foreign)
				if foreign.R == nil {
					foreign.R = &attributePageR{}
				}
				foreign.R.PageType = local
				break
			}
		}
	}

	return nil
}

// LoadPages allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (pageTypeL) LoadPages(e boil.Executor, singular bool, maybePageType interface{}, mods queries.Applicator) error {
	var slice []*PageType
	var object *PageType

	if singular {
		var ok bool
		object, ok = maybePageType.(*PageType)
		if !ok {
			object = new(PageType)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePageType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePageType))
			}
		}
	} else {
		s, ok := maybePageType.(*[]*PageType)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePageType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePageType))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &pageTypeR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &pageTypeR{}
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
		qm.From(`pages`),
		qm.WhereIn(`pages.page_type_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load pages")
	}

	var resultSlice []*Page
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice pages")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on pages")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for pages")
	}

	if singular {
		object.R.Pages = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &pageR{}
			}
			foreign.R.PageType = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.PageTypeID {
				local.R.Pages = append(local.R.Pages, foreign)
				if foreign.R == nil {
					foreign.R = &pageR{}
				}
				foreign.R.PageType = local
				break
			}
		}
	}

	return nil
}

// AddAttributePages adds the given related objects to the existing relationships
// of the page_type, optionally inserting them as new records.
// Appends related to o.R.AttributePages.
// Sets related.R.PageType appropriately.
func (o *PageType) AddAttributePages(exec boil.Executor, insert bool, related ...*AttributePage) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.PageTypeID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"attribute_pages\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"page_type_id"}),
				strmangle.WhereClause("\"", "\"", 2, attributePagePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.PageTypeID = o.ID
		}
	}

	if o.R == nil {
		o.R = &pageTypeR{
			AttributePages: related,
		}
	} else {
		o.R.AttributePages = append(o.R.AttributePages, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &attributePageR{
				PageType: o,
			}
		} else {
			rel.R.PageType = o
		}
	}
	return nil
}

// AddPages adds the given related objects to the existing relationships
// of the page_type, optionally inserting them as new records.
// Appends related to o.R.Pages.
// Sets related.R.PageType appropriately.
func (o *PageType) AddPages(exec boil.Executor, insert bool, related ...*Page) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.PageTypeID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"pages\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"page_type_id"}),
				strmangle.WhereClause("\"", "\"", 2, pagePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.PageTypeID = o.ID
		}
	}

	if o.R == nil {
		o.R = &pageTypeR{
			Pages: related,
		}
	} else {
		o.R.Pages = append(o.R.Pages, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &pageR{
				PageType: o,
			}
		} else {
			rel.R.PageType = o
		}
	}
	return nil
}

// PageTypes retrieves all the records using an executor.
func PageTypes(mods ...qm.QueryMod) pageTypeQuery {
	mods = append(mods, qm.From("\"page_types\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"page_types\".*"})
	}

	return pageTypeQuery{q}
}

// FindPageType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPageType(exec boil.Executor, iD string, selectCols ...string) (*PageType, error) {
	pageTypeObj := &PageType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"page_types\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, pageTypeObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from page_types")
	}

	return pageTypeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PageType) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no page_types provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(pageTypeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	pageTypeInsertCacheMut.RLock()
	cache, cached := pageTypeInsertCache[key]
	pageTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			pageTypeAllColumns,
			pageTypeColumnsWithDefault,
			pageTypeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(pageTypeType, pageTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(pageTypeType, pageTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"page_types\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"page_types\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into page_types")
	}

	if !cached {
		pageTypeInsertCacheMut.Lock()
		pageTypeInsertCache[key] = cache
		pageTypeInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the PageType.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PageType) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	pageTypeUpdateCacheMut.RLock()
	cache, cached := pageTypeUpdateCache[key]
	pageTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			pageTypeAllColumns,
			pageTypePrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update page_types, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"page_types\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, pageTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(pageTypeType, pageTypeMapping, append(wl, pageTypePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update page_types row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for page_types")
	}

	if !cached {
		pageTypeUpdateCacheMut.Lock()
		pageTypeUpdateCache[key] = cache
		pageTypeUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q pageTypeQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for page_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for page_types")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PageTypeSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pageTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"page_types\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, pageTypePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in pageType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all pageType")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PageType) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no page_types provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(pageTypeColumnsWithDefault, o)

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

	pageTypeUpsertCacheMut.RLock()
	cache, cached := pageTypeUpsertCache[key]
	pageTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			pageTypeAllColumns,
			pageTypeColumnsWithDefault,
			pageTypeColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			pageTypeAllColumns,
			pageTypePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert page_types, could not build update column list")
		}

		ret := strmangle.SetComplement(pageTypeAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(pageTypePrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert page_types, could not build conflict column list")
			}

			conflict = make([]string, len(pageTypePrimaryKeyColumns))
			copy(conflict, pageTypePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"page_types\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(pageTypeType, pageTypeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(pageTypeType, pageTypeMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert page_types")
	}

	if !cached {
		pageTypeUpsertCacheMut.Lock()
		pageTypeUpsertCache[key] = cache
		pageTypeUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single PageType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PageType) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no PageType provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), pageTypePrimaryKeyMapping)
	sql := "DELETE FROM \"page_types\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from page_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for page_types")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q pageTypeQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no pageTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from page_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for page_types")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PageTypeSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pageTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"page_types\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, pageTypePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from pageType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for page_types")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PageType) Reload(exec boil.Executor) error {
	ret, err := FindPageType(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PageTypeSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PageTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pageTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"page_types\".* FROM \"page_types\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, pageTypePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in PageTypeSlice")
	}

	*o = slice

	return nil
}

// PageTypeExists checks if the PageType row exists.
func PageTypeExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"page_types\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if page_types exists")
	}

	return exists, nil
}

// Exists checks if the PageType row exists.
func (o *PageType) Exists(exec boil.Executor) (bool, error) {
	return PageTypeExists(exec, o.ID)
}
