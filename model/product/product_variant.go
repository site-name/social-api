package product

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/measurement"
)

type ProductVariant struct {
	Id                   string          `json:"id"`
	Name                 string          `json:"name"`
	ProductID            string          `json:"product_id"`
	Sku                  string          `json:"sku"`
	Weight               *float32        `json:"weight"`
	WeightUnit           string          `json:"weight_unit"`
	TrackInventory       *bool           `json:"track_inventory"`
	Medias               []*ProductMedia `json:"medias" db:"-"`
	*model.Sortable      `db:"-"`
	*model.ModelMetadata `db:"-"`
}

func (p *ProductVariant) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.product_variant.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "product_variant_id=" + p.Id
	}

	return model.NewAppError("ProductVariant.IsValid", id, nil, details, http.StatusBadRequest)
}

func (p *ProductVariant) IsValid() *model.AppError {
	if !model.IsValidId(p.Id) {
		return p.createAppError("id")
	}
	if !model.IsValidId(p.ProductID) {
		return p.createAppError("product_id")
	}
	if len(p.Sku) > PRODUCT_VARIANT_SKU_MAX_LENGTH {
		return p.createAppError("sku")
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_VARIANT_NAME_MAX_LENGTH {
		return p.createAppError("name")
	}
	if p.Weight != nil && *p.Weight < 0 {
		return p.createAppError("weight")
	}
	if p.WeightUnit != "" {
		if _, ok := measurement.WEIGHT_UNIT_CONVERSION[strings.ToLower(p.WeightUnit)]; !ok {
			return p.createAppError("weight_unit")
		}
	}

	return nil
}

func (p *ProductVariant) String() string {
	return p.Name
}

func (p *ProductVariant) ToJson() string {
	b, _ := json.JSON.Marshal(p)
	return string(b)
}

func (p *ProductVariant) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
	p.Name = model.SanitizeUnicode(p.Name)
	if p.TrackInventory == nil {
		p.TrackInventory = model.NewBool(true)
	}
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
}

func (p *ProductVariant) PreUpdate() {
	p.Name = model.SanitizeUnicode(p.Name)
	if p.Weight != nil && p.WeightUnit == "" {
		p.WeightUnit = measurement.STANDARD_WEIGHT_UNIT
	}
}

func ProductVariantFromJson(data io.Reader) *ProductVariant {
	var prd ProductVariant
	err := json.JSON.NewDecoder(data).Decode(&prd)
	if err != nil {
		return nil
	}
	return &prd
}

func (p *ProductVariant) GetPrice(
	product *Product,
	collection []*Collection,
	channel *channel.Channel,
	channelListing *ProductChannelListing,
	discounts *[]*model.DiscountInfo) {

}
