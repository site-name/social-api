package attribute

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type SqlAssignedPageAttributeStore struct {
	store.Store
}

func NewSqlAssignedPageAttributeStore(s store.Store) store.AssignedPageAttributeStore {
	return &SqlAssignedPageAttributeStore{
		Store: s,
	}
}

func (as *SqlAssignedPageAttributeStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{"Id", "PageID", "AssignmentID"}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (as *SqlAssignedPageAttributeStore) Save(pageAttr *attribute.AssignedPageAttribute) (*attribute.AssignedPageAttribute, error) {
	pageAttr.PreSave()
	if err := pageAttr.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.AssignedPageAttributeTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
	if _, err := as.GetMasterX().NamedExec(query, pageAttr); err != nil {
		if as.IsUniqueConstraintError(err, []string{"PageID", "AssignmentID", "assignedpageattributes_pageid_assignmentid_key"}) {
			return nil, store.NewErrInvalidInput(store.AssignedPageAttributeTableName, "PageID/AssignmentID", pageAttr.PageID+"/"+pageAttr.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned page attribute with id=%s", pageAttr.Id)
	}

	return pageAttr, nil
}

func (as *SqlAssignedPageAttributeStore) Get(id string) (*attribute.AssignedPageAttribute, error) {
	var res attribute.AssignedPageAttribute

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+store.AssignedPageAttributeTableName+" WHERE Id = ?", id)
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
		query = query.Where(option.AssignmentID)
	}
	if option.PageID != nil {
		query = query.Where(option.PageID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var res attribute.AssignedPageAttribute

	err = as.GetReplicaX().Get(
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
