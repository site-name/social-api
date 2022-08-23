package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedVariantAttributeValueStore struct {
	store.Store
}

var assignedVariantAttrValueDuplicateKeys = []string{
	"ValueID",
	"AssignmentID",
	"assignedvariantattributevalues_valueid_assignmentid_key",
}

func NewSqlAssignedVariantAttributeValueStore(s store.Store) store.AssignedVariantAttributeValueStore {
	return &SqlAssignedVariantAttributeValueStore{s}
}

func (as *SqlAssignedVariantAttributeValueStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"ValueID",
		"AssignmentID",
		"SortOrder",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAssignedVariantAttributeValueStore) ScanFields(assignedVariantAttributeValue attribute.AssignedVariantAttributeValue) []interface{} {
	return []interface{}{
		&assignedVariantAttributeValue.Id,
		&assignedVariantAttributeValue.ValueID,
		&assignedVariantAttributeValue.AssignmentID,
		&assignedVariantAttributeValue.SortOrder,
	}
}

func (as *SqlAssignedVariantAttributeValueStore) Save(assignedVariantAttrValue *attribute.AssignedVariantAttributeValue) (*attribute.AssignedVariantAttributeValue, error) {
	assignedVariantAttrValue.PreSave()
	if err := assignedVariantAttrValue.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.AssignedVariantAttributeValueTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
	if _, err := as.GetMasterX().NamedExec(query, assignedVariantAttrValue); err != nil {
		if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", assignedVariantAttrValue.ValueID+"/"+assignedVariantAttrValue.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned variant attribute value with id=%s", assignedVariantAttrValue.Id)
	}

	return assignedVariantAttrValue, nil
}

func (as *SqlAssignedVariantAttributeValueStore) Get(assignedVariantAttrValueID string) (*attribute.AssignedVariantAttributeValue, error) {
	var res attribute.AssignedVariantAttributeValue

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AssignedVariantAttributeValueTableName+" WHERE Id = :ID", map[string]interface{}{"ID": assignedVariantAttrValueID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedVariantAttributeValueTableName, assignedVariantAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute value with id=%s", assignedVariantAttrValueID)
	}

	return &res, nil
}

func (as *SqlAssignedVariantAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedVariantAttributeValue, error) {
	tx, err := as.GetMasterX().Beginx()
	if err != nil {
		return nil, errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// return value:
	res := []*attribute.AssignedVariantAttributeValue{}

	query := "INSERT INTO " + store.AssignedVariantAttributeValueTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"

	for _, id := range attributeValueIDs {
		newValue := &attribute.AssignedVariantAttributeValue{
			ValueID:      id,
			AssignmentID: assignmentID,
		}
		newValue.PreSave()
		if appErr := newValue.IsValid(); appErr != nil {
			return nil, appErr
		}

		_, err = tx.NamedExec(query, newValue)
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

func (as *SqlAssignedVariantAttributeValueStore) SelectForSort(assignmentID string) ([]*attribute.AssignedVariantAttributeValue, []*attribute.AttributeValue, error) {
	query, args, err := as.GetQueryBuilder().
		Select(append(as.ModelFields(store.AssignedVariantAttributeValueTableName+"."), as.AttributeValue().ModelFields(store.AttributeValueTableName+".")...)...).
		From(store.AssignedVariantAttributeValueTableName).
		InnerJoin(store.AttributeValueTableName + " ON (AttributeValues.Id = AssignedVariantAttributeValues.ValueID)").
		Where(squirrel.Eq{"AssignedVariantAttributeValues.AssignmentID": assignmentID}).
		ToSql()

	if err != nil {
		return nil, nil, errors.Wrap(err, "SelectForSort_ToSql")
	}

	rows, err := as.GetReplicaX().QueryX(query, args...)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to find values with AssignmentID=%s", assignmentID)
	}
	var (
		assignedVariantAttributeValues []*attribute.AssignedVariantAttributeValue
		attributeValues                []*attribute.AttributeValue
		assignedVariantAttributeValue  attribute.AssignedVariantAttributeValue
		attributeValue                 attribute.AttributeValue
		scanFields                     = append(as.ScanFields(assignedVariantAttributeValue), as.AttributeValue().ScanFields(attributeValue)...)
	)
	for rows.Next() {
		scanErr := rows.Scan(scanFields...)
		if scanErr != nil {
			return nil, nil, errors.Wrapf(scanErr, "error scanning values")
		}

		assignedVariantAttributeValues = append(assignedVariantAttributeValues, assignedVariantAttributeValue.DeepCopy())
		attributeValues = append(attributeValues, attributeValue.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, nil, errors.Wrap(err, "error closing rows")
	}

	if err = rows.Err(); err != nil {
		return nil, nil, errors.Wrap(err, "error occured during scanning iteration")
	}

	return assignedVariantAttributeValues, attributeValues, nil
}

func (as *SqlAssignedVariantAttributeValueStore) UpdateInBulk(attributeValues []*attribute.AssignedVariantAttributeValue) error {
	tx, err := as.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	query := "UPDATE " + store.AssignedVariantAttributeValueTableName + " SET " +
		as.ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

	for _, value := range attributeValues {
		// try validating if the value exist:
		_, err := as.Get(value.Id)
		if err != nil {
			return err
		}

		result, err := tx.NamedExec(query, value)
		if err != nil {
			// check if error is duplicate conflict error:
			if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
				return store.NewErrInvalidInput(store.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", value.ValueID+"/"+value.AssignmentID)
			}
			return errors.Wrapf(err, "failed to update value with id=%s", value.Id)
		}
		if numUpdated, _ := result.RowsAffected(); numUpdated > 1 {
			return errors.Errorf("more than one value with id=%s were updated(%d)", value.Id, numUpdated)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}

	return nil
}
