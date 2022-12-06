package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

// -------------------------- Checkout ----------------------------------

type Checkout struct {
	ID                     string           `json:"id"`
	Created                DateTime         `json:"created"`
	LastChange             DateTime         `json:"lastChange"`
	Note                   string           `json:"note"`
	Discount               *Money           `json:"discount"`
	DiscountName           *string          `json:"discountName"`
	TranslatedDiscountName *string          `json:"translatedDiscountName"`
	VoucherCode            *string          `json:"voucherCode"`
	PrivateMetadata        []*MetadataItem  `json:"privateMetadata"`
	Metadata               []*MetadataItem  `json:"metadata"`
	Email                  string           `json:"email"`
	Quantity               int32            `json:"quantity"`
	Token                  string           `json:"token"`
	LanguageCode           LanguageCodeEnum `json:"languageCode"`
	shippingAddressID      *string
	billingAddressID       *string
	userID                 *string

	// ShippingPrice          *TaxedMoney      `json:"shippingPrice"`
	// SubtotalPrice          *TaxedMoney      `json:"subtotalPrice"`
	// TotalPrice             *TaxedMoney      `json:"totalPrice"``
	// IsShippingRequired        bool              `json:"isShippingRequired"`
	// User                      *User             `json:"user"`
	// Channel                   *Channel          `json:"channel"`
	// BillingAddress            *Address          `json:"billingAddress"`
	// ShippingAddress           *Address          `json:"shippingAddress"`
	// GiftCards                 []*GiftCard       `json:"giftCards"`
	// AvailableShippingMethods  []*ShippingMethod `json:"availableShippingMethods"`
	// AvailableCollectionPoints []*Warehouse      `json:"availableCollectionPoints"`
	// AvailablePaymentGateways  []*PaymentGateway `json:"availablePaymentGateways"`
	// Lines                     []*CheckoutLine   `json:"lines"`
	// DeliveryMethod     			 DeliveryMethod    `json:"deliveryMethod"`
}

func SystemCheckoutToGraphqlCheckout(ckout *model.Checkout) *Checkout {
	if ckout == nil {
		return nil
	}

	ckout.PopulateNonDbFields()

	res := &Checkout{
		ID:                     ckout.Token,
		Created:                DateTime{model.GetTimeForMillis(ckout.CreateAt)},
		LastChange:             DateTime{model.GetTimeForMillis(ckout.UpdateAt)},
		Note:                   ckout.Note,
		Discount:               SystemMoneyToGraphqlMoney(ckout.Discount),
		DiscountName:           ckout.DiscountName,
		TranslatedDiscountName: ckout.TranslatedDiscountName,
		VoucherCode:            ckout.VoucherCode,
		PrivateMetadata:        MetadataToSlice(ckout.PrivateMetadata),
		Metadata:               MetadataToSlice(ckout.Metadata),
		Email:                  ckout.Email,
		Quantity:               int32(ckout.Quantity),
		Token:                  ckout.Token,
		LanguageCode:           SystemLanguageToGraphqlLanguageCodeEnum(ckout.LanguageCode),

		shippingAddressID: ckout.ShippingAddressID,
		billingAddressID:  ckout.BillingAddressID,
		userID:            ckout.UserID,
	}
	return res
}

func (c *Checkout) User(ctx context.Context) (*User, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if (c.userID != nil && embedCtx.AppContext.Session().UserId == *c.userID) ||
		embedCtx.App.Srv().
			AccountService().
			SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageUsers) {

		if c.userID != nil {
			user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, *c.userID)
			if appErr != nil {
				return nil, appErr
			}

			return SystemUserToGraphqlUser(user), nil
		}

		return nil, nil
	}

	return nil, model.NewAppError("checkout.User", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}

