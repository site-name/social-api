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

// Wishlist is an object representing the database table.
type Wishlist struct {
	ID        string `boil:"id" json:"id" toml:"id" yaml:"id"`
	Token     string `boil:"token" json:"token" toml:"token" yaml:"token"`
	UserID    string `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	CreatedAt int64  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *wishlistR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L wishlistL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var WishlistColumns = struct {
	ID        string
	Token     string
	UserID    string
	CreatedAt string
}{
	ID:        "id",
	Token:     "token",
	UserID:    "user_id",
	CreatedAt: "created_at",
}

var WishlistTableColumns = struct {
	ID        string
	Token     string
	UserID    string
	CreatedAt string
}{
	ID:        "wishlists.id",
	Token:     "wishlists.token",
	UserID:    "wishlists.user_id",
	CreatedAt: "wishlists.created_at",
}

// Generated where

var WishlistWhere = struct {
	ID        whereHelperstring
	Token     whereHelperstring
	UserID    whereHelperstring
	CreatedAt whereHelperint64
}{
	ID:        whereHelperstring{field: "\"wishlists\".\"id\""},
	Token:     whereHelperstring{field: "\"wishlists\".\"token\""},
	UserID:    whereHelperstring{field: "\"wishlists\".\"user_id\""},
	CreatedAt: whereHelperint64{field: "\"wishlists\".\"created_at\""},
}

// WishlistRels is where relationship names are stored.
var WishlistRels = struct {
	User          string
	WishlistItems string
}{
	User:          "User",
	WishlistItems: "WishlistItems",
}

// wishlistR is where relationships are stored.
type wishlistR struct {
	User          *User             `boil:"User" json:"User" toml:"User" yaml:"User"`
	WishlistItems WishlistItemSlice `boil:"WishlistItems" json:"WishlistItems" toml:"WishlistItems" yaml:"WishlistItems"`
}

// NewStruct creates a new relationship struct
func (*wishlistR) NewStruct() *wishlistR {
	return &wishlistR{}
}

func (r *wishlistR) GetUser() *User {
	if r == nil {
		return nil
	}
	return r.User
}

func (r *wishlistR) GetWishlistItems() WishlistItemSlice {
	if r == nil {
		return nil
	}
	return r.WishlistItems
}

// wishlistL is where Load methods for each relationship are stored.
type wishlistL struct{}

var (
	wishlistAllColumns            = []string{"id", "token", "user_id", "created_at"}
	wishlistColumnsWithoutDefault = []string{"id", "token", "user_id", "created_at"}
	wishlistColumnsWithDefault    = []string{}
	wishlistPrimaryKeyColumns     = []string{"id"}
	wishlistGeneratedColumns      = []string{}
)

type (
	// WishlistSlice is an alias for a slice of pointers to Wishlist.
	// This should almost always be used instead of []Wishlist.
	WishlistSlice []*Wishlist

	wishlistQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	wishlistType                 = reflect.TypeOf(&Wishlist{})
	wishlistMapping              = queries.MakeStructMapping(wishlistType)
	wishlistPrimaryKeyMapping, _ = queries.BindMapping(wishlistType, wishlistMapping, wishlistPrimaryKeyColumns)
	wishlistInsertCacheMut       sync.RWMutex
	wishlistInsertCache          = make(map[string]insertCache)
	wishlistUpdateCacheMut       sync.RWMutex
	wishlistUpdateCache          = make(map[string]updateCache)
	wishlistUpsertCacheMut       sync.RWMutex
	wishlistUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single wishlist record from the query.
func (q wishlistQuery) One(exec boil.Executor) (*Wishlist, error) {
	o := &Wishlist{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for wishlists")
	}

	return o, nil
}

// All returns all Wishlist records from the query.
func (q wishlistQuery) All(exec boil.Executor) (WishlistSlice, error) {
	var o []*Wishlist

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to Wishlist slice")
	}

	return o, nil
}

// Count returns the count of all Wishlist records in the query.
func (q wishlistQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count wishlists rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q wishlistQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if wishlists exists")
	}

	return count > 0, nil
}

// User pointed to by the foreign key.
func (o *Wishlist) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	return Users(queryMods...)
}

// WishlistItems retrieves all the wishlist_item's WishlistItems with an executor.
func (o *Wishlist) WishlistItems(mods ...qm.QueryMod) wishlistItemQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"wishlist_items\".\"wishlist_id\"=?", o.ID),
	)

	return WishlistItems(queryMods...)
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (wishlistL) LoadUser(e boil.Executor, singular bool, maybeWishlist interface{}, mods queries.Applicator) error {
	var slice []*Wishlist
	var object *Wishlist

	if singular {
		var ok bool
		object, ok = maybeWishlist.(*Wishlist)
		if !ok {
			object = new(Wishlist)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeWishlist)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeWishlist))
			}
		}
	} else {
		s, ok := maybeWishlist.(*[]*Wishlist)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeWishlist)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeWishlist))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &wishlistR{}
		}
		args[object.UserID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &wishlistR{}
			}

			args[obj.UserID] = struct{}{}

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
		qm.From(`users`),
		qm.WhereIn(`users.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for users")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for users")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &userR{}
		}
		foreign.R.Wishlist = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.Wishlist = local
				break
			}
		}
	}

	return nil
}

