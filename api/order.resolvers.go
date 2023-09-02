package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

// NOTE: order events are sorted by CreateAt
// NOTE: please refer to ./schemas/order.graphqls for details on directives used.
func (r *Resolver) HomepageEvents(ctx context.Context, args GraphqlParams) (*OrderEventCountableConnection, error) {
	paginationValues, appErr := args.Parse("HomePageEvents")
	if appErr != nil {
		return nil, appErr
	}

	orderEventFilterOpts := model.OrderEventFilterOptions{
		Conditions: squirrel.Eq{model.OrderEventTableName + ".Type": []model.OrderEventType{
			model.ORDER_EVENT_TYPE_ORDER_FULLY_PAID,
			model.ORDER_EVENT_TYPE_PLACED,
			model.ORDER_EVENT_TYPE_PLACED_FROM_DRAFT,
		}},
		GraphqlPaginationValues: *paginationValues,
		CountTotal:              true,
	}

	if orderEventFilterOpts.GraphqlPaginationValues.OrderBy == "" {
		orderDirection := args.orderDirection()
		orderFields := util.AnyArray[string]{model.OrderEventTableName + ".CreateAt"}
		orderEventFilterOpts.GraphqlPaginationValues.OrderBy = orderFields.
			Map(func(_ int, item string) string {
				return item + " " + orderDirection
			}).
			Join(",")
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	totalCount, events, appErr := embedCtx.App.Srv().OrderService().FilterOrderEventsByOptions(&orderEventFilterOpts)
	if appErr != nil {
		return nil, appErr
	}

	keyFunc := func(e *model.OrderEvent) []any { return []any{model.OrderEventTableName + ".CreateAt", e.CreateAt} }
	hasNextPage, hasPrevPage := args.checkNextPageAndPreviousPage(len(events))
	res := constructCountableConnection(events, totalCount, hasNextPage, hasPrevPage, keyFunc, SystemOrderEventToGraphqlOrderEvent)
	return (*OrderEventCountableConnection)(unsafe.Pointer(res)), nil
}

// NOTE: please refer to ./schemas/order.graphqls for details on directives used.
func (r *Resolver) OrderAddNote(ctx context.Context, args struct {
	Order string
	Input OrderAddNoteInput
}) (*OrderAddNote, error) {
	if !model.IsValidId(args.Order) {
		return nil, model.NewAppError("OrderAddNote", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderId"}, "please provide valid order id to add note", http.StatusBadRequest)
	}
	message := args.Input.Message
	message = strings.TrimSpace(message)
	if len(message) == 0 {
		return nil, model.NewAppError("OrderAddNote", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Message"}, "please provide message", http.StatusBadRequest)
	}
	message = model.SanitizeUnicode(message)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	order, appErr := embedCtx.App.Srv().OrderService().OrderById(args.Order)
	if appErr != nil {
		return nil, appErr
	}

	// begin transaction
	tran := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tran.Error != nil {
		return nil, model.NewAppError("OrderAddNote", model.ErrorCreatingTransactionErrorID, nil, tran.Error.Error(), http.StatusInternalServerError)
	}

	event, appErr := embedCtx.App.Srv().OrderService().CommonCreateOrderEvent(tran, &model.OrderEventOption{
		OrderID: order.Id,
		UserID:  &embedCtx.AppContext.Session().UserId,
		Parameters: model.StringInterface{
			"message": message,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	// commit
	err := tran.Commit().Error
	if err != nil {
		return nil, model.NewAppError("OrderAddNote", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &OrderAddNote{
		Order: SystemOrderToGraphqlOrder(order),
		Event: SystemOrderEventToGraphqlOrderEvent(event),
	}, nil
}

func (r *Resolver) OrderCancel(ctx context.Context, args struct{ Id string }) (*OrderCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderCapture(ctx context.Context, args struct {
	Amount *decimal.Decimal
	Id     string
}) (*OrderCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderConfirm(ctx context.Context, args struct{ Id string }) (*OrderConfirm, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfill(ctx context.Context, args struct {
	Input OrderFulfillInput
	Order *string
}) (*OrderFulfill, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentCancel(ctx context.Context, args struct {
	Id    string
	Input *FulfillmentCancelInput
}) (*FulfillmentCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentApprove(ctx context.Context, args struct {
	AllowStockToBeExceeded *bool
	Id                     string
	NotifyCustomer         bool
}) (*FulfillmentApprove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentUpdateTracking(ctx context.Context, args struct {
	Id    string
	Input FulfillmentUpdateTrackingInput
}) (*FulfillmentUpdateTracking, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentRefundProducts(ctx context.Context, args struct {
	Input OrderRefundProductsInput
	Order string
}) (*FulfillmentRefundProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderFulfillmentReturnProducts(ctx context.Context, args struct {
	Input OrderReturnProductsInput
	Order string
}) (*FulfillmentReturnProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderMarkAsPaid(ctx context.Context, args struct {
	Id                   string
	TransactionReference *string
}) (*OrderMarkAsPaid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderRefund(ctx context.Context, args struct {
	Amount *decimal.Decimal
	Id     string
}) (*OrderRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderUpdate(ctx context.Context, args struct {
	Id    string
	Input OrderUpdateInput
}) (*OrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderUpdateShipping(ctx context.Context, args struct {
	Order string
	Input OrderUpdateShippingInput
}) (*OrderUpdateShipping, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderVoid(ctx context.Context, args struct{ Id string }) (*OrderVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderBulkCancel(ctx context.Context, args struct{ Ids []string }) (*OrderBulkCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderSettings(ctx context.Context) (*OrderSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Order(ctx context.Context, args struct{ Id string }) (*Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Orders(ctx context.Context, args struct {
	SortBy  *OrderSortingInput
	Filter  *OrderFilterInput
	Channel *string
	GraphqlParams
}) (*OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrders(ctx context.Context, args struct {
	SortBy *OrderSortingInput
	Filter *OrderDraftFilterInput
	GraphqlParams
}) (*OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrdersTotal(ctx context.Context, args struct {
	Period  *ReportingPeriod
	Channel *string
}) (*TaxedMoney, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) OrderByToken(ctx context.Context, args struct{ Token string }) (*Order, error) {
	panic(fmt.Errorf("not implemented"))
}
