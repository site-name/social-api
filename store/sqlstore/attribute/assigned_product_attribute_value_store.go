package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedProductAttributeValueStore struct {
	store.Store
}

var assignedProductAttrValueDuplicateKeys = []string{
	"ValueID",
	"AssignmentID",
	"assignedproductattributevalues_valueid_assignmentid_key",
}

func NewSqlAssignedProductAttributeValueStore(s store.Store) store.AssignedProductAttributeValueStore {
	return &SqlAssignedProductAttributeValueStore{s}
}

func (as *SqlAssignedProductAttributeValueStore) ModelFields(prefix string) model.StringArray {
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

func (as *SqlAssignedProductAttributeValueStore) ScanFields(assignedProductAttributeValue attribute.AssignedProductAttributeValue) []interface{} {
	return []interface{}{
		&assignedProductAttributeValue.Id,
		&assignedProductAttributeValue.ValueID,
		&assignedProductAttributeValue.AssignmentID,
		&assignedProductAttributeValue.SortOrder,
	}
}

func (as *SqlAssignedProductAttributeValueStore) Save(assignedProductAttrValue *attribute.AssignedProductAttributeValue) (*attribute.AssignedProductAttributeValue, error) {
	assignedProductAttrValue.PreSave()
	if err := assignedProductAttrValue.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.AssignedProductAttributeValueTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"

	if _, err := as.GetMasterX().NamedExec(query, assignedProductAttrValue); err != nil {
		if as.IsUniqueConstraintError(err, assignedProductAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.AssignedProductAttributeValueTableName, "ValueID/AssignmentID", assignedProductAttrValue.ValueID+"/"+assignedProductAttrValue.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned product attribute value with id=%s", assignedProductAttrValue.Id)
	}

	return assignedProductAttrValue, nil
}

func (as *SqlAssignedProductAttributeValueStore) Get(assignedProductAttrValueID string) (*attribute.AssignedProductAttributeValue, error) {
	var res attribute.AssignedProductAttributeValue

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AssignedProductAttributeValueTableName+" WHERE Id = ?", assignedProductAttrValueID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedProductAttributeValueTableName, assignedProductAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute value with id=%s", assignedProductAttrValueID)
	}

	return &res, nil
}

func (as *SqlAssignedProductAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedProductAttributeValue, error) {
	tx, err := as.GetMasterX().Beginx()
	if err != nil {
		return nil, errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// return value:
	res := []*attribute.AssignedProductAttributeValue{}

	insertQuery := "INSERT INTO " + store.AssignedProductAttributeValueTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"

	for _, id := range attributeValueIDs {
		newValue := &attribute.AssignedProductAttributeValue{
			ValueID:      id,
			AssignmentID: assignmentID,
		}
		newValue.PreSave()

		if appErr := newValue.IsValid(); appErr != nil {
			return nil, appErr
		}

		_, err = tx.NamedExec(insertQuery, newValue)
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
	tx, err := as.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	updateQuery := "UPDATE " + store.AssignedProductAttributeValueTableName + " SET " +
		as.ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).Join(",") + " WHERE Id=:Id"

	for _, value := range attributeValues {
		// try validating if the value exist:
		_, err := as.Get(value.Id)
		if err != nil {
			return err
		}

		result, err := tx.NamedExec(updateQuery, value)
		if err != nil {
			// check if error is duplicate conflict error:
			if as.IsUniqueConstraintError(err, assignedProductAttrValueDuplicateKeys) {
				return store.NewErrInvalidInput(store.AssignedProductAttributeValueTableName, "ValueID/AssignmentID", value.ValueID+"/"+value.AssignmentID)
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

func (as *SqlAssignedProductAttributeValueStore) SelectForSort(assignmentID string) ([]*attribute.AssignedProductAttributeValue, []*attribute.AttributeValue, error) {
	query, args, err := as.GetQueryBuilder().
		Select(append(as.ModelFields(store.AssignedProductAttributeValueTableName+"."), as.AttributeValue().ModelFields(store.AttributeValueTableName+".")...)...).
		From(store.AssignedProductAttributeValueTableName).
		InnerJoin(store.AttributeValueTableName + " ON (AssignedProductAttributeValues.Id = AttributeValues.ValueID)").
		Where(squirrel.Eq{"AssignedProductAttributeValues.AssignmentID": assignmentID}).
		ToSql()

	if err != nil {
		return nil, nil, errors.Wrap(err, "SelectForSort_ToSql")
	}

	rows, err := as.GetReplicaX().QueryX(query, args...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find assigned product attribute values")
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
