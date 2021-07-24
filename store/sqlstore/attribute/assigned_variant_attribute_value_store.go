package attribute

import (
	"bytes"
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedVariantAttributeValueStore struct {
	store.Store
}

var (
	assignedVariantAttrValueDuplicateKeys = []string{"ValueID", "AssignmentID", strings.ToLower(store.AssignedVariantAttributeValueTableName) + "_valueid_assignmentid_key"}
	// "AVAV" part is acronym for the table name. Make sure to make alias when building query with table name
	AssignedVariantAttributeValueSelectList = []string{
		"AVAV.Id",
		"AVAV.ValueID",
		"AVAV.AssignmentID",
		"AVAV.SortOrder",
	}
)

func NewSqlAssignedVariantAttributeValueStore(s store.Store) store.AssignedVariantAttributeValueStore {
	as := &SqlAssignedVariantAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedVariantAttributeValue{}, store.AssignedVariantAttributeValueTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedVariantAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedVariantAttributeValueTableName, "ValueID", store.AttributeValueTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedVariantAttributeValueTableName, "AssignmentID", store.AssignedVariantAttributeTableName, "Id", true)
}

func (as *SqlAssignedVariantAttributeValueStore) Save(assignedVariantAttrValue *attribute.AssignedVariantAttributeValue) (*attribute.AssignedVariantAttributeValue, error) {
	assignedVariantAttrValue.PreSave()
	if err := assignedVariantAttrValue.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(assignedVariantAttrValue); err != nil {
		if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", assignedVariantAttrValue.ValueID+"/"+assignedVariantAttrValue.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned variant attribute value with id=%s", assignedVariantAttrValue.Id)
	}

	return assignedVariantAttrValue, nil
}

func (as *SqlAssignedVariantAttributeValueStore) Get(assignedVariantAttrValueID string) (*attribute.AssignedVariantAttributeValue, error) {
	res, err := as.GetReplica().Get(attribute.AssignedVariantAttributeValue{}, assignedVariantAttrValueID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedVariantAttributeValueTableName, assignedVariantAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned variant attribute value with id=%s", assignedVariantAttrValueID)
	}

	return res.(*attribute.AssignedVariantAttributeValue), nil
}

func (as *SqlAssignedVariantAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedVariantAttributeValue, error) {
	tx, err := as.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// return value:
	res := []*attribute.AssignedVariantAttributeValue{}

	for _, id := range attributeValueIDs {
		newValue := &attribute.AssignedVariantAttributeValue{
			ValueID:      id,
			AssignmentID: assignmentID,
		}
		newValue.PreSave()
		if appErr := newValue.IsValid(); appErr != nil {
			return nil, appErr
		}

		err = tx.Insert(newValue)
		if err != nil {
			if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
				return nil, store.NewErrInvalidInput(store.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", newValue.ValueID+"/"+newValue.AssignmentID)
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

func (as *SqlAssignedVariantAttributeValueStore) SelectForSort(assignmentID string) (assignedVariantAttributeValues []*attribute.AssignedVariantAttributeValue, attributeValues []*attribute.AttributeValue, err error) {
	selectValues := strings.Join(
		append(AssignedVariantAttributeValueSelectList, AttributeValueSelect...),
		", ",
	)
	query := `SELECT ` + selectValues + ` FROM ` +
		store.AssignedVariantAttributeValueTableName + ` AS AVAV INNER JOIN ` +
		store.AttributeValueTableName + ` AS AV ON(
			AVAV.ValueID = AV.Id
		)
		WHERE (
			AVAV.AssignmentID = :AssignmentID
		)`

	rows, err := as.GetReplica().Query(query, map[string]interface{}{"AssignmentID": assignmentID})
	if err != nil {
		if err == sql.ErrNoRows {
			err = store.NewErrNotFound(store.AssignedVariantAttributeValueTableName, "AssignmentID="+assignmentID)
			return
		}
		err = errors.Wrapf(err, "failed to find values with AssignmentID=%s", assignmentID)
		return
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if err == nil {
				err = errors.Wrap(closeErr, "error closing rows")
				return
			}
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			if err == nil {
				err = errors.Wrap(rowsErr, "rows error")
				return
			}
		}
	}()

	var (
		assignedVariantAttributeValue attribute.AssignedVariantAttributeValue
		attributeValue                attribute.AttributeValue
	)

	for rows.Next() {
		var richText []byte

		scanErr := rows.Scan(
			&assignedVariantAttributeValue.Id,
			&assignedVariantAttributeValue.ValueID,
			&assignedVariantAttributeValue.AssignmentID,
			&assignedVariantAttributeValue.SortOrder,

			&attributeValue.Id,
			&attributeValue.Name,
			&attributeValue.Value,
			&attributeValue.Slug,
			&attributeValue.FileUrl,
			&attributeValue.ContentType,
			&attributeValue.AttributeID,
			&richText, // NOTE this is because Scan() may not support parsing map[string]interface{}
			&attributeValue.Boolean,
			&attributeValue.SortOrder,
		)
		if scanErr != nil {
			err = errors.Wrapf(scanErr, "error scanning values")
			return
		}

		parseErr := model.ModelFromJson(&attributeValue.RichText, bytes.NewReader(richText))
		if parseErr != nil {
			err = parseErr
			return
		}

		assignedVariantAttributeValues = append(assignedVariantAttributeValues, &assignedVariantAttributeValue)
		attributeValues = append(attributeValues, &attributeValue)
	}

	return
}

func (as *SqlAssignedVariantAttributeValueStore) UpdateInBulk(attributeValues []*attribute.AssignedVariantAttributeValue) error {
	tx, err := as.GetMaster().Begin()
	if err != nil {
		return errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	for _, value := range attributeValues {
		// try validating if the value exist:
		_, err := tx.Get(attribute.AssignedVariantAttributeValue{}, value.Id)
		if err != nil {
			return errors.Wrapf(err, "failed to find value with id=%s", value.Id)
		}
		numUpdated, err := tx.Update(value)
		if err != nil {
			// check if error is duplicate conflict error:
			if as.IsUniqueConstraintError(err, assignedVariantAttrValueDuplicateKeys) {
				return store.NewErrInvalidInput(store.AssignedVariantAttributeValueTableName, "ValueID/AssignmentID", value.ValueID+"/"+value.AssignmentID)
			}
			return errors.Wrapf(err, "failed to update value with id=%s", value.Id)
		}
		if numUpdated > 1 {
			return errors.Errorf("more than one value with id=%s were updated(%d)", value.Id, numUpdated)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}

	return nil
}
