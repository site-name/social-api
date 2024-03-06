package product

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	return model.Products(conditions...).All(ps.GetReplica())
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
			channelActiveExpr = fmt.Sprintf("%s.%s AND ", model.ChannelTableName, model.ChannelColumnIsActive)
		} else {
			channelActiveExpr = fmt.Sprintf("NOT %s.%s AND ", model.ChannelTableName, model.ChannelColumnIsActive)
		}
	}
	return ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.ChannelTableName).
		Where(channelActiveExpr+"(Channels.Slug = ? OR Channels.Id = ?) AND Channels.Id = ?.ChannelID", channelSlugOrID, channelSlugOrID, compareToTable).
		Suffix(")").
		Limit(1)
}

// FilterPublishedProducts finds and returns products that belong to given channel slug and are published
//
// refer to ./product_store_doc.md (line 1)
func (ps *SqlProductStore) PublishedProducts(channelSlug string) (model.ProductSlice, error) {
	channelQuery := ps.channelQuery(channelSlug, model.GetPointerOfValue(true), model.ProductChannelListingTableName)

	today := util.StartOfDay(time.Now())

	productChannelListingQuery := ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.ProductChannelListingTableName).
		Where(
			fmt.Sprintf(`(
					%[1]s.%[2]s <= ?
					OR %[1]s.%[2]s IS NULL
				)
				AND %[1]s.%[3]s
				AND %[1]s.%[4]s = %[5]s.%[6]s`,

				model.ProductChannelListingTableName,       // 1
				model.PublishableColumnPublicationDate,     // 2
				model.PublishableColumnIsPublished,         // 3
				model.ProductChannelListingColumnProductID, // 4
				model.ProductTableName,                     // 5
				model.ProductColumnId,                      // 6
			),
			today,
		).
		Where(channelQuery).
		Suffix(")").
		Limit(1)

	query := ps.
		GetQueryBuilder().
		Select(model.ProductTableName + ".*").
		From(model.ProductTableName).
		Where(productChannelListingQuery)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterPublishedProducts_ToSql")
	}

	var res model.Products
	err = ps.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find published products with channel slug=%s", channelSlug)
	}

	return res, nil
}

// FilterNotPublishedProducts finds all not published products belong to given channel
//
// refer to ./product_store_doc.md (line 45)
func (ps *SqlProductStore) NotPublishedProducts(channelID string) (model.Products, error) {
	today := util.StartOfDay(time.Now()) // start of day

	queryString, args, err := ps.GetQueryBuilder().
		Select(model.ProductTableName + ".*").
		From(model.ProductTableName).
		InnerJoin(
			fmt.Sprintf(
				"%[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.ProductChannelListingTableName,       // 1
				model.ProductTableName,                     // 2
				model.ProductChannelListingColumnProductID, // 3
				model.ProductColumnId,                      // 4
			),
		).
		Where(squirrel.Or{
			squirrel.And{
				squirrel.Expr(model.ProductChannelListingTableName+"."+model.PublishableColumnPublicationDate+" > ?", today),
				squirrel.Expr(model.ProductChannelListingTableName + "." + model.PublishableColumnIsPublished),
			},
			squirrel.Expr("NOT " + model.ProductChannelListingTableName + "." + model.PublishableColumnIsPublished),
			squirrel.Expr(model.ProductChannelListingTableName + "." + model.PublishableColumnIsPublished + " IS NULL"),
		}).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "NotPublishedProducts_ToSql")
	}

	var res model.Products
	err = ps.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find not published product with channel id=%s", channelID)
	}

	return res, nil
}

