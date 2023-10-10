package product

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlCollectionStore struct {
	store.Store
}

func NewSqlCollectionStore(s store.Store) store.CollectionStore {
	return &SqlCollectionStore{s}
}

func (ps *SqlCollectionStore) ScanFields(col *model.Collection) []interface{} {
	return []interface{}{
		&col.Id,
		&col.Name,
		&col.Slug,
		&col.BackgroundImage,
		&col.BackgroundImageAlt,
		&col.Description,
		&col.Metadata,
		&col.PrivateMetadata,
		&col.SeoTitle,
		&col.SeoDescription,
	}
}

// Upsert depends on given collection's Id property to decide update or insert the collection
func (cs *SqlCollectionStore) Upsert(collection *model.Collection) (*model.Collection, error) {
	err := cs.GetMaster().Save(collection).Error
	if err != nil {
		if cs.IsUniqueConstraintError(err, []string{"slug_unique_key", "slug"}) {
			return nil, store.NewErrInvalidInput(model.CollectionTableName, "Slug", collection.Slug)
		}
		return nil, errors.Wrap(err, "failed to upsert collection")
	}
	return collection, nil
}

// Get finds and returns collection with given collectionID
func (cs *SqlCollectionStore) Get(collectionID string) (*model.Collection, error) {
	var res model.Collection
	err := cs.GetReplica().First(&res, "Id = ?", collectionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.CollectionTableName, collectionID)
		}
		return nil, errors.Wrapf(err, "failed to find collection with id=%s", collectionID)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of collections satisfy the given option.
//
// NOTE: make sure to provide `ShopID` before calling me.
func (cs *SqlCollectionStore) FilterByOption(option *model.CollectionFilterOption) (int64, []*model.Collection, error) {
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

func (s *SqlCollectionStore) Delete(ids ...string) error {
	err := s.GetMaster().Delete(&model.Collection{}, "Id IN ?", ids).Error
	if err != nil {
		errors.Wrap(err, "failed to delete collection(s) by given ids")
	}

	return nil
}
