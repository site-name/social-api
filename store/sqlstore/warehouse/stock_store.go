package warehouse

import (
	"database/sql"
	"fmt"

	"github.com/mattermost/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gorm.io/gorm"
)

type SqlStockStore struct {
	store.Store
}

func NewSqlStockStore(s store.Store) store.StockStore {
	return &SqlStockStore{Store: s}
}

func (ss *SqlStockStore) BulkUpsert(transaction *gorm.DB, stocks model.StockSlice) (model.StockSlice, error) {
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
	stock, err := model.FindStock(ss.GetReplica(), stockID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Stocks, stockID)
		}
		return nil, err
	}

	return stock, nil
}

func (ss *SqlStockStore) FilterForChannel(options model_helper.StockFilterForChannelOption) (squirrel.Sqlizer, model.StockSlice, error) {
	channelQuery := ss.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.Channels).
		Where(squirrel.Eq{
			model.ChannelTableColumns.ID: model.ShippingZoneChannelTableColumns.ChannelID,
			model.ChannelTableColumns.ID: options.ChannelID,
		}).
		Limit(1).
		Suffix(")")

	shippingZoneChannelQuery := ss.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.ShippingZoneChannels).
		Where(channelQuery).
		Where(squirrel.Eq{
			model.ShippingZoneChannelTableColumns.ShippingZoneID: model.WarehouseShippingZoneTableColumns.ShippingZoneID,
		}).
		Limit(1).
		Suffix(")")

	warehouseShippingZoneQuery := ss.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.WarehouseShippingZones).
		Where(shippingZoneChannelQuery).
		Where(squirrel.Eq{
			model.WarehouseShippingZoneTableColumns.WarehouseID: model.StockTableColumns.WarehouseID,
		}).
		Limit(1).
		Suffix(")")

	// ---
	if options.ReturnQueryOnly {
		return ss.GetQueryBuilder().
			Select(model.TableNames.Stocks).
			From(model.TableNames.Stocks).
			Where(warehouseShippingZoneQuery).
			Where(options.Conditions), nil, nil
	}

	query := ss.GetQueryBuilder().
		Select(model.TableNames.Stocks).
		From(model.TableNames.Stocks).
		Where(warehouseShippingZoneQuery).
		Where(options.Conditions)

	// parse options
	if options.SelectRelatedProductVariant {
		query = query.
			InnerJoin(fmt.Sprintf("%s ON (%s = %s)", model.TableNames.ProductVariants, model.StockTableColumns.ProductVariantID, model.ProductVariantTableColumns.ID))
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

	var returningStocks model.StockSlice

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

func (ss *SqlStockStore) commonQueryBuilder(options model_helper.StockFilterOption) []qm.QueryMod {
	conds := options.Conditions
	for _, load := range options.Preloads {
		conds = append(conds, qm.Load(load))
	}
	if options.Warehouse_ShippingZone_ChannelID != nil ||
		options.Warehouse_ShippingZone_countries != nil ||
		options.Search != "" {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Warehouses, model.WarehouseTableColumns.ID, model.StockTableColumns.WarehouseID)),
		)

		if options.Warehouse_ShippingZone_ChannelID != nil ||
			options.Warehouse_ShippingZone_countries != nil {
			conds = append(
				conds,
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.WarehouseShippingZones, model.WarehouseShippingZoneTableColumns.WarehouseID, model.WarehouseTableColumns.ID)),
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZones, model.ShippingZoneTableColumns.ID, model.WarehouseShippingZoneTableColumns.ShippingZoneID)),
			)
		}

		if options.Warehouse_ShippingZone_ChannelID != nil {
			conds = append(conds, options.Warehouse_ShippingZone_ChannelID)
		}
		if options.Warehouse_ShippingZone_countries != nil {
			conds = append(conds, options.Warehouse_ShippingZone_countries)
		}

		if options.Search != "" {

			searchExpr := fmt.Sprintf("%%%s%%", options.Search)

			conds = append(
				conds,
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Addresses, model.AddressTableColumns.ID, model.WarehouseTableColumns.AddressID)),
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariants, model.ProductVariantTableColumns.ID, model.StockTableColumns.ProductVariantID)),
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Products, model.ProductTableColumns.ID, model.ProductVariantTableColumns.ProductID)),
				model_helper.Or{
					squirrel.ILike{model.ProductTableColumns.Name: searchExpr},
					squirrel.ILike{model.ProductVariantTableColumns.Name: searchExpr},
					squirrel.ILike{model.WarehouseTableColumns.Name: searchExpr},
					squirrel.ILike{model.AddressTableColumns.CompanyName: searchExpr},
				},
			)
		}
	}

	if options.AnnotateAvailableQuantity {
		var annotations = model_helper.AnnotationAggregator{
			model_helper.StockAnnotationKeys.AvailableQuantity: fmt.Sprintf("%s - COALESCE(SUM(%s), 0)", model.StockTableColumns.Quantity, model.AllocationTableColumns.QuantityAllocated),
		}
		conds = append(
			conds,
			qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Allocations, model.AllocationTableColumns.StockID, model.StockTableColumns.ID)),
			qm.Select(model.TableNames.Stocks+".*"), // this is needed
			annotations,
			qm.GroupBy(model.StockTableColumns.ID),
		)
	}

	return conds
}

