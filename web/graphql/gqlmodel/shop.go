package gqlmodel

type Shop struct {
	AvailablePaymentGatewayIDs          []string                      `json:"availablePaymentGateways"`         // changed, PaymentGateway
	AvailableExternalAuthenticationIDs  []string                      `json:"availableExternalAuthentications"` // changed: ExternalAuthentication
	AvailableShippingMethodIDs          []*string                     `json:"availableShippingMethods"`         // changed ShippingMethod
	Countries                           []CountryDisplay              `json:"countries"`
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
	TranslationID                       *string                       `json:"translation"`                         // changed, ShopTranslation
	AutomaticFulfillmentDigitalProducts *bool                         `json:"automaticFulfillmentDigitalProducts"` //
	DefaultDigitalMaxDownloads          *int                          `json:"defaultDigitalMaxDownloads"`          //
	DefaultDigitalURLValidDays          *int                          `json:"defaultDigitalUrlValidDays"`          //
	CompanyAddressID                    *string                       `json:"companyAddress"`                      // changed, Address
	CustomerSetPasswordURL              *string                       `json:"customerSetPasswordUrl"`
	StaffNotificationRecipients         []*StaffNotificationRecipient `json:"staffNotificationRecipients"`
	Limits                              *LimitInfo                    `json:"limits"`
	Version                             string                        `json:"version"`
}
