package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
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
	result, err := as.GetReplica().Get(attribute.AssignedProductAttribute{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedProductAttributeTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute with id=%s", id)
	}

	return result.(*attribute.AssignedProductAttribute), nil
}

func (as *SqlAssignedProductAttributeStore) GetWithOption(option *attribute.AssignedProductAttributeFilterOption) (*attribute.AssignedProductAttribute, error) {
	if option == nil || !model.IsValidId(option.ProductID) || !model.IsValidId(option.AssignmentID) {
		return nil, store.NewErrInvalidInput(store.AssignedProductAttributeTableName, "option", option)
	}

	var res *attribute.AssignedProductAttribute

	tx, err := as.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// try finding first:
	err = tx.SelectOne(
		&res,
		"SELECT * FROM "+store.AssignedProductAttributeTableName+" WHERE (ProductID = :ProductID AND AssignmentID = :AssignmentID)",
		map[string]interface{}{
			"ProductID":    option.ProductID,
			"AssignmentID": option.AssignmentID,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows { // this mean we need to create one
			res = new(attribute.AssignedProductAttribute)
			res.AssignmentID = option.AssignmentID
			res.ProductID = option.ProductID
			res.PreSave()
			if appErr := res.IsValid(); appErr != nil {
				return nil, appErr
			}

			if err = tx.Insert(res); err != nil {
				if as.IsUniqueConstraintError(err, []string{"ProductID", "AssignmentID", strings.ToLower(store.AssignedProductAttributeTableName) + "_productid_assignmentid_key"}) {
					return nil, store.NewErrInvalidInput(store.AssignedProductAttributeTableName, "ProductID/AssignmentID", option.ProductID+"/"+option.AssignmentID)
				}
				return nil, errors.Wrapf(err, "failed to insert new assigned product attribute with id=%s", res.Id)
			}
		}
		// system error
		return nil, errors.Wrapf(err, "failed to find assigned product attribute with ProductID = %s, AssignmentID = %s", option.ProductID, option.AssignmentID)
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit_transaction")
	}

	return res, nil
}
