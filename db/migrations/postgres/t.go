package main

// import (
// 	"errors"
// 	"fmt"
// 	"log"
// 	"os"
// 	"sort"
// 	"strings"
// )

// var tables = []string{
// 	"Categories",
// 	"CategoryTranslations",
// 	"ProductChannelListings",
// 	"CollectionChannelListings",
// 	"ProductCollections",
// 	"Collections",
// 	"CollectionTranslations",
// 	"DigitalContents",
// 	"DigitalContentURLs",
// 	"ProductMedias",
// 	"Products",
// 	"ProductTranslations",
// 	"ProductTypes",
// 	"ProductVariantChannelListings",
// 	"VariantMedias",
// 	"ProductVariants",
// 	"ProductVariantTranslations",
// 	"WishlistItems",
// 	"WishlistItemProductVariants",
// 	"Wishlists",
// 	"Stocks",
// 	"Warehouses",
// 	"WarehouseShippingZones",
// 	"Allocations",
// 	"PreorderAllocations",
// 	"CheckoutLines",
// 	"Checkouts",
// 	"Orderlines",
// 	"Orders",
// 	"FulfillmentLines",
// 	"Fulfillments",
// 	"OrderEvents",
// 	"Addresses",
// 	"Users",
// 	"CustomerEvents",
// 	"StaffNotificationRecipients",
// 	"CustomerNotes",
// 	"Tokens",
// 	"UserAddresses",
// 	"TermsOfServices",
// 	"Status",
// 	"Channels",
// 	"GiftCards",
// 	"GiftcardEvents",
// 	"OrderGiftCards",
// 	"GiftcardCheckouts",
// 	"Payments",
// 	"Transactions",
// 	"PluginKeyValueStore",
// 	"Preferences",
// 	"Roles",
// 	"ExportEvents",
// 	"ExportFiles",
// 	"BaseAssignedAttributes",
// 	"Attributes",
// 	"AttributeTranslations",
// 	"AttributeValues",
// 	"AttributeValueTranslations",
// 	"AssignedPageAttributeValues",
// 	"AssignedPageAttributes",
// 	"AttributePages",
// 	"AssignedVariantAttributeValues",
// 	"AssignedVariantAttributes",
// 	"AttributeVariants",
// 	"AssignedProductAttributeValues",
// 	"AssignedProductAttributes",
// 	"AttributeProducts",
// 	"Vouchers",
// 	"VoucherCategories",
// 	"VoucherCollections",
// 	"VoucherProducts",
// 	"VoucherChannelListings",
// 	"VoucherCustomers",
// 	"SaleChannelListings",
// 	"Sales",
// 	"SaleTranslations",
// 	"VoucherTranslations",
// 	"SaleCategories",
// 	"SaleProducts",
// 	"SaleCollections",
// 	"OrderDiscounts",
// 	"VoucherProductVariants",
// 	"SaleProductVariants",
// 	"Shops",
// 	"ShopTranslations",
// 	"ShopStaffs",
// 	"Menus",
// 	"MenuItems",
// 	"MenuItemTranslations",
// 	"ShippingMethods",
// 	"ShippingZones",
// 	"ShippingZoneChannels",
// 	"ShippingMethodTranslations",
// 	"ShippingMethodPostalCodeRules",
// 	"ShippingMethodChannelListings",
// 	"ShippingMethodExcludedProducts",
// 	"Jobs",
// 	"FileInfos",
// 	"UploadSessions",
// 	"Pages",
// 	"PageTranslations",
// 	"PageTypes",
// 	"InvoiceEvents",
// 	"Invoices",
// 	"OpenExchangeRates",
// 	"PluginConfigurations",
// }

// func main() {

// 	sort.Strings(tables)

// 	for index, name := range tables {
// 		var (
// 			lowerName = strings.ToLower(name)
// 			prefix    = fmt.Sprintf("%06d", index+1)
// 			upName    = fmt.Sprintf("%s_create_%s.up.sql", prefix, lowerName)
// 			downName  = fmt.Sprintf("%s_create_%s.down.sql", prefix, lowerName)
// 		)

// 		// check up file exist
// 		if _, err := os.Stat(upName); err != nil && errors.Is(err, os.ErrNotExist) {
// 			upFile, err := os.Create(upName)
// 			if err != nil {
// 				log.Fatalln(err)
// 			}
// 			defer upFile.Close()

// 			upFile.WriteString("CREATE TABLE IF NOT EXISTS " + lowerName + " ();")
// 		}

// 		if _, err := os.Stat(downName); err != nil && errors.Is(err, os.ErrNotExist) {
// 			downFile, err := os.Create(downName)
// 			if err != nil {
// 				log.Fatalln(err)
// 			}
// 			defer downFile.Close()

// 			downFile.WriteString("DROP TABLE IF EXISTS " + lowerName + ";")
// 		}

// 	}
// }
