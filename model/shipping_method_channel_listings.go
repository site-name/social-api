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
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// ShippingMethodChannelListing is an object representing the database table.
type ShippingMethodChannelListing struct {
	ID                      string            `boil:"id" json:"id" toml:"id" yaml:"id"`
	ShippingMethodID        string            `boil:"shipping_method_id" json:"shipping_method_id" toml:"shipping_method_id" yaml:"shipping_method_id"`
	ChannelID               string            `boil:"channel_id" json:"channel_id" toml:"channel_id" yaml:"channel_id"`
	MinimumOrderPriceAmount types.Decimal     `boil:"minimum_order_price_amount" json:"minimum_order_price_amount" toml:"minimum_order_price_amount" yaml:"minimum_order_price_amount"`
	Currency                Currency          `boil:"currency" json:"currency" toml:"currency" yaml:"currency"`
	MaximumOrderPriceAmount types.NullDecimal `boil:"maximum_order_price_amount" json:"maximum_order_price_amount,omitempty" toml:"maximum_order_price_amount" yaml:"maximum_order_price_amount,omitempty"`
	PriceAmount             types.Decimal     `boil:"price_amount" json:"price_amount" toml:"price_amount" yaml:"price_amount"`
	CreatedAt               int64             `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *shippingMethodChannelListingR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L shippingMethodChannelListingL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ShippingMethodChannelListingColumns = struct {
	ID                      string
	ShippingMethodID        string
	ChannelID               string
	MinimumOrderPriceAmount string
	Currency                string
	MaximumOrderPriceAmount string
	PriceAmount             string
	CreatedAt               string
}{
	ID:                      "id",
	ShippingMethodID:        "shipping_method_id",
	ChannelID:               "channel_id",
	MinimumOrderPriceAmount: "minimum_order_price_amount",
	Currency:                "currency",
	MaximumOrderPriceAmount: "maximum_order_price_amount",
	PriceAmount:             "price_amount",
	CreatedAt:               "created_at",
}

var ShippingMethodChannelListingTableColumns = struct {
	ID                      string
	ShippingMethodID        string
	ChannelID               string
	MinimumOrderPriceAmount string
	Currency                string
	MaximumOrderPriceAmount string
	PriceAmount             string
	CreatedAt               string
}{
	ID:                      "shipping_method_channel_listings.id",
	ShippingMethodID:        "shipping_method_channel_listings.shipping_method_id",
	ChannelID:               "shipping_method_channel_listings.channel_id",
	MinimumOrderPriceAmount: "shipping_method_channel_listings.minimum_order_price_amount",
	Currency:                "shipping_method_channel_listings.currency",
	MaximumOrderPriceAmount: "shipping_method_channel_listings.maximum_order_price_amount",
	PriceAmount:             "shipping_method_channel_listings.price_amount",
	CreatedAt:               "shipping_method_channel_listings.created_at",
}

// Generated where

var ShippingMethodChannelListingWhere = struct {
	ID                      whereHelperstring
	ShippingMethodID        whereHelperstring
	ChannelID               whereHelperstring
	MinimumOrderPriceAmount whereHelpertypes_Decimal
	Currency                whereHelperCurrency
	MaximumOrderPriceAmount whereHelpertypes_NullDecimal
	PriceAmount             whereHelpertypes_Decimal
	CreatedAt               whereHelperint64
}{
	ID:                      whereHelperstring{field: "\"shipping_method_channel_listings\".\"id\""},
	ShippingMethodID:        whereHelperstring{field: "\"shipping_method_channel_listings\".\"shipping_method_id\""},
	ChannelID:               whereHelperstring{field: "\"shipping_method_channel_listings\".\"channel_id\""},
	MinimumOrderPriceAmount: whereHelpertypes_Decimal{field: "\"shipping_method_channel_listings\".\"minimum_order_price_amount\""},
	Currency:                whereHelperCurrency{field: "\"shipping_method_channel_listings\".\"currency\""},
	MaximumOrderPriceAmount: whereHelpertypes_NullDecimal{field: "\"shipping_method_channel_listings\".\"maximum_order_price_amount\""},
	PriceAmount:             whereHelpertypes_Decimal{field: "\"shipping_method_channel_listings\".\"price_amount\""},
	CreatedAt:               whereHelperint64{field: "\"shipping_method_channel_listings\".\"created_at\""},
}

// ShippingMethodChannelListingRels is where relationship names are stored.
var ShippingMethodChannelListingRels = struct {
	Channel        string
	ShippingMethod string
}{
	Channel:        "Channel",
	ShippingMethod: "ShippingMethod",
}

// shippingMethodChannelListingR is where relationships are stored.
type shippingMethodChannelListingR struct {
	Channel        *Channel        `boil:"Channel" json:"Channel" toml:"Channel" yaml:"Channel"`
	ShippingMethod *ShippingMethod `boil:"ShippingMethod" json:"ShippingMethod" toml:"ShippingMethod" yaml:"ShippingMethod"`
}

// NewStruct creates a new relationship struct
func (*shippingMethodChannelListingR) NewStruct() *shippingMethodChannelListingR {
	return &shippingMethodChannelListingR{}
}

func (r *shippingMethodChannelListingR) GetChannel() *Channel {
	if r == nil {
		return nil
	}
	return r.Channel
}

func (r *shippingMethodChannelListingR) GetShippingMethod() *ShippingMethod {
	if r == nil {
		return nil
	}
	return r.ShippingMethod
}

// shippingMethodChannelListingL is where Load methods for each relationship are stored.
type shippingMethodChannelListingL struct{}

var (
	shippingMethodChannelListingAllColumns            = []string{"id", "shipping_method_id", "channel_id", "minimum_order_price_amount", "currency", "maximum_order_price_amount", "price_amount", "created_at"}
	shippingMethodChannelListingColumnsWithoutDefault = []string{"shipping_method_id", "channel_id", "currency", "created_at"}
	shippingMethodChannelListingColumnsWithDefault    = []string{"id", "minimum_order_price_amount", "maximum_order_price_amount", "price_amount"}
	shippingMethodChannelListingPrimaryKeyColumns     = []string{"id"}
	shippingMethodChannelListingGeneratedColumns      = []string{}
)

type (
	// ShippingMethodChannelListingSlice is an alias for a slice of pointers to ShippingMethodChannelListing.
	// This should almost always be used instead of []ShippingMethodChannelListing.
	ShippingMethodChannelListingSlice []*ShippingMethodChannelListing

	shippingMethodChannelListingQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	shippingMethodChannelListingType                 = reflect.TypeOf(&ShippingMethodChannelListing{})
	shippingMethodChannelListingMapping              = queries.MakeStructMapping(shippingMethodChannelListingType)
	shippingMethodChannelListingPrimaryKeyMapping, _ = queries.BindMapping(shippingMethodChannelListingType, shippingMethodChannelListingMapping, shippingMethodChannelListingPrimaryKeyColumns)
	shippingMethodChannelListingInsertCacheMut       sync.RWMutex
	shippingMethodChannelListingInsertCache          = make(map[string]insertCache)
	shippingMethodChannelListingUpdateCacheMut       sync.RWMutex
	shippingMethodChannelListingUpdateCache          = make(map[string]updateCache)
	shippingMethodChannelListingUpsertCacheMut       sync.RWMutex
	shippingMethodChannelListingUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single shippingMethodChannelListing record from the query.
func (q shippingMethodChannelListingQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ShippingMethodChannelListing, error) {
	o := &ShippingMethodChannelListing{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for shipping_method_channel_listings")
	}

	return o, nil
}

// All returns all ShippingMethodChannelListing records from the query.
func (q shippingMethodChannelListingQuery) All(ctx context.Context, exec boil.ContextExecutor) (ShippingMethodChannelListingSlice, error) {
	var o []*ShippingMethodChannelListing

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to ShippingMethodChannelListing slice")
	}

	return o, nil
}

// Count returns the count of all ShippingMethodChannelListing records in the query.
func (q shippingMethodChannelListingQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count shipping_method_channel_listings rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q shippingMethodChannelListingQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if shipping_method_channel_listings exists")
	}

	return count > 0, nil
}

// Channel pointed to by the foreign key.
func (o *ShippingMethodChannelListing) Channel(mods ...qm.QueryMod) channelQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ChannelID),
	}

	queryMods = append(queryMods, mods...)

	return Channels(queryMods...)
}

// ShippingMethod pointed to by the foreign key.
func (o *ShippingMethodChannelListing) ShippingMethod(mods ...qm.QueryMod) shippingMethodQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ShippingMethodID),
	}

	queryMods = append(queryMods, mods...)

	return ShippingMethods(queryMods...)
}

// LoadChannel allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (shippingMethodChannelListingL) LoadChannel(ctx context.Context, e boil.ContextExecutor, singular bool, maybeShippingMethodChannelListing interface{}, mods queries.Applicator) error {
	var slice []*ShippingMethodChannelListing
	var object *ShippingMethodChannelListing

	if singular {
		var ok bool
		object, ok = maybeShippingMethodChannelListing.(*ShippingMethodChannelListing)
		if !ok {
			object = new(ShippingMethodChannelListing)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeShippingMethodChannelListing)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeShippingMethodChannelListing))
			}
		}
	} else {
		s, ok := maybeShippingMethodChannelListing.(*[]*ShippingMethodChannelListing)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeShippingMethodChannelListing)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeShippingMethodChannelListing))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &shippingMethodChannelListingR{}
		}
		args = append(args, object.ChannelID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &shippingMethodChannelListingR{}
			}

			for _, a := range args {
				if a == obj.ChannelID {
					continue Outer
				}
			}

			args = append(args, obj.ChannelID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`channels`),
		qm.WhereIn(`channels.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Channel")
	}

	var resultSlice []*Channel
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Channel")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for channels")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for channels")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Channel = foreign
		if foreign.R == nil {
			foreign.R = &channelR{}
		}
		foreign.R.ShippingMethodChannelListings = append(foreign.R.ShippingMethodChannelListings, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ChannelID == foreign.ID {
				local.R.Channel = foreign
				if foreign.R == nil {
					foreign.R = &channelR{}
				}
				foreign.R.ShippingMethodChannelListings = append(foreign.R.ShippingMethodChannelListings, local)
				break
			}
		}
	}

	return nil
}

