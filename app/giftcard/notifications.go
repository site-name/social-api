package giftcard

import (
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/giftcard"
)

// SendGiftcardNotification Trigger sending a gift card notification for the given recipient
func (s *ServiceGiftcard) SendGiftcardNotification(requesterUser account.User, _ interface{}, customerUser account.User, giftCard giftcard.GiftCard, manager interfaces.PluginManagerInterface, channelID string, resending bool) {

}

func (s *ServiceGiftcard) GetDefaultGiftcardPayload(giftCard giftcard.GiftCard) model.StringInterface {
	return model.StringInterface{
		"id":       giftCard.Id,
		"code":     giftCard.Code,
		"balance":  giftCard.CurrentBalanceAmount,
		"currency": giftCard.Currency,
	}
}
