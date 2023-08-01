package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitename/sitename/model"
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
	err := embedCtx.App.Srv().Store.DiscountVoucher().ToggleVoucherRelations(nil, model.Vouchers{newVoucher}, args.Input.Collections, args.Input.Products, args.Input.Variants, args.Input.Categories, false)
	if err != nil {
		return nil, model.NewAppError("VoucherCreate", "app.discount.save_voucher_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &VoucherCreate{
		Voucher: systemVoucherToGraphqlVoucher(newVoucher),
	}, nil
}

// NOTE: Refer to ./schemas/voucher.graphqls for details on directives used
func (r *Resolver) VoucherDelete(ctx context.Context, args struct{ Id string }) (*VoucherDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherBulkDelete(ctx context.Context, args struct{ Ids []string }) (*VoucherBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherUpdate(ctx context.Context, args struct {
	Id    string
	Input VoucherInput
}) (*VoucherUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherCataloguesAdd(ctx context.Context, args struct {
	Id    string
	Input CatalogueInput
}) (*VoucherAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherCataloguesRemove(ctx context.Context, args struct {
	Id    string
	Input CatalogueInput
}) (*VoucherRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherTranslate(ctx context.Context, args struct {
	Id           string
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*VoucherTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input VoucherChannelListingInput
}) (*VoucherChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Voucher(ctx context.Context, args struct {
	Id      string
	Channel *string
}) (*Voucher, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Vouchers(ctx context.Context, args struct {
	Filter  *VoucherFilterInput
	SortBy  *VoucherSortingInput
	Query   *string
	Channel *string
	GraphqlParams
}) (*VoucherCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
