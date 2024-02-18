package attribute

// import (
// 	"github.com/Masterminds/squirrel"
// 	"github.com/pkg/errors"
// 	"github.com/sitename/sitename/model"
// 	"github.com/sitename/sitename/store"
// 	"gorm.io/gorm"
// )

// type SqlAttributeVariantStore struct {
// 	store.Store
// }

// func NewSqlAttributeVariantStore(s store.Store) store.AttributeVariantStore {
// 	return &SqlAttributeVariantStore{s}
// }

// func (as *SqlAttributeVariantStore) Save(attributeVariant *model.AttributeVariant) (*model.AttributeVariant, error) {
// 	err := as.GetMaster().Create(attributeVariant).Error
// 	if err != nil {
// 		if as.IsUniqueConstraintError(err, []string{"AttributeID", "ProductTypeID", "attributevariants_attributeid_producttypeid_key"}) {
// 			return nil, store.NewErrInvalidInput(model.AttributeVariantTableName, "AttributeID/ProductTypeID", attributeVariant.AttributeID+"/"+attributeVariant.ProductTypeID)
// 		}
// 		return nil, errors.Wrapf(err, "failed to save attribute variant with id=%s", attributeVariant.Id)
// 	}

// 	return attributeVariant, nil
// }

// func (as *SqlAttributeVariantStore) Get(id string) (*model.AttributeVariant, error) {
// 	var res model.AttributeVariant
// 	err := as.GetReplica().First(&res, "Id = ?", id).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, store.NewErrNotFound(model.AttributeVariantTableName, id)
// 		}
// 		return nil, errors.Wrapf(err, "failed to find attribute variant with id=%s", id)
// 	}

// 	return &res, nil
// }

// func (s *SqlAttributeVariantStore) commonQueryBuilder(options *model.AttributeVariantFilterOption) squirrel.SelectBuilder {
// 	query := s.GetQueryBuilder().Select("*").From(model.AttributeVariantTableName).Where(options.Conditions)

// 	// parse option
// 	if value := options.AttributeVisibleInStoreFront; value != nil {
// 		query = query.
// 			InnerJoin(model.AttributeTableName + " ON Attributes.Id = AttributeVariants.AttributeID").
// 			Where(squirrel.Eq{model.AttributeTableName + ".VisibleInStoreFront": *value})
// 	}

// 	return query
// }

// func (as *SqlAttributeVariantStore) GetByOption(option *model.AttributeVariantFilterOption) (*model.AttributeVariant, error) {
// 	queryString, args, err := as.commonQueryBuilder(option).ToSql()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "GetByOption_ToSql")
// 	}
// 	var res model.AttributeVariant

// 	err = as.GetReplica().Raw(queryString, args...).Scan(&res).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, store.NewErrNotFound(model.AttributeVariantTableName, "")
// 		}
// 		return nil, errors.Wrap(err, "failed to find attribute variant with given options")
// 	}

// 	return &res, nil
// }

// func (s *SqlAttributeVariantStore) FilterByOptions(options *model.AttributeVariantFilterOption) ([]*model.AttributeVariant, error) {
// 	queryString, args, err := s.commonQueryBuilder(options).ToSql()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
// 	}

// 	var res []*model.AttributeVariant
// 	err = s.GetReplica().Raw(queryString, args...).Scan(&res).Error
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to find attribute variant by given options")
// 	}
// 	return res, nil
// }
