package product

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattermost/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlProductStore struct {
	store.Store
}

func NewSqlProductStore(s store.Store) store.ProductStore {
	return &SqlProductStore{s}
}

func (ps *SqlProductStore) Save(tx boil.ContextTransactor, product model.Product) (*model.Product, error) {
	if tx == nil {
		tx = ps.GetMaster()
	}

	isSaving := product.ID == ""
	if isSaving {
		model_helper.ProductPreSave(&product)
	} else {
		model_helper.ProductPreUpdate(&product)
	}

	if err := model_helper.ProductIsValid(product); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = product.Insert(tx, boil.Infer())
	} else {
		_, err = product.Update(tx, boil.Blacklist(model.ProductColumns.CreatedAt))
	}

	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{model.ProductColumns.Name, "products_name_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Products, model.ProductColumns.Name, product.Name)
		}
		if ps.IsUniqueConstraintError(err, []string{model.ProductColumns.Slug, "products_slug_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Products, model.ProductColumns.Slug, product.Slug)
		}
		return nil, err
	}

	return &product, nil
}

func (ps *SqlProductStore) commonQueryBuilder(option model_helper.ProductFilterOption) []qm.QueryMod {
	conds := option.Conditions
	if option.ProductVariantID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariants, model.ProductVariantTableColumns.ProductID, model.ProductTableColumns.ID)),
			option.ProductVariantID,
		)
	} else if option.HasNoProductVariants {
		conds = append(
			conds,
			qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariants, model.ProductVariantTableColumns.ProductID, model.ProductTableColumns.ID)),
			qm.Where(fmt.Sprintf("%s IS NULL", model.ProductVariantTableColumns.ProductID)),
		)
	}
	if option.VoucherID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.VoucherProducts, model.VoucherProductTableColumns.ProductID, model.ProductTableColumns.ID)),
			option.VoucherID,
		)
	}
	if option.SaleID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.SaleProducts, model.SaleProductTableColumns.ProductID, model.ProductTableColumns.ID)),
			option.SaleID,
		)
	}
	if option.CollectionID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductCollections, model.ProductCollectionTableColumns.ProductID, model.ProductTableColumns.ID)),
			option.CollectionID,
		)
	}
	for _, load := range option.Preloads {
		conds = append(conds, qm.Load(load))
	}

	return conds
}

func (ps *SqlProductStore) FilterByOption(option model_helper.ProductFilterOption) (model.ProductSlice, error) {
	conditions := ps.commonQueryBuilder(option)
	return model.ProductSlice(conditions...).All(ps.GetReplica())
}

// channelQuery is a utility function to compose a filter query on `Channels` table.
//
// `channelSlugOrID` is to filter attribute Channels.Slug = ... OR Channels.Id = ....
//
// `compareToTable` is database table that has property `ChannelID`.
// This argument can be `ProductChannelListings` or `ProductVariantChannelListings`
func (ps *SqlProductStore) channelQuery(channelSlugOrID string, isActive *bool, compareToTable string) squirrel.SelectBuilder {
	var channelActiveExpr string
	if isActive != nil {
		if *isActive {
			channelActiveExpr = fmt.Sprintf("%s AND ", model.ChannelTableColumns.IsActive)
		} else {
			channelActiveExpr = fmt.Sprintf("NOT %s AND ", model.ChannelTableColumns.IsActive)
		}
	}
	return ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.Channels).
		Where(squirrel.Expr(channelActiveExpr)).
		Where(squirrel.And{
			squirrel.Or{
				squirrel.Eq{model.ChannelTableColumns.Slug: channelSlugOrID},
				squirrel.Eq{model.ChannelTableColumns.ID: channelSlugOrID},
			},
			squirrel.Eq{model.ChannelTableColumns.ID: compareToTable + ".channel_id"},
		}).
		Suffix(")").
		Limit(1)
}

