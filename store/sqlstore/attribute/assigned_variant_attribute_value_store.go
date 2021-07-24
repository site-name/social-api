package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedVariantAttributeValueStore struct {
	store.Store
}

var (
	assignedVariantAttrValueDuplicateKeys = []string{"ValueID", "AssignmentID", strings.ToLower(store.AssignedVariantAttributeValueTableName) + "_valueid_assignmentid_key"}
)

func NewSqlAssignedVariantAttributeValueStore(s store.Store) store.AssignedVariantAttributeValueStore {
	as := &SqlAssignedVariantAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedVariantAttributeValue{}, store.AssignedVariantAttributeValueTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedVariantAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedVariantAttributeValueTableName, "ValueID", store.AttributeValueTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedVariantAttributeValueTableName, "AssignmentID", store.AssignedVariantAttributeTableName, "Id", true)
}

func (as *SqlAssignedVariantAttributeValueStore) Save(assignedVariantAttrValue *attribute.AssignedVariantAttributeValue) (*attribute.AssignedVariantAttributeValue, error) {
	assignedVariantAttrValue.PreSave()
	if err := assignedVariantAttrValue.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(assignedVariantAttrValue); err != nil {
		if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", assignedVariantAttrValue.ValueID+"/"+assignedVariantAttrValue.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned variant attribute value with id=%s", assignedVariantAttrValue.Id)
	}

	return assignedVariantAttrValue, nil
}

func (as *SqlAssignedVariantAttributeValueStore) Get(assignedVariantAttrValueID string) (*attribute.AssignedVariantAttributeValue, error) {
	res, err := as.GetReplica().Get(attribute.AssignedVariantAttributeValue{}, assignedVariantAttrValueID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedVariantAttributeValueTableName, assignedVariantAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute value with id=%s", assignedVariantAttrValueID)
	}

	return res.(*attribute.AssignedVariantAttributeValue), nil
}

func (as *SqlAssignedVariantAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedVariantAttributeValue, error) {
	tx, err := as.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// return value:
	res := []*attribute.AssignedVariantAttributeValue{}

	for _, id := range attributeValueIDs {
		newValue := &attribute.AssignedVariantAttributeValue{
			ValueID:      id,
			AssignmentID: assignmentID,
		}
		newValue.PreSave()
		if appErr := newValue.IsValid(); appErr != nil {
			return nil, appErr
		}

		err = tx.Insert(newValue)
		if err != nil {
			if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
				return nil, store.NewErrInvalidInput(store.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", newValue.ValueID+"/"+newValue.AssignmentID)
			}
			return nil, errors.Wrapf(err, "failed to save assigned variant attribute value with id=%s", newValue.Id)
		}
		// append to return value if success
		res = append(res, newValue)
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "commit_transaction")
	}

	return res, nil
}
