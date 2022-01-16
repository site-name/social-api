package product

import (
	"database/sql"
	"time"
	timemodule "time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlProductStore struct {
	store.Store
}

func NewSqlProductStore(s store.Store) store.ProductStore {
	ps := &SqlProductStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Product{}, store.ProductTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductTypeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("DefaultVariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CategoryID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(product_and_discount.PRODUCT_SLUG_MAX_LENGTH).SetUnique(true)

		s.CommonSeoMaxLength(table)
	}
	return ps
}

func (ps *SqlProductStore) TableName(withField string) string {
	name := "Products"
	if withField != "" {
		name += "." + withField
	}
	return name
}

func (ps *SqlProductStore) ModelFields() []string {
	return []string{
		"Products.Id",
		"Products.ProductTypeID",
		"Products.Name",
		"Products.Slug",
		"Products.Description",
		"Products.DescriptionPlainText",
		"Products.CategoryID",
		"Products.CreateAt",
		"Products.UpdateAt",
		"Products.ChargeTaxes",
		"Products.Weight",
		"Products.WeightUnit",
		"Products.DefaultVariantID",
		"Products.Rating",
		"Products.Metadata",
		"Products.PrivateMetadata",
		"Products.SeoTitle",
		"Products.SeoDescription",
	}
}

func (ps *SqlProductStore) ScanFields(prd product_and_discount.Product) []interface{} {
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

func (ps *SqlProductStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_products_name", store.ProductTableName, "Name")
	ps.CreateIndexIfNotExists("idx_products_slug", store.ProductTableName, "Slug")
	ps.CreateIndexIfNotExists("idx_products_name_lower_textpattern", store.ProductTableName, "lower(Name) text_pattern_ops")

	ps.CommonMetaDataIndex(store.ProductTableName)
}

// Save inserts given product into database then returns it
func (ps *SqlProductStore) Save(prd *product_and_discount.Product) (*product_and_discount.Product, error) {
	prd.PreSave()
	if err := prd.IsValid(); err != nil {
		return nil, err
	}

	if err := ps.GetMaster().Insert(prd); err != nil {
		if ps.IsUniqueConstraintError(err, []string{"Name", "products_name_key", "idx_products_name_unique"}) {
			return nil, store.NewErrInvalidInput("Product", "name", prd.Name)
		}
		if ps.IsUniqueConstraintError(err, []string{"Slug", "products_slug_key", "idx_products_slug_unique"}) {
			return nil, store.NewErrInvalidInput("Product", "slug", prd.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save Product with productId=%s", prd.Id)
	}

	return prd, nil
}

func (ps *SqlProductStore) commonQueryBuilder(option *product_and_discount.ProductFilterOption) (string, []interface{}, error) {
	query := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductTableName).
		OrderBy(store.TableOrderingMap[store.ProductTableName])

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.ProductVariantID != nil {
		// decide which type of join to use (LEFT or INNER)
		var whichJoinFunc func(join string, rest ...interface{}) squirrel.SelectBuilder = query.InnerJoin

		if val, ok := option.ProductVariantID.(squirrel.Eq); ok {
			// squirrel.Eq{"": nil}
			for _, v := range val {
				if v == nil {
					whichJoinFunc = query.LeftJoin
					break
				}
			}
		}

		query = whichJoinFunc(store.ProductVariantTableName + " ON (Products.Id = ProductVariants.ProductID)").
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
func (ps *SqlProductStore) FilterByOption(option *product_and_discount.ProductFilterOption) ([]*product_and_discount.Product, error) {
	queryString, args, err := ps.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.Product
	_, err = ps.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find products by given option")
	}

	return res, nil
}

// GetByOption finds and returns 1 product that satisfies given option
func (ps *SqlProductStore) GetByOption(option *product_and_discount.ProductFilterOption) (*product_and_discount.Product, error) {
	queryString, args, err := ps.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res product_and_discount.Product
	err = ps.GetReplica().SelectOne(&res, queryString, args...)
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
			channelActiveExpr = "Channels.IsActive "
		} else {
			channelActiveExpr = "NOT Channels.IsActive "
		}
	}
	return ps.
		GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ChannelTableName).
		Where(channelActiveExpr+"AND Channels.Slug = ? AND Channels.Id = ?.ChannelID", channelSlug, compareToTable).
		Suffix(")").
		Limit(1)
}

