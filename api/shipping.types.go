package api

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
)

func shippingZonesByChannelIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*ShippingZone] {
	panic("not implemented")
}
