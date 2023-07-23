package attribute

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAttributeProductStore struct {
	store.Store
}

func NewSqlAttributeProductStore(s store.Store) store.AttributeProductStore {
	return &SqlAttributeProductStore{s}
}

func (as *SqlAttributeProductStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id", "AttributeID", "ProductTypeID", "SortOrder",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAttributeProductStore) Save(attributeProduct *model.AttributeProduct) (*model.AttributeProduct, error) {
	err := as.GetMaster().Create(attributeProduct).Error
	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"attributeproducts_attributeid_producttypeid_key", "AttributeID", "ProductTypeID"}) {
			return nil, store.NewErrInvalidInput(model.AttributeProductTableName, "AttributeID/ProductTypeID", attributeProduct.AttributeID+"/"+attributeProduct.ProductTypeID)
		}
		return nil, errors.Wrap(err, "failed to save new attributeProduct")
	}

	return attributeProduct, nil
}

func (as *SqlAttributeProductStore) Get(id string) (*model.AttributeProduct, error) {
	var res model.AttributeProduct
	err := as.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AttributeProductTableName, id)
		}
		return nil, errors.Wrap(err, "failed to find attribute product by given id")
	}

	return &res, nil
}

func (s *SqlAttributeProductStore) commonQueryBuilder(option *model.AttributeProductFilterOption) squirrel.SelectBuilder {
	query := s.GetQueryBuilder().
		Select(s.ModelFields(model.AttributeProductTableName + ".")...).
		From(model.AttributeProductTableName).
		Where(option.Conditions)

	// parse option
	if option.AttributeVisibleInStoreFront != nil {
		query = query.
			InnerJoin(model.AttributeTableName + " ON Attributes.Id = AttributeProducts.AttributeID").
			Where(squirrel.Eq{model.AttributeTableName + ".VisibleInStoreFront": *option.AttributeVisibleInStoreFront})
	}
	return query
}

func (as *SqlAttributeProductStore) GetByOption(option *model.AttributeProductFilterOption) (*model.AttributeProduct, error) {
	queryString, args, err := as.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var attributeProduct model.AttributeProduct
	err = as.GetReplica().Raw(queryString, args...).Scan(&attributeProduct).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AttributeProductTableName, "")
		}
		return nil, errors.Wrapf(err, "failed to find attribute product with given options")
	}

	return &attributeProduct, nil
}

func (s *SqlAttributeProductStore) FilterByOptions(option *model.AttributeProductFilterOption) ([]*model.AttributeProduct, error) {
	queryString, args, err := s.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res []*model.AttributeProduct
	err = s.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find attribute products by given options")
	}

	return res, nil
}
