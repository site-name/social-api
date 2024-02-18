package attribute

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlAssignedProductAttributeValueStore struct {
	store.Store
}

func NewSqlAssignedProductAttributeValueStore(s store.Store) store.AssignedProductAttributeValueStore {
	return &SqlAssignedProductAttributeValueStore{s}
}

func (as *SqlAssignedProductAttributeValueStore) Save(assignedProductAttrValue model.AssignedProductAttributeValue) (*model.AssignedProductAttributeValue, error) {
	model_helper.AssignedProductAttributeValuePreSave(&assignedProductAttrValue)
	if err := model_helper.AssignedProductAttributeValueIsValid(assignedProductAttrValue); err != nil {
		return nil, err
	}

	err := assignedProductAttrValue.Insert(as.GetMaster(), boil.Infer())
	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"assigned_product_attribute_values_value_id_assignment_id_key", model.AssignedProductAttributeValueColumns.ValueID, model.AssignedProductAttributeValueColumns.AssignmentID}) {
			return nil, store.NewErrInvalidInput(model.TableNames.AssignedProductAttributeValues, "ValueID/AssignmentID", "unique")
		}
		return nil, err
	}

	return &assignedProductAttrValue, nil
}

func (as *SqlAssignedProductAttributeValueStore) Get(id string) (*model.AssignedProductAttributeValue, error) {
	value, err := model.FindAssignedProductAttributeValue(as.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.AssignedProductAttributeValues, id)
		}
		return nil, err
	}

	return value, nil
}

func (as *SqlAssignedProductAttributeValueStore) SelectForSort(assignmentID string) (model.AssignedProductAttributeValueSlice, model.AttributeValueSlice, error) {
	// rows, err := as.GetReplica().
	// 	Raw("SELECT AssignedProductAttributeValues.*, AttributeValues.* FROM "+model.AssignedProductAttributeValueTableName+" INNER JOIN "+model.AttributeValueTableName+" ON AssignedProductAttributeValues.ValueID = AttributeValues.Id WHERE AssignedProductAttributeValues.AssignmentID = ?", assignmentID).
	// 	Rows()
	// if err != nil {
	// 	return nil, nil, errors.Wrap(err, "failed to find assigned product attribute values")
	// }
	// defer rows.Close()

	// var (
	// 	assignedProductAttributeValues []*model.AssignedProductAttributeValue
	// 	attributeValues                []*model.AttributeValue
	// )

	// for rows.Next() {
	// 	var (
	// 		assignedProductAttributeValue model.AssignedProductAttributeValue
	// 		attributeValue                model.AttributeValue
	// 		scanFields                    = append(as.ScanFields(&assignedProductAttributeValue), as.AttributeValue().ScanFields(&attributeValue)...)
	// 	)

	// 	scanErr := rows.Scan(scanFields...)
	// 	if scanErr != nil {
	// 		return nil, nil, errors.Wrapf(scanErr, "error scanning values")
	// 	}

	// 	assignedProductAttributeValues = append(assignedProductAttributeValues, &assignedProductAttributeValue)
	// 	attributeValues = append(attributeValues, &attributeValue)
	// }

	// return assignedProductAttributeValues, attributeValues, nil
	panic("not implemented")
}

func (s *SqlAssignedProductAttributeValueStore) FilterByOptions(options model_helper.AssignedProductAttributeValueFilterOptions) (model.AssignedProductAttributeValueSlice, error) {
	return model.AssignedProductAttributeValues(options.Conditions...).All(s.GetReplica())
}
