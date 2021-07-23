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
	res, err := as.GetReplica().Get(attribute.AttributePage{}, pageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributePageTableName, pageID)
		}
		return nil, errors.Wrapf(err, "failed to find attribute page with id=%s", pageID)
	}

	return res.(*attribute.AttributePage), nil
}

func (as *SqlAttributePageStore) GetByOption(option *attribute.AttributePageFilterOption) (*attribute.AttributePage, error) {
	if option == nil || !model.IsValidId(option.AttributeID) || !model.IsValidId(option.PageTypeID) {
		return nil, store.NewErrInvalidInput(store.AttributePageTableName, "option", option)
	}

	var res *attribute.AttributePage
	err := as.GetReplica().SelectOne(
		&res,
		"SELECT * FROM "+store.AttributePageTableName+" WHERE (AttributeID = :AttributeID AND PageTypeID = :PageTypeID)",
		map[string]interface{}{
			"AttributeID": option.AttributeID,
			"PageTypeID":  option.PageTypeID,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AttributePageTableName, "")
		}
		return nil, errors.Wrapf(err, "failed to find attribute product with AttributeID = %s, PageTypeID = %s", option.AttributeID, option.PageTypeID)
	}

	return res, nil
}
