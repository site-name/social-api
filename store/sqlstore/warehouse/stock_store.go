package warehouse

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

type SqlStockStore struct {
	store.Store
}

func NewSqlStockStore(s store.Store) store.StockStore {
	ss := &SqlStockStore{
		Store: s,
	}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.Stock{}, store.StockTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("WarehouseID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("WarehouseID", "ProductVariantID")
	}
	return ss
}

func (ss *SqlStockStore) ModelFields() []string {
	return []string{
		"Stocks.Id",
		"Stocks.CreateAt",
		"Stocks.WarehouseID",
		"Stocks.ProductVariantID",
		"Stocks.Quantity",
	}
}

func (ss *SqlStockStore) CreateIndexesIfNotExists() {
	ss.CreateForeignKeyIfNotExists(store.StockTableName, "WarehouseID", store.WarehouseTableName, "Id", true)
	ss.CreateForeignKeyIfNotExists(store.StockTableName, "ProductVariantID", store.ProductVariantTableName, "Id", true)
}

// BulkUpsert performs upserts or inserts given stocks, then returns them
func (ss *SqlStockStore) BulkUpsert(transaction *gorp.Transaction, stocks []*warehouse.Stock) ([]*warehouse.Stock, error) {

	var (
		isSaving   bool
		insertFunc func(list ...interface{}) error          = ss.GetMaster().Insert
		updateFunc func(list ...interface{}) (int64, error) = ss.GetMaster().Update
	)
	if transaction != nil {
		insertFunc = transaction.Insert
		updateFunc = transaction.Update
	}

	for _, stock := range stocks {
		isSaving = false // reset

		if stock.Id == "" {
			isSaving = true
			stock.PreSave()
		} else {
			stock.PreUpdate()
		}

		if err := stock.IsValid(); err != nil {
			return nil, err
		}

		var (
			err       error
			numUpdate int64
			oldStock  *warehouse.Stock
		)
		if isSaving {
			err = insertFunc(stock)
		} else {
			// try finding a stock with id:
			oldStock, err = ss.Get(stock.Id)
			if err != nil {
				return nil, err
			}

			stock.CreateAt = oldStock.CreateAt

			numUpdate, err = updateFunc(stock)
		}

		if err != nil {
			if ss.IsUniqueConstraintError(err, []string{"WarehouseID", "ProductVariantID", "stocks_warehouseid_productvariantid_key"}) {
				return nil, store.NewErrInvalidInput(store.StockTableName, "WarehouseID/ProductVariantID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert a stock with id=%s", stock.Id)
		}
		if numUpdate > 1 {
			return nil, errors.Errorf("multiple stocks with id=%d were updated: %d instead of 1", stock.Id, numUpdate)
		}
	}

	return stocks, nil
}

func (ss *SqlStockStore) Get(stockID string) (*warehouse.Stock, error) {
	var res warehouse.Stock
	if err := ss.GetReplica().SelectOne(&res, "SELECT * FROM "+store.StockTableName+" WHERE Id = :ID", map[string]interface{}{"ID": stockID}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.StockTableName, stockID)
		}
		return nil, errors.Wrapf(err, "failed to find stock with id=%s", stockID)
	} else {
		return &res, nil
	}
}

// commonLookup is not exported
func (ss *SqlStockStore) commonLookup(transaction *gorp.Transaction, query squirrel.SelectBuilder) ([]*warehouse.Stock, error) {
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "commonLookup.ToSql")
	}

	var (
		returningStocks []*warehouse.Stock
		stock           warehouse.Stock
		wareHouse       warehouse.WareHouse
		productVariant  product_and_discount.ProductVariant
		queryFunc       func(query string, args ...interface{}) (*sql.Rows, error) = ss.GetReplica().Query
	)
	if transaction != nil {
		queryFunc = transaction.Query
	}

	rows, err := queryFunc(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks with given options")
	}
	for rows.Next() {
		err = rows.Scan(
			&stock.Id,
			&stock.CreateAt,
			&stock.WarehouseID,
			&stock.ProductVariantID,
			&stock.Quantity,

			&wareHouse.Id,
			&wareHouse.Name,
			&wareHouse.Slug,
			&wareHouse.AddressID,
			&wareHouse.Email,
			&wareHouse.Metadata,
			&wareHouse.PrivateMetadata,

			&productVariant.Id,
			&productVariant.Name,
			&productVariant.ProductID,
			&productVariant.Sku,
			&productVariant.Weight,
			&productVariant.WeightUnit,
			&productVariant.TrackInventory,
			&productVariant.SortOrder,
			&productVariant.Metadata,
			&productVariant.PrivateMetadata,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of stock, warehouse, product variant")
		}

		stock.Warehouse = &wareHouse
		stock.ProductVariant = &productVariant
		returningStocks = append(returningStocks, &stock)
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return returningStocks, nil
}

