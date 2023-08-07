package model

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/site-name/decimal"
)

func init() {
	initConsts()
	initPermissions()
	initRoles()
}

const (
	Decimal reflect.Kind = iota + reflect.UnsafePointer + 1
	Time
)

// SystemModels holds all database struct models of the wholte system.
// When a new struct model is defined, it must be added into this.
var SystemModels = [...]TableModel{
	&Attribute{},                    // attribute
	&AttributeValue{},               //
	&AttributeTranslation{},         //
	&AttributeValueTranslation{},    //
	&AttributeVariant{},             //
	&AttributePage{},                //
	&AttributeProduct{},             //
	&Audit{},                        // audit
	&Channel{},                      // channel
	&CheckoutLine{},                 // checkout
	&Checkout{},                     //
	&ClusterDiscovery{},             // cluster
	&Compliance{},                   // compliance
	&ExportEvent{},                  // csv
	&ExportFile{},                   //
	&OrderDiscount{},                // discount
	&Sale{},                         //
	&VoucherChannelListing{},        //
	&SaleChannelListing{},           //
	&Voucher{},                      //
	&VoucherCustomer{},              //
	&OpenExchangeRate{},             // 3rd party
	&FileInfo{},                     // file info
	&UploadSession{},                //
	&GiftCardEvent{},                // giftcard
	&GiftCard{},                     //
	&InvoiceEvent{},                 // invoice
	&Invoice{},                      //
	&Job{},                          // job
	&MenuItemTranslation{},          // menu
	&MenuItem{},                     //
	&Menu{},                         //
	&Order{},                        // order
	&OrderLine{},                    //
	&FulfillmentLine{},              //
	&Fulfillment{},                  //
	&OrderEvent{},                   //
	&PageType{},                     // page
	&PageTranslation{},              //
	&Page{},                         //
	&Payment{},                      // payment
	&PaymentTransaction{},           //
	&PluginConfiguration{},          // plugin
	&PluginKeyValue{},               //
	&Preference{},                   // preference
	&ShippingMethod{},               // shipping
	&ShippingMethodChannelListing{}, //
	&ShippingMethodPostalCodeRule{}, //
	&ShippingMethodTranslation{},    //
	&ShippingZone{},                 //
	&System{},                       // system
	&TermsOfService{},               //
	&Token{},                        //
	&Vat{},                          // vat
	&Allocation{},                   // warehouse
	&Stock{},                        //
	&WareHouse{},                    //
	&PreorderAllocation{},           //
	&Wishlist{},                     // wishlist
	&WishlistItem{},                 //
	&ShopTranslation{},              // shop
	&ShopStaff{},                    //
	&Address{},                      // account
	&Status{},                       //
	&UserAccessToken{},              //
	&CustomerEvent{},                //
	&CustomerNote{},                 //
	&AppToken{},                     //
	&Session{},                      //
	&Role{},                         //
	&User{},                         //
	&Category{},                     // product
	&CategoryTranslation{},          //
	&Product{},                      //
	&ProductVariant{},               //
	&Collection{},                   //
	&CollectionChannelListing{},     //
	&CollectionProduct{},            //
	&CollectionTranslation{},        //
	&DigitalContent{},               //
	&DigitalContentUrl{},            //
	&ProductMedia{},                 //
	&ProductTranslation{},           //
	&ProductType{},                  //
	&ProductVariantChannelListing{}, //
	&ProductVariantTranslation{},    //
	&ProductChannelListing{},        //
}

func modelFieldKindInit() func(key string) (reflect.Kind, bool) {
	// modelFieldsTypeMap contains models' field types.
	// E.g
	//
	//	"Products.Id": reflect.String
	//	"Users.Metadata": reflect.Map
	var modelFieldsTypeMap = map[string]reflect.Kind{}

	for _, model := range SystemModels {
		typeOf := reflect.TypeOf(model)
		if typeOf.Kind() == reflect.Pointer {
			typeOf = typeOf.Elem()
		}

		for i := 0; i < typeOf.NumField(); i++ {
			fieldAtIdx := typeOf.Field(i)
			gormTag := fieldAtIdx.Tag.Get("gorm")
			jsonTag := fieldAtIdx.Tag.Get("json")
			fieldType := fieldAtIdx.Type

			if fieldType.Kind() == reflect.Pointer {
				fieldType = fieldType.Elem()
			}

			if fieldAtIdx.IsExported() && gormTag == "" && jsonTag == "" && fieldType.Kind() == reflect.Struct { // embed struct
				for j := 0; j < fieldType.NumField(); j++ {
					subField := fieldType.Field(j)
					kind, err := inspectField(model.TableName(), subField)
					if err != nil {
						panic(err)
					}

					modelFieldsTypeMap[model.TableName()+"."+subField.Name] = kind
				}
				continue
			}

			if fieldAtIdx.IsExported() && gormTag != "-" && jsonTag != "-" {
				kind, err := inspectField(model.TableName(), fieldAtIdx)
				if err != nil {
					panic(err)
				}

				modelFieldsTypeMap[model.TableName()+"."+fieldAtIdx.Name] = kind
			}
		}
	}

	return func(key string) (reflect.Kind, bool) {
		kind, ok := modelFieldsTypeMap[key]
		return kind, ok
	}
}

// E.g
//
//	GetModelFieldKind("Users.Id") // => reflect.String, true
var GetModelFieldKind = modelFieldKindInit()

func inspectField(modelName string, aField reflect.StructField) (reflect.Kind, error) {
	gormTag := aField.Tag.Get("gorm")
	splitGormTag := strings.Split(gormTag, ";")
	fieldType := aField.Type
	if fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	columnAttr, found := lo.Find(splitGormTag, func(s string) bool { return strings.HasPrefix(s, "column:") })
	if found {
		columnName := columnAttr[len("column:"):]
		if columnName == aField.Name {
			var kind reflect.Kind
			switch fieldType {
			case reflect.TypeOf(time.Time{}):
				kind = Time
			case reflect.TypeOf(decimal.Decimal{}):
				kind = Decimal
			default:
				kind = fieldType.Kind()
			}

			return kind, nil
		}

		return reflect.Invalid, fmt.Errorf("column: %s != field: %s", columnName, aField.Name)
	}

	return reflect.Invalid, fmt.Errorf("model: %s, field: %s gorm column attribute not found", modelName, aField.Name)
}
