package menu

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlMenuStore struct {
	store.Store
}

func NewSqlMenuStore(sqlStore store.Store) store.MenuStore {
	return &SqlMenuStore{sqlStore}
}

func (s *SqlMenuStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Name",
		"Slug",
		"CreateAt",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (ms *SqlMenuStore) Save(mnu *model.Menu) (*model.Menu, error) {
	mnu.PreSave()
	if err := mnu.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.MenuTableName + "(" + ms.ModelFields("").Join(",") + ") VALUES (" + ms.ModelFields(":").Join(",") + ")"

	if _, err := ms.GetMasterX().NamedExec(query, mnu); err != nil {
		if ms.IsUniqueConstraintError(err, []string{"Name", "menus_name_key", "idx_menus_name_unique"}) {
			return nil, store.NewErrInvalidInput(model.MenuTableName, "Name", mnu.Name)
		}
		if ms.IsUniqueConstraintError(err, []string{"Slug", "menus_slug_key", "idx_menus_slug_unique"}) {
			return nil, store.NewErrInvalidInput(model.MenuTableName, "Slug", mnu.Slug)
		}
		return nil, errors.Wrapf(err, "failed to save menu with id=%s", mnu.Id)
	}

	return mnu, nil
}

func (ms *SqlMenuStore) commonQueryBuilder(options *model.MenuFilterOptions) squirrel.SelectBuilder {
	query := ms.GetQueryBuilder().
		Select(ms.ModelFields(model.MenuTableName + ".")...).
		From(model.MenuTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Name != nil {
		query = query.Where(options.Name)
	}
	if options.Slug != nil {
		query = query.Where(options.Slug)
	}

	return query
}

func (ms *SqlMenuStore) GetByOptions(options *model.MenuFilterOptions) (*model.Menu, error) {
	queryString, args, err := ms.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOptions_ToSql")
	}

	var res model.Menu
	err = ms.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.MenuTableName, "")
		}
		return nil, errors.Wrap(err, "failed to find menu with given options")
	}

	return &res, nil
}

func (ms *SqlMenuStore) FilterByOptions(options *model.MenuFilterOptions) ([]*model.Menu, error) {
	queryString, args, err := ms.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.Menu
	err = ms.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find menus with given options")
	}

	return res, nil
}
