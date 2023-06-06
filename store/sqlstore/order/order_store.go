package order

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
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
func (os *SqlOrderStore) BulkUpsert(transaction store_iface.SqlxTxExecutor, orders []*model.Order) ([]*model.Order, error) {
	var (
		saveQuery   = "INSERT INTO " + store.OrderTableName + "(" + os.ModelFields("").Join(",") + ") VALUES (" + os.ModelFields(":").Join(",") + ")"
		updateQuery = "UPDATE " + store.OrderTableName + " SET " +
			os.ModelFields("").
				Map(func(_ int, s string) string {
					return s + "=:" + s
				}).
				Join(",") + " WHERE Id=:Id"
		runner = os.GetMasterX()
	)
	if transaction != nil {
		runner = transaction
	}

	for _, ord := range orders {
		var (
			err      error
			isSaving bool
		)

		if !model.IsValidId(ord.Id) {
			ord.Id = ""
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
				_, err = runner.NamedExec(saveQuery, ord)
				if err != nil {
					if os.IsUniqueConstraintError(err, []string{"Token", "orders_token_key"}) {
						ord.NewToken()
						continue
					}
					break // system error, break right now to return
				}
				break // no error, break
			}

		} else {
			var oldOrder model.Order
			// try finding if order exist
			err = runner.Get(&oldOrder, "SELECT * FROM "+store.OrderTableName+" WHERE Id = ?", ord.Id)
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

			_, err = runner.NamedExec(updateQuery, ord)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to upsert order with id=%s", ord.Id)
		}
	}

	return orders, nil
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

// FilterByOption returns a list of orders, filtered by given option
func (os *SqlOrderStore) FilterByOption(option *model.OrderFilterOption) ([]*model.Order, error) {
	query := os.GetQueryBuilder().
		Select(os.ModelFields(store.OrderTableName + ".")...).
		From(store.OrderTableName)

		// parse options:
	for _, cond := range []squirrel.Sqlizer{
		option.Id,
		option.Status,
		option.CheckoutToken,
		option.UserEmail,
		option.UserID,
		option.ChannelID,
		option.ShippingMethodID,
	} {
		if cond != nil {
			query = query.Where(cond)
		}
	}

	if option.ChannelSlug != nil {
		query = query.
			InnerJoin(store.ChannelTableName + " ON (Channels.Id = Orders.ChannelID)").
			Where(option.ChannelSlug)
	}
	if option.SelectForUpdate {
		query = query.Suffix("FOR UPDATE")
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
