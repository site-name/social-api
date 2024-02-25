package discount

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
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

func (vs *SqlVoucherStore) Upsert(voucher model.Voucher) (*model.Voucher, error) {
	isSaving := voucher.ID == ""
	if isSaving {
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

func (vs *SqlVoucherStore) commonQueryOptionsBuilder(option model_helper.VoucherFilterOption) ([]qm.QueryMod, *model_helper.AppError) {
	appErr := option.Validate()
	if appErr != nil {
		return nil, appErr
	}

	result := option.Conditions
	if option.Annotate_MinValues {
		result = append(
			result,
			qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.VoucherChannelListings, model.VoucherChannelListingTableColumns.VoucherID, model.VoucherTableColumns.ID)),
			qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.VoucherChannelListingTableColumns.ChannelID)),
			qm.GroupBy(model.VoucherTableColumns.ID),
		)

		minDiscountValueColumn := fmt.Sprintf(
			`MIN (%s) FILTER (WHERE %s = '%s' OR %s = '%s') AS %s`,
			model.VoucherChannelListingTableColumns.DiscountValue,
			model.VoucherChannelListingTableColumns.ChannelID,
			option.ChannelIdOrSlug,
			model.ChannelTableColumns.Slug,
			option.ChannelIdOrSlug,
			model_helper.CustomVoucherTableColumns.MinDiscountValue,
		)
		result = append(result, qm.Select(minDiscountValueColumn))

		minSpentAmountColumn := fmt.Sprintf(
			`MIN (%s) FILTER (WHERE %s = '%s' OR %s = '%s') AS %s`,
			model.VoucherChannelListingTableColumns.MinSpendAmount,
			model.VoucherChannelListingTableColumns.ChannelID,
			option.ChannelIdOrSlug,
			model.ChannelTableColumns.Slug,
			option.ChannelIdOrSlug,
			model_helper.CustomVoucherTableColumns.MinSpentAmount,
		)
		result = append(result, qm.Select(minSpentAmountColumn))
	}

	return result, nil
}

func (vs *SqlVoucherStore) FilterVouchersByOption(option model_helper.VoucherFilterOption) (model_helper.CustomVoucherSlice, error) {
	queryOptions, appErr := vs.commonQueryOptionsBuilder(option)
	if appErr != nil {
		return nil, appErr
	}

	rows, err := model.Vouchers(queryOptions...).Query.Query(vs.GetReplica())
	if err != nil {
		return nil, errors.Wrap(err, "failed to find vouchers with given conditions")
	}

	var vouchers model_helper.CustomVoucherSlice
	for rows.Next() {
		var customVoucher model_helper.CustomVoucher
		var scanValues []any

		if option.Annotate_MinValues {
			scanValues = model_helper.CustomVoucherScanValues(&customVoucher)
		} else {
			scanValues = model_helper.VoucherScanValues(&customVoucher.Voucher)
		}

		err = rows.Scan(scanValues...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of voucher")
		}

		vouchers = append(vouchers, &customVoucher)
	}

	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows has error")
	}

	return vouchers, nil
}

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
