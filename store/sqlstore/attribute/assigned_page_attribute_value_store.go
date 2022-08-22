package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedPageAttributeValueStore struct {
	store.Store
}

var (
	assignedPageAttrValueDuplicateKeys = []string{"ValueID", "AssignmentID", "assignedpageattributevalues_valueid_assignmentid_key"}
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

func (as *SqlAssignedPageAttributeValueStore) ModelFields() []string {
	return []string{
		"AssignedPageAttributeValues.Id",
		"AssignedPageAttributeValues.ValueID",
		"AssignedPageAttributeValues.AssignmentID",
		"AssignedPageAttributeValues.SortOrder",
	}
}

func (as *SqlAssignedPageAttributeValueStore) ScanFields(attributeValue attribute.AssignedPageAttributeValue) []interface{} {
	return []interface{}{
		&attributeValue.Id,
		&attributeValue.ValueID,
		&attributeValue.AssignmentID,
		&attributeValue.SortOrder,
	}
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
	var res attribute.AssignedPageAttributeValue
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AssignedPageAttributeValueTableName+" WHERE Id = :ID", map[string]interface{}{"ID": assignedPageAttrValueID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedPageAttributeValueTableName, assignedPageAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute value with id=%s", assignedPageAttrValueID)
	}

	return &res, nil
}

func (as *SqlAssignedPageAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedPageAttributeValue, error) {
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

		err := as.GetMaster().Insert(newValue)
		if err != nil {
			if as.IsUniqueConstraintError(err, assignedPageAttrValueDuplicateKeys) {
				return nil, store.NewErrInvalidInput(store.AssignedPageAttributeValueTableName, "ValueID/AssignmentID", newValue.ValueID+"/"+newValue.AssignmentID)
			}
			return nil, errors.Wrapf(err, "failed to save assigned page attribute value with id=%s", newValue.Id)
		}
		// append to return value if success
		res = append(res, newValue)
	}

	return res, nil
}

func (as *SqlAssignedPageAttributeValueStore) SelectForSort(assignmentID string) ([]*attribute.AssignedPageAttributeValue, []*attribute.AttributeValue, error) {
	rows, err := as.GetQueryBuilder().
		Select(append(as.ModelFields(), as.AttributeValue().ModelFields()...)...).
		From(store.AssignedPageAttributeValueTableName).
		InnerJoin(store.AttributeValueTableName + " ON (AttributeValues.Id = AssignedPageAttributeValues.ValueID)").
		Where(squirrel.Eq{"AssignedPageAttributeValues.AssignmentID": assignmentID}).
		RunWith(as.GetReplica()).
		Query()

	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to find values with AssignmentID=%s", assignmentID)
	}

	var (
		assignedPageAttributeValues []*attribute.AssignedPageAttributeValue
		attributeValues             []*attribute.AttributeValue
		assignedPageAttributeValue  attribute.AssignedPageAttributeValue
		attributeValue              attribute.AttributeValue
		scanFields                  = append(as.ScanFields(assignedPageAttributeValue), as.AttributeValue().ScanFields(attributeValue)...)
	)

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "error scanning a row of assigned page attribute value")
		}

		assignedPageAttributeValues = append(assignedPageAttributeValues, assignedPageAttributeValue.DeepCopy())
		attributeValues = append(attributeValues, attributeValue.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, nil, errors.Wrap(err, "error closing rows")
	}

	return assignedPageAttributeValues, attributeValues, nil
}

func (as *SqlAssignedPageAttributeValueStore) UpdateInBulk(attributeValues []*attribute.AssignedPageAttributeValue) error {

	for _, value := range attributeValues {
		// try validating if the value exist:
		_, err := as.Get(value.Id)
		if err != nil {
			return errors.Wrapf(err, "failed to find value with id=%s", value.Id)
		}
		numUpdated, err := as.GetMaster().Update(value)
		if err != nil {
			// check if error is duplicate conflict error:
			if as.IsUniqueConstraintError(err, assignedPageAttrValueDuplicateKeys) {
				return store.NewErrInvalidInput(store.AssignedPageAttributeValueTableName, "ValueID/AssignmentID", value.ValueID+"/"+value.AssignmentID)
			}
			return errors.Wrapf(err, "failed to update value with id=%s", value.Id)
		}
		if numUpdated > 1 {
			return errors.Errorf("more than one value with id=%s were updated(%d)", value.Id, numUpdated)
		}
	}

	return nil
}
