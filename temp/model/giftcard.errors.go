package model

import "fmt"

type InvalidPromoCode struct {
	Where   string
	Message string
	Code    GiftcardErrorCode
}

type GiftcardErrorCode string

const (
	GIFT_CARD_ERROR_CODE_ALREADY_EXISTS GiftcardErrorCode = "already_exists"
	GIFT_CARD_ERROR_CODE_GRAPHQL_ERROR  GiftcardErrorCode = "graphql_error"
	GIFT_CARD_ERROR_CODE_INVALID        GiftcardErrorCode = "invalid"
	GIFT_CARD_ERROR_CODE_NOT_FOUND      GiftcardErrorCode = "not_found"
	GIFT_CARD_ERROR_CODE_REQUIRED       GiftcardErrorCode = "required"
	GIFT_CARD_ERROR_CODE_UNIQUE         GiftcardErrorCode = "unique"
)

// NewInvalidPromoCode is common function to create invalid promo code with code = "invalid"
func NewInvalidPromoCode(where, message string) *InvalidPromoCode {
	return &InvalidPromoCode{
		Where:   where,
		Message: message,
		Code:    GIFT_CARD_ERROR_CODE_INVALID,
	}
}

func (e *InvalidPromoCode) Error() string {
	return fmt.Sprintf("InvalidPromoCode<where: %s, message: %s, code: %s>", e.Where, e.Message, e.Code)
}
