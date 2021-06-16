package gqlmodel

type Shop struct {
	AvailableExternalAuthenticationIDs  []string                                      `json:"availableExternalAuthentications"` // ExternalAuthentication
	DefaultCountry                      *CountryDisplay                               `json:"defaultCountry"`
	DefaultMailSenderName               *string                                       `json:"defaultMailSenderName"`
	DefaultMailSenderAddress            *string                                       `json:"defaultMailSenderAddress"`
	Description                         *string                                       `json:"description"`
	Domain                              *Domain                                       `json:"domain"`
	Languages                           []*LanguageDisplay                            `json:"languages"`
	Name                                string                                        `json:"name"`
	Permissions                         []*Permission                                 `json:"permissions"`
	PhonePrefixes                       []*string                                     `json:"phonePrefixes"`
	HeaderText                          *string                                       `json:"headerText"`
	IncludeTaxesInPrices                bool                                          `json:"includeTaxesInPrices"`
	DisplayGrossPrices                  bool                                          `json:"displayGrossPrices"`
	ChargeTaxesOnShipping               bool                                          `json:"chargeTaxesOnShipping"`
	TrackInventoryByDefault             *bool                                         `json:"trackInventoryByDefault"`
	DefaultWeightUnit                   *WeightUnitsEnum                              `json:"defaultWeightUnit"`
	AutomaticFulfillmentDigitalProducts *bool                                         `json:"automaticFulfillmentDigitalProducts"` //
	DefaultDigitalMaxDownloads          *int                                          `json:"defaultDigitalMaxDownloads"`          //
	DefaultDigitalURLValidDays          *int                                          `json:"defaultDigitalUrlValidDays"`          //
	CompanyAddressID                    *string                                       `json:"companyAddress"`                      // *Address
	CustomerSetPasswordURL              *string                                       `json:"customerSetPasswordUrl"`
	StaffNotificationRecipients         []*StaffNotificationRecipient                 `json:"staffNotificationRecipients"`
	Limits                              *LimitInfo                                    `json:"limits"`
	Version                             string                                        `json:"version"`
	AvailablePaymentGateways            func(*string, *string) []PaymentGateway       `json:"availablePaymentGateways"`
	Translation                         func(*LanguageCodeEnum) *ShopTranslation      `json:"translation"`
	AvailableShippingMethods            func(string, *AddressInput) []*ShippingMethod `json:"availableShippingMethods"`
	Countries                           func(*LanguageCodeEnum) []CountryDisplay      `json:"countries"`
}
