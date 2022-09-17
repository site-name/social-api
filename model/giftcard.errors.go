package model

type InvalidPromoCode struct {
	Where   string
	Message string
	Code    GiftcardErrorCode
}

type GiftcardErrorCode string

const (
	ALREADY_EXISTS GiftcardErrorCode = "already_exists"
	GRAPHQL_ERROR_ GiftcardErrorCode = "graphql_error"
	INVALID_       GiftcardErrorCode = "invalid"
	NOT_FOUND_     GiftcardErrorCode = "not_found"
	REQUIRED_      GiftcardErrorCode = "required"
	UNIQUE_        GiftcardErrorCode = "unique"
)

// NewInvalidPromoCode is common function to create invalid promo code with code = "invalid"
func NewInvalidPromoCode(where, message string) *InvalidPromoCode {
	return &InvalidPromoCode{
		Where:   where,
		Message: message,
		Code:    INVALID_,
	}
}
