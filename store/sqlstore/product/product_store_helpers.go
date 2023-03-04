package product

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

func (ps *SqlProductStore) filterCategories(query squirrel.SelectBuilder, categoryIDs []string) squirrel.SelectBuilder {
	if len(categoryIDs) == 0 {
		return query
	}

	panic("not implemented")
}

func (ps *SqlProductStore) filterCollections(query squirrel.SelectBuilder, collectionIDs []string) squirrel.SelectBuilder {
	if len(collectionIDs) == 0 {
		return query
	}

	condition := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) as "a"`).
		Prefix("EXISTS (").
		Suffix(")").
		From(store.CollectionProductRelationTableName).
		Where(squirrel.Eq{"ProductCollections.CollectionID": collectionIDs}).
		Where("ProductCollections.ProductID = Products.Id").
		Limit(1)

	return query.Where(condition)
}

func (ps *SqlProductStore) filterIsPublished(query squirrel.SelectBuilder, isPublished bool, channelID string) squirrel.SelectBuilder {
	return query.Where(`
		EXISTS (
			SELECT
				(1) AS "a"
			FROM 
				`+store.ProductChannelListingTableName+`
			WHERE
				(
					EXISTS (
						SELECT 
							(1) AS "a"
						FROM
							`+store.ChannelTableName+`
						WHERE
							(
								ProductChannelListings.ChannelID = Channels.Id AND
								Channels.Id = ?
							)
						LIMIT 1
					)
					AND ProductChannelListings.IsPublished = ?
					AND ProductChannelListings.ProductID = Products.Id
				)
			LIMIT 1
		)
		AND EXISTS (
			SELECT
				(1) AS "a"
			FROM
				`+store.ProductVariantTableName+`
			WHERE
				(
					EXISTS (
						SELECT
							(1) AS "a"
						FROM
							`+store.ProductVariantChannelListingTableName+`
						WHERE
							(
								EXISTS (
									SELECT
										(1) AS "a"
									FROM
										`+store.ChannelTableName+`
									WHERE
										(
											Channels.Id = ?
											AND Channels.Id = ProductVariantChannelListings.ChannelID
										)
									LIMIT 1
								)
								AND ProductVariantChannelListings.PriceAmount IS NOT NULL
								AND ProductVariantChannelListings.VariantID = ProductVariants.Id
							)
						LIMIT 1
					)
					AND ProductVariants.ProductID = Products.Id
				)
			LIMIT 1
		)`,
		channelID,
		isPublished,
		channelID,
	)
}

func (ps *SqlProductStore) filterVariantPrice(
	query squirrel.SelectBuilder,
	priceRange struct {
		Gte *float64
		Lte *float64
	}, channelID string,
) squirrel.SelectBuilder {
	channelQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ChannelTableName).
		Where("Channels.Id = ? AND Channels.Id = ProductVariantChannelListings.ChannelID", channelID).
		Limit(1).
		Suffix(")")

	productVariantChannelListingQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ProductVariantChannelListingTableName).
		Where(channelQuery).
		Where("ProductVariantChannelListings.VariantID = ProductVariants.Id").
		Limit(1).
		Suffix(")")

	if priceRange.Lte != nil {
		productVariantChannelListingQuery = productVariantChannelListingQuery.
			Where("ProductVariantChannelListings.PriceAmount <= ? OR ProductVariantChannelListings.PriceAmount IS NULL", *priceRange.Gte)
	}
	if priceRange.Gte != nil {
		productVariantChannelListingQuery = productVariantChannelListingQuery.
			Where("ProductVariantChannelListings.PriceAmount >= ? OR ProductVariantChannelListings.PriceAmount IS NULL", *priceRange.Lte)
	}

	productVariantQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		From(store.ProductVariantTableName).
		Prefix("EXISTS (").
		Where(productVariantChannelListingQuery).
		Where("ProductVariants.ProductID = Products.Id").
		Limit(1).
		Suffix(")")

	return query.Where(productVariantQuery)
}

func (ps *SqlProductStore) filterMinimalPrice(
	query squirrel.SelectBuilder,
	priceRange struct {
		Gte *float64
		Lte *float64
	},
	channelID string,
) squirrel.SelectBuilder {
	channelQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ChannelTableName).
		Where("Channels.Id = ? AND Channels.Id = ProductChannelListings.ChannelID", channelID).
		Limit(1).
		Suffix(")")

	productChannelListingQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ProductChannelListingTableName).
		Where(channelQuery).
		Where("ProductChannelListings.ProductID = Products.Id").
		Limit(1).
		Suffix(")")

	if priceRange.Lte != nil {
		productChannelListingQuery = productChannelListingQuery.
			Where("ProductChannelListings.DiscountedPriceAmount IS NULL OR ProductChannelListings.DiscountedPriceAmount <= ?", *priceRange.Lte)
	}
	if priceRange.Gte != nil {
		productChannelListingQuery = productChannelListingQuery.
			Where("ProductChannelListings.DiscountedPriceAmount IS NULL OR ProductChannelListings.DiscountedPriceAmount >= ?", *priceRange.Gte)
	}

	return query.Where(productChannelListingQuery)
}

type (
	value struct {
		Slug   string
		Values []string
	}
	booleanRange struct {
		Slug    string
		Boolean bool
	}
	valueRange struct {
		Slug  string
		Range struct {
			Gte *int32
			Lte *int32
		}
	}
	timeRange struct {
		Slug string
		Date struct {
			Gte *time.Time
			Lte *time.Time
		}
	}

	valueList      []value
	booleanList    []booleanRange
	valueRangeList []valueRange
	timeRangeList  []timeRange
)

func (t timeRangeList) Slugs() []string {
	return lo.Map(t, func(item timeRange, _ int) string { return item.Slug })
}

func (t booleanList) Slugs() []string {
	return lo.Map(t, func(item booleanRange, _ int) string { return item.Slug })
}

type safeMap struct {
	mu sync.Mutex
	m  map[string][]string
}

func (m *safeMap) write(key string, value []string) {
	if m.mu.TryLock() {
		m.m[key] = append(m.m[key], value...)
		m.mu.Unlock()
	}
}

func (ps *SqlProductStore) filterAttributes(
	query squirrel.SelectBuilder,
	attributes []*model.AttributeFilter,
) squirrel.SelectBuilder {

	// filter out nil values
	nonNilAttributes := lo.Filter(attributes, func(v *model.AttributeFilter, _ int) bool { return v != nil })

	if len(nonNilAttributes) == 0 {
		return query
	}

	var (
		value_list           valueList
		boolean_list         booleanList
		value_range_list     valueRangeList
		date_range_list      timeRangeList
		date_time_range_list timeRangeList
	)

	for _, input := range nonNilAttributes {
		if len(input.Values) > 0 {
			value_list = append(value_list, value{input.Slug, input.Values})
		} else if input.ValuesRange != nil {
			value_range_list = append(value_range_list, valueRange{input.Slug, *input.ValuesRange})
		} else if input.Date != nil {
			date_range_list = append(date_range_list, timeRange{input.Slug, *input.Date})
		} else if input.DateTime != nil {
			date_time_range_list = append(date_time_range_list, timeRange{input.Slug, *input.Date})
		} else if input.Boolean != nil {
			boolean_list = append(boolean_list, booleanRange{input.Slug, *input.Boolean})
		}
	}

	var queries = &safeMap{
		m: map[string][]string{},
	}
	var wg sync.WaitGroup
	var interError error
	var errorSyncGuard sync.Mutex

	syncSetErr := func(err error) {
		if errorSyncGuard.TryLock() {
			if err != nil && interError == nil {
				interError = err
			}
			errorSyncGuard.Unlock()
		}
	}

	if len(value_list) > 0 {
		wg.Add(1)

		go func() {
			err := ps.cleanProductAttributesFilterInput(value_list, queries)
			syncSetErr(err)
			wg.Done()
		}()
	}

	if len(value_range_list) > 0 {
		wg.Add(1)

		go func() {
			err := ps.cleanProductAttributesRangeFilterInput(value_range_list, queries)
			syncSetErr(err)
			wg.Done()
		}()
	}

	if len(date_range_list) > 0 {
		wg.Add(1)

		go func() {
			err := ps.cleanProductAttributesDateTimeRangeFilterInput(date_range_list, queries, true)
			syncSetErr(err)
			wg.Done()
		}()
	}

	if len(date_time_range_list) > 0 {
		wg.Add(1)

		go func() {
			err := ps.cleanProductAttributesDateTimeRangeFilterInput(date_time_range_list, queries, false)
			syncSetErr(err)
			wg.Done()
		}()
	}

	if len(boolean_list) > 0 {
		wg.Add(1)

		go func() {
			err := ps.cleanProductAttributesBooleanFilterInput(boolean_list, queries)
			syncSetErr(err)
			wg.Done()
		}()
	}

	wg.Wait()

	if interError != nil {
		slog.Error("Filter product attributes error", slog.Err(interError))
		return query
	}

	return ps.filterProductsByAttributesValues(query, queries)
}

func (ps *SqlProductStore) filterProductsByAttributesValues(query squirrel.SelectBuilder, queries *safeMap) squirrel.SelectBuilder {
	for _, values := range queries.m {
		orExpr := squirrel.Or{}

		assignedProductAttributeValues := ps.GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			From(store.AssignedProductAttributeValueTableName).
			Where(squirrel.Eq{"AssignedProductAttributeValues.ValueID": values}).
			Where("AssignedProductAttributeValues.AssignmentID = AssignedProductAttributes.Id").
			Limit(1).
			Suffix(")")

		assignedProductAttributes := ps.GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			From(store.AssignedProductAttributeTableName).
			Where(assignedProductAttributeValues).
			Where("AssignedProductAttributes.ProductID = Products.Id").
			Limit(1).
			Suffix(")")

		orExpr = append(orExpr, assignedProductAttributes)

		ssignedVariantAttributeValues := ps.GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			From(store.AssignedVariantAttributeValueTableName).
			Prefix("EXISTS (").
			Where(squirrel.Eq{"AssignedVariantAttributeValues.ValueID": values}).
			Where("AssignedVariantAttributeValues.AssignmentID = AssignedVariantAttributes.Id").
			Limit(1).
			Suffix(")")

		assignedVariantAttributes := ps.GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			From(store.AssignedVariantAttributeTableName).
			Prefix("EXISTS (").
			Where(ssignedVariantAttributeValues).
			Where("AssignedVariantAttributes.VariantID = ProductVariants.Id").
			Limit(1).
			Suffix(")")

		productVariants := ps.GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			From(store.ProductVariantTableName).
			Prefix("EXISTS (").
			Where(assignedVariantAttributes).
			Where("ProductVariants.ProductID = Products.Id").
			Limit(1).
			Suffix(")")

		orExpr = append(orExpr, productVariants)
		query = query.Where(orExpr)
	}

	return query
}

func (ps *SqlProductStore) cleanProductAttributesFilterInput(filterValue valueList, queries *safeMap) error {
	attributes, err := ps.Attribute().FilterbyOption(&model.AttributeFilterOption{})
	if err != nil {
		return errors.Wrap(err, "failed to find all attributes")
	}
	var (
		attributesSlugPkMap = map[string]string{}
		attributesPkSlugMap = map[string]string{}
		valuesMap           = map[string]map[string]string{}
	)

	for _, attr := range attributes {
		attributesSlugPkMap[attr.Slug] = attr.Id
		attributesPkSlugMap[attr.Id] = attr.Slug
	}

	var attributeValues model.AttributeValues
	err = ps.GetReplicaX().Select(&attributeValues, "SELECT * FROM "+store.AttributeValueTableName)
	if err != nil {
		return errors.Wrap(err, "failed to find all attribute values")
	}

	for _, attrValue := range attributeValues {
		valuesMap[attributesPkSlugMap[attrValue.AttributeID]][attrValue.Id] = attrValue.Slug
	}

	// Convert attribute:value pairs into a dictionary where
	// attributes are keys and values are grouped in lists
	for _, value := range filterValue {
		attrPk, ok := attributesSlugPkMap[value.Slug]
		if !ok {
			return fmt.Errorf("unknown attribute name: %s", value.Slug)
		}

		attrvaluePk := []string{}

		for _, valueSlug := range value.Values {
			if item, ok := valuesMap[value.Slug]; ok {
				attrvaluePk = append(attrvaluePk, item[valueSlug])
			}
		}

		queries.write(attrPk, attrvaluePk)
	}

	return nil
}

func (ps *SqlProductStore) cleanProductAttributesRangeFilterInput(filterValue valueRangeList, queries *safeMap) error {
	attributeQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.AttributeTableName).
		Where(squirrel.Eq{store.AttributeTableName + ".InputType": model.NUMERIC}).
		Where("Attributes.Id = AttributeValues.AttributeID").
		Suffix(")").
		Limit(1)

	attributeValues, err := ps.AttributeValue().FilterByOptions(model.AttributeValueFilterOptions{
		SelectRelatedAttribute: true,
		Extra:                  attributeQuery,
	})
	if err != nil {
		return err
	}

	var (
		// attributesMap has keys are attribute slugs, values are attribute ids
		attributesMap = model.StringMap{}
		valuesMap     = map[string]map[float64]string{}
	)
	for _, attrValue := range attributeValues {
		attributesMap[attrValue.GetAttribute().Slug] = attrValue.AttributeID

		// we can parse strings into float64 here since:
		// all found attribute values have parent attributes's input type is 'numeric'
		numericName, err := strconv.ParseFloat(attrValue.Name, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse attribute value's name to float64")
		}
		valuesMap[attrValue.GetAttribute().Slug][numericName] = attrValue.Id
	}

	for _, vlRange := range filterValue {
		attrPk, ok := attributesMap[vlRange.Slug]
		if !ok {
			return fmt.Errorf("unknown numeric attribute name: %v", vlRange.Slug)
		}

		var (
			gte float64 = 0
			lte float64 = math.MaxInt64
		)
		if vlRange.Range.Gte != nil {
			gte = float64(*vlRange.Range.Gte)
		}
		if vlRange.Range.Lte != nil {
			lte = float64(*vlRange.Range.Lte)
		}

		attrValues := valuesMap[vlRange.Slug]

		attrValPks := []string{}
		for key, value := range attrValues {
			if gte <= key && key <= lte {
				attrValPks = append(attrValPks, value)
			}
		}

		queries.write(attrPk, attrValPks)
	}

	return nil
}

func (ps *SqlProductStore) cleanProductAttributesDateTimeRangeFilterInput(filterRange timeRangeList, queries *safeMap, isDate bool) error {
	attributes, err := ps.Attribute().FilterbyOption(&model.AttributeFilterOption{
		Slug:                           squirrel.Eq{store.AttributeTableName + ".Slug": filterRange.Slugs()},
		PrefetchRelatedAttributeValues: true,
	})
	if err != nil {
		return err
	}

	type aMap struct {
		pk     string
		values map[*time.Time]string
	}

	var valuesMap = map[string]aMap{}

	for _, attr := range attributes {
		values := map[*time.Time]string{}

		for _, attrValue := range attr.GetAttributeValues() {
			values[attrValue.Datetime] = attrValue.Id
		}

		valuesMap[attr.Slug] = aMap{attr.Id, values}
	}

	for _, item := range filterRange {
		var (
			attrPK           = valuesMap[item.Slug].pk
			gte              = item.Date.Gte
			lte              = item.Date.Lte
			matchingValuesID = []string{}
		)

		for value, pk := range valuesMap[item.Slug].values {
			if value == nil {
				continue
			}

			realValue := *value
			if isDate {
				realValue = util.StartOfDay(realValue)
			}

			if gte != nil && lte != nil {
				if (gte.Equal(realValue) || gte.Before(realValue)) && (lte.After(realValue) || lte.Equal(realValue)) {
					matchingValuesID = append(matchingValuesID, pk)
				}
			} else if gte != nil && (gte.Equal(realValue) || gte.Before(realValue)) {
				matchingValuesID = append(matchingValuesID, pk)
			} else if lte != nil && (lte.After(realValue) || lte.Equal(realValue)) {
				matchingValuesID = append(matchingValuesID, pk)
			}

			queries.write(attrPK, matchingValuesID)
		}
	}

	return nil
}

func (ps *SqlProductStore) cleanProductAttributesBooleanFilterInput(filterValue booleanList, queries *safeMap) error {
	attributes, err := ps.Attribute().FilterbyOption(&model.AttributeFilterOption{
		PrefetchRelatedAttributeValues: true,
		Slug:                           squirrel.Eq{store.AttributeTableName + ".Slug": filterValue.Slugs()},
		InputType:                      squirrel.Eq{store.AttributeTableName + ".InputType": model.BOOLEAN},
	})
	if err != nil {
		return err
	}

	type aMap struct {
		pk     string
		values map[bool]string
	}

	var valuesMap = map[string]aMap{}

	for _, attr := range attributes {
		values := map[bool]string{}

		for _, attrValue := range attr.GetAttributeValues() {
			if attrValue.Boolean != nil {
				values[*attrValue.Boolean] = attrValue.Id
			}
		}

		valuesMap[attr.Slug] = aMap{attr.Id, values}
	}

	for _, item := range filterValue {
		attrPK := valuesMap[item.Slug].pk

		if valuePK, ok := valuesMap[item.Slug].values[item.Boolean]; ok {
			queries.write(attrPK, []string{valuePK})
		}
	}

	return nil
}

func (ps *SqlProductStore) filterStockAvailability(query squirrel.SelectBuilder, value string, channelID interface{}) squirrel.SelectBuilder {
	var (
		validValue = true
		prefix     string
		channelID_ string
	)
	if channelID != nil {
		channelID_ = channelID.(string)
	}

	switch value {
	case model.StockAvailabilityInStock:
		prefix = "EXISTS ("
	case model.StockAvailabilityOutOfStock:
		prefix = "NOT EXISTS ("

	default:
		validValue = false
	}

	if !validValue {
		return query
	}

	channelQuery, _, _ := ps.Stock().FilterForChannel(&model.StockFilterForChannelOption{
		ChannelID:       channelID_,
		ReturnQueryOnly: true,
	})

	productVariantIDsQuery := ps.GetQueryBuilder(squirrel.Question).
		Select("Stocks.ProductVariantID").
		Prefix("ProductVariants.Id IN (").
		From(store.StockTableName).
		Where(channelQuery).
		Where(`Stocks.Quantity > COALESCE (
			(
        SELECT
          SUM( Allocations.QuantityAllocated )
        FROM
					Allocations
        WHERE
            Allocations.QuantityAllocated > 0
            AND Allocations.StockID = Stocks.Id
        GROUP BY
					Allocations.StockID
      ), 0
		)`).
		Suffix(")")

	productVariantSelect := ps.GetQueryBuilder(squirrel.Question).
		Select("ProductVariants.ProductID").
		Prefix(prefix).
		From(store.ProductVariantTableName).
		Where(productVariantIDsQuery).
		Where("ProductVariants.ProductID = Products.Id").
		Limit(1).
		Suffix(")")

	return query.Where(productVariantSelect)
}

func (ps *SqlProductStore) filterProductTypes(query squirrel.SelectBuilder, value []string) squirrel.SelectBuilder {
	if len(value) == 0 {
		return query
	}

	return query.Where(squirrel.Eq{"Products.ProductTypeID": value})
}

func (ps *SqlProductStore) filterStocks(
	query squirrel.SelectBuilder,
	value struct {
		WarehouseIds []string
		Quantity     *struct {
			Gte *int32
			Lte *int32
		}
	},
) squirrel.SelectBuilder {
	if len(value.WarehouseIds) > 0 && value.Quantity == nil {
		return query.
			InnerJoin(store.ProductVariantTableName + " ON ProductVariants.ProductID = Products.Id").
			InnerJoin(store.StockTableName + " ON Stocks.ProductVariantID = ProductVariants.Id").
			Where(squirrel.Eq{"Stocks.WarehouseID": value.WarehouseIds}).
			Distinct()
	}

	if len(value.WarehouseIds) == 0 && value.Quantity != nil {
		return ps.filterQuantity(query, *value.Quantity, nil).Distinct()
	}

	if value.Quantity != nil && len(value.WarehouseIds) != 0 {
		return ps.filterQuantity(query, *value.Quantity, value.WarehouseIds).Distinct()
	}

	return query
}

func (ps *SqlProductStore) filterQuantity(
	query squirrel.SelectBuilder,
	quantity struct {
		Gte *int32
		Lte *int32
	},
	warehouseIDs []string,
) squirrel.SelectBuilder {

	queryString, args, err := query.ToSql()
	if err != nil {
		slog.Error("failed to build query string for products", slog.Err(err))
		return query
	}

	var products model.Products
	err = ps.GetReplicaX().Select(&products, queryString, args...)
	if err != nil {
		slog.Error("failed to find products", slog.Err(err))
		return query
	}

	productVariantQuery := ps.GetQueryBuilder().
		Select(ps.ProductVariant().ModelFields(store.ProductVariantTableName + ".")...).
		From(store.ProductVariantTableName).
		Where(squirrel.Eq{"ProductVariants.ProductID": products.IDs()})

	if len(warehouseIDs) > 0 {
		productVariantQuery = productVariantQuery.
			Column(`SUM (Stocks.Quantity) FILTER (
				WHERE	Stocks.WarehouseID IN (`+squirrel.Placeholders(len(warehouseIDs))+`)
			) AS TotalQuantity`, lo.Map(warehouseIDs, func(item string, _ int) any { return item })...)

	} else {
		productVariantQuery = productVariantQuery.Column(`SUM (Stocks.Quantity) AS TotalQuantity`)
	}

	productVariantQuery = productVariantQuery.
		LeftJoin(store.StockTableName + " ON Stocks.ProductVariantID = ProductVariants.Id").
		GroupBy("ProductVariants.Id")

	// parse quantity range
	if quantity.Gte != nil {
		productVariantQuery = productVariantQuery.Where("TotalQuantity >= ?", *quantity.Gte)
	}
	if quantity.Lte != nil {
		productVariantQuery = productVariantQuery.Where("TotalQuantity <= ?", *quantity.Lte)
	}

	queryString, args, err = productVariantQuery.ToSql()
	if err != nil {
		slog.Error("failed to build query string for product variants", slog.Err(err))
		return query
	}

	var variants model.ProductVariants
	err = ps.GetReplicaX().Select(&variants, queryString, args...)
	if err != nil {
		slog.Error("failed to find product variants", slog.Err(err))
		return query
	}

	return query.
		InnerJoin(store.ProductVariantTableName + " ON ProductVariants.ProductID = Products.Id").
		Where(squirrel.Eq{"ProductVariants.Id": variants.IDs()})
}

func (ps *SqlProductStore) filterGiftCard(query squirrel.SelectBuilder, value bool) squirrel.SelectBuilder {
	prefix := "EXISTS ("
	if !value {
		prefix = "NOT EXISTS ("
	}
	productTypeFilter := ps.GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix(prefix).
		From(store.ProductTypeTableName).
		Where("ProductTypes.Kind = ?", model.GIFT_CARD).
		Where("ProductTypes.Id = Products.ProductTypeID").
		Limit(1).
		Suffix(")")

	return query.Where(productTypeFilter)
}

func (ps *SqlProductStore) filterProductIDs(query squirrel.SelectBuilder, productIDs []string) squirrel.SelectBuilder {
	if len(productIDs) == 0 {
		return query
	}

	return query.Where(squirrel.Eq{"Products.Id": productIDs})
}

func (ps *SqlProductStore) filterHasPreorderedVariants(query squirrel.SelectBuilder, value bool) squirrel.SelectBuilder {
	prefix := "EXISTS ("
	if !value {
		prefix = "NOT EXISTS ("
	}
	variantQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix(prefix).
		From(store.ProductVariantTableName).
		Where(
			`ProductVariants.IsPreOrder
			AND (
				ProductVariants.PreorderEndDate IS NULL 
				OR ProductVariants.PreorderEndDate > ? 
			)
			AND ProductVariants.ProductID = Products.Id`,
			model.GetMillis(),
		).
		Limit(1).
		Suffix(")")

	return query.Where(variantQuery)
}

func (ps *SqlProductStore) filterSearch(query squirrel.SelectBuilder, value string) squirrel.SelectBuilder {
	slog.Warn("this method is not implemented", slog.String("method", "SqlProductStore.filterSearch"))
	return query
}
