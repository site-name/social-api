package giftcard

import "github.com/sitename/sitename/model"

// The different gift card event types
type GiftCardEvents string

const (
	ISSUED              GiftCardEvents = "issued"
	BOUGHT              GiftCardEvents = "bought"
	UPDATED             GiftCardEvents = "updated"
	ACTIVATED           GiftCardEvents = "activated"
	DEACTIVATED         GiftCardEvents = "deactivated"
	BALANCE_RESET       GiftCardEvents = "balance_reset"
	EXPIRY_DATE_UPDATED GiftCardEvents = "expiry_date_updated"
	TAG_UPDATED         GiftCardEvents = "tag_updated"
	SENT_TO_CUSTOMER    GiftCardEvents = "sent_to_customer"
	RESENT              GiftCardEvents = "resent"
	NOTE_ADDED          GiftCardEvents = "note_added"
	USED_IN_ORDER       GiftCardEvents = "used_in_order"
)

var GiftCardEventsString = map[GiftCardEvents]string{
	ISSUED:              "The gift card was created be staff user or app.",
	BOUGHT:              "The gift card was bought by customer.",
	UPDATED:             "The gift card was updated.",
	ACTIVATED:           "The gift card was activated.",
	DEACTIVATED:         "The gift card was deactivated.",
	BALANCE_RESET:       "The gift card balance was reset.",
	EXPIRY_DATE_UPDATED: "The gift card expiry date was updated.",
	TAG_UPDATED:         "The gift card tag was updated.",
	SENT_TO_CUSTOMER:    "The gift card was sent to the customer.",
	RESENT:              "The gift card was resent to the customer.",
	NOTE_ADDED:          "A note was added to the gift card.",
	USED_IN_ORDER:       "The gift card was used in order.",
}

// max lengths for some fields of giftcard event
const (
	GiftCardEventTypeMaxLength = 255
)

type GiftCardEvent struct {
	Id         string                `json:"id"`
	Date       int64                 `json:"date"` // not editable
	Type       GiftCardEvents        `json:"type"`
	Parameters model.StringInterface `json:"parameters"`  // default map[stirng]string{}
	UserID     *string               `json:"user_id"`     // ON DELETE SET NULL
	GiftcardID string                `json:"giftcard_id"` // ON DELETE CASCADE
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
	if len(g.Type) > GiftCardEventTypeMaxLength || GiftCardEventsString[g.Type] == "" {
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
		g.Parameters = make(model.StringInterface)
	}
}