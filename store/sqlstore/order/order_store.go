package order

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlOrderStore struct {
	store.Store
}

func NewSqlOrderStore(sqlStore store.Store) store.OrderStore {
	return &SqlOrderStore{sqlStore}
}

func (os *SqlOrderStore) ScanFields(holder *model.Order) []interface{} {
	return []interface{}{
		&holder.Id,
		&holder.CreateAt,
		&holder.Status,
		&holder.UserID,
		&holder.LanguageCode,
		&holder.TrackingClientID,
		&holder.BillingAddressID,
		&holder.ShippingAddressID,
		&holder.UserEmail,
		&holder.OriginalID,
		&holder.Origin,
		&holder.Currency,
		&holder.ShippingMethodID,
		&holder.CollectionPointID,
		&holder.ShippingMethodName,
		&holder.CollectionPointName,
		&holder.ChannelID,
		&holder.ShippingPriceNetAmount,
		&holder.ShippingPriceGrossAmount,
		&holder.ShippingTaxRate,
		&holder.Token,
		&holder.CheckoutToken,
		&holder.TotalNetAmount,
		&holder.UnDiscountedTotalNetAmount,
		&holder.TotalGrossAmount,
		&holder.UnDiscountedTotalGrossAmount,
		&holder.TotalPaidAmount,
		&holder.VoucherID,
		&holder.DisplayGrossPrices,
		&holder.CustomerNote,
		&holder.WeightAmount,
		&holder.WeightUnit,
		&holder.Weight,
		&holder.RedirectUrl,
		&holder.Metadata,
		&holder.PrivateMetadata,
	}
}

// BulkUpsert performs bulk upsert given orders
func (os *SqlOrderStore) BulkUpsert(transaction *gorm.DB, orders []*model.Order) ([]*model.Order, error) {
	if transaction == nil {
		transaction = os.GetMaster()
	}

	for _, ord := range orders {
		var err error
		if ord.Id == "" {
			err = transaction.Create(ord).Error
		} else {
			// prevent update non-editable fields
			ord.CreateAt = 0
			ord.TrackingClientID = ""
			ord.BillingAddressID = nil
			ord.ShippingAddressID = nil
			ord.CollectionPointName = nil
			ord.ShippingMethodName = nil
			ord.ShippingPriceNetAmount = nil
			ord.ShippingPriceGrossAmount = nil

			err = transaction.Model(ord).Updates(ord).Error
		}

		if err != nil {
			if os.IsUniqueConstraintError(err, []string{"Token", "orders_token_key"}) {
				return nil, store.NewErrInvalidInput(model.OrderTableName, "Token", ord.Token)
			}
			return nil, errors.Wrap(err, "failed to upsert order")
		}
	}

	return orders, nil
}

// Get finds and returns 1 order with given id
func (os *SqlOrderStore) Get(id string) (*model.Order, error) {
	var order model.Order
	err := os.GetReplica().First(&order, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.OrderTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find order with Id=%s", id)
	}
	return &order, nil
}

// FilterByOption returns a list of orders, filtered by given option
func (os *SqlOrderStore) FilterByOption(option *model.OrderFilterOption) ([]*model.Order, error) {
	query := os.GetQueryBuilder().
		Select(model.OrderTableName + ".*").
		From(model.OrderTableName).
		Where(option.Conditions)

	if option.ChannelSlug != nil {
		query = query.
			InnerJoin(model.ChannelTableName + " ON (Channels.Id = Orders.ChannelID)").
			Where(option.ChannelSlug)
	}
	if option.SelectForUpdate && option.Transaction != nil {
		query = query.Suffix("FOR UPDATE")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	runner := os.GetReplica()
	if option.Transaction != nil {
		runner = option.Transaction
	}
	for _, preload := range option.Preload {
		runner = runner.Preload(preload)
	}
	var res model.Orders
	err = runner.Raw(queryString, args...).Scan(&res).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find orders with given option")
	}

	return res, nil
}

func (s *SqlOrderStore) Delete(transaction *gorm.DB, ids []string) (int64, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	result := transaction.Raw("DELETE FROM "+model.OrderTableName+" WHERE Id IN ?", ids)
	if result.Error != nil {
		return 0, errors.Wrap(result.Error, "failed to delete orders")
	}

	return result.RowsAffected, nil
}
