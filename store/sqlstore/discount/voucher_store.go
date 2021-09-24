package discount

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlVoucherStore struct {
	store.Store
}

func NewSqlDiscountVoucherStore(sqlStore store.Store) store.DiscountVoucherStore {
	vs := &SqlVoucherStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Voucher{}, store.VoucherTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShopID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(product_and_discount.VOUCHER_TYPE_MAX_LENGTH)
		table.ColMap("Code").SetMaxSize(product_and_discount.VOUCHER_CODE_MAX_LENGTH).SetUnique(true)
		table.ColMap("Name").SetMaxSize(product_and_discount.VOUCHER_NAME_MAX_LENGTH)
		table.ColMap("Countries").SetMaxSize(model.MULTIPLE_COUNTRIES_MAX_LENGTH)
		table.ColMap("DiscountValueType").SetMaxSize(product_and_discount.VOUCHER_DISCOUNT_VALUE_TYPE_MAX_LENGTH)
	}

	return vs
}

func (vs *SqlVoucherStore) ModelFields() []string {
	return []string{
		"Vouchers.Id",
		"Vouchers.ShopID",
		"Vouchers.Type",
		"Vouchers.Name",
		"Vouchers.Code",
		"Vouchers.UsageLimit",
		"Vouchers.Used",
		"Vouchers.StartDate",
		"Vouchers.EndDate",
		"Vouchers.ApplyOncePerOrder",
		"Vouchers.ApplyOncePerCustomer",
		"Vouchers.OnlyForStaff",
		"Vouchers.DiscountValueType",
		"Vouchers.Countries",
		"Vouchers.MinCheckoutItemsQuantity",
		"Vouchers.CreateAt",
		"Vouchers.UpdateAt",
		"Vouchers.Metadata",
		"Vouchers.PrivateMetadata",
	}
}

func (vs *SqlVoucherStore) ScanFields(voucher product_and_discount.Voucher) []interface{} {
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

func (vs *SqlVoucherStore) CreateIndexesIfNotExists() {
	vs.CreateIndexIfNotExists("idx_vouchers_code", store.VoucherTableName, "Code")
	vs.CreateForeignKeyIfNotExists(store.VoucherTableName, "ShopID", store.ShopTableName, "Id", true)
}

// Upsert saves or updates given voucher then returns it with an error
func (vs *SqlVoucherStore) Upsert(voucher *product_and_discount.Voucher) (*product_and_discount.Voucher, error) {
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
		oldVoucher *product_and_discount.Voucher
		err        error
		numUpdated int64
	)

	if saving {
		err = vs.GetMaster().Insert(voucher)
	} else {
		oldVoucher, err = vs.Get(voucher.Id)
		if err != nil {
			return nil, err
		}

		voucher.Used = oldVoucher.Used
		numUpdated, err = vs.GetMaster().Update(voucher)
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
func (vs *SqlVoucherStore) Get(voucherID string) (*product_and_discount.Voucher, error) {
	var res product_and_discount.Voucher
	err := vs.GetReplica().SelectOne(&res, "SELECT * FROM "+store.VoucherTableName+" WHERE Id = :ID", map[string]interface{}{"ID": voucherID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.VoucherTableName, voucherID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher with id=%s", voucherID)
	}

	return &res, nil
}

// FilterVouchersByOption finds vouchers bases on given option.
func (vs *SqlVoucherStore) FilterVouchersByOption(option *product_and_discount.VoucherFilterOption) ([]*product_and_discount.Voucher, error) {
	query := vs.
		GetQueryBuilder().
		Select(vs.ModelFields()...).
		From(store.VoucherTableName).
		OrderBy(store.TableOrderingMap[store.VoucherTableName])

	// parse options:
	if option.UsageLimit != nil {
		query = query.Where(option.UsageLimit.ToSquirrel("Vouchers.UsageLimit"))
	}
	if option.EndDate != nil {
		query = query.Where(option.EndDate.ToSquirrel("Vouchers.EndDate"))
	}
	if option.StartDate != nil {
		query = query.Where(option.StartDate.ToSquirrel("Vouchers.StartDate"))
	}
	if option.Code != nil {
		query = query.Where(option.Code.ToSquirrel("Vouchers.Code"))
	}
	if option.ChannelListingSlug != nil || option.ChannelListingActive != nil {
		query = query.
			InnerJoin(store.VoucherChannelListingTableName + " ON (VoucherChannelListings.VoucherID = Vouchers.Id)").
			InnerJoin(store.ChannelTableName + " ON (Channels.Id = VoucherChannelListings.ChannelID)")

		if option.ChannelListingSlug != nil {
			query = query.Where(option.ChannelListingSlug.ToSquirrel("Channels.Slug"))
		}

		if option.ChannelListingActive != nil {
			query = query.Where(squirrel.Eq{"Channels.IsActive": *option.ChannelListingActive})
		}
	}
	if option.WithLook {
		query = query.Suffix("FOR UPDATE") // SELECT ... FOR UPDATE
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterVouchersByOption_tosql")
	}

	var vouchers []*product_and_discount.Voucher
	_, err = vs.GetReplica().Select(&vouchers, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find vouchers based on given option")
	}

	return vouchers, nil
}

// ExpiredVouchers finds and returns vouchers that are expired before given date
func (vs *SqlVoucherStore) ExpiredVouchers(date *time.Time) ([]*product_and_discount.Voucher, error) {
	if date == nil {
		date = model.NewTime(time.Now())
	}
	beginOfDate := util.StartOfDay(*date)

	var res []*product_and_discount.Voucher
	_, err := vs.GetReplica().Select(&res, "SELECT * FROM "+store.VoucherTableName+" WHERE (Used >= UsageLimit OR EndDate < $1) AND StartDate < $2", beginOfDate, beginOfDate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find expired vouchers with given date")
	}

	return res, nil
}
