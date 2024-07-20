package model_helper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// OrderFilterStatus mostly has the same values as model.OrderStatus
// but it also has some custom values that are used for filtering
type OrderFilterStatus model.OrderStatus

const (
	OrderStatusFilterReadyToFulfill     OrderFilterStatus = "ready_to_fulfill"
	OrderStatusFilterReadyToCapture     OrderFilterStatus = "ready_to_capture"
	OrderStatusFilterUnfulfilled        OrderFilterStatus = OrderFilterStatus(model.OrderStatusUnfulfilled)
	OrderStatusFilterUnconfirmed        OrderFilterStatus = OrderFilterStatus(model.OrderStatusUnconfirmed)
	OrderStatusFilterPartiallyFulfilled OrderFilterStatus = OrderFilterStatus(model.OrderStatusPartiallyFulfilled)
	OrderStatusFilterFulfilled          OrderFilterStatus = OrderFilterStatus(model.OrderStatusFulfilled)
	OrderStatusFilterCanceled           OrderFilterStatus = OrderFilterStatus(model.OrderStatusCanceled)
)

func OrderCommonPre(order *model.Order) {
	if order.Token == "" {
		order.Token = NewId()
	}
	order.CustomerNote = SanitizeUnicode(order.CustomerNote)
	if order.Currency.IsValid() != nil {
		order.Currency = DEFAULT_CURRENCY
	}
	if order.WeightAmount < 0 {
		order.WeightAmount = 0
	}
	if order.Status.IsValid() != nil {
		order.Status = model.OrderStatusUnfulfilled
	}
	if order.Currency.IsValid() != nil {
		order.Currency = DEFAULT_CURRENCY
	}
	if order.LanguageCode.IsValid() != nil {
		order.LanguageCode = DEFAULT_LOCALE
	}
	order.WeightUnit = strings.ToLower(order.WeightUnit)
}

func OrderPreSave(o *model.Order) {
	OrderCommonPre(o)
	if o.ID == "" {
		o.ID = NewId()
	}
	if o.CreatedAt == 0 {
		o.CreatedAt = GetMillis()
	}

	// NOTE: those 2 fields are not placed inside commonPre because they are not editable
	if !o.ShippingMethodName.IsNil() {
		o.ShippingMethodName.String = GetPointerOfValue(SanitizeUnicode(*o.ShippingMethodName.String))
	}
	if !o.CollectionPointName.IsNil() {
		o.CollectionPointName.String = GetPointerOfValue(SanitizeUnicode(*o.CollectionPointName.String))
	}
}

