package product

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlDigitalContentUrlStore struct {
	store.Store
}

func NewSqlDigitalContentUrlStore(s store.Store) store.DigitalContentUrlStore {
	return &SqlDigitalContentUrlStore{s}
}

func (ps *SqlDigitalContentUrlStore) Upsert(contentURL model.DigitalContentURL) (*model.DigitalContentURL, error) {
	isSaving := contentURL.ID == ""
	if isSaving {
		model_helper.DigitalContentUrlPreSave(&contentURL)
	}

	if err := model_helper.DigitalContentUrlIsValid(contentURL); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = contentURL.Insert(ps.GetMaster(), boil.Infer())
	} else {
		_, err = contentURL.Update(ps.GetMaster(), boil.Blacklist(
			model.DigitalContentURLColumns.CreatedAt,
			model.DigitalContentURLColumns.Token,
		))
	}

	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{model.DigitalContentURLColumns.Token, "digital_content_urls_token_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.DigitalContentUrls, model.DigitalContentURLColumns.Token, contentURL.Token)
		}
		if ps.IsUniqueConstraintError(err, []string{model.DigitalContentURLColumns.LineID, "digital_content_urls_line_id_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.DigitalContentUrls, model.DigitalContentURLColumns.LineID, contentURL.LineID)
		}
		return nil, err
	}

	return &contentURL, nil
}

func (ps *SqlDigitalContentUrlStore) Get(id string) (*model.DigitalContentURL, error) {
	contentURL, err := model.FindDigitalContentURL(ps.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.DigitalContentUrls, id)
		}
		return nil, err
	}

	return contentURL, nil
}

func (s *SqlDigitalContentUrlStore) FilterByOptions(options model_helper.DigitalContentUrlFilterOptions) (model.DigitalContentURLSlice, error) {
	return model.DigitalContentUrls(options.Conditions...).All(s.GetReplica())
}
