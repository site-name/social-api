// Code generated by SQLBoiler 4.17.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// GiftcardTag is an object representing the database table.
type GiftcardTag struct {
	ID   string `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name string `boil:"name" json:"name" toml:"name" yaml:"name"`

	R *giftcardTagR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L giftcardTagL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GiftcardTagColumns = struct {
	ID   string
	Name string
}{
	ID:   "id",
	Name: "name",
}

var GiftcardTagTableColumns = struct {
	ID   string
	Name string
}{
	ID:   "giftcard_tags.id",
	Name: "giftcard_tags.name",
}

// Generated where

var GiftcardTagWhere = struct {
	ID   whereHelperstring
	Name whereHelperstring
}{
	ID:   whereHelperstring{field: "\"giftcard_tags\".\"id\""},
	Name: whereHelperstring{field: "\"giftcard_tags\".\"name\""},
}

// GiftcardTagRels is where relationship names are stored.
var GiftcardTagRels = struct {
	TagGiftcardTagGiftcards string
}{
	TagGiftcardTagGiftcards: "TagGiftcardTagGiftcards",
}

// giftcardTagR is where relationships are stored.
type giftcardTagR struct {
	TagGiftcardTagGiftcards GiftcardTagGiftcardSlice `boil:"TagGiftcardTagGiftcards" json:"TagGiftcardTagGiftcards" toml:"TagGiftcardTagGiftcards" yaml:"TagGiftcardTagGiftcards"`
}

// NewStruct creates a new relationship struct
func (*giftcardTagR) NewStruct() *giftcardTagR {
	return &giftcardTagR{}
}

func (r *giftcardTagR) GetTagGiftcardTagGiftcards() GiftcardTagGiftcardSlice {
	if r == nil {
		return nil
	}
	return r.TagGiftcardTagGiftcards
}

// giftcardTagL is where Load methods for each relationship are stored.
type giftcardTagL struct{}

var (
	giftcardTagAllColumns            = []string{"id", "name"}
	giftcardTagColumnsWithoutDefault = []string{"id", "name"}
	giftcardTagColumnsWithDefault    = []string{}
	giftcardTagPrimaryKeyColumns     = []string{"id"}
	giftcardTagGeneratedColumns      = []string{}
)

type (
	// GiftcardTagSlice is an alias for a slice of pointers to GiftcardTag.
	// This should almost always be used instead of []GiftcardTag.
	GiftcardTagSlice []*GiftcardTag

	giftcardTagQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	giftcardTagType                 = reflect.TypeOf(&GiftcardTag{})
	giftcardTagMapping              = queries.MakeStructMapping(giftcardTagType)
	giftcardTagPrimaryKeyMapping, _ = queries.BindMapping(giftcardTagType, giftcardTagMapping, giftcardTagPrimaryKeyColumns)
	giftcardTagInsertCacheMut       sync.RWMutex
	giftcardTagInsertCache          = make(map[string]insertCache)
	giftcardTagUpdateCacheMut       sync.RWMutex
	giftcardTagUpdateCache          = make(map[string]updateCache)
	giftcardTagUpsertCacheMut       sync.RWMutex
	giftcardTagUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single giftcardTag record from the query.
func (q giftcardTagQuery) One(exec boil.Executor) (*GiftcardTag, error) {
	o := &GiftcardTag{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for giftcard_tags")
	}

	return o, nil
}

// All returns all GiftcardTag records from the query.
func (q giftcardTagQuery) All(exec boil.Executor) (GiftcardTagSlice, error) {
	var o []*GiftcardTag

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to GiftcardTag slice")
	}

	return o, nil
}

// Count returns the count of all GiftcardTag records in the query.
func (q giftcardTagQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count giftcard_tags rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q giftcardTagQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if giftcard_tags exists")
	}

	return count > 0, nil
}

// TagGiftcardTagGiftcards retrieves all the giftcard_tag_giftcard's GiftcardTagGiftcards with an executor via tag_id column.
func (o *GiftcardTag) TagGiftcardTagGiftcards(mods ...qm.QueryMod) giftcardTagGiftcardQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"giftcard_tag_giftcards\".\"tag_id\"=?", o.ID),
	)

	return GiftcardTagGiftcards(queryMods...)
}

// LoadTagGiftcardTagGiftcards allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (giftcardTagL) LoadTagGiftcardTagGiftcards(e boil.Executor, singular bool, maybeGiftcardTag interface{}, mods queries.Applicator) error {
	var slice []*GiftcardTag
	var object *GiftcardTag

	if singular {
		var ok bool
		object, ok = maybeGiftcardTag.(*GiftcardTag)
		if !ok {
			object = new(GiftcardTag)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeGiftcardTag)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeGiftcardTag))
			}
		}
	} else {
		s, ok := maybeGiftcardTag.(*[]*GiftcardTag)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeGiftcardTag)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeGiftcardTag))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &giftcardTagR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &giftcardTagR{}
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
		qm.From(`giftcard_tag_giftcards`),
		qm.WhereIn(`giftcard_tag_giftcards.tag_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load giftcard_tag_giftcards")
	}

	var resultSlice []*GiftcardTagGiftcard
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice giftcard_tag_giftcards")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on giftcard_tag_giftcards")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for giftcard_tag_giftcards")
	}

	if singular {
		object.R.TagGiftcardTagGiftcards = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &giftcardTagGiftcardR{}
			}
			foreign.R.Tag = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.TagID {
				local.R.TagGiftcardTagGiftcards = append(local.R.TagGiftcardTagGiftcards, foreign)
				if foreign.R == nil {
					foreign.R = &giftcardTagGiftcardR{}
				}
				foreign.R.Tag = local
				break
			}
		}
	}

	return nil
}

