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

// SaleProduct is an object representing the database table.
type SaleProduct struct {
	ID        string `boil:"id" json:"id" toml:"id" yaml:"id"`
	SaleID    string `boil:"sale_id" json:"sale_id" toml:"sale_id" yaml:"sale_id"`
	ProductID string `boil:"product_id" json:"product_id" toml:"product_id" yaml:"product_id"`
	CreatedAt int64  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *saleProductR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L saleProductL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SaleProductColumns = struct {
	ID        string
	SaleID    string
	ProductID string
	CreatedAt string
}{
	ID:        "id",
	SaleID:    "sale_id",
	ProductID: "product_id",
	CreatedAt: "created_at",
}

var SaleProductTableColumns = struct {
	ID        string
	SaleID    string
	ProductID string
	CreatedAt string
}{
	ID:        "sale_products.id",
	SaleID:    "sale_products.sale_id",
	ProductID: "sale_products.product_id",
	CreatedAt: "sale_products.created_at",
}

// Generated where

var SaleProductWhere = struct {
	ID        whereHelperstring
	SaleID    whereHelperstring
	ProductID whereHelperstring
	CreatedAt whereHelperint64
}{
	ID:        whereHelperstring{field: "\"sale_products\".\"id\""},
	SaleID:    whereHelperstring{field: "\"sale_products\".\"sale_id\""},
	ProductID: whereHelperstring{field: "\"sale_products\".\"product_id\""},
	CreatedAt: whereHelperint64{field: "\"sale_products\".\"created_at\""},
}

// SaleProductRels is where relationship names are stored.
var SaleProductRels = struct {
	Product string
	Sale    string
}{
	Product: "Product",
	Sale:    "Sale",
}

// saleProductR is where relationships are stored.
type saleProductR struct {
	Product *Product `boil:"Product" json:"Product" toml:"Product" yaml:"Product"`
	Sale    *Sale    `boil:"Sale" json:"Sale" toml:"Sale" yaml:"Sale"`
}

// NewStruct creates a new relationship struct
func (*saleProductR) NewStruct() *saleProductR {
	return &saleProductR{}
}

func (r *saleProductR) GetProduct() *Product {
	if r == nil {
		return nil
	}
	return r.Product
}

func (r *saleProductR) GetSale() *Sale {
	if r == nil {
		return nil
	}
	return r.Sale
}

// saleProductL is where Load methods for each relationship are stored.
type saleProductL struct{}

var (
	saleProductAllColumns            = []string{"id", "sale_id", "product_id", "created_at"}
	saleProductColumnsWithoutDefault = []string{"sale_id", "product_id", "created_at"}
	saleProductColumnsWithDefault    = []string{"id"}
	saleProductPrimaryKeyColumns     = []string{"id"}
	saleProductGeneratedColumns      = []string{}
)

type (
	// SaleProductSlice is an alias for a slice of pointers to SaleProduct.
	// This should almost always be used instead of []SaleProduct.
	SaleProductSlice []*SaleProduct

	saleProductQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	saleProductType                 = reflect.TypeOf(&SaleProduct{})
	saleProductMapping              = queries.MakeStructMapping(saleProductType)
	saleProductPrimaryKeyMapping, _ = queries.BindMapping(saleProductType, saleProductMapping, saleProductPrimaryKeyColumns)
	saleProductInsertCacheMut       sync.RWMutex
	saleProductInsertCache          = make(map[string]insertCache)
	saleProductUpdateCacheMut       sync.RWMutex
	saleProductUpdateCache          = make(map[string]updateCache)
	saleProductUpsertCacheMut       sync.RWMutex
	saleProductUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single saleProduct record from the query.
func (q saleProductQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SaleProduct, error) {
	o := &SaleProduct{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for sale_products")
	}

	return o, nil
}

// All returns all SaleProduct records from the query.
func (q saleProductQuery) All(ctx context.Context, exec boil.ContextExecutor) (SaleProductSlice, error) {
	var o []*SaleProduct

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to SaleProduct slice")
	}

	return o, nil
}

// Count returns the count of all SaleProduct records in the query.
func (q saleProductQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count sale_products rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q saleProductQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if sale_products exists")
	}

	return count > 0, nil
}

