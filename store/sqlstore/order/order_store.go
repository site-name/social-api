package order

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

type SqlOrderStore struct {
	store.Store
}

func NewSqlOrderStore(sqlStore store.Store) store.OrderStore {
	os := &SqlOrderStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(order.Order{}, store.OrderTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("BillingAddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingAddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OriginalID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VoucherID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Status").SetMaxSize(order.ORDER_STATUS_MAX_LENGTH)
		table.ColMap("TrackingClientID").SetMaxSize(order.ORDER_TRACKING_CLIENT_ID_MAX_LENGTH)
		table.ColMap("Origin").SetMaxSize(order.ORDER_ORIGIN_MAX_LENGTH)
		table.ColMap("ShippingMethodName").SetMaxSize(order.ORDER_SHIPPING_METHOD_NAME_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(order.ORDER_TOKEN_MAX_LENGTH).SetUnique(true)
		table.ColMap("CheckoutToken").SetMaxSize(order.ORDER_CHECKOUT_TOKEN_MAX_LENGTH)
		table.ColMap("UserEmail").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.URL_LINK_MAX_LENGTH)
	}

	return os
}

func (os *SqlOrderStore) CreateIndexesIfNotExists() {
	os.CommonMetaDataIndex(store.OrderTableName)
	os.CreateIndexIfNotExists("idx_orders_user_email", store.OrderTableName, "UserEmail")
	os.CreateIndexIfNotExists("idx_orders_status", store.OrderTableName, "Status")
}

func (os *SqlOrderStore) Save(order *order.Order) (*order.Order, error) {
	order.PreSave()
	if err := order.IsValid(); err != nil {
		return nil, err
	}

	for {
		if err := os.GetMaster().Insert(order); err != nil {
			if os.IsUniqueConstraintError(err, []string{"Token", "orders_token_key", "idx_orders_token_unique"}) {
				order.Token = model.NewId()
				continue
			}
			return nil, errors.Wrapf(err, "failed to save order with Id=%s", order.Id)
		}
		break
	}
	order.PopulateNonDbFields()
	return order, nil
}

func (os *SqlOrderStore) Get(id string) (*order.Order, error) {
	var order order.Order

	if err := os.GetReplica().SelectOne(&order, "SELECT * FROM "+store.OrderTableName+" WHERE Id = :id", map[string]interface{}{"id": id}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.OrderTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find order with Id=%s", id)
	}
	order.PopulateNonDbFields()
	return &order, nil
}

func (os *SqlOrderStore) Update(newOrder *order.Order) (*order.Order, error) {
	if err := newOrder.IsValid(); err != nil {
		return nil, err
	}

	oldOrderResult, err := os.GetMaster().Get(order.Order{}, newOrder.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get order with Id=%s", newOrder.Id)
	}

	if oldOrderResult == nil {
		return nil, store.NewErrInvalidInput(store.OrderTableName, "id", newOrder.Id)
	}

	// set all NOT editable fields for newOrder:
	// NOTE: order's Token can be updated too
	oldOrder := oldOrderResult.(*order.Order)
	newOrder.CreateAt = oldOrder.CreateAt
	newOrder.TrackingClientID = oldOrder.TrackingClientID
	newOrder.BillingAddressID = oldOrder.BillingAddressID
	newOrder.ShippingAddressID = oldOrder.ShippingAddressID
	newOrder.ShippingMethodName = oldOrder.ShippingMethodName
	newOrder.ShippingPriceNetAmount = oldOrder.ShippingPriceNetAmount
	newOrder.ShippingPriceGrossAmount = oldOrder.ShippingPriceGrossAmount

	count, err := os.GetMaster().Update(newOrder)
	if err != nil {
		if os.IsUniqueConstraintError(err, []string{"Token", "orders_token_key", "idx_orders_token_unique"}) {
			// this is user's intension to update token, he/she must be notified
			return nil, store.NewErrInvalidInput(store.OrderTableName, "token", newOrder.Token)
		}
		return nil, errors.Wrapf(err, "failed to update order with id=%s", newOrder.Id)
	}

	if count > 1 {
		return nil, fmt.Errorf("multiple orders were updated: orderId=%s, count=%d", newOrder.Id, count)
	}

	newOrder.PopulateNonDbFields()
	return newOrder, nil
}

func (os *SqlOrderStore) UpdateTotalPaid(orderId string, newTotalPaid *decimal.Decimal) error {
	result, err := os.GetMaster().Exec("UPDATE "+store.OrderTableName+" SET TotalPaidAmount = :newTotalPaidAmount WHERE Id = :id",
		map[string]interface{}{"newTotalPaidAmount": *newTotalPaid, "id": orderId})
	if err != nil {
		return errors.Wrapf(err, "failed to update total paid amount for order with id=%s", orderId)
	}
	if rows, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "failed to fetch number of order updated")
	} else if rows > 1 {
		return fmt.Errorf("multiple orders updated, orderId=%s", orderId)
	}

	return nil
}
