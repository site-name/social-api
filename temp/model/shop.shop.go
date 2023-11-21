package model

import (
	"net/mail"
	"net/url"
	"regexp"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/modules/measurement"
)

type GiftCardSettingsExpiryType string

func (g GiftCardSettingsExpiryType) IsValid() bool {
	return GiftCardSettingsExpiryTypeValues[g] != ""
}

// valid values for shop's giftcard expiry type
const (
	NEVER_EXPIRE  GiftCardSettingsExpiryType = "never_expire"
	EXPIRY_PERIOD GiftCardSettingsExpiryType = "expiry_period"
)

var GiftCardSettingsExpiryTypeValues = map[GiftCardSettingsExpiryType]string{
	NEVER_EXPIRE:  "Never expire",
	EXPIRY_PERIOD: "Expiry period",
}

// max length values for some shop fields
const (
	SHOP_NAME_MAX_LENGTH                        = 100
	SHOP_HEADER_TEXT_MAX_LENGTH                 = 100
	SHOP_DESCRIPTION_MAX_LENGTH                 = 200
	SHOP_DEFAULT_WEIGHT_UNIT_MAX_LENGTH         = 10
	SHOP_DEFAULT_MAX_EMAIL_DISPLAY_NAME_LENGTH  = 78
	SHOP_GIFTCARD_EXPIRY_TYPE_MAX_LENGTH        = 32
	SHOP_GIFTCARD_EXPIRY_PERIOD_TYPE_MAX_LENGTH = 32
)

// Shop represents a selling unit
type Shop struct {
	Id                                       string                     `json:"id"`
	CreateAt                                 int64                      `json:"create_at"`
	UpdateAt                                 int64                      `json:"update_at"`
	Name                                     string                     `json:"name"`
	HeaderText                               string                     `json:"header_text"`
	Description                              string                     `json:"description"`
	TopMenuID                                *string                    `json:"top_menu_id"`
	BottomMenuID                             *string                    `json:"bottom_menu_id"`
	IncludeTaxesInPrice                      *bool                      `json:"include_taxes_in_prices"`                // default true
	DisplayGrossPrices                       *bool                      `json:"display_gross_prices"`                   // default true
	ChargeTaxesOnShipping                    *bool                      `json:"charge_taxes_on_shipping"`               // default true
	TrackInventoryByDefault                  *bool                      `json:"track_inventory_by_default"`             // default true
	DefaultWeightUnit                        string                     `json:"default_weight_unit"`                    // default kg
	AutomaticFulfillmentDigitalProducts      *bool                      `json:"automatic_fulfillment_digital_products"` // default true
	DefaultDigitalMaxDownloads               *int                       `json:"default_digital_max_downloads"`
	DefaultDigitalUrlValidDays               *int                       `json:"default_digital_url_valid_days"`
	AddressID                                *string                    `json:"address_id"`
	CompanyAddressID                         *string                    `json:"company_address_id"`
	DefaultMailSenderName                    string                     `json:"default_mail_sender_name"`
	DefaultMailSenderAddress                 string                     `json:"default_mail_sender_address"`
	CustomerSetPasswordUrl                   *string                    `json:"customer_set_password_url"`
	AutomaticallyConfirmAllNewOrders         *bool                      `json:"automatically_confirm_all_new_orders"` // default true
	FulfillmentAutoApprove                   *bool                      `json:"fulfillment_auto_approve"`             // default *true
	FulfillmentAllowUnPaid                   *bool                      `json:"fulfillment_allow_unpaid"`             // default *true
	GiftcardExpiryType                       GiftCardSettingsExpiryType `json:"gift_card_expiry_type"`                // default "never_expire"
	GiftcardExpiryPeriodType                 TimePeriodType             `json:"gift_card_expiry_period_type"`
	GiftcardExpiryPeriod                     *int                       `json:"gift_card_expiry_period"`
	AutomaticallyFulfillNonShippableGiftcard *bool                      `json:"automatically_fulfill_non_shippable_gift_card"` // default *true

	companyAddress *Address `db:"-"`
}

func (s *Shop) GetCompanyAddress() *Address {
	return s.companyAddress
}

func (s *Shop) SetCompanyAddress(a *Address) {
	s.companyAddress = a
}

type ShopFilterOptions struct {
	Id   squirrel.Sqlizer
	Name squirrel.Sqlizer

	SelectRelatedCompanyAddress bool
}

