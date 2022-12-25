package product

import (
	"database/sql"
	"strings"
	timemodule "time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlProductStore struct {
	store.Store
}

func NewSqlProductStore(s store.Store) store.ProductStore {
	return &SqlProductStore{s}
}

func (ps *SqlProductStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"ProductTypeID",
		"Name",
		"Slug",
		"Description",
		"DescriptionPlainText",
		"CategoryID",
		"CreateAt",
		"UpdateAt",
		"ChargeTaxes",
		"Weight",
		"WeightUnit",
		"DefaultVariantID",
		"Rating",
		"Metadata",
		"PrivateMetadata",
		"SeoTitle",
		"SeoDescription",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
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
	product.PreSave()
	if err := product.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.ProductTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
	if _, err := ps.GetMasterX().NamedExec(query, product); err != nil {
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
		Select(ps.ModelFields(store.ProductTableName + ".")...).
		From(store.ProductTableName).
		OrderBy(store.TableOrderingMap[store.ProductTableName])

	// parse option
	if option.Limit != nil {
		query = query.Limit(*option.Limit)
	}
	if option.CreateAt != nil {
		query = query.Where(option.CreateAt)
	}
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.ProductVariantID != nil {
		// decide which type of join to use (LEFT or INNER)
		var joinFunc func(join string, rest ...interface{}) squirrel.SelectBuilder = query.InnerJoin

		strExpr, _, _ := option.ProductVariantID.ToSql()
		if strings.Contains(strings.ToUpper(strExpr), "IS NULL") {
			joinFunc = query.LeftJoin
		}

		query = joinFunc(store.ProductVariantTableName + " ON (Products.Id = ProductVariants.ProductID)").
			Where(option.ProductVariantID)
	}
	if option.VoucherID != nil {
		query = query.
			InnerJoin(store.VoucherProductTableName + " ON Products.Id = VoucherProducts.ProductID").
			Where(option.VoucherID)
	}
	if option.SaleID != nil {
		query = query.
			InnerJoin(store.SaleProductRelationTableName + " ON Products.Id = SaleProducts.ProductID").
			Where(option.SaleID)
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
	err = ps.GetReplicaX().Select(&products, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find products by given option")
	}

	var (
		productIDs  = products.IDs()
		productsMap = map[string]*model.Product{} // productsMap has keys are product ids
	)
	for _, product := range products {
		productsMap[product.Id] = product
	}

	// check if need prefetch related assigned product attribute
	if option.PrefetchRelatedAssignedProductAttributes && len(productIDs) > 0 {
		assignedAttributes, err := ps.AssignedProductAttribute().FilterByOptions(&model.AssignedProductAttributeFilterOption{
			ProductID: squirrel.Eq{store.AssignedProductAttributeTableName + ".ProductID": productIDs},
		})
		if err != nil {
			return nil, err
		}

		for _, attr := range assignedAttributes {
			product, ok := productsMap[attr.ProductID]
			if ok && product != nil {
				product.AssignedProductAttributes = append(product.AssignedProductAttributes, attr)
			}
		}
	}

	// check if need prefetch related categories
	if option.PrefetchRelatedCategory && len(productIDs) > 0 {
		categories, err := ps.Category().FilterByOption(&model.CategoryFilterOption{
			Id: squirrel.Eq{store.CategoryTableName + ".Id": products.CategoryIDs()},
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
				prd.Category = categoriesMap[*prd.CategoryID]
			}
		}
	}

	// check if need prefetch related collections
	if option.PrefetchRelatedCollections && len(productIDs) > 0 {
		collectionProducts, err := ps.CollectionProduct().FilterByOptions(&model.CollectionProductFilterOptions{
			ProductID:               squirrel.Eq{store.CollectionProductRelationTableName + ".ProductID": productIDs},
			SelectRelatedCollection: true,
		})
		if err != nil {
			return nil, err
		}

		for _, rel := range collectionProducts {
			product, ok := productsMap[rel.ProductID]
			if ok && product != nil {
				product.Collections = append(product.Collections, rel.Collection)
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
			product.ProductType = productTypesMap[product.ProductTypeID]
		}
	}

	// check if we need to prefetch related file infos
	if option.PrefetchRelatedMedia && len(productIDs) > 0 {
		fileInfos, err := ps.FileInfo().GetWithOptions(nil, nil, &model.GetFileInfosOptions{
			ParentID: productIDs,
		})
		if err != nil {
			return nil, err
		}

		for _, info := range fileInfos {
			product, ok := productsMap[info.ParentID]
			if ok && product != nil {
				product.Medias = append(product.Medias, info)
			}
		}
	}

	return products, nil
}

// GetByOption finds and returns 1 product that satisfies given option
func (ps *SqlProductStore) GetByOption(option *model.ProductFilterOption) (*model.Product, error) {
	queryString, args, err := ps.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res model.Product
	err = ps.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find product by given option")
	}

	return &res, nil
}

// channelQuery is a utility function to compose a filter query on `Channels` table.
//
// `channelSlug` is to filter attribute Channels.Slug.
//
// `compareToTable` is database table that has property `ChannelID`.
// This argument can be `ProductChannelListings` or `ProductVariantChannelListings`
func (ps *SqlProductStore) channelQuery(channelSlug string, isActive *bool, compareToTable string) squirrel.SelectBuilder {
	var channelActiveExpr string
	if isActive != nil {
		if *isActive {
			channelActiveExpr = "Channels.IsActive"
		} else {
			channelActiveExpr = "NOT Channels.IsActive"
		}
	}
	return ps.
		GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ChannelTableName).
		Where(channelActiveExpr+" AND Channels.Slug = ? AND Channels.Id = ?.ChannelID", channelSlug, compareToTable).
		Suffix(")").
		Limit(1)
}

// FilterPublishedProducts finds and returns products that belong to given channel slug and are published
//
// refer to ./product_store_doc.md (line 1)
func (ps *SqlProductStore) PublishedProducts(channelSlug string) ([]*model.Product, error) {

	channelQuery := ps.channelQuery(channelSlug, model.NewBool(true), store.ProductChannelListingTableName)

	today := util.StartOfDay(timemodule.Now().UTC())

	productChannelListingQuery := ps.
		GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ProductChannelListingTableName).
		Where(squirrel.And{
			squirrel.Or{
				squirrel.LtOrEq{"ProductChannelListings.PublicationDate": today},
				squirrel.Eq{"ProductChannelListings.PublicationDate": nil},
			},
			squirrel.Eq{"ProductChannelListings.IsPublished": true},
			squirrel.Expr("ProductChannelListings.ProductID = Products.Id"),
			channelQuery,
		}).
		Suffix(")").
		Limit(1)

	query := ps.
		GetQueryBuilder().
		Select("*").
		From(store.ProductTableName).
		Where(productChannelListingQuery).
		OrderBy(store.TableOrderingMap[store.ProductTableName])

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterPublishedProducts_ToSql")
	}

	var res model.Products
	err = ps.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find published products with channel slug=%s", channelSlug)
	}

	return res, nil
}

