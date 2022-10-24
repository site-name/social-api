package api

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
)

// --------------------------- Order line -----------------------------

func SystemOrderLineToGraphqlOrderLine(line *model.OrderLine) *OrderLine {
	if line == nil {
		return nil
	}

	res := new(OrderLine)
	panic("not implemented")
	return res
}

func graphqlOrderLinesByIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[*OrderLine] {
	panic("not implemented")

}
