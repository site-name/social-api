package model_helper

import (
	"net/http"
	"strings"

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
	if !IsValidId(v.VoucherID) {
		return NewAppError("VoucherCustomerIsValid", "model.voucher_customer.is_valid.voucher_id.app_error", nil, "invalid voucher id", http.StatusBadRequest)
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
	SaleChannelListing_ChannelSlug qm.QueryMod
}

type VoucherFilterOption struct {
	CommonQueryOptions

	Annotate_MinDiscountValue bool
	Annotate_MinSpentAmount   bool
	ChannelIdOrSlug           string
}

func (v VoucherFilterOption) Validate() error {
	if (v.Annotate_MinDiscountValue || v.Annotate_MinSpentAmount) && v.ChannelIdOrSlug == "" {
		return NewAppError("VoucherFilterOption.Validate", "ChannelIdOrSlug", nil, "please provide related channel id or slug", http.StatusBadRequest)
	}
	return nil
}

type CustomVoucher struct {
	model.Voucher
	MinDiscountValue *decimal.Decimal
	MinSpentAmount   *decimal.Decimal
}

type CustomVoucherSlice []*CustomVoucher

func SaleIsValid(s model.Sale) *AppError {
	if !IsValidId(s.ID) {
		return NewAppError("SaleIsValid", "model.sale.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if s.Name == "" {
		return NewAppError("SaleIsValid", "model.sale.is_valid.name.app_error", nil, "invalid name", http.StatusBadRequest)
	}
	if s.EndDate.IsNotNilAndEqual(0) {
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
	s.CreatedAt = GetMillis()
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

func VoucherIsValid(v model.Voucher) *AppError {
	if !IsValidId(v.ID) {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if v.Name.IsNotNilAndEqual("") {
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
	if !PromoCodeRegex.MatchString(v.Code) {
		return NewAppError("VoucherIsValid", "model.voucher.is_valid.code.app_error", nil, "code must look like 78GH-UJKI-90RD", http.StatusBadRequest)
	}
	for _, country := range strings.Fields(v.Countries) {
		if model.CountryCode(country).IsValid() != nil {
			return NewAppError("VoucherIsValid", "model.voucher.is_valid.countries.app_error", nil, "please provide valid countries", http.StatusBadRequest)
		}

	}
	return nil
}

func VoucherPreSave(v *model.Voucher) {
	v.CreatedAt = GetMillis()
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
	if v.OnlyForStaff.IsNil() {
		v.OnlyForStaff = model_types.NewNullBool(false)
	}
	if v.Name.IsNotNilAndNotEqual("") {
		*v.Name.String = SanitizeUnicode(*v.Name.String)
	}
	if v.DiscountValueType.IsValid() != nil {
		v.DiscountValueType = model.DiscountValueTypeFixed
	}
	if v.Type.IsValid() != nil {
		v.Type = model.VoucherTypeEntireOrder
	}
	if v.UsageLimit < 0 {
		v.UsageLimit = 0
	}
	v.Countries = strings.ToUpper(v.Countries)
	if v.Code == "" {
		v.Code = NewPromoCode()
	}
}

func SaleChannelListingPreSave(s *model.SaleChannelListing) {
	s.CreatedAt = GetMillis()
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
	v.CreatedAt = GetMillis()
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
	return goprices.Money{
		Amount:   v.DiscountValue,
		Currency: v.Currency.String(),
	}
}

func VoucherChannelListingGetMinSpent(v model.VoucherChannelListing) goprices.Money {
	return goprices.Money{
		Amount:   v.MinSpendAmount,
		Currency: v.Currency.String(),
	}
}

func VoucherChannelListingSetDiscount(v *model.VoucherChannelListing, discount goprices.Money) {
	v.DiscountValue = discount.Amount
	v.Currency = model.Currency(discount.Currency)
}

func VoucherChannelListingSetMinSpent(v *model.VoucherChannelListing, minSpent goprices.Money) {
	v.MinSpendAmount = minSpent.Amount
	v.Currency = model.Currency(minSpent.Currency)
}

type SaleTranslationFilterOption struct {
	CommonQueryOptions
}

type VoucherChannelListingFilterOption struct {
	CommonQueryOptions
}

func VoucherTranslationPreSave(v *model.VoucherTranslation) {
	v.CreatedAt = GetMillis()
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
