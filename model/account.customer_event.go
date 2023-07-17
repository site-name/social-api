package model

import (
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

type CustomerEventType string

func (t CustomerEventType) IsValid() bool {
	return CustomerEventTypes[t]
}

// some available types for CustomerEvent's Type attribute
const (
	CUSTOMER_EVENT_TYPE_ACCOUNT_CREATED          CustomerEventType = "account_created"
	CUSTOMER_EVENT_TYPE_PASSWORD_RESET_LINK_SENT CustomerEventType = "password_reset_link_sent"
	CUSTOMER_EVENT_TYPE_PASSWORD_RESET           CustomerEventType = "password_reset"
	CUSTOMER_EVENT_TYPE_PASSWORD_CHANGED         CustomerEventType = "password_changed"
	CUSTOMER_EVENT_TYPE_EMAIL_CHANGE_REQUEST     CustomerEventType = "email_changed_request"
	CUSTOMER_EVENT_TYPE_EMAIL_CHANGED            CustomerEventType = "email_changed"
	CUSTOMER_EVENT_TYPE_PLACED_ORDER             CustomerEventType = "placed_order"            // created an order
	CUSTOMER_EVENT_TYPE_NOTE_ADDED_TO_ORDER      CustomerEventType = "note_added_to_order"     // added a note to one of their orders
	CUSTOMER_EVENT_TYPE_DIGITAL_LINK_DOWNLOADED  CustomerEventType = "digital_link_downloaded" // downloaded a digital good
	CUSTOMER_EVENT_TYPE_CUSTOMER_DELETED         CustomerEventType = "customer_deleted"        // staff user deleted a customer
	CUSTOMER_EVENT_TYPE_EMAIL_ASSIGNED           CustomerEventType = "email_assigned"          // the staff user assigned a email to the customer
	CUSTOMER_EVENT_TYPE_NAME_ASSIGNED            CustomerEventType = "name_assigned"           // the staff user added set a name to the customer
	CUSTOMER_EVENT_TYPE_CUSTOMER_NOTE_ADDED      CustomerEventType = "note_added"              // the staff user added a note to the customer
)

var CustomerEventTypes = map[CustomerEventType]bool{
	CUSTOMER_EVENT_TYPE_ACCOUNT_CREATED:          true,
	CUSTOMER_EVENT_TYPE_PASSWORD_RESET_LINK_SENT: true,
	CUSTOMER_EVENT_TYPE_PASSWORD_RESET:           true,
	CUSTOMER_EVENT_TYPE_PASSWORD_CHANGED:         true,
	CUSTOMER_EVENT_TYPE_EMAIL_CHANGE_REQUEST:     true,
	CUSTOMER_EVENT_TYPE_EMAIL_CHANGED:            true,
	CUSTOMER_EVENT_TYPE_PLACED_ORDER:             true,
	CUSTOMER_EVENT_TYPE_NOTE_ADDED_TO_ORDER:      true,
	CUSTOMER_EVENT_TYPE_DIGITAL_LINK_DOWNLOADED:  true,
	CUSTOMER_EVENT_TYPE_CUSTOMER_DELETED:         true,
	CUSTOMER_EVENT_TYPE_EMAIL_ASSIGNED:           true,
	CUSTOMER_EVENT_TYPE_NAME_ASSIGNED:            true,
	CUSTOMER_EVENT_TYPE_CUSTOMER_NOTE_ADDED:      true,
}

type CustomerEvent struct {
	Id      string            `json:"id" gorm:"primaryKey;type:uuid;defautl:gen_random_uuid();column:Id"`
	Date    int64             `json:"date" gorm:"type:bigint;autoCreateTime:milli;column:Date"`
	Type    CustomerEventType `json:"type" gorm:"type:varchar(255);column:Type"`
	OrderID *string           `json:"order_id" gorm:"type:uuid;index;column:OrderID"`
	UserID  *string           `json:"user_id" gorm:"type:uuid;index:customerevents_userid_key;column:UserID"`
	// To reduce number of type checking steps,
	// below are possible keys and their according values's Types you must follow
	//  "message": string
	//  "count": int
	//  "order_line_pk": string
	Parameters StringInterface `json:"parameters" gorm:"type:jsonb;column:Parameters"`
}

func (c *CustomerEvent) BeforeCreate(_ *gorm.DB) error {
	c.commonPre()
	return c.IsValid()
}

func (c *CustomerEvent) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	return c.IsValid()
}

func (*CustomerEvent) TableName() string {
	return CustomerEventTableName
}

func (ce *CustomerEvent) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.customer_event.is_valid.%s.app_error",
		"customer_event_id=",
		"CustomerEvent.IsValid",
	)
	if ce.Date == 0 {
		return outer("date", &ce.Id)
	}
	if ce.UserID != nil && !IsValidId(*ce.UserID) {
		return outer("usder_id", &ce.Id)
	}
	if ce.OrderID != nil && !IsValidId(*ce.OrderID) {
		return outer("order_id", &ce.Id)
	}
	if !CustomerEventTypes[ce.Type] {
		return outer("type", &ce.Id)
	}

	return nil
}

func (c *CustomerEvent) commonPre() {
	_, ok1 := c.Parameters["currency"]
	_, ok2 := c.Parameters["amount"]
	if ok1 && ok2 {
		c.Parameters["_type"] = "Money"
	}
}

type StaffNotificationRecipient struct {
	Id         string  `json:"id" gorm:"primaryKey;type:uuid;defautl:gen_random_uuid();column:Id"`
	UserID     *string `json:"user_id" gorm:"type:uuid;column:UserID;index:staffnotificationrecipients_userid_key"`
	StaffEmail *string `json:"staff_email" gorm:"uniqueIndex:staff_notification_recipients_staff_email_unique_key;column:StaffEmail"`
	Active     *bool   `json:"active" gorm:"default:true;column:Active"`
}

type StaffNotificationRecipientFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (ce *StaffNotificationRecipient) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.staff_notification_recipient.is_valid.%s.app_error",
		"staff_notification_recipient_id=",
		"CustomerEvent.IsValid",
	)
	if ce.UserID != nil && !IsValidId(*ce.UserID) {
		return outer("usder_id", &ce.Id)
	}
	if ce.StaffEmail != nil && !IsValidEmail(*ce.StaffEmail) {
		return outer("staff_email", &ce.Id)
	}

	return nil
}

func (c *StaffNotificationRecipient) BeforeCreate(_ *gorm.DB) error {
	c.commonPre()
	return c.IsValid()
}

func (c *StaffNotificationRecipient) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	return c.IsValid()
}

func (*StaffNotificationRecipient) TableName() string {
	return StaffNotificationRecipientTableName
}

func (c *StaffNotificationRecipient) commonPre() {
	if c.StaffEmail != nil {
		c.StaffEmail = NewPrimitive(NormalizeEmail(*c.StaffEmail))
	}
}
