package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedProductAttributeValueStore struct {
	store.Store
}

var (
	assignedProductAttrValueDuplicateKeys = []string{
		"ValueID",
		"AssignmentID",
		"assignedproductattributevalues_valueid_assignmentid_key",
	}
)

func NewSqlAssignedProductAttributeValueStore(s store.Store) store.AssignedProductAttributeValueStore {
	as := &SqlAssignedProductAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedProductAttributeValue{}, store.AssignedProductAttributeValueTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedProductAttributeValueStore) ModelFields() []string {
	return []string{
		"AssignedProductAttributeValues.Id",
		"AssignedProductAttributeValues.ValueID",
		"AssignedProductAttributeValues.AssignmentID",
		"AssignedProductAttributeValues.SortOrder",
	}
}

func (as *SqlAssignedProductAttributeValueStore) ScanFields(assignedProductAttributeValue attribute.AssignedProductAttributeValue) []interface{} {
	return []interface{}{
		&assignedProductAttributeValue.Id,
		&assignedProductAttributeValue.ValueID,
		&assignedProductAttributeValue.AssignmentID,
		&assignedProductAttributeValue.SortOrder,
	}
}

func (as *SqlAssignedProductAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedProductAttributeValueTableName, "ValueID", store.AttributeValueTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedProductAttributeValueTableName, "AssignmentID", store.AssignedProductAttributeTableName, "Id", true)
}

func (as *SqlAssignedProductAttributeValueStore) Save(assignedProductAttrValue *attribute.AssignedProductAttributeValue) (*attribute.AssignedProductAttributeValue, error) {
	assignedProductAttrValue.PreSave()
	if err := assignedProductAttrValue.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(assignedProductAttrValue); err != nil {
		if as.IsUniqueConstraintError(err, assignedProductAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.AssignedProductAttributeValueTableName, "ValueID/AssignmentID", assignedProductAttrValue.ValueID+"/"+assignedProductAttrValue.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned product attribute value with id=%s", assignedProductAttrValue.Id)
	}

	return assignedProductAttrValue, nil
}

func (as *SqlAssignedProductAttributeValueStore) Get(assignedProductAttrValueID string) (*attribute.AssignedProductAttributeValue, error) {
	var res attribute.AssignedProductAttributeValue
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AssignedProductAttributeValueTableName+" WHERE Id = :ID", map[string]interface{}{"ID": assignedProductAttrValueID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedProductAttributeValueTableName, assignedProductAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute value with id=%s", assignedProductAttrValueID)
	}

	return &res, nil
}

func (as *SqlAssignedProductAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedProductAttributeValue, error) {
	tx, err := as.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// return value:
	res := []*attribute.AssignedProductAttributeValue{}

	for _, id := range attributeValueIDs {
		newValue := &attribute.AssignedProductAttributeValue{
			ValueID:      id,
			AssignmentID: assignmentID,
		}
		newValue.PreSave()
		if appErr := newValue.IsValid(); appErr != nil {
			return nil, appErr
		}

		err = tx.Insert(newValue)
		if err != nil {
			if as.IsUniqueConstraintError(err, assignedProductAttrValueDuplicateKeys) {
				return nil, store.NewErrInvalidInput(store.AssignedProductAttributeValueTableName, "ValueID/AssignmentID", newValue.ValueID+"/"+newValue.AssignmentID)
			}
			return nil, errors.Wrapf(err, "failed to save assigned product attribute value with id=%s", newValue.Id)
		}
		// append to return value if success
		res = append(res, newValue)
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "commit_transaction")
	}

	return res, nil
}

func (as *SqlAssignedProductAttributeValueStore) UpdateInBulk(attributeValues []*attribute.AssignedProductAttributeValue) error {
	tx, err := as.GetMaster().Begin()
	if err != nil {
		return errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	for _, value := range attributeValues {
		// try validating if the value exist:
		_, err := tx.Get(attribute.AssignedProductAttributeValue{}, value.Id)
		if err != nil {
			return errors.Wrapf(err, "failed to find value with id=%s", value.Id)
		}
		numUpdated, err := tx.Update(value)
		if err != nil {
			// check if error is duplicate conflict error:
			if as.IsUniqueConstraintError(err, assignedProductAttrValueDuplicateKeys) {
				return store.NewErrInvalidInput(store.AssignedProductAttributeValueTableName, "ValueID/AssignmentID", value.ValueID+"/"+value.AssignmentID)
			}
			return errors.Wrapf(err, "failed to update value with id=%s", value.Id)
		}
		if numUpdated > 1 {
			return errors.Errorf("more than one value with id=%s were updated(%d)", value.Id, numUpdated)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}

	return nil
}

func (as *SqlAssignedProductAttributeValueStore) SelectForSort(assignmentID string) ([]*attribute.AssignedProductAttributeValue, []*attribute.AttributeValue, error) {
	rows, err := as.GetQueryBuilder().
		Select(append(as.ModelFields(), as.AttributeValue().ModelFields()...)...).
		From(store.AssignedProductAttributeValueTableName).
		InnerJoin(store.AttributeValueTableName + " ON (AssignedProductAttributeValues.Id = AttributeValues.ValueID)").
		Where(squirrel.Eq{"AssignedProductAttributeValues.AssignmentID": assignmentID}).
		RunWith(as.GetReplica()).
		Query()

	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to find values with AssignmentID=%s", assignmentID)
	}
	var (
		assignedProductAttributeValues []*attribute.AssignedProductAttributeValue
		attributeValues                []*attribute.AttributeValue
		assignedProductAttributeValue  attribute.AssignedProductAttributeValue
		attributeValue                 attribute.AttributeValue
		scanFields                     = append(as.ScanFields(assignedProductAttributeValue), as.AttributeValue().ScanFields(attributeValue)...)
	)
	for rows.Next() {
		scanErr := rows.Scan(scanFields...)
		if scanErr != nil {
			return nil, nil, errors.Wrapf(scanErr, "error scanning values")
		}

		assignedProductAttributeValues = append(assignedProductAttributeValues, assignedProductAttributeValue.DeepCopy())
		attributeValues = append(attributeValues, attributeValue.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, nil, errors.Wrap(err, "error closing rows")
	}

	return assignedProductAttributeValues, attributeValues, nil
}
