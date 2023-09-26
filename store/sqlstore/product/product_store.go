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

func (ps *SqlProductStore) ScanFields(prd *model.Product) []interface{} {
	return []interface{}{
		&prd.Id,
		&prd.ProductTypeID,
		&prd.Name,
		&prd.Slug,
		&prd.Description,
		&prd.DescriptionPlainText,
		&prd.CategoryID,
		&prd.CreateAt,
		&prd.UpdateAt,
		&prd.ChargeTaxes,
		&prd.Weight,
		&prd.WeightUnit,
		&prd.DefaultVariantID,
		&prd.Rating,
		&prd.Metadata,
		&prd.PrivateMetadata,
		&prd.SeoTitle,
		&prd.SeoDescription,
	}
}

// Save inserts given product into database then returns it
func (ps *SqlProductStore) Save(product *model.Product) (*model.Product, error) {
	if err := ps.GetMaster().Save(product).Error; err != nil {
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

func (ps *SqlProductStore) commonQueryBuilder(option *model.ProductFilterOption) (string, []interface{}, error) {
	query := ps.GetQueryBuilder().
		Select(model.ProductTableName + ".*").
		From(model.ProductTableName).Where(option.Conditions)

	// parse option
	if option.Limit > 0 {
		query = query.Limit(option.Limit)
	}

	if option.ProductVariantID != nil {
		query = query.
			InnerJoin(model.ProductVariantTableName + " ON Products.Id = ProductVariants.ProductID").
			Where(option.ProductVariantID)
	} else if option.HasNoProductVariants {
		query = query.
			LeftJoin(model.ProductVariantTableName + " ON ProductVariants.ProductID = Products.Id").
			Where(model.ProductVariantTableName + ".ProductID IS NULL")
	}

	if option.VoucherID != nil {
		query = query.
			InnerJoin(model.VoucherProductTableName + " ON Products.Id = VoucherProducts.ProductID").
			Where(option.VoucherID)
	}
	if option.SaleID != nil {
		query = query.
			InnerJoin(model.SaleProductTableName + " ON Products.Id = SaleProducts.ProductID").
			Where(option.SaleID)
	}
	if option.CollectionID != nil {
		query = query.
			InnerJoin(model.CollectionProductRelationTableName + " ON ProductCollections.ProductID = Products.Id").
			Where(option.CollectionID)
	}

	return query.ToSql()
}

// FilterByOption finds and returns all products that satisfy given option
func (ps *SqlProductStore) FilterByOption(option *model.ProductFilterOption) ([]*model.Product, error) {
	queryString, args, err := ps.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var products model.Products
	err = ps.GetReplica().Raw(queryString, args...).Scan(&products).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find products by given option")
	}

	var (
		productIDs  = make([]string, 0, len(products))
		productsMap = map[string]*model.Product{} // productsMap has keys are product ids
	)
	for _, product := range products {
		_, met := productsMap[product.Id]
		if !met {
			productsMap[product.Id] = product
			productIDs = append(productIDs, product.Id)
		}
	}

	// check if need prefetch related assigned product attribute
	if option.PrefetchRelatedAssignedProductAttributes && len(productIDs) > 0 {
		assignedAttributes, err := ps.AssignedProductAttribute().FilterByOptions(&model.AssignedProductAttributeFilterOption{
			Conditions: squirrel.Eq{model.AssignedProductAttributeTableName + ".ProductID": productIDs},
		})
		if err != nil {
			return nil, err
		}

		for _, attr := range assignedAttributes {
			product, ok := productsMap[attr.ProductID]
			if ok && product != nil {
				product.Attributes = append(product.Attributes, attr)
			}
		}
	}

	// check if need prefetch related categories
	if option.PrefetchRelatedCategory && len(productIDs) > 0 {
		categories, err := ps.Category().FilterByOption(&model.CategoryFilterOption{
			Conditions: squirrel.Eq{model.CategoryTableName + ".id": products.CategoryIDs()},
		})
		if err != nil {
			return nil, err
		}

		var categoriesMap = map[string]*model.Category{}
		for _, cate := range categories {
			categoriesMap[cate.Id] = cate
		}

		for _, prd := range products {
			if prd.CategoryID != nil {
				prd.SetCategory(categoriesMap[*prd.CategoryID])
			}
		}
	}

	// check if need prefetch related collections
	if option.PrefetchRelatedCollections && len(productIDs) > 0 {
		collectionProducts, err := ps.CollectionProduct().FilterByOptions(&model.CollectionProductFilterOptions{
			Conditions:              squirrel.Eq{model.CollectionProductRelationTableName + ".ProductID": productIDs},
			SelectRelatedCollection: true,
		})
		if err != nil {
			return nil, err
		}

		for _, rel := range collectionProducts {
			product, ok := productsMap[rel.ProductID]
			if ok && product != nil {
				product.Collections = append(product.Collections, rel.GetCollection())
			}
		}
	}

	// check if we need to prefetch related product type
	if option.PrefetchRelatedProductType && len(productIDs) > 0 {
		productTypes, err := ps.ProductType().ProductTypesByProductIDs(productIDs)
		if err != nil {
			return nil, err
		}

		var productTypesMap = map[string]*model.ProductType{}
		for _, prdType := range productTypes {
			productTypesMap[prdType.Id] = prdType
		}

		for _, product := range products {
			product.SetProductType(productTypesMap[product.ProductTypeID])
		}
	}

	// check if we need to prefetch related file infos
	if option.PrefetchRelatedMedia && len(productIDs) > 0 {
		fileInfos, err := ps.FileInfo().GetWithOptions(&model.GetFileInfosOptions{
			Conditions: squirrel.Eq{model.FileInfoTableName + ".ParentID": productIDs},
		})
		if err != nil {
			return nil, err
		}

		for _, info := range fileInfos {
			product, ok := productsMap[info.ParentID]
			if ok && product != nil {
				product.SetMedias(append(product.GetMedias(), info))
			}
		}
	}

	return products, nil
}

// GetByOption finds and returns 1 product that satisfies given option
func (ps *SqlProductStore) GetByOption(option *model.ProductFilterOption) (*model.Product, error) {
	option.Limit = 0
	queryString, args, err := ps.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res model.Product
	err = ps.GetReplica().Raw(queryString, args...).Scan(&res).Error
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
// `channel_Slug_or_ID` is to filter attribute Channels.Slug = ... OR Channels.Id = ....
//
// `compareToTable` is database table that has property `ChannelID`.
// This argument can be `ProductChannelListings` or `ProductVariantChannelListings`
func (ps *SqlProductStore) channelQuery(channel_Slug_or_ID string, isActive *bool, compareToTable string) squirrel.SelectBuilder {
	var channelActiveExpr string
	if isActive != nil {
		if *isActive {
			channelActiveExpr = "Channels.IsActive AND "
		} else {
			channelActiveExpr = "NOT Channels.IsActive AND "
		}
	}
	return ps.
		GetQueryBuilder(squirrel.Question).
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(model.ChannelTableName).
		Where(channelActiveExpr+"(Channels.Slug = ? OR Channels.Id = ?) AND Channels.Id = ?.ChannelID", channel_Slug_or_ID, channel_Slug_or_ID, compareToTable).
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
		Where(`(ProductChannelListings.PublicationDate <= ? OR 
			ProductChannelListings.PublicationDate IS NULL)
			AND ProductChannelListings.IsPublished
			AND ProductChannelListings.ProductID = Products.Id`, today).
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
func (ps *SqlProductStore) NotPublishedProducts(channelID string) (
	[]*struct {
		model.Product
		IsPublished     bool
		PublicationDate *time.Time
	},
	error,
) {
	today := util.StartOfDay(time.Now()) // start of day

	isPublishedColumnSelect := ps.GetQueryBuilder(squirrel.Question).
		Select("ProductChannelListings.IsPublished").
		From(model.ProductChannelListingTableName).
		Where("ProductChannelListings.ProductID = Products.Id AND ProductChannelListings.ChannelID = ?", channelID).
		Limit(1)

	publicationDateColumnSelect := ps.GetQueryBuilder(squirrel.Question).
		Select("ProductChannelListings.PublicationDate").
		From(model.ProductChannelListingTableName).
		Where("ProductChannelListings.ProductID = Products.Id AND ProductChannelListings.ChannelID = ?", channelID).
		Limit(1)

	queryString, args, err := ps.GetQueryBuilder().
		Select(model.ProductTableName + ".*").
		Column(squirrel.Alias(isPublishedColumnSelect, "IsPublished")).
		Column(squirrel.Alias(publicationDateColumnSelect, "PublicationDate")).
		From(model.ProductTableName).
		Where(squirrel.Or{
			squirrel.And{
				squirrel.Expr("Products.PublicationDate > ?", today),
				squirrel.Expr("Products.IsPublished"),
			},
			squirrel.Expr("NOT Products.IsPublished"),
			squirrel.Expr("Products.IsPublished IS NULL"),
		}).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "NotPublishedProducts_ToSql")
	}

	var res []*struct {
		model.Product
		IsPublished     bool
		PublicationDate *time.Time
	}

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

	if transaction == nil {
		transaction = ps.GetMaster()
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
	if len(input.Filter.Categories) != 0 {
		query = ps.filterCategories(query, input.Filter.Categories)
	}
	if input.Filter.HasCategory != nil {
		// default to has no category
		condition := "Products.CategoryID IS NULL"

		if *input.Filter.HasCategory {
			condition = "Products.CategoryID IS NOT NULL"
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
	err := s.GetMaster().Raw("SELECT P.CategoryID, COUNT(p.Id) AS ProductCount FROM Products P WHERE P.CategoryID IN ? GROUP BY P.CategoryID", categoryIDs).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to count products by given category ids")
	}

	return res, nil
}
