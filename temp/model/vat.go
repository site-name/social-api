package model

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"gorm.io/gorm"
)

type Vat struct {
	Id          string      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CountryCode CountryCode `json:"country_code" gorm:"type:varchar(5);column:CountryCode;index:countrycode_key"` // db index
	// NOTE:
	//
	// Data may contains keys of:
	//  "standard_rate": float64,
	//  "reduced_rates": map[string]float64{
	//    "pharmaceuticals": float64,
	//    "medical": float64,
	//    "passenger transport": float64,
	//    "newspapers": float64,
	//    "hotels": float64,
	//    "restaurants": float64,
	//    "admission to cultural events": float64,
	//    "admission to sporting events": float64,
	//    "admission to entertainment events": float64,
	//    "foodstuffs": float64,
	//  }
	Data StringInterface `json:"data" gorm:"type:jsonb;column:Data"`
}

func (t *Vat) TableName() string             { return VatTableName }
func (t *Vat) BeforeCreate(_ *gorm.DB) error { t.commonPre(); return t.IsValid() }
func (t *Vat) BeforeUpdate(_ *gorm.DB) error { t.commonPre(); return t.IsValid() }

type VatFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (v *Vat) commonPre() {
	if v.Data == nil {
		v.Data = StringInterface{}
	}
}

func (v *Vat) String() string {
	return v.CountryCode.String()
}

func (v *Vat) IsValid() *AppError {
	if !v.CountryCode.IsValid() {
		return NewAppError("Vat.IsValid", "model.vat.is_valid.country_code.app_error", nil, v.Id+" is invalid id", http.StatusBadRequest)
	}

	return nil
}
