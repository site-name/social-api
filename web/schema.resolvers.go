package web

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/sitename/sitename/web/model"
)

func (r *mutationResolver) WebhookCreate(ctx context.Context, input model.WebhookCreateInput) (*model.WebhookCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) WebhookDelete(ctx context.Context, id string) (*model.WebhookDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) WebhookUpdate(ctx context.Context, id string, input model.WebhookUpdateInput) (*model.WebhookUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateWarehouse(ctx context.Context, input model.WarehouseCreateInput) (*model.WarehouseCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateWarehouse(ctx context.Context, id string, input model.WarehouseUpdateInput) (*model.WarehouseUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteWarehouse(ctx context.Context, id string) (*model.WarehouseDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AssignWarehouseShippingZone(ctx context.Context, id string, shippingZoneIds []string) (*model.WarehouseShippingZoneAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnassignWarehouseShippingZone(ctx context.Context, id string, shippingZoneIds []string) (*model.WarehouseShippingZoneUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffNotificationRecipientCreate(ctx context.Context, input model.StaffNotificationRecipientInput) (*model.StaffNotificationRecipientCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffNotificationRecipientUpdate(ctx context.Context, id string, input model.StaffNotificationRecipientInput) (*model.StaffNotificationRecipientUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffNotificationRecipientDelete(ctx context.Context, id string) (*model.StaffNotificationRecipientDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopDomainUpdate(ctx context.Context, input *model.SiteDomainInput) (*model.ShopDomainUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopSettingsUpdate(ctx context.Context, input model.ShopSettingsInput) (*model.ShopSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopFetchTaxRates(ctx context.Context) (*model.ShopFetchTaxRates, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopSettingsTranslate(ctx context.Context, input model.ShopSettingsTranslationInput, languageCode model.LanguageCodeEnum) (*model.ShopSettingsTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShopAddressUpdate(ctx context.Context, input *model.AddressInput) (*model.ShopAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderSettingsUpdate(ctx context.Context, input model.OrderSettingsUpdateInput) (*model.OrderSettingsUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingMethodChannelListingUpdate(ctx context.Context, id string, input model.ShippingMethodChannelListingInput) (*model.ShippingMethodChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceCreate(ctx context.Context, input model.ShippingPriceInput) (*model.ShippingPriceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceDelete(ctx context.Context, id string) (*model.ShippingPriceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceBulkDelete(ctx context.Context, ids []*string) (*model.ShippingPriceBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceUpdate(ctx context.Context, id string, input model.ShippingPriceInput) (*model.ShippingPriceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceTranslate(ctx context.Context, id string, input model.ShippingPriceTranslationInput, languageCode model.LanguageCodeEnum) (*model.ShippingPriceTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceExcludeProducts(ctx context.Context, id string, input model.ShippingPriceExcludeProductsInput) (*model.ShippingPriceExcludeProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceRemoveProductFromExclude(ctx context.Context, id string, products []*string) (*model.ShippingPriceRemoveProductFromExclude, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneCreate(ctx context.Context, input model.ShippingZoneCreateInput) (*model.ShippingZoneCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneDelete(ctx context.Context, id string) (*model.ShippingZoneDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneBulkDelete(ctx context.Context, ids []*string) (*model.ShippingZoneBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneUpdate(ctx context.Context, id string, input model.ShippingZoneUpdateInput) (*model.ShippingZoneUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductAttributeAssign(ctx context.Context, operations []*model.ProductAttributeAssignInput, productTypeID string) (*model.ProductAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductAttributeUnassign(ctx context.Context, attributeIds []*string, productTypeID string) (*model.ProductAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryCreate(ctx context.Context, input model.CategoryInput, parent *string) (*model.CategoryCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryDelete(ctx context.Context, id string) (*model.CategoryDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryBulkDelete(ctx context.Context, ids []*string) (*model.CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryUpdate(ctx context.Context, id string, input model.CategoryInput) (*model.CategoryUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryTranslate(ctx context.Context, id string, input model.TranslationInput, languageCode model.LanguageCodeEnum) (*model.CategoryTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionAddProducts(ctx context.Context, collectionID string, products []*string) (*model.CollectionAddProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionCreate(ctx context.Context, input model.CollectionCreateInput) (*model.CollectionCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionDelete(ctx context.Context, id string) (*model.CollectionDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionReorderProducts(ctx context.Context, collectionID string, moves []*model.MoveProductInput) (*model.CollectionReorderProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionBulkDelete(ctx context.Context, ids []*string) (*model.CollectionBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionRemoveProducts(ctx context.Context, collectionID string, products []*string) (*model.CollectionRemoveProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionUpdate(ctx context.Context, id string, input model.CollectionInput) (*model.CollectionUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionTranslate(ctx context.Context, id string, input model.TranslationInput, languageCode model.LanguageCodeEnum) (*model.CollectionTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionChannelListingUpdate(ctx context.Context, id string, input model.CollectionChannelListingUpdateInput) (*model.CollectionChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductCreate(ctx context.Context, input model.ProductCreateInput) (*model.ProductCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductDelete(ctx context.Context, id string) (*model.ProductDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductBulkDelete(ctx context.Context, ids []*string) (*model.ProductBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductUpdate(ctx context.Context, id string, input model.ProductInput) (*model.ProductUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTranslate(ctx context.Context, id string, input model.TranslationInput, languageCode model.LanguageCodeEnum) (*model.ProductTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductChannelListingUpdate(ctx context.Context, id string, input model.ProductChannelListingUpdateInput) (*model.ProductChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaCreate(ctx context.Context, input model.ProductMediaCreateInput) (*model.ProductMediaCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantReorder(ctx context.Context, moves []*model.ReorderInput, productID string) (*model.ProductVariantReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaDelete(ctx context.Context, id string) (*model.ProductMediaDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaBulkDelete(ctx context.Context, ids []*string) (*model.ProductMediaBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaReorder(ctx context.Context, mediaIds []*string, productID string) (*model.ProductMediaReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaUpdate(ctx context.Context, id string, input model.ProductMediaUpdateInput) (*model.ProductMediaUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeCreate(ctx context.Context, input model.ProductTypeInput) (*model.ProductTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeDelete(ctx context.Context, id string) (*model.ProductTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeBulkDelete(ctx context.Context, ids []*string) (*model.ProductTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeUpdate(ctx context.Context, id string, input model.ProductTypeInput) (*model.ProductTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeReorderAttributes(ctx context.Context, moves []*model.ReorderInput, productTypeID string, typeArg model.ProductAttributeType) (*model.ProductTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductReorderAttributeValues(ctx context.Context, attributeID string, moves []*model.ReorderInput, productID string) (*model.ProductReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentCreate(ctx context.Context, input model.DigitalContentUploadInput, variantID string) (*model.DigitalContentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentDelete(ctx context.Context, variantID string) (*model.DigitalContentDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentUpdate(ctx context.Context, input model.DigitalContentInput, variantID string) (*model.DigitalContentUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentURLCreate(ctx context.Context, input model.DigitalContentURLCreateInput) (*model.DigitalContentURLCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantCreate(ctx context.Context, input model.ProductVariantCreateInput) (*model.ProductVariantCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantDelete(ctx context.Context, id string) (*model.ProductVariantDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantBulkCreate(ctx context.Context, product string, variants []*model.ProductVariantBulkCreateInput) (*model.ProductVariantBulkCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantBulkDelete(ctx context.Context, ids []*string) (*model.ProductVariantBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksCreate(ctx context.Context, stocks []model.StockInput, variantID string) (*model.ProductVariantStocksCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksDelete(ctx context.Context, variantID string, warehouseIds []string) (*model.ProductVariantStocksDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksUpdate(ctx context.Context, stocks []model.StockInput, variantID string) (*model.ProductVariantStocksUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantUpdate(ctx context.Context, id string, input model.ProductVariantInput) (*model.ProductVariantUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantSetDefault(ctx context.Context, productID string, variantID string) (*model.ProductVariantSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantTranslate(ctx context.Context, id string, input model.NameTranslationInput, languageCode model.LanguageCodeEnum) (*model.ProductVariantTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantChannelListingUpdate(ctx context.Context, id string, input []model.ProductVariantChannelListingAddInput) (*model.ProductVariantChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantReorderAttributeValues(ctx context.Context, attributeID string, moves []*model.ReorderInput, variantID string) (*model.ProductVariantReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VariantMediaAssign(ctx context.Context, mediaID string, variantID string) (*model.VariantMediaAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VariantMediaUnassign(ctx context.Context, mediaID string, variantID string) (*model.VariantMediaUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentCapture(ctx context.Context, amount *string, paymentID string) (*model.PaymentCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentRefund(ctx context.Context, amount *string, paymentID string) (*model.PaymentRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentVoid(ctx context.Context, paymentID string) (*model.PaymentVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentInitialize(ctx context.Context, channel *string, gateway string, paymentData *string) (*model.PaymentInitialize, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageCreate(ctx context.Context, input model.PageCreateInput) (*model.PageCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageDelete(ctx context.Context, id string) (*model.PageDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageBulkDelete(ctx context.Context, ids []*string) (*model.PageBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageBulkPublish(ctx context.Context, ids []*string, isPublished bool) (*model.PageBulkPublish, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageUpdate(ctx context.Context, id string, input model.PageInput) (*model.PageUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTranslate(ctx context.Context, id string, input model.PageTranslationInput, languageCode model.LanguageCodeEnum) (*model.PageTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeCreate(ctx context.Context, input model.PageTypeCreateInput) (*model.PageTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeUpdate(ctx context.Context, id *string, input model.PageTypeUpdateInput) (*model.PageTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeDelete(ctx context.Context, id string) (*model.PageTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeBulkDelete(ctx context.Context, ids []string) (*model.PageTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageAttributeAssign(ctx context.Context, attributeIds []string, pageTypeID string) (*model.PageAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageAttributeUnassign(ctx context.Context, attributeIds []string, pageTypeID string) (*model.PageAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeReorderAttributes(ctx context.Context, moves []model.ReorderInput, pageTypeID string) (*model.PageTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageReorderAttributeValues(ctx context.Context, attributeID string, moves []*model.ReorderInput, pageID string) (*model.PageReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderComplete(ctx context.Context, id string) (*model.DraftOrderComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderCreate(ctx context.Context, input model.DraftOrderCreateInput) (*model.DraftOrderCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderDelete(ctx context.Context, id string) (*model.DraftOrderDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderBulkDelete(ctx context.Context, ids []*string) (*model.DraftOrderBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderLinesBulkDelete(ctx context.Context, ids []*string) (*model.DraftOrderLinesBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderUpdate(ctx context.Context, id string, input model.DraftOrderInput) (*model.DraftOrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderAddNote(ctx context.Context, order string, input model.OrderAddNoteInput) (*model.OrderAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderCancel(ctx context.Context, id string) (*model.OrderCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderCapture(ctx context.Context, amount string, id string) (*model.OrderCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderConfirm(ctx context.Context, id string) (*model.OrderConfirm, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderFulfill(ctx context.Context, input model.OrderFulfillInput, order *string) (*model.OrderFulfill, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderFulfillmentCancel(ctx context.Context, id string, input model.FulfillmentCancelInput) (*model.FulfillmentCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderFulfillmentUpdateTracking(ctx context.Context, id string, input model.FulfillmentUpdateTrackingInput) (*model.FulfillmentUpdateTracking, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderFulfillmentRefundProducts(ctx context.Context, input model.OrderRefundProductsInput, order string) (*model.FulfillmentRefundProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderFulfillmentReturnProducts(ctx context.Context, input model.OrderReturnProductsInput, order string) (*model.FulfillmentReturnProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLinesCreate(ctx context.Context, id string, input []*model.OrderLineCreateInput) (*model.OrderLinesCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineDelete(ctx context.Context, id string) (*model.OrderLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineUpdate(ctx context.Context, id string, input model.OrderLineInput) (*model.OrderLineUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderDiscountAdd(ctx context.Context, input model.OrderDiscountCommonInput, orderID string) (*model.OrderDiscountAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderDiscountUpdate(ctx context.Context, discountID string, input model.OrderDiscountCommonInput) (*model.OrderDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderDiscountDelete(ctx context.Context, discountID string) (*model.OrderDiscountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineDiscountUpdate(ctx context.Context, input model.OrderDiscountCommonInput, orderLineID string) (*model.OrderLineDiscountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderLineDiscountRemove(ctx context.Context, orderLineID string) (*model.OrderLineDiscountRemove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderMarkAsPaid(ctx context.Context, id string, transactionReference *string) (*model.OrderMarkAsPaid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderRefund(ctx context.Context, amount string, id string) (*model.OrderRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderUpdate(ctx context.Context, id string, input model.OrderUpdateInput) (*model.OrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderUpdateShipping(ctx context.Context, order string, input *model.OrderUpdateShippingInput) (*model.OrderUpdateShipping, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderVoid(ctx context.Context, id string) (*model.OrderVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) OrderBulkCancel(ctx context.Context, ids []*string) (*model.OrderBulkCancel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteMetadata(ctx context.Context, id string, keys []string) (*model.DeleteMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePrivateMetadata(ctx context.Context, id string, keys []string) (*model.DeletePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateMetadata(ctx context.Context, id string, input []model.MetadataInput) (*model.UpdateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdatePrivateMetadata(ctx context.Context, id string, input []model.MetadataInput) (*model.UpdatePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AssignNavigation(ctx context.Context, menu *string, navigationType model.NavigationType) (*model.AssignNavigation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuCreate(ctx context.Context, input model.MenuCreateInput) (*model.MenuCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuDelete(ctx context.Context, id string) (*model.MenuDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuBulkDelete(ctx context.Context, ids []*string) (*model.MenuBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuUpdate(ctx context.Context, id string, input model.MenuInput) (*model.MenuUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemCreate(ctx context.Context, input model.MenuItemCreateInput) (*model.MenuItemCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemDelete(ctx context.Context, id string) (*model.MenuItemDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemBulkDelete(ctx context.Context, ids []*string) (*model.MenuItemBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemUpdate(ctx context.Context, id string, input model.MenuItemInput) (*model.MenuItemUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemTranslate(ctx context.Context, id string, input model.NameTranslationInput, languageCode model.LanguageCodeEnum) (*model.MenuItemTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemMove(ctx context.Context, menu string, moves []*model.MenuItemMoveInput) (*model.MenuItemMove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceRequest(ctx context.Context, number *string, orderID string) (*model.InvoiceRequest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceRequestDelete(ctx context.Context, id string) (*model.InvoiceRequestDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceCreate(ctx context.Context, input model.InvoiceCreateInput, orderID string) (*model.InvoiceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceDelete(ctx context.Context, id string) (*model.InvoiceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceUpdate(ctx context.Context, id string, input model.UpdateInvoiceInput) (*model.InvoiceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceSendNotification(ctx context.Context, id string) (*model.InvoiceSendNotification, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardActivate(ctx context.Context, id string) (*model.GiftCardActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardCreate(ctx context.Context, input model.GiftCardCreateInput) (*model.GiftCardCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardDeactivate(ctx context.Context, id string) (*model.GiftCardDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardUpdate(ctx context.Context, id string, input model.GiftCardUpdateInput) (*model.GiftCardUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PluginUpdate(ctx context.Context, channel *string, id string, input model.PluginUpdateInput) (*model.PluginUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleCreate(ctx context.Context, input model.SaleInput) (*model.SaleCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleDelete(ctx context.Context, id string) (*model.SaleDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleBulkDelete(ctx context.Context, ids []*string) (*model.SaleBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleUpdate(ctx context.Context, id string, input model.SaleInput) (*model.SaleUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleCataloguesAdd(ctx context.Context, id string, input model.CatalogueInput) (*model.SaleAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleCataloguesRemove(ctx context.Context, id string, input model.CatalogueInput) (*model.SaleRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleTranslate(ctx context.Context, id string, input model.NameTranslationInput, languageCode model.LanguageCodeEnum) (*model.SaleTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleChannelListingUpdate(ctx context.Context, id string, input model.SaleChannelListingInput) (*model.SaleChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherCreate(ctx context.Context, input model.VoucherInput) (*model.VoucherCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherDelete(ctx context.Context, id string) (*model.VoucherDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherBulkDelete(ctx context.Context, ids []*string) (*model.VoucherBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherUpdate(ctx context.Context, id string, input model.VoucherInput) (*model.VoucherUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherCataloguesAdd(ctx context.Context, id string, input model.CatalogueInput) (*model.VoucherAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherCataloguesRemove(ctx context.Context, id string, input model.CatalogueInput) (*model.VoucherRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherTranslate(ctx context.Context, id string, input model.NameTranslationInput, languageCode model.LanguageCodeEnum) (*model.VoucherTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherChannelListingUpdate(ctx context.Context, id string, input model.VoucherChannelListingInput) (*model.VoucherChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExportProducts(ctx context.Context, input model.ExportProductsInput) (*model.ExportProducts, error) {
	return r.exportProducts(ctx, input) // done
}

func (r *mutationResolver) FileUpload(ctx context.Context, file graphql.Upload) (*model.FileUpload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutAddPromoCode(ctx context.Context, checkoutID string, promoCode string) (*model.CheckoutAddPromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutBillingAddressUpdate(ctx context.Context, billingAddress model.AddressInput, checkoutID string) (*model.CheckoutBillingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutComplete(ctx context.Context, checkoutID string, paymentData *string, redirectURL *string, storeSource *bool) (*model.CheckoutComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCreate(ctx context.Context, input model.CheckoutCreateInput) (*model.CheckoutCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerAttach(ctx context.Context, checkoutID string) (*model.CheckoutCustomerAttach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerDetach(ctx context.Context, checkoutID string) (*model.CheckoutCustomerDetach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutEmailUpdate(ctx context.Context, checkoutID *string, email string) (*model.CheckoutEmailUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLineDelete(ctx context.Context, checkoutID string, lineID *string) (*model.CheckoutLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLinesAdd(ctx context.Context, checkoutID string, lines []*model.CheckoutLineInput) (*model.CheckoutLinesAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLinesUpdate(ctx context.Context, checkoutID string, lines []*model.CheckoutLineInput) (*model.CheckoutLinesUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutRemovePromoCode(ctx context.Context, checkoutID string, promoCode string) (*model.CheckoutRemovePromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutPaymentCreate(ctx context.Context, checkoutID string, input model.PaymentInput) (*model.CheckoutPaymentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutShippingAddressUpdate(ctx context.Context, checkoutID string, shippingAddress model.AddressInput) (*model.CheckoutShippingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutShippingMethodUpdate(ctx context.Context, checkoutID *string, shippingMethodID string) (*model.CheckoutShippingMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLanguageCodeUpdate(ctx context.Context, checkoutID string, languageCode model.LanguageCodeEnum) (*model.CheckoutLanguageCodeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelCreate(ctx context.Context, input model.ChannelCreateInput) (*model.ChannelCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelUpdate(ctx context.Context, id string, input model.ChannelUpdateInput) (*model.ChannelUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelDelete(ctx context.Context, id string, input *model.ChannelDeleteInput) (*model.ChannelDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelActivate(ctx context.Context, id string) (*model.ChannelActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelDeactivate(ctx context.Context, id string) (*model.ChannelDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeCreate(ctx context.Context, input model.AttributeCreateInput) (*model.AttributeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeDelete(ctx context.Context, id string) (*model.AttributeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeUpdate(ctx context.Context, id string, input model.AttributeUpdateInput) (*model.AttributeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeTranslate(ctx context.Context, id string, input model.NameTranslationInput, languageCode model.LanguageCodeEnum) (*model.AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeBulkDelete(ctx context.Context, ids []*string) (*model.AttributeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueBulkDelete(ctx context.Context, ids []*string) (*model.AttributeValueBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueCreate(ctx context.Context, attribute string, input model.AttributeValueCreateInput) (*model.AttributeValueCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueDelete(ctx context.Context, id string) (*model.AttributeValueDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueUpdate(ctx context.Context, id string, input model.AttributeValueCreateInput) (*model.AttributeValueUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueTranslate(ctx context.Context, id string, input model.AttributeValueTranslationInput, languageCode model.LanguageCodeEnum) (*model.AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeReorderValues(ctx context.Context, attributeID string, moves []*model.ReorderInput) (*model.AttributeReorderValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppCreate(ctx context.Context, input model.AppInput) (*model.AppCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppUpdate(ctx context.Context, id string, input model.AppInput) (*model.AppUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDelete(ctx context.Context, id string) (*model.AppDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenCreate(ctx context.Context, input model.AppTokenInput) (*model.AppTokenCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenDelete(ctx context.Context, id string) (*model.AppTokenDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenVerify(ctx context.Context, token string) (*model.AppTokenVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppInstall(ctx context.Context, input model.AppInstallInput) (*model.AppInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppRetryInstall(ctx context.Context, activateAfterInstallation *bool, id string) (*model.AppRetryInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDeleteFailedInstallation(ctx context.Context, id string) (*model.AppDeleteFailedInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppFetchManifest(ctx context.Context, manifestURL string) (*model.AppFetchManifest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppActivate(ctx context.Context, id string) (*model.AppActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDeactivate(ctx context.Context, id string) (*model.AppDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokenCreate(ctx context.Context, email string, password string) (*model.CreateToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokenRefresh(ctx context.Context, csrfToken *string, refreshToken *string) (*model.RefreshToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokenVerify(ctx context.Context, token string) (*model.VerifyToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokensDeactivateAll(ctx context.Context) (*model.DeactivateAllUserTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalAuthenticationURL(ctx context.Context, input string, pluginID string) (*model.ExternalAuthenticationURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalObtainAccessTokens(ctx context.Context, input string, pluginID string) (*model.ExternalObtainAccessTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalRefresh(ctx context.Context, input string, pluginID string) (*model.ExternalRefresh, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalLogout(ctx context.Context, input string, pluginID string) (*model.ExternalLogout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalVerify(ctx context.Context, input string, pluginID string) (*model.ExternalVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RequestPasswordReset(ctx context.Context, channel *string, email string, redirectURL string) (*model.RequestPasswordReset, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ConfirmAccount(ctx context.Context, email string, token string) (*model.ConfirmAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetPassword(ctx context.Context, email string, password string, token string) (*model.SetPassword, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PasswordChange(ctx context.Context, newPassword string, oldPassword string) (*model.PasswordChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RequestEmailChange(ctx context.Context, channel *string, newEmail string, password string, redirectURL string) (*model.RequestEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ConfirmEmailChange(ctx context.Context, channel *string, token string) (*model.ConfirmEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressCreate(ctx context.Context, input model.AddressInput, typeArg *model.AddressTypeEnum) (*model.AccountAddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressUpdate(ctx context.Context, id string, input model.AddressInput) (*model.AccountAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressDelete(ctx context.Context, id string) (*model.AccountAddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountSetDefaultAddress(ctx context.Context, id string, typeArg model.AddressTypeEnum) (*model.AccountSetDefaultAddress, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountRegister(ctx context.Context, input model.AccountRegisterInput) (*model.AccountRegister, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountUpdate(ctx context.Context, input model.AccountInput) (*model.AccountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountRequestDeletion(ctx context.Context, channel *string, redirectURL string) (*model.AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountDelete(ctx context.Context, token string) (*model.AccountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressCreate(ctx context.Context, input model.AddressInput, userID string) (*model.AddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressUpdate(ctx context.Context, id string, input model.AddressInput) (*model.AddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressDelete(ctx context.Context, id string) (*model.AddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressSetDefault(ctx context.Context, addressID string, typeArg model.AddressTypeEnum, userID string) (*model.AddressSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerCreate(ctx context.Context, input model.UserCreateInput) (*model.CustomerCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerUpdate(ctx context.Context, id string, input model.CustomerInput) (*model.CustomerUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerDelete(ctx context.Context, id string) (*model.CustomerDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerBulkDelete(ctx context.Context, ids []*string) (*model.CustomerBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffCreate(ctx context.Context, input model.StaffCreateInput) (*model.StaffCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffUpdate(ctx context.Context, id string, input model.StaffUpdateInput) (*model.StaffUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffDelete(ctx context.Context, id string) (*model.StaffDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffBulkDelete(ctx context.Context, ids []*string) (*model.StaffBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserAvatarUpdate(ctx context.Context, image graphql.Upload) (*model.UserAvatarUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserAvatarDelete(ctx context.Context) (*model.UserAvatarDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserBulkSetActive(ctx context.Context, ids []*string, isActive bool) (*model.UserBulkSetActive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupCreate(ctx context.Context, input model.PermissionGroupCreateInput) (*model.PermissionGroupCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupUpdate(ctx context.Context, id string, input model.PermissionGroupUpdateInput) (*model.PermissionGroupUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupDelete(ctx context.Context, id string) (*model.PermissionGroupDelete, error) {
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
	return r.payment(ctx, id) // done
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

func (r *queryResolver) OrderByToken(ctx context.Context, token uuid.UUID) (*model.Order, error) {
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

func (r *queryResolver) Checkout(ctx context.Context, token *uuid.UUID) (*model.Checkout, error) {
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

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
