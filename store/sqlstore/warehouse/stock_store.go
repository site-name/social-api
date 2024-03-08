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
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
)

type SqlStockStore struct {
	store.Store
}

func NewSqlStockStore(s store.Store) store.StockStore {
	return &SqlStockStore{Store: s}
}

func (ss *SqlStockStore) Upsert(transaction boil.ContextTransactor, stocks model.StockSlice) (model.StockSlice, error) {
	if transaction == nil {
		transaction = ss.GetMaster()
	}

	for _, stock := range stocks {
		if stock == nil {
			continue
		}

		isSaving := stock.ID == ""
		if isSaving {
			model_helper.StockPreSave(stock)
		}

		if err := model_helper.StockIsValid(*stock); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = stock.Insert(transaction, boil.Infer())
		} else {
			_, err = stock.Update(transaction, boil.Blacklist(
				model.StockColumns.CreatedAt,
				model.StockColumns.WarehouseID,
			))
		}

		if err != nil {
			if ss.IsUniqueConstraintError(err, []string{"stocks_warehouse_id_product_variant_id_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.Stocks, "WarehouseID/ProductVariantID", "duplicate")
			}
			return nil, err
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

func (ss *SqlStockStore) commonFilterForChannelQuery(options model_helper.StockFilterForChannelOption) squirrel.SelectBuilder {
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

	return ss.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.WarehouseShippingZones).
		Where(shippingZoneChannelQuery).
		Where(squirrel.Eq{
			model.WarehouseShippingZoneTableColumns.WarehouseID: model.StockTableColumns.WarehouseID,
		}).
		Limit(1).
		Suffix(")")
}

func (ss *SqlStockStore) GetFilterForChannelQuery(options model_helper.StockFilterForChannelOption) squirrel.SelectBuilder {
	warehouseShippingZoneQuery := ss.commonFilterForChannelQuery(options)

	return ss.GetQueryBuilder().
		Select(model.TableNames.Stocks + ".*").
		From(model.TableNames.Stocks).
		Where(warehouseShippingZoneQuery).
		Where(options.Conditions)
}

func (ss *SqlStockStore) FilterForChannel(options model_helper.StockFilterForChannelOption) (model.StockSlice, error) {
	query, args, err := ss.commonFilterForChannelQuery(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create commonFilterForChannelQuery")
	}

	conds := []qm.QueryMod{
		qm.Load(model.StockRels.ProductVariant),
		qmhelper.WhereQueryMod{
			Clause: query,
			Args:   args,
		},
	}

	return model.Stocks(conds...).All(ss.GetReplica())
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
	warehouseIDQuery, args, err := ss.
		warehouseIdSelectQuery(options.CountryCode, options.ChannelSlug).
		PlaceholderFormat(squirrel.Question).
		Prefix("(").
		Suffix(")").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create warehouseIDQuery")
	}

	conds := []qm.QueryMod{
		qm.Load(model.StockRels.Warehouse),
		qm.Load(model.StockRels.ProductVariant),
		qmhelper.WhereQueryMod{
			Clause: model.StockTableColumns.WarehouseID + " IN " + warehouseIDQuery,
			Args:   args,
		},
	}

	return model.Stocks(conds...).All(ss.GetReplica())
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

func (ss *SqlStockStore) FilterVariantStocksForCountry(options model_helper.StockFilterVariantStocksForCountryFilterOptions) (model.StockSlice, error) {
	stocks, err := ss.FilterForCountryAndChannel(options.StockFilterOptionsForCountryAndChannel)
	if err != nil {
		return nil, err
	}

	var stocksOfGivenVariant = make(model.StockSlice, 0, len(stocks))
	for _, stock := range stocks {
		if stock != nil &&
			stock.R != nil &&
			stock.R.ProductVariant != nil &&
			stock.R.ProductVariant.ID == options.ProductVariantID {
			stocksOfGivenVariant = append(stocksOfGivenVariant, stock)
		}
	}

	return stocksOfGivenVariant, nil
}

func (ss *SqlStockStore) FilterProductStocksForCountryAndChannel(options model_helper.StockFilterProductStocksForCountryAndChannelFilterOptions) (model.StockSlice, error) {
	stocks, err := ss.FilterForCountryAndChannel(options.StockFilterOptionsForCountryAndChannel)
	if err != nil {
		return nil, err
	}

	var resultStocks = make(model.StockSlice, 0, len(stocks))
	for _, stock := range stocks {
		if stock != nil &&
			stock.R != nil &&
			stock.R.ProductVariant != nil &&
			stock.R.ProductVariant.ProductID == options.ProductID {
			resultStocks = append(resultStocks, stock)
		}
	}

	return resultStocks, nil
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
