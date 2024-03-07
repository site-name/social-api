package product

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/mattermost/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (ps *SqlProductStore) filterCategories(query squirrel.SelectBuilder, categoryIDs []string) squirrel.SelectBuilder {
	if len(categoryIDs) == 0 {
		return query
	}

	return query.Where(squirrel.Eq{model.ProductTableColumns.CategoryID: categoryIDs})
}

func (ps *SqlProductStore) filterCollections(query squirrel.SelectBuilder, collectionIDs []string) squirrel.SelectBuilder {
	if len(collectionIDs) == 0 {
		return query
	}

	condition := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) as "a"`).
		Prefix("EXISTS (").
		Suffix(")").
		From(model.TableNames.ProductCollections).
		Where(squirrel.Eq{model.ProductCollectionTableColumns.CollectionID: collectionIDs}).
		Where(fmt.Sprintf("%s = %s", model.ProductCollectionTableColumns.ProductID, model.ProductTableColumns.ID)).
		Limit(1)

	return query.Where(condition)
}

func (ps *SqlProductStore) filterIsPublished(query squirrel.SelectBuilder, isPublished bool, channelIdOrSlug string) squirrel.SelectBuilder {
	return query.Where(
		fmt.Sprintf(
			`EXISTS (
			SELECT
				(1) AS "a"
			FROM 
				%[1]s
			WHERE
				(
					EXISTS (
						SELECT 
							(1) AS "a"
						FROM
							%[2]s
						WHERE
							(
								%[3]s = %[4]s
								AND (%[4]s = ? OR %[5]s = ?)
							)
						LIMIT 1
					)
					AND %[6]s = ?
					AND %[7]s = %[8]s
				)
			LIMIT 1
			)
			AND EXISTS (
				SELECT
					(1) AS "a"
				FROM
					%[9]s
				WHERE
					(
						EXISTS (
							SELECT
								(1) AS "a"
							FROM
								%[10]s PVCL
							WHERE
								(
									EXISTS (
										SELECT
											(1) AS "a"
										FROM
											%[2]s
										WHERE
											(
												(%[4]s = ? OR %[5]s = ?)
												AND %[4]s = %[11]s
											)
										LIMIT 1
									)
									AND %[12]s IS NOT NULL
									AND %[13]s = %[14]s
								)
							LIMIT 1
						)
						AND %[15]s = %[8]s
					)
				LIMIT 1
			)`,

			model.TableNames.ProductChannelListings,                    // 1
			model.TableNames.Channels,                                  // 2
			model.ProductChannelListingTableColumns.ChannelID,          // 3
			model.ChannelTableColumns.ID,                               // 4
			model.ChannelTableColumns.Slug,                             // 5
			model.ProductChannelListingTableColumns.IsPublished,        // 6
			model.ProductChannelListingTableColumns.ProductID,          // 7
			model.ProductTableColumns.ID,                               // 8
			model.TableNames.ProductVariants,                           // 9
			model.TableNames.ProductVariantChannelListings,             // 10
			model.ProductVariantChannelListingTableColumns.ChannelID,   // 11
			model.ProductVariantChannelListingTableColumns.PriceAmount, // 12
			model.ProductVariantChannelListingTableColumns.VariantID,   // 13
			model.ProductVariantTableColumns.ID,                        // 14
			model.ProductVariantTableColumns.ProductID,                 // 15
		),
		channelIdOrSlug,
		channelIdOrSlug,
		isPublished,
		channelIdOrSlug,
		channelIdOrSlug,
	)
}

func (ps *SqlProductStore) filterVariantPrice(
	query squirrel.SelectBuilder,
	priceRange struct {
		Gte *float64
		Lte *float64
	}, channelIdOrSlug string,
) squirrel.SelectBuilder {
	channelQuery := ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.Channels).
		Where(squirrel.And{
			squirrel.Or{
				squirrel.Eq{model.ChannelTableColumns.ID: channelIdOrSlug},
				squirrel.Eq{model.ChannelTableColumns.Slug: channelIdOrSlug},
			},
			squirrel.Eq{model.ProductVariantChannelListingTableColumns.ChannelID: model.ChannelTableColumns.ID},
		}).
		Limit(1).
		Suffix(")")

	productVariantChannelListingQuery := ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.ProductVariantChannelListings).
		Where(channelQuery).
		Where(squirrel.Eq{
			model.ProductVariantChannelListingTableColumns.VariantID: model.ProductVariantTableColumns.ID,
		}).
		Limit(1).
		Suffix(")")

	if priceRange.Lte != nil {
		productVariantChannelListingQuery = productVariantChannelListingQuery.
			Where(squirrel.Or{
				squirrel.Eq{model.ProductVariantChannelListingTableColumns.PriceAmount: nil},
				squirrel.LtOrEq{model.ProductVariantChannelListingTableColumns.PriceAmount: *priceRange.Lte},
			})
	}
	if priceRange.Gte != nil {
		productVariantChannelListingQuery = productVariantChannelListingQuery.
			Where(squirrel.Or{
				squirrel.Eq{model.ProductVariantChannelListingTableColumns.PriceAmount: nil},
				squirrel.GtOrEq{model.ProductVariantChannelListingTableColumns.PriceAmount: *priceRange.Gte},
			})
	}

	productVariantQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		From(model.TableNames.ProductVariants).
		Prefix("EXISTS (").
		Where(productVariantChannelListingQuery).
		Where(squirrel.Eq{
			model.ProductVariantTableColumns.ProductID: model.ProductTableColumns.ID,
		}).
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
	channelIdOrSlug string,
) squirrel.SelectBuilder {
	channelQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.Channels).
		Where(squirrel.And{
			squirrel.Eq{model.ChannelTableColumns.ID: model.ProductChannelListingTableColumns.ChannelID},
			squirrel.Or{
				squirrel.Eq{model.ChannelTableColumns.ID: channelIdOrSlug},
				squirrel.Eq{model.ChannelTableColumns.Slug: channelIdOrSlug},
			},
		}).
		Limit(1).
		Suffix(")")

	productChannelListingQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.ProductChannelListings).
		Where(channelQuery).
		Where(squirrel.Eq{
			model.ProductChannelListingTableColumns.ProductID: model.ProductTableColumns.ID,
		}).
		Limit(1).
		Suffix(")")

	if priceRange.Lte != nil {
		productChannelListingQuery = productChannelListingQuery.
			Where(squirrel.Or{
				squirrel.Eq{model.ProductChannelListingTableColumns.DiscountedPriceAmount: nil},
				squirrel.LtOrEq{model.ProductChannelListingTableColumns.DiscountedPriceAmount: *priceRange.Lte},
			})
	}
	if priceRange.Gte != nil {
		productChannelListingQuery = productChannelListingQuery.
			Where(squirrel.Or{
				squirrel.Eq{model.ProductChannelListingTableColumns.DiscountedPriceAmount: nil},
				squirrel.GtOrEq{model.ProductChannelListingTableColumns.DiscountedPriceAmount: *priceRange.Gte},
			})
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
	m.mu.Lock()
	m.m[key] = append(m.m[key], value...)
	m.mu.Unlock()
}

var wg sync.WaitGroup
var errorSyncGuard sync.Mutex

func (ps *SqlProductStore) filterAttributes(
	query squirrel.SelectBuilder,
	attributes []*model_helper.AttributeFilter,
) squirrel.SelectBuilder {
	// filter out nil values
	nonNilAttributes := lo.Filter(attributes, func(v *model_helper.AttributeFilter, _ int) bool { return v != nil })

	length := len(nonNilAttributes)
	if length == 0 {
		return query
	}

	var (
		value_list           = make(valueList, 0, length)
		boolean_list         = make(booleanList, 0, length)
		value_range_list     = make(valueRangeList, 0, length)
		date_range_list      = make(timeRangeList, 0, length)
		date_time_range_list = make(timeRangeList, 0, length)
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
	var interError error

	syncSetErr := func(err error) {
		errorSyncGuard.Lock()
		if err != nil && interError == nil {
			interError = err
		}
		errorSyncGuard.Unlock()
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
			From(model.TableNames.AssignedProductAttributeValues).
			Where(squirrel.Eq{model.AssignedProductAttributeValueTableColumns.ValueID: values}).
			Where(squirrel.Eq{model.AssignedProductAttributeValueTableColumns.AssignmentID: model.AssignedProductAttributeTableColumns.ID}).
			Limit(1).
			Suffix(")")

		assignedProductAttributes := ps.GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			From(model.TableNames.AssignedProductAttributes).
			Where(assignedProductAttributeValues).
			Where(squirrel.Eq{model.AssignedProductAttributeTableColumns.ProductID: model.ProductTableColumns.ID}).
			Limit(1).
			Suffix(")")

		//

		// TODO: consider the following code.
		// Currently we only have custom variant attribute value and custom product attribute

		// assignedVariantAttributeValues := ps.GetQueryBuilder(squirrel.Question).
		// 	Select(`(1) AS "a"`).
		// 	From(model.TableNames.AssignedVariantAttributeValues).
		// 	Prefix("EXISTS (").
		// 	Where(squirrel.Eq{"AssignedVariantAttributeValues.ValueID": values}).
		// 	Where("AssignedVariantAttributeValues.AssignmentID = AssignedVariantAttributes.Id").
		// 	Limit(1).
		// 	Suffix(")")

		// assignedVariantAttributes := ps.GetQueryBuilder(squirrel.Question).
		// 	Select(`(1) AS "a"`).
		// 	From(model.AssignedVariantAttributeTableName).
		// 	Prefix("EXISTS (").
		// 	Where(assignedVariantAttributeValues).
		// 	Where("AssignedVariantAttributes.VariantID = ProductVariants.Id").
		// 	Limit(1).
		// 	Suffix(")")

		// productVariants := ps.GetQueryBuilder(squirrel.Question).
		// 	Select(`(1) AS "a"`).
		// 	From(model.ProductVariantTableName).
		// 	Prefix("EXISTS (").
		// 	Where(assignedVariantAttributes).
		// 	Where("ProductVariants.ProductID = Products.Id").
		// 	Limit(1).
		// 	Suffix(")")

		orExpr = append(orExpr, assignedProductAttributes)
		query = query.Where(orExpr)
	}

	return query
}

func (ps *SqlProductStore) cleanProductAttributesFilterInput(filterValue valueList, queries *safeMap) error {
	attributes, err := ps.Attribute().FilterbyOption(model_helper.AttributeFilterOption{})
	if err != nil {
		return err
	}
	var (
		attributesSlugPkMap = map[string]string{}            // keys are attribute slugs, values are attribute ids
		attributesPkSlugMap = map[string]string{}            // keys are attribute ids, values are attribute slugs
		valuesMap           = map[string]map[string]string{} // keys are attribute slugs, values are maps with keys are attribute value ids and values are attribute value slugs
		attributeIds        = make([]string, len(attributes))
	)

	for idx, attr := range attributes {
		attributesSlugPkMap[attr.Slug] = attr.ID
		attributesPkSlugMap[attr.ID] = attr.Slug
		attributeIds[idx] = attr.ID
	}

	attributeValues, err := model.AttributeValues(
		model.AttributeValueWhere.AttributeID.IN(attributeIds),
	).All(ps.GetReplica())
	if err != nil {
		return err
	}

	for _, attrValue := range attributeValues {
		if valuesMap[attributesPkSlugMap[attrValue.AttributeID]] == nil {
			valuesMap[attributesPkSlugMap[attrValue.AttributeID]] = map[string]string{}
		}
		valuesMap[attributesPkSlugMap[attrValue.AttributeID]][attrValue.ID] = attrValue.Slug
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
	attributeValues, err := ps.AttributeValue().FilterByOptions(model_helper.AttributeValueFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			qm.Where(
				fmt.Sprintf(
					`EXISTS (
						SELECT (1) AS "a"
						FROM %s
						WHERE %s = ?
						AND %s = %s
						LIMIT 1
					)`,
					model.TableNames.Attributes,
					model.AttributeTableColumns.InputType,
					model.AttributeTableColumns.ID,
					model.AttributeValueTableColumns.AttributeID,
				),
				model.AttributeInputTypeNumeric,
			),
		),
		Preloads: []string{
			model.AttributeValueRels.Attribute,
		},
	})
	if err != nil {
		return err
	}

	var (
		// attributesMap has keys are attribute slugs, values are attribute ids
		attributesMap = model_helper.StringMap{}
		valuesMap     = map[string]map[float64]string{}
	)
	for _, attrValue := range attributeValues {
		attributesMap[attrValue.R.Attribute.Slug] = attrValue.AttributeID

		// we can parse strings into float64 here since:
		// all found attribute values have parent attributes's input type is 'numeric'
		numericName, err := strconv.ParseFloat(attrValue.Name, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse attribute value's name to float64")
		}
		if valuesMap[attrValue.R.Attribute.Slug] == nil {
			valuesMap[attrValue.R.Attribute.Slug] = map[float64]string{}
		}
		valuesMap[attrValue.R.Attribute.Slug][numericName] = attrValue.ID
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
	attributes, err := ps.Attribute().FilterbyOption(model_helper.AttributeFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.AttributeWhere.Slug.IN(filterRange.Slugs()),
		),
		Preload: []string{
			model.AttributeRels.AttributeValues,
		},
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

		for _, attrValue := range attr.R.AttributeValues {
			values[attrValue.Datetime.Time] = attrValue.ID
		}

		valuesMap[attr.Slug] = aMap{attr.ID, values}
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
	attributes, err := ps.Attribute().FilterbyOption(model_helper.AttributeFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.AttributeWhere.Slug.IN(filterValue.Slugs()),
			model.AttributeWhere.InputType.EQ(model.AttributeInputTypeBoolean),
		),
		Preload: []string{model.AttributeRels.AttributeValues},
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

		for _, attrValue := range attr.R.AttributeValues {
			if !attrValue.Boolean.IsNil() {
				values[*attrValue.Boolean.Bool] = attrValue.ID
			}
		}

		valuesMap[attr.Slug] = aMap{attr.ID, values}
	}

	for _, item := range filterValue {
		attrPK := valuesMap[item.Slug].pk

		if valuePK, ok := valuesMap[item.Slug].values[item.Boolean]; ok {
			queries.write(attrPK, []string{valuePK})
		}
	}

	return nil
}

func (ps *SqlProductStore) filterStockAvailability(query squirrel.SelectBuilder, value model_helper.StockAvailability, channelIdOrSlug string) squirrel.SelectBuilder {
	var prefix string

	switch value {
	case model_helper.StockAvailabilityInStock:
		prefix = "EXISTS ("
	case model_helper.StockAvailabilityOutOfStock:
		prefix = "NOT EXISTS ("

	default:
		return query
	}

	channelQuery, _, _ := ps.Stock().FilterForChannel(model_helper.StockFilterForChannelOption{
		ChannelID:       channelIdOrSlug,
		ReturnQueryOnly: true,
	})

	productVariantIDsQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(model.StockTableColumns.ProductVariantID).
		Prefix(model.ProductVariantTableColumns.ID + " IN (").
		From(model.TableNames.Stocks).
		Where(channelQuery).
		Where(fmt.Sprintf(
			`%[1]s > COALESCE (
				(
					SELECT
						SUM (%[2]s)
					FROM
						%[3]s
					WHERE
						%[2]s > 0
						AND %[4]s = %[5]s
					GROUP BY
						%[4]s
				),
				0
			)`,
			model.StockTableColumns.Quantity,               // 1
			model.AllocationTableColumns.QuantityAllocated, // 2
			model.TableNames.Allocations,                   // 3
			model.AllocationTableColumns.StockID,           // 4
			model.StockTableColumns.ID,                     // 5
		)).
		Suffix(")")

	productVariantSelect := ps.GetQueryBuilder(squirrel.Question).
		Select(model.ProductVariantTableColumns.ProductID).
		Prefix(prefix).
		From(model.TableNames.ProductVariants).
		Where(productVariantIDsQuery).
		Where(squirrel.Eq{model.ProductVariantTableColumns.ProductID: model.ProductTableColumns.ID}).
		Limit(1).
		Suffix(")")

	return query.Where(productVariantSelect)
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
			InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariants, model.ProductVariantTableColumns.ProductID, model.ProductTableColumns.ID)).
			InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Stocks, model.StockTableColumns.ProductVariantID, model.ProductVariantTableColumns.ID)).
			Where(squirrel.Eq{model.StockTableColumns.WarehouseID: value.WarehouseIds}).
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
	productVariantQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(model.ProductVariantTableColumns.ID).
		From(model.TableNames.ProductVariants).
		Where(squirrel.Eq{model.ProductVariantTableColumns.ProductID: model.ProductTableColumns.ID})

	if len(warehouseIDs) > 0 {
		warehouseIDs := *(*[]any)(unsafe.Pointer(&warehouseIDs))

		productVariantQuery = productVariantQuery.
			Column(
				fmt.Sprintf(
					`SUM (%s) FILTER (
						WHERE %s IN (%s)
					) AS TotalQuantity`,
					model.StockTableColumns.Quantity,
					model.StockTableColumns.WarehouseID,
					squirrel.Placeholders(len(warehouseIDs)),
				),
				warehouseIDs...,
			)
	} else {
		productVariantQuery = productVariantQuery.Column(fmt.Sprintf(`SUM (%s) AS TotalQuantity`, model.StockTableColumns.Quantity))
	}

	productVariantQuery = productVariantQuery.
		LeftJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Stocks, model.StockTableColumns.ProductVariantID, model.ProductVariantTableColumns.ID)).
		GroupBy(model.ProductVariantTableColumns.ID)

	// parse quantity range
	if quantity.Gte != nil {
		productVariantQuery = productVariantQuery.Where("TotalQuantity >= ?", *quantity.Gte)
	}
	if quantity.Lte != nil {
		productVariantQuery = productVariantQuery.Where("TotalQuantity <= ?", *quantity.Lte)
	}

	return query.
		InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariants, model.ProductVariantTableColumns.ProductID, model.ProductTableColumns.ID)).
		Where(squirrel.Expr(model.ProductVariantTableColumns.ID+" IN ?", productVariantQuery))
}

func (ps *SqlProductStore) filterGiftCard(query squirrel.SelectBuilder, value bool) squirrel.SelectBuilder {
	// prefix := "EXISTS ("
	// if !value {
	// 	prefix = "NOT EXISTS ("
	// }
	// productTypeFilter := ps.GetQueryBuilder().
	// 	Select(`(1) AS "a"`).
	// 	Prefix(prefix).
	// 	From(model.ProductTypeTableName).
	// 	Where("ProductTypes.Kind = ?", model.GIFT_CARD).
	// 	Where("ProductTypes.Id = Products.ProductTypeID").
	// 	Limit(1).
	// 	Suffix(")")

	// return query.Where(productTypeFilter)

	// TODO: investigate this
	return query
}

func (ps *SqlProductStore) filterProductIDs(query squirrel.SelectBuilder, productIDs []string) squirrel.SelectBuilder {
	if len(productIDs) == 0 {
		return query
	}

	return query.Where(squirrel.Eq{model.ProductTableColumns.ID: productIDs})
}

func (ps *SqlProductStore) filterHasPreorderedVariants(query squirrel.SelectBuilder, value bool) squirrel.SelectBuilder {
	prefix := "EXISTS ("
	if !value {
		prefix = "NOT EXISTS ("
	}
	variantQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix(prefix).
		From(model.TableNames.ProductVariants).
		Where(
			fmt.Sprintf(
				`%[1]s
				AND (
					%[2]s IS NULL
					OR %[2]s > ?
				)
				AND %[3]s = %[4]s`,
				model.ProductVariantTableColumns.IsPreorder,      // 1
				model.ProductVariantTableColumns.PreorderEndDate, // 2
				model.ProductVariantTableColumns.ProductID,       // 3
				model.ProductTableColumns.ID,                     // 4
			),
			model_helper.GetMillis(),
		).
		Limit(1).
		Suffix(")")

	return query.Where(variantQuery)
}

// TODO: add search by search vector.
func (ps *SqlProductStore) filterSearch(query squirrel.SelectBuilder, value string) squirrel.SelectBuilder {
	variantQuery := ps.GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.ProductVariants).
		Where(fmt.Sprintf(`%s = ? AND %s = %s`, model.ProductVariantTableColumns.Sku, model.ProductVariantTableColumns.ProductID, model.ProductTableColumns.ID), value).
		Suffix(")")

	return query.Where(variantQuery)
}