// FilterPublishedProducts finds and returns products that belong to given channel slug and are published
//
// refer to ./product_store_doc.md (line 1)
func (ps *SqlProductStore) PublishedProducts(channelSlug string) (model.ProductSlice, error) {
	channelQuery := ps.channelQuery(channelSlug, model_helper.GetPointerOfValue(true), model.TableNames.ProductChannelListings)

	today := util.MillisFromTime(util.StartOfDay(model_helper.GetTimeUTCNow()))

	productChannelListingQuery := ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.ProductChannelListings).
		Where(squirrel.Expr(
			fmt.Sprintf(
				`(%[1]s <= ? OR %[1]s IS NULL)
				AND %[2]s
				AND %[3]s = %[4]s`,
				model.ProductChannelListingTableColumns.PublicationDate, // 1
				model.ProductChannelListingTableColumns.IsPublished,     // 2
				model.ProductChannelListingTableColumns.ProductID,       // 3
				model.ProductTableColumns.ID,                            // 4
			),
			today,
		)).
		Where(channelQuery).
		Suffix(")").
		Limit(1)

	query := ps.
		GetQueryBuilder().
		Select(model.TableNames.Products + ".*").
		From(model.TableNames.Products).
		Where(productChannelListingQuery)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterPublishedProducts_ToSql")
	}

	var res model.ProductSlice
	err = queries.Raw(queryString, args...).Bind(context.Background(), ps.GetReplica(), &res)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find published products with channel slug=%s", channelSlug)
	}

	return res, nil
}

// FilterNotPublishedProducts finds all not published products belong to given channel
//
// refer to ./product_store_doc.md (line 45)
func (ps *SqlProductStore) NotPublishedProducts(channelID string) (model.ProductSlice, error) {
	today := util.MillisFromTime(util.StartOfDay(model_helper.GetTimeUTCNow()))

	return model.ProductSlice(
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductChannelListings, model.ProductTableColumns.ID, model.ProductChannelListingTableColumns.ProductID)),
		model_helper.Or{
			model_helper.And{
				squirrel.Gt{model.ProductChannelListingTableColumns.PublicationDate: today},
				squirrel.Expr(model.ProductChannelListingTableColumns.IsPublished),
			},
			squirrel.Expr("NOT " + model.ProductChannelListingTableColumns.IsPublished),
			squirrel.Expr(model.ProductChannelListingTableColumns.IsPublished + " IS NULL"),
		},
	).All(ps.GetReplica())
}

// PublishedWithVariants finds and returns products.
//
// refer to ./product_store_doc.md (line 157)
func (ps *SqlProductStore) PublishedWithVariants(channelIdOrSlug string) squirrel.SelectBuilder {
	channelQuery1 := ps.channelQuery(channelIdOrSlug, model_helper.GetPointerOfValue(true), model.TableNames.ProductChannelListings)
	today := util.MillisFromTime(util.StartOfDay(model_helper.GetTimeUTCNow()))

	productChannelListingQuery := ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.ProductChannelListings).
		Where(squirrel.Or{
			squirrel.Expr(model.ProductChannelListingTableColumns.PublicationDate + " IS NULL"),
			squirrel.LtOrEq{model.ProductChannelListingTableColumns.PublicationDate: today},
		}).
		Where(squirrel.And{
			squirrel.Expr(model.ProductChannelListingTableColumns.IsPublished),
			squirrel.Eq{model.ProductChannelListingTableColumns.ProductID: model.ProductTableColumns.ID},
		}).
		Where(channelQuery1).
		Suffix(")").
		Limit(1)

	channelQuery2 := ps.channelQuery(channelIdOrSlug, model_helper.GetPointerOfValue(true), model.TableNames.ProductVariantChannelListings)

	productVariantChannelListingQuery := ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.ProductVariantChannelListings).
		Where(channelQuery2).
		Where(squirrel.And{
			squirrel.Expr(model.ProductVariantChannelListingTableColumns.PriceAmount + " IS NOT NULL"),
			squirrel.Eq{model.ProductVariantChannelListingTableColumns.VariantID: model.ProductVariantTableColumns.ID},
		}).
		Suffix(")").
		Limit(1)

	productVariantQuery := ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.TableNames.ProductVariants).
		Where(productVariantChannelListingQuery).
		Where(squirrel.Eq{model.ProductVariantTableColumns.ProductID: model.ProductTableColumns.ID}).
		Suffix(")").
		Limit(1)

	return ps.GetQueryBuilder().
		Select(model.TableNames.Products + ".*").
		From(model.TableNames.Products).
		Where(productChannelListingQuery).
		Where(productVariantQuery)
}

