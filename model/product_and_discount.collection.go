package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"golang.org/x/text/language"
)

// Max lengths for some fields
const (
	COLLECTION_NAME_MAX_LENGTH           = 250
	COLLECTION_SLUG_MAX_LENGTH           = 255
	COLLECTION_BACKGROUND_ALT_MAX_LENGTH = 128
)

type Collection struct {
	Id                 string          `json:"id"`
	ShopID             string          `json:"shop_id"` // shop that owns this collection
	Name               string          `json:"name"`
	Slug               string          `json:"slug"`
	BackgroundImage    *string         `json:"background_image"`
	BackgroundImageAlt string          `json:"background_image_alt"`
	Description        StringInterface `json:"description"`
	ModelMetadata
	Seo
}

// CollectionFilterOption is used to build sql queries.
//
// if `SelectAll` is set to true, it finds all collections of given shop, ignores other options too
type CollectionFilterOption struct {
	ShopID    string // single string since we can only view collections of ONLY 1 shop at a time
	SelectAll bool   // if this is true, ignore every other options and find all collections by shop

	Id   squirrel.Sqlizer
	Name squirrel.Sqlizer
	Slug squirrel.Sqlizer

	ProductID squirrel.Sqlizer // SELECT * FROM Collections INNER JOIN ProductCollections ON (...) WHERE ProductCollections.ProductID ...
	VoucherID squirrel.Sqlizer // SELECT * FROM Collections INNER JOIN VoucherCollections ON (...) WHERE VoucherCollections.VoucherID ...
	SaleID    squirrel.Sqlizer // SELECT * FROM Collections INNER JOIN SaleCollections ON (Collections.Id = SaleCollections.CollectionID) WHERE SaleCollections.SaleID ...

	ChannelListingPublicationDate squirrel.Sqlizer // INNER JOIN `CollectionChannelListings`
	ChannelListingChannelSlug     squirrel.Sqlizer // INNER JOIN `CollectionChannelListings` INNER JOIN `Channels`
	ChannelListingChannelIsActive *bool            // INNER JOIN `CollectionChannelListing` INNER JOIN `Channels`
	ChannelListingIsPublished     *bool            // INNER JOIN `CollectionChannelListing`
}

type Collections []*Collection

func (c *Collection) DeepCopy() *Collection {
	if c == nil {
		return nil
	}

	res := *c
	res.Description = c.Description.DeepCopy()
	res.ModelMetadata = c.ModelMetadata.DeepCopy()
	if c.BackgroundImage != nil {
		res.BackgroundImage = NewPrimitive(*c.BackgroundImage)
	}

	return &res
}

func (c Collections) IDs() []string {
	return lo.Map(c, func(o *Collection, _ int) string { return o.Id })
}

func (c Collections) DeepCopy() Collections {
	return lo.Map(c, func(o *Collection, _ int) *Collection { return o.DeepCopy() })
}

func (c *Collection) String() string {
	return c.Name
}

func (c *Collection) IsValid() *AppError {
	outer := CreateAppErrorForModel("collection.is_valid.%s.app_error", "collection_id=", "Collection.IsValid")
	if !IsValidId(c.Id) {
		return outer("id", nil)
	}
	if !IsValidId(c.ShopID) {
		return outer("shop_id", &c.Id)
	}
	if utf8.RuneCountInString(c.Name) > COLLECTION_NAME_MAX_LENGTH {
		return outer("name", &c.Id)
	}
	if utf8.RuneCountInString(c.Slug) > COLLECTION_SLUG_MAX_LENGTH {
		return outer("slug", &c.Id)
	}
	if utf8.RuneCountInString(c.BackgroundImageAlt) > COLLECTION_BACKGROUND_ALT_MAX_LENGTH {
		return outer("background_image_alt", &c.Id)
	}

	return nil
}

func (c *Collection) ToJSON() string {
	return ModelToJson(c)
}

func (c *Collection) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
	c.Name = SanitizeUnicode(c.Name)
	c.Slug = slug.Make(c.Name)
}

func (c *Collection) PreUpdate() {
	c.Name = SanitizeUnicode(c.Name)
	c.Slug = slug.Make(c.Name) // ?
}

// CollectionTranslation
type CollectionTranslation struct {
	Id           string  `json:"id"`
	LanguageCode string  `json:"language_code"`
	CollectionID string  `json:"collection_id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	SeoTranslation
}

func (c *CollectionTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel("collection_translation.is_valid.%s.app_error", "collection_translation_id=", "CollectionTranslation.IsValid")
	if !IsValidId(c.Id) {
		return outer("id", nil)
	}
	if !IsValidId(c.CollectionID) {
		return outer("collection_id", &c.Id)
	}
	if utf8.RuneCountInString(c.Name) > COLLECTION_NAME_MAX_LENGTH {
		return outer("name", &c.Id)
	}
	if tag, err := language.Parse(c.LanguageCode); err != nil || !strings.EqualFold(tag.String(), c.LanguageCode) {
		return outer("language_code", &c.Id)
	}

	return nil
}

func (c *CollectionTranslation) String() string {
	return c.Name
}

func (c *CollectionTranslation) ToJSON() string {
	return ModelToJson(c)
}

func (c *CollectionTranslation) PreSave() {
	if c.Id == "" {
		c.Id = NewId()
	}
	c.Name = SanitizeUnicode(c.Name)
}

func (c *CollectionTranslation) PreUpdate() {
	c.Name = SanitizeUnicode(c.Name)
}
