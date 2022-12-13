package api

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
)

func DiscountsByDateTimeLoader(ctx context.Context, dateTimes []*time.Time) []*dataloader.Result[*model.DiscountInfo] {
	panic("not implemented")
}