type ShopDefaultDigitalContentSettings struct {
	AutomaticFulfillmentDigitalProducts *bool
	DefaultDigitalMaxDownloads          *int
	DefaultDigitalUrlValidDays          *int
}

func (s *Shop) DefaultFromEmail() (string, error) {
	if s.DefaultMailSenderAddress == "" {
		return s.DefaultMailSenderAddress, nil
	}

	address, err := mail.ParseAddress(s.DefaultMailSenderAddress)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}

func (s *Shop) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.shop.is_valid.%s.app_error",
		"shop_id=",
		"Shop.IsValid",
	)
	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.TopMenuID != nil && !IsValidId(*s.TopMenuID) {
		return outer("top_menu_id", &s.Id)
	}
	if s.BottomMenuID != nil && !IsValidId(*s.BottomMenuID) {
		return outer("top_menu_id", &s.Id)
	}
	if s.AddressID != nil && !IsValidId(*s.AddressID) {
		return outer("address_id", &s.Id)
	}
	if s.CompanyAddressID != nil && !IsValidId(*s.CompanyAddressID) {
		return outer("company_address_id", &s.Id)
	}
	if len(s.HeaderText) > SHOP_HEADER_TEXT_MAX_LENGTH {
		return outer("header_text", &s.Id)
	}
	if utf8.RuneCountInString(s.Name) > SHOP_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}
	if utf8.RuneCountInString(s.Description) > SHOP_DESCRIPTION_MAX_LENGTH {
		return outer("description", &s.Id)
	}
	if len(s.DefaultWeightUnit) > SHOP_DEFAULT_WEIGHT_UNIT_MAX_LENGTH || measurement.WEIGHT_UNIT_STRINGS[measurement.WeightUnit(s.DefaultWeightUnit)] == "" {
		return outer("default_weight_unit", &s.Id)
	}
	if s.CustomerSetPasswordUrl != nil {
		if len(*s.CustomerSetPasswordUrl) > URL_LINK_MAX_LENGTH {
			return outer("customer_set_password_url", &s.Id)
		}
		_, err := url.Parse(*s.CustomerSetPasswordUrl)
		if err != nil {
			return outer("customer_set_password_url", &s.Id)
		}
	}
	if matched, err := regexp.MatchString(`[\n\r]`, s.DefaultMailSenderName); err == nil && matched {
		return outer("default_mail_sender_name", &s.Id)
	}
	if utf8.RuneCountInString(s.DefaultMailSenderName) > SHOP_DEFAULT_MAX_EMAIL_DISPLAY_NAME_LENGTH {
		return outer("default_mail_sender_name", &s.Id)
	}
	if s.FulfillmentAutoApprove == nil {
		return outer("fulfillment_auto_approve", &s.Id)
	}
	if s.FulfillmentAllowUnPaid == nil {
		return outer("fulfillment_allow_unpaid", &s.Id)
	}
	if len(s.GiftcardExpiryType) > SHOP_GIFTCARD_EXPIRY_TYPE_MAX_LENGTH || GiftCardSettingsExpiryTypeValues[s.GiftcardExpiryType] == "" {
		return outer("gift_card_expiry_type", &s.Id)
	}
	if !s.GiftcardExpiryPeriodType.IsValid() {
		return outer("gift_card_expiry_period_type", &s.Id)
	}
	if s.AutomaticallyFulfillNonShippableGiftcard == nil {
		return outer("automatically_fulfill_non_shippable_gift_card", &s.Id)
	}
	if s.GiftcardExpiryPeriod != nil && *s.GiftcardExpiryPeriod < 0 {
		return outer("giftcard_expiry_period", &s.Id)
	}

	return nil
}

func (s *Shop) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	if s.CreateAt == 0 {
		s.CreateAt = GetMillis()
	}
	s.UpdateAt = s.CreateAt
	s.commonPre()
}

