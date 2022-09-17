package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlAttributeValueStore struct {
	store.Store
}

func NewSqlAttributeValueStore(s store.Store) store.AttributeValueStore {
	return &SqlAttributeValueStore{s}
}

func (as *SqlAttributeValueStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"Name",
		"Value",
		"Slug",
		"FileUrl",
		"ContentType",
		"AttributeID",
		"RichText",
		"Boolean",
		"Datetime",
		"SortOrder",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAttributeValueStore) ScanFields(attributeValue model.AttributeValue) []interface{} {
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
	var isSaving bool

	if !model.IsValidId(av.Id) {
		av.Id = ""
		isSaving = true
		av.PreSave()
	} else {
		av.PreUpdate()
	}

	if err := av.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + store.AttributeValueTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
		_, err = as.GetMasterX().NamedExec(query, av)

	} else {
		query := "UPDATE " + store.AttributeValueTableName + " SET " +
			as.ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = as.GetMasterX().NamedExec(query, av)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"Slug", "AttributeID", "attributevalues_slug_attributeid_key"}) {
			return nil, store.NewErrInvalidInput(store.AttributeValueTableName, "Slug/AttributeID", av.Slug+"/"+av.AttributeID)
		}
		return nil, errors.Wrapf(err, "failed to upsert attribute value with id=%s", av.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("%d attribute values were/was updated instead of 1", numUpdated)
	}

	return av, nil
}

func (as *SqlAttributeValueStore) Get(id string) (*model.AttributeValue, error) {
	var res model.AttributeValue

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AttributeValueTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeValueTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find attribute value with id=%s", id)
	}

	return &res, nil
}

// FilterByOptions finds and returns all matched attribute values based on given options
func (as *SqlAttributeValueStore) FilterByOptions(options model.AttributeValueFilterOptions) (model.AttributeValues, error) {
	var executor store_iface.SqlxExecutor = as.GetReplicaX()
	if options.Transaction != nil {
		executor = options.Transaction
	}

	selectFields := as.ModelFields(store.AttributeValueTableName + ".")
	if options.SelectRelatedAttribute {
		selectFields = append(selectFields, as.Attribute().ModelFields(store.AttributeTableName+".")...)
	}

	query := as.GetQueryBuilder().
		Select(selectFields...).
		From(store.AttributeValueTableName)

	// parse options
	if options.Limit != 0 {
		query = query.Limit(options.Limit)
	}
	if options.SelectForUpdate {
		query = query.Suffix("FOR UPDATE")
	}
	if options.OrderBy != "" {
		query = query.OrderBy(options.OrderBy)
	}
	if options.SelectRelatedAttribute {
		query = query.InnerJoin(store.AttributeTableName + " ON AttributeValues.AttributeID = Attributes.Id")
	}
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.AttributeID != nil {
		query = query.Where(options.AttributeID)
	}
	if options.Extra != nil {
		query = query.Where(options.Extra)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	rows, err := executor.QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find attribute values")
	}

	var (
		res            model.AttributeValues
		attributeValue model.AttributeValue
		attr           model.Attribute
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

func (as *SqlAttributeValueStore) Delete(ids ...string) (int64, error) {
	res, err := as.GetMasterX().Exec("DELETE FROM "+store.AttributeValueTableName+" WHERE Id IN ?", ids)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete attribute values")
	}

	numDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of deleted attribute values")
	}

	return numDeleted, nil
}

func (as *SqlAttributeValueStore) BulkUpsert(transaction store_iface.SqlxTxExecutor, values model.AttributeValues) (model.AttributeValues, error) {
	var executor store_iface.SqlxExecutor = as.GetMasterX()
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
			value.PreSave()
		} else {
			value.PreUpdate()
		}

		if err := value.IsValid(); err != nil {
			return nil, err
		}

		if isSaving {
			query := "INSERT INTO " + store.AttributeValueTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
			_, err = executor.NamedExec(query, value)
		} else {
			query := "UPDATE " + store.AttributeValueTableName + " SET " + as.ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id=:Id"

			var result sql.Result

			result, err = executor.NamedExec(query, value)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
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

func (as *SqlAttributeValueStore) Count(options *model.AttributeValueFilterOptions) (int64, error) {
	query := as.GetQueryBuilder().
		Select("COUNT (*)").
		From(store.AttributeValueTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.AttributeID != nil {
		query = query.Where(options.AttributeID)
	}
	if options.Extra != nil {
		query = query.Where(options.Extra)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "AttributeValue.Count.ToSql")
	}

	var count int64
	err = as.GetReplicaX().Get(&count, queryStr, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of attribute value with given options")
	}

	return count, nil
}
