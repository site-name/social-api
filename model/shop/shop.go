package shop

import (
	"net/mail"
	"net/url"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
)

// max length values for some shop fields
const (
	SHOP_NAME_MAX_LENGTH                       = 100
	SHOP_DESCRIPTION_MAX_LENGTH                = 200
	SHOP_DEFAULT_WEIGHT_UNIT_MAX_LENGTH        = 10
	SHOP_DEFAULT_MAX_EMAIL_DISPLAY_NAME_LENGTH = 78
)

// Shop represents a selling unit
type Shop struct {
	Id                                  string  `json:"id"`
	OwnerID                             string  `json:"owner_id"`
	CreateAt                            int64   `json:"create_at"`
	UpdateAt                            int64   `json:"update_at"`
	Name                                string  `json:"name"`
	Description                         string  `json:"description"`
	TopMenuID                           string  `json:"top_menu_id"`
	IncludeTaxesInPrice                 *bool   `json:"include_taxes_in_prices"`                // default true
	DisplayGrossPrices                  *bool   `json:"display_gross_prices"`                   // default true
	ChargeTaxesOnShipping               *bool   `json:"charge_taxes_on_shipping"`               // default true
	TrackInventoryByDefault             *bool   `json:"track_inventory_by_default"`             // default true
	DefaultWeightUnit                   string  `json:"default_weight_unit"`                    // default kg
	AutomaticFulfillmentDigitalProducts *bool   `json:"automatic_fulfillment_digital_products"` // default true
	DefaultDigitalMaxDownloads          *uint   `json:"default_digital_max_downloads"`
	DefaultDigitalUrlValidDays          *uint   `json:"default_digital_url_valid_days"`
	AddressID                           *string `json:"address_id"`
	DefaultMailSenderName               *string `json:"default_mail_sender_name"`
	DefaultMailSenderAddress            string  `json:"default_mail_sender_address"`
	CustomerSetPasswordUrl              *string `json:"customer_set_password_url"`
	AutomaticallyConfirmAllNewOrders    *bool   `json:"automatically_confirm_all_new_orders"` // default true
}

type ShopDefaultDigitalContentSettings struct {
	AutomaticFulfillmentDigitalProducts *bool
	DefaultDigitalMaxDownloads          *uint
	DefaultDigitalUrlValidDays          *uint
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

func (s *Shop) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.shop.is_valid.%s.app_error",
		"shop_id=",
		"Shop.IsValid",
	)
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(s.OwnerID) {
		return outer("owner_id", nil)
	}
	if !model.IsValidId(s.TopMenuID) {
		return outer("top_menu_id", nil)
	}
	if s.AddressID != nil && !model.IsValidId(*s.AddressID) {
		return outer("address_id", nil)
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
		var err bool
		if len(*s.CustomerSetPasswordUrl) > model.URL_LINK_MAX_LENGTH {
			err = true
		}
		_, parseErr := url.Parse(*s.CustomerSetPasswordUrl)
		if parseErr != nil {
			err = true
		}
		if err {
			return outer("customer_set_password_url", &s.Id)
		}
	}

	if s.DefaultMailSenderName != nil && utf8.RuneCountInString(*s.DefaultMailSenderName) > SHOP_DEFAULT_MAX_EMAIL_DISPLAY_NAME_LENGTH {
		return outer("default_mail_sender_name", &s.Id)
	}

	return nil
}

func (s *Shop) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	if s.CreateAt == 0 {
		s.CreateAt = model.GetMillis()
	}
	s.UpdateAt = s.CreateAt
	s.UpsertCommon()
}

func (s *Shop) UpsertCommon() {
	s.Name = model.SanitizeUnicode(s.Name)
	s.Description = model.SanitizeUnicode(s.Description)
	if s.IncludeTaxesInPrice == nil {
		s.IncludeTaxesInPrice = model.NewBool(true)
	}
	if s.DisplayGrossPrices == nil {
		s.DisplayGrossPrices = model.NewBool(true)
	}
	if s.ChargeTaxesOnShipping == nil {
		s.ChargeTaxesOnShipping = model.NewBool(true)
	}
	if s.TrackInventoryByDefault == nil {
		s.TrackInventoryByDefault = model.NewBool(true)
	}
	if s.AutomaticallyConfirmAllNewOrders == nil {
		s.AutomaticallyConfirmAllNewOrders = model.NewBool(true)
	}
}

func (s *Shop) PreUpdate() {
	s.UpdateAt = model.GetMillis()
	s.UpsertCommon()
}
