package shipping

import (
	"strings"
	"testing"

	"github.com/mattermost/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/stretchr/testify/require"
)

func TestApplicableShippingMethods(t *testing.T) {
	money, _ := goprices.NewMoney(56.78, "USD")
	_, err := ApplicableShippingMethods(
		money,
		model_helper.NewId(),
		&measurement.ZeroWeight,
		"US",
		[]string{
			model_helper.NewId(),
			model_helper.NewId(),
		},
	)

	require.NoError(t, err)
}

func ApplicableShippingMethods(price *goprices.Money, channelID string, weight *measurement.Weight, countryCode string, productIDs []string) (string, error) {
	selects := []string{
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

		"ShippingZones.Id",
		"ShippingZones.Name",
		"ShippingZones.Contries",
		"ShippingZones.Default",
		"ShippingZones.Description",
		"ShippingZones.Metadata",
		"ShippingZones.PrivateMetadata",

		"ShippingMethodPostalCodeRules.Id",
		"ShippingMethodPostalCodeRules.ShippingMethodID",
		"ShippingMethodPostalCodeRules.Start",
		"ShippingMethodPostalCodeRules.End",
		"ShippingMethodPostalCodeRules.InclusionType",
	}

	priceAmount := price.GetAmount().InexactFloat64()

	params := map[string]any{
		"ChannelID":               channelID,
		"Currency":                price.GetCurrency(),
		"CountryCode":             "%" + countryCode + "%",
		"MinimumOrderPriceAmount": priceAmount,
		"MaximumOrderPriceAmount": priceAmount,
		"MinimumOrderWeight":      weight.Amount,
		"MaximumOrderWeight":      weight.Amount,
		"WeightBasedShippingType": model.ShippingMethodTypeWeight,
		"PriceBasedShipType":      model.ShippingMethodTypePrice,
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

	query := `SELECT ` + strings.Join(selects, ", ") + `,
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

	return squirrel.Dollar.ReplacePlaceholders(query)
}
