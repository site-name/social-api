package warehouse

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlStockStore struct {
	store.Store
}

func NewSqlStockStore(s store.Store) store.StockStore {
	return &SqlStockStore{Store: s}
}

func (ss *SqlStockStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"CreateAt",
		"WarehouseID",
		"ProductVariantID",
		"Quantity",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}
func (ss *SqlStockStore) ScanFields(stock *model.Stock) []interface{} {
	return []interface{}{
		&stock.Id,
		&stock.CreateAt,
		&stock.WarehouseID,
		&stock.ProductVariantID,
		&stock.Quantity,
	}
}

// BulkUpsert performs upserts or inserts given stocks, then returns them
func (ss *SqlStockStore) BulkUpsert(transaction store_iface.SqlxTxExecutor, stocks []*model.Stock) ([]*model.Stock, error) {
	var executor store_iface.SqlxExecutor = ss.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	var (
		saveQuery   = "INSERT INTO " + store.StockTableName + "(" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"
		updateQuery = "UPDATE " + store.StockTableName + " SET " + ss.
				ModelFields("").
				Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"
	)

	for _, stock := range stocks {
		isSaving := false // reset

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
		)
		if isSaving {
			_, err = executor.NamedExec(saveQuery, stock)

		} else {
			var result sql.Result
			result, err = executor.NamedExec(updateQuery, stock)
			if err == nil && result != nil {
				numUpdate, _ = result.RowsAffected()
			}
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

func (ss *SqlStockStore) Get(stockID string) (*model.Stock, error) {
	var res model.Stock
	if err := ss.GetReplicaX().Get(&res, "SELECT * FROM "+store.StockTableName+" WHERE Id = ?", stockID); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.StockTableName, stockID)
		}
		return nil, errors.Wrapf(err, "failed to find stock with id=%s", stockID)
	}
	return &res, nil
}

// FilterForChannel finds and returns stocks that satisfy given options
func (ss *SqlStockStore) FilterForChannel(options *model.StockFilterForChannelOption) (squirrel.Sqlizer, []*model.Stock, error) {
	channelQuery := ss.GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ChannelTableName).
		Where("Channels.Id = ?", options.ChannelID).
		Where("Channels.Id = ShippingZoneChannels.ChannelID").
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

	selectFields := ss.ModelFields(store.StockTableName + ".")
	// check if we need select related data:
	if options.SelectRelatedProductVariant {
		selectFields = append(selectFields, ss.ProductVariant().ModelFields(store.ProductVariantTableName+".")...)
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
		query = query.Where(options.Id)
	}
	if options.WarehouseID != nil {
		query = query.Where(options.WarehouseID)
	}
	if options.ProductVariantID != nil {
		query = query.Where(options.ProductVariantID)
	}

	if options.ReturnQueryOnly {
		return query, nil, nil
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, nil, errors.Wrap(err, "FilterForChannel_ToSql")
	}

	rows, err := ss.GetReplicaX().QueryX(queryString, args...)

	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find stocks with given channel slug")
	}

	var (
		returningStocks []*model.Stock
		stock           model.Stock
		productVariant  model.ProductVariant
		scanFields      = ss.ScanFields(&stock)
	)
	if options.SelectRelatedProductVariant {
		scanFields = append(scanFields, ss.ProductVariant().ScanFields(&productVariant)...)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to scan a row contains stock")
		}

		if options.SelectRelatedProductVariant {
			stock.SetProductVariant(&productVariant)
		}
		returningStocks = append(returningStocks, stock.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, nil, errors.Wrap(err, "failed to close rows of stocks")
	}

	return nil, returningStocks, nil
}

