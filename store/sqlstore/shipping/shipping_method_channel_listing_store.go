package shipping

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlShippingMethodChannelListingStore struct {
	store.Store
}

func NewSqlShippingMethodChannelListingStore(s store.Store) store.ShippingMethodChannelListingStore {
	return &SqlShippingMethodChannelListingStore{s}
}

// Upsert depends on given listing's Id to decide whether to save or update the listing
func (s *SqlShippingMethodChannelListingStore) Upsert(transaction boil.ContextTransactor, listings model.ShippingMethodChannelListingSlice) (model.ShippingMethodChannelListingSlice, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	for _, listing := range listings {
		if listing == nil {
			continue
		}

		isSaving := listing.ID == ""
		if isSaving {
			model_helper.ShippingMethodChannelListingPreSave(listing)
		} else {
			model_helper.ShippingMethodChannelListingCommonPre(listing)
		}

		if err := model_helper.ShippingMethodChannelListingIsValid(*listing); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = listing.Insert(transaction, boil.Infer())
		} else {
			_, err = listing.Update(transaction, boil.Blacklist(model.ShippingMethodChannelListingColumns.CreatedAt))
		}

		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"shipping_method_channel_listings_shipping_method_id_channel_id_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.ShippingMethodChannelListings, model.ShippingMethodChannelListingColumns.ShippingMethodID+"/"+model.ShippingMethodChannelListingColumns.ChannelID, "unique")
			}
			return nil, err
		}
	}

	return listings, nil
}

// Get finds a shipping method channel listing with given listingID
func (s *SqlShippingMethodChannelListingStore) Get(listingID string) (*model.ShippingMethodChannelListing, error) {
	listing, err := model.FindShippingMethodChannelListing(s.GetReplica(), listingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ShippingMethodChannelListings, listingID)
		}
		return nil, err
	}

	return listing, nil
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

func (s *SqlShippingMethodChannelListingStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	_, err := model.ShippingMethodChannelListings(model.ShippingMethodChannelListingWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
