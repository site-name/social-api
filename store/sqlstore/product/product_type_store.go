package product

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlProductTypeStore struct {
	store.Store
}

func NewSqlProductTypeStore(s store.Store) store.ProductTypeStore {
	return &SqlProductTypeStore{s}
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

func (ps *SqlProductTypeStore) Save(tx *gorm.DB, productType *model.ProductType) (*model.ProductType, error) {
	if tx == nil {
		tx = ps.GetMaster()
	}
	if err := tx.Save(productType).Error; err != nil {
		if ps.IsUniqueConstraintError(err, []string{"slug_key", "slug"}) {
			return nil, store.NewErrInvalidInput(model.ProductTypeTableName, "slug", productType.Slug)
		}
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
		Select(model.ProductTypeTableName+".*").
		From(model.ProductTypeTableName).
		InnerJoin(model.ProductTableName+" ON (ProductTypes.Id = Products.ProductTypeID)").
		InnerJoin(model.ProductVariantTableName+" ON (ProductVariants.ProductID = Products.Id)").
		InnerJoin(model.CheckoutLineTableName+" ON (CheckoutLines.VariantID = ProductVariants.Id)").
		InnerJoin(model.CheckoutTableName+" ON (Checkouts.Token = CheckoutLines.CheckoutID)").
		Where("Checkouts.Token = ?", checkoutToken).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "FilterProductTypesByCheckoutToken_ToSql")
	}

	var productTypes []*model.ProductType

	err = ps.GetReplica().Raw(queryString, args...).Scan(&productTypes).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find product types related to checkout with token=%s", checkoutToken)
	}
	return productTypes, nil
}

func (pts *SqlProductTypeStore) ProductTypesByProductIDs(productIDs []string) ([]*model.ProductType, error) {
	var productTypes []*model.ProductType

	err := pts.GetReplica().Raw("SELECT * FROM "+model.ProductTypeTableName+" INNER JOIN "+model.ProductTableName+" ON ProductTypes.Id = Products.ProductTypeID WHERE Products.Id IN ?", productIDs).Scan(&productTypes).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product types with given product ids")
	}

	return productTypes, nil
}

// ProductTypeByProductVariantID finds and returns 1 product type that is related to given product variant
func (pts *SqlProductTypeStore) ProductTypeByProductVariantID(variantID string) (*model.ProductType, error) {
	query := pts.GetQueryBuilder().
		Select(model.ProductTypeTableName+".*").
		From(model.ProductTypeTableName).
		InnerJoin(model.ProductTableName+" ON (Products.ProductTypeID = ProductTypes.Id)").
		InnerJoin(model.ProductVariantTableName+" ON (Products.Id = ProductVariants.ProductID)").
		Where("ProductVariants.Id = ?", variantID)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "ProductTypeByProductVariantID_ToSql")
	}

	var res model.ProductType
	err = pts.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ProductTypeTableName, "variantID="+variantID)
		}
		return nil, errors.Wrapf(err, "failed to find product type with product variant id=%s", variantID)
	}

	return &res, nil
}

func (pts *SqlProductTypeStore) commonQueryBuilder(options *model.ProductTypeFilterOption) squirrel.SelectBuilder {
	query := pts.GetQueryBuilder().
		Select(model.ProductTypeTableName + ".*").
		From(model.ProductTypeTableName).
		Where(options.Conditions)

	if options.AttributeProducts_AttributeID != nil {
		query = query.
			InnerJoin(model.AttributeProductTableName + " ON AttributeProducts.ProductTypeID = ProductTypes.Id").
			Where(options.AttributeProducts_AttributeID)
	}
	if options.AttributeVariants_AttributeID != nil {
		query = query.
			InnerJoin(model.AttributeVariantTableName + " ON AttributeVariants.ProductTypeID = ProductTypes.Id").
			Where(options.AttributeVariants_AttributeID)
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
	err = pts.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ProductTypeTableName, "options")
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
	err = pts.GetReplica().Raw(queryString, args...).Scan(&res).Error
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
	err = pts.GetReplica().Raw(queryStr, args...).Scan(&count).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of product types by options")
	}

	return count, nil
}

func (s *SqlProductTypeStore) Delete(tx *gorm.DB, ids []string) (int64, error) {
	if tx == nil {
		tx = s.GetMaster()
	}

	res := tx.Where("Id IN ?", ids).Delete(&model.ProductType{})
	if res.Error != nil {
		return 0, errors.Wrap(res.Error, "failed to delete product types")
	}
	return res.RowsAffected, nil
}

func (s *SqlProductTypeStore) ToggleProductTypeRelations(tx *gorm.DB, productTypeID string, productAttributes, variantAttributes model.Attributes, isDelete bool) error {
	if tx == nil {
		tx = s.GetMaster()
	}

	relationsMap := map[string]model.Attributes{
		"ProductAttributes": productAttributes,
		"VariantAttributes": variantAttributes,
	}

	for assocName, relations := range relationsMap {
		if len(relations) > 0 {
			var err error
			if isDelete {
				err = tx.Model(&model.ProductType{Id: productTypeID}).Association(assocName).Delete(relations)
			} else {
				err = tx.Model(&model.ProductType{Id: productTypeID}).Association(assocName).Append(relations)
			}

			if err != nil {
				action := "add"
				if isDelete {
					action = "delete"
				}
				return errors.Wrapf(err, "failed to %s %s relations", action, assocName)
			}
		}
	}

	return nil
}
