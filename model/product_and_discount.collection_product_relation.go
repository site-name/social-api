package model

import (
	"github.com/Masterminds/squirrel"
)

type CollectionProduct struct {
	Id           string `json:"id"`
	CollectionID string `json:"collection_id"`
	ProductID    string `json:"product_id"`

	collection *Collection `db:"-"` // get populated if CollectionProductFilterOptions.SelectRelatedCollection is true
	product    *Product    `db:"-"` // get populated if CollectionProductFilterOptions.SelectRelatedProduct is true
}

func (c *CollectionProduct) GetCollection() *Collection {
	return c.collection
}

func (c *CollectionProduct) SetCollection(col *Collection) {
	c.collection = col
}

func (c *CollectionProduct) GetProduct() *Product {
	return c.product
}

func (c *CollectionProduct) SetProduct(p *Product) {
	c.product = p
}

type CollectionProductFilterOptions struct {
	CollectionID squirrel.Sqlizer
	ProductID    squirrel.Sqlizer

	SelectRelatedCollection bool
	SelectRelatedProduct    bool
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
	if c.GetCollection() != nil {
		res.collection = c.GetCollection().DeepCopy()
	}
	if c.GetProduct() != nil {
		res.product = c.GetProduct().DeepCopy()
	}

	return &res
}
