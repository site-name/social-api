package app

const (
	InvalidArgumentAppErrorID         = "app.invalid_arguments.app_error" // InvalidArgumentAppErrorID used when creating app errors on invalid argument
	InternalServerErrorID             = "app.internal_error.app_error"
	ProductNotPublishedAppErrID       = "app.checkout.product_unpublished.app_error"
	NewMoneyCreationAppErrorID        = "app.money_creation_error.app_error"
	InvalidPromoCodeAppErrorID        = "app.invalid_promo_code.app_error" // use this to make invalid promo code *AppError(s) in checkout app and giftcard app
	ErrorCalculatingMoneyErrorID      = "app.error_calculating_money.app_error"
	ErrorCreatingTransactionErrorID   = "app.error_creating_transaction.app_error"
	ErrorCommittingTransactionErrorID = "app.error_committing_transaction.app_error"
	ErrorCalculatingMeasurementID     = "app.error_calculating_measurement.app_error"
	ErrorMarshallingDataID            = "app.error_marshalling_data.app_error"
	ErrorUnMarshallingDataID          = "app.error_unmarshalling_data.app_error"
)
