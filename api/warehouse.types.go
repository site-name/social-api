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

func SystemWarehouseToGraphqlWarehouse(wh *model.WareHouse) *Warehouse {
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
		TotalCount: model.NewPrimitive(int32(count)),
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
	ID    string `json:"id"`
	stock *model.Stock
}

func SystemStockToGraphqlStock(s *model.Stock) *Stock {
	if s == nil {
		return nil
	}

	return &Stock{
		ID:    s.Id,
		stock: s,
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

	return SystemWarehouseToGraphqlWarehouse(warehouse), nil
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

	return int32(s.stock.Quantity), nil
}

func (s *Stock) QuantityAllocated(ctx context.Context) (int32, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return 0, err
	}

	if !embedCtx.App.Srv().
		AccountService().
		SessionHasPermissionToAny(embedCtx.AppContext.Session(), model.PermissionManageProducts, model.PermissionManageOrders) {
		allocations, err := AllocationsByStockIDLoader.Load(ctx, s.ID)()
		if err != nil {
			return 0, err
		}

		var sum int
		for _, allocation := range allocations {
			sum += allocation.QuantityAllocated
		}

		if sum < 0 {
			sum = 0
		}

		return int32(sum), nil
	}

	return 0, model.NewAppError("Stock.QuantityAllocated", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}

func allocationsByStockIDLoader(ctx context.Context, stockIDs []string) []*dataloader.Result[[]*model.Allocation] {
	var (
		res            = make([]*dataloader.Result[[]*model.Allocation], len(stockIDs))
		appErr         *model.AppError
		allocations    model.Allocations
		allocationsMap = map[string]model.Allocations{} // keys are stock ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	allocations, appErr = embedCtx.App.Srv().WarehouseService().AllocationsByOption(nil, &model.AllocationFilterOption{
		StockID: squirrel.Eq{store.AllocationTableName + ".StockID": stockIDs},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, allocation := range allocations {
		allocationsMap[allocation.StockID] = append(allocationsMap[allocation.StockID], allocation)
	}

	for idx, id := range stockIDs {
		res[idx] = &dataloader.Result[[]*model.Allocation]{Data: allocationsMap[id]}
	}
	return res

errorLabel:
	for idx := range stockIDs {
		res[idx] = &dataloader.Result[[]*model.Allocation]{Error: err}
	}
	return res
}

func (s *Stock) ProductVariant(ctx context.Context) (*ProductVariant, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	variant, appErr := embedCtx.App.Srv().ProductService().ProductVariantById(s.stock.ProductVariantID)
	if appErr != nil {
		return nil, appErr
	}

	return SystemProductVariantToGraphqlProductVariant(variant), nil
}

func stocksByIDLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Stock] {
	var (
		res      = make([]*dataloader.Result[*model.Stock], len(ids))
		stocks   model.Stocks
		appErr   *model.AppError
		stockMap = map[string]*model.Stock{} // keys are stock ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	stocks, appErr = embedCtx.App.Srv().
		WarehouseService().
		StocksByOption(nil, &model.StockFilterOption{
			Id:                     squirrel.Eq{store.StockTableName + ".Id": ids},
			SelectRelatedWarehouse: true,
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, st := range stocks {
		stockMap[st.Id] = st
	}

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Stock]{Data: stockMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Stock]{Error: err}
	}
	return res
}

// ----------------- allocation ----------------

type Allocation struct {
	ID string `json:"id"`

	// Quantity  int32      `json:"quantity"`
	// Warehouse *Warehouse `json:"warehouse"`
}

func systemAllocationToGraphqlAllocation(a *model.Allocation) *Allocation {
	if a == nil {
		return nil
	}

	return &Allocation{
		ID: a.Id,
	}
}

func (a *Allocation) Quantity(ctx context.Context) (int32, error) {
	panic("not implemented")
}

func (a *Allocation) Warehouse(ctx context.Context) (Warehouse, error) {
	panic("not implemented")
}

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
