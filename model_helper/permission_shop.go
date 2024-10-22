package model_helper

var (
	PermissionCreateWarehouse *Permission
	PermissionReadWarehouse   *Permission
	PermissionUpdateWarehouse *Permission
	PermissionDeleteWarehouse *Permission

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

	PermissionCreateShippingZone *Permission
	PermissionReadShippingZone   *Permission
	PermissionUpdateShippingZone *Permission
	PermissionDeleteShippingZone *Permission

	PermissionCreateTransaction *Permission
	PermissionReadTransaction   *Permission
	PermissionUpdateTransaction *Permission
	PermissionDeleteTransaction *Permission

	// PermissionCreateAssignedPageAttribute *Permission
	// PermissionReadAssignedPageAttribute   *Permission
	// PermissionUpdateAssignedPageAttribute *Permission
	// PermissionDeleteAssignedPageAttribute *Permission

	// PermissionCreateAssignedVariantAttributeValue *Permission
	// PermissionReadAssignedVariantAttributeValue   *Permission
	// PermissionUpdateAssignedVariantAttributeValue *Permission
	// PermissionDeleteAssignedVariantAttributeValue *Permission

	PermissionCreateAttributeValue *Permission
	PermissionReadAttributeValue   *Permission
	PermissionUpdateAttributeValue *Permission // only administrators have this permission
	PermissionDeleteAttributeValue *Permission // only administrators have this permission

	// PermissionCreateAttributePage *Permission
	// PermissionReadAttributePage   *Permission
	// PermissionUpdateAttributePage *Permission
	// PermissionDeleteAttributePage *Permission

	// PermissionCreateAttributeProduct *Permission
	// PermissionReadAttributeProduct   *Permission
	// PermissionUpdateAttributeProduct *Permission
	// PermissionDeleteAttributeProduct *Permission

	PermissionCreateAttributeTranslation *Permission
	PermissionReadAttributeTranslation   *Permission
	PermissionUpdateAttributeTranslation *Permission
	PermissionDeleteAttributeTranslation *Permission

	PermissionCreateAttributeValueTranslation *Permission
	PermissionReadAttributeValueTranslation   *Permission
	PermissionUpdateAttributeValueTranslation *Permission
	PermissionDeleteAttributeValueTranslation *Permission

	PermissionReadAttribute   *Permission
	PermissionCreateAttribute *Permission // system manager and system admin can do this
	PermissionUpdateAttribute *Permission // system manager and system admin can do this
	PermissionDeleteAttribute *Permission // system manager and system admin can do this

	// PermissionCreateAssignedProductAttribute *Permission
	// PermissionReadAssignedProductAttribute   *Permission
	// PermissionUpdateAssignedProductAttribute *Permission
	// PermissionDeleteAssignedProductAttribute *Permission

	PermissionCreateAttributeVariant *Permission
	PermissionReadAttributeVariant   *Permission
	PermissionUpdateAttributeVariant *Permission
	PermissionDeleteAttributeVariant *Permission

	// PermissionCreateAssignedPageAttributeValue *Permission
	// PermissionReadAssignedPageAttributeValue   *Permission
	// PermissionUpdateAssignedPageAttributeValue *Permission
	// PermissionDeleteAssignedPageAttributeValue *Permission

	// PermissionCreateOrderEvent *Permission
	// PermissionDeleteOrderEvent *Permission
	// PermissionUpdateOrderEvent *Permission

	PermissionReadOrderEvent *Permission

	PermissionCreatePreOrderAllocation *Permission
	PermissionReadPreOrderAllocation   *Permission
	PermissionUpdatePreOrderAllocation *Permission
	PermissionDeletePreOrderAllocation *Permission

	PermissionReadCustomerEvent *Permission
	// PermissionCreateCustomerEvent *Permission
	// PermissionUpdateCustomerEvent *Permission
	// PermissionDeleteCustomerEvent *Permission

	PermissionReadCsvExportEvent *Permission
	// PermissionCreateCsvExportEvent *Permission
	// PermissionUpdateCsvExportEvent *Permission
	// PermissionDeleteCsvExportEvent *Permission

	PermissionReadCsvExportFile *Permission
	// PermissionCreateCsvExportFile *Permission
	// PermissionUpdateCsvExportFile *Permission
	// PermissionDeleteCsvExportFile *Permission

	PermissionCreateCustomerNote *Permission
	PermissionReadCustomerNote   *Permission
	PermissionUpdateCustomerNote *Permission
	PermissionDeleteCustomerNote *Permission

	PermissionCreateMenuItemTranslation *Permission
	PermissionReadMenuItemTranslation   *Permission
	PermissionUpdateMenuItemTranslation *Permission
	PermissionDeleteMenuItemTranslation *Permission

	PermissionCreateCheckoutLine *Permission
	PermissionReadCheckoutLine   *Permission
	PermissionUpdateCheckoutLine *Permission
	PermissionDeleteCheckoutLine *Permission

	PermissionCreateShopStaff *Permission
	PermissionReadShopStaff   *Permission
	PermissionUpdateShopStaff *Permission
	// PermissionDeleteShopStaff *Permission

	PermissionCreateProductTranslation *Permission
	PermissionReadProductTranslation   *Permission
	PermissionUpdateProductTranslation *Permission
	PermissionDeleteProductTranslation *Permission

	PermissionCreatePluginConfiguration *Permission // each shop has their own configuration
	PermissionReadPluginConfiguration   *Permission // each shop has their own configuration
	PermissionUpdatePluginConfiguration *Permission // each shop has their own configuration
	PermissionDeletePluginConfiguration *Permission // each shop has their own configuration

	PermissionCreateAudit *Permission // system scoped
	PermissionReadAudit   *Permission // system scoped
	PermissionUpdateAudit *Permission // system scoped
	PermissionDeleteAudit *Permission // system scoped

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

	PermissionCreateClusterDiscovery *Permission // system scoped
	PermissionReadClusterDiscovery   *Permission // system scoped
	PermissionUpdateClusterDiscovery *Permission // system scoped
	PermissionDeleteClusterDiscovery *Permission // system scoped

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

	PermissionCreateCompliance *Permission // system scoped
	PermissionReadCompliance   *Permission // system scoped
	PermissionUpdateCompliance *Permission // system scoped
	PermissionDeleteCompliance *Permission // system scoped

	PermissionCreateStaffNotificationRecipient *Permission
	PermissionReadStaffNotificationRecipient   *Permission
	PermissionUpdateStaffNotificationRecipient *Permission
	PermissionDeleteStaffNotificationRecipient *Permission

	// PermissionCreatePluginKeyValueStore *Permission
	// PermissionReadPluginKeyValueStore   *Permission
	// PermissionUpdatePluginKeyValueStore *Permission
	// PermissionDeletePluginKeyValueStore *Permission

	PermissionCreateChannel *Permission // system admin can do this only
	PermissionReadChannel   *Permission // system admin can do this only
	PermissionUpdateChannel *Permission // system admin can do this only
	PermissionDeleteChannel *Permission // system admin can do this only

	PermissionCreateFulfillmentLine *Permission
	PermissionReadFulfillmentLine   *Permission
	PermissionUpdateFulfillmentLine *Permission
	PermissionDeleteFulfillmentLine *Permission

	PermissionCreateFulfillment *Permission
	PermissionReadFulfillment   *Permission
	PermissionUpdateFulfillment *Permission
	PermissionDeleteFulfillment *Permission

	PermissionCreateProduct *Permission
	PermissionReadProduct   *Permission
	PermissionUpdateProduct *Permission
	PermissionDeleteProduct *Permission

	PermissionCreateTermsOfService *Permission // system scoped
	PermissionReadTermsOfService   *Permission // system scoped
	PermissionUpdateTermsOfService *Permission // system scoped
	PermissionDeleteTermsOfService *Permission // system scoped

	PermissionCreateOrderDiscount *Permission
	PermissionReadOrderDiscount   *Permission
	PermissionUpdateOrderDiscount *Permission
	PermissionDeleteOrderDiscount *Permission

	PermissionCreateProductVariantMedia *Permission
	PermissionReadProductVariantMedia   *Permission
	PermissionUpdateProductVariantMedia *Permission
	PermissionDeleteProductVariantMedia *Permission

	PermissionCreateGiftcard *Permission
	PermissionReadGiftcard   *Permission
	PermissionUpdateGiftcard *Permission
	PermissionDeleteGiftcard *Permission

	PermissionCreatePayment *Permission
	PermissionReadPayment   *Permission
	// PermissionUpdatePayment *Permission
	// PermissionDeletePayment *Permission

	PermissionCreateToken *Permission
	PermissionReadToken   *Permission
	// PermissionUpdateToken *Permission
	PermissionDeleteToken *Permission

	PermissionCreateSale *Permission
	PermissionReadSale   *Permission
	PermissionUpdateSale *Permission
	PermissionDeleteSale *Permission

	PermissionCreateShippingMethod *Permission
	PermissionReadShippingMethod   *Permission
	PermissionUpdateShippingMethod *Permission
	PermissionDeleteShippingMethod *Permission

	// PermissionCreateShippingMethodPostalCodeRule *Permission
	PermissionReadShippingMethodPostalCodeRule *Permission
	// PermissionUpdateShippingMethodPostalCodeRule *Permission
	// PermissionDeleteShippingMethodPostalCodeRule *Permission

	PermissionCreateCheckout *Permission
	PermissionReadCheckout   *Permission
	PermissionUpdateCheckout *Permission
	// PermissionDeleteCheckout *Permission

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

	// PermissionCreateFileInfo *Permission
	// PermissionReadFileInfo   *Permission
	// PermissionUpdateFileInfo *Permission
	// PermissionDeleteFileInfo *Permission

	PermissionCreatePage *Permission
	PermissionReadPage   *Permission
	PermissionUpdatePage *Permission
	PermissionDeletePage *Permission

	// PermissionCreateInvoiceEvent *Permission
	// PermissionReadInvoiceEvent   *Permission
	// PermissionUpdateInvoiceEvent *Permission
	// PermissionDeleteInvoiceEvent *Permission

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

	PermissionCreatePageType *Permission
	PermissionReadPageType   *Permission
	PermissionUpdatePageType *Permission
	PermissionDeletePageType *Permission

	PermissionCreateDigitalContentURL *Permission
	PermissionReadDigitalContentURL   *Permission
	PermissionUpdateDigitalContentURL *Permission
	PermissionDeleteDigitalContentURL *Permission

	PermissionCreateShippingZoneChannel *Permission
	PermissionReadShippingZoneChannel   *Permission
	PermissionUpdateShippingZoneChannel *Permission
	PermissionDeleteShippingZoneChannel *Permission

	// PermissionCreateGiftcardEvent *Permission
	// PermissionReadGiftcardEvent   *Permission
	// PermissionUpdateGiftcardEvent *Permission
	// PermissionDeleteGiftcardEvent *Permission
)

