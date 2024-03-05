package product

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

func (cs *SqlCollectionStore) FilterByOption(option model_helper.CollectionFilterOptions) (model.CollectionSlice, error) {
	query := cs.GetQueryBuilder().
		Select(model.CollectionTableName + ".*").
		From(model.CollectionTableName).
		Where(option.Conditions)

	// parse options
	if option.ProductID != nil {
		query = query.
			InnerJoin(
				fmt.Sprintf(
					"%[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.CollectionProductRelationTableName,  // 1
					model.CollectionTableName,                 // 2
					model.CollectionProductColumnCollectionID, // 3
					model.CollectionColumnId,                  // 4
				),
			).
			Where(option.ProductID)
	}
	if option.VoucherID != nil {
		query = query.
			InnerJoin(
				fmt.Sprintf(
					"%[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.VoucherCollectionTableName, // 1
					model.CollectionTableName,        // 2
					"collection_id",                  // 3
					model.CollectionColumnId,         // 4
				),
			).
			Where(option.VoucherID)
	}
	if option.SaleID != nil {
		query = query.
			InnerJoin(
				fmt.Sprintf(
					"%[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.SaleCollectionTableName, // 1
					model.CollectionTableName,     // 2
					"collection_id",               // 3
					model.CollectionColumnId,      // 4
				),
			).
			Where(option.SaleID)
	}

	if option.ChannelListingPublicationDate != nil ||
		option.ChannelListingIsPublished != nil ||

		option.ChannelListingChannelSlug != nil ||
		option.ChannelListingChannelIsActive != nil {
		query = query.
			InnerJoin(
				fmt.Sprintf(
					"%[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.CollectionChannelListingTableName,          // 1
					model.CollectionTableName,                        // 2
					model.CollectionChannelListingColumnCollectionID, // 3
					model.CollectionColumnId,                         // 4
				),
			).
			Where(option.ChannelListingPublicationDate).
			Where(option.ChannelListingIsPublished)

		if option.ChannelListingChannelSlug != nil ||
			option.ChannelListingChannelIsActive != nil {
			query = query.
				InnerJoin(
					fmt.Sprintf(
						"%[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
						model.ChannelTableName,                        // 1
						model.CollectionChannelListingTableName,       // 2
						model.ChannelColumnId,                         // 3
						model.CollectionChannelListingColumnChannelID, // 4
					),
				).
				Where(option.ChannelListingChannelSlug).
				Where(option.ChannelListingChannelIsActive)
		}
	}

	// annotate for sorting
	if option.AnnotateProductCount {
		query = query.
			Column(fmt.Sprintf(`COUNT (%s.Id) AS "%s.ProductCount"`, model.CollectionProductRelationTableName, model.CollectionTableName)).
			LeftJoin(
				fmt.Sprintf(
					"%[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.CollectionProductRelationTableName,  // 1
					model.CollectionTableName,                 // 2
					model.CollectionProductColumnCollectionID, // 3
					model.CollectionColumnId,                  // 4
				),
			).
			GroupBy(model.CollectionTableName + ".Id")

	} else if option.AnnotateIsPublished && option.ChannelSlugForIsPublishedAndPublicationDateAnnotation != "" {
		isPublishedExpr := fmt.Sprintf(
			`(
			SELECT
				%[1]s.%[4]s
			FROM %[1]s
			INNER JOIN
				%[2]s ON %[2]s.%[5]s = %[1]s.%[6]s
			WHERE (
				%[2]s.%[7]s = ?
				AND %[1]s.%[8]s = %[3]s.%[9]s
			)
			LIMIT 1
		) AS "%[3]s.IsPublished"`,
			model.CollectionChannelListingTableName,          // 1
			model.ChannelTableName,                           // 2
			model.CollectionTableName,                        // 3
			model.PublishableColumnIsPublished,               // 4
			model.ChannelColumnId,                            // 5
			model.CollectionChannelListingColumnChannelID,    // 6
			model.ChannelColumnSlug,                          // 7
			model.CollectionChannelListingColumnCollectionID, // 8
			model.CollectionColumnId,                         // 9
		)

		query = query.Column(isPublishedExpr, option.ChannelSlugForIsPublishedAndPublicationDateAnnotation)

	} else if option.AnnotatePublicationDate && option.ChannelSlugForIsPublishedAndPublicationDateAnnotation != "" {
		publicationDateExpr := fmt.Sprintf(`(
			SELECT
				%[1]s.%[4]s
			FROM
				%[1]s
			INNER JOIN
				%[2]s ON %[2]s.%[5]s = %[1]s.%[6]s
			WHERE (
				%[2]s.%[7]s = ?
				AND %[1]s.%[8]s = %[3]s.%[9]s
			)
			LIMIT 1
		) AS "%[3]s.PublicationDate"`,
			model.CollectionChannelListingTableName,          // 1
			model.ChannelTableName,                           // 2
			model.CollectionTableName,                        // 3
			model.PublishableColumnPublicationDate,           // 4
			model.ChannelColumnId,                            // 5
			model.CollectionChannelListingColumnChannelID,    // 6
			model.ChannelColumnSlug,                          // 7
			model.CollectionChannelListingColumnCollectionID, // 8
			model.CollectionColumnId,                         // 9
		)

		query = query.Column(publicationDateExpr, option.ChannelSlugForIsPublishedAndPublicationDateAnnotation)
	}

	var runner = cs.GetReplica()

	// NOTE: count total must be applied right before pagination like this
	var totalCount int64
	if option.CountTotal {
		countQuery, args, err := cs.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "FilterByOption_CountTotal_ToSql")
		}

		err = runner.Raw(countQuery, args...).Scan(&totalCount).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total collections by options")
		}
	}

	// apply pagination
	option.GraphqlPaginationValues.AddPaginationToSelectBuilderIfNeeded(&query)

	queryString, args, err := query.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	for _, preload := range option.Preload {
		runner = runner.Preload(preload)
	}

	rows, err := runner.Raw(queryString, args...).Rows()
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to find collections with given options")
	}
	defer rows.Close()

	var res model.Collections

	for rows.Next() {
		var col model.Collection
		scanFields := cs.ScanFields(&col)
		// check if we have annotation here:
		if option.AnnotateProductCount {
			scanFields = append(scanFields, &col.ProductCount)
		} else if option.AnnotateIsPublished {
			scanFields = append(scanFields, &col.IsPublished)
		} else if option.AnnotatePublicationDate {
			scanFields = append(scanFields, &col.PublicationDate)
		}

		err := rows.Scan(scanFields...)
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to scan a row of collection")
		}

		res = append(res, &col)
	}

	return totalCount, res, nil
}

func (s *SqlCollectionStore) Delete(tx boil.ContextTransactor, ids []string) error {
	if tx == nil {
		tx = s.GetMaster()
	}

	_, err := model.Collections(model.CollectionWhere.ID.IN(ids)).DeleteAll(tx)
	return err
}
