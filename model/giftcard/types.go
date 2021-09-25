package giftcard

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