// PublishedWithVariants finds and returns products.
//
// refer to ./product_store_doc.md (line 157)
func (ps *SqlProductStore) PublishedWithVariants(channelIdOrSlug string) squirrel.SelectBuilder {
	channelQuery1 := ps.channelQuery(channelIdOrSlug, model.GetPointerOfValue(true), model.ProductChannelListingTableName)
	today := util.StartOfDay(time.Now())

	productChannelListingQuery := ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.ProductChannelListingTableName).
		Where("ProductChannelListings.PublicationDate IS NULL OR ProductChannelListings.PublicationDate <= ?", today).
		Where("ProductChannelListings.IsPublished AND ProductChannelListings.ProductID = Products.Id").
		Where(channelQuery1).
		Suffix(")").
		Limit(1)

	channelQuery2 := ps.channelQuery(channelIdOrSlug, model.GetPointerOfValue(true), model.ProductVariantChannelListingTableName)

	productVariantChannelListingQuery := ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.ProductVariantChannelListingTableName).
		Where(channelQuery2).
		Where("ProductVariantChannelListings.PriceAmount IS NOT NULL AND ProductVariantChannelListings.VariantID = ProductVariants.Id").
		Suffix(")").
		Limit(1)

	productVariantQuery := ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.ProductVariantTableName).
		Where(productVariantChannelListingQuery).
		Where("Products.Id = ProductVariants.ProductID").
		Suffix(")").
		Limit(1)

	return ps.GetQueryBuilder().
		Select(model.ProductTableName + ".*").
		From(model.ProductTableName).
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
			return ps.GetQueryBuilder().Select(model.ProductTableName + ".*").From(model.ProductTableName) // find all
		}

		channelQuery := ps.channelQuery(channelSlugOrID, nil, model.ProductChannelListingTableName)
		productChannelListingQuery := ps.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			From(model.ProductChannelListingTableName).
			Where(channelQuery).
			Where("ProductChannelListings.ProductID = Products.Id").
			Suffix(")").
			Limit(1)

		return ps.
			GetQueryBuilder().
			Select(model.ProductTableName + ".*").
			From(model.ProductTableName).
			Where(productChannelListingQuery)
	}

	return ps.PublishedWithVariants(channelSlugOrID)
}

// SelectForUpdateDiscountedPricesOfCatalogues finds and returns product based on given ids lists.
func (ps *SqlProductStore) SelectForUpdateDiscountedPricesOfCatalogues(transaction boil.ContextTransactor, productIDs, categoryIDs, collectionIDs, variantIDs []string) (model.ProductSlice, error) {
	if transaction == nil {
		return nil, store.NewErrInvalidInput("SelectForUpdateDiscountedPricesOfCatalogues", "transaction", nil)
	}
	query := ps.GetQueryBuilder().
		Select(model.ProductTableName + ".*").
		Distinct().
		From(model.ProductTableName)

	if transaction != nil {
		query = query.Suffix("FOR UPDATE")
	}

	orCondition := squirrel.Or{}

	if len(productIDs) > 0 {
		orCondition = append(orCondition, squirrel.Eq{"Products.Id": productIDs})
	}
	if len(categoryIDs) > 0 {
		orCondition = append(orCondition, squirrel.Eq{"Products.CategoryID": categoryIDs})
	}
	if len(collectionIDs) > 0 {
		query = query.InnerJoin(model.CollectionProductRelationTableName + " ON (Products.Id = ProductCollections.ProductID)")
		orCondition = append(orCondition, squirrel.Eq{"ProductCollections.CollectionID": collectionIDs})
	}
	if len(variantIDs) > 0 {
		query = query.InnerJoin(model.ProductVariantTableName + " ON Products.Id = ProductVariants.ProductID")
		orCondition = append(orCondition, squirrel.Eq{model.ProductVariantTableName + ".Id": variantIDs})
	}

	if len(orCondition) > 0 {
		query = query.Where(orCondition)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "SelectForUpdateDiscountedPricesOfCatalogues_ToSql")
	}

	var products model.Products
	err = transaction.Raw(queryString, args...).Scan(&products).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find products by given params")
	}

	return products, nil
}

