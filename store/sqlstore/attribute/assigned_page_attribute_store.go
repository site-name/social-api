package attribute

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlAssignedPageAttributeStore struct {
	store.Store
}

func NewSqlAssignedPageAttributeStore(s store.Store) store.AssignedPageAttributeStore {
	return &SqlAssignedPageAttributeStore{
		Store: s,
	}
}

func (as *SqlAssignedPageAttributeStore) Upsert(pageAttr model.AssignedPageAttribute) (*model.AssignedPageAttribute, error) {
	if err := model_helper.AssignedPageAttributeIsValid(pageAttr); err != nil {
		return nil, err
	}

	var err error
	if pageAttr.ID == "" {
		err = pageAttr.Insert(as.GetMaster(), boil.Infer())
	} else {
		_, err = pageAttr.Update(as.GetMaster(), boil.Infer())
	}

	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"assigned_page_attributes_page_id_assignment_id_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.AssignedPageAttributes, "PageID/AssignmentID", pageAttr.PageID+"/"+pageAttr.AssignmentID)
		}
		return nil, err
	}

	return &pageAttr, nil
}

func (as *SqlAssignedPageAttributeStore) Get(id string) (*model.AssignedPageAttribute, error) {
	attr, err := model.FindAssignedPageAttribute(as.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.AssignedPageAttributes, id)
		}
		return nil, err
	}

	return attr, nil
}

func (as *SqlAssignedPageAttributeStore) FilterByOptions(options model_helper.AssignedPageAttributeFilterOption) (model.AssignedPageAttributeSlice, error) {
	return model.AssignedPageAttributes(options.Conditions...).All(as.GetReplica())
}
