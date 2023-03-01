package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlAttributeStore struct {
	store.Store
}

func NewSqlAttributeStore(s store.Store) store.AttributeStore {
	return &SqlAttributeStore{s}
}

func (as *SqlAttributeStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"Slug",
		"Name",
		"Type",
		"InputType",
		"EntityType",
		"Unit",
		"ValueRequired",
		"IsVariantOnly",
		"VisibleInStoreFront",
		"FilterableInStorefront",
		"FilterableInDashboard",
		"StorefrontSearchPosition",
		"AvailableInGrid",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}
	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAttributeStore) ScanFields(v *model.Attribute) []interface{} {
	return []interface{}{
		&v.Id,
		&v.Slug,
		&v.Name,
		&v.Type,
		&v.InputType,
		&v.EntityType,
		&v.Unit,
		&v.ValueRequired,
		&v.IsVariantOnly,
		&v.VisibleInStoreFront,
		&v.FilterableInStorefront,
		&v.FilterableInDashboard,
		&v.StorefrontSearchPosition,
		&v.AvailableInGrid,
		&v.Metadata,
		&v.PrivateMetadata,
	}
}

// Upsert inserts or updates given attribute then returns it
func (as *SqlAttributeStore) Upsert(attr *model.Attribute) (*model.Attribute, error) {
	var isSaving bool

	if !model.IsValidId(attr.Id) {
		attr.Id = ""
		isSaving = true
		attr.PreSave()
	} else {
		attr.PreUpdate()
	}

	if err := attr.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)

	if isSaving {
		query := "INSERT INTO " + store.AttributeTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
		_, err = as.GetMasterX().NamedExec(query, attr)

	} else {
		query := "UPDATE " + store.AttributeTableName + " SET " +
			as.ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = as.GetMasterX().NamedExec(query, attr)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"Slug", "attributes_slug_key", "idx_attributes_slug_unique"}) {
			return nil, store.NewErrInvalidInput(store.AttributeTableName, "slug", attr.Slug)
		}

		return nil, errors.Wrap(err, "failed to upsert attribute")
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("%d attribute(s) was/were updated instead of 1", numUpdated)
	}

	return attr, nil
}

func (as *SqlAttributeStore) commonQueryBuilder(option *model.AttributeFilterOption) (string, []interface{}, error) {
	query := as.GetQueryBuilder().
		Select(as.ModelFields(store.AttributeTableName + ".")...).
		From(store.AttributeTableName)

	// parse options
	if option.OrderBy != "" {
		query = query.OrderBy(option.OrderBy)
	} else {
		query = query.OrderBy(store.TableOrderingMap[store.AttributeTableName])
	}
	if option.Id != nil {
		query = query.Where(option.Id)
	}

	fieldValuesMap := map[string]*bool{
		"Attributes.VisibleInStoreFront":    option.VisibleInStoreFront,
		"Attributes.ValueRequired":          option.ValueRequired,
		"Attributes.IsVariantOnly":          option.IsVariantOnly,
		"Attributes.FilterableInStorefront": option.FilterableInStorefront,
		"Attributes.FilterableInDashboard":  option.FilterableInDashboard,
		"Attributes.AvailableInGrid":        option.AvailableInGrid,
	}

	for fieldName, value := range fieldValuesMap {
		if value != nil {
			query = query.Where(squirrel.Eq{fieldName: *value})
		}
	}

	if option.Type != nil {
		query = query.Where(option.Type)
	}
	if option.Distinct {
		query = query.Distinct()
	}
	if option.Slug != nil {
		query = query.Where(option.Slug)
	}
	if option.InputType != nil {
		query = query.Where(option.InputType)
	}
	if option.Extra != nil {
		query = query.Where(option.Extra)
	}
	if option.ProductTypes != nil {
		query = query.
			InnerJoin(store.AttributeProductTableName + " ON (AttributeProducts.AttributeID = Attributes.Id)").
			Where(option.ProductTypes)
	}
	if option.ProductVariantTypes != nil {
		query = query.
			InnerJoin(store.AttributeVariantTableName + " ON (AttributeVariants.AttributeID = Attributes.Id)").
			Where(option.ProductVariantTypes)
	}

	return query.ToSql()
}

func (as *SqlAttributeStore) GetByOption(option *model.AttributeFilterOption) (*model.Attribute, error) {
	queryString, args, err := as.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res model.Attribute
	err = as.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find attribute by option")
	}

	// check if we need to prefetch related attribute values of found attributes
	if option.PrefetchRelatedAttributeValues {
		attributeValues, err := as.AttributeValue().FilterByOptions(model.AttributeValueFilterOptions{
			AttributeID: squirrel.Eq{store.AttributeValueTableName + ".AttributeID": res.Id},
		})
		if err != nil {
			return nil, err
		}

		res.SetAttributeValues(attributeValues)
	}

	return &res, nil
}