func (s *Shop) commonPre() {
	s.Name = SanitizeUnicode(s.Name)
	s.Description = SanitizeUnicode(s.Description)
	if s.IncludeTaxesInPrice == nil {
		s.IncludeTaxesInPrice = GetPointerOfValue(true)
	}
	if s.DisplayGrossPrices == nil {
		s.DisplayGrossPrices = GetPointerOfValue(true)
	}
	if s.ChargeTaxesOnShipping == nil {
		s.ChargeTaxesOnShipping = GetPointerOfValue(true)
	}
	if s.TrackInventoryByDefault == nil {
		s.TrackInventoryByDefault = GetPointerOfValue(true)
	}
	if s.AutomaticallyConfirmAllNewOrders == nil {
		s.AutomaticallyConfirmAllNewOrders = GetPointerOfValue(true)
	}
	if s.FulfillmentAllowUnPaid == nil {
		s.FulfillmentAllowUnPaid = GetPointerOfValue(true)
	}
	if s.FulfillmentAutoApprove == nil {
		s.FulfillmentAutoApprove = GetPointerOfValue(true)
	}
	if s.AutomaticallyFulfillNonShippableGiftcard == nil {
		s.AutomaticallyFulfillNonShippableGiftcard = GetPointerOfValue(true)
	}
	if len(s.GiftcardExpiryType) == 0 {
		s.GiftcardExpiryType = NEVER_EXPIRE
	}
	if s.GiftcardExpiryPeriod != nil && *s.GiftcardExpiryPeriod < 0 {
		s.GiftcardExpiryPeriod = GetPointerOfValue(0)
	}
}

func (s *Shop) PreUpdate() {
	s.UpdateAt = GetMillis()
	s.commonPre()
}

func (s *Shop) DeepCopy() *Shop {
	res := *s

	if s.TopMenuID != nil {
		res.TopMenuID = GetPointerOfValue(*s.TopMenuID)
	}
	if s.BottomMenuID != nil {
		res.BottomMenuID = GetPointerOfValue(*s.BottomMenuID)
	}
	if s.IncludeTaxesInPrice != nil {
		res.IncludeTaxesInPrice = GetPointerOfValue(*s.IncludeTaxesInPrice)
	}
	if s.DisplayGrossPrices != nil {
		res.DisplayGrossPrices = GetPointerOfValue(*s.DisplayGrossPrices)
	}
	if s.ChargeTaxesOnShipping != nil {
		res.ChargeTaxesOnShipping = GetPointerOfValue(*s.ChargeTaxesOnShipping)
	}
	if s.TrackInventoryByDefault != nil {
		res.TrackInventoryByDefault = GetPointerOfValue(*s.TrackInventoryByDefault)
	}
	if s.AutomaticFulfillmentDigitalProducts != nil {
		res.AutomaticFulfillmentDigitalProducts = GetPointerOfValue(*s.AutomaticFulfillmentDigitalProducts)
	}
	if s.DefaultDigitalMaxDownloads != nil {
		res.DefaultDigitalMaxDownloads = GetPointerOfValue(*s.DefaultDigitalMaxDownloads)
	}
	if s.DefaultDigitalUrlValidDays != nil {
		res.DefaultDigitalUrlValidDays = GetPointerOfValue(*s.DefaultDigitalUrlValidDays)
	}
	if s.AddressID != nil {
		res.AddressID = GetPointerOfValue(*s.AddressID)
	}
	if s.CompanyAddressID != nil {
		res.CompanyAddressID = GetPointerOfValue(*s.CompanyAddressID)
	}
	if s.CustomerSetPasswordUrl != nil {
		res.CustomerSetPasswordUrl = GetPointerOfValue(*s.CustomerSetPasswordUrl)
	}
	if s.AutomaticallyConfirmAllNewOrders != nil {
		res.AutomaticallyConfirmAllNewOrders = GetPointerOfValue(*s.AutomaticallyConfirmAllNewOrders)
	}
	if s.FulfillmentAutoApprove != nil {
		res.FulfillmentAutoApprove = GetPointerOfValue(*s.FulfillmentAutoApprove)
	}
	if s.FulfillmentAllowUnPaid != nil {
		res.FulfillmentAllowUnPaid = GetPointerOfValue(*s.FulfillmentAllowUnPaid)
	}
	if s.GiftcardExpiryPeriod != nil {
		res.GiftcardExpiryPeriod = GetPointerOfValue(*s.GiftcardExpiryPeriod)
	}
	if s.AutomaticallyFulfillNonShippableGiftcard != nil {
		res.AutomaticallyFulfillNonShippableGiftcard = GetPointerOfValue(*s.AutomaticallyFulfillNonShippableGiftcard)
	}
	if s.companyAddress != nil {
		res.companyAddress = s.companyAddress.DeepCopy()
	}

	return &res
}
