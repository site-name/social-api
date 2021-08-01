package warehouse

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
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
	if res, err := ss.GetReplica().Get(warehouse.Stock{}, stockID); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.StockTableName, stockID)
		}
		return nil, errors.Wrapf(err, "failed to find stock with id=%s", stockID)
	} else {
		return res.(*warehouse.Stock), nil
	}
}

// queryBuildHelperWithOptions common method for building sql query
func queryBuildHelperWithOptions(options *warehouse.ForCountryAndChannelFilter) (string, map[string]interface{}, error) {
	// check if valid country code is provided and valid
	_, exist := model.Countries[options.CountryCode]
	if !exist {
		return "", nil, store.NewErrInvalidInput(store.StockTableName, "countryCode", options.CountryCode)
	}

	subQueryCondition := `Sz.Countries :: text ILIKE :CountryCode`
	query := `SELECT Wh.Id FROM ` + store.WarehouseTableName + ` AS Wh
		INNER JOIN ` + store.WarehouseShippingZoneTableName + ` AS WhSz ON (
			WhSz.WarehouseID = Wh.Id
		)
		INNER JOIN ` + store.ShippingZoneTableName + ` AS Sz ON (
			Sz.Id = WhSz.ShippingZoneID
		)`
	params := map[string]interface{}{
		"CountryCode": "%" + options.CountryCode + "%",
	}

	// if channel slug is provided and valid
	if options.ChannelSlug != "" {
		subQueryCondition += ` AND Cn.Slug = :ChannelSlug`
		query += ` INNER JOIN ` + store.ShippingZoneChannelTableName + ` AS SzCn ON (
			SzCn.ShippingZoneID = Sz.Id
		)
		INNER JOIN ` + store.ChannelTableName + ` AS Cn ON (
			Cn.Id = SzCn.ChannelID
		)`
		params["ChannelSlug"] = options.ChannelSlug
	}
	query += ` WHERE (` + subQueryCondition + `)`

	return query, params, nil
}

// commonLookup is not exported
func (ss *SqlStockStore) commonLookup(query string, params map[string]interface{}) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	rows, err := ss.GetReplica().Query(query, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, store.NewErrNotFound(fmt.Sprintf("%s/%s/%s", store.StockTableName, store.WarehouseTableName, store.ProductVariantTableName), "")
		}
		return nil, nil, nil, errors.Wrapf(err, "failed to perform database lookup operation")
	}

	// defines returning values
	var (
		stocks          []*warehouse.Stock
		warehouses      []*warehouse.WareHouse
		productVariants []*product_and_discount.ProductVariant
	)
	// scan rows
	for rows.Next() {
		var (
			st warehouse.Stock
			wh warehouse.WareHouse
			pv product_and_discount.ProductVariant
		)
		err := rows.Scan(
			// scan for stock:
			&st.Id, &st.WarehouseID, &st.ProductVariantID, &st.Quantity,
			// scan for warehouse:
			&wh.Id, &wh.Name, &wh.Slug, &wh.AddressID, &wh.Email, &wh.Metadata, &wh.PrivateMetadata,
			// scan for product variant:
			&pv.Id, &pv.Name, &pv.ProductID, &pv.Sku, &pv.Weight, &pv.WeightUnit,
			&pv.TrackInventory, &pv.SortOrder, &pv.Metadata, &pv.PrivateMetadata,
		)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "failed to parse a row")
		}
		stocks = append(stocks, &st)
		warehouses = append(warehouses, &wh)
		productVariants = append(productVariants, &pv)
	}

	rows.Close()
	if rows.Err() != nil {
		return nil, nil, nil, errors.Wrap(rows.Err(), "failed to parse rows")
	}

	return stocks, warehouses, productVariants, nil
}

func (ss *SqlStockStore) FilterForCountryAndChannel(options *warehouse.ForCountryAndChannelFilter) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	options.CountryCode = strings.TrimSpace(strings.ToUpper(options.CountryCode))
	options.ChannelSlug = strings.TrimSpace(options.ChannelSlug)

	subQuery, params, err := queryBuildHelperWithOptions(options)
	if err != nil {
		return nil, nil, nil, err
	}

	selects := ss.ModelFields()
	selects = append(selects, ss.Warehouse().ModelFields()...)
	selects = append(selects, ss.ProductVariant().ModelFields()...)
	selectStr := strings.Join(selects, ", ")

	mainQuery := `SELECT ` + selectStr + ` FROM ` + store.StockTableName + ` 
		INNER JOIN ` + store.WarehouseTableName + ` ON (
			Stocks.WarehouseID = Warehouses.Id
		)
		INNER JOIN ` + store.ProductVariantTableName + ` ON (
			ProductVariants.Id = Stocks.ProductVariantID
		)
		WHERE (
			Stocks.WarehouseID IN (` + subQuery + `)
		)
		ORDER BY Stocks.Id ASC`

	return ss.commonLookup(mainQuery, params)
}

func (ss *SqlStockStore) FilterVariantStocksForCountry(options *warehouse.ForCountryAndChannelFilter, productVariantID string) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	options.CountryCode = strings.TrimSpace(strings.ToUpper(options.CountryCode))
	options.ChannelSlug = strings.TrimSpace(options.ChannelSlug)

	subQuery, params, err := queryBuildHelperWithOptions(options)
	if err != nil {
		return nil, nil, nil, err
	}

	selects := ss.ModelFields()
	selects = append(selects, ss.Warehouse().ModelFields()...)
	selects = append(selects, ss.ProductVariant().ModelFields()...)
	selectStr := strings.Join(selects, ", ")

	mainQuery := `SELECT ` + selectStr + ` FROM ` + store.StockTableName + ` 
		INNER JOIN ` + store.WarehouseTableName + ` ON (
			Stocks.WarehouseID = Warehouses.Id
		)
		INNER JOIN ` + store.ProductVariantTableName + ` ON (
			ProductVariants.Id = Stocks.ProductVariantID
		)
		WHERE (
			Stocks.WarehouseID IN (` + subQuery + `)
			AND Stocks.ProductVariantID = :ProductVariantID
		)
		ORDER BY Stocks.Id ASC`

	params["ProductVariantID"] = productVariantID

	return ss.commonLookup(mainQuery, params)
}

func (ss *SqlStockStore) FilterProductStocksForCountryAndChannel(options *warehouse.ForCountryAndChannelFilter, productID string) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	options.CountryCode = strings.TrimSpace(strings.ToUpper(options.CountryCode))
	options.ChannelSlug = strings.TrimSpace(options.ChannelSlug)

	subQuery, params, err := queryBuildHelperWithOptions(options)
	if err != nil {
		return nil, nil, nil, err
	}

	selects := ss.ModelFields()
	selects = append(selects, ss.Warehouse().ModelFields()...)
	selects = append(selects, ss.ProductVariant().ModelFields()...)
	selectStr := strings.Join(selects, ", ")

	mainQuery := `SELECT ` + selectStr + ` FROM ` + store.StockTableName + ` 
		INNER JOIN ` + store.WarehouseTableName + ` ON (
			Stocks.WarehouseID = Warehouses.Id
		)
		INNER JOIN ` + store.ProductVariantTableName + ` ON (
			ProductVariants.Id = Stocks.ProductVariantID
		)
		WHERE (
			Stocks.WarehouseID IN (` + subQuery + `)
			AND ProductVariants.ProductID = :ProductID
		)
		ORDER BY Stocks.Id ASC`

	params["ProductID"] = productID

	return ss.commonLookup(mainQuery, params)
}
