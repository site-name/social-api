package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

func (r *mutationResolver) ShippingMethodChannelListingUpdate(ctx context.Context, id string, input ShippingMethodChannelListingInput) (*ShippingMethodChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceCreate(ctx context.Context, input ShippingPriceInput) (*ShippingPriceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceDelete(ctx context.Context, id string) (*ShippingPriceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceBulkDelete(ctx context.Context, ids []*string) (*ShippingPriceBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceUpdate(ctx context.Context, id string, input ShippingPriceInput) (*ShippingPriceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceTranslate(ctx context.Context, id string, input ShippingPriceTranslationInput, languageCode LanguageCodeEnum) (*ShippingPriceTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceExcludeProducts(ctx context.Context, id string, input ShippingPriceExcludeProductsInput) (*ShippingPriceExcludeProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingPriceRemoveProductFromExclude(ctx context.Context, id string, products []*string) (*ShippingPriceRemoveProductFromExclude, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneCreate(ctx context.Context, input ShippingZoneCreateInput) (*ShippingZoneCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneDelete(ctx context.Context, id string) (*ShippingZoneDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneBulkDelete(ctx context.Context, ids []*string) (*ShippingZoneBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShippingZoneUpdate(ctx context.Context, id string, input ShippingZoneUpdateInput) (*ShippingZoneUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductAttributeAssign(ctx context.Context, operations []*ProductAttributeAssignInput, productTypeID string) (*ProductAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductAttributeUnassign(ctx context.Context, attributeIds []*string, productTypeID string) (*ProductAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryCreate(ctx context.Context, input CategoryInput, parent *string) (*CategoryCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryDelete(ctx context.Context, id string) (*CategoryDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryBulkDelete(ctx context.Context, ids []*string) (*CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryUpdate(ctx context.Context, id string, input CategoryInput) (*CategoryUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryTranslate(ctx context.Context, id string, input TranslationInput, languageCode LanguageCodeEnum) (*CategoryTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionAddProducts(ctx context.Context, collectionID string, products []*string) (*CollectionAddProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionCreate(ctx context.Context, input CollectionCreateInput) (*CollectionCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionDelete(ctx context.Context, id string) (*CollectionDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionReorderProducts(ctx context.Context, collectionID string, moves []*MoveProductInput) (*CollectionReorderProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionBulkDelete(ctx context.Context, ids []*string) (*CollectionBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionRemoveProducts(ctx context.Context, collectionID string, products []*string) (*CollectionRemoveProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionUpdate(ctx context.Context, id string, input CollectionInput) (*CollectionUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionTranslate(ctx context.Context, id string, input TranslationInput, languageCode LanguageCodeEnum) (*CollectionTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionChannelListingUpdate(ctx context.Context, id string, input CollectionChannelListingUpdateInput) (*CollectionChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductCreate(ctx context.Context, input ProductCreateInput) (*ProductCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductDelete(ctx context.Context, id string) (*ProductDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductBulkDelete(ctx context.Context, ids []*string) (*ProductBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductUpdate(ctx context.Context, id string, input ProductInput) (*ProductUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTranslate(ctx context.Context, id string, input TranslationInput, languageCode LanguageCodeEnum) (*ProductTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductChannelListingUpdate(ctx context.Context, id string, input ProductChannelListingUpdateInput) (*ProductChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaCreate(ctx context.Context, input ProductMediaCreateInput) (*ProductMediaCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantReorder(ctx context.Context, moves []*ReorderInput, productID string) (*ProductVariantReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaDelete(ctx context.Context, id string) (*ProductMediaDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaBulkDelete(ctx context.Context, ids []*string) (*ProductMediaBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaReorder(ctx context.Context, mediaIds []*string, productID string) (*ProductMediaReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductMediaUpdate(ctx context.Context, id string, input ProductMediaUpdateInput) (*ProductMediaUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeCreate(ctx context.Context, input ProductTypeInput) (*ProductTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeDelete(ctx context.Context, id string) (*ProductTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeBulkDelete(ctx context.Context, ids []*string) (*ProductTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeUpdate(ctx context.Context, id string, input ProductTypeInput) (*ProductTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTypeReorderAttributes(ctx context.Context, moves []*ReorderInput, productTypeID string, typeArg ProductAttributeType) (*ProductTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductReorderAttributeValues(ctx context.Context, attributeID string, moves []*ReorderInput, productID string) (*ProductReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentCreate(ctx context.Context, input DigitalContentUploadInput, variantID string) (*DigitalContentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentDelete(ctx context.Context, variantID string) (*DigitalContentDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentUpdate(ctx context.Context, input DigitalContentInput, variantID string) (*DigitalContentUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DigitalContentURLCreate(ctx context.Context, input DigitalContentURLCreateInput) (*DigitalContentURLCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantCreate(ctx context.Context, input ProductVariantCreateInput) (*ProductVariantCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantDelete(ctx context.Context, id string) (*ProductVariantDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantBulkCreate(ctx context.Context, product string, variants []*ProductVariantBulkCreateInput) (*ProductVariantBulkCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantBulkDelete(ctx context.Context, ids []*string) (*ProductVariantBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksCreate(ctx context.Context, stocks []StockInput, variantID string) (*ProductVariantStocksCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksDelete(ctx context.Context, variantID string, warehouseIds []string) (*ProductVariantStocksDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksUpdate(ctx context.Context, stocks []StockInput, variantID string) (*ProductVariantStocksUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantUpdate(ctx context.Context, id string, input ProductVariantInput) (*ProductVariantUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantSetDefault(ctx context.Context, productID string, variantID string) (*ProductVariantSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantTranslate(ctx context.Context, id string, input NameTranslationInput, languageCode LanguageCodeEnum) (*ProductVariantTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantChannelListingUpdate(ctx context.Context, id string, input []ProductVariantChannelListingAddInput) (*ProductVariantChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantReorderAttributeValues(ctx context.Context, attributeID string, moves []*ReorderInput, variantID string) (*ProductVariantReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VariantMediaAssign(ctx context.Context, mediaID string, variantID string) (*VariantMediaAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VariantMediaUnassign(ctx context.Context, mediaID string, variantID string) (*VariantMediaUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentCapture(ctx context.Context, amount *string, paymentID string) (*PaymentCapture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentRefund(ctx context.Context, amount *string, paymentID string) (*PaymentRefund, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentVoid(ctx context.Context, paymentID string) (*PaymentVoid, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PaymentInitialize(ctx context.Context, channel *string, gateway string, paymentData *string) (*PaymentInitialize, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageCreate(ctx context.Context, input PageCreateInput) (*PageCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageDelete(ctx context.Context, id string) (*PageDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageBulkDelete(ctx context.Context, ids []*string) (*PageBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageBulkPublish(ctx context.Context, ids []*string, isPublished bool) (*PageBulkPublish, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageUpdate(ctx context.Context, id string, input PageInput) (*PageUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTranslate(ctx context.Context, id string, input PageTranslationInput, languageCode LanguageCodeEnum) (*PageTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeCreate(ctx context.Context, input PageTypeCreateInput) (*PageTypeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeUpdate(ctx context.Context, id *string, input PageTypeUpdateInput) (*PageTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeDelete(ctx context.Context, id string) (*PageTypeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeBulkDelete(ctx context.Context, ids []string) (*PageTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageAttributeAssign(ctx context.Context, attributeIds []string, pageTypeID string) (*PageAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageAttributeUnassign(ctx context.Context, attributeIds []string, pageTypeID string) (*PageAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTypeReorderAttributes(ctx context.Context, moves []ReorderInput, pageTypeID string) (*PageTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageReorderAttributeValues(ctx context.Context, attributeID string, moves []*ReorderInput, pageID string) (*PageReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderComplete(ctx context.Context, id string) (*DraftOrderComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderCreate(ctx context.Context, input DraftOrderCreateInput) (*DraftOrderCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderDelete(ctx context.Context, id string) (*DraftOrderDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderBulkDelete(ctx context.Context, ids []*string) (*DraftOrderBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderLinesBulkDelete(ctx context.Context, ids []*string) (*DraftOrderLinesBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DraftOrderUpdate(ctx context.Context, id string, input DraftOrderInput) (*DraftOrderUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteMetadata(ctx context.Context, id string, keys []string) (*DeleteMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePrivateMetadata(ctx context.Context, id string, keys []string) (*DeletePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateMetadata(ctx context.Context, id string, input []MetadataInput) (*UpdateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdatePrivateMetadata(ctx context.Context, id string, input []MetadataInput) (*UpdatePrivateMetadata, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AssignNavigation(ctx context.Context, menu *string, navigationType NavigationType) (*AssignNavigation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuCreate(ctx context.Context, input MenuCreateInput) (*MenuCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuDelete(ctx context.Context, id string) (*MenuDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuBulkDelete(ctx context.Context, ids []*string) (*MenuBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuUpdate(ctx context.Context, id string, input MenuInput) (*MenuUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemCreate(ctx context.Context, input MenuItemCreateInput) (*MenuItemCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemDelete(ctx context.Context, id string) (*MenuItemDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemBulkDelete(ctx context.Context, ids []*string) (*MenuItemBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemUpdate(ctx context.Context, id string, input MenuItemInput) (*MenuItemUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemTranslate(ctx context.Context, id string, input NameTranslationInput, languageCode LanguageCodeEnum) (*MenuItemTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MenuItemMove(ctx context.Context, menu string, moves []*MenuItemMoveInput) (*MenuItemMove, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceRequest(ctx context.Context, number *string, orderID string) (*InvoiceRequest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceRequestDelete(ctx context.Context, id string) (*InvoiceRequestDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceCreate(ctx context.Context, input InvoiceCreateInput, orderID string) (*InvoiceCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceDelete(ctx context.Context, id string) (*InvoiceDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceUpdate(ctx context.Context, id string, input UpdateInvoiceInput) (*InvoiceUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) InvoiceSendNotification(ctx context.Context, id string) (*InvoiceSendNotification, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardActivate(ctx context.Context, id string) (*GiftCardActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardCreate(ctx context.Context, input GiftCardCreateInput) (*GiftCardCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardDeactivate(ctx context.Context, id string) (*GiftCardDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) GiftCardUpdate(ctx context.Context, id string, input GiftCardUpdateInput) (*GiftCardUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PluginUpdate(ctx context.Context, channelID *string, id string, input PluginUpdateInput) (*PluginUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleCreate(ctx context.Context, input SaleInput) (*SaleCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleDelete(ctx context.Context, id string) (*SaleDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleBulkDelete(ctx context.Context, ids []*string) (*SaleBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleUpdate(ctx context.Context, id string, input SaleInput) (*SaleUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleCataloguesAdd(ctx context.Context, id string, input CatalogueInput) (*SaleAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleCataloguesRemove(ctx context.Context, id string, input CatalogueInput) (*SaleRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleTranslate(ctx context.Context, id string, input NameTranslationInput, languageCode LanguageCodeEnum) (*SaleTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleChannelListingUpdate(ctx context.Context, id string, input SaleChannelListingInput) (*SaleChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherCreate(ctx context.Context, input VoucherInput) (*VoucherCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherDelete(ctx context.Context, id string) (*VoucherDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherBulkDelete(ctx context.Context, ids []*string) (*VoucherBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherUpdate(ctx context.Context, id string, input VoucherInput) (*VoucherUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherCataloguesAdd(ctx context.Context, id string, input CatalogueInput) (*VoucherAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherCataloguesRemove(ctx context.Context, id string, input CatalogueInput) (*VoucherRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherTranslate(ctx context.Context, id string, input NameTranslationInput, languageCode LanguageCodeEnum) (*VoucherTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherChannelListingUpdate(ctx context.Context, id string, input VoucherChannelListingInput) (*VoucherChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExportProducts(ctx context.Context, input ExportProductsInput) (*ExportProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) FileUpload(ctx context.Context, file graphql.Upload) (*FileUpload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutAddPromoCode(ctx context.Context, checkoutID string, promoCode string) (*CheckoutAddPromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutBillingAddressUpdate(ctx context.Context, billingAddress AddressInput, checkoutID string) (*CheckoutBillingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutComplete(ctx context.Context, checkoutID string, paymentData *string, redirectURL *string, storeSource *bool) (*CheckoutComplete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCreate(ctx context.Context, input CheckoutCreateInput) (*CheckoutCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerAttach(ctx context.Context, checkoutID string) (*CheckoutCustomerAttach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutCustomerDetach(ctx context.Context, checkoutID string) (*CheckoutCustomerDetach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutEmailUpdate(ctx context.Context, checkoutID *string, email string) (*CheckoutEmailUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLineDelete(ctx context.Context, checkoutID string, lineID *string) (*CheckoutLineDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLinesAdd(ctx context.Context, checkoutID string, lines []*CheckoutLineInput) (*CheckoutLinesAdd, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLinesUpdate(ctx context.Context, checkoutID string, lines []*CheckoutLineInput) (*CheckoutLinesUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutRemovePromoCode(ctx context.Context, checkoutID string, promoCode string) (*CheckoutRemovePromoCode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutPaymentCreate(ctx context.Context, checkoutID string, input PaymentInput) (*CheckoutPaymentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutShippingAddressUpdate(ctx context.Context, checkoutID string, shippingAddress AddressInput) (*CheckoutShippingAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutShippingMethodUpdate(ctx context.Context, checkoutID *string, shippingMethodID string) (*CheckoutShippingMethodUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CheckoutLanguageCodeUpdate(ctx context.Context, checkoutID string, languageCode LanguageCodeEnum) (*CheckoutLanguageCodeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelCreate(ctx context.Context, input ChannelCreateInput) (*ChannelCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelUpdate(ctx context.Context, id string, input ChannelUpdateInput) (*ChannelUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelDelete(ctx context.Context, id string, input *ChannelDeleteInput) (*ChannelDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelActivate(ctx context.Context, id string) (*ChannelActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelDeactivate(ctx context.Context, id string) (*ChannelDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeCreate(ctx context.Context, input AttributeCreateInput) (*AttributeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeDelete(ctx context.Context, id string) (*AttributeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeUpdate(ctx context.Context, id string, input AttributeUpdateInput) (*AttributeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeTranslate(ctx context.Context, id string, input NameTranslationInput, languageCode LanguageCodeEnum) (*AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeBulkDelete(ctx context.Context, ids []*string) (*AttributeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueBulkDelete(ctx context.Context, ids []*string) (*AttributeValueBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueCreate(ctx context.Context, attribute string, input AttributeValueCreateInput) (*AttributeValueCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueDelete(ctx context.Context, id string) (*AttributeValueDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueUpdate(ctx context.Context, id string, input AttributeValueCreateInput) (*AttributeValueUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueTranslate(ctx context.Context, id string, input AttributeValueTranslationInput, languageCode LanguageCodeEnum) (*AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeReorderValues(ctx context.Context, attributeID string, moves []*ReorderInput) (*AttributeReorderValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppCreate(ctx context.Context, input AppInput) (*AppCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppUpdate(ctx context.Context, id string, input AppInput) (*AppUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDelete(ctx context.Context, id string) (*AppDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenCreate(ctx context.Context, input AppTokenInput) (*AppTokenCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenDelete(ctx context.Context, id string) (*AppTokenDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppTokenVerify(ctx context.Context, token string) (*AppTokenVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppInstall(ctx context.Context, input AppInstallInput) (*AppInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppRetryInstall(ctx context.Context, activateAfterInstallation *bool, id string) (*AppRetryInstall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDeleteFailedInstallation(ctx context.Context, id string) (*AppDeleteFailedInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppFetchManifest(ctx context.Context, manifestURL string) (*AppFetchManifest, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppActivate(ctx context.Context, id string) (*AppActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AppDeactivate(ctx context.Context, id string) (*AppDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokenCreate(ctx context.Context, email string, password string) (*CreateToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokenRefresh(ctx context.Context, csrfToken *string, refreshToken *string) (*RefreshToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokenVerify(ctx context.Context, token string) (*VerifyToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokensDeactivateAll(ctx context.Context) (*DeactivateAllUserTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalAuthenticationURL(ctx context.Context, input string, pluginID string) (*ExternalAuthenticationURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalObtainAccessTokens(ctx context.Context, input string, pluginID string) (*ExternalObtainAccessTokens, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalRefresh(ctx context.Context, input string, pluginID string) (*ExternalRefresh, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalLogout(ctx context.Context, input string, pluginID string) (*ExternalLogout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalVerify(ctx context.Context, input string, pluginID string) (*ExternalVerify, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RequestPasswordReset(ctx context.Context, channel *string, email string, redirectURL string) (*RequestPasswordReset, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ConfirmAccount(ctx context.Context, email string, token string) (*ConfirmAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetPassword(ctx context.Context, email string, password string, token string) (*SetPassword, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PasswordChange(ctx context.Context, newPassword string, oldPassword string) (*PasswordChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RequestEmailChange(ctx context.Context, channel *string, newEmail string, password string, redirectURL string) (*RequestEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ConfirmEmailChange(ctx context.Context, channel *string, token string) (*ConfirmEmailChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressCreate(ctx context.Context, input AddressInput, typeArg *AddressTypeEnum) (*AccountAddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressUpdate(ctx context.Context, id string, input AddressInput) (*AccountAddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountAddressDelete(ctx context.Context, id string) (*AccountAddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountSetDefaultAddress(ctx context.Context, id string, typeArg AddressTypeEnum) (*AccountSetDefaultAddress, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountRegister(ctx context.Context, input AccountRegisterInput) (*AccountRegister, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountUpdate(ctx context.Context, input AccountInput) (*AccountUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountRequestDeletion(ctx context.Context, channel *string, redirectURL string) (*AccountRequestDeletion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AccountDelete(ctx context.Context, token string) (*AccountDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressCreate(ctx context.Context, input AddressInput, userID string) (*AddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressUpdate(ctx context.Context, id string, input AddressInput) (*AddressUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressDelete(ctx context.Context, id string) (*AddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddressSetDefault(ctx context.Context, addressID string, typeArg AddressTypeEnum, userID string) (*AddressSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerCreate(ctx context.Context, input UserCreateInput) (*CustomerCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerUpdate(ctx context.Context, id string, input CustomerInput) (*CustomerUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerDelete(ctx context.Context, id string) (*CustomerDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CustomerBulkDelete(ctx context.Context, ids []*string) (*CustomerBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupCreate(ctx context.Context, input PermissionGroupCreateInput) (*PermissionGroupCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupUpdate(ctx context.Context, id string, input PermissionGroupUpdateInput) (*PermissionGroupUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PermissionGroupDelete(ctx context.Context, id string) (*PermissionGroupDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Translations(ctx context.Context, kind TranslatableKinds, before *string, after *string, first *int, last *int) (*TranslatableItemConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Translation(ctx context.Context, id string, kind TranslatableKinds) (TranslatableItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Stock(ctx context.Context, id string) (*Stock, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Stocks(ctx context.Context, filter *StockFilterInput, before *string, after *string, first *int, last *int) (*StockCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ShippingZone(ctx context.Context, id string, channel *string) (*ShippingZone, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ShippingZones(ctx context.Context, filter *ShippingZoneFilterInput, channel *string, before *string, after *string, first *int, last *int) (*ShippingZoneCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DigitalContent(ctx context.Context, id string) (*DigitalContent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DigitalContents(ctx context.Context, before *string, after *string, first *int, last *int) (*DigitalContentCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Categories(ctx context.Context, filter *CategoryFilterInput, sortBy *CategorySortingInput, level *int, before *string, after *string, first *int, last *int) (*CategoryCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Category(ctx context.Context, id *string, slug *string) (*Category, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Collection(ctx context.Context, id *string, slug *string, channel *string) (*Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Collections(ctx context.Context, filter *CollectionFilterInput, sortBy *CollectionSortingInput, channel *string, before *string, after *string, first *int, last *int) (*CollectionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Product(ctx context.Context, id *string, slug *string, channel *string) (*Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context, filter *ProductFilterInput, sortBy *ProductOrder, channel *string, before *string, after *string, first *int, last *int) (*ProductCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductType(ctx context.Context, id string) (*ProductType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductTypes(ctx context.Context, filter *ProductTypeFilterInput, sortBy *ProductTypeSortingInput, before *string, after *string, first *int, last *int) (*ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductVariant(ctx context.Context, id *string, sku *string, channel *string) (*ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductVariants(ctx context.Context, ids []*string, channel *string, filter *ProductVariantFilterInput, before *string, after *string, first *int, last *int) (*ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ReportProductSales(ctx context.Context, period ReportingPeriod, channel string, before *string, after *string, first *int, last *int) (*ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Payment(ctx context.Context, id string) (*Payment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Payments(ctx context.Context, filter *PaymentFilterInput, before *string, after *string, first *int, last *int) (*PaymentCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Page(ctx context.Context, id *string, slug *string) (*Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Pages(ctx context.Context, sortBy *PageSortingInput, filter *PageFilterInput, before *string, after *string, first *int, last *int) (*PageCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PageType(ctx context.Context, id string) (*PageType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PageTypes(ctx context.Context, sortBy *PageTypeSortingInput, filter *PageTypeFilterInput, before *string, after *string, first *int, last *int) (*PageTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) HomepageEvents(ctx context.Context, before *string, after *string, first *int, last *int) (*OrderEventCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Menu(ctx context.Context, channel *string, id *string, name *string, slug *string) (*Menu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Menus(ctx context.Context, channel *string, sortBy *MenuSortingInput, filter *MenuFilterInput, before *string, after *string, first *int, last *int) (*MenuCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MenuItem(ctx context.Context, id string, channel *string) (*MenuItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MenuItems(ctx context.Context, channel *string, sortBy *MenuItemSortingInput, filter *MenuItemFilterInput, before *string, after *string, first *int, last *int) (*MenuItemCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GiftCard(ctx context.Context, id string) (*GiftCard, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GiftCards(ctx context.Context, before *string, after *string, first *int, last *int) (*GiftCardCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Plugin(ctx context.Context, id string) (*Plugin, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Plugins(ctx context.Context, filter *PluginFilterInput, sortBy *PluginSortingInput, before *string, after *string, first *int, last *int) (*PluginCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sale(ctx context.Context, id string, channel *string) (*Sale, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sales(ctx context.Context, filter *SaleFilterInput, sortBy *SaleSortingInput, query *string, channel *string, before *string, after *string, first *int, last *int) (*SaleCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Voucher(ctx context.Context, id string, channel *string) (*Voucher, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Vouchers(ctx context.Context, filter *VoucherFilterInput, sortBy *VoucherSortingInput, query *string, channel *string, before *string, after *string, first *int, last *int) (*VoucherCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ExportFile(ctx context.Context, id string) (*ExportFile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ExportFiles(ctx context.Context, filter *ExportFileFilterInput, sortBy *ExportFileSortingInput, before *string, after *string, first *int, last *int) (*ExportFileCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TaxTypes(ctx context.Context) ([]*TaxType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkout(ctx context.Context, token *uuid.UUID) (*Checkout, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Checkouts(ctx context.Context, channel *string, before *string, after *string, first *int, last *int) (*CheckoutCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CheckoutLine(ctx context.Context, id *string) (*CheckoutLine, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CheckoutLines(ctx context.Context, before *string, after *string, first *int, last *int) (*CheckoutLineCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Channel(ctx context.Context, id *string) (*Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Channels(ctx context.Context) ([]Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attributes(ctx context.Context, filter *AttributeFilterInput, sortBy *AttributeSortingInput, before *string, after *string, first *int, last *int) (*AttributeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attribute(ctx context.Context, id *string, slug *string) (*Attribute, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AppsInstallations(ctx context.Context) ([]AppInstallation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Apps(ctx context.Context, filter *AppFilterInput, sortBy *AppSortingInput, before *string, after *string, first *int, last *int) (*AppCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) App(ctx context.Context, id *string) (*App, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AddressValidationRules(ctx context.Context, countryCode CountryCode, countryArea *string, city *string, cityArea *string) (*AddressValidationData, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Address(ctx context.Context, id string) (*Address, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Customers(ctx context.Context, filter *CustomerFilterInput, sortBy *UserSortingInput, before *string, after *string, first *int, last *int) (*UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PermissionGroups(ctx context.Context, filter *PermissionGroupFilterInput, sortBy *PermissionGroupSortingInput, before *string, after *string, first *int, last *int) (*GroupCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PermissionGroup(ctx context.Context, id string) (*Group, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id *string, email *string) (*User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