// LoadShippingMethod allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (shippingMethodChannelListingL) LoadShippingMethod(ctx context.Context, e boil.ContextExecutor, singular bool, maybeShippingMethodChannelListing interface{}, mods queries.Applicator) error {
	var slice []*ShippingMethodChannelListing
	var object *ShippingMethodChannelListing

	if singular {
		var ok bool
		object, ok = maybeShippingMethodChannelListing.(*ShippingMethodChannelListing)
		if !ok {
			object = new(ShippingMethodChannelListing)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeShippingMethodChannelListing)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeShippingMethodChannelListing))
			}
		}
	} else {
		s, ok := maybeShippingMethodChannelListing.(*[]*ShippingMethodChannelListing)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeShippingMethodChannelListing)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeShippingMethodChannelListing))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &shippingMethodChannelListingR{}
		}
		args = append(args, object.ShippingMethodID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &shippingMethodChannelListingR{}
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
		foreign.R.ShippingMethodChannelListings = append(foreign.R.ShippingMethodChannelListings, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ShippingMethodID == foreign.ID {
				local.R.ShippingMethod = foreign
				if foreign.R == nil {
					foreign.R = &shippingMethodR{}
				}
				foreign.R.ShippingMethodChannelListings = append(foreign.R.ShippingMethodChannelListings, local)
				break
			}
		}
	}

	return nil
}

