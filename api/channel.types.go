package api

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
)

func channelByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*Channel] {
	panic("not implemented")
}

func channelBySlugLoader(ctx context.Context, slugs []string) []*dataloader.Result[*Channel] {
	panic("not implemented")
}

func channelByCheckoutLineIDLoader(ctx context.Context, checkoutLineIDs []string) []*dataloader.Result[*Channel] {
	panic("not implemented")
}

func channelByOrderLineIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[[]*Channel] {
	panic("not implemented")
}

func channelWithHasOrdersByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*Channel] {
	panic("not implemented")
}