// 1) If requesting user has any of product-related permissions
//
//	+) if `channelSlugOrID` is empty string: returns all products. refer to ./product_store_doc.md (line 241, CASE 2)
//
//	+) if `channelSlugOrID` is provided: refer to ./product_store_doc.md (line 241, CASE 1)
//
// 2) If requesting user is shop visitor: Refer to ./product_store_doc.md (line 241, case 3)
func (ps *SqlProductStore) VisibleToUserProductsQuery(channelSlugOrID string, userHasOneOfProductpermissions bool) squirrel.SelectBuilder {
	// check if requesting user has right to view products
	if userHasOneOfProductpermissions {
		if channelSlugOrID == "" {
			return ps.GetQueryBuilder().Select(model.TableNames.Products + ".*").From(model.TableNames.Products) // find all
		}

		channelQuery := ps.channelQuery(channelSlugOrID, nil, model.TableNames.ProductChannelListings)
		productChannelListingQuery := ps.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			From(model.TableNames.ProductChannelListings).
			Where(channelQuery).
			Where(squirrel.Eq{model.ProductChannelListingTableColumns.ProductID: model.ProductTableColumns.ID}).
			Suffix(")").
			Limit(1)

		return ps.
			GetQueryBuilder().
			Select(model.TableNames.Products + ".*").
			From(model.TableNames.Products).
			Where(productChannelListingQuery)
	}

	return ps.PublishedWithVariants(channelSlugOrID)
}

// SelectForUpdateDiscountedPricesOfCatalogues finds and returns product based on given ids lists.
func (ps *SqlProductStore) SelectForUpdateDiscountedPricesOfCatalogues(transaction boil.ContextTransactor, productIDs, categoryIDs, collectionIDs, variantIDs []string) (model.ProductSlice, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	conds := []qm.QueryMod{}
	orConds := model_helper.Or{}

	if len(productIDs) > 0 {
		orConds = append(orConds, squirrel.Eq{model.ProductTableColumns.ID: productIDs})
	}
	if len(categoryIDs) > 0 {
		orConds = append(orConds, squirrel.Eq{model.ProductTableColumns.CategoryID: categoryIDs})
	}
	if len(collectionIDs) > 0 {
		conds = append(conds, qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductCollections, model.ProductTableColumns.ID, model.ProductCollectionTableColumns.ProductID)))
		orConds = append(orConds, squirrel.Eq{model.ProductCollectionTableColumns.CollectionID: collectionIDs})
	}
	if len(variantIDs) > 0 {
		conds = append(conds, qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariants, model.ProductTableColumns.ID, model.ProductVariantTableColumns.ProductID)))
		orConds = append(orConds, squirrel.Eq{model.ProductVariantTableColumns.ID: variantIDs})
	}

	conds = append(conds, orConds)
	return model.ProductSlice(conds...).All(transaction)
}

