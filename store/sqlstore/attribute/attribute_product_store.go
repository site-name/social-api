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
	as := &SqlAttributeProductStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeProduct{}, store.AttributeProductTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductTypeID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("AttributeID", "ProductTypeID")
	}
	return as
}

func (as *SqlAttributeProductStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AttributeProductTableName, "AttributeID", store.AttributeTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AttributeProductTableName, "ProductTypeID", store.ProductTypeTableName, "Id", true)
}

func (as *SqlAttributeProductStore) Save(attributeProduct *attribute.AttributeProduct) (*attribute.AttributeProduct, error) {
	attributeProduct.PreSave()
	if err := attributeProduct.IsValid(); err != nil {
		return nil, err
	}

	err := as.GetMaster().Insert(attributeProduct)
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
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AttributeProductTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeProductTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find attribute product with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAttributeProductStore) GetByOption(option *attribute.AttributeProductGetOption) (*attribute.AttributeProduct, error) {

	if option == nil || !model.IsValidId(option.AttributeID) || !model.IsValidId(option.ProductTypeID) {
		return nil, store.NewErrInvalidInput(store.AttributeProductTableName, "option", option.AttributeID+"/"+option.ProductTypeID)
	}

	var attributeProduct *attribute.AttributeProduct

	err := as.GetReplica().SelectOne(
		&attributeProduct,
		"SELECT * FROM "+store.AttributeProductTableName+" WHERE (AttributeID = :AttributeID AND ProductTypeID = :ProductTypeID)",
		map[string]interface{}{
			"AttributeID":   option.AttributeID,
			"ProductTypeID": option.ProductTypeID,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeProductTableName, "")
		}
		return nil, errors.Wrapf(err, "failed to find attribute product with AttributeID = %s, ProductTypeID = %s", option.AttributeID, option.ProductTypeID)
	}

	return attributeProduct, nil
}
