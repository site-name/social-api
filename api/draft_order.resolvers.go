package api

import (
	"context"
	"fmt"
)

// NOTE: Refer to ./schemas/draft_order.graphqls for details on directives used.
func (r *Resolver) DraftOrderComplete(ctx context.Context, args struct{ Id string }) (*DraftOrderComplete, error) {
	// if !model.IsValidId(args.Id) {
	// 	return nil, model.NewAppError("DraftOrderComplete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, args.Id + " is not a valid order id", http.StatusBadRequest)
	// }

	// embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// orders, appErr := embedCtx.App.Srv().
	// OrderService().
	// FilterOrdersByOptions(&model.OrderFilterOption{

	// })
	// if appErr != nil {
	// 	return nil, appErr
	// }

	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrderCreate(ctx context.Context, args struct {
	Input DraftOrderCreateInput
}) (*DraftOrderCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrderDelete(ctx context.Context, args struct{ Id string }) (*DraftOrderDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrderBulkDelete(ctx context.Context, args struct{ Ids []string }) (*DraftOrderBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DraftOrderUpdate(ctx context.Context, args struct {
	Id    string
	Input DraftOrderInput
}) (*DraftOrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}
