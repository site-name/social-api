package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlAssignedPageAttributeValueStore struct {
	store.Store
}

var assignedPageAttrValueDuplicateKeys = []string{"ValueID", "AssignmentID", "assignedpageattributevalues_valueid_assignmentid_key"}

func NewSqlAssignedPageAttributeValueStore(s store.Store) store.AssignedPageAttributeValueStore {
	return &SqlAssignedPageAttributeValueStore{s}
}

func (as *SqlAssignedPageAttributeValueStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"ValueID",
		"AssignmentID",
		"SortOrder",
	}

	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (as *SqlAssignedPageAttributeValueStore) ScanFields(attributeValue *model.AssignedPageAttributeValue) []interface{} {
	return []interface{}{
		&attributeValue.Id,
		&attributeValue.ValueID,
		&attributeValue.AssignmentID,
		&attributeValue.SortOrder,
	}
}

func (as *SqlAssignedPageAttributeValueStore) Save(assignedPageAttrValue *model.AssignedPageAttributeValue) (*model.AssignedPageAttributeValue, error) {
	assignedPageAttrValue.PreSave()
	if err := assignedPageAttrValue.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.AssignedPageAttributeValueTableName + "(" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
	if _, err := as.GetMasterX().NamedExec(query, assignedPageAttrValue); err != nil {
		if as.IsUniqueConstraintError(err, assignedPageAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.AssignedPageAttributeValueTableName, "ValueID/AssignmentID", assignedPageAttrValue.ValueID+"/"+assignedPageAttrValue.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned page attribute value with id=%s", assignedPageAttrValue.Id)
	}

	return assignedPageAttrValue, nil
}

func (as *SqlAssignedPageAttributeValueStore) Get(assignedPageAttrValueID string) (*model.AssignedPageAttributeValue, error) {
	var res model.AssignedPageAttributeValue

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AssignedPageAttributeValueTableName+" WHERE Id = ?", assignedPageAttrValueID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedPageAttributeValueTableName, assignedPageAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute value with id=%s", assignedPageAttrValueID)
	}

	return &res, nil
}

func (as *SqlAssignedPageAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*model.AssignedPageAttributeValue, error) {
	// return value:
	res := []*model.AssignedPageAttributeValue{}

	for _, id := range attributeValueIDs {
		newValue, err := as.Save(&model.AssignedPageAttributeValue{
			ValueID:      id,
			AssignmentID: assignmentID,
		})

		if err != nil {
			return nil, err
		}
		// append to return value if success
		res = append(res, newValue)
	}

	return res, nil
}

func (as *SqlAssignedPageAttributeValueStore) SelectForSort(assignmentID string) ([]*model.AssignedPageAttributeValue, []*model.AttributeValue, error) {
	query, args, err := as.GetQueryBuilder().
		Select(append(as.ModelFields(store.AssignedPageAttributeValueTableName+"."), as.AttributeValue().ModelFields(store.AttributeValueTableName+".")...)...).
		From(store.AssignedPageAttributeValueTableName).
		InnerJoin(store.AttributeValueTableName + " ON (AttributeValues.Id = AssignedPageAttributeValues.ValueID)").
		Where(squirrel.Eq{"AssignedPageAttributeValues.AssignmentID": assignmentID}).
		ToSql()

	if err != nil {
		return nil, nil, errors.Wrap(err, "SelectForSort_ToSql")
	}

	rows, err := as.GetReplicaX().QueryX(query, args...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find assignment attribute values with given assignment Id")
	}

	var (
		assignedPageAttributeValues []*model.AssignedPageAttributeValue
		attributeValues             []*model.AttributeValue
		assignedPageAttributeValue  model.AssignedPageAttributeValue
		attributeValue              model.AttributeValue
		scanFields                  = append(as.ScanFields(&assignedPageAttributeValue), as.AttributeValue().ScanFields(&attributeValue)...)
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

func (as *SqlAssignedPageAttributeValueStore) UpdateInBulk(attributeValues []*model.AssignedPageAttributeValue) error {

	query := "UPDATE " + store.AssignedPageAttributeValueTableName + " SET " + as.
		ModelFields("").
		Map(func(_ int, s string) string {
			return s + "=:" + s
		}).
		Join(",") + " WHERE Id=:Id"

	for _, value := range attributeValues {
		result, err := as.GetMasterX().NamedExec(query, value)
		if err != nil {
			// check if error is duplicate conflict error:
			if as.IsUniqueConstraintError(err, assignedPageAttrValueDuplicateKeys) {
				return store.NewErrInvalidInput(store.AssignedPageAttributeValueTableName, "ValueID/AssignmentID", value.ValueID+"/"+value.AssignmentID)
			}
			return errors.Wrapf(err, "failed to update value with id=%s", value.Id)
		}

		numUpdated, _ := result.RowsAffected()
		if numUpdated > 1 {
			return errors.Errorf("more than one value with id=%s were updated(%d)", value.Id, numUpdated)
		}
	}

	return nil
}
