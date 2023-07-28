package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlProductMediaStore struct {
	store.Store
}

func NewSqlProductMediaStore(s store.Store) store.ProductMediaStore {
	return &SqlProductMediaStore{s}
}

// Upsert depends on given media's Id property to decide insert or update it
func (ps *SqlProductMediaStore) Upsert(media *model.ProductMedia) (*model.ProductMedia, error) {
	err := ps.GetMaster().Save(media).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert product media with id=%s", media.Id)
	}

	return media, nil
}

// Get finds and returns 1 product media with given id
func (ps *SqlProductMediaStore) Get(id string) (*model.ProductMedia, error) {
	var res model.ProductMedia
	err := ps.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ProductMediaTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find product media with id=%s", id)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of product medias with given id
func (ps *SqlProductMediaStore) FilterByOption(option *model.ProductMediaFilterOption) ([]*model.ProductMedia, error) {
	query := ps.GetQueryBuilder().
		Select(model.ProductMediaTableName + ".*").
		From(model.ProductMediaTableName).Where(option.Conditions)

	// parse options
	if option.VariantID != nil {
		query = query.
			LeftJoin("VariantMedias ON VariantMedias.MediaID = ProductMedias.Id").
			Where(option.VariantID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.ProductMedia
	err = ps.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product medias by given option")
	}

	return res, nil
}
