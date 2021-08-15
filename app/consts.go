package app

const (
	InvalidArgumentAppErrorID    = "app.invalid_arguments.app_error" // InvalidArgumentAppErrorID used when creating app errors on invalid argument
	InternalServerErrorID        = "app.internal_error.app_error"
	ProductNotPublishedAppErrID  = "app.checkout.product_unpublished.app_error"
	NewMoneyCreationAppErrorID   = "app.money_creation_error.app_error"
	InvalidPromoCodeAppErrorID   = "app.invalid_promo_code.app_error" // use this to make invalid promo code *AppError(s) in checkout app and giftcard app
	ErrorCalculatingMoneyErrorID = "app.error_calculating_money.app_error"
)
