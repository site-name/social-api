package product

import (
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

func (ps *SqlProductChannelListingStore) ScanFields(prd *model.ProductChannelListing) []interface{} {
	return []interface{}{
		&prd.Id,
		&prd.ProductID,
		&prd.ChannelID,
		&prd.VisibleInListings,
		&prd.AvailableForPurchase,
		&prd.Currency,
		&prd.DiscountedPriceAmount,
		&prd.CreateAt,
		&prd.PublicationDate,
		&prd.IsPublished,
	}
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
	selectFields := []string{model.ProductChannelListingTableName + ".*"}
	if option.PrefetchChannel {
		selectFields = append(selectFields, model.ChannelTableName+".*")
	}
	query := ps.
		GetQueryBuilder().
		Select(selectFields...).
		From(model.ProductChannelListingTableName).
		Where(option.Conditions)

	// parse option
	if option.ChannelSlug != nil || option.PrefetchChannel {
		query = query.InnerJoin(model.ChannelTableName + " ON Channels.Id = ProductChannelListings.ChannelID")

		if option.ChannelSlug != nil {
			query = query.Where(option.ChannelSlug)
		}
	}
	if option.ProductVariantsId != nil {
		query = query.
			InnerJoin(model.ProductTableName + " ON Products.Id = ProductChannelListings.ProductID").
			InnerJoin(model.ProductVariantTableName + " ON ProductVariants.ProductID = Products.Id").
			Where(option.ProductVariantsId)
	}

	sqlString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	rows, err := ps.GetReplica().Raw(sqlString, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product channel listings with given option")
	}
	defer rows.Close()

	var listings model.ProductChannelListings

	for rows.Next() {
		var (
			productChannelListing model.ProductChannelListing
			channel               model.Channel
			scanFields            = ps.ScanFields(&productChannelListing)
		)
		if option.PrefetchChannel {
			scanFields = append(scanFields, ps.Channel().ScanFields(&channel)...)
		}

		err := rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of product channel listing")
		}

		if option.PrefetchChannel {
			productChannelListing.SetChannel(&channel)
		}

		listings = append(listings, &productChannelListing)
	}

	return listings, nil
}
