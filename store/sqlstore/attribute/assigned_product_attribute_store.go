package attribute

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAssignedProductAttributeStore struct {
	store.Store
}

func (s *SqlAssignedProductAttributeStore) ScanFields(prdAttr *model.AssignedProductAttribute) []any {
	return []any{
		&prdAttr.Id,
		&prdAttr.ProductID,
		&prdAttr.AssignmentID,
	}
}

func NewSqlAssignedProductAttributeStore(s store.Store) store.AssignedProductAttributeStore {
	return &SqlAssignedProductAttributeStore{s}
}

func (as *SqlAssignedProductAttributeStore) Save(newInstance *model.AssignedProductAttribute) (*model.AssignedProductAttribute, error) {
	if err := as.GetMaster().Save(newInstance).Error; err != nil {
		if as.IsUniqueConstraintError(err, []string{"ProductID", "AssignmentID", strings.ToLower(model.AssignedProductAttributeTableName) + "_productid_assignmentid_key"}) {
			return nil, store.NewErrInvalidInput(model.AssignedProductAttributeTableName, "ProductID/AssignmentID", newInstance.ProductID+"/"+newInstance.AssignmentID)
		}
		return nil, errors.Wrap(err, "failed to insert new assigned product attribute with")
	}

	return newInstance, nil
}

func (as *SqlAssignedProductAttributeStore) Get(id string) (*model.AssignedProductAttribute, error) {
	var res model.AssignedProductAttribute

	err := as.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedProductAttributeTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAssignedProductAttributeStore) commonQueryBuilder(options *model.AssignedProductAttributeFilterOption) squirrel.SelectBuilder {
	selectFields := []string{model.AssignedProductAttributeTableName + ".*"}

	query := as.GetQueryBuilder().
		Select(selectFields...).
		From(model.AssignedProductAttributeTableName).
		Where(options.Conditions)

	if options.AttributeProduct_Attribute_VisibleInStoreFront != nil {
		query = query.
			InnerJoin(model.AttributeProductTableName + " ON AttributeProducts.Id = AssignedProductAttributes.AssignmentID").
			InnerJoin(model.AttributeTableName + " ON AttributeProducts.AttributeID = Attributes.Id").
			Where(squirrel.Eq{model.AttributeTableName + ".VisibleInStoreFront": *options.AttributeProduct_Attribute_VisibleInStoreFront})
	}

	return query
}

func (as *SqlAssignedProductAttributeStore) GetWithOption(option *model.AssignedProductAttributeFilterOption) (*model.AssignedProductAttribute, error) {
	queryString, args, err := as.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetWithOption_ToSql")
	}

	var res model.AssignedProductAttribute
	db := as.GetReplica()

	for _, preload := range option.Preloads {
		db = db.Preload(preload)
	}

	err = db.Raw(queryString, args...).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedProductAttributeTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute with given options")
	}

	return &res, nil
}

func (as *SqlAssignedProductAttributeStore) FilterByOptions(options *model.AssignedProductAttributeFilterOption) ([]*model.AssignedProductAttribute, error) {
	queryString, args, err := as.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res model.AssignedProductAttributes
	db := as.GetReplica()
	for _, preload := range options.Preloads {
		db = db.Preload(preload)
	}
	err = db.Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find assigned product attributes with given options")
	}

	return res, nil
}
