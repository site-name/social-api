package discount

import (
	"strings"

	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlDiscountSaleStore struct {
	store.Store
}

func NewSqlDiscountSaleStore(sqlStore store.Store) store.DiscountSaleStore {
	return &SqlDiscountSaleStore{sqlStore}
}

func (s *SqlDiscountSaleStore) ScanFields(sale *model.Sale) []any {
	return []any{
		&sale.Id,
		&sale.Name,
		&sale.Type,
		&sale.StartDate,
		&sale.EndDate,
		&sale.CreateAt,
		&sale.UpdateAt,
		&sale.Metadata,
		&sale.PrivateMetadata,
	}
}

// Upsert bases on sale's Id to decide to update or insert given sale
func (ss *SqlDiscountSaleStore) Upsert(transaction *gorm.DB, sale *model.Sale) (*model.Sale, error) {
	if transaction == nil {
		transaction = ss.GetMaster()
	}

	var err error
	if sale.Id == "" {
		err = transaction.Create(sale).Error
	} else {
		sale.CreateAt = 0 // prevent update
		err = transaction.Model(sale).Updates(sale).Error
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert sale")
	}

	return sale, nil
}

// Get finds and returns a sale with given saleID
func (ss *SqlDiscountSaleStore) Get(saleID string) (*model.Sale, error) {
	var sale model.Sale
	err := ss.GetReplica().First(&sale, "Id = ?", saleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("sales", saleID)
		}
		return nil, errors.Wrap(err, "failed to find sale by id")
	}
	return &sale, nil
}

// FilterSalesByOption filter sales by option
func (ss *SqlDiscountSaleStore) FilterSalesByOption(option *model.SaleFilterOption) (int64, []*model.Sale, error) {
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

	if option.GraphqlPaginationValues.PaginationApplicable() {
		query = query.
			Where(option.GraphqlPaginationValues.Condition).
			OrderBy(option.GraphqlPaginationValues.OrderBy)
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

	// NOTICE:
	// we add limit to the query after counting.
	if option.GraphqlPaginationValues.Limit > 0 {
		query = query.Limit(option.GraphqlPaginationValues.Limit)
	}

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

func (s *SqlDiscountSaleStore) Delete(transaction *gorm.DB, options *model.SaleFilterOption) (int64, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	result := transaction.Raw("DELETE FROM "+model.SaleTableName, store.BuildSqlizer(options.Conditions)...)
	if result.Error != nil {
		return 0, errors.Wrap(result.Error, "failed to delete sale(s) by given options")
	}

	return result.RowsAffected, nil
}

func (s *SqlDiscountSaleStore) ToggleSaleRelations(transaction *gorm.DB, sales model.Sales, collectionIds, productIds, variantIds, categoryIds []string, isDelete bool) error {
	if len(sales) == 0 {
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
		for _, sale := range sales {
			if sale != nil {
				switch {
				case isDelete:
					err := transaction.Model(sale).Association(associationName).Delete(relations)
					if err != nil {
						return errors.Wrap(err, "failed to delete sale "+strings.ToLower(associationName)+" relations")
					}
				default:
					err := transaction.Model(sale).Association(associationName).Append(relations)
					if err != nil {
						return errors.Wrap(err, "failed to insert sale "+strings.ToLower(associationName)+" relations")
					}
				}
			}
		}
	}

	return nil
}
