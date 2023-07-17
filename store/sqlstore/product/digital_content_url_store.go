package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlDigitalContentUrlStore struct {
	store.Store
}

func NewSqlDigitalContentUrlStore(s store.Store) store.DigitalContentUrlStore {
	return &SqlDigitalContentUrlStore{s}
}

func (s *SqlDigitalContentUrlStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Token",
		"ContentID",
		"CreateAt",
		"DownloadNum",
		"LineID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert inserts or updates given digital content url into database then returns it
func (ps *SqlDigitalContentUrlStore) Upsert(contentURL *model.DigitalContentUrl) (*model.DigitalContentUrl, error) {

	var isSaving bool
	if contentURL.Id == "" {
		isSaving = true
		contentURL.PreSave()
	}

	if err := contentURL.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	for {
		if isSaving {
			query := "INSERT INTO " + model.DigitalContentURLTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
			_, err = ps.GetMasterX().NamedExec(query, contentURL)

		} else {

			query := "UPDATE " + model.DigitalContentURLTableName + " SET " + ps.
				ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id=:Id"

			var result sql.Result
			result, err = ps.GetMasterX().NamedExec(query, contentURL)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
		}

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"Token", "digitalcontenturls_token_key"}) {
				contentURL.NewToken(true)
				continue
			}
			if ps.IsUniqueConstraintError(err, []string{"LineID", "digitalcontenturls_lineid_key"}) {
				return nil, store.NewErrInvalidInput(model.DigitalContentURLTableName, "LineID", contentURL.LineID)
			}
			return nil, errors.Wrapf(err, "failed to upsert content url with id=%s", contentURL.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("multiple content urls were updated for content url with id=%s: %d instead of 1", contentURL.Id, numUpdated)
		}

		return contentURL, nil
	}
}

// Get finds and returns a digital content url with given id
func (ps *SqlDigitalContentUrlStore) Get(id string) (*model.DigitalContentUrl, error) {
	var res model.DigitalContentUrl

	err := ps.GetReplicaX().Get(&res, "SELECT * FROM "+model.DigitalContentURLTableName+" WHERE Id = ?", id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.DigitalContentURLTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find digital content url with id=%s", id)
	}

	return &res, nil
}

func (s *SqlDigitalContentUrlStore) FilterByOptions(options *model.DigitalContentUrlFilterOptions) ([]*model.DigitalContentUrl, error) {
	query := s.GetQueryBuilder().Select("*").From(model.DigitalContentURLTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.Token != nil {
		query = query.Where(options.Token)
	}
	if options.ContentID != nil {
		query = query.Where(options.ContentID)
	}
	if options.LineID != nil {
		query = query.Where(options.LineID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.DigitalContentUrl
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find digital content urls by options")
	}
	return res, nil
}
