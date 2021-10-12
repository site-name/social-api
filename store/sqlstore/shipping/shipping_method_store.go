package shipping

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodStore struct {
	store.Store
}

func NewSqlShippingMethodStore(s store.Store) store.ShippingMethodStore {
	smls := &SqlShippingMethodStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethod{}, store.ShippingMethodTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingZoneID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(shipping.SHIPPING_METHOD_NAME_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(shipping.SHIPPING_METHOD_TYPE_MAX_LENGTH)
		table.ColMap("WeightUnit").SetMaxSize(model.WEIGHT_UNIT_MAX_LENGTH)
	}
	return smls
}

func (s *SqlShippingMethodStore) ModelFields() []string {
	return []string{
		"ShippingMethods.Id",
		"ShippingMethods.Name",
		"ShippingMethods.Type",
		"ShippingMethods.ShippingZoneID",
		"ShippingMethods.MinimumOrderWeight",
		"ShippingMethods.MaximumOrderWeight",
		"ShippingMethods.WeightUnit",
		"ShippingMethods.MaximumDeliveryDays",
		"ShippingMethods.MinimumDeliveryDays",
		"ShippingMethods.Description",
		"ShippingMethods.Metadata",
		"ShippingMethods.PrivateMetadata",
	}
}

func (s *SqlShippingMethodStore) ScanFields(shippingMethod shipping.ShippingMethod) []interface{} {
	return []interface{}{
		&shippingMethod.Id,
		&shippingMethod.Name,
		&shippingMethod.Type,
		&shippingMethod.ShippingZoneID,
		&shippingMethod.MinimumOrderWeight,
		&shippingMethod.MaximumOrderWeight,
		&shippingMethod.WeightUnit,
		&shippingMethod.MaximumDeliveryDays,
		&shippingMethod.MinimumDeliveryDays,
		&shippingMethod.Description,
		&shippingMethod.Metadata,
		&shippingMethod.PrivateMetadata,
	}
}

// TableName returns current store's table name. If non-empty `withField` is provided, it returns the field name of the store, in standard sql style
//
// E.g:
//  TableName("Id") == "ShippingMethods.Id"
//  TableName("") == "ShippingMethods"
func (s *SqlShippingMethodStore) TableName(withField string) string {
	if withField == "" {
		return "ShippingMethods"
	} else {
		return "ShippingMethods." + withField
	}
}

func (s *SqlShippingMethodStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_shipping_methods_name", store.ShippingMethodTableName, "Name")
	s.CreateIndexIfNotExists("idx_shipping_methods_name_lower_textpattern", store.ShippingMethodTableName, "lower(Name) text_pattern_ops")

	s.CreateForeignKeyIfNotExists(store.ShippingMethodTableName, "ShippingZoneID", store.ShippingZoneTableName, "Id", true)
}

// Upsert bases on given method's Id to decide update or insert it
func (s *SqlShippingMethodStore) Upsert(method *shipping.ShippingMethod) (*shipping.ShippingMethod, error) {
	method.PreSave()
	if err := method.IsValid(); err != nil {
		return nil, err
	}

	err := s.GetMaster().Insert(method)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert shipping method with id=%s", method.Id)
	}

	return method, nil
}

// Get finds and returns a shipping method with given id
func (s *SqlShippingMethodStore) Get(methodID string) (*shipping.ShippingMethod, error) {
	var res shipping.ShippingMethod
	err := s.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ShippingMethodTableName+" WHERE Id = :ID", map[string]interface{}{"ID": methodID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShippingMethodTableName, methodID)
		}
		return nil, errors.Wrapf(err, "failed to find shipping method with id=%s", methodID)
	}

	return &res, nil
}