// Product pointed to by the foreign key.
func (o *SaleProduct) Product(mods ...qm.QueryMod) productQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ProductID),
	}

	queryMods = append(queryMods, mods...)

	return Products(queryMods...)
}

// Sale pointed to by the foreign key.
func (o *SaleProduct) Sale(mods ...qm.QueryMod) saleQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SaleID),
	}

	queryMods = append(queryMods, mods...)

	return Sales(queryMods...)
}

// LoadProduct allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (saleProductL) LoadProduct(ctx context.Context, e boil.ContextExecutor, singular bool, maybeSaleProduct interface{}, mods queries.Applicator) error {
	var slice []*SaleProduct
	var object *SaleProduct

	if singular {
		var ok bool
		object, ok = maybeSaleProduct.(*SaleProduct)
		if !ok {
			object = new(SaleProduct)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeSaleProduct)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeSaleProduct))
			}
		}
	} else {
		s, ok := maybeSaleProduct.(*[]*SaleProduct)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeSaleProduct)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeSaleProduct))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &saleProductR{}
		}
		args = append(args, object.ProductID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &saleProductR{}
			}

			for _, a := range args {
				if a == obj.ProductID {
					continue Outer
				}
			}

			args = append(args, obj.ProductID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`products`),
		qm.WhereIn(`products.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
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
		foreign.R.SaleProducts = append(foreign.R.SaleProducts, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ProductID == foreign.ID {
				local.R.Product = foreign
				if foreign.R == nil {
					foreign.R = &productR{}
				}
				foreign.R.SaleProducts = append(foreign.R.SaleProducts, local)
				break
			}
		}
	}

	return nil
}

// LoadSale allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (saleProductL) LoadSale(ctx context.Context, e boil.ContextExecutor, singular bool, maybeSaleProduct interface{}, mods queries.Applicator) error {
	var slice []*SaleProduct
	var object *SaleProduct

	if singular {
		var ok bool
		object, ok = maybeSaleProduct.(*SaleProduct)
		if !ok {
			object = new(SaleProduct)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeSaleProduct)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeSaleProduct))
			}
		}
	} else {
		s, ok := maybeSaleProduct.(*[]*SaleProduct)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeSaleProduct)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeSaleProduct))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &saleProductR{}
		}
		args = append(args, object.SaleID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &saleProductR{}
			}

			for _, a := range args {
				if a == obj.SaleID {
					continue Outer
				}
			}

			args = append(args, obj.SaleID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`sales`),
		qm.WhereIn(`sales.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Sale")
	}

	var resultSlice []*Sale
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Sale")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for sales")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for sales")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Sale = foreign
		if foreign.R == nil {
			foreign.R = &saleR{}
		}
		foreign.R.SaleProducts = append(foreign.R.SaleProducts, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.SaleID == foreign.ID {
				local.R.Sale = foreign
				if foreign.R == nil {
					foreign.R = &saleR{}
				}
				foreign.R.SaleProducts = append(foreign.R.SaleProducts, local)
				break
			}
		}
	}

	return nil
}

// SetProduct of the saleProduct to the related item.
// Sets o.R.Product to related.
// Adds o to related.R.SaleProducts.
func (o *SaleProduct) SetProduct(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Product) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"sale_products\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"product_id"}),
		strmangle.WhereClause("\"", "\"", 2, saleProductPrimaryKeyColumns),
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

	o.ProductID = related.ID
	if o.R == nil {
		o.R = &saleProductR{
			Product: related,
		}
	} else {
		o.R.Product = related
	}

	if related.R == nil {
		related.R = &productR{
			SaleProducts: SaleProductSlice{o},
		}
	} else {
		related.R.SaleProducts = append(related.R.SaleProducts, o)
	}

	return nil
}

