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
			InnerJoin(model.CollectionProductRelationTableName + " ON Collections.Id = ProductCollections.CollectionID").
			Where(option.ProductID)
	}
	if option.VoucherID != nil {
		query = query.
			InnerJoin(model.VoucherCollectionTableName + " ON Collections.Id = VoucherCollections.CollectionID").
			Where(option.VoucherID)
	}
	if option.SaleID != nil {
		query = query.
			InnerJoin(model.SaleCollectionTableName + " ON Collections.Id = SaleCollections.CollectionID").
			Where(option.SaleID)
	}

	if option.ChannelListingPublicationDate != nil ||
		option.ChannelListingIsPublished != nil ||
		option.ChannelListingChannelSlug != nil ||
		option.ChannelListingChannelIsActive != nil {
		query = query.
			InnerJoin(model.CollectionChannelListingTableName + " ON Collections.Id = CollectionChannelListings.CollectionID").
			Where(option.ChannelListingPublicationDate).
			Where(option.ChannelListingIsPublished)

		if option.ChannelListingChannelSlug != nil ||
			option.ChannelListingChannelIsActive != nil {
			query = query.
				InnerJoin(model.ChannelTableName + " ON Channels.Id = CollectionChannelListings.ChannelID").
				Where(option.ChannelListingChannelSlug).
				Where(option.ChannelListingChannelIsActive)
		}
	}

	// annotate for sorting
	if option.AnnotateProductCount {
		query = query.
			Column(fmt.Sprintf(`COUNT (%s.Id) AS "%s.ProductCount"`, model.CollectionProductRelationTableName, model.CollectionTableName)).
			LeftJoin(model.CollectionProductRelationTableName + " ON ProductCollections.CollectionID = Collections.Id").
			GroupBy(model.CollectionTableName + ".Id")

	} else if option.AnnotateIsPublished && option.ChannelSlugForIsPublishedAndPublicationDateAnnotation != "" {
		isPublishedExpr := fmt.Sprintf(`(
			SELECT
				CCL.IsPublished
			FROM %[1]s CCL
			INNER JOIN %[2]s C ON C.Id = CCL.ChannelID
			WHERE (
				C.Slug = ?
				AND CCL.CollectionID = %[3]s.Id
			)
			LIMIT 1
		) AS "%[3]s.IsPublished"`,
			model.CollectionChannelListingTableName,
			model.ChannelTableName,
			model.CollectionTableName,
		)
		query = query.Column(isPublishedExpr, option.ChannelSlugForIsPublishedAndPublicationDateAnnotation)

	} else if option.AnnotatePublicationDate && option.ChannelSlugForIsPublishedAndPublicationDateAnnotation != "" {
		publicationDateExpr := fmt.Sprintf(`(
			SELECT
				CCL.PublicationDate
			FROM
				%[1]s CCL
			INNER JOIN %[2]s C ON C.Id = CCL.ChannelID
			WHERE (
				C.Slug = ?
				AND CCL.CollectionID = %[3]s.Id
			)
			LIMIT 1
		) AS "%[3]s.PublicationDate"`,
			model.CollectionChannelListingTableName,
			model.ChannelTableName,
			model.CollectionTableName,
		)

		query = query.Column(publicationDateExpr, option.ChannelSlugForIsPublishedAndPublicationDateAnnotation)
	}

	// NOTE: count total must be applied right before pagination like this
	var totalCount int64
	if option.CountTotal {
		countQuery, args, err := cs.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "FilterByOption_CountTotal_ToSql")
		}

		err = cs.GetReplica().Raw(countQuery, args...).Scan(&totalCount).Error
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

	var runner = cs.GetReplica()
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
	err := s.GetMaster().Raw("DELETE FROM "+model.CollectionTableName+" WHERE Id IN ?", ids).Error
	if err != nil {
		errors.Wrap(err, "failed to delete collection(s) by given ids")
	}

	return nil
}
