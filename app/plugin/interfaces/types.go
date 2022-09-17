package interfaces

import "github.com/sitename/sitename/model"

// PaymentMethod is type for some methods of PluginManager.
// They are:
//
// 1) AuthorizePayment
//
// 2) CapturePayment
//
// 3) ConfirmPayment
//
// 4) ProcessPayment
//
// 5) RefundPayment
//
// 6) VoidPayment
type PaymentMethod func(gateway string, paymentInformation model.PaymentData, channelID string) (*model.GatewayResponse, error)
