package model

import (
	"github.com/sitename/sitename/modules/util"
)

// max length values
const (
	CUSTOMER_EVENT_TYPE_MAX_LENGTH = 255
)

// some available types for CustomerEvent's Type attribute
const (
	ACCOUNT_CREATED          = "account_created"
	PASSWORD_RESET_LINK_SENT = "password_reset_link_sent"
	PASSWORD_RESET           = "password_reset"
	PASSWORD_CHANGED         = "password_changed"
	EMAIL_CHANGE_REQUEST     = "email_changed_request"
	EMAIL_CHANGED            = "email_changed"

	// Order related events
	PLACED_ORDER            = "placed_order"            // created an order
	NOTE_ADDED_TO_ORDER     = "note_added_to_order"     // added a note to one of their orders
	DIGITAL_LINK_DOWNLOADED = "digital_link_downloaded" // downloaded a digital good

	// Staff actions over customers events
	CUSTOMER_DELETED    = "customer_deleted" // staff user deleted a customer
	EMAIL_ASSIGNED      = "email_assigned"   // the staff user assigned a email to the customer
	NAME_ASSIGNED       = "name_assigned"    // the staff user added set a name to the customer
	CUSTOMER_NOTE_ADDED = "note_added"       // the staff user added a note to the customer
)

var CustomerEventTypes = []string{
	ACCOUNT_CREATED,
	PASSWORD_RESET_LINK_SENT,
	PASSWORD_RESET,
	PASSWORD_CHANGED,
	EMAIL_CHANGE_REQUEST,
	EMAIL_CHANGED,
	PLACED_ORDER,
	NOTE_ADDED_TO_ORDER,
	DIGITAL_LINK_DOWNLOADED,
	CUSTOMER_DELETED,
	EMAIL_ASSIGNED,
	NAME_ASSIGNED,
	CUSTOMER_NOTE_ADDED,
}

// Model used to store events that happened during the customer lifecycle
type CustomerEvent struct {
	Id         string          `json:"id"`
	Date       int64           `json:"date"`
	Type       string          `json:"type"`
	OrderID    *string         `json:"order_id"`
	UserID     *string         `json:"user_id"`
	Parameters StringInterface `json:"parameters"`
}

func (c *CustomerEvent) ToJSON() string {
	return ModelToJson(c)
}

func (ce *CustomerEvent) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"customer_event.is_valid.%s.app_error",
		"customer_event_id=",
		"CustomerEvent.IsValid",
	)
	if !IsValidId(ce.Id) {
		return outer("id", nil)
	}
	if ce.Date == 0 {
		return outer("date", &ce.Id)
	}
	if ce.UserID != nil && !IsValidId(*ce.UserID) {
		return outer("usder_id", &ce.Id)
	}
	if ce.OrderID != nil && !IsValidId(*ce.OrderID) {
		return outer("order_id", &ce.Id)
	}
	if len(ce.Type) > CUSTOMER_EVENT_TYPE_MAX_LENGTH ||
		!util.ItemInSlice(ce.Type, CustomerEventTypes) {
		return outer("type", &ce.Id)
	}

	return nil
}

func (c *CustomerEvent) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
	if c.Date == 0 {
		c.Date = GetMillis()
	}
	_, ok1 := c.Parameters["currency"]
	_, ok2 := c.Parameters["amount"]
	if ok1 && ok2 {
		c.Parameters["_type"] = "Money"
	}
}

type StaffNotificationRecipient struct {
	Id         string  `json:"id"`
	UserID     *string `json:"user_id"`
	StaffEmail *string `json:"staff_email"`
	Active     *bool   `json:"active"`
}

func (c *StaffNotificationRecipient) ToJSON() string {
	return ModelToJson(c)
}

func (ce *StaffNotificationRecipient) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"staff_notification_recipient.is_valid.%s.app_error",
		"staff_notification_recipient_id=",
		"CustomerEvent.IsValid",
	)
	if !IsValidId(ce.Id) {
		return outer("id", nil)
	}
	if ce.UserID != nil && !IsValidId(*ce.UserID) {
		return outer("usder_id", &ce.Id)
	}
	if ce.StaffEmail != nil && !IsValidEmail(*ce.StaffEmail) {
		return outer("staff_email", &ce.Id)
	}

	return nil
}

func (c *StaffNotificationRecipient) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
	if c.Active == nil {
		c.Active = NewBool(true)
	}
	if c.StaffEmail != nil {
		c.StaffEmail = NewString(NormalizeEmail(*c.StaffEmail))
	}
}