// FilterByOption finds and returns a slice of stocks that satisfy given option
func (ss *SqlStockStore) FilterByOption(transaction store_iface.SqlxTxExecutor, options *model.StockFilterOption) ([]*model.Stock, error) {
	selectFields := ss.ModelFields(store.StockTableName + ".")
	if options.SelectRelatedWarehouse {
		selectFields = append(selectFields, ss.Warehouse().ModelFields(store.WarehouseTableName+".")...)
	}
	if options.SelectRelatedProductVariant {
		selectFields = append(selectFields, ss.ProductVariant().ModelFields(store.ProductVariantTableName+".")...)
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
		query = query.Where(options.Id)
	}
	if options.WarehouseID != nil {
		query = query.Where(options.WarehouseID)
	}
	if options.ProductVariantID != nil {
		query = query.Where(options.ProductVariantID)
	}
	if options.LockForUpdate {
		query = query.Suffix("FOR UPDATE")
	}
	if options.ForUpdateOf != "" && options.LockForUpdate {
		query = query.Suffix("OF " + options.ForUpdateOf)
	}
	if options.Warehouse_ShippingZone_countries != nil ||
		options.Warehouse_ShippingZone_ChannelID != nil {
		query = query.
			InnerJoin(store.WarehouseTableName + " ON Warehouses.Id = Stocks.WarehouseID").
			InnerJoin(store.WarehouseShippingZoneTableName + " ON WarehouseShippingZones.WarehouseID = Warehouses.Id").
			InnerJoin(store.ShippingZoneTableName + " ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID")
	}
	if options.Warehouse_ShippingZone_countries != nil {
		query = query.Where(options.Warehouse_ShippingZone_countries)
	}
	if options.Warehouse_ShippingZone_ChannelID != nil {
		query = query.
			InnerJoin(store.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id").
			Where(options.Warehouse_ShippingZone_ChannelID)
	}

	var groupBy string

	if options.AnnotateAvailabeQuantity {
		query = query.
			Column(squirrel.Alias(squirrel.Expr("Stocks.Quantity - COALESCE( SUM ( Allocations.QuantityAllocated ), 0 )"), "AvailableQuantity")).
			LeftJoin(store.AllocationTableName + " ON (Stocks.Id = Allocations.StockID)")
		groupBy = "Stocks.Id"
	}

	if len(groupBy) > 0 {
		query = query.GroupBy(groupBy)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var (
		returningStocks   []*model.Stock
		stock             model.Stock
		variant           model.ProductVariant
		wareHouse         model.WareHouse
		queryer           store_iface.SqlxExecutor = ss.GetReplicaX()
		availableQuantity int
		scanFields        = ss.ScanFields(&stock)
	)
	if transaction != nil {
		queryer = transaction
	}

	if options.SelectRelatedWarehouse {
		scanFields = append(scanFields, ss.Warehouse().ScanFields(&wareHouse)...)
	}
	if options.SelectRelatedProductVariant {
		scanFields = append(scanFields, ss.ProductVariant().ScanFields(&variant)...)
	}
	if options.AnnotateAvailabeQuantity {
		scanFields = append(scanFields, &availableQuantity)
	}

	rows, err := queryer.QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks by given options")
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find stocks with related warehouses and product variants")
		}

		if options.SelectRelatedWarehouse {
			stock.SetWarehouse(&wareHouse)
		}
		if options.SelectRelatedProductVariant {
			stock.SetProductVariant(&variant)
		}
		if options.AnnotateAvailabeQuantity {
			stock.AvailableQuantity = availableQuantity
		}
		returningStocks = append(returningStocks, stock.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return returningStocks, nil
}

// FilterForCountryAndChannel finds and returns stocks with given options
func (ss *SqlStockStore) FilterForCountryAndChannel(transaction store_iface.SqlxTxExecutor, options *model.StockFilterForCountryAndChannel) ([]*model.Stock, error) {
	options.CountryCode = strings.ToUpper(options.CountryCode)

	warehouseIDQuery := ss.
		warehouseIdSelectQuery(options.CountryCode, options.ChannelSlug).
		PlaceholderFormat(squirrel.Question)

	// remember the order when scan
	selectFields := ss.ModelFields(store.StockTableName + ".")
	selectFields = append(selectFields, ss.Warehouse().ModelFields(store.WarehouseTableName+".")...)
	selectFields = append(selectFields, ss.ProductVariant().ModelFields(store.ProductVariantTableName+".")...)

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(store.StockTableName).
		InnerJoin(store.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)").
		InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		Where(squirrel.Expr("Stocks.WarehouseID IN ?", warehouseIDQuery)).
		OrderBy(store.TableOrderingMap[store.StockTableName])

	// parse option for FilterVariantStocksForCountry
	// parse additional options
	if options.AnnotateAvailabeQuantity {
		query = query.
			Column("Stocks.Quantity - COALESCE(SUM(Allocations.QuantityAllocated), 0) AS AvailableQuantity").
			LeftJoin(store.AllocationTableName + " ON (Stocks.Id = Allocations.StockID)")
	}

	if options.ProductVariantID != "" {
		query = query.Where("Stocks.ProductVariantID = ?", options.ProductVariantID)
	}

	// parse option for FilterProductStocksForCountryAndChannel
	if options.ProductID != "" {
		query = query.
			InnerJoin(store.ProductTableName+" ON (Products.Id = ProductVariants.ProductID)").
			Where("Products.Id = ?", options.ProductID)
	}
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.WarehouseIDFilter != nil {
		query = query.Where(options.WarehouseIDFilter)
	}
	if options.ProductVariantIDFilter != nil {
		query = query.Where(options.ProductVariantIDFilter)
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
		returningStocks   []*model.Stock
		stock             model.Stock
		wareHouse         model.WareHouse
		productVariant    model.ProductVariant
		queryer           store_iface.SqlxExecutor = ss.GetReplicaX()
		availableQuantity int
		scanFields        = ss.ScanFields(&stock)
	)
	// add some more fields to scan
	scanFields = append(scanFields, ss.Warehouse().ScanFields(&wareHouse)...)
	scanFields = append(scanFields, ss.ProductVariant().ScanFields(&productVariant)...)
	// decide which query to use
	if transaction != nil {
		queryer = transaction
	}

	if options.AnnotateAvailabeQuantity {
		scanFields = append(scanFields, &availableQuantity)
	}

	rows, err := queryer.QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks with given options")
	}
	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of stock, warehouse, product variant")
		}

		stock.SetWarehouse(&wareHouse)
		stock.SetProductVariant(&productVariant)
		returningStocks = append(returningStocks, stock.DeepCopy())

		if options.AnnotateAvailabeQuantity {
			stock.AvailableQuantity = availableQuantity
		}
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return returningStocks, nil
}

