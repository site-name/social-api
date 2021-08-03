package shipping

import (
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
	result, err := s.GetReplica().Get(shipping.ShippingMethod{}, methodID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find shipping method with id=%s", methodID)
	}

	return result.(*shipping.ShippingMethod), nil
}

// ApplicableShippingMethodsForCheckout finds all shipping method for given checkout
//
// sql queries here are borrowed. Please check the file .md
func (s *SqlShippingMethodStore) ApplicableShippingMethods(price *goprices.Money, channelID string, weight *measurement.Weight, countryCode string, productIDs []string) ([]*shipping.ShippingMethod, []*shipping.ShippingZone, []*shipping.ShippingMethodPostalCodeRule, error) {
	/*
		NOTE: we also prefetch postal_code_rules, shipping zones for later use
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
		"MinimumOrderWeight":      weight.Amount,
		"MaximumOrderWeight":      weight.Amount,
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

	var results []*struct {

		// postal code
		Id               string
		ShippingMethodID string
		Start            string
		End              string
		InclusionType    string
	}
	_, err := s.GetReplica().Select(
		&results,
		query,
		params,
	)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to finds shipping methods and related values")
	}

	return nil, nil, nil, nil
}
