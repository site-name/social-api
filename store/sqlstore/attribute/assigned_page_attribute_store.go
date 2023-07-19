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
	err := as.GetMaster().Create(pageAttr).Error
	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"PageID", "AssignmentID", "assignedpageattributes_pageid_assignmentid_key"}) {
			return nil, store.NewErrInvalidInput(model.AssignedPageAttributeTableName, "PageID/AssignmentID", pageAttr.PageID+"/"+pageAttr.AssignmentID)
		}
		return nil, errors.Wrap(err, "failed to save assigned page attribute with")
	}

	return pageAttr, nil
}

func (as *SqlAssignedPageAttributeStore) Get(id string) (*model.AssignedPageAttribute, error) {
	var res model.AssignedPageAttribute
	err := as.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedPageAttributeTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAssignedPageAttributeStore) GetByOption(option *model.AssignedPageAttributeFilterOption) (*model.AssignedPageAttribute, error) {
	var res model.AssignedPageAttribute
	err := as.GetReplica().First(&res, store.BuildSqlizer(option.Conditions)...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AssignedPageAttributeTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find assigned page attribute with given option")
	}

	return &res, nil
}
