package gqlmodel

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

// original implementation

// type Collection struct {
// 	ID              string                      `json:"id"`
// 	SeoTitle        *string                     `json:"seoTitle"`
// 	SeoDescription  *string                     `json:"seoDescription"`
// 	Name            string                      `json:"name"`
// 	Description     model.StringInterface       `json:"description"`
// 	Slug            string                      `json:"slug"`
// 	PrivateMetadata []*MetadataItem             `json:"privateMetadata"`
// 	Metadata        []*MetadataItem             `json:"metadata"`
// 	Channel         *string                     `json:"channel"`
// 	Products        *ProductCountableConnection `json:"products"`
// 	BackgroundImage *Image                      `json:"backgroundImage"`
// 	Translation     *CollectionTranslation      `json:"translation"`
// 	ChannelListings []*CollectionChannelListing `json:"channelListings"`
// }

// func (Collection) IsNode()               {}
// func (Collection) IsObjectWithMetadata() {}

type Collection struct {
	ID                string                             `json:"id"`
	SeoTitle          *string                            `json:"seoTitle"`
	SeoDescription    *string                            `json:"seoDescription"`
	Name              string                             `json:"name"`
	Description       model.StringInterface              `json:"description"`
	Slug              string                             `json:"slug"`
	PrivateMetadata   []*MetadataItem                    `json:"privateMetadata"`
	Metadata          []*MetadataItem                    `json:"metadata"`
	Channel           *string                            `json:"channel"`
	Products          func() *ProductCountableConnection `json:"products"`
	BackgroundImage   func() *Image                      `json:"backgroundImage"`
	TranslationID     *string                            `json:"translation"`
	ChannelListingIDs []string                           `json:"channelListings"`
}

func (Collection) IsNode()               {}
func (Collection) IsObjectWithMetadata() {}

func SystemCollectionToGraphqlCollection(c *product_and_discount.Collection) *Collection {
	if c == nil {
		return nil
	}

	res := &Collection{
		ID:              c.Id,
		SeoTitle:        c.SeoTitle,
		SeoDescription:  c.SeoDescription,
		Name:            c.Name,
		Description:     c.Description,
		Slug:            c.Slug,
		PrivateMetadata: MapToGraphqlMetaDataItems(c.PrivateMetadata),
		Metadata:        MapToGraphqlMetaDataItems(c.Metadata),
	}

	return res
}
