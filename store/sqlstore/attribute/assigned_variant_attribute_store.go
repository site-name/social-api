package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedVariantAttributeStore struct {
	store.Store
}

func NewSqlAssignedVariantAttributeStore(s store.Store) store.AssignedVariantAttributeStore {
	as := &SqlAssignedVariantAttributeStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedVariantAttribute{}, store.AssignedVariantAttributeTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("VariantID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedVariantAttributeStore) ModelFields() []string {
	return []string{
		"AssignedVariantAttributes.Id",
		"AssignedVariantAttributes.VariantID",
		"AssignedVariantAttributes.AssignmentID",
	}
}

func (as *SqlAssignedVariantAttributeStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedVariantAttributeTableName, "VariantID", store.ProductVariantTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedVariantAttributeTableName, "AssignmentID", store.AttributeVariantTableName, "Id", true)
}

func (as *SqlAssignedVariantAttributeStore) Save(variant *attribute.AssignedVariantAttribute) (*attribute.AssignedVariantAttribute, error) {
	variant.PreSave()
	if err := variant.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(variant); err != nil {
		if as.IsUniqueConstraintError(err, []string{"VariantID", "AssignmentID", strings.ToLower(store.AssignedVariantAttributeTableName) + "_variantid_assignmentid_key"}) {
			return nil, store.NewErrInvalidInput(store.AssignedVariantAttributeTableName, "VariantID/AssignmentID", variant.VariantID+"/"+variant.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned variant attribute with id=%s", variant.Id)
	}

	return variant, nil
}

func (as *SqlAssignedVariantAttributeStore) Get(variantID string) (*attribute.AssignedVariantAttribute, error) {
	var res attribute.AssignedVariantAttribute
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AssignedVariantAttributeTableName+" WHERE Id = :ID", map[string]interface{}{"ID": variantID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedVariantAttributeTableName, variantID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute with id=%s", variantID)
	}

	return &res, nil
}

// builFilterQuery is common method for building filter queries
func (as *SqlAssignedVariantAttributeStore) builFilterQuery(option *attribute.AssignedVariantAttributeFilterOption) (string, []interface{}, error) {
	query := as.GetQueryBuilder().
		Select(as.ModelFields()...).
		From(store.AssignedVariantAttributeTableName)

	// parse option
	if option.AssignmentID != nil {
		query = query.Where(option.AssignmentID.ToSquirrel("AssignmentID"))
	}
	if option.VariantID != nil {
		query = query.Where(option.VariantID.ToSquirrel("VariantID"))
	}
	var joined_AssignedVariantAttributes_and_Attributes_tables bool
	if option.AssignmentAttributeInputType != nil {
		query = query.
			InnerJoin(store.AttributeVariantTableName + " ON (AssignedVariantAttributes.AssignmentID = AttributeVariants.Id)").
			InnerJoin(store.AttributeTableName + " ON (AttributeVariants.AttributeID = Attributes.Id)").
			Where(option.AssignmentAttributeInputType.ToSquirrel("Attributes.InputType"))

		joined_AssignedVariantAttributes_and_Attributes_tables = true // indicate that already joined 2 tables
	}
	if option.AssignmentAttributeType != nil {
		if !joined_AssignedVariantAttributes_and_Attributes_tables {
			query = query.
				InnerJoin(store.AttributeVariantTableName + " ON (AssignedVariantAttributes.AssignmentID = AttributeVariants.Id)").
				InnerJoin(store.AttributeTableName + " ON (AttributeVariants.AttributeID = Attributes.Id)")
		}
		query = query.Where(option.AssignmentAttributeType.ToSquirrel("Attributes.Type"))
	}

	return query.ToSql()
}

// GetWithOption finds and returns 1 assigned variant attribute with given option
func (as *SqlAssignedVariantAttributeStore) GetWithOption(option *attribute.AssignedVariantAttributeFilterOption) (*attribute.AssignedVariantAttribute, error) {
	queryString, args, err := as.builFilterQuery(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetWithOption_ToSql")
	}

	var res attribute.AssignedVariantAttribute
	err = as.GetReplica().SelectOne(
		&res,
		queryString,
		args...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedVariantAttributeTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute with VariantID = %s, AssignmentID = %s", option.VariantID, option.AssignmentID)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of assigned variant attributes filtered by given options
func (as *SqlAssignedVariantAttributeStore) FilterByOption(option *attribute.AssignedVariantAttributeFilterOption) ([]*attribute.AssignedVariantAttribute, error) {

	queryString, args, err := as.builFilterQuery(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*attribute.AssignedVariantAttribute
	_, err = as.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find assigned variant attributes by given option")
	}

	return res, nil
}