// FilterPublishedProducts finds and returns products that belong to given channel slug and are published
//
// refer to ./product_store_doc.md (line 1)
func (ps *SqlProductStore) PublishedProducts(channelSlug string) ([]*product_and_discount.Product, error) {

	channelQuery := ps.channelQuery(channelSlug, model.NewBool(true), store.ProductChannelListingTableName)

	today := util.StartOfDay(time.Now().UTC())

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

	var res []*product_and_discount.Product
	_, err = ps.GetReplica().Select(&res, queryString, args...)
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
		product_and_discount.Product
		IsPublished     bool
		PublicationDate *timemodule.Time
	},
	error,
) {
	today := util.StartOfDay(timemodule.Now().UTC()) // start of day

	isPublishedColumnSelect := ps.GetQueryBuilder().
		Select("PCL.IsPublished").
		From(store.ProductChannelListingTableName + " AS PCL").
		InnerJoin(store.ChannelTableName + " AS C ON (PCL.ChannelID = C.Id)").
		Where(squirrel.Expr("C.Slug = ?", channelSlug)).
		Where(squirrel.Expr("PCL.ProductID = Products.Id")).
		OrderBy(store.TableOrderingMap[store.ProductChannelListingTableName]).
		Limit(1)
	isPublishedColumnSelectString, args_1, err := isPublishedColumnSelect.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "isPublishedColumnSelect_ToSql")
	}

	publicationDateColumnSelect := ps.GetQueryBuilder().
		Select("PCL.PublicationDate").
		From(store.ProductChannelListingTableName + " AS PCL").
		InnerJoin(store.ChannelTableName + " AS C ON (C.Id = PCL.ChannelID)").
		Where(squirrel.Expr("C.Slug = ?", channelSlug)).
		Where(squirrel.Expr("PCL.ProductID = Products.Id")).
		OrderBy(store.TableOrderingMap[store.ProductChannelListingTableName]).
		Limit(1)
	publicationDateColumnSelectString, args_2, err := publicationDateColumnSelect.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "publicationDateColumnSelect_ToSql")
	}

	queryString, args, err := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		Column(squirrel.Alias(isPublishedColumnSelect, "IsPublished")).
		Column(squirrel.Alias(publicationDateColumnSelect, "PublicationDate")).
		From(store.ProductTableName).
		Where(squirrel.Or{
			squirrel.And{
				squirrel.Expr(publicationDateColumnSelectString+" > ?", append(args_2, today)...),
				squirrel.Expr(isPublishedColumnSelectString, args_1...),
			},
			squirrel.Expr("NOT "+isPublishedColumnSelectString, args_1...),
			squirrel.Expr(isPublishedColumnSelectString+" IS NULL", args_1...),
		}).
		OrderBy(store.TableOrderingMap[store.ProductTableName]).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "NotPublishedProducts_ToSql")
	}

	var res []*struct {
		product_and_discount.Product
		IsPublished     bool
		PublicationDate *time.Time
	}

	_, err = ps.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find not published product with channel slug=%s", channelSlug)
	}

	return res, nil
}

