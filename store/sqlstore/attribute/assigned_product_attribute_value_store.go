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

type SqlAssignedProductAttributeValueStore struct {
	store.Store
}

var (
	assignedProductAttrValueDuplicateKeys = []string{
		"ValueID",
		"AssignmentID",
		strings.ToLower(store.AssignedProductAttributeValueTableName) + "_valueid_assignmentid_key",
	}
	// prefixes "APAV" stand for the table name. When building queries using this variable, please make acronyms correctly
	AssignedProductAttributeValueSelectList = []string{
		"APAV.Id",
		"APAV.ValueID",
		"APAV.AssignmentID",
		"APAV.SortOrder",
	}
)

func NewSqlAssignedProductAttributeValueStore(s store.Store) store.AssignedProductAttributeValueStore {
	as := &SqlAssignedProductAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedProductAttributeValue{}, store.AssignedProductAttributeValueTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ValueID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ValueID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedProductAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedProductAttributeValueTableName, "ValueID", store.AttributeValueTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedProductAttributeValueTableName, "AssignmentID", store.AssignedProductAttributeTableName, "Id", true)
}

func (as *SqlAssignedProductAttributeValueStore) Save(assignedProductAttrValue *attribute.AssignedProductAttributeValue) (*attribute.AssignedProductAttributeValue, error) {
	assignedProductAttrValue.PreSave()
	if err := assignedProductAttrValue.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(assignedProductAttrValue); err != nil {
		if as.IsUniqueConstraintError(err, assignedProductAttrValueDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.AssignedProductAttributeValueTableName, "ValueID/AssignmentID", assignedProductAttrValue.ValueID+"/"+assignedProductAttrValue.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned product attribute value with id=%s", assignedProductAttrValue.Id)
	}

	return assignedProductAttrValue, nil
}

func (as *SqlAssignedProductAttributeValueStore) Get(assignedProductAttrValueID string) (*attribute.AssignedProductAttributeValue, error) {
	res, err := as.GetReplica().Get(attribute.AssignedProductAttributeValue{}, assignedProductAttrValueID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedProductAttributeValueTableName, assignedProductAttrValueID)
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute value with id=%s", assignedProductAttrValueID)
	}

	return res.(*attribute.AssignedProductAttributeValue), nil
}

func (as *SqlAssignedProductAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedProductAttributeValue, error) {
	tx, err := as.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrapf(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	// return value:
	res := []*attribute.AssignedProductAttributeValue{}

	for _, id := range attributeValueIDs {
		newValue := &attribute.AssignedProductAttributeValue{
			ValueID:      id,
			AssignmentID: assignmentID,
		}
		newValue.PreSave()
		if appErr := newValue.IsValid(); appErr != nil {
			return nil, appErr
		}

		err = tx.Insert(newValue)
		if err != nil {
			if as.IsUniqueConstraintError(err, assignedProductAttrValueDuplicateKeys) {
				return nil, store.NewErrInvalidInput(store.AssignedProductAttributeValueTableName, "ValueID/AssignmentID", newValue.ValueID+"/"+newValue.AssignmentID)
			}
			return nil, errors.Wrapf(err, "failed to save assigned product attribute value with id=%s", newValue.Id)
		}
		// append to return value if success
		res = append(res, newValue)
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "commit_transaction")
	}

	return res, nil
}

func (as *SqlAssignedProductAttributeValueStore) SelectForSort(assignmentID string) (assignedProductAttributeValues []*attribute.AssignedProductAttributeValue, attributeValues []*attribute.AttributeValue, err error) {
	selectValues := strings.Join(
		append(AssignedProductAttributeValueSelectList, AttributeValueSelect...),
		", ",
	)
	query := `SELECT ` + selectValues + ` FROM ` +
		store.AssignedProductAttributeValueTableName + ` AS APAV INNER JOIN ` +
		store.AttributeValueTableName + ` AS AV ON(
			APAV.ValueID = AV.Id
		)
		WHERE (
			APAV.AssignmentID = :AssignmentID
		)`

	rows, err := as.GetReplica().Query(query, map[string]interface{}{"AssignmentID": assignmentID})
	if err != nil {
		if err == sql.ErrNoRows {
			err = store.NewErrNotFound(store.AssignedProductAttributeValueTableName, "AssignmentID="+assignmentID)
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
		assignedProductAttributeValue attribute.AssignedProductAttributeValue
		attributeValue                attribute.AttributeValue
	)

	for rows.Next() {
		var richText []byte

		scanErr := rows.Scan(
			&assignedProductAttributeValue.Id,
			&assignedProductAttributeValue.ValueID,
			&assignedProductAttributeValue.AssignmentID,
			&assignedProductAttributeValue.SortOrder,

			&attributeValue.Id,
			&attributeValue.Name,
			&attributeValue.Value,
			&attributeValue.Slug,
			&attributeValue.FileUrl,
			&attributeValue.ContentType,
			&attributeValue.AttributeID,
			&richText, // NOTE this is because Scan() may not supports parsing map[string]interface{}
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

		assignedProductAttributeValues = append(assignedProductAttributeValues, &assignedProductAttributeValue)
		attributeValues = append(attributeValues, &attributeValue)
	}

	return
}