// AdvancedFilterQueryBuilder advancedly finds products, filtered using given options
func (ps *SqlProductStore) AdvancedFilterQueryBuilder(input *model.ExportProductsFilterOptions) squirrel.SelectBuilder {
	query := ps.GetQueryBuilder().
		Select(model.ProductTableName + ".*").
		From(model.ProductTableName)

	if input.Scope == "all" {
		return query
	}
	if input.Scope == "ids" {
		return query.Where(squirrel.Eq{"Products.Id": input.Ids})
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
		condition := fmt.Sprintf("%s.%s IS NULL", model.ProductTableName, model.ProductColumnCategoryID)

		if *input.Filter.HasCategory {
			condition = fmt.Sprintf("%s.%s IS NOT NULL", model.ProductTableName, model.ProductColumnCategoryID)
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
	if len(input.Filter.ProductTypes) > 0 {
		query = ps.filterProductTypes(query, input.Filter.ProductTypes)
	}
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
	if meta := input.Filter.Metadata; len(meta) > 0 {
		conditions := []string{}

		for _, pair := range meta {
			if pair != nil && pair.Key != "" {
				if pair.Value == "" {
					expr := fmt.Sprintf(`Products.Metadata::jsonb ? '%s'`, pair.Key)
					conditions = append(conditions, expr)
					continue
				}
				expr := fmt.Sprintf(`Products.Metadata::jsonb @> '{%q:%q}'`, pair.Key, pair.Value)
				conditions = append(conditions, expr)
			}
		}
		query = query.Where(strings.Join(conditions, " AND "))
	}

	// filter by SortBy
	if input.SortBy != nil && input.SortBy.Field != nil {
		switch *input.SortBy.Field {
		case model.ProductOrderFieldPrice:
			query = query.
				Column(`MIN(
					ProductVariantChannelListings.PriceAmount
				) FILTER (
					WHERE (
						(Channel.Id = ? OR Channel.Slug = ?)
						AND ProductVariantChannelListings.PriceAmount IS NOT NULL
					)
				) AS MinVariantsPriceAmount`, channelIdOrSlug, channelIdOrSlug).
				LeftJoin(model.ProductVariantTableName + " ON Products.Id = ProductVariants.ProductID").
				LeftJoin(model.ProductVariantChannelListingTableName + " ON ProductVariants.Id = ProductVariantChannelListings.VariantID").
				LeftJoin(model.ChannelTableName + " ON Channels.Id = ProductVariantChannelListings.ChannelID").
				GroupBy("Products.Id").
				OrderBy("MinVariantsPriceAmount " + string(input.SortBy.Direction))

		case model.ProductOrderFieldMinimalPrice:
			query = query.
				Column(`MIN(
				ProductChannelListings.DiscountedPriceAmount
			) FILTER (
				WHERE Channels.Slug = ? OR Channels.Id = ?
			) AS DiscountedPriceAmount`, channelIdOrSlug, channelIdOrSlug).
				LeftJoin(model.ProductChannelListingTableName + " ON ProductChannelListings.ProductID = Products.Id").
				LeftJoin(model.ChannelTableName + " ON Channels.Id = ProductChannelListings.ChannelID").
				OrderBy("DiscountedPriceAmount " + string(input.SortBy.Direction)).
				GroupBy("Products.Id")

		case model.ProductOrderFieldPublished:
			query = query.
				Column(`(
					SELECT PC.IsPublished
					FROM ProductChannelListings PC
					INNER JOIN Channels C ON C.Id = PC.ChannelID
					WHERE (
						(C.Slug = ? OR C.Id = ?)
						AND PC.ProductID = Products.Id
					)
					ORDER BY PC.Id ASC
					LIMIT 1
				) AS IsPublished`, channelIdOrSlug, channelIdOrSlug).
				OrderBy("IsPublished " + string(input.SortBy.Direction))

		case model.ProductOrderFieldPublicationDate:
			query = query.
				Column(`(
					SELECT PC.PublicationDate
					FROM ProductChannelListings PC
					INNER JOIN Channels C ON C.Id = PC.ChannelID
					WHERE (
						(C.Id = ? OR C.Slug = ?)
						AND PC.ProductID = Products.Id
					)
					ORDER BY PC.Id ASC
					LIMIT 1
				) AS PublicationDate`, channelIdOrSlug, channelIdOrSlug).
				OrderBy("PublicationDate " + string(input.SortBy.Direction))

		case model.ProductOrderFieldCollection:
			// model.CollectionProductRelationTableName
			query = query.
				Column(`DENSE_RANK() OVER (
					ORDER BY ProductCollections.SortOrder ASC NULLS LAST,
					ProductCollections.Id
				) AS SortOrder`).
				LeftJoin(model.CollectionProductRelationTableName + " ON Products.Id = ProductCollections.ProductID").
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
	err := model.Products(
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
