package discount

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlVoucherStore struct {
	store.Store
}

func NewSqlDiscountVoucherStore(sqlStore store.Store) store.DiscountVoucherStore {
	return &SqlVoucherStore{sqlStore}
}

// Upsert saves or updates given voucher then returns it with an error
func (vs *SqlVoucherStore) Upsert(voucher model.Voucher) (*model.Voucher, error) {
	isSaving := false
	if voucher.ID == "" {
		isSaving = true
		model_helper.VoucherPreSave(&voucher)
	} else {
		model_helper.VoucherPreUpdate(&voucher)
	}

	if err := model_helper.VoucherIsValid(voucher); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = voucher.Insert(vs.GetMaster(), boil.Infer())
	} else {
		_, err = voucher.Update(vs.GetMaster(), boil.Blacklist(model.VoucherColumns.CreatedAt))
	}

	if err != nil {
		if vs.IsUniqueConstraintError(err, []string{model.VoucherColumns.Code, "vouchers_code_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Vouchers, model.VoucherColumns.Code, "unique")
		}
		return nil, err
	}

	return &voucher, nil
}

// Get finds a voucher with given id, then returns it with an error
func (vs *SqlVoucherStore) Get(voucherID string) (*model.Voucher, error) {
	voucher, err := model.FindVoucher(vs.GetReplica(), voucherID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Vouchers, voucherID)
		}
		return nil, err
	}

	return voucher, nil
}

func (vs *SqlVoucherStore) FilterVouchersByOption(option model_helper.VoucherFilterOption) (model_helper.CustomVoucherSlice, error) {
	appErr := option.Validate()
	if appErr != nil {
		return nil, appErr
	}

	conds := option.Conditions
	selectColumns := []string{model.TableNames.Vouchers + ".*"}

	if option.Annotate_MinDiscountValue || option.Annotate_MinSpentAmount {
		conds = append(
			conds,
			qm.LeftOuterJoin(
				fmt.Sprintf("%s ON %s = %s", model.TableNames.VoucherChannelListings, model.VoucherChannelListingTableColumns.VoucherID, model.VoucherTableColumns.ID),
			),
			qm.LeftOuterJoin(
				fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.VoucherChannelListingTableColumns.ChannelID),
			),
			qm.GroupBy(model.VoucherTableColumns.ID),
		)

		if option.Annotate_MinDiscountValue {
			minDiscountValueColumn := fmt.Sprintf(
				`MIN (%s) FILTER (WHERE %s = '%s' OR %s = '%s') AS "minDiscountValue"`,
				model.VoucherChannelListingTableColumns.DiscountValue,
				model.VoucherChannelListingTableColumns.ChannelID,
				option.ChannelIdOrSlug,
				model.ChannelTableColumns.Slug,
				option.ChannelIdOrSlug,
			)
			selectColumns = append(selectColumns, minDiscountValueColumn)
		}
		if option.Annotate_MinSpentAmount {
			minSpentAmountColumn := fmt.Sprintf(
				`MIN (%s) FILTER (WHERE %s = '%s' OR %s = '%s') AS "minSpentAmount"`,
				model.VoucherChannelListingTableColumns.MinSpendAmount,
				model.VoucherChannelListingTableColumns.ChannelID,
				option.ChannelIdOrSlug,
				model.ChannelTableColumns.Slug,
				option.ChannelIdOrSlug,
			)
			selectColumns = append(selectColumns, minSpentAmountColumn)
		}
	} else if option.Filter_channelIsActive != nil {
		conds = append(conds,
			qm.InnerJoin(
				fmt.Sprintf("%s ON %s = %s", model.TableNames.VoucherChannelListings, model.VoucherChannelListingTableColumns.VoucherID, model.VoucherTableColumns.ID),
			),
			qm.InnerJoin(
				fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.VoucherChannelListingTableColumns.ChannelID),
			),
			option.Filter_channelIsActive,
		)
	}

	conds = append(conds, qm.Select(selectColumns...))

	var result model_helper.CustomVoucherSlice
	return result, model.Vouchers(conds...).Bind(nil, vs.GetReplica(), &result)
}

// ExpiredVouchers finds and returns vouchers that are expired before given date
func (vs *SqlVoucherStore) ExpiredVouchers(date time.Time) (model.VoucherSlice, error) {
	milisecond := util.MillisFromTime(date)

	return model.Vouchers(
		model_helper.Or{
			squirrel.GtOrEq{model.VoucherColumns.Used: model.VoucherColumns.UsageLimit},
			squirrel.Lt{model.VoucherColumns.EndDate: milisecond},
		},
		model.VoucherWhere.StartDate.LT(milisecond),
	).All(vs.GetReplica())

}

func (s *SqlVoucherStore) Delete(transaction boil.ContextTransactor, ids []string) (int64, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	return model.Vouchers(model.VoucherWhere.ID.IN(ids)).DeleteAll(transaction)
}
