package model

var (
	PermissionCreateWarehouse *Permission
	PermissionReadWarehouse   *Permission
	PermissionUpdateWarehouse *Permission
	PermissionDeleteWarehouse *Permission

	PermissionCreateAssignedPageAttribute *Permission
	PermissionReadAssignedPageAttribute   *Permission
	PermissionUpdateAssignedPageAttribute *Permission
	PermissionDeleteAssignedPageAttribute *Permission

	PermissionCreateSaleChannelListing *Permission
	PermissionReadSaleChannelListing   *Permission
	PermissionUpdateSaleChannelListing *Permission
	PermissionDeleteSaleChannelListing *Permission

	PermissionCreateShop *Permission // system user can create shops
	PermissionReadShop   *Permission
	PermissionUpdateShop *Permission
	PermissionDeleteShop *Permission

	PermissionCreatePageTranslation *Permission
	PermissionReadPageTranslation   *Permission
	PermissionUpdatePageTranslation *Permission
	PermissionDeletePageTranslation *Permission

	PermissionCreateCategoryTranslation *Permission // system scoped since only system admin can perform CRUD on categories
	PermissionReadCategoryTranslation   *Permission // system scoped since only system admin can perform CRUD on categories
	PermissionUpdateCategoryTranslation *Permission // system scoped since only system admin can perform CRUD on categories
	PermissionDeleteCategoryTranslation *Permission // system scoped since only system admin can perform CRUD on categories

	PermissionCreateStock *Permission
	PermissionReadStock   *Permission
	PermissionUpdateStock *Permission
	PermissionDeleteStock *Permission

	PermissionCreateOpenExchangeRate *Permission // system scoped since only system admin can perform CRUD
	PermissionReadOpenExchangeRate   *Permission // system scoped since only system admin can perform CRUD
	PermissionUpdateOpenExchangeRate *Permission // system scoped since only system admin can perform CRUD
	PermissionDeleteOpenExchangeRate *Permission // system scoped since only system admin can perform CRUD

	PermissionCreateAssignedVariantAttributeValue *Permission
	PermissionReadAssignedVariantAttributeValue   *Permission
	PermissionUpdateAssignedVariantAttributeValue *Permission
	PermissionDeleteAssignedVariantAttributeValue *Permission

	PermissionCreateShippingZone *Permission
	PermissionReadShippingZone   *Permission
	PermissionUpdateShippingZone *Permission
	PermissionDeleteShippingZone *Permission

	PermissionCreateWishlistItemProductVariant *Permission
	PermissionReadWishlistItemProductVariant   *Permission
	PermissionUpdateWishlistItemProductVariant *Permission
	PermissionDeleteWishlistItemProductVariant *Permission

	PermissionCreateTransaction *Permission
	PermissionReadTransaction   *Permission
	PermissionUpdateTransaction *Permission
	PermissionDeleteTransaction *Permission

	PermissionCreateAttributeValue *Permission
	PermissionReadAttributeValue   *Permission
	PermissionUpdateAttributeValue *Permission
	PermissionDeleteAttributeValue *Permission

	PermissionCreateAttributePage *Permission
	PermissionReadAttributePage   *Permission
	PermissionUpdateAttributePage *Permission
	PermissionDeleteAttributePage *Permission

	PermissionCreateSaleProductVariant *Permission
	PermissionReadSaleProductVariant   *Permission
	PermissionUpdateSaleProductVariant *Permission
	PermissionDeleteSaleProductVariant *Permission

	PermissionCreateOrderEvent *Permission
	PermissionReadOrderEvent   *Permission
	PermissionUpdateOrderEvent *Permission
	PermissionDeleteOrderEvent *Permission

	PermissionCreatePreOrderAllocation *Permission
	PermissionReadPreOrderAllocation   *Permission
	PermissionUpdatePreOrderAllocation *Permission
	PermissionDeletePreOrderAllocation *Permission

	PermissionCreateCustomerEvent *Permission
	PermissionReadCustomerEvent   *Permission
	PermissionUpdateCustomerEvent *Permission
	PermissionDeleteCustomerEvent *Permission

	PermissionCreateAttributeProduct *Permission
	PermissionReadAttributeProduct   *Permission
	PermissionUpdateAttributeProduct *Permission
	PermissionDeleteAttributeProduct *Permission

	PermissionCreateCsvExportEvent *Permission
	PermissionReadCsvExportEvent   *Permission
	PermissionUpdateCsvExportEvent *Permission
	PermissionDeleteCsvExportEvent *Permission

	PermissionCreateCsvExportFile *Permission
	PermissionReadCsvExportFile   *Permission
	PermissionUpdateCsvExportFile *Permission
	PermissionDeleteCsvExportFile *Permission

	PermissionCreateCollectionProductRelation *Permission
	PermissionReadCollectionProductRelation   *Permission
	PermissionUpdateCollectionProductRelation *Permission
	PermissionDeleteCollectionProductRelation *Permission

	PermissionCreateCustomerNote *Permission
	PermissionReadCustomerNote   *Permission
	PermissionUpdateCustomerNote *Permission
	PermissionDeleteCustomerNote *Permission

	PermissionCreateMenuItemTranslation *Permission
	PermissionReadMenuItemTranslation   *Permission
	PermissionUpdateMenuItemTranslation *Permission
	PermissionDeleteMenuItemTranslation *Permission

	PermissionCreateOrderGiftCard *Permission
	PermissionReadOrderGiftCard   *Permission
	PermissionUpdateOrderGiftCard *Permission
	PermissionDeleteOrderGiftCard *Permission

	PermissionCreateCheckoutLine *Permission
	PermissionReadCheckoutLine   *Permission
	PermissionUpdateCheckoutLine *Permission
	PermissionDeleteCheckoutLine *Permission

	PermissionCreateSaleCollectionRelation *Permission
	PermissionReadSaleCollectionRelation   *Permission
	PermissionUpdateSaleCollectionRelation *Permission
	PermissionDeleteSaleCollectionRelation *Permission

	PermissionCreateSaleProductRelation *Permission
	PermissionReadSaleProductRelation   *Permission
	PermissionUpdateSaleProductRelation *Permission
	PermissionDeleteSaleProductRelation *Permission

	PermissionCreateShopStaff *Permission
	PermissionReadShopStaff   *Permission
	PermissionUpdateShopStaff *Permission
	PermissionDeleteShopStaff *Permission

	PermissionCreateProductTranslation *Permission
	PermissionReadProductTranslation   *Permission
	PermissionUpdateProductTranslation *Permission
	PermissionDeleteProductTranslation *Permission

	PermissionCreateWarehouseShippingZone *Permission
	PermissionReadWarehouseShippingZone   *Permission
	PermissionUpdateWarehouseShippingZone *Permission
	PermissionDeleteWarehouseShippingZone *Permission

	PermissionCreatePluginConfiguration *Permission // each shop has their own configuration
	PermissionReadPluginConfiguration   *Permission // each shop has their own configuration
	PermissionUpdatePluginConfiguration *Permission // each shop has their own configuration
	PermissionDeletePluginConfiguration *Permission // each shop has their own configuration

	PermissionCreateAudit *Permission
	PermissionReadAudit   *Permission
	PermissionUpdateAudit *Permission
	PermissionDeleteAudit *Permission

	PermissionCreateProductChannelListing *Permission
	PermissionReadProductChannelListing   *Permission
	PermissionUpdateProductChannelListing *Permission
	PermissionDeleteProductChannelListing *Permission

	PermissionCreateCollectionChannelListing *Permission
	PermissionReadCollectionChannelListing   *Permission
	PermissionUpdateCollectionChannelListing *Permission
	PermissionDeleteCollectionChannelListing *Permission

	PermissionCreateVoucherTranslation *Permission
	PermissionReadVoucherTranslation   *Permission
	PermissionUpdateVoucherTranslation *Permission
	PermissionDeleteVoucherTranslation *Permission

	PermissionCreateClusterDiscovery *Permission
	PermissionReadClusterDiscovery   *Permission
	PermissionUpdateClusterDiscovery *Permission
	PermissionDeleteClusterDiscovery *Permission

	PermissionCreateProductVariantTranslation *Permission
	PermissionReadProductVariantTranslation   *Permission
	PermissionUpdateProductVariantTranslation *Permission
	PermissionDeleteProductVariantTranslation *Permission

	PermissionCreateShopTranslation *Permission
	PermissionReadShopTranslation   *Permission
	PermissionUpdateShopTranslation *Permission
	PermissionDeleteShopTranslation *Permission

	PermissionCreateShippingMethodChannelListing *Permission
	PermissionReadShippingMethodChannelListing   *Permission
	PermissionUpdateShippingMethodChannelListing *Permission
	PermissionDeleteShippingMethodChannelListing *Permission

	PermissionCreateRole *Permission
	PermissionReadRole   *Permission
	PermissionUpdateRole *Permission
	PermissionDeleteRole *Permission

	PermissionCreateAssignedVariantAttribute *Permission
	PermissionReadAssignedVariantAttribute   *Permission
	PermissionUpdateAssignedVariantAttribute *Permission
	PermissionDeleteAssignedVariantAttribute *Permission

	PermissionCreateCompliance                 *Permission
	PermissionReadCompliance                   *Permission
	PermissionUpdateCompliance                 *Permission
	PermissionDeleteCompliance                 *Permission
	PermissionCreateStaffNotificationRecipient *Permission
	PermissionReadStaffNotificationRecipient   *Permission
	PermissionUpdateStaffNotificationRecipient *Permission
	PermissionDeleteStaffNotificationRecipient *Permission
	PermissionCreatePluginKeyValueStore        *Permission
	PermissionReadPluginKeyValueStore          *Permission
	PermissionUpdatePluginKeyValueStore        *Permission
	PermissionDeletePluginKeyValueStore        *Permission

	PermissionCreateChannel *Permission // system admin can do this only
	PermissionReadChannel   *Permission // system admin can do this only
	PermissionUpdateChannel *Permission // system admin can do this only
	PermissionDeleteChannel *Permission // system admin can do this only

	PermissionCreateFulfillmentLine *Permission
	PermissionReadFulfillmentLine   *Permission
	PermissionUpdateFulfillmentLine *Permission
	PermissionDeleteFulfillmentLine *Permission

	PermissionCreateVoucherCollection *Permission
	PermissionReadVoucherCollection   *Permission
	PermissionUpdateVoucherCollection *Permission
	PermissionDeleteVoucherCollection *Permission

	PermissionCreateVoucherProduct *Permission
	PermissionReadVoucherProduct   *Permission
	PermissionUpdateVoucherProduct *Permission
	PermissionDeleteVoucherProduct *Permission

	PermissionCreateVoucherProductVariant *Permission
	PermissionReadVoucherProductVariant   *Permission
	PermissionUpdateVoucherProductVariant *Permission
	PermissionDeleteVoucherProductVariant *Permission

	PermissionCreateFulfillment *Permission
	PermissionReadFulfillment   *Permission
	PermissionUpdateFulfillment *Permission
	PermissionDeleteFulfillment *Permission

	PermissionCreateProduct *Permission
	PermissionReadProduct   *Permission
	PermissionUpdateProduct *Permission
	PermissionDeleteProduct *Permission

	PermissionCreateTermsOfService *Permission
	PermissionReadTermsOfService   *Permission
	PermissionUpdateTermsOfService *Permission
	PermissionDeleteTermsOfService *Permission

	PermissionCreateAssignedProductAttributeValue *Permission
	PermissionReadAssignedProductAttributeValue   *Permission
	PermissionUpdateAssignedProductAttributeValue *Permission
	PermissionDeleteAssignedProductAttributeValue *Permission

	PermissionCreateOrderDiscount *Permission
	PermissionReadOrderDiscount   *Permission
	PermissionUpdateOrderDiscount *Permission
	PermissionDeleteOrderDiscount *Permission

	PermissionCreateProductVariantMedia *Permission
	PermissionReadProductVariantMedia   *Permission
	PermissionUpdateProductVariantMedia *Permission
	PermissionDeleteProductVariantMedia *Permission

	PermissionCreateAttributeTranslation *Permission
	PermissionReadAttributeTranslation   *Permission
	PermissionUpdateAttributeTranslation *Permission
	PermissionDeleteAttributeTranslation *Permission

	PermissionCreateAttributeValueTranslation *Permission
	PermissionReadAttributeValueTranslation   *Permission
	PermissionUpdateAttributeValueTranslation *Permission
	PermissionDeleteAttributeValueTranslation *Permission

	PermissionCreateGiftcard *Permission
	PermissionReadGiftcard   *Permission
	PermissionUpdateGiftcard *Permission
	PermissionDeleteGiftcard *Permission

	PermissionCreatePayment *Permission
	PermissionReadPayment   *Permission
	PermissionUpdatePayment *Permission
	PermissionDeletePayment *Permission

	PermissionCreateToken *Permission
	PermissionReadToken   *Permission
	PermissionUpdateToken *Permission
	PermissionDeleteToken *Permission

	PermissionCreateAttribute *Permission
	PermissionReadAttribute   *Permission
	PermissionUpdateAttribute *Permission
	PermissionDeleteAttribute *Permission

	PermissionCreateSale *Permission
	PermissionReadSale   *Permission
	PermissionUpdateSale *Permission
	PermissionDeleteSale *Permission

	PermissionCreateShippingMethod *Permission
	PermissionReadShippingMethod   *Permission
	PermissionUpdateShippingMethod *Permission
	PermissionDeleteShippingMethod *Permission

	PermissionCreateShippingMethodPostalCodeRule *Permission
	PermissionReadShippingMethodPostalCodeRule   *Permission
	PermissionUpdateShippingMethodPostalCodeRule *Permission
	PermissionDeleteShippingMethodPostalCodeRule *Permission

	PermissionCreateCheckout *Permission
	PermissionReadCheckout   *Permission
	PermissionUpdateCheckout *Permission
	PermissionDeleteCheckout *Permission

	PermissionCreateAllocation *Permission
	PermissionReadAllocation   *Permission
	PermissionUpdateAllocation *Permission
	PermissionDeleteAllocation *Permission

	PermissionCreateVoucher *Permission
	PermissionReadVoucher   *Permission
	PermissionUpdateVoucher *Permission
	PermissionDeleteVoucher *Permission

	PermissionCreateMenuItem *Permission
	PermissionReadMenuItem   *Permission
	PermissionUpdateMenuItem *Permission
	PermissionDeleteMenuItem *Permission

	PermissionCreateProductMedia *Permission
	PermissionReadProductMedia   *Permission
	PermissionUpdateProductMedia *Permission
	PermissionDeleteProductMedia *Permission

	PermissionCreateProductType *Permission
	PermissionReadProductType   *Permission
	PermissionUpdateProductType *Permission
	PermissionDeleteProductType *Permission

	PermissionCreateSaleTranslation *Permission
	PermissionReadSaleTranslation   *Permission
	PermissionUpdateSaleTranslation *Permission
	PermissionDeleteSaleTranslation *Permission

	PermissionCreateShippingMethodExcludedProduct *Permission
	PermissionReadShippingMethodExcludedProduct   *Permission
	PermissionUpdateShippingMethodExcludedProduct *Permission
	PermissionDeleteShippingMethodExcludedProduct *Permission

	PermissionCreateOrderLine *Permission
	PermissionReadOrderLine   *Permission
	PermissionUpdateOrderLine *Permission
	PermissionDeleteOrderLine *Permission

	PermissionCreateUser *Permission
	PermissionReadUser   *Permission
	PermissionDeleteUser *Permission

	PermissionCreateVoucherCustomer *Permission
	PermissionReadVoucherCustomer   *Permission
	PermissionUpdateVoucherCustomer *Permission
	PermissionDeleteVoucherCustomer *Permission

	PermissionCreateCollection *Permission
	PermissionReadCollection   *Permission
	PermissionUpdateCollection *Permission
	PermissionDeleteCollection *Permission

	PermissionCreateWishlistItem *Permission
	PermissionReadWishlistItem   *Permission
	PermissionUpdateWishlistItem *Permission
	PermissionDeleteWishlistItem *Permission

	PermissionCreateDigitalContent *Permission
	PermissionReadDigitalContent   *Permission
	PermissionUpdateDigitalContent *Permission
	PermissionDeleteDigitalContent *Permission

	PermissionCreateVoucherChannelListing *Permission
	PermissionReadVoucherChannelListing   *Permission
	PermissionUpdateVoucherChannelListing *Permission
	PermissionDeleteVoucherChannelListing *Permission

	PermissionCreateProductVariant *Permission
	PermissionReadProductVariant   *Permission
	PermissionUpdateProductVariant *Permission
	PermissionDeleteProductVariant *Permission

	PermissionCreateAddress *Permission
	PermissionReadAddress   *Permission
	PermissionUpdateAddress *Permission
	PermissionDeleteAddress *Permission

	PermissionCreateVoucherCategory *Permission
	PermissionReadVoucherCategory   *Permission
	PermissionUpdateVoucherCategory *Permission
	PermissionDeleteVoucherCategory *Permission

	PermissionCreateMenu *Permission
	PermissionReadMenu   *Permission
	PermissionUpdateMenu *Permission
	PermissionDeleteMenu *Permission

	PermissionCreateCollectionTranslation *Permission
	PermissionReadCollectionTranslation   *Permission
	PermissionUpdateCollectionTranslation *Permission
	PermissionDeleteCollectionTranslation *Permission

	PermissionCreateProductVariantChannelListing *Permission
	PermissionReadProductVariantChannelListing   *Permission
	PermissionUpdateProductVariantChannelListing *Permission
	PermissionDeleteProductVariantChannelListing *Permission

	PermissionCreateSaleCategoryRelation *Permission
	PermissionReadSaleCategoryRelation   *Permission
	PermissionUpdateSaleCategoryRelation *Permission
	PermissionDeleteSaleCategoryRelation *Permission

	PermissionCreateOrder *Permission
	PermissionReadOrder   *Permission
	PermissionUpdateOrder *Permission
	PermissionDeleteOrder *Permission

	PermissionCreateCategory *Permission // system_admin manages categories
	PermissionReadCategory   *Permission // system_admin manages categories
	PermissionUpdateCategory *Permission // system_admin manages categories
	PermissionDeleteCategory *Permission // system_admin manages categories

	PermissionCreateShippingMethodTranslation *Permission
	PermissionReadShippingMethodTranslation   *Permission
	PermissionUpdateShippingMethodTranslation *Permission
	PermissionDeleteShippingMethodTranslation *Permission

	PermissionCreateFileInfo *Permission
	PermissionReadFileInfo   *Permission
	PermissionUpdateFileInfo *Permission
	PermissionDeleteFileInfo *Permission

	PermissionCreatePage *Permission
	PermissionReadPage   *Permission
	PermissionUpdatePage *Permission
	PermissionDeletePage *Permission

	PermissionCreateInvoiceEvent *Permission
	PermissionReadInvoiceEvent   *Permission
	PermissionUpdateInvoiceEvent *Permission
	PermissionDeleteInvoiceEvent *Permission

	PermissionCreateInvoice *Permission
	PermissionReadInvoice   *Permission
	PermissionUpdateInvoice *Permission
	PermissionDeleteInvoice *Permission

	PermissionCreateWishlist *Permission
	PermissionReadWishlist   *Permission
	PermissionUpdateWishlist *Permission
	PermissionDeleteWishlist *Permission

	PermissionCreatePreference *Permission
	PermissionReadPreference   *Permission
	PermissionUpdatePreference *Permission
	PermissionDeletePreference *Permission

	PermissionCreateAssignedProductAttribute *Permission
	PermissionReadAssignedProductAttribute   *Permission
	PermissionUpdateAssignedProductAttribute *Permission
	PermissionDeleteAssignedProductAttribute *Permission

	PermissionCreatePageType *Permission
	PermissionReadPageType   *Permission
	PermissionUpdatePageType *Permission
	PermissionDeletePageType *Permission

	PermissionCreateDigitalContentURL *Permission
	PermissionReadDigitalContentURL   *Permission
	PermissionUpdateDigitalContentURL *Permission
	PermissionDeleteDigitalContentURL *Permission

	PermissionCreateAttributeVariant *Permission
	PermissionReadAttributeVariant   *Permission
	PermissionUpdateAttributeVariant *Permission
	PermissionDeleteAttributeVariant *Permission

	PermissionCreateAssignedPageAttributeValue *Permission
	PermissionReadAssignedPageAttributeValue   *Permission
	PermissionUpdateAssignedPageAttributeValue *Permission
	PermissionDeleteAssignedPageAttributeValue *Permission

	PermissionCreateShippingZoneChannel *Permission
	PermissionReadShippingZoneChannel   *Permission
	PermissionUpdateShippingZoneChannel *Permission
	PermissionDeleteShippingZoneChannel *Permission

	PermissionCreateGiftcardEvent *Permission
	PermissionReadGiftcardEvent   *Permission
	PermissionUpdateGiftcardEvent *Permission
	PermissionDeleteGiftcardEvent *Permission

	PermissionCreateGiftcardCheckout *Permission
	PermissionReadGiftcardCheckout   *Permission
	PermissionUpdateGiftcardCheckout *Permission
	PermissionDeleteGiftcardCheckout *Permission
)

