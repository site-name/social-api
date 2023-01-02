package api

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type Warehouse struct {
	ID                    string                             `json:"id"`
	Name                  string                             `json:"name"`
	Slug                  string                             `json:"slug"`
	Email                 string                             `json:"email"`
	IsPrivate             bool                               `json:"isPrivate"`
	PrivateMetadata       []*MetadataItem                    `json:"privateMetadata"`
	Metadata              []*MetadataItem                    `json:"metadata"`
	ClickAndCollectOption WarehouseClickAndCollectOptionEnum `json:"clickAndCollectOption"`

	addressID *string
	// ShippingZones         *ShippingZoneCountableConnection   `json:"shippingZones"`
	// Address               *Address                           `json:"address"`
}

func SystemWarehouseTpGraphqlWarehouse(wh *model.WareHouse) *Warehouse {
	if wh == nil {
		return nil
	}

	return &Warehouse{
		ID:                    wh.Id,
		Name:                  wh.Name,
		Slug:                  wh.Slug,
		Email:                 wh.Email,
		IsPrivate:             *wh.IsPrivate,
		Metadata:              MetadataToSlice(wh.Metadata),
		PrivateMetadata:       MetadataToSlice(wh.Metadata),
		ClickAndCollectOption: WarehouseClickAndCollectOptionEnum(wh.ClickAndCollectOption),

		addressID: wh.AddressID,
	}
}

func (w *Warehouse) ShippingZones(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*ShippingZoneCountableConnection, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	filterOpts := &model.ShippingZoneFilterOption{
		PaginationOptions: model.PaginationOptions{
			Before: args.Before,
			After:  args.After,
			First:  args.First,
			Last:   args.Last,
		},
	}

	zones, appErr := embedCtx.App.Srv().
		ShippingService().
		ShippingZonesByOption(filterOpts)
	if appErr != nil {
		return nil, appErr
	}

	count, err := embedCtx.App.Srv().Store.ShippingZone().CountByOptions(filterOpts)
	if err != nil {
		return nil, err
	}

	hasNextPage := len(zones) == int(filterOpts.Limit())
	edgeLength := len(zones)
	if hasNextPage {
		edgeLength--
	}

	res := &ShippingZoneCountableConnection{
		TotalCount: model.NewInt32(int32(count)),
		PageInfo: &PageInfo{
			HasPreviousPage: filterOpts.HasPreviousPage(),
			HasNextPage:     hasNextPage,
		},
		Edges: make([]*ShippingZoneCountableEdge, edgeLength),
	}

	for i := 0; i < edgeLength; i++ {
		res.Edges[i] = &ShippingZoneCountableEdge{
			Node:   SystemShippingZoneToGraphqlShippingZone(zones[i]),
			Cursor: base64.StdEncoding.EncodeToString([]byte(zones[i].Name)),
		}
	}

	res.PageInfo.StartCursor = &res.Edges[0].Cursor
	res.PageInfo.EndCursor = &res.Edges[edgeLength-1].Cursor

	return res, nil
}

func (w *Warehouse) Address(ctx context.Context) (*Address, error) {
	if w.addressID == nil {
		return nil, nil
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	address, appErr := embedCtx.App.Srv().AccountService().AddressById(*w.addressID)
	if appErr != nil {
		return nil, err
	}

	return SystemAddressToGraphqlAddress(address), nil
}

func warehouseByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.WareHouse] {
	var (
		res          = make([]*dataloader.Result[*model.WareHouse], len(ids))
		appErr       *model.AppError
		warehouses   model.Warehouses
		warehouseMap = map[string]*model.WareHouse{} // keys are warehouse ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	warehouses, appErr = embedCtx.App.Srv().
		WarehouseService().
		WarehousesByOption(&model.WarehouseFilterOption{
			Id: squirrel.Eq{store.WarehouseTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, wh := range warehouses {
		warehouseMap[wh.Id] = wh
	}

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.WareHouse]{Data: warehouseMap[id]}
	}
	return res

errorLabel:
	for i := range ids {
		res[i] = &dataloader.Result[*model.WareHouse]{Error: err}
	}
	return res
}

// ---------------------- stock --------------------

type Stock struct {
	quantity          int32  `json:"quantity"`
	ID                string `json:"id"`
	quantityAllocated int32  `json:"quantityAllocated"`

	warehouseID      string
	productVariantID string
}

func SystemStockToGraphqlStock(s *model.Stock) *Stock {
	if s == nil {
		return nil
	}

	return &Stock{
		ID:               s.Id,
		quantity:         int32(s.Quantity),
		warehouseID:      s.WarehouseID,
		productVariantID: s.ProductVariantID,
	}
}

func (s *Stock) Warehouse(ctx context.Context) (*Warehouse, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	warehouse, appErr := embedCtx.App.Srv().WarehouseService().WarehouseByStockID(s.ID)
	if appErr != nil {
		return nil, appErr
	}

	return SystemWarehouseTpGraphqlWarehouse(warehouse), nil
}

func (s *Stock) Quantity(ctx context.Context) (int32, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return 0, err
	}

	if !embedCtx.App.Srv().
		AccountService().
		SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.PermissionManageProducts, model.PermissionManageOrders) {
		return 0, model.NewAppError("stock.Wuantity", ErrorUnauthorized, nil, "You are not authorized to perform this action", http.StatusUnauthorized)
	}

	return s.quantity, nil
}

func (s *Stock) QuantityAllocated(ctx context.Context) (int32, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return 0, err
	}

	if !embedCtx.App.Srv().
		AccountService().
		SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.PermissionManageProducts, model.PermissionManageOrders) {
		return 0, model.NewAppError("stock.Wuantity", ErrorUnauthorized, nil, "You are not authorized to perform this action", http.StatusUnauthorized)
	}

	return s.quantityAllocated, nil
}

func (s *Stock) ProductVariant(ctx context.Context) (*ProductVariant, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	variant, appErr := embedCtx.App.Srv().ProductService().ProductVariantById(s.productVariantID)
	if appErr != nil {
		return nil, appErr
	}

	return SystemProductVariantToGraphqlProductVariant(variant), nil
}

// ----------------- allocation ----------------

func allocationsByOrderLineIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[[]*model.Allocation] {
	var (
		res           = make([]*dataloader.Result[[]*model.Allocation], len(orderLineIDs))
		appErr        *model.AppError
		allocationMap = map[string]model.Allocations{}
		allocations   model.Allocations
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	allocations, appErr = embedCtx.App.Srv().WarehouseService().AllocationsByOption(nil, &model.AllocationFilterOption{
		OrderLineID: squirrel.Eq{store.AllocationTableName + ".OrderLineID": orderLineIDs},
	})

	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, all := range allocations {
		allocationMap[all.OrderLineID] = append(allocationMap[all.OrderLineID], all)
	}

	for idx, id := range orderLineIDs {
		res[idx] = &dataloader.Result[[]*model.Allocation]{Data: allocationMap[id]}
	}
	return res

errorLabel:
	for idx := range orderLineIDs {
		res[idx] = &dataloader.Result[[]*model.Allocation]{Error: err}
	}
	return res
}
