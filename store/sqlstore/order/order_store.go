package order

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
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
func (os *SqlOrderStore) FilterByOption(option *model.OrderFilterOption) (int64, []*model.Order, error) {
	query := os.GetQueryBuilder().
		Select(model.OrderTableName + ".*").
		From(model.OrderTableName).
		Where(option.Conditions)

	if option.ChannelSlug != nil {
		query = query.
			InnerJoin(fmt.Sprintf(`%[1]s ON %[1]s.Id = %[2]s.ChannelID`, model.ChannelTableName, model.OrderTableName)).
			Where(option.ChannelSlug)
	}
	if option.SelectForUpdate && option.Transaction != nil {
		query = query.Suffix("FOR UPDATE")
	}
	if option.Customer != "" {
		query = query.
			InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Id = %[2]s.UserID", model.UserTableName, model.OrderTableName)).
			Where(squirrel.Or{
				squirrel.Expr(model.OrderTableName+".UserEmail % ?", option.Customer),

				squirrel.Expr(model.UserTableName+".Email % ?", option.Customer),
				squirrel.Expr(model.UserTableName+".FirstName % ?", option.Customer),
				squirrel.Expr(model.UserTableName+".LastName % ?", option.Customer),
			})
	}
	if option.Search != "" {
		orConditions := squirrel.Or{
			squirrel.Expr(model.OrderTableName+".UserEmail % ?", option.Search),
			squirrel.Expr(
				fmt.Sprintf(`EXISTS (
					SELECT (1) AS "a"
					FROM %[1]s
					WHERE (
						(
							%[1]s.Email % ?
							OR %[1]s.FirstName % ?
							OR %[1]s.LastName % ?
						)
						AND %[1]s.Id = %[2]s.UserID
					)
					LIMIT 1
				)`,
					model.UserTableName,
					model.OrderTableName,
				),
				option.Search,
				option.Search,
				option.Search,
			),
			squirrel.Expr(
				fmt.Sprintf(`EXISTS (
				SELECT (1) AS "a"
				FROM %[1]s
				WHERE (
					%[1]s.PspReference = ?
					AND %[1]s.OrderID = %[2]s.Id
				)
				LIMIT 1
			)`, model.PaymentTableName, model.OrderTableName),
				option.Search,
			),
			squirrel.Expr(
				fmt.Sprintf(`EXISTS (
					SELECT (1) AS "a"
					FROM %[1]s
					WHERE (
						(
							%[1]s.Name % ?
							OR %[1]s.TranslatedName % ?
						)
						AND %[1]s.OrderID = %[2]s.Id
					)
					LIMIT 1
				)`, model.OrderDiscountTableName, model.OrderTableName),
				option.Search,
				option.Search,
			),
			squirrel.Expr(
				fmt.Sprintf(`EXISTS (
					SELECT (1) AS "a"
					FROM %[1]s
					WHERE (
						%[1]s.ProductSku = ?
						AND %[1]s.OrderID = %[2]s.Id
					)
					LIMIT 1
				)`, model.OrderLineTableName, model.OrderTableName),
				option.Search,
			),
		}

		query = query.Where(orConditions)
	}

	if option.PaymentChargeStatus != nil {
		query = query.
			InnerJoin(fmt.Sprintf(`%[1]s ON %[1]s.OrderID = %[2]s.Id`, model.PaymentTableName, model.OrderTableName)).
			Where(squirrel.And{
				squirrel.Expr(model.PaymentTableName + ".IsActive"),
				option.PaymentChargeStatus,
			})
	}
	if len(option.Statuses) > 0 {
		statusOrConditions := squirrel.Or{}

		nativeOrderStatuses := lo.Filter(option.Statuses, func(item model.OrderFilterStatus, _ int) bool { return model.OrderStatus(item).IsValid() })
		if len(nativeOrderStatuses) > 0 {
			statusOrConditions = append(statusOrConditions, squirrel.Eq{model.OrderTableName + ".Status": nativeOrderStatuses})
		}

		if lo.Contains(option.Statuses, model.OrderStatusFilterReadyToFulfill) {
			query = query.
				LeftJoin(fmt.Sprintf(`%[1]s ON %[1]s.OrderID = %[2]s.Id`, model.PaymentTableName, model.OrderTableName)).
				Column(fmt.Sprintf(`SUM ( %[1]s.CapturedAmount ) AS AmountPaid`, model.PaymentTableName)).
				GroupBy(model.OrderTableName + ".Id")

			statusOrConditions = append(statusOrConditions, squirrel.And{
				squirrel.Eq{model.OrderTableName + ".Status": []model.OrderStatus{model.ORDER_STATUS_UNFULFILLED, model.ORDER_STATUS_PARTIALLY_FULFILLED}},
				squirrel.Expr(model.OrderTableName + ".TotalGrossAmount <= AmountPaid"),
				squirrel.Expr(fmt.Sprintf(`EXISTS (
					SELECT (1) AS "a"
					FROM %[1]s
					WHERE
						%[1]s.IsActive
						AND %[1]s.OrderID = %[2]s.Id
					LIMIT 1
				)`, model.PaymentTableName, model.OrderTableName)),
			})
		}

		if lo.Contains(option.Statuses, model.OrderStatusFilterReadyToCapture) {
			statusOrConditions = append(statusOrConditions, squirrel.And{
				squirrel.NotEq{model.OrderTableName + ".Status": []model.OrderStatus{model.ORDER_STATUS_DRAFT, model.ORDER_STATUS_CANCELED}},
				squirrel.Expr(
					fmt.Sprintf(`EXISTS (
					SELECT (1) AS "a"
					FROM %[1]s
					WHERE
						%[1]s.IsActive
						AND %[1]s.ChargeStatus = ?
						AND %[1]s.OrderID = %[2]s.Id
					LIMIT 1
				)`, model.PaymentTableName, model.OrderTableName),
					model.PAYMENT_CHARGE_STATUS_NOT_CHARGED,
				),
			})
		}

		query = query.Where(statusOrConditions)
	}

	// annotation for graphql pagination
	if option.AnnotateBillingAddressNames {
		query = query.
			InnerJoin(fmt.Sprintf("%[1]s ON %[1]s.Id = %[2]s.BillingAddressID", model.AddressTableName, model.OrderTableName)).
			Column(fmt.Sprintf(`%[1]s.LastName AS "%[2]s.BillingAddressLastName"`, model.AddressTableName, model.OrderTableName)).
			Column(fmt.Sprintf(`%[1]s.FirstName AS "%[2]s.BillingAddressFirstName"`, model.AddressTableName, model.OrderTableName))
	}
	if option.AnnotateLastPaymentChargeStatus {
		query = query.
			Column(fmt.Sprintf(`(
				SELECT
					%[1]s.ChargeStatus
				FROM %[1]s
				WHERE
					%[1]s.OrderID = %[2]s.Id
				ORDER BY
					%[1]s.CreateAt DESC
				LIMIT 1
			)
			AS "%[2]s.LastPaymentChargeStatus"`, model.PaymentTableName, model.OrderTableName))
	}

	// check count total
	var totalCount int64
	if option.CountTotal {
		countQuery, args, err := os.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "CountTotal_ToSql")
		}

		err = os.GetReplica().Raw(countQuery, args...).Scan(&totalCount).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total number of orders that satisfy given conditions")
		}
	}

	// apply pagination
	option.GraphqlPaginationValues.AddPaginationToSelectBuilderIfNeeded(&query)

	queryString, args, err := query.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	runner := os.GetReplica()
	if option.Transaction != nil {
		runner = option.Transaction
	}
	for _, preload := range option.Preload {
		runner = runner.Preload(preload)
	}

	rows, err := runner.Raw(queryString, args...).Rows()
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to find orders with given option")
	}
	defer rows.Close()

	var res model.Orders
	for rows.Next() {
		var (
			order      model.Order
			scanFields = os.ScanFields(&order)
		)
		if option.AnnotateBillingAddressNames {
			scanFields = append(scanFields, &order.BillingAddressLastName, &order.BillingAddressFirstName)
		}
		if option.AnnotateLastPaymentChargeStatus {
			scanFields = append(scanFields, &order.LastPaymentChargeStatus)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to scan a row of order")
		}

		res = append(res, &order)
	}

	return totalCount, res, nil
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
