package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeVariantStore struct {
	store.Store
}

func NewSqlAttributeVariantStore(s store.Store) store.AttributeVariantStore {
	as := &SqlAttributeVariantStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeVariant{}, store.AttributeVariantTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductTypeID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("AttributeID", "ProductTypeID")
	}
	return as
}

func (as *SqlAttributeVariantStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AttributeVariantTableName, "AttributeID", store.AttributeTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AttributeVariantTableName, "ProductTypeID", store.ProductTypeTableName, "Id", true)
}

func (as *SqlAttributeVariantStore) Save(attributeVariant *attribute.AttributeVariant) (*attribute.AttributeVariant, error) {
	attributeVariant.PreSave()
	if err := attributeVariant.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(attributeVariant); err != nil {
		if as.IsUniqueConstraintError(err, []string{"AttributeID", "ProductTypeID", strings.ToLower(store.AttributeVariantTableName) + "_attributeid_producttypeid_key"}) {
			return nil, store.NewErrInvalidInput(store.AttributeVariantTableName, "AttributeID/ProductTypeID", attributeVariant.AttributeID+"/"+attributeVariant.ProductTypeID)
		}
		return nil, errors.Wrapf(err, "failed to save attribute variant with id=%s", attributeVariant.Id)
	}

	return attributeVariant, nil
}

func (as *SqlAttributeVariantStore) Get(attributeVariantID string) (*attribute.AttributeVariant, error) {

	res, err := as.GetReplica().Get(attribute.AttributeVariant{}, attributeVariantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeVariantTableName, attributeVariantID)
		}
		return nil, errors.Wrapf(err, "failed to find attribute variant with id=%s", attributeVariantID)
	}

	return res.(*attribute.AttributeVariant), nil
}

func (as *SqlAttributeVariantStore) GetByOption(option *attribute.AttributeVariantFilterOption) (*attribute.AttributeVariant, error) {
	if option == nil || option.AttributeID == "" || option.ProductID == "" {
		return nil, store.NewErrInvalidInput(store.AttributeVariantTableName, "option", "")
	}

	var res *attribute.AttributeVariant
	err := as.GetReplica().SelectOne(
		&res,
		`SELECT * FROM `+store.AttributeVariantTableName+` AS av
		INNER JOIN `+store.ProductTypeTableName+` AS pdt ON (
			pdt.Id = av.ProductTypeID
		) 
		INNER JOIN `+store.ProductTableName+` AS pd ON (
			pd.ProductTypeID = pdt.Id
		)
		WHERE (
			av.AttributeID = :AttributeID AND pd.Id = :ProductID
		)`,
		map[string]interface{}{
			"AttributeID": option.AttributeID,
			"ProductID":   option.ProductID,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeVariantTableName, "")
		}
		return nil, errors.Wrapf(err, "failed to find attribute variant with ProductID = %s, AttributeID = %s", option.ProductID, option.AttributeID)
	}

	return res, nil
}
