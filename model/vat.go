package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
)

type Vat struct {
	Id          string      `json:"id"`
	CountryCode CountryCode `json:"country_code"` // db index
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
	Data StringInterface `json:"data"`
}

type VatFilterOptions struct {
	Id          squirrel.Sqlizer
	CountryCode squirrel.Sqlizer
}

func (v *Vat) PreSave() {
	if !IsValidId(v.Id) {
		v.Id = NewId()
	}
	if v.Data == nil {
		v.Data = StringInterface{}
	}
}

func (v *Vat) PreUpdate() {
	if v.Data == nil {
		v.Data = StringInterface{}
	}
}

func (v *Vat) String() string {
	return v.CountryCode.String()
}

func (v *Vat) IsValid() *AppError {
	if !IsValidId(v.Id) {
		return NewAppError("Vat.IsValid", "model.vat.is_valid.id.app_error", nil, v.Id+" is invalid id", http.StatusBadRequest)
	}
	if !v.CountryCode.IsValid() {
		return NewAppError("Vat.IsValid", "model.vat.is_valid.country_code.app_error", nil, v.Id+" is invalid id", http.StatusBadRequest)
	}

	return nil
}
