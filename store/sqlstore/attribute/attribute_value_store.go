package attribute

import (
	"database/sql"
	"strings"

	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeValueStore struct {
	store.Store
}

func NewSqlAttributeValueStore(s store.Store) store.AttributeValueStore {
	as := &SqlAttributeValueStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributeValue{}, store.AttributeValueTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(attribute.ATTRIBUTE_VALUE_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(attribute.ATTRIBUTE_VALUE_SLUG_MAX_LENGTH)
		table.ColMap("Value").SetMaxSize(attribute.ATTRIBUTE_VALUE_VALUE_MAX_LENGTH)
		table.ColMap("ContentType").SetMaxSize(attribute.ATTRIBUTE_VALUE_CONTENT_TYPE_MAX_LENGTH)
		table.ColMap("FileUrl").SetMaxSize(model.URL_LINK_MAX_LENGTH)

		table.SetUniqueTogether("Slug", "AttributeID")
	}
	return as
}

func (as *SqlAttributeValueStore) ModelFields() []string {
	return []string{
		"AttributeValues.Id",
		"AttributeValues.Name",
		"AttributeValues.Value",
		"AttributeValues.Slug",
		"AttributeValues.FileUrl",
		"AttributeValues.ContentType",
		"AttributeValues.AttributeID",
		"AttributeValues.RichText",
		"AttributeValues.Boolean",
		"AttributeValues.Datetime",
		"AttributeValues.SortOrder",
	}
}

func (as *SqlAttributeValueStore) ScanFields(attributeValue attribute.AttributeValue) []interface{} {
	return []interface{}{
		&attributeValue.Id,
		&attributeValue.Name,
		&attributeValue.Value,
		&attributeValue.Slug,
		&attributeValue.FileUrl,
		&attributeValue.ContentType,
		&attributeValue.AttributeID,
		&attributeValue.RichText,
		&attributeValue.Boolean,
		&attributeValue.Datetime,
		&attributeValue.SortOrder,
	}
}

func (as *SqlAttributeValueStore) CreateIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_attributevalues_name", store.AttributeValueTableName, "Name")
	as.CreateIndexIfNotExists("idx_attributevalues_slug", store.AttributeValueTableName, "Slug")
	as.CreateIndexIfNotExists("idx_attributevalues_name_lower_textpattern", store.AttributeValueTableName, "lower(Name) text_pattern_ops")

	as.CreateForeignKeyIfNotExists(store.AttributeValueTableName, "AttributeID", store.AttributeTableName, "Id", true)
}

func (as *SqlAttributeValueStore) Upsert(av *attribute.AttributeValue) (*attribute.AttributeValue, error) {
	var isSaving bool

	if !model.IsValidId(av.Id) {
		av.Id = ""
		isSaving = true
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		av.PreSave()
		if err := av.IsValid(); err != nil {
			return nil, err
		}

		err = as.GetMaster().Insert(av)

	} else {
		av.PreUpdate()
		if err := av.IsValid(); err != nil {
			return nil, err
		}

		numUpdated, err = as.GetMaster().Update(av)
	}

	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"Slug", "AttributeID", strings.ToLower(store.AttributeValueTableName) + "_slug_attributeid_key"}) {
			return nil, store.NewErrInvalidInput(store.AttributeValueTableName, "Slug/AttributeID", av.Slug+"/"+av.AttributeID)
		}
		return nil, errors.Wrapf(err, "failed to upsert attribute value with id=%s", av.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("%d attribute values were/was updated instead of 1", numUpdated)
	}

	return av, nil
}

func (as *SqlAttributeValueStore) Get(id string) (*attribute.AttributeValue, error) {
	var res attribute.AttributeValue
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AttributeValueTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeValueTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find attribute value with id=%s", id)
	}

	return &res, nil
}

// FilterByOptions finds and returns all matched attribute values based on given options
func (as *SqlAttributeValueStore) FilterByOptions(options attribute.AttributeValueFilterOptions) (attribute.AttributeValues, error) {
	var executor gorp.SqlExecutor = as.GetReplica()
	if options.Transaction != nil {
		executor = options.Transaction
	}

	selectFields := as.ModelFields()
	if options.SelectRelatedAttribute {
		selectFields = append(selectFields, as.Attribute().ModelFields()...)
	}

	query := as.GetQueryBuilder().
		Select(selectFields...).
		From(store.AttributeValueTableName)

	// parse options
	if options.SelectForUpdate {
		query = query.Suffix("FOR UPDATE")
	}
	if options.OrderBy != "" {
		query = query.OrderBy(options.OrderBy)
	}
	if options.SelectRelatedAttribute {
		query = query.InnerJoin(store.AttributeTableName + " ON AttributeValues.AttributeID = Attributes.Id")
	}

	if !options.All { // check if we need to filter further
		if options.Id != nil {
			query = query.Where(options.Id)
		}
		if options.AttributeID != nil {
			query = query.Where(options.AttributeID)
		}
		if options.Extra != nil {
			query = query.Where(options.Extra)
		}
	}

	rows, err := query.RunWith(executor).Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find attribute values")
	}

	var (
		res attribute.AttributeValues

		attributeValue attribute.AttributeValue
		attr           attribute.Attribute
		scanFields     = as.ScanFields(attributeValue)
	)
	if options.SelectRelatedAttribute {
		scanFields = append(scanFields, as.Attribute().ScanFields(attr)...)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of attribute value")
		}

		// don't worry when we assign directly value here.
		// The Attribute will be deep copied later
		if options.SelectRelatedAttribute {
			attributeValue.Attribute = &attr
		}

		res = append(res, attributeValue.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows of attribute values")
	}

	return res, nil
}

func (as *SqlAttributeValueStore) Delete(id string) error {
	res, err := as.GetMaster().Exec("DELETE FROM "+store.AttributeValueTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		return errors.Wrapf(err, "failed to delete attribute value with id=%s", id)
	}

	numDeleted, _ := res.RowsAffected()
	if numDeleted != 1 {
		return errors.Errorf("%d attribute value(s) was/were deleted instead of 1", numDeleted)
	}

	return nil
}

func (as *SqlAttributeValueStore) BulkUpsert(transaction *gorp.Transaction, values attribute.AttributeValues) (attribute.AttributeValues, error) {
	var executor gorp.SqlExecutor = as.GetMaster()
	if transaction != nil {
		executor = transaction
	}

	for _, value := range values {

		var (
			isSaving   bool
			err        error
			numUpdated int64
		)

		if !model.IsValidId(value.Id) {
			isSaving = true
			value.Id = ""
		}

		if isSaving {
			value.PreSave()
			if err := value.IsValid(); err != nil {
				return nil, err
			}

			err = executor.Insert(value)
		} else {
			value.PreUpdate()
			if err := value.IsValid(); err != nil {
				return nil, err
			}

			numUpdated, err = executor.Update(value)
		}

		if err != nil {
			if as.IsUniqueConstraintError(err, []string{"Slug", "AttributeID", strings.ToLower(store.AttributeValueTableName) + "_slug_attributeid_key"}) {
				return nil, store.NewErrInvalidInput(store.AttributeValueTableName, "Slug/AttributeID", value.Slug+"/"+value.AttributeID)
			}
			return nil, errors.Wrapf(err, "failed to upsert attribute value with id=%s", value.Id)
		}

		if numUpdated != 1 {
			return nil, errors.Errorf("%d attribute value(1) was/were updated instead of 1", numUpdated)
		}
	}

	return values, nil
}