// ShopScopedAllPermissions contains all shop-related permissions
var ShopScopedAllPermissions Permissions
var ShopStaffPermissions Permissions
var SystemUserPermissions Permissions
var SystemGuestPermissions Permissions

func initializeShopScopedPermissions() {
	PermissionCreateWarehouse = &Permission{"create_warehouse", "", "", PermissionScopeShop}
	PermissionReadWarehouse = &Permission{"read_warehouse", "", "", PermissionScopeShop}
	PermissionUpdateWarehouse = &Permission{"update_warehouse", "", "", PermissionScopeShop}
	PermissionDeleteWarehouse = &Permission{"delete_warehouse", "", "", PermissionScopeShop}

	PermissionCreateAssignedPageAttribute = &Permission{"create_assignedpageattribute", "", "", PermissionScopeShop}
	PermissionReadAssignedPageAttribute = &Permission{"read_assignedpageattribute", "", "", PermissionScopeShop}
	PermissionUpdateAssignedPageAttribute = &Permission{"update_assignedpageattribute", "", "", PermissionScopeShop}
	PermissionDeleteAssignedPageAttribute = &Permission{"delete_assignedpageattribute", "", "", PermissionScopeShop}

	PermissionCreateSaleChannelListing = &Permission{"create_salechannellisting", "", "", PermissionScopeShop}
	PermissionReadSaleChannelListing = &Permission{"read_salechannellisting", "", "", PermissionScopeShop}
	PermissionUpdateSaleChannelListing = &Permission{"update_salechannellisting", "", "", PermissionScopeShop}
	PermissionDeleteSaleChannelListing = &Permission{"delete_salechannellisting", "", "", PermissionScopeShop}

	PermissionCreateShop = &Permission{"create_shop", "", "", PermissionScopeSystem}
	PermissionReadShop = &Permission{"read_shop", "", "", PermissionScopeShop}
	PermissionUpdateShop = &Permission{"update_shop", "", "", PermissionScopeShop}
	PermissionDeleteShop = &Permission{"delete_shop", "", "", PermissionScopeShop}

	PermissionCreatePageTranslation = &Permission{"create_pagetranslation", "", "", PermissionScopeShop}
	PermissionReadPageTranslation = &Permission{"read_pagetranslation", "", "", PermissionScopeShop}
	PermissionUpdatePageTranslation = &Permission{"update_pagetranslation", "", "", PermissionScopeShop}
	PermissionDeletePageTranslation = &Permission{"delete_pagetranslation", "", "", PermissionScopeShop}

	PermissionCreateCategoryTranslation = &Permission{"create_categorytranslation", "", "", PermissionScopeSystem}
	PermissionReadCategoryTranslation = &Permission{"read_categorytranslation", "", "", PermissionScopeSystem}
	PermissionUpdateCategoryTranslation = &Permission{"update_categorytranslation", "", "", PermissionScopeSystem}
	PermissionDeleteCategoryTranslation = &Permission{"delete_categorytranslation", "", "", PermissionScopeSystem}

	PermissionCreateStock = &Permission{"create_stock", "", "", PermissionScopeShop}
	PermissionReadStock = &Permission{"read_stock", "", "", PermissionScopeShop}
	PermissionUpdateStock = &Permission{"update_stock", "", "", PermissionScopeShop}
	PermissionDeleteStock = &Permission{"delete_stock", "", "", PermissionScopeShop}

	PermissionCreateOpenExchangeRate = &Permission{"create_openexchangerate", "", "", PermissionScopeSystem}
	PermissionReadOpenExchangeRate = &Permission{"read_openexchangerate", "", "", PermissionScopeSystem}
	PermissionUpdateOpenExchangeRate = &Permission{"update_openexchangerate", "", "", PermissionScopeSystem}
	PermissionDeleteOpenExchangeRate = &Permission{"delete_openexchangerate", "", "", PermissionScopeSystem}

	PermissionCreateAssignedVariantAttributeValue = &Permission{"create_assignedvariantattributevalue", "", "", PermissionScopeShop}
	PermissionReadAssignedVariantAttributeValue = &Permission{"read_assignedvariantattributevalue", "", "", PermissionScopeShop}
	PermissionUpdateAssignedVariantAttributeValue = &Permission{"update_assignedvariantattributevalue", "", "", PermissionScopeShop}
	PermissionDeleteAssignedVariantAttributeValue = &Permission{"delete_assignedvariantattributevalue", "", "", PermissionScopeShop}

	PermissionCreateShippingZone = &Permission{"create_shippingzone", "", "", PermissionScopeShop}
	PermissionReadShippingZone = &Permission{"read_shippingzone", "", "", PermissionScopeShop}
	PermissionUpdateShippingZone = &Permission{"update_shippingzone", "", "", PermissionScopeShop}
	PermissionDeleteShippingZone = &Permission{"delete_shippingzone", "", "", PermissionScopeShop}

	PermissionCreateWishlistItemProductVariant = &Permission{"create_wishlistitemproductvariant", "", "", PermissionScopeSystem}
	PermissionReadWishlistItemProductVariant = &Permission{"read_wishlistitemproductvariant", "", "", PermissionScopeSystem}
	PermissionUpdateWishlistItemProductVariant = &Permission{"update_wishlistitemproductvariant", "", "", PermissionScopeSystem}
	PermissionDeleteWishlistItemProductVariant = &Permission{"delete_wishlistitemproductvariant", "", "", PermissionScopeSystem}

	PermissionCreateTransaction = &Permission{"create_transaction", "", "", PermissionScopeSystem}
	PermissionReadTransaction = &Permission{"read_transaction", "", "", PermissionScopeSystem}
	PermissionUpdateTransaction = &Permission{"update_transaction", "", "", PermissionScopeSystem}
	PermissionDeleteTransaction = &Permission{"delete_transaction", "", "", PermissionScopeSystem}

	PermissionCreateAttributeValue = &Permission{"create_attributevalue", "", "", PermissionScopeShop}
	PermissionReadAttributeValue = &Permission{"read_attributevalue", "", "", PermissionScopeShop}
	PermissionUpdateAttributeValue = &Permission{"update_attributevalue", "", "", PermissionScopeShop}
	PermissionDeleteAttributeValue = &Permission{"delete_attributevalue", "", "", PermissionScopeShop}

	PermissionCreateAttributePage = &Permission{"create_attributepage", "", "", PermissionScopeShop}
	PermissionReadAttributePage = &Permission{"read_attributepage", "", "", PermissionScopeShop}
	PermissionUpdateAttributePage = &Permission{"update_attributepage", "", "", PermissionScopeShop}
	PermissionDeleteAttributePage = &Permission{"delete_attributepage", "", "", PermissionScopeShop}

	PermissionCreateSaleProductVariant = &Permission{"create_saleproductvariant", "", "", PermissionScopeShop}
	PermissionReadSaleProductVariant = &Permission{"read_saleproductvariant", "", "", PermissionScopeShop}
	PermissionUpdateSaleProductVariant = &Permission{"update_saleproductvariant", "", "", PermissionScopeShop}
	PermissionDeleteSaleProductVariant = &Permission{"delete_saleproductvariant", "", "", PermissionScopeShop}

	PermissionCreateOrderEvent = &Permission{"create_orderevent", "", "", PermissionScopeShop}
	PermissionReadOrderEvent = &Permission{"read_orderevent", "", "", PermissionScopeShop}
	PermissionUpdateOrderEvent = &Permission{"update_orderevent", "", "", PermissionScopeShop}
	PermissionDeleteOrderEvent = &Permission{"delete_orderevent", "", "", PermissionScopeShop}

	PermissionCreatePreOrderAllocation = &Permission{"create_preorderallocation", "", "", PermissionScopeShop}
	PermissionReadPreOrderAllocation = &Permission{"read_preorderallocation", "", "", PermissionScopeShop}
	PermissionUpdatePreOrderAllocation = &Permission{"update_preorderallocation", "", "", PermissionScopeShop}
	PermissionDeletePreOrderAllocation = &Permission{"delete_preorderallocation", "", "", PermissionScopeShop}

	PermissionCreateCustomerEvent = &Permission{"create_customerevent", "", "", PermissionScopeShop}
	PermissionReadCustomerEvent = &Permission{"read_customerevent", "", "", PermissionScopeShop}
	PermissionUpdateCustomerEvent = &Permission{"update_customerevent", "", "", PermissionScopeShop}
	PermissionDeleteCustomerEvent = &Permission{"delete_customerevent", "", "", PermissionScopeShop}

	PermissionCreateAttributeProduct = &Permission{"create_attributeproduct", "", "", PermissionScopeShop}
	PermissionReadAttributeProduct = &Permission{"read_attributeproduct", "", "", PermissionScopeShop}
	PermissionUpdateAttributeProduct = &Permission{"update_attributeproduct", "", "", PermissionScopeShop}
	PermissionDeleteAttributeProduct = &Permission{"delete_attributeproduct", "", "", PermissionScopeShop}

	PermissionCreateCsvExportEvent = &Permission{"create_csvexportevent", "", "", PermissionScopeShop}
	PermissionReadCsvExportEvent = &Permission{"read_csvexportevent", "", "", PermissionScopeShop}
	PermissionUpdateCsvExportEvent = &Permission{"update_csvexportevent", "", "", PermissionScopeShop}
	PermissionDeleteCsvExportEvent = &Permission{"delete_csvexportevent", "", "", PermissionScopeShop}

	PermissionCreateCsvExportFile = &Permission{"create_csvexportfile", "", "", PermissionScopeShop}
	PermissionReadCsvExportFile = &Permission{"read_csvexportfile", "", "", PermissionScopeShop}
	PermissionUpdateCsvExportFile = &Permission{"update_csvexportfile", "", "", PermissionScopeShop}
	PermissionDeleteCsvExportFile = &Permission{"delete_csvexportfile", "", "", PermissionScopeShop}

	PermissionCreateCollectionProductRelation = &Permission{"create_collectionproductrelation", "", "", PermissionScopeShop}
	PermissionReadCollectionProductRelation = &Permission{"read_collectionproductrelation", "", "", PermissionScopeShop}
	PermissionUpdateCollectionProductRelation = &Permission{"update_collectionproductrelation", "", "", PermissionScopeShop}
	PermissionDeleteCollectionProductRelation = &Permission{"delete_collectionproductrelation", "", "", PermissionScopeShop}

	PermissionCreateCustomerNote = &Permission{"create_customernote", "", "", PermissionScopeShop}
	PermissionReadCustomerNote = &Permission{"read_customernote", "", "", PermissionScopeShop}
	PermissionUpdateCustomerNote = &Permission{"update_customernote", "", "", PermissionScopeShop}
	PermissionDeleteCustomerNote = &Permission{"delete_customernote", "", "", PermissionScopeShop}

	PermissionCreateMenuItemTranslation = &Permission{"create_menuitemtranslation", "", "", PermissionScopeShop}
	PermissionReadMenuItemTranslation = &Permission{"read_menuitemtranslation", "", "", PermissionScopeShop}
	PermissionUpdateMenuItemTranslation = &Permission{"update_menuitemtranslation", "", "", PermissionScopeShop}
	PermissionDeleteMenuItemTranslation = &Permission{"delete_menuitemtranslation", "", "", PermissionScopeShop}

	PermissionCreateOrderGiftCard = &Permission{"create_ordergiftcard", "", "", PermissionScopeShop}
	PermissionReadOrderGiftCard = &Permission{"read_ordergiftcard", "", "", PermissionScopeShop}
	PermissionUpdateOrderGiftCard = &Permission{"update_ordergiftcard", "", "", PermissionScopeShop}
	PermissionDeleteOrderGiftCard = &Permission{"delete_ordergiftcard", "", "", PermissionScopeShop}

	PermissionCreateCheckoutLine = &Permission{"create_checkoutline", "", "", PermissionScopeShop}
	PermissionReadCheckoutLine = &Permission{"read_checkoutline", "", "", PermissionScopeShop}
	PermissionUpdateCheckoutLine = &Permission{"update_checkoutline", "", "", PermissionScopeShop}
	PermissionDeleteCheckoutLine = &Permission{"delete_checkoutline", "", "", PermissionScopeShop}

	PermissionCreateSaleCollectionRelation = &Permission{"create_salecollectionrelation", "", "", PermissionScopeShop}
	PermissionReadSaleCollectionRelation = &Permission{"read_salecollectionrelation", "", "", PermissionScopeShop}
	PermissionUpdateSaleCollectionRelation = &Permission{"update_salecollectionrelation", "", "", PermissionScopeShop}
	PermissionDeleteSaleCollectionRelation = &Permission{"delete_salecollectionrelation", "", "", PermissionScopeShop}

	PermissionCreateSaleProductRelation = &Permission{"create_saleproductrelation", "", "", PermissionScopeShop}
	PermissionReadSaleProductRelation = &Permission{"read_saleproductrelation", "", "", PermissionScopeShop}
	PermissionUpdateSaleProductRelation = &Permission{"update_saleproductrelation", "", "", PermissionScopeShop}
	PermissionDeleteSaleProductRelation = &Permission{"delete_saleproductrelation", "", "", PermissionScopeShop}

	PermissionCreateShopStaff = &Permission{"create_shopstaff", "", "", PermissionScopeShop}
	PermissionReadShopStaff = &Permission{"read_shopstaff", "", "", PermissionScopeShop}
	PermissionUpdateShopStaff = &Permission{"update_shopstaff", "", "", PermissionScopeShop}
	PermissionDeleteShopStaff = &Permission{"delete_shopstaff", "", "", PermissionScopeShop}

	PermissionCreateProductTranslation = &Permission{"create_producttranslation", "", "", PermissionScopeShop}
	PermissionReadProductTranslation = &Permission{"read_producttranslation", "", "", PermissionScopeShop}
	PermissionUpdateProductTranslation = &Permission{"update_producttranslation", "", "", PermissionScopeShop}
	PermissionDeleteProductTranslation = &Permission{"delete_producttranslation", "", "", PermissionScopeShop}

	PermissionCreateWarehouseShippingZone = &Permission{"create_warehouseshippingzone", "", "", PermissionScopeShop}
	PermissionReadWarehouseShippingZone = &Permission{"read_warehouseshippingzone", "", "", PermissionScopeShop}
	PermissionUpdateWarehouseShippingZone = &Permission{"update_warehouseshippingzone", "", "", PermissionScopeShop}
	PermissionDeleteWarehouseShippingZone = &Permission{"delete_warehouseshippingzone", "", "", PermissionScopeShop}

	PermissionCreatePluginConfiguration = &Permission{"create_pluginconfiguration", "", "", PermissionScopeShop}
	PermissionReadPluginConfiguration = &Permission{"read_pluginconfiguration", "", "", PermissionScopeShop}
	PermissionUpdatePluginConfiguration = &Permission{"update_pluginconfiguration", "", "", PermissionScopeShop}
	PermissionDeletePluginConfiguration = &Permission{"delete_pluginconfiguration", "", "", PermissionScopeShop}

	PermissionCreateAudit = &Permission{"create_audit", "", "", PermissionScopeSystem}
	PermissionReadAudit = &Permission{"read_audit", "", "", PermissionScopeSystem}
	PermissionUpdateAudit = &Permission{"update_audit", "", "", PermissionScopeSystem}
	PermissionDeleteAudit = &Permission{"delete_audit", "", "", PermissionScopeSystem}

	PermissionCreateProductChannelListing = &Permission{"create_productchannellisting", "", "", PermissionScopeShop}
	PermissionReadProductChannelListing = &Permission{"read_productchannellisting", "", "", PermissionScopeShop}
	PermissionUpdateProductChannelListing = &Permission{"update_productchannellisting", "", "", PermissionScopeShop}
	PermissionDeleteProductChannelListing = &Permission{"delete_productchannellisting", "", "", PermissionScopeShop}

	PermissionCreateCollectionChannelListing = &Permission{"create_collectionchannellisting", "", "", PermissionScopeShop}
	PermissionReadCollectionChannelListing = &Permission{"read_collectionchannellisting", "", "", PermissionScopeShop}
	PermissionUpdateCollectionChannelListing = &Permission{"update_collectionchannellisting", "", "", PermissionScopeShop}
	PermissionDeleteCollectionChannelListing = &Permission{"delete_collectionchannellisting", "", "", PermissionScopeShop}

	PermissionCreateVoucherTranslation = &Permission{"create_vouchertranslation", "", "", PermissionScopeShop}
	PermissionReadVoucherTranslation = &Permission{"read_vouchertranslation", "", "", PermissionScopeShop}
	PermissionUpdateVoucherTranslation = &Permission{"update_vouchertranslation", "", "", PermissionScopeShop}
	PermissionDeleteVoucherTranslation = &Permission{"delete_vouchertranslation", "", "", PermissionScopeShop}

	PermissionCreateClusterDiscovery = &Permission{"create_clusterdiscovery", "", "", PermissionScopeSystem}
	PermissionReadClusterDiscovery = &Permission{"read_clusterdiscovery", "", "", PermissionScopeSystem}
	PermissionUpdateClusterDiscovery = &Permission{"update_clusterdiscovery", "", "", PermissionScopeSystem}
	PermissionDeleteClusterDiscovery = &Permission{"delete_clusterdiscovery", "", "", PermissionScopeSystem}

	PermissionCreateProductVariantTranslation = &Permission{"create_productvarianttranslation", "", "", PermissionScopeShop}
	PermissionReadProductVariantTranslation = &Permission{"read_productvarianttranslation", "", "", PermissionScopeShop}
	PermissionUpdateProductVariantTranslation = &Permission{"update_productvarianttranslation", "", "", PermissionScopeShop}
	PermissionDeleteProductVariantTranslation = &Permission{"delete_productvarianttranslation", "", "", PermissionScopeShop}

	PermissionCreateShopTranslation = &Permission{"create_shoptranslation", "", "", PermissionScopeShop}
	PermissionReadShopTranslation = &Permission{"read_shoptranslation", "", "", PermissionScopeShop}
	PermissionUpdateShopTranslation = &Permission{"update_shoptranslation", "", "", PermissionScopeShop}
	PermissionDeleteShopTranslation = &Permission{"delete_shoptranslation", "", "", PermissionScopeShop}

	PermissionCreateShippingMethodChannelListing = &Permission{"create_shippingmethodchannellisting", "", "", PermissionScopeShop}
	PermissionReadShippingMethodChannelListing = &Permission{"read_shippingmethodchannellisting", "", "", PermissionScopeShop}
	PermissionUpdateShippingMethodChannelListing = &Permission{"update_shippingmethodchannellisting", "", "", PermissionScopeShop}
	PermissionDeleteShippingMethodChannelListing = &Permission{"delete_shippingmethodchannellisting", "", "", PermissionScopeShop}

	PermissionCreateRole = &Permission{"create_role", "", "", PermissionScopeSystem}
	PermissionReadRole = &Permission{"read_role", "", "", PermissionScopeSystem}
	PermissionUpdateRole = &Permission{"update_role", "", "", PermissionScopeSystem}
	PermissionDeleteRole = &Permission{"delete_role", "", "", PermissionScopeSystem}

	PermissionCreateAssignedVariantAttribute = &Permission{"create_assignedvariantattribute", "", "", PermissionScopeShop}
	PermissionReadAssignedVariantAttribute = &Permission{"read_assignedvariantattribute", "", "", PermissionScopeShop}
	PermissionUpdateAssignedVariantAttribute = &Permission{"update_assignedvariantattribute", "", "", PermissionScopeShop}
	PermissionDeleteAssignedVariantAttribute = &Permission{"delete_assignedvariantattribute", "", "", PermissionScopeShop}

	PermissionCreateCompliance = &Permission{"create_compliance", "", "", PermissionScopeSystem}
	PermissionReadCompliance = &Permission{"read_compliance", "", "", PermissionScopeSystem}
	PermissionUpdateCompliance = &Permission{"update_compliance", "", "", PermissionScopeSystem}
	PermissionDeleteCompliance = &Permission{"delete_compliance", "", "", PermissionScopeSystem}

	PermissionCreateStaffNotificationRecipient = &Permission{"create_staffnotificationrecipient", "", "", PermissionScopeShop}
	PermissionReadStaffNotificationRecipient = &Permission{"read_staffnotificationrecipient", "", "", PermissionScopeShop}
	PermissionUpdateStaffNotificationRecipient = &Permission{"update_staffnotificationrecipient", "", "", PermissionScopeShop}
	PermissionDeleteStaffNotificationRecipient = &Permission{"delete_staffnotificationrecipient", "", "", PermissionScopeShop}

	PermissionCreatePluginKeyValueStore = &Permission{"create_pluginkeyvaluestore", "", "", PermissionScopeSystem}
	PermissionReadPluginKeyValueStore = &Permission{"read_pluginkeyvaluestore", "", "", PermissionScopeSystem}
	PermissionUpdatePluginKeyValueStore = &Permission{"update_pluginkeyvaluestore", "", "", PermissionScopeSystem}
	PermissionDeletePluginKeyValueStore = &Permission{"delete_pluginkeyvaluestore", "", "", PermissionScopeSystem}

	PermissionCreateChannel = &Permission{"create_channel", "", "", PermissionScopeSystem}
	PermissionReadChannel = &Permission{"read_channel", "", "", PermissionScopeSystem}
	PermissionUpdateChannel = &Permission{"update_channel", "", "", PermissionScopeSystem}
	PermissionDeleteChannel = &Permission{"delete_channel", "", "", PermissionScopeSystem}

	PermissionCreateFulfillmentLine = &Permission{"create_fulfillmentline", "", "", PermissionScopeShop}
	PermissionReadFulfillmentLine = &Permission{"read_fulfillmentline", "", "", PermissionScopeShop}
	PermissionUpdateFulfillmentLine = &Permission{"update_fulfillmentline", "", "", PermissionScopeShop}
	PermissionDeleteFulfillmentLine = &Permission{"delete_fulfillmentline", "", "", PermissionScopeShop}

	PermissionCreateVoucherCollection = &Permission{"create_vouchercollection", "", "", PermissionScopeShop}
	PermissionReadVoucherCollection = &Permission{"read_vouchercollection", "", "", PermissionScopeShop}
	PermissionUpdateVoucherCollection = &Permission{"update_vouchercollection", "", "", PermissionScopeShop}
	PermissionDeleteVoucherCollection = &Permission{"delete_vouchercollection", "", "", PermissionScopeShop}

	PermissionCreateVoucherProduct = &Permission{"create_voucherproduct", "", "", PermissionScopeShop}
	PermissionReadVoucherProduct = &Permission{"read_voucherproduct", "", "", PermissionScopeShop}
	PermissionUpdateVoucherProduct = &Permission{"update_voucherproduct", "", "", PermissionScopeShop}
	PermissionDeleteVoucherProduct = &Permission{"delete_voucherproduct", "", "", PermissionScopeShop}

	PermissionCreateVoucherProductVariant = &Permission{"create_voucherproductvariant", "", "", PermissionScopeShop}
	PermissionReadVoucherProductVariant = &Permission{"read_voucherproductvariant", "", "", PermissionScopeShop}
	PermissionUpdateVoucherProductVariant = &Permission{"update_voucherproductvariant", "", "", PermissionScopeShop}
	PermissionDeleteVoucherProductVariant = &Permission{"delete_voucherproductvariant", "", "", PermissionScopeShop}

	PermissionCreateFulfillment = &Permission{"create_fulfillment", "", "", PermissionScopeShop}
	PermissionReadFulfillment = &Permission{"read_fulfillment", "", "", PermissionScopeShop}
	PermissionUpdateFulfillment = &Permission{"update_fulfillment", "", "", PermissionScopeShop}
	PermissionDeleteFulfillment = &Permission{"delete_fulfillment", "", "", PermissionScopeShop}

	PermissionCreateProduct = &Permission{"create_product", "", "", PermissionScopeShop}
	PermissionReadProduct = &Permission{"read_product", "", "", PermissionScopeShop}
	PermissionUpdateProduct = &Permission{"update_product", "", "", PermissionScopeShop}
	PermissionDeleteProduct = &Permission{"delete_product", "", "", PermissionScopeShop}

	PermissionCreateTermsOfService = &Permission{"create_termsofservice", "", "", PermissionScopeSystem}
	PermissionReadTermsOfService = &Permission{"read_termsofservice", "", "", PermissionScopeSystem}
	PermissionUpdateTermsOfService = &Permission{"update_termsofservice", "", "", PermissionScopeSystem}
	PermissionDeleteTermsOfService = &Permission{"delete_termsofservice", "", "", PermissionScopeSystem}

	PermissionCreateAssignedProductAttributeValue = &Permission{"create_assignedproductattributevalue", "", "", PermissionScopeShop}
	PermissionReadAssignedProductAttributeValue = &Permission{"read_assignedproductattributevalue", "", "", PermissionScopeShop}
	PermissionUpdateAssignedProductAttributeValue = &Permission{"update_assignedproductattributevalue", "", "", PermissionScopeShop}
	PermissionDeleteAssignedProductAttributeValue = &Permission{"delete_assignedproductattributevalue", "", "", PermissionScopeShop}

	PermissionCreateOrderDiscount = &Permission{"create_orderdiscount", "", "", PermissionScopeShop}
	PermissionReadOrderDiscount = &Permission{"read_orderdiscount", "", "", PermissionScopeShop}
	PermissionUpdateOrderDiscount = &Permission{"update_orderdiscount", "", "", PermissionScopeShop}
	PermissionDeleteOrderDiscount = &Permission{"delete_orderdiscount", "", "", PermissionScopeShop}

	PermissionCreateProductVariantMedia = &Permission{"create_productvariantmedia", "", "", PermissionScopeShop}
	PermissionReadProductVariantMedia = &Permission{"read_productvariantmedia", "", "", PermissionScopeShop}
	PermissionUpdateProductVariantMedia = &Permission{"update_productvariantmedia", "", "", PermissionScopeShop}
	PermissionDeleteProductVariantMedia = &Permission{"delete_productvariantmedia", "", "", PermissionScopeShop}

	PermissionCreateAttributeTranslation = &Permission{"create_attributetranslation", "", "", PermissionScopeShop}
	PermissionReadAttributeTranslation = &Permission{"read_attributetranslation", "", "", PermissionScopeShop}
	PermissionUpdateAttributeTranslation = &Permission{"update_attributetranslation", "", "", PermissionScopeShop}
	PermissionDeleteAttributeTranslation = &Permission{"delete_attributetranslation", "", "", PermissionScopeShop}

	PermissionCreateAttributeValueTranslation = &Permission{"create_attributevaluetranslation", "", "", PermissionScopeShop}
	PermissionReadAttributeValueTranslation = &Permission{"read_attributevaluetranslation", "", "", PermissionScopeShop}
	PermissionUpdateAttributeValueTranslation = &Permission{"update_attributevaluetranslation", "", "", PermissionScopeShop}
	PermissionDeleteAttributeValueTranslation = &Permission{"delete_attributevaluetranslation", "", "", PermissionScopeShop}

	PermissionCreateGiftcard = &Permission{"create_giftcard", "", "", PermissionScopeShop}
	PermissionReadGiftcard = &Permission{"read_giftcard", "", "", PermissionScopeShop}
	PermissionUpdateGiftcard = &Permission{"update_giftcard", "", "", PermissionScopeShop}
	PermissionDeleteGiftcard = &Permission{"delete_giftcard", "", "", PermissionScopeShop}

	PermissionCreatePayment = &Permission{"create_payment", "", "", PermissionScopeShop}
	PermissionReadPayment = &Permission{"read_payment", "", "", PermissionScopeShop}
	PermissionUpdatePayment = &Permission{"update_payment", "", "", PermissionScopeShop}
	PermissionDeletePayment = &Permission{"delete_payment", "", "", PermissionScopeShop}

	PermissionCreateToken = &Permission{"create_token", "", "", PermissionScopeSystem}
	PermissionReadToken = &Permission{"read_token", "", "", PermissionScopeSystem}
	PermissionUpdateToken = &Permission{"update_token", "", "", PermissionScopeSystem}
	PermissionDeleteToken = &Permission{"delete_token", "", "", PermissionScopeSystem}

	PermissionCreateAttribute = &Permission{"create_attribute", "", "", PermissionScopeShop}
	PermissionReadAttribute = &Permission{"read_attribute", "", "", PermissionScopeShop}
	PermissionUpdateAttribute = &Permission{"update_attribute", "", "", PermissionScopeShop}
	PermissionDeleteAttribute = &Permission{"delete_attribute", "", "", PermissionScopeShop}

	PermissionCreateSale = &Permission{"create_sale", "", "", PermissionScopeShop}
	PermissionReadSale = &Permission{"read_sale", "", "", PermissionScopeShop}
	PermissionUpdateSale = &Permission{"update_sale", "", "", PermissionScopeShop}
	PermissionDeleteSale = &Permission{"delete_sale", "", "", PermissionScopeShop}

	PermissionCreateShippingMethod = &Permission{"create_shippingmethod", "", "", PermissionScopeShop}
	PermissionReadShippingMethod = &Permission{"read_shippingmethod", "", "", PermissionScopeShop}
	PermissionUpdateShippingMethod = &Permission{"update_shippingmethod", "", "", PermissionScopeShop}
	PermissionDeleteShippingMethod = &Permission{"delete_shippingmethod", "", "", PermissionScopeShop}

	PermissionCreateShippingMethodPostalCodeRule = &Permission{"create_shippingmethodpostalcoderule", "", "", PermissionScopeShop}
	PermissionReadShippingMethodPostalCodeRule = &Permission{"read_shippingmethodpostalcoderule", "", "", PermissionScopeShop}
	PermissionUpdateShippingMethodPostalCodeRule = &Permission{"update_shippingmethodpostalcoderule", "", "", PermissionScopeShop}
	PermissionDeleteShippingMethodPostalCodeRule = &Permission{"delete_shippingmethodpostalcoderule", "", "", PermissionScopeShop}

	PermissionCreateCheckout = &Permission{"create_checkout", "", "", PermissionScopeShop}
	PermissionReadCheckout = &Permission{"read_checkout", "", "", PermissionScopeShop}
	PermissionUpdateCheckout = &Permission{"update_checkout", "", "", PermissionScopeShop}
	PermissionDeleteCheckout = &Permission{"delete_checkout", "", "", PermissionScopeShop}

	PermissionCreateAllocation = &Permission{"create_allocation", "", "", PermissionScopeShop}
	PermissionReadAllocation = &Permission{"read_allocation", "", "", PermissionScopeShop}
	PermissionUpdateAllocation = &Permission{"update_allocation", "", "", PermissionScopeShop}
	PermissionDeleteAllocation = &Permission{"delete_allocation", "", "", PermissionScopeShop}

	PermissionCreateVoucher = &Permission{"create_voucher", "", "", PermissionScopeShop}
	PermissionReadVoucher = &Permission{"read_voucher", "", "", PermissionScopeShop}
	PermissionUpdateVoucher = &Permission{"update_voucher", "", "", PermissionScopeShop}
	PermissionDeleteVoucher = &Permission{"delete_voucher", "", "", PermissionScopeShop}

	PermissionCreateMenuItem = &Permission{"create_menuitem", "", "", PermissionScopeShop}
	PermissionReadMenuItem = &Permission{"read_menuitem", "", "", PermissionScopeShop}
	PermissionUpdateMenuItem = &Permission{"update_menuitem", "", "", PermissionScopeShop}
	PermissionDeleteMenuItem = &Permission{"delete_menuitem", "", "", PermissionScopeShop}

	PermissionCreateProductMedia = &Permission{"create_productmedia", "", "", PermissionScopeShop}
	PermissionReadProductMedia = &Permission{"read_productmedia", "", "", PermissionScopeShop}
	PermissionUpdateProductMedia = &Permission{"update_productmedia", "", "", PermissionScopeShop}
	PermissionDeleteProductMedia = &Permission{"delete_productmedia", "", "", PermissionScopeShop}

	PermissionCreateProductType = &Permission{"create_producttype", "", "", PermissionScopeShop}
	PermissionReadProductType = &Permission{"read_producttype", "", "", PermissionScopeShop}
	PermissionUpdateProductType = &Permission{"update_producttype", "", "", PermissionScopeShop}
	PermissionDeleteProductType = &Permission{"delete_producttype", "", "", PermissionScopeShop}

	PermissionCreateSaleTranslation = &Permission{"create_saletranslation", "", "", PermissionScopeShop}
	PermissionReadSaleTranslation = &Permission{"read_saletranslation", "", "", PermissionScopeShop}
	PermissionUpdateSaleTranslation = &Permission{"update_saletranslation", "", "", PermissionScopeShop}
	PermissionDeleteSaleTranslation = &Permission{"delete_saletranslation", "", "", PermissionScopeShop}

	PermissionCreateShippingMethodExcludedProduct = &Permission{"create_shippingmethodexcludedproduct", "", "", PermissionScopeShop}
	PermissionReadShippingMethodExcludedProduct = &Permission{"read_shippingmethodexcludedproduct", "", "", PermissionScopeShop}
	PermissionUpdateShippingMethodExcludedProduct = &Permission{"update_shippingmethodexcludedproduct", "", "", PermissionScopeShop}
	PermissionDeleteShippingMethodExcludedProduct = &Permission{"delete_shippingmethodexcludedproduct", "", "", PermissionScopeShop}

	PermissionCreateOrderLine = &Permission{"create_orderline", "", "", PermissionScopeShop}
	PermissionReadOrderLine = &Permission{"read_orderline", "", "", PermissionScopeShop}
	PermissionUpdateOrderLine = &Permission{"update_orderline", "", "", PermissionScopeShop}
	PermissionDeleteOrderLine = &Permission{"delete_orderline", "", "", PermissionScopeShop}

	PermissionCreateUser = &Permission{"create_user", "", "", PermissionScopeSystem}
	PermissionReadUser = &Permission{"read_user", "", "", PermissionScopeSystem}
	PermissionDeleteUser = &Permission{"delete_user", "", "", PermissionScopeSystem}

	PermissionCreateVoucherCustomer = &Permission{"create_vouchercustomer", "", "", PermissionScopeShop}
	PermissionReadVoucherCustomer = &Permission{"read_vouchercustomer", "", "", PermissionScopeShop}
	PermissionUpdateVoucherCustomer = &Permission{"update_vouchercustomer", "", "", PermissionScopeShop}
	PermissionDeleteVoucherCustomer = &Permission{"delete_vouchercustomer", "", "", PermissionScopeShop}

	PermissionCreateCollection = &Permission{"create_collection", "", "", PermissionScopeShop}
	PermissionReadCollection = &Permission{"read_collection", "", "", PermissionScopeShop}
	PermissionUpdateCollection = &Permission{"update_collection", "", "", PermissionScopeShop}
	PermissionDeleteCollection = &Permission{"delete_collection", "", "", PermissionScopeShop}

	PermissionCreateWishlistItem = &Permission{"create_wishlistitem", "", "", PermissionScopeShop}
	PermissionReadWishlistItem = &Permission{"read_wishlistitem", "", "", PermissionScopeShop}
	PermissionUpdateWishlistItem = &Permission{"update_wishlistitem", "", "", PermissionScopeShop}
	PermissionDeleteWishlistItem = &Permission{"delete_wishlistitem", "", "", PermissionScopeShop}

	PermissionCreateDigitalContent = &Permission{"create_digitalcontent", "", "", PermissionScopeShop}
	PermissionReadDigitalContent = &Permission{"read_digitalcontent", "", "", PermissionScopeShop}
	PermissionUpdateDigitalContent = &Permission{"update_digitalcontent", "", "", PermissionScopeShop}
	PermissionDeleteDigitalContent = &Permission{"delete_digitalcontent", "", "", PermissionScopeShop}

	PermissionCreateVoucherChannelListing = &Permission{"create_voucherchannellisting", "", "", PermissionScopeShop}
	PermissionReadVoucherChannelListing = &Permission{"read_voucherchannellisting", "", "", PermissionScopeShop}
	PermissionUpdateVoucherChannelListing = &Permission{"update_voucherchannellisting", "", "", PermissionScopeShop}
	PermissionDeleteVoucherChannelListing = &Permission{"delete_voucherchannellisting", "", "", PermissionScopeShop}

	PermissionCreateProductVariant = &Permission{"create_productvariant", "", "", PermissionScopeShop}
	PermissionReadProductVariant = &Permission{"read_productvariant", "", "", PermissionScopeShop}
	PermissionUpdateProductVariant = &Permission{"update_productvariant", "", "", PermissionScopeShop}
	PermissionDeleteProductVariant = &Permission{"delete_productvariant", "", "", PermissionScopeShop}

	PermissionCreateAddress = &Permission{"create_address", "", "", PermissionScopeSystem}
	PermissionReadAddress = &Permission{"read_address", "", "", PermissionScopeSystem}
	PermissionUpdateAddress = &Permission{"update_address", "", "", PermissionScopeSystem}
	PermissionDeleteAddress = &Permission{"delete_address", "", "", PermissionScopeSystem}

	PermissionCreateVoucherCategory = &Permission{"create_vouchercategory", "", "", PermissionScopeShop}
	PermissionReadVoucherCategory = &Permission{"read_vouchercategory", "", "", PermissionScopeShop}
	PermissionUpdateVoucherCategory = &Permission{"update_vouchercategory", "", "", PermissionScopeShop}
	PermissionDeleteVoucherCategory = &Permission{"delete_vouchercategory", "", "", PermissionScopeShop}

	PermissionCreateMenu = &Permission{"create_menu", "", "", PermissionScopeShop}
	PermissionReadMenu = &Permission{"read_menu", "", "", PermissionScopeShop}
	PermissionUpdateMenu = &Permission{"update_menu", "", "", PermissionScopeShop}
	PermissionDeleteMenu = &Permission{"delete_menu", "", "", PermissionScopeShop}

	PermissionCreateCollectionTranslation = &Permission{"create_collectiontranslation", "", "", PermissionScopeShop}
	PermissionReadCollectionTranslation = &Permission{"read_collectiontranslation", "", "", PermissionScopeShop}
	PermissionUpdateCollectionTranslation = &Permission{"update_collectiontranslation", "", "", PermissionScopeShop}
	PermissionDeleteCollectionTranslation = &Permission{"delete_collectiontranslation", "", "", PermissionScopeShop}

	PermissionCreateProductVariantChannelListing = &Permission{"create_productvariantchannellisting", "", "", PermissionScopeShop}
	PermissionReadProductVariantChannelListing = &Permission{"read_productvariantchannellisting", "", "", PermissionScopeShop}
	PermissionUpdateProductVariantChannelListing = &Permission{"update_productvariantchannellisting", "", "", PermissionScopeShop}
	PermissionDeleteProductVariantChannelListing = &Permission{"delete_productvariantchannellisting", "", "", PermissionScopeShop}

	PermissionCreateSaleCategoryRelation = &Permission{"create_salecategoryrelation", "", "", PermissionScopeShop}
	PermissionReadSaleCategoryRelation = &Permission{"read_salecategoryrelation", "", "", PermissionScopeShop}
	PermissionUpdateSaleCategoryRelation = &Permission{"update_salecategoryrelation", "", "", PermissionScopeShop}
	PermissionDeleteSaleCategoryRelation = &Permission{"delete_salecategoryrelation", "", "", PermissionScopeShop}

	PermissionCreateOrder = &Permission{"create_order", "", "", PermissionScopeShop}
	PermissionReadOrder = &Permission{"read_order", "", "", PermissionScopeShop}
	PermissionUpdateOrder = &Permission{"update_order", "", "", PermissionScopeShop}
	PermissionDeleteOrder = &Permission{"delete_order", "", "", PermissionScopeShop}

	PermissionCreateCategory = &Permission{"create_category", "", "", PermissionScopeSystem}
	PermissionReadCategory = &Permission{"read_category", "", "", PermissionScopeSystem}
	PermissionUpdateCategory = &Permission{"update_category", "", "", PermissionScopeSystem}
	PermissionDeleteCategory = &Permission{"delete_category", "", "", PermissionScopeSystem}

	PermissionCreateShippingMethodTranslation = &Permission{"create_shippingmethodtranslation", "", "", PermissionScopeShop}
	PermissionReadShippingMethodTranslation = &Permission{"read_shippingmethodtranslation", "", "", PermissionScopeShop}
	PermissionUpdateShippingMethodTranslation = &Permission{"update_shippingmethodtranslation", "", "", PermissionScopeShop}
	PermissionDeleteShippingMethodTranslation = &Permission{"delete_shippingmethodtranslation", "", "", PermissionScopeShop}

	PermissionCreateFileInfo = &Permission{"create_fileinfo", "", "", PermissionScopeShop}
	PermissionReadFileInfo = &Permission{"read_fileinfo", "", "", PermissionScopeShop}
	PermissionUpdateFileInfo = &Permission{"update_fileinfo", "", "", PermissionScopeShop}
	PermissionDeleteFileInfo = &Permission{"delete_fileinfo", "", "", PermissionScopeShop}

	PermissionCreatePage = &Permission{"create_page", "", "", PermissionScopeShop}
	PermissionReadPage = &Permission{"read_page", "", "", PermissionScopeShop}
	PermissionUpdatePage = &Permission{"update_page", "", "", PermissionScopeShop}
	PermissionDeletePage = &Permission{"delete_page", "", "", PermissionScopeShop}

	PermissionCreateInvoiceEvent = &Permission{"create_invoiceevent", "", "", PermissionScopeShop}
	PermissionReadInvoiceEvent = &Permission{"read_invoiceevent", "", "", PermissionScopeShop}
	PermissionUpdateInvoiceEvent = &Permission{"update_invoiceevent", "", "", PermissionScopeShop}
	PermissionDeleteInvoiceEvent = &Permission{"delete_invoiceevent", "", "", PermissionScopeShop}

	PermissionCreateInvoice = &Permission{"create_invoice", "", "", PermissionScopeShop}
	PermissionReadInvoice = &Permission{"read_invoice", "", "", PermissionScopeShop}
	PermissionUpdateInvoice = &Permission{"update_invoice", "", "", PermissionScopeShop}
	PermissionDeleteInvoice = &Permission{"delete_invoice", "", "", PermissionScopeShop}

	PermissionCreateWishlist = &Permission{"create_wishlist", "", "", PermissionScopeShop}
	PermissionReadWishlist = &Permission{"read_wishlist", "", "", PermissionScopeShop}
	PermissionUpdateWishlist = &Permission{"update_wishlist", "", "", PermissionScopeShop}
	PermissionDeleteWishlist = &Permission{"delete_wishlist", "", "", PermissionScopeShop}

	PermissionCreatePreference = &Permission{"create_preference", "", "", PermissionScopeSystem}
	PermissionReadPreference = &Permission{"read_preference", "", "", PermissionScopeSystem}
	PermissionUpdatePreference = &Permission{"update_preference", "", "", PermissionScopeSystem}
	PermissionDeletePreference = &Permission{"delete_preference", "", "", PermissionScopeSystem}

	PermissionCreateAssignedProductAttribute = &Permission{"create_assignedproductattribute", "", "", PermissionScopeShop}
	PermissionReadAssignedProductAttribute = &Permission{"read_assignedproductattribute", "", "", PermissionScopeShop}
	PermissionUpdateAssignedProductAttribute = &Permission{"update_assignedproductattribute", "", "", PermissionScopeShop}
	PermissionDeleteAssignedProductAttribute = &Permission{"delete_assignedproductattribute", "", "", PermissionScopeShop}

	PermissionCreatePageType = &Permission{"create_pagetype", "", "", PermissionScopeShop}
	PermissionReadPageType = &Permission{"read_pagetype", "", "", PermissionScopeShop}
	PermissionUpdatePageType = &Permission{"update_pagetype", "", "", PermissionScopeShop}
	PermissionDeletePageType = &Permission{"delete_pagetype", "", "", PermissionScopeShop}

	PermissionCreateDigitalContentURL = &Permission{"create_digitalcontenturl", "", "", PermissionScopeShop}
	PermissionReadDigitalContentURL = &Permission{"read_digitalcontenturl", "", "", PermissionScopeShop}
	PermissionUpdateDigitalContentURL = &Permission{"update_digitalcontenturl", "", "", PermissionScopeShop}
	PermissionDeleteDigitalContentURL = &Permission{"delete_digitalcontenturl", "", "", PermissionScopeShop}

	PermissionCreateAttributeVariant = &Permission{"create_attributevariant", "", "", PermissionScopeShop}
	PermissionReadAttributeVariant = &Permission{"read_attributevariant", "", "", PermissionScopeShop}
	PermissionUpdateAttributeVariant = &Permission{"update_attributevariant", "", "", PermissionScopeShop}
	PermissionDeleteAttributeVariant = &Permission{"delete_attributevariant", "", "", PermissionScopeShop}

	PermissionCreateAssignedPageAttributeValue = &Permission{"create_assignedpageattributevalue", "", "", PermissionScopeShop}
	PermissionReadAssignedPageAttributeValue = &Permission{"read_assignedpageattributevalue", "", "", PermissionScopeShop}
	PermissionUpdateAssignedPageAttributeValue = &Permission{"update_assignedpageattributevalue", "", "", PermissionScopeShop}
	PermissionDeleteAssignedPageAttributeValue = &Permission{"delete_assignedpageattributevalue", "", "", PermissionScopeShop}

	PermissionCreateShippingZoneChannel = &Permission{"create_shippingzonechannel", "", "", PermissionScopeShop}
	PermissionReadShippingZoneChannel = &Permission{"read_shippingzonechannel", "", "", PermissionScopeShop}
	PermissionUpdateShippingZoneChannel = &Permission{"update_shippingzonechannel", "", "", PermissionScopeShop}
	PermissionDeleteShippingZoneChannel = &Permission{"delete_shippingzonechannel", "", "", PermissionScopeShop}

	PermissionCreateGiftcardEvent = &Permission{"create_giftcardevent", "", "", PermissionScopeShop}
	PermissionReadGiftcardEvent = &Permission{"read_giftcardevent", "", "", PermissionScopeShop}
	PermissionUpdateGiftcardEvent = &Permission{"update_giftcardevent", "", "", PermissionScopeShop}
	PermissionDeleteGiftcardEvent = &Permission{"delete_giftcardevent", "", "", PermissionScopeShop}

	PermissionCreateGiftcardCheckout = &Permission{"create_giftcardcheckout", "", "", PermissionScopeShop}
	PermissionReadGiftcardCheckout = &Permission{"read_giftcardcheckout", "", "", PermissionScopeShop}
	PermissionUpdateGiftcardCheckout = &Permission{"update_giftcardcheckout", "", "", PermissionScopeShop}
	PermissionDeleteGiftcardCheckout = &Permission{"delete_giftcardcheckout", "", "", PermissionScopeShop}

	SystemGuestPermissions = Permissions{
		PermissionReadShop,
		PermissionReadPageTranslation,
		PermissionReadCategoryTranslation,
		PermissionReadProductTranslation,
		PermissionReadProductChannelListing,
		PermissionReadCollectionChannelListing,
		PermissionReadVoucherTranslation,
		PermissionReadProductVariantTranslation,
		PermissionReadShopTranslation,
		PermissionReadChannel,
		PermissionReadProduct,
		PermissionReadProductVariantMedia,
		PermissionReadSale,
		PermissionReadVoucher,
		PermissionReadSaleTranslation,
		PermissionReadCollection,
		PermissionReadVoucherChannelListing,
		PermissionReadProductVariant,
		PermissionReadMenu,
		PermissionReadVoucherCategory,
		PermissionReadCollectionTranslation,
		PermissionReadProductVariantChannelListing,
		PermissionReadCategory,
		PermissionReadPage,
	}

	SystemUserPermissions = append(
		SystemGuestPermissions, //

		PermissionReadSaleCategoryRelation,
		PermissionInviteUser,
		PermissionCreateShop,
		PermissionCreatePayment,
		PermissionReadPayment,
		PermissionCreateOrderLine,
		PermissionReadOrderLine,
		PermissionUpdateOrderLine,
		PermissionDeleteOrderLine,
		PermissionReadSaleChannelListing,
		PermissionReadShippingZone,
		PermissionCreateTransaction,
		PermissionReadTransaction,
		PermissionUpdateTransaction,
		PermissionCreateCheckoutLine,
		PermissionReadCheckoutLine,
		PermissionReadShippingMethodChannelListing,
		PermissionReadAttributeTranslation,
		PermissionReadAttributeValueTranslation,
		PermissionReadGiftcard,
		PermissionCreatePayment,
		PermissionReadPayment,
		PermissionReadAttribute,
		PermissionReadShippingMethod,
		PermissionReadShippingMethodPostalCodeRule,
		PermissionCreateCheckout,
		PermissionReadMenuItem,
		PermissionReadProductMedia,
		PermissionReadUser,
		PermissionCreateWishlistItem,
		PermissionReadWishlistItem,
		PermissionUpdateWishlistItem,
		PermissionDeleteWishlistItem,
		PermissionCreateAddress,
		PermissionReadAddress,
		PermissionUpdateAddress,
		PermissionDeleteAddress,
		PermissionCreateOrder,
		PermissionReadOrder,
		PermissionUpdateOrder,
		PermissionDeleteOrder,
		PermissionReadShippingMethodTranslation,
		PermissionCreateWishlist,
		PermissionReadWishlist,
		PermissionUpdateWishlist,
		PermissionDeleteWishlist,
		PermissionCreatePreference,
		PermissionReadPreference,
		PermissionUpdatePreference,
		PermissionDeletePreference)

	ShopStaffPermissions = append(
		SystemUserPermissions,
	)

	ShopScopedAllPermissions = append(ShopStaffPermissions)
}
