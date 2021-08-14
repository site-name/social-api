package warehouse

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/measurement"
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

// queryBuildHelperWithOptions common method for building sql query
func queryBuildHelperWithOptions(options *warehouse.ForCountryAndChannelFilter) (string, map[string]interface{}, error) {
	// check if valid country code is provided and valid
	_, exist := model.Countries[options.CountryCode]
	if !exist {
		return "", nil, store.NewErrInvalidInput(store.StockTableName, "countryCode", options.CountryCode)
	}

	subQueryCondition := `ShippingZones.Countries :: text ILIKE :CountryCode`
	query := `SELECT Warehouses.Id FROM ` + store.WarehouseTableName + `
		INNER JOIN ` + store.WarehouseShippingZoneTableName + ` ON (
			WarehouseShippingZones.WarehouseID = Warehouses.Id
		)
		INNER JOIN ` + store.ShippingZoneTableName + ` ON (
			ShippingZones.Id = WarehouseShippingZones.ShippingZoneID
		)`
	params := map[string]interface{}{
		"CountryCode": "%" + options.CountryCode + "%",
		"OrderBy":     store.TableOrderingMap[store.StockTableName], // Orderby added here
	}

	// if channel slug is provided and valid
	if options.ChannelSlug != "" {
		subQueryCondition += ` AND Channels.Slug = :ChannelSlug`
		query += ` INNER JOIN ` + store.ShippingZoneChannelTableName + ` ON (
			ShippingZoneChannels.ShippingZoneID = ShippingZones.Id
		)
		INNER JOIN ` + store.ChannelTableName + ` ON (
			ChannelTableName.Id = ShippingZoneChannels.ChannelID
		)`
		params["ChannelSlug"] = options.ChannelSlug
	}
	query += ` WHERE (` + subQueryCondition + `)`

	return query, params, nil
}

// commonLookup is not exported
func (ss *SqlStockStore) commonLookup(query string, params map[string]interface{}) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	var selectedRows []*struct {
		Id               string
		CreateAt         int64
		WarehouseID      string
		ProductVariantID string
		Quantity         uint

		WareHouseID     string
		Name            string
		Slug            string
		AddressID       *string
		Email           string
		Metadata        model.StringMap
		PrivateMetadata model.StringMap

		VariantId              string
		VariantName            string
		ProductID              string
		Sku                    string
		Weight                 *float32
		WeightUnit             measurement.WeightUnit
		TrackInventory         *bool
		SortOrder              int
		VariantMetadata        model.StringMap
		VariantPrivateMetadata model.StringMap
	}
	_, err := ss.GetReplica().Select(&selectedRows, query, params)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to perform database lookup operation")
	}

	// defines returning values
	var (
		stocks          []*warehouse.Stock
		warehouses      []*warehouse.WareHouse
		productVariants []*product_and_discount.ProductVariant
	)

	for _, row := range selectedRows {
		stocks = append(stocks, &warehouse.Stock{
			Id:               row.Id,
			CreateAt:         row.CreateAt,
			WarehouseID:      row.WarehouseID,
			ProductVariantID: row.ProductVariantID,
			Quantity:         row.Quantity,
		})
		warehouses = append(warehouses, &warehouse.WareHouse{
			Id:        row.WareHouseID,
			Name:      row.Name,
			Slug:      row.Slug,
			AddressID: row.AddressID,
			Email:     row.Email,
			ModelMetadata: model.ModelMetadata{
				Metadata:        row.Metadata,
				PrivateMetadata: row.PrivateMetadata,
			},
		})
		productVariants = append(productVariants, &product_and_discount.ProductVariant{
			Id:             row.VariantId,
			Name:           row.VariantName,
			ProductID:      row.ProductID,
			Sku:            row.Sku,
			Weight:         row.Weight,
			WeightUnit:     row.WeightUnit,
			TrackInventory: row.TrackInventory,
			Sortable: model.Sortable{
				SortOrder: row.SortOrder,
			},
			ModelMetadata: model.ModelMetadata{
				Metadata:        row.VariantMetadata,
				PrivateMetadata: row.VariantPrivateMetadata,
			},
		})
	}

	return stocks, warehouses, productVariants, nil
}

func (ss *SqlStockStore) FilterForCountryAndChannel(options *warehouse.ForCountryAndChannelFilter) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	options.CountryCode = strings.ToUpper(options.CountryCode)

	subQuery, params, err := queryBuildHelperWithOptions(options)
	if err != nil {
		return nil, nil, nil, err
	}

	selectingFields := ss.ModelFields()
	selectingFields = append(selectingFields, ss.Warehouse().ModelFields()...)
	selectingFields = append(selectingFields, ss.ProductVariant().ModelFields()...)

	mainQuery := `SELECT ` + strings.Join(selectingFields, ", ") + `
	FROM ` + store.StockTableName + `
	INNER JOIN ` + store.WarehouseTableName + ` ON (
		Stocks.WarehouseID = Warehouses.Id
	)
	INNER JOIN ` + store.ProductVariantTableName + ` ON (
		ProductVariants.Id = Stocks.ProductVariantID
	)
	WHERE (
		Stocks.WarehouseID IN (` + subQuery + `)
	)
	ORDER BY :OrderBy`

	return ss.commonLookup(mainQuery, params)
}

func (ss *SqlStockStore) FilterVariantStocksForCountry(options *warehouse.ForCountryAndChannelFilter, productVariantID string) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	options.CountryCode = strings.ToUpper(options.CountryCode)

	subQuery, params, err := queryBuildHelperWithOptions(options)
	if err != nil {
		return nil, nil, nil, err
	}

	selectingFields := ss.ModelFields()
	selectingFields = append(selectingFields, ss.Warehouse().ModelFields()...)
	selectingFields = append(selectingFields, ss.ProductVariant().ModelFields()...)

	mainQuery := `SELECT ` + strings.Join(selectingFields, ", ") + `
	FROM ` + store.StockTableName + ` 
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
	ORDER BY :OrderBy`

	params["ProductVariantID"] = productVariantID

	return ss.commonLookup(mainQuery, params)
}

func (ss *SqlStockStore) FilterProductStocksForCountryAndChannel(options *warehouse.ForCountryAndChannelFilter, productID string) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) {
	options.CountryCode = strings.ToUpper(options.CountryCode)

	subQuery, params, err := queryBuildHelperWithOptions(options)
	if err != nil {
		return nil, nil, nil, err
	}

	selectingFields := ss.ModelFields()
	selectingFields = append(selectingFields, ss.Warehouse().ModelFields()...)
	selectingFields = append(selectingFields, ss.ProductVariant().ModelFields()...)

	mainQuery := `SELECT ` + strings.Join(selectingFields, ", ") + ` FROM ` + store.StockTableName + ` 
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
		ORDER BY :OrderBy`

	params["ProductID"] = productID

	return ss.commonLookup(mainQuery, params)
}
