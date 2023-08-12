package warehouse

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlStockStore struct {
	store.Store
}

func NewSqlStockStore(s store.Store) store.StockStore {
	return &SqlStockStore{Store: s}
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
func (ss *SqlStockStore) BulkUpsert(transaction *gorm.DB, stocks []*model.Stock) ([]*model.Stock, error) {
	if transaction == nil {
		transaction = ss.GetMaster()
	}

	for _, stock := range stocks {
		err := transaction.Save(stock).Error
		if err != nil {
			if ss.IsUniqueConstraintError(err, []string{"WarehouseID", "ProductVariantID", "warehouseid_productvariantid_key"}) {
				return nil, store.NewErrInvalidInput(model.StockTableName, "WarehouseID/ProductVariantID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert a stock with id=%s", stock.Id)
		}
	}

	return stocks, nil
}

func (ss *SqlStockStore) Get(stockID string) (*model.Stock, error) {
	var res model.Stock
	if err := ss.GetReplica().First(&res, "Id = ?", stockID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	selectFields := []string{model.StockTableName + ".*"}
	// check if we need select related data:
	if options.SelectRelatedProductVariant {
		selectFields = append(selectFields, model.ProductVariantTableName+".*")
	}

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(model.StockTableName).
		Where(warehouseShippingZoneQuery).Where(options.Conditions)

	// parse options
	if options.SelectRelatedProductVariant {
		query = query.InnerJoin(model.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)")
	}

	if options.ReturnQueryOnly {
		return query, nil, nil
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, nil, errors.Wrap(err, "FilterForChannel_ToSql")
	}

	rows, err := ss.GetReplica().Raw(queryString, args...).Rows()
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

// func (s *SqlStockStore) CountByOptions(options *model.StockFilterOption) (int32, error) {
// 	query := s.GetQueryBuilder().Select("COUNT(DISTINCT Stocks.Id)").From(model.StockTableName)

// 	var stockSearchOpts squirrel.Sqlizer = nil
// 	if options.Search != "" {
// 		expr := "%" + options.Search + "%"

// 		stockSearchOpts = squirrel.Or{
// 			squirrel.ILike{model.ProductTableName + ".Name": expr},
// 			squirrel.ILike{model.ProductVariantTableName + ".Name": expr},
// 			squirrel.ILike{model.WarehouseTableName + ".Name": expr},
// 			squirrel.ILike{model.AddressTableName + ".CompanyName": expr},
// 		}
// 	}

// 	// parse options:
// 	for _, opt := range []squirrel.Sqlizer{
// 		options.Conditions,
// 		options.Warehouse_ShippingZone_countries,
// 		options.Warehouse_ShippingZone_ChannelID,
// 		stockSearchOpts, //
// 	} {
// 		query = query.Where(opt)
// 	}

// 	if options.Search != "" ||
// 		options.Warehouse_ShippingZone_countries != nil ||
// 		options.Warehouse_ShippingZone_ChannelID != nil {

// 		query = query.InnerJoin(model.WarehouseTableName + " ON Warehouses.Id = Stocks.WarehouseID")

// 		if options.Warehouse_ShippingZone_countries != nil ||
// 			options.Warehouse_ShippingZone_ChannelID != nil {
// 			query = query.
// 				InnerJoin(model.WarehouseShippingZoneTableName + " ON WarehouseShippingZones.WarehouseID = Warehouses.Id").
// 				InnerJoin(model.ShippingZoneTableName + " ON ShippingZones.Id = WarehouseShippingZones.ShippingZoneID")

// 			if options.Warehouse_ShippingZone_ChannelID != nil {
// 				query = query.InnerJoin(model.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id")
// 			}
// 		}

// 		if options.Search != "" {
// 			query = query.InnerJoin(model.AddressTableName + " ON Addresses.Id = Warehouses.AddressID")
// 		}
// 	}

// 	queryStr, args, err := query.ToSql()
// 	if err != nil {
// 		return 0, errors.Wrap(err, "CountByOptions_ToSql")
// 	}

// 	var res int32
// 	err = s.GetReplica().Raw(queryStr, args...).Scan(&res).Error
// 	if err != nil {
// 		return 0, errors.Wrap(err, "failed to count stocks by given options")
// 	}

// 	return res, nil
// }

// FilterByOption finds and returns a slice of stocks that satisfy given option
func (ss *SqlStockStore) FilterByOption(options *model.StockFilterOption) (int64, []*model.Stock, error) {
	selectFields := []string{model.StockTableName + ".*"}
	if options.SelectRelatedProductVariant {
		selectFields = append(selectFields, model.ProductVariantTableName+".*")
	}
	if options.SelectRelatedWarehouse {
		selectFields = append(selectFields, model.WarehouseTableName+".*")
	}

	query := ss.GetQueryBuilder().
		Select(selectFields...). // this selecting fields differ the query from `if` caluse
		From(model.StockTableName)

	if options.Distinct {
		query = query.Distinct()
	}

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
		query = query.Where(opt)
	}

	if options.LockForUpdate && options.Transaction != nil {
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
	if options.AnnotateAvailableQuantity {
		query = query.
			Column("Stocks.Quantity - COALESCE( SUM ( Allocations.QuantityAllocated ), 0 ) AS AvailableQuantity").
			LeftJoin(model.AllocationTableName + " ON Stocks.Id = Allocations.StockID")
		groupBy = "Stocks.Id"
	}

	if groupBy != "" {
		query = query.GroupBy(groupBy)
	}

	// check if graphql pagination provided:
	if options.GraphqlPaginationValues.PaginationApplicable() {
		query = query.
			Where(options.GraphqlPaginationValues.Condition).
			OrderBy(options.GraphqlPaginationValues.OrderBy)
	}

	// NOTE: we have to construct the count query here before pagination limit is applied
	var totalCount int64
	if options.CountTotal {
		countQuery, countArgs, err := ss.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "CountTotal_ToSql")
		}
		err = ss.GetReplica().Raw(countQuery, countArgs...).Scan(&totalCount).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total number of stocks by given options")
		}
	}

	// query = options.GraphqlPaginationValues.AddPaginationToSelectBuilderIfNeeded(query)
	if options.GraphqlPaginationValues.Limit > 0 {
		query = query.Limit(options.GraphqlPaginationValues.Limit)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	runner := ss.GetReplica()
	if options.Transaction != nil {
		runner = options.Transaction
	}

	rows, err := runner.Raw(queryString, args...).Rows()
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to find stocks by given options")
	}
	defer rows.Close()

	var returningStocks model.Stocks

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
		if options.AnnotateAvailableQuantity {
			scanFields = append(scanFields, &availableQuantity)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to find stocks with related warehouses and product variants")
		}

		if options.SelectRelatedProductVariant {
			stock.SetProductVariant(&variant)
		}
		if options.SelectRelatedWarehouse {
			stock.SetWarehouse(&wareHouse)
		}
		if options.AnnotateAvailableQuantity {
			stock.AvailableQuantity = availableQuantity
		}
		returningStocks = append(returningStocks, &stock)
	}

	return totalCount, returningStocks, nil
}

