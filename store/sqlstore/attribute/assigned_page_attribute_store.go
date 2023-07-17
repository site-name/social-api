package attribute

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAssignedPageAttributeStore struct {
	store.Store
}

func NewSqlAssignedPageAttributeStore(s store.Store) store.AssignedPageAttributeStore {
	return &SqlAssignedPageAttributeStore{
		Store: s,
	}
}

func (as *SqlAssignedPageAttributeStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{"Id", "PageID", "AssignmentID"}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (as *SqlAssignedPageAttributeStore) Save(pageAttr *model.AssignedPageAttribute) (*model.AssignedPageAttribute, error) {
	pageAttr.PreSave()
	if err := pageAttr.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.AssignedPageAttributeTableName + " (" + as.ModelFields("").Join(",") + ") VALUES (" + as.ModelFields(":").Join(",") + ")"
	if _, err := as.GetMasterX().NamedExec(query, pageAttr); err != nil {
		if as.IsUniqueConstraintError(err, []string{"PageID", "AssignmentID", "assignedpageattributes_pageid_assignmentid_key"}) {
			return nil, store.NewErrInvalidInput(model.AssignedPageAttributeTableName, "PageID/AssignmentID", pageAttr.PageID+"/"+pageAttr.AssignmentID)
		}
		return nil, errors.Wrapf(err, "failed to save assigned page attribute with id=%s", pageAttr.Id)
	}

	return pageAttr, nil
}

func (as *SqlAssignedPageAttributeStore) Get(id string) (*model.AssignedPageAttribute, error) {
	var res model.AssignedPageAttribute

	err := as.GetReplicaX().Get(&res, "SELECT * FROM "+model.AssignedPageAttributeTableName+" WHERE Id = ?", id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedPageAttributeTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAssignedPageAttributeStore) GetByOption(option *model.AssignedPageAttributeFilterOption) (*model.AssignedPageAttribute, error) {
	query := as.GetQueryBuilder().Select("*").From(model.AssignedPageAttributeTableName)

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

	var res model.AssignedPageAttribute

	err = as.GetReplicaX().Get(
		&res,
		queryString,
		args...,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedPageAttributeTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute with given option")
	}

	return &res, nil
}
