package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeStore struct {
	store.Store
}

func NewSqlAttributeStore(s store.Store) store.AttributeStore {
	as := &SqlAttributeStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.Attribute{}, store.AttributeTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(attribute.ATTRIBUTE_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("Name").SetMaxSize(attribute.ATTRIBUTE_NAME_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(attribute.ATTRIBUTE_TYPE_MAX_LENGTH)
		table.ColMap("InputType").SetMaxSize(attribute.ATTRIBUTE_INPUT_TYPE_MAX_LENGTH)
		table.ColMap("EntityType").SetMaxSize(attribute.ATTRIBUTE_ENTITY_TYPE_MAX_LENGTH)
		table.ColMap("Unit").SetMaxSize(attribute.ATTRIBUTE_UNIT_MAX_LENGTH)

	}
	return as
}

func (as *SqlAttributeStore) ModelFields() []string {
	return []string{
		"Attributes.Id",
		"Attributes.Slug",
		"Attributes.Name",
		"Attributes.Type",
		"Attributes.InputType",
		"Attributes.EntityType",
		"Attributes.Unit",
		"Attributes.ValueRequired",
		"Attributes.IsVariantOnly",
		"Attributes.VisibleInStoreFront",
		"Attributes.FilterableInStorefront",
		"Attributes.FilterableInDashboard",
		"Attributes.StorefrontSearchPosition",
		"Attributes.AvailableInGrid",
		"Attributes.Metadata",
		"Attributes.PrivateMetadata",
	}
}

func (as *SqlAttributeStore) ScanFields(v attribute.Attribute) []interface{} {
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

func (as *SqlAttributeStore) CreateIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_attributes_name", store.AttributeTableName, "Name")
	as.CreateIndexIfNotExists("idx_attributes_name_lower_textpattern", store.AttributeTableName, "lower(Name) text_pattern_ops")
	as.CreateIndexIfNotExists("idx_attributes_slug", store.AttributeTableName, "Slug")
}

// Upsert inserts or updates given attribute then returns it
func (as *SqlAttributeStore) Upsert(attr *attribute.Attribute) (*attribute.Attribute, error) {
	var isSaving bool

	if !model.IsValidId(attr.Id) {
		attr.Id = ""
		isSaving = true
	}

	var (
		err        error
		numUpdated int64
	)

	if isSaving {
		attr.PreSave()
		if err := attr.IsValid(); err != nil {
			return nil, err
		}

		err = as.GetMaster().Insert(attr)
	} else {
		attr.PreUpdate()
		if err := attr.IsValid(); err != nil {
			return nil, err
		}

		numUpdated, err = as.GetMaster().Update(attr)
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

func (as *SqlAttributeStore) commonQueryBuilder(option *attribute.AttributeFilterOption) (string, []interface{}, error) {
	query := as.GetQueryBuilder().
		Select(as.ModelFields()...).
		From(store.AttributeTableName).
		OrderBy(store.TableOrderingMap[store.AttributeTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.VisibleInStoreFront != nil {
		expr := store.AttributeTableName + ".VisibleInStoreFront"
		if !*option.VisibleInStoreFront {
			expr = "NOT " + expr
		}
		query = query.Where(expr)
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

func (as *SqlAttributeStore) GetByOption(option *attribute.AttributeFilterOption) (*attribute.Attribute, error) {
	queryString, args, err := as.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res attribute.Attribute
	err = as.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find attribute by option")
	}

	// check if we need to prefetch related attribute values of found attributes
	if option.PrefetchRelatedAttributeValues {
		attributeValues, err := as.AttributeValue().FilterByOptions(attribute.AttributeValueFilterOptions{
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
func (as *SqlAttributeStore) FilterbyOption(option *attribute.AttributeFilterOption) (attribute.Attributes, error) {

	queryString, args, err := as.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var attributes attribute.Attributes
	_, err = as.GetReplica().Select(&attributes, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find attributes with given option")
	}

	// check if we need to prefetch related attribute values of found attributes
	if option.PrefetchRelatedAttributeValues && len(attributes) > 0 {
		attributeValues, err := as.AttributeValue().FilterByOptions(attribute.AttributeValueFilterOptions{
			AttributeID: squirrel.Eq{store.AttributeValueTableName + ".AttributeID": attributes.IDs()},
		})
		if err != nil {
			return nil, err
		}

		// attributesMap has keys are attribute ids
		var attributesMap = map[string]*attribute.Attribute{}

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
	result, err := as.GetMaster().Exec("DELETE FROM "+store.AttributeTableName+" WHERE Id IN :IDs", map[string]interface{}{"IDs": ids})
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete attributes")
	}

	numDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of deleted attributes")
	}

	return numDeleted, nil
}
