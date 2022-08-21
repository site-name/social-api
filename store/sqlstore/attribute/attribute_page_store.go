package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributePageStore struct {
	store.Store
}

func NewSqlAttributePageStore(s store.Store) store.AttributePageStore {
	return &SqlAttributePageStore{s}
}

func (as *SqlAttributePageStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
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

func (as *SqlAttributePageStore) Save(page *attribute.AttributePage) (*attribute.AttributePage, error) {
	page.PreSave()
	if err := page.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.AttributePageTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"

	if _, err := as.GetMasterX().NamedExec(query, page); err != nil {
		if as.IsUniqueConstraintError(err, []string{"AttributeID", "PageTypeID", strings.ToLower(store.AttributePageTableName) + "_attributeid_pagetypeid_key"}) {
			return nil, store.NewErrInvalidInput(store.AttributePageTableName, "AttributeID/PageTypeID", page.AttributeID+"/"+page.PageTypeID)
		}
		return nil, errors.Wrapf(err, "failed to save attribute page with id=%s", page.Id)
	}

	return page, nil
}

func (as *SqlAttributePageStore) Get(pageID string) (*attribute.AttributePage, error) {
	var res attribute.AttributePage
	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AttributePageTableName+" WHERE Id = ?", pageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributePageTableName, pageID)
		}
		return nil, errors.Wrapf(err, "failed to find attribute page with id=%s", pageID)
	}

	return &res, nil
}

func (as *SqlAttributePageStore) GetByOption(option *attribute.AttributePageFilterOption) (*attribute.AttributePage, error) {
	query := as.GetQueryBuilder().Select("*").From(store.AttributePageTableName)

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

	var res attribute.AttributePage
	err = as.GetReplicaX().Get(
		&res,
		queryString,
		args...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributePageTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find attribute product with given option")
	}

	return &res, nil
}
