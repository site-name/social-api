package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlDigitalContentUrlStore struct {
	store.Store
}

func NewSqlDigitalContentUrlStore(s store.Store) store.DigitalContentUrlStore {
	return &SqlDigitalContentUrlStore{s}
}

func (s *SqlDigitalContentUrlStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
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
			query := "INSERT INTO " + store.DigitalContentURLTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
			_, err = ps.GetMasterX().NamedExec(query, contentURL)

		} else {

			query := "UPDATE " + store.DigitalContentURLTableName + " SET " + ps.
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
				return nil, store.NewErrInvalidInput(store.DigitalContentURLTableName, "LineID", contentURL.LineID)
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

	err := ps.GetReplicaX().Get(&res, "SELECT * FROM "+store.DigitalContentURLTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.DigitalContentURLTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find digital content url with id=%s", id)
	}

	return &res, nil
}
