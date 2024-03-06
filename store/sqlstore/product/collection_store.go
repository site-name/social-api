package product

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlCollectionStore struct {
	store.Store
}

func NewSqlCollectionStore(s store.Store) store.CollectionStore {
	return &SqlCollectionStore{s}
}

func (cs *SqlCollectionStore) Upsert(collection model.Collection) (*model.Collection, error) {
	isSaving := collection.ID == ""
	if isSaving {
		model_helper.CollectionPreSave(&collection)
	} else {
		model_helper.CollectionCommonPre(&collection)
	}

	if err := model_helper.CollectionIsValid(collection); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = collection.Insert(cs.GetMaster(), boil.Infer())
	} else {
		_, err = collection.Update(cs.GetMaster(), boil.Infer())
	}

	if err != nil {
		if cs.IsUniqueConstraintError(err, []string{"slug_unique_key", model.CollectionColumns.Slug}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Collections, model.CollectionColumns.Slug, collection.Slug)
		}
		return nil, err
	}

	return &collection, nil
}

func (cs *SqlCollectionStore) Get(collectionID string) (*model.Collection, error) {
	collection, err := model.FindCollection(cs.GetReplica(), collectionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Collections, collectionID)
		}
		return nil, err
	}

	return collection, nil
}

func (cs *SqlCollectionStore) commonQueryBuilder(option model_helper.CollectionFilterOptions) []qm.QueryMod {
	conds := option.Conditions
	conds = append(conds, qm.Select(model.TableNames.Collections+".*"))

	if option.ProductID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductCollections, model.ProductCollectionTableColumns.CollectionID, model.CollectionTableColumns.ID)),
			option.ProductID,
		)
	}
	if option.VoucherID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.VoucherCollections, model.VoucherCollectionTableColumns.CollectionID, model.CollectionTableColumns.ID)),
			option.VoucherID,
		)
	}
	if option.SaleID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.SaleCollections, model.SaleCollectionTableColumns.CollectionID, model.CollectionTableColumns.ID)),
			option.SaleID,
		)
	}
	if option.RelatedCollectionChannelListingConds != nil ||
		option.RelatedCollectionChannelListingChannelConds != nil {

		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.CollectionChannelListings, model.CollectionChannelListingTableColumns.CollectionID, model.CollectionTableColumns.ID)),
			option.RelatedCollectionChannelListingConds,
		)

		if option.RelatedCollectionChannelListingChannelConds != nil {
			conds = append(
				conds,
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.CollectionChannelListingTableColumns.ChannelID)),
				option.RelatedCollectionChannelListingChannelConds,
			)
		}
	}

	if option.AnnotateProductCount {
		conds = append(
			conds,
			qm.Select(fmt.Sprintf("COUNT %s AS %q", model.ProductCollectionTableColumns.ID, model_helper.CustomCollectionColumns.ProductCount)),
			qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductCollections, model.ProductCollectionTableColumns.CollectionID, model.CollectionTableColumns.ID)),
			qm.GroupBy(model.CollectionTableColumns.ID),
		)
	}
	if option.AnnotateIsPublished && slug.IsSlug(option.ChannelSlugForIsPublishedAndPublicationDateAnnotation) {
		conds = append(
			conds,
			qm.Select(
				fmt.Sprintf(
					`(
						SELECT %[1]s
						FROM %[2]s
						INNER JOIN %[3]s ON %[4]s = %[5]s
						WHERE %[6]s = '%[7]s'
						AND %[8]s = %[9]s
						LIMIT 1
					) AS "%[10]s"`,
					model.CollectionChannelListingTableColumns.IsPublished,       // 1
					model.TableNames.CollectionChannelListings,                   // 2
					model.TableNames.Channels,                                    // 3
					model.ChannelTableColumns.ID,                                 // 4
					model.CollectionChannelListingTableColumns.ChannelID,         // 5
					model.ChannelTableColumns.Slug,                               // 6
					option.ChannelSlugForIsPublishedAndPublicationDateAnnotation, // 7
					model.CollectionChannelListingTableColumns.CollectionID,      // 8
					model.CollectionTableColumns.ID,                              // 9
					model_helper.CustomCollectionColumns.IsPublished,             // 10
				),
			),
		)
	}
	if option.AnnotatePublicationDate && slug.IsSlug(option.ChannelSlugForIsPublishedAndPublicationDateAnnotation) {
		conds = append(
			conds,
			qm.Select(
				fmt.Sprintf(
					`(
						SELECT %[1]s
						FROM %[2]s
						INNER JOIN %[3]s ON %[4]s = %[5]s
						WHERE %[6]s = '%[7]s'
						AND %[8]s = %[9]s
						LIMIT 1
					) AS "%[10]s"`,
					model.CollectionChannelListingTableColumns.PublicationDate,   // 1
					model.TableNames.CollectionChannelListings,                   // 2
					model.TableNames.Channels,                                    // 3
					model.ChannelTableColumns.ID,                                 // 4
					model.CollectionChannelListingTableColumns.ChannelID,         // 5
					model.ChannelTableColumns.Slug,                               // 6
					option.ChannelSlugForIsPublishedAndPublicationDateAnnotation, // 7
					model.CollectionChannelListingTableColumns.CollectionID,      // 8
					model.CollectionTableColumns.ID,                              // 9
					model_helper.CustomCollectionColumns.PublicationDate,         // 10
				),
			),
		)
	}

	return conds
}

func (cs *SqlCollectionStore) FilterByOption(option model_helper.CollectionFilterOptions) (model_helper.CustomCollectionSlice, error) {
	conds := cs.commonQueryBuilder(option)
	var customCollections model_helper.CustomCollectionSlice

	err := model.Collections(conds...).Bind(context.Background(), cs.GetReplica(), &customCollections)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find collections with given options")
	}

	return customCollections, nil
}

func (s *SqlCollectionStore) Delete(tx boil.ContextTransactor, ids []string) error {
	if tx == nil {
		tx = s.GetMaster()
	}

	_, err := model.Collections(model.CollectionWhere.ID.IN(ids)).DeleteAll(tx)
	return err
}
