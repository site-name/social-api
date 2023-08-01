package discount

import (
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlVoucherStore struct {
	store.Store
}

func NewSqlDiscountVoucherStore(sqlStore store.Store) store.DiscountVoucherStore {
	return &SqlVoucherStore{sqlStore}

}

func (vs *SqlVoucherStore) ScanFields(voucher *model.Voucher) []interface{} {
	return []interface{}{
		&voucher.Id,
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
	var err error
	if voucher.Id == "" {
		err = vs.GetMaster().Create(voucher).Error
	} else {
		voucher.Used = 0 // prevent updates by gorm
		voucher.CreateAt = 0
		err = vs.GetMaster().Model(voucher).Updates(voucher).Error
	}
	if err != nil {
		if vs.IsUniqueConstraintError(err, []string{"Code", "vouchers_code_key"}) {
			return nil, store.NewErrInvalidInput(model.VoucherTableName, "code", voucher.Code)
		}
		return nil, errors.Wrapf(err, "failed to upsert voucher with id=%s", voucher.Id)
	}

	return voucher, nil
}

// Get finds a voucher with given id, then returns it with an error
func (vs *SqlVoucherStore) Get(voucherID string) (*model.Voucher, error) {
	var res model.Voucher
	err := vs.GetReplica().First(&res, "Id = ?", voucherID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.VoucherTableName, voucherID)
		}
		return nil, errors.Wrapf(err, "failed to find voucher with id=%s", voucherID)
	}

	return &res, nil
}

func (vs *SqlVoucherStore) commonQueryBuilder(option *model.VoucherFilterOption) squirrel.SelectBuilder {
	query := vs.
		GetQueryBuilder().
		Select(model.VoucherTableName + ".*").
		From(model.VoucherTableName).
		Where(option.Conditions)

	// parse options:
	if option.VoucherChannelListing_ChannelSlug != nil || option.VoucherChannelListing_ChannelIsActive != nil {
		query = query.
			InnerJoin(model.VoucherChannelListingTableName + " ON (VoucherChannelListings.VoucherID = Vouchers.Id)").
			InnerJoin(model.ChannelTableName + " ON (Channels.Id = VoucherChannelListings.ChannelID)")

		if option.VoucherChannelListing_ChannelSlug != nil {
			query = query.Where(option.VoucherChannelListing_ChannelSlug)
		}
		if option.VoucherChannelListing_ChannelIsActive != nil {
			query = query.Where(option.VoucherChannelListing_ChannelIsActive)
		}
	}
	if option.ForUpdate && option.Trnsaction != nil {
		query = query.Suffix("FOR UPDATE")
	}

	return query
}

// FilterVouchersByOption finds vouchers bases on given option.
func (vs *SqlVoucherStore) FilterVouchersByOption(option *model.VoucherFilterOption) ([]*model.Voucher, error) {
	queryString, args, err := vs.commonQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterVouchersByOption_tosql")
	}

	runner := vs.GetReplica()
	if option.Trnsaction != nil {
		runner = option.Trnsaction
	}

	var vouchers []*model.Voucher
	err = runner.Raw(queryString, args...).Scan(&vouchers).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find vouchers based on given option")
	}

	return vouchers, nil
}

// GetByOptions finds and returns 1 voucher filtered using given options
func (vs *SqlVoucherStore) GetByOptions(options *model.VoucherFilterOption) (*model.Voucher, error) {
	vouchers, err := vs.FilterVouchersByOption(options)
	if err != nil {
		return nil, err
	}
	if len(vouchers) == 0 {
		return nil, store.NewErrNotFound(model.VoucherTableName, "options")
	}
	return vouchers[0], nil
}

// ExpiredVouchers finds and returns vouchers that are expired before given date
func (vs *SqlVoucherStore) ExpiredVouchers(date *time.Time) ([]*model.Voucher, error) {
	if date == nil {
		date = util.NewTime(time.Now())
	}
	beginOfDate := util.StartOfDay(*date)

	var res []*model.Voucher
	err := vs.GetReplica().Raw("SELECT * FROM "+model.VoucherTableName+" WHERE (Used >= UsageLimit OR EndDate < $1) AND StartDate < $2", beginOfDate, beginOfDate).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find expired vouchers with given date")
	}

	return res, nil
}

func (s *SqlVoucherStore) ToggleVoucherRelations(transaction *gorm.DB, vouchers model.Vouchers, collectionIds, productIds, variantIds, categoryIds []string, isDelete bool) error {
	if len(vouchers) == 0 {
		return errors.New("please speficy relations")
	}
	if transaction == nil {
		transaction = s.GetMaster()
	}

	relationsMap := map[string]any{
		"Products":        lo.Map(productIds, func(id string, _ int) *model.Product { return &model.Product{Id: id} }),
		"Collections":     lo.Map(collectionIds, func(id string, _ int) *model.Collection { return &model.Collection{Id: id} }),
		"ProductVariants": lo.Map(variantIds, func(id string, _ int) *model.ProductVariant { return &model.ProductVariant{Id: id} }),
		"Categories":      lo.Map(categoryIds, func(id string, _ int) *model.Category { return &model.Category{Id: id} }),
	}

	for associationName, relations := range relationsMap {
		for _, voucher := range vouchers {
			if voucher != nil {
				switch {
				case isDelete:
					err := transaction.Model(voucher).Association(associationName).Delete(relations)
					if err != nil {
						return errors.Wrap(err, "failed to delete voucher "+strings.ToLower(associationName)+" relations")
					}
				default:
					err := transaction.Model(voucher).Association(associationName).Append(relations)
					if err != nil {
						return errors.Wrap(err, "failed to insert voucher "+strings.ToLower(associationName)+" relations")
					}
				}
			}
		}
	}

	return nil
}
