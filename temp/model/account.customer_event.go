package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

type CustomerEventType string

func (t CustomerEventType) IsValid() bool {
	return CustomerEventTypes[t]
}

// some available types for CustomerEvent's Type attribute
const (
	CUSTOMER_EVENT_TYPE_ACCOUNT_CREATED          CustomerEventType = "ACCOUNT_CREATED"
	CUSTOMER_EVENT_TYPE_PASSWORD_RESET_LINK_SENT CustomerEventType = "PASSWORD_RESET_LINK_SENT"
	CUSTOMER_EVENT_TYPE_PASSWORD_RESET           CustomerEventType = "PASSWORD_RESET"
	CUSTOMER_EVENT_TYPE_PASSWORD_CHANGED         CustomerEventType = "PASSWORD_CHANGED"
	CUSTOMER_EVENT_TYPE_EMAIL_CHANGE_REQUEST     CustomerEventType = "EMAIL_CHANGED_REQUEST"
	CUSTOMER_EVENT_TYPE_EMAIL_CHANGED            CustomerEventType = "EMAIL_CHANGED"
	CUSTOMER_EVENT_TYPE_PLACED_ORDER             CustomerEventType = "PLACED_ORDER"            // created an order
	CUSTOMER_EVENT_TYPE_NOTE_ADDED_TO_ORDER      CustomerEventType = "NOTE_ADDED_TO_ORDER"     // added a note to one of their orders
	CUSTOMER_EVENT_TYPE_DIGITAL_LINK_DOWNLOADED  CustomerEventType = "DIGITAL_LINK_DOWNLOADED" // downloaded a digital good
	CUSTOMER_EVENT_TYPE_CUSTOMER_DELETED         CustomerEventType = "CUSTOMER_DELETED"        // staff user deleted a customer
	CUSTOMER_EVENT_TYPE_EMAIL_ASSIGNED           CustomerEventType = "EMAIL_ASSIGNED"          // the staff user assigned a email to the customer
	CUSTOMER_EVENT_TYPE_NAME_ASSIGNED            CustomerEventType = "NAME_ASSIGNED"           // the staff user added set a name to the customer
	CUSTOMER_EVENT_TYPE_CUSTOMER_NOTE_ADDED      CustomerEventType = "NOTE_ADDED"              // the staff user added a note to the customer
	CUSTOMER_EVENT_TYPE_ACCOUNT_ACTIVATED        CustomerEventType = "ACCOUNT_ACTIVATED"
	CUSTOMER_EVENT_TYPE_ACCOUNT_DEACTIVATED      CustomerEventType = "ACCOUNT_DEACTIVATED"
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
	CUSTOMER_EVENT_TYPE_ACCOUNT_ACTIVATED:        true,
	CUSTOMER_EVENT_TYPE_ACCOUNT_DEACTIVATED:      true,
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
	if ce.Date == 0 {
		return NewAppError("CustomerEvent.IsValid", "model.customer_event.is_valid.date.app_error", nil, "please provide valid date", http.StatusBadRequest)
	}
	if ce.UserID != nil && !IsValidId(*ce.UserID) {
		return NewAppError("CustomerEvent.IsValid", "model.customer_event.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if ce.OrderID != nil && !IsValidId(*ce.OrderID) {
		return NewAppError("CustomerEvent.IsValid", "model.customer_event.is_valid.order_id.app_error", nil, "please provide order id", http.StatusBadRequest)
	}
	if !CustomerEventTypes[ce.Type] {
		return NewAppError("CustomerEvent.IsValid", "model.customer_event.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
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
	StaffEmail *string `json:"staff_email" gorm:"unique:staff_notification_recipients_staff_email_key;column:StaffEmail"`
	Active     *bool   `json:"active" gorm:"default:true;column:Active"`
}

type StaffNotificationRecipientFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (ce *StaffNotificationRecipient) IsValid() *AppError {
	if ce.UserID != nil && !IsValidId(*ce.UserID) {
		return NewAppError("CustomerEvent.IsValid", "model.staff_notification_recipient.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if ce.StaffEmail != nil && !IsValidEmail(*ce.StaffEmail) {
		return NewAppError("CustomerEvent.IsValid", "model.staff_notification_recipient.is_valid.staff_email.app_error", nil, "please provide valid staff email", http.StatusBadRequest)
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
		c.StaffEmail = GetPointerOfValue(NormalizeEmail(*c.StaffEmail))
	}
}
