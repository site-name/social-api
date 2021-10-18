package gqlmodel

type Category struct {
	ID              string                       `json:"id"`
	SeoTitle        *string                      `json:"seoTitle"`
	SeoDescription  *string                      `json:"seoDescription"`
	Name            string                       `json:"name"`
	Description     *string                      `json:"description"`
	Slug            string                       `json:"slug"`
	ParentID        *string                      `json:"parent"` // Category
	Level           int                          `json:"level"`
	PrivateMetadata []*MetadataItem              `json:"privateMetadata"`
	Metadata        []*MetadataItem              `json:"metadata"`
	Ancestors       *CategoryCountableConnection `json:"ancestors"`
	Products        *ProductCountableConnection  `json:"products"`
	Children        *CategoryCountableConnection `json:"children"`
	BackgroundImage *Image                       `json:"backgroundImage"`
	TranslationID   *string                      `json:"translation"` // CategoryTranslation
}

func (Category) IsNode()               {}
func (Category) IsObjectWithMetadata() {}