// AddTagGiftcardTagGiftcards adds the given related objects to the existing relationships
// of the giftcard_tag, optionally inserting them as new records.
// Appends related to o.R.TagGiftcardTagGiftcards.
// Sets related.R.Tag appropriately.
func (o *GiftcardTag) AddTagGiftcardTagGiftcards(exec boil.Executor, insert bool, related ...*GiftcardTagGiftcard) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.TagID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"giftcard_tag_giftcards\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"tag_id"}),
				strmangle.WhereClause("\"", "\"", 2, giftcardTagGiftcardPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.TagID = o.ID
		}
	}

	if o.R == nil {
		o.R = &giftcardTagR{
			TagGiftcardTagGiftcards: related,
		}
	} else {
		o.R.TagGiftcardTagGiftcards = append(o.R.TagGiftcardTagGiftcards, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &giftcardTagGiftcardR{
				Tag: o,
			}
		} else {
			rel.R.Tag = o
		}
	}
	return nil
}

// GiftcardTags retrieves all the records using an executor.
func GiftcardTags(mods ...qm.QueryMod) giftcardTagQuery {
	mods = append(mods, qm.From("\"giftcard_tags\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"giftcard_tags\".*"})
	}

	return giftcardTagQuery{q}
}

// FindGiftcardTag retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGiftcardTag(exec boil.Executor, iD string, selectCols ...string) (*GiftcardTag, error) {
	giftcardTagObj := &GiftcardTag{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"giftcard_tags\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, giftcardTagObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from giftcard_tags")
	}

	return giftcardTagObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *GiftcardTag) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no giftcard_tags provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(giftcardTagColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	giftcardTagInsertCacheMut.RLock()
	cache, cached := giftcardTagInsertCache[key]
	giftcardTagInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			giftcardTagAllColumns,
			giftcardTagColumnsWithDefault,
			giftcardTagColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(giftcardTagType, giftcardTagMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(giftcardTagType, giftcardTagMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"giftcard_tags\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"giftcard_tags\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into giftcard_tags")
	}

	if !cached {
		giftcardTagInsertCacheMut.Lock()
		giftcardTagInsertCache[key] = cache
		giftcardTagInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the GiftcardTag.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *GiftcardTag) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	giftcardTagUpdateCacheMut.RLock()
	cache, cached := giftcardTagUpdateCache[key]
	giftcardTagUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			giftcardTagAllColumns,
			giftcardTagPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update giftcard_tags, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"giftcard_tags\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, giftcardTagPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(giftcardTagType, giftcardTagMapping, append(wl, giftcardTagPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update giftcard_tags row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for giftcard_tags")
	}

	if !cached {
		giftcardTagUpdateCacheMut.Lock()
		giftcardTagUpdateCache[key] = cache
		giftcardTagUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q giftcardTagQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for giftcard_tags")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for giftcard_tags")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GiftcardTagSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), giftcardTagPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"giftcard_tags\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, giftcardTagPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in giftcardTag slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all giftcardTag")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *GiftcardTag) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no giftcard_tags provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(giftcardTagColumnsWithDefault, o)

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

	giftcardTagUpsertCacheMut.RLock()
	cache, cached := giftcardTagUpsertCache[key]
	giftcardTagUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			giftcardTagAllColumns,
			giftcardTagColumnsWithDefault,
			giftcardTagColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			giftcardTagAllColumns,
			giftcardTagPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert giftcard_tags, could not build update column list")
		}

		ret := strmangle.SetComplement(giftcardTagAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(giftcardTagPrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert giftcard_tags, could not build conflict column list")
			}

			conflict = make([]string, len(giftcardTagPrimaryKeyColumns))
			copy(conflict, giftcardTagPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"giftcard_tags\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(giftcardTagType, giftcardTagMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(giftcardTagType, giftcardTagMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert giftcard_tags")
	}

	if !cached {
		giftcardTagUpsertCacheMut.Lock()
		giftcardTagUpsertCache[key] = cache
		giftcardTagUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single GiftcardTag record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *GiftcardTag) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no GiftcardTag provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), giftcardTagPrimaryKeyMapping)
	sql := "DELETE FROM \"giftcard_tags\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from giftcard_tags")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for giftcard_tags")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q giftcardTagQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no giftcardTagQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from giftcard_tags")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for giftcard_tags")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GiftcardTagSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), giftcardTagPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"giftcard_tags\" WHERE " +
		strmangle.WhereInClause(string(dialect.LQ), string(dialect.RQ), 1, giftcardTagPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from giftcardTag slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for giftcard_tags")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *GiftcardTag) Reload(exec boil.Executor) error {
	ret, err := FindGiftcardTag(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GiftcardTagSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GiftcardTagSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), giftcardTagPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"giftcard_tags\".* FROM \"giftcard_tags\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, giftcardTagPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in GiftcardTagSlice")
	}

	*o = slice

	return nil
}

// GiftcardTagExists checks if the GiftcardTag row exists.
func GiftcardTagExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"giftcard_tags\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if giftcard_tags exists")
	}

	return exists, nil
}

// Exists checks if the GiftcardTag row exists.
func (o *GiftcardTag) Exists(exec boil.Executor) (bool, error) {
	return GiftcardTagExists(exec, o.ID)
}
