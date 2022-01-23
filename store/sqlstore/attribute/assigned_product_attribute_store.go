package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedProductAttributeStore struct {
	store.Store
}

func NewSqlAssignedProductAttributeStore(s store.Store) store.AssignedProductAttributeStore {
	as := &SqlAssignedProductAttributeStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedProductAttribute{}, store.AssignedProductAttributeTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ProductID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedProductAttributeStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedProductAttributeTableName, "ProductID", store.ProductTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedProductAttributeTableName, "AssignmentID", store.AttributeProductTableName, "Id", true)
}

func (as *SqlAssignedProductAttributeStore) Save(newInstance *attribute.AssignedProductAttribute) (*attribute.AssignedProductAttribute, error) {
	newInstance.PreSave()
	if err := newInstance.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(newInstance); err != nil {
		if as.IsUniqueConstraintError(err, []string{"ProductID", "AssignmentID", strings.ToLower(store.AssignedProductAttributeTableName) + "_productid_assignmentid_key"}) {
			return nil, store.NewErrInvalidInput(store.AssignedProductAttributeTableName, "ProductID/AssignmentID", newInstance.ProductID+"/"+newInstance.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to insert new assigned product attribute with id=%s", newInstance.Id)
	}

	return newInstance, nil
}

func (as *SqlAssignedProductAttributeStore) Get(id string) (*attribute.AssignedProductAttribute, error) {
	var res attribute.AssignedProductAttribute
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AssignedProductAttributeTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedProductAttributeTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAssignedProductAttributeStore) GetWithOption(option *attribute.AssignedProductAttributeFilterOption) (*attribute.AssignedProductAttribute, error) {
	query := as.GetQueryBuilder().Select("*").From(store.AssignedProductAttributeTableName)

	// parse option
	if option.AssignmentID != nil {
		query = query.Where(option.AssignmentID)
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetWithOption_ToSql")
	}

	var res attribute.AssignedProductAttribute

	err = as.GetReplica().SelectOne(
		&res,
		queryString,
		args...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedProductAttributeTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute with given options")
	}

	return &res, nil
}