// FilterForCountryAndChannel finds and returns stocks with given options
func (ss *SqlStockStore) FilterForCountryAndChannel(options *model.StockFilterForCountryAndChannel) ([]*model.Stock, error) {
	warehouseIDQuery := ss.
		warehouseIdSelectQuery(options.CountryCode, options.ChannelSlug).
		PlaceholderFormat(squirrel.Question)

	// remember the order when scan

	query := ss.GetQueryBuilder().
		Select(model.StockTableName+".*", model.WarehouseTableName+".*", model.ProductVariantTableName+".*").
		From(model.StockTableName).
		InnerJoin(model.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)").
		InnerJoin(model.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		Where(squirrel.Expr("Stocks.WarehouseID IN ?", warehouseIDQuery))

	// parse option for FilterVariantStocksForCountry
	// parse additional options
	if options.AnnotateAvailableQuantity {
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
	if options.LockForUpdate && options.Transaction != nil {
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

	runner := ss.GetReplica()
	if options.Transaction != nil {
		runner = options.Transaction
	}

	rows, err := runner.Raw(queryString, args...).Rows()
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

		if options.AnnotateAvailableQuantity {
			scanFields = append(scanFields, &availableQuantity)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of stock, warehouse, product variant")
		}

		stock.SetWarehouse(&wareHouse)
		stock.SetProductVariant(&productVariant)

		returningStocks = append(returningStocks, &stock)

		if options.AnnotateAvailableQuantity {
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
	err := ss.GetMaster().Raw("UPDATE Stocks SET Quantity = Quantity + ? WHERE Id = ?", quantity, stockID).Error
	if err != nil {
		return errors.Wrapf(err, "failed to change stock quantity for stock with id=%s", stockID)
	}

	return nil
}