// FilterbyOption returns a list of attributes by given option
func (as *SqlAttributeStore) FilterbyOption(option *model.AttributeFilterOption) (model.Attributes, error) {
	queryString, args, err := as.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var attributes model.Attributes
	err = as.GetReplicaX().Select(&attributes, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find attributes with given option")
	}

	// check if we need to prefetch related attribute values of found attributes
	if option.PrefetchRelatedAttributeValues && len(attributes) > 0 {
		attributeValues, err := as.AttributeValue().FilterByOptions(model.AttributeValueFilterOptions{
			AttributeID: squirrel.Eq{store.AttributeValueTableName + ".AttributeID": attributes.IDs()},
		})
		if err != nil {
			return nil, err
		}

		var attributeValueMap = map[string]model.AttributeValues{} // keys are attribute ids
		for _, value := range attributeValues {
			attributeValueMap[value.AttributeID] = append(attributeValueMap[value.AttributeID], value)
		}

		for _, attr := range attributes {
			values, ok := attributeValueMap[attr.Id]
			if ok {
				attr.SetAttributeValues(values)
			}
		}
	}

	return attributes, nil
}

func (as *SqlAttributeStore) Delete(ids ...string) (int64, error) {
	query, args, err := as.GetQueryBuilder().
		Delete("*").
		From(store.AttributeTableName).
		Where(squirrel.Eq{store.AttributeTableName + ".Id": ids}).
		ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "SqlAttributeStore.Delete_ToSql")
	}

	result, err := as.GetMasterX().Exec(query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete attributes")
	}

	numDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of deleted attributes")
	}

	return numDeleted, nil
}

func (s *SqlAttributeStore) GetProductTypeAttributes(productTypeID string, unassigned bool) (model.Attributes, error) {
	// unassigned query:
	query := `SELECT * FROM ` +
		store.AttributeTableName + `A WHERE A.Type = $1 AND NOT (
	(
		EXISTS(
			SELECT (1) AS "a"
			FROM ` + store.AttributeProductTableName + ` AP
			WHERE 
				AP.ProductTypeID = $2 AND AP.AttributeID = A.Id
			LIMIT 1
		)
		OR EXISTS(
			SELECT (1) AS "a"
			FROM ` + store.AttributeVariantTableName + ` AV
			WHERE
				AV.ProductTypeID = $3 AND AV.AttributeID = A.Id
			LIMIT 1
		)
	)
)`

	// in case select assigned attributes only:
	if !unassigned {
		query = `SELECT ` +
			s.ModelFields(store.AttributeTableName+".").Join(",") +
			` FROM ` + store.AttributeTableName +
			` A LEFT OUTER JOIN ` + store.AttributeProductTableName + ` AP ON A.Id = AP.AttributeID` +
			` LEFT OUTER JOIN ` + store.AttributeVariantTableName + ` AV ON AV.AttributeID = A.Id` +
			` WHERE (
					A.Type = $1 AND (
						AP.ProductTypeID = $2
						OR AV.ProductTypeID = $3
					)
			)`
	}

	var res model.Attributes
	err := s.GetReplicaX().Select(&res, query, model.PRODUCT_TYPE, productTypeID, productTypeID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product type attributes with given product type id")
	}

	return res, nil
}

func (s *SqlAttributeStore) GetPageTypeAttributes(pageTypeID string, unassigned bool) (model.Attributes, error) {
	// unassigned
	query := `SELECT * FROM ` + store.AttributeTableName +
		` A WHERE
			A.Type = $1
		AND
			NOT EXISTS(
				SELECT (1) AS "a"
				FROM ` + store.AttributePageTableName + ` AP
				WHERE (
					AP.PageTypeID = $2
					AND AP.AttributeID = A.Id
				)
				LIMIT 1
			)`

	if !unassigned {
		query = `SELECT ` + s.ModelFields("A.").Join(",") +
			` FROM ` + store.AttributeTableName +
			` A INNER JOIN ` + store.AttributePageTableName +
			` AP ON AP.AttributeID = A.Id
			WHERE A.Type = $1
			AND AP.PageTypeID = $2`
	}

	var res model.Attributes
	err := s.GetReplicaX().Select(&res, query, model.PAGE_TYPE, pageTypeID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find page type attribute with given page type id")
	}

	return res, nil
}
