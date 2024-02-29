package shipping

import (
	"database/sql"
	"fmt"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlShippingMethodChannelListingStore struct {
	store.Store
}

func NewSqlShippingMethodChannelListingStore(s store.Store) store.ShippingMethodChannelListingStore {
	return &SqlShippingMethodChannelListingStore{s}
}

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

func (s *SqlShippingMethodChannelListingStore) FilterByOption(option model_helper.ShippingMethodChannelListingFilterOption) (model.ShippingMethodChannelListingSlice, error) {
	conds := option.Conditions
	if option.ChannelSlug != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.ShippingMethodChannelListingTableColumns.ChannelID)),
			option.ChannelSlug,
		)
	}
	if option.ShippingMethod_ShippingZoneID_Inner != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingMethods, model.ShippingMethodTableColumns.ID, model.ShippingMethodChannelListingTableColumns.ShippingMethodID)),
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZones, model.ShippingZoneTableColumns.ID, model.ShippingMethodTableColumns.ShippingZoneID)),
			option.ShippingMethod_ShippingZoneID_Inner,
		)
	}

	return model.ShippingMethodChannelListings(conds...).All(s.GetReplica())
}

func (s *SqlShippingMethodChannelListingStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	_, err := model.ShippingMethodChannelListings(model.ShippingMethodChannelListingWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
