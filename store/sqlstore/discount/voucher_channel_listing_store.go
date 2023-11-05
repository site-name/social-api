package discount

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlVoucherChannelListingStore struct {
	store.Store
}

var voucherChannelListingDuplicateList = []string{
	"VoucherID", "ChannelID", "voucherid_channelid_key",
}

func NewSqlVoucherChannelListingStore(sqlStore store.Store) store.VoucherChannelListingStore {
	return &SqlVoucherChannelListingStore{sqlStore}
}

// upsert check given listing's Id to decide whether to create or update it. Then returns a listing with an error
func (vcls *SqlVoucherChannelListingStore) Upsert(transaction *gorm.DB, voucherChannelListings []*model.VoucherChannelListing) ([]*model.VoucherChannelListing, error) {
	if transaction == nil {
		transaction = vcls.GetMaster()
	}

	for _, listing := range voucherChannelListings {
		var err error
		if listing.Id == "" {
			err = transaction.Create(listing).Error
		} else {
			// keep non-editable fields intact
			// Refer to https://gorm.io/docs/update.html#Updates-multiple-columns
			listing.CreateAt = 0
			err = transaction.Model(listing).Updates(listing).Error
		}
		if err != nil {
			if vcls.IsUniqueConstraintError(err, voucherChannelListingDuplicateList) {
				return nil, store.NewErrInvalidInput(model.VoucherChannelListingTableName, "VoucherID/ChannelID", "duplicate values")
			}
			return nil, errors.Wrap(err, "failed to upsert voucher channel listing")
		}
	}

	return voucherChannelListings, nil
}

// Get finds a listing with given id, then returns it with an error
func (vcls *SqlVoucherChannelListingStore) Get(voucherChannelListingID string) (*model.VoucherChannelListing, error) {
	var res model.VoucherChannelListing
	err := vcls.GetReplica().First(&res, "Id = ?", voucherChannelListingID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.VoucherChannelListingTableName, voucherChannelListingID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher channel listing with id=%s", voucherChannelListingID)
	}

	res.PopulateNonDbFields()
	return &res, nil
}

// FilterbyOption finds and returns a list of voucher channel listing relationship instances filtered by given option
func (vcls *SqlVoucherChannelListingStore) FilterbyOption(option *model.VoucherChannelListingFilterOption) ([]*model.VoucherChannelListing, error) {
	args, err := store.BuildSqlizer(option.Conditions, "VoucherChannelListingsByOptions")
	if err != nil {
		return nil, err
	}
	var res []*model.VoucherChannelListing
	err = vcls.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher channel listing relationship instances with given option")
	}

	return res, nil
}

func (s *SqlVoucherChannelListingStore) Delete(transaction *gorm.DB, option *model.VoucherChannelListingFilterOption) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	args, err := store.BuildSqlizer(option.Conditions, "VoucherChannelListingDelete")
	if err != nil {
		return err
	}

	err = transaction.Delete(model.VoucherChannelListingTableName, args...).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete voucher channel listing")
	}
	return nil
}
