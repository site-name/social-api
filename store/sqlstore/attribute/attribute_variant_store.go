package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlAttributeVariantStore struct {
	store.Store
}

func NewSqlAttributeVariantStore(s store.Store) store.AttributeVariantStore {
	return &SqlAttributeVariantStore{s}
}

func (as *SqlAttributeVariantStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id", "AttributeID", "ProductTypeID", "VariantSelection", "SortOrder",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAttributeVariantStore) Save(attributeVariant *model.AttributeVariant) (*model.AttributeVariant, error) {
	attributeVariant.PreSave()
	if err := attributeVariant.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.AttributeVariantTableName + "(" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
	if _, err := as.GetMasterX().NamedExec(query, attributeVariant); err != nil {
		if as.IsUniqueConstraintError(err, []string{"AttributeID", "ProductTypeID", "attributevariants_attributeid_producttypeid_key"}) {
			return nil, store.NewErrInvalidInput(store.AttributeVariantTableName, "AttributeID/ProductTypeID", attributeVariant.AttributeID+"/"+attributeVariant.ProductTypeID)
		}
		return nil, errors.Wrapf(err, "failed to save attribute variant with id=%s", attributeVariant.Id)
	}

	return attributeVariant, nil
}

func (as *SqlAttributeVariantStore) Get(attributeVariantID string) (*model.AttributeVariant, error) {
	var res model.AttributeVariant

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AttributeVariantTableName+" WHERE Id = :ID", map[string]interface{}{"ID": attributeVariantID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeVariantTableName, attributeVariantID)
		}
		return nil, errors.Wrapf(err, "failed to find attribute variant with id=%s", attributeVariantID)
	}

	return &res, nil
}

func (s *SqlAttributeVariantStore) commonQueryBuilder(options *model.AttributeVariantFilterOption) squirrel.SelectBuilder {
	query := s.GetQueryBuilder().Select("*").From(store.AttributeVariantTableName)

	// parse option
	if options.AttributeID != nil {
		query = query.Where(options.AttributeID)
	}
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.ProductTypeID != nil {
		query = query.Where(options.ProductTypeID)
	}
	if value := options.AttributeVisibleInStoreFront; value != nil {
		query = query.
			InnerJoin(store.AttributeTableName + " ON Attributes.Id = AttributeVariants.AttributeID").
			Where(squirrel.Eq{store.AttributeTableName + ".VisibleInStoreFront": *value})
	}

	return query
}

func (as *SqlAttributeVariantStore) GetByOption(option *model.AttributeVariantFilterOption) (*model.AttributeVariant, error) {
	queryString, args, err := as.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}
	var res model.AttributeVariant

	err = as.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributeVariantTableName, "")
		}
		return nil, errors.Wrap(err, "failed to find attribute variant with given options")
	}

	return &res, nil
}

func (s *SqlAttributeVariantStore) FilterByOptions(options *model.AttributeVariantFilterOption) ([]*model.AttributeVariant, error) {
	queryString, args, err := s.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.AttributeVariant
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find attribute variant by given options")
	}
	return res, nil
}
