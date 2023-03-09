package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlProductTypeStore struct {
	store.Store
}

func NewSqlProductTypeStore(s store.Store) store.ProductTypeStore {
	return &SqlProductTypeStore{s}
}

func (ps *SqlProductTypeStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Name",
		"Slug",
		"Kind",
		"HasVariants",
		"IsShippingRequired",
		"IsDigital",
		"Weight",
		"WeightUnit",
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

func (ps *SqlProductTypeStore) ScanFields(productType *model.ProductType) []interface{} {
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

func (ps *SqlProductTypeStore) Save(productType *model.ProductType) (*model.ProductType, error) {
	productType.PreSave()
	if err := productType.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.ProductTypeTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
	if _, err := ps.GetMasterX().NamedExec(query, productType); err != nil {
		return nil, errors.Wrapf(err, "failed to save product type withh id=%s", productType.Id)
	}

	return productType, nil
}

func (ps *SqlProductTypeStore) FilterProductTypesByCheckoutToken(checkoutToken string) ([]*model.ProductType, error) {
	/*
					checkout
					|      |
		...<--|		   |--> checkoutLine <-- productVariant <-- product <-- productType
																							|												     |
													 ...checkoutLine <--|              ...product <--|
	*/

	queryString, args, err := ps.GetQueryBuilder().
		Select(ps.ModelFields(store.ProductTypeTableName+".")...).
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

	var productTypes []*model.ProductType

	err = ps.GetReplicaX().Select(&productTypes, queryString, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find product types related to checkout with token=%s", checkoutToken)
	}
	return productTypes, nil
}

func (pts *SqlProductTypeStore) ProductTypesByProductIDs(productIDs []string) ([]*model.ProductType, error) {
	var productTypes []*model.ProductType
	queryString, args, err := pts.
		GetQueryBuilder().
		Select("pt.*").
		From(store.ProductTypeTableName + " pt").
		InnerJoin(store.ProductTableName + " p ON pt.Id = p.ProductTypeID").
		Where(squirrel.Eq{"p.Id": productIDs}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "ProductTypesByProductIDs_ToSql")
	}

	err = pts.GetReplicaX().Select(&productTypes, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product types with given product ids")
	}

	return productTypes, nil
}

// ProductTypeByProductVariantID finds and returns 1 product type that is related to given product variant
func (pts *SqlProductTypeStore) ProductTypeByProductVariantID(variantID string) (*model.ProductType, error) {
	query := pts.GetQueryBuilder().
		Select(pts.ModelFields(store.ProductTypeTableName+".")...).
		From(store.ProductTypeTableName).
		OrderBy(store.TableOrderingMap[store.ProductTypeTableName]).
		InnerJoin(store.ProductTableName+" ON (Products.ProductTypeID = ProductTypes.Id)").
		InnerJoin(store.ProductVariantTableName+" ON (Products.Id = ProductVariants.ProductID)").
		Where("ProductVariants.Id = ?", variantID)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "ProductTypeByProductVariantID_ToSql")
	}

	var res model.ProductType
	err = pts.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductTypeTableName, "variantID="+variantID)
		}
		return nil, errors.Wrapf(err, "failed to find product type with product variant id=%s", variantID)
	}

	return &res, nil
}

func (pts *SqlProductTypeStore) commonQueryBuilder(options *model.ProductTypeFilterOption) squirrel.SelectBuilder {
	query := pts.GetQueryBuilder().
		Select(pts.ModelFields(store.ProductTypeTableName + ".")...).
		From(store.ProductTypeTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Name != nil {
		query = query.Where(options.Name)
	}
	if options.AttributeProducts_AttributeID != nil {
		query = query.InnerJoin(store.AttributeProductTableName + " ON AttributeProducts.ProductTypeID = ProductTypes.Id").
			Where(options.AttributeProducts_AttributeID)
	}
	if options.AttributeVariants_AttributeID != nil {
		query = query.InnerJoin(store.AttributeVariantTableName + " ON AttributeVariants.ProductTypeID = ProductTypes.Id").
			Where(options.AttributeVariants_AttributeID)
	}
	if options.Extra != nil {
		return query.Where(options.Extra)
	}

	return query
}

// GetByOption finds and returns a product type with given options
func (pts *SqlProductTypeStore) GetByOption(options *model.ProductTypeFilterOption) (*model.ProductType, error) {
	queryString, args, err := pts.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res model.ProductType
	err = pts.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductTypeTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find product type with given options")
	}

	return &res, nil
}

// FilterbyOption finds and returns a slice of product types filtered using given options
func (pts *SqlProductTypeStore) FilterbyOption(options *model.ProductTypeFilterOption) ([]*model.ProductType, error) {
	queryString, args, err := pts.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var res []*model.ProductType
	err = pts.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product types with given option")
	}

	return res, nil
}

func (pts *SqlProductTypeStore) Count(options *model.ProductTypeFilterOption) (int64, error) {
	countQuery := pts.commonQueryBuilder(options)

	queryStr, args, err := pts.GetQueryBuilder().Select("COUNT(*)").FromSelect(countQuery, "c").ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Count_ToSql")
	}

	var count int64
	err = pts.GetReplicaX().Get(&count, queryStr, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of product types by options")
	}

	return count, nil
}
