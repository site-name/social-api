package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlAttributeProductStore struct {
	store.Store
}

func NewSqlAttributeProductStore(s store.Store) store.AttributeProductStore {
	return &SqlAttributeProductStore{s}
}

func (as *SqlAttributeProductStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
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
	attributeProduct.PreSave()
	if err := attributeProduct.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.AttributeProductTableName + "(" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"

	_, err := as.GetMasterX().NamedExec(query, attributeProduct)
	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"attributeproducts_attributeid_producttypeid_key", "AttributeID", "ProductTypeID"}) {
			return nil, store.NewErrInvalidInput(store.AttributeProductTableName, "AttributeID/ProductTypeID", attributeProduct.AttributeID+"/"+attributeProduct.ProductTypeID)
		}
		return nil, errors.Wrapf(err, "failed to save new attributeProduct with id=%s", attributeProduct.Id)
	}

	return attributeProduct, nil
}

func (as *SqlAttributeProductStore) Get(id string) (*model.AttributeProduct, error) {
	var res model.AttributeProduct

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AttributeProductTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeProductTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find attribute product with id=%s", id)
	}

	return &res, nil
}

func (s *SqlAttributeProductStore) commonQueryBuilder(option *model.AttributeProductFilterOption) squirrel.SelectBuilder {
	query := s.GetQueryBuilder().
		Select("*").
		From(store.AttributeProductTableName)

	// parse option
	if option.AttributeID != nil {
		query = query.Where(option.AttributeID)
	}
	if option.ProductTypeID != nil {
		query = query.Where(option.ProductTypeID)
	}
	if option.AttributeVisibleInStoreFront != nil {
		query = query.
			InnerJoin(store.AttributeTableName + " ON Attributes.Id = AttributeProducts.AttributeID").
			Where(squirrel.Eq{store.AttributeTableName + ".VisibleInStoreFront": *option.AttributeVisibleInStoreFront})
	}
	return query
}

func (as *SqlAttributeProductStore) GetByOption(option *model.AttributeProductFilterOption) (*model.AttributeProduct, error) {
	queryString, args, err := as.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var attributeProduct model.AttributeProduct
	err = as.GetReplicaX().Get(
		&attributeProduct,
		queryString,
		args...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeProductTableName, "")
		}
		return nil, errors.Wrapf(err, "failed to find attribute product with AttributeID = %s, ProductTypeID = %s", option.AttributeID, option.ProductTypeID)
	}

	return &attributeProduct, nil
}

func (s *SqlAttributeProductStore) FilterByOptions(option *model.AttributeProductFilterOption) ([]*model.AttributeProduct, error) {
	queryString, args, err := s.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res []*model.AttributeProduct
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find attribute products by given options")
	}

	return res, nil
}
