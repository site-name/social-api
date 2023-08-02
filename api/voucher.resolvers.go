package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/voucher.graphqls for details on directives used
func (r *Resolver) VoucherCreate(ctx context.Context, args struct{ Input VoucherInput }) (*VoucherCreate, error) {
	// validate params
	appErr := args.Input.Validate("VoucherCreate")
	if appErr != nil {
		return nil, appErr
	}

	voucher := &model.Voucher{}
	args.Input.PatchVoucher(voucher)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	newVoucher, appErr := embedCtx.App.Srv().DiscountService().UpsertVoucher(voucher)
	if appErr != nil {
		return nil, appErr
	}

	// save relations
	appErr = embedCtx.App.Srv().DiscountService().ToggleVoucherRelations(nil, model.Vouchers{newVoucher}, args.Input.Collections, args.Input.Products, args.Input.Variants, args.Input.Categories, false)
	if appErr != nil {
		return nil, appErr
	}

	return &VoucherCreate{
		Voucher: systemVoucherToGraphqlVoucher(newVoucher),
	}, nil
}

// NOTE: Refer to ./schemas/voucher.graphqls for details on directives used
func (r *Resolver) VoucherDelete(ctx context.Context, args struct{ Id string }) (*VoucherDelete, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("VoucherDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid voucher id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	err := embedCtx.App.Srv().Store.DiscountVoucher().Delete(nil, []string{args.Id})
	if err != nil {
		return nil, model.NewAppError("VoucherDelete", "app.discount.delete_voucher.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &VoucherDelete{
		Voucher: &Voucher{ID: args.Id},
	}, nil
}

// NOTE: Refer to ./schemas/voucher.graphqls for details on directives used
func (r *Resolver) VoucherBulkDelete(ctx context.Context, args struct{ Ids []string }) (*VoucherBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/voucher.graphqls for details on directives used
func (r *Resolver) VoucherUpdate(ctx context.Context, args struct {
	Id    string
	Input VoucherInput
}) (*VoucherUpdate, error) {
	// validate params
	appErr := args.Input.Validate("VoucherUpdate")
	if appErr != nil {
		return nil, appErr
	}
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("VoucherUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid voucher id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	voucher, appErr := embedCtx.App.Srv().DiscountService().VoucherById(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	// update voucher in database
	args.Input.PatchVoucher(voucher)

	voucher, appErr = embedCtx.App.Srv().DiscountService().UpsertVoucher(voucher)
	if appErr != nil {
		return nil, appErr
	}

	// save relations
	appErr = embedCtx.App.Srv().DiscountService().ToggleVoucherRelations(nil, model.Vouchers{voucher}, args.Input.Collections, args.Input.Products, args.Input.Variants, args.Input.Categories, false)
	if appErr != nil {
		return nil, appErr
	}

	return &VoucherUpdate{
		Voucher: systemVoucherToGraphqlVoucher(voucher),
	}, nil
}

// NOTE: Refer to ./schemas/voucher.graphqls for details on directives used
func (r *Resolver) VoucherCataloguesAdd(ctx context.Context, args struct {
	Id    string
	Input CatalogueInput
}) (*VoucherAddCatalogues, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("VoucherCataloguesAdd", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid voucher id", http.StatusBadRequest)
	}
	appErr := args.Input.Validate("VoucherCataloguesAdd")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// NOTE: only products that have variants can be added to voucher.
	// So we have to verify that every given products have variant(s)
	if len(args.Input.Products) > 0 {
		productsWithNovariants, appErr := embedCtx.App.Srv().ProductService().ProductsByOption(&model.ProductFilterOption{
			Conditions:           squirrel.Eq{model.ProductTableName + ".Id": args.Input.Products},
			HasNoProductVariants: true,
		})
		if appErr != nil {
			return nil, appErr
		}
		if len(productsWithNovariants) > 0 {
			return nil, model.NewAppError("VoucherCataloguesAdd", "app.discount.add_products_with_no_variants_to_voucher.app_error", nil, "cant add products that have no variants to sale", http.StatusNotAcceptable)
		}
	}

	// save relations
	appErr = embedCtx.App.Srv().
		DiscountService().
		ToggleVoucherRelations(nil, model.Vouchers{{Id: args.Id}}, args.Input.Products, []string{}, args.Input.Categories, args.Input.Collections, false)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().ProductService().UpdateProductsDiscountedPricesOfCatalogues(nil, args.Input.Products, args.Input.Categories, args.Input.Collections, nil)
	if appErr != nil {
		return nil, appErr
	}

	return &VoucherAddCatalogues{
		Voucher: &Voucher{ID: args.Id},
	}, nil
}

// NOTE: Refer to ./schemas/voucher.graphqls for details on directives used
func (r *Resolver) VoucherCataloguesRemove(ctx context.Context, args struct {
	Id    string
	Input CatalogueInput
}) (*VoucherRemoveCatalogues, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("VoucherCataloguesRemove", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid voucher id", http.StatusBadRequest)
	}
	appErr := args.Input.Validate("VoucherCataloguesRemove")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// save relations
	appErr = embedCtx.App.Srv().
		DiscountService().
		ToggleVoucherRelations(nil, model.Vouchers{{Id: args.Id}}, args.Input.Products, []string{}, args.Input.Categories, args.Input.Collections, true)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().ProductService().UpdateProductsDiscountedPricesOfCatalogues(nil, args.Input.Products, args.Input.Categories, args.Input.Collections, nil)
	if appErr != nil {
		return nil, appErr
	}

	return &VoucherRemoveCatalogues{
		Voucher: &Voucher{ID: args.Id},
	}, nil
}

// NOTE: Refer to ./schemas/voucher.graphqls for details on directives used
func (r *Resolver) VoucherTranslate(ctx context.Context, args struct {
	Id           string
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*VoucherTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/voucher.graphqls for details on directives used
func (r *Resolver) VoucherChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input VoucherChannelListingInput
}) (*VoucherChannelListingUpdate, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("VoucherChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid voucher id", http.StatusBadRequest)
	}
	appErr := args.Input.Validate()
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	voucher, appErr := embedCtx.App.Srv().DiscountService().VoucherById(args.Id)
	if appErr != nil {
		return nil, appErr
	}
	channelsAssignedToVoucher, appErr := embedCtx.App.Srv().ChannelService().ChannelsByOption(&model.ChannelFilterOption{
		VoucherChannelListing_VoucherID: squirrel.Expr(model.VoucherChannelListingTableName+".VoucherID = ?", args.Id),
	})
	if appErr != nil {
		return nil, appErr
	}
	// keys are channel ids
	channelsAssignedToVoucherMap := lo.SliceToMap(channelsAssignedToVoucher, func(c *model.Channel) (string, *model.Channel) { return c.Id, c })

	// clean discount values
	for _, addChannelObj := range args.Input.AddChannels {
		if addChannelObj == nil {
			continue
		}

		channelToAdd := channelsAssignedToVoucherMap[addChannelObj.ChannelID]
		channelIsNOTAssignedToVoucher := channelToAdd == nil
		if channelIsNOTAssignedToVoucher && addChannelObj.DiscountValue == nil {
			return nil, model.NewAppError("VoucherChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "discountValue"}, "please provide discount value to assign channel to voucher", http.StatusBadRequest)
		}

		switch voucher.DiscountValueType {
		case model.DISCOUNT_VALUE_TYPE_FIXED:
			// check if discount value is provided and has valid precision.
			// If not valid, correct it, not raise error
			if addChannelObj.DiscountValue == nil {
				return nil, model.NewAppError("VoucherChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "discountValue"}, "please provide discount value to assign channel to voucher", http.StatusBadRequest)
			}
			precision, _ := goprices.GetCurrencyPrecision(channelToAdd.Currency)
			roundedValue := decimal.Decimal(*addChannelObj.DiscountValue).RoundUp(int32(precision))
			*addChannelObj.DiscountValue = *(*PositiveDecimal)(unsafe.Pointer(&roundedValue))

		case model.DISCOUNT_VALUE_TYPE_PERCENTAGE:
			// discount can't > 100
			if decimal.Decimal(*addChannelObj.DiscountValue).GreaterThan(decimal.NewFromInt(100)) {
				return nil, model.NewAppError("VoucherChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "discountValue"}, "discount value can greater than 100", http.StatusBadRequest)
			}
		}

		// validate min spent amount
		if addChannelObj.MinAmountSpent != nil {
			precision, _ := goprices.GetCurrencyPrecision(channelToAdd.Currency)
			roundedValue := decimal.Decimal(*addChannelObj.MinAmountSpent).RoundUp(int32(precision))
			*addChannelObj.MinAmountSpent = *(*PositiveDecimal)(unsafe.Pointer(&roundedValue))
		}
	}

	// init transaction
	tran := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tran.Error != nil {
		return nil, model.NewAppError("VoucherChannelListingUpdate", model.ErrorCreatingTransactionErrorID, nil, tran.Error.Error(), http.StatusInternalServerError)
	}
	defer tran.Rollback()

	// perform database mutation:
	listingsToAdd := lo.Map(args.Input.AddChannels, func(item *VoucherChannelListingAddInput, _ int) *model.VoucherChannelListing {
		return &model.VoucherChannelListing{
			VoucherID:      args.Id,
			ChannelID:      item.ChannelID,
			DiscountValue:  (*decimal.Decimal)(unsafe.Pointer(item.DiscountValue)),
			MinSpentAmount: (*decimal.Decimal)(unsafe.Pointer(item.MinAmountSpent)),
		}
	})
	_, err := embedCtx.App.Srv().Store.VoucherChannelListing().Upsert(tran, listingsToAdd)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("VoucherChannelListingUpdate", "app.discount.add_voucher_channel_listings.app_error", nil, err.Error(), statusCode)
	}

	err = embedCtx.App.Srv().Store.VoucherChannelListing().Delete(tran, &model.VoucherChannelListingFilterOption{
		Conditions: squirrel.Eq{
			model.VoucherChannelListingTableName + ".VoucherID": args.Id,
			model.VoucherChannelListingTableName + ".ChannelID": args.Input.RemoveChannels,
		},
	})
	if err != nil {
		return nil, model.NewAppError("VoucherChannelListingUpdate", "app.discount.delete_voucher_channel_listings.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// commit transaction
	if err := tran.Commit().Error; err != nil {
		return nil, model.NewAppError("VoucherChannelListingUpdate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &VoucherChannelListingUpdate{
		Voucher: systemVoucherToGraphqlVoucher(voucher),
	}, nil
}

func (r *Resolver) Voucher(ctx context.Context, args struct {
	Id      string
	Channel *string // NOTE: this is channel id
}) (*Voucher, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("Voucher", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid voucher id", http.StatusBadRequest)
	}

	voucher, err := VoucherByIDLoader.Load(ctx, args.Id)()
	if err != nil {
		return nil, err
	}

	return systemVoucherToGraphqlVoucher(voucher), nil
}

func (r *Resolver) Vouchers(ctx context.Context, args struct {
	Filter  *VoucherFilterInput
	SortBy  *VoucherSortingInput
	Channel *string
	GraphqlParams
}) (*VoucherCountableConnection, error) {
	// validate params
	appErr := args.GraphqlParams.Validate("Vouchers")
	if appErr != nil {
		return nil, appErr
	}

	var voucherFilter = &model.VoucherFilterOption{}
	if args.Filter != nil {
		voucherFilter, appErr = args.Filter.Parse()
		if appErr != nil {
			return nil, appErr
		}
	}

	if args.Channel != nil && model.IsValidId(*args.Channel) {
		voucherFilter.VoucherChannelListing_ChannelSlug = squirrel.Expr(model.ChannelTableName+".Slug = ?", *args.Channel)
	}

}
