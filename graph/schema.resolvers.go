package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sitename/sitename/graph/generated"
	"github.com/sitename/sitename/graph/model"
)

func (m *mutationResolver) WebhookCreate(ctx context.Context, input model.WebhookCreateInput) (*model.WebhookCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) WebhookDelete(ctx context.Context, id string) (*model.WebhookDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) WebhookUpdate(ctx context.Context, id string, input model.WebhookUpdateInput) (*model.WebhookUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CreateWarehouse(ctx context.Context, input model.WarehouseCreateInput) (*model.WarehouseCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) UpdateWarehouse(ctx context.Context, id string, input model.WarehouseUpdateInput) (*model.WarehouseUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DeleteWarehouse(ctx context.Context, id string) (*model.WarehouseDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AssignWarehouseShippingZone(ctx context.Context, id string, shippingZoneIds []string) (*model.WarehouseShippingZoneAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) UnassignWarehouseShippingZone(ctx context.Context, id string, shippingZoneIds []string) (*model.WarehouseShippingZoneUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) StaffNotificationRecipientCreate(ctx context.Context, input model.StaffNotificationRecipientInput) (*model.StaffNotificationRecipientCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) StaffNotificationRecipientUpdate(ctx context.Context, id string, input model.StaffNotificationRecipientInput) (*model.StaffNotificationRecipientUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) StaffNotificationRecipientDelete(ctx context.Context, id string) (*model.StaffNotificationRecipientDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShopDomainUpdate(ctx context.Context, input *model.SiteDomainInput) (*model.ShopDomainUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShopSettingsUpdate(ctx context.Context, input model.ShopSettingsInput) (*model.ShopSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShopFetchTaxRates(ctx context.Context) (*model.ShopFetchTaxRates, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShopSettingsTranslate(ctx context.Context, input model.ShopSettingsTranslationInput, languageCode model.LanguageCodeEnum) (*model.ShopSettingsTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShopAddressUpdate(ctx context.Context, input *model.AddressInput) (*model.ShopAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderSettingsUpdate(ctx context.Context, input model.OrderSettingsUpdateInput) (*model.OrderSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingMethodChannelListingUpdate(ctx context.Context, id string, input model.ShippingMethodChannelListingInput) (*model.ShippingMethodChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingPriceCreate(ctx context.Context, input model.ShippingPriceInput) (*model.ShippingPriceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingPriceDelete(ctx context.Context, id string) (*model.ShippingPriceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingPriceBulkDelete(ctx context.Context, ids []*string) (*model.ShippingPriceBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingPriceUpdate(ctx context.Context, id string, input model.ShippingPriceInput) (*model.ShippingPriceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingPriceTranslate(ctx context.Context, id string, input model.ShippingPriceTranslationInput, languageCode model.LanguageCodeEnum) (*model.ShippingPriceTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingPriceExcludeProducts(ctx context.Context, id string, input model.ShippingPriceExcludeProductsInput) (*model.ShippingPriceExcludeProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingPriceRemoveProductFromExclude(ctx context.Context, id string, products []*string) (*model.ShippingPriceRemoveProductFromExclude, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingZoneCreate(ctx context.Context, input model.ShippingZoneCreateInput) (*model.ShippingZoneCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingZoneDelete(ctx context.Context, id string) (*model.ShippingZoneDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingZoneBulkDelete(ctx context.Context, ids []*string) (*model.ShippingZoneBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ShippingZoneUpdate(ctx context.Context, id string, input model.ShippingZoneUpdateInput) (*model.ShippingZoneUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductAttributeAssign(ctx context.Context, operations []*model.ProductAttributeAssignInput, productTypeID string) (*model.ProductAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductAttributeUnassign(ctx context.Context, attributeIds []*string, productTypeID string) (*model.ProductAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CategoryCreate(ctx context.Context, input model.CategoryInput, parent *string) (*model.CategoryCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CategoryDelete(ctx context.Context, id string) (*model.CategoryDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CategoryBulkDelete(ctx context.Context, ids []*string) (*model.CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CategoryUpdate(ctx context.Context, id string, input model.CategoryInput) (*model.CategoryUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CategoryTranslate(ctx context.Context, id string, input model.TranslationInput, languageCode model.LanguageCodeEnum) (*model.CategoryTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CollectionAddProducts(ctx context.Context, collectionID string, products []*string) (*model.CollectionAddProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CollectionCreate(ctx context.Context, input model.CollectionCreateInput) (*model.CollectionCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CollectionDelete(ctx context.Context, id string) (*model.CollectionDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CollectionReorderProducts(ctx context.Context, collectionID string, moves []*model.MoveProductInput) (*model.CollectionReorderProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CollectionBulkDelete(ctx context.Context, ids []*string) (*model.CollectionBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CollectionRemoveProducts(ctx context.Context, collectionID string, products []*string) (*model.CollectionRemoveProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CollectionUpdate(ctx context.Context, id string, input model.CollectionInput) (*model.CollectionUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CollectionTranslate(ctx context.Context, id string, input model.TranslationInput, languageCode model.LanguageCodeEnum) (*model.CollectionTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CollectionChannelListingUpdate(ctx context.Context, id string, input model.CollectionChannelListingUpdateInput) (*model.CollectionChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductCreate(ctx context.Context, input model.ProductCreateInput) (*model.ProductCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductDelete(ctx context.Context, id string) (*model.ProductDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductBulkDelete(ctx context.Context, ids []*string) (*model.ProductBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductUpdate(ctx context.Context, id string, input model.ProductInput) (*model.ProductUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductTranslate(ctx context.Context, id string, input model.TranslationInput, languageCode model.LanguageCodeEnum) (*model.ProductTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductChannelListingUpdate(ctx context.Context, id string, input model.ProductChannelListingUpdateInput) (*model.ProductChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductMediaCreate(ctx context.Context, input model.ProductMediaCreateInput) (*model.ProductMediaCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantReorder(ctx context.Context, moves []*model.ReorderInput, productID string) (*model.ProductVariantReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductMediaDelete(ctx context.Context, id string) (*model.ProductMediaDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductMediaBulkDelete(ctx context.Context, ids []*string) (*model.ProductMediaBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductMediaReorder(ctx context.Context, mediaIds []*string, productID string) (*model.ProductMediaReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductMediaUpdate(ctx context.Context, id string, input model.ProductMediaUpdateInput) (*model.ProductMediaUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductTypeCreate(ctx context.Context, input model.ProductTypeInput) (*model.ProductTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductTypeDelete(ctx context.Context, id string) (*model.ProductTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductTypeBulkDelete(ctx context.Context, ids []*string) (*model.ProductTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductTypeUpdate(ctx context.Context, id string, input model.ProductTypeInput) (*model.ProductTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductTypeReorderAttributes(ctx context.Context, moves []*model.ReorderInput, productTypeID string, typeArg model.ProductAttributeType) (*model.ProductTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductReorderAttributeValues(ctx context.Context, attributeID string, moves []*model.ReorderInput, productID string) (*model.ProductReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DigitalContentCreate(ctx context.Context, input model.DigitalContentUploadInput, variantID string) (*model.DigitalContentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DigitalContentDelete(ctx context.Context, variantID string) (*model.DigitalContentDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DigitalContentUpdate(ctx context.Context, input model.DigitalContentInput, variantID string) (*model.DigitalContentUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DigitalContentURLCreate(ctx context.Context, input model.DigitalContentURLCreateInput) (*model.DigitalContentURLCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantCreate(ctx context.Context, input model.ProductVariantCreateInput) (*model.ProductVariantCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantDelete(ctx context.Context, id string) (*model.ProductVariantDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantBulkCreate(ctx context.Context, product string, variants []*model.ProductVariantBulkCreateInput) (*model.ProductVariantBulkCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantBulkDelete(ctx context.Context, ids []*string) (*model.ProductVariantBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantStocksCreate(ctx context.Context, stocks []model.StockInput, variantID string) (*model.ProductVariantStocksCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantStocksDelete(ctx context.Context, variantID string, warehouseIds []string) (*model.ProductVariantStocksDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantStocksUpdate(ctx context.Context, stocks []model.StockInput, variantID string) (*model.ProductVariantStocksUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantUpdate(ctx context.Context, id string, input model.ProductVariantInput) (*model.ProductVariantUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantSetDefault(ctx context.Context, productID string, variantID string) (*model.ProductVariantSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantTranslate(ctx context.Context, id string, input model.NameTranslationInput, languageCode model.LanguageCodeEnum) (*model.ProductVariantTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantChannelListingUpdate(ctx context.Context, id string, input []model.ProductVariantChannelListingAddInput) (*model.ProductVariantChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ProductVariantReorderAttributeValues(ctx context.Context, attributeID string, moves []*model.ReorderInput, variantID string) (*model.ProductVariantReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) VariantMediaAssign(ctx context.Context, mediaID string, variantID string) (*model.VariantMediaAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) VariantMediaUnassign(ctx context.Context, mediaID string, variantID string) (*model.VariantMediaUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PaymentCapture(ctx context.Context, amount *string, paymentID string) (*model.PaymentCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PaymentRefund(ctx context.Context, amount *string, paymentID string) (*model.PaymentRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PaymentVoid(ctx context.Context, paymentID string) (*model.PaymentVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PaymentInitialize(ctx context.Context, gateway string, paymentData *string) (*model.PaymentInitialize, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageCreate(ctx context.Context, input model.PageCreateInput) (*model.PageCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageDelete(ctx context.Context, id string) (*model.PageDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageBulkDelete(ctx context.Context, ids []*string) (*model.PageBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageBulkPublish(ctx context.Context, ids []*string, isPublished bool) (*model.PageBulkPublish, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageUpdate(ctx context.Context, id string, input model.PageInput) (*model.PageUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageTranslate(ctx context.Context, id string, input model.PageTranslationInput, languageCode model.LanguageCodeEnum) (*model.PageTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageTypeCreate(ctx context.Context, input model.PageTypeCreateInput) (*model.PageTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageTypeUpdate(ctx context.Context, id *string, input model.PageTypeUpdateInput) (*model.PageTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageTypeDelete(ctx context.Context, id string) (*model.PageTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageTypeBulkDelete(ctx context.Context, ids []string) (*model.PageTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageAttributeAssign(ctx context.Context, attributeIds []string, pageTypeID string) (*model.PageAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageAttributeUnassign(ctx context.Context, attributeIds []string, pageTypeID string) (*model.PageAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageTypeReorderAttributes(ctx context.Context, moves []model.ReorderInput, pageTypeID string) (*model.PageTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PageReorderAttributeValues(ctx context.Context, attributeID string, moves []*model.ReorderInput, pageID string) (*model.PageReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DraftOrderComplete(ctx context.Context, id string) (*model.DraftOrderComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DraftOrderCreate(ctx context.Context, input model.DraftOrderCreateInput) (*model.DraftOrderCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DraftOrderDelete(ctx context.Context, id string) (*model.DraftOrderDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DraftOrderBulkDelete(ctx context.Context, ids []*string) (*model.DraftOrderBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DraftOrderLinesBulkDelete(ctx context.Context, ids []*string) (*model.DraftOrderLinesBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DraftOrderUpdate(ctx context.Context, id string, input model.DraftOrderInput) (*model.DraftOrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderAddNote(ctx context.Context, order string, input model.OrderAddNoteInput) (*model.OrderAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderCancel(ctx context.Context, id string) (*model.OrderCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderCapture(ctx context.Context, amount string, id string) (*model.OrderCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderConfirm(ctx context.Context, id string) (*model.OrderConfirm, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderFulfill(ctx context.Context, input model.OrderFulfillInput, order *string) (*model.OrderFulfill, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderFulfillmentCancel(ctx context.Context, id string, input model.FulfillmentCancelInput) (*model.FulfillmentCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderFulfillmentUpdateTracking(ctx context.Context, id string, input model.FulfillmentUpdateTrackingInput) (*model.FulfillmentUpdateTracking, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderFulfillmentRefundProducts(ctx context.Context, input model.OrderRefundProductsInput, order string) (*model.FulfillmentRefundProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderFulfillmentReturnProducts(ctx context.Context, input model.OrderReturnProductsInput, order string) (*model.FulfillmentReturnProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderLinesCreate(ctx context.Context, id string, input []*model.OrderLineCreateInput) (*model.OrderLinesCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderLineDelete(ctx context.Context, id string) (*model.OrderLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderLineUpdate(ctx context.Context, id string, input model.OrderLineInput) (*model.OrderLineUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderDiscountAdd(ctx context.Context, input model.OrderDiscountCommonInput, orderID string) (*model.OrderDiscountAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderDiscountUpdate(ctx context.Context, discountID string, input model.OrderDiscountCommonInput) (*model.OrderDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderDiscountDelete(ctx context.Context, discountID string) (*model.OrderDiscountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderLineDiscountUpdate(ctx context.Context, input model.OrderDiscountCommonInput, orderLineID string) (*model.OrderLineDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderLineDiscountRemove(ctx context.Context, orderLineID string) (*model.OrderLineDiscountRemove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderMarkAsPaid(ctx context.Context, id string, transactionReference *string) (*model.OrderMarkAsPaid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderRefund(ctx context.Context, amount string, id string) (*model.OrderRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderUpdate(ctx context.Context, id string, input model.OrderUpdateInput) (*model.OrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderUpdateShipping(ctx context.Context, order string, input *model.OrderUpdateShippingInput) (*model.OrderUpdateShipping, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderVoid(ctx context.Context, id string) (*model.OrderVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) OrderBulkCancel(ctx context.Context, ids []*string) (*model.OrderBulkCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DeleteMetadata(ctx context.Context, id string, keys []string) (*model.DeleteMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) DeletePrivateMetadata(ctx context.Context, id string, keys []string) (*model.DeletePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) UpdateMetadata(ctx context.Context, id string, input []model.MetadataInput) (*model.UpdateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) UpdatePrivateMetadata(ctx context.Context, id string, input []model.MetadataInput) (*model.UpdatePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AssignNavigation(ctx context.Context, menu *string, navigationType model.NavigationType) (*model.AssignNavigation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) MenuCreate(ctx context.Context, input model.MenuCreateInput) (*model.MenuCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) MenuDelete(ctx context.Context, id string) (*model.MenuDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) MenuBulkDelete(ctx context.Context, ids []*string) (*model.MenuBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) MenuUpdate(ctx context.Context, id string, input model.MenuInput) (*model.MenuUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) MenuItemCreate(ctx context.Context, input model.MenuItemCreateInput) (*model.MenuItemCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) MenuItemDelete(ctx context.Context, id string) (*model.MenuItemDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) MenuItemBulkDelete(ctx context.Context, ids []*string) (*model.MenuItemBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) MenuItemUpdate(ctx context.Context, id string, input model.MenuItemInput) (*model.MenuItemUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) MenuItemTranslate(ctx context.Context, id string, input model.NameTranslationInput, languageCode model.LanguageCodeEnum) (*model.MenuItemTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) MenuItemMove(ctx context.Context, menu string, moves []*model.MenuItemMoveInput) (*model.MenuItemMove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) InvoiceRequest(ctx context.Context, number *string, orderID string) (*model.InvoiceRequest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) InvoiceRequestDelete(ctx context.Context, id string) (*model.InvoiceRequestDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) InvoiceCreate(ctx context.Context, input model.InvoiceCreateInput, orderID string) (*model.InvoiceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) InvoiceDelete(ctx context.Context, id string) (*model.InvoiceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) InvoiceUpdate(ctx context.Context, id string, input model.UpdateInvoiceInput) (*model.InvoiceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) InvoiceSendNotification(ctx context.Context, id string) (*model.InvoiceSendNotification, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) GiftCardActivate(ctx context.Context, id string) (*model.GiftCardActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) GiftCardCreate(ctx context.Context, input model.GiftCardCreateInput) (*model.GiftCardCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) GiftCardDeactivate(ctx context.Context, id string) (*model.GiftCardDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) GiftCardUpdate(ctx context.Context, id string, input model.GiftCardUpdateInput) (*model.GiftCardUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PluginUpdate(ctx context.Context, id string, input model.PluginUpdateInput) (*model.PluginUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) SaleCreate(ctx context.Context, input model.SaleInput) (*model.SaleCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) SaleDelete(ctx context.Context, id string) (*model.SaleDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) SaleBulkDelete(ctx context.Context, ids []*string) (*model.SaleBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) SaleUpdate(ctx context.Context, id string, input model.SaleInput) (*model.SaleUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) SaleCataloguesAdd(ctx context.Context, id string, input model.CatalogueInput) (*model.SaleAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) SaleCataloguesRemove(ctx context.Context, id string, input model.CatalogueInput) (*model.SaleRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) SaleTranslate(ctx context.Context, id string, input model.NameTranslationInput, languageCode model.LanguageCodeEnum) (*model.SaleTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) SaleChannelListingUpdate(ctx context.Context, id string, input model.SaleChannelListingInput) (*model.SaleChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) VoucherCreate(ctx context.Context, input model.VoucherInput) (*model.VoucherCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) VoucherDelete(ctx context.Context, id string) (*model.VoucherDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) VoucherBulkDelete(ctx context.Context, ids []*string) (*model.VoucherBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) VoucherUpdate(ctx context.Context, id string, input model.VoucherInput) (*model.VoucherUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) VoucherCataloguesAdd(ctx context.Context, id string, input model.CatalogueInput) (*model.VoucherAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) VoucherCataloguesRemove(ctx context.Context, id string, input model.CatalogueInput) (*model.VoucherRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) VoucherTranslate(ctx context.Context, id string, input model.NameTranslationInput, languageCode model.LanguageCodeEnum) (*model.VoucherTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) VoucherChannelListingUpdate(ctx context.Context, id string, input model.VoucherChannelListingInput) (*model.VoucherChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ExportProducts(ctx context.Context, input model.ExportProductsInput) (*model.ExportProducts, error) {
	return m.exportProducts(ctx, input)
}

func (m *mutationResolver) FileUpload(ctx context.Context, file graphql.Upload) (*model.FileUpload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutAddPromoCode(ctx context.Context, checkoutID string, promoCode string) (*model.CheckoutAddPromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutBillingAddressUpdate(ctx context.Context, billingAddress model.AddressInput, checkoutID string) (*model.CheckoutBillingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutComplete(ctx context.Context, checkoutID string, paymentData *string, redirectURL *string, storeSource *bool) (*model.CheckoutComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutCreate(ctx context.Context, input model.CheckoutCreateInput) (*model.CheckoutCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutCustomerAttach(ctx context.Context, checkoutID string) (*model.CheckoutCustomerAttach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutCustomerDetach(ctx context.Context, checkoutID string) (*model.CheckoutCustomerDetach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutEmailUpdate(ctx context.Context, checkoutID *string, email string) (*model.CheckoutEmailUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutLineDelete(ctx context.Context, checkoutID string, lineID *string) (*model.CheckoutLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutLinesAdd(ctx context.Context, checkoutID string, lines []*model.CheckoutLineInput) (*model.CheckoutLinesAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutLinesUpdate(ctx context.Context, checkoutID string, lines []*model.CheckoutLineInput) (*model.CheckoutLinesUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutRemovePromoCode(ctx context.Context, checkoutID string, promoCode string) (*model.CheckoutRemovePromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutPaymentCreate(ctx context.Context, checkoutID string, input model.PaymentInput) (*model.CheckoutPaymentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutShippingAddressUpdate(ctx context.Context, checkoutID string, shippingAddress model.AddressInput) (*model.CheckoutShippingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutShippingMethodUpdate(ctx context.Context, checkoutID *string, shippingMethodID string) (*model.CheckoutShippingMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CheckoutLanguageCodeUpdate(ctx context.Context, checkoutID string, languageCode model.LanguageCodeEnum) (*model.CheckoutLanguageCodeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ChannelCreate(ctx context.Context, input model.ChannelCreateInput) (*model.ChannelCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ChannelUpdate(ctx context.Context, id string, input model.ChannelUpdateInput) (*model.ChannelUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ChannelDelete(ctx context.Context, id string, input *model.ChannelDeleteInput) (*model.ChannelDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ChannelActivate(ctx context.Context, id string) (*model.ChannelActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ChannelDeactivate(ctx context.Context, id string) (*model.ChannelDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeCreate(ctx context.Context, input model.AttributeCreateInput) (*model.AttributeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeDelete(ctx context.Context, id string) (*model.AttributeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeUpdate(ctx context.Context, id string, input model.AttributeUpdateInput) (*model.AttributeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeTranslate(ctx context.Context, id string, input model.NameTranslationInput, languageCode model.LanguageCodeEnum) (*model.AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeBulkDelete(ctx context.Context, ids []*string) (*model.AttributeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeValueBulkDelete(ctx context.Context, ids []*string) (*model.AttributeValueBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeValueCreate(ctx context.Context, attribute string, input model.AttributeValueCreateInput) (*model.AttributeValueCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeValueDelete(ctx context.Context, id string) (*model.AttributeValueDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeValueUpdate(ctx context.Context, id string, input model.AttributeValueCreateInput) (*model.AttributeValueUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeValueTranslate(ctx context.Context, id string, input model.AttributeValueTranslationInput, languageCode model.LanguageCodeEnum) (*model.AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AttributeReorderValues(ctx context.Context, attributeID string, moves []*model.ReorderInput) (*model.AttributeReorderValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppCreate(ctx context.Context, input model.AppInput) (*model.AppCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppUpdate(ctx context.Context, id string, input model.AppInput) (*model.AppUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppDelete(ctx context.Context, id string) (*model.AppDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppTokenCreate(ctx context.Context, input model.AppTokenInput) (*model.AppTokenCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppTokenDelete(ctx context.Context, id string) (*model.AppTokenDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppTokenVerify(ctx context.Context, token string) (*model.AppTokenVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppInstall(ctx context.Context, input model.AppInstallInput) (*model.AppInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppRetryInstall(ctx context.Context, activateAfterInstallation *bool, id string) (*model.AppRetryInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppDeleteFailedInstallation(ctx context.Context, id string) (*model.AppDeleteFailedInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppFetchManifest(ctx context.Context, manifestURL string) (*model.AppFetchManifest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppActivate(ctx context.Context, id string) (*model.AppActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AppDeactivate(ctx context.Context, id string) (*model.AppDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) TokenCreate(ctx context.Context, email string, password string) (*model.CreateToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) TokenRefresh(ctx context.Context, csrfToken *string, refreshToken *string) (*model.RefreshToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) TokenVerify(ctx context.Context, token string) (*model.VerifyToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) TokensDeactivateAll(ctx context.Context) (*model.DeactivateAllUserTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ExternalAuthenticationURL(ctx context.Context, input string, pluginID string) (*model.ExternalAuthenticationURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ExternalObtainAccessTokens(ctx context.Context, input string, pluginID string) (*model.ExternalObtainAccessTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ExternalRefresh(ctx context.Context, input string, pluginID string) (*model.ExternalRefresh, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ExternalLogout(ctx context.Context, input string, pluginID string) (*model.ExternalLogout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ExternalVerify(ctx context.Context, input string, pluginID string) (*model.ExternalVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) RequestPasswordReset(ctx context.Context, email string, redirectURL string) (*model.RequestPasswordReset, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ConfirmAccount(ctx context.Context, email string, token string) (*model.ConfirmAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) SetPassword(ctx context.Context, email string, password string, token string) (*model.SetPassword, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PasswordChange(ctx context.Context, newPassword string, oldPassword string) (*model.PasswordChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) RequestEmailChange(ctx context.Context, newEmail string, password string, redirectURL string) (*model.RequestEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) ConfirmEmailChange(ctx context.Context, token string) (*model.ConfirmEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AccountAddressCreate(ctx context.Context, input model.AddressInput, typeArg *model.AddressTypeEnum) (*model.AccountAddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AccountAddressUpdate(ctx context.Context, id string, input model.AddressInput) (*model.AccountAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AccountAddressDelete(ctx context.Context, id string) (*model.AccountAddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AccountSetDefaultAddress(ctx context.Context, id string, typeArg model.AddressTypeEnum) (*model.AccountSetDefaultAddress, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AccountRegister(ctx context.Context, input model.AccountRegisterInput) (*model.AccountRegister, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AccountUpdate(ctx context.Context, input model.AccountInput) (*model.AccountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AccountRequestDeletion(ctx context.Context, redirectURL string) (*model.AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AccountDelete(ctx context.Context, token string) (*model.AccountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AddressCreate(ctx context.Context, input model.AddressInput, userID string) (*model.AddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AddressUpdate(ctx context.Context, id string, input model.AddressInput) (*model.AddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AddressDelete(ctx context.Context, id string) (*model.AddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) AddressSetDefault(ctx context.Context, addressID string, typeArg model.AddressTypeEnum, userID string) (*model.AddressSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CustomerCreate(ctx context.Context, input model.UserCreateInput) (*model.CustomerCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CustomerUpdate(ctx context.Context, id string, input model.CustomerInput) (*model.CustomerUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CustomerDelete(ctx context.Context, id string) (*model.CustomerDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) CustomerBulkDelete(ctx context.Context, ids []*string) (*model.CustomerBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) StaffCreate(ctx context.Context, input model.StaffCreateInput) (*model.StaffCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) StaffUpdate(ctx context.Context, id string, input model.StaffUpdateInput) (*model.StaffUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) StaffDelete(ctx context.Context, id string) (*model.StaffDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) StaffBulkDelete(ctx context.Context, ids []*string) (*model.StaffBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) UserAvatarUpdate(ctx context.Context, image graphql.Upload) (*model.UserAvatarUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) UserAvatarDelete(ctx context.Context) (*model.UserAvatarDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) UserBulkSetActive(ctx context.Context, ids []*string, isActive bool) (*model.UserBulkSetActive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PermissionGroupCreate(ctx context.Context, input model.PermissionGroupCreateInput) (*model.PermissionGroupCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PermissionGroupUpdate(ctx context.Context, id string, input model.PermissionGroupUpdateInput) (*model.PermissionGroupUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (m *mutationResolver) PermissionGroupDelete(ctx context.Context, id string) (*model.PermissionGroupDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Webhook(ctx context.Context, id string) (*model.Webhook, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) WebhookEvents(ctx context.Context) ([]*model.WebhookEvent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) WebhookSamplePayload(ctx context.Context, eventType model.WebhookSampleEventTypeEnum) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Warehouse(ctx context.Context, id string) (*model.Warehouse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Warehouses(ctx context.Context, filter *model.WarehouseFilterInput, sortBy *model.WarehouseSortingInput, before *string, after *string, first *int, last *int) (*model.WarehouseCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Translations(ctx context.Context, kind model.TranslatableKinds, before *string, after *string, first *int, last *int) (*model.TranslatableItemConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Translation(ctx context.Context, id string, kind model.TranslatableKinds) (model.TranslatableItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Stock(ctx context.Context, id string) (*model.Stock, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Stocks(ctx context.Context, filter *model.StockFilterInput, before *string, after *string, first *int, last *int) (*model.StockCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Shop(ctx context.Context) (*model.Shop, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OrderSettings(ctx context.Context) (*model.OrderSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ShippingZone(ctx context.Context, id string, channel *string) (*model.ShippingZone, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ShippingZones(ctx context.Context, filter *model.ShippingZoneFilterInput, channel *string, before *string, after *string, first *int, last *int) (*model.ShippingZoneCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DigitalContent(ctx context.Context, id string) (*model.DigitalContent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DigitalContents(ctx context.Context, before *string, after *string, first *int, last *int) (*model.DigitalContentCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Categories(ctx context.Context, filter *model.CategoryFilterInput, sortBy *model.CategorySortingInput, level *int, before *string, after *string, first *int, last *int) (*model.CategoryCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Category(ctx context.Context, id *string, slug *string) (*model.Category, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Collection(ctx context.Context, id *string, slug *string, channel *string) (*model.Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Collections(ctx context.Context, filter *model.CollectionFilterInput, sortBy *model.CollectionSortingInput, channel *string, before *string, after *string, first *int, last *int) (*model.CollectionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Product(ctx context.Context, id *string, slug *string, channel *string) (*model.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context, filter *model.ProductFilterInput, sortBy *model.ProductOrder, channel *string, before *string, after *string, first *int, last *int) (*model.ProductCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductType(ctx context.Context, id string) (*model.ProductType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductTypes(ctx context.Context, filter *model.ProductTypeFilterInput, sortBy *model.ProductTypeSortingInput, before *string, after *string, first *int, last *int) (*model.ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductVariant(ctx context.Context, id *string, sku *string, channel *string) (*model.ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductVariants(ctx context.Context, ids []*string, channel *string, filter *model.ProductVariantFilterInput, before *string, after *string, first *int, last *int) (*model.ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ReportProductSales(ctx context.Context, period model.ReportingPeriod, channel string, before *string, after *string, first *int, last *int) (*model.ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Payment(ctx context.Context, id string) (*model.Payment, error) {
	return r.payment(ctx, id)
}

func (r *queryResolver) Payments(ctx context.Context, filter *model.PaymentFilterInput, before *string, after *string, first *int, last *int) (*model.PaymentCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Page(ctx context.Context, id *string, slug *string) (*model.Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Pages(ctx context.Context, sortBy *model.PageSortingInput, filter *model.PageFilterInput, before *string, after *string, first *int, last *int) (*model.PageCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PageType(ctx context.Context, id string) (*model.PageType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PageTypes(ctx context.Context, sortBy *model.PageTypeSortingInput, filter *model.PageTypeFilterInput, before *string, after *string, first *int, last *int) (*model.PageTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) HomepageEvents(ctx context.Context, before *string, after *string, first *int, last *int) (*model.OrderEventCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Order(ctx context.Context, id string) (*model.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Orders(ctx context.Context, sortBy *model.OrderSortingInput, filter *model.OrderFilterInput, channel *string, before *string, after *string, first *int, last *int) (*model.OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DraftOrders(ctx context.Context, sortBy *model.OrderSortingInput, filter *model.OrderDraftFilterInput, before *string, after *string, first *int, last *int) (*model.OrderCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OrdersTotal(ctx context.Context, period *model.ReportingPeriod, channel *string) (*model.TaxedMoney, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OrderByToken(ctx context.Context, token string) (*model.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Menu(ctx context.Context, channel *string, id *string, name *string, slug *string) (*model.Menu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Menus(ctx context.Context, channel *string, sortBy *model.MenuSortingInput, filter *model.MenuFilterInput, before *string, after *string, first *int, last *int) (*model.MenuCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MenuItem(ctx context.Context, id string, channel *string) (*model.MenuItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MenuItems(ctx context.Context, channel *string, sortBy *model.MenuItemSortingInput, filter *model.MenuItemFilterInput, before *string, after *string, first *int, last *int) (*model.MenuItemCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GiftCard(ctx context.Context, id string) (*model.GiftCard, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GiftCards(ctx context.Context, before *string, after *string, first *int, last *int) (*model.GiftCardCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Plugin(ctx context.Context, id string) (*model.Plugin, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Plugins(ctx context.Context, filter *model.PluginFilterInput, sortBy *model.PluginSortingInput, before *string, after *string, first *int, last *int) (*model.PluginCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sale(ctx context.Context, id string, channel *string) (*model.Sale, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sales(ctx context.Context, filter *model.SaleFilterInput, sortBy *model.SaleSortingInput, query *string, channel *string, before *string, after *string, first *int, last *int) (*model.SaleCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Voucher(ctx context.Context, id string, channel *string) (*model.Voucher, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Vouchers(ctx context.Context, filter *model.VoucherFilterInput, sortBy *model.VoucherSortingInput, query *string, channel *string, before *string, after *string, first *int, last *int) (*model.VoucherCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ExportFile(ctx context.Context, id string) (*model.ExportFile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ExportFiles(ctx context.Context, filter *model.ExportFileFilterInput, sortBy *model.ExportFileSortingInput, before *string, after *string, first *int, last *int) (*model.ExportFileCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TaxTypes(ctx context.Context) ([]*model.TaxType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkout(ctx context.Context, token *string) (*model.Checkout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkouts(ctx context.Context, channel *string, before *string, after *string, first *int, last *int) (*model.CheckoutCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CheckoutLine(ctx context.Context, id *string) (*model.CheckoutLine, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CheckoutLines(ctx context.Context, before *string, after *string, first *int, last *int) (*model.CheckoutLineCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Channel(ctx context.Context, id *string) (*model.Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Channels(ctx context.Context) ([]model.Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attributes(ctx context.Context, filter *model.AttributeFilterInput, sortBy *model.AttributeSortingInput, before *string, after *string, first *int, last *int) (*model.AttributeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attribute(ctx context.Context, id *string, slug *string) (*model.Attribute, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AppsInstallations(ctx context.Context) ([]model.AppInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Apps(ctx context.Context, filter *model.AppFilterInput, sortBy *model.AppSortingInput, before *string, after *string, first *int, last *int) (*model.AppCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) App(ctx context.Context, id *string) (*model.App, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AddressValidationRules(ctx context.Context, countryCode model.CountryCode, countryArea *string, city *string, cityArea *string) (*model.AddressValidationData, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Address(ctx context.Context, id string) (*model.Address, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Customers(ctx context.Context, filter *model.CustomerFilterInput, sortBy *model.UserSortingInput, before *string, after *string, first *int, last *int) (*model.UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PermissionGroups(ctx context.Context, filter *model.PermissionGroupFilterInput, sortBy *model.PermissionGroupSortingInput, before *string, after *string, first *int, last *int) (*model.GroupCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PermissionGroup(ctx context.Context, id string) (*model.Group, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) StaffUsers(ctx context.Context, filter *model.StaffUserInput, sortBy *model.UserSortingInput, before *string, after *string, first *int, last *int) (*model.UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id *string, email *string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
