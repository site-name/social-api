package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
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
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AssignedVariantAttributeTableName+" WHERE Id=:ID", map[string]interface{}{"ID": variantID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedVariantAttributeTableName, variantID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute with id=%s", variantID)
	}

	return &res, nil
}

func (as *SqlAssignedVariantAttributeStore) GetWithOption(option *attribute.AssignedVariantAttributeFilterOption) (*attribute.AssignedVariantAttribute, error) {
	if option == nil || !model.IsValidId(option.VariantID) || !model.IsValidId(option.AssignmentID) {
		return nil, store.NewErrInvalidInput(store.AssignedVariantAttributeTableName, "option", option)
	}

	var res *attribute.AssignedVariantAttribute

	tx, err := as.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// try finding first:
	err = tx.SelectOne(
		&res,
		"SELECT * FROM "+store.AssignedVariantAttributeTableName+" WHERE (VariantID = :VariantId AND AssignmentID = :AssignmentID)",
		map[string]interface{}{
			"VariantId":    option.VariantID,
			"AssignmentID": option.AssignmentID,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows { // this mean we need to create one
			res = new(attribute.AssignedVariantAttribute)
			res.AssignmentID = option.AssignmentID
			res.VariantID = option.VariantID
			res.PreSave()
			if appErr := res.IsValid(); appErr != nil {
				return nil, appErr
			}

			if err = tx.Insert(res); err != nil {
				if as.IsUniqueConstraintError(err, []string{"VariantID", "AssignmentID", strings.ToLower(store.AssignedVariantAttributeTableName) + "_variantid_assignmentid_key"}) {
					return nil, store.NewErrInvalidInput(store.AssignedProductAttributeTableName, "VariantID/AssignmentID", option.VariantID+"/"+option.AssignmentID)
				}
				return nil, errors.Wrapf(err, "failed to insert new assigned variant attribute with id=%s", res.Id)
			}
		}
		// system error
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute with VariantID = %s, AssignmentID = %s", option.VariantID, option.AssignmentID)
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit_transaction")
	}

	return res, nil
}
