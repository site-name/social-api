package discount

import (
	"github.com/pkg/errors"
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
	var sales []*model.Sale
	err := ss.GetReplica().Find(&sales, store.BuildSqlizer(option.Conditions)...).Error
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

func (s *SqlDiscountSaleStore) AddSaleRelations(transaction *gorm.DB, sales model.Sales, relations any) error {
	if relations == nil || len(sales) == 0 {
		return errors.New("please specify relations")
	}

	if transaction == nil {
		transaction = s.GetMaster()
	}

	var association string

	switch relations.(type) {
	case model.Products:
		association = "Products"
	case model.Categories:
		association = "Categories"
	case model.Collections:
		association = "Collections"
	case model.ProductVariants:
		association = "ProductVariants"
	default:
		return errors.New("only *model.(Product|Category|ProductVariant|Collection) types are supported")
	}

	for _, sale := range sales {
		if sale != nil && sale.Id != "" {
			err := transaction.Model(sale).Association(association).Append(relations)
			if err != nil {
				return errors.Wrap(err, "failed to insert sale-collection relations")
			}
		}
	}

	return nil
}
