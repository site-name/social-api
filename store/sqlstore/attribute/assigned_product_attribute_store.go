package attribute

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedProductAttributeStore struct {
	store.Store
}

func NewSqlAssignedProductAttributeStore(s store.Store) store.AssignedProductAttributeStore {
	return &SqlAssignedProductAttributeStore{s}
}

func (as *SqlAssignedProductAttributeStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{"Id", "ProductID", "AssignmentID"}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAssignedProductAttributeStore) Save(newInstance *attribute.AssignedProductAttribute) (*attribute.AssignedProductAttribute, error) {
	newInstance.PreSave()
	if err := newInstance.IsValid(); err != nil {
		return nil, err
	}

	if _, err := as.GetMasterX().Exec(
		"INSERT INTO "+store.AssignedPageAttributeTableName+" (Id, ProductID, AssignmentID) VALUES (?, ?, ?)",
		newInstance.Id, newInstance.ProductID, newInstance.AssignmentID,
	); err != nil {
		if as.IsUniqueConstraintError(err, []string{"ProductID", "AssignmentID", strings.ToLower(store.AssignedProductAttributeTableName) + "_productid_assignmentid_key"}) {
			return nil, store.NewErrInvalidInput(store.AssignedProductAttributeTableName, "ProductID/AssignmentID", newInstance.ProductID+"/"+newInstance.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to insert new assigned product attribute with id=%s", newInstance.Id)
	}

	return newInstance, nil
}

func (as *SqlAssignedProductAttributeStore) Get(id string) (*attribute.AssignedProductAttribute, error) {
	var res attribute.AssignedProductAttribute

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AssignedProductAttributeTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedProductAttributeTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAssignedProductAttributeStore) commonQueryBuilder(options *attribute.AssignedProductAttributeFilterOption) squirrel.SelectBuilder {
	query := as.GetQueryBuilder().Select("*").From(store.AssignedProductAttributeTableName)

	// parse option
	if options.AssignmentID != nil {
		query = query.Where(options.AssignmentID)
	}
	if options.ProductID != nil {
		query = query.Where(options.ProductID)
	}

	return query
}

func (as *SqlAssignedProductAttributeStore) GetWithOption(option *attribute.AssignedProductAttributeFilterOption) (*attribute.AssignedProductAttribute, error) {
	queryString, args, err := as.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetWithOption_ToSql")
	}

	var res attribute.AssignedProductAttribute
	err = as.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedProductAttributeTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find assigned product attribute with given options")
	}

	return &res, nil
}

func (as *SqlAssignedProductAttributeStore) FilterByOptions(options *attribute.AssignedProductAttributeFilterOption) ([]*attribute.AssignedProductAttribute, error) {
	queryString, args, err := as.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetWithOption_ToSql")
	}

	var res []*attribute.AssignedProductAttribute
	err = as.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find assigned product attributes with given options")
	}

	return res, nil
}
