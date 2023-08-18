package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
)

type Collection struct {
	Id                 string          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:Id"`
	Name               string          `json:"name" gorm:"type:varchar(250);unique;not null;column:Name"`
	Slug               string          `json:"slug" gorm:"type:varchar(255);uniqueIndex:slug_unique_key;column:Slug"`
	BackgroundImage    *string         `json:"background_image" gorm:"type:varchar(200);column:BackgroundImage"`
	BackgroundImageAlt string          `json:"background_image_alt" gorm:"type:varchar(128);column:BackgroundImageAlt"`
	Description        StringInterface `json:"description" gorm:"type:jsonb;column:Description"`
	ModelMetadata
	Seo

	Sales              Sales                `json:"-" gorm:"many2many:SaleCollections"`
	Vouchers           Vouchers             `json:"-" gorm:"many2many:VoucherCollections"`
	Products           Products             `json:"-" gorm:"many2many:ProductCollections"`
	CollectionProducts []*CollectionProduct `json:"-" gorm:"foreignKey:CollectionID"`
}

func (c *Collection) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *Collection) BeforeUpdate(_ *gorm.DB) error { c.PreUpdate(); return c.IsValid() }
func (c *Collection) TableName() string             { return CollectionTableName }

// CollectionFilterOption is used to build sql queries.
//
// if `SelectAll` is set to true, it finds all collections of given shop, ignores other options too
type CollectionFilterOption struct {
	Conditions squirrel.Sqlizer
	// E.g:
	//  []string{"Sales", "CollectionProducts"} // etc...
	Preload []string

	ProductID squirrel.Sqlizer // INNER JOIN ProductCollections ON ... WHERE ProductCollections.ProductID ...
	VoucherID squirrel.Sqlizer // INNER JOIN VoucherCollections ON ... WHERE VoucherCollections.VoucherID ...
	SaleID    squirrel.Sqlizer // INNER JOIN SaleCollections    ON ... WHERE SaleCollections.SaleID ...

	ChannelListingPublicationDate squirrel.Sqlizer // INNER JOIN `CollectionChannelListings` ON ... WHERE CollectionChannelListings.PublicationDate ...
	ChannelListingChannelSlug     squirrel.Sqlizer // INNER JOIN `CollectionChannelListings` ON ... INNER JOIN `Channels` ON ... WHERE Channels.Slug ...
	ChannelListingChannelIsActive squirrel.Sqlizer // INNER JOIN `CollectionChannelListings` ON ... INNER JOIN `Channels` ON ... WHERE Channels.IsActive ...
	ChannelListingIsPublished     squirrel.Sqlizer // INNER JOIN `CollectionChannelListings` ON ... WHERE CollectionChannelListings.IsPublished ...
}

type Collections []*Collection

func (ps Collections) Contains(c *Collection) bool {
	return c != nil && lo.SomeBy(ps, func(ct *Collection) bool { return ct != nil && ct.Id == c.Id })
}

func (c *Collection) DeepCopy() *Collection {
	if c == nil {
		return nil
	}

	res := *c
	if c.Description != nil {
		res.Description = c.Description.DeepCopy()
	}
	res.ModelMetadata = c.ModelMetadata.DeepCopy()
	if c.BackgroundImage != nil {
		res.BackgroundImage = NewPrimitive(*c.BackgroundImage)
	}

	return &res
}

func (c Collections) IDs() util.AnyArray[string] {
	return lo.Map(c, func(o *Collection, _ int) string { return o.Id })
}
func (c Collections) DeepCopy() Collections {
	return lo.Map(c, func(o *Collection, _ int) *Collection { return o.DeepCopy() })
}
func (c *Collection) String() string     { return c.Name }
func (c *Collection) IsValid() *AppError { return nil }
func (c *Collection) PreSave()           { c.commonPre(); c.Slug = slug.Make(c.Name) }
func (c *Collection) commonPre() {
	c.Seo.commonPre()
	if c.Description == nil {
		c.Description = StringInterface{}
	}
	c.Name = SanitizeUnicode(c.Name)
}
func (c *Collection) PreUpdate() {
	c.commonPre()
}

// CollectionTranslation
type CollectionTranslation struct {
	Id           string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode LanguageCodeEnum `json:"language_code" gorm:"type:varchar(5);column:LanguageCode"`
	CollectionID string           `json:"collection_id" gorm:"type:uuid;column:CollectionID"`
	Name         string           `json:"name" gorm:"type:varchar(250);column:Name"`
	Description  StringInterface  `json:"description" gorm:"type:jsonb;column:Description"`
	SeoTranslation
}

func (c *CollectionTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *CollectionTranslation) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *CollectionTranslation) TableName() string             { return CollectionTranslationTableName }

func (c *CollectionTranslation) IsValid() *AppError {
	if !IsValidId(c.CollectionID) {
		return NewAppError("CollectionTranslation.IsValid", "model.collection_translation.is_valid.collection_id.app_error", nil, "please provide valid collection id", http.StatusBadRequest)
	}
	if !c.LanguageCode.IsValid() {
		return NewAppError("CollectionTranslation.IsValid", "model.collection_translation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}

func (c *CollectionTranslation) String() string {
	return c.Name
}

func (c *CollectionTranslation) commonPre() {
	c.SeoTranslation.commonPre()
	c.Name = SanitizeUnicode(c.Name)
	if c.Description == nil {
		c.Description = StringInterface{}
	}
}
