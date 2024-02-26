package order

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlOrderStore struct {
	store.Store
}

func NewSqlOrderStore(sqlStore store.Store) store.OrderStore {
	return &SqlOrderStore{sqlStore}
}

func (os *SqlOrderStore) BulkUpsert(transaction boil.ContextTransactor, orders model.OrderSlice) (model.OrderSlice, error) {
	if transaction == nil {
		transaction = os.GetMaster()
	}

	for _, order := range orders {
		if order == nil {
			continue
		}

		isSaving := order.ID == ""
		if isSaving {
			model_helper.OrderPreSave(order)
		} else {
			model_helper.OrderCommonPre(order)
		}

		if err := model_helper.OrderIsValid(*order); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = order.Insert(transaction, boil.Infer())
		} else {
			_, err = order.Update(transaction, boil.Blacklist(
				model.OrderColumns.CreatedAt,
				model.OrderColumns.TrackingClientID,
				model.OrderColumns.BillingAddressID,
				model.OrderColumns.ShippingAddressID,
				model.OrderColumns.CollectionPointName,
				model.OrderColumns.ShippingMethodName,
				model.OrderColumns.ShippingPriceNetAmount,
				model.OrderColumns.ShippingPriceGrossAmount,
			))
		}

		if err != nil {
			if os.IsUniqueConstraintError(err, []string{model.OrderColumns.Token, "orders_token_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.Orders, model.OrderColumns.Token, order.Token)
			}
			return nil, err
		}
	}

	return orders, nil
}

func (os *SqlOrderStore) Get(id string) (*model.Order, error) {
	order, err := model.FindOrder(os.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Orders, id)
		}
		return nil, err
	}

	return order, nil
}

