package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductTypeStore struct {
	store.Store
}

func NewSqlProductTypeStore(s store.Store) store.ProductTypeStore {
	pts := &SqlProductTypeStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductType{}, store.ProductTypeTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_TYPE_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(product_and_discount.PRODUCT_TYPE_SLUG_MAX_LENGTH)
		table.ColMap("Kind").SetMaxSize(product_and_discount.PRODUCT_TYPE_KIND_MAX_LENGTH)
	}
	return pts
}

func (ps *SqlProductTypeStore) ModelFields() []string {
	return []string{
		"ProductTypes.Id",
		"ProductTypes.Name",
		"ProductTypes.Slug",
		"ProductTypes.Kind",
		"ProductTypes.HasVariants",
		"ProductTypes.IsShippingRequired",
		"ProductTypes.IsDigital",
		"ProductTypes.Weight",
		"ProductTypes.WeightUnit",
		"ProductTypes.Metadata",
		"ProductTypes.PrivateMetadata",
	}
}

func (ps *SqlProductTypeStore) ScanFields(productType product_and_discount.ProductType) []interface{} {
	return []interface{}{
		&productType.Id,
		&productType.Name,
		&productType.Slug,
		&productType.Kind,
		&productType.HasVariants,
		&productType.IsShippingRequired,
		&productType.IsDigital,
		&productType.Weight,
		&productType.WeightUnit,
		&productType.Metadata,
		&productType.PrivateMetadata,
	}
}

func (ps *SqlProductTypeStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_product_types_name", store.ProductTypeTableName, "Name")
	ps.CreateIndexIfNotExists("idx_product_types_name_lower_textpattern", store.ProductTypeTableName, "lower(Name) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_product_types_slug", store.ProductTypeTableName, "Slug")
}

func (ps *SqlProductTypeStore) Save(productType *product_and_discount.ProductType) (*product_and_discount.ProductType, error) {
	productType.PreSave()
	if err := productType.IsValid(); err != nil {
		return nil, err
	}

	if err := ps.GetMaster().Insert(productType); err != nil {
		return nil, errors.Wrapf(err, "failed to save product type withh id=%s", productType.Id)
	}

	return productType, nil
}

func (ps *SqlProductTypeStore) FilterProductTypesByCheckoutToken(checkoutToken string) ([]*product_and_discount.ProductType, error) {
	/*
					checkout
					|      |
		...<--|		   |--> checkoutLine <-- productVariant <-- product <-- productType
																							|												     |
													 ...checkoutLine <--|              ...product <--|
	*/

	queryString, args, err := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductTypeTableName).
		InnerJoin(store.ProductTableName+" ON (ProductTypes.Id = Products.ProductTypeID)").
		InnerJoin(store.ProductVariantTableName+" ON (ProductVariants.ProductID = Products.Id)").
		InnerJoin(store.CheckoutLineTableName+" ON (CheckoutLines.VariantID = ProductVariants.Id)").
		InnerJoin(store.CheckoutTableName+" ON (Checkouts.Token = CheckoutLines.CheckoutID)").
		Where("Checkouts.Token = ?", checkoutToken).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "FilterProductTypesByCheckoutToken_ToSql")
	}

	var productTypes []*product_and_discount.ProductType

	_, err = ps.GetReplica().Select(&productTypes, queryString, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find product types related to checkout with token=%s", checkoutToken)
	}
	return productTypes, nil
}

func (pts *SqlProductTypeStore) ProductTypesByProductIDs(productIDs []string) ([]*product_and_discount.ProductType, error) {
	var productTypes []*product_and_discount.ProductType
	_, err := pts.GetReplica().Select(
		&productTypes,
		`SELECT * FROM `+store.ProductTypeTableName+` 
		INNER JOIN `+store.ProductTableName+` ON (
			ProductTypes.Id = Products.ProductTypeID
		) 
		WHERE Products.Id IN :IDs`,
		map[string]interface{}{
			"IDs": productIDs,
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to find product types with given product ids")
	}

	return productTypes, nil
}

// ProductTypeByProductVariantID finds and returns 1 product type that is related to given product variant
func (pts *SqlProductTypeStore) ProductTypeByProductVariantID(variantID string) (*product_and_discount.ProductType, error) {
	query := pts.GetQueryBuilder().
		Select(pts.ModelFields()...).
		From(store.ProductTypeTableName).
		OrderBy(store.TableOrderingMap[store.ProductTypeTableName]).
		InnerJoin(store.ProductTableName+" ON (Products.ProductTypeID = ProductTypes.Id)").
		InnerJoin(store.ProductVariantTableName+" ON (Products.Id = ProductVariants.ProductID)").
		Where("ProductVariants.Id = ?", variantID)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "ProductTypeByProductVariantID_ToSql")
	}

	var res product_and_discount.ProductType
	err = pts.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductTypeTableName, "variantID="+variantID)
		}
		return nil, errors.Wrapf(err, "failed to find product type with product variant id=%s", variantID)
	}

	return &res, nil
}

func (pts *SqlProductTypeStore) commonQueryBuilder(options *product_and_discount.ProductTypeFilterOption) squirrel.SelectBuilder {
	query := pts.GetQueryBuilder().
		Select(pts.ModelFields()...).
		From(store.ProductTypeTableName).
		OrderBy(store.TableOrderingMap[store.ProductTypeTableName])

	// parse options
	if options.Limit > 0 {
		query = query.Limit(uint64(options.Limit))
	}
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Name != nil {
		query = query.Where(options.Name)
	}
	if options.AttributeID != nil {
		query = query.InnerJoin(store.AttributeProductTableName + " ON AttributeProducts.ProductTypeID = ProductTypes.Id").
			Where(options.AttributeID)
	}
	if options.Extra != nil {
		query = query.Where(options.Extra)
	}

	return query
}

// GetByOption finds and returns a product type with given options
func (pts *SqlProductTypeStore) GetByOption(options *product_and_discount.ProductTypeFilterOption) (*product_and_discount.ProductType, error) {
	queryString, args, err := pts.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res product_and_discount.ProductType
	err = pts.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductTypeTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find product type with given options")
	}

	return &res, nil
}

// FilterbyOption finds and returns a slice of product types filtered using given options
func (pts *SqlProductTypeStore) FilterbyOption(options *product_and_discount.ProductTypeFilterOption) ([]*product_and_discount.ProductType, error) {
	queryString, args, err := pts.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var res []*product_and_discount.ProductType
	_, err = pts.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product types with given option")
	}

	return res, nil
}

func (pts *SqlProductTypeStore) Count(options *product_and_discount.ProductTypeFilterOption) (int64, error) {
	options.Limit = 0 // unset limit

	query := pts.commonQueryBuilder(options)

	queryStr, args, err := pts.GetQueryBuilder().Select("COUNT(*)").FromSelect(query, "c").ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Count_ToSql")
	}

	count, err := pts.GetReplica().SelectInt(queryStr, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of product types by options")
	}

	return count, nil
}
