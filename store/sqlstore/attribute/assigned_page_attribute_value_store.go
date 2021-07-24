package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedPageAttributeValueStore struct {
	store.Store
}

var (
	assignedPageAttrValueDuplicateKeys = []string{"ValueID", "AssignmentID", strings.ToLower(store.AssignedPageAttributeValueTableName) + "_valueid_assignmentid_key"}
)

func NewSqlAssignedPageAttributeValueStore(s store.Store) store.AssignedPageAttributeValueStore {
	as := &SqlAssignedPageAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedPageAttributeValue{}, store.AssignedPageAttributeValueTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedPageAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedPageAttributeValueTableName, "ValueID", store.AttributeValueTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedPageAttributeValueTableName, "AssignmentID", store.AssignedPageAttributeTableName, "Id", true)
}

func (as *SqlAssignedPageAttributeValueStore) Save(assignedPageAttrValue *attribute.AssignedPageAttributeValue) (*attribute.AssignedPageAttributeValue, error) {
	assignedPageAttrValue.PreSave()
	if err := assignedPageAttrValue.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(assignedPageAttrValue); err != nil {
		if as.IsUniqueConstraintError(err, assignedPageAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.AssignedPageAttributeValueTableName, "ValueID/AssignmentID", assignedPageAttrValue.ValueID+"/"+assignedPageAttrValue.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned page attribute value with id=%s", assignedPageAttrValue.Id)
	}

	return assignedPageAttrValue, nil
}

func (as *SqlAssignedPageAttributeValueStore) Get(assignedPageAttrValueID string) (*attribute.AssignedPageAttributeValue, error) {
	res, err := as.GetReplica().Get(attribute.AssignedPageAttributeValue{}, assignedPageAttrValueID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedPageAttributeValueTableName, assignedPageAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute value with id=%s", assignedPageAttrValueID)
	}

	return res.(*attribute.AssignedPageAttributeValue), nil
}

func (as *SqlAssignedPageAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedPageAttributeValue, error) {
	tx, err := as.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// return value:
	res := []*attribute.AssignedPageAttributeValue{}

	for _, id := range attributeValueIDs {
		newValue := &attribute.AssignedPageAttributeValue{
			ValueID:      id,
			AssignmentID: assignmentID,
		}
		newValue.PreSave()
		if appErr := newValue.IsValid(); appErr != nil {
			return nil, appErr
		}

		err = tx.Insert(newValue)
		if err != nil {
			if as.IsUniqueConstraintError(err, assignedPageAttrValueDuplicateKeys) {
				return nil, store.NewErrInvalidInput(store.AssignedPageAttributeValueTableName, "ValueID/AssignmentID", newValue.ValueID+"/"+newValue.AssignmentID)
			}
			return nil, errors.Wrapf(err, "failed to save assigned page attribute value with id=%s", newValue.Id)
		}
		// append to return value if success
		res = append(res, newValue)
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "commit_transaction")
	}

	return res, nil
}
