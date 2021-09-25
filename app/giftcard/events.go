package giftcard

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
)

// CommonCreateGiftcardEvent is common method for creating giftcard events
func (s *ServiceGiftcard) CommonCreateGiftcardEvent(giftcardID, userID string, parameters model.StringMap, Type string) (*giftcard.GiftCardEvent, *model.AppError) {
	panic("not implemented")
}
