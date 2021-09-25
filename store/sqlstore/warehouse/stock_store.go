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
func (ss *SqlStockStore) ScanFields(stock warehouse.Stock) []interface{} {
	return []interface{}{
		&stock.Id,
		&stock.CreateAt,
		&stock.WarehouseID,
		&stock.ProductVariantID,
		&stock.Quantity,
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

// FilterForChannel finds and returns stocks that satisfy given options
func (ss *SqlStockStore) FilterForChannel(options *warehouse.StockFilterForChannelOption) ([]*warehouse.Stock, error) {
	channelQuery := ss.GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ChannelTableName).
		Where(squirrel.Expr("Channels.Slug = ?", options.ChannelSlug)).
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

	selectFields := ss.ModelFields()
	// check if we need select related data:
	if options.SelectRelatedProductVariant {
		selectFields = append(selectFields, ss.ProductVariant().ModelFields()...)
	}

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(store.StockTableName).
		Where(warehouseShippingZoneQuery).
		OrderBy(store.TableOrderingMap[store.StockTableName])

	// parse options

	if options.SelectRelatedProductVariant {
		query = query.InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)")
	}
	if options.Id != nil {
		query = query.Where(options.Id.ToSquirrel("Stocks.Id"))
	}
	if options.WarehouseID != nil {
		query = query.Where(options.WarehouseID.ToSquirrel("Stocks.WarehouseID"))
	}
	if options.ProductVariantID != nil {
		query = query.Where(options.ProductVariantID.ToSquirrel("Stocks.ProductVariantID"))
	}

	rows, err := query.RunWith(ss.GetReplica()).Query()

	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks with given channel slug")
	}

	var (
		returningStocks []*warehouse.Stock
		stock           warehouse.Stock
		productVariant  product_and_discount.ProductVariant
		scanFields      = ss.ScanFields(stock)
	)
	if options.SelectRelatedProductVariant {
		scanFields = append(scanFields, ss.ProductVariant().ScanFields(productVariant)...)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row contains stock, product variant")
		}

		if options.SelectRelatedProductVariant {
			stock.ProductVariant = &productVariant
		}
		returningStocks = append(returningStocks, &stock)
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return returningStocks, nil
}