// PublishedWithVariants finds and returns products.
//
// refer to ./product_store_doc.md (line 157)
func (ps *SqlProductStore) PublishedWithVariants(channelSlug string) ([]*product_and_discount.Product, error) {

	channelQuery := ps.channelQuery(channelSlug, model.NewBool(true), store.ProductChannelListingTableName)
	today := util.StartOfDay(time.Now().UTC())

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
		Where(squirrel.And{
			channelQuery,
			squirrel.NotEq{"ProductVariantChannelListings.PriceAmount": nil},
			squirrel.Expr("ProductVariantChannelListings.VariantID = ProductVariants.Id"),
		}).
		Suffix(")").
		Limit(1)

	productVariantQuery := ps.
		GetQueryBuilder().
		Select(`(1) AS "a"`).
		Prefix("EXISTS (").
		From(store.ProductVariantTableName).
		Where(squirrel.And{
			productVariantChannelListingQuery,
			squirrel.Expr("Products.Id = ProductVariants.ProductID"),
		}).
		Suffix(")").
		Limit(1)

	queryString, args, err := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductTableName).
		Where(squirrel.And{
			productChannelListingQuery,
			productVariantQuery,
		}).
		OrderBy(store.TableOrderingMap[store.ProductTableName]).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "PublishedWithVariants_ToSql")
	}
	var res []*product_and_discount.Product
	_, err = ps.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find published with variants product with channelSlug=%s", channelSlug)
	}

	return res, nil
}

// FilterVisibleToUserProduct finds and returns all products that are visible to requesting user.
//
// 1) If requesting user is shop staff:
//
// 	+) if `channelSlug` is empty string: returns all products. refer to ./product_store_doc.md (line 241, CASE 2)
//
// 	+) if `channelSlug` is provided: refer to ./product_store_doc.md (line 241, CASE 1)
//
// 2) If requesting user is shop visitor: Refer to ./product_store_doc.md (line 241, case 3)
func (ps *SqlProductStore) VisibleToUserProducts(channelSlug string, requesterIsStaff bool) ([]*product_and_discount.Product, error) {
	var (
		res []*product_and_discount.Product
		err error
	)
	// check if requesting user has right to view products
	if requesterIsStaff {
		if channelSlug == "" {
			_, err = ps.GetReplica().Select(&res, "SELECT * FROM "+store.ProductTableName)
		} else {
			channelQuery := ps.channelQuery(channelSlug, nil, store.ProductChannelListingTableName)
			productChannelListingQuery := ps.
				GetQueryBuilder().
				Select(`(1) AS "a"`).
				Prefix("EXISTS (").
				From(store.ProductChannelListingTableName).
				Where(squirrel.And{
					squirrel.Expr("ProductChannelListings.ProductID = Products.Id"),
					channelQuery,
				}).
				Suffix(")").
				Limit(1)

			productQueryString, args, er := ps.
				GetQueryBuilder().
				Select(ps.ModelFields()...).
				From(store.ProductTableName).
				Where(productChannelListingQuery).
				OrderBy(store.TableOrderingMap[store.ProductTableName]).
				ToSql()
			if er != nil {
				return nil, errors.Wrap(er, "VisibleToUserProducts_ToSql") // return immediately since this is system error
			}
			_, err = ps.GetReplica().Select(&res, productQueryString, args...)
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
func (ps *SqlProductStore) SelectForUpdateDiscountedPricesOfCatalogues(productIDs []string, categoryIDs []string, collectionIDs []string) ([]*product_and_discount.Product, error) {
	query := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		Distinct().
		From(store.ProductTableName).
		OrderBy(store.TableOrderingMap[store.ProductTableName])

	if len(productIDs) > 0 {
		query = query.Where("Products.Id IN ?", productIDs)
	}
	if len(categoryIDs) > 0 {
		query = query.Where("OR Products.CategoryID IN ?", categoryIDs)
	}
	if len(collectionIDs) > 0 {
		query = query.
			LeftJoin(store.CollectionProductRelationTableName+" ON (Products.Id = ProductCollections.ProductID)").
			Where("OR ProductCollections.CollectionID IN ?", collectionIDs)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "SelectForUpdateDiscountedPricesOfCatalogues_ToSql")
	}

	var products []*product_and_discount.Product
	_, err = ps.GetReplica().Select(&products, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find products by given params")
	}

	return products, nil
}