// FilterForChannel
func (ss *SqlStockStore) FilterForChannel(channelSlug string) ([]*warehouse.Stock, error) {

	channelQuery := ss.GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ChannelTableName).
		Where(squirrel.Expr("Channels.Slug = ?", channelSlug)).
		Where(squirrel.Expr("Channels.Id = ShippingZoneChannels.ChannelID")).
		Limit(1).
		Suffix(")")

	shippingZoneChannelQuery := ss.GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ShippingZoneChannelTableName).
		Where(channelQuery).
		Where("ShippingZoneChannels.ShippingZoneID = WarehouseShippingZones.ShippingZoneID").
		Limit(1).
		Suffix(")")

	warehouseShippingZoneQuery := ss.GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.WarehouseShippingZoneTableName).
		Where(shippingZoneChannelQuery).
		Where("WarehouseShippingZones.WarehouseID = Stocks.WarehouseID").
		Limit(1).
		Suffix(")")

	selectFields := append(
		ss.ModelFields(),
		ss.ProductVariant().ModelFields()...,
	)
	rows, err := ss.GetQueryBuilder().
		Select(selectFields...).
		From(store.StockTableName).
		InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		Where(warehouseShippingZoneQuery).
		OrderBy(store.TableOrderingMap[store.StockTableName]).
		RunWith(ss.GetReplica()).Query()

	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks with given channel slug")
	}

	var (
		returningStocks []*warehouse.Stock
		stock           warehouse.Stock
		productVariant  product_and_discount.ProductVariant
	)

	for rows.Next() {
		err = rows.Scan(
			&stock.Id,
			&stock.CreateAt,
			&stock.WarehouseID,
			&stock.ProductVariantID,
			&stock.Quantity,

			&productVariant.Id,
			&productVariant.Name,
			&productVariant.ProductID,
			&productVariant.Sku,
			&productVariant.Weight,
			&productVariant.WeightUnit,
			&productVariant.TrackInventory,
			&productVariant.SortOrder,
			&productVariant.Metadata,
			&productVariant.PrivateMetadata,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row contains stock, product variant")
		}

		stock.ProductVariant = &productVariant
		returningStocks = append(returningStocks, &stock)
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return returningStocks, nil
}

// FilterByOption finds and returns a slice of stocks that satisfy given option
func (ss *SqlStockStore) FilterByOption(transaction *gorp.Transaction, options *warehouse.StockFilterOption) ([]*warehouse.Stock, error) {
	// decide which query to use
	var (
		query                        squirrel.SelectBuilder
		useForCountryAndChannelQuery bool
	)
	if options.ForCountryAndChannel != nil {
		query = ss.FilterForCountryAndChannel(options.ForCountryAndChannel)
		useForCountryAndChannelQuery = true // indicate that we are using query build by method `FilterForCountryAndChannel()`
	} else {
		query = ss.GetQueryBuilder().
			Select(ss.ModelFields()...). // this selecting fields differ the query from `if` caluse
			From(store.StockTableName).
			OrderBy(store.TableOrderingMap[store.StockTableName])
	}

	// parse options:
	if options.Id != nil {
		query = query.Where(options.Id.ToSquirrel("Stocks.Id"))
	}
	if options.WarehouseID != nil {
		query = query.Where(options.WarehouseID.ToSquirrel("Stocks.WarehouseID"))
	}
	if options.ProductVariantID != nil {
		query = query.Where(options.ProductVariantID.ToSquirrel("Stocks.ProductVariantID"))
	}
	if options.LockForUpdate {
		query = query.Suffix("FOR UPDATE")
	}
	if options.ForUpdateOf != "" && options.LockForUpdate {
		query = query.Suffix("OF " + options.ForUpdateOf)
	}

	// check which query is used
	if !useForCountryAndChannelQuery {
		queryString, args, err := query.ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "FilterbyOption_ToSql")
		}

		var res []*warehouse.Stock
		_, err = transaction.Select(&res, queryString, args...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find stocks by given options")
		}
		return res, nil // these stocks does not contains related data
	}

	stocks, err := ss.commonLookup(transaction, query) // these stocks contains related data suc as `Warehouse`, `ProductVariant`
	return stocks, err
}