// ShopScopedAllPermissions contains all shop-related permissions
var ShopScopedAllPermissions Permissions
var ShopStaffPermissions Permissions
var SystemUserPermissions Permissions

func initializeShopScopedPermissions() {
	PermissionCreateWarehouse = &Permission{"create_warehouse", "", "", PermissionScopeShop}
	PermissionReadWarehouse = &Permission{"read_warehouse", "", "", PermissionScopeShop}
	PermissionUpdateWarehouse = &Permission{"update_warehouse", "", "", PermissionScopeShop}
	PermissionDeleteWarehouse = &Permission{"delete_warehouse", "", "", PermissionScopeShop}

	// PermissionCreateAssignedPageAttribute = &Permission{"create_assignedpageattribute", "", "", PermissionScopeShop}
	// PermissionReadAssignedPageAttribute = &Permission{"read_assignedpageattribute", "", "", PermissionScopeShop}
	// PermissionUpdateAssignedPageAttribute = &Permission{"update_assignedpageattribute", "", "", PermissionScopeShop}
	// PermissionDeleteAssignedPageAttribute = &Permission{"delete_assignedpageattribute", "", "", PermissionScopeShop}

	PermissionCreateSaleChannelListing = &Permission{"create_salechannellisting", "", "", PermissionScopeShop}
	PermissionReadSaleChannelListing = &Permission{"read_salechannellisting", "", "", PermissionScopeShop}
	PermissionUpdateSaleChannelListing = &Permission{"update_salechannellisting", "", "", PermissionScopeShop}
	PermissionDeleteSaleChannelListing = &Permission{"delete_salechannellisting", "", "", PermissionScopeShop}

	PermissionAddReaction = &Permission{"add_reaction", "authentication.permissions.add_reaction.name", "authentication.permissions.add_reaction.description", PermissionScopeShop}
	PermissionRemoveReaction = &Permission{"remove_reaction", "authentication.permissions.remove_reaction.name", "authentication.permissions.remove_reaction.description", PermissionScopeShop}
	PermissionUploadFile = &Permission{"upload_file", "authentication.permissions.upload_file.name", "authentication.permissions.upload_file.description", PermissionScopeShop}

	PermissionReadShop = &Permission{"read_shop", "", "", PermissionScopeShop}
	PermissionUpdateShop = &Permission{"update_shop", "", "", PermissionScopeShop}
	PermissionDeleteShop = &Permission{"delete_shop", "", "", PermissionScopeShop}

	PermissionCreatePageTranslation = &Permission{"create_pagetranslation", "", "", PermissionScopeShop}
	PermissionReadPageTranslation = &Permission{"read_pagetranslation", "", "", PermissionScopeShop}
	PermissionUpdatePageTranslation = &Permission{"update_pagetranslation", "", "", PermissionScopeShop}
	PermissionDeletePageTranslation = &Permission{"delete_pagetranslation", "", "", PermissionScopeShop}

	PermissionCreateStock = &Permission{"create_stock", "", "", PermissionScopeShop}
	PermissionReadStock = &Permission{"read_stock", "", "", PermissionScopeShop}
	PermissionUpdateStock = &Permission{"update_stock", "", "", PermissionScopeShop}
	PermissionDeleteStock = &Permission{"delete_stock", "", "", PermissionScopeShop}

	// PermissionCreateAssignedVariantAttributeValue = &Permission{"create_assignedvariantattributevalue", "", "", PermissionScopeShop}
	// PermissionReadAssignedVariantAttributeValue = &Permission{"read_assignedvariantattributevalue", "", "", PermissionScopeShop}
	// PermissionUpdateAssignedVariantAttributeValue = &Permission{"update_assignedvariantattributevalue", "", "", PermissionScopeShop}
	// PermissionDeleteAssignedVariantAttributeValue = &Permission{"delete_assignedvariantattributevalue", "", "", PermissionScopeShop}

	PermissionCreateShippingZone = &Permission{"create_shippingzone", "", "", PermissionScopeShop}
	PermissionReadShippingZone = &Permission{"read_shippingzone", "", "", PermissionScopeShop}
	PermissionUpdateShippingZone = &Permission{"update_shippingzone", "", "", PermissionScopeShop}
	PermissionDeleteShippingZone = &Permission{"delete_shippingzone", "", "", PermissionScopeShop}

	PermissionCreateAttributeValue = &Permission{"create_attributevalue", "", "", PermissionScopeShop}
	PermissionReadAttributeValue = &Permission{"read_attributevalue", "", "", PermissionScopeShop}

	// PermissionCreateAttributePage = &Permission{"create_attributepage", "", "", PermissionScopeShop}
	// PermissionReadAttributePage = &Permission{"read_attributepage", "", "", PermissionScopeShop}
	// PermissionUpdateAttributePage = &Permission{"update_attributepage", "", "", PermissionScopeShop}
	// PermissionDeleteAttributePage = &Permission{"delete_attributepage", "", "", PermissionScopeShop}

	// PermissionCreateOrderEvent = &Permission{"create_orderevent", "", "", PermissionScopeShop}
	PermissionReadOrderEvent = &Permission{"read_orderevent", "", "", PermissionScopeShop}
	// PermissionUpdateOrderEvent = &Permission{"update_orderevent", "", "", PermissionScopeShop}
	// PermissionDeleteOrderEvent = &Permission{"delete_orderevent", "", "", PermissionScopeShop}

	PermissionCreatePreOrderAllocation = &Permission{"create_preorderallocation", "", "", PermissionScopeShop}
	PermissionReadPreOrderAllocation = &Permission{"read_preorderallocation", "", "", PermissionScopeShop}
	PermissionUpdatePreOrderAllocation = &Permission{"update_preorderallocation", "", "", PermissionScopeShop}
	PermissionDeletePreOrderAllocation = &Permission{"delete_preorderallocation", "", "", PermissionScopeShop}

	// PermissionCreateCustomerEvent = &Permission{"create_customerevent", "", "", PermissionScopeShop}
	PermissionReadCustomerEvent = &Permission{"read_customerevent", "", "", PermissionScopeShop}
	// PermissionUpdateCustomerEvent = &Permission{"update_customerevent", "", "", PermissionScopeShop}
	// PermissionDeleteCustomerEvent = &Permission{"delete_customerevent", "", "", PermissionScopeShop}

	// PermissionCreateAttributeProduct = &Permission{"create_attributeproduct", "", "", PermissionScopeShop}
	// PermissionReadAttributeProduct = &Permission{"read_attributeproduct", "", "", PermissionScopeShop}
	// PermissionUpdateAttributeProduct = &Permission{"update_attributeproduct", "", "", PermissionScopeShop}
	// PermissionDeleteAttributeProduct = &Permission{"delete_attributeproduct", "", "", PermissionScopeShop}

	// PermissionCreateCsvExportEvent = &Permission{"create_csvexportevent", "", "", PermissionScopeShop}
	PermissionReadCsvExportEvent = &Permission{"read_csvexportevent", "", "", PermissionScopeShop}
	// PermissionUpdateCsvExportEvent = &Permission{"update_csvexportevent", "", "", PermissionScopeShop}
	// PermissionDeleteCsvExportEvent = &Permission{"delete_csvexportevent", "", "", PermissionScopeShop}

	// PermissionCreateCsvExportFile = &Permission{"create_csvexportfile", "", "", PermissionScopeShop}
	PermissionReadCsvExportFile = &Permission{"read_csvexportfile", "", "", PermissionScopeShop}
	// PermissionUpdateCsvExportFile = &Permission{"update_csvexportfile", "", "", PermissionScopeShop}
	// PermissionDeleteCsvExportFile = &Permission{"delete_csvexportfile", "", "", PermissionScopeShop}

	PermissionCreateCustomerNote = &Permission{"create_customernote", "", "", PermissionScopeShop}
	PermissionReadCustomerNote = &Permission{"read_customernote", "", "", PermissionScopeShop}
	PermissionUpdateCustomerNote = &Permission{"update_customernote", "", "", PermissionScopeShop}
	PermissionDeleteCustomerNote = &Permission{"delete_customernote", "", "", PermissionScopeShop}

	PermissionCreateMenuItemTranslation = &Permission{"create_menuitemtranslation", "", "", PermissionScopeShop}
	PermissionReadMenuItemTranslation = &Permission{"read_menuitemtranslation", "", "", PermissionScopeShop}
	PermissionUpdateMenuItemTranslation = &Permission{"update_menuitemtranslation", "", "", PermissionScopeShop}
	PermissionDeleteMenuItemTranslation = &Permission{"delete_menuitemtranslation", "", "", PermissionScopeShop}

	PermissionCreateCheckoutLine = &Permission{"create_checkoutline", "", "", PermissionScopeShop}
	PermissionReadCheckoutLine = &Permission{"read_checkoutline", "", "", PermissionScopeShop}
	PermissionUpdateCheckoutLine = &Permission{"update_checkoutline", "", "", PermissionScopeShop}
	PermissionDeleteCheckoutLine = &Permission{"delete_checkoutline", "", "", PermissionScopeShop}

	PermissionCreateShopStaff = &Permission{"create_shopstaff", "", "", PermissionScopeShop}
	PermissionReadShopStaff = &Permission{"read_shopstaff", "", "", PermissionScopeShop}
	PermissionUpdateShopStaff = &Permission{"update_shopstaff", "", "", PermissionScopeShop}
	// PermissionDeleteShopStaff = &Permission{"delete_shopstaff", "", "", PermissionScopeShop}

	PermissionCreateProductTranslation = &Permission{"create_producttranslation", "", "", PermissionScopeShop}
	PermissionReadProductTranslation = &Permission{"read_producttranslation", "", "", PermissionScopeShop}
	PermissionUpdateProductTranslation = &Permission{"update_producttranslation", "", "", PermissionScopeShop}
	PermissionDeleteProductTranslation = &Permission{"delete_producttranslation", "", "", PermissionScopeShop}

	PermissionCreatePluginConfiguration = &Permission{"create_pluginconfiguration", "", "", PermissionScopeShop}
	PermissionReadPluginConfiguration = &Permission{"read_pluginconfiguration", "", "", PermissionScopeShop}
	PermissionUpdatePluginConfiguration = &Permission{"update_pluginconfiguration", "", "", PermissionScopeShop}
	PermissionDeletePluginConfiguration = &Permission{"delete_pluginconfiguration", "", "", PermissionScopeShop}

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

	PermissionCreateStaffNotificationRecipient = &Permission{"create_staffnotificationrecipient", "", "", PermissionScopeShop}
	PermissionReadStaffNotificationRecipient = &Permission{"read_staffnotificationrecipient", "", "", PermissionScopeShop}
	PermissionUpdateStaffNotificationRecipient = &Permission{"update_staffnotificationrecipient", "", "", PermissionScopeShop}
	PermissionDeleteStaffNotificationRecipient = &Permission{"delete_staffnotificationrecipient", "", "", PermissionScopeShop}

	// PermissionCreatePluginKeyValueStore = &Permission{"create_pluginkeyvaluestore", "", "", PermissionScopeSystem}
	// PermissionReadPluginKeyValueStore = &Permission{"read_pluginkeyvaluestore", "", "", PermissionScopeSystem}
	// PermissionUpdatePluginKeyValueStore = &Permission{"update_pluginkeyvaluestore", "", "", PermissionScopeSystem}
	// PermissionDeletePluginKeyValueStore = &Permission{"delete_pluginkeyvaluestore", "", "", PermissionScopeSystem}

	PermissionCreateFulfillmentLine = &Permission{"create_fulfillmentline", "", "", PermissionScopeShop}
	PermissionReadFulfillmentLine = &Permission{"read_fulfillmentline", "", "", PermissionScopeShop}
	PermissionUpdateFulfillmentLine = &Permission{"update_fulfillmentline", "", "", PermissionScopeShop}
	PermissionDeleteFulfillmentLine = &Permission{"delete_fulfillmentline", "", "", PermissionScopeShop}

	PermissionCreateFulfillment = &Permission{"create_fulfillment", "", "", PermissionScopeShop}
	PermissionReadFulfillment = &Permission{"read_fulfillment", "", "", PermissionScopeShop}
	PermissionUpdateFulfillment = &Permission{"update_fulfillment", "", "", PermissionScopeShop}
	PermissionDeleteFulfillment = &Permission{"delete_fulfillment", "", "", PermissionScopeShop}

	PermissionCreateProduct = &Permission{"create_product", "", "", PermissionScopeShop}
	PermissionReadProduct = &Permission{"read_product", "", "", PermissionScopeShop}
	PermissionUpdateProduct = &Permission{"update_product", "", "", PermissionScopeShop}
	PermissionDeleteProduct = &Permission{"delete_product", "", "", PermissionScopeShop}

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
	// PermissionUpdatePayment = &Permission{"update_payment", "", "", PermissionScopeShop}
	// PermissionDeletePayment = &Permission{"delete_payment", "", "", PermissionScopeShop}

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

	// PermissionCreateShippingMethodPostalCodeRule = &Permission{"create_shippingmethodpostalcoderule", "", "", PermissionScopeShop}
	PermissionReadShippingMethodPostalCodeRule = &Permission{"read_shippingmethodpostalcoderule", "", "", PermissionScopeShop}
	// PermissionUpdateShippingMethodPostalCodeRule = &Permission{"update_shippingmethodpostalcoderule", "", "", PermissionScopeShop}
	// PermissionDeleteShippingMethodPostalCodeRule = &Permission{"delete_shippingmethodpostalcoderule", "", "", PermissionScopeShop}

	PermissionCreateCheckout = &Permission{"create_checkout", "", "", PermissionScopeShop}
	PermissionReadCheckout = &Permission{"read_checkout", "", "", PermissionScopeShop}
	PermissionUpdateCheckout = &Permission{"update_checkout", "", "", PermissionScopeShop}
	// PermissionDeleteCheckout = &Permission{"delete_checkout", "", "", PermissionScopeShop}

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

	PermissionCreateOrder = &Permission{"create_order", "", "", PermissionScopeShop}
	PermissionReadOrder = &Permission{"read_order", "", "", PermissionScopeShop}
	PermissionUpdateOrder = &Permission{"update_order", "", "", PermissionScopeShop}
	PermissionDeleteOrder = &Permission{"delete_order", "", "", PermissionScopeShop}

	PermissionCreateShippingMethodTranslation = &Permission{"create_shippingmethodtranslation", "", "", PermissionScopeShop}
	PermissionReadShippingMethodTranslation = &Permission{"read_shippingmethodtranslation", "", "", PermissionScopeShop}
	PermissionUpdateShippingMethodTranslation = &Permission{"update_shippingmethodtranslation", "", "", PermissionScopeShop}
	PermissionDeleteShippingMethodTranslation = &Permission{"delete_shippingmethodtranslation", "", "", PermissionScopeShop}

	// PermissionCreateFileInfo = &Permission{"create_fileinfo", "", "", PermissionScopeShop}
	// PermissionReadFileInfo = &Permission{"read_fileinfo", "", "", PermissionScopeShop}
	// PermissionUpdateFileInfo = &Permission{"update_fileinfo", "", "", PermissionScopeShop}
	// PermissionDeleteFileInfo = &Permission{"delete_fileinfo", "", "", PermissionScopeShop}

	PermissionCreatePage = &Permission{"create_page", "", "", PermissionScopeShop}
	PermissionReadPage = &Permission{"read_page", "", "", PermissionScopeShop}
	PermissionUpdatePage = &Permission{"update_page", "", "", PermissionScopeShop}
	PermissionDeletePage = &Permission{"delete_page", "", "", PermissionScopeShop}

	// PermissionCreateInvoiceEvent = &Permission{"create_invoiceevent", "", "", PermissionScopeShop}
	// PermissionReadInvoiceEvent = &Permission{"read_invoiceevent", "", "", PermissionScopeShop}
	// PermissionUpdateInvoiceEvent = &Permission{"update_invoiceevent", "", "", PermissionScopeShop}
	// PermissionDeleteInvoiceEvent = &Permission{"delete_invoiceevent", "", "", PermissionScopeShop}

	PermissionCreateInvoice = &Permission{"create_invoice", "", "", PermissionScopeShop}
	PermissionReadInvoice = &Permission{"read_invoice", "", "", PermissionScopeShop}
	PermissionUpdateInvoice = &Permission{"update_invoice", "", "", PermissionScopeShop}
	PermissionDeleteInvoice = &Permission{"delete_invoice", "", "", PermissionScopeShop}

	PermissionCreateWishlist = &Permission{"create_wishlist", "", "", PermissionScopeShop}
	PermissionReadWishlist = &Permission{"read_wishlist", "", "", PermissionScopeShop}
	PermissionUpdateWishlist = &Permission{"update_wishlist", "", "", PermissionScopeShop}
	PermissionDeleteWishlist = &Permission{"delete_wishlist", "", "", PermissionScopeShop}

	// PermissionCreateAssignedProductAttribute = &Permission{"create_assignedproductattribute", "", "", PermissionScopeShop}
	// PermissionReadAssignedProductAttribute = &Permission{"read_assignedproductattribute", "", "", PermissionScopeShop}
	// PermissionUpdateAssignedProductAttribute = &Permission{"update_assignedproductattribute", "", "", PermissionScopeShop}
	// PermissionDeleteAssignedProductAttribute = &Permission{"delete_assignedproductattribute", "", "", PermissionScopeShop}

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

	// PermissionCreateAssignedPageAttributeValue = &Permission{"create_assignedpageattributevalue", "", "", PermissionScopeShop}
	// PermissionReadAssignedPageAttributeValue = &Permission{"read_assignedpageattributevalue", "", "", PermissionScopeShop}
	// PermissionUpdateAssignedPageAttributeValue = &Permission{"update_assignedpageattributevalue", "", "", PermissionScopeShop}
	// PermissionDeleteAssignedPageAttributeValue = &Permission{"delete_assignedpageattributevalue", "", "", PermissionScopeShop}

	PermissionCreateShippingZoneChannel = &Permission{"create_shippingzonechannel", "", "", PermissionScopeShop}
	PermissionReadShippingZoneChannel = &Permission{"read_shippingzonechannel", "", "", PermissionScopeShop}
	PermissionUpdateShippingZoneChannel = &Permission{"update_shippingzonechannel", "", "", PermissionScopeShop}
	PermissionDeleteShippingZoneChannel = &Permission{"delete_shippingzonechannel", "", "", PermissionScopeShop}

	// PermissionCreateGiftcardEvent = &Permission{"create_giftcardevent", "", "", PermissionScopeShop}
	// PermissionReadGiftcardEvent = &Permission{"read_giftcardevent", "", "", PermissionScopeShop}
	// PermissionUpdateGiftcardEvent = &Permission{"update_giftcardevent", "", "", PermissionScopeShop}
	// PermissionDeleteGiftcardEvent = &Permission{"delete_giftcardevent", "", "", PermissionScopeShop}
}

