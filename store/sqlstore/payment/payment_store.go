package payment

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

type SqlPaymentStore struct {
	store.Store
}

func NewSqlPaymentStore(s store.Store) store.PaymentStore {
	ps := &SqlPaymentStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(payment.Payment{}, "Payments").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CheckoutID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OrderID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("GateWay").SetMaxSize(payment.MAX_LENGTH_PAYMENT_GATEWAY)
		table.ColMap("ChargeStatus").SetMaxSize(payment.MAX_LENGTH_PAYMENT_CHARGE_STATUS).
			SetDefaultConstraint(model.NewString(payment.NOT_CHARGED))
		table.ColMap("Token").SetMaxSize(payment.MAX_LENGTH_PAYMENT_TOKEN)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)
		table.ColMap("BillingEmail").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
		table.ColMap("BillingFirstName").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingLastName").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingCompanyName").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingAddress1").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingAddress2").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingCity").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("BillingCityArea").SetMaxSize(account.ADDRESS_CITY_AREA_MAX_LENGTH)
		table.ColMap("BillingPostalCode").SetMaxSize(account.ADDRESS_POSTAL_CODE_MAX_LENGTH)
		table.ColMap("BillingCountryCode").SetMaxSize(model.SINGLE_COUNTRY_CODE_MAX_LENGTH)
		table.ColMap("BillingCountryArea").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)

		table.ColMap("CcFirstDigits").SetMaxSize(payment.MAX_LENGTH_CC_FIRST_DIGITS)
		table.ColMap("CcLastDigits").SetMaxSize(payment.MAX_LENGTH_CC_LAST_DIGITS)
		table.ColMap("CcBrand").SetMaxSize(payment.MAX_LENGTH_CC_BRAND)
		table.ColMap("PaymentMethodType").SetMaxSize(payment.MAX_LENGTH_PAYMENT_COMMON_256)
		table.ColMap("CustomerIpAddress").SetMaxSize(model.IP_ADDRESS_MAX_LENGTH)
		table.ColMap("ReturnUrl").SetMaxSize(model.URL_LINK_MAX_LENGTH)
		table.ColMap("PspReference").SetMaxSize(payment.PAYMENT_PSP_REFERENCE_MAX_LENGTH)
	}
	return ps
}

func (ps *SqlPaymentStore) CreateIndexesIfNotExists() {
	// NOTE: need more investigation in the future
	ps.CreateIndexIfNotExists("idx_payments_billing_email", "Payments", "BillingEmail")
	ps.CreateIndexIfNotExists("idx_payments_billing_first_name", "Payments", "BillingFirstName")
	ps.CreateIndexIfNotExists("idx_payments_billing_last_name", "Payments", "BillingLastName")
	ps.CreateIndexIfNotExists("idx_payments_billing_company_name", "Payments", "BillingCompanyName")
	ps.CreateIndexIfNotExists("idx_payments_billing_address_1", "Payments", "BillingAddress1")
	ps.CreateIndexIfNotExists("idx_payments_billing_city", "Payments", "BillingCity")
	ps.CreateIndexIfNotExists("idx_payments_billing_city_area", "Payments", "BillingCityArea")

	ps.CreateIndexIfNotExists("idx_payments_billing_email_lower_textpattern", "Payments", "lower(BillingEmail) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_payments_billing_first_name_lower_textpattern", "Payments", "lower(BillingFirstName) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_payments_billing_last_name_lower_textpattern", "Payments", "lower(BillingLastName) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_payments_billing_city_area_lower_textpattern", "Payments", "lower(BillingCityArea) text_pattern_ops")

	ps.CreateIndexIfNotExists("idx_payments_psp_reference", "Payments", "PspReference")
}
