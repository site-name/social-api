package product

import (
	"database/sql"
	"fmt"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlProductMediaStore struct {
	store.Store
}

func NewSqlProductMediaStore(s store.Store) store.ProductMediaStore {
	return &SqlProductMediaStore{s}
}

func (ps *SqlProductMediaStore) Upsert(tx boil.ContextTransactor, medias model.ProductMediumSlice) (model.ProductMediumSlice, error) {
	if tx == nil {
		tx = ps.GetMaster()
	}

	for _, media := range medias {
		if media == nil {
			continue
		}

		isSaving := media.ID == ""
		if isSaving {
			model_helper.ProductMediaPreSave(media)
		} else {
			model_helper.ProductMediaCommonPre(media)
		}

		if err := model_helper.ProductMediaIsValid(*media); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = media.Insert(tx, boil.Infer())
		} else {
			_, err = media.Update(tx, boil.Blacklist(
				model.ProductMediumColumns.CreatedAt,
			))
		}

		if err != nil {
			return nil, err
		}
	}

	return medias, nil
}

func (ps *SqlProductMediaStore) Get(id string) (*model.ProductMedium, error) {
	media, err := model.FindProductMedium(ps.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ProductMedia, id)
		}
		return nil, err
	}

	return media, nil
}

func (ps *SqlProductMediaStore) FilterByOption(option model_helper.ProductMediaFilterOption) (model.ProductMediumSlice, error) {
	conds := option.Conditions
	if option.VariantID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.VariantMedia, model.VariantMediumTableColumns.MediaID, model.ProductMediumTableColumns.ID)),
			option.VariantID,
		)
	}

	for _, load := range option.Preloads {
		conds = append(conds, qm.Load(load))
	}

	return model.ProductMedia(conds...).All(ps.GetReplica())
}

func (p *SqlProductMediaStore) Delete(tx boil.ContextTransactor, ids []string) (int64, error) {
	if tx == nil {
		tx = p.GetMaster()
	}
	return model.ProductMedia(model.ProductMediumWhere.ID.IN(ids)).DeleteAll(tx)
}
