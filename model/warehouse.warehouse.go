package model

import (
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
)

// max lengths for some warehouse's fields
const (
	WAREHOUSE_NAME_MAX_LENGTH                     = 250
	WAREHOUSE_SLUG_MAX_LENGTH                     = 255
	WAREHOUSE_CLICK_AND_COLLECT_OPTION_MAX_LENGTH = 30
)

type WarehouseClickAndCollectOption string

// default values for warehouse's click_and_collect_option field
const (
	DISABLED       WarehouseClickAndCollectOption = "disabled"
	LOCAL_STOCK    WarehouseClickAndCollectOption = "local"
	ALL_WAREHOUSES WarehouseClickAndCollectOption = "all"
)

func (w WarehouseClickAndCollectOption) IsValid() bool {
	return ValidWarehouseClickAndCollectOptionMap[w] != ""
}

var ValidWarehouseClickAndCollectOptionMap = map[WarehouseClickAndCollectOption]string{
	DISABLED:       "Disabled",
	LOCAL_STOCK:    "Local stock only",
	ALL_WAREHOUSES: "Al warehouses",
}

type WareHouse struct {
	Id                    string                         `json:"id"`
	Name                  string                         `json:"name"`                     // unique
	Slug                  string                         `json:"slug"`                     // unique
	AddressID             *string                        `json:"address_id"`               // nullable
	Email                 string                         `json:"email"`                    //
	ClickAndCollectOption WarehouseClickAndCollectOption `json:"click_and_collect_option"` // default to "disabled"
	IsPrivate             *bool                          `json:"is_private"`               // default *true
	ModelMetadata

	address       *Address      `json:"-" gorm:"-"` // this field hold data from select related queries
	ShippingZones ShippingZones `json:"-" gorm:"many2many:WarehouseShippingZones"`
}

// WarehouseFilterOption is used to build squirrel queries
type WarehouseFilterOption struct {
	Conditions squirrel.Sqlizer

	ShippingZonesCountries squirrel.Sqlizer // INNER JOIN warehouseShippingZones ON (...) INNER JOIN shippingZones ON (...) WHERE ShippingZones.Countries ...
	ShippingZonesId        squirrel.Sqlizer // INNER JOIN warehouseShippingZones ON (...) INNER JOIN shippingZones ON (...) WHERE ShippingZones.Id ...

	// NOTE: If set, store will use OR ILIKE to check it against:
	//
	// warehouse's name, email, warehouse's relevant address's company name,
	// street address (1/2), city, postal code, phone
	Search string

	SelectRelatedAddress  bool // set true if you want it to attach the `Address` property to returning warehouse(s)
	PrefetchShippingZones bool // set true if you want it to find all shipping zones of found warehouses also
	Distinct              bool // SELECT DISTINCT
}

func (w *WareHouse) GetAddress() *Address {
	return w.address
}

func (w *WareHouse) SetAddress(a *Address) {
	w.address = a
}

type Warehouses []*WareHouse

func (ws Warehouses) IDs() []string {
	return lo.Map(ws, func(w *WareHouse, _ int) string { return w.Id })
}

func (w *WareHouse) String() string {
	return w.Name
}

func (w *WareHouse) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.warehouse.is_valid.%s.app_error",
		"warehouse_id=",
		"WareHouse.IsValid",
	)
	if !IsValidId(w.Id) {
		return outer("id", nil)
	}
	if w.AddressID != nil && !IsValidId(*w.AddressID) {
		return outer("address_id", nil)
	}
	if utf8.RuneCountInString(w.Name) > WAREHOUSE_NAME_MAX_LENGTH {
		return outer("name", &w.Id)
	}
	if len(w.Slug) > WAREHOUSE_SLUG_MAX_LENGTH {
		return outer("slug", &w.Id)
	}
	if w.Email != "" && !IsValidEmail(w.Email) {
		return outer("email", &w.Id)
	}
	if len(w.ClickAndCollectOption) > WAREHOUSE_CLICK_AND_COLLECT_OPTION_MAX_LENGTH ||
		ValidWarehouseClickAndCollectOptionMap[w.ClickAndCollectOption] == "" {
		return outer("click_and_collect_option", &w.Id)
	}
	if w.IsPrivate == nil { // this must be set to true if is left not set
		return outer("is_private", &w.Id)
	}

	return nil
}

func (w *WareHouse) PreSave() {
	if w.Id == "" {
		w.Id = NewId()
	}
	w.Slug = slug.Make(w.Name)
	w.ModelMetadata.PopulateFields()
	w.commonPre()
}

func (w *WareHouse) commonPre() {
	w.Name = SanitizeUnicode(w.Name)
	if w.ClickAndCollectOption == "" {
		w.ClickAndCollectOption = DISABLED
	}
	if w.IsPrivate == nil {
		w.IsPrivate = NewPrimitive(true)
	}
}

func (w *WareHouse) PreUpdate() {
	w.ModelMetadata.PopulateFields()
	w.commonPre()
}

func (w *WareHouse) ToJSON() string {
	return ModelToJson(w)
}

func (w *WareHouse) DeepCopy() *WareHouse {
	res := *w

	if w.AddressID != nil {
		res.AddressID = NewPrimitive(*w.AddressID)
	}
	if w.IsPrivate != nil {
		res.IsPrivate = NewPrimitive(*w.IsPrivate)
	}
	res.address = w.address.DeepCopy()
	res.ShippingZones = w.ShippingZones.DeepCopy()
	return &res
}
