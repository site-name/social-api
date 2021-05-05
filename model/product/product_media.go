package product

import (
	"io"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
)

type ProductMedia struct {
	Id        string `json:"id"`
	ProductID string `json:"product_id"`
	*model.Sortable
}

func (p *ProductMedia) ToJson() string {
	b, _ := json.JSON.Marshal(p)
	return string(b)
}

func ProductMediaFromJson(data io.Reader) *ProductMedia {
	var prd ProductMedia
	err := json.JSON.NewDecoder(data).Decode(&prd)
	if err != nil {
		return nil
	}
	return &prd
}