// ApplicableShippingMethods finds all shipping method for given checkout
//
// sql queries here are borrowed. Please check the file shipping_method_store.md
func (s *SqlShippingMethodStore) ApplicableShippingMethods(price *goprices.Money, channelID string, weight *measurement.Weight, countryCode string, productIDs []string) ([]*shipping.ShippingMethod, error) {
	/*
		NOTE: we also prefetch postal_code_rules, shipping zones for later use
		please refer to saleor/shipping/models for details
	*/
	selectFields := append(s.ModelFields(), s.ShippingZone().ModelFields()...)
	selectFields = append(selectFields, s.ShippingMethodPostalCodeRule().ModelFields()...)

	priceAmount, _ := price.Amount.Float64()

	params := map[string]interface{}{
		"ChannelID":               channelID,
		"Currency":                price.Currency,
		"CountryCode":             "%" + countryCode + "%",
		"MinimumOrderPriceAmount": priceAmount,
		"MaximumOrderPriceAmount": priceAmount,
		"MinimumOrderWeight":      *weight.Amount,
		"MaximumOrderWeight":      *weight.Amount,
		"WeightBasedShippingType": shipping.WEIGHT_BASED,
		"PriceBasedShipType":      shipping.PRICE_BASED,
	}

	// check if productIDs is provided:
	var forExcludedProductQuery string
	if len(productIDs) > 0 {
		forExcludedProductQuery = `
		AND NOT (
			EXISTS(
				SELECT
					(1) AS a
				FROM
					ShippingMethodExcludedProducts
				WHERE (
					ShippingMethodExcludedProducts.ProductID IN :ExcludedProductIDs
					AND ShippingMethodExcludedProducts.ShippingMethodID = ShippingMethods.Id
				)
				LIMIT 1
			)
		)`
		// update params also
		params["ExcludedProductIDs"] = productIDs
	}

	query := `SELECT ` + strings.Join(selectFields, ", ") + `,
	(
		SELECT
			ShippingMethodChannelListings.PriceAmount
		FROM
			ShippingMethodChannelListings
		WHERE (
			ShippingMethodChannelListings.ChannelID = :ChannelID
			AND ShippingMethodChannelListings.ShippingMethodID = ShippingMethods.Id
		)
	) AS PriceAmount
	FROM
		ShippingMethods
	INNER JOIN ShippingMethodChannelListings ON (
		ShippingMethodChannelListings.ShippingMethodID = ShippingMethods.Id
	)
	INNER JOIN ShippingZones ON (
		ShippingZones.Id = ShippingMethods.ShippingZoneID
	)
	INNER JOIN ShippingZoneChannels ON (
		ShippingZones.Id = ShippingZoneChannels.ShippingZoneID
	)
	INNER JOIN ShippingMethodPostalCodeRules ON (
		ShippingMethodPostalCodeRules.ShippingMethodID = ShippingMethods.Id
	)
	WHERE
		(
			(
				ShippingMethodChannelListings.ChannelID = :ChannelID
				AND ShippingMethodChannelListings.Currency = :Currency
				AND ShippingZoneChannels.ChannelID = :ChannelID
				AND ShippingZones.Countries :: text LIKE :CountryCode ` + forExcludedProductQuery + `
				AND ShippingMethods.Type = :PriceBasedShipType
				AND ShippingMethods.Id IN (
				SELECT
					ShippingMethodID
				FROM
					ShippingMethodChannelListings
				WHERE (
					ShippingMethodChannelListings.ChannelID = :ChannelID
					AND ShippingMethodChannelListings.ShippingMethodID IN (
						SELECT
							Id
						FROM
							ShippingMethods
						INNER JOIN ShippingMethodChannelListings ON (
							ShippingMethodChannelListings.ShippingMethodID = ShippingMethods.Id
						)
						INNER JOIN ShippingZones ON (
							ShippingMethods.ShippingZoneID = ShippingZones.Id
						)
						INNER JOIN ShippingZoneChannels ON (
							ShippingZoneChannels.ShippingZoneID = ShippingZones.Id
						)
						WHERE (
							ShippingMethodChannelListings.ChannelID = :ChannelID
							AND ShippingMethodChannelListings.Currency = :Currency
							AND ShippingZoneChannels.ChannelID = :ChannelID
							AND ShippingZones.Countries :: text LIKE :CountryCode
							AND ShippingMethods.Type = :PriceBasedShipType ` + forExcludedProductQuery + `
						)
					)
					AND ShippingMethodChannelListings.MinimumOrderPriceAmount <= :MinimumOrderPriceAmount
					AND (
						ShippingMethodChannelListings.MaximumOrderPriceAmount IS NULL
						OR ShippingMethodChannelListings.MaximumOrderPriceAmount >= :MaximumOrderPriceAmount
					)
				)
			)
			OR (
				ShippingMethodChannelListings.ChannelID = :ChannelID
				AND ShippingMethodChannelListings.Currency = :Currency
				AND ShippingZoneChannels.ChannelID = :ChannelID
				AND ShippingZones.Countries :: text LIKE :CountryCode ` + forExcludedProductQuery + `
				AND ShippingMethods.Type = :WeightBasedShippingType
				AND (
					ShippingMethods.MinimumOrderWeight <= :MinimumOrderWeight
					OR ShippingMethods.MinimumOrderWeight IS NULL
				)
				AND (
					ShippingMethods.MaximumOrderWeight >= :MaximumOrderWeight
					OR ShippingMethods.MaximumOrderWeight IS NULL
				)
			)
		)
	ORDER BY PriceAmount ASC`

	// this struct inherits entirely from `ShippingMethod`, `ShippingZone`, `ShippingMethodPostalCodeRule`
	var holder []*struct {
		// fields of shipping method
		MethodID              string
		MethodName            string
		Type                  string
		ShippingZoneID        string
		MinimumOrderWeight    float32
		MaximumOrderWeight    *float32
		WeightUnit            measurement.WeightUnit
		MaximumDeliveryDays   *uint
		MinimumDeliveryDays   *uint
		MethodDescription     *model.StringInterface
		MethodMetadata        model.StringMap
		MethodPrivateMetadata model.StringMap

		// fields of shipping zone
		ShippingZoneId              string
		ShippingZoneName            string
		Countries                   string
		Default                     *bool
		Description                 string
		ShippingZoneMetadata        model.StringMap
		ShippingZonePrivateMetadata model.StringMap

		// fields for postal code rules
		PostalCodeId     string
		ShippingMethodID string
		Start            string
		End              string
		InclusionType    string
	}

	// use Select() here since it can inteprets map[string]interface{} value mapping
	// Query() cannot understand.
	_, err := s.GetReplica().Select(&holder, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to finds shipping methods for given conditions")
	}

	var (
		shippingMethods []*shipping.ShippingMethod

		// since each model instance has got a unique UUID, so 1 meetMap is ok.
		//
		// 1 shipping method can have multiple shipping zones.
		// 1 shipping method can have multiple shipping method postal code rules.
		meetMap = map[string]bool{}
	)

	for _, item := range holder {
		// check if item is not nil
		if item != nil {
			var shippingMethod *shipping.ShippingMethod

			// check if this method is already parsed
			if _, exist := meetMap[item.MethodID]; !exist {
				shippingMethod = &shipping.ShippingMethod{
					Id:                  item.MethodID,
					Name:                item.MethodName,
					Type:                item.Type,
					ShippingZoneID:      item.ShippingZoneID,
					MinimumOrderWeight:  item.MinimumOrderWeight,
					MaximumOrderWeight:  item.MaximumOrderWeight,
					WeightUnit:          item.WeightUnit,
					MaximumDeliveryDays: item.MaximumDeliveryDays,
					MinimumDeliveryDays: item.MinimumDeliveryDays,
					Description:         item.MethodDescription,
					ModelMetadata: model.ModelMetadata{
						Metadata:        item.MethodMetadata,
						PrivateMetadata: item.MethodPrivateMetadata,
					},
				}

				// append shipping method to the defined slice
				shippingMethods = append(shippingMethods, shippingMethod)

				// let meepMap hold parsed value for later check
				meetMap[item.MethodID] = true
			}

			// check if this zone is already parsed:
			if _, exist := meetMap[item.ShippingZoneId]; !exist && shippingMethod != nil {
				shippingMethod.ShippingZones = append(
					shippingMethod.ShippingZones,
					&shipping.ShippingZone{
						Id:          item.ShippingZoneId,
						Name:        item.ShippingZoneName,
						Countries:   item.Countries,
						Default:     item.Default,
						Description: item.Description,
						ModelMetadata: model.ModelMetadata{
							Metadata:        item.ShippingZoneMetadata,
							PrivateMetadata: item.ShippingZonePrivateMetadata,
						},
					},
				)

				// let meetMap hold the id for later check
				meetMap[item.ShippingZoneId] = true
			}

			// check if this postal code rule is already parsed:
			if _, exist := meetMap[item.PostalCodeId]; !exist && shippingMethod != nil {
				shippingMethod.ShippingMethodPostalCodeRules = append(
					shippingMethod.ShippingMethodPostalCodeRules,
					&shipping.ShippingMethodPostalCodeRule{
						Id:               item.PostalCodeId,
						ShippingMethodID: item.ShippingMethodID,
						Start:            item.Start,
						End:              item.End,
						InclusionType:    item.InclusionType,
					},
				)

				// let meetMap hold the id for later check:
				meetMap[item.PostalCodeId] = true
			}
		}
	}

	return shippingMethods, nil
}

