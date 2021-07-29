package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherChannelListingStore struct {
	store.Store
}

var (
	VoucherChannelListingDuplicateList = []string{
		"VoucherID", "ChannelID", "voucherchannellistings_voucherid_channelid_key",
	}
)

func NewSqlVoucherChannelListingStore(sqlStore store.Store) store.VoucherChannelListingStore {
	vcls := &SqlVoucherChannelListingStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.VoucherChannelListing{}, store.VoucherChannelListingTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("VoucherID", "ChannelID")
	}

	return vcls
}

func (vcls *SqlVoucherChannelListingStore) CreateIndexesIfNotExists() {
	vcls.CreateForeignKeyIfNotExists(store.VoucherChannelListingTableName, "VoucherID", store.VoucherTableName, "Id", true)
	vcls.CreateForeignKeyIfNotExists(store.VoucherChannelListingTableName, "ChannelID", store.ChannelTableName, "Id", true)
}

// upsert check given listing's Id to decide whether to create or update it. Then returns a listing with an error
func (vcls *SqlVoucherChannelListingStore) Upsert(voucherChannelListing *product_and_discount.VoucherChannelListing) (*product_and_discount.VoucherChannelListing, error) {
	var saving bool
	if voucherChannelListing.Id == "" {
		saving = true
		voucherChannelListing.PreSave()
	}

	if err := voucherChannelListing.IsValid(); err != nil {
		return nil, err
	}

	var err error
	var numUpdated int64
	if saving {
		err = vcls.GetMaster().Insert(voucherChannelListing)
	} else {
		// validate if the listing does exist:
		_, err = vcls.Get(voucherChannelListing.Id)
		if err != nil {
			return nil, err
		}
		numUpdated, err = vcls.GetMaster().Update(voucherChannelListing)
	}

	if err != nil {
		if vcls.IsUniqueConstraintError(err, VoucherChannelListingDuplicateList) {
			return nil, store.NewErrInvalidInput(store.VoucherChannelListingTableName, "VoucherID/ChannelID", "duplicate values")
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher channel listing with id=%s", voucherChannelListing.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple voucher channel listings updated: %d instead of 1", numUpdated)
	}

	voucherChannelListing.PopulateNonDbFields()
	return voucherChannelListing, nil
}

// Get finds a listing with given id, then returns it with an error
func (vcls *SqlVoucherChannelListingStore) Get(voucherChannelListingID string) (*product_and_discount.VoucherChannelListing, error) {
	result, err := vcls.GetReplica().Get(product_and_discount.VoucherChannelListing{}, voucherChannelListingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherChannelListingTableName, voucherChannelListingID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher channel listing with id=%s", voucherChannelListingID)
	}

	listing := result.(*product_and_discount.VoucherChannelListing)
	listing.PopulateNonDbFields()
	return listing, nil
}

// FilterByVoucherAndChannel finds a list of listings that belong to given voucher and own given channel
func (vcls *SqlVoucherChannelListingStore) FilterByVoucherAndChannel(voucherID string, channelID string) ([]*product_and_discount.VoucherChannelListing, error) {
	var listings []*product_and_discount.VoucherChannelListing
	_, err := vcls.GetReplica().Select(
		&listings,
		`SELECT * FROM `+store.VoucherChannelListingTableName+`
		WHERE (
			VoucherID = :VoucherID AND ChannelID = :ChannelID
		)
		ORDER BY CreateAt`, // since ids are UUIDs, not number so order by creation time is used instead
		map[string]interface{}{
			"VoucherID": voucherID,
			"ChannelID": channelID,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherChannelListingTableName, "voucherID="+voucherID+", channelID="+channelID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher channel listing with VoucherID = %s, ChannelID = %s", voucherID, channelID)
	}

	product_and_discount.VoucherChannelListingList(listings).PopulateNonDbFields()
	return listings, nil
}
