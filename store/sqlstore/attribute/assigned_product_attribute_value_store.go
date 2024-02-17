package attribute

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAssignedProductAttributeValueStore struct {
	store.Store
}

var assignedProductAttrValueDuplicateKeys = []string{
	"ValueID",
	"AssignmentID",
	"valueid_assignmentid_key",
}

func NewSqlAssignedProductAttributeValueStore(s store.Store) store.AssignedProductAttributeValueStore {
	return &SqlAssignedProductAttributeValueStore{s}
}

func (as *SqlAssignedProductAttributeValueStore) Save(assignedProductAttrValue *model.AssignedProductAttributeValue) (*model.AssignedProductAttributeValue, error) {
	if err := as.GetMaster().Create(assignedProductAttrValue).Error; err != nil {
		if as.IsUniqueConstraintError(err, assignedProductAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(model.AssignedProductAttributeValueTableName, "ValueID/AssignmentID", assignedProductAttrValue.ValueID+"/"+assignedProductAttrValue.AssignmentID)
		}
		return nil, errors.Wrap(err, "failed to save assigned product attribute value")
	}

	return assignedProductAttrValue, nil
}

func (as *SqlAssignedProductAttributeValueStore) Get(id string) (*model.AssignedProductAttributeValue, error) {
	var res model.AssignedProductAttributeValue

	err := as.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedProductAttributeValueTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute value with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAssignedProductAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*model.AssignedProductAttributeValue, error) {
	relations := lo.Map(attributeValueIDs, func(item string, _ int) *model.AssignedProductAttributeValue {
		return &model.AssignedProductAttributeValue{
			ValueID:      item,
			AssignmentID: assignmentID,
		}
	})

	err := as.GetMaster().Create(relations).Error
	if err != nil {
		if as.IsUniqueConstraintError(err, assignedProductAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(model.AssignedProductAttributeValueTableName, "ValueID/AssignmentID", "")
		}
		return nil, errors.Wrap(err, "failed to save assigned product attribute value")
	}

	return relations, nil
}

func (as *SqlAssignedProductAttributeValueStore) UpdateInBulk(attributeValues []*model.AssignedProductAttributeValue) error {
	for _, value := range attributeValues {
		err := as.GetMaster().Save(value).Error
		if err != nil {
			if as.IsUniqueConstraintError(err, assignedProductAttrValueDuplicateKeys) {
				return store.NewErrInvalidInput(model.AssignedProductAttributeValueTableName, "ValueID/AssignmentID", value.ValueID+"/"+value.AssignmentID)
			}
			return errors.Wrapf(err, "failed to update value with id=%s", value.Id)
		}
	}

	return nil
}

func (as *SqlAssignedProductAttributeValueStore) SelectForSort(assignmentID string) ([]*model.AssignedProductAttributeValue, []*model.AttributeValue, error) {
	rows, err := as.GetReplica().
		Raw("SELECT AssignedProductAttributeValues.*, AttributeValues.* FROM "+model.AssignedProductAttributeValueTableName+" INNER JOIN "+model.AttributeValueTableName+" ON AssignedProductAttributeValues.ValueID = AttributeValues.Id WHERE AssignedProductAttributeValues.AssignmentID = ?", assignmentID).
		Rows()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find assigned product attribute values")
	}
	defer rows.Close()

	var (
		assignedProductAttributeValues []*model.AssignedProductAttributeValue
		attributeValues                []*model.AttributeValue
	)

	for rows.Next() {
		var (
			assignedProductAttributeValue model.AssignedProductAttributeValue
			attributeValue                model.AttributeValue
			scanFields                    = append(as.ScanFields(&assignedProductAttributeValue), as.AttributeValue().ScanFields(&attributeValue)...)
		)

		scanErr := rows.Scan(scanFields...)
		if scanErr != nil {
			return nil, nil, errors.Wrapf(scanErr, "error scanning values")
		}

		assignedProductAttributeValues = append(assignedProductAttributeValues, &assignedProductAttributeValue)
		attributeValues = append(attributeValues, &attributeValue)
	}

	return assignedProductAttributeValues, attributeValues, nil
}

func (s *SqlAssignedProductAttributeValueStore) FilterByOptions(options *model.AssignedProductAttributeValueFilterOptions) ([]*model.AssignedProductAttributeValue, error) {
	args, err := store.BuildSqlizer(options.Conditions, "FilterByOptions")
	if err != nil {
		return nil, err
	}

	var res []*model.AssignedProductAttributeValue
	err = s.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find assigned product attributes by given options")
	}
	return res, nil
}
