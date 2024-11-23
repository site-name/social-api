package model_helper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gosimple/slug"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type VoucherCustomerFilterOption struct {
	CommonQueryOptions
}

func VoucherCustomerIsValid(v model.VoucherCustomer) *AppError {
	if !IsValidId(v.ID) {
		return NewAppError("VoucherCustomerIsValid", "model.voucher_customer.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if !IsValidId(v.VoucherCodeID) {
		return NewAppError("VoucherCustomerIsValid", "model.voucher_customer.is_valid.voucher_code_id.app_error", nil, "invalid voucher id", http.StatusBadRequest)
	}
	if !IsValidEmail(v.CustomerEmail) {
		return NewAppError("VoucherCustomerIsValid", "model.voucher_customer.is_valid.email.app_error", nil, "invalid email", http.StatusBadRequest)
	}
	return nil
}

type VoucherTranslationFilterOption struct {
	CommonQueryOptions
}

type SaleFilterOption struct {
	CommonQueryOptions
	SaleChannelListing_ChannelSlug qm.QueryMod // INNER JOIN sale_channel_listings ON ... INNER JOIN channels ON ... WHERE channels.slug ...

	AnnotateSaleDiscountValue bool   // LEFT JOIN sale_channel_listings ON ... LEFT JOIN channels ON ... WHERE channels.slug ...
	ChannelSlug               string // This param goes with `AnnotateSaleDiscountValue`
}

func (s SaleFilterOption) Validate() *AppError {
	if s.AnnotateSaleDiscountValue && !slug.IsSlug(s.ChannelSlug) {
		return NewAppError("SaleFilterOption.Validate", InvalidArgumentAppErrorID, map[string]any{"Fields": "channel_slug"}, "please provide related channel slug", http.StatusBadRequest)
	}
	return nil
}

type VoucherFilterOption struct {
	CommonQueryOptions

	Annotate_MinValues bool // this options tell store whether to annotate `min_discount_value` and `min_spent_amount` to the result
	ChannelIdOrSlug    string
}

func (v VoucherFilterOption) Validate() *AppError {
	if v.Annotate_MinValues &&
		(!slug.IsSlug(v.ChannelIdOrSlug) && !IsValidId(v.ChannelIdOrSlug)) {
		return NewAppError("VoucherFilterOption.Validate", InvalidArgumentAppErrorID, map[string]any{"Fields": "channel_id_or_slug"}, "please provide related channel id or slug", http.StatusBadRequest)
	}
	return nil
}

type CustomVoucher struct {
	model.Voucher
	MinDiscountValue *decimal.Decimal `boil:"min_discount_value" json:"min_discount_value" toml:"min_discount_value" yaml:"min_discount_value"`
	MinSpentAmount   *decimal.Decimal `boil:"min_spent_amount" json:"min_spent_amount" toml:"min_spent_amount" yaml:"min_spent_amount"`
}

var CustomVoucherTableColumns = struct {
	MinDiscountValue string
	MinSpentAmount   string
}{
	MinDiscountValue: `"vouchers.min_discount_value"`,
	MinSpentAmount:   `"vouchers.min_spent_amount"`,
}

// NOTE: this function's return value MUST BE updated when fields of `model.Voucher` are updated
func VoucherScanValues(v *model.Voucher) []any {
	return []any{
		&v.ID,
		&v.Type,
		&v.Name,
		&v.UsageLimit,
		&v.StartDate,
		&v.EndDate,
		&v.ApplyOncePerOrder,
		&v.ApplyOncePerCustomer,
		&v.SingleUse,
		&v.OnlyForStaff,
		&v.DiscountValueType,
		&v.Countries,
		&v.MinCheckoutItemsQuantity,
		&v.CreatedAt,
		&v.UpdatedAt,
		&v.Metadata,
		&v.PrivateMetadata,
	}
}

func CustomVoucherScanValues(v *CustomVoucher) []any {
	return append(VoucherScanValues(&v.Voucher), &v.MinDiscountValue, &v.MinSpentAmount)
}

type CustomVoucherSlice []*CustomVoucher

func SaleIsValid(s model.Sale) *AppError {
	if !IsValidId(s.ID) {
		return NewAppError("SaleIsValid", "model.sale.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if s.Name == "" {
		return NewAppError("SaleIsValid", "model.sale.is_valid.name.app_error", nil, "invalid name", http.StatusBadRequest)
	}
	if model_types.PrimitiveIsNotNilAndEqual(s.EndDate.Int64, 0) {
		return NewAppError("SaleIsValid", "model.sale.is_valid.end_date.app_error", nil, "invalid end date", http.StatusBadRequest)
	}
	if s.CreatedAt <= 0 {
		return NewAppError("SaleIsValid", "model.sale.is_valid.created_at.app_error", nil, "invalid created at", http.StatusBadRequest)
	}
	if s.UpdatedAt <= 0 {
		return NewAppError("SaleIsValid", "model.sale.is_valid.updated_at.app_error", nil, "invalid updated at", http.StatusBadRequest)
	}
	if s.EndDate.Int64 != nil && *s.EndDate.Int64 < s.StartDate {
		return NewAppError("SaleIsValid", "model.sale.is_valid.end_date.app_error", nil, "start date must be before end date", http.StatusBadRequest)
	}
	if s.Type.IsValid() != nil {
		return NewAppError("SaleIsValid", "model.sale.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}
	return nil
}

func SalePreSave(s *model.Sale) {
	saleCommonPre(s)
	if s.CreatedAt == 0 {
		s.CreatedAt = GetMillis()
	}
	s.UpdatedAt = s.CreatedAt
	if s.ID == "" {
		s.ID = NewId()
	}
}

func saleCommonPre(s *model.Sale) {
	s.Name = SanitizeUnicode(s.Name)
	if s.Type.IsValid() != nil {
		s.Type = model.DiscountValueTypeFixed
	}
}

func SalePreUpdate(s *model.Sale) {
	saleCommonPre(s)
	s.UpdatedAt = GetMillis()
}

func VoucherIsValid(v *model.Voucher) *AppError {
	if v == nil {
		return nil
	}
	if !IsValidId(v.ID) {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if model_types.PrimitiveIsNotNilAndEqual(v.Name.String, "") {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.name.app_error", nil, "invalid name", http.StatusBadRequest)
	}
	if v.CreatedAt <= 0 {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.created_at.app_error", nil, "invalid created at", http.StatusBadRequest)
	}
	if v.UpdatedAt <= 0 {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.updated_at.app_error", nil, "invalid updated at", http.StatusBadRequest)
	}
	if v.Type.IsValid() != nil {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}
	if v.DiscountValueType.IsValid() != nil {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.discount_value_type.app_error", nil, "please provide valid discount value type", http.StatusBadRequest)
	}
	if v.EndDate.Int64 != nil && *v.EndDate.Int64 <= v.StartDate {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.end_date.app_error", nil, "invalid end date", http.StatusBadRequest)
	}
	if v.StartDate == 0 {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.start_date.app_error", nil, "invalid start date", http.StatusBadRequest)
	}
	if !v.UsageLimit.IsNil() && *v.UsageLimit.Int <= 0 {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.usage_limit.app_error", nil, "invalid usage limit", http.StatusBadRequest)
	}
	for _, country := range strings.Fields(v.Countries) {
		if model.CountryCode(country).IsValid() != nil {
			return NewAppError("VoucherIsValid", "model.voucher.is_valid.countries.app_error", nil, "please provide valid countries", http.StatusBadRequest)
		}

	}
	return nil
}

func VoucherPreSave(v *model.Voucher) {
	if v.CreatedAt == 0 {
		v.CreatedAt = GetMillis()
	}
	v.UpdatedAt = v.CreatedAt
	if v.ID == "" {
		v.ID = NewId()
	}
	voucherCommonPre(v)
}

func VoucherPreUpdate(v *model.Voucher) {
	v.UpdatedAt = GetMillis()
	voucherCommonPre(v)
}

func voucherCommonPre(v *model.Voucher) {
	if model_types.PrimitiveIsNotNilAndNotEqual(v.Name.String, "") {
		*v.Name.String = SanitizeUnicode(*v.Name.String)
	}
	if v.DiscountValueType.IsValid() != nil {
		v.DiscountValueType = model.DiscountValueTypeFixed
	}
	if v.Type.IsValid() != nil {
		v.Type = model.VoucherTypeEntireOrder
	}
	if !v.UsageLimit.IsNil() && *v.UsageLimit.Int < 0 {
		v.UsageLimit.Int = GetPointerOfValue(0)
	}
	v.Countries = strings.ToUpper(v.Countries)
}

func VoucherValidateMinCheckoutItemsQuantity(v model.Voucher, quantity int) *NotApplicable {
	if v.MinCheckoutItemsQuantity > quantity {
		return &NotApplicable{
			Where:                    "ValidateMinCheckoutItemsQuantity",
			Message:                  fmt.Sprintf("This offer is onlyvalid for orders with a minimum of %d in quantity", v.MinCheckoutItemsQuantity),
			MinCheckoutItemsQuantity: v.MinCheckoutItemsQuantity,
		}
	}
	return nil
}

func SaleChannelListingPreSave(s *model.SaleChannelListing) {
	if s.CreatedAt == 0 {
		s.CreatedAt = GetMillis()
	}
	if s.ID == "" {
		s.ID = NewId()
	}
	saleChannelListingCommonPre(s)
}

func SaleChannelListingPreUpdate(s *model.SaleChannelListing) {
	saleChannelListingCommonPre(s)
}

func saleChannelListingCommonPre(s *model.SaleChannelListing) {
	if s.Currency.IsValid() != nil {
		s.Currency = DEFAULT_CURRENCY
	}
}

func SaleChannelListingIsValid(s model.SaleChannelListing) *AppError {
	if !IsValidId(s.ID) {
		return NewAppError("SaleChannelListingIsValid", "model.sale_channel_listing.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if s.CreatedAt <= 0 {
		return NewAppError("SaleChannelListingIsValid", "model.sale_channel_listing.is_valid.created_at.app_error", nil, "invalid created at", http.StatusBadRequest)
	}
	if !IsValidId(s.SaleID) {
		return NewAppError("SaleChannelListingIsValid", "model.sale_channel_listing.is_valid.sale_id.app_error", nil, "invalid sale id", http.StatusBadRequest)
	}
	if !IsValidId(s.ChannelID) {
		return NewAppError("SaleChannelListingIsValid", "model.sale_channel_listing.is_valid.channel_slug.app_error", nil, "invalid channel slug", http.StatusBadRequest)
	}
	if s.Currency.IsValid() != nil {
		return NewAppError("SaleChannelListingIsValid", "model.sale_channel_listing.is_valid.currency.app_error", nil, "invalid currency", http.StatusBadRequest)
	}
	return nil
}

type SaleChannelListingFilterOption struct {
	CommonQueryOptions
}

func VoucherChannelListingPreSave(v *model.VoucherChannelListing) {
	if v.CreatedAt == 0 {
		v.CreatedAt = GetMillis()
	}
	if v.ID == "" {
		v.ID = NewId()
	}
	voucherChannelListingCommonPre(v)
}

func voucherChannelListingCommonPre(v *model.VoucherChannelListing) {
	if v.Currency.IsValid() != nil {
		v.Currency = DEFAULT_CURRENCY
	}
}

func VoucherChannelListingIsValid(v model.VoucherChannelListing) *AppError {
	if !IsValidId(v.ID) {
		return NewAppError("VoucherChannelListingIsValid", "model.voucher_channel_listing.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if v.CreatedAt <= 0 {
		return NewAppError("VoucherChannelListingIsValid", "model.voucher_channel_listing.is_valid.created_at.app_error", nil, "invalid created at", http.StatusBadRequest)
	}
	if !IsValidId(v.VoucherID) {
		return NewAppError("VoucherChannelListingIsValid", "model.voucher_channel_listing.is_valid.voucher_id.app_error", nil, "invalid voucher id", http.StatusBadRequest)
	}
	if !IsValidId(v.ChannelID) {
		return NewAppError("VoucherChannelListingIsValid", "model.voucher_channel_listing.is_valid.channel_id.app_error", nil, "invalid channel id", http.StatusBadRequest)
	}
	if v.Currency.IsValid() != nil {
		return NewAppError("VoucherChannelListingIsValid", "model.voucher_channel_listing.is_valid.currency.app_error", nil, "invalid currency", http.StatusBadRequest)
	}
	return nil
}

func VoucherChannelListingGetDiscount(v model.VoucherChannelListing) goprices.Money {
	money, _ := goprices.NewMoneyFromDecimal(v.DiscountValue, v.Currency.String())
	return *money
}

func VoucherChannelListingGetMinSpent(v model.VoucherChannelListing) goprices.Money {
	money, _ := goprices.NewMoneyFromDecimal(v.MinSpendAmount, v.Currency.String())
	return *money
}

func VoucherChannelListingSetDiscount(v *model.VoucherChannelListing, discount goprices.Money) {
	v.DiscountValue = discount.GetAmount()
	v.Currency = model.Currency(discount.GetCurrency())
}

func VoucherChannelListingSetMinSpent(v *model.VoucherChannelListing, minSpent goprices.Money) {
	v.MinSpendAmount = minSpent.GetAmount()
	v.Currency = model.Currency(minSpent.GetCurrency())
}

type SaleTranslationFilterOption struct {
	CommonQueryOptions
}

type VoucherChannelListingFilterOption struct {
	CommonQueryOptions
}

func VoucherTranslationPreSave(v *model.VoucherTranslation) {
	if v.CreatedAt == 0 {
		v.CreatedAt = GetMillis()
	}
	if v.ID == "" {
		v.ID = NewId()
	}
	VoucherTranslationCommonPre(v)
}

func VoucherTranslationCommonPre(v *model.VoucherTranslation) {
	if v.LanguageCode.IsValid() != nil {
		v.LanguageCode = DEFAULT_LOCALE
	}
	v.Name = SanitizeUnicode(v.Name)
}

func VoucherTranslationIsValid(v model.VoucherTranslation) *AppError {
	if !IsValidId(v.ID) {
		return NewAppError("VoucherTranslationIsValid", "model.voucher_translation.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if v.CreatedAt <= 0 {
		return NewAppError("VoucherTranslationIsValid", "model.voucher_translation.is_valid.created_at.app_error", nil, "invalid created at", http.StatusBadRequest)
	}
	if v.LanguageCode.IsValid() != nil {
		return NewAppError("VoucherTranslationIsValid", "model.voucher_translation.is_valid.language_code.app_error", nil, "invalid language code", http.StatusBadRequest)
	}
	if v.Name == "" {
		return NewAppError("VoucherTranslationIsValid", "model.voucher_translation.is_valid.name.app_error", nil, "invalid name", http.StatusBadRequest)
	}
	return nil
}

type OrderDiscountFilterOption struct {
	CommonQueryOptions
}

func OrderDiscountPreSave(o *model.OrderDiscount) {
	orderDiscountCommonPre(o)
	if o.ID == "" {
		o.ID = NewId()
	}
}

func orderDiscountCommonPre(o *model.OrderDiscount) {
	if o.Type.IsValid() != nil {
		o.Type = model.OrderDiscountTypeManual
	}
	if o.ValueType.IsValid() != nil {
		o.ValueType = model.DiscountValueTypeFixed
	}
	if !o.Name.IsNil() {
		*o.Name.String = SanitizeUnicode(*o.Name.String)
	}
	if !o.TranslatedName.IsNil() {
		*o.TranslatedName.String = SanitizeUnicode(*o.TranslatedName.String)
	}
	if !o.Reason.IsNil() {
		*o.Reason.String = SanitizeUnicode(*o.Reason.String)
	}
	if o.Currency.IsValid() != nil {
		o.Currency = DEFAULT_CURRENCY
	}
}

func OrderDiscountPreUpdate(o *model.OrderDiscount) {
	orderDiscountCommonPre(o)
}

func OrderDiscountIsValid(o model.OrderDiscount) *AppError {
	if !IsValidId(o.ID) {
		return NewAppError("OrderDiscountIsValid", "model.order_discount.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if o.Type.IsValid() != nil {
		return NewAppError("OrderDiscountIsValid", "model.order_discount.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}
	if o.ValueType.IsValid() != nil {
		return NewAppError("OrderDiscountIsValid", "model.order_discount.is_valid.value_type.app_error", nil, "please provide valid value type", http.StatusBadRequest)
	}
	if o.Currency.IsValid() != nil {
		return NewAppError("OrderDiscountIsValid", "model.order_discount.is_valid.currency.app_error", nil, "please provide valid currency", http.StatusBadRequest)
	}
	if !o.OrderID.IsNil() && !IsValidId(*o.OrderID.String) {
		return NewAppError("OrderDiscountIsValid", "model.order_discount.is_valid.order_id.app_error", nil, "please provide valid order id", http.StatusBadRequest)
	}
	return nil
}

type CustomSale struct {
	model.Sale
	DiscountValue *decimal.Decimal `boil:"discount_value" json:"discount_value" toml:"discount_value" yaml:"discount_value"`
}

var CustomSaleTableColumns = struct {
	DiscountValue string
}{
	DiscountValue: `"sales.discount_value"`,
}

type CustomSaleSlice []*CustomSale

// NOTE: this function's return value MUST BE updated when fields of `model.Sale` are updated
func SaleScanValues(s *model.Sale) []any {
	return []any{
		&s.ID,
		&s.Name,
		&s.Type,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.Metadata,
		&s.PrivateMetadata,
	}
}

func CustomSaleScanValues(s *CustomSale) []any {
	return append(SaleScanValues(&s.Sale), &s.DiscountValue)
}
