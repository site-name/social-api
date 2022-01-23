package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAttributePageStore struct {
	store.Store
}

func NewSqlAttributePageStore(s store.Store) store.AttributePageStore {
	as := &SqlAttributePageStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AttributePage{}, store.AttributePageTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AttributeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("PageTypeID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("AttributeID", "PageTypeID")
	}
	return as
}

func (as *SqlAttributePageStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AttributePageTableName, "AttributeID", store.AttributeTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AttributePageTableName, "PageTypeID", "PageTypes", "Id", true)
}

func (as *SqlAttributePageStore) Save(page *attribute.AttributePage) (*attribute.AttributePage, error) {
	page.PreSave()
	if err := page.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(page); err != nil {
		if as.IsUniqueConstraintError(err, []string{"AttributeID", "PageTypeID", strings.ToLower(store.AttributePageTableName) + "_attributeid_pagetypeid_key"}) {
			return nil, store.NewErrInvalidInput(store.AttributePageTableName, "AttributeID/PageTypeID", page.AttributeID+"/"+page.PageTypeID)
		}
		return nil, errors.Wrapf(err, "failed to save attribute page with id=%s", page.Id)
	}

	return page, nil
}

func (as *SqlAttributePageStore) Get(pageID string) (*attribute.AttributePage, error) {
	var res attribute.AttributePage
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AttributePageTableName+" WHERE Id = :ID", map[string]interface{}{"ID": pageID})
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
	err = as.GetReplica().SelectOne(
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
