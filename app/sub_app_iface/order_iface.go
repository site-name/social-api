// Code generated by "make app-layers"
// DO NOT EDIT

package sub_app_iface

import (
	"github.com/mattermost/gorp"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/order/types"
	"github.com/sitename/sitename/exception"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/model/warehouse"
)

// OrderService contains methods for working with orders
type OrderService interface {
	// AddGiftcardsToOrder
	AddGiftcardsToOrder(transaction *gorp.Transaction, checkoutInfo *checkout.CheckoutInfo, orDer *order.Order, totalPriceLeft *goprices.Money, user *account.User, _ interface{}) *model.AppError
	// AddVariantToOrder Add total_quantity of variant to order.
	//
	// Returns an order line the variant was added to.
	AddVariantToOrder(orDer *order.Order, variant *product_and_discount.ProductVariant, quantity int, user *account.User, _ interface{}, manager interface{}, discounts []*product_and_discount.DiscountInfo, allocateStock bool) (*order.OrderLine, *exception.InsufficientStock, *model.AppError)
	// AllDigitalOrderLinesOfOrder finds all order lines belong to given order, and are digital products
	AllDigitalOrderLinesOfOrder(orderID string) ([]*order.OrderLine, *model.AppError)
	// AnAddressOfOrder returns shipping address of given order if presents
	AnAddressOfOrder(orderID string, whichAddressID order.WhichOrderAddressID) (*account.Address, *model.AppError)
	// ApplyDiscountToValue Calculate the price based on the provided values
	ApplyDiscountToValue(value *decimal.Decimal, valueType string, currency string, priceToDiscount interface{}) (interface{}, error)
	// AutomaticallyFulfillDigitalLines
	// Fulfill all digital lines which have enabled automatic fulfillment setting. Send confirmation email afterward.
	AutomaticallyFulfillDigitalLines(ord *order.Order, manager interface{}) *model.AppError
	// BulkUpsertFulfillmentLines performs bulk upsert given fulfillment lines and returns them
	BulkUpsertFulfillmentLines(transaction *gorp.Transaction, fulfillmentLines []*order.FulfillmentLine) ([]*order.FulfillmentLine, *model.AppError)
	// BulkUpsertOrderLines perform bulk upsert given order lines
	BulkUpsertOrderLines(transaction *gorp.Transaction, orderLines []*order.OrderLine) ([]*order.OrderLine, *model.AppError)
	// BulkUpsertOrders performs bulk upsert given orders
	BulkUpsertOrders(orders []*order.Order) ([]*order.Order, *model.AppError)
	// Calculate discount value depending on voucher and discount types.
	//
	// Raise NotApplicable if voucher of given type cannot be applied.
	GetVoucherDiscountForOrder(ord *order.Order) (result interface{}, notApplicableErr *product_and_discount.NotApplicable, appErr *model.AppError)
	// CanMarkOrderAsPaid checks if given order can be marked as paid.
	CanMarkOrderAsPaid(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError)
	// CancelFulfillment Return products to corresponding stocks.
	CancelFulfillment(fulfillment *order.Fulfillment, user *account.User, _ interface{}, warehouse *warehouse.WareHouse, manager interface{}) (*order.Fulfillment, *model.AppError)
	// CancelOrder Release allocation of unfulfilled order items.
	CancelOrder(orDer *order.Order, user *account.User, _, manager interface{}) *model.AppError
	// CancelWaitingFulfillment cancels fulfillments which is in waiting for approval state.
	CancelWaitingFulfillment(fulfillment *order.Fulfillment, user *account.User, _ interface{}, manager interface{}) *model.AppError
	// ChangeOrderLineQuantity Change the quantity of ordered items in a order line.
	//
	// NOTE: userID can be empty
	ChangeOrderLineQuantity(transaction *gorp.Transaction, userID string, _ interface{}, lineInfo *order.OrderLineData, oldQuantity int, newQuantity int, channelSlug string, manager interface{}, sendEvent bool) (*exception.InsufficientStock, *model.AppError)
	// CleanMarkOrderAsPaid Check if an order can be marked as paid.
	CleanMarkOrderAsPaid(ord *order.Order) (*payment.PaymentError, *model.AppError)
	// CommonCreateOrderEvent is common method for creating desired order event instance
	CommonCreateOrderEvent(transaction *gorp.Transaction, option *order.OrderEventOption) (*order.OrderEvent, *model.AppError)
	// CreateGiftcardsWhenApprovingFulfillment
	CreateGiftcardsWhenApprovingFulfillment(orDer *order.Order, linesData []*order.OrderLineData, user *account.User, _ interface{}, manager interface{}, settings *shop.Shop) *model.AppError
	// CreateOrderDiscountForOrder Add new order discount and update the prices
	CreateOrderDiscountForOrder(transaction *gorp.Transaction, ord *order.Order, reason string, valueType string, value *decimal.Decimal) (*product_and_discount.OrderDiscount, *model.AppError)
	// CreateReplaceOrder Create draft order with lines to replace
	CreateReplaceOrder(user *account.User, _ interface{}, originalOrder *order.Order, orderLinesToReplace []*order.OrderLineData, fulfillmentLinesToReplace []*order.FulfillmentLineData) (*order.Order, *model.AppError)
	// CustomerEmail try finding order's owner's email. If order has no user or error occured during the finding process, returns order's UserEmail property instead
	CustomerEmail(ord *order.Order) (string, *model.AppError)
	// DeleteFulfillmentLinesByOption tells store to delete fulfillment lines filtered by given option
	DeleteFulfillmentLinesByOption(transaction *gorp.Transaction, option *order.FulfillmentLineFilterOption) *model.AppError
	// DeleteFulfillmentsByOption tells store to delete fulfillments that satisfy given option
	DeleteFulfillmentsByOption(transaction *gorp.Transaction, options *order.FulfillmentFilterOption) *model.AppError
	// DeleteOrderLine Delete an order line from an order.
	DeleteOrderLine(lineInfo *order.OrderLineData, manager interface{}) (*exception.InsufficientStock, *model.AppError)
	// DeleteOrderLines perform bulk delete given order lines
	DeleteOrderLines(orderLineIDs []string) *model.AppError
	// FilterOrdersByOptions is common method for filtering orders by given option
	FilterOrdersByOptions(option *order.OrderFilterOption) ([]*order.Order, *model.AppError)
	// Fulfill order.
	//
	//     Function create fulfillments with lines.
	//     Next updates Order based on created fulfillments.
	//
	//     Args:
	//         requester (User): Requester who trigger this action.
	//         order (Order): Order to fulfill
	//         fulfillment_lines_for_warehouses (Dict): Dict with information from which
	//             system create fulfillments. Example:
	//                 {
	//                     (Warehouse.pk): [
	//                         {
	//                             "order_line": (OrderLine),
	//                             "quantity": (int),
	//                         },
	//                         ...
	//                     ]
	//                 }
	//         manager (PluginsManager): Base manager for handling plugins logic.
	//         notify_customer (bool): If `True` system send email about
	//             fulfillments to customer.
	//
	//     Return:
	//         List[Fulfillment]: Fulfillmet with lines created for this order
	//             based on information form `fulfillment_lines_for_warehouses`
	//
	//
	//     Raise:
	//         InsufficientStock: If system hasn't containt enough item in stock for any line.
	CreateFulfillments(user *account.User, _ interface{}, orDer *order.Order, fulfillmentLinesForWarehouses map[string][]*order.QuantityOrderLine, manager interface{}, notifyCustomer bool, approved bool, allowStockTobeExceeded bool) ([]*order.Fulfillment, *exception.InsufficientStock, *model.AppError)
	// FulfillOrderLines Fulfill order line with given quantity
	FulfillOrderLines(orderLineInfos []*order.OrderLineData, manager interface{}, allowStockTobeExceeded bool) (*exception.InsufficientStock, *model.AppError)
	// FulfillmentByOption returns 1 fulfillment filtered using given options
	FulfillmentByOption(transaction *gorp.Transaction, option *order.FulfillmentFilterOption) (*order.Fulfillment, *model.AppError)
	// FulfillmentLinesByOption returns all fulfillment lines by option
	FulfillmentLinesByOption(option *order.FulfillmentLineFilterOption) (order.FulfillmentLines, *model.AppError)
	// FulfillmentTrackingUpdated
	FulfillmentTrackingUpdated(fulfillment *order.Fulfillment, user *account.User, _ interface{}, trackingNumber string, manager interface{}) *model.AppError
	// FulfillmentsByOption returns a list of fulfillments be given options
	FulfillmentsByOption(transaction *gorp.Transaction, option *order.FulfillmentFilterOption) (order.Fulfillments, *model.AppError)
	// Get prices of variants belonging to the discounted specific products.
	//
	// Specific products are products, collections and categories.
	// Product must be assigned directly to the discounted category, assigning
	// product to child category won't work
	GetPricesOfDiscountedSpecificProduct(orderLines []*order.OrderLine, voucher *product_and_discount.Voucher) ([]*goprices.Money, *model.AppError)
	// GetDiscountedLines returns a list of discounted order lines, filterd from given orderLines
	GetDiscountedLines(orderLines []*order.OrderLine, voucher *product_and_discount.Voucher) ([]*order.OrderLine, *model.AppError)
	// GetOrCreateFulfillment take a filtering option, trys finding a fulfillment with given option.
	// If a fulfillment found, returns it. Otherwise, creates a new one then returns it.
	GetOrCreateFulfillment(transaction *gorp.Transaction, option *order.FulfillmentFilterOption) (*order.Fulfillment, *model.AppError)
	// GetOrderCountry Return country to which order will be shipped
	GetOrderCountry(ord *order.Order) (string, *model.AppError)
	// GetOrderDiscounts Return all discounts applied to the order by staff user
	GetOrderDiscounts(ord *order.Order) ([]*product_and_discount.OrderDiscount, *model.AppError)
	// GetProductsVoucherDiscountForOrder Calculate products discount value for a voucher, depending on its type.
	GetProductsVoucherDiscountForOrder(ord *order.Order) (*goprices.Money, *model.AppError)
	// GetTotalOrderDiscount Return total order discount assigned to the order
	GetTotalOrderDiscount(ord *order.Order) (*goprices.Money, *model.AppError)
	// GetValidShippingMethodsForOrder returns a list of valid shipping methods for given order
	GetValidShippingMethodsForOrder(ord *order.Order) ([]*shipping.ShippingMethod, *model.AppError)
	// HandleFullyPaidOrder
	//
	// user can be nil
	HandleFullyPaidOrder(manager interface{}, orDer *order.Order, user *account.User, _ interface{}) *model.AppError
	// Mark order as paid.
	//
	// Allows to create a payment for an order without actually performing any
	// payment by the gateway.
	//
	// externalReference can be empty
	MarkOrderAsPaid(orDer *order.Order, requestUser *account.User, _ interface{}, manager interface{}, externalReference string) (*payment.PaymentError, *model.AppError)
	// OrderAuthorized
	OrderAuthorized(ord *order.Order, user *account.User, _ interface{}, amount *decimal.Decimal, payMent *payment.Payment, manager interface{}) *model.AppError
	// OrderAwaitsFulfillmentApproval
	OrderAwaitsFulfillmentApproval(fulfillments []*order.Fulfillment, user *account.User, _ interface{}, fulfillmentLines order.FulfillmentLines, mnager interface{}, notifyCustomer bool) *model.AppError
	// OrderById retuns an order with given id
	OrderById(id string) (*order.Order, *model.AppError)
	// OrderCanCalcel checks if given order can be canceled
	OrderCanCancel(ord *order.Order) (bool, *model.AppError)
	// OrderCanCapture
	OrderCanCapture(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)
	// OrderCanRefund checks if order can refund
	OrderCanRefund(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)
	// OrderCanVoid
	OrderCanVoid(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)
	// OrderCaptured
	OrderCaptured(ord *order.Order, user *account.User, _ interface{}, amount *decimal.Decimal, payMent *payment.Payment, manager interface{}) *model.AppError
	// OrderConfirmed Trigger event, plugin hooks and optionally confirmation email.
	OrderConfirmed(ord *order.Order, user *account.User, _ interface{}, manager interface{}, sendConfirmationEmail bool) *model.AppError
	// OrderCreated. `fromDraft` is default to false
	OrderCreated(ord *order.Order, user *account.User, _, manager interface{}, fromDraft bool) *model.AppError
	// OrderFulfilled
	OrderFulfilled(fulfillments []*order.Fulfillment, user *account.User, _ interface{}, fulfillmentLines []*order.FulfillmentLine, manager interface{}, notifyCustomer bool) *model.AppError
	// OrderIsCaptured checks if given order is captured
	OrderIsCaptured(orderID string) (bool, *model.AppError)
	// OrderIsPreAuthorized checks if order is pre-authorized
	OrderIsPreAuthorized(orderID string) (bool, *model.AppError)
	// OrderLineById returns an order line byt given orderLineID
	OrderLineById(orderLineID string) (*order.OrderLine, *model.AppError)
	// OrderLineIsDigital Check if a variant is digital and contains digital content.
	OrderLineIsDigital(orderLine *order.OrderLine) (bool, *model.AppError)
	// OrderLineNeedsAutomaticFulfillment Check if given line is digital and should be automatically fulfilled.
	//
	// NOTE: before calling this, caller can attach related data into `orderLine` so this function does not have to call the database
	OrderLineNeedsAutomaticFulfillment(orderLine *order.OrderLine, shopDigitalSettings *shop.ShopDefaultDigitalContentSettings) (bool, *model.AppError)
	// OrderLinesByOption returns a list of order lines by given option
	OrderLinesByOption(option *order.OrderLineFilterOption) ([]*order.OrderLine, *model.AppError)
	// OrderNeedsAutomaticFulfillment checks if given order has digital products which shoul be automatically fulfilled.
	OrderNeedsAutomaticFulfillment(ord *order.Order) (bool, *model.AppError)
	// OrderRefunded
	OrderRefunded(ord *order.Order, user *account.User, _ interface{}, amount *decimal.Decimal, payMent *payment.Payment, manager interface{}) *model.AppError
	// OrderReturned
	OrderReturned(transaction *gorp.Transaction, ord *order.Order, user *account.User, _ interface{}, returnedLines []*order.QuantityOrderLine) *model.AppError
	// OrderShippingIsRequired returns a boolean value indicating that given order requires shipping or not
	OrderShippingIsRequired(orderID string) (bool, *model.AppError)
	// OrderShippingUpdated
	OrderShippingUpdated(ord *order.Order, manager interface{}) *model.AppError
	// OrderSubTotal returns sum of TotalPrice of all order lines that belong to given order
	OrderSubTotal(ord *order.Order) (*goprices.TaxedMoney, *model.AppError)
	// OrderTotalAuthorized returns order's total authorized amount
	OrderTotalAuthorized(ord *order.Order) (*goprices.Money, *model.AppError)
	// OrderTotalQuantity return total quantity of given order
	OrderTotalQuantity(orderID string) (int, *model.AppError)
	// OrderVoided
	OrderVoided(ord *order.Order, user *account.User, _ interface{}, payMent *payment.Payment, manager interface{}) *model.AppError
	// Proceed with all steps required for refunding products.
	//
	// Calculate refunds for products based on the order's lines and fulfillment
	// lines.  The logic takes the list of order lines, fulfillment lines, and their
	// quantities which is used to create the refund fulfillment. The stock for
	// unfulfilled lines will be deallocated.
	//
	// NOTE: `refundShippingCosts` default to false
	CreateRefundFulfillment(requester *account.User, _ interface{}, ord *order.Order, payMent *payment.Payment, orderLinesToRefund []*order.OrderLineData, fulfillmentLinesToRefund []*order.FulfillmentLineData, manager interface{}, amount *decimal.Decimal, refundShippingCosts bool) (interface{}, *model.AppError)
	// Process the request for replacing or returning the products.
	//
	// Process the refund when the refund is set to True. The amount of refund will be
	// calculated for all lines with statuses different from refunded.  The lines which
	// are set to replace will not be included in the refund amount.
	//
	// If the amount is provided, the refund will be used for this amount.
	//
	// If refund_shipping_costs is True, the calculated refund amount will include
	// shipping costs.
	//
	// All lines with replace set to True will be used to create a new draft order, with
	// the same order details as the original order.  These lines will be moved to
	// fulfillment with status replaced. The events with relation to new order will be
	// created.
	//
	// All lines with replace set to False will be moved to fulfillment with status
	// returned/refunded_and_returned - depends on refund flag and current line status.
	// If the fulfillment line has refunded status it will be moved to
	// returned_and_refunded
	//
	// NOTE: `payMent`, `amount` , `user` are optional.
	//
	// `refund` and `refundShippingCosts` default to false.
	//
	CreateFulfillmentsForReturnedProducts(user *account.User, _ interface{}, ord *order.Order, payMent *payment.Payment, orderLineDatas []*order.OrderLineData, fulfillmentLineDatas []*order.FulfillmentLineData, manager interface{}, refund bool, amount *decimal.Decimal, refundShippingCosts bool) (*order.Fulfillment, *order.Fulfillment, *order.Order, *model.AppError)
	// ProcessReplace Create replace fulfillment and new draft order.
	//
	// Move all requested lines to fulfillment with status replaced. Based on original
	// order create the draft order with all user details, and requested lines.
	ProcessReplace(requester *account.User, ord *order.Order, orderLineDatas []*order.OrderLineData, fulfillmentLineDatas []*order.FulfillmentLineData, manager interface{}) (*order.Fulfillment, *order.Order, *model.AppError)
	// ReCalculateOrderWeight
	ReCalculateOrderWeight(transaction *gorp.Transaction, ord *order.Order) *model.AppError
	// Recalculate all order discounts assigned to order.
	//
	// It returns the list of tuples which contains order discounts where the amount has been changed.
	RecalculateOrderDiscounts(transaction *gorp.Transaction, ord *order.Order) ([][2]*product_and_discount.OrderDiscount, *model.AppError)
	// Recalculate and assign total price of order.
	//
	// Total price is a sum of items in order and order shipping price minus
	// discount amount.
	//
	// Voucher discount amount is recalculated by default. To avoid this, pass
	// update_voucher_discount argument set to False.
	//
	// NOTE: `kwargs` can be nil
	RecalculateOrder(transaction *gorp.Transaction, ord *order.Order, kwargs map[string]interface{}) *model.AppError
	// RemoveDiscountFromOrderLine Drop discount applied to order line. Restore undiscounted price
	RemoveDiscountFromOrderLine(orderLine *order.OrderLine, ord *order.Order, manager interface{}, taxIncluded bool) *model.AppError
	// RemoveOrderDiscountFromOrder Remove the order discount from order and update the prices.
	RemoveOrderDiscountFromOrder(transaction *gorp.Transaction, ord *order.Order, orderDiscount *product_and_discount.OrderDiscount) *model.AppError
	// RestockFulfillmentLines Return fulfilled products to corresponding stocks.
	//
	// Return products to stocks and update order lines quantity fulfilled values.
	RestockFulfillmentLines(transaction *gorp.Transaction, fulfillment *order.Fulfillment, warehouse *warehouse.WareHouse) (appErr *model.AppError)
	// RestockOrderLines Return ordered products to corresponding stocks
	RestockOrderLines(ord *order.Order, manager interface{}) *model.AppError
	// SendOrderConfirmation sends notification with order confirmation
	SendOrderConfirmation(orDer *order.Order, redirectURL string, manager interface{}) *model.AppError
	// SendPaymentConfirmation sends notification with the payment confirmation
	SendPaymentConfirmation(orDer *order.Order, manager interface{}) *model.AppError
	// SetGiftcardUser Set user when the gift card is used for the first time.
	SetGiftcardUser(giftCard *giftcard.GiftCard, usedByUser *account.User, usedByEmail string)
	// UpdateDiscountForOrderLine Update discount fields for order line. Apply discount to the price
	//
	// `reason`, `valueType` can be empty. `value` can be nil
	UpdateDiscountForOrderLine(orderLine *order.OrderLine, ord *order.Order, reason string, valueType string, value *decimal.Decimal, manager interface{}, taxIncluded bool) *model.AppError
	// UpdateOrderDiscountForOrder Update the order_discount for an order and recalculate the order's prices
	//
	// `reason`, `valueType` and `value` can be nil
	UpdateOrderDiscountForOrder(transaction *gorp.Transaction, ord *order.Order, orderDiscountToUpdate *product_and_discount.OrderDiscount, reason string, valueType string, value *decimal.Decimal) *model.AppError
	// UpdateOrderPrices Update prices in order with given discounts and proper taxes.
	UpdateOrderPrices(ord *order.Order, manager interface{}, taxIncluded bool) *model.AppError
	// UpdateOrderStatus Update order status depending on fulfillments
	UpdateOrderStatus(transaction *gorp.Transaction, ord *order.Order) *model.AppError
	// UpdateOrderTotalPaid update given order's total paid amount
	UpdateOrderTotalPaid(transaction *gorp.Transaction, orDer *order.Order) *model.AppError
	// UpdateVoucherDiscount Recalculate order discount amount based on order voucher
	UpdateVoucherDiscount(fun types.RecalculateOrderPricesFunc) types.RecalculateOrderPricesFunc
	// UpsertFulfillment performs some actions then save given fulfillment
	UpsertFulfillment(transaction *gorp.Transaction, fulfillment *order.Fulfillment) (*order.Fulfillment, *model.AppError)
	// UpsertOrder depends on given order's Id property to decide update/save it
	UpsertOrder(transaction *gorp.Transaction, ord *order.Order) (*order.Order, *model.AppError)
	// UpsertOrderLine depends on given orderLine's Id property to decide update order save it
	UpsertOrderLine(transaction *gorp.Transaction, orderLine *order.OrderLine) (*order.OrderLine, *model.AppError)
	ApproveFulfillment(fulfillment *order.Fulfillment, user *account.User, _ interface{}, manager interface{}, settings *shop.Shop, notifyCustomer bool, allowStockTobeExceeded bool) (*order.Fulfillment, *exception.InsufficientStock, *model.AppError)
	CreateOrderEvent(transaction *gorp.Transaction, orderLine *order.OrderLine, userID string, quantityDiff int) *model.AppError
	CreateReturnFulfillment(requester *account.User, ord *order.Order, orderLineDatas []*order.OrderLineData, fulfillmentLineDatas []*order.FulfillmentLineData, totalRefundAmount *decimal.Decimal, shippingRefundAmount *decimal.Decimal, manager interface{}) (*order.Fulfillment, *model.AppError)
	DraftOrderCreatedFromReplaceEvent(transaction *gorp.Transaction, draftOrder *order.Order, originalOrder *order.Order, user *account.User, _ interface{}, lines []*order.QuantityOrderLine) (*order.OrderEvent, *model.AppError)
	FulfillmentAwaitsApprovalEvent(transaction *gorp.Transaction, orDer *order.Order, user *account.User, _ interface{}, fulfillmentLines order.FulfillmentLines) (*order.OrderEvent, *model.AppError)
	FulfillmentCanceledEvent(transaction *gorp.Transaction, orDer *order.Order, user *account.User, _ interface{}, fulfillment *order.Fulfillment) (*order.OrderEvent, *model.AppError)
	FulfillmentFulfilledItemsEvent(transaction *gorp.Transaction, orDer *order.Order, user *account.User, _ interface{}, fulfillmentLines order.FulfillmentLines) (*order.OrderEvent, *model.AppError)
	FulfillmentReplacedEvent(transaction *gorp.Transaction, orDer *order.Order, user *account.User, _ interface{}, replacedLines []*order.QuantityOrderLine) (*order.OrderEvent, *model.AppError)
	FulfillmentTrackingUpdatedEvent(orDer *order.Order, user *account.User, _ interface{}, trackingNumber string, fulfillment *order.Fulfillment) (*order.OrderEvent, *model.AppError)
	GetVoucherDiscountAssignedToOrder(ord *order.Order) (*product_and_discount.OrderDiscount, *model.AppError)
	MatchOrdersWithNewUser(user *account.User) *model.AppError
	OrderConfirmedEvent(orDer *order.Order, user *account.User, _ interface{}) (*order.OrderEvent, *model.AppError)
	OrderCreatedEvent(orDer *order.Order, user *account.User, _ interface{}, fromDraft bool) (*order.OrderEvent, *model.AppError)
	OrderDiscountAutomaticallyUpdatedEvent(transaction *gorp.Transaction, ord *order.Order, orderDiscount *product_and_discount.OrderDiscount, oldOrderDiscount *product_and_discount.OrderDiscount) (*order.OrderEvent, *model.AppError)
	OrderDiscountEvent(transaction *gorp.Transaction, eventType order.OrderEvents, ord *order.Order, user *account.User, orderDiscount *product_and_discount.OrderDiscount, oldOrderDiscount *product_and_discount.OrderDiscount) (*order.OrderEvent, *model.AppError)
	OrderDiscountsAutomaticallyUpdatedEvent(transaction *gorp.Transaction, ord *order.Order, changedOrderDiscounts [][2]*product_and_discount.OrderDiscount) *model.AppError
	OrderLineDiscountEvent(eventType order.OrderEvents, ord *order.Order, user *account.User, line *order.OrderLine, lineBeforeUpdate *order.OrderLine) (*order.OrderEvent, *model.AppError)
	OrderManuallyMarkedAsPaidEvent(transaction *gorp.Transaction, orDer *order.Order, user *account.User, _ interface{}, transactionReference string) (*order.OrderEvent, *model.AppError)
	OrderReplacementCreated(transaction *gorp.Transaction, originalOrder *order.Order, replaceOrder *order.Order, user *account.User, _ interface{}) (*order.OrderEvent, *model.AppError)
	SendFulfillmentConfirmationToCustomer(orDer *order.Order, fulfillment *order.Fulfillment, user *account.User, _, manager interface{}) *model.AppError
	SendOrderCancelledConfirmation(orDer *order.Order, user *account.User, _, manager interface{}) *model.AppError
	SumOrderTotals(orders []*order.Order, currencyCode string) (*goprices.TaxedMoney, *model.AppError)
	UpdateGiftcardBalance(giftCard *giftcard.GiftCard, totalPriceLeft *goprices.Money) giftcard.BalanceObject
	UpdateTaxesForOrderLine(line *order.OrderLine, ord *order.Order, manager interface{}, taxIncluded bool) *model.AppError
	UpdateTaxesForOrderLines(lines []*order.OrderLine, ord *order.Order, manager interface{}, taxIncludeed bool) *model.AppError
}