package api

import (
	"context"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type Payment struct {
	ID                string                  `json:"id"`
	Gateway           string                  `json:"gateway"`
	IsActive          bool                    `json:"isActive"`
	Created           DateTime                `json:"created"`
	Modified          DateTime                `json:"modified"`
	Token             string                  `json:"token"`
	PaymentMethodType string                  `json:"paymentMethodType"`
	PrivateMetadata   []*MetadataItem         `json:"privateMetadata"`
	ChargeStatus      PaymentChargeStatusEnum `json:"chargeStatus"`
	Total             *Money                  `json:"total"`
	CapturedAmount    *Money                  `json:"capturedAmount"`

	p *model.Payment

	// Checkout          *Checkout               `json:"checkout"`
	// Order             *Order                  `json:"order"``
	// Metadata               []*MetadataItem         `json:"metadata"`
	// CustomerIPAddress      *string                 `json:"customerIpAddress"`
	// Actions                []*OrderAction          `json:"actions"`
	// Transactions           []*Transaction          `json:"transactions"`
	// AvailableCaptureAmount *Money                  `json:"availableCaptureAmount"`
	// AvailableRefundAmount  *Money                  `json:"availableRefundAmount"`
	// CreditCard             *CreditCard             `json:"creditCard"`
}

func SystemPaymentToGraphqlPayment(p *model.Payment) *Payment {
	if p == nil {
		return nil
	}

	return &Payment{
		ID:                p.Id,
		Gateway:           p.GateWay,
		IsActive:          *p.IsActive,
		Created:           DateTime{util.TimeFromMillis(p.CreateAt)},
		Modified:          DateTime{util.TimeFromMillis(p.UpdateAt)},
		Token:             p.Token,
		PaymentMethodType: p.PaymentMethodType,
		PrivateMetadata:   MetadataToSlice(p.PrivateMetadata),
		ChargeStatus:      PaymentChargeStatusEnum(p.ChargeStatus),
		Total:             SystemMoneyToGraphqlMoney(p.GetTotal()),
		CapturedAmount:    SystemMoneyToGraphqlMoney(p.GetCapturedAmount()),
		p:                 p,
	}
}

func (p *Payment) Checkout(ctx context.Context) (*Checkout, error) {
	if p.p.CheckoutID == nil {
		return nil, nil
	}

	checkout, err := CheckoutByTokenLoader.Load(ctx, *p.p.CheckoutID)()
	if err != nil {
		return nil, err
	}
	return SystemCheckoutToGraphqlCheckout(checkout), nil
}

func (p *Payment) Order(ctx context.Context) (*Order, error) {
	if p.p.OrderID == nil {
		return nil, nil
	}

	order, err := OrderByIdLoader.Load(ctx, *p.p.OrderID)()
	if err != nil {
		return nil, err
	}

	return SystemOrderToGraphqlOrder(order), nil
}

// requester must be owner of payment or a shop's member.
// Refer to ./schemas/payment.graphqls for details on directive used
func (p *Payment) Metadata(ctx context.Context) ([]*MetadataItem, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	requesterIsShopStaff := embedCtx.AppContext.
		Session().
		GetUserRoles().
		InterSection([]string{model.ShopStaffRoleId, model.ShopAdminRoleId}).
		Len() > 0

	if !requesterIsShopStaff {
		// TODO: checks if we need a dataloader for this
		currentUserOwnPayment, err := embedCtx.App.Srv().Store.Payment().PaymentOwnedByUser(embedCtx.AppContext.Session().UserId, p.p.Id)
		if err != nil {
			return nil, model.NewAppError("Payment.Metadata", "app.payment.checking_user_own_payment.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		if currentUserOwnPayment {
			goto returnLabel
		}
		return nil, MakeUnauthorizedError("Payment.Metadata")
	}

returnLabel:
	return MetadataToSlice(p.p.Metadata), nil
}

// NOTE: Refer to ./schemas/payment.graphqls for derective used on this method.
func (p *Payment) CustomerIPAddress(ctx context.Context) (*string, error) {
	return p.p.CustomerIpAddress, nil
}

// NOTE: Refer to ./schemas/payment.graphqls for derective used on this method.
func (p *Payment) Actions(ctx context.Context) ([]OrderAction, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	actions := []OrderAction{}
	if p.p.CanCapture() {
		actions = append(actions, OrderActionCapture)
	}
	if p.p.CanRefund() {
		actions = append(actions, OrderActionRefund)
	}

	// TODO: check if we need a dataloader for this
	canVoid, appErr := embedCtx.App.Srv().PaymentService().PaymentCanVoid(p.p)
	if appErr != nil {
		return nil, appErr
	}
	if canVoid {
		actions = append(actions, OrderActionVoid)
	}

	return actions, nil
}

// NOTE: Refer to ./schemas/payment.graphqls for derective used on this method.
func (p *Payment) Transactions(ctx context.Context) ([]*Transaction, error) {
	transactions, err := TransactionsByPaymentIdLoader.Load(ctx, p.p.Id)()
	if err != nil {
		return nil, err
	}
	return systemRecordsToGraphql(transactions, systemTransactionToGraphqlTransaction), nil
}

// NOTE: permissions/roles are checked by directives. Refer to ./schemas/payment.graphqls for details
func (p *Payment) AvailableCaptureAmount(ctx context.Context) (*Money, error) {
	if p.p.CanCapture() {
		return &Money{
			Amount:   p.p.GetChargeAmount().InexactFloat64(),
			Currency: p.p.Currency,
		}, nil
	}

	return nil, nil
}

// NOTE: permissions/roles are checked by directives. Refer to ./schemas/payment.graphqls for details
func (p *Payment) AvailableRefundAmount(ctx context.Context) (*Money, error) {
	if p.p.CanRefund() {
		return SystemMoneyToGraphqlMoney(p.p.GetCapturedAmount()), nil
	}

	return nil, nil
}

func (p *Payment) CreditCard(ctx context.Context) (*CreditCard, error) {
	// check if payment has no credit card-related information, return nil
	if p.p.CcBrand == "" &&
		(p.p.CcExpMonth == nil || *p.p.CcExpMonth == 0) &&
		(p.p.CcExpYear == nil || *p.p.CcExpYear == 0) &&
		p.p.CcFirstDigits == "" &&
		p.p.CcLastDigits == "" {
		return nil, nil
	}

	res := &CreditCard{
		Brand:       p.p.CcBrand,
		FirstDigits: &p.p.CcFirstDigits,
		LastDigits:  p.p.CcLastDigits,
	}
	if m := p.p.CcExpMonth; m != nil {
		res.ExpMonth = model.NewPrimitive(int32(*m))
	}
	if m := p.p.CcExpYear; m != nil {
		res.ExpYear = model.NewPrimitive(int32(*m))
	}

	return res, nil
}

func paymentsByOrderIdLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.Payment] {
	var (
		res        = make([]*dataloader.Result[[]*model.Payment], len(orderIDs))
		paymentMap = map[string][]*model.Payment{} // keys are order ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	payments, appErr := embedCtx.App.Srv().
		PaymentService().
		PaymentsByOption(&model.PaymentFilterOption{
			OrderID: squirrel.Eq{store.PaymentTableName + ".OrderID": orderIDs},
		})
	if appErr != nil {
		for idx := range orderIDs {
			res[idx] = &dataloader.Result[[]*model.Payment]{Error: appErr}
		}
		return res
	}

	for _, payment := range payments {
		if payment.OrderID != nil {
			paymentMap[*payment.OrderID] = append(paymentMap[*payment.OrderID], payment)
		}
	}
	for idx, id := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.Payment]{Data: paymentMap[id]}
	}
	return res
}

func paymentsByTokenLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Payment] {
	var (
		res        = make([]*dataloader.Result[*model.Payment], len(ids))
		paymentMap = map[string]*model.Payment{} // keys are payment ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	payments, appErr := embedCtx.App.Srv().PaymentService().PaymentsByOption(&model.PaymentFilterOption{
		Id: squirrel.Eq{store.PaymentTableName + ".Token": ids},
	})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.Payment]{Error: appErr}
		}
		return res
	}

	for _, pm := range payments {
		paymentMap[pm.Token] = pm
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Payment]{Data: paymentMap[id]}
	}
	return res
}

// ----------------- Transaction -----------------

type Transaction struct {
	ID              string          `json:"id"`
	Created         DateTime        `json:"created"`
	Token           string          `json:"token"`
	Kind            TransactionKind `json:"kind"`
	IsSuccess       bool            `json:"isSuccess"`
	Error           *string         `json:"error"`
	GatewayResponse JSONString      `json:"gatewayResponse"`
	Amount          *Money          `json:"amount"`

	t *model.PaymentTransaction

	// Payment         *Payment        `json:"payment"`
}

func systemTransactionToGraphqlTransaction(t *model.PaymentTransaction) *Transaction {
	if t == nil {
		return nil
	}

	return &Transaction{
		ID:              t.Id,
		Created:         DateTime{util.TimeFromMillis(t.CreateAt)},
		Token:           t.Token,
		Kind:            TransactionKind(t.Kind),
		IsSuccess:       t.IsSuccess,
		Error:           t.Error,
		GatewayResponse: JSONString(t.GatewayResponse),
		Amount:          SystemMoneyToGraphqlMoney(t.GetAmount()),
	}
}

func (t *Transaction) Payment(ctx context.Context) (*Payment, error) {
	payment, err := PaymentsByTokensLoader.Load(ctx, t.t.PaymentID)()
	if err != nil {
		return nil, err
	}
	return SystemPaymentToGraphqlPayment(payment), nil
}

func transactionsByPaymentIdLoader(ctx context.Context, paymentIDs []string) []*dataloader.Result[[]*model.PaymentTransaction] {
	var (
		res            = make([]*dataloader.Result[[]*model.PaymentTransaction], len(paymentIDs))
		transactionMap = map[string][]*model.PaymentTransaction{} // keys are payment ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	transactions, appErr := embedCtx.App.Srv().PaymentService().TransactionsByOption(&model.PaymentTransactionFilterOpts{
		PaymentID: squirrel.Eq{store.TransactionTableName + ".PaymentID": paymentIDs},
	})
	if appErr != nil {
		for idx := range paymentIDs {
			res[idx] = &dataloader.Result[[]*model.PaymentTransaction]{Error: appErr}
		}
		return res
	}

	for _, tran := range transactions {
		transactionMap[tran.PaymentID] = append(transactionMap[tran.PaymentID], tran)
	}
	for idx, id := range paymentIDs {
		res[idx] = &dataloader.Result[[]*model.PaymentTransaction]{Data: transactionMap[id]}
	}
	return res
}
