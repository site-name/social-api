package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
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

func (as *SqlAssignedVariantAttributeValueStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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

func (as *SqlAssignedVariantAttributeValueStore) ScanFields(assignedVariantAttributeValue *model.AssignedVariantAttributeValue) []interface{} {
	return []interface{}{
		&assignedVariantAttributeValue.Id,
		&assignedVariantAttributeValue.ValueID,
		&assignedVariantAttributeValue.AssignmentID,
		&assignedVariantAttributeValue.SortOrder,
	}
}

func (as *SqlAssignedVariantAttributeValueStore) Save(assignedVariantAttrValue *model.AssignedVariantAttributeValue) (*model.AssignedVariantAttributeValue, error) {
	assignedVariantAttrValue.PreSave()
	if err := assignedVariantAttrValue.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.AssignedVariantAttributeValueTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
	if _, err := as.GetMasterX().NamedExec(query, assignedVariantAttrValue); err != nil {
		if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(model.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", assignedVariantAttrValue.ValueID+"/"+assignedVariantAttrValue.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned variant attribute value with id=%s", assignedVariantAttrValue.Id)
	}

	return assignedVariantAttrValue, nil
}

func (as *SqlAssignedVariantAttributeValueStore) Get(assignedVariantAttrValueID string) (*model.AssignedVariantAttributeValue, error) {
	var res model.AssignedVariantAttributeValue

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+model.AssignedVariantAttributeValueTableName+" WHERE Id = :ID", map[string]interface{}{"ID": assignedVariantAttrValueID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.AssignedVariantAttributeValueTableName, assignedVariantAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute value with id=%s", assignedVariantAttrValueID)
	}

	return &res, nil
}

func (as *SqlAssignedVariantAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*model.AssignedVariantAttributeValue, error) {
	tx, err := as.GetMasterX().Beginx()
	if err != nil {
		return nil, errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// return value:
	res := []*model.AssignedVariantAttributeValue{}

	query := "INSERT INTO " + model.AssignedVariantAttributeValueTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"

	for _, id := range attributeValueIDs {
		newValue := &model.AssignedVariantAttributeValue{
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
				return nil, store.NewErrInvalidInput(model.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", newValue.ValueID+"/"+newValue.AssignmentID)
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

func (as *SqlAssignedVariantAttributeValueStore) SelectForSort(assignmentID string) ([]*model.AssignedVariantAttributeValue, []*model.AttributeValue, error) {
	query, args, err := as.GetQueryBuilder().
		Select(append(as.ModelFields(model.AssignedVariantAttributeValueTableName+"."), as.AttributeValue().ModelFields(model.AttributeValueTableName+".")...)...).
		From(model.AssignedVariantAttributeValueTableName).
		InnerJoin(model.AttributeValueTableName + " ON (AttributeValues.Id = AssignedVariantAttributeValues.ValueID)").
		Where(squirrel.Eq{"AssignedVariantAttributeValues.AssignmentID": assignmentID}).
		ToSql()

	if err != nil {
		return nil, nil, errors.Wrap(err, "SelectForSort_ToSql")
	}

	rows, err := as.GetReplicaX().QueryX(query, args...)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to find values with AssignmentID=%s", assignmentID)
	}
	defer rows.Close()

	var (
		assignedVariantAttributeValues []*model.AssignedVariantAttributeValue
		attributeValues                []*model.AttributeValue
	)
	for rows.Next() {
		var (
			assignedVariantAttributeValue model.AssignedVariantAttributeValue
			attributeValue                model.AttributeValue
			scanFields                    = append(as.ScanFields(&assignedVariantAttributeValue), as.AttributeValue().ScanFields(&attributeValue)...)
		)
		scanErr := rows.Scan(scanFields...)
		if scanErr != nil {
			return nil, nil, errors.Wrapf(scanErr, "error scanning values")
		}

		assignedVariantAttributeValues = append(assignedVariantAttributeValues, &assignedVariantAttributeValue)
		attributeValues = append(attributeValues, &attributeValue)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, errors.Wrap(err, "error occured during scanning iteration")
	}

	return assignedVariantAttributeValues, attributeValues, nil
}

func (as *SqlAssignedVariantAttributeValueStore) UpdateInBulk(attributeValues []*model.AssignedVariantAttributeValue) error {
	query := "UPDATE " + model.AssignedVariantAttributeValueTableName + " SET " +
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

		result, err := as.GetMasterX().NamedExec(query, value)
		if err != nil {
			// check if error is duplicate conflict error:
			if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
				return store.NewErrInvalidInput(model.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", value.ValueID+"/"+value.AssignmentID)
			}
			return errors.Wrapf(err, "failed to update value with id=%s", value.Id)
		}
		if numUpdated, _ := result.RowsAffected(); numUpdated > 1 {
			return errors.Errorf("more than one value with id=%s were updated(%d)", value.Id, numUpdated)
		}
	}

	return nil
}

func (s *SqlAssignedVariantAttributeValueStore) FilterByOptions(options *model.AssignedVariantAttributeValueFilterOptions) ([]*model.AssignedVariantAttributeValue, error) {
	query := s.GetQueryBuilder().Select("*").From(model.AssignedVariantAttributeValueTableName)
	if options.AssignmentID != nil {
		query = query.Where(options.AssignmentID)
	}
	if options.ValueID != nil {
		query = query.Where(options.ValueID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.AssignedVariantAttributeValue
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find assigned variant attribute values by given options")
	}

	return res, nil
}
