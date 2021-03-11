package models

import (
	"errors"
	"net/mail"
)

var (
	// ErrEmailAddressNotExist email address not exist
	ErrEmailAddressNotExist = errors.New("Email address does not exist")
)

// EmailAddress is the list of all email addresses of a user. Can contain the
// primary email address, but is not obligatory.
type EmailAddress struct {
	ID          int64  `xorm:"pk autoincr"`
	UID         int64  `xorm:"INDEX NOT NULL"`
	Email       string `xorm:"UNIQUE NOT NULL"`
	IsActivated bool
	IsPrimary   bool `xorm:"-"`
}

// ValidateEmail check if email is a allowed address
func ValidateEmail(email string) error {
	if len(email) == 0 {
		return nil
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ErrEmailInvalid{email}
	}

	// TODO: add an email allow/block list

	return nil
}

func isEmailUsed(e Engine, email string) (bool, error) {
	if len(email) == 0 {
		return true, nil
	}

	return e.Where("email=?", email).Get(&EmailAddress{})
}
