package discount

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlVoucherStore struct {
	store.Store
}

func NewSqlDiscountVoucherStore(sqlStore store.Store) store.DiscountVoucherStore {
	return &SqlVoucherStore{sqlStore}

}

func (vs *SqlVoucherStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"ShopID",
		"Type",
		"Name",
		"Code",
		"UsageLimit",
		"Used",
		"StartDate",
		"EndDate",
		"ApplyOncePerOrder",
		"ApplyOncePerCustomer",
		"OnlyForStaff",
		"DiscountValueType",
		"Countries",
		"MinCheckoutItemsQuantity",
		"CreateAt",
		"UpdateAt",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (vs *SqlVoucherStore) ScanFields(voucher model.Voucher) []interface{} {
	return []interface{}{
		&voucher.Id,
		&voucher.ShopID,
		&voucher.Type,
		&voucher.Name,
		&voucher.Code,
		&voucher.UsageLimit,
		&voucher.Used,
		&voucher.StartDate,
		&voucher.EndDate,
		&voucher.ApplyOncePerOrder,
		&voucher.ApplyOncePerCustomer,
		&voucher.OnlyForStaff,
		&voucher.DiscountValueType,
		&voucher.Countries,
		&voucher.MinCheckoutItemsQuantity,
		&voucher.CreateAt,
		&voucher.UpdateAt,
		&voucher.Metadata,
		&voucher.PrivateMetadata,
	}
}

// Upsert saves or updates given voucher then returns it with an error
func (vs *SqlVoucherStore) Upsert(voucher *model.Voucher) (*model.Voucher, error) {
	var saving bool

	if voucher.Id == "" {
		voucher.PreSave()
		saving = true
	} else {
		voucher.PreUpdate()
	}
	if appErr := voucher.IsValid(); appErr != nil {
		return nil, appErr
	}

	var (
		oldVoucher *model.Voucher
		err        error
		numUpdated int64
	)

	if saving {
		query := "INSERT INTO " + store.VoucherTableName + "(" + vs.ModelFields("").Join(",") + ") VALUES (" + vs.ModelFields(":").Join(",") + ")"
		_, err = vs.GetMasterX().NamedExec(query, voucher)
	} else {

		oldVoucher, err = vs.Get(voucher.Id)
		if err != nil {
			return nil, err
		}

		voucher.Used = oldVoucher.Used

		query := "UPDATE " + store.VoucherTableName + " SET " + vs.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = vs.GetMasterX().NamedExec(query, voucher)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if vs.IsUniqueConstraintError(err, []string{"Code", "vouchers_code_key"}) {
			return nil, store.NewErrInvalidInput(store.VoucherTableName, "code", voucher.Code)
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher with id=%s", voucher.Id)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple vouchers were updated: %d instead of 1", numUpdated)
	}

	return voucher, nil
}

// Get finds a voucher with given id, then returns it with an error
func (vs *SqlVoucherStore) Get(voucherID string) (*model.Voucher, error) {
	var res model.Voucher
	err := vs.GetReplicaX().Get(&res, "SELECT * FROM "+store.VoucherTableName+" WHERE Id = ?", voucherID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherTableName, voucherID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher with id=%s", voucherID)
	}

	return &res, nil
}

func (vs *SqlVoucherStore) commonQueryBuilder(option *model.VoucherFilterOption) squirrel.SelectBuilder {
	query := vs.
		GetQueryBuilder().
		Select(vs.ModelFields(store.VoucherTableName + ".")...).
		From(store.VoucherTableName).
		OrderBy(store.TableOrderingMap[store.VoucherTableName])

	// parse options:
	if option.UsageLimit != nil {
		query = query.Where(option.UsageLimit)
	}
	if option.EndDate != nil {
		query = query.Where(option.EndDate)
	}
	if option.StartDate != nil {
		query = query.Where(option.StartDate)
	}
	if option.Code != nil {
		query = query.Where(option.Code)
	}
	if option.ChannelListingSlug != nil || option.ChannelListingActive != nil {
		query = query.
			InnerJoin(store.VoucherChannelListingTableName + " ON (VoucherChannelListings.VoucherID = Vouchers.Id)").
			InnerJoin(store.ChannelTableName + " ON (Channels.Id = VoucherChannelListings.ChannelID)")

		if option.ChannelListingSlug != nil {
			query = query.Where(option.ChannelListingSlug)
		}

		if option.ChannelListingActive != nil {
			query = query.Where(squirrel.Eq{"Channels.IsActive": *option.ChannelListingActive})
		}
	}
	if option.WithLook {
		query = query.Suffix("FOR UPDATE") // SELECT ... FOR UPDATE
	}

	return query
}

// FilterVouchersByOption finds vouchers bases on given option.
func (vs *SqlVoucherStore) FilterVouchersByOption(option *model.VoucherFilterOption) ([]*model.Voucher, error) {
	queryString, args, err := vs.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterVouchersByOption_tosql")
	}

	var vouchers []*model.Voucher
	err = vs.GetReplicaX().Select(&vouchers, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find vouchers based on given option")
	}

	return vouchers, nil
}

// GetByOptions finds and returns 1 voucher filtered using given options
func (vs *SqlVoucherStore) GetByOptions(options *model.VoucherFilterOption) (*model.Voucher, error) {
	queryString, args, err := vs.commonQueryBuilder(options).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOptions_tosql")
	}

	var res model.Voucher
	err = vs.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherTableName, "options")
		}
		return nil, errors.Wrap(err, "failed voucher with given options")
	}

	return &res, nil
}

// ExpiredVouchers finds and returns vouchers that are expired before given date
func (vs *SqlVoucherStore) ExpiredVouchers(date *time.Time) ([]*model.Voucher, error) {
	if date == nil {
		date = util.NewTime(time.Now())
	}
	beginOfDate := util.StartOfDay(*date)

	var res []*model.Voucher
	err := vs.GetReplicaX().Select(&res, "SELECT * FROM "+store.VoucherTableName+" WHERE (Used >= UsageLimit OR EndDate < $1) AND StartDate < $2", beginOfDate, beginOfDate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find expired vouchers with given date")
	}

	return res, nil
}
