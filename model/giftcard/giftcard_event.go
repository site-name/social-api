package giftcard

import "github.com/sitename/sitename/model"

// max lengths for some fields of giftcard event
const (
	GiftCardEventTypeMaxLength = 255
)

// valid values for giftcard event type
const (
	GiftcardEventType_Issued                = "issued"
	GiftcardEventType_Bought                = "bought"
	GiftcardEventType_Updated               = "updated"
	GiftcardEventType_Activated             = "activated"
	GiftcardEventType_Deactivated           = "deactivated"
	GiftcardEventType_BalanceReset          = "balance_reset"
	GiftcardEventType_ExpirySettingsUpdated = "expiry_settings_updated"
	GiftcardEventType_SentToCustomer        = "sent_to_customer"
	GiftcardEventType_Resent                = "resent"
)

var GiftCardEventTypeMap = map[string]string{
	GiftcardEventType_Issued:                "The gift card was created be staff user or app.",
	GiftcardEventType_Bought:                "The gift card was bought by customer.",
	GiftcardEventType_Updated:               "The gift card was updated.",
	GiftcardEventType_Activated:             "The gift card was activated.",
	GiftcardEventType_Deactivated:           "The gift card was deactivated.",
	GiftcardEventType_BalanceReset:          "The gift card balance was reset.",
	GiftcardEventType_ExpirySettingsUpdated: "The gift card expiry settings was updated.",
	GiftcardEventType_SentToCustomer:        "The gift card was sent to the customer.",
	GiftcardEventType_Resent:                "The gift card was resent to the customer.",
}

type GiftCardEvent struct {
	Id         string          `json:"id"`
	Date       int64           `json:"date"` // not editable
	Type       string          `json:"type"`
	Parameters model.StringMap `json:"parameters"`  // default map[stirng]string{}
	UserID     *string         `json:"user_id"`     // ON DELETE SET NULL
	GiftcardID string          `json:"giftcard_id"` // ON DELETE CASCADE
}

func (g *GiftCardEvent) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.giftcard_event.is_valid.%s.app_error",
		"giftcard_event_id=",
		"GiftcardEvent.IsValid",
	)

	if !model.IsValidId(g.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(g.GiftcardID) {
		return outer("giftcard_id", &g.Id)
	}
	if g.UserID != nil && !model.IsValidId(*g.UserID) {
		return outer("user_id", &g.Id)
	}
	if len(g.Type) > GiftCardEventTypeMaxLength || GiftCardEventTypeMap[g.Type] == "" {
		return outer("type", &g.Id)
	}
	if g.Date <= 0 {
		return outer("date", &g.Id)
	}

	return nil
}

func (g *GiftCardEvent) PreSave() {
	if !model.IsValidId(g.Id) {
		g.Id = model.NewId()
	}
	g.Date = model.GetMillis()
	if g.Parameters == nil {
		g.Parameters = make(model.StringMap)
	}
}
