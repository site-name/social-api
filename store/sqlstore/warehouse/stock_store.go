package warehouse

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
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

func (ss *SqlStockStore) Save(stock *warehouse.Stock) (*warehouse.Stock, error) {
	stock.PreSave()
	if err := stock.IsValid(); err != nil {
		return nil, err
	}

	if err := ss.GetMaster().Insert(stock); err != nil {
		if ss.IsUniqueConstraintError(err, []string{"WarehouseID", "ProductVariantID", "stocks_warehouseid_productvariantid_key"}) {
			return nil, store.NewErrInvalidInput(store.StockTableName, "WarehouseID/ProductVariantID", stock.WarehouseID+"/"+stock.ProductVariantID)
		}
		return nil, errors.Wrapf(err, "failed to save stock object with id=%s", stock.Id)
	}

	return stock, nil
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

// GetbyOption finds 1 stock by given option then returns it
// func (ss *SqlStockStore) GetbyOption(option *warehouse.StockFilterOption) (*warehouse.Stock, error) {
// 	query := ss.GetQueryBuilder().
// 		Select(ss.ModelFields()...).
// 		From(store.StockTableName)

// 	// parse option
// 	if option.Id != nil {
// 		query = query.Where(option.Id.ToSquirrel("Stocks.Id"))
// 	}
// 	if option.WarehouseID != nil {
// 		query = query.
// 			InnerJoin(store.WarehouseTableName + " ON (Stocks.WarehouseID = Warehouses.Id)").
// 			Where(option.WarehouseID.ToSquirrel("Warehouses.Id"))
// 	}
// 	if option.ProductVariantID != nil {
// 		query = query.
// 			InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
// 			Where(option.ProductVariantID.ToSquirrel("ProductVariants.Id"))
// 	}

// 	queryString, args, err := query.ToSql()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "GetbyOption_ToSql")
// 	}

// 	var res *warehouse.Stock
// 	err = ss.GetReplica().SelectOne(&res, queryString, args...)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, store.NewErrNotFound(store.StockTableName, "option")
// 		}
// 		return nil, errors.Wrap(err, "failed to find a stock with given option")
// 	}

// 	return res, nil
// }

// commonLookup is not exported
func (ss *SqlStockStore) commonLookup(query squirrel.SelectBuilder) ([]*warehouse.Stock, error) {
	var (
		returningStocks []*warehouse.Stock
		stock           warehouse.Stock
		wareHouse       warehouse.WareHouse
		productVariant  product_and_discount.ProductVariant
	)

	rows, err := query.RunWith(ss.GetReplica()).Query()
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

		stock.WareHouse = &wareHouse
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

func (ss *SqlStockStore) FilterForCountryAndChannel(options *warehouse.StockFilterOption) ([]*warehouse.Stock, error) {
	options.CountryCode = strings.ToUpper(options.CountryCode)

	subQuery := ss.warehouseSubQuery(options.CountryCode, options.ChannelSlug)

	selectFields := ss.ModelFields()
	selectFields = append(selectFields, ss.Warehouse().ModelFields()...)
	selectFields = append(selectFields, ss.ProductVariant().ModelFields()...)

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(store.StockTableName).
		InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		InnerJoin(store.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)").
		Where(squirrel.Expr("Stocks.WarehouseID IN ?", subQuery)).
		OrderBy(store.TableOrderingMap[store.StockTableName])

	return ss.commonLookup(query)
}

func (ss *SqlStockStore) FilterVariantStocksForCountry(options *warehouse.StockFilterOption) ([]*warehouse.Stock, error) {
	options.CountryCode = strings.ToUpper(options.CountryCode)

	subQuery := ss.warehouseSubQuery(options.CountryCode, options.ChannelSlug)

	selectFields := ss.ModelFields()
	selectFields = append(selectFields, ss.Warehouse().ModelFields()...)
	selectFields = append(selectFields, ss.ProductVariant().ModelFields()...)

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(store.StockTableName).
		InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		InnerJoin(store.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)").
		Where(squirrel.Expr("Stocks.WarehouseID IN ?", subQuery)).
		Where(squirrel.Expr("Stocks.ProductVariantID = ?", options.ProductVariantID)).
		OrderBy(store.TableOrderingMap[store.StockTableName])

	return ss.commonLookup(query)
}

func (ss *SqlStockStore) FilterProductStocksForCountryAndChannel(options *warehouse.StockFilterOption) ([]*warehouse.Stock, error) {
	options.CountryCode = strings.ToUpper(options.CountryCode)

	subQuery := ss.warehouseSubQuery(options.CountryCode, options.ChannelSlug)

	selectingFields := ss.ModelFields()
	selectingFields = append(selectingFields, ss.Warehouse().ModelFields()...)
	selectingFields = append(selectingFields, ss.ProductVariant().ModelFields()...)

	query := ss.GetQueryBuilder().
		Select(selectingFields...).
		From(store.StockTableName).
		InnerJoin(store.WarehouseTableName + " ON (Warehouses.Id = Stocks.WarehouseID)").
		InnerJoin(store.ProductVariantTableName + " ON (Stocks.ProductVariantID = ProductVariants.Id)").
		InnerJoin(store.ProductTableName + " ON (Products.Id = ProductVariants.ProductID)").
		Where(squirrel.Expr("Stocks.WarehouseID IN ?", subQuery)).
		Where(squirrel.Expr("Products.Id = ?", options.ProductID))

	return ss.commonLookup(query)
}

func (ss *SqlStockStore) warehouseSubQuery(countryCode string, channelSlug string) squirrel.SelectBuilder {
	query := ss.GetQueryBuilder().
		Select("*").
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
