package attribute

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlAttributePageStore struct {
	store.Store
}

func NewSqlAttributePageStore(s store.Store) store.AttributePageStore {
	return &SqlAttributePageStore{s}
}

func (as *SqlAttributePageStore) Save(record model.AttributePage) (*model.AttributePage, error) {
	model_helper.AttributePagePreSave(&record)
	if err := model_helper.AttributePageIsValid(record); err != nil {
		return nil, err
	}

	err := record.Insert(as.GetMaster(), boil.Infer())
	if err != nil {
		if as.IsUniqueConstraintError(err, []string{"attribute_pages_attribute_id_page_type_id_key", model.AttributePageColumns.PageTypeID, model.AttributePageColumns.AttributeID}) {
			return nil, store.NewErrInvalidInput(model.TableNames.AttributePages, "page_type_id/attribute_id", "unique")
		}
		return nil, err
	}

	return &record, nil
}

func (as *SqlAttributePageStore) Get(id string) (*model.AttributePage, error) {
	record, err := model.FindAttributePage(as.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.AttributePages, id)
		}
		return nil, err
	}

	return record, nil
}

func (as *SqlAttributePageStore) GetByOption(option model_helper.AttributePageFilterOption) (*model.AttributePage, error) {
	record, err := model.AttributePages(option.Conditions...).One(as.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.AttributePages, "options")
		}
		return nil, err
	}

	return record, nil
}
