package attribute

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedPageAttributeStore struct {
	store.Store
}

func NewSqlAssignedPageAttributeStore(s store.Store) store.AssignedPageAttributeStore {
	as := &SqlAssignedPageAttributeStore{
		Store: s,
	}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(attribute.AssignedPageAttribute{}, store.AssignedPageAttributeTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("PageID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AssignmentID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("PageID", "AssignmentID")
	}
	return as
}

func (as *SqlAssignedPageAttributeStore) CreateIndexesIfNotExists() {
	as.CreateForeignKeyIfNotExists(store.AssignedPageAttributeTableName, "AssignmentID", store.AttributePageTableName, "Id", true)
	as.CreateForeignKeyIfNotExists(store.AssignedPageAttributeTableName, "PageID", "Pages", "Id", true)
}

func (as *SqlAssignedPageAttributeStore) Save(pageAttr *attribute.AssignedPageAttribute) (*attribute.AssignedPageAttribute, error) {
	pageAttr.PreSave()
	if err := pageAttr.IsValid(); err != nil {
		return nil, err
	}

	if err := as.GetMaster().Insert(pageAttr); err != nil {
		if as.IsUniqueConstraintError(err, []string{"PageID", "AssignmentID", strings.ToLower(store.AssignedPageAttributeTableName) + "_pageid_assignmentid_key"}) {
			return nil, store.NewErrInvalidInput(store.AssignedPageAttributeTableName, "PageID/AssignmentID", pageAttr.PageID+"/"+pageAttr.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned page attribute with id=%s", pageAttr.Id)
	}

	return pageAttr, nil
}

func (as *SqlAssignedPageAttributeStore) Get(id string) (*attribute.AssignedPageAttribute, error) {
	var res attribute.AssignedPageAttribute
	err := as.GetReplica().SelectOne(&res, "SELECT * FROM "+store.AssignedPageAttributeTableName+" WHERE Id = :Id", map[string]interface{}{"Id": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedPageAttributeTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAssignedPageAttributeStore) GetByOption(option *attribute.AssignedPageAttributeFilterOption) (*attribute.AssignedPageAttribute, error) {
	query := as.GetQueryBuilder().Select("*").From(store.AssignedPageAttributeTableName)

	// parse option
	if option.AssignmentID != nil {
		query = query.Where(option.AssignmentID.ToSquirrel("AssignmentID"))
	}
	if option.PageID != nil {
		query = query.Where(option.PageID.ToSquirrel("PageID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res attribute.AssignedPageAttribute
	err = as.GetReplica().SelectOne(
		&res,
		queryString,
		args...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.AssignedPageAttributeTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute with given option")
	}

	return &res, nil
}
