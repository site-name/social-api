package attribute

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
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

func (as *SqlAssignedVariantAttributeValueStore) ScanFields(assignedVariantAttributeValue *model.AssignedVariantAttributeValue) []interface{} {
	return []interface{}{
		&assignedVariantAttributeValue.Id,
		&assignedVariantAttributeValue.ValueID,
		&assignedVariantAttributeValue.AssignmentID,
		&assignedVariantAttributeValue.SortOrder,
	}
}

func (as *SqlAssignedVariantAttributeValueStore) Save(assignedVariantAttrValue *model.AssignedVariantAttributeValue) (*model.AssignedVariantAttributeValue, error) {
	err := as.GetMaster().Create(assignedVariantAttrValue).Error
	if err != nil {
		if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(model.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", assignedVariantAttrValue.ValueID+"/"+assignedVariantAttrValue.AssignmentID)
		}
		return nil, errors.Wrap(err, "failed to save assigned variant attribute")
	}

	return assignedVariantAttrValue, nil
}

func (as *SqlAssignedVariantAttributeValueStore) Get(id string) (*model.AssignedVariantAttributeValue, error) {
	var res model.AssignedVariantAttributeValue

	if err := as.GetReplica().First(&res, "Id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedVariantAttributeValueTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute value with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAssignedVariantAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*model.AssignedVariantAttributeValue, error) {
	relations := lo.Map(attributeValueIDs, func(item string, _ int) *model.AssignedVariantAttributeValue {
		return &model.AssignedVariantAttributeValue{
			ValueID:      item,
			AssignmentID: assignmentID,
		}
	})
	err := as.GetMaster().Create(relations).Error
	if err != nil {
		if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(model.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", "")
		}
		return nil, errors.Wrapf(err, "failed to save assigned variant attribute values")
	}

	return relations, nil
}

func (as *SqlAssignedVariantAttributeValueStore) SelectForSort(assignmentID string) ([]*model.AssignedVariantAttributeValue, []*model.AttributeValue, error) {
	rows, err := as.GetReplica().
		Raw("SELECT AssignedVariantAttributeValues.*, AttributeValues.* FROM "+model.AssignedVariantAttributeValueTableName+" INNER JOIN "+model.AttributeValueTableName+" ON (AttributeValues.Id = AssignedVariantAttributeValues.ValueID) WHERE AssignedVariantAttributeValues.AssignmentID = ?", assignmentID).Rows()
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

	return assignedVariantAttributeValues, attributeValues, nil
}

func (as *SqlAssignedVariantAttributeValueStore) UpdateInBulk(attributeValues []*model.AssignedVariantAttributeValue) error {
	for _, value := range attributeValues {
		err := as.GetMaster().Save(value).Error
		if err != nil {
			if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
				return store.NewErrInvalidInput(model.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", value.ValueID+"/"+value.AssignmentID)
			}
			return errors.Wrapf(err, "failed to update value with id=%s", value.Id)
		}
	}

	return nil
}

func (s *SqlAssignedVariantAttributeValueStore) FilterByOptions(options *model.AssignedVariantAttributeValueFilterOptions) ([]*model.AssignedVariantAttributeValue, error) {
	args, err := store.BuildSqlizer(options.Conditions, "FilterByOptions")
	if err != nil {
		return nil, err
	}

	var res []*model.AssignedVariantAttributeValue
	err = s.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find assigned variant attribute values by given options")
	}

	return res, nil
}
