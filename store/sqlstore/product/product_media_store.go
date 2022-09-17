package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlProductMediaStore struct {
	store.Store
}

func NewSqlProductMediaStore(s store.Store) store.ProductMediaStore {
	return &SqlProductMediaStore{s}
}

func (s *SqlProductMediaStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"CreateAt",
		"ProductID",
		"Ppoi",
		"Image",
		"Alt",
		"Type",
		"ExternalUrl",
		"OembedData",
		"SortOrder",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert depends on given media's Id property to decide insert or update it
func (ps *SqlProductMediaStore) Upsert(media *model.ProductMedia) (*model.ProductMedia, error) {
	var isSaving bool

	if media.Id == "" {
		media.PreSave()
		isSaving = true
	} else {
		media.PreUpdate()
	}

	if err := media.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + store.ProductMediaTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
		_, err = ps.GetMasterX().NamedExec(query, media)

	} else {
		query := "UPDATE " + store.ProductMediaTableName + " SET " + ps.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = ps.GetMasterX().NamedExec(query, media)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert product media with id=%s", media.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple product medias were updated: %d instead of 1", numUpdated)
	}

	return media, nil
}

// Get finds and returns 1 product media with given id
func (ps *SqlProductMediaStore) Get(id string) (*model.ProductMedia, error) {
	var res model.ProductMedia
	err := ps.GetReplicaX().Get(
		&res,
		"SELECT * FROM "+store.ProductMediaTableName+" WHERE Id = ?",
		id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductMediaTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find product media with id=%s", id)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of product medias with given id
func (ps *SqlProductMediaStore) FilterByOption(option *model.ProductMediaFilterOption) ([]*model.ProductMedia, error) {
	query := ps.GetQueryBuilder().
		Select("*").
		From(store.ProductMediaTableName).
		OrderBy(store.TableOrderingMap[store.ProductMediaTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID)
	}
	if option.Type != nil {
		query = query.Where(option.Type)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.ProductMedia
	err = ps.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product medias by given option")
	}

	return res, nil
}
