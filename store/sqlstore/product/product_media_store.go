package product

import (
	"fmt"

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
func (ps *SqlProductMediaStore) Upsert(tx *gorm.DB, medias model.ProductMedias) (model.ProductMedias, error) {
	if tx == nil {
		tx = ps.GetMaster()
	}
	err := tx.Save(medias).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert product medias")
	}

	return medias, nil
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
		From(model.ProductMediaTableName).
		Where(option.Conditions)

	// parse options
	if option.VariantID != nil {
		query = query.
			InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.MediaID = %[2]s.Id", model.ProductVariantMediaTableName, model.ProductMediaTableName)).
			Where(option.VariantID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	db := ps.GetReplica()
	if len(option.Preloads) > 0 {
		for _, preload := range option.Preloads {
			db = db.Preload(preload)
		}
	}

	var res model.ProductMedias
	err = db.Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product medias by given option")
	}

	return res, nil
}

func (p *SqlProductMediaStore) Delete(tx *gorm.DB, ids []string) (int64, error) {
	if tx == nil {
		tx = p.GetMaster()
	}

	result := tx.Where("Id IN ?", ids).Delete(&model.ProductMedia{})
	if result.Error != nil {
		return 0, errors.Wrap(result.Error, "failed to delete product medias by given ids")
	}

	return result.RowsAffected, nil
}
