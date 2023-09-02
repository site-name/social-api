package attribute

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAssignedPageAttributeValueStore struct {
	store.Store
}

var assignedPageAttrValueDuplicateKeys = []string{"ValueID", "AssignmentID", "valueid_assignmentid_key"}

func NewSqlAssignedPageAttributeValueStore(s store.Store) store.AssignedPageAttributeValueStore {
	return &SqlAssignedPageAttributeValueStore{s}
}

func (as *SqlAssignedPageAttributeValueStore) ScanFields(attributeValue *model.AssignedPageAttributeValue) []interface{} {
	return []interface{}{
		&attributeValue.ValueID,
		&attributeValue.AssignmentID,
		&attributeValue.SortOrder,
	}
}

func (as *SqlAssignedPageAttributeValueStore) Save(assignedPageAttrValue *model.AssignedPageAttributeValue) (*model.AssignedPageAttributeValue, error) {
	if err := as.GetMaster().Create(assignedPageAttrValue).Error; err != nil {
		if as.IsUniqueConstraintError(err, assignedPageAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(model.AssignedPageAttributeValueTableName, "ValueID/AssignmentID", assignedPageAttrValue.ValueID+"/"+assignedPageAttrValue.AssignmentID)
		}
		return nil, errors.Wrap(err, "failed to save assigned page attribute value with")
	}

	return assignedPageAttrValue, nil
}

func (as *SqlAssignedPageAttributeValueStore) Get(id string) (*model.AssignedPageAttributeValue, error) {
	var res model.AssignedPageAttributeValue

	err := as.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedPageAttributeValueTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute value with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAssignedPageAttributeValueStore) SaveInBulk(assignmentID model.UUID, attributeValueIDs []model.UUID) ([]*model.AssignedPageAttributeValue, error) {
	relations := lo.Map(attributeValueIDs, func(item model.UUID, _ int) *model.AssignedPageAttributeValue {
		return &model.AssignedPageAttributeValue{AssignmentID: assignmentID, ValueID: item}
	})

	err := as.GetMaster().Create(relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

func (as *SqlAssignedPageAttributeValueStore) SelectForSort(assignmentID string) ([]*model.AssignedPageAttributeValue, []*model.AttributeValue, error) {
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

func (as *SqlAssignedPageAttributeValueStore) UpdateInBulk(attributeValues []*model.AssignedPageAttributeValue) error {
	for _, value := range attributeValues {
		err := as.GetMaster().Raw("UPDATE "+model.AssignedPageAttributeValueTableName+" SET SortOrder = ? WHERE ValueID = ? AND AssignmentID = ?", value.SortOrder, value.ValueID, value.AssignmentID).Error
		if err != nil {
			return errors.Wrapf(err, "failed to update AssignedPageAttributeValue with ValueID = %s and AssignmentID = %s", value.ValueID, value.AssignmentID)
		}
	}

	return nil
}