func (ss *SqlStockStore) FilterByOption(options model_helper.StockFilterOption) (model.StockSlice, error) {
	conds := ss.commonQueryBuilder(options)
	return model.Stocks(conds...).All(ss.GetReplica())
}

func (ss *SqlStockStore) FilterForCountryAndChannel(options model_helper.StockFilterOptionsForCountryAndChannel) (model.StockSlice, error) {
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

func (ss *SqlStockStore) FilterVariantStocksForCountry(options model_helper.StockFilterOptionsForCountryAndChannel) (model.StockSlice, error) {
	return ss.FilterForCountryAndChannel(options)
}

func (ss *SqlStockStore) FilterProductStocksForCountryAndChannel(options model_helper.StockFilterOptionsForCountryAndChannel) (model.StockSlice, error) {
	// TODO: finish me
	return ss.FilterForCountryAndChannel(options)
}

func (ss *SqlStockStore) warehouseIdSelectQuery(countryCode model.CountryCode, channelSlug string) squirrel.SelectBuilder {
	query := ss.GetQueryBuilder().
		Select(model.WarehouseTableColumns.ID).
		From(model.TableNames.Warehouses)

	if countryCode != "" {
		query = query.
			InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.WarehouseShippingZones, model.WarehouseTableColumns.ID, model.WarehouseShippingZoneTableColumns.WarehouseID)).
			InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZones, model.ShippingZoneTableColumns.ID, model.WarehouseShippingZoneTableColumns.ShippingZoneID)).
			Where(squirrel.Like{model.ShippingZoneTableColumns.Countries: fmt.Sprintf("%%%s%%", countryCode)})
	}
	if channelSlug != "" {
		query = query.
			InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZoneChannels, model.ShippingZoneChannelTableColumns.ShippingZoneID, model.ShippingZoneTableColumns.ID)).
			InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.ShippingZoneChannelTableColumns.ChannelID)).
			Where(squirrel.Eq{model.ChannelTableColumns.Slug: channelSlug})
	}

	return query
}

func (ss *SqlStockStore) ChangeQuantity(stockID string, quantityDelta int) error {
	query := fmt.Sprintf("UPDATE %s SET %s = %s + ? WHERE %s = ?", model.TableNames.Stocks, model.StockColumns.Quantity, model.StockColumns.Quantity, model.StockColumns.ID)
	_, err := queries.Raw(query, quantityDelta, stockID).Exec(ss.GetMaster())
	return err
}

func (s *SqlStockStore) Delete(tx boil.ContextTransactor, ids []string) (int64, error) {
	if tx == nil {
		tx = s.GetMaster()
	}
	return model.Stocks(model.StockWhere.ID.IN(ids)).DeleteAll(tx)
}
