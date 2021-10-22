package product_and_discount

import (
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/seo"
	"golang.org/x/text/language"
)

// Max lengths for some fields
const (
	COLLECTION_NAME_MAX_LENGTH           = 250
	COLLECTION_SLUG_MAX_LENGTH           = 255
	COLLECTION_BACKGROUND_ALT_MAX_LENGTH = 128
)

type Collection struct {
	Id                 string  `json:"id"`
	ShopID             string  `json:"shop_id"` // shop that owns this collection
	Name               string  `json:"name"`
	Slug               string  `json:"slug"`
	BackgroundImage    *string `json:"background_image"`
	BackgroundImageAlt string  `json:"background_image_alt"`
	Description        *string `json:"description"`
	model.ModelMetadata
	seo.Seo
}

// CollectionFilterOption is used to build sql queries.
//
// if `SelectAll` is set to true, it finds all collections of given shop, ignores other options too
type CollectionFilterOption struct {
	ShopID    string // single string since we can only view collections of ONLY 1 shop at a time
	SelectAll bool   // if this is true, ignore every other options and find all collections by shop

	Id   *model.StringFilter //
	Name *model.StringFilter //
	Slug *model.StringFilter

	ProductIDs                    []string            // use sub query SELECT ... FROM Collections WHERE Id IN (SELECT CollectionID FROM ... WHERE OtherID = ?)
	VoucherIDs                    []string            // use sub query SELECT ... FROM Collections WHERE Id IN (SELECT CollectionID FROM ... WHERE OtherID = ?)
	SaleIDs                       []string            //
	ChannelListingPublicationDate *model.TimeFilter   // INNER JOIN `CollectionChannelListings`
	ChannelListingChannelSlug     *model.StringFilter // INNER JOIN `CollectionChannelListings` INNER JOIN `Channels`
	ChannelListingChannelIsActive *bool               // INNER JOIN `CollectionChannelListing` INNER JOIN `Channels`
	ChannelListingIsPublished     *bool               // INNER JOIN `CollectionChannelListing`
}

type Collections []*Collection

func (c Collections) IDs() []string {
	var res []string
	for _, item := range c {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (c *Collection) String() string {
	return c.Name
}

func (c *Collection) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel("model.collection.is_valid.%s.app_error", "collection_id=", "Collection.IsValid")
	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(c.ShopID) {
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
	return model.ModelToJson(c)
}

func (c *Collection) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}
	c.Name = model.SanitizeUnicode(c.Name)
	c.Slug = slug.Make(c.Name)
}

func (c *Collection) PreUpdate() {
	c.Name = model.SanitizeUnicode(c.Name)
	c.Slug = slug.Make(c.Name) // ?
}

// CollectionTranslation
type CollectionTranslation struct {
	Id           string  `json:"id"`
	LanguageCode string  `json:"language_code"`
	CollectionID string  `json:"collection_id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	seo.SeoTranslation
}

func (c *CollectionTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel("model.collection_translation.is_valid.%s.app_error", "collection_translation_id=", "CollectionTranslation.IsValid")
	if !model.IsValidId(c.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(c.CollectionID) {
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
	return model.ModelToJson(c)
}

func (c *CollectionTranslation) PreSave() {
	if c.Id == "" {
		c.Id = model.NewId()
	}
	c.Name = model.SanitizeUnicode(c.Name)
}

func (c *CollectionTranslation) PreUpdate() {
	c.Name = model.SanitizeUnicode(c.Name)
}
