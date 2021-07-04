package warehouse

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore/channel"
	"github.com/sitename/sitename/store/sqlstore/product"
	"github.com/sitename/sitename/store/sqlstore/shipping"
)

const (
	StockTableName = "Stocks"
)

type SqlStockStore struct {
	store.Store
}

var StockQuery = []string{
	"St.Id",
	"St.WarehouseID",
	"St.ProductVariantID",
	"St.Quantity",
}

func NewSqlStockStore(s store.Store) store.StockStore {
	ws := &SqlStockStore{
		Store: s,
	}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(warehouse.Stock{}, StockTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("WarehouseID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("WarehouseID", "ProductVariantID")
	}
	return ws
}

func (ws *SqlStockStore) CreateIndexesIfNotExists() {
	ws.CreateForeignKeyIfNotExists(StockTableName, "WarehouseID", WarehouseTableName, "Id", true)
	ws.CreateForeignKeyIfNotExists(StockTableName, "ProductVariantID", product.ProductVariantTableName, "Id", true)
}

func (ws *SqlStockStore) Save(stock *warehouse.Stock) (*warehouse.Stock, error) {
	stock.PreSave()
	if err := stock.IsValid(); err != nil {
		return nil, err
	}

	if err := ws.GetMaster().Insert(stock); err != nil {
		if ws.IsUniqueConstraintError(err, []string{"WarehouseID", "ProductVariantID"}) {
			return nil, store.NewErrInvalidInput(StockTableName, "WarehouseID/ProductVariantID", stock.WarehouseID+"/"+stock.ProductVariantID)
		}
		return nil, errors.Wrapf(err, "failed to save stock object with id=%s", stock.Id)
	}

	return stock, nil
}

func (ws *SqlStockStore) Get(stockID string) (*warehouse.Stock, error) {
	if res, err := ws.GetReplica().Get(warehouse.Stock{}, stockID); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(StockTableName, stockID)
		}
		return nil, errors.Wrapf(err, "failed to find stock with id=%s", stockID)
	} else {
		return res.(*warehouse.Stock), nil
	}
}

// queryBuildHelperWithOptions common method for building sql query
func queryBuildHelperWithOptions(options *warehouse.ForCountryAndChannelFilter) (string, error) {
	// check if valid country code is provided
	_, exist := model.Countries[options.CountryCode]
	if !exist {
		return "", store.NewErrInvalidInput(StockTableName, "countryCode", options.CountryCode)
	}

	subQueryCondition := `Sz.Countries :: text ILIKE :CountryCode`
	subQuery := `SELECT Wh.Id FROM ` + WarehouseTableName + ` AS Wh
		INNER JOIN ` + WarehouseShippingZoneTableName + ` AS WhSz ON (
			WhSz.WarehouseID = Wh.Id
		)
		INNER JOIN ` + shipping.ShippingZoneTableName + ` AS Sz ON (
			Sz.Id = WhSz.ShippingZoneID
		)`

	// if channel slug is provided
	if options.ChannelSlug != "" {
		subQueryCondition += ` AND Cn.Slug = :ChannelSlug`
		subQuery += ` INNER JOIN ` + shipping.ShippingZoneChannelTableName + ` AS SzCn ON (
			SzCn.ShippingZoneID = Sz.Id
		)
		INNER JOIN ` + channel.ChannelTableName + ` AS Cn ON (
			Cn.Id = SzCn.ChannelID
		)`
	}
	subQuery += ` WHERE (` + subQueryCondition + `)`

	return subQuery, nil
}

// commonLookup is not exported
func (ss *SqlStockStore) commonLookup(query string, params map[string]interface{}) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	rows, err := ss.GetReplica().Query(query, params)
	if err != nil {
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
			&st.Id, &st.WarehouseID, &st.ProductVariantID, &st.Quantity, // scan for stock
			&wh.Id, &wh.Name, &wh.Slug, &wh.AddressID, &wh.Email, &wh.Metadata, &wh.PrivateMetadata, // scan for warehouse
			&pv.Id, &pv.Name, &pv.ProductID, &pv.Sku, &pv.Weight, &pv.WeightUnit, // scan for product variant
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

func (ss *SqlStockStore) FilterVariantStocksForCountry(options *warehouse.ForCountryAndChannelFilter, productVariantID string) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	options.CountryCode = strings.TrimSpace(strings.ToUpper(options.CountryCode))
	options.ChannelSlug = strings.TrimSpace(options.ChannelSlug)

	subQuery, err := queryBuildHelperWithOptions(options)
	if err != nil {
		return nil, nil, nil, err
	}

	selects := StockQuery
	selects = append(selects, WarehouseQuery...)
	selects = append(selects, product.ProductVariantQuery...)
	selectStr := strings.Join(selects, ", ")

	mainQuery := `SELECT ` + selectStr +
		` FROM ` + StockTableName + ` AS St 
		INNER JOIN ` + WarehouseTableName + ` AS Wh ON (
			St.WarehouseID = Wh.Id
		)
		INNER JOIN ` + product.ProductVariantTableName + ` AS Pv ON (
			Pv.Id = St.ProductVariantID
		)
		WHERE (
			St.WarehouseID IN (` + subQuery + `)
			AND St.ProductVariantID = :ProductVariantID
		) ORDER BY St.Id ASC`

	params := map[string]interface{}{
		"CountryCode":      "%" + options.CountryCode + "%",
		"ProductVariantID": productVariantID,
	}
	if options.ChannelSlug != "" {
		params["ChannelSlug"] = options.ChannelSlug
	}

	return ss.commonLookup(mainQuery, params)
}

func (ss *SqlStockStore) FilterProductStocksForCountryAndChannel(options *warehouse.ForCountryAndChannelFilter, productID string) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	options.CountryCode = strings.TrimSpace(strings.ToUpper(options.CountryCode))
	options.ChannelSlug = strings.TrimSpace(options.ChannelSlug)

	subQuery, err := queryBuildHelperWithOptions(options)
	if err != nil {
		return nil, nil, nil, err
	}

	selects := StockQuery
	selects = append(selects, WarehouseQuery...)
	selects = append(selects, product.ProductVariantQuery...)
	selectStr := strings.Join(selects, ", ")

	mainQuery := `SELECT ` + selectStr +
		` FROM ` + StockTableName + ` AS St 
		INNER JOIN ` + WarehouseTableName + ` AS Wh ON (
			St.WarehouseID = Wh.Id
		)
		INNER JOIN ` + product.ProductVariantTableName + ` AS Pv ON (
			Pv.Id = St.ProductVariantID
		)
		WHERE (
			St.WarehouseID IN (` + subQuery + `)
			AND Pv.ProductID = :ProductID
		) ORDER BY St.Id ASC`

	params := map[string]interface{}{
		"CountryCode": "%" + options.CountryCode + "%",
		"ProductID":   productID,
	}
	if options.ChannelSlug != "" {
		params["ChannelSlug"] = options.ChannelSlug
	}

	return ss.commonLookup(mainQuery, params)
}
