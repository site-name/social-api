package warehouse

import (
	"io"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/shipping"
)

// max lengths for some warehouse's fields
const (
	WAREHOUSE_NAME_MAX_LENGTH = 250
	WAREHOUSE_SLUG_MAX_LENGTH = 255
)

type WareHouse struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`       // unique
	Slug      string  `json:"slug"`       // unique
	AddressID *string `json:"address_id"` // nullable
	Email     string  `json:"email"`
	model.ModelMetadata

	Address       *account.Address         `json:"-" db:"-"` // this field hold data from select related queries
	ShippingZones []*shipping.ShippingZone `json:"-" db:"-"` // this field hold data from prefetch_related queries
}

// WarehouseFilterOption is used to build squirrel queries
type WarehouseFilterOption struct {
	Id                     *model.StringFilter
	Name                   *model.StringFilter
	Slug                   *model.StringFilter
	AddressID              *model.StringFilter
	Email                  *model.StringFilter
	ShippingZonesCountries *model.StringFilter // join shipping zone table

	SelectRelatedAddress  bool // set true if you want it to attach the `Address` property also also
	PrefetchShippingZones bool // set true if you want it to find all shipping zones of found warehouses also
}

type Warehouses []*WareHouse

func (ws Warehouses) IDs() []string {
	res := []string{}
	meetMap := map[string]bool{}

	for _, warehouse := range ws {
		if warehouse != nil && !meetMap[warehouse.Id] {
			res = append(res, warehouse.Id)
			meetMap[warehouse.Id] = true
		}
	}

	return nil
}

func (w *WareHouse) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.warehouse.is_valid.%s.app_error",
		"warehouse_id=",
		"WareHouse.IsValid",
	)
	if !model.IsValidId(w.Id) {
		return outer("id", nil)
	}
	if w.AddressID != nil && !model.IsValidId(*w.AddressID) {
		return outer("address_id", nil)
	}
	if utf8.RuneCountInString(w.Name) > WAREHOUSE_NAME_MAX_LENGTH {
		return outer("name", &w.Id)
	}
	if len(w.Slug) > WAREHOUSE_SLUG_MAX_LENGTH {
		return outer("slug", &w.Id)
	}
	if w.Email != "" && !model.IsValidEmail(w.Email) {
		return outer("email", &w.Id)
	}

	return nil
}

func (w *WareHouse) PreSave() {
	if w.Id == "" {
		w.Id = model.NewId()
	}
	w.Name = model.SanitizeUnicode(w.Name)
	w.Slug = slug.Make(w.Name)
	w.ModelMetadata.PreSave()
}

func (w *WareHouse) PreUpdate() {
	w.Name = model.SanitizeUnicode(w.Name)
	w.Slug = slug.Make(w.Name)
	w.ModelMetadata.PreUpdate()
}

func (w *WareHouse) ToJson() string {
	return model.ModelToJson(w)
}

func WareHouseFromJson(data io.Reader) *WareHouse {
	var w WareHouse
	model.ModelFromJson(&w, data)
	return &w
}
