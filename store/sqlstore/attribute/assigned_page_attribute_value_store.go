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

type SqlAssignedPageAttributeValueStore struct {
	store.Store
}

var (
	assignedPageAttrValueDuplicateKeys = []string{"ValueID", "AssignmentID", strings.ToLower(store.AssignedPageAttributeValueTableName) + "_valueid_assignmentid_key"}
	// "APAV" is acronym for table name. Make sure to turn `AssignedPageAttributeValueTableName` to "APAV" when building queries
	AssignedPageAttributeValueSelectList = []string{
		"APAV.Id",
		"APAV.ValueID",
		"APAV.AssignmentID",
		"APAV.SortOrder",
	}
)

func NewSqlAssignedPageAttributeValueStore(s store.Store) store.AssignedPageAttributeValueStore {
	as := &SqlAssignedPageAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedPageAttributeValue{}, store.AssignedPageAttributeValueTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedPageAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedPageAttributeValueTableName, "ValueID", store.AttributeValueTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedPageAttributeValueTableName, "AssignmentID", store.AssignedPageAttributeTableName, "Id", true)
}

func (as *SqlAssignedPageAttributeValueStore) Save(assignedPageAttrValue *attribute.AssignedPageAttributeValue) (*attribute.AssignedPageAttributeValue, error) {
	assignedPageAttrValue.PreSave()
	if err := assignedPageAttrValue.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(assignedPageAttrValue); err != nil {
		if as.IsUniqueConstraintError(err, assignedPageAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.AssignedPageAttributeValueTableName, "ValueID/AssignmentID", assignedPageAttrValue.ValueID+"/"+assignedPageAttrValue.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned page attribute value with id=%s", assignedPageAttrValue.Id)
	}

	return assignedPageAttrValue, nil
}

func (as *SqlAssignedPageAttributeValueStore) Get(assignedPageAttrValueID string) (*attribute.AssignedPageAttributeValue, error) {
	res, err := as.GetReplica().Get(attribute.AssignedPageAttributeValue{}, assignedPageAttrValueID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedPageAttributeValueTableName, assignedPageAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute value with id=%s", assignedPageAttrValueID)
	}

	return res.(*attribute.AssignedPageAttributeValue), nil
}

func (as *SqlAssignedPageAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedPageAttributeValue, error) {
	tx, err := as.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// return value:
	res := []*attribute.AssignedPageAttributeValue{}

	for _, id := range attributeValueIDs {
		newValue := &attribute.AssignedPageAttributeValue{
			ValueID:      id,
			AssignmentID: assignmentID,
		}
		newValue.PreSave()
		if appErr := newValue.IsValid(); appErr != nil {
			return nil, appErr
		}

		err = tx.Insert(newValue)
		if err != nil {
			if as.IsUniqueConstraintError(err, assignedPageAttrValueDuplicateKeys) {
				return nil, store.NewErrInvalidInput(store.AssignedPageAttributeValueTableName, "ValueID/AssignmentID", newValue.ValueID+"/"+newValue.AssignmentID)
			}
			return nil, errors.Wrapf(err, "failed to save assigned page attribute value with id=%s", newValue.Id)
		}
		// append to return value if success
		res = append(res, newValue)
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "commit_transaction")
	}

	return res, nil
}

func (as *SqlAssignedPageAttributeValueStore) SelectForSort(assignmentID string) (assignedPageAttributeValues []*attribute.AssignedPageAttributeValue, attributeValues []*attribute.AttributeValue, err error) {
	selectValues := strings.Join(
		append(AssignedPageAttributeValueSelectList, AttributeValueSelect...),
		", ",
	)
	query := `SELECT ` + selectValues + ` FROM ` +
		store.AssignedPageAttributeValueTableName + ` AS APAV INNER JOIN ` +
		store.AttributeValueTableName + ` AS AV ON(
			APAV.ValueID = AV.Id
		)
		WHERE (
			APAV.AssignmentID = :AssignmentID
		)`

	rows, err := as.GetReplica().Query(query, map[string]interface{}{"AssignmentID": assignmentID})
	if err != nil {
		if err == sql.ErrNoRows {
			err = store.NewErrNotFound(store.AssignedPageAttributeValueTableName, "AssignmentID="+assignmentID)
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
		assignedPageAttributeValue attribute.AssignedPageAttributeValue
		attributeValue             attribute.AttributeValue
	)

	for rows.Next() {
		var richText []byte

		scanErr := rows.Scan(
			&assignedPageAttributeValue.Id,
			&assignedPageAttributeValue.ValueID,
			&assignedPageAttributeValue.AssignmentID,
			&assignedPageAttributeValue.SortOrder,

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

		assignedPageAttributeValues = append(assignedPageAttributeValues, &assignedPageAttributeValue)
		attributeValues = append(attributeValues, &attributeValue)
	}

	return
}

func (as *SqlAssignedPageAttributeValueStore) UpdateInBulk(attributeValues []*attribute.AssignedPageAttributeValue) error {
	tx, err := as.GetMaster().Begin()
	if err != nil {
		return errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	for _, value := range attributeValues {
		// try validating if the value exist:
		_, err := tx.Get(attribute.AssignedPageAttributeValue{}, value.Id)
		if err != nil {
			return errors.Wrapf(err, "failed to find value with id=%s", value.Id)
		}
		numUpdated, err := tx.Update(value)
		if err != nil {
			// check if error is duplicate conflict error:
			if as.IsUniqueConstraintError(err, assignedPageAttrValueDuplicateKeys) {
				return store.NewErrInvalidInput(store.AssignedPageAttributeValueTableName, "ValueID/AssignmentID", value.ValueID+"/"+value.AssignmentID)
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