// FilterNotPublishedProducts finds all not published products belong to given channel
//
// refer to ./product_store_doc.md (line 45)
func (ps *SqlProductStore) NotPublishedProducts(channelSlug string) (
	[]*struct {
		model.Product
		IsPublished     bool
		PublicationDate *timemodule.Time
	},
	error,
) {
	today := util.StartOfDay(timemodule.Now().UTC()) // start of day

	isPublishedColumnSelect := ps.GetQueryBuilder().
		Select("ProductChannelListings.IsPublished").
		From(store.ProductChannelListingTableName).
		InnerJoin(store.ChannelTableName+" ON (ProductChannelListings.ChannelID = Channels.Id)").
		Where("ProductChannelListings.ProductID = Products.Id AND Channels.Slug = ?", channelSlug).
		OrderBy(store.TableOrderingMap[store.ProductChannelListingTableName]).
		Limit(1)

	publicationDateColumnSelect := ps.GetQueryBuilder().
		Select("ProductChannelListings.PublicationDate").
		From(store.ProductChannelListingTableName).
		InnerJoin(store.ChannelTableName+" ON (Channels.Id = ProductChannelListings.ChannelID)").
		Where("ProductChannelListings.ProductID = Products.Id AND Channels.Slug = ?", channelSlug).
		OrderBy(store.TableOrderingMap[store.ProductChannelListingTableName]).
		Limit(1)

	queryString, args, err := ps.GetQueryBuilder().
		Select(ps.ModelFields(store.ProductTableName + ".")...).
		Column(squirrel.Alias(isPublishedColumnSelect, "IsPublished")).
		Column(squirrel.Alias(publicationDateColumnSelect, "PublicationDate")).
		From(store.ProductTableName).
		Where(squirrel.Or{
			squirrel.And{
				squirrel.Expr("PublicationDate :: date > ?", today),
				squirrel.Expr("IsPublished"),
			},
			squirrel.Expr("NOT IsPublished"),
			squirrel.Expr("IsPublished IS NULL"),
		}).
		OrderBy(store.TableOrderingMap[store.ProductTableName]).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "NotPublishedProducts_ToSql")
	}

	var res []*struct {
		model.Product
		IsPublished     bool
		PublicationDate *timemodule.Time
	}

	err = ps.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find not published product with channel slug=%s", channelSlug)
	}

	return res, nil
}

