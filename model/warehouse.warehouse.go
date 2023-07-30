package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"gorm.io/gorm"
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
	Id                    string                         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name                  string                         `json:"name" gorm:"type:varchar(250);column:Name"`
	Slug                  string                         `json:"slug" gorm:"type:varchar(255);column:Slug;unique"`                              // unique
	AddressID             *string                        `json:"address_id" gorm:"type:uuid;column:AddressID"`                                  // nullable
	Email                 string                         `json:"email" gorm:"type:varchar(254);column:Email"`                                   //
	ClickAndCollectOption WarehouseClickAndCollectOption `json:"click_and_collect_option" gorm:"column:ClickAndCollectOption;type:varchar(30)"` // default to "disabled"
	IsPrivate             *bool                          `json:"is_private" gorm:"default:true;column:IsPrivate"`                               // default *true
	CreateAt              int64                          `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	ModelMetadata

	address       *Address      `json:"-" gorm:"-"` // this field hold data from select related queries
	ShippingZones ShippingZones `json:"-" gorm:"many2many:WarehouseShippingZones"`
}

func (c *WareHouse) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *WareHouse) BeforeUpdate(_ *gorm.DB) error { c.PreUpdate(); return c.IsValid() }
func (c *WareHouse) TableName() string             { return WarehouseTableName }

// WarehouseFilterOption is used to build squirrel queries
type WarehouseFilterOption struct {
	Conditions squirrel.Sqlizer

	ShippingZonesCountries squirrel.Sqlizer // INNER JOIN WarehouseShippingZones ON (...) INNER JOIN shippingZones ON (...) WHERE ShippingZones.Countries ...
	ShippingZonesId        squirrel.Sqlizer // INNER JOIN WarehouseShippingZones ON (...) INNER JOIN shippingZones ON (...) WHERE ShippingZones.Id ...

	// NOTE: If set, store will use OR ILIKE to check it against:
	//
	// warehouse's name, email, warehouse's relevant address's company name,
	// street address (1/2), city, postal code, phone
	Search string

	SelectRelatedAddress  bool // set true if you want it to attach the `Address` property to returning warehouse(s)
	PrefetchShippingZones bool // set true if you want it to find all shipping zones of found warehouses also
	Distinct              bool // SELECT DISTINCT
}

type WarehouseShippingZone struct {
	WarehouseID    string
	ShippingZoneID string
}

func (w *WareHouse) GetAddress() *Address  { return w.address }
func (w *WareHouse) SetAddress(a *Address) { w.address = a }

type Warehouses []*WareHouse

func (ws Warehouses) IDs() []string {
	return lo.Map(ws, func(w *WareHouse, _ int) string { return w.Id })
}

func (w *WareHouse) String() string {
	return w.Name
}

func (w *WareHouse) IsValid() *AppError {
	if w.AddressID != nil && !IsValidId(*w.AddressID) {
		return NewAppError("Warehouse.IsValid", "model.warehouse.is_valid.address_id.app_error", nil, "please provide valid address id", http.StatusBadRequest)
	}
	if w.Email != "" && !IsValidEmail(w.Email) {
		return NewAppError("Warehouse.IsValid", "model.warehouse.is_valid.email.app_error", nil, "please provide valid email", http.StatusBadRequest)
	}
	if !w.ClickAndCollectOption.IsValid() {
		return NewAppError("Warehouse.IsValid", "model.warehouse.is_valid.click_and_collect_option.app_error", nil, "please provide valid click and collect option", http.StatusBadRequest)
	}

	return nil
}

func (w *WareHouse) PreSave() {
	w.commonPre()
	w.Slug = slug.Make(w.Name)
}

func (w *WareHouse) commonPre() {
	w.Name = SanitizeUnicode(w.Name)
	if w.ClickAndCollectOption == "" {
		w.ClickAndCollectOption = DISABLED
	}
	if w.IsPrivate == nil {
		w.IsPrivate = NewPrimitive(true)
	}
	w.ModelMetadata.PopulateFields()

}

func (w *WareHouse) PreUpdate() {
	w.commonPre()
	w.CreateAt = 0 // prevent updating
}

func (w *WareHouse) DeepCopy() *WareHouse {
	res := *w

	res.AddressID = CopyPointer(w.AddressID)
	res.IsPrivate = CopyPointer(w.IsPrivate)

	if w.address != nil {
		res.address = w.address.DeepCopy()
	}
	res.ShippingZones = w.ShippingZones.DeepCopy()
	return &res
}