func (ss *SqlStockStore) FilterForCountryAndChannel(options *warehouse.StockFilterForCountryAndChannel) squirrel.SelectBuilder {
	options.CountryCode = strings.ToUpper(options.CountryCode)

	warehouseIDQuery := ss.warehouseIdSelectQuery(options.CountryCode, options.ChannelSlug)

	selectFields := ss.ModelFields()
	selectFields = append(selectFields, ss.Warehouse().ModelFields()...)
	selectFields = append(selectFields, ss.ProductVariant().ModelFields()...)

	return ss.GetQueryBuilder().
		Select(selectFields...).
		From(store.StockTableName).
		InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		InnerJoin(store.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)").
		Where(squirrel.Expr("Stocks.WarehouseID IN ?", warehouseIDQuery)).
		OrderBy(store.TableOrderingMap[store.StockTableName])
}

func (ss *SqlStockStore) FilterVariantStocksForCountry(options *warehouse.StockFilterForCountryAndChannel) ([]*warehouse.Stock, error) {
	transaction, err := ss.GetReplica().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}

	options.CountryCode = strings.ToUpper(options.CountryCode)

	warehouseIDQuery := ss.warehouseIdSelectQuery(options.CountryCode, options.ChannelSlug)

	selectFields := ss.ModelFields()
	selectFields = append(selectFields, ss.Warehouse().ModelFields()...)
	selectFields = append(selectFields, ss.ProductVariant().ModelFields()...)

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(store.StockTableName).
		InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		InnerJoin(store.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)").
		Where(squirrel.Expr("Stocks.WarehouseID IN ?", warehouseIDQuery)).
		Where(squirrel.Expr("Stocks.ProductVariantID = ?", options.ProductVariantID)).
		OrderBy(store.TableOrderingMap[store.StockTableName])

	return ss.commonLookup(transaction, query)
}

func (ss *SqlStockStore) FilterProductStocksForCountryAndChannel(options *warehouse.StockFilterForCountryAndChannel) ([]*warehouse.Stock, error) {
	transaction, err := ss.GetReplica().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}

	options.CountryCode = strings.ToUpper(options.CountryCode)
	warehouseIDQuery := ss.warehouseIdSelectQuery(options.CountryCode, options.ChannelSlug)

	selectingFields := ss.ModelFields()
	selectingFields = append(selectingFields, ss.Warehouse().ModelFields()...)
	selectingFields = append(selectingFields, ss.ProductVariant().ModelFields()...)

	query := ss.GetQueryBuilder().
		Select(selectingFields...).
		From(store.StockTableName).
		InnerJoin(store.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)").
		InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		InnerJoin(store.ProductTableName + " ON (Products.Id = ProductVariants.ProductID)").
		Where(squirrel.Expr("Stocks.WarehouseID IN ?", warehouseIDQuery)).
		Where(squirrel.Expr("Products.Id = ?", options.ProductID)).
		OrderBy(store.TableOrderingMap[store.StockTableName])

	return ss.commonLookup(transaction, query)
}

func (ss *SqlStockStore) warehouseIdSelectQuery(countryCode string, channelSlug string) squirrel.SelectBuilder {
	query := ss.GetQueryBuilder().
		Select("Warehouses.Id").
		From(store.WarehouseTableName)

	if countryCode != "" {
		query = query.
			InnerJoin(store.WarehouseShippingZoneTableName+" ON (Warehouses.Id = WarehouseShippingZones.WarehouseID)").
			InnerJoin(store.ShippingZoneTableName+" ON (ShippingZones.Id = WarehouseShippingZones.ShippingZoneID)").
			Where("ShippingZones.Countries :: text LIKE ?", "%"+countryCode+"%")
	}
	if channelSlug != "" {
		query = query.
			InnerJoin(store.ShippingZoneChannelTableName+" ON (ShippingZoneChannels.ShippingZoneID = ShippingZones.Id)").
			InnerJoin(store.ChannelTableName+" ON (Channels.Id = ShippingZoneChannels.ChannelID)").
			Where("Channels.Slug = ?", channelSlug)
	}

	return query
}

// ChangeQuantity reduce or increase the quantity of given stock
func (ss *SqlStockStore) ChangeQuantity(stockID string, quantity int) error {
	_, err := ss.GetMaster().Exec("UPDATE Stocks SET Quantity = Quantity + $1 WHERE Id = $2", quantity, stockID)
	if err != nil {
		return errors.Wrapf(err, "failed to change stock quantity for stock with id=%s", stockID)
	}

	return nil
}
