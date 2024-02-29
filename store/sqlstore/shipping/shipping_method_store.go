package shipping

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlShippingMethodStore struct {
	store.Store
}

func NewSqlShippingMethodStore(s store.Store) store.ShippingMethodStore {
	return &SqlShippingMethodStore{s}
}

func (s *SqlShippingMethodStore) Upsert(transaction boil.ContextTransactor, method model.ShippingMethod) (*model.ShippingMethod, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	isSaving := method.ID == ""
	if isSaving {
		model_helper.ShippingMethodPreSave(&method)
	} else {
		model_helper.ShippingMethodCommonPre(&method)
	}

	if err := model_helper.ShippingMethodIsValid(method); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = method.Insert(transaction, boil.Infer())
	} else {
		_, err = method.Update(transaction, boil.Infer())
	}
	if err != nil {
		return nil, err
	}

	return &method, nil
}

func (s *SqlShippingMethodStore) Get(methodID string) (*model.ShippingMethod, error) {
	method, err := model.FindShippingMethod(s.GetReplica(), methodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ShippingMethods, methodID)
		}
		return nil, err
	}

	return method, nil
}

func (s *SqlShippingMethodStore) ApplicableShippingMethods(price goprices.Money, channelID string, weight measurement.Weight, countryCode model.CountryCode, productIDs []string) ([]*model.ShippingMethod, error) {
	selectFields := []string{
		model.ShippingMethodTableName + ".*",
		model.ShippingZoneTableName + ".*",
		model.ShippingMethodPostalCodeRuleTableName + ".*",
	}

	priceAmount := price.Amount.InexactFloat64()

	params := map[string]interface{}{
		"ChannelID":               channelID,
		"Currency":                price.Currency,
		"CountryCode":             "%" + countryCode + "%",
		"MinimumOrderPriceAmount": priceAmount,
		"MaximumOrderPriceAmount": priceAmount,
		"MinimumOrderWeight":      weight.Amount,
		"MaximumOrderWeight":      weight.Amount,
		"WeightBasedShippingType": model.WEIGHT_BASED,
		"PriceBasedShipType":      model.PRICE_BASED,
	}

	// check if productIDs is provided:
	var forExcludedProductQuery string
	if len(productIDs) > 0 {
		forExcludedProductQuery = `AND NOT (
			EXISTS(
				SELECT
					(1) AS "a"
				FROM
					ShippingMethodExcludedProducts
				WHERE (
					ShippingMethodExcludedProducts.ProductID IN @ExcludedProductIDs
					AND ShippingMethodExcludedProducts.ShippingMethodID = ShippingMethods.Id
				)
				LIMIT 1
			)
		)`
		// update params also
		params["ExcludedProductIDs"] = productIDs
	}

	query := `SELECT ` + strings.Join(selectFields, ",") + `,
	(
		SELECT
			ShippingMethodChannelListings.PriceAmount
		FROM
			ShippingMethodChannelListings
		WHERE (
			ShippingMethodChannelListings.ChannelID = @ChannelID
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
				ShippingMethodChannelListings.ChannelID = @ChannelID
				AND ShippingMethodChannelListings.Currency = @Currency
				AND ShippingZoneChannels.ChannelID = @ChannelID
				AND ShippingZones.Countries::text LIKE @CountryCode ` + forExcludedProductQuery + `
				AND ShippingMethods.Type = @PriceBasedShipType
				AND ShippingMethods.Id IN (
				SELECT
					ShippingMethodID
				FROM
					ShippingMethodChannelListings
				WHERE (
					ShippingMethodChannelListings.ChannelID = @ChannelID
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
							ShippingMethodChannelListings.ChannelID = @ChannelID
							AND ShippingMethodChannelListings.Currency = @Currency
							AND ShippingZoneChannels.ChannelID = @ChannelID
							AND ShippingZones.Countries::text LIKE @CountryCode
							AND ShippingMethods.Type = @PriceBasedShipType ` + forExcludedProductQuery + `
						)
					)
					AND ShippingMethodChannelListings.MinimumOrderPriceAmount <= @MinimumOrderPriceAmount
					AND (
						ShippingMethodChannelListings.MaximumOrderPriceAmount IS NULL
						OR ShippingMethodChannelListings.MaximumOrderPriceAmount >= @MaximumOrderPriceAmount
					)
				)
			)
			OR (
				ShippingMethodChannelListings.ChannelID = @ChannelID
				AND ShippingMethodChannelListings.Currency = @Currency
				AND ShippingZoneChannels.ChannelID = @ChannelID
				AND ShippingZones.Countries::text LIKE @CountryCode ` + forExcludedProductQuery + `
				AND ShippingMethods.Type = @WeightBasedShippingType
				AND (
					ShippingMethods.MinimumOrderWeight <= @MinimumOrderWeight
					OR ShippingMethods.MinimumOrderWeight IS NULL
				)
				AND (
					ShippingMethods.MaximumOrderWeight >= @MaximumOrderWeight
					OR ShippingMethods.MaximumOrderWeight IS NULL
				)
			)
		)
	ORDER BY PriceAmount ASC`

	// use Select() here since it can inteprets map[string]interface{} value mapping
	// Query() cannot understand.
	rows, err := s.GetReplica().Raw(query, params).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to finds shipping methods for given conditions")
	}
	defer rows.Close()

	var shippingMethodMeetMap = map[string]*model.ShippingMethod{}

	for rows.Next() {
		var (
			shippingMethod model.ShippingMethod
			shippingZone   model.ShippingZone
			postalCodeRule model.ShippingMethodPostalCodeRule
			scanFields     = s.ScanFields(&shippingMethod)
		)
		scanFields = append(scanFields, s.ShippingZone().ScanFields(&shippingZone)...)
		scanFields = append(scanFields, s.ShippingMethodPostalCodeRule().ScanFields(&postalCodeRule))

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		if _, exist := shippingMethodMeetMap[shippingMethod.Id]; !exist {
			shippingMethodMeetMap[shippingMethod.Id] = &shippingMethod
		}
		shippingMethodMeetMap[shippingMethod.Id].AppendShippingMethodPostalCodeRule(&postalCodeRule)
		shippingMethodMeetMap[shippingMethod.Id].SetShippingZone(&shippingZone)
	}

	return lo.Values(shippingMethodMeetMap), nil
}

func (ss *SqlShippingMethodStore) commonQueryBuilder(options model_helper.ShippingMethodFilterOption) []qm.QueryMod {
	conds := options.Conditions
	if options.ShippingZoneChannelSlug != nil || options.ShippingZoneCountries != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZones, model.ShippingZoneTableColumns.ID, model.ShippingMethodTableColumns.ShippingZoneID)),
		)

		if options.ShippingZoneCountries != nil {
			conds = append(conds, options.ShippingZoneCountries)
		}

		if options.ShippingZoneChannelSlug != nil {
			conds = append(
				conds,
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZoneChannels, model.ShippingZoneChannelTableColumns.ShippingZoneID, model.ShippingZoneTableColumns.ID)),
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.ShippingZoneChannelTableColumns.ChannelID)),
				options.ShippingZoneChannelSlug,
			)
		}
	}

	for _, load := range options.Load {
		conds = append(conds, qm.Load(load))
	}

	return conds
}

func (ss *SqlShippingMethodStore) FilterByOptions(options model_helper.ShippingMethodFilterOption) (model.ShippingMethodSlice, error) {
	conds := ss.commonQueryBuilder(options)
	return model.ShippingMethods(conds...).All(ss.GetReplica())
}

func (s *SqlShippingMethodStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	_, err := model.ShippingMethods(model.ShippingMethodWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
