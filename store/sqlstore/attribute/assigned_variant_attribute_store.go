package attribute

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAssignedVariantAttributeStore struct {
	store.Store
}

func NewSqlAssignedVariantAttributeStore(s store.Store) store.AssignedVariantAttributeStore {
	return &SqlAssignedVariantAttributeStore{s}
}

func (as *SqlAssignedVariantAttributeStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"VariantID",
		"AssignmentID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (as *SqlAssignedVariantAttributeStore) Save(variant *model.AssignedVariantAttribute) (*model.AssignedVariantAttribute, error) {
	if err := as.GetMaster().Create(variant).Error; err != nil {
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
func (as *SqlAssignedVariantAttributeStore) builFilterQuery(option *model.AssignedVariantAttributeFilterOption) (string, []interface{}, error) {
	query := as.GetQueryBuilder().
		Select("AssignedVariantAttributes.*").
		From(model.AssignedVariantAttributeTableName).
		Where(option.Conditions)

	// parse option
	if option.AssignmentAttributeInputType != nil ||
		option.Assignment_Attribute_VisibleInStoreFront != nil ||
		option.AssignmentAttributeType != nil {
		query = query.
			InnerJoin(model.AttributeVariantTableName + " ON (AssignedVariantAttributes.AssignmentID = AttributeVariants.Id)").
			InnerJoin(model.AttributeTableName + " ON (AttributeVariants.AttributeID = Attributes.Id)")

		if option.AssignmentAttributeInputType != nil {
			query = query.Where(option.AssignmentAttributeInputType)
		}
		if option.AssignmentAttributeType != nil {
			query = query.Where(option.AssignmentAttributeType)
		}
		if value := option.Assignment_Attribute_VisibleInStoreFront; value != nil {
			query = query.Where(squirrel.Eq{model.AttributeTableName + ".VisibleInStoreFront": *value})
		}
	}

	return query.ToSql()
}

// GetWithOption finds and returns 1 assigned variant attribute with given option
func (as *SqlAssignedVariantAttributeStore) GetWithOption(option *model.AssignedVariantAttributeFilterOption) (*model.AssignedVariantAttribute, error) {
	queryString, args, err := as.builFilterQuery(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetWithOption_ToSql")
	}

	var res model.AssignedVariantAttribute
	err = as.GetReplica().Raw(queryString, args...).Scan(&res).Error
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
	queryString, args, err := as.builFilterQuery(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.AssignedVariantAttribute
	err = as.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find assigned variant attributes by given option")
	}

	return res, nil
}