// FilterByOption finds and returns a slice of stocks that satisfy given option
func (ss *SqlStockStore) FilterByOption(transaction *gorp.Transaction, options *warehouse.StockFilterOption) ([]*warehouse.Stock, error) {
	selectFields := ss.ModelFields()
	if options.SelectRelatedWarehouse {
		selectFields = append(selectFields, ss.Warehouse().ModelFields()...)
	}
	if options.SelectRelatedProductVariant {
		selectFields = append(selectFields, ss.ProductVariant().ModelFields()...)
	}

	query := ss.GetQueryBuilder().
		Select(selectFields...). // this selecting fields differ the query from `if` caluse
		From(store.StockTableName).
		OrderBy(store.TableOrderingMap[store.StockTableName])

	// parse options:
	if options.SelectRelatedProductVariant {
		query = query.InnerJoin(store.ProductVariantTableName + " ON (ProductVariants.Id = Stocks.ProductVariantID)")
	}
	if options.SelectRelatedWarehouse {
		query = query.InnerJoin(store.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)")
	}
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
	if options.AnnotateAvailabeQuantity {
		query = query.
			Column(squirrel.Alias(squirrel.Expr("Stocks.Quantity - COALESCE(SUM(Allocations.QuantityAllocated), 0)"), "AvailableQuantity")).
			LeftJoin(store.AllocationTableName + " ON (Stocks.Id = Allocations.StockID)")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var (
		returningStocks   []*warehouse.Stock
		stock             warehouse.Stock
		variant           product_and_discount.ProductVariant
		wareHouse         warehouse.WareHouse
		selectFunc        func(query string, args ...interface{}) (*sql.Rows, error) = ss.GetReplica().Query
		availableQuantity int
		scanFields        = ss.ScanFields(stock)
	)
	if transaction != nil {
		selectFunc = transaction.Query
	}

	if options.SelectRelatedWarehouse {
		scanFields = append(scanFields, ss.Warehouse().ScanFields(wareHouse)...)
	}
	if options.SelectRelatedProductVariant {
		scanFields = append(scanFields, ss.ProductVariant().ScanFields(variant)...)
	}
	if options.AnnotateAvailabeQuantity {
		scanFields = append(scanFields, &availableQuantity)
	}

	rows, err := selectFunc(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks by given options")
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find stocks with related warehouses and product variants")
		}

		stock.Warehouse = &wareHouse
		stock.ProductVariant = &variant
		returningStocks = append(returningStocks, &stock)

		if options.AnnotateAvailabeQuantity {
			stock.AvailableQuantity = availableQuantity
		}
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return returningStocks, nil
}

// FilterForCountryAndChannel finds and returns stocks with given options
func (ss *SqlStockStore) FilterForCountryAndChannel(transaction *gorp.Transaction, options *warehouse.StockFilterForCountryAndChannel) ([]*warehouse.Stock, error) {
	options.CountryCode = strings.ToUpper(options.CountryCode)

	warehouseIDQuery := ss.warehouseIdSelectQuery(options.CountryCode, options.ChannelSlug)

	// remember the order when scan
	selectFields := ss.ModelFields()
	selectFields = append(selectFields, ss.Warehouse().ModelFields()...)
	selectFields = append(selectFields, ss.ProductVariant().ModelFields()...)

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(store.StockTableName).
		InnerJoin(store.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)").
		InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		Where(squirrel.Expr("Stocks.WarehouseID IN ?", warehouseIDQuery)).
		OrderBy(store.TableOrderingMap[store.StockTableName])

	// parse option for FilterVariantStocksForCountry
	if options.ProductVariantID != "" {
		query = query.Where(squirrel.Expr("Stocks.ProductVariantID = ?", options.ProductVariantID))
	}
	// parse option for FilterProductStocksForCountryAndChannel
	if options.ProductID != "" {
		query = query.
			InnerJoin(store.ProductTableName + " ON (Products.Id = ProductVariants.ProductID)").
			Where(squirrel.Expr("Products.Id = ?", options.ProductID))
	}
	// parse additional options
	if options.AnnotateAvailabeQuantity {
		query = query.
			Column(squirrel.Alias(squirrel.Expr("Stocks.Quantity - COALESCE(SUM(Allocations.QuantityAllocated), 0)"), "AvailableQuantity")).
			LeftJoin(store.AllocationTableName + " ON (Stocks.Id = Allocations.StockID)")
	}
	if options.Id != nil {
		query = query.Where(options.Id.ToSquirrel("Stocks.Id"))
	}
	if options.WarehouseIDFilter != nil {
		query = query.Where(options.WarehouseIDFilter.ToSquirrel("Stocks.WarehouseID"))
	}
	if options.ProductVariantIDFilter != nil {
		query = query.Where(options.ProductVariantIDFilter.ToSquirrel("Stocks.ProductVariantID"))
	}
	if options.LockForUpdate {
		query = query.Suffix("FOR UPDATE")
	}
	if options.ForUpdateOf != "" && options.LockForUpdate {
		query = query.Suffix("OF " + options.ForUpdateOf)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterForCountryAndChannel_ToSql")
	}

	var (
		returningStocks   []*warehouse.Stock
		stock             warehouse.Stock
		wareHouse         warehouse.WareHouse
		productVariant    product_and_discount.ProductVariant
		queryer           squirrel.Queryer = ss.GetReplica()
		availableQuantity int
		scanFields        = ss.ScanFields(stock)
	)
	// add some more fields to scan
	scanFields = append(scanFields, ss.Warehouse().ScanFields(wareHouse)...)
	scanFields = append(scanFields, ss.ProductVariant().ScanFields(productVariant)...)
	// decide which query to use
	if transaction != nil {
		queryer = transaction
	}

	if options.AnnotateAvailabeQuantity {
		scanFields = append(scanFields, &availableQuantity)
	}

	rows, err := queryer.Query(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks with given options")
	}
	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of stock, warehouse, product variant")
		}

		stock.Warehouse = &wareHouse
		stock.ProductVariant = &productVariant
		returningStocks = append(returningStocks, &stock)

		if options.AnnotateAvailabeQuantity {
			stock.AvailableQuantity = availableQuantity
		}
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return returningStocks, nil
}

func (ss *SqlStockStore) FilterVariantStocksForCountry(transaction *gorp.Transaction, options *warehouse.StockFilterForCountryAndChannel) ([]*warehouse.Stock, error) {
	return ss.FilterForCountryAndChannel(transaction, options)
}

func (ss *SqlStockStore) FilterProductStocksForCountryAndChannel(transaction *gorp.Transaction, options *warehouse.StockFilterForCountryAndChannel) ([]*warehouse.Stock, error) {
	return ss.FilterForCountryAndChannel(transaction, options)
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
