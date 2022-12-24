package attribute

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlAttributeVariantStore struct {
	store.Store
}

func NewSqlAttributeVariantStore(s store.Store) store.AttributeVariantStore {
	return &SqlAttributeVariantStore{s}
}

func (as *SqlAttributeVariantStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
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

func (as *SqlAttributeVariantStore) GetByOption(option *model.AttributeVariantFilterOption) (*model.AttributeVariant, error) {
	query := as.GetQueryBuilder().Select("*").From(store.AttributeVariantTableName)

	// parse option
	if option.AttributeID != nil {
		query = query.Where(option.AttributeID)
	}
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.ProductTypeID != nil {
		query = query.Where(option.ProductTypeID)
	}
	if len(option.ProductIDs) > 0 {
		subQuery := as.GetQueryBuilder().
			Select("AttributeID").
			From(store.ProductTableName).
			Where(squirrel.Eq{store.ProductTableName + ".Id": option.ProductIDs})
		query = query.Where(squirrel.Expr("AttributeID IN ?", subQuery))
	}

	queryString, args, err := query.ToSql()
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
