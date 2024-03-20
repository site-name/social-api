package order

import (
	"database/sql"
	"fmt"

	"github.com/mattermost/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
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
			model.PaymentWhere.IsActive.EQ(true),
			option.PaymentChargeStatus,
		)
	}

	if option.Statuses.Len() > 0 {
		statusOrConditions := model_helper.Or{
			squirrel.Eq{model.OrderTableColumns.Status: option.Statuses},
		}

		if option.Statuses.Contains(model_helper.OrderStatusFilterReadyToFulfill) {
			conds = append(
				conds,
				qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Payments, model.PaymentTableColumns.OrderID, model.OrderTableColumns.ID)),
				qm.Select(fmt.Sprintf("SUM ( %s ) AS AmountPaid", model.PaymentTableColumns.CapturedAmount)),
				qm.GroupBy(model.OrderTableColumns.ID),
			)

			statusOrConditions = append(
				statusOrConditions,
				model_helper.And{
					squirrel.Eq{model.OrderTableColumns.Status: []model.OrderStatus{model.OrderStatusUnfulfilled, model.OrderStatusPartiallyFulfilled}},
					squirrel.Expr(fmt.Sprintf("%s <= AmountPaid", model.OrderTableColumns.TotalGrossAmount)),
					squirrel.Expr(
						fmt.Sprintf(
							`EXISTS (
								SELECT (1) AS "a"
								FROM %s
								WHERE
									%s
									AND %s = %s
								LIMIT 1
							)`,
							model.TableNames.Payments,
							model.PaymentTableColumns.IsActive,
							model.PaymentTableColumns.OrderID,
							model.OrderTableColumns.ID,
						),
					),
				},
			)
		}

		if option.Statuses.Contains(model_helper.OrderStatusFilterReadyToCapture) {
			statusOrConditions = append(
				statusOrConditions,
				model_helper.And{
					squirrel.NotEq{model.OrderTableColumns.Status: []model.OrderStatus{model.OrderStatusDraft, model.OrderStatusCanceled}},
					// model.OrderWhere.Status.NOT_IN([]model.OrderStatus{model.ORDER_STATUS_DRAFT, model.ORDER_STATUS_CANCELED}),
					squirrel.Expr(
						fmt.Sprintf(
							`EXISTS (
								SELECT (1) AS "a"
								FROM %s
								WHERE
									%s
									AND %s = ?
									AND %s = %s
								LIMIT 1
							)`,
							model.TableNames.Payments,
							model.PaymentTableColumns.IsActive,
							model.PaymentTableColumns.ChargeStatus,
							model.PaymentTableColumns.OrderID,
							model.OrderTableColumns.ID,
						),
						model.PaymentChargeStatusNotCharged,
					),
				},
			)
		}

		conds = append(conds, &statusOrConditions)
	}

	if option.AnnotateBillingAddressNames {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Addresses, model.AddressTableColumns.ID, model.OrderTableColumns.BillingAddressID)),
			qm.Select(fmt.Sprintf(`%s AS "%s"`, model.AddressTableColumns.LastName, model_helper.CustomOrderTableColumns.OrderBillingAddressLastName)),
			qm.Select(fmt.Sprintf(`%s AS "%s"`, model.AddressTableColumns.FirstName, model_helper.CustomOrderTableColumns.OrderBillingAddressFirstName)),
		)
	}

	if option.AnnotateLastPaymentChargeStatus {
		conds = append(
			conds,
			qm.Select(fmt.Sprintf(
				`(
					SELECT
						%s
					FROM %s
					WHERE
						%s = %s
					ORDER BY
						%s DESC
					LIMIT 1
				) AS "%s"`,
				model.PaymentTableColumns.ChargeStatus,
				model.TableNames.Payments,
				model.PaymentTableColumns.OrderID,
				model.OrderTableColumns.ID,
				model.PaymentTableColumns.CreatedAt,
				model_helper.CustomOrderTableColumns.OrderLastPaymentChargeStatus,
			)),
		)
	}

	return conds
}

func (os *SqlOrderStore) FilterByOption(option model_helper.OrderFilterOption) (model_helper.CustomOrderSlice, error) {
	conds := os.commonQueryBuilder(option)
	rows, err := model.Orders(conds...).Query.Query(os.GetReplica())
	if err != nil {
		return nil, errors.Wrap(err, "failed to find orders with given option")
	}
	defer rows.Close()

	var result model_helper.CustomOrderSlice
	for rows.Next() {
		var order model_helper.CustomOrder
		var scanValues = model_helper.OrderScanValues(&order.Order)

		if option.AnnotateBillingAddressNames {
			scanValues = append(scanValues, &order.OrderBillingAddressLastName, &order.OrderBillingAddressFirstName)
		}
		if option.AnnotateLastPaymentChargeStatus {
			scanValues = append(scanValues, &order.OrderLastPaymentChargeStatus)
		}

		err := rows.Scan(scanValues...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of order")
		}

		result = append(result, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "order rows has error")
	}

	return result, nil
}

func (s *SqlOrderStore) Delete(transaction boil.ContextTransactor, ids []string) (int64, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	return model.Orders(model.OrderWhere.ID.IN(ids)).DeleteAll(transaction)
}
