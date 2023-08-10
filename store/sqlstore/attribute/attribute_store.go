package attribute

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

type SqlAttributeStore struct {
	store.Store
}

func NewSqlAttributeStore(s store.Store) store.AttributeStore {
	return &SqlAttributeStore{s}
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
	err := as.GetMaster().Save(attr).Error

	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"Slug", "attributes_slug_key", "idx_attributes_slug_unique", "slug_unique_key"}) {
			return nil, store.NewErrInvalidInput(model.AttributeTableName, "Slug", attr.Slug)
		}
		return nil, errors.Wrap(err, "failed to upsert attribute")
	}

	return attr, nil
}

func (as *SqlAttributeStore) commonQueryBuilder(option *model.AttributeFilterOption) squirrel.SelectBuilder {
	query := as.GetQueryBuilder().
		Select(model.AttributeTableName + ".*").
		From(model.AttributeTableName).
		Where(option.Conditions)

	// parse options
	if option.OrderBy != "" {
		query = query.OrderBy(option.OrderBy)
	}
	if option.Distinct {
		query = query.Distinct()
	}
	if option.Limit > 0 {
		query = query.Limit(uint64(option.Limit))
	}

	if option.AttributeProduct_ProductTypeID != nil {
		query = query.
			InnerJoin(model.AttributeProductTableName + " ON (AttributeProducts.AttributeID = Attributes.Id)").
			Where(option.AttributeProduct_ProductTypeID)
	}
	if option.AttributeVariant_ProductTypeID != nil {
		query = query.
			InnerJoin(model.AttributeVariantTableName + " ON (AttributeVariants.AttributeID = Attributes.Id)").
			Where(option.AttributeVariant_ProductTypeID)
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
		if option.ChannelSlug != nil {
			channelIdOrSlug = *option.ChannelSlug
		}

		productQuery := as.
			Product().
			VisibleToUserProductsQuery(channelIdOrSlug, option.UserIsShopStaff)

		if option.InCategory != nil {
			productQuery = productQuery.Where("Products.CategoryID = ?", *option.InCategory)

			if !option.UserIsShopStaff {
				productQuery = productQuery.Column(`(
					SELECT ProductChannelListings.VisibleInListings
					FROM ProductChannelListings
					INNER JOIN Channels ON Channels.Id = ProductChannelListings.ChannelID
					WHERE (
						(Channels.Id = ? OR Channels.Slug = ?)
						AND ProductChannelListings.ProductID = Products.Id
					)
					LIMIT 1
				) AS VisibleInListings`, channelIdOrSlug, channelIdOrSlug).
					Where("VisibleInListings")
			}

		} else if option.InCollection != nil {
			productQuery = productQuery.
				InnerJoin(model.CollectionProductRelationTableName+" ON ProductCollections.ProductID = Products.Id").
				Where("ProductCollections.CollectionID = ?", *option.InCollection)
		}
		//
		products, err := as.Product().FilterByQuery(productQuery)
		if err != nil {
			slog.Error("failed to find products for filtering attributes", slog.Err(err))
			return query
		}

		productTypeIDs := products.ProductTypeIDs()
		query = query.
			LeftJoin(model.AttributeVariantTableName + " ON AttributeVariants.AttributeID = Attributes.Id").
			LeftJoin(model.AttributeProductTableName + " ON AttributeProducts.AttributeID = Attributes.Id").
			Where(squirrel.Or{
				squirrel.Eq{model.AttributeVariantTableName + ".ProductTypeID": productTypeIDs},
				squirrel.Eq{model.AttributeProductTableName + ".ProductTypeID": productTypeIDs},
			})
	}

	return query
}