func OrderIsValid(o model.Order) *AppError {
	if !IsValidId(o.ID) {
		return NewAppError("OrderIsValid", "model.order.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(o.Token) {
		return NewAppError("OrderIsValid", "model.order.is_valid.token.app_error", nil, "", http.StatusBadRequest)
	}
	if o.Currency.IsValid() != nil {
		return NewAppError("OrderIsValid", "model.order.is_valid.currency.app_error", nil, "", http.StatusBadRequest)
	}
	if o.TotalNetAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.total_net_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.TotalGrossAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.total_gross_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.UndiscountedTotalNetAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.undiscounted_total_net_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.UndiscountedTotalGrossAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.undiscounted_total_gross_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.WeightAmount < 0 {
		return NewAppError("OrderIsValid", "model.order.is_valid.weight_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.ShippingPriceNetAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.shipping_price_net_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.ShippingPriceGrossAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderIsValid", "model.order.is_valid.shipping_price_gross_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.CreatedAt <= 0 {
		return NewAppError("OrderIsValid", "model.order.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if o.LanguageCode.IsValid() != nil {
		return NewAppError("OrderIsValid", "model.order.is_valid.language_code.app_error", nil, "", http.StatusBadRequest)
	}
	if o.Status.IsValid() != nil {
		return NewAppError("OrderIsValid", "model.order.is_valid.status.app_error", nil, "", http.StatusBadRequest)
	}
	if o.Currency.IsValid() != nil {
		return NewAppError("OrderIsValid", "model.order.is_valid.currency.app_error", nil, "", http.StatusBadRequest)
	}
	if o.Origin.Valid && o.Origin.Val.IsValid() != nil {
		return NewAppError("OrderIsValid", "model.order.is_valid.origin.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.ShippingAddressID.IsNil() && !IsValidId(*o.ShippingAddressID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.shipping_address_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.BillingAddressID.IsNil() && !IsValidId(*o.BillingAddressID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.billing_address_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.UserID.IsNil() && !IsValidId(*o.UserID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.VoucherID.IsNil() && !IsValidId(*o.VoucherID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.voucher_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.OriginalID.IsNil() && !IsValidId(*o.OriginalID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.original_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.ShippingMethodID.IsNil() && !IsValidId(*o.ShippingMethodID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.shipping_method_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.CollectionPointID.IsNil() && !IsValidId(*o.CollectionPointID.String) {
		return NewAppError("OrderIsValid", "model.order.is_valid.collection_point_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(o.CheckoutToken) {
		return NewAppError("OrderIsValid", "model.order.is_valid.checkout_token.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(o.ChannelID) {
		return NewAppError("OrderIsValid", "model.order.is_valid.channel_id.app_error", nil, "", http.StatusBadRequest)
	}
	_, ok := measurement.WEIGHT_UNIT_STRINGS[measurement.WeightUnit(o.WeightUnit)]
	if !ok {
		return NewAppError("OrderIsValid", "model.order.is_valid.weight_unit.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func OrderGetShippingPrice(order model.Order) goprices.TaxedMoney {
	shippingPriceNet, _ := goprices.NewMoneyFromDecimal(order.ShippingPriceNetAmount, string(order.Currency))
	shippingPriceGross, _ := goprices.NewMoneyFromDecimal(order.ShippingPriceGrossAmount, string(order.Currency))

	taxedMoney, _ := goprices.NewTaxedMoney(*shippingPriceNet, *shippingPriceGross)
	return *taxedMoney
}

func OrderGetTotalPrice(order model.Order) goprices.TaxedMoney {
	totalPriceNet, _ := goprices.NewMoneyFromDecimal(order.TotalNetAmount, string(order.Currency))
	totalPriceGross, _ := goprices.NewMoneyFromDecimal(order.TotalGrossAmount, string(order.Currency))

	taxedMoney, _ := goprices.NewTaxedMoney(*totalPriceNet, *totalPriceGross)
	return *taxedMoney
}

func OrderGetTotalPaidPrice(o model.Order) goprices.Money {
	money, _ := goprices.NewMoneyFromDecimal(o.TotalPaidAmount, o.Currency.String())
	return *money
}

func OrderGetUnDiscountedTotalPrice(order model.Order) goprices.TaxedMoney {
	undiscountedTotalPriceNet, _ := goprices.NewMoneyFromDecimal(order.UndiscountedTotalNetAmount, string(order.Currency))
	undiscountedTotalPriceGross, _ := goprices.NewMoneyFromDecimal(order.UndiscountedTotalGrossAmount, string(order.Currency))

	taxedMoney, _ := goprices.NewTaxedMoney(*undiscountedTotalPriceNet, *undiscountedTotalPriceGross)
	return *taxedMoney
}

func OrderIsDraft(o model.Order) bool {
	return o.Status == model.OrderStatusDraft
}

func OrderIsUnConfirmed(o model.Order) bool {
	return o.Status == model.OrderStatusUnconfirmed
}

func OrderIsOpen(o model.Order) bool {
	return o.Status == model.OrderStatusDraft || o.Status == model.OrderStatusPartiallyFulfilled
}

func OrderIsFullyPaid(o model.Order) bool {
	return o.TotalGrossAmount.LessThanOrEqual(o.TotalPaidAmount)
}

func OrderIsPartlyPaid(o model.Order) bool {
	return o.TotalPaidAmount.GreaterThan(decimal.Zero)
}

var OrderGetTotalCapturedPrice = OrderGetTotalPaidPrice

func OrderGetTotalBalance(o model.Order) goprices.Money {
	var capturedPrice = OrderGetTotalCapturedPrice(o)
	var totalPrice = OrderGetTotalPrice(o)
	result, _ := capturedPrice.Sub(totalPrice.GetGross())
	return *result
}

func OrderGetTotalWeight(o model.Order) measurement.Weight {
	return measurement.Weight{
		Amount: float64(o.WeightAmount),
		Unit:   measurement.WeightUnit(o.WeightUnit),
	}
}

type OrderFilterOption struct {
	CommonQueryOptions

	// INNER JOIN user ON ... WHERE user.email % ... OR user.first_name % ... OR user.last_name % ...
	// NOTE: this use trigram search feature
	Customer string

	Search          string
	ChannelIdOrSlug string

	PaymentChargeStatus qm.QueryMod // INNER JOIN payment ON ... WHERE payment.charge_status = ...
	Statuses            util.AnyArray[OrderFilterStatus]

	AnnotateBillingAddressNames     bool
	AnnotateLastPaymentChargeStatus bool

	ChannelSlug string
}

type CustomOrder struct {
	model.Order
	OrderBillingAddressLastName  *string                    `boil:"billing_address_last_name" json:"billing_address_last_name"`
	OrderBillingAddressFirstName *string                    `boil:"billing_address_first_name" json:"billing_address_first_name"`
	OrderLastPaymentChargeStatus *model.PaymentChargeStatus `boil:"last_payment_charge_status" json:"last_payment_charge_status"`
}

var CustomOrderTableColumns = struct {
	OrderBillingAddressLastName  string
	OrderBillingAddressFirstName string
	OrderLastPaymentChargeStatus string
}{
	OrderBillingAddressLastName:  "billing_address_last_name",
	OrderBillingAddressFirstName: "billing_address_first_name",
	OrderLastPaymentChargeStatus: "last_payment_charge_status",
}

type CustomOrderSlice []*CustomOrder

// NOTE: when model.Order is updated, this function should be updated too
func OrderScanValues(o *model.Order) []any {
	return []any{
		&o.ID,
		&o.CreatedAt,
		&o.Status,
		&o.UserID,
		&o.LanguageCode,
		&o.TrackingClientID,
		&o.BillingAddressID,
		&o.ShippingAddressID,
		&o.UserEmail,
		&o.OriginalID,
		&o.Origin,
		&o.Currency,
		&o.ShippingMethodID,
		&o.CollectionPointID,
		&o.ShippingMethodName,
		&o.CollectionPointName,
		&o.ChannelID,
		&o.ShippingPriceNetAmount,
		&o.ShippingPriceGrossAmount,
		&o.ShippingTaxRate,
		&o.Token,
		&o.CheckoutToken,
		&o.TotalNetAmount,
		&o.UndiscountedTotalNetAmount,
		&o.TotalGrossAmount,
		&o.UndiscountedTotalGrossAmount,
		&o.TotalPaidAmount,
		&o.VoucherID,
		&o.DisplayGrossPrices,
		&o.CustomerNote,
		&o.WeightAmount,
		&o.WeightUnit,
		&o.RedirectURL,
		&o.Metadata,
		&o.PrivateMetadata,
	}
}

func OrderLinePreSave(ol *model.OrderLine) {
	if ol.CreatedAt == 0 {
		ol.CreatedAt = GetMillis()
	}
	if ol.ID == "" {
		ol.ID = NewId()
	}
	OrderLineCommonPre(ol)
}

func OrderLineCommonPre(ol *model.OrderLine) {
	ol.ProductName = SanitizeUnicode(ol.ProductName)
	ol.VariantName = SanitizeUnicode(ol.VariantName)
	ol.TranslatedProductName = SanitizeUnicode(ol.TranslatedProductName)
	ol.TranslatedVariantName = SanitizeUnicode(ol.TranslatedVariantName)

	if !ol.UnitDiscountReason.IsNil() {
		*ol.UnitDiscountReason.String = SanitizeUnicode(*ol.UnitDiscountReason.String)
	}
	if ol.UnitDiscountType.IsValid() != nil {
		ol.UnitDiscountType = model.DiscountValueTypeFixed
	}
	if ol.UnitDiscountValue.LessThan(decimal.Zero) {
		ol.UnitDiscountValue = decimal.NewFromInt(0)
	}
	if ol.TaxRate.IsNil() {
		ol.TaxRate = model_types.NewNullDecimal(decimal.NewFromInt(0))
	}
}

func OrderLineQuantityUnFulfilled(o model.OrderLine) int {
	return o.Quantity - o.QuantityFulfilled
}

func OrderLineSliceGetProductVariantIDs(ols model.OrderLineSlice) []string {
	var ids []string
	for _, ol := range ols {
		if ol != nil && !ol.VariantID.IsNil() {
			ids = append(ids, *ol.VariantID.String)
		}
	}
	return ids
}

func OrderLineIsValid(o model.OrderLine) *AppError {
	if !IsValidId(o.ID) {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(o.OrderID) {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.order_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !o.VariantID.IsNil() && !IsValidId(*o.VariantID.String) {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.variant_id.app_error", nil, "", http.StatusBadRequest)
	}
	if o.ProductName == "" {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.product_name.app_error", nil, "", http.StatusBadRequest)
	}
	if o.Quantity < 0 {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.quantity.app_error", nil, "", http.StatusBadRequest)
	}
	if o.QuantityFulfilled < 0 {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.quantity_fulfilled.app_error", nil, "", http.StatusBadRequest)
	}
	if o.QuantityFulfilled > o.Quantity {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.quantity_fulfilled.app_error", nil, "", http.StatusBadRequest)
	}
	if o.UnitPriceNetAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.unit_price_net_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.UnitPriceGrossAmount.LessThan(decimal.Zero) {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.unit_price_gross_amount.app_error", nil, "", http.StatusBadRequest)
	}
	if o.UnitDiscountValue.LessThan(decimal.Zero) {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.unit_discount_value.app_error", nil, "", http.StatusBadRequest)
	}
	if o.CreatedAt <= 0 {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if o.Currency.IsValid() != nil {
		return NewAppError("OrderLineIsValid", "model.order_line.is_valid.currency.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type OrderLineFilterOptions struct {
	CommonQueryOptions

	RelatedOrderConditions qm.QueryMod // INNER JOIN order ON ... WHERE ...
	Preload                []string
	VariantProductID       qm.QueryMod // INNER JOIN product_variant ON ... WHERE product_variants.product_id ...
}

func FulfillmentPreSave(f *model.Fulfillment) {
	if f.ID == "" {
		f.ID = NewId()
	}
	if f.CreatedAt == 0 {
		f.CreatedAt = GetMillis()
	}
	FulfillmentCommonPre(f)
}

func FulfillmentCommonPre(f *model.Fulfillment) {
	if f.Status.IsValid() != nil {
		f.Status = model.FulfillmentStatusFulfilled
	}
}

func FulfillmentIsValid(f model.Fulfillment) *AppError {
	if !IsValidId(f.ID) {
		return NewAppError("FulfillmentIsValid", "model.fulfillment.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(f.OrderID) {
		return NewAppError("FulfillmentIsValid", "model.fulfillment.is_valid.order_id.app_error", nil, "", http.StatusBadRequest)
	}
	if f.CreatedAt <= 0 {
		return NewAppError("FulfillmentIsValid", "model.fulfillment.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if f.Status.IsValid() != nil {
		return NewAppError("FulfillmentIsValid", "model.fulfillment.is_valid.status.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type FulfillmentFilterOption struct {
	CommonQueryOptions
	Preload                []string
	HaveNoFulfillmentLines bool
	FulfillmentLineID      qm.QueryMod
}

func FulfillmentLinePreSave(fl *model.FulfillmentLine) {
	if fl.ID == "" {
		fl.ID = NewId()
	}
}

func FulfillmentLineIsValid(fl model.FulfillmentLine) *AppError {
	if !IsValidId(fl.ID) {
		return NewAppError("FulfillmentLineIsValid", "model.fulfillment_line.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(fl.OrderLineID) {
		return NewAppError("FulfillmentLineIsValid", "model.fulfillment_line.is_valid.order_line_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(fl.FulfillmentID) {
		return NewAppError("FulfillmentLineIsValid", "model.fulfillment_line.is_valid.fulfillment_id.app_error", nil, "", http.StatusBadRequest)
	}
	if fl.Quantity <= 0 {
		return NewAppError("FulfillmentLineIsValid", "model.fulfillment_line.is_valid.quantity.app_error", nil, "", http.StatusBadRequest)
	}
	if !fl.StockID.IsNil() && !IsValidId(*fl.StockID.String) {
		return NewAppError("FulfillmentLineIsValid", "model.fulfillment_line.is_valid.stock_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type FulfillmentLineFilterOption struct {
	CommonQueryOptions
	RelatedFulfillmentConds qm.QueryMod
	Preload                 []string
}

func OrderLineGetUnitPrice(o model.OrderLine) goprices.TaxedMoney {
	unitPriceNet, _ := goprices.NewMoneyFromDecimal(o.UnitPriceNetAmount, string(o.Currency))
	unitPriceGross, _ := goprices.NewMoneyFromDecimal(o.UnitPriceGrossAmount, string(o.Currency))

	taxedMoney, _ := goprices.NewTaxedMoney(*unitPriceNet, *unitPriceGross)
	return *taxedMoney
}

type OrderLineData struct {
	Line        model.OrderLine
	Quantity    int
	Variant     *model.ProductVariant // can be nil
	Replace     bool                  // default false
	WarehouseID *string               // can be nil
}

func OrderLineString(o model.OrderLine) string {
	if o.VariantName != "" {
		return fmt.Sprintf("%s (%s)", o.ProductName, o.VariantName)
	}
	return o.ProductName
}

// if either order line's total price net amount or total price gross amount is nil, then the result will be nil.
func OrderLineGetTotalPrice(o model.OrderLine) *goprices.TaxedMoney {
	if o.TotalPriceNetAmount.IsNil() || o.TotalPriceGrossAmount.IsNil() {
		return nil
	}

	totalPriceNet, _ := goprices.NewMoneyFromDecimal(*o.TotalPriceNetAmount.Decimal, string(o.Currency))
	totalPriceGross, _ := goprices.NewMoneyFromDecimal(*o.TotalPriceGrossAmount.Decimal, string(o.Currency))

	taxedMoney, _ := goprices.NewTaxedMoney(*totalPriceNet, *totalPriceGross)
	return taxedMoney
}

func OrderLineSetTotalPrice(o *model.OrderLine, price goprices.TaxedMoney) {
	o.TotalPriceNetAmount = model_types.NewNullDecimal(price.GetNet().GetAmount())
	o.TotalPriceGrossAmount = model_types.NewNullDecimal(price.GetGross().GetAmount())
	o.Currency = model.Currency(strings.ToUpper(price.GetCurrency()))
}

func OrderLineSetUnitPrice(o *model.OrderLine, price goprices.TaxedMoney) {
	o.UnitPriceNetAmount = price.GetNet().GetAmount()
	o.UnitPriceGrossAmount = price.GetGross().GetAmount()
	o.Currency = model.Currency(strings.ToUpper(price.GetCurrency()))
}