// PublishedWithVariants finds and returns products.
//
// refer to ./product_store_doc.md (line 157)
func (ps *SqlProductStore) PublishedWithVariants(channelSlug string) ([]*model.Product, error) {

	channelQuery := ps.channelQuery(channelSlug, model.NewBool(true), store.ProductChannelListingTableName)
	today := util.StartOfDay(timemodule.Now().UTC())

	productChannelListingQuery := ps.
		GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ProductChannelListingTableName).
		Where(squirrel.And{
			squirrel.Or{
				squirrel.LtOrEq{"ProductChannelListings.PublicationDate": today},
				squirrel.Eq{"ProductChannelListings.PublicationDate": nil},
			},
			squirrel.Eq{"ProductChannelListings.IsPublished": true},
			squirrel.Expr("ProductChannelListings.ProductID = Products.Id"),
			channelQuery,
		}).
		Suffix(")").
		Limit(1)

	channelQuery = ps.channelQuery(channelSlug, model.NewBool(true), store.ProductVariantChannelListingTableName)

	productVariantChannelListingQuery := ps.
		GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ProductVariantChannelListingTableName).
		Where(channelQuery).
		Where("ProductVariantChannelListings.PriceAmount IS NOT NULL AND ProductVariantChannelListings.VariantID = ProductVariants.Id").
		Suffix(")").
		Limit(1)

	productVariantQuery := ps.
		GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ProductVariantTableName).
		Where(productVariantChannelListingQuery).
		Where("Products.Id = ProductVariants.ProductID").
		Suffix(")").
		Limit(1)

	queryString, args, err := ps.GetQueryBuilder().
		Select(ps.ModelFields(store.ProductTableName + ".")...).
		From(store.ProductTableName).
		Where(productChannelListingQuery).
		Where(productVariantQuery).
		OrderBy(store.TableOrderingMap[store.ProductTableName]).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "PublishedWithVariants_ToSql")
	}
	var res model.Products
	err = ps.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find published with variants product with channelSlug=%s", channelSlug)
	}

	return res, nil
}

// FilterVisibleToUserProduct finds and returns all products that are visible to requesting user.
//
// 1) If requesting user is shop staff:
//
//	+) if `channelSlug` is empty string: returns all products. refer to ./product_store_doc.md (line 241, CASE 2)
//
//	+) if `channelSlug` is provided: refer to ./product_store_doc.md (line 241, CASE 1)
//
// 2) If requesting user is shop visitor: Refer to ./product_store_doc.md (line 241, case 3)
func (ps *SqlProductStore) VisibleToUserProducts(channelSlug string, requesterIsStaff bool) ([]*model.Product, error) {
	var (
		res model.Products
		err error
	)
	// check if requesting user has right to view products
	if requesterIsStaff {
		if channelSlug == "" {
			err = ps.GetReplicaX().Select(&res, "SELECT * FROM "+store.ProductTableName)
		} else {
			channelQuery := ps.channelQuery(channelSlug, nil, store.ProductChannelListingTableName)
			productChannelListingQuery := ps.
				GetQueryBuilder().
				Select(`(1) AS "a"`).
				Prefix("EXISTS (").
				From(store.ProductChannelListingTableName).
				Where(channelQuery).
				Where("ProductChannelListings.ProductID = Products.Id").
				Suffix(")").
				Limit(1)

			productQueryString, args, er := ps.
				GetQueryBuilder().
				Select(ps.ModelFields(store.ProductTableName + ".")...).
				From(store.ProductTableName).
				Where(productChannelListingQuery).
				OrderBy(store.TableOrderingMap[store.ProductTableName]).
				ToSql()
			if er != nil {
				return nil, errors.Wrap(er, "VisibleToUserProducts_ToSql") // return immediately since this is system error
			}
			err = ps.GetReplicaX().Select(&res, productQueryString, args...)
		}
	} else {
		res, err = ps.PublishedWithVariants(channelSlug)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to find visible to users products")
	}

	return res, nil
}

