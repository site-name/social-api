package gqlmodel

import (
	"time"

	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/modules/util"
)

// original implementation
//
// type GiftCard struct {
// 	Code            string           `json:"code"`
// 	IsActive        bool             `json:"isActive"`
// 	ExpiryDate      *time.Time       `json:"expiryDate"`
// 	Tag             *string          `json:"tag"`
// 	Created         time.Time        `json:"created"`
// 	LastUsedOn      *time.Time       `json:"lastUsedOn"`
// 	InitialBalance  *Money           `json:"initialBalance"`
// 	CurrentBalance  *Money           `json:"currentBalance"`
// 	ID              string           `json:"id"`
// 	PrivateMetadata []*MetadataItem  `json:"privateMetadata"`
// 	Metadata        []*MetadataItem  `json:"metadata"`
// 	DisplayCode     string           `json:"displayCode"`
// 	CreatedBy       *User            `json:"createdBy"`
// 	UsedBy          *User            `json:"usedBy"`
// 	CreatedByEmail  *string          `json:"createdByEmail"`
// 	UsedByEmail     *string          `json:"usedByEmail"`
// 	App             *App             `json:"app"`
// 	Product         *Product         `json:"product"`
// 	Events          []*GiftCardEvent `json:"events"`
// 	BoughtInChannel *string          `json:"boughtInChannel"`
// }

// func (GiftCard) IsNode()               {}
// func (GiftCard) IsObjectWithMetadata() {}

type GiftCard struct {
	ID              string          `json:"id"`
	Code            string          `json:"code"`
	IsActive        bool            `json:"isActive"`
	ExpiryDate      *time.Time      `json:"expiryDate"`
	Tag             *string         `json:"tag"`
	Created         time.Time       `json:"created"`
	LastUsedOn      *time.Time      `json:"lastUsedOn"`
	InitialBalance  *Money          `json:"initialBalance"`
	CurrentBalance  *Money          `json:"currentBalance"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	DisplayCode     string          `json:"displayCode"`
	CreatedByID     *string         `json:"createdBy"`
	UsedByID        *string         `json:"usedBy"`
	CreatedByEmail  *string         `json:"createdByEmail"`
	UsedByEmail     *string         `json:"usedByEmail"`
	ProductID       *string         `json:"product"`
	EventIDs        []string        `json:"events"`
	BoughtInChannel *string         `json:"boughtInChannel"`
}

func (GiftCard) IsNode()               {}
func (GiftCard) IsObjectWithMetadata() {}

func SystemGiftcardToGraphqlGiftcard(g *giftcard.GiftCard) *GiftCard {
	if g == nil {
		return nil
	}

	res := &GiftCard{
		ID:              g.Id,
		Code:            g.Code,
		IsActive:        *g.IsActive,
		ExpiryDate:      g.ExpiryDate,
		Tag:             g.Tag,
		Created:         util.TimeFromMillis(g.CreateAt),
		Metadata:        MapToGraphqlMetaDataItems(g.Metadata),
		PrivateMetadata: MapToGraphqlMetaDataItems(g.PrivateMetadata),
		DisplayCode:     g.DisplayCode(),
		CreatedByID:     g.CreatedByID,
		UsedByID:        g.UsedByID,
		CreatedByEmail:  g.CreatedByEmail,
		UsedByEmail:     g.UsedByEmail,
		ProductID:       g.ProductID,
	}
	if g.LastUsedOn != nil {
		res.LastUsedOn = util.TimePointerFromMillis(*g.LastUsedOn)
	}
	if g.InitialBalanceAmount != nil {
		res.InitialBalance = new(Money)
		res.InitialBalance.Amount, _ = g.InitialBalanceAmount.Float64()
		res.InitialBalance.Currency = g.Currency
	}
	if g.CurrentBalanceAmount != nil {
		res.CurrentBalance = new(Money)
		res.CurrentBalance.Amount, _ = g.CurrentBalanceAmount.Float64()
		res.CurrentBalance.Currency = g.Currency
	}

	return res
}

func SystemGiftcardsToGraphqlGiftcards(gs giftcard.Giftcards) []*GiftCard {
	res := []*GiftCard{}
	for _, g := range gs {
		if g != nil {
			res = append(res, SystemGiftcardToGraphqlGiftcard(g))
		}
	}
	return res
}
