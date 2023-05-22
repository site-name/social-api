package giftcard

import (
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
)

// SendGiftcardNotification Trigger sending a gift card notification for the given recipient
func (s *ServiceGiftcard) SendGiftcardNotification(requesterUser *model.User, _ interface{}, customerUser *model.User, email string, giftCard model.GiftCard, manager interfaces.PluginManagerInterface, channelID string, resending bool) *model.AppError {
	var (
		userPayload interface{}
		userID      *string
	)
	if requesterUser != nil {
		userPayload = s.srv.AccountService().GetDefaultUserPayload(requesterUser)
		userID = &requesterUser.Id
	}

	payload := model.StringInterface{
		"gift_card":         s.GetDefaultGiftcardPayload(giftCard),
		"user":              userPayload,
		"requester_user_id": userID,
		"requester_app_id":  nil,
		"recipient_email":   email,
		"resending":         resending,
		"domain":            s.srv.Config().ServiceSettings.SiteURL,
		"site_name":         s.srv.Config().ServiceSettings.SiteName,
	}

	_, appErr := manager.Notify(model.SEND_GIFT_CARD, payload, channelID, "")
	return appErr
}

func (s *ServiceGiftcard) GetDefaultGiftcardPayload(giftCard model.GiftCard) model.StringInterface {
	return model.StringInterface{
		"id":       giftCard.Id,
		"code":     giftCard.Code,
		"balance":  giftCard.CurrentBalanceAmount,
		"currency": giftCard.Currency,
	}
}
