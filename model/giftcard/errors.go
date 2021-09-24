package giftcard

type InvalidPromoCode struct {
	Where   string
	Message string
	Code    GiftcardErrorCode
}

type GiftcardErrorCode string

const (
	ALREADY_EXISTS GiftcardErrorCode = "already_exists"
	GRAPHQL_ERROR  GiftcardErrorCode = "graphql_error"
	INVALID        GiftcardErrorCode = "invalid"
	NOT_FOUND      GiftcardErrorCode = "not_found"
	REQUIRED       GiftcardErrorCode = "required"
	UNIQUE         GiftcardErrorCode = "unique"
)

// NewInvalidPromoCode is common function to create invalid promo code with code = "invalid"
func NewInvalidPromoCode(where, message string) *InvalidPromoCode {
	return &InvalidPromoCode{
		Where:   where,
		Message: message,
		Code:    INVALID,
	}
}
