package discount

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlVoucherChannelListingStore struct {
	store.Store
}

func NewSqlVoucherChannelListingStore(sqlStore store.Store) store.VoucherChannelListingStore {
	return &SqlVoucherChannelListingStore{sqlStore}
}

// upsert check given listing's Id to decide whether to create or update it. Then returns a listing with an error
func (vcls *SqlVoucherChannelListingStore) Upsert(transaction boil.ContextTransactor, voucherChannelListings model.VoucherChannelListingSlice) (model.VoucherChannelListingSlice, error) {
	if transaction == nil {
		transaction = vcls.GetMaster()
	}

	for _, listing := range voucherChannelListings {
		if listing == nil {
			continue
		}

		isSaving := false
		if listing.ID == "" {
			isSaving = true
			model_helper.VoucherChannelListingPreSave(listing)
		}

		if err := model_helper.VoucherChannelListingIsValid(*listing); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = listing.Insert(transaction, boil.Infer())
		} else {
			_, err = listing.Update(transaction, boil.Blacklist(model.VoucherChannelListingColumns.CreatedAt))
		}

		if err != nil {
			if vcls.IsUniqueConstraintError(err, []string{model.VoucherChannelListingColumns.VoucherID, model.VoucherChannelListingColumns.ChannelID, "voucher_channel_listings_voucher_id_channel_id_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.VoucherChannelListings, "voucherid_channelid_key", "unique")
			}
			return nil, err
		}
	}

	return voucherChannelListings, nil
}

// Get finds a listing with given id, then returns it with an error
func (vcls *SqlVoucherChannelListingStore) Get(voucherChannelListingID string) (*model.VoucherChannelListing, error) {
	record, err := model.FindVoucherChannelListing(vcls.GetReplica(), voucherChannelListingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.VoucherChannelListings, voucherChannelListingID)
		}
		return nil, err
	}

	return record, nil
}

// FilterbyOption finds and returns a list of voucher channel listing relationship instances filtered by given option
func (vcls *SqlVoucherChannelListingStore) FilterbyOption(option model_helper.VoucherChannelListingFilterOption) (model.VoucherChannelListingSlice, error) {
	return model.VoucherChannelListings(option.Conditions...).All(vcls.GetReplica())
}

func (s *SqlVoucherChannelListingStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}
	_, err := model.VoucherChannelListings(model.VoucherChannelListingWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
