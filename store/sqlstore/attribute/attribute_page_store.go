package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlAttributePageStore struct {
	store.Store
}

func NewSqlAttributePageStore(s store.Store) store.AttributePageStore {
	return &SqlAttributePageStore{s}
}

func (as *SqlAttributePageStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"AttributeID",
		"PageTypeID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAttributePageStore) Save(page *model.AttributePage) (*model.AttributePage, error) {
	page.PreSave()
	if err := page.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.AttributePageTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"

	if _, err := as.GetMasterX().NamedExec(query, page); err != nil {
		if as.IsUniqueConstraintError(err, []string{"AttributeID", "PageTypeID", strings.ToLower(model.AttributePageTableName) + "_attributeid_pagetypeid_key"}) {
			return nil, store.NewErrInvalidInput(model.AttributePageTableName, "AttributeID/PageTypeID", page.AttributeID+"/"+page.PageTypeID)
		}
		return nil, errors.Wrapf(err, "failed to save attribute page with id=%s", page.Id)
	}

	return page, nil
}

func (as *SqlAttributePageStore) Get(pageID string) (*model.AttributePage, error) {
	var res model.AttributePage
	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+model.AttributePageTableName+" WHERE Id = ?", pageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.AttributePageTableName, pageID)
		}
		return nil, errors.Wrapf(err, "failed to find attribute page with id=%s", pageID)
	}

	return &res, nil
}

func (as *SqlAttributePageStore) GetByOption(option *model.AttributePageFilterOption) (*model.AttributePage, error) {
	query := as.GetQueryBuilder().Select("*").From(model.AttributePageTableName)

	// parse option
	if option.PageTypeID != nil {
		query = query.Where(option.PageTypeID)
	}
	if option.AttributeID != nil {
		query = query.Where(option.AttributeID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByoption_ToSql")
	}

	var res model.AttributePage
	err = as.GetReplicaX().Get(
		&res,
		queryString,
		args...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.AttributePageTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find attribute product with given option")
	}

	return &res, nil
}
