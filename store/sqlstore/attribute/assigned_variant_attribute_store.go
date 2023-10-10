package attribute

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAssignedVariantAttributeStore struct {
	store.Store
}

func NewSqlAssignedVariantAttributeStore(s store.Store) store.AssignedVariantAttributeStore {
	return &SqlAssignedVariantAttributeStore{s}
}

func (as *SqlAssignedVariantAttributeStore) Save(variant *model.AssignedVariantAttribute) (*model.AssignedVariantAttribute, error) {
	if err := as.GetMaster().Save(variant).Error; err != nil {
		if as.IsUniqueConstraintError(err, []string{"VariantID", "AssignmentID", strings.ToLower(model.AssignedVariantAttributeTableName) + "_variantid_assignmentid_key"}) {
			return nil, store.NewErrInvalidInput(model.AssignedVariantAttributeTableName, "VariantID/AssignmentID", variant.VariantID+"/"+variant.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned variant attribute with id=%s", variant.Id)
	}

	return variant, nil
}

func (as *SqlAssignedVariantAttributeStore) Get(variantID string) (*model.AssignedVariantAttribute, error) {
	var res model.AssignedVariantAttribute

	err := as.GetReplica().First(&res, "Id = ?", variantID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedVariantAttributeTableName, variantID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute with id=%s", variantID)
	}

	return &res, nil
}

// builFilterQuery is common method for building filter queries
func (as *SqlAssignedVariantAttributeStore) builFilterQuery(option *model.AssignedVariantAttributeFilterOption) (*gorm.DB, squirrel.Sqlizer) {
	db := as.GetReplica()
	if option == nil {
		return db, nil
	}

	conditions := squirrel.And{}
	if option.Conditions != nil {
		conditions = append(conditions, option.Conditions)
	}

	for _, preload := range option.Preloads {
		db = db.Preload(preload)
	}

	if option.Assignment_Conditions != nil || option.Assignment_Attribute_Conditions != nil {
		conditions = append(conditions, option.Assignment_Conditions)
		db = db.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.AttributeVariantTableName,                  // 1
				model.AssignedVariantAttributeTableName,          // 2
				model.AttributeVariantColumnId,                   // 3
				model.AssignedVariantAttributeColumnAssignmentID, // 4
			),
		)

		if option.Assignment_Attribute_Conditions != nil {
			conditions = append(conditions, option.Assignment_Attribute_Conditions)
			db = db.Joins(
				fmt.Sprintf(
					"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.AttributeTableName,                // 1
					model.AttributeVariantTableName,         // 2
					model.AttributeColumnId,                 // 3
					model.AttributeVariantColumnAttributeID, // 4
				),
			)
		}
	}

	return db, conditions
}

// GetWithOption finds and returns 1 assigned variant attribute with given option
func (as *SqlAssignedVariantAttributeStore) GetWithOption(option *model.AssignedVariantAttributeFilterOption) (*model.AssignedVariantAttribute, error) {
	db, conditions := as.builFilterQuery(option)

	var res model.AssignedVariantAttribute
	err := db.First(&res, store.BuildSqlizer(conditions)...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedVariantAttributeTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute with given options")
	}

	return &res, nil
}

// FilterByOption finds and returns a list of assigned variant attributes filtered by given options
func (as *SqlAssignedVariantAttributeStore) FilterByOption(option *model.AssignedVariantAttributeFilterOption) ([]*model.AssignedVariantAttribute, error) {
	db, conditions := as.builFilterQuery(option)

	var res []*model.AssignedVariantAttribute
	err := db.Find(&res, store.BuildSqlizer(conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find assigned variant attributes by given option")
	}

	return res, nil
}
