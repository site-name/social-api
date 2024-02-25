package discount

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlDiscountSaleStore struct {
	store.Store
}

func NewSqlDiscountSaleStore(sqlStore store.Store) store.DiscountSaleStore {
	return &SqlDiscountSaleStore{sqlStore}
}

func (ss *SqlDiscountSaleStore) Upsert(transaction boil.ContextTransactor, sale model.Sale) (*model.Sale, error) {
	if transaction == nil {
		transaction = ss.GetMaster()
	}

	var isSaving bool
	if sale.ID == "" {
		isSaving = true
		model_helper.SalePreSave(&sale)
	} else {
		model_helper.SalePreUpdate(&sale)
	}

	if err := model_helper.SaleIsValid(sale); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = sale.Insert(transaction, boil.Infer())
	} else {
		_, err = sale.Update(transaction, boil.Blacklist(model.SaleColumns.CreatedAt))
	}

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

func (ss *SqlDiscountSaleStore) Get(saleID string) (*model.Sale, error) {
	sale, err := model.FindSale(ss.GetReplica(), saleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Sales, saleID)
		}
		return nil, err
	}
	return sale, nil
}

func (ss *SqlDiscountSaleStore) commonQueryOptionsBuilder(option model_helper.SaleFilterOption) ([]qm.QueryMod, *model_helper.AppError) {
	appErr := option.Validate()
	if appErr != nil {
		return nil, appErr
	}

	result := option.CommonQueryOptions.Conditions
	if option.SaleChannelListing_ChannelSlug != nil {
		result = append(
			result,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.SaleChannelListings, model.SaleChannelListingTableColumns.SaleID, model.SaleTableColumns.ID)),
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.SaleChannelListingTableColumns.ChannelID)),
			option.SaleChannelListing_ChannelSlug,
		)
	}

	if option.AnnotateSaleDiscountValue {
		result = append(
			result,
			qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.SaleChannelListings, model.SaleChannelListingTableColumns.SaleID, model.SaleTableColumns.ID)),
			qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.SaleChannelListingTableColumns.ChannelID)),
			qm.Select(fmt.Sprintf("MIN (%s) FILTER (WHERE %s = '%s') AS %s", model.SaleChannelListingTableColumns.DiscountValue, model.ChannelTableColumns.Slug, option.ChannelSlug, model_helper.CustomSaleTableColumns.DiscountValue)),
		)
	}

	return result, nil
}

func (ss *SqlDiscountSaleStore) FilterSalesByOption(option model_helper.SaleFilterOption) (model_helper.CustomSaleSlice, error) {
	conds, appErr := ss.commonQueryOptionsBuilder(option)
	if appErr != nil {
		return nil, appErr
	}

	rows, err := model.Sales(conds...).Query.Query(ss.GetReplica())
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sales with given conditions")
	}

	var result model_helper.CustomSaleSlice
	for rows.Next() {
		var customSale model_helper.CustomSale
		var scanFields []any

		if option.AnnotateSaleDiscountValue {
			scanFields = model_helper.CustomSaleScanValues(&customSale)
		} else {
			scanFields = model_helper.SaleScanValues(&customSale.Sale)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of sale")
		}

		result = append(result, &customSale)
	}

	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows has error")
	}

	return result, nil
}

func (s *SqlDiscountSaleStore) Delete(transaction boil.ContextTransactor, ids []string) (int64, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	return model.Sales(model.SaleWhere.ID.IN(ids)).DeleteAll(transaction)
}