func (ss *SqlStockStore) FilterVariantStocksForCountry(transaction store_iface.SqlxTxExecutor, options *model.StockFilterForCountryAndChannel) ([]*model.Stock, error) {
	return ss.FilterForCountryAndChannel(transaction, options)
}

func (ss *SqlStockStore) FilterProductStocksForCountryAndChannel(transaction store_iface.SqlxTxExecutor, options *model.StockFilterForCountryAndChannel) ([]*model.Stock, error) {
	return ss.FilterForCountryAndChannel(transaction, options)
}

func (ss *SqlStockStore) warehouseIdSelectQuery(countryCode string, channelSlug string) squirrel.SelectBuilder {
	query := ss.GetQueryBuilder().
		Select("Warehouses.Id").
		From(store.WarehouseTableName)

	if countryCode != "" {
		query = query.
			InnerJoin(store.WarehouseShippingZoneTableName+" ON Warehouses.Id = WarehouseShippingZones.WarehouseID").
			InnerJoin(store.ShippingZoneTableName+" ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID").
			Where("ShippingZones.Countries::text LIKE ?", "%"+countryCode+"%")
	}
	if channelSlug != "" {
		query = query.
			InnerJoin(store.ShippingZoneChannelTableName+" ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id").
			InnerJoin(store.ChannelTableName+" ON Channels.Id = ShippingZoneChannels.ChannelID").
			Where("Channels.Slug = ?", channelSlug)
	}

	return query
}

// ChangeQuantity reduce or increase the quantity of given stock
func (ss *SqlStockStore) ChangeQuantity(stockID string, quantity int) error {
	_, err := ss.GetMasterX().Exec("UPDATE Stocks SET Quantity = Quantity + ? WHERE Id = ?", quantity, stockID)
	if err != nil {
		return errors.Wrapf(err, "failed to change stock quantity for stock with id=%s", stockID)
	}

	return nil
}
