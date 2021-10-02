package gqlmodel

// --------------------------- original implementation ----------------------

// type Shop struct {
// 	AvailablePaymentGateways            []*PaymentGateway             `json:"availablePaymentGateways"`
// 	AvailableExternalAuthentications    []*ExternalAuthentication     `json:"availableExternalAuthentications"`
// 	AvailableShippingMethods            []*ShippingMethod             `json:"availableShippingMethods"`
// 	ChannelCurrencies                   []string                      `json:"channelCurrencies"`
// 	Countries                           []*CountryDisplay             `json:"countries"`
// 	DefaultCountry                      *CountryDisplay               `json:"defaultCountry"`
// 	DefaultMailSenderName               *string                       `json:"defaultMailSenderName"`
// 	DefaultMailSenderAddress            *string                       `json:"defaultMailSenderAddress"`
// 	Description                         *string                       `json:"description"`
// 	Domain                              *Domain                       `json:"domain"`
// 	Languages                           []*LanguageDisplay            `json:"languages"`
// 	Name                                string                        `json:"name"`
// 	Permissions                         []*Permission                 `json:"permissions"`
// 	PhonePrefixes                       []*string                     `json:"phonePrefixes"`
// 	HeaderText                          *string                       `json:"headerText"`
// 	IncludeTaxesInPrices                bool                          `json:"includeTaxesInPrices"`
// 	FulfillmentAutoApprove              bool                          `json:"fulfillmentAutoApprove"`
// 	FulfillmentAllowUnpaid              bool                          `json:"fulfillmentAllowUnpaid"`
// 	DisplayGrossPrices                  bool                          `json:"displayGrossPrices"`
// 	ChargeTaxesOnShipping               bool                          `json:"chargeTaxesOnShipping"`
// 	TrackInventoryByDefault             *bool                         `json:"trackInventoryByDefault"`
// 	DefaultWeightUnit                   *WeightUnitsEnum              `json:"defaultWeightUnit"`
// 	Translation                         *ShopTranslation              `json:"translation"`
// 	AutomaticFulfillmentDigitalProducts *bool                         `json:"automaticFulfillmentDigitalProducts"`
// 	DefaultDigitalMaxDownloads          *int                          `json:"defaultDigitalMaxDownloads"`
// 	DefaultDigitalURLValidDays          *int                          `json:"defaultDigitalUrlValidDays"`
// 	CompanyAddress                      *Address                      `json:"companyAddress"`
// 	CustomerSetPasswordURL              *string                       `json:"customerSetPasswordUrl"`
// 	StaffNotificationRecipients         []*StaffNotificationRecipient `json:"staffNotificationRecipients"`
// 	Limits                              *LimitInfo                    `json:"limits"`
// 	Version                             string                        `json:"version"`
// }

type Shop struct {
	ChannelCurrencies                   []string                      `json:"channelCurrencies"`
	FulfillmentAutoApprove              bool                          `json:"fulfillmentAutoApprove"`
	FulfillmentAllowUnpaid              bool                          `json:"fulfillmentAllowUnpaid"`
	AvailableExternalAuthenticationIDs  []string                      `json:"availableExternalAuthentications"` // ExternalAuthentication
	DefaultCountry                      *CountryDisplay               `json:"defaultCountry"`
	DefaultMailSenderName               *string                       `json:"defaultMailSenderName"`
	DefaultMailSenderAddress            *string                       `json:"defaultMailSenderAddress"`
	Description                         *string                       `json:"description"`
	Domain                              *Domain                       `json:"domain"`
	Languages                           []*LanguageDisplay            `json:"languages"`
	Name                                string                        `json:"name"`
	Permissions                         []*Permission                 `json:"permissions"`
	PhonePrefixes                       []*string                     `json:"phonePrefixes"`
	HeaderText                          *string                       `json:"headerText"`
	IncludeTaxesInPrices                bool                          `json:"includeTaxesInPrices"`
	DisplayGrossPrices                  bool                          `json:"displayGrossPrices"`
	ChargeTaxesOnShipping               bool                          `json:"chargeTaxesOnShipping"`
	TrackInventoryByDefault             *bool                         `json:"trackInventoryByDefault"`
	DefaultWeightUnit                   *WeightUnitsEnum              `json:"defaultWeightUnit"`
	AutomaticFulfillmentDigitalProducts *bool                         `json:"automaticFulfillmentDigitalProducts"` //
	DefaultDigitalMaxDownloads          *int                          `json:"defaultDigitalMaxDownloads"`          //
	DefaultDigitalURLValidDays          *int                          `json:"defaultDigitalUrlValidDays"`          //
	CompanyAddressID                    *string                       `json:"companyAddress"`                      // *Address
	CustomerSetPasswordURL              *string                       `json:"customerSetPasswordUrl"`
	StaffNotificationRecipients         []*StaffNotificationRecipient `json:"staffNotificationRecipients"`
	Limits                              *LimitInfo                    `json:"limits"`
	Version                             string                        `json:"version"`
	AvailablePaymentGateways            func() []PaymentGateway       `json:"availablePaymentGateways"`
	Translation                         func() *ShopTranslation       `json:"translation"`
	AvailableShippingMethods            func() []*ShippingMethod      `json:"availableShippingMethods"`
	Countries                           func() []CountryDisplay       `json:"countries"`
}