func initShopPermissionGroups() {
	SystemUserPermissions = Permissions{
		PermissionReadShop, PermissionReadShopTranslation,
		PermissionReadPage, PermissionReadPageTranslation,
		PermissionReadCategory, PermissionReadCategoryTranslation,
		PermissionReadProduct, PermissionReadProductTranslation,
		PermissionReadCollection, PermissionReadCollectionTranslation,
		PermissionReadVoucher, PermissionReadVoucherTranslation,
		PermissionReadSale, PermissionReadSaleTranslation,
		PermissionReadProductVariant, PermissionReadProductVariantTranslation,
		PermissionReadMenuItem, PermissionReadMenuItemTranslation,
		PermissionReadVoucherChannelListing, PermissionReadSaleChannelListing,
		PermissionReadAttribute, PermissionReadAttributeTranslation,
		PermissionReadAttributeValue, PermissionReadAttributeValueTranslation,
		PermissionReadShippingMethod, PermissionReadShippingMethodTranslation,
		PermissionCreatePayment, PermissionReadPayment,
		PermissionCreateOrderLine, PermissionReadOrderLine, PermissionUpdateOrderLine, PermissionDeleteOrderLine,
		PermissionCreateTransaction, PermissionReadTransaction, PermissionUpdateTransaction,
		PermissionCreateCheckoutLine, PermissionReadCheckoutLine,
		PermissionCreateAddress, PermissionReadAddress, PermissionUpdateAddress, PermissionDeleteAddress,
		PermissionCreateOrder, PermissionReadOrder, PermissionUpdateOrder, PermissionDeleteOrder,
		PermissionCreateWishlistItem, PermissionReadWishlistItem, PermissionUpdateWishlistItem, PermissionDeleteWishlistItem,
		PermissionCreateWishlist, PermissionReadWishlist, PermissionUpdateWishlist, PermissionDeleteWishlist,
		PermissionCreatePreference, PermissionReadPreference, PermissionUpdatePreference, PermissionDeletePreference,
		PermissionCreateCheckout, PermissionReadCheckout,
		PermissionCreateToken, PermissionReadToken, PermissionDeleteToken,

		PermissionReadMenu,
		PermissionReadProductChannelListing,
		PermissionReadCollectionChannelListing,
		PermissionReadChannel,
		PermissionReadProductVariantMedia,
		PermissionReadProductVariantChannelListing,
		PermissionReadShippingZone,
		PermissionReadShippingMethodChannelListing,
		PermissionReadProductMedia,
		PermissionInviteUser,
		PermissionCreateShop,
		PermissionReadGiftcard,
		PermissionReadShippingMethodPostalCodeRule,
		PermissionReadInvoice,
	}

	ShopStaffPermissions = append(
		SystemUserPermissions,

		PermissionReadWarehouse, PermissionUpdateWarehouse,
		PermissionCreateSaleChannelListing, PermissionUpdateSaleChannelListing,
		PermissionCreatePageTranslation, PermissionUpdatePageTranslation, PermissionDeletePageTranslation,
		PermissionCreateShippingZone, PermissionUpdateShippingZone, PermissionDeleteShippingZone,
		PermissionReadAttributeValue,
		PermissionCreatePreOrderAllocation, PermissionReadPreOrderAllocation, PermissionUpdatePreOrderAllocation, PermissionDeletePreOrderAllocation,
		PermissionCreateCustomerNote, PermissionReadCustomerNote, PermissionUpdateCustomerNote, PermissionDeleteCustomerNote,
		PermissionCreateMenuItemTranslation, PermissionUpdateMenuItemTranslation, PermissionDeleteMenuItemTranslation,
		PermissionUpdateCheckoutLine, PermissionDeleteCheckoutLine,
		PermissionCreateProductTranslation, PermissionUpdateProductTranslation, PermissionDeleteProductTranslation,
		PermissionCreatePluginConfiguration, PermissionReadPluginConfiguration, PermissionUpdatePluginConfiguration, PermissionDeletePluginConfiguration,
		PermissionCreateProductChannelListing, PermissionUpdateProductChannelListing, PermissionDeleteProductChannelListing,
		PermissionCreateCollectionChannelListing, PermissionUpdateCollectionChannelListing, PermissionDeleteCollectionChannelListing,
		PermissionCreateVoucherTranslation, PermissionUpdateVoucherTranslation, PermissionDeleteVoucherTranslation,
		PermissionCreateProductVariantTranslation, PermissionUpdateProductVariantTranslation, PermissionDeleteProductVariantTranslation,
		PermissionCreateShopTranslation, PermissionUpdateShopTranslation, PermissionDeleteShopTranslation,
		PermissionCreateShippingMethodChannelListing, PermissionUpdateShippingMethodChannelListing, PermissionDeleteShippingMethodChannelListing,
		PermissionCreateStaffNotificationRecipient, PermissionReadStaffNotificationRecipient, PermissionUpdateStaffNotificationRecipient, PermissionDeleteStaffNotificationRecipient,
		PermissionCreateFulfillmentLine, PermissionReadFulfillmentLine, PermissionUpdateFulfillmentLine,
		PermissionCreateFulfillment, PermissionReadFulfillment, PermissionUpdateFulfillment,
		PermissionCreateProduct, PermissionReadProduct, PermissionUpdateProduct, PermissionDeleteProduct,
		PermissionCreateOrderDiscount, PermissionReadOrderDiscount, PermissionUpdateOrderDiscount, PermissionDeleteOrderDiscount,
		PermissionCreateProductVariantMedia, PermissionReadProductVariantMedia, PermissionUpdateProductVariantMedia, PermissionDeleteProductVariantMedia,
		PermissionCreateAttributeTranslation, PermissionUpdateAttributeTranslation, PermissionDeleteAttributeTranslation,
		PermissionCreateAttributeValueTranslation, PermissionUpdateAttributeValueTranslation, PermissionDeleteAttributeValueTranslation,
		PermissionCreateGiftcard, PermissionUpdateGiftcard, PermissionDeleteGiftcard,
		PermissionCreateSale, PermissionUpdateSale,
		PermissionCreateShippingMethod, PermissionUpdateShippingMethod, PermissionDeleteShippingMethod,
		PermissionCreateCheckout, PermissionUpdateCheckout,
		PermissionCreateAllocation, PermissionReadAllocation, PermissionUpdateAllocation, PermissionDeleteAllocation,
		PermissionCreateVoucher, PermissionUpdateVoucher,
		PermissionCreateMenuItem, PermissionUpdateMenuItem, PermissionDeleteMenuItem,
		PermissionCreateProductMedia, PermissionUpdateProductMedia, PermissionDeleteProductMedia,
		PermissionCreateProductType, PermissionReadProductType, PermissionUpdateProductType, PermissionDeleteProductType,
		PermissionCreateSaleTranslation, PermissionUpdateSaleTranslation, PermissionDeleteSaleTranslation,
		PermissionCreateShippingMethodExcludedProduct, PermissionReadShippingMethodExcludedProduct, PermissionUpdateShippingMethodExcludedProduct, PermissionDeleteShippingMethodExcludedProduct,
		PermissionCreateCollection, PermissionUpdateCollection, PermissionDeleteCollection,
		PermissionCreateDigitalContent, PermissionReadDigitalContent, PermissionUpdateDigitalContent, PermissionDeleteDigitalContent,
		PermissionCreateVoucherChannelListing, PermissionUpdateVoucherChannelListing, PermissionDeleteVoucherChannelListing,
		PermissionCreateProductVariant, PermissionUpdateProductVariant, PermissionDeleteProductVariant,
		PermissionCreateMenu, PermissionUpdateMenu, PermissionDeleteMenu,
		PermissionCreateCollectionTranslation, PermissionUpdateCollectionTranslation, PermissionDeleteCollectionTranslation,
		PermissionCreateProductVariantChannelListing, PermissionUpdateProductVariantChannelListing, PermissionDeleteProductVariantChannelListing,
		PermissionCreateShippingMethodTranslation, PermissionUpdateShippingMethodTranslation, PermissionDeleteShippingMethodTranslation,
		PermissionCreatePage, PermissionUpdatePage, PermissionDeletePage,
		PermissionCreateInvoice, PermissionUpdateInvoice, PermissionDeleteInvoice,
		PermissionCreatePageType, PermissionReadPageType, PermissionUpdatePageType, PermissionDeletePageType,
		PermissionCreateDigitalContentURL, PermissionReadDigitalContentURL, PermissionUpdateDigitalContentURL, PermissionDeleteDigitalContentURL,
		PermissionCreateAttributeVariant, PermissionReadAttributeVariant, PermissionUpdateAttributeVariant, PermissionDeleteAttributeVariant,
		PermissionCreateShippingZoneChannel, PermissionReadShippingZoneChannel, PermissionUpdateShippingZoneChannel, PermissionDeleteShippingZoneChannel,

		PermissionReadStock,
		PermissionReadOrderEvent,
		PermissionReadCustomerEvent,
		PermissionReadCsvExportEvent,
		PermissionReadCsvExportFile,
		PermissionReadShopStaff,
		PermissionReadUser,
	)

	ShopScopedAllPermissions = append(
		ShopStaffPermissions,

		PermissionCreateWarehouse, PermissionDeleteWarehouse,
		PermissionDeleteSaleChannelListing,
		PermissionUpdateShop, PermissionDeleteShop,
		PermissionUpdateStock, PermissionCreateStock, PermissionDeleteStock,
		PermissionCreateShopStaff, PermissionUpdateShopStaff,

		PermissionDeleteTransaction,
		PermissionDeleteFulfillmentLine,
		PermissionDeleteFulfillment,
		PermissionDeleteSale,
		PermissionDeleteVoucher,
	)
}
