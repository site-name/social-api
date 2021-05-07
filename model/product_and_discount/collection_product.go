package product_and_discount

import (
	"io"

	"github.com/sitename/sitename/model"
)

type CollectionProduct struct {
	Id           string `json:"id"`
	CollectionID string `json:"collection_id"`
	ProductID    string `json:"product_id"`
}

func (c *CollectionProduct) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel("model.collection_product.is_valid.%s.app_error", "collection_product_id=", "CollectionProduct.IsValid")
	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(c.CollectionID) {
		return outer("collection_id", &c.Id)
	}
	if !model.IsValidId(c.ProductID) {
		return outer("product_id", &c.Id)
	}

	return nil
}

func (c *CollectionProduct) ToJson() string {
	return model.ModelToJson(c)
}

func CollectionProductFromJson(data io.Reader) *CollectionProduct {
	var c CollectionProduct
	model.ModelFromJson(&c, data)
	return &c
}

func (c *CollectionProduct) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}
}
