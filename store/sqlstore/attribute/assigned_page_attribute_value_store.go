package attribute

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlAssignedPageAttributeValueStore struct {
	store.Store
}

func NewSqlAssignedPageAttributeValueStore(s store.Store) store.AssignedPageAttributeValueStore {
	return &SqlAssignedPageAttributeValueStore{s}
}

func (as *SqlAssignedPageAttributeValueStore) Upsert(tx boil.ContextTransactor, assignedPageAttrValue model.AssignedPageAttributeValueSlice) (model.AssignedPageAttributeValueSlice, error) {
	if tx == nil {
		tx = as.GetMaster()
	}

	for _, value := range assignedPageAttrValue {
		if value == nil {
			continue
		}

		isSaving := value.ID == ""
		if isSaving {
			model_helper.AssignedPageAttributeValuePreSave(value)
		}

		if err := model_helper.AssignedPageAttributeValueIsValid(*value); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = value.Insert(tx, boil.Infer())
		} else {
			_, err = value.Update(tx, boil.Infer())
		}

		if err != nil {
			if as.IsUniqueConstraintError(err, []string{"assigned_page_attributes_page_id_assignment_id_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.AssignedPageAttributeValues, model.AssignedPageAttributeColumns.PageID+"/"+model.AssignedPageAttributeColumns.AssignmentID, "unique")
			}
			return nil, err
		}
	}

	return assignedPageAttrValue, nil
}

func (as *SqlAssignedPageAttributeValueStore) Get(id string) (*model.AssignedPageAttributeValue, error) {
	value, err := model.FindAssignedPageAttributeValue(as.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.AssignedPageAttributeValues, id)
		}
		return nil, err
	}
	return value, nil
}

func (as *SqlAssignedPageAttributeValueStore) SelectForSort(assignmentID string) (model.AssignedPageAttributeValueSlice, model.AttributeValueSlice, error) {
	rows, err := as.GetReplica().
		Raw("SELECT AssignedPageAttributeValues.*, AttributeValues.* FROM "+
			model.AssignedPageAttributeValueTableName+
			" INNER JOIN "+model.AttributeValueTableName+
			" ON AttributeValues.Id = AssignedPageAttributeValues.ValueID WHERE AssignedPageAttributeValues.AssignmentID = ?", assignmentID).
		Rows()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find assignment attribute values with given assignment Id")
	}
	defer rows.Close()

	var (
		assignedPageAttributeValues []*model.AssignedPageAttributeValue
		attributeValues             []*model.AttributeValue
	)

	for rows.Next() {
		var (
			assignedPageAttributeValue model.AssignedPageAttributeValue
			attributeValue             model.AttributeValue
			scanFields                 = append(as.ScanFields(&assignedPageAttributeValue), as.AttributeValue().ScanFields(&attributeValue)...)
		)
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "error scanning a row of assigned page attribute value")
		}

		assignedPageAttributeValues = append(assignedPageAttributeValues, &assignedPageAttributeValue)
		attributeValues = append(attributeValues, &attributeValue)
	}

	return assignedPageAttributeValues, attributeValues, nil
}
