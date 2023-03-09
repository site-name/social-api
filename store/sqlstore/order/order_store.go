package order

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlOrderStore struct {
	store.Store
}

func NewSqlOrderStore(sqlStore store.Store) store.OrderStore {
	return &SqlOrderStore{sqlStore}
}

func (os *SqlOrderStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"CreateAt",
		"Status",
		"UserID",
		"ShopID",
		"LanguageCode",
		"TrackingClientID",
		"BillingAddressID",
		"ShippingAddressID",
		"UserEmail",
		"OriginalID",
		"Origin",
		"Currency",
		"ShippingMethodID",
		"CollectionPointID",
		"ShippingMethodName",
		"CollectionPointName",
		"ChannelID",
		"ShippingPriceNetAmount",
		"ShippingPriceGrossAmount",
		"ShippingTaxRate",
		"Token",
		"CheckoutToken",
		"TotalNetAmount",
		"UnDiscountedTotalNetAmount",
		"TotalGrossAmount",
		"UnDiscountedTotalGrossAmount",
		"TotalPaidAmount",
		"VoucherID",
		"DisplayGrossPrices",
		"CustomerNote",
		"WeightAmount",
		"WeightUnit",
		"Weight",
		"RedirectUrl",
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

func (os *SqlOrderStore) ScanFields(holder *model.Order) []interface{} {
	return []interface{}{
		&holder.Id,
		&holder.CreateAt,
		&holder.Status,
		&holder.UserID,
		&holder.ShopID,
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
func (os *SqlOrderStore) BulkUpsert(orders []*model.Order) ([]*model.Order, error) {
	var (
		saveQuery   = "INSERT INTO " + store.OrderTableName + "(" + os.ModelFields("").Join(",") + ") VALUES (" + os.ModelFields(":").Join(",") + ")"
		updateQuery = "UPDATE " + store.OrderTableName + " SET " +
			os.ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id=:Id"
	)

	for _, ord := range orders {
		var (
			err        error
			numUpdated int64
			isSaving   bool
		)

		if ord.Id == "" {
			isSaving = true
			ord.PreSave()
		} else {
			ord.PreUpdate()
		}

		if err := ord.IsValid(); err != nil {
			return nil, err
		}

		if isSaving {
			for {
				_, err = os.GetMasterX().NamedExec(saveQuery, ord)
				if err != nil {
					if os.IsUniqueConstraintError(err, []string{"Token", "orders_token_key"}) {
						ord.NewToken()
						continue
					}
					break
				}
				break
			}

		} else {
			var oldOrder model.Order
			// try finding if order exist
			err = os.GetReplicaX().Get(&oldOrder, "SELECT * FROM "+store.OrderTableName+" WHERE Id = ?", ord.Id)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, err
				}
				return nil, errors.Wrapf(err, "failed to find order with id=%s", ord.Id)
			}

			// set all NOT editable fields for newOrder:
			// NOTE: order's Token can be updated too
			ord.CreateAt = oldOrder.CreateAt
			ord.TrackingClientID = oldOrder.TrackingClientID
			ord.BillingAddressID = oldOrder.BillingAddressID
			ord.ShippingAddressID = oldOrder.ShippingAddressID
			ord.CollectionPointName = oldOrder.CollectionPointName
			ord.ShippingMethodName = oldOrder.ShippingMethodName
			ord.ShippingPriceNetAmount = oldOrder.ShippingPriceNetAmount
			ord.ShippingPriceGrossAmount = oldOrder.ShippingPriceGrossAmount

			var result sql.Result
			result, err = os.GetMasterX().NamedExec(updateQuery, ord)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert order with id=%s", ord.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("multiple orders were updated: %d instead of 1", numUpdated)
		}
	}

	return orders, nil
}

func (os *SqlOrderStore) Save(transaction store_iface.SqlxTxExecutor, order *model.Order) (*model.Order, error) {
	var executor store_iface.SqlxExecutor = os.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	order.PreSave()
	if err := order.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.OrderTableName + "(" + os.ModelFields("").Join(",") + ") VALUES (" + os.ModelFields(":").Join(",") + ")"
	for {
		if _, err := executor.NamedExec(query, order); err != nil {
			if os.IsUniqueConstraintError(err, []string{"Token", "orders_token_key", "idx_orders_token_unique"}) {
				order.NewToken()
				continue
			}
			return nil, errors.Wrapf(err, "failed to save order with Id=%s", order.Id)
		}
		break
	}
	return order, nil
}

// Get finds and returns 1 order with given id
func (os *SqlOrderStore) Get(id string) (*model.Order, error) {
	var order model.Order
	err := os.GetReplicaX().Get(&order, "SELECT * FROM "+store.OrderTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.OrderTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find order with Id=%s", id)
	}
	return &order, nil
}

func (os *SqlOrderStore) Update(transaction store_iface.SqlxTxExecutor, newOrder *model.Order) (*model.Order, error) {
	var executor store_iface.SqlxExecutor = os.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	newOrder.PreUpdate()
	if err := newOrder.IsValid(); err != nil {
		return nil, err
	}

	// check if order exist
	var oldOrder model.Order
	err := os.GetMasterX().Get(&oldOrder, newOrder.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get order with Id=%s", newOrder.Id)
	}

	// set all NOT editable fields for newOrder:
	// NOTE: order's Token can be updated too
	newOrder.CreateAt = oldOrder.CreateAt
	newOrder.TrackingClientID = oldOrder.TrackingClientID
	newOrder.BillingAddressID = oldOrder.BillingAddressID
	newOrder.ShippingAddressID = oldOrder.ShippingAddressID
	newOrder.CollectionPointName = oldOrder.CollectionPointName
	newOrder.ShippingMethodName = oldOrder.ShippingMethodName
	newOrder.ShippingPriceNetAmount = oldOrder.ShippingPriceNetAmount
	newOrder.ShippingPriceGrossAmount = oldOrder.ShippingPriceGrossAmount

	query := "UPDATE " + store.OrderTableName + " SET " + os.
		ModelFields("").
		Map(func(_ int, s string) string {
			return s + "=:" + s
		}).
		Join(",") + " WHERE Id=:Id"

	result, err := executor.NamedExec(query, newOrder)
	if err != nil {
		if os.IsUniqueConstraintError(err, []string{"Token", "orders_token_key", "idx_orders_token_unique"}) {
			// this is user's intension to update token, he/she must be notified
			return nil, store.NewErrInvalidInput(store.OrderTableName, "token", newOrder.Token)
		}
		return nil, errors.Wrapf(err, "failed to update order with id=%s", newOrder.Id)
	}

	if numberOfUpdatedOrder, _ := result.RowsAffected(); numberOfUpdatedOrder > 1 {
		return nil, fmt.Errorf("multiple orders were updated: orderId=%s, count=%d", newOrder.Id, numberOfUpdatedOrder)
	}

	return newOrder, nil
}

// FilterByOption returns a list of orders, filtered by given option
func (os *SqlOrderStore) FilterByOption(option *model.OrderFilterOption) ([]*model.Order, error) {
	query := os.GetQueryBuilder().
		Select(os.ModelFields(store.OrderTableName + ".")...).
		From(store.OrderTableName)

	// parse options:
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.Status != nil {
		query = query.Where(option.Status)
	}
	if option.CheckoutToken != nil {
		query = query.Where(option.CheckoutToken)
	}
	if option.ChannelSlug != nil {
		query = query.
			InnerJoin(store.ChannelTableName + " ON (Channels.Id = Orders.ChannelID)").
			Where(option.ChannelSlug)
	}
	if option.UserEmail != nil {
		query = query.Where(option.UserEmail)
	}
	if option.UserID != nil {
		query = query.Where(option.UserID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res model.Orders
	err = os.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find orders with given option")
	}

	return res, nil
}
