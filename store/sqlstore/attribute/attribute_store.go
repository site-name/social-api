package attribute

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlAttributeStore struct {
	store.Store
}

func NewSqlAttributeStore(s store.Store) store.AttributeStore {
	return &SqlAttributeStore{s}
}

var attributeFieldNames = util.AnyArray[string]{
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

func (as *SqlAttributeStore) ModelFields(prefix string) util.AnyArray[string] {
	if prefix == "" {
		return attributeFieldNames
	}
	return attributeFieldNames.Map(func(_ int, s string) string {
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
			attr.Slug = attr.Slug + model.NewRandomString(5)
			return as.Upsert(attr)
		}
		return nil, errors.Wrap(err, "failed to upsert attribute")
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("%d attribute(s) was/were updated instead of 1", numUpdated)
	}

	return attr, nil
}

func (as *SqlAttributeStore) commonQueryBuilder(option *model.AttributeFilterOption) squirrel.SelectBuilder {
	query := as.GetQueryBuilder().
		Select(as.ModelFields(store.AttributeTableName + ".")...).
		From(store.AttributeTableName)

	// parse options
	if option.OrderBy != "" {
		query = query.OrderBy(option.OrderBy)
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
	if option.Metadata != nil && len(option.Metadata) > 0 {
		delete(option.Metadata, "")
		conditions := []string{}

		for key, value := range option.Metadata {
			if value != "" {
				expr := fmt.Sprintf("Attributes.Metadata::jsonb @> '{%q:%q}'", key, value)
				conditions = append(conditions, expr)
				continue
			}
			expr := fmt.Sprintf("Attributes.Metadata::jsonb ? '%s'", key)
			conditions = append(conditions, expr)
		}
		query = query.Where(strings.Join(conditions, " AND "))
	}
	if option.Search != "" {
		expr := "%" + option.Search + "%"
		query = query.Where("Attributes.Name ILIKE ? OR Attributes.Slug ILIKE ?", expr, expr)
	}

	if option.InCategory != nil || option.InCollection != nil {

		var channelIdOrSlug string
		if option.Channel != nil {
			channelIdOrSlug = *option.Channel
		}

		productQuery := as.
			Product().
			VisibleToUserProductsQuery(channelIdOrSlug, option.UserIsShopStaff)

		if option.InCategory != nil {
			productQuery = productQuery.Where("Products.CategoryID = ?", *option.InCategory)

			if !option.UserIsShopStaff {
				productQuery = productQuery.Column(`(
					SELECT PC.VisibleInListings
					FROM ProductChannelListings PC
					INNER JOIN Channels C ON C.Id = PC.ChannelID
					WHERE (
						(C.Id = ? OR C.Slug = ?)
						AND PC.ProductID = Products.Id
					)
					LIMIT 1
				) AS VisibleInListings`, channelIdOrSlug, channelIdOrSlug).
					Where("VisibleInListings")
			}

		} else if option.InCollection != nil {
			productQuery = productQuery.
				InnerJoin(store.CollectionProductRelationTableName+" PC ON PC.ProductID = Products.Id").
				Where("PC.CollectionID = ?", *option.InCollection)
		}
		//
		products, err := as.Product().FilterByQuery(productQuery)
		if err != nil {
			slog.Error("failed to find products for filtering attributes", slog.Err(err))
			return query
		}

		productTypeIDs := products.ProductTypeIDs()
		query = query.
			LeftJoin(store.AttributeVariantTableName + " AV ON AV.AttributeID = Attributes.Id").
			LeftJoin(store.AttributeProductTableName + " AP ON AP.AttributeID = Attributes.Id").
			Where(squirrel.Or{
				squirrel.Eq{"AV.ProductTypeID": productTypeIDs},
				squirrel.Eq{"AP.ProductTypeID": productTypeIDs},
			})
	}

	return query
}

func (as *SqlAttributeStore) GetByOption(option *model.AttributeFilterOption) (*model.Attribute, error) {
	queryString, args, err := as.commonQueryBuilder(option).ToSql()
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
	queryString, args, err := as.commonQueryBuilder(option).ToSql()
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
	args := lo.Map(ids, func(id string, _ int) any { return id })
	queryStr := "DELETE FROM " + store.AttributeTableName + " WHERE Id IN (" + squirrel.Placeholders(len(ids)) + ")"

	result, err := as.GetMasterX().Exec(queryStr, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete attributes")
	}

	numDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of deleted attributes")
	}

	return numDeleted, nil
}

func (s *SqlAttributeStore) GetProductTypeAttributes(productTypeID string, unassigned bool, filter *model.AttributeFilterOption) (model.Attributes, error) {
	if filter == nil {
		filter = new(model.AttributeFilterOption)
	}
	filter.Type = squirrel.Eq{store.AttributeTableName + ".Type": model.PRODUCT_TYPE}
	filter.Distinct = true
	sqQuery := s.commonQueryBuilder(filter)

	if unassigned {
		sqQuery = sqQuery.Where(`NOT (
	EXISTS(
		SELECT (1) AS "a"
		FROM `+store.AttributeProductTableName+` WHERE
			AttributeProducts.ProductTypeID = ? AND AttributeProducts.AttributeID = Attributes.Id
		LIMIT 1
	)
	OR EXISTS(
		SELECT (1) AS "a"
		FROM `+store.AttributeVariantTableName+` WHERE
			AttributeVariants.ProductTypeID = ? AND AttributeVariants.AttributeID = Attributes.Id
		LIMIT 1
	)
)`, productTypeID, productTypeID)

	} else {
		sqQuery = sqQuery.
			LeftJoin(store.AttributeProductTableName+" ON AttributeProducts.AttributeID = Attributes.Id").
			LeftJoin(store.AttributeVariantTableName+" ON Attributes.Id = AttributeVariants.AttributeID").
			Where("Attributes.Type = ?", model.PRODUCT_TYPE).
			Where("AttributeProducts.ProductTypeID = ? OR AttributeVariants.ProductTypeID = ?", productTypeID)
	}

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetProductTypeAttributes_ToSql")
	}

	var res model.Attributes
	err = s.GetReplicaX().Select(&res, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product type attributes with given product type id")
	}

	return res, nil
}

func (s *SqlAttributeStore) GetPageTypeAttributes(pageTypeID string, unassigned bool) (model.Attributes, error) {
	// unassigned
	query := `SELECT * FROM ` +
		store.AttributeTableName +
		` A WHERE A.Type = $1
		AND NOT EXISTS(
			SELECT (1) AS "a"
			FROM ` + store.AttributePageTableName +
		` AP WHERE (
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
