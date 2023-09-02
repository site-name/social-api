package discount

import (
	"strings"
	"time"

	"github.com/gosimple/slug"
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

// FilterVouchersByOption finds vouchers bases on given option.
func (vs *SqlVoucherStore) FilterVouchersByOption(option *model.VoucherFilterOption) (int64, []*model.Voucher, error) {
	query := vs.
		GetQueryBuilder().
		Select(model.VoucherTableName + ".*").
		From(model.VoucherTableName).
		Where(option.Conditions)

	if option.Annotate_DiscountValue ||
		option.Annotate_MinimumSpentAmount {

		// check if channel provided:
		if !slug.IsSlug(option.ChannelSlug) {
			return 0, nil, store.NewErrInvalidInput("SqlVoucherStore.commonQueryBuilder", "option.ChannelSlug", option.ChannelSlug)
		}

		query = query.
			LeftJoin(model.VoucherChannelListingTableName + " ON (VoucherChannelListings.VoucherID = Vouchers.Id)").
			LeftJoin(model.ChannelTableName + " ON (Channels.Id = VoucherChannelListings.ChannelID)").
			GroupBy(model.VoucherTableName + ".Id")

		// NOTE: these annotation are for sorting voucher
		if option.Annotate_MinimumSpentAmount {
			query = query.
				Column(`MIN(
					VoucherChannelListings.MinSpentAmount
				) FILTER (
					WHERE Channels.Slug = ?
				) AS "Vouchers.MinSpentAmount"`, option.ChannelSlug)
		} else if option.Annotate_DiscountValue {
			query = query.
				Column(`MIN(
					VoucherChannelListings.DiscountValue
				) FILTER (
					WHERE Channels.Slug = ?
				) AS "Vouchers.DiscountValue"`, option.ChannelSlug)
		}

	} else if option.VoucherChannelListing_ChannelSlug != nil ||
		option.VoucherChannelListing_ChannelIsActive != nil {
		query = query.
			InnerJoin(model.VoucherChannelListingTableName + " ON (VoucherChannelListings.VoucherID = Vouchers.Id)").
			InnerJoin(model.ChannelTableName + " ON (Channels.Id = VoucherChannelListings.ChannelID)").
			Where(option.VoucherChannelListing_ChannelSlug).
			Where(option.VoucherChannelListing_ChannelIsActive)
	}

	if option.ForUpdate && option.Transaction != nil {
		query = query.Suffix("FOR UPDATE")
	}

	// count total vouchers if required
	var totalVoucher int64
	if option.CountTotal {
		queryStr, args, err := vs.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "FilterByOptions_Count_ToSql")
		}

		err = vs.GetReplica().Raw(queryStr, args...).Scan(&totalVoucher).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total vouchers that satisfy given options")
		}
	}

	option.GraphqlPaginationValues.AddPaginationToSelectBuilderIfNeeded(&query)

	querystr, args, err := query.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	runner := vs.GetReplica()
	if option.Transaction != nil {
		runner = option.Transaction
	}

	rows, err := runner.Raw(querystr, args...).Rows()
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to find vouchers by given options")
	}
	defer rows.Close()

	var vouchers model.Vouchers

	for rows.Next() {
		var (
			voucher    model.Voucher
			scanFields = vs.ScanFields(&voucher)
		)
		if option.Annotate_MinimumSpentAmount {
			scanFields = append(scanFields, &voucher.MinSpentAmount)
		} else if option.Annotate_DiscountValue {
			scanFields = append(scanFields, &voucher.DiscountValue)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to scan a row of voucher")
		}

		vouchers = append(vouchers, &voucher)
	}

	return totalVoucher, vouchers, nil
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

func (s *SqlVoucherStore) ToggleVoucherRelations(transaction *gorm.DB, vouchers model.Vouchers, collectionIds, productIds, variantIds, categoryIds []model.UUID, isDelete bool) error {
	if len(vouchers) == 0 {
		return errors.New("please speficy relations")
	}
	if transaction == nil {
		transaction = s.GetMaster()
	}

	relationsMap := map[string]any{
		"Products":        lo.Map(productIds, func(id model.UUID, _ int) *model.Product { return &model.Product{Id: id} }),
		"Collections":     lo.Map(collectionIds, func(id model.UUID, _ int) *model.Collection { return &model.Collection{Id: id} }),
		"ProductVariants": lo.Map(variantIds, func(id model.UUID, _ int) *model.ProductVariant { return &model.ProductVariant{Id: id} }),
		"Categories":      lo.Map(categoryIds, func(id model.UUID, _ int) *model.Category { return &model.Category{Id: id} }),
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

func (s *SqlVoucherStore) Delete(transaction *gorm.DB, ids []string) (int64, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	res := transaction.Raw("DELETE FROM "+model.VoucherTableName+" WHERE Id IN ?", ids)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
