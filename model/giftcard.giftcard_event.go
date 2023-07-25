package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

// The different gift card event types
type GiftcardEventType string

func (t GiftcardEventType) IsValid() bool {
	return GiftCardEventsString[t] != ""
}

// valid types for giftcard
const (
	GIFT_CARD_EVENT_TYPE_ISSUED              GiftcardEventType = "issued"
	GIFT_CARD_EVENT_TYPE_BOUGHT              GiftcardEventType = "bought"
	GIFT_CARD_EVENT_TYPE_UPDATED             GiftcardEventType = "updated"
	GIFT_CARD_EVENT_TYPE_ACTIVATED           GiftcardEventType = "activated"
	GIFT_CARD_EVENT_TYPE_DEACTIVATED         GiftcardEventType = "deactivated"
	GIFT_CARD_EVENT_TYPE_BALANCE_RESET       GiftcardEventType = "balance_reset"
	GIFT_CARD_EVENT_TYPE_EXPIRY_DATE_UPDATED GiftcardEventType = "expiry_date_updated"
	GIFT_CARD_EVENT_TYPE_TAG_UPDATED         GiftcardEventType = "tag_updated"
	GIFT_CARD_EVENT_TYPE_SENT_TO_CUSTOMER    GiftcardEventType = "sent_to_customer"
	GIFT_CARD_EVENT_TYPE_RESENT              GiftcardEventType = "resent"
	GIFT_CARD_EVENT_TYPE_NOTE_ADDED          GiftcardEventType = "note_added"
	GIFT_CARD_EVENT_TYPE_USED_IN_ORDER       GiftcardEventType = "used_in_order"
)

var GiftCardEventsString = map[GiftcardEventType]string{
	GIFT_CARD_EVENT_TYPE_ISSUED:              "The gift card was created be staff user or app.",
	GIFT_CARD_EVENT_TYPE_BOUGHT:              "The gift card was bought by customer.",
	GIFT_CARD_EVENT_TYPE_UPDATED:             "The gift card was updated.",
	GIFT_CARD_EVENT_TYPE_ACTIVATED:           "The gift card was activated.",
	GIFT_CARD_EVENT_TYPE_DEACTIVATED:         "The gift card was deactivated.",
	GIFT_CARD_EVENT_TYPE_BALANCE_RESET:       "The gift card balance was reset.",
	GIFT_CARD_EVENT_TYPE_EXPIRY_DATE_UPDATED: "The gift card expiry date was updated.",
	GIFT_CARD_EVENT_TYPE_TAG_UPDATED:         "The gift card tag was updated.",
	GIFT_CARD_EVENT_TYPE_SENT_TO_CUSTOMER:    "The gift card was sent to the customer.",
	GIFT_CARD_EVENT_TYPE_RESENT:              "The gift card was resent to the customer.",
	GIFT_CARD_EVENT_TYPE_NOTE_ADDED:          "A note was added to the gift card.",
	GIFT_CARD_EVENT_TYPE_USED_IN_ORDER:       "The gift card was used in order.",
}

type GiftCardEvent struct {
	Id   string            `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Date int64             `json:"date" gorm:"type:bigint;column:Date"` // not editable
	Type GiftcardEventType `json:"type" gorm:"type:varchar(255);column:Type"`
	// if "expiry_date" presents, its value should has format of "2006-01-02" or of type time.Time
	// To reduce number of type checking steps, below are possible keys and their according value Types you must follow:
	//  "message": string
	//  "email": string
	//  "order_id": string
	//  "tag": *string
	//  "old_tag": *string
	//  "balance": map[string]any{}
	//  "expiry_date": *time.Time
	//  "old_expiry_date": *time.Time
	Parameters StringInterface `json:"parameters" gorm:"type:jsonb;column:Parameters"` // default map[stirng]string{}
	UserID     *string         `json:"user_id" gorm:"type:uuid;column:UserID"`         // ON DELETE SET NULL
	GiftcardID string          `json:"giftcard_id" gorm:"type:uuid;column:GiftcardID"` // ON DELETE CASCADE
}

func (c *GiftCardEvent) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *GiftCardEvent) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *GiftCardEvent) TableName() string             { return GiftcardEventTableName }

// GiftCardEventFilterOption is used for building squirrel queries.
type GiftCardEventFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (g *GiftCardEvent) IsValid() *AppError {
	if !IsValidId(g.GiftcardID) {
		return NewAppError("GiftcardEvent.IsValid", "model.giftcard_event.is_valid.giftcard_id.app_error", nil, "", http.StatusBadRequest)
	}
	if g.UserID != nil && !IsValidId(*g.UserID) {
		return NewAppError("GiftcardEvent.IsValid", "model.giftcard_event.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !g.Type.IsValid() {
		return NewAppError("GiftcardEvent.IsValid", "model.giftcard_event.is_valid.type.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (g *GiftCardEvent) commonPre() {
	if g.Parameters == nil {
		g.Parameters = make(StringInterface)
	}
}