func (c *Checkout) IsShippingRequired(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (c *Checkout) Channel(ctx context.Context) (*Channel, error) {
	panic("not implemented")
}

func (c *Checkout) BillingAddress(ctx context.Context) (*Address, error) {
	if c.billingAddressID == nil {
		return nil, nil
	}

	return dataloaders.addressesByIDs.Load(ctx, *c.billingAddressID)()
}

func (c *Checkout) ShippingAddress(ctx context.Context) (*Address, error) {
	if c.shippingAddressID == nil {
		return nil, nil
	}

	return dataloaders.addressesByIDs.Load(ctx, *c.shippingAddressID)()
}

func (c *Checkout) GiftCards(ctx context.Context) ([]*GiftCard, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	giftcards, appErr := embedCtx.App.Srv().GiftcardService().GiftcardsByOption(nil, &model.GiftCardFilterOption{
		CheckoutToken: squirrel.Eq{store.GiftcardCheckoutTableName + ".CheckoutID": c.Token},
	})
	if appErr != nil {
		return nil, appErr
	}

	res := []*GiftCard{}
	for _, gc := range giftcards {
		if gc != nil {
			res = append(res, SystemGiftcardToGraphqlGiftcard(gc))
		}
	}

	return res, nil
}

func (c *Checkout) AvailableShippingMethods(ctx context.Context) ([]*ShippingMethod, error) {
	panic("not implemented")
}

func (c *Checkout) AvailableCollectionPoints(ctx context.Context) ([]*Warehouse, error) {
	panic("not implemented")
}

func (c *Checkout) AvailablePaymentGateways(ctx context.Context) ([]*PaymentGateway, error) {
	panic("not implemented")
}

func (c *Checkout) Lines(ctx context.Context) ([]*CheckoutLine, error) {
	panic("not implemented")
}

func (c *Checkout) DeliveryMethod(ctx context.Context) (any, error) {
	panic("not implemented")
}

// NOTE:
// keys are strings that have format uuid__uuid.
// The first uuid part is userID, send is channelID
func graphqlCheckoutsByUserAndChannelLoader(ctx context.Context, keys []string) []*dataloader.Result[[]*Checkout] {
	var (
		appErr     *model.AppError
		res        []*dataloader.Result[[]*Checkout]
		userIDs    model.AnyArray[string]
		channelIDs model.AnyArray[string]
		checkouts  []*model.Checkout
		// checkoutsMap has keys are each items of given param keys
		checkoutsMap = map[string][]*Checkout{}
	)

	for _, item := range keys {
		sepIndex := strings.Index(item, "__")
		if sepIndex == -1 {
			continue
		}

		userIDs = append(userIDs, item[:sepIndex])
		channelIDs = append(channelIDs, item[sepIndex+2:])
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}
	checkouts, appErr = embedCtx.
		App.
		Srv().
		CheckoutService().
		CheckoutsByOption(&model.CheckoutFilterOption{
			ChannelIsActive: model.NewBool(true),
			UserID:          squirrel.Eq{store.CheckoutTableName + ".UserID": userIDs},
			ChannelID:       squirrel.Eq{store.CheckoutTableName + ".ChannelID": channelIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, checkout := range checkouts {
		if checkout.UserID != nil {
			key := *checkout.UserID + "__" + checkout.ChannelID
			checkoutsMap[key] = append(checkoutsMap[key], SystemCheckoutToGraphqlCheckout(checkout))
		}
	}

	for _, key := range keys {
		res = append(res, &dataloader.Result[[]*Checkout]{Data: checkoutsMap[key]})
	}
	return res

errorLabel:
	for range keys {
		res = append(res, &dataloader.Result[[]*Checkout]{Error: err})
	}
	return res
}

func graphqlCheckoutByUserLoader(ctx context.Context, userIDs []string) []*dataloader.Result[[]*Checkout] {
	var (
		appErr       *model.AppError
		checkouts    []*model.Checkout
		res          []*dataloader.Result[[]*Checkout]
		checkoutsMap = map[string][]*Checkout{} // keys are user ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	checkouts, appErr = embedCtx.
		App.
		Srv().
		CheckoutService().
		CheckoutsByOption(&model.CheckoutFilterOption{
			ChannelIsActive: model.NewBool(true),
			UserID:          squirrel.Eq{store.CheckoutTableName + ".UserID": userIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, checkout := range checkouts {
		if checkout.UserID != nil {
			checkoutsMap[*checkout.UserID] = append(checkoutsMap[*checkout.UserID], SystemCheckoutToGraphqlCheckout(checkout))
		}
	}

	for _, key := range userIDs {
		res = append(res, &dataloader.Result[[]*Checkout]{Data: checkoutsMap[key]})
	}
	return res

errorLabel:
	for range userIDs {
		res = append(res, &dataloader.Result[[]*Checkout]{Error: err})
	}
	return res
}
