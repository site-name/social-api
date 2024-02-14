package product

import (
	"fmt"

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
		InnerJoin(fmt.Sprintf("%[1]s ON %[2]s.Id = %[1]s.ProductTypeID", model.ProductTableName, model.ProductTypeTableName)).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.ProductID = %[2]s.Id", model.ProductVariantTableName, model.ProductTableName)).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.VariantID = %[2]s.Id", model.CheckoutLineTableName, model.ProductVariantTableName)).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Token = %[2]s.CheckoutID", model.CheckoutTableName, model.CheckoutLineTableName)).
		Where(model.CheckoutTableName+".Token = ?", checkoutToken).
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
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.ProductTypeID = %[2]s.Id", model.ProductTableName, model.ProductTypeTableName)).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.ProductID = %[2]s.Id", model.ProductVariantTableName, model.ProductTableName)).
		Where(model.ProductVariantTableName+".Id = ?", variantID)

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
			InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.ProductTypeID = %[2]s.Id", model.AttributeProductTableName, model.ProductTypeTableName)).
			Where(options.AttributeProducts_AttributeID)
	}
	if options.AttributeVariants_AttributeID != nil {
		query = query.
			InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.ProductTypeID = %[2]s.Id", model.AttributeVariantTableName, model.ProductTypeTableName)).
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
func (pts *SqlProductTypeStore) FilterbyOption(options *model.ProductTypeFilterOption) (int64, []*model.ProductType, error) {
	query := pts.commonQueryBuilder(options)

	// count if needed
	var totalCount int64
	if options.CountTotal {
		countQuery, args, err := pts.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "FilterByOption_CountTotal_ToSql")
		}
		err = pts.GetReplica().Raw(countQuery, args...).Scan(totalCount).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total product types by given options")
		}
	}

	// apply pagination if needed
	options.GraphqlPaginationValues.AddPaginationToSelectBuilderIfNeeded(&query)

	queryString, args, err := query.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var res []*model.ProductType
	err = pts.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to find product types with given option")
	}

	return totalCount, res, nil
}

// func (pts *SqlProductTypeStore) Count(options *model.ProductTypeFilterOption) (int64, error) {
// 	countQuery := pts.commonQueryBuilder(options)

// 	queryStr, args, err := pts.GetQueryBuilder().Select("COUNT(*)").FromSelect(countQuery, "c").ToSql()
// 	if err != nil {
// 		return 0, errors.Wrap(err, "Count_ToSql")
// 	}

// 	var count int64
// 	err = pts.GetReplica().Raw(queryStr, args...).Scan(&count).Error
// 	if err != nil {
// 		return 0, errors.Wrap(err, "failed to count number of product types by options")
// 	}

// 	return count, nil
// }

func (s *SqlProductTypeStore) Delete(tx *gorm.DB, ids []string) (int64, error) {
	if tx == nil {
		tx = s.GetMaster()
	}

	// delete attribute values that are in relations with given product types
	attributeValueQuery, args, err := s.GetQueryBuilder().
		Select(model.AttributeValueTableName + ".Id").
		From(model.AttributeValueTableName).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Id = %[2]s.AttributeID", model.AttributeTableName, model.AttributeValueTableName)).
		//
		InnerJoin(fmt.Sprintf("%[1]s ON %[2]s.Id = %[1]s.ValueID", model.AssignedProductAttributeValueTableName, model.AttributeValueTableName)).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Id = %[2]s.AssignmentID", model.AssignedProductAttributeTableName, model.AssignedProductAttributeValueTableName)).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Id = %[2]s.AssignmentID", model.AttributeProductTableName, model.AssignedProductAttributeTableName)).
		//
		InnerJoin(fmt.Sprintf("%[1]s ON %[2]s.Id = %[1]s.ValueID", model.AssignedVariantAttributeValueTableName, model.AttributeValueTableName)).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Id = %[2]s.AssignmentID", model.AssignedVariantAttributeTableName, model.AssignedVariantAttributeValueTableName)).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Id = %[2]s.AssignmentID", model.AttributeVariantTableName, model.AssignedVariantAttributeTableName)).
		//
		Where(squirrel.Eq{model.AttributeTableName + ".InputType": model.TYPES_WITH_UNIQUE_VALUES}).
		Where(squirrel.Or{
			squirrel.Eq{model.AttributeVariantTableName + ".ProductTypeID": ids},
			squirrel.Eq{model.AttributeProductTableName + ".ProductTypeID": ids},
		}).
		ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Delete_AttributeValues_ToSql")
	}

	var attributeValueIDs []string
	err = s.GetReplica().Raw(attributeValueQuery, args...).Scan(&attributeValueIDs).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to find related attribute values of given product types")
	}

	_, err = s.AttributeValue().Delete(tx, attributeValueIDs...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete related attribute values of given product types")
	}

	// delete relate DRAFT order lines

	orderLineQuery, args, err := s.GetQueryBuilder().
		Select(model.OrderLineTableName+".Id").
		From(model.OrderLineTableName).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Id = %[2]s.VariantID", model.ProductVariantTableName, model.OrderLineTableName)).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Id = %[2]s.OrderID", model.OrderTableName, model.OrderLineTableName)).
		InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Id = %[2]s.ProductID", model.ProductTableName, model.ProductVariantTableName)).
		Where(model.OrderTableName+".Status = ?", model.ORDER_STATUS_DRAFT).
		Where(squirrel.Eq{model.ProductTableName + ".ProductTypeID": ids}).
		ToSql()

	if err != nil {
		return 0, errors.Wrap(err, "Delete_OrderLine_ToSql")
	}

	var orderLineIDs []string
	err = s.GetReplica().Raw(orderLineQuery, args...).Scan(&orderLineIDs).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to find related order line ids of given product types")
	}

	err = s.OrderLine().BulkDelete(tx, orderLineIDs)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete related draft order lines of given product types")
	}

	// delete product types
	delRes := tx.Raw("DELETE FROM "+model.ProductTypeTableName+" WHERE Id IN ?", ids)
	if delRes.Error != nil {
		return 0, errors.Wrap(delRes.Error, "failed to delete product types by given ids")
	}

	return delRes.RowsAffected, nil
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
