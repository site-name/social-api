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

var VoucherChannelListingDuplicateList = []string{
	"VoucherID", "ChannelID", "voucherchannellistings_voucherid_channelid_key",
}

func NewSqlVoucherChannelListingStore(sqlStore store.Store) store.VoucherChannelListingStore {
	return &SqlVoucherChannelListingStore{sqlStore}
}

// upsert check given listing's Id to decide whether to create or update it. Then returns a listing with an error
func (vcls *SqlVoucherChannelListingStore) Upsert(voucherChannelListing *model.VoucherChannelListing) (*model.VoucherChannelListing, error) {
	var err error
	if voucherChannelListing.Id == "" {
		err = vcls.GetMaster().Create(voucherChannelListing).Error
	} else {
		// keep non-editable fields intact
		// Refer to https://gorm.io/docs/update.html#Updates-multiple-columns
		voucherChannelListing.CreateAt = 0
		err = vcls.GetMaster().Table(model.VoucherChannelListingTableName).Updates(voucherChannelListing).Error
	}
	if err != nil {
		if vcls.IsUniqueConstraintError(err, VoucherChannelListingDuplicateList) {
			return nil, store.NewErrInvalidInput(model.VoucherChannelListingTableName, "VoucherID/ChannelID", "duplicate values")
		}
		return nil, errors.Wrap(err, "failed to upsert voucher channel listing")
	}
	return voucherChannelListing, nil
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
	var res []*model.VoucherChannelListing
	err := vcls.GetReplica().Find(&res, store.BuildSqlizer(option.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher channel listing relationship instances with given option")
	}

	return res, nil
}
