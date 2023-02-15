package api

import (
	"context"

	"github.com/sitename/sitename/model"
)

type DigitalContent struct {
	UseDefaultSettings   bool            `json:"useDefaultSettings"`
	AutomaticFulfillment bool            `json:"automaticFulfillment"`
	ContentFile          string          `json:"contentFile"`
	MaxDownloads         *int32          `json:"maxDownloads"`
	URLValidDays         *int32          `json:"urlValidDays"`
	ID                   string          `json:"id"`
	PrivateMetadata      []*MetadataItem `json:"privateMetadata"`
	Metadata             []*MetadataItem `json:"metadata"`
	d                    *model.DigitalContent

	// ProductVariant       *ProductVariant      `json:"productVariant"`
	// Urls                 []*DigitalContentURL `json:"urls"`
}

func systemDigitalContentToGraphqlDigitalContent(d *model.DigitalContent) *DigitalContent {
	if d == nil {
		return nil
	}

	res := &DigitalContent{
		ID:                   d.Id,
		Metadata:             MetadataToSlice(d.Metadata),
		PrivateMetadata:      MetadataToSlice(d.PrivateMetadata),
		UseDefaultSettings:   *d.UseDefaultSettings,
		AutomaticFulfillment: *d.AutomaticFulfillment,
		ContentFile:          d.ContentFile,
	}
	if d.MaxDownloads != nil {
		res.MaxDownloads = model.NewPrimitive(int32(*d.MaxDownloads))
	}
	if d.UrlValidDays != nil {
		res.URLValidDays = model.NewPrimitive(int32(*d.UrlValidDays))
	}

	return res
}

func (d *DigitalContent) Urls(ctx context.Context) ([]*DigitalContentURL, error) {
	panic("not implemented")
}

func (d *DigitalContent) ProductVariant(ctx context.Context) (*ProductVariant, error) {
	variant, err := ProductVariantByIdLoader.Load(ctx, d.d.ProductVariantID)()
	if err != nil {
		return nil, err
	}
	return SystemProductVariantToGraphqlProductVariant(variant), nil
}
