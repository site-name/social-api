package sqlstore

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributeStore struct {
	*SqlStore
}

func newSqlAttributeStore(s *SqlStore) store.AttributeStore {
	as := &SqlAttributeStore{
		SqlStore: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.Attribute{}, "Attributes").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(attribute.ATTRIBUTE_SLUG_MAX_LENGTH).SetUnique(true)
		table.ColMap("Name").SetMaxSize(attribute.ATTRIBUTE_NAME_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(attribute.ATTRIBUTE_TYPE_MAX_LENGTH)
		table.ColMap("InputType").SetMaxSize(attribute.ATTRIBUTE_INPUT_TYPE_MAX_LENGTH).SetDefaultConstraint(model.NewString(attribute.DROPDOWN))
		table.ColMap("EntityType").SetMaxSize(attribute.ATTRIBUTE_ENTITY_TYPE_MAX_LENGTH)
		table.ColMap("Unit").SetMaxSize(attribute.ATTRIBUTE_UNIT_MAX_LENGTH)

	}
	return as
}

func (as *SqlAttributeStore) createIndexesIfNotExists() {
	as.CreateIndexIfNotExists("idx_attributes_name", "Attributes", "Name")
	as.CreateIndexIfNotExists("idx_attributes_name_lower_textpattern", "Attributes", "lower(Name) text_pattern_ops")
	as.CreateIndexIfNotExists("idx_attributes_slug", "Attributes", "Slug")
}

func (as *SqlAttributeStore) Save(attr *attribute.Attribute) (*attribute.Attribute, error) {
	attr.PreSave()
	if err := attr.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(attr); err != nil {
		if IsUniqueConstraintError(err, []string{"Slug", "attributes_slug_key", "idx_attributes_slug_unique"}) {
			return nil, store.NewErrInvalidInput("Attribute", "slug", attr.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save Attribute with attributeId=%s", attr.Id)
	}

	return attr, nil
}

func (as *SqlAttributeStore) Get(id string) (*attribute.Attribute, error) {
	fetchedRow, err := as.GetMaster().Get(attribute.Attribute{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Attribute", id)
		}
		return nil, errors.Wrapf(err, "failed to get Attribute with Id=%s", id)
	}

	return fetchedRow.(*attribute.Attribute), nil
}

func (as *SqlAttributeStore) GetAttributesByIds(ids []string) ([]*attribute.Attribute, error) {
	query, args, err := as.getQueryBuilder().Select("*").From("Attributes").Where(squirrel.Eq{"Id": ids}).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_attributes_by_ids")
	}
	var attrs []*attribute.Attribute
	if _, err := as.GetMaster().Select(&attrs, query, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find Attributes")
	}

	return attrs, nil
}

// GetProductAndVariantHeaders is used for get headers for csv export preparation.
func (as *SqlAttributeStore) GetProductAndVariantHeaders(ids []string) ([]string, error) {
	tx, err := as.GetReplica().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "begin_transaction")
	}
	defer finalizeTransaction(tx)

	var productHeaders []string
	_, err = tx.Select(
		&productHeaders,
		`SELECT DISTINCT
		CONCAT(a.Slug, ' (product attribute)')
		AS
			header
		FROM
			Attributes AS a
		INNER JOIN
			AttributeProducts AS ap
		ON
			(ap.AttributeID = a.Id)
		WHERE
			a.Id IN :IDS
		AND 
			ap.ProductTypeID IS NOT NULL
		ORDER BY
			a.Slug`,
		map[string]interface{}{"IDS": ids},
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to find product attributes")
	}

	var variantHeaders []string
	_, err = tx.Select(
		&variantHeaders,
		`SELECT DISTINCT
		CONCAT(a.Slug, ' (variant attribute)') 
		AS 
			header 
		FROM 
			Attributes AS a 
		INNER JOIN 
			AttributeVariants AS av
		ON 
			(av.AttributeID = a.Id)
		WHERE
			a.Id IN :IDS
		AND 
			av.ProductTypeID IS NOT NULL
		ORDER BY
			a.Slug`,
		map[string]interface{}{"IDS": ids},
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to find variant attributes")
	}

	return append(productHeaders, variantHeaders...), nil
}
