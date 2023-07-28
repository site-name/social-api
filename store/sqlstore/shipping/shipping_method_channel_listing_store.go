package shipping

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlShippingMethodChannelListingStore struct {
	store.Store
}

func NewSqlShippingMethodChannelListingStore(s store.Store) store.ShippingMethodChannelListingStore {
	return &SqlShippingMethodChannelListingStore{s}
}

// Upsert depends on given listing's Id to decide whether to save or update the listing
func (s *SqlShippingMethodChannelListingStore) Upsert(transaction *gorm.DB, listings model.ShippingMethodChannelListings) (model.ShippingMethodChannelListings, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	for _, listing := range listings {
		err := transaction.Save(listing).Error
		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"ShippingMethodID", "ChannelID", "shippingmethodid_channelid_key"}) {
				return nil, store.NewErrInvalidInput(model.ShippingMethodChannelListingTableName, "ShippingMethodID/ChannelID", listing.ShippingMethodID+"/"+listing.ChannelID)
			}
			return nil, errors.Wrap(err, "failed to upsert shipping method channel listing")
		}
	}

	return listings, nil
}

// Get finds a shipping method channel listing with given listingID
func (s *SqlShippingMethodChannelListingStore) Get(listingID string) (*model.ShippingMethodChannelListing, error) {
	var res model.ShippingMethodChannelListing
	err := s.GetReplica().First(&res, "Id = ?", listingID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ShippingMethodChannelListingTableName, listingID)
		}
		return nil, errors.Wrapf(err, "failed to find shipping method channel listing with id=%s", listingID)
	}

	res.PopulateNonDbFields()
	return &res, nil
}

// FilterByOption returns a list of shipping method channel listings based on given option. result sorted by creation time ASC
func (s *SqlShippingMethodChannelListingStore) FilterByOption(option *model.ShippingMethodChannelListingFilterOption) ([]*model.ShippingMethodChannelListing, error) {
	query := s.GetQueryBuilder().
		Select(model.ShippingMethodChannelListingTableName + ".*").
		From(model.ShippingMethodChannelListingTableName).Where(option.Conditions)

	// parse filter option
	if option.ChannelSlug != nil {
		query = query.
			InnerJoin(model.ChannelTableName + " ON Channels.Id = ShippingMethodChannelListings.ChannelID").
			Where(option.ChannelSlug)
	}
	if option.ShippingMethod_ShippingZoneID_Inner != nil {
		query = query.
			InnerJoin(model.ShippingMethodTableName + " ON ShippingMethods.Id = ShippingMethodChannelListings.ShippingMethodID").
			InnerJoin(model.ShippingZoneTableName + " ON ShippingZones.Id = ShippingMethods.ShippingZoneID").
			Where(option.ShippingMethod_ShippingZoneID_Inner)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_tosql")
	}

	var res []*model.ShippingMethodChannelListing
	err = s.GetReplica().Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping method channel listings by option")
	}

	return res, nil
}

func (s *SqlShippingMethodChannelListingStore) BulkDelete(transaction *gorm.DB, options *model.ShippingMethodChannelListingFilterOption) error {
	query := s.GetQueryBuilder().Delete(model.ShippingMethodChannelListingTableName).Where(options.Conditions)

	queryStr, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "BulkDelete_ToSql")
	}

	if transaction == nil {
		transaction = s.GetMaster()
	}

	err = transaction.Raw(queryStr, args...).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete shipping method channel listings")
	}

	return nil
}
