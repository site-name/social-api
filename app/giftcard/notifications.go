package giftcard

import (
	"github.com/sitename/sitename/app/plugin"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/giftcard"
)

// SendGiftcardNotification Trigger sending a gift card notification for the given recipient
func (s *ServiceGiftcard) SendGiftcardNotification(requesterUser *account.User, _ interface{}, customerUser *account.User, email string, giftCard giftcard.GiftCard, manager interfaces.PluginManagerInterface, channelID string, resending bool) *model.AppError {
	var (
		userPayload interface{}
		userID      *string
	)
	if requesterUser != nil {
		userPayload = s.srv.AccountService().GetDefaultUserPayload(*requesterUser)
		userID = &requesterUser.Id
	}

	shop, appErr := s.srv.ShopService().ShopById(manager.GetShopID())
	if appErr != nil {
		return appErr
	}

	payload := model.StringInterface{
		"gift_card":         s.GetDefaultGiftcardPayload(giftCard),
		"user":              userPayload,
		"requester_user_id": userID,
		"requester_app_id":  nil,
		"recipient_email":   email,
		"resending":         resending,
		"domain":            s.srv.Config().ServiceSettings.SiteURL,
		"site_name":         shop.Name,
	}

	_, appErr = manager.Notify(plugin.SEND_GIFT_CARD, payload, channelID, "")
	return appErr
}

func (s *ServiceGiftcard) GetDefaultGiftcardPayload(giftCard giftcard.GiftCard) model.StringInterface {
	return model.StringInterface{
		"id":       giftCard.Id,
		"code":     giftCard.Code,
		"balance":  giftCard.CurrentBalanceAmount,
		"currency": giftCard.Currency,
	}
}
