package discount

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlVoucherStore struct {
	store.Store
}

var (
	VoucherUniqueList = []string{"Code", "vouchers_code_key"}
)

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
		if vs.IsUniqueConstraintError(err, VoucherUniqueList) {
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
		Select("*").
		From(store.VoucherTableName).
		OrderBy(store.TableOrderingMap[store.VoucherTableName])

	// check usage limit
	if option.UsageLimit != nil {
		query = query.Where(option.UsageLimit.ToSquirrel("Vouchers.UsageLimit"))
	}

	// check end date
	if option.EndDate != nil {
		query = query.Where(option.EndDate.ToSquirrel("Vouchers.EndDate"))
	}

	// check start date
	if option.StartDate != nil {
		query = query.Where(option.StartDate.ToSquirrel("Vouchers.StartDate"))
	}

	// check code
	if option.Code != nil {
		query = query.Where(option.Code.ToSquirrel("Vouchers.Code"))
	}

	// check channel listing channel slug
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
		query = query.Suffix("FOR UPDATE")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sql_tosql")
	}

	var vouchers []*product_and_discount.Voucher
	_, err = vs.GetReplica().Select(&vouchers, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find vouchers based on given option")
	}

	return vouchers, nil
}
