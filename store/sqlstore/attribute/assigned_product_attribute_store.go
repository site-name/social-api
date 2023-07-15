package attribute

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlAssignedProductAttributeStore struct {
	store.Store
}

func NewSqlAssignedProductAttributeStore(s store.Store) store.AssignedProductAttributeStore {
	return &SqlAssignedProductAttributeStore{s}
}

func (as *SqlAssignedProductAttributeStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{"Id", "ProductID", "AssignmentID"}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAssignedProductAttributeStore) Save(newInstance *model.AssignedProductAttribute) (*model.AssignedProductAttribute, error) {
	newInstance.PreSave()
	if err := newInstance.IsValid(); err != nil {
		return nil, err
	}

	if _, err := as.GetMasterX().Exec(
		"INSERT INTO "+model.AssignedPageAttributeTableName+" (Id, ProductID, AssignmentID) VALUES (?, ?, ?)",
		newInstance.Id, newInstance.ProductID, newInstance.AssignmentID,
	); err != nil {
		if as.IsUniqueConstraintError(err, []string{"ProductID", "AssignmentID", strings.ToLower(model.AssignedProductAttributeTableName) + "_productid_assignmentid_key"}) {
			return nil, store.NewErrInvalidInput(model.AssignedProductAttributeTableName, "ProductID/AssignmentID", newInstance.ProductID+"/"+newInstance.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to insert new assigned product attribute with id=%s", newInstance.Id)
	}

	return newInstance, nil
}

func (as *SqlAssignedProductAttributeStore) Get(id string) (*model.AssignedProductAttribute, error) {
	var res model.AssignedProductAttribute

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+model.AssignedProductAttributeTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.AssignedProductAttributeTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAssignedProductAttributeStore) commonQueryBuilder(options *model.AssignedProductAttributeFilterOption) squirrel.SelectBuilder {
	query := as.GetQueryBuilder().Select("*").From(model.AssignedProductAttributeTableName)

	// parse option
	if options.AssignmentID != nil {
		query = query.Where(options.AssignmentID)
	}
	if options.ProductID != nil {
		query = query.Where(options.ProductID)
	}
	if value := options.AttributeProduct_Attribute_VisibleInStoreFront; value != nil {
		query = query.
			InnerJoin(model.AttributeProductTableName + " ON AttributeProducts.Id = AssignedProductAttributes.AssignmentID").
			InnerJoin(model.AttributeTableName + " ON AttributeProducts.AttributeID = Attributes.Id").
			Where(squirrel.Eq{model.AttributeTableName + ".VisibleInStoreFront": *value})
	}

	return query
}

func (as *SqlAssignedProductAttributeStore) GetWithOption(option *model.AssignedProductAttributeFilterOption) (*model.AssignedProductAttribute, error) {
	queryString, args, err := as.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetWithOption_ToSql")
	}

	var res model.AssignedProductAttribute
	err = as.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.AssignedProductAttributeTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute with given options")
	}

	return &res, nil
}

func (as *SqlAssignedProductAttributeStore) FilterByOptions(options *model.AssignedProductAttributeFilterOption) ([]*model.AssignedProductAttribute, error) {
	queryString, args, err := as.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetWithOption_ToSql")
	}

	var res []*model.AssignedProductAttribute
	err = as.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find assigned product attributes with given options")
	}

	return res, nil
}
