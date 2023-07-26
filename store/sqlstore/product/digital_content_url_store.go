package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlDigitalContentUrlStore struct {
	store.Store
}

func NewSqlDigitalContentUrlStore(s store.Store) store.DigitalContentUrlStore {
	return &SqlDigitalContentUrlStore{s}
}

// Upsert inserts or updates given digital content url into database then returns it
func (ps *SqlDigitalContentUrlStore) Upsert(contentURL *model.DigitalContentUrl) (*model.DigitalContentUrl, error) {
	err := ps.GetMaster().Save(contentURL).Error
	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"Token", "digitalcontenturls_token_key"}) {
			return nil, store.NewErrInvalidInput(model.DigitalContentURLTableName, "Token", contentURL.Token)
		}
		if ps.IsUniqueConstraintError(err, []string{"LineID", "digitalcontenturls_lineid_key"}) {
			return nil, store.NewErrInvalidInput(model.DigitalContentURLTableName, "LineID", contentURL.LineID)
		}
		return nil, errors.Wrapf(err, "failed to upsert content url with id=%s", contentURL.Id)
	}

	return contentURL, nil
}

// Get finds and returns a digital content url with given id
func (ps *SqlDigitalContentUrlStore) Get(id string) (*model.DigitalContentUrl, error) {
	var res model.DigitalContentUrl

	err := ps.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.DigitalContentURLTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find digital content url with id=%s", id)
	}

	return &res, nil
}

func (s *SqlDigitalContentUrlStore) FilterByOptions(options *model.DigitalContentUrlFilterOptions) ([]*model.DigitalContentUrl, error) {
	var res []*model.DigitalContentUrl
	err := s.GetReplica().Find(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find digital content urls by options")
	}
	return res, nil
}