// AdvancedFilterQueryBuilder advancedly finds products, filtered using given options
func (ps *SqlProductStore) AdvancedFilterQueryBuilder(input model_helper.ExportProductsFilterOptions) squirrel.SelectBuilder {
	query := ps.GetQueryBuilder().
		Select(model.TableNames.Products + ".*").
		From(model.TableNames.Products)

	if input.Scope == "all" {
		return query
	}
	if input.Scope == "ids" {
		return query.Where(squirrel.Eq{model.ProductTableColumns.ID: input.Ids})
	}

	var channelIdOrSlug string
	if cID := input.Filter.Channel; cID != nil {
		channelIdOrSlug = *cID
	}

	// parse input.Filter
	if input.Filter.IsPublished != nil {
		query = ps.filterIsPublished(query, *input.Filter.IsPublished, channelIdOrSlug)
	}
	if len(input.Filter.Collections) > 0 {
		query = ps.filterCollections(query, input.Filter.Collections)
	}
	if len(input.Filter.Categories) > 0 {
		query = ps.filterCategories(query, input.Filter.Categories)
	}
	if input.Filter.HasCategory != nil {
		// default to has no category
		condition := fmt.Sprintf("%s IS NULL", model.ProductTableColumns.CategoryID)

		if *input.Filter.HasCategory {
			condition = fmt.Sprintf("%s IS NOT NULL", model.ProductTableColumns.CategoryID)
		}
		query = query.Where(condition)
	}
	if input.Filter.Price != nil {
		query = ps.filterVariantPrice(query, *input.Filter.Price, channelIdOrSlug)
	}
	if input.Filter.MinimalPrice != nil {
		query = ps.filterMinimalPrice(query, *input.Filter.MinimalPrice, channelIdOrSlug)
	}
	if len(input.Filter.Attributes) > 0 {
		query = ps.filterAttributes(query, input.Filter.Attributes)
	}
	if input.Filter.StockAvailability != nil {
		query = ps.filterStockAvailability(query, *input.Filter.StockAvailability, channelIdOrSlug)
	}
	// if len(input.Filter.ProductTypes) > 0 {
	// 	query = ps.filterProductTypes(query, input.Filter.ProductTypes)
	// }
	if input.Filter.Stocks != nil {
		query = ps.filterStocks(query, *input.Filter.Stocks)
	}
	if input.Filter.GiftCard != nil {
		query = ps.filterGiftCard(query, *input.Filter.GiftCard)
	}
	if len(input.Filter.Ids) != 0 {
		query = ps.filterProductIDs(query, input.Filter.Ids)
	}
	if input.Filter.HasPreorderedVariants != nil {
		query = ps.filterHasPreorderedVariants(query, *input.Filter.HasPreorderedVariants)
	}
	if input.Filter.Search != nil {
		query = ps.filterSearch(query, *input.Filter.Search)
	}
	if len(input.Filter.Metadata) > 0 {
		conditions := []string{}

		for _, pair := range input.Filter.Metadata {
			if pair != nil && pair.Key != "" {
				if pair.Value == "" {
					expr := fmt.Sprintf(`%s::jsonb ? '%s'`, model.ProductTableColumns.Metadata, pair.Key)
					conditions = append(conditions, expr)
					continue
				}
				expr := fmt.Sprintf(`%s::jsonb @> '{%q:%q}'`, model.ProductTableColumns.Metadata, pair.Key, pair.Value)
				conditions = append(conditions, expr)
			}
		}
		query = query.Where(strings.Join(conditions, " AND "))
	}

	// filter by SortBy
	if input.SortBy != nil && input.SortBy.Field != nil {
		switch *input.SortBy.Field {
		case model_helper.ProductOrderFieldPrice:
			query = query.
				Column(
					fmt.Sprintf(
						`MIN (%s) FILTER (
							WHERE (
								(%s = ? OR %s = ?)
								AND %s IS NOT NULL
							)
						) AS MinVariantsPriceAmount`,
						model.ProductVariantChannelListingTableColumns.PriceAmount,
						model.ChannelTableColumns.ID,
						model.ChannelTableColumns.Slug,
						model.ProductVariantChannelListingTableColumns.PriceAmount,
					),
					channelIdOrSlug,
					channelIdOrSlug,
				).
				LeftJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariants, model.ProductTableColumns.ID, model.ProductVariantTableColumns.ProductID)).
				LeftJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariantChannelListings, model.ProductVariantTableColumns.ID, model.ProductVariantChannelListingTableColumns.VariantID)).
				LeftJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.ProductVariantChannelListingTableColumns.ChannelID)).
				GroupBy(model.ProductTableColumns.ID).
				OrderBy("MinVariantsPriceAmount " + string(input.SortBy.Direction))

		case model_helper.ProductOrderFieldMinimalPrice:
			query = query.
				Column(
					fmt.Sprintf(
						`MIN (%s) FILTER (
							WHERE %s = ? OR %s = ?
						) AS DiscountedPriceAmount`,
						model.ProductChannelListingTableColumns.DiscountedPriceAmount,
						model.ChannelTableColumns.ID,
						model.ChannelTableColumns.Slug,
					),
					channelIdOrSlug,
					channelIdOrSlug,
				).
				LeftJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductChannelListings, model.ProductTableColumns.ID, model.ProductChannelListingTableColumns.ProductID)).
				LeftJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.ProductChannelListingTableColumns.ChannelID)).
				OrderBy("DiscountedPriceAmount " + string(input.SortBy.Direction)).
				GroupBy(model.ProductTableColumns.ID)

		case model_helper.ProductOrderFieldPublished:
			query = query.
				Column(
					fmt.Sprintf(
						`(
							SELECT %[1]s
							FROM %[2]s
							INNER JOIN %[3]s ON %[4]s = %[5]s
							WHERE (
								(%[6]s = ? OR %[7]s = ?)
								AND %[8]s = %[9]s
							)
							ORDER BY %[10]s ASC
							LIMIT 1
						) AS IsPublished`,
						model.ProductChannelListingTableColumns.IsPublished, // 1
						model.TableNames.ProductChannelListings,             // 2
						model.TableNames.Channels,                           // 3
						model.ChannelTableColumns.ID,                        // 4
						model.ProductChannelListingTableColumns.ChannelID,   // 5
						model.ChannelTableColumns.ID,                        // 6
						model.ChannelTableColumns.Slug,                      // 7
						model.ProductChannelListingTableColumns.ProductID,   // 8
						model.ProductTableColumns.ID,                        // 9
						model.ProductChannelListingTableColumns.ID,          // 10
					),
					channelIdOrSlug,
					channelIdOrSlug,
				).
				OrderBy("IsPublished " + string(input.SortBy.Direction))

		case model_helper.ProductOrderFieldPublicationDate:
			query = query.
				Column(
					fmt.Sprintf(
						`(
							SELECT %[1]s
							FROM %[2]s
							INNER JOIN %[3]s ON %[4]s = %[5]s
							WHERE (
								(%[6]s = ? OR %[7]s = ?)
								AND %[8]s = %[9]s
							)
							ORDER BY %[10]s ASC
							LIMIT 1
						) AS PublicationDate`,
						model.ProductChannelListingTableColumns.PublicationDate, // 1
						model.TableNames.ProductChannelListings,                 // 2
						model.TableNames.Channels,                               // 3
						model.ChannelTableColumns.ID,                            // 4
						model.ProductChannelListingTableColumns.ChannelID,       // 5
						model.ChannelTableColumns.ID,                            // 6
						model.ChannelTableColumns.Slug,                          // 7
						model.ProductChannelListingTableColumns.ProductID,       // 8
						model.ProductTableColumns.ID,                            // 9
						model.ProductChannelListingTableColumns.PublicationDate, // 10
					),
					channelIdOrSlug,
					channelIdOrSlug,
				).
				OrderBy("PublicationDate " + string(input.SortBy.Direction))

		case model_helper.ProductOrderFieldCollection:
			// model.CollectionProductRelationTableName
			query = query.
				Column(
					fmt.Sprintf(
						`DENSE_RANK() OVER (
							ORDER BY %s ASC NULLS LAST,
						) AS SortOrder,`,
						model.ProductCollectionTableColumns.SortOrder,
					),
				).
				LeftJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductCollections, model.ProductTableColumns.ID, model.ProductCollectionTableColumns.ProductID)).
				OrderBy("SortOrder " + string(input.SortBy.Direction))
		}
	}

	return query
}

func (ps *SqlProductStore) FilterByQuery(query squirrel.SelectBuilder) (model.ProductSlice, error) {
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByQuery_ToSql")
	}

	var products model.ProductSlice
	err = queries.Raw(queryString, args...).Bind(context.Background(), ps.GetReplica(), &products)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find products by given query")
	}

	return products, nil
}

func (s *SqlProductStore) CountByCategoryIDs(categoryIDs []string) ([]*model_helper.ProductCountByCategoryID, error) {
	var res []*model_helper.ProductCountByCategoryID
	err := model.ProductSlice(
		qm.Select(
			model.ProductTableColumns.CategoryID,
			fmt.Sprintf("COUNT (%s) as %q", model.ProductTableColumns.ID, "product_count"),
		),
		qm.GroupBy(model.ProductTableColumns.CategoryID),
	).Bind(context.Background(), s.GetReplica(), &res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count products by given category ids")
	}

	return res, nil
}
