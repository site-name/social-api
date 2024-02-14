package product

import (
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlProductStore struct {
	store.Store
}

func NewSqlProductStore(s store.Store) store.ProductStore {
	return &SqlProductStore{s}
}

// Save inserts given product into database then returns it
func (ps *SqlProductStore) Save(tx *gorm.DB, product *model.Product) (*model.Product, error) {
	if tx == nil {
		tx = ps.GetMaster()
	}
	if err := tx.Save(product).Error; err != nil {
		if ps.IsUniqueConstraintError(err, []string{"Name", "products_name_key", "idx_products_name_unique"}) {
			return nil, store.NewErrInvalidInput("Product", "name", product.Name)
		}
		if ps.IsUniqueConstraintError(err, []string{"Slug", "products_slug_key", "idx_products_slug_unique"}) {
			return nil, store.NewErrInvalidInput("Product", "slug", product.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save Product with productId=%s", product.Id)
	}

	return product, nil
}

func (ps *SqlProductStore) commonQueryBuilder(option *model.ProductFilterOption) (*gorm.DB, squirrel.Sqlizer) {
	db := ps.GetReplica()
	conditions := squirrel.And{}

	if option.Conditions != nil {
		conditions = append(conditions, option.Conditions)
	}
	if option.Limit > 0 {
		db = db.Limit(int(option.Limit))
	}
	for _, preload := range option.Preloads {
		db = db.Preload(preload)
	}

	if option.ProductVariantID != nil {
		db = db.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.ProductVariantTableName,       // 1
				model.ProductTableName,              // 2
				model.ProductVariantColumnProductID, // 3
				model.ProductColumnId,               // 4
			),
		)
		conditions = append(conditions, option.ProductVariantID)
	} else if option.HasNoProductVariants {
		db = db.Joins(
			fmt.Sprintf(
				"LEFT JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.ProductVariantTableName,       // 1
				model.ProductTableName,              // 2
				model.ProductVariantColumnProductID, // 3
				model.ProductColumnId,               // 4
			),
		)
		conditions = append(conditions, squirrel.Expr(model.ProductVariantTableName+"."+model.ProductVariantColumnProductID+" IS NULL"))
	}
	if option.VoucherID != nil {
		db = db.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.VoucherProductTableName, // 1
				model.ProductTableName,        // 2
				"product_id",                  // 3
				model.ProductColumnId,         // 4
			),
		)
		conditions = append(conditions, option.VoucherID)
	}
	if option.SaleID != nil {
		db = db.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.SaleProductTableName, // 1
				model.ProductTableName,     // 2
				"product_id",               // 3
				model.ProductColumnId,      // 4
			),
		)
		conditions = append(conditions, option.SaleID)
	}
	if option.CollectionID != nil {
		db = db.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.CollectionProductRelationTableName, // 1
				model.ProductTableName,                   // 2
				model.CollectionProductColumnProductID,   // 3
				model.ProductColumnId,                    // 4
			),
		)
		conditions = append(conditions, option.CollectionID)
	}

	return db, conditions
}

// FilterByOption finds and returns all products that satisfy given option
func (ps *SqlProductStore) FilterByOption(option *model.ProductFilterOption) ([]*model.Product, error) {
	db, conditions := ps.commonQueryBuilder(option)
	args, err := store.BuildSqlizer(conditions, "Product_FilterByOption")
	if err != nil {
		return nil, err
	}
	var products model.Products
	err = db.Find(&products, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find productswith given options")
	}

	return products, nil
}

// GetByOption finds and returns 1 product that satisfies given option
func (ps *SqlProductStore) GetByOption(option *model.ProductFilterOption) (*model.Product, error) {
	option.Limit = 0
	db, conditions := ps.commonQueryBuilder(option)
	args, err := store.BuildSqlizer(conditions, "Product_GetByOption")
	if err != nil {
		return nil, err
	}
	var res model.Product
	err = db.First(&res, args...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ProductTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find product by given option")
	}

	return &res, nil
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
func (ps *SqlProductStore) PublishedProducts(channelSlug string) ([]*model.Product, error) {
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
func (ps *SqlProductStore) SelectForUpdateDiscountedPricesOfCatalogues(transaction *gorm.DB, productIDs, categoryIDs, collectionIDs, variantIDs []string) ([]*model.Product, error) {
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

// FilterByQuery finds and returns products with given query, limit, createdAtGt
func (ps *SqlProductStore) FilterByQuery(query squirrel.SelectBuilder) (model.Products, error) {
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByQuery_ToSql")
	}

	var products model.Products
	err = ps.GetReplica().Raw(queryString, args...).Scan(&products).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find products with given query and conditions")
	}
	return products, nil
}

func (s *SqlProductStore) CountByCategoryIDs(categoryIDs []string) ([]*model.ProductCountByCategoryID, error) {
	var res []*model.ProductCountByCategoryID
	query := fmt.Sprintf(
		`SELECT
			%[1]s.%[2]s,
			COUNT(%[1]s.%[3]s) AS ProductCount
		FROM
			%[1]s
		WHERE
			%[1]s.%[2]s IN ?
		GROUP BY %[1]s.%[2]s`,

		model.ProductTableName,        // 1
		model.ProductColumnCategoryID, // 2
		model.ProductColumnId,         // 3
	)
	err := s.GetMaster().Raw(query, categoryIDs).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to count products by given category ids")
	}

	return res, nil
}
