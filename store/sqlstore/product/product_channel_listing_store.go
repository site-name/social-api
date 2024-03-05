package product

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlProductChannelListingStore struct {
	store.Store
}

func NewSqlProductChannelListingStore(s store.Store) store.ProductChannelListingStore {
	return &SqlProductChannelListingStore{s}
}

func (ps *SqlProductChannelListingStore) BulkUpsert(transaction boil.ContextTransactor, listings model.ProductChannelListingSlice) (model.ProductChannelListingSlice, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	for _, listing := range listings {
		if listing == nil {
			continue
		}

		isSaving := listing.ID == ""
		if isSaving {
			model_helper.ProductChannelListingPreSave(listing)
		} else {
			model_helper.ProductChannelListingCommonPre(listing)
		}

		if err := model_helper.ProductChannelListingIsValid(*listing); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = listing.Insert(transaction, boil.Infer())
		} else {
			_, err = listing.Update(transaction, boil.Blacklist(
				model.ProductChannelListingColumns.CreatedAt,
			))
		}

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{model.ProductChannelListingTableColumns.ProductID, model.ProductChannelListingColumns.ChannelID, "product_channel_listings_product_id_channel_id_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.ProductChannelListings, "ProductID/ChannelID", "duplicate")
			}
			return nil, err
		}
	}

	return listings, nil
}

func (ps *SqlProductChannelListingStore) Get(listingID string) (*model.ProductChannelListing, error) {
	listing, err := model.FindProductChannelListing(ps.GetReplica(), listingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ProductChannelListings, listingID)
		}
		return nil, err
	}

	return listing, nil
}

func (ps *SqlProductChannelListingStore) FilterByOption(option model_helper.ProductChannelListingFilterOption) (model.ProductChannelListingSlice, error) {
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
