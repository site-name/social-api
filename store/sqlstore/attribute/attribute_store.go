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
		Select(as.ModelFields(store.AttributeTableName + ".")...). // SELECT Attributes.Id, Attributes.Slug, ...
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
	if option.VisibleInStoreFront != nil {
		cond := store.AttributeTableName + ".VisibleInStoreFront"
		if !*option.VisibleInStoreFront {
			cond = "NOT " + cond
		}
		query = query.Where(cond)
	}
	if option.ValueRequired != nil {
		cond := store.AttributeTableName + ".ValueRequired"
		if !*option.ValueRequired {
			cond = "NOT " + cond
		}
		query = query.Where(cond)
	}
	if option.IsVariantOnly != nil {
		cond := store.AttributeTableName + ".IsVariantOnly"
		if !*option.IsVariantOnly {
			cond = "NOT " + cond
		}
		query = query.Where(cond)
	}
	if option.FilterableInStorefront != nil {
		cond := store.AttributeTableName + ".FilterableInStorefront"
		if !*option.FilterableInStorefront {
			cond = "NOT " + cond
		}
		query = query.Where(cond)
	}
	if option.FilterableInDashboard != nil {
		cond := store.AttributeTableName + ".FilterableInDashboard"
		if !*option.FilterableInDashboard {
			cond = "NOT " + cond
		}
		query = query.Where(cond)
	}
	if option.AvailableInGrid != nil {
		cond := store.AttributeTableName + ".AvailableInGrid"
		if !*option.AvailableInGrid {
			cond = "NOT " + cond
		}
		query = query.Where(cond)
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

		res.AttributeValues = attributeValues
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

		// attributesMap has keys are attribute ids
		var attributesMap = map[string]*model.Attribute{}

		for _, attr := range attributes {
			attributesMap[attr.Id] = attr
		}

		for _, attrVl := range attributeValues {
			attributesMap[attrVl.AttributeID].AttributeValues = append(attributesMap[attrVl.AttributeID].AttributeValues, attrVl)
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
