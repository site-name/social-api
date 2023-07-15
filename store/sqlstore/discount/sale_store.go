package discount

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlDiscountSaleStore struct {
	store.Store
}

func NewSqlDiscountSaleStore(sqlStore store.Store) store.DiscountSaleStore {
	return &SqlDiscountSaleStore{sqlStore}
}

func (s *SqlDiscountSaleStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Name",
		"Type",
		"StartDate",
		"EndDate",
		"CreateAt",
		"UpdateAt",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert bases on sale's Id to decide to update or insert given sale
func (ss *SqlDiscountSaleStore) Upsert(transaction *gorm.DB, sale *model.Sale) (*model.Sale, error) {
	if transaction == nil {
		transaction = ss.GetMaster()
	}

	result := transaction.Save(sale)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to upsert sale")
	}

	return sale, nil
}

// Get finds and returns a sale with given saleID
func (ss *SqlDiscountSaleStore) Get(saleID string) (*model.Sale, error) {
	var sale model.Sale
	result := ss.GetReplica().First(&sale, "id = ?", saleID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("sales", saleID)
		}
		return nil, errors.Wrap(result.Error, "failed to find sale by id")
	}
	return &sale, nil
}

// FilterSalesByOption filter sales by option
func (ss *SqlDiscountSaleStore) FilterSalesByOption(option *model.SaleFilterOption) ([]*model.Sale, error) {
	query := ss.
		GetQueryBuilder().
		Select(ss.ModelFields(model.SaleTableName + ".")...).
		From(model.SaleTableName)

	// check sale start date
	if option.StartDate != nil {
		query = query.Where(option.StartDate)
	}
	// check sale end date
	if option.EndDate != nil {
		query = query.Where(option.EndDate)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterSalesByOption_ToSql")
	}

	var sales []*model.Sale
	err = ss.GetReplicaX().Select(&sales, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sales with given condition.")
	}

	return sales, nil
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
		err := transaction.Model(sale).Association(association).Append(relations)
		if err != nil {
			return errors.Wrap(err, "failed to insert sale-collection relations")
		}
	}

	return nil
}
