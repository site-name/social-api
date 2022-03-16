package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
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

// Save inserts given attribute into database then returns it
func (as *SqlAttributeStore) Save(attr *attribute.Attribute) (*attribute.Attribute, error) {
	attr.PreSave()
	if err := attr.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(attr); err != nil {
		if as.IsUniqueConstraintError(err, []string{"Slug", "attributes_slug_key", "idx_attributes_slug_unique"}) {
			return nil, store.NewErrInvalidInput(store.AttributeTableName, "slug", attr.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save Attribute with attributeId=%s", attr.Id)
	}

	return attr, nil
}

func (as *SqlAttributeStore) Get(id string) (*attribute.Attribute, error) {
	var res attribute.Attribute
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AttributeTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Attribute", id)
		}
		return nil, errors.Wrapf(err, "failed to get Attribute with Id=%s", id)
	}

	return &res, nil
}

func (as *SqlAttributeStore) GetBySlug(slug string) (*attribute.Attribute, error) {
	var attr *attribute.Attribute
	err := as.GetReplica().SelectOne(&attr, "SELECT * FROM "+store.AttributeTableName+" WHERE Slug = :Slug", map[string]interface{}{"Slug": slug})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeTableName, "slug="+slug)
		}
		return nil, errors.Wrapf(err, "failed to find attribute with slug=%s", slug)
	}

	return attr, nil
}

// FilterbyOption returns a list of attributes by given option
func (as *SqlAttributeStore) FilterbyOption(option *attribute.AttributeFilterOption) (attribute.Attributes, error) {
	query := as.GetQueryBuilder().
		Select(as.ModelFields()...).
		From(store.AttributeTableName).
		OrderBy(store.TableOrderingMap[store.AttributeTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
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

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var res attribute.Attributes
	_, err = as.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find attributes with given option")
	}

	// check if we need to prefetch related attribute values of found attributes
	if option.PrefetchRelatedAttributeValues && len(res) > 0 {
		attributeValues, err := as.AttributeValue().FilterByOptions(attribute.AttributeValueFilterOptions{
			AttributeID: squirrel.Eq{store.AttributeValueTableName + ".AttributeID": res.IDs()},
		})
		if err != nil {
			return nil, err
		}

		var (
			// attributesMap has keys are attribute ids
			attributesMap = map[string]*attribute.Attribute{}
		)

		for _, attr := range res {
			attributesMap[attr.Id] = attr
		}

		for _, attrVl := range attributeValues {
			attributesMap[attrVl.AttributeID].AttributeValues = append(attributesMap[attrVl.AttributeID].AttributeValues, attrVl)
		}
	}

	return res, nil
}
