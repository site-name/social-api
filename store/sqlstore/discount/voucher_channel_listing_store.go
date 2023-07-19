package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlVoucherChannelListingStore struct {
	store.Store
}

var VoucherChannelListingDuplicateList = []string{
	"VoucherID", "ChannelID", "voucherchannellistings_voucherid_channelid_key",
}

func (s *SqlVoucherChannelListingStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id", "CreateAt", "VoucherID", "ChannelID", "DiscountValue", "Currency", "MinSpentAmount",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func NewSqlVoucherChannelListingStore(sqlStore store.Store) store.VoucherChannelListingStore {
	return &SqlVoucherChannelListingStore{sqlStore}
}

// upsert check given listing's Id to decide whether to create or update it. Then returns a listing with an error
func (vcls *SqlVoucherChannelListingStore) Upsert(voucherChannelListing *model.VoucherChannelListing) (*model.VoucherChannelListing, error) {
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
		query := "INSERT INTO " + model.VoucherChannelListingTableName + "(" + vcls.ModelFields("").Join(",") + ") VALUES (" + vcls.ModelFields(":").Join(",") + ")"
		_, err = vcls.GetMaster().NamedExec(query, voucherChannelListing)

	} else {
		query := "UPDATE " + model.VoucherChannelListingTableName + " SET " + vcls.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = vcls.GetMaster().NamedExec(query, voucherChannelListing)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if vcls.IsUniqueConstraintError(err, VoucherChannelListingDuplicateList) {
			return nil, store.NewErrInvalidInput(model.VoucherChannelListingTableName, "VoucherID/ChannelID", "duplicate values")
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
func (vcls *SqlVoucherChannelListingStore) Get(voucherChannelListingID string) (*model.VoucherChannelListing, error) {
	var res model.VoucherChannelListing

	err := vcls.GetReplica().Get(&res, "SELECT * FROM "+model.VoucherChannelListingTableName+" WHERE Id = ?", voucherChannelListingID)
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
	query := vcls.GetQueryBuilder().
		Select("*").
		From(model.VoucherChannelListingTableName)

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.ChannelID != nil {
		query = query.Where(option.ChannelID)
	}
	if option.VoucherID != nil {
		query = query.Where(option.VoucherID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var res []*model.VoucherChannelListing
	err = vcls.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher channel listing relationship instances with given option")
	}

	return res, nil
}