// SetChannel of the shippingMethodChannelListing to the related item.
// Sets o.R.Channel to related.
// Adds o to related.R.ShippingMethodChannelListings.
func (o *ShippingMethodChannelListing) SetChannel(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Channel) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"shipping_method_channel_listings\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"channel_id"}),
		strmangle.WhereClause("\"", "\"", 2, shippingMethodChannelListingPrimaryKeyColumns),
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

	o.ChannelID = related.ID
	if o.R == nil {
		o.R = &shippingMethodChannelListingR{
			Channel: related,
		}
	} else {
		o.R.Channel = related
	}

	if related.R == nil {
		related.R = &channelR{
			ShippingMethodChannelListings: ShippingMethodChannelListingSlice{o},
		}
	} else {
		related.R.ShippingMethodChannelListings = append(related.R.ShippingMethodChannelListings, o)
	}

	return nil
}

// SetShippingMethod of the shippingMethodChannelListing to the related item.
// Sets o.R.ShippingMethod to related.
// Adds o to related.R.ShippingMethodChannelListings.
func (o *ShippingMethodChannelListing) SetShippingMethod(ctx context.Context, exec boil.ContextExecutor, insert bool, related *ShippingMethod) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"shipping_method_channel_listings\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"shipping_method_id"}),
		strmangle.WhereClause("\"", "\"", 2, shippingMethodChannelListingPrimaryKeyColumns),
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
		o.R = &shippingMethodChannelListingR{
			ShippingMethod: related,
		}
	} else {
		o.R.ShippingMethod = related
	}

	if related.R == nil {
		related.R = &shippingMethodR{
			ShippingMethodChannelListings: ShippingMethodChannelListingSlice{o},
		}
	} else {
		related.R.ShippingMethodChannelListings = append(related.R.ShippingMethodChannelListings, o)
	}

	return nil
}

// ShippingMethodChannelListings retrieves all the records using an executor.
func ShippingMethodChannelListings(mods ...qm.QueryMod) shippingMethodChannelListingQuery {
	mods = append(mods, qm.From("\"shipping_method_channel_listings\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"shipping_method_channel_listings\".*"})
	}

	return shippingMethodChannelListingQuery{q}
}