// SetSale of the saleProduct to the related item.
// Sets o.R.Sale to related.
// Adds o to related.R.SaleProducts.
func (o *SaleProduct) SetSale(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Sale) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"sale_products\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"sale_id"}),
		strmangle.WhereClause("\"", "\"", 2, saleProductPrimaryKeyColumns),
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

	o.SaleID = related.ID
	if o.R == nil {
		o.R = &saleProductR{
			Sale: related,
		}
	} else {
		o.R.Sale = related
	}

	if related.R == nil {
		related.R = &saleR{
			SaleProducts: SaleProductSlice{o},
		}
	} else {
		related.R.SaleProducts = append(related.R.SaleProducts, o)
	}

	return nil
}

// SaleProducts retrieves all the records using an executor.
func SaleProducts(mods ...qm.QueryMod) saleProductQuery {
	mods = append(mods, qm.From("\"sale_products\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"sale_products\".*"})
	}

	return saleProductQuery{q}
}

// FindSaleProduct retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSaleProduct(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*SaleProduct, error) {
	saleProductObj := &SaleProduct{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"sale_products\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, saleProductObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from sale_products")
	}

	return saleProductObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SaleProduct) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no sale_products provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(saleProductColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	saleProductInsertCacheMut.RLock()
	cache, cached := saleProductInsertCache[key]
	saleProductInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			saleProductAllColumns,
			saleProductColumnsWithDefault,
			saleProductColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(saleProductType, saleProductMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(saleProductType, saleProductMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"sale_products\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"sale_products\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into sale_products")
	}

	if !cached {
		saleProductInsertCacheMut.Lock()
		saleProductInsertCache[key] = cache
		saleProductInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the SaleProduct.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SaleProduct) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	saleProductUpdateCacheMut.RLock()
	cache, cached := saleProductUpdateCache[key]
	saleProductUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			saleProductAllColumns,
			saleProductPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update sale_products, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"sale_products\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, saleProductPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(saleProductType, saleProductMapping, append(wl, saleProductPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update sale_products row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for sale_products")
	}

	if !cached {
		saleProductUpdateCacheMut.Lock()
		saleProductUpdateCache[key] = cache
		saleProductUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q saleProductQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for sale_products")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for sale_products")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SaleProductSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), saleProductPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"sale_products\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, saleProductPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in saleProduct slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all saleProduct")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SaleProduct) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no sale_products provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(saleProductColumnsWithDefault, o)

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

	saleProductUpsertCacheMut.RLock()
	cache, cached := saleProductUpsertCache[key]
	saleProductUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			saleProductAllColumns,
			saleProductColumnsWithDefault,
			saleProductColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			saleProductAllColumns,
			saleProductPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert sale_products, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(saleProductPrimaryKeyColumns))
			copy(conflict, saleProductPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"sale_products\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(saleProductType, saleProductMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(saleProductType, saleProductMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert sale_products")
	}

	if !cached {
		saleProductUpsertCacheMut.Lock()
		saleProductUpsertCache[key] = cache
		saleProductUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single SaleProduct record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SaleProduct) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no SaleProduct provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), saleProductPrimaryKeyMapping)
	sql := "DELETE FROM \"sale_products\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from sale_products")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for sale_products")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q saleProductQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no saleProductQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from sale_products")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for sale_products")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SaleProductSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), saleProductPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"sale_products\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, saleProductPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from saleProduct slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for sale_products")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SaleProduct) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSaleProduct(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SaleProductSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SaleProductSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), saleProductPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"sale_products\".* FROM \"sale_products\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, saleProductPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in SaleProductSlice")
	}

	*o = slice

	return nil
}

// SaleProductExists checks if the SaleProduct row exists.
func SaleProductExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"sale_products\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if sale_products exists")
	}

	return exists, nil
}

// Exists checks if the SaleProduct row exists.
func (o *SaleProduct) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return SaleProductExists(ctx, exec, o.ID)
}