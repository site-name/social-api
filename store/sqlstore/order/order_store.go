package order

import (
	"database/sql"
	"fmt"

	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
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
		table.ColMap("ShopID").SetMaxSize(store.UUID_MAX_LENGTH)
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
	os.CreateIndexIfNotExists("idx_orders_user_email_lower_textpattern", store.OrderTableName, "lower(UserEmail) text_pattern_ops")
	os.CreateForeignKeyIfNotExists(store.OrderTableName, "UserID", store.UserTableName, "Id", false)
	os.CreateForeignKeyIfNotExists(store.OrderTableName, "BillingAddressID", store.AddressTableName, "Id", false)
	os.CreateForeignKeyIfNotExists(store.OrderTableName, "ShippingAddressID", store.AddressTableName, "Id", false)
	os.CreateForeignKeyIfNotExists(store.OrderTableName, "OriginalID", store.OrderTableName, "Id", false)
	os.CreateForeignKeyIfNotExists(store.OrderTableName, "ShippingMethodID", store.ShippingMethodTableName, "Id", false)
	os.CreateForeignKeyIfNotExists(store.OrderTableName, "ChannelID", store.ChannelTableName, "Id", false)
	os.CreateForeignKeyIfNotExists(store.OrderTableName, "VoucherID", store.VoucherTableName, "Id", false)
	os.CreateForeignKeyIfNotExists(store.OrderTableName, "ShopID", store.ShopTableName, "Id", false)
}

func (os *SqlOrderStore) ModelFields() []string {
	return []string{
		"Orders.Id",
		"Orders.CreateAt",
		"Orders.Status",
		"Orders.UserID",
		"Orders.ShopID",
		"Orders.LanguageCode",
		"Orders.TrackingClientID",
		"Orders.BillingAddressID",
		"Orders.ShippingAddressID",
		"Orders.UserEmail",
		"Orders.OriginalID",
		"Orders.Origin",
		"Orders.Currency",
		"Orders.ShippingMethodID",
		"Orders.ShippingMethodName",
		"Orders.ChannelID",
		"Orders.ShippingPriceNetAmount",
		"Orders.ShippingPriceGrossAmount",
		"Orders.ShippingTaxRate",
		"Orders.Token",
		"Orders.CheckoutToken",
		"Orders.TotalNetAmount",
		"Orders.UnDiscountedTotalNetAmount",
		"Orders.TotalGrossAmount",
		"Orders.UnDiscountedTotalGrossAmount",
		"Orders.TotalPaidAmount",
		"Orders.VoucherID",
		"Orders.DisplayGrossPrices",
		"Orders.CustomerNote",
		"Orders.WeightAmount",
		"Orders.WeightUnit",
		"Orders.Weight",
		"Orders.RedirectUrl",
		"Orders.Metadata",
		"Orders.PrivateMetadata",
	}
}

func (os *SqlOrderStore) ScanFields(holder order.Order) []interface{} {
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
		&holder.ShippingMethodName,
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
func (os *SqlOrderStore) BulkUpsert(orders []*order.Order) ([]*order.Order, error) {
	var (
		isSaving   bool
		err        error
		oldOrder   *order.Order
		numUpdated int64
	)

	tx, err := os.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	defer store.FinalizeTransaction(tx)

	for _, ord := range orders {
		isSaving = false // reset

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
				err = tx.Insert(ord)
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
			// try finding if order exist
			err = tx.SelectOne(&oldOrder, "SELECT * FROM "+store.OrderTableName+" WHERE Id = :ID", map[string]interface{}{"ID": ord.Id})
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
			ord.ShippingMethodName = oldOrder.ShippingMethodName
			ord.ShippingPriceNetAmount = oldOrder.ShippingPriceNetAmount
			ord.ShippingPriceGrossAmount = oldOrder.ShippingPriceGrossAmount

			numUpdated, err = tx.Update(ord)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert order with id=%s", ord.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("multiple orders were updated: %d instead of 1", numUpdated)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "transaction_commit")
	}

	return orders, nil
}

func (os *SqlOrderStore) Save(transaction *gorp.Transaction, order *order.Order) (*order.Order, error) {
	var insertFunc func(list ...interface{}) error = os.GetMaster().Insert
	if transaction != nil {
		insertFunc = transaction.Insert
	}

	order.PreSave()
	if err := order.IsValid(); err != nil {
		return nil, err
	}

	for {
		if err := insertFunc(order); err != nil {
			if os.IsUniqueConstraintError(err, []string{"Token", "orders_token_key", "idx_orders_token_unique"}) {
				order.NewToken()
				continue
			}
			return nil, errors.Wrapf(err, "failed to save order with Id=%s", order.Id)
		}
		break
	}
	order.PopulateNonDbFields()
	return order, nil
}

// Get finds and returns 1 order with given id
func (os *SqlOrderStore) Get(id string) (*order.Order, error) {
	var order order.Order
	err := os.GetReplica().SelectOne(&order, "SELECT * FROM "+store.OrderTableName+" WHERE Id = :id", map[string]interface{}{"id": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.OrderTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find order with Id=%s", id)
	}
	order.PopulateNonDbFields()
	return &order, nil
}

func (os *SqlOrderStore) Update(transaction *gorp.Transaction, newOrder *order.Order) (*order.Order, error) {
	var updateFunc func(list ...interface{}) (int64, error) = os.GetMaster().Update
	if transaction != nil {
		updateFunc = transaction.Update
	}

	newOrder.PreUpdate()
	if err := newOrder.IsValid(); err != nil {
		return nil, err
	}

	oldOrderResult, err := os.GetMaster().Get(order.Order{}, newOrder.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get order with Id=%s", newOrder.Id)
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

	numberOfUpdatedOrder, err := updateFunc(newOrder)
	if err != nil {
		if os.IsUniqueConstraintError(err, []string{"Token", "orders_token_key", "idx_orders_token_unique"}) {
			// this is user's intension to update token, he/she must be notified
			return nil, store.NewErrInvalidInput(store.OrderTableName, "token", newOrder.Token)
		}
		return nil, errors.Wrapf(err, "failed to update order with id=%s", newOrder.Id)
	}

	if numberOfUpdatedOrder > 1 {
		return nil, fmt.Errorf("multiple orders were updated: orderId=%s, count=%d", newOrder.Id, numberOfUpdatedOrder)
	}

	newOrder.PopulateNonDbFields()
	return newOrder, nil
}

// FilterByOption returns a list of orders, filtered by given option
func (os *SqlOrderStore) FilterByOption(option *order.OrderFilterOption) ([]*order.Order, error) {
	query := os.GetQueryBuilder().
		Select(os.ModelFields()...).
		From(store.OrderTableName).
		OrderBy(store.TableOrderingMap[store.OrderTableName])

	// parse options:
	if option.Status != nil {
		query = query.Where(option.Status.ToSquirrel("Orders.Status"))
	}
	if option.CheckoutToken != nil {
		query = query.Where(option.CheckoutToken.ToSquirrel("Orders.CheckoutToken"))
	}
	if option.ChannelSlug != nil {
		query = query.
			InnerJoin(store.ChannelTableName + " ON (Channels.Id = Orders.ChannelID)").
			Where(option.ChannelSlug.ToSquirrel("Channels.Slug"))
	}
	if option.UserEmail != nil {
		query = query.Where(option.UserEmail.ToSquirrel("Orders.UserEmail"))
	}
	if option.UserID != nil {
		query = query.Where(option.UserID.ToSquirrel("Orders.UserID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*order.Order
	_, err = os.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find orders with given option")
	}

	return res, nil
}
