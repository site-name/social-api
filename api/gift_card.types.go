package api

import (
	"context"
	"net/http"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

func SystemGiftcardToGraphqlGiftcard(gc *model.GiftCard) *GiftCard {
	res := new(GiftCard)

	gc.PopulateNonDbFields()

	res.ID = gc.Id
	res.IsActive = *gc.IsActive
	res.Tag = gc.Tag
	res.DisplayCode = gc.DisplayCode()

	res.createdByEmail = gc.CreatedByEmail
	res.usedByEmail = gc.UsedByEmail
	res.code = gc.Code
	res.usedByID = gc.UsedByID
	res.createdByID = gc.CreatedByID

	if gc.ExpiryDate != nil {
		res.ExpiryDate = &Date{
			DateTime{util.StartOfDay(*gc.ExpiryDate)},
		}
	}
	res.Created = DateTime{util.TimeFromMillis(gc.CreateAt)}
	res.Tag = gc.Tag
	if gc.LastUsedOn != nil {
		res.LastUsedOn = &DateTime{util.TimeFromMillis(*gc.LastUsedOn)}
	}
	if gc.CurrentBalance != nil { // these fields may got populated above
		flt, _ := gc.CurrentBalance.Amount.Float64()
		res.CurrentBalance = &Money{
			Amount:   flt,
			Currency: gc.CurrentBalance.Currency,
		}
	}
	if gc.InitialBalance != nil { // these fields may got populated above
		flt, _ := gc.InitialBalance.Amount.Float64()
		res.InitialBalance = &Money{
			Amount:   flt,
			Currency: gc.InitialBalance.Currency,
		}
	}
	res.Metadata = MetadataToSlice[string](gc.Metadata)
	res.PrivateMetadata = MetadataToSlice[string](gc.PrivateMetadata)

	return res
}

func (g *GiftCard) CreatedByEmail(ctx context.Context) (string, error) {

}

func (g *GiftCard) UsedByEmail(ctx context.Context) (string, error) {

}

func (g *GiftCard) UsedBy(ctx context.Context) (*User, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if g.usedByID == nil {
		if embedCtx.App.Srv().
			AccountService().
			SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageUsers) {

			return SystemUserToGraphqlUser(&model.User{
				Id: embedCtx.AppContext.Session().UserId,
			}), nil
		}

		return nil, model.NewAppError("UsedBy", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
	}

	user, err := dataloaders.users.Load(ctx, *g.usedByID)()
	if err != nil {
		return nil, err
	}

}

func (g *GiftCard) CreatedBy(ctx context.Context) (*User, error) {

}

func (g *GiftCard) Code(ctx context.Context) (string, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return "", err
	}

	if g.usedByID == nil {
		if g.usedByEmail == nil &&
			embedCtx.App.Srv().
				AccountService().
				SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageGiftcard) {

			return g.code, nil
		}

		return "", model.NewAppError("Code", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}

	user, err := dataloaders.users.Load(ctx, *g.usedByID)()
	if err != nil {
		return "", err
	}

	if user.ID == embedCtx.AppContext.Session().UserId {
		return g.code, nil
	}

	return "", model.NewAppError("Code", ErrorUnauthorized, nil, "you don't own this giftcard", http.StatusUnauthorized)
}

func graphqlGiftcardsLoader(ctx context.Context, userIDs []string) []*dataloader.Result[*GiftCard] {

}
