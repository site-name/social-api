package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) HomepageEvents(ctx context.Context, args GraphqlParams) (*OrderEventCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

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

	return lo.Map(taxTypes, func(item *model_helper.TaxType, _ int) *TaxType {
		return &TaxType{
			Description: &item.Descriptiton,
			TaxCode:     &item.Code,
		}
	}), nil
}
