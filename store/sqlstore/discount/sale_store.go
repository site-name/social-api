package discount

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
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
func (ss *SqlDiscountSaleStore) FilterSalesByOption(option *model.SaleFilterOption) ([]*model.Sale, error) {
	query := ss.GetQueryBuilder().
		Select(model.SaleTableName + ".*").
		From(model.SaleTableName).
		Where(option.Conditions)

	if option.SaleChannelListing_ChannelID != nil {
		query = query.
			InnerJoin(model.SaleChannelListingTableName + " ON SaleChannelListings.SaleID = Sales.Id").
			Where(option.SaleChannelListing_ChannelID)
	}
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterSalesByOption_ToSql")
	}

	var sales model.Sales
	err = ss.GetReplica().Raw(queryString, args...).Scan(&sales).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sales with given condition.")
	}

	return sales, nil
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