// GetbyOption finds and returns a shipping method that satisfy given options
func (ss *SqlShippingMethodStore) GetbyOption(options *shipping.ShippingMethodFilterOption) (*shipping.ShippingMethod, error) {
	query := ss.GetQueryBuilder().
		Select(ss.ModelFields()...).
		From(store.ShippingMethodTableName)

	// parse options
	if options.Id != nil {
		// query = query.Where(options.Id.ToSquirrel("ShippingMethods.Id"))
		query = query.Where(options.Id)
	}
	if options.Type != nil {
		// query = query.Where(options.Type.ToSquirrel("ShippingMethods.Type"))
		query = query.Where(options.Type)
	}
	if options.MinimumOrderWeight != nil {
		// query = query.Where(options.MinimumOrderWeight.ToSquirrel("ShippingMethods.MinimumOrderWeight"))
		query = query.Where(options.MinimumOrderWeight)
	}
	if options.MaximumOrderWeight != nil {
		// query = query.Where(options.MaximumOrderWeight.ToSquirrel("ShippingMethods.MaximumOrderWeight"))
		query = query.Where(options.MaximumOrderWeight)
	}
	if options.ShippingZoneChannelSlug != nil {
		// query = query.Where(options.ShippingZoneChannelSlug.ToSquirrel("ShippingMethods.ShippingZoneChannelSlug"))
		query = query.Where(options.ShippingZoneChannelSlug)
	}
	if options.ChannelListingsChannelSlug != nil {
		// query = query.Where(options.ChannelListingsChannelSlug.ToSquirrel("ShippingMethods.ChannelListingsChannelSlug"))
		query = query.Where(options.ChannelListingsChannelSlug)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}

	var res shipping.ShippingMethod
	err = ss.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShippingMethodTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find shipping method by given options")
	}

	return &res, nil
}
