package product

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlDigitalContentStore struct {
	store.Store
}

func NewSqlDigitalContentStore(s store.Store) store.DigitalContentStore {
	return &SqlDigitalContentStore{s}
}

func (ds *SqlDigitalContentStore) Save(content model.DigitalContent) (*model.DigitalContent, error) {
	model_helper.DigitalContentPreSave(&content)

	if err := model_helper.DigitalContentIsValid(content); err != nil {
		return nil, err
	}

	if err := content.Insert(ds.GetMaster(), boil.Infer()); err != nil {
		return nil, err
	}

	return &content, nil
}

func (ds *SqlDigitalContentStore) GetByOption(option model_helper.DigitalContentFilterOption) (*model.DigitalContent, error) {
	content, err := model.DigitalContents(option.Conditions...).One(ds.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.DigitalContents, "options")
		}
		return nil, err
	}

	return content, nil
}

func (ds *SqlDigitalContentStore) FilterByOption(option model_helper.DigitalContentFilterOption) (model.DigitalContentSlice, error) {
	return model.DigitalContents(option.Conditions...).All(ds.GetReplica())
}

func (s *SqlDigitalContentStore) Delete(tx boil.ContextTransactor, ids []string) error {
	if tx == nil {
		tx = s.GetMaster()
	}

	_, err := model.DigitalContents(model.DigitalContentWhere.ID.IN(ids)).DeleteAll(tx)
	return err
}