func (os *SqlOrderStore) commonQueryBuilder(option model_helper.OrderFilterOption) []qm.QueryMod {
	conds := option.Conditions

	if len(option.Customer) > 0 {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Users, model.UserTableColumns.ID, model.OrderTableColumns.UserID)),
			model_helper.Or{
				squirrel.Expr(fmt.Sprintf("%s %% ?", model.OrderTableColumns.UserEmail), option.Customer),
				squirrel.Expr(fmt.Sprintf("%s %% ?", model.UserTableColumns.Email), option.Customer),
				squirrel.Expr(fmt.Sprintf("%s %% ?", model.UserTableColumns.FirstName), option.Customer),
				squirrel.Expr(fmt.Sprintf("%s %% ?", model.UserTableColumns.LastName), option.Customer),
			},
		)
	}

	if len(option.Search) > 0 {
		conds = append(
			conds,
			model_helper.Or{
				squirrel.Expr(fmt.Sprintf("%s %% ?", model.OrderTableColumns.UserEmail), option.Search),
				squirrel.Expr(
					fmt.Sprintf(
						`EXISTS (
							SELECT (1) AS "a"
							FROM %s
							WHERE (
								(
									%s %% ?
									OR %s %% ?
									OR %s %% ?
								)
								AND %s = %s
							)
							LIMIT 1
						)`,
						model.TableNames.Users,
						model.UserTableColumns.Email,
						model.UserTableColumns.FirstName,
						model.UserTableColumns.LastName,
						model.UserTableColumns.ID,
						model.OrderTableColumns.UserID,
					),
					option.Search,
					option.Search,
					option.Search,
				),
				squirrel.Expr(
					fmt.Sprintf(
						`EXISTS (
							SELECT (1) AS "a"
							FROM %s
							WHERE (
								%s = ?
								AND %s = %s
							)
							LIMIT 1
						)`,
						model.TableNames.Payments,
						model.PaymentTableColumns.PSPReference,
						model.PaymentTableColumns.OrderID,
						model.OrderTableColumns.ID,
					),
					option.Search,
				),
				squirrel.Expr(
					fmt.Sprintf(
						`EXISTS (
							SELECT (1) AS "a"
							FROM %s
							WHERE (
								(
									%s %% ?
									OR %s %% ?
								)
								AND %s = %s
							)
							LIMIT 1
						)`,
						model.TableNames.OrderDiscounts,
						model.OrderDiscountTableColumns.Name,
						model.OrderDiscountTableColumns.TranslatedName,
						model.OrderDiscountTableColumns.OrderID,
						model.OrderTableColumns.ID,
					),
					option.Search,
					option.Search,
				),
				squirrel.Expr(
					fmt.Sprintf(
						`EXISTS (
							SELECT (1) AS "a"
							FROM %s
							WHERE (
								%s = ?
								AND %s = %s
							)
							LIMIT 1
						)`,
						model.TableNames.OrderLines,
						model.OrderLineTableColumns.ProductSku,
						model.OrderLineTableColumns.OrderID,
						model.OrderTableColumns.ID,
					),
					option.Search,
				),
			},
		)
	}

	if option.PaymentChargeStatus != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Payments, model.PaymentTableColumns.OrderID, model.OrderTableColumns.ID)),
			model.PaymentWhere.IsActive.EQ(model_types.NewNullBool(true)),
			option.PaymentChargeStatus,
		)
	}

	if option.Statuses.Len() > 0 {
		statusOrConditions := model_helper.Or{
			squirrel.Eq{model.OrderTableColumns.Status: option.Statuses},
		}

		if option.Statuses.Contains(model.OrderStatusCanceled) {
			statusOrConditions = append(
				statusOrConditions,

				qm.LeftJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Payments, model.PaymentTableColumns.OrderID, model.OrderTableColumns.ID)),
				qm.Select(fmt.Sprintf("SUM ( %s ) AS AmountPaid", model.PaymentTableColumns.CapturedAmount)),
				qm.GroupBy(model.OrderTableColumns.ID),
				model_helper.And{
					model.OrderWhere.Status.IN([]model.OrderStatus{model.ORDER_STATUS_UNFULFILLED, model.ORDER_STATUS_PARTIALLY_FULFILLED}),
					squirrel.Expr(fmt.Sprintf("%s <= AmountPaid", model.OrderTableColumns.TotalGrossAmount)),
					squirrel.Expr(
						fmt.Sprintf(
							`EXISTS (
								SELECT (1) AS "a"
								FROM %s
								WHERE
									%s
									AND %s
								LIMIT 1
							)`,
							model.TableNames.Payments,
							model.PaymentWhere.IsActive,
							model.PaymentWhere.OrderID.EQ(model.OrderTableColumns.ID),
						),
					),
				},
			)
		}

		// if option.Statuses.Contains(model.OrderStatusFilterReadyToCapture) {
		// 	conds = append(
		// 		conds,
		// 		model_helper.And{
		// 			model.OrderWhere.Status.NOT_IN([]model.OrderStatus{model.ORDER_STATUS_DRAFT, model.ORDER_STATUS_CANCELED}),
		// 			squirrel.Expr(
		// 				fmt.Sprintf(
		// 					`EXISTS (
		// 						SELECT (1) AS "a"
		// 						FROM %s
		// 						WHERE
		// 							%s
		// 							AND %s = ?
		// 							AND %s
		// 						LIMIT 1
		// 					)`,
		// 					model.TableNames.Payments,
		// 					model.PaymentWhere.IsActive,
		// 					model.PaymentWhere.ChargeStatus.EQ(model.PAYMENT_CHARGE_STATUS_NOT_CHARGED),
		// 					model.PaymentWhere.OrderID.EQ(model.OrderTableColumns.ID),
		// 				),
		// 				model.PAYMENT_CHARGE_STATUS_NOT_CHARGED,
		// 			),
		// 		},
		// 	)
	}

	return conds
}

func (os *SqlOrderStore) FilterByOption(option *model.OrderFilterOption) (int64, model.OrderSlice, error) {
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

func (s *SqlOrderStore) Delete(transaction boil.ContextTransactor, ids []string) (int64, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	return model.Orders(model.OrderWhere.ID.IN(ids)).DeleteAll(transaction)
}
