package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) FileUpload(ctx context.Context, args struct{ File graphql.Upload }) (*FileUpload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalNotificationTrigger(ctx context.Context, args struct {
	Channel  string
	Input    ExternalNotificationTriggerInput
	PluginID *string
}) (*ExternalNotificationTrigger, error) {
	panic(fmt.Errorf("not implemented"))
}

type TaxType struct {
	Description *string `json:"description"`
	TaxCode     *string `json:"taxCode"`
}

func (r *Resolver) TaxTypes(ctx context.Context) ([]*TaxType, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	taxTypes, appErr := pluginMng.GetTaxRateTypeChoices()
	if appErr != nil {
		return nil, appErr
	}

	return lo.Map(taxTypes, func(item *model.TaxType, _ int) *TaxType {
		return &TaxType{
			Description: &item.Descriptiton,
			TaxCode:     &item.Code,
		}
	}), nil
}

// NOTE: please refer to ./schemas/order.graphqls for details on directives used.
func (r *Resolver) OrderSettingsUpdate(ctx context.Context, args struct {
	Input OrderSettingsUpdateInput
}) (*OrderSettingsUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentConfig := embedCtx.App.Config().ShopSettings

	embedCtx.App.UpdateConfig(func(c *model.Config) {
		if val := args.Input.AutomaticallyConfirmAllNewOrders; val != nil &&
			*currentConfig.AutomaticallyConfirmAllNewOrders != *val {
			*c.ShopSettings.AutomaticallyConfirmAllNewOrders = *val
		}
		if val := args.Input.AutomaticallyFulfillNonShippableGiftCard; val != nil &&
			*currentConfig.AutomaticallyFulfillNonShippableGiftcard != *val {
			*c.ShopSettings.AutomaticallyFulfillNonShippableGiftcard = *val
		}
	})

	updatedShopSettings := embedCtx.App.Config().ShopSettings

	return &OrderSettingsUpdate{
		OrderSettings: &OrderSettings{
			AutomaticallyConfirmAllNewOrders:         *updatedShopSettings.AutomaticallyConfirmAllNewOrders,
			AutomaticallyFulfillNonShippableGiftCard: *updatedShopSettings.AutomaticallyFulfillNonShippableGiftcard,
		},
	}, nil
}

// NOTE: please refer to ./schemas/order.graphqls for details on directives used.
func (r *Resolver) GiftCardSettingsUpdate(ctx context.Context, args struct {
	Input GiftCardSettingsUpdateInput
}) (*GiftCardSettingsUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSettings := embedCtx.App.Config().ShopSettings

	embedCtx.App.UpdateConfig(func(c *model.Config) {
		var (
			expiryType   = args.Input.ExpiryType
			expiryPeriod = args.Input.ExpiryPeriod
		)

		if expiryPeriod != nil {
			if *currentSettings.GiftcardExpiryPeriodType != expiryPeriod.Type {
				*c.ShopSettings.GiftcardExpiryPeriodType = expiryPeriod.Type
			}
			if *currentSettings.GiftcardExpiryPeriod != int(expiryPeriod.Amount) {
				*c.ShopSettings.GiftcardExpiryPeriod = int(expiryPeriod.Amount)
			}
		}

		if expiryType != nil && *expiryType != *currentSettings.GiftcardExpiryType {
			*c.ShopSettings.GiftcardExpiryType = *expiryType
		}
	})

	updatedSettings := embedCtx.App.Config().ShopSettings
	return &GiftCardSettingsUpdate{
		GiftCardSettings: &GiftCardSettings{
			ExpiryType: *updatedSettings.GiftcardExpiryType,
			ExpiryPeriod: &TimePeriod{
				Amount: int32(*updatedSettings.GiftcardExpiryPeriod),
				Type:   *updatedSettings.GiftcardExpiryPeriodType,
			},
		},
	}, nil
}
