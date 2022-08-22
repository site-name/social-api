package attribute

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeProductStore struct {
	store.Store
}

func NewSqlAttributeProductStore(s store.Store) store.AttributeProductStore {
	return &SqlAttributeProductStore{s}
}

func (as *SqlAttributeProductStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id", "AttributeID", "ProductTypeID", "SortOrder",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAttributeProductStore) Save(attributeProduct *attribute.AttributeProduct) (*attribute.AttributeProduct, error) {
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

func (as *SqlAttributeProductStore) Get(id string) (*attribute.AttributeProduct, error) {
	var res attribute.AttributeProduct

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AttributeProductTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeProductTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find attribute product with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAttributeProductStore) GetByOption(option *attribute.AttributeProductFilterOption) (*attribute.AttributeProduct, error) {
	query := as.GetQueryBuilder().
		Select("*").
		From(store.AttributeProductTableName)

	// parse option
	if option.AttributeID != nil {
		query = query.Where(option.AttributeID)
	}
	if option.ProductTypeID != nil {
		query = query.Where(option.ProductTypeID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var attributeProduct attribute.AttributeProduct
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