// SelectForUpdateDiscountedPricesOfCatalogues finds and returns product based on given ids lists.
func (ps *SqlProductStore) SelectForUpdateDiscountedPricesOfCatalogues(productIDs []string, categoryIDs []string, collectionIDs []string) ([]*model.Product, error) {
	query := ps.GetQueryBuilder().
		Select(ps.ModelFields(store.ProductTableName + ".")...).
		Distinct().
		From(store.ProductTableName).
		OrderBy(store.TableOrderingMap[store.ProductTableName])

	if len(productIDs) > 0 {
		query = query.Where(squirrel.Eq{"Products.Id": productIDs})
	}
	if len(categoryIDs) > 0 {
		query = query.Where(squirrel.Eq{"Products.CategoryID": categoryIDs})
	}
	if len(collectionIDs) > 0 {
		query = query.
			LeftJoin(store.CollectionProductRelationTableName + " ON (Products.Id = ProductCollections.ProductID)").
			Where(squirrel.Eq{"ProductCollections.CollectionID": collectionIDs})
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "SelectForUpdateDiscountedPricesOfCatalogues_ToSql")
	}

	var products model.Products
	err = ps.GetReplicaX().Select(&products, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find products by given params")
	}

	return products, nil
}

// AdvancedFilterQueryBuilder advancedly finds products, filtered using given options
func (ps *SqlProductStore) AdvancedFilterQueryBuilder(input *model.ExportProductsFilterOptions) squirrel.SelectBuilder {
	query := ps.GetQueryBuilder().
		Select(ps.ModelFields(store.ProductTableName + ".")...).
		From(store.ProductTableName).
		OrderBy("Products.CreateAt ASC")

	if strings.EqualFold(input.Scope, "all") {
		return query
	}
	if strings.EqualFold(input.Scope, "ids") {
		return query.Where(squirrel.Eq{"Products.Id": input.Ids})
	}

	var options = input.Filter
	var channelID interface{}
	if options.Channel != nil {
		channelID = *options.Channel
	}

	// parse options
	if options.IsPublished != nil {
		query = ps.filterIsPublished(query, *options.IsPublished, channelID)
	}
	if len(options.Collections) > 0 {
		query = ps.filterCollections(query, options.Collections)
	}
	if len(options.Categories) != 0 {
		query = ps.filterCategories(query, options.Categories)
	}
	if options.HasCategory != nil {
		// default to has no category
		condition := "Products.CategoryID IS NULL"

		if *options.HasCategory {
			condition = "Products.CategoryID IS NOT NULL"
		}
		query = query.Where(condition)
	}
	if options.Price != nil {
		query = ps.filterVariantPrice(query, *options.Price, channelID)
	}
	if options.MinimalPrice != nil {
		query = ps.filterMinimalPrice(query, *options.MinimalPrice, channelID)
	}
	if len(options.Attributes) > 0 {
		arg := lo.Map(
			options.Attributes,
			func(v *struct {
				Slug        string
				Values      []string
				ValuesRange *struct {
					Gte *int32
					Lte *int32
				}
				DateTime *struct {
					Gte *timemodule.Time
					Lte *timemodule.Time
				}
				Date *struct {
					Gte *timemodule.Time
					Lte *timemodule.Time
				}
				Boolean *bool
			}, _ int) *attributeFilterInput {

				return &attributeFilterInput{
					Slug:        v.Slug,
					Values:      v.Values,
					ValuesRange: v.ValuesRange,
					DateTime:    v.DateTime,
					Date:        v.Date,
					Boolean:     v.Boolean,
				}
			})

		query = ps.filterAttributes(query, arg)
	}
	if options.StockAvailability != nil {
		query = ps.filterStockAvailability(query, *options.StockAvailability, channelID)
	}
	if len(options.ProductTypes) > 0 {
		query = ps.filterProductTypes(query, options.ProductTypes)
	}
	if options.Stocks != nil {
		query = ps.filterStocks(query, *options.Stocks)
	}
	if options.GiftCard != nil {
		query = ps.filterGiftCard(query, *options.GiftCard)
	}
	if len(options.Ids) != 0 {
		query = ps.filterProductIDs(query, options.Ids)
	}
	if options.HasPreorderedVariants != nil {
		query = ps.filterHasPreorderedVariants(query, *options.HasPreorderedVariants)
	}
	if options.Search != nil {
		query = ps.filterSearch(query, *options.Search)
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
	err = ps.GetReplicaX().Select(&products, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find products with given query and conditions")
	}

	return products, nil
}
