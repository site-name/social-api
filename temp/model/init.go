package model

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/site-name/decimal"
)

func init() {
	initConsts()
	initPermissions()
	initRoles()
}

// ModelFieldKind is similar to reflect.Kind, but simplified.
// There are 2 types of kinds: normals and pointers.
// normal types are odd numbers (1, 3, 5, 7, 9, ...),
// their according pointer versions are doubled of them (2, 6, 10, 14, ...)
type ModelFieldKind uint8

// NOTE: Never change position of those constants below
const (
	Decimal ModelFieldKind = (iota * 2) + 1
	Time
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	String
	Map
	Slice
	Struct
)

// every values here is x2 of according types above
const (
	DecimalPtr ModelFieldKind = ((iota * 2) + 1) * 2
	TimePtr
	BoolPtr
	IntPtr
	Int8Ptr
	Int16Ptr
	Int32Ptr
	Int64Ptr
	UintPtr
	Uint8Ptr
	Uint16Ptr
	Uint32Ptr
	Uint64Ptr
	Float32Ptr
	Float64Ptr
	StringPtr
	MapPtr
	SlicePtr
	StructPtr
)

// NormalModelKindsToAccordingPointerKindsMap holds:
// normal kind <-> according pointer kinds
// E.g
//
//	Int8: Int8Ptr
//	Bool: BoolPtr
var NormalModelKindsToAccordingPointerKindsMap = map[ModelFieldKind]ModelFieldKind{
	Decimal: DecimalPtr,
	Time:    TimePtr,
	Bool:    BoolPtr,
	Int:     IntPtr,
	Int8:    Int8Ptr,
	Int16:   Int16Ptr,
	Int32:   Int32Ptr,
	Int64:   Int64Ptr,
	Uint:    UintPtr,
	Uint8:   Uint8Ptr,
	Uint16:  Uint16Ptr,
	Uint32:  Uint32Ptr,
	Uint64:  Uint64Ptr,
	Float32: Float32Ptr,
	Float64: Float64Ptr,
	String:  StringPtr,
	Map:     MapPtr,
	Slice:   SlicePtr,
	Struct:  StructPtr,
}

// SystemModels holds all database struct models of the wholte system.
// When a new struct model is defined, it must be added into this.
var SystemModels = []Modeler{
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

func modelFieldKindInit() func(key string) (ModelFieldKind, bool) {
	// modelFieldsTypeMap contains models' field types.
	// E.g
	//
	//	"Products.Id": reflect.String
	//	"Users.Metadata": reflect.Map
	var modelFieldsTypeMap = map[string]ModelFieldKind{}

	for _, model := range SystemModels {

		typeOf := reflect.TypeOf(model).Elem()

		for i := 0; i < typeOf.NumField(); i++ {
			var (
				fieldAtIdx = typeOf.Field(i)
				gormTag    = fieldAtIdx.Tag.Get("gorm")
				jsonTag    = fieldAtIdx.Tag.Get("json")
				fieldType  = fieldAtIdx.Type
			)

			if !fieldAtIdx.IsExported() {
				continue
			}

			if fieldType.Kind() == reflect.Pointer {
				fieldType = fieldType.Elem()
			}

			if fieldType.Kind() == reflect.Struct {
				if gormTag == "" && jsonTag == "" { // embed struct
					for j := 0; j < fieldType.NumField(); j++ {
						subField := fieldType.Field(j)
						kind := inspectField(model.TableName(), subField)
						modelFieldsTypeMap[model.TableName()+"."+subField.Name] = kind
					}
					continue
				}

				kind := inspectField(model.TableName(), fieldAtIdx)
				modelFieldsTypeMap[model.TableName()+"."+fieldAtIdx.Name] = kind
				continue
			}

			kind := inspectField(model.TableName(), fieldAtIdx)
			modelFieldsTypeMap[model.TableName()+"."+fieldAtIdx.Name] = kind
		}
	}

	return func(key string) (ModelFieldKind, bool) {
		kind, ok := modelFieldsTypeMap[key]
		return kind, ok
	}
}

// E.g
//
//	GetModelFieldKind("Users.Id") // => reflect.String, true
var GetModelFieldKind = modelFieldKindInit()

func inspectField(modelName string, aField reflect.StructField) ModelFieldKind {
	var (
		gormTag              = aField.Tag.Get("gorm")
		splitGormTag         = strings.Split(gormTag, ";")
		reflectFieldType     = aField.Type
		realReflectFieldKind = aField.Type.Kind()
	)

	if reflectFieldType.Kind() == reflect.Pointer {
		reflectFieldType = reflectFieldType.Elem()
	}

	for _, gormAttr := range splitGormTag {
		if strings.HasPrefix(gormAttr, "column:") {
			columnName := gormAttr[len("column:"):]
			if columnName != aField.Name {
				panic(fmt.Errorf("model: %s, column: %s != field: %s", modelName, columnName, aField.Name))
			}
		}
	}

	var kind ModelFieldKind

	switch {
	case reflectFieldType == reflect.TypeOf(time.Time{}):
		kind = Time
	case reflectFieldType == reflect.TypeOf(decimal.Decimal{}):
		kind = Decimal
	case reflectFieldType.Kind() == reflect.Bool:
		kind = Bool
	case reflectFieldType.Kind() == reflect.Int:
		kind = Int
	case reflectFieldType.Kind() == reflect.Int8:
		kind = Int8
	case reflectFieldType.Kind() == reflect.Int16:
		kind = Int16
	case reflectFieldType.Kind() == reflect.Int32:
		kind = Int32
	case reflectFieldType.Kind() == reflect.Int64:
		kind = Int64
	case reflectFieldType.Kind() == reflect.Uint:
		kind = Uint
	case reflectFieldType.Kind() == reflect.Uint8:
		kind = Uint8
	case reflectFieldType.Kind() == reflect.Uint16:
		kind = Uint16
	case reflectFieldType.Kind() == reflect.Uint32:
		kind = Uint32
	case reflectFieldType.Kind() == reflect.Uint64:
		kind = Uint64
	case reflectFieldType.Kind() == reflect.Float32:
		kind = Float32
	case reflectFieldType.Kind() == reflect.Float64:
		kind = Float64
	case reflectFieldType.Kind() == reflect.String:
		kind = String
	case reflectFieldType.Kind() == reflect.Map:
		kind = Map
	case reflectFieldType.Kind() == reflect.Slice:
		kind = Slice
	case reflectFieldType.Kind() == reflect.Struct:
		kind = Struct
	default:
		panic(errors.Errorf("model: %s, field: %s got unexpected model field type: %s", modelName, aField.Name, reflectFieldType.Kind().String()))
	}

	if realReflectFieldKind == reflect.Pointer {
		kind = NormalModelKindsToAccordingPointerKindsMap[kind]
	}

	return kind

	// panic(fmt.Errorf("model: %s, field: %s gorm column attribute not found", modelName, aField.Name))
}
