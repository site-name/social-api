package giftcard

import (
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
)

// SendGiftcardNotification Trigger sending a gift card notification for the given recipient
func (s *ServiceGiftcard) SendGiftcardNotification(requesterUser *model.User, _ any, customerUser *model.User, email string, giftCard model.GiftCard, manager interfaces.PluginManagerInterface, channelID string, resending bool) *model_helper.AppError {
	var (
		userPayload any
		userID      *string
	)
	if requesterUser != nil {
		userPayload = s.srv.Account.GetDefaultUserPayload(requesterUser)
		userID = &requesterUser.ID
	}

	payload := model_types.JSONString{
		"gift_card":         s.GetDefaultGiftcardPayload(giftCard),
		"user":              userPayload,
		"requester_user_id": model_helper.GetValueOfpointerOrNil(userID),
		"requester_app_id":  nil,
		"recipient_email":   email,
		"resending":         resending,
		"domain":            s.srv.Config().ServiceSettings.SiteURL,
		"site_name":         s.srv.Config().ServiceSettings.SiteName,
	}

	_, appErr := manager.Notify(model_helper.SEND_GIFT_CARD, payload, channelID, "")
	return appErr
}

func (s *ServiceGiftcard) GetDefaultGiftcardPayload(giftCard model.Giftcard) model_types.JSONString {
	return model_types.JSONString{
		"id":       giftCard.ID,
		"code":     giftCard.Code,
		"balance":  giftCard.CurrentBalanceAmount,
		"currency": giftCard.Currency,
	}
}
