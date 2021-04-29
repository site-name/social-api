package model

// max length for some fields
const ()

type Category struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	ParentID           string `json:"parent_id"`
	BackgroundImage    string `json:"background_image"`
	BackgroundImageAlt string `json:"background_image_alt"`
	*Seo
	*ModelMetadata
}

type CategoryTranslation struct {
	Id           string `json:"id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	*SeoTranslation
}
