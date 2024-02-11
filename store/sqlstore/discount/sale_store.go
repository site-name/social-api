package discount

import (
	"database/sql"

	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlDiscountSaleStore struct {
	store.Store
}

func NewSqlDiscountSaleStore(sqlStore store.Store) store.DiscountSaleStore {
	return &SqlDiscountSaleStore{sqlStore}
}

// Upsert bases on sale's Id to decide to update or insert given sale
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

// Get finds and returns a sale with given saleID
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

// FilterSalesByOption filter sales by option
func (ss *SqlDiscountSaleStore) FilterSalesByOption(option model_helper.SaleFilterOption) (model.SaleSlice, error) {
	query := ss.GetQueryBuilder().
		Select(model.SaleTableName + ".*").
		From(model.SaleTableName).
		Where(option.Conditions)

	if option.SaleChannelListing_ChannelSlug != nil {
		query = query.
			InnerJoin(model.SaleChannelListingTableName + " ON SaleChannelListings.SaleID = Sales.Id").
			InnerJoin(model.ChannelTableName + " ON Channels.Id = SaleChannelListings.ChannelID").
			Where(option.SaleChannelListing_ChannelSlug)
	} else if option.Annotate_Value {
		// check if channel provided:
		if !slug.IsSlug(option.ChannelSlug) {
			return 0, nil, store.NewErrInvalidInput("FilterSalesByOption", "option.ChannelSlug", option.ChannelSlug)
		}

		query = query.
			LeftJoin(model.SaleChannelListingTableName+" ON SaleChannelListings.SaleID = Sales.Id").
			LeftJoin(model.ChannelTableName+" ON Channels.Id = SaleChannelListings.ChannelID").
			Column(`MIN (
				SaleChannelListings.DiscountValue
			) FILTER (
				WHERE Channels.Slug = ?
			) AS "Sales.Value"`, option.ChannelSlug).
			GroupBy(model.SaleTableName + ".Id")
	}

	var totalSale int64
	if option.CountTotal {
		query, args, err := ss.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "FilterSalesByOptions_Count_ToSql")
		}
		err = ss.GetReplica().Raw(query, args...).Scan(&totalSale).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total number of sales by given options")
		}
	}

	option.GraphqlPaginationValues.AddPaginationToSelectBuilderIfNeeded(&query)

	queryString, args, err := query.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "FilterSalesByOption_ToSql")
	}

	rows, err := ss.GetReplica().Raw(queryString, args...).Rows()
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to find sales with given condition.")
	}
	defer rows.Close()

	var sales model.Sales
	for rows.Next() {
		var (
			sale       model.Sale
			scanFields = ss.ScanFields(&sale)
			value      decimal.Decimal
		)
		if option.Annotate_Value {
			scanFields = append(scanFields, &value)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to scan a row of sale")
		}

		if option.Annotate_Value {
			sale.Value = &value
		}
	}

	return totalSale, sales, nil
}

func (s *SqlDiscountSaleStore) Delete(transaction boil.ContextTransactor, ids []string) (int64, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	return model.Sales(model.SaleWhere.ID.IN(ids)).DeleteAll(transaction)
}
