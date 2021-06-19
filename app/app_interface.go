package app

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
)

type AppInterface interface {
	GiftCard() GiftCardInterface
}

type GiftCardInterface interface {
	GetAllByUserId(userID string) ([]*giftcard.GiftCard, *model.AppError)
	GetAll() ([]*giftcard.GiftCard, *model.AppError)
}