// FilterbyOption returns a list of attributes by given option
func (as *SqlAttributeStore) FilterbyOption(option *model.AttributeFilterOption) (model.Attributes, error) {
	queryString, args, err := as.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var attributes model.Attributes
	err = as.GetReplica().Raw(queryString, args...).Scan(&attributes).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find attributes with given option")
	}

	// check if we need to prefetch related attribute values of found attributes
	if option.PrefetchRelatedAttributeValues && len(attributes) > 0 {
		attributeValues, err := as.AttributeValue().FilterByOptions(model.AttributeValueFilterOptions{
			Conditions: squirrel.Eq{model.AttributeValueTableName + ".AttributeID": attributes.IDs()},
		})
		if err != nil {
			return nil, err
		}

		var attributeValueMap = map[string]model.AttributeValues{} // keys are attribute ids
		for _, value := range attributeValues {
			attributeValueMap[value.AttributeID] = append(attributeValueMap[value.AttributeID], value)
		}

		for _, attr := range attributes {
			attr.AttributeValues = attributeValueMap[attr.Id]
		}
	}

	return attributes, nil
}

func (as *SqlAttributeStore) Delete(ids ...string) (int64, error) {
	result := as.GetMaster().Raw("DELETE FROM "+model.AttributeTableName+" WHERE Id IN ?", ids)
	if result.Error != nil {
		return 0, errors.Wrap(result.Error, "failed to delete attributes")
	}

	return result.RowsAffected, nil
}

func (s *SqlAttributeStore) GetProductTypeAttributes(productTypeID string, unassigned bool, filter *model.AttributeFilterOption) (model.Attributes, error) {
	if filter == nil {
		filter = new(model.AttributeFilterOption)
	}
	filter.Conditions = squirrel.Eq{model.AttributeTableName + ".Type": model.PRODUCT_TYPE}
	filter.Distinct = true
	sqQuery := s.commonQueryBuilder(filter)

	if unassigned {
		sqQuery = sqQuery.Where(`NOT (
	EXISTS(
		SELECT (1) AS "a"
		FROM `+model.AttributeProductTableName+` WHERE
			AttributeProducts.ProductTypeID = ? AND AttributeProducts.AttributeID = Attributes.Id
		LIMIT 1
	)
	OR EXISTS(
		SELECT (1) AS "a"
		FROM `+model.AttributeVariantTableName+` WHERE
			AttributeVariants.ProductTypeID = ? AND AttributeVariants.AttributeID = Attributes.Id
		LIMIT 1
	)
)`, productTypeID, productTypeID)

	} else {
		sqQuery = sqQuery.
			LeftJoin(model.AttributeProductTableName+" ON AttributeProducts.AttributeID = Attributes.Id").
			LeftJoin(model.AttributeVariantTableName+" ON Attributes.Id = AttributeVariants.AttributeID").
			Where("Attributes.Type = ?", model.PRODUCT_TYPE).
			Where("AttributeProducts.ProductTypeID = ? OR AttributeVariants.ProductTypeID = ?", productTypeID)
	}

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetProductTypeAttributes_ToSql")
	}

	var res model.Attributes
	err = s.GetReplica().Raw(query, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product type attributes with given product type id")
	}

	return res, nil
}

func (s *SqlAttributeStore) GetPageTypeAttributes(pageTypeID string, unassigned bool) (model.Attributes, error) {
	// unassigned
	query := `SELECT * FROM ` +
		model.AttributeTableName +
		` A WHERE A.Type = $1
		AND NOT EXISTS(
			SELECT (1) AS "a"
			FROM ` + model.AttributePageTableName +
		` AP WHERE (
				AP.PageTypeID = $2
				AND AP.AttributeID = A.Id
			)
			LIMIT 1
		)`

	if !unassigned {
		query = `SELECT A.* FROM ` + model.AttributeTableName +
			` A INNER JOIN ` + model.AttributePageTableName +
			` AP ON AP.AttributeID = A.Id
			WHERE A.Type = $1
			AND AP.PageTypeID = $2`
	}

	var res model.Attributes
	err := s.GetReplica().Raw(query, model.PAGE_TYPE, pageTypeID).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find page type attribute with given page type id")
	}

	return res, nil
}
