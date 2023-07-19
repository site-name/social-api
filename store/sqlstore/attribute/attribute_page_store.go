package attribute

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAttributePageStore struct {
	store.Store
}

func NewSqlAttributePageStore(s store.Store) store.AttributePageStore {
	return &SqlAttributePageStore{s}
}

func (as *SqlAttributePageStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"AttributeID",
		"PageTypeID",
		"SortOrder",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (as *SqlAttributePageStore) Save(page *model.AttributePage) (*model.AttributePage, error) {
	if err := as.GetMaster().Create(page).Error; err != nil {
		if as.IsUniqueConstraintError(err, []string{"AttributeID", "PageTypeID", "attributeid_pagetypeid_key"}) {
			return nil, store.NewErrInvalidInput(model.AttributePageTableName, "AttributeID/PageTypeID", page.AttributeID+"/"+page.PageTypeID)
		}
		return nil, errors.Wrapf(err, "failed to save attribute page with id=%s", page.Id)
	}

	return page, nil
}

func (as *SqlAttributePageStore) Get(id string) (*model.AttributePage, error) {
	var res model.AttributePage
	err := as.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AttributePageTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find attribute page with id=%s", id)
	}

	return &res, nil
}

func (as *SqlAttributePageStore) GetByOption(option *model.AttributePageFilterOption) (*model.AttributePage, error) {
	var res model.AttributePage
	err := as.GetReplica().First(&res, store.BuildSqlizer(option.Conditions)...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AttributePageTableName, "option")
		}
		return nil, errors.Wrapf(err, "failed to find attribute product with given option")
	}

	return &res, nil
}