// LoadWishlistItems allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (wishlistL) LoadWishlistItems(e boil.Executor, singular bool, maybeWishlist interface{}, mods queries.Applicator) error {
	var slice []*Wishlist
	var object *Wishlist

	if singular {
		var ok bool
		object, ok = maybeWishlist.(*Wishlist)
		if !ok {
			object = new(Wishlist)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeWishlist)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeWishlist))
			}
		}
	} else {
		s, ok := maybeWishlist.(*[]*Wishlist)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeWishlist)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeWishlist))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &wishlistR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &wishlistR{}
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
		qm.From(`wishlist_items`),
		qm.WhereIn(`wishlist_items.wishlist_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load wishlist_items")
	}

	var resultSlice []*WishlistItem
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice wishlist_items")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on wishlist_items")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for wishlist_items")
	}

	if singular {
		object.R.WishlistItems = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &wishlistItemR{}
			}
			foreign.R.Wishlist = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.WishlistID {
				local.R.WishlistItems = append(local.R.WishlistItems, foreign)
				if foreign.R == nil {
					foreign.R = &wishlistItemR{}
				}
				foreign.R.Wishlist = local
				break
			}
		}
	}

	return nil
}

// SetUser of the wishlist to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Wishlist.
func (o *Wishlist) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"wishlists\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, wishlistPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID = related.ID
	if o.R == nil {
		o.R = &wishlistR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			Wishlist: o,
		}
	} else {
		related.R.Wishlist = o
	}

	return nil
}

// AddWishlistItems adds the given related objects to the existing relationships
// of the wishlist, optionally inserting them as new records.
// Appends related to o.R.WishlistItems.
// Sets related.R.Wishlist appropriately.
func (o *Wishlist) AddWishlistItems(exec boil.Executor, insert bool, related ...*WishlistItem) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.WishlistID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"wishlist_items\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"wishlist_id"}),
				strmangle.WhereClause("\"", "\"", 2, wishlistItemPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.WishlistID = o.ID
		}
	}

	if o.R == nil {
		o.R = &wishlistR{
			WishlistItems: related,
		}
	} else {
		o.R.WishlistItems = append(o.R.WishlistItems, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &wishlistItemR{
				Wishlist: o,
			}
		} else {
			rel.R.Wishlist = o
		}
	}
	return nil
}

// Wishlists retrieves all the records using an executor.
func Wishlists(mods ...qm.QueryMod) wishlistQuery {
	mods = append(mods, qm.From("\"wishlists\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"wishlists\".*"})
	}

	return wishlistQuery{q}
}

// FindWishlist retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindWishlist(exec boil.Executor, iD string, selectCols ...string) (*Wishlist, error) {
	wishlistObj := &Wishlist{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"wishlists\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, wishlistObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from wishlists")
	}

	return wishlistObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Wishlist) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no wishlists provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(wishlistColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	wishlistInsertCacheMut.RLock()
	cache, cached := wishlistInsertCache[key]
	wishlistInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			wishlistAllColumns,
			wishlistColumnsWithDefault,
			wishlistColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(wishlistType, wishlistMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(wishlistType, wishlistMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"wishlists\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"wishlists\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into wishlists")
	}

	if !cached {
		wishlistInsertCacheMut.Lock()
		wishlistInsertCache[key] = cache
		wishlistInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Wishlist.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Wishlist) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	wishlistUpdateCacheMut.RLock()
	cache, cached := wishlistUpdateCache[key]
	wishlistUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			wishlistAllColumns,
			wishlistPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update wishlists, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"wishlists\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, wishlistPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(wishlistType, wishlistMapping, append(wl, wishlistPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update wishlists row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for wishlists")
	}

	if !cached {
		wishlistUpdateCacheMut.Lock()
		wishlistUpdateCache[key] = cache
		wishlistUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q wishlistQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for wishlists")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for wishlists")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o WishlistSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), wishlistPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"wishlists\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, wishlistPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in wishlist slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all wishlist")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Wishlist) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no wishlists provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(wishlistColumnsWithDefault, o)

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

	wishlistUpsertCacheMut.RLock()
	cache, cached := wishlistUpsertCache[key]
	wishlistUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			wishlistAllColumns,
			wishlistColumnsWithDefault,
			wishlistColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			wishlistAllColumns,
			wishlistPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert wishlists, could not build update column list")
		}

		ret := strmangle.SetComplement(wishlistAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(wishlistPrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert wishlists, could not build conflict column list")
			}

			conflict = make([]string, len(wishlistPrimaryKeyColumns))
			copy(conflict, wishlistPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"wishlists\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(wishlistType, wishlistMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(wishlistType, wishlistMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert wishlists")
	}

	if !cached {
		wishlistUpsertCacheMut.Lock()
		wishlistUpsertCache[key] = cache
		wishlistUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Wishlist record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Wishlist) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no Wishlist provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), wishlistPrimaryKeyMapping)
	sql := "DELETE FROM \"wishlists\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from wishlists")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for wishlists")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q wishlistQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no wishlistQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from wishlists")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for wishlists")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o WishlistSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), wishlistPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"wishlists\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, wishlistPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from wishlist slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for wishlists")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Wishlist) Reload(exec boil.Executor) error {
	ret, err := FindWishlist(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *WishlistSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := WishlistSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), wishlistPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"wishlists\".* FROM \"wishlists\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, wishlistPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in WishlistSlice")
	}

	*o = slice

	return nil
}

// WishlistExists checks if the Wishlist row exists.
func WishlistExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"wishlists\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if wishlists exists")
	}

	return exists, nil
}

// Exists checks if the Wishlist row exists.
func (o *Wishlist) Exists(exec boil.Executor) (bool, error) {
	return WishlistExists(exec, o.ID)
}