// FindShippingMethodChannelListing retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindShippingMethodChannelListing(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*ShippingMethodChannelListing, error) {
	shippingMethodChannelListingObj := &ShippingMethodChannelListing{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"shipping_method_channel_listings\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, shippingMethodChannelListingObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from shipping_method_channel_listings")
	}

	return shippingMethodChannelListingObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ShippingMethodChannelListing) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no shipping_method_channel_listings provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(shippingMethodChannelListingColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	shippingMethodChannelListingInsertCacheMut.RLock()
	cache, cached := shippingMethodChannelListingInsertCache[key]
	shippingMethodChannelListingInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			shippingMethodChannelListingAllColumns,
			shippingMethodChannelListingColumnsWithDefault,
			shippingMethodChannelListingColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(shippingMethodChannelListingType, shippingMethodChannelListingMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(shippingMethodChannelListingType, shippingMethodChannelListingMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"shipping_method_channel_listings\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"shipping_method_channel_listings\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into shipping_method_channel_listings")
	}

	if !cached {
		shippingMethodChannelListingInsertCacheMut.Lock()
		shippingMethodChannelListingInsertCache[key] = cache
		shippingMethodChannelListingInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the ShippingMethodChannelListing.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ShippingMethodChannelListing) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	shippingMethodChannelListingUpdateCacheMut.RLock()
	cache, cached := shippingMethodChannelListingUpdateCache[key]
	shippingMethodChannelListingUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			shippingMethodChannelListingAllColumns,
			shippingMethodChannelListingPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update shipping_method_channel_listings, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"shipping_method_channel_listings\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, shippingMethodChannelListingPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(shippingMethodChannelListingType, shippingMethodChannelListingMapping, append(wl, shippingMethodChannelListingPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update shipping_method_channel_listings row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for shipping_method_channel_listings")
	}

	if !cached {
		shippingMethodChannelListingUpdateCacheMut.Lock()
		shippingMethodChannelListingUpdateCache[key] = cache
		shippingMethodChannelListingUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q shippingMethodChannelListingQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for shipping_method_channel_listings")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for shipping_method_channel_listings")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ShippingMethodChannelListingSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), shippingMethodChannelListingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"shipping_method_channel_listings\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, shippingMethodChannelListingPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in shippingMethodChannelListing slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all shippingMethodChannelListing")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ShippingMethodChannelListing) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no shipping_method_channel_listings provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(shippingMethodChannelListingColumnsWithDefault, o)

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

	shippingMethodChannelListingUpsertCacheMut.RLock()
	cache, cached := shippingMethodChannelListingUpsertCache[key]
	shippingMethodChannelListingUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			shippingMethodChannelListingAllColumns,
			shippingMethodChannelListingColumnsWithDefault,
			shippingMethodChannelListingColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			shippingMethodChannelListingAllColumns,
			shippingMethodChannelListingPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert shipping_method_channel_listings, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(shippingMethodChannelListingPrimaryKeyColumns))
			copy(conflict, shippingMethodChannelListingPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"shipping_method_channel_listings\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(shippingMethodChannelListingType, shippingMethodChannelListingMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(shippingMethodChannelListingType, shippingMethodChannelListingMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert shipping_method_channel_listings")
	}

	if !cached {
		shippingMethodChannelListingUpsertCacheMut.Lock()
		shippingMethodChannelListingUpsertCache[key] = cache
		shippingMethodChannelListingUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single ShippingMethodChannelListing record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ShippingMethodChannelListing) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no ShippingMethodChannelListing provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), shippingMethodChannelListingPrimaryKeyMapping)
	sql := "DELETE FROM \"shipping_method_channel_listings\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from shipping_method_channel_listings")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for shipping_method_channel_listings")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q shippingMethodChannelListingQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no shippingMethodChannelListingQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from shipping_method_channel_listings")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for shipping_method_channel_listings")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ShippingMethodChannelListingSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), shippingMethodChannelListingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"shipping_method_channel_listings\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, shippingMethodChannelListingPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from shippingMethodChannelListing slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for shipping_method_channel_listings")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ShippingMethodChannelListing) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindShippingMethodChannelListing(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ShippingMethodChannelListingSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ShippingMethodChannelListingSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), shippingMethodChannelListingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"shipping_method_channel_listings\".* FROM \"shipping_method_channel_listings\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, shippingMethodChannelListingPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in ShippingMethodChannelListingSlice")
	}

	*o = slice

	return nil
}

// ShippingMethodChannelListingExists checks if the ShippingMethodChannelListing row exists.
func ShippingMethodChannelListingExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"shipping_method_channel_listings\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if shipping_method_channel_listings exists")
	}

	return exists, nil
}

// Exists checks if the ShippingMethodChannelListing row exists.
func (o *ShippingMethodChannelListing) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return ShippingMethodChannelListingExists(ctx, exec, o.ID)
}