package warehouse

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlStockStore struct {
	store.Store
}

func NewSqlStockStore(s store.Store) store.StockStore {
	return &SqlStockStore{Store: s}
}

func (ss *SqlStockStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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
func (ss *SqlStockStore) BulkUpsert(transaction store_iface.SqlxExecutor, stocks []*model.Stock) ([]*model.Stock, error) {
	var executor = ss.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	var (
		saveQuery   = "INSERT INTO " + model.StockTableName + "(" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"
		updateQuery = "UPDATE " + model.StockTableName + " SET " + ss.
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
				return nil, store.NewErrInvalidInput(model.StockTableName, "WarehouseID/ProductVariantID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert a stock with id=%s", stock.Id)
		}
		if numUpdate > 1 {
			return nil, errors.Errorf("multiple stocks with id=%s were updated: %d instead of 1", stock.Id, numUpdate)
		}
	}

	return stocks, nil
}

func (ss *SqlStockStore) Get(stockID string) (*model.Stock, error) {
	var res model.Stock
	if err := ss.GetReplicaX().Get(&res, "SELECT * FROM "+model.StockTableName+" WHERE Id = ?", stockID); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.StockTableName, stockID)
		}
		return nil, errors.Wrapf(err, "failed to find stock with id=%s", stockID)
	}
	return &res, nil
}

// FilterForChannel finds and returns stocks that satisfy given options
func (ss *SqlStockStore) FilterForChannel(options *model.StockFilterForChannelOption) (squirrel.Sqlizer, []*model.Stock, error) {
	channelQuery := ss.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.ChannelTableName).
		Where("Channels.Id = ?", options.ChannelID).
		Where("Channels.Id = ShippingZoneChannels.ChannelID").
		Limit(1).
		Suffix(")")

	shippingZoneChannelQuery := ss.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.ShippingZoneChannelTableName).
		Where(channelQuery).
		Where("ShippingZoneChannels.ShippingZoneID = WarehouseShippingZones.ShippingZoneID").
		Limit(1).
		Suffix(")")

	warehouseShippingZoneQuery := ss.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.WarehouseShippingZoneTableName).
		Where(shippingZoneChannelQuery).
		Where("WarehouseShippingZones.WarehouseID = Stocks.WarehouseID").
		Limit(1).
		Suffix(")")

	selectFields := ss.ModelFields(model.StockTableName + ".")
	// check if we need select related data:
	if options.SelectRelatedProductVariant {
		selectFields = append(selectFields, ss.ProductVariant().ModelFields(model.ProductVariantTableName+".")...)
	}

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(model.StockTableName).
		Where(warehouseShippingZoneQuery)

	// parse options
	if options.SelectRelatedProductVariant {
		query = query.InnerJoin(model.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)")
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
	defer rows.Close()

	var returningStocks []*model.Stock

	for rows.Next() {
		var (
			stock          model.Stock
			productVariant model.ProductVariant
			scanFields     = ss.ScanFields(&stock)
		)
		if options.SelectRelatedProductVariant {
			scanFields = append(scanFields, ss.ProductVariant().ScanFields(&productVariant)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to scan a row contains stock")
		}

		if options.SelectRelatedProductVariant {
			stock.SetProductVariant(&productVariant)
		}

		returningStocks = append(returningStocks, &stock)
	}

	return nil, returningStocks, nil
}

func (s *SqlStockStore) CountByOptions(options *model.StockFilterOption) (int32, error) {
	query := s.GetQueryBuilder().Select("COUNT(DISTINCT Stocks.Id)").From(model.StockTableName)

	var stockSearchOpts squirrel.Sqlizer = nil
	if options.Search != "" {
		expr := "%" + options.Search + "%"

		stockSearchOpts = squirrel.Or{
			squirrel.ILike{model.ProductTableName + ".Name": expr},
			squirrel.ILike{model.ProductVariantTableName + ".Name": expr},
			squirrel.ILike{model.WarehouseTableName + ".Name": expr},
			squirrel.ILike{model.AddressTableName + ".CompanyName": expr},
		}
	}

	// parse options:
	for _, opt := range []squirrel.Sqlizer{
		options.Conditions,
		options.Warehouse_ShippingZone_countries,
		options.Warehouse_ShippingZone_ChannelID,
		stockSearchOpts, //
	} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	if options.Search != "" ||
		options.Warehouse_ShippingZone_countries != nil ||
		options.Warehouse_ShippingZone_ChannelID != nil {

		query = query.InnerJoin(model.WarehouseTableName + " ON Warehouses.Id = Stocks.WarehouseID")

		if options.Warehouse_ShippingZone_countries != nil ||
			options.Warehouse_ShippingZone_ChannelID != nil {
			query = query.
				InnerJoin(model.WarehouseShippingZoneTableName + " ON WarehouseShippingZones.WarehouseID = Warehouses.Id").
				InnerJoin(model.ShippingZoneTableName + " ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID")

			if options.Warehouse_ShippingZone_ChannelID != nil {
				query = query.InnerJoin(model.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id")
			}
		}

		if options.Search != "" {
			query = query.InnerJoin(model.AddressTableName + " ON Addresses.Id = Warehouses.AddressID")
		}
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "CountByOptions_ToSql")
	}

	var res int32
	err = s.GetReplicaX().Select(&res, queryStr, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count stocks by given options")
	}

	return res, nil
}

// FilterByOption finds and returns a slice of stocks that satisfy given option
func (ss *SqlStockStore) FilterByOption(options *model.StockFilterOption) ([]*model.Stock, error) {
	selectFields := ss.ModelFields(model.StockTableName + ".")
	if options.SelectRelatedProductVariant {
		selectFields = append(selectFields, ss.ProductVariant().ModelFields(model.ProductVariantTableName+".")...)
	}
	if options.SelectRelatedWarehouse {
		selectFields = append(selectFields, ss.Warehouse().ModelFields(model.WarehouseTableName+".")...)
	}

	query := ss.GetQueryBuilder().
		Select(selectFields...). // this selecting fields differ the query from `if` caluse
		From(model.StockTableName)

	var stockSearchOpts squirrel.Sqlizer = nil
	if options.Search != "" {
		expr := "%" + options.Search + "%"

		stockSearchOpts = squirrel.Or{
			squirrel.ILike{model.ProductTableName + ".Name": expr},
			squirrel.ILike{model.ProductVariantTableName + ".Name": expr},
			squirrel.ILike{model.WarehouseTableName + ".Name": expr},
			squirrel.ILike{model.AddressTableName + ".CompanyName": expr},
		}
	}

	// parse options:
	for _, opt := range []squirrel.Sqlizer{
		options.Conditions,
		options.Warehouse_ShippingZone_countries,
		options.Warehouse_ShippingZone_ChannelID,
		stockSearchOpts, //
	} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	if options.LockForUpdate {
		forUpdate := "FOR UPDATE"
		if options.ForUpdateOf != "" {
			forUpdate += " OF " + options.ForUpdateOf
		}

		query = query.Suffix(forUpdate)
	}

	// NOTE: The order of join must similar to order of select above
	if options.SelectRelatedProductVariant || options.Search != "" {
		query = query.InnerJoin(model.ProductVariantTableName + " ON ProductVariants.Id = Stocks.ProductVariantID")

		if options.Search != "" {
			query = query.InnerJoin(model.ProductTableName + " ON Products.Id = ProductVariants.ProductID")
		}
	}

	if options.SelectRelatedWarehouse ||
		options.Search != "" ||
		options.Warehouse_ShippingZone_countries != nil ||
		options.Warehouse_ShippingZone_ChannelID != nil {

		query = query.InnerJoin(model.WarehouseTableName + " ON Warehouses.Id = Stocks.WarehouseID")

		if options.Warehouse_ShippingZone_countries != nil ||
			options.Warehouse_ShippingZone_ChannelID != nil {
			query = query.
				InnerJoin(model.WarehouseShippingZoneTableName + " ON WarehouseShippingZones.WarehouseID = Warehouses.Id").
				InnerJoin(model.ShippingZoneTableName + " ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID")

			if options.Warehouse_ShippingZone_ChannelID != nil {
				query = query.InnerJoin(model.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id")
			}
		}

		if options.Search != "" {
			query = query.InnerJoin(model.AddressTableName + " ON Addresses.Id = Warehouses.AddressID")
		}
	}

	var groupBy string

	if options.AnnotateAvailabeQuantity {
		query = query.
			Column(squirrel.Alias(squirrel.Expr("Stocks.Quantity - COALESCE( SUM ( Allocations.QuantityAllocated ), 0 )"), "AvailableQuantity")).
			LeftJoin(model.AllocationTableName + " ON (Stocks.Id = Allocations.StockID)")
		groupBy = "Stocks.Id"
	}

	if groupBy != "" {
		query = query.GroupBy(groupBy)
	}

	query = options.PaginationValues.AddPaginationToSelectBuilderIfNeeded(query)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	rows, err := ss.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks by given options")
	}
	defer rows.Close()

	returningStocks := make(model.Stocks, 0)

	for rows.Next() {
		var (
			stock             model.Stock
			variant           model.ProductVariant
			wareHouse         model.WareHouse
			availableQuantity int
			scanFields        = ss.ScanFields(&stock)
		)

		// NOTE: The order of scan fields must similar to order of select above
		if options.SelectRelatedProductVariant {
			scanFields = append(scanFields, ss.ProductVariant().ScanFields(&variant)...)
		}
		if options.SelectRelatedWarehouse {
			scanFields = append(scanFields, ss.Warehouse().ScanFields(&wareHouse)...)
		}
		if options.AnnotateAvailabeQuantity {
			scanFields = append(scanFields, &availableQuantity)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find stocks with related warehouses and product variants")
		}

		if options.SelectRelatedProductVariant {
			stock.SetProductVariant(&variant)
		}
		if options.SelectRelatedWarehouse {
			stock.SetWarehouse(&wareHouse)
		}
		if options.AnnotateAvailabeQuantity {
			stock.AvailableQuantity = availableQuantity
		}
		returningStocks = append(returningStocks, &stock)
	}

	return returningStocks, nil
}

// FilterForCountryAndChannel finds and returns stocks with given options
func (ss *SqlStockStore) FilterForCountryAndChannel(options *model.StockFilterForCountryAndChannel) ([]*model.Stock, error) {
	warehouseIDQuery := ss.
		warehouseIdSelectQuery(options.CountryCode, options.ChannelSlug).
		PlaceholderFormat(squirrel.Question)

	// remember the order when scan
	selectFields := ss.ModelFields(model.StockTableName + ".")
	selectFields = append(selectFields, ss.Warehouse().ModelFields(model.WarehouseTableName+".")...)
	selectFields = append(selectFields, ss.ProductVariant().ModelFields(model.ProductVariantTableName+".")...)

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(model.StockTableName).
		InnerJoin(model.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)").
		InnerJoin(model.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		Where(squirrel.Expr("Stocks.WarehouseID IN ?", warehouseIDQuery))

	// parse option for FilterVariantStocksForCountry
	// parse additional options
	if options.AnnotateAvailabeQuantity {
		query = query.
			Column("Stocks.Quantity - COALESCE(SUM(Allocations.QuantityAllocated), 0) AS AvailableQuantity").
			LeftJoin(model.AllocationTableName + " ON (Stocks.Id = Allocations.StockID)")
	}

	if options.ProductVariantID != "" {
		query = query.Where("Stocks.ProductVariantID = ?", options.ProductVariantID)
	}

	// parse option for FilterProductStocksForCountryAndChannel
	if options.ProductID != "" {
		query = query.
			InnerJoin(model.ProductTableName+" ON (roducts.Id = ProductVariants.ProductID").
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
		suffix := "FPR UPDATE"
		if options.ForUpdateOf != "" {
			suffix += " OF " + options.ForUpdateOf
		}

		query = query.Suffix(suffix)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterForCountryAndChannel_ToSql")
	}

	var returningStocks model.Stocks

	rows, err := ss.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks with given options")
	}
	defer rows.Close()

	for rows.Next() {
		var (
			stock             model.Stock
			wareHouse         model.WareHouse
			productVariant    model.ProductVariant
			availableQuantity int
			scanFields        = ss.ScanFields(&stock)
		)
		scanFields = append(scanFields, ss.Warehouse().ScanFields(&wareHouse)...)
		scanFields = append(scanFields, ss.ProductVariant().ScanFields(&productVariant)...)

		if options.AnnotateAvailabeQuantity {
			scanFields = append(scanFields, &availableQuantity)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of stock, warehouse, product variant")
		}

		stock.SetWarehouse(&wareHouse)
		stock.SetProductVariant(&productVariant)

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

func (ss *SqlStockStore) FilterVariantStocksForCountry(options *model.StockFilterForCountryAndChannel) ([]*model.Stock, error) {
	return ss.FilterForCountryAndChannel(options)
}

func (ss *SqlStockStore) FilterProductStocksForCountryAndChannel(options *model.StockFilterForCountryAndChannel) ([]*model.Stock, error) {
	// TODO: finish me
	return ss.FilterForCountryAndChannel(options)
}

func (ss *SqlStockStore) warehouseIdSelectQuery(countryCode model.CountryCode, channelSlug string) squirrel.SelectBuilder {
	query := ss.GetQueryBuilder().
		Select("Warehouses.Id").
		From(model.WarehouseTableName)

	if countryCode != "" {
		query = query.
			InnerJoin(model.WarehouseShippingZoneTableName+" ON Warehouses.Id = WarehouseShippingZones.WarehouseID").
			InnerJoin(model.ShippingZoneTableName+" ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID").
			Where("ShippingZones.Countries::text LIKE ?", "%"+countryCode+"%")
	}
	if channelSlug != "" {
		query = query.
			InnerJoin(model.ShippingZoneChannelTableName+" ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id").
			InnerJoin(model.ChannelTableName+" ON Channels.Id = ShippingZoneChannels.ChannelID").
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
