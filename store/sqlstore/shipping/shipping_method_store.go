package shipping

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlShippingMethodStore struct {
	store.Store
}

func NewSqlShippingMethodStore(s store.Store) store.ShippingMethodStore {
	return &SqlShippingMethodStore{s}
}

func (s *SqlShippingMethodStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Name",
		"Type",
		"ShippingZoneID",
		"MinimumOrderWeight",
		"MaximumOrderWeight",
		"WeightUnit",
		"MaximumDeliveryDays",
		"MinimumDeliveryDays",
		"Description",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (s *SqlShippingMethodStore) ScanFields(shippingMethod *model.ShippingMethod) []interface{} {
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

// Upsert bases on given method's Id to decide update or insert it
func (s *SqlShippingMethodStore) Upsert(transaction store_iface.SqlxTxExecutor, method *model.ShippingMethod) (*model.ShippingMethod, error) {
	isSaving := false
	if method.Id == "" {
		isSaving = true
		method.PreSave()
	}

	if err := method.IsValid(); err != nil {
		return nil, err
	}

	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	var (
		result sql.Result
		err    error
	)

	if isSaving {
		query := "INSERT INTO " + store.ShippingMethodTableName + "(" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
		result, err = runner.NamedExec(query, method)

	} else {
		query := "UPDATE " + store.ShippingMethodTableName + " SET " +
			s.ModelFields("").
				Map(func(_ int, item string) string {
					return item + ":=" + item
				}).
				Join(",") + " WHERE Id:=Id"

		result, err = runner.NamedExec(query, method)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert shipping method with id=%s", method.Id)
	}
	numUpserted, _ := result.RowsAffected()
	if numUpserted != 1 {
		return nil, errors.Errorf("%d shipping method(s) upserted instead of 1", numUpserted)
	}

	return method, nil
}

// Get finds and returns a shipping method with given id
func (s *SqlShippingMethodStore) Get(methodID string) (*model.ShippingMethod, error) {
	return s.GetbyOption(&model.ShippingMethodFilterOption{
		Id: squirrel.Eq{store.ShippingMethodTableName + ".Id": methodID},
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
	selectFields := append(s.ModelFields(store.ShippingMethodTableName+"."), s.ShippingZone().ModelFields(store.ShippingZoneTableName+".")...)
	selectFields = append(selectFields, s.ShippingMethodPostalCodeRule().ModelFields(store.ShippingMethodPostalCodeRuleTableName+".")...)

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
					ShippingMethodExcludedProducts.ProductID IN :ExcludedProductIDs
					AND ShippingMethodExcludedProducts.ShippingMethodID = ShippingMethods.Id
				)
				LIMIT 1
			)
		)`
		// update params also
		params["ExcludedProductIDs"] = productIDs
	}

	query := `SELECT ` + selectFields.Join(",") + `,
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
				AND ShippingZones.Countries::text LIKE :CountryCode ` + forExcludedProductQuery + `
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
							AND ShippingZones.Countries::text LIKE :CountryCode
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
				AND ShippingZones.Countries::text LIKE :CountryCode ` + forExcludedProductQuery + `
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

	// use Select() here since it can inteprets map[string]interface{} value mapping
	// Query() cannot understand.
	rows, err := s.GetReplicaX().NamedQuery(query, params)
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
	selectFields := ss.ModelFields(store.ShippingMethodTableName + ".")
	if options.SelectRelatedShippingZone {
		selectFields = append(selectFields, ss.ShippingZone().ModelFields(store.ShippingZoneTableName+".")...)
	}

	query := ss.GetQueryBuilder().
		Select(selectFields...).
		From(store.ShippingMethodTableName)

	for _, opt := range []squirrel.Sqlizer{
		options.Id,
		options.Type,
		options.MinimumOrderWeight,
		options.MaximumOrderWeight,
		options.ShippingZoneID,
		options.ShippingZoneChannelSlug,
		options.ShippingZoneCountries,
		options.ChannelListingsChannelSlug,
	} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	if options.ShippingZoneChannelSlug != nil ||
		options.ShippingZoneCountries != nil ||
		options.SelectRelatedShippingZone {
		query = query.InnerJoin(store.ShippingZoneTableName + " ON ShippingZones.Id = ShippingMethods.ShippingZoneID")

		if options.ShippingZoneChannelSlug != nil {
			query = query.
				InnerJoin(store.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ShippingZoneID = ShippingZones.Id").
				InnerJoin(store.ChannelTableName + " ON Channels.Id = ShippingZoneChannels.ChannelID")
		}
	}
	if options.ChannelListingsChannelSlug != nil {
		query = query.
			InnerJoin(store.ShippingMethodChannelListingTableName + " ON ShippingMethodChannelListings.ShippingMethodID = ShippingMethods.Id").
			InnerJoin(store.ChannelTableName + " ON Channels.Id = ShippingMethodChannelListings.ChannelID")
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

	err = ss.GetReplicaX().QueryRowX(queryString, args...).Scan(scanFields...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShippingMethodTableName, "options")
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

	rows, err := ss.GetReplicaX().QueryX(queryString, args...)
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

func (s *SqlShippingMethodStore) Delete(transaction store_iface.SqlxTxExecutor, ids ...string) error {
	query := "DELETE FROM " + store.ShippingMethodTableName + " WHERE Id IN (" + squirrel.Placeholders(len(ids)) + ")"
	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	result, err := runner.Exec(query, lo.Map(ids, func(id string, _ int) any { return id })...)
	if err != nil {
		return errors.Wrap(err, "failed to delete shipping methods")
	}
	rowsAfft, _ := result.RowsAffected()
	if int(rowsAfft) != len(ids) {
		return errors.Errorf("%d shipping methods deleted instead of %d", rowsAfft, len(ids))
	}

	return nil
}
