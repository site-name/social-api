package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

type CollectionProduct struct {
	Id           string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CollectionID string `json:"collection_id" gorm:"type:uuid;column:CollectionID"`
	ProductID    string `json:"product_id" gorm:"type:uuid;column:ProductID"`
	Sortable

	collection *Collection `gorm:"-"` // get populated if CollectionProductFilterOptions.SelectRelatedCollection is true
	product    *Product    `gorm:"-"` // get populated if CollectionProductFilterOptions.SelectRelatedProduct is true
}

// column names of collection-product table
const (
	CollectionProductColumnId           = "Id"
	CollectionProductColumnCollectionID = "CollectionID"
	CollectionProductColumnProductID    = "ProductID"
)

func (c *CollectionProduct) GetCollection() *Collection    { return c.collection }
func (c *CollectionProduct) SetCollection(col *Collection) { c.collection = col }
func (c *CollectionProduct) GetProduct() *Product          { return c.product }
func (c *CollectionProduct) SetProduct(p *Product)         { c.product = p }
func (c *CollectionProduct) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *CollectionProduct) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *CollectionProduct) TableName() string             { return CollectionProductRelationTableName }

type CollectionProductFilterOptions struct {
	Conditions              squirrel.Sqlizer
	SelectRelatedCollection bool
	SelectRelatedProduct    bool
}

func (c *CollectionProduct) IsValid() *AppError {
	if !IsValidId(c.CollectionID) {
		return NewAppError("CollectionProduct.IsValid", "model.collection_product.is_valid.collection_id.app_error", nil, "please provide valid collection id", http.StatusBadRequest)
	}
	if !IsValidId(c.ProductID) {
		return NewAppError("CollectionProduct.IsValid", "model.collection_product.is_valid.product_id.app_error", nil, "please provide valid product id", http.StatusBadRequest)
	}

	return nil
}

func (c *CollectionProduct) DeepCopy() *CollectionProduct {
	if c == nil {
		return nil
	}

	res := *c
	if c.SortOrder != nil {
		*res.SortOrder = *c.SortOrder
	}
	if c.collection != nil {
		res.collection = c.collection.DeepCopy()
	}
	if c.product != nil {
		res.product = c.product.DeepCopy()
	}

	return &res
}
