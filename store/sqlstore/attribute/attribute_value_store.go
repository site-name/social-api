package attribute

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAttributeValueStore struct {
	store.Store
}

func NewSqlAttributeValueStore(s store.Store) store.AttributeValueStore {
	return &SqlAttributeValueStore{s}
}

func (as *SqlAttributeValueStore) ScanFields(attributeValue *model.AttributeValue) []interface{} {
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

func (as *SqlAttributeValueStore) Upsert(av *model.AttributeValue) (*model.AttributeValue, error) {
	err := as.GetMaster().Save(av).Error
	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"Slug", "AttributeID", "attributevalues_slug_attributeid_key"}) {
			return nil, store.NewErrInvalidInput(model.AttributeValueTableName, "Slug/AttributeID", av.Slug+"/"+av.AttributeID)
		}
		return nil, errors.Wrapf(err, "failed to upsert attribute value with id=%s", av.Id)
	}

	return av, nil
}

func (as *SqlAttributeValueStore) Get(id string) (*model.AttributeValue, error) {
	var res model.AttributeValue
	err := as.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AttributeValueTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find attribute value with id=%s", id)
	}

	return &res, nil
}

// FilterByOptions finds and returns all matched attribute values based on given options
func (as *SqlAttributeValueStore) FilterByOptions(options model.AttributeValueFilterOptions) (model.AttributeValues, error) {
	var executor *gorm.DB = as.GetMaster()
	if options.Transaction != nil {
		executor = options.Transaction
	}

	selectFields := []string{model.AttributeValueTableName + ".*"}
	if options.SelectRelatedAttribute {
		selectFields = append(selectFields, model.AttributeTableName+".*")
	}

	query := as.GetQueryBuilder().
		Select(selectFields...).
		From(model.AttributeValueTableName).Where(options.Conditions)

	if options.SelectForUpdate && options.Transaction != nil {
		query = query.Suffix("FOR UPDATE")
	}
	if options.Ordering != "" {
		query = query.OrderBy(options.Ordering)
	}
	if options.SelectRelatedAttribute {
		query = query.InnerJoin(model.AttributeTableName + " ON AttributeValues.AttributeID = Attributes.Id")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	rows, err := executor.Raw(queryString, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find attribute values")
	}
	defer rows.Close()

	var res model.AttributeValues

	for rows.Next() {
		var (
			attributeValue model.AttributeValue
			attribute      model.Attribute
			scanFields     = as.ScanFields(&attributeValue)
		)
		if options.SelectRelatedAttribute {
			scanFields = append(scanFields, as.Attribute().ScanFields(&attribute)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of attribute value")
		}

		if options.SelectRelatedAttribute {
			attributeValue.Attribute = &attribute
		}

		res = append(res, &attributeValue)
	}

	return res, nil
}

func (as *SqlAttributeValueStore) Delete(ids ...string) (int64, error) {
	result := as.GetMaster().Raw("DELETE FROM "+model.AttributeValueTableName+" WHERE Id IN ?", ids)
	if result.Error != nil {
		return 0, errors.Wrap(result.Error, "failed to delete attribute values")
	}

	return result.RowsAffected, nil
}

func (as *SqlAttributeValueStore) BulkUpsert(transaction *gorm.DB, values model.AttributeValues) (model.AttributeValues, error) {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	for _, value := range values {
		err := transaction.Save(value).Error
		if err != nil {
			if as.IsUniqueConstraintError(err, []string{"Slug", "AttributeID", strings.ToLower(model.AttributeValueTableName) + "_slug_attributeid_key"}) {
				return nil, store.NewErrInvalidInput(model.AttributeValueTableName, "Slug/AttributeID", value.Slug+"/"+value.AttributeID)
			}
			return nil, errors.Wrapf(err, "failed to upsert attribute value with id=%s", value.Id)
		}
	}

	return values, nil
}

func (as *SqlAttributeValueStore) Count(options *model.AttributeValueFilterOptions) (int64, error) {
	query, args, err := as.GetQueryBuilder().
		Select("COUNT (*)").
		From(model.AttributeValueTableName).Where(options.Conditions).ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "AttributeValue.Count.ToSql")
	}

	var count int64
	err = as.GetReplica().Raw(query, args...).Scan(&count).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of attribute value with given options")
	}

	return count, nil
}
