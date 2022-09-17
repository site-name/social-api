package model

import (
	"github.com/Masterminds/squirrel"
)

type CollectionProduct struct {
	Id           string `json:"id"`
	CollectionID string `json:"collection_id"`
	ProductID    string `json:"product_id"`

	Collection *Collection `json:"-" db:"-"`
	Product    *Product    `json:"-" db:"-"`
}

type CollectionProductFilterOptions struct {
	CollectionID squirrel.Sqlizer
	ProductID    squirrel.Sqlizer

	SelectRelatedCollection bool
}

func (c *CollectionProduct) IsValid() *AppError {
	outer := CreateAppErrorForModel("collection_product.is_valid.%s.app_error", "collection_product_id=", "CollectionProduct.IsValid")
	if !IsValidId(c.Id) {
		return outer("id", nil)
	}
	if !IsValidId(c.CollectionID) {
		return outer("collection_id", &c.Id)
	}
	if !IsValidId(c.ProductID) {
		return outer("product_id", &c.Id)
	}

	return nil
}

func (c *CollectionProduct) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
}

func (c *CollectionProduct) DeepCopy() *CollectionProduct {
	if c == nil {
		return nil
	}

	res := *c
	if c.Collection != nil {
		res.Collection = c.Collection.DeepCopy()
	}
	if c.Product != nil {
		res.Product = c.Product.DeepCopy()
	}

	return &res
}
