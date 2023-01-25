package model

import (
	"github.com/Masterminds/squirrel"
)

// The different gift card event types
type GiftcardEventType string

const (
	ISSUED              GiftcardEventType = "issued"
	BOUGHT              GiftcardEventType = "bought"
	UPDATED             GiftcardEventType = "updated"
	ACTIVATED           GiftcardEventType = "activated"
	DEACTIVATED         GiftcardEventType = "deactivated"
	BALANCE_RESET       GiftcardEventType = "balance_reset"
	EXPIRY_DATE_UPDATED GiftcardEventType = "expiry_date_updated"
	TAG_UPDATED         GiftcardEventType = "tag_updated"
	SENT_TO_CUSTOMER    GiftcardEventType = "sent_to_customer"
	RESENT              GiftcardEventType = "resent"
	NOTE_ADDED_         GiftcardEventType = "note_added"
	USED_IN_ORDER       GiftcardEventType = "used_in_order"
)

var GiftCardEventsString = map[GiftcardEventType]string{
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
	NOTE_ADDED_:         "A note was added to the gift card.",
	USED_IN_ORDER:       "The gift card was used in order.",
}

// max lengths for some fields of giftcard event
const (
	GiftCardEventTypeMaxLength = 255
)

type GiftCardEvent struct {
	Id   string            `json:"id"`
	Date int64             `json:"date"` // not editable
	Type GiftcardEventType `json:"type"`

	// if "expiry_date" presents, its value should has format of "2006-01-02" or of type time.Time
	// To reduce number of type checking steps, below are possible keys and their according value Types you must follow:
	//  "message": string
	//  "email": string
	//  "order_id": string
	//  "tag": string
	//  "old_tag": string
	Parameters StringInterface `json:"parameters"`  // default map[stirng]string{}
	UserID     *string         `json:"user_id"`     // ON DELETE SET NULL
	GiftcardID string          `json:"giftcard_id"` // ON DELETE CASCADE
}

// GiftCardEventFilterOption is used for building squirrel queries.
type GiftCardEventFilterOption struct {
	Id         squirrel.Sqlizer
	Type       squirrel.Sqlizer
	Parameters squirrel.Sqlizer
	GiftcardID squirrel.Sqlizer
}

func (g *GiftCardEvent) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"giftcard_event.is_valid.%s.app_error",
		"giftcard_event_id=",
		"GiftcardEvent.IsValid",
	)

	if !IsValidId(g.Id) {
		return outer("id", nil)
	}
	if !IsValidId(g.GiftcardID) {
		return outer("giftcard_id", &g.Id)
	}
	if g.UserID != nil && !IsValidId(*g.UserID) {
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
	if !IsValidId(g.Id) {
		g.Id = NewId()
	}
	g.Date = GetMillis()
	if g.Parameters == nil {
		g.Parameters = make(StringInterface)
	}
}
