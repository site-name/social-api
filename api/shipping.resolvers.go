package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) ShippingMethodChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingMethodChannelListingInput
}) (*ShippingMethodChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceCreate(ctx context.Context, args struct{ Input ShippingPriceInput }) (*ShippingPriceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingPriceDelete(ctx context.Context, args struct{ Id string }) (*ShippingPriceDelete, error) {
	// validate params
	args.Id = decodeBase64String(args.Id)
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingPriceDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	shippingMethod, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodByOption(&model.ShippingMethodFilterOption{
		Id:                        squirrel.Eq{store.ShippingMethodTableName + ".Id": args.Id},
		SelectRelatedShippingZone: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	err := embedCtx.App.Srv().Store.ShippingMethod().Delete(nil, args.Id)
	if err != nil {
		return nil, model.NewAppError("ShippingPriceDelete", "app.shipping.delete_shipping_method.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &ShippingPriceDelete{
		ShippingMethod: SystemShippingMethodToGraphqlShippingMethod(shippingMethod),
		ShippingZone:   SystemShippingZoneToGraphqlShippingZone(shippingMethod.GetShippingZone()),
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingPriceBulkDelete(ctx context.Context, args struct{ Ids []string }) (*ShippingPriceBulkDelete, error) {
	// validate params
	args.Ids = decodeBase64Strings(args.Ids...)
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("ShippingPriceBulkDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	err := embedCtx.App.Srv().Store.ShippingMethod().Delete(nil, args.Ids...)
	if err != nil {
		return nil, model.NewAppError("ShippingPriceDelete", "app.shipping.delete_shipping_methods.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &ShippingPriceBulkDelete{
		Count: int32(len(args.Ids)),
	}, nil
}

func (r *Resolver) ShippingPriceUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingPriceInput
}) (*ShippingPriceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceTranslate(ctx context.Context, args struct {
	Id           string
	Input        ShippingPriceTranslationInput
	LanguageCode LanguageCodeEnum
}) (*ShippingPriceTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceExcludeProducts(ctx context.Context, args struct {
	Id    string
	Input ShippingPriceExcludeProductsInput
}) (*ShippingPriceExcludeProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingPriceRemoveProductFromExclude(ctx context.Context, args struct {
	Id       string
	Products []string
}) (*ShippingPriceRemoveProductFromExclude, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneCreate(ctx context.Context, args struct {
	Input ShippingZoneCreateInput
}) (*ShippingZoneCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneDelete(ctx context.Context, args struct{ Id string }) (*ShippingZoneDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneBulkDelete(ctx context.Context, args struct{ Ids []string }) (*ShippingZoneBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingZoneUpdateInput
}) (*ShippingZoneUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZone(ctx context.Context, args struct {
	Id      string
	Channel *string
}) (*ShippingZone, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZones(ctx context.Context, args struct {
	Filter  *ShippingZoneFilterInput
	Channel *string
	GraphqlParams
}) (*ShippingZoneCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
