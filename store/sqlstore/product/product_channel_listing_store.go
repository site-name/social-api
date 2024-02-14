package product

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlProductChannelListingStore struct {
	store.Store
}

func NewSqlProductChannelListingStore(s store.Store) store.ProductChannelListingStore {
	return &SqlProductChannelListingStore{s}
}

// BulkUpsert performs bulk upsert on given product channel listings
func (ps *SqlProductChannelListingStore) BulkUpsert(transaction *gorm.DB, listings []*model.ProductChannelListing) ([]*model.ProductChannelListing, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	for _, listing := range listings {
		var err error

		if listing.Id == "" {
			err = transaction.Create(listing).Error
		} else {
			err = transaction.Model(listing).Updates(listing).Error
		}

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"ProductID", "ChannelID", "productid_channelid_key"}) {
				return nil, store.NewErrInvalidInput(model.ProductChannelListingTableName, "ProductID/ChannelID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert product channel listing with id=%s", listing.Id)
		}
	}

	return listings, nil
}

func (ps *SqlProductChannelListingStore) Get(listingID string) (*model.ProductChannelListing, error) {
	var res model.ProductChannelListing

	err := ps.GetReplica().First(&res, "Id = ?", listingID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ProductChannelListingTableName, listingID)
		}
		return nil, errors.Wrapf(err, "failed to find product channel listing with id=%s", listingID)
	}

	return &res, nil
}

// FilterByOption finds and returns all ProductChannelListings filtered by given option
func (ps *SqlProductChannelListingStore) FilterByOption(option *model.ProductChannelListingFilterOption) ([]*model.ProductChannelListing, error) {
	db := ps.GetReplica()
	conditions := squirrel.And{}
	for _, preload := range option.Preloads {
		db = db.Preload(preload)
	}

	if option.Conditions != nil {
		conditions = append(conditions, option.Conditions)
	}
	if option.RelatedChannelConditions != nil {
		db = db.Joins(fmt.Sprintf(
			"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
			model.ChannelTableName,                     // 1
			model.ProductChannelListingTableName,       // 2
			model.ChannelColumnId,                      // 3
			model.ProductChannelListingColumnChannelID, // 4
		))
		conditions = append(conditions, option.RelatedChannelConditions)
	}
	if option.ProductVariantsId != nil {
		conditions = append(conditions, option.ProductVariantsId)
		db = db.
			Joins(fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.ProductTableName,                     // 1
				model.ProductChannelListingTableName,       // 2
				model.ProductColumnId,                      // 3
				model.ProductChannelListingColumnProductID, // 4
			)).
			Joins(fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.ProductVariantTableName,       // 1
				model.ProductTableName,              // 2
				model.ProductVariantColumnProductID, // 3
				model.ProductColumnId,               // 4
			))
	}

	args, err := store.BuildSqlizer(conditions, "ProductChannelListing_FilterByOptions")
	if err != nil {
		return nil, err
	}
	var res model.ProductChannelListings
	err = db.Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product channel listings with given option")
	}

	return res, nil
}
