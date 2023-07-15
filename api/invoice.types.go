package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

type Invoice struct {
	ID              string          `json:"id"`
	Metadata        []*MetadataItem `json:"metadata"`
	Number          *string         `json:"number"`
	ExternalURL     *string         `json:"externalUrl"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	CreatedAt       DateTime        `json:"createdAt"`
	URL             *string         `json:"url"`

	// UpdatedAt       DateTime        `json:"updatedAt"`
	// Message         *string         `json:"message"`
	// Status          JobStatusEnum   `json:"status"`
}

func SystemInvoiceToGraphqlInvoice(i *model.Invoice) *Invoice {
	if i == nil {
		return nil
	}

	return &Invoice{
		ID:              i.Id,
		Metadata:        MetadataToSlice(i.Metadata),
		PrivateMetadata: MetadataToSlice(i.PrivateMetadata),
		Number:          &i.Number,
		ExternalURL:     &i.ExternalUrl,
		CreatedAt:       DateTime{util.TimeFromMillis(i.CreateAt)},
		URL:             nil,
	}
}

func invoicesByOrderIDLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.Invoice] {
	var (
		res        = make([]*dataloader.Result[[]*model.Invoice], len(orderIDs))
		invoiceMap = map[string][]*model.Invoice{} // keys are order ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	invoices, appErr := embedCtx.App.Srv().
		InvoiceService().FilterInvoicesByOptions(&model.InvoiceFilterOptions{
		OrderID: squirrel.Eq{model.InvoiceTableName + ".OrderID": orderIDs},
	})
	if appErr != nil {
		for idx := range orderIDs {
			res[idx] = &dataloader.Result[[]*model.Invoice]{Error: appErr}
		}
		return res
	}

	for _, iv := range invoices {
		if iv.OrderID == nil {
			continue
		}
		invoiceMap[*iv.OrderID] = append(invoiceMap[*iv.OrderID], iv)
	}

	for idx, id := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.Invoice]{Data: invoiceMap[id]}
	}
	return res
}
