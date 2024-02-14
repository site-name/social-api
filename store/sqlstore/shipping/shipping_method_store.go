package shipping

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gorm.io/gorm"
)

type SqlShippingMethodStore struct {
	store.Store
}

func NewSqlShippingMethodStore(s store.Store) store.ShippingMethodStore {
	return &SqlShippingMethodStore{s}
}

// Upsert bases on given method's Id to decide update or insert it
func (s *SqlShippingMethodStore) Upsert(transaction *gorm.DB, method *model.ShippingMethod) (*model.ShippingMethod, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	err := transaction.Save(method).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert shipping method with id=%s", method.Id)
	}

	return method, nil
}

// Get finds and returns a shipping method with given id
func (s *SqlShippingMethodStore) Get(methodID string) (*model.ShippingMethod, error) {
	return s.GetbyOption(&model.ShippingMethodFilterOption{
		Conditions: squirrel.Eq{model.ShippingMethodTableName + ".Id": methodID},
	})
}

// ApplicableShippingMethods finds all shipping method for given checkout
//
// sql queries here are borrowed. Please check the file shipping_method_store.md
func (s *SqlShippingMethodStore) ApplicableShippingMethods(price *goprices.Money, channelID string, weight *measurement.Weight, countryCode model.CountryCode, productIDs []string) ([]*model.ShippingMethod, error) {
	/*
		NOTE: we also prefetch postal_code_rules, shipping zones for later use
		please refer to saleor/shipping/models for details
	*/

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

func (ss *SqlShippingMethodStore) commonQueryBuilder(options *model.ShippingMethodFilterOption) squirrel.SelectBuilder {
	selectFields := []string{model.ShippingMethodTableName + ".*"}
	if options.SelectRelatedShippingZone {
		selectFields = append(selectFields, model.ShippingZoneTableName+".*")
	}

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(model.ShippingMethodTableName).Where(options.Conditions)

	for _, opt := range []squirrel.Sqlizer{
		options.ShippingZoneChannelSlug,
		options.ShippingZoneCountries,
		options.ChannelListingsChannelSlug,
	} {
		query = query.Where(opt)
	}

	if options.ShippingZoneChannelSlug != nil ||
		options.ShippingZoneCountries != nil ||
		options.SelectRelatedShippingZone {
		query = query.InnerJoin(model.ShippingZoneTableName + " ON ShippingZones.Id = ShippingMethods.ShippingZoneID")

		if options.ShippingZoneChannelSlug != nil {
			query = query.
				InnerJoin(model.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id").
				InnerJoin(model.ChannelTableName + " ON Channels.Id = ShippingZoneChannels.ChannelID")
		}
	}
	if options.ChannelListingsChannelSlug != nil {
		query = query.
			InnerJoin(model.ShippingMethodChannelListingTableName + " ON ShippingMethodChannelListings.ShippingMethodID = ShippingMethods.Id").
			InnerJoin(model.ChannelTableName + " ON Channels.Id = ShippingMethodChannelListings.ChannelID")
	}

	return query
}

// GetbyOption finds and returns a shipping method that satisfy given options
func (ss *SqlShippingMethodStore) GetbyOption(options *model.ShippingMethodFilterOption) (*model.ShippingMethod, error) {
	queryString, args, err := ss.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}

	var (
		res          model.ShippingMethod
		shippingZone model.ShippingZone
		scanFields   = ss.ScanFields(&res)
	)
	if options.SelectRelatedShippingZone {
		scanFields = append(scanFields, ss.ShippingZone().ScanFields(&shippingZone)...)
	}

	err = ss.GetReplica().Raw(queryString, args...).Row().Scan(scanFields...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.NewErrNotFound(model.ShippingMethodTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find shipping method by given options")
	}

	if options.SelectRelatedShippingZone {
		res.SetShippingZone(&shippingZone)
	}

	return &res, nil
}

func (ss *SqlShippingMethodStore) FilterByOptions(options *model.ShippingMethodFilterOption) ([]*model.ShippingMethod, error) {
	queryString, args, err := ss.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}

	rows, err := ss.GetReplica().Raw(queryString, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping methods with given options.")
	}
	defer rows.Close()

	var shippingMethodMap = map[string]*model.ShippingMethod{}

	for rows.Next() {
		var (
			method     model.ShippingMethod
			zone       model.ShippingZone
			scanFields = ss.ScanFields(&method)
		)
		if options.SelectRelatedShippingZone {
			scanFields = append(scanFields, ss.ShippingZone().ScanFields(&zone)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row of shipping method")
		}

		if options.SelectRelatedShippingZone {
			method.SetShippingZone(&zone)
		}

		shippingMethodMap[method.Id] = &method
	}

	return lo.Values(shippingMethodMap), nil
}

func (s *SqlShippingMethodStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	_, err := model.ShippingMethods(model.ShippingMethodWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
