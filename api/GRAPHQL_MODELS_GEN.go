package api

import (
	"fmt"
	"net/http"
	"unsafe"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type AccountAddressCreate struct {
	User    *User    `json:"user"`
	Address *Address `json:"address"`
}

type AccountAddressDelete struct {
	Address *Address
	User    *User
}

type AccountAddressUpdate struct {
	User    *User    `json:"user"`
	Address *Address `json:"address"`
}

type AccountDelete struct {
	User *User `json:"user"`
}

type AccountError struct {
	Field   *string          `json:"field"`
	Message *string          `json:"message"`
	Code    AccountErrorCode `json:"code"`
}

type AccountInput struct {
	FirstName              *string           `json:"firstName"`
	LastName               *string           `json:"lastName"`
	DefaultBillingAddress  *AddressInput     `json:"defaultBillingAddress"`
	DefaultShippingAddress *AddressInput     `json:"defaultShippingAddress"`
	LanguageCode           *LanguageCodeEnum `json:"languageCode"`
}

type AccountRegister struct {
	RequiresConfirmation *bool `json:"requiresConfirmation"`
	User                 *User `json:"user"`
}

type AccountRegisterInput struct {
	FirstName    *string           `json:"firstName"`
	LastName     *string           `json:"lastName"`
	UserName     string            `json:"userName"`
	Email        string            `json:"email"`
	Password     string            `json:"password"`
	RedirectURL  *string           `json:"redirectUrl"`
	LanguageCode *LanguageCodeEnum `json:"languageCode"`
	Metadata     []*MetadataInput  `json:"metadata"`
	Channel      *string           `json:"channel"`
}

type AccountRequestDeletion struct {
	Ok bool `json:"ok"`
}

type AccountSetDefaultAddress struct {
	User *User `json:"user"`
}

type AccountUpdate struct {
	User *User `json:"user"`
}

type AddressCreate struct {
	User    *User           `json:"user"`
	Errors  []*AccountError `json:"errors"`
	Address *Address        `json:"address"`
}

type AddressDelete struct {
	User    *User           `json:"user"`
	Errors  []*AccountError `json:"errors"`
	Address *Address        `json:"address"`
}

type AddressInput struct {
	FirstName      *string      `json:"firstName"`
	LastName       *string      `json:"lastName"`
	CompanyName    *string      `json:"companyName"`
	StreetAddress1 *string      `json:"streetAddress1"`
	StreetAddress2 *string      `json:"streetAddress2"`
	City           *string      `json:"city"`
	CityArea       *string      `json:"cityArea"`
	PostalCode     *string      `json:"postalCode"`
	Country        *CountryCode `json:"country"`
	CountryArea    *string      `json:"countryArea"`
	Phone          *string      `json:"phone"`
}

func (a *AddressInput) Validate(where string) *model.AppError {
	// validate input country
	if country := a.Country; country == nil || !country.IsValid() {
		return model.NewAppError(where, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "country"}, "country field is required", http.StatusBadRequest)
	}

	// validate input phone
	if phone := a.Phone; phone != nil {
		_, ok := util.ValidatePhoneNumber(*phone, a.Country.String())
		if !ok {
			return model.NewAppError(where, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "phone"}, fmt.Sprintf("phone number value %v is invalid", *phone), http.StatusBadRequest)
		}
	}

	return nil
}

// NOTE: must be called after calling Validate().
//
// The returned boolean value indicates if the given address is really changed
func (a *AddressInput) PatchAddress(addr *model.Address) bool {
	changed := true

	switch {
	case a.Phone != nil && *a.Phone != "" && addr.Phone != *a.Phone:
		addr.Phone = *a.Phone
		fallthrough

	case *a.Country != addr.Country:
		addr.Country = *a.Country
		fallthrough

	case a.FirstName != nil && *a.FirstName != addr.FirstName:
		addr.FirstName = *a.FirstName
		fallthrough

	case a.LastName != nil && *a.LastName != addr.LastName:
		addr.LastName = *a.LastName
		fallthrough

	case a.CompanyName != nil && *a.CompanyName != addr.CompanyName:
		addr.CompanyName = *a.CompanyName
		fallthrough

	case a.StreetAddress1 != nil && *a.StreetAddress1 != addr.StreetAddress1:
		addr.StreetAddress1 = *a.StreetAddress1
		fallthrough

	case a.StreetAddress2 != nil && *a.StreetAddress2 != addr.StreetAddress2:
		addr.StreetAddress2 = *a.StreetAddress2
		fallthrough

	case a.City != nil && *a.City != addr.City:
		addr.City = *a.City
		fallthrough

	case a.CityArea != nil && *a.CityArea != addr.CityArea:
		addr.CityArea = *a.CityArea
		fallthrough

	case a.PostalCode != nil && *a.PostalCode != addr.PostalCode:
		addr.PostalCode = *a.PostalCode
		fallthrough

	case a.CountryArea != nil && *a.CountryArea != addr.CountryArea:
		addr.CountryArea = *a.CountryArea

	default:
		changed = false
	}

	return changed
}

type AddressSetDefault struct {
	User   *User           `json:"user"`
	Errors []*AccountError `json:"errors"`
}

type AddressUpdate struct {
	User    *User           `json:"user"`
	Errors  []*AccountError `json:"errors"`
	Address *Address        `json:"address"`
}

type AddressValidationData struct {
	CountryCode        *string        `json:"countryCode"`
	CountryName        *string        `json:"countryName"`
	AddressFormat      *string        `json:"addressFormat"`
	AddressLatinFormat *string        `json:"addressLatinFormat"`
	AllowedFields      []string       `json:"allowedFields"`
	RequiredFields     []string       `json:"requiredFields"`
	UpperFields        []string       `json:"upperFields"`
	CountryAreaType    *string        `json:"countryAreaType"`
	CountryAreaChoices []*ChoiceValue `json:"countryAreaChoices"`
	CityType           *string        `json:"cityType"`
	CityChoices        []*ChoiceValue `json:"cityChoices"`
	CityAreaType       *string        `json:"cityAreaType"`
	CityAreaChoices    []*ChoiceValue `json:"cityAreaChoices"`
	PostalCodeType     *string        `json:"postalCodeType"`
	PostalCodeMatchers []string       `json:"postalCodeMatchers"`
	PostalCodeExamples []string       `json:"postalCodeExamples"`
	PostalCodePrefix   *string        `json:"postalCodePrefix"`
}

type App struct {
	ID               string          `json:"id"`
	Name             *string         `json:"name"`
	Created          *DateTime       `json:"created"`
	IsActive         *bool           `json:"isActive"`
	Permissions      []*Permission   `json:"permissions"`
	Tokens           []*AppToken     `json:"tokens"`
	PrivateMetadata  []*MetadataItem `json:"privateMetadata"`
	Metadata         []*MetadataItem `json:"metadata"`
	Type             *AppTypeEnum    `json:"type"`
	Webhooks         []*Webhook      `json:"webhooks"`
	AboutApp         *string         `json:"aboutApp"`
	DataPrivacy      *string         `json:"dataPrivacy"`
	DataPrivacyURL   *string         `json:"dataPrivacyUrl"`
	HomepageURL      *string         `json:"homepageUrl"`
	SupportURL       *string         `json:"supportUrl"`
	ConfigurationURL *string         `json:"configurationUrl"`
	AppURL           *string         `json:"appUrl"`
	Version          *string         `json:"version"`
	AccessToken      *string         `json:"accessToken"`
	Extensions       []*AppExtension `json:"extensions"`
}

type AppActivate struct {
	Errors []*AppError `json:"errors"`
	App    *App        `json:"app"`
}

type AppCountableConnection struct {
	PageInfo   *PageInfo           `json:"pageInfo"`
	Edges      []*AppCountableEdge `json:"edges"`
	TotalCount *int32              `json:"totalCount"`
}

type AppCountableEdge struct {
	Node   *App   `json:"node"`
	Cursor string `json:"cursor"`
}

type AppCreate struct {
	AuthToken *string     `json:"authToken"`
	Errors    []*AppError `json:"errors"`
	App       *App        `json:"app"`
}

type AppDeactivate struct {
	Errors []*AppError `json:"errors"`
	App    *App        `json:"app"`
}

type AppDelete struct {
	Errors []*AppError `json:"errors"`
	App    *App        `json:"app"`
}

type AppDeleteFailedInstallation struct {
	Errors          []*AppError      `json:"errors"`
	AppInstallation *AppInstallation `json:"appInstallation"`
}

type AppError struct {
	Field       *string          `json:"field"`
	Message     *string          `json:"message"`
	Code        AppErrorCode     `json:"code"`
	Permissions []PermissionEnum `json:"permissions"`
}

type AppExtension struct {
	ID          string                 `json:"id"`
	App         *App                   `json:"app"`
	Label       string                 `json:"label"`
	URL         string                 `json:"url"`
	View        AppExtensionViewEnum   `json:"view"`
	Type        AppExtensionTypeEnum   `json:"type"`
	Target      AppExtensionTargetEnum `json:"target"`
	Permissions []*Permission          `json:"permissions"`
	AccessToken *string                `json:"accessToken"`
}

type AppExtensionCountableConnection struct {
	PageInfo   *PageInfo                    `json:"pageInfo"`
	Edges      []*AppExtensionCountableEdge `json:"edges"`
	TotalCount *int32                       `json:"totalCount"`
}

type AppExtensionCountableEdge struct {
	Node   *AppExtension `json:"node"`
	Cursor string        `json:"cursor"`
}

type AppExtensionFilterInput struct {
	View   *AppExtensionViewEnum   `json:"view"`
	Type   *AppExtensionTypeEnum   `json:"type"`
	Target *AppExtensionTargetEnum `json:"target"`
}

type AppFetchManifest struct {
	Manifest *Manifest   `json:"manifest"`
	Errors   []*AppError `json:"errors"`
}

type AppFilterInput struct {
	Search   *string      `json:"search"`
	IsActive *bool        `json:"isActive"`
	Type     *AppTypeEnum `json:"type"`
}

type AppInput struct {
	Name        *string           `json:"name"`
	Permissions []*PermissionEnum `json:"permissions"`
}

type AppInstall struct {
	Errors          []*AppError      `json:"errors"`
	AppInstallation *AppInstallation `json:"appInstallation"`
}

type AppInstallInput struct {
	AppName                   *string           `json:"appName"`
	ManifestURL               *string           `json:"manifestUrl"`
	ActivateAfterInstallation *bool             `json:"activateAfterInstallation"`
	Permissions               []*PermissionEnum `json:"permissions"`
}

type AppInstallation struct {
	AppName     string        `json:"appName"`
	ManifestURL string        `json:"manifestUrl"`
	ID          string        `json:"id"`
	Status      JobStatusEnum `json:"status"`
	CreatedAt   DateTime      `json:"createdAt"`
	UpdatedAt   DateTime      `json:"updatedAt"`
	Message     *string       `json:"message"`
}

type AppManifestExtension struct {
	Permissions []*Permission          `json:"permissions"`
	Label       string                 `json:"label"`
	URL         string                 `json:"url"`
	View        AppExtensionViewEnum   `json:"view"`
	Type        AppExtensionTypeEnum   `json:"type"`
	Target      AppExtensionTargetEnum `json:"target"`
}

type AppRetryInstall struct {
	Errors          []*AppError      `json:"errors"`
	AppInstallation *AppInstallation `json:"appInstallation"`
}

type AppSortingInput struct {
	Direction OrderDirection `json:"direction"`
	Field     AppSortField   `json:"field"`
}

type AppToken struct {
	Name      *string `json:"name"`
	AuthToken *string `json:"authToken"`
	ID        string  `json:"id"`
}

type AppTokenCreate struct {
	AuthToken *string     `json:"authToken"`
	Errors    []*AppError `json:"errors"`
	AppToken  *AppToken   `json:"appToken"`
}

type AppTokenDelete struct {
	Errors   []*AppError `json:"errors"`
	AppToken *AppToken   `json:"appToken"`
}

type AppTokenInput struct {
	Name *string `json:"name"`
	App  string  `json:"app"`
}

type AppTokenVerify struct {
	Valid  bool        `json:"valid"`
	Errors []*AppError `json:"errors"`
}

type AppUpdate struct {
	Errors []*AppError `json:"errors"`
	App    *App        `json:"app"`
}

type AssignNavigation struct {
	Menu   *Menu        `json:"menu"`
	Errors []*MenuError `json:"errors"`
}

type AttributeBulkDelete struct {
	Count  int32             `json:"count"`
	Errors []*AttributeError `json:"errors"`
}

type AttributeChoicesSortingInput struct {
	Direction OrderDirection            `json:"direction"`
	Field     AttributeChoicesSortField `json:"field"`
}

type AttributeCountableConnection struct {
	PageInfo   *PageInfo                 `json:"pageInfo"`
	Edges      []*AttributeCountableEdge `json:"edges"`
	TotalCount *int32                    `json:"totalCount"`
}

type AttributeCountableEdge struct {
	Node   *Attribute `json:"node"`
	Cursor string     `json:"cursor"`
}

type AttributeCreate struct {
	Attribute *Attribute        `json:"attribute"`
	Errors    []*AttributeError `json:"errors"`
}

type AttributeDelete struct {
	Errors    []*AttributeError `json:"errors"`
	Attribute *Attribute        `json:"attribute"`
}

type AttributeError struct {
	Field   *string            `json:"field"`
	Message *string            `json:"message"`
	Code    AttributeErrorCode `json:"code"`
}

type AttributeFilterInput struct {
	ValueRequired          *bool                `json:"valueRequired"`
	IsVariantOnly          *bool                `json:"isVariantOnly"`
	VisibleInStorefront    *bool                `json:"visibleInStorefront"`
	FilterableInStorefront *bool                `json:"filterableInStorefront"`
	FilterableInDashboard  *bool                `json:"filterableInDashboard"`
	AvailableInGrid        *bool                `json:"availableInGrid"`
	Metadata               []*MetadataInput     `json:"metadata"`
	Search                 *string              `json:"search"`
	Ids                    []string             `json:"ids"`
	Type                   *model.AttributeType `json:"type"`
	InCollection           *string              `json:"inCollection"`
	InCategory             *string              `json:"inCategory"`
	Channel                *string              `json:"channel"`
}

func (a *AttributeFilterInput) toSystemAttributeFilterOption() *model.AttributeFilterOption {
	res := &model.AttributeFilterOption{
		VisibleInStoreFront:    a.VisibleInStorefront,
		ValueRequired:          a.ValueRequired,
		IsVariantOnly:          a.IsVariantOnly,
		FilterableInStorefront: a.FilterableInStorefront,
		FilterableInDashboard:  a.FilterableInDashboard,
		AvailableInGrid:        a.AvailableInGrid,
		InCollection:           a.InCollection,
		InCategory:             a.InCategory,
		Channel:                a.Channel,
	}

	if len(a.Metadata) > 0 {
		res.Metadata = model.StringMAP{}
		for _, meta := range a.Metadata {
			if meta != nil && meta.Key != "" {
				res.Metadata[meta.Key] = meta.Value
			}
		}
	}
	if len(a.Ids) > 0 {
		res.Id = squirrel.Eq{store.AttributeTableName + ".Id": a.Ids}
	}
	if a.Type != nil && a.Type.IsValid() {
		res.Type = squirrel.Eq{store.AttributeTableName + ".Type": *a.Type}
	}

	return res
}

type AttributeInput struct {
	Slug        string              `json:"slug"`
	Values      []string            `json:"values"`
	ValuesRange *IntRangeInput      `json:"valuesRange"`
	DateTime    *DateTimeRangeInput `json:"dateTime"`
	Date        *DateRangeInput     `json:"date"`
	Boolean     *bool               `json:"boolean"`
}

type AttributeReorderValues struct {
	Attribute *Attribute        `json:"attribute"`
	Errors    []*AttributeError `json:"errors"`
}

type AttributeSortingInput struct {
	Direction OrderDirection     `json:"direction"`
	Field     AttributeSortField `json:"field"`
}

type AttributeTranslatableContent struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Translation *AttributeTranslation `json:"translation"`
}

type AttributeTranslate struct {
	Errors    []*TranslationError `json:"errors"`
	Attribute *Attribute          `json:"attribute"`
}

type AttributeTranslation struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Language *LanguageDisplay `json:"language"`
}

type AttributeUpdate struct {
	Attribute *Attribute        `json:"attribute"`
	Errors    []*AttributeError `json:"errors"`
}

type AttributeCreateInput struct {
	InputType                *model.AttributeInputType    `json:"inputType"`
	EntityType               *model.AttributeEntityType   `json:"entityType"`
	Name                     string                       `json:"name"`
	Slug                     *string                      `json:"slug"`
	Type                     model.AttributeType          `json:"type"`
	Unit                     *MeasurementUnitsEnum        `json:"unit"`
	Values                   []*AttributeValueCreateInput `json:"values"`
	ValueRequired            *bool                        `json:"valueRequired"`
	IsVariantOnly            *bool                        `json:"isVariantOnly"`
	VisibleInStorefront      *bool                        `json:"visibleInStorefront"`
	FilterableInStorefront   *bool                        `json:"filterableInStorefront"`
	FilterableInDashboard    *bool                        `json:"filterableInDashboard"`
	StorefrontSearchPosition *int32                       `json:"storefrontSearchPosition"`
	AvailableInGrid          *bool                        `json:"availableInGrid"`
}

func (a *AttributeCreateInput) getInputType() model.AttributeInputType {
	if a.InputType != nil {
		return *a.InputType
	}
	return model.AttributeInputType("")
}

func (a *AttributeCreateInput) getFieldValueByString(field string) any {
	switch field {
	case "filterable_in_storefront":
		return a.FilterableInStorefront
	case "filterable_in_dashboard":
		return a.FilterableInDashboard
	case "available_in_grid":
		return a.AvailableInGrid
	case "storefront_search_position":
		return a.StorefrontSearchPosition
	case "values":
		return a.Values
	case "input_type":
		return a.InputType
	default:
		return nil
	}
}

type AttributeUpdateInput struct {
	Name                     *string                      `json:"name"`
	Slug                     *string                      `json:"slug"`
	Unit                     *MeasurementUnitsEnum        `json:"unit"`
	RemoveValues             []string                     `json:"removeValues"`
	AddValues                []*AttributeValueUpdateInput `json:"addValues"`
	ValueRequired            *bool                        `json:"valueRequired"`
	IsVariantOnly            *bool                        `json:"isVariantOnly"`
	VisibleInStorefront      *bool                        `json:"visibleInStorefront"`
	FilterableInStorefront   *bool                        `json:"filterableInStorefront"`
	FilterableInDashboard    *bool                        `json:"filterableInDashboard"`
	StorefrontSearchPosition *int32                       `json:"storefrontSearchPosition"`
	AvailableInGrid          *bool                        `json:"availableInGrid"`
}

func (a *AttributeUpdateInput) getInputType() model.AttributeInputType {
	return model.AttributeInputType("")
}

func (i *AttributeUpdateInput) getFieldValueByString(name string) any {
	switch name {
	case "add_values":
		return i.AddValues
	default:
		return nil
	}
}

var (
	_ attributeValueInputIface = (*AttributeValueCreateInput)(nil)
	_ attributeValueInputIface = (*AttributeValueUpdateInput)(nil)
)

type AttributeValueCreateInput struct {
	Name        string     `json:"name"`
	Value       *string    `json:"value"`
	RichText    JSONString `json:"richText"`
	FileURL     *string    `json:"fileUrl"`
	ContentType *string    `json:"contentType"`
}

func (a *AttributeValueCreateInput) getName() string           { return a.Name }
func (a *AttributeValueCreateInput) getFileURL() *string       { return a.FileURL }
func (a *AttributeValueCreateInput) getContentType() *string   { return a.ContentType }
func (a *AttributeValueCreateInput) getValue() *string         { return a.Value }
func (a *AttributeValueCreateInput) getJsonString() JSONString { return a.RichText }

type AttributeValueUpdateInput struct {
	Value       *string    `json:"value"`
	RichText    JSONString `json:"richText"`
	FileURL     *string    `json:"fileUrl"`
	ContentType *string    `json:"contentType"`
	Name        string     `json:"name"`
}

func (a *AttributeValueUpdateInput) getName() string           { return a.Name }
func (a *AttributeValueUpdateInput) getFileURL() *string       { return a.FileURL }
func (a *AttributeValueUpdateInput) getContentType() *string   { return a.ContentType }
func (a *AttributeValueUpdateInput) getValue() *string         { return a.Value }
func (a *AttributeValueUpdateInput) getJsonString() JSONString { return a.RichText }

type AttributeValueBulkDelete struct {
	Count  int32             `json:"count"`
	Errors []*AttributeError `json:"errors"`
}

type AttributeValueCountableConnection struct {
	PageInfo   *PageInfo                      `json:"pageInfo"`
	Edges      []*AttributeValueCountableEdge `json:"edges"`
	TotalCount *int32                         `json:"totalCount"`
}

type AttributeValueCountableEdge struct {
	Node   *AttributeValue `json:"node"`
	Cursor string          `json:"cursor"`
}

type AttributeValueCreate struct {
	Attribute      *Attribute        `json:"attribute"`
	Errors         []*AttributeError `json:"errors"`
	AttributeValue *AttributeValue   `json:"attributeValue"`
}

type AttributeValueDelete struct {
	Attribute      *Attribute        `json:"attribute"`
	Errors         []*AttributeError `json:"errors"`
	AttributeValue *AttributeValue   `json:"attributeValue"`
}

type AttributeValueFilterInput struct {
	Search *string `json:"search"` // find attribute values with Name ILIKE %...% OR Slug ILIKE %...%
}

type AttributeValueInput struct {
	ID          *string    `json:"id"`
	Values      []string   `json:"values"`
	File        *string    `json:"file"`
	ContentType *string    `json:"contentType"`
	References  []string   `json:"references"`
	RichText    JSONString `json:"richText"`
	Boolean     *bool      `json:"boolean"`
	Date        *Date      `json:"date"`
	DateTime    *DateTime  `json:"dateTime"`
}

type AttributeValueTranslatableContent struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	RichText    JSONString                 `json:"richText"`
	Translation *AttributeValueTranslation `json:"translation"`
}

type AttributeValueTranslate struct {
	Errors         []*TranslationError `json:"errors"`
	AttributeValue *AttributeValue     `json:"attributeValue"`
}

type AttributeValueTranslation struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	RichText JSONString       `json:"richText"`
	Language *LanguageDisplay `json:"language"`
}

type AttributeValueTranslationInput struct {
	Name     *string    `json:"name"`
	RichText JSONString `json:"richText"`
}

type AttributeValueUpdate struct {
	Attribute      *Attribute        `json:"attribute"`
	Errors         []*AttributeError `json:"errors"`
	AttributeValue *AttributeValue   `json:"attributeValue"`
}

type BulkAttributeValueInput struct {
	ID      *string  `json:"id"`
	Values  []string `json:"values"`
	Boolean *bool    `json:"boolean"`
}

type BulkProductError struct {
	Field      *string          `json:"field"`
	Message    *string          `json:"message"`
	Code       ProductErrorCode `json:"code"`
	Attributes []string         `json:"attributes"`
	Values     []string         `json:"values"`
	Index      *int32           `json:"index"`
	Warehouses []string         `json:"warehouses"`
	Channels   []string         `json:"channels"`
}

type BulkStockError struct {
	Field      *string          `json:"field"`
	Message    *string          `json:"message"`
	Code       ProductErrorCode `json:"code"`
	Attributes []string         `json:"attributes"`
	Values     []string         `json:"values"`
	Index      *int32           `json:"index"`
}

type CatalogueInput struct {
	Products    []string `json:"products"`
	Categories  []string `json:"categories"`
	Collections []string `json:"collections"`
}

type CategoryBulkDelete struct {
	Count  int32           `json:"count"`
	Errors []*ProductError `json:"errors"`
}

type CategoryCountableConnection struct {
	PageInfo   *PageInfo                `json:"pageInfo"`
	Edges      []*CategoryCountableEdge `json:"edges"`
	TotalCount *int32                   `json:"totalCount"`
}

type CategoryCountableEdge struct {
	Node   *Category `json:"node"`
	Cursor string    `json:"cursor"`
}

type CategoryCreate struct {
	Errors   []*ProductError `json:"errors"`
	Category *Category       `json:"category"`
}

type CategoryDelete struct {
	Errors   []*ProductError `json:"errors"`
	Category *Category       `json:"category"`
}

type CategoryFilterInput struct {
	Search   *string          `json:"search"` // categories.Slug ILIKE ... OR categories.Name ILIKE ...
	Ids      []string         `json:"ids"`
	Metadata []*MetadataInput `json:"metadata"`
}

type CategoryInput struct {
	Description        JSONString `json:"description"`
	Name               string     `json:"name"`
	Slug               *string    `json:"slug"`
	Seo                *SeoInput  `json:"seo"`
	BackgroundImage    *string    `json:"backgroundImage"`
	BackgroundImageAlt *string    `json:"backgroundImageAlt"`
}

func (c *CategoryInput) Validate() *model.AppError {
	if c.Slug != nil && !slug.IsSlug(*c.Slug) {
		return model.NewAppError("CategoryInput.Validate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "slug"}, fmt.Sprintf("%s is not a slug", *c.Slug), http.StatusBadRequest)
	}
	return nil
}

// PatchCategory must be called after calling Validate()
func (c *CategoryInput) PatchCategory(category *model.Category) {
	category.Name = c.Name
	category.BackgroundImage = c.BackgroundImage
	if c.Description != nil {
		category.Description = model.StringInterface(c.Description)
	}
	if c.Slug != nil {
		category.Slug = *c.Slug
	}
	if c.Seo != nil {
		category.Seo = model.Seo{
			SeoTitle:       c.Seo.Title,
			SeoDescription: c.Seo.Description,
		}
	}
	if c.BackgroundImageAlt != nil {
		category.BackgroundImageAlt = *c.BackgroundImageAlt
	}
}

type CategorySortingInput struct {
	Direction OrderDirection    `json:"direction"`
	Field     CategorySortField `json:"field"`
	// Channel   *string           `json:"channel"`
}

type CategoryTranslatableContent struct {
	ID             string               `json:"id"`
	SeoTitle       *string              `json:"seoTitle"`
	SeoDescription *string              `json:"seoDescription"`
	Name           string               `json:"name"`
	Description    JSONString           `json:"description"`
	Translation    *CategoryTranslation `json:"translation"`
}

type CategoryTranslate struct {
	Errors   []*TranslationError `json:"errors"`
	Category *Category           `json:"category"`
}

type CategoryTranslation struct {
	ID             string           `json:"id"`
	SeoTitle       *string          `json:"seoTitle"`
	SeoDescription *string          `json:"seoDescription"`
	Name           *string          `json:"name"`
	Description    JSONString       `json:"description"`
	Language       *LanguageDisplay `json:"language"`
}

type CategoryUpdate struct {
	Errors   []*ProductError `json:"errors"`
	Category *Category       `json:"category"`
}

type ChannelActivate struct {
	Channel *Channel        `json:"channel"`
	Errors  []*ChannelError `json:"errors"`
}

type ChannelCreate struct {
	Errors  []*ChannelError `json:"errors"`
	Channel *Channel        `json:"channel"`
}

type ChannelCreateInput struct {
	IsActive         *bool       `json:"isActive"`
	Name             string      `json:"name"`
	Slug             string      `json:"slug"`
	CurrencyCode     string      `json:"currencyCode"`
	DefaultCountry   CountryCode `json:"defaultCountry"`
	AddShippingZones []string    `json:"addShippingZones"`
}

type ChannelDeactivate struct {
	Channel *Channel        `json:"channel"`
	Errors  []*ChannelError `json:"errors"`
}

type ChannelDelete struct {
	Errors  []*ChannelError `json:"errors"`
	Channel *Channel        `json:"channel"`
}

type ChannelDeleteInput struct {
	ChannelID string `json:"channelId"`
}

type ChannelError struct {
	Field         *string          `json:"field"`
	Message       *string          `json:"message"`
	Code          ChannelErrorCode `json:"code"`
	ShippingZones []string         `json:"shippingZones"`
}

type ChannelUpdate struct {
	Errors  []*ChannelError `json:"errors"`
	Channel *Channel        `json:"channel"`
}

type ChannelUpdateInput struct {
	IsActive            *bool        `json:"isActive"`
	Name                *string      `json:"name"`
	Slug                *string      `json:"slug"`
	DefaultCountry      *CountryCode `json:"defaultCountry"`
	AddShippingZones    []string     `json:"addShippingZones"`
	RemoveShippingZones []string     `json:"removeShippingZones"`
}

type CheckoutAddPromoCode struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutBillingAddressUpdate struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutComplete struct {
	Order              *Order           `json:"order"`
	ConfirmationNeeded bool             `json:"confirmationNeeded"`
	ConfirmationData   JSONString       `json:"confirmationData"`
	Errors             []*CheckoutError `json:"errors"`
}

type CheckoutCountableConnection struct {
	PageInfo   *PageInfo                `json:"pageInfo"`
	Edges      []*CheckoutCountableEdge `json:"edges"`
	TotalCount *int32                   `json:"totalCount"`
}

type CheckoutCountableEdge struct {
	Node   *Checkout `json:"node"`
	Cursor string    `json:"cursor"`
}

type CheckoutCreate struct {
	Errors   []*CheckoutError `json:"errors"`
	Checkout *Checkout        `json:"checkout"`
}

type CheckoutCreateInput struct {
	Channel         *string              `json:"channel"`
	Lines           []*CheckoutLineInput `json:"lines"`
	Email           *string              `json:"email"`
	ShippingAddress *AddressInput        `json:"shippingAddress"`
	BillingAddress  *AddressInput        `json:"billingAddress"`
	LanguageCode    *LanguageCodeEnum    `json:"languageCode"`
}

type CheckoutCustomerAttach struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutCustomerDetach struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutDeliveryMethodUpdate struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutEmailUpdate struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutError struct {
	Field       *string                `json:"field"`
	Message     *string                `json:"message"`
	Code        CheckoutErrorCode      `json:"code"`
	Variants    []string               `json:"variants"`
	AddressType *model.AddressTypeEnum `json:"addressType"`
}

type CheckoutLanguageCodeUpdate struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutLineCountableConnection struct {
	PageInfo   *PageInfo                    `json:"pageInfo"`
	Edges      []*CheckoutLineCountableEdge `json:"edges"`
	TotalCount *int32                       `json:"totalCount"`
}

type CheckoutLineCountableEdge struct {
	Node   *CheckoutLine `json:"node"`
	Cursor string        `json:"cursor"`
}

type CheckoutLineDelete struct {
	Checkout *Checkout `json:"checkout"`
	// Errors   []*CheckoutError `json:"errors"`
}

type CheckoutLineInput struct {
	Quantity  int32  `json:"quantity"`
	VariantID string `json:"variantId"`
}

type CheckoutLinesAdd struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutLinesUpdate struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutPaymentCreate struct {
	Checkout *Checkout       `json:"checkout"`
	Payment  *Payment        `json:"payment"`
	Errors   []*PaymentError `json:"errors"`
}

type CheckoutRemovePromoCode struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutShippingAddressUpdate struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type CheckoutShippingMethodUpdate struct {
	Checkout *Checkout        `json:"checkout"`
	Errors   []*CheckoutError `json:"errors"`
}

type ChoiceValue struct {
	Raw     *string `json:"raw"`
	Verbose *string `json:"verbose"`
}

type CollectionAddProducts struct {
	Collection *Collection        `json:"collection"`
	Errors     []*CollectionError `json:"errors"`
}

type CollectionBulkDelete struct {
	Count  int32              `json:"count"`
	Errors []*CollectionError `json:"errors"`
}

type CollectionChannelListingError struct {
	Field      *string          `json:"field"`
	Message    *string          `json:"message"`
	Code       ProductErrorCode `json:"code"`
	Attributes []string         `json:"attributes"`
	Values     []string         `json:"values"`
	Channels   []string         `json:"channels"`
}

type CollectionChannelListingUpdate struct {
	Collection *Collection                      `json:"collection"`
	Errors     []*CollectionChannelListingError `json:"errors"`
}

type CollectionChannelListingUpdateInput struct {
	AddChannels    []*PublishableChannelListingInput `json:"addChannels"`
	RemoveChannels []string                          `json:"removeChannels"`
}

type CollectionCountableConnection struct {
	PageInfo   *PageInfo                  `json:"pageInfo"`
	Edges      []*CollectionCountableEdge `json:"edges"`
	TotalCount *int32                     `json:"totalCount"`
}

type CollectionCountableEdge struct {
	Node   *Collection `json:"node"`
	Cursor string      `json:"cursor"`
}

type CollectionCreate struct {
	Errors     []*CollectionError `json:"errors"`
	Collection *Collection        `json:"collection"`
}

type CollectionCreateInput struct {
	IsPublished        *bool           `json:"isPublished"`
	Name               *string         `json:"name"`
	Slug               *string         `json:"slug"`
	Description        JSONString      `json:"description"`
	BackgroundImage    *graphql.Upload `json:"backgroundImage"`
	BackgroundImageAlt *string         `json:"backgroundImageAlt"`
	Seo                *SeoInput       `json:"seo"`
	PublicationDate    *Date           `json:"publicationDate"`
	Products           []string        `json:"products"`
}

type CollectionDelete struct {
	Errors     []*CollectionError `json:"errors"`
	Collection *Collection        `json:"collection"`
}

type CollectionError struct {
	Field    *string             `json:"field"`
	Message  *string             `json:"message"`
	Products []string            `json:"products"`
	Code     CollectionErrorCode `json:"code"`
}

type CollectionFilterInput struct {
	Published *CollectionPublished `json:"published"`
	Search    *string              `json:"search"`
	Metadata  []*MetadataInput     `json:"metadata"`
	Ids       []string             `json:"ids"`
	Channel   *string              `json:"channel"`
}

type CollectionInput struct {
	IsPublished        *bool           `json:"isPublished"`
	Name               *string         `json:"name"`
	Slug               *string         `json:"slug"`
	Description        JSONString      `json:"description"`
	BackgroundImage    *graphql.Upload `json:"backgroundImage"`
	BackgroundImageAlt *string         `json:"backgroundImageAlt"`
	Seo                *SeoInput       `json:"seo"`
	PublicationDate    *Date           `json:"publicationDate"`
}

type CollectionRemoveProducts struct {
	Collection *Collection        `json:"collection"`
	Errors     []*CollectionError `json:"errors"`
}

type CollectionReorderProducts struct {
	Collection *Collection        `json:"collection"`
	Errors     []*CollectionError `json:"errors"`
}

type CollectionSortingInput struct {
	Direction OrderDirection      `json:"direction"`
	Channel   *string             `json:"channel"`
	Field     CollectionSortField `json:"field"`
}

type CollectionTranslatableContent struct {
	ID             string                 `json:"id"`
	SeoTitle       *string                `json:"seoTitle"`
	SeoDescription *string                `json:"seoDescription"`
	Name           string                 `json:"name"`
	Description    JSONString             `json:"description"`
	Translation    *CollectionTranslation `json:"translation"`
}

type CollectionTranslate struct {
	Errors     []*TranslationError `json:"errors"`
	Collection *Collection         `json:"collection"`
}

type CollectionTranslation struct {
	ID             string           `json:"id"`
	SeoTitle       *string          `json:"seoTitle"`
	SeoDescription *string          `json:"seoDescription"`
	Name           *string          `json:"name"`
	Description    JSONString       `json:"description"`
	Language       *LanguageDisplay `json:"language"`
}

type CollectionUpdate struct {
	Errors     []*CollectionError `json:"errors"`
	Collection *Collection        `json:"collection"`
}

type ConfigurationItem struct {
	Name     string                      `json:"name"`
	Value    *string                     `json:"value"`
	Type     *ConfigurationTypeFieldEnum `json:"type"`
	HelpText *string                     `json:"helpText"`
	Label    *string                     `json:"label"`
}

type ConfigurationItemInput struct {
	Name  string  `json:"name"`
	Value *string `json:"value"`
}

type ConfirmAccount struct {
	User *User `json:"user"`
}

type ConfirmEmailChange struct {
	User *User `json:"user"`
}

type CountryDisplay struct {
	Code    string `json:"code"`
	Country string `json:"country"`
	Vat     *Vat   `json:"vat"`
}

type CreateToken struct {
	Token        *string `json:"token"`
	RefreshToken *string `json:"refreshToken"`
	CsrfToken    *string `json:"csrfToken"`
	User         *User   `json:"user"`
}

type CreditCard struct {
	Brand       string  `json:"brand"`
	FirstDigits *string `json:"firstDigits"`
	LastDigits  string  `json:"lastDigits"`
	ExpMonth    *int32  `json:"expMonth"`
	ExpYear     *int32  `json:"expYear"`
}

type CustomerBulkDelete struct {
	Count  int32           `json:"count"`
	Errors []*AccountError `json:"errors"`
}

type CustomerCreate struct {
	Errors []*AccountError `json:"errors"`
	User   *User           `json:"user"`
}

type CustomerDelete struct {
	Errors []*AccountError `json:"errors"`
	User   *User           `json:"user"`
}

type CustomerFilterInput struct {
	DateJoined     *DateRangeInput  `json:"dateJoined"`
	NumberOfOrders *IntRangeInput   `json:"numberOfOrders"`
	PlacedOrders   *DateRangeInput  `json:"placedOrders"`
	Search         *string          `json:"search"`
	Metadata       []*MetadataInput `json:"metadata"`
}

type CustomerInput struct {
	DefaultBillingAddress  *AddressInput     `json:"defaultBillingAddress"`
	DefaultShippingAddress *AddressInput     `json:"defaultShippingAddress"`
	FirstName              *string           `json:"firstName"`
	LastName               *string           `json:"lastName"`
	Email                  *string           `json:"email"`
	IsActive               *bool             `json:"isActive"`
	Note                   *string           `json:"note"`
	LanguageCode           *LanguageCodeEnum `json:"languageCode"`
}

type CustomerUpdate struct {
	Errors []*AccountError `json:"errors"`
	User   *User           `json:"user"`
}

type DateRangeInput struct {
	Gte *Date `json:"gte"`
	Lte *Date `json:"lte"`
}

type DateTimeRangeInput struct {
	Gte *DateTime `json:"gte"`
	Lte *DateTime `json:"lte"`
}

type DeactivateAllUserTokens struct {
	Ok bool `json:"ok"`
}

type ObjectWithMetadata struct {
	PrivateMetadata []*MetadataItem
	Metadata        []*MetadataItem
}

type DeleteMetadata struct {
	Errors []*MetadataError   `json:"errors"`
	Item   ObjectWithMetadata `json:"item"`
}

type DeletePrivateMetadata struct {
	Errors []*MetadataError   `json:"errors"`
	Item   ObjectWithMetadata `json:"item"`
}

type DigitalContentCountableConnection struct {
	PageInfo   *PageInfo                      `json:"pageInfo"`
	Edges      []*DigitalContentCountableEdge `json:"edges"`
	TotalCount *int32                         `json:"totalCount"`
}

type DigitalContentCountableEdge struct {
	Node   *DigitalContent `json:"node"`
	Cursor string          `json:"cursor"`
}

type DigitalContentCreate struct {
	Variant *ProductVariant `json:"variant"`
	Content *DigitalContent `json:"content"`
	Errors  []*ProductError `json:"errors"`
}

type DigitalContentDelete struct {
	Variant *ProductVariant `json:"variant"`
	Errors  []*ProductError `json:"errors"`
}

type DigitalContentInput struct {
	UseDefaultSettings   bool   `json:"useDefaultSettings"`
	MaxDownloads         *int32 `json:"maxDownloads"`
	URLValidDays         *int32 `json:"urlValidDays"`
	AutomaticFulfillment *bool  `json:"automaticFulfillment"`
}

type DigitalContentUpdate struct {
	Variant *ProductVariant `json:"variant"`
	Content *DigitalContent `json:"content"`
	Errors  []*ProductError `json:"errors"`
}

type DigitalContentUploadInput struct {
	UseDefaultSettings   bool           `json:"useDefaultSettings"`
	MaxDownloads         *int32         `json:"maxDownloads"`
	URLValidDays         *int32         `json:"urlValidDays"`
	AutomaticFulfillment *bool          `json:"automaticFulfillment"`
	ContentFile          graphql.Upload `json:"contentFile"`
}

type DigitalContentURLCreate struct {
	Errors            []*ProductError    `json:"errors"`
	DigitalContentURL *DigitalContentURL `json:"digitalContentUrl"`
}

type DigitalContentURLCreateInput struct {
	Content string `json:"content"`
}

type DiscountError struct {
	Field    *string           `json:"field"`
	Message  *string           `json:"message"`
	Products []string          `json:"products"`
	Code     DiscountErrorCode `json:"code"`
	Channels []string          `json:"channels"`
}

type Domain struct {
	Host       string `json:"host"`
	SslEnabled bool   `json:"sslEnabled"`
	URL        string `json:"url"`
}

type DraftOrderBulkDelete struct {
	Count  int32         `json:"count"`
	Errors []*OrderError `json:"errors"`
}

type DraftOrderComplete struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type DraftOrderCreate struct {
	Errors []*OrderError `json:"errors"`
	Order  *Order        `json:"order"`
}

type DraftOrderCreateInput struct {
	BillingAddress  *AddressInput           `json:"billingAddress"`
	User            *string                 `json:"user"`
	UserEmail       *string                 `json:"userEmail"`
	Discount        *PositiveDecimal        `json:"discount"`
	ShippingAddress *AddressInput           `json:"shippingAddress"`
	ShippingMethod  *string                 `json:"shippingMethod"`
	Voucher         *string                 `json:"voucher"`
	CustomerNote    *string                 `json:"customerNote"`
	ChannelID       *string                 `json:"channelId"`
	RedirectURL     *string                 `json:"redirectUrl"`
	Lines           []*OrderLineCreateInput `json:"lines"`
}

type DraftOrderDelete struct {
	Errors []*OrderError `json:"errors"`
	Order  *Order        `json:"order"`
}

type DraftOrderInput struct {
	BillingAddress  *AddressInput    `json:"billingAddress"`
	User            *string          `json:"user"`
	UserEmail       *string          `json:"userEmail"`
	Discount        *PositiveDecimal `json:"discount"`
	ShippingAddress *AddressInput    `json:"shippingAddress"`
	ShippingMethod  *string          `json:"shippingMethod"`
	Voucher         *string          `json:"voucher"`
	CustomerNote    *string          `json:"customerNote"`
	ChannelID       *string          `json:"channelId"`
	RedirectURL     *string          `json:"redirectUrl"`
}

type DraftOrderLinesBulkDelete struct {
	Count  int32         `json:"count"`
	Errors []*OrderError `json:"errors"`
}

type DraftOrderUpdate struct {
	Errors []*OrderError `json:"errors"`
	Order  *Order        `json:"order"`
}

type ExportEvent struct {
	ID      string           `json:"id"`
	Date    DateTime         `json:"date"`
	Type    ExportEventsEnum `json:"type"`
	User    *User            `json:"user"`
	Message string           `json:"message"`
}

type ExportFile struct {
	ID        string         `json:"id"`
	User      *User          `json:"user"`
	Status    JobStatusEnum  `json:"status"`
	CreatedAt DateTime       `json:"createdAt"`
	UpdatedAt DateTime       `json:"updatedAt"`
	Message   *string        `json:"message"`
	URL       *string        `json:"url"`
	Events    []*ExportEvent `json:"events"`
}

type ExportFileCountableConnection struct {
	PageInfo   *PageInfo                  `json:"pageInfo"`
	Edges      []*ExportFileCountableEdge `json:"edges"`
	TotalCount *int32                     `json:"totalCount"`
}

type ExportFileCountableEdge struct {
	Node   *ExportFile `json:"node"`
	Cursor string      `json:"cursor"`
}

type ExportFileFilterInput struct {
	CreatedAt *DateTimeRangeInput `json:"createdAt"`
	UpdatedAt *DateTimeRangeInput `json:"updatedAt"`
	Status    *JobStatusEnum      `json:"status"`
	User      *string             `json:"user"`
}

type ExportFileSortingInput struct {
	Direction OrderDirection      `json:"direction"`
	Field     ExportFileSortField `json:"field"`
}

type ExportInfoInput struct {
	Attributes []string           `json:"attributes"`
	Warehouses []string           `json:"warehouses"`
	Channels   []string           `json:"channels"`
	Fields     []ProductFieldEnum `json:"fields"`
}

type ExportProducts struct {
	ExportFile *ExportFile `json:"exportFile"`
}

type ExportProductsInput struct {
	Scope      ExportScope         `json:"scope"`
	Filter     *ProductFilterInput `json:"filter"`
	Ids        []string            `json:"ids"`
	ExportInfo *ExportInfoInput    `json:"exportInfo"`
	FileType   FileTypesEnum       `json:"fileType"`
}

type ExternalAuthentication struct {
	ID   string  `json:"id"`
	Name *string `json:"name"`
}

type ExternalAuthenticationURL struct {
	AuthenticationData JSONString `json:"authenticationData"`
}

type ExternalLogout struct {
	LogoutData JSONString `json:"logoutData"`
}

type ExternalNotificationError struct {
	Field   *string                        `json:"field"`
	Message *string                        `json:"message"`
	Code    ExternalNotificationErrorCodes `json:"code"`
}

type ExternalNotificationTrigger struct {
	Errors []*ExternalNotificationError `json:"errors"`
}

type ExternalNotificationTriggerInput struct {
	Ids               []string   `json:"ids"`
	ExtraPayload      JSONString `json:"extraPayload"`
	ExternalEventType string     `json:"externalEventType"`
}

type ExternalObtainAccessTokens struct {
	Token        *string `json:"token"`
	RefreshToken *string `json:"refreshToken"`
	CsrfToken    *string `json:"csrfToken"`
	User         *User   `json:"user"`
}

type ExternalRefresh struct {
	Token        *string `json:"token"`
	RefreshToken *string `json:"refreshToken"`
	CsrfToken    *string `json:"csrfToken"`
	User         *User   `json:"user"`
}

type ExternalVerify struct {
	User       *User      `json:"user"`
	IsValid    bool       `json:"isValid"`
	VerifyData JSONString `json:"verifyData"`
}

type File struct {
	URL         string  `json:"url"`
	ContentType *string `json:"contentType"`
}

type FileUpload struct {
	UploadedFile *File          `json:"uploadedFile"`
	Errors       []*UploadError `json:"errors"`
}

type FulfillmentApprove struct {
	Fulfillment *Fulfillment  `json:"fulfillment"`
	Order       *Order        `json:"order"`
	Errors      []*OrderError `json:"errors"`
}

type FulfillmentCancel struct {
	Fulfillment *Fulfillment  `json:"fulfillment"`
	Order       *Order        `json:"order"`
	Errors      []*OrderError `json:"errors"`
}

type FulfillmentCancelInput struct {
	WarehouseID *string `json:"warehouseId"`
}

type FulfillmentRefundProducts struct {
	Fulfillment *Fulfillment  `json:"fulfillment"`
	Order       *Order        `json:"order"`
	Errors      []*OrderError `json:"errors"`
}

type FulfillmentReturnProducts struct {
	ReturnFulfillment  *Fulfillment  `json:"returnFulfillment"`
	ReplaceFulfillment *Fulfillment  `json:"replaceFulfillment"`
	Order              *Order        `json:"order"`
	ReplaceOrder       *Order        `json:"replaceOrder"`
	Errors             []*OrderError `json:"errors"`
}

type FulfillmentUpdateTracking struct {
	Fulfillment *Fulfillment  `json:"fulfillment"`
	Order       *Order        `json:"order"`
	Errors      []*OrderError `json:"errors"`
}

type FulfillmentUpdateTrackingInput struct {
	TrackingNumber *string `json:"trackingNumber"`
	NotifyCustomer *bool   `json:"notifyCustomer"`
}

type GatewayConfigLine struct {
	Field string  `json:"field"`
	Value *string `json:"value"`
}

type GiftCardActivate struct {
	GiftCard *GiftCard        `json:"giftCard"`
	Errors   []*GiftCardError `json:"errors"`
}

type GiftCardAddNote struct {
	GiftCard *GiftCard        `json:"giftCard"`
	Event    *GiftCardEvent   `json:"event"`
	Errors   []*GiftCardError `json:"errors"`
}

type GiftCardAddNoteInput struct {
	Message string `json:"message"`
}

type GiftCardBulkActivate struct {
	Count  int32            `json:"count"`
	Errors []*GiftCardError `json:"errors"`
}

type GiftCardBulkDeactivate struct {
	Count  int32            `json:"count"`
	Errors []*GiftCardError `json:"errors"`
}

type GiftCardBulkDelete struct {
	Count  int32            `json:"count"`
	Errors []*GiftCardError `json:"errors"`
}

type GiftCardCountableConnection struct {
	PageInfo   *PageInfo                `json:"pageInfo"`
	Edges      []*GiftCardCountableEdge `json:"edges"`
	TotalCount *int32                   `json:"totalCount"`
}

type GiftCardCountableEdge struct {
	Node   *GiftCard `json:"node"`
	Cursor string    `json:"cursor"`
}

type GiftCardCreate struct {
	Errors   []*GiftCardError `json:"errors"`
	GiftCard *GiftCard        `json:"giftCard"`
}

type GiftCardCreateInput struct {
	Tag        *string     `json:"tag"`
	ExpiryDate *Date       `json:"expiryDate"`
	StartDate  *Date       `json:"startDate"`
	EndDate    *Date       `json:"endDate"`
	Balance    *PriceInput `json:"balance"`
	UserEmail  *string     `json:"userEmail"`
	Channel    *string     `json:"channel"`
	IsActive   bool        `json:"isActive"`
	Code       *string     `json:"code"`
	Note       *string     `json:"note"`
}

type GiftCardDeactivate struct {
	GiftCard *GiftCard        `json:"giftCard"`
	Errors   []*GiftCardError `json:"errors"`
}

type GiftCardDelete struct {
	Errors   []*GiftCardError `json:"errors"`
	GiftCard *GiftCard        `json:"giftCard"`
}

type GiftCardError struct {
	Field   *string           `json:"field"`
	Message *string           `json:"message"`
	Code    GiftCardErrorCode `json:"code"`
}

type GiftCardEventBalance struct {
	InitialBalance    *Money `json:"initialBalance"`
	CurrentBalance    *Money `json:"currentBalance"`
	OldInitialBalance *Money `json:"oldInitialBalance"`
	OldCurrentBalance *Money `json:"oldCurrentBalance"`
}

type GiftCardResend struct {
	GiftCard *GiftCard        `json:"giftCard"`
	Errors   []*GiftCardError `json:"errors"`
}

type GiftCardResendInput struct {
	ID      string  `json:"id"`
	Email   *string `json:"email"`
	Channel string  `json:"channel"`
}

type GiftCardSettings struct {
	ExpiryType   GiftCardSettingsExpiryTypeEnum `json:"expiryType"`
	ExpiryPeriod *TimePeriod                    `json:"expiryPeriod"`
}

type GiftCardSettingsError struct {
	Field   *string                   `json:"field"`
	Message *string                   `json:"message"`
	Code    GiftCardSettingsErrorCode `json:"code"`
}

type GiftCardSettingsUpdate struct {
	GiftCardSettings *GiftCardSettings        `json:"giftCardSettings"`
	Errors           []*GiftCardSettingsError `json:"errors"`
}

type GiftCardSettingsUpdateInput struct {
	ExpiryType   *GiftCardSettingsExpiryTypeEnum `json:"expiryType"`
	ExpiryPeriod *TimePeriodInputType            `json:"expiryPeriod"`
}

type GiftCardSortingInput struct {
	Direction OrderDirection    `json:"direction"`
	Field     GiftCardSortField `json:"field"`
}

type GiftCardUpdate struct {
	Errors   []*GiftCardError `json:"errors"`
	GiftCard *GiftCard        `json:"giftCard"`
}

type GiftCardUpdateInput struct {
	Tag           *string          `json:"tag"`
	ExpiryDate    *Date            `json:"expiryDate"`
	BalanceAmount *PositiveDecimal `json:"balanceAmount"`

	// StartDate     *Date            `json:"startDate"`
	// EndDate       *Date            `json:"endDate"`
}

type Group struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Permissions   []*Permission `json:"permissions"`
	Users         []*User       `json:"users"`
	UserCanManage bool          `json:"userCanManage"`
}

type GroupCountableConnection struct {
	PageInfo   *PageInfo             `json:"pageInfo"`
	Edges      []*GroupCountableEdge `json:"edges"`
	TotalCount *int32                `json:"totalCount"`
}

type GroupCountableEdge struct {
	Node   *Group `json:"node"`
	Cursor string `json:"cursor"`
}

type Image struct {
	URL string  `json:"url"`
	Alt *string `json:"alt"`
}

type IntRangeInput struct {
	Gte *int32 `json:"gte"`
	Lte *int32 `json:"lte"`
}

type InvoiceCreate struct {
	Errors  []*InvoiceError `json:"errors"`
	Invoice *Invoice        `json:"invoice"`
}

type InvoiceCreateInput struct {
	Number string `json:"number"`
	URL    string `json:"url"`
}

type InvoiceDelete struct {
	Errors  []*InvoiceError `json:"errors"`
	Invoice *Invoice        `json:"invoice"`
}

type InvoiceError struct {
	Field   *string          `json:"field"`
	Message *string          `json:"message"`
	Code    InvoiceErrorCode `json:"code"`
}

type InvoiceRequest struct {
	Order   *Order          `json:"order"`
	Errors  []*InvoiceError `json:"errors"`
	Invoice *Invoice        `json:"invoice"`
}

type InvoiceRequestDelete struct {
	Errors  []*InvoiceError `json:"errors"`
	Invoice *Invoice        `json:"invoice"`
}

type InvoiceSendNotification struct {
	Errors  []*InvoiceError `json:"errors"`
	Invoice *Invoice        `json:"invoice"`
}

type InvoiceUpdate struct {
	Errors  []*InvoiceError `json:"errors"`
	Invoice *Invoice        `json:"invoice"`
}

type LanguageDisplay struct {
	Code     LanguageCodeEnum `json:"code"`
	Language string           `json:"language"`
}

type LimitInfo struct {
	CurrentUsage *Limits `json:"currentUsage"`
	AllowedUsage *Limits `json:"allowedUsage"`
}

type Limits struct {
	Channels        *int32 `json:"channels"`
	Orders          *int32 `json:"orders"`
	ProductVariants *int32 `json:"productVariants"`
	StaffUsers      *int32 `json:"staffUsers"`
	Warehouses      *int32 `json:"warehouses"`
}

type LoginError struct {
	Field   *string        `json:"field"`
	Message *string        `json:"message"`
	Code    LoginErrorCode `json:"code"`
}

type LoginInput struct {
	ID       string `json:"id"`
	LoginID  string `json:"loginId"`
	Password string `json:"password"`
	MfaToken string `json:"mfaToken"`
	DeviceID string `json:"deviceId"`
	LdapOnly bool   `json:"ldapOnly"`
}

type LoginResponse struct {
	Error *LoginError `json:"error"`
	User  *User       `json:"user"`
}

type Manifest struct {
	Identifier       string                  `json:"identifier"`
	Version          string                  `json:"version"`
	Name             string                  `json:"name"`
	About            *string                 `json:"about"`
	Permissions      []*Permission           `json:"permissions"`
	AppURL           *string                 `json:"appUrl"`
	ConfigurationURL *string                 `json:"configurationUrl"`
	TokenTargetURL   *string                 `json:"tokenTargetUrl"`
	DataPrivacy      *string                 `json:"dataPrivacy"`
	DataPrivacyURL   *string                 `json:"dataPrivacyUrl"`
	HomepageURL      *string                 `json:"homepageUrl"`
	SupportURL       *string                 `json:"supportUrl"`
	Extensions       []*AppManifestExtension `json:"extensions"`
}

type Margin struct {
	Start *int32 `json:"start"`
	Stop  *int32 `json:"stop"`
}

type MenuBulkDelete struct {
	Count  int32        `json:"count"`
	Errors []*MenuError `json:"errors"`
}

type MenuCountableConnection struct {
	PageInfo   *PageInfo            `json:"pageInfo"`
	Edges      []*MenuCountableEdge `json:"edges"`
	TotalCount *int32               `json:"totalCount"`
}

type MenuCountableEdge struct {
	Node   *Menu  `json:"node"`
	Cursor string `json:"cursor"`
}

type MenuCreate struct {
	Errors []*MenuError `json:"errors"`
	Menu   *Menu        `json:"menu"`
}

type MenuCreateInput struct {
	Name  string           `json:"name"`
	Slug  *string          `json:"slug"`
	Items []*MenuItemInput `json:"items"`
}

type MenuDelete struct {
	Errors []*MenuError `json:"errors"`
	Menu   *Menu        `json:"menu"`
}

type MenuError struct {
	Field   *string       `json:"field"`
	Message *string       `json:"message"`
	Code    MenuErrorCode `json:"code"`
}

type MenuFilterInput struct {
	Search   *string          `json:"search"`
	Slug     []string         `json:"slug"`
	Metadata []*MetadataInput `json:"metadata"`
}

type MenuInput struct {
	Name *string `json:"name"`
	Slug *string `json:"slug"`
}

type MenuItemBulkDelete struct {
	Count  int32        `json:"count"`
	Errors []*MenuError `json:"errors"`
}

type MenuItemCountableConnection struct {
	PageInfo   *PageInfo                `json:"pageInfo"`
	Edges      []*MenuItemCountableEdge `json:"edges"`
	TotalCount *int32                   `json:"totalCount"`
}

type MenuItemCountableEdge struct {
	Node   *MenuItem `json:"node"`
	Cursor string    `json:"cursor"`
}

type MenuItemCreate struct {
	Errors   []*MenuError `json:"errors"`
	MenuItem *MenuItem    `json:"menuItem"`
}

type MenuItemCreateInput struct {
	Name       string  `json:"name"`
	URL        *string `json:"url"`
	Category   *string `json:"category"`
	Collection *string `json:"collection"`
	Page       *string `json:"page"`
	Menu       string  `json:"menu"`
	Parent     *string `json:"parent"`
}

type MenuItemDelete struct {
	Errors   []*MenuError `json:"errors"`
	MenuItem *MenuItem    `json:"menuItem"`
}

type MenuItemFilterInput struct {
	Search   *string          `json:"search"`
	Metadata []*MetadataInput `json:"metadata"`
}

type MenuItemInput struct {
	Name       *string `json:"name"`
	URL        *string `json:"url"`
	Category   *string `json:"category"`
	Collection *string `json:"collection"`
	Page       *string `json:"page"`
}

type MenuItemMove struct {
	Menu   *Menu        `json:"menu"`
	Errors []*MenuError `json:"errors"`
}

type MenuItemMoveInput struct {
	ItemID    string  `json:"itemId"`
	ParentID  *string `json:"parentId"`
	SortOrder *int32  `json:"sortOrder"`
}

type MenuItemSortingInput struct {
	Direction OrderDirection     `json:"direction"`
	Field     MenuItemsSortField `json:"field"`
}

type MenuItemTranslatableContent struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Translation *MenuItemTranslation `json:"translation"`
}

type MenuItemTranslate struct {
	Errors   []*TranslationError `json:"errors"`
	MenuItem *MenuItem           `json:"menuItem"`
}

type MenuItemTranslation struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Language *LanguageDisplay `json:"language"`
}

type MenuItemUpdate struct {
	Errors   []*MenuError `json:"errors"`
	MenuItem *MenuItem    `json:"menuItem"`
}

type MenuSortingInput struct {
	Direction OrderDirection `json:"direction"`
	Field     MenuSortField  `json:"field"`
}

type MenuUpdate struct {
	Errors []*MenuError `json:"errors"`
	Menu   *Menu        `json:"menu"`
}

type MetadataError struct {
	Field   *string           `json:"field"`
	Message *string           `json:"message"`
	Code    MetadataErrorCode `json:"code"`
}

type MetadataFilter struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
}

type MetadataInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MetadataItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Money struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type MoneyRange struct {
	Start *Money `json:"start"`
	Stop  *Money `json:"stop"`
}

type MoveProductInput struct {
	ProductID string `json:"productId"`
	SortOrder *int32 `json:"sortOrder"`
}

type NameTranslationInput struct {
	Name *string `json:"name"`
}

type OrderAddNote struct {
	Order  *Order        `json:"order"`
	Event  *OrderEvent   `json:"event"`
	Errors []*OrderError `json:"errors"`
}

type OrderAddNoteInput struct {
	Message string `json:"message"`
}

type OrderBulkCancel struct {
	Count  int32         `json:"count"`
	Errors []*OrderError `json:"errors"`
}

type OrderCancel struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type OrderCapture struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type OrderConfirm struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type OrderCountableConnection struct {
	PageInfo   *PageInfo             `json:"pageInfo"`
	Edges      []*OrderCountableEdge `json:"edges"`
	TotalCount *int32                `json:"totalCount"`
}

type OrderCountableEdge struct {
	Node   *Order `json:"node"`
	Cursor string `json:"cursor"` // string format of order's createAt
}

type OrderDiscountAdd struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type OrderDiscountCommonInput struct {
	ValueType DiscountValueTypeEnum `json:"valueType"`
	Value     PositiveDecimal       `json:"value"`
	Reason    *string               `json:"reason"`
}

type OrderDiscountDelete struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type OrderDiscountUpdate struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type OrderDraftFilterInput struct {
	Customer *string          `json:"customer"`
	Created  *DateRangeInput  `json:"created"`
	Search   *string          `json:"search"`
	Metadata []*MetadataInput `json:"metadata"`
	Channels []string         `json:"channels"`
}

type OrderError struct {
	Field       *string                `json:"field"`
	Message     *string                `json:"message"`
	Code        OrderErrorCode         `json:"code"`
	Warehouse   *string                `json:"warehouse"`
	OrderLine   *string                `json:"orderLine"`
	Variants    []string               `json:"variants"`
	AddressType *model.AddressTypeEnum `json:"addressType"`
}

type OrderEventCountableConnection struct {
	PageInfo   *PageInfo                  `json:"pageInfo"`
	Edges      []*OrderEventCountableEdge `json:"edges"`
	TotalCount *int32                     `json:"totalCount"`
}

type OrderEventCountableEdge struct {
	Node   *OrderEvent `json:"node"`
	Cursor string      `json:"cursor"`
}

type OrderEventDiscountObject struct {
	ValueType    DiscountValueTypeEnum  `json:"valueType"`
	Value        PositiveDecimal        `json:"value"`
	Reason       *string                `json:"reason"`
	Amount       *Money                 `json:"amount"`
	OldValueType *DiscountValueTypeEnum `json:"oldValueType"`
	OldValue     *PositiveDecimal       `json:"oldValue"`
	OldAmount    *Money                 `json:"oldAmount"`
}

type OrderEventOrderLineObject struct {
	Quantity  *int32                    `json:"quantity"`
	OrderLine *OrderLine                `json:"orderLine"`
	ItemName  *string                   `json:"itemName"`
	Discount  *OrderEventDiscountObject `json:"discount"`
}

type OrderFilterInput struct {
	PaymentStatus []*PaymentChargeStatusEnum `json:"paymentStatus"`
	Status        []*OrderStatusFilter       `json:"status"`
	Customer      *string                    `json:"customer"`
	Created       *DateRangeInput            `json:"created"`
	Search        *string                    `json:"search"`
	Metadata      []*MetadataInput           `json:"metadata"`
	Channels      []string                   `json:"channels"`
}

type OrderFulfill struct {
	Fulfillments []*Fulfillment `json:"fulfillments"`
	Order        *Order         `json:"order"`
	Errors       []*OrderError  `json:"errors"`
}

type OrderFulfillInput struct {
	Lines                  []*OrderFulfillLineInput `json:"lines"`
	NotifyCustomer         *bool                    `json:"notifyCustomer"`
	AllowStockToBeExceeded *bool                    `json:"allowStockToBeExceeded"`
}

type OrderFulfillLineInput struct {
	OrderLineID *string                   `json:"orderLineId"`
	Stocks      []*OrderFulfillStockInput `json:"stocks"`
}

type OrderFulfillStockInput struct {
	Quantity  int32  `json:"quantity"`
	Warehouse string `json:"warehouse"`
}

type OrderLineCreateInput struct {
	Quantity  int32  `json:"quantity"`
	VariantID string `json:"variantId"`
}

type OrderLineDelete struct {
	Order     *Order        `json:"order"`
	OrderLine *OrderLine    `json:"orderLine"`
	Errors    []*OrderError `json:"errors"`
}

type OrderLineDiscountRemove struct {
	OrderLine *OrderLine    `json:"orderLine"`
	Order     *Order        `json:"order"`
	Errors    []*OrderError `json:"errors"`
}

type OrderLineDiscountUpdate struct {
	OrderLine *OrderLine    `json:"orderLine"`
	Order     *Order        `json:"order"`
	Errors    []*OrderError `json:"errors"`
}

type OrderLineInput struct {
	Quantity int32 `json:"quantity"`
}

type OrderLineUpdate struct {
	Order     *Order        `json:"order"`
	Errors    []*OrderError `json:"errors"`
	OrderLine *OrderLine    `json:"orderLine"`
}

type OrderLinesCreate struct {
	Order      *Order        `json:"order"`
	OrderLines []*OrderLine  `json:"orderLines"`
	Errors     []*OrderError `json:"errors"`
}

type OrderMarkAsPaid struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type OrderRefund struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type OrderRefundFulfillmentLineInput struct {
	FulfillmentLineID string `json:"fulfillmentLineId"`
	Quantity          int32  `json:"quantity"`
}

type OrderRefundLineInput struct {
	OrderLineID string `json:"orderLineId"`
	Quantity    int32  `json:"quantity"`
}

type OrderRefundProductsInput struct {
	OrderLines           []*OrderRefundLineInput            `json:"orderLines"`
	FulfillmentLines     []*OrderRefundFulfillmentLineInput `json:"fulfillmentLines"`
	AmountToRefund       *PositiveDecimal                   `json:"amountToRefund"`
	IncludeShippingCosts *bool                              `json:"includeShippingCosts"`
}

type OrderReturnFulfillmentLineInput struct {
	FulfillmentLineID string `json:"fulfillmentLineId"`
	Quantity          int32  `json:"quantity"`
	Replace           *bool  `json:"replace"`
}

type OrderReturnLineInput struct {
	OrderLineID string `json:"orderLineId"`
	Quantity    int32  `json:"quantity"`
	Replace     *bool  `json:"replace"`
}

type OrderReturnProductsInput struct {
	OrderLines           []*OrderReturnLineInput            `json:"orderLines"`
	FulfillmentLines     []*OrderReturnFulfillmentLineInput `json:"fulfillmentLines"`
	AmountToRefund       *PositiveDecimal                   `json:"amountToRefund"`
	IncludeShippingCosts *bool                              `json:"includeShippingCosts"`
	Refund               *bool                              `json:"refund"`
}

type OrderSettings struct {
	AutomaticallyConfirmAllNewOrders         bool `json:"automaticallyConfirmAllNewOrders"`
	AutomaticallyFulfillNonShippableGiftCard bool `json:"automaticallyFulfillNonShippableGiftCard"`
}

type OrderSettingsError struct {
	Field   *string                `json:"field"`
	Message *string                `json:"message"`
	Code    OrderSettingsErrorCode `json:"code"`
}

type OrderSettingsUpdate struct {
	OrderSettings *OrderSettings        `json:"orderSettings"`
	Errors        []*OrderSettingsError `json:"errors"`
}

type OrderSettingsUpdateInput struct {
	AutomaticallyConfirmAllNewOrders         *bool `json:"automaticallyConfirmAllNewOrders"`
	AutomaticallyFulfillNonShippableGiftCard *bool `json:"automaticallyFulfillNonShippableGiftCard"`
}

type OrderSortingInput struct {
	Direction OrderDirection `json:"direction"`
	Field     OrderSortField `json:"field"`
}

type OrderUpdate struct {
	Errors []*OrderError `json:"errors"`
	Order  *Order        `json:"order"`
}

type OrderUpdateInput struct {
	BillingAddress  *AddressInput `json:"billingAddress"`
	UserEmail       *string       `json:"userEmail"`
	ShippingAddress *AddressInput `json:"shippingAddress"`
}

type OrderUpdateShipping struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type OrderUpdateShippingInput struct {
	ShippingMethod *string `json:"shippingMethod"`
}

type OrderVoid struct {
	Order  *Order        `json:"order"`
	Errors []*OrderError `json:"errors"`
}

type PageAttributeAssign struct {
	PageType *PageType    `json:"pageType"`
	Errors   []*PageError `json:"errors"`
}

type PageAttributeUnassign struct {
	PageType *PageType    `json:"pageType"`
	Errors   []*PageError `json:"errors"`
}

type PageBulkDelete struct {
	Count  int32        `json:"count"`
	Errors []*PageError `json:"errors"`
}

type PageBulkPublish struct {
	Count  int32        `json:"count"`
	Errors []*PageError `json:"errors"`
}

type PageCountableConnection struct {
	PageInfo   *PageInfo            `json:"pageInfo"`
	Edges      []*PageCountableEdge `json:"edges"`
	TotalCount *int32               `json:"totalCount"`
}

type PageCountableEdge struct {
	Node   *Page  `json:"node"`
	Cursor string `json:"cursor"`
}

type PageCreate struct {
	Errors []*PageError `json:"errors"`
	Page   *Page        `json:"page"`
}

type PageCreateInput struct {
	Slug            *string                `json:"slug"`
	Title           *string                `json:"title"`
	Content         JSONString             `json:"content"`
	Attributes      []*AttributeValueInput `json:"attributes"`
	IsPublished     *bool                  `json:"isPublished"`
	PublicationDate *string                `json:"publicationDate"`
	Seo             *SeoInput              `json:"seo"`
	PageType        string                 `json:"pageType"`
}

type PageDelete struct {
	Errors []*PageError `json:"errors"`
	Page   *Page        `json:"page"`
}

type PageError struct {
	Field      *string       `json:"field"`
	Message    *string       `json:"message"`
	Code       PageErrorCode `json:"code"`
	Attributes []string      `json:"attributes"`
	Values     []string      `json:"values"`
}

type PageFilterInput struct {
	Search    *string          `json:"search"`
	Metadata  []*MetadataInput `json:"metadata"`
	PageTypes []string         `json:"pageTypes"`
	Ids       []string         `json:"ids"`
}

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}

type PageInput struct {
	Slug            *string                `json:"slug"`
	Title           *string                `json:"title"`
	Content         JSONString             `json:"content"`
	Attributes      []*AttributeValueInput `json:"attributes"`
	IsPublished     *bool                  `json:"isPublished"`
	PublicationDate *string                `json:"publicationDate"`
	Seo             *SeoInput              `json:"seo"`
}

type PageReorderAttributeValues struct {
	Page   *Page        `json:"page"`
	Errors []*PageError `json:"errors"`
}

type PageSortingInput struct {
	Direction OrderDirection `json:"direction"`
	Field     PageSortField  `json:"field"`
}

type PageTranslatableContent struct {
	ID             string           `json:"id"`
	SeoTitle       *string          `json:"seoTitle"`
	SeoDescription *string          `json:"seoDescription"`
	Title          string           `json:"title"`
	Content        JSONString       `json:"content"`
	Translation    *PageTranslation `json:"translation"`
}

type PageTranslate struct {
	Errors []*TranslationError      `json:"errors"`
	Page   *PageTranslatableContent `json:"page"`
}

type PageTranslation struct {
	ID             string           `json:"id"`
	SeoTitle       *string          `json:"seoTitle"`
	SeoDescription *string          `json:"seoDescription"`
	Title          *string          `json:"title"`
	Content        JSONString       `json:"content"`
	Language       *LanguageDisplay `json:"language"`
}

type PageTranslationInput struct {
	SeoTitle       *string    `json:"seoTitle"`
	SeoDescription *string    `json:"seoDescription"`
	Title          *string    `json:"title"`
	Content        JSONString `json:"content"`
}

type PageType struct {
	ID                  string                        `json:"id"`
	Name                string                        `json:"name"`
	Slug                string                        `json:"slug"`
	PrivateMetadata     []*MetadataItem               `json:"privateMetadata"`
	Metadata            []*MetadataItem               `json:"metadata"`
	Attributes          []*Attribute                  `json:"attributes"`
	AvailableAttributes *AttributeCountableConnection `json:"availableAttributes"`
	HasPages            *bool                         `json:"hasPages"`
}

type PageTypeBulkDelete struct {
	Count  int32        `json:"count"`
	Errors []*PageError `json:"errors"`
}

type PageTypeCountableConnection struct {
	PageInfo   *PageInfo                `json:"pageInfo"`
	Edges      []*PageTypeCountableEdge `json:"edges"`
	TotalCount *int32                   `json:"totalCount"`
}

type PageTypeCountableEdge struct {
	Node   *PageType `json:"node"`
	Cursor string    `json:"cursor"`
}

type PageTypeCreate struct {
	Errors   []*PageError `json:"errors"`
	PageType *PageType    `json:"pageType"`
}

type PageTypeCreateInput struct {
	Name          *string  `json:"name"`
	Slug          *string  `json:"slug"`
	AddAttributes []string `json:"addAttributes"`
}

type PageTypeDelete struct {
	Errors   []*PageError `json:"errors"`
	PageType *PageType    `json:"pageType"`
}

type PageTypeFilterInput struct {
	Search *string `json:"search"`
}

type PageTypeReorderAttributes struct {
	PageType *PageType    `json:"pageType"`
	Errors   []*PageError `json:"errors"`
}

type PageTypeSortingInput struct {
	Direction OrderDirection    `json:"direction"`
	Field     PageTypeSortField `json:"field"`
}

type PageTypeUpdate struct {
	Errors   []*PageError `json:"errors"`
	PageType *PageType    `json:"pageType"`
}

type PageTypeUpdateInput struct {
	Name             *string  `json:"name"`
	Slug             *string  `json:"slug"`
	AddAttributes    []string `json:"addAttributes"`
	RemoveAttributes []string `json:"removeAttributes"`
}

type PageUpdate struct {
	Errors []*PageError `json:"errors"`
	Page   *Page        `json:"page"`
}

type PasswordChange struct {
	User *User `json:"user"`
}

type PaymentCapture struct {
	Payment *Payment        `json:"payment"`
	Errors  []*PaymentError `json:"errors"`
}

type PaymentCountableConnection struct {
	PageInfo   *PageInfo               `json:"pageInfo"`
	Edges      []*PaymentCountableEdge `json:"edges"`
	TotalCount *int32                  `json:"totalCount"`
}

type PaymentCountableEdge struct {
	Node   *Payment `json:"node"`
	Cursor string   `json:"cursor"`
}

type PaymentError struct {
	Field   *string          `json:"field"`
	Message *string          `json:"message"`
	Code    PaymentErrorCode `json:"code"`
}

type PaymentFilterInput struct {
	Checkouts []string `json:"checkouts"`
}

type PaymentInitialize struct {
	InitializedPayment *PaymentInitialized `json:"initializedPayment"`
	Errors             []*PaymentError     `json:"errors"`
}

type PaymentInitialized struct {
	Gateway string     `json:"gateway"`
	Name    string     `json:"name"`
	Data    JSONString `json:"data"`
}

type PaymentInput struct {
	Gateway            string                  `json:"gateway"`
	Token              *string                 `json:"token"`
	Amount             *PositiveDecimal        `json:"amount"`
	ReturnURL          *string                 `json:"returnUrl"`
	StorePaymentMethod *StorePaymentMethodEnum `json:"storePaymentMethod"`
	Metadata           []*MetadataInput        `json:"metadata"`
}

type PaymentRefund struct {
	Payment *Payment        `json:"payment"`
	Errors  []*PaymentError `json:"errors"`
}

type PaymentSource struct {
	Gateway         string          `json:"gateway"`
	PaymentMethodID *string         `json:"paymentMethodId"`
	CreditCardInfo  *CreditCard     `json:"creditCardInfo"`
	Metadata        []*MetadataItem `json:"metadata"`
}

type PaymentVoid struct {
	Payment *Payment        `json:"payment"`
	Errors  []*PaymentError `json:"errors"`
}

type Permission struct {
	Code PermissionEnum `json:"code"`
	Name string         `json:"name"`
}

type PermissionGroupCreate struct {
	Errors []*PermissionGroupError `json:"errors"`
	Group  *Group                  `json:"group"`
}

type PermissionGroupCreateInput struct {
	AddPermissions []PermissionEnum `json:"addPermissions"`
	AddUsers       []string         `json:"addUsers"`
	Name           string           `json:"name"`
}

type PermissionGroupDelete struct {
	Errors []*PermissionGroupError `json:"errors"`
	Group  *Group                  `json:"group"`
}

type PermissionGroupError struct {
	Field       *string                  `json:"field"`
	Message     *string                  `json:"message"`
	Code        PermissionGroupErrorCode `json:"code"`
	Permissions []PermissionEnum         `json:"permissions"`
	Users       []string                 `json:"users"`
}

type PermissionGroupFilterInput struct {
	Search *string `json:"search"`
}

type PermissionGroupSortingInput struct {
	Direction OrderDirection           `json:"direction"`
	Field     PermissionGroupSortField `json:"field"`
}

type PermissionGroupUpdate struct {
	Errors []*PermissionGroupError `json:"errors"`
	Group  *Group                  `json:"group"`
}

type PermissionGroupUpdateInput struct {
	AddPermissions    []PermissionEnum `json:"addPermissions"`
	AddUsers          []string         `json:"addUsers"`
	Name              *string          `json:"name"`
	RemovePermissions []PermissionEnum `json:"removePermissions"`
	RemoveUsers       []string         `json:"removeUsers"`
}

type Plugin struct {
	ID                    string                 `json:"id"`
	Name                  string                 `json:"name"`
	Description           string                 `json:"description"`
	GlobalConfiguration   *PluginConfiguration   `json:"globalConfiguration"`
	ChannelConfigurations []*PluginConfiguration `json:"channelConfigurations"`
}

type PluginConfiguration struct {
	Active        bool                 `json:"active"`
	Channel       *Channel             `json:"channel"`
	Configuration []*ConfigurationItem `json:"configuration"`
}

type PluginCountableConnection struct {
	PageInfo   *PageInfo              `json:"pageInfo"`
	Edges      []*PluginCountableEdge `json:"edges"`
	TotalCount *int32                 `json:"totalCount"`
}

type PluginCountableEdge struct {
	Node   *Plugin `json:"node"`
	Cursor string  `json:"cursor"`
}

type PluginError struct {
	Field   *string         `json:"field"`
	Message *string         `json:"message"`
	Code    PluginErrorCode `json:"code"`
}

type PluginFilterInput struct {
	StatusInChannels *PluginStatusInChannelsInput `json:"statusInChannels"`
	Search           *string                      `json:"search"`
	Type             *PluginConfigurationType     `json:"type"`
}

type PluginSortingInput struct {
	Direction OrderDirection  `json:"direction"`
	Field     PluginSortField `json:"field"`
}

type PluginStatusInChannelsInput struct {
	Active   bool     `json:"active"`
	Channels []string `json:"channels"`
}

type PluginUpdate struct {
	Plugin *Plugin        `json:"plugin"`
	Errors []*PluginError `json:"errors"`
}

type PluginUpdateInput struct {
	Active        *bool                     `json:"active"`
	Configuration []*ConfigurationItemInput `json:"configuration"`
}

type PriceInput struct {
	Currency string          `json:"currency"`
	Amount   PositiveDecimal `json:"amount"`
}

type PriceRangeInput struct {
	Gte *float64 `json:"gte"`
	Lte *float64 `json:"lte"`
}

type ProductAttributeAssign struct {
	ProductType *ProductType    `json:"productType"`
	Errors      []*ProductError `json:"errors"`
}

type ProductAttributeAssignInput struct {
	ID   string               `json:"id"`
	Type ProductAttributeType `json:"type"`
}

type ProductAttributeUnassign struct {
	ProductType *ProductType    `json:"productType"`
	Errors      []*ProductError `json:"errors"`
}

type ProductBulkDelete struct {
	Count  int32           `json:"count"`
	Errors []*ProductError `json:"errors"`
}

type ProductChannelListingAddInput struct {
	ChannelID                string   `json:"channelId"`
	IsPublished              *bool    `json:"isPublished"`
	PublicationDate          *Date    `json:"publicationDate"`
	VisibleInListings        *bool    `json:"visibleInListings"`
	IsAvailableForPurchase   *bool    `json:"isAvailableForPurchase"`
	AvailableForPurchaseDate *Date    `json:"availableForPurchaseDate"`
	AddVariants              []string `json:"addVariants"`
	RemoveVariants           []string `json:"removeVariants"`
}

type ProductChannelListingError struct {
	Field      *string          `json:"field"`
	Message    *string          `json:"message"`
	Code       ProductErrorCode `json:"code"`
	Attributes []string         `json:"attributes"`
	Values     []string         `json:"values"`
	Channels   []string         `json:"channels"`
	Variants   []string         `json:"variants"`
}

type ProductChannelListingUpdate struct {
	Product *Product                      `json:"product"`
	Errors  []*ProductChannelListingError `json:"errors"`
}

type ProductChannelListingUpdateInput struct {
	UpdateChannels []*ProductChannelListingAddInput `json:"updateChannels"`
	RemoveChannels []string                         `json:"removeChannels"`
}

type ProductCountableConnection struct {
	PageInfo   *PageInfo               `json:"pageInfo"`
	Edges      []*ProductCountableEdge `json:"edges"`
	TotalCount *int32                  `json:"totalCount"`
}

type ProductCountableEdge struct {
	Node   *Product `json:"node"`
	Cursor string   `json:"cursor"`
}

type ProductCreate struct {
	Errors  []*ProductError `json:"errors"`
	Product *Product        `json:"product"`
}

type ProductCreateInput struct {
	Attributes  []*AttributeValueInput `json:"attributes"`
	Category    *string                `json:"category"`
	ChargeTaxes *bool                  `json:"chargeTaxes"`
	Collections []string               `json:"collections"`
	Description JSONString             `json:"description"`
	Name        *string                `json:"name"`
	Slug        *string                `json:"slug"`
	TaxCode     *string                `json:"taxCode"`
	Seo         *SeoInput              `json:"seo"`
	Weight      *string                `json:"weight"`
	Rating      *float64               `json:"rating"`
	ProductType string                 `json:"productType"`
}

type ProductDelete struct {
	Errors  []*ProductError `json:"errors"`
	Product *Product        `json:"product"`
}

type ProductError struct {
	Field      *string          `json:"field"`
	Message    *string          `json:"message"`
	Code       ProductErrorCode `json:"code"`
	Attributes []string         `json:"attributes"`
	Values     []string         `json:"values"`
}

type ProductFilterInput struct {
	IsPublished           *bool                    `json:"isPublished"`
	Collections           []string                 `json:"collections"`
	Categories            []string                 `json:"categories"`
	HasCategory           *bool                    `json:"hasCategory"`
	Attributes            []*AttributeInput        `json:"attributes"`
	StockAvailability     *StockAvailability       `json:"stockAvailability"`
	Stocks                *ProductStockFilterInput `json:"stocks"`
	Search                *string                  `json:"search"`
	Metadata              []*MetadataInput         `json:"metadata"`
	Price                 *PriceRangeInput         `json:"price"`
	MinimalPrice          *PriceRangeInput         `json:"minimalPrice"`
	ProductTypes          []string                 `json:"productTypes"`
	GiftCard              *bool                    `json:"giftCard"`
	Ids                   []string                 `json:"ids"`
	HasPreorderedVariants *bool                    `json:"hasPreorderedVariants"`
	Channel               *string                  `json:"channel"` // can be either channel id or channel slug
}

func (p *ProductFilterInput) ToSystemProductFilterInput() *model.ProductFilterInput {
	systemAttributeFilter := lo.Map(p.Attributes, func(item *AttributeInput, _ int) *model.AttributeFilter {
		res := &model.AttributeFilter{
			Slug:   item.Slug,
			Values: item.Values,
			ValuesRange: &struct {
				Gte *int32
				Lte *int32
			}{
				Gte: item.ValuesRange.Gte,
				Lte: item.ValuesRange.Lte,
			},
			Boolean: item.Boolean,
		}

		if item.DateTime != nil && item.DateTime.Gte != nil {
			res.DateTime.Gte = &item.DateTime.Gte.Time
		}
		if item.DateTime != nil && item.DateTime.Lte != nil {
			res.DateTime.Lte = &item.DateTime.Lte.Time
		}

		if item.Date != nil && item.Date.Gte != nil {
			res.Date.Gte = &item.Date.Gte.Time
		}
		if item.Date != nil && item.Date.Lte != nil {
			res.Date.Lte = &item.Date.Lte.Time
		}
		return res
	})

	metadata := lo.Map(p.Metadata, func(m *MetadataInput, _ int) *struct {
		Key   string
		Value string
	} {
		res := &struct {
			Key   string
			Value string
		}{}
		if m != nil {
			res.Key = m.Key
			res.Value = m.Value
		}
		return res
	})

	res := &model.ProductFilterInput{
		IsPublished:           p.IsPublished,
		Collections:           p.Collections,
		Categories:            p.Categories,
		HasCategory:           p.HasCategory,
		Attributes:            systemAttributeFilter,
		Search:                p.Search,
		Metadata:              metadata,
		ProductTypes:          p.ProductTypes,
		GiftCard:              p.GiftCard,
		Ids:                   p.Ids,
		HasPreorderedVariants: p.HasPreorderedVariants,
		Channel:               p.Channel,
		StockAvailability:     p.StockAvailability,
	}

	if p.Stocks != nil {
		res.Stocks = &struct {
			WarehouseIds []string
			Quantity     *struct {
				Gte *int32
				Lte *int32
			}
		}{
			WarehouseIds: p.Stocks.WarehouseIds,
		}
		if p.Stocks.Quantity != nil {
			res.Stocks.Quantity = &struct {
				Gte *int32
				Lte *int32
			}{
				Gte: p.Stocks.Quantity.Gte,
				Lte: p.Stocks.Quantity.Lte,
			}
		}
	}
	if p.Price != nil {
		res.Price = &struct {
			Gte *float64
			Lte *float64
		}{
			Gte: p.Price.Gte,
			Lte: p.Price.Lte,
		}
	}
	if p.MinimalPrice != nil {
		res.MinimalPrice = &struct {
			Gte *float64
			Lte *float64
		}{
			Gte: p.MinimalPrice.Gte,
			Lte: p.MinimalPrice.Lte,
		}
	}

	return res
}

type ProductInput struct {
	Attributes  []*AttributeValueInput `json:"attributes"`
	Category    *string                `json:"category"`
	ChargeTaxes *bool                  `json:"chargeTaxes"`
	Collections []string               `json:"collections"`
	Description JSONString             `json:"description"`
	Name        *string                `json:"name"`
	Slug        *string                `json:"slug"`
	TaxCode     *string                `json:"taxCode"`
	Seo         *SeoInput              `json:"seo"`
	Weight      *string                `json:"weight"`
	Rating      *float64               `json:"rating"`
}

type ProductMediaBulkDelete struct {
	Count  int32           `json:"count"`
	Errors []*ProductError `json:"errors"`
}

type ProductMediaCreate struct {
	Product *Product        `json:"product"`
	Media   *ProductMedia   `json:"media"`
	Errors  []*ProductError `json:"errors"`
}

type ProductMediaCreateInput struct {
	Alt      *string         `json:"alt"`
	Image    *graphql.Upload `json:"image"`
	Product  string          `json:"product"`
	MediaURL *string         `json:"mediaUrl"`
}

type ProductMediaDelete struct {
	Product *Product        `json:"product"`
	Media   *ProductMedia   `json:"media"`
	Errors  []*ProductError `json:"errors"`
}

type ProductMediaReorder struct {
	Product *Product        `json:"product"`
	Media   []*ProductMedia `json:"media"`
	Errors  []*ProductError `json:"errors"`
}

type ProductMediaUpdate struct {
	Product *Product        `json:"product"`
	Media   *ProductMedia   `json:"media"`
	Errors  []*ProductError `json:"errors"`
}

type ProductMediaUpdateInput struct {
	Alt *string `json:"alt"`
}

type ProductOrder struct {
	Direction   OrderDirection     `json:"direction"`
	Channel     *string            `json:"channel"` // DEPRECATED, don't use this field
	AttributeID *string            `json:"attributeId"`
	Field       *ProductOrderField `json:"field"`
}

func (o *ProductOrder) ToSystemProductOrder() *model.ProductOrder {
	if o == nil {
		return nil
	}

	res := &model.ProductOrder{
		Direction:   o.Direction,
		Field:       o.Field,
		AttributeID: o.AttributeID,
	}

	return res
}

type ProductPricingInfo struct {
	OnSale                  *bool            `json:"onSale"`
	Discount                *TaxedMoney      `json:"discount"`
	DiscountLocalCurrency   *TaxedMoney      `json:"discountLocalCurrency"`
	PriceRange              *TaxedMoneyRange `json:"priceRange"`
	PriceRangeUndiscounted  *TaxedMoneyRange `json:"priceRangeUndiscounted"`
	PriceRangeLocalCurrency *TaxedMoneyRange `json:"priceRangeLocalCurrency"`
}

type ProductReorderAttributeValues struct {
	Product *Product        `json:"product"`
	Errors  []*ProductError `json:"errors"`
}

type ProductStockFilterInput struct {
	WarehouseIds []string       `json:"warehouseIds"`
	Quantity     *IntRangeInput `json:"quantity"`
}

type ProductTranslatableContent struct {
	ID             string              `json:"id"`
	SeoTitle       *string             `json:"seoTitle"`
	SeoDescription *string             `json:"seoDescription"`
	Name           string              `json:"name"`
	Description    JSONString          `json:"description"`
	Translation    *ProductTranslation `json:"translation"`
}

type ProductTranslate struct {
	Errors  []*TranslationError `json:"errors"`
	Product *Product            `json:"product"`
}

type ProductTranslation struct {
	ID             string           `json:"id"`
	SeoTitle       *string          `json:"seoTitle"`
	SeoDescription *string          `json:"seoDescription"`
	Name           *string          `json:"name"`
	Description    JSONString       `json:"description"`
	Language       *LanguageDisplay `json:"language"`
}

type ProductTypeBulkDelete struct {
	Count  int32           `json:"count"`
	Errors []*ProductError `json:"errors"`
}

type ProductTypeCountableConnection struct {
	PageInfo   *PageInfo                   `json:"pageInfo"`
	Edges      []*ProductTypeCountableEdge `json:"edges"`
	TotalCount *int32                      `json:"totalCount"`
}

type ProductTypeCountableEdge struct {
	Node   *ProductType `json:"node"`
	Cursor string       `json:"cursor"`
}

type ProductTypeCreate struct {
	Errors      []*ProductError `json:"errors"`
	ProductType *ProductType    `json:"productType"`
}

type ProductTypeDelete struct {
	Errors      []*ProductError `json:"errors"`
	ProductType *ProductType    `json:"productType"`
}

type ProductTypeFilterInput struct {
	Search       *string                  `json:"search"`
	Configurable *ProductTypeConfigurable `json:"configurable"`
	ProductType  *ProductTypeEnum         `json:"productType"`
	Metadata     []*MetadataInput         `json:"metadata"`
	Kind         *ProductTypeKindEnum     `json:"kind"`
	Ids          []string                 `json:"ids"`
}

type ProductTypeInput struct {
	Name               *string              `json:"name"`
	Slug               *string              `json:"slug"`
	Kind               *ProductTypeKindEnum `json:"kind"`
	HasVariants        *bool                `json:"hasVariants"`
	ProductAttributes  []string             `json:"productAttributes"`
	VariantAttributes  []string             `json:"variantAttributes"`
	IsShippingRequired *bool                `json:"isShippingRequired"`
	IsDigital          *bool                `json:"isDigital"`
	Weight             *string              `json:"weight"`
	TaxCode            *string              `json:"taxCode"`
}

type ProductTypeReorderAttributes struct {
	ProductType *ProductType    `json:"productType"`
	Errors      []*ProductError `json:"errors"`
}

type ProductTypeSortingInput struct {
	Direction OrderDirection       `json:"direction"`
	Field     ProductTypeSortField `json:"field"`
}

type ProductTypeUpdate struct {
	Errors      []*ProductError `json:"errors"`
	ProductType *ProductType    `json:"productType"`
}

type ProductUpdate struct {
	Errors  []*ProductError `json:"errors"`
	Product *Product        `json:"product"`
}

type ProductVariantBulkCreate struct {
	Count           int32               `json:"count"`
	ProductVariants []*ProductVariant   `json:"productVariants"`
	Errors          []*BulkProductError `json:"errors"`
}

type ProductVariantBulkCreateInput struct {
	Attributes      []*BulkAttributeValueInput              `json:"attributes"`
	Sku             *string                                 `json:"sku"`
	TrackInventory  *bool                                   `json:"trackInventory"`
	Weight          *string                                 `json:"weight"`
	Stocks          []*StockInput                           `json:"stocks"`
	ChannelListings []*ProductVariantChannelListingAddInput `json:"channelListings"`
}

type ProductVariantBulkDelete struct {
	Count  int32           `json:"count"`
	Errors []*ProductError `json:"errors"`
}

type PreorderThreshold struct {
	Quantity  *int32 `json:"quantity"`
	SoldUnits int32  `json:"soldUnits"`
}

type ProductVariantChannelListingAddInput struct {
	ChannelID string           `json:"channelId"`
	Price     PositiveDecimal  `json:"price"`
	CostPrice *PositiveDecimal `json:"costPrice"`
}

type ProductVariantChannelListingUpdate struct {
	Variant *ProductVariant               `json:"variant"`
	Errors  []*ProductChannelListingError `json:"errors"`
}

type ProductVariantCountableConnection struct {
	PageInfo   *PageInfo                      `json:"pageInfo"`
	Edges      []*ProductVariantCountableEdge `json:"edges"`
	TotalCount *int32                         `json:"totalCount"`
}

type ProductVariantCountableEdge struct {
	Node   *ProductVariant `json:"node"`
	Cursor string          `json:"cursor"`
}

type ProductVariantCreate struct {
	Errors         []*ProductError `json:"errors"`
	ProductVariant *ProductVariant `json:"productVariant"`
}

type ProductVariantCreateInput struct {
	Attributes     []*AttributeValueInput `json:"attributes"`
	Sku            *string                `json:"sku"`
	TrackInventory *bool                  `json:"trackInventory"`
	Weight         *string                `json:"weight"`
	Product        string                 `json:"product"`
	Stocks         []*StockInput          `json:"stocks"`
}

type ProductVariantDelete struct {
	Errors         []*ProductError `json:"errors"`
	ProductVariant *ProductVariant `json:"productVariant"`
}

type ProductVariantFilterInput struct {
	Search   *string          `json:"search"`
	Sku      []string         `json:"sku"`
	Metadata []*MetadataInput `json:"metadata"`
}

type ProductVariantInput struct {
	Attributes     []*AttributeValueInput `json:"attributes"`
	Sku            *string                `json:"sku"`
	TrackInventory *bool                  `json:"trackInventory"`
	Weight         *string                `json:"weight"`
}

type ProductVariantReorder struct {
	Product *Product        `json:"product"`
	Errors  []*ProductError `json:"errors"`
}

type ProductVariantReorderAttributeValues struct {
	ProductVariant *ProductVariant `json:"productVariant"`
	Errors         []*ProductError `json:"errors"`
}

type ProductVariantSetDefault struct {
	Product *Product        `json:"product"`
	Errors  []*ProductError `json:"errors"`
}

type ProductVariantStocksCreate struct {
	ProductVariant *ProductVariant   `json:"productVariant"`
	Errors         []*BulkStockError `json:"errors"`
}

type ProductVariantStocksDelete struct {
	ProductVariant *ProductVariant `json:"productVariant"`
	Errors         []*StockError   `json:"errors"`
}

type ProductVariantStocksUpdate struct {
	ProductVariant *ProductVariant   `json:"productVariant"`
	Errors         []*BulkStockError `json:"errors"`
}

type ProductVariantTranslatableContent struct {
	ID              string                               `json:"id"`
	Name            string                               `json:"name"`
	Translation     *ProductVariantTranslation           `json:"translation"`
	AttributeValues []*AttributeValueTranslatableContent `json:"attributeValues"`
}

type ProductVariantTranslate struct {
	Errors         []*TranslationError `json:"errors"`
	ProductVariant *ProductVariant     `json:"productVariant"`
}

type ProductVariantTranslation struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Language *LanguageDisplay `json:"language"`
}

type ProductVariantUpdate struct {
	Errors         []*ProductError `json:"errors"`
	ProductVariant *ProductVariant `json:"productVariant"`
}

type PublishableChannelListingInput struct {
	ChannelID       string `json:"channelId"`
	IsPublished     *bool  `json:"isPublished"`
	PublicationDate *Date  `json:"publicationDate"`
}

type RefreshToken struct {
	Token *string `json:"token"`
	User  *User   `json:"user"`
}

type ReorderInput struct {
	ID        string `json:"id"`
	SortOrder *int32 `json:"sortOrder"`
}

type RequestEmailChange struct {
	User *User `json:"user"`
}

type RequestPasswordReset struct {
	Ok bool `json:"ok"`
}

type SaleAddCatalogues struct {
	Sale   *Sale            `json:"sale"`
	Errors []*DiscountError `json:"errors"`
}

type SaleBulkDelete struct {
	Count  int32            `json:"count"`
	Errors []*DiscountError `json:"errors"`
}

type SaleChannelListingAddInput struct {
	ChannelID     string          `json:"channelId"`
	DiscountValue PositiveDecimal `json:"discountValue"`
}

type SaleChannelListingInput struct {
	AddChannels    []*SaleChannelListingAddInput `json:"addChannels"`
	RemoveChannels []string                      `json:"removeChannels"`
}

type SaleChannelListingUpdate struct {
	Sale   *Sale            `json:"sale"`
	Errors []*DiscountError `json:"errors"`
}

type SaleCountableConnection struct {
	PageInfo   *PageInfo            `json:"pageInfo"`
	Edges      []*SaleCountableEdge `json:"edges"`
	TotalCount *int32               `json:"totalCount"`
}

type SaleCountableEdge struct {
	Node   *Sale  `json:"node"`
	Cursor string `json:"cursor"`
}

type SaleCreate struct {
	Errors []*DiscountError `json:"errors"`
	Sale   *Sale            `json:"sale"`
}

type SaleDelete struct {
	Errors []*DiscountError `json:"errors"`
	Sale   *Sale            `json:"sale"`
}

type SaleFilterInput struct {
	Status   []*DiscountStatusEnum  `json:"status"`
	SaleType *DiscountValueTypeEnum `json:"saleType"`
	Started  *DateTimeRangeInput    `json:"started"`
	Search   *string                `json:"search"`
	Metadata []*MetadataFilter      `json:"metadata"`
}

type SaleInput struct {
	Name        *string                `json:"name"`
	Type        *DiscountValueTypeEnum `json:"type"`
	Value       *PositiveDecimal       `json:"value"`
	Products    []string               `json:"products"`
	Variants    []string               `json:"variants"`
	Categories  []string               `json:"categories"`
	Collections []string               `json:"collections"`
	StartDate   *DateTime              `json:"startDate"`
	EndDate     *DateTime              `json:"endDate"`
}

type SaleRemoveCatalogues struct {
	Sale   *Sale            `json:"sale"`
	Errors []*DiscountError `json:"errors"`
}

type SaleSortingInput struct {
	Direction OrderDirection `json:"direction"`
	Channel   *string        `json:"channel"`
	Field     SaleSortField  `json:"field"`
}

type SaleTranslatableContent struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Translation *SaleTranslation `json:"translation"`
}

type SaleTranslate struct {
	Errors []*TranslationError `json:"errors"`
	Sale   *Sale               `json:"sale"`
}

type SaleTranslation struct {
	ID       string           `json:"id"`
	Name     *string          `json:"name"`
	Language *LanguageDisplay `json:"language"`
}

type SaleUpdate struct {
	Errors []*DiscountError `json:"errors"`
	Sale   *Sale            `json:"sale"`
}

type SelectedAttribute struct {
	Attribute *Attribute        `json:"attribute"`
	Values    []*AttributeValue `json:"values"`
}

type SeoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type SetPassword struct {
	Token        *string `json:"token"`
	RefreshToken *string `json:"refreshToken"`
	CsrfToken    *string `json:"csrfToken"`
	User         *User   `json:"user"`
}

type ShippingError struct {
	Field      *string           `json:"field"`
	Message    *string           `json:"message"`
	Code       ShippingErrorCode `json:"code"`
	Warehouses []string          `json:"warehouses"`
	Channels   []string          `json:"channels"`
}

type ShippingMethodChannelListingAddInput struct {
	ChannelID         string           `json:"channelId"`
	Price             *PositiveDecimal `json:"price"`
	MinimumOrderPrice *PositiveDecimal `json:"minimumOrderPrice"`
	MaximumOrderPrice *PositiveDecimal `json:"maximumOrderPrice"`
}

type ShippingMethodChannelListingInput struct {
	AddChannels    []*ShippingMethodChannelListingAddInput `json:"addChannels"`
	RemoveChannels []string                                `json:"removeChannels"`
}

type ShippingMethodChannelListingUpdate struct {
	ShippingMethod *ShippingMethod  `json:"shippingMethod"`
	Errors         []*ShippingError `json:"errors"`
}

type ShippingMethodPostalCodeRule struct {
	Start         *string                          `json:"start"`
	End           *string                          `json:"end"`
	InclusionType *PostalCodeRuleInclusionTypeEnum `json:"inclusionType"`
	ID            string                           `json:"id"`
}

type ShippingMethodTranslatableContent struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description JSONString                 `json:"description"`
	Translation *ShippingMethodTranslation `json:"translation"`
}

type ShippingMethodTranslation struct {
	ID          string           `json:"id"`
	Name        *string          `json:"name"`
	Description JSONString       `json:"description"`
	Language    *LanguageDisplay `json:"language"`
}

type ShippingPostalCodeRulesCreateInputRange struct {
	Start string  `json:"start"`
	End   *string `json:"end"`
}

type ShippingPriceBulkDelete struct {
	Count  int32            `json:"count"`
	Errors []*ShippingError `json:"errors"`
}

type ShippingPriceCreate struct {
	ShippingZone   *ShippingZone    `json:"shippingZone"`
	ShippingMethod *ShippingMethod  `json:"shippingMethod"`
	Errors         []*ShippingError `json:"errors"`
}

type ShippingPriceDelete struct {
	ShippingMethod *ShippingMethod  `json:"shippingMethod"`
	ShippingZone   *ShippingZone    `json:"shippingZone"`
	Errors         []*ShippingError `json:"errors"`
}

type ShippingPriceExcludeProducts struct {
	ShippingMethod *ShippingMethod  `json:"shippingMethod"`
	Errors         []*ShippingError `json:"errors"`
}

type ShippingPriceExcludeProductsInput struct {
	Products []string `json:"products"`
}

type ShippingPriceInput struct {
	Name                  *string                                    `json:"name"`
	Description           JSONString                                 `json:"description"`
	MinimumOrderWeight    *Weight                                    `json:"minimumOrderWeight"`
	MaximumOrderWeight    *Weight                                    `json:"maximumOrderWeight"`
	MaximumDeliveryDays   *int32                                     `json:"maximumDeliveryDays"`
	MinimumDeliveryDays   *int32                                     `json:"minimumDeliveryDays"`
	Type                  *ShippingMethodTypeEnum                    `json:"type"`
	ShippingZone          *string                                    `json:"shippingZone"`
	AddPostalCodeRules    []*ShippingPostalCodeRulesCreateInputRange `json:"addPostalCodeRules"`
	DeletePostalCodeRules []string                                   `json:"deletePostalCodeRules"`
	InclusionType         *PostalCodeRuleInclusionTypeEnum           `json:"inclusionType"`
}

// NOTE: Patch must be called after calling Validate().
//
// returned `updated` boolean value indicates wether given `method` is modified.
func (s *ShippingPriceInput) Patch(method *model.ShippingMethod) (updated bool) {
	updated = true

	switch {
	case s.Name != nil && *s.Name != method.Name:
		method.Name = *s.Name
		fallthrough

	case s.Description != nil:
		for key, value := range s.Description {
			method.Description[key] = value
		}
		fallthrough

	case s.MinimumOrderWeight != nil:
		method.MinimumOrderWeight = float32(s.MinimumOrderWeight.Value)
		fallthrough

	case s.MaximumOrderWeight != nil:
		method.MaximumOrderWeight = (*float32)(unsafe.Pointer(&s.MaximumOrderWeight.Value))
		fallthrough

	case s.MinimumDeliveryDays != nil:
		method.MinimumDeliveryDays = (*int)(unsafe.Pointer(s.MinimumDeliveryDays))
		fallthrough

	case s.MaximumDeliveryDays != nil:
		method.MaximumDeliveryDays = (*int)(unsafe.Pointer(s.MaximumDeliveryDays))
		fallthrough

	case s.Type != nil && s.Type.IsValid() && *s.Type != method.Type:
		method.Type = *s.Type
		fallthrough

	case s.ShippingZone != nil && *s.ShippingZone != method.ShippingZoneID: // NOTE: s.ShippingZone is already converted and validated
		method.ShippingZoneID = *s.ShippingZone

	default:
		updated = false
	}

	return updated
}

func (s *ShippingPriceInput) Validate(api string) *model.AppError {
	// clean weights:
	if s.MinimumOrderWeight != nil {
		if s.MinimumOrderWeight.Value < 0 {
			return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MinimumOrderWeight"}, "shipping cannot have negative weight", http.StatusBadRequest)
		}
		if measurement.WEIGHT_UNIT_STRINGS[s.MinimumOrderWeight.Unit] == "" { // invalid unit
			return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MinimumOrderWeight"}, "weight unit is invalid", http.StatusBadRequest)
		}
	}
	if s.MaximumOrderWeight != nil {
		if s.MaximumOrderWeight.Value < 0 {
			return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MaximumOrderWeight"}, "shipping cannot have negative weight", http.StatusBadRequest)
		}
		if measurement.WEIGHT_UNIT_STRINGS[s.MaximumOrderWeight.Unit] == "" { // invalid unit
			return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MaximumOrderWeight"}, "weight unit is invalid", http.StatusBadRequest)
		}
	}

	if s.MinimumOrderWeight != nil &&
		s.MaximumOrderWeight != nil &&
		(s.MinimumOrderWeight.Unit == s.MaximumOrderWeight.Unit || s.MinimumOrderWeight.Value >= s.MaximumOrderWeight.Value) {
		return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MaximumOrderWeight / MinimumOrderWeight"}, "weight units must be the same and min weight must less than (<) max weight", http.StatusBadRequest)
	}

	// clean delivery time
	// - check if minimum_delivery_days is not higher than maximum_delivery_days
	// - check if minimum_delivery_days and maximum_delivery_days are positive values
	if s.MinimumDeliveryDays != nil && *s.MinimumDeliveryDays < 0 {
		return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MinimumDeliveryDays"}, "delivery days cannot be negative", http.StatusBadRequest)
	}
	if s.MaximumDeliveryDays != nil && *s.MaximumDeliveryDays < 0 {
		return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MaximumDeliveryDays"}, "delivery days cannot be negative", http.StatusBadRequest)
	}
	if s.MinimumDeliveryDays != nil && s.MaximumDeliveryDays != nil && *s.MinimumDeliveryDays >= *s.MaximumDeliveryDays {
		return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MinimumDeliveryDays, MaximumDeliveryDays"}, "min delivery day must less than max delivery days", http.StatusBadRequest)
	}

	// clean postal code rules
	s.DeletePostalCodeRules = decodeBase64Strings(s.DeletePostalCodeRules...)
	if !lo.EveryBy(s.DeletePostalCodeRules, model.IsValidId) {
		return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "delete postal code rules"}, "please provide valid delete postal code rule ids", http.StatusBadRequest)
	}
	if len(s.AddPostalCodeRules) > 0 && (s.InclusionType == nil || !s.InclusionType.IsValid()) {
		return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "inclusion type"}, "inclusion type is required when add postal code rules are provided", http.StatusBadRequest)
	}

	if s.Type != nil && !s.Type.IsValid() {
		return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "type"}, "please provide valid type", http.StatusBadRequest)
	}

	if s.ShippingZone != nil {
		shippingZoneID := decodeBase64String(*s.ShippingZone)
		if !model.IsValidId(shippingZoneID) {
			return model.NewAppError(api, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "shipping zone"}, "please provide valid shipping zone id", http.StatusBadRequest)
		}
		s.ShippingZone = &shippingZoneID // NOTE: no need to convert later (in case nil error is returned)
	}

	return nil
}

type ShippingPriceRemoveProductFromExclude struct {
	ShippingMethod *ShippingMethod  `json:"shippingMethod"`
	Errors         []*ShippingError `json:"errors"`
}

type ShippingPriceTranslate struct {
	Errors         []*TranslationError `json:"errors"`
	ShippingMethod *ShippingMethod     `json:"shippingMethod"`
}

type ShippingPriceTranslationInput struct {
	Name        *string    `json:"name"`
	Description JSONString `json:"description"`
}

type ShippingPriceUpdate struct {
	ShippingZone   *ShippingZone    `json:"shippingZone"`
	ShippingMethod *ShippingMethod  `json:"shippingMethod"`
	Errors         []*ShippingError `json:"errors"`
}

type ShippingZoneBulkDelete struct {
	Count  int32            `json:"count"`
	Errors []*ShippingError `json:"errors"`
}

type ShippingZoneCountableConnection struct {
	PageInfo   *PageInfo                    `json:"pageInfo"`
	Edges      []*ShippingZoneCountableEdge `json:"edges"`
	TotalCount *int32                       `json:"totalCount"`
}

type ShippingZoneCountableEdge struct {
	Node   *ShippingZone `json:"node"`
	Cursor string        `json:"cursor"`
}

type ShippingZoneCreate struct {
	Errors       []*ShippingError `json:"errors"`
	ShippingZone *ShippingZone    `json:"shippingZone"`
}

type ShippingZoneCreateInput struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	Countries     []string `json:"countries"`
	Default       *bool    `json:"default"`
	AddWarehouses []string `json:"addWarehouses"`
	AddChannels   []string `json:"addChannels"`
}

type ShippingZoneDelete struct {
	Errors       []*ShippingError `json:"errors"`
	ShippingZone *ShippingZone    `json:"shippingZone"`
}

type ShippingZoneFilterInput struct {
	Search   *string  `json:"search"`
	Channels []string `json:"channels"`
}

type ShippingZoneUpdate struct {
	Errors       []*ShippingError `json:"errors"`
	ShippingZone *ShippingZone    `json:"shippingZone"`
}

type ShippingZoneUpdateInput struct {
	Name             *string  `json:"name"`
	Description      *string  `json:"description"`
	Countries        []string `json:"countries"`
	Default          *bool    `json:"default"`
	AddWarehouses    []string `json:"addWarehouses"`
	AddChannels      []string `json:"addChannels"`
	RemoveWarehouses []string `json:"removeWarehouses"`
	RemoveChannels   []string `json:"removeChannels"`
}

type ShopAddressUpdate struct {
	Shop   *Shop        `json:"shop"`
	Errors []*ShopError `json:"errors"`
}

type ShopDomainUpdate struct {
	Shop   *Shop        `json:"shop"`
	Errors []*ShopError `json:"errors"`
}

type ShopError struct {
	Field   *string       `json:"field"`
	Message *string       `json:"message"`
	Code    ShopErrorCode `json:"code"`
}

type ShopFetchTaxRates struct {
	Shop   *Shop        `json:"shop"`
	Errors []*ShopError `json:"errors"`
}

type ShopSettingsInput struct {
	HeaderText                          *string          `json:"headerText"`
	Description                         *string          `json:"description"`
	IncludeTaxesInPrices                *bool            `json:"includeTaxesInPrices"`
	DisplayGrossPrices                  *bool            `json:"displayGrossPrices"`
	ChargeTaxesOnShipping               *bool            `json:"chargeTaxesOnShipping"`
	TrackInventoryByDefault             *bool            `json:"trackInventoryByDefault"`
	DefaultWeightUnit                   *WeightUnitsEnum `json:"defaultWeightUnit"`
	AutomaticFulfillmentDigitalProducts *bool            `json:"automaticFulfillmentDigitalProducts"`
	FulfillmentAutoApprove              *bool            `json:"fulfillmentAutoApprove"`
	FulfillmentAllowUnpaid              *bool            `json:"fulfillmentAllowUnpaid"`
	DefaultDigitalMaxDownloads          *int32           `json:"defaultDigitalMaxDownloads"`
	DefaultDigitalURLValidDays          *int32           `json:"defaultDigitalUrlValidDays"`
	DefaultMailSenderName               *string          `json:"defaultMailSenderName"`
	DefaultMailSenderAddress            *string          `json:"defaultMailSenderAddress"`
	CustomerSetPasswordURL              *string          `json:"customerSetPasswordUrl"`
}

type ShopSettingsTranslate struct {
	Shop   *Shop               `json:"shop"`
	Errors []*TranslationError `json:"errors"`
}

type ShopSettingsTranslationInput struct {
	HeaderText  *string `json:"headerText"`
	Description *string `json:"description"`
}

type ShopSettingsUpdate struct {
	Shop   *Shop        `json:"shop"`
	Errors []*ShopError `json:"errors"`
}

type ShopTranslation struct {
	ID          string           `json:"id"`
	HeaderText  string           `json:"headerText"`
	Description string           `json:"description"`
	Language    *LanguageDisplay `json:"language"`
}

type SiteDomainInput struct {
	Domain *string `json:"domain"`
	Name   *string `json:"name"`
}

type StaffBulkDelete struct {
	Count  int32         `json:"count"`
	Errors []*StaffError `json:"errors"`
}

type StaffCreate struct {
	Errors []*StaffError `json:"errors"`
	User   *User         `json:"user"`
}

type StaffCreateInput struct {
	FirstName   *string  `json:"firstName"`
	LastName    *string  `json:"lastName"`
	Email       *string  `json:"email"`
	IsActive    *bool    `json:"isActive"`
	Note        *string  `json:"note"`
	AddGroups   []string `json:"addGroups"`
	RedirectURL *string  `json:"redirectUrl"`
}

type StaffDelete struct {
	Errors []*StaffError `json:"errors"`
	User   *User         `json:"user"`
}

type StaffError struct {
	Field       *string                `json:"field"`
	Message     *string                `json:"message"`
	Code        AccountErrorCode       `json:"code"`
	AddressType *model.AddressTypeEnum `json:"addressType"`
	Permissions []PermissionEnum       `json:"permissions"`
	Groups      []string               `json:"groups"`
	Users       []string               `json:"users"`
}

type StaffNotificationRecipientCreate struct {
	Errors                     []*ShopError                `json:"errors"`
	StaffNotificationRecipient *StaffNotificationRecipient `json:"staffNotificationRecipient"`
}

type StaffNotificationRecipientDelete struct {
	Errors                     []*ShopError                `json:"errors"`
	StaffNotificationRecipient *StaffNotificationRecipient `json:"staffNotificationRecipient"`
}

type StaffNotificationRecipientInput struct {
	User   *string `json:"user"`
	Email  *string `json:"email"`
	Active *bool   `json:"active"`
}

type StaffNotificationRecipientUpdate struct {
	Errors                     []*ShopError                `json:"errors"`
	StaffNotificationRecipient *StaffNotificationRecipient `json:"staffNotificationRecipient"`
}

type StaffUpdate struct {
	Errors []*StaffError `json:"errors"`
	User   *User         `json:"user"`
}

type StaffUpdateInput struct {
	FirstName    *string  `json:"firstName"`
	LastName     *string  `json:"lastName"`
	Email        *string  `json:"email"`
	IsActive     *bool    `json:"isActive"`
	Note         *string  `json:"note"`
	AddGroups    []string `json:"addGroups"`
	RemoveGroups []string `json:"removeGroups"`
}

type StaffUserInput struct {
	Status *StaffMemberStatus `json:"status"`
	Search *string            `json:"search"`
}

type StockCountableConnection struct {
	PageInfo   *PageInfo             `json:"pageInfo"`
	Edges      []*StockCountableEdge `json:"edges"`
	TotalCount *int32                `json:"totalCount"`
}

type StockCountableEdge struct {
	Node   *Stock `json:"node"`
	Cursor string `json:"cursor"`
}

type StockError struct {
	Field   *string        `json:"field"`
	Message *string        `json:"message"`
	Code    StockErrorCode `json:"code"`
}

type StockFilterInput struct {
	Quantity *int32  `json:"quantity"`
	Search   *string `json:"search"`
}

type StockInput struct {
	Warehouse string `json:"warehouse"`
	Quantity  int32  `json:"quantity"`
}

type TaxType struct {
	Description *string `json:"description"`
	TaxCode     *string `json:"taxCode"`
}

type TaxedMoney struct {
	Currency string `json:"currency"`
	Gross    *Money `json:"gross"`
	Net      *Money `json:"net"`
	Tax      *Money `json:"tax"`
}

type TaxedMoneyRange struct {
	Start *TaxedMoney `json:"start"`
	Stop  *TaxedMoney `json:"stop"`
}

type TimePeriod struct {
	Amount int32                `json:"amount"`
	Type   model.TimePeriodType `json:"type"`
}

type TimePeriodInputType struct {
	Amount int32                `json:"amount"`
	Type   model.TimePeriodType `json:"type"`
}

type TokenCreateInput struct {
	ID       string `json:"id"`
	LoginID  string `json:"loginId"`
	Password string `json:"password"`
	Token    string `json:"token"`
	DeviceID string `json:"deviceId"`
	LdapOnly string `json:"ldapOnly"`
}

type TranslatableItemConnection struct {
	PageInfo   *PageInfo               `json:"pageInfo"`
	Edges      []*TranslatableItemEdge `json:"edges"`
	TotalCount *int32                  `json:"totalCount"`
}

type TranslatableItemEdge struct {
	Node   TranslatableItem `json:"node"`
	Cursor string           `json:"cursor"`
}

type TranslationError struct {
	Field   *string              `json:"field"`
	Message *string              `json:"message"`
	Code    TranslationErrorCode `json:"code"`
}

type TranslationInput struct {
	SeoTitle       *string    `json:"seoTitle"`
	SeoDescription *string    `json:"seoDescription"`
	Name           *string    `json:"name"`
	Description    JSONString `json:"description"`
}

type UpdateInvoiceInput struct {
	Number *string `json:"number"`
	URL    *string `json:"url"`
}

type UpdateMetadata struct {
	Errors []*MetadataError   `json:"errors"`
	Item   ObjectWithMetadata `json:"item"`
}

type UpdatePrivateMetadata struct {
	Errors []*MetadataError   `json:"errors"`
	Item   ObjectWithMetadata `json:"item"`
}

type UploadError struct {
	Field   *string         `json:"field"`
	Message *string         `json:"message"`
	Code    UploadErrorCode `json:"code"`
}

type UserAvatarDelete struct {
	User   *User           `json:"user"`
	Errors []*AccountError `json:"errors"`
}

type UserAvatarUpdate struct {
	User   *User           `json:"user"`
	Errors []*AccountError `json:"errors"`
}

type UserBulkSetActive struct {
	Count  int32           `json:"count"`
	Errors []*AccountError `json:"errors"`
}

type UserCountableConnection struct {
	PageInfo   *PageInfo            `json:"pageInfo"`
	Edges      []*UserCountableEdge `json:"edges"`
	TotalCount *int32               `json:"totalCount"`
}

type UserCountableEdge struct {
	Node   *User  `json:"node"`
	Cursor string `json:"cursor"`
}

type UserCreateInput struct {
	DefaultBillingAddress  *AddressInput     `json:"defaultBillingAddress"`
	DefaultShippingAddress *AddressInput     `json:"defaultShippingAddress"`
	FirstName              *string           `json:"firstName"`
	LastName               *string           `json:"lastName"`
	Email                  *string           `json:"email"`
	IsActive               *bool             `json:"isActive"`
	Note                   *string           `json:"note"`
	LanguageCode           *LanguageCodeEnum `json:"languageCode"`
	RedirectURL            *string           `json:"redirectUrl"`
	Channel                *string           `json:"channel"`
}

type UserPermission struct {
	Code                   PermissionEnum `json:"code"`
	Name                   string         `json:"name"`
	SourcePermissionGroups []*Group       `json:"sourcePermissionGroups"`
}

type UserSortingInput struct {
	Direction OrderDirection `json:"direction"`
	Field     UserSortField  `json:"field"`
}

type VariantMediaAssign struct {
	ProductVariant *ProductVariant `json:"productVariant"`
	Media          *ProductMedia   `json:"media"`
	Errors         []*ProductError `json:"errors"`
}

type VariantMediaUnassign struct {
	ProductVariant *ProductVariant `json:"productVariant"`
	Media          *ProductMedia   `json:"media"`
	Errors         []*ProductError `json:"errors"`
}

type VariantPricingInfo struct {
	OnSale                *bool       `json:"onSale"`
	Discount              *TaxedMoney `json:"discount"`
	DiscountLocalCurrency *TaxedMoney `json:"discountLocalCurrency"`
	Price                 *TaxedMoney `json:"price"`
	PriceUndiscounted     *TaxedMoney `json:"priceUndiscounted"`
	PriceLocalCurrency    *TaxedMoney `json:"priceLocalCurrency"`
}

type VerifyToken struct {
	User    *User      `json:"user"`
	IsValid bool       `json:"isValid"`
	Payload JSONString `json:"payload"`
}

type VoucherAddCatalogues struct {
	Voucher *Voucher         `json:"voucher"`
	Errors  []*DiscountError `json:"errors"`
}

type VoucherBulkDelete struct {
	Count  int32            `json:"count"`
	Errors []*DiscountError `json:"errors"`
}

type VoucherChannelListingAddInput struct {
	ChannelID      string           `json:"channelId"`
	DiscountValue  *PositiveDecimal `json:"discountValue"`
	MinAmountSpent *PositiveDecimal `json:"minAmountSpent"`
}

type VoucherChannelListingInput struct {
	AddChannels    []*VoucherChannelListingAddInput `json:"addChannels"`
	RemoveChannels []string                         `json:"removeChannels"`
}

type VoucherChannelListingUpdate struct {
	Voucher *Voucher         `json:"voucher"`
	Errors  []*DiscountError `json:"errors"`
}

type VoucherCountableConnection struct {
	PageInfo   *PageInfo               `json:"pageInfo"`
	Edges      []*VoucherCountableEdge `json:"edges"`
	TotalCount *int32                  `json:"totalCount"`
}

type VoucherCountableEdge struct {
	Node   *Voucher `json:"node"`
	Cursor string   `json:"cursor"`
}

type VoucherCreate struct {
	Errors  []*DiscountError `json:"errors"`
	Voucher *Voucher         `json:"voucher"`
}

type VoucherDelete struct {
	Errors  []*DiscountError `json:"errors"`
	Voucher *Voucher         `json:"voucher"`
}

type VoucherFilterInput struct {
	Status       []*DiscountStatusEnum  `json:"status"`
	TimesUsed    *IntRangeInput         `json:"timesUsed"`
	DiscountType []*VoucherDiscountType `json:"discountType"`
	Started      *DateTimeRangeInput    `json:"started"`
	Search       *string                `json:"search"`
	Metadata     []*MetadataFilter      `json:"metadata"`
}

type VoucherInput struct {
	Type                     *VoucherTypeEnum       `json:"type"`
	Name                     *string                `json:"name"`
	Code                     *string                `json:"code"`
	StartDate                *DateTime              `json:"startDate"`
	EndDate                  *DateTime              `json:"endDate"`
	DiscountValueType        *DiscountValueTypeEnum `json:"discountValueType"`
	Products                 []string               `json:"products"`
	Variants                 []string               `json:"variants"`
	Collections              []string               `json:"collections"`
	Categories               []string               `json:"categories"`
	MinCheckoutItemsQuantity *int32                 `json:"minCheckoutItemsQuantity"`
	Countries                []string               `json:"countries"`
	ApplyOncePerOrder        *bool                  `json:"applyOncePerOrder"`
	ApplyOncePerCustomer     *bool                  `json:"applyOncePerCustomer"`
	UsageLimit               *int32                 `json:"usageLimit"`
}

type VoucherRemoveCatalogues struct {
	Voucher *Voucher         `json:"voucher"`
	Errors  []*DiscountError `json:"errors"`
}

type VoucherSortingInput struct {
	Direction OrderDirection   `json:"direction"`
	Channel   *string          `json:"channel"`
	Field     VoucherSortField `json:"field"`
}

type VoucherTranslatableContent struct {
	ID          string              `json:"id"`
	Name        *string             `json:"name"`
	Translation *VoucherTranslation `json:"translation"`
}

type VoucherTranslate struct {
	Errors  []*TranslationError `json:"errors"`
	Voucher *Voucher            `json:"voucher"`
}

type VoucherTranslation struct {
	ID       string           `json:"id"`
	Name     *string          `json:"name"`
	Language *LanguageDisplay `json:"language"`
}

type VoucherUpdate struct {
	Errors  []*DiscountError `json:"errors"`
	Voucher *Voucher         `json:"voucher"`
}

type WarehouseCountableConnection struct {
	PageInfo   *PageInfo                 `json:"pageInfo"`
	Edges      []*WarehouseCountableEdge `json:"edges"`
	TotalCount *int32                    `json:"totalCount"`
}

type WarehouseCountableEdge struct {
	Node   *Warehouse `json:"node"`
	Cursor string     `json:"cursor"`
}

type WarehouseCreate struct {
	Errors    []*WarehouseError `json:"errors"`
	Warehouse *Warehouse        `json:"warehouse"`
}

type WarehouseCreateInput struct {
	Slug          *string       `json:"slug"`
	CompanyName   *string       `json:"companyName"`
	Email         *string       `json:"email"`
	Name          string        `json:"name"`
	Address       *AddressInput `json:"address"`
	ShippingZones []string      `json:"shippingZones"`
}

type WarehouseDelete struct {
	Errors    []*WarehouseError `json:"errors"`
	Warehouse *Warehouse        `json:"warehouse"`
}

type WarehouseError struct {
	Field   *string            `json:"field"`
	Message *string            `json:"message"`
	Code    WarehouseErrorCode `json:"code"`
}

type WarehouseFilterInput struct {
	ClickAndCollectOption *model.WarehouseClickAndCollectOption `json:"clickAndCollectOption"`
	Search                *string                               `json:"search"`
	Ids                   []string                              `json:"ids"`
	IsPrivate             *bool                                 `json:"isPrivate"`
}

type WarehouseShippingZoneAssign struct {
	Errors    []*WarehouseError `json:"errors"`
	Warehouse *Warehouse        `json:"warehouse"`
}

type WarehouseShippingZoneUnassign struct {
	Errors    []*WarehouseError `json:"errors"`
	Warehouse *Warehouse        `json:"warehouse"`
}

type WarehouseSortingInput struct {
	Direction OrderDirection     `json:"direction"`
	Field     WarehouseSortField `json:"field"`
}

type WarehouseUpdate struct {
	Errors    []*WarehouseError `json:"errors"`
	Warehouse *Warehouse        `json:"warehouse"`
}

type WarehouseUpdateInput struct {
	Slug                  *string                               `json:"slug"`
	Email                 *string                               `json:"email"`
	Name                  *string                               `json:"name"`
	Address               *AddressInput                         `json:"address"`
	ClickAndCollectOption *model.WarehouseClickAndCollectOption `json:"clickAndCollectOption"`
	IsPrivate             *bool                                 `json:"isPrivate"`
}

type Webhook struct {
	Name      string          `json:"name"`
	TargetURL string          `json:"targetUrl"`
	IsActive  bool            `json:"isActive"`
	SecretKey *string         `json:"secretKey"`
	ID        string          `json:"id"`
	Events    []*WebhookEvent `json:"events"`
	App       *App            `json:"app"`
}

type WebhookCreate struct {
	Errors  []*WebhookError `json:"errors"`
	Webhook *Webhook        `json:"webhook"`
}

type WebhookCreateInput struct {
	Name      *string                 `json:"name"`
	TargetURL *string                 `json:"targetUrl"`
	Events    []*WebhookEventTypeEnum `json:"events"`
	App       *string                 `json:"app"`
	IsActive  *bool                   `json:"isActive"`
	SecretKey *string                 `json:"secretKey"`
}

type WebhookDelete struct {
	Errors  []*WebhookError `json:"errors"`
	Webhook *Webhook        `json:"webhook"`
}

type WebhookError struct {
	Field   *string          `json:"field"`
	Message *string          `json:"message"`
	Code    WebhookErrorCode `json:"code"`
}

type WebhookEvent struct {
	EventType WebhookEventTypeEnum `json:"eventType"`
	Name      string               `json:"name"`
}

type WebhookUpdate struct {
	Errors  []*WebhookError `json:"errors"`
	Webhook *Webhook        `json:"webhook"`
}

type WebhookUpdateInput struct {
	Name      *string                 `json:"name"`
	TargetURL *string                 `json:"targetUrl"`
	Events    []*WebhookEventTypeEnum `json:"events"`
	App       *string                 `json:"app"`
	IsActive  *bool                   `json:"isActive"`
	SecretKey *string                 `json:"secretKey"`
}

type Weight struct {
	Unit  WeightUnitsEnum `json:"unit"`
	Value float64         `json:"value"`
}

type Wishlist struct {
	ID       string          `json:"id"`
	Token    string          `json:"token"`
	CreateAt DateTime        `json:"createAt"`
	Items    []*WishlistItem `json:"items"`
}

type WishlistItem struct {
	ID       string            `json:"id"`
	Product  *Product          `json:"product"`
	CreateAt DateTime          `json:"createAt"`
	Variants []*ProductVariant `json:"variants"`
}

type AccountErrorCode string

const (
	AccountErrorCodeActivateOwnAccount          AccountErrorCode = "ACTIVATE_OWN_ACCOUNT"
	AccountErrorCodeActivateSuperuserAccount    AccountErrorCode = "ACTIVATE_SUPERUSER_ACCOUNT"
	AccountErrorCodeDuplicatedInputItem         AccountErrorCode = "DUPLICATED_INPUT_ITEM"
	AccountErrorCodeDeactivateOwnAccount        AccountErrorCode = "DEACTIVATE_OWN_ACCOUNT"
	AccountErrorCodeDeactivateSuperuserAccount  AccountErrorCode = "DEACTIVATE_SUPERUSER_ACCOUNT"
	AccountErrorCodeDeleteNonStaffUser          AccountErrorCode = "DELETE_NON_STAFF_USER"
	AccountErrorCodeDeleteOwnAccount            AccountErrorCode = "DELETE_OWN_ACCOUNT"
	AccountErrorCodeDeleteStaffAccount          AccountErrorCode = "DELETE_STAFF_ACCOUNT"
	AccountErrorCodeDeleteSuperuserAccount      AccountErrorCode = "DELETE_SUPERUSER_ACCOUNT"
	AccountErrorCodeGraphqlError                AccountErrorCode = "GRAPHQL_ERROR"
	AccountErrorCodeInactive                    AccountErrorCode = "INACTIVE"
	AccountErrorCodeInvalid                     AccountErrorCode = "INVALID"
	AccountErrorCodeInvalidPassword             AccountErrorCode = "INVALID_PASSWORD"
	AccountErrorCodeLeftNotManageablePermission AccountErrorCode = "LEFT_NOT_MANAGEABLE_PERMISSION"
	AccountErrorCodeInvalidCredentials          AccountErrorCode = "INVALID_CREDENTIALS"
	AccountErrorCodeNotFound                    AccountErrorCode = "NOT_FOUND"
	AccountErrorCodeOutOfScopeUser              AccountErrorCode = "OUT_OF_SCOPE_USER"
	AccountErrorCodeOutOfScopeGroup             AccountErrorCode = "OUT_OF_SCOPE_GROUP"
	AccountErrorCodeOutOfScopePermission        AccountErrorCode = "OUT_OF_SCOPE_PERMISSION"
	AccountErrorCodePasswordEntirelyNumeric     AccountErrorCode = "PASSWORD_ENTIRELY_NUMERIC"
	AccountErrorCodePasswordTooCommon           AccountErrorCode = "PASSWORD_TOO_COMMON"
	AccountErrorCodePasswordTooShort            AccountErrorCode = "PASSWORD_TOO_SHORT"
	AccountErrorCodePasswordTooSimilar          AccountErrorCode = "PASSWORD_TOO_SIMILAR"
	AccountErrorCodeRequired                    AccountErrorCode = "REQUIRED"
	AccountErrorCodeUnique                      AccountErrorCode = "UNIQUE"
	AccountErrorCodeJwtSignatureExpired         AccountErrorCode = "JWT_SIGNATURE_EXPIRED"
	AccountErrorCodeJwtInvalidToken             AccountErrorCode = "JWT_INVALID_TOKEN"
	AccountErrorCodeJwtDecodeError              AccountErrorCode = "JWT_DECODE_ERROR"
	AccountErrorCodeJwtMissingToken             AccountErrorCode = "JWT_MISSING_TOKEN"
	AccountErrorCodeJwtInvalidCsrfToken         AccountErrorCode = "JWT_INVALID_CSRF_TOKEN"
	AccountErrorCodeChannelInactive             AccountErrorCode = "CHANNEL_INACTIVE"
	AccountErrorCodeMissingChannelSlug          AccountErrorCode = "MISSING_CHANNEL_SLUG"
)

func (e AccountErrorCode) IsValid() bool {
	switch e {
	case AccountErrorCodeActivateOwnAccount, AccountErrorCodeActivateSuperuserAccount, AccountErrorCodeDuplicatedInputItem, AccountErrorCodeDeactivateOwnAccount, AccountErrorCodeDeactivateSuperuserAccount, AccountErrorCodeDeleteNonStaffUser, AccountErrorCodeDeleteOwnAccount, AccountErrorCodeDeleteStaffAccount, AccountErrorCodeDeleteSuperuserAccount, AccountErrorCodeGraphqlError, AccountErrorCodeInactive, AccountErrorCodeInvalid, AccountErrorCodeInvalidPassword, AccountErrorCodeLeftNotManageablePermission, AccountErrorCodeInvalidCredentials, AccountErrorCodeNotFound, AccountErrorCodeOutOfScopeUser, AccountErrorCodeOutOfScopeGroup, AccountErrorCodeOutOfScopePermission, AccountErrorCodePasswordEntirelyNumeric, AccountErrorCodePasswordTooCommon, AccountErrorCodePasswordTooShort, AccountErrorCodePasswordTooSimilar, AccountErrorCodeRequired, AccountErrorCodeUnique, AccountErrorCodeJwtSignatureExpired, AccountErrorCodeJwtInvalidToken, AccountErrorCodeJwtDecodeError, AccountErrorCodeJwtMissingToken, AccountErrorCodeJwtInvalidCsrfToken, AccountErrorCodeChannelInactive, AccountErrorCodeMissingChannelSlug:
		return true
	}
	return false
}

type AppErrorCode string

const (
	AppErrorCodeForbidden              AppErrorCode = "FORBIDDEN"
	AppErrorCodeGraphqlError           AppErrorCode = "GRAPHQL_ERROR"
	AppErrorCodeInvalid                AppErrorCode = "INVALID"
	AppErrorCodeInvalidStatus          AppErrorCode = "INVALID_STATUS"
	AppErrorCodeInvalidPermission      AppErrorCode = "INVALID_PERMISSION"
	AppErrorCodeInvalidURLFormat       AppErrorCode = "INVALID_URL_FORMAT"
	AppErrorCodeInvalidManifestFormat  AppErrorCode = "INVALID_MANIFEST_FORMAT"
	AppErrorCodeManifestURLCantConnect AppErrorCode = "MANIFEST_URL_CANT_CONNECT"
	AppErrorCodeNotFound               AppErrorCode = "NOT_FOUND"
	AppErrorCodeRequired               AppErrorCode = "REQUIRED"
	AppErrorCodeUnique                 AppErrorCode = "UNIQUE"
	AppErrorCodeOutOfScopeApp          AppErrorCode = "OUT_OF_SCOPE_APP"
	AppErrorCodeOutOfScopePermission   AppErrorCode = "OUT_OF_SCOPE_PERMISSION"
)

func (e AppErrorCode) IsValid() bool {
	switch e {
	case AppErrorCodeForbidden, AppErrorCodeGraphqlError, AppErrorCodeInvalid, AppErrorCodeInvalidStatus, AppErrorCodeInvalidPermission, AppErrorCodeInvalidURLFormat, AppErrorCodeInvalidManifestFormat, AppErrorCodeManifestURLCantConnect, AppErrorCodeNotFound, AppErrorCodeRequired, AppErrorCodeUnique, AppErrorCodeOutOfScopeApp, AppErrorCodeOutOfScopePermission:
		return true
	}
	return false
}

type AppExtensionTargetEnum string

const (
	AppExtensionTargetEnumMoreActions AppExtensionTargetEnum = "MORE_ACTIONS"
	AppExtensionTargetEnumCreate      AppExtensionTargetEnum = "CREATE"
)

func (e AppExtensionTargetEnum) IsValid() bool {
	switch e {
	case AppExtensionTargetEnumMoreActions, AppExtensionTargetEnumCreate:
		return true
	}
	return false
}

type AppExtensionTypeEnum string

const (
	AppExtensionTypeEnumOverview AppExtensionTypeEnum = "OVERVIEW"
	AppExtensionTypeEnumDetails  AppExtensionTypeEnum = "DETAILS"
)

func (e AppExtensionTypeEnum) IsValid() bool {
	switch e {
	case AppExtensionTypeEnumOverview, AppExtensionTypeEnumDetails:
		return true
	}
	return false
}

type AppExtensionViewEnum string

const (
	AppExtensionViewEnumProduct AppExtensionViewEnum = "PRODUCT"
)

func (e AppExtensionViewEnum) IsValid() bool {
	switch e {
	case AppExtensionViewEnumProduct:
		return true
	}
	return false
}

type AppSortField string

const (
	AppSortFieldName         AppSortField = "NAME"
	AppSortFieldCreationDate AppSortField = "CREATION_DATE"
)

func (e AppSortField) IsValid() bool {
	switch e {
	case AppSortFieldName, AppSortFieldCreationDate:
		return true
	}
	return false
}

type AppTypeEnum string

const (
	AppTypeEnumLocal      AppTypeEnum = "LOCAL"
	AppTypeEnumThirdparty AppTypeEnum = "THIRDPARTY"
)

func (e AppTypeEnum) IsValid() bool {
	switch e {
	case AppTypeEnumLocal, AppTypeEnumThirdparty:
		return true
	}
	return false
}

type AreaUnitsEnum string

const (
	AreaUnitsEnumSqCm   AreaUnitsEnum = measurement.SQ_CM
	AreaUnitsEnumSqM    AreaUnitsEnum = measurement.SQ_M
	AreaUnitsEnumSqKm   AreaUnitsEnum = measurement.SQ_KM
	AreaUnitsEnumSqFt   AreaUnitsEnum = measurement.SQ_FT
	AreaUnitsEnumSqYd   AreaUnitsEnum = measurement.SQ_YD
	AreaUnitsEnumSqInch AreaUnitsEnum = measurement.SQ_INCH
)

func (e AreaUnitsEnum) IsValid() bool {
	switch e {
	case AreaUnitsEnumSqCm, AreaUnitsEnumSqM, AreaUnitsEnumSqKm, AreaUnitsEnumSqFt, AreaUnitsEnumSqYd, AreaUnitsEnumSqInch:
		return true
	}
	return false
}

type AttributeChoicesSortField string

const (
	AttributeChoicesSortFieldName AttributeChoicesSortField = "NAME"
	AttributeChoicesSortFieldSlug AttributeChoicesSortField = "SLUG"
)

func (e AttributeChoicesSortField) IsValid() bool {
	switch e {
	case AttributeChoicesSortFieldName, AttributeChoicesSortFieldSlug:
		return true
	}
	return false
}

type AttributeErrorCode string

const (
	AttributeErrorCodeAlreadyExists AttributeErrorCode = "already_exists"
	AttributeErrorCodeGraphqlError  AttributeErrorCode = "graphql_error"
	AttributeErrorCodeInvalid       AttributeErrorCode = "invalid"
	AttributeErrorCodeNotFound      AttributeErrorCode = "not_found"
	AttributeErrorCodeRequired      AttributeErrorCode = "required"
	AttributeErrorCodeUnique        AttributeErrorCode = "unique"
)

func (e AttributeErrorCode) IsValid() bool {
	switch e {
	case AttributeErrorCodeAlreadyExists, AttributeErrorCodeGraphqlError, AttributeErrorCodeInvalid, AttributeErrorCodeNotFound, AttributeErrorCodeRequired, AttributeErrorCodeUnique:
		return true
	}
	return false
}

type AttributeSortField string

const (
	AttributeSortFieldName                     AttributeSortField = "NAME"
	AttributeSortFieldSlug                     AttributeSortField = "SLUG"
	AttributeSortFieldValueRequired            AttributeSortField = "VALUE_REQUIRED"
	AttributeSortFieldIsVariantOnly            AttributeSortField = "IS_VARIANT_ONLY"
	AttributeSortFieldVisibleInStorefront      AttributeSortField = "VISIBLE_IN_STOREFRONT"
	AttributeSortFieldFilterableInStorefront   AttributeSortField = "FILTERABLE_IN_STOREFRONT"
	AttributeSortFieldFilterableInDashboard    AttributeSortField = "FILTERABLE_IN_DASHBOARD"
	AttributeSortFieldStorefrontSearchPosition AttributeSortField = "STOREFRONT_SEARCH_POSITION"
	AttributeSortFieldAvailableInGrid          AttributeSortField = "AVAILABLE_IN_GRID"
)

func (e AttributeSortField) IsValid() bool {
	switch e {
	case AttributeSortFieldName, AttributeSortFieldSlug, AttributeSortFieldValueRequired, AttributeSortFieldIsVariantOnly, AttributeSortFieldVisibleInStorefront, AttributeSortFieldFilterableInStorefront, AttributeSortFieldFilterableInDashboard, AttributeSortFieldStorefrontSearchPosition, AttributeSortFieldAvailableInGrid:
		return true
	}
	return false
}

type CategorySortField string

const (
	CategorySortFieldName             CategorySortField = "NAME"
	CategorySortFieldProductCount     CategorySortField = "PRODUCT_COUNT"
	CategorySortFieldSubcategoryCount CategorySortField = "SUBCATEGORY_COUNT"
)

func (e CategorySortField) IsValid() bool {
	switch e {
	case CategorySortFieldName, CategorySortFieldProductCount, CategorySortFieldSubcategoryCount:
		return true
	}
	return false
}

type ChannelErrorCode string

const (
	ChannelErrorCodeAlreadyExists                 ChannelErrorCode = "ALREADY_EXISTS"
	ChannelErrorCodeGraphqlError                  ChannelErrorCode = "GRAPHQL_ERROR"
	ChannelErrorCodeInvalid                       ChannelErrorCode = "INVALID"
	ChannelErrorCodeNotFound                      ChannelErrorCode = "NOT_FOUND"
	ChannelErrorCodeRequired                      ChannelErrorCode = "REQUIRED"
	ChannelErrorCodeUnique                        ChannelErrorCode = "UNIQUE"
	ChannelErrorCodeChannelsCurrencyMustBeTheSame ChannelErrorCode = "CHANNELS_CURRENCY_MUST_BE_THE_SAME"
	ChannelErrorCodeChannelWithOrders             ChannelErrorCode = "CHANNEL_WITH_ORDERS"
	ChannelErrorCodeDuplicatedInputItem           ChannelErrorCode = "DUPLICATED_INPUT_ITEM"
)

func (e ChannelErrorCode) IsValid() bool {
	switch e {
	case ChannelErrorCodeAlreadyExists, ChannelErrorCodeGraphqlError, ChannelErrorCodeInvalid, ChannelErrorCodeNotFound, ChannelErrorCodeRequired, ChannelErrorCodeUnique, ChannelErrorCodeChannelsCurrencyMustBeTheSame, ChannelErrorCodeChannelWithOrders, ChannelErrorCodeDuplicatedInputItem:
		return true
	}
	return false
}

type CheckoutErrorCode string

const (
	CheckoutErrorCodeBillingAddressNotSet          CheckoutErrorCode = "BILLING_ADDRESS_NOT_SET"
	CheckoutErrorCodeCheckoutNotFullyPaid          CheckoutErrorCode = "CHECKOUT_NOT_FULLY_PAID"
	CheckoutErrorCodeGraphqlError                  CheckoutErrorCode = "GRAPHQL_ERROR"
	CheckoutErrorCodeProductNotPublished           CheckoutErrorCode = "PRODUCT_NOT_PUBLISHED"
	CheckoutErrorCodeProductUnavailableForPurchase CheckoutErrorCode = "PRODUCT_UNAVAILABLE_FOR_PURCHASE"
	CheckoutErrorCodeInsufficientStock             CheckoutErrorCode = "INSUFFICIENT_STOCK"
	CheckoutErrorCodeInvalid                       CheckoutErrorCode = "INVALID"
	CheckoutErrorCodeInvalidShippingMethod         CheckoutErrorCode = "INVALID_SHIPPING_METHOD"
	CheckoutErrorCodeNotFound                      CheckoutErrorCode = "NOT_FOUND"
	CheckoutErrorCodePaymentError                  CheckoutErrorCode = "PAYMENT_ERROR"
	CheckoutErrorCodeQuantityGreaterThanLimit      CheckoutErrorCode = "QUANTITY_GREATER_THAN_LIMIT"
	CheckoutErrorCodeRequired                      CheckoutErrorCode = "REQUIRED"
	CheckoutErrorCodeShippingAddressNotSet         CheckoutErrorCode = "SHIPPING_ADDRESS_NOT_SET"
	CheckoutErrorCodeShippingMethodNotApplicable   CheckoutErrorCode = "SHIPPING_METHOD_NOT_APPLICABLE"
	CheckoutErrorCodeDeliveryMethodNotApplicable   CheckoutErrorCode = "DELIVERY_METHOD_NOT_APPLICABLE"
	CheckoutErrorCodeShippingMethodNotSet          CheckoutErrorCode = "SHIPPING_METHOD_NOT_SET"
	CheckoutErrorCodeShippingNotRequired           CheckoutErrorCode = "SHIPPING_NOT_REQUIRED"
	CheckoutErrorCodeTaxError                      CheckoutErrorCode = "TAX_ERROR"
	CheckoutErrorCodeUnique                        CheckoutErrorCode = "UNIQUE"
	CheckoutErrorCodeVoucherNotApplicable          CheckoutErrorCode = "VOUCHER_NOT_APPLICABLE"
	CheckoutErrorCodeGiftCardNotApplicable         CheckoutErrorCode = "GIFT_CARD_NOT_APPLICABLE"
	CheckoutErrorCodeZeroQuantity                  CheckoutErrorCode = "ZERO_QUANTITY"
	CheckoutErrorCodeMissingChannelSlug            CheckoutErrorCode = "MISSING_CHANNEL_SLUG"
	CheckoutErrorCodeChannelInactive               CheckoutErrorCode = "CHANNEL_INACTIVE"
	CheckoutErrorCodeUnavailableVariantInChannel   CheckoutErrorCode = "UNAVAILABLE_VARIANT_IN_CHANNEL"
)

func (e CheckoutErrorCode) IsValid() bool {
	switch e {
	case CheckoutErrorCodeBillingAddressNotSet, CheckoutErrorCodeCheckoutNotFullyPaid, CheckoutErrorCodeGraphqlError, CheckoutErrorCodeProductNotPublished, CheckoutErrorCodeProductUnavailableForPurchase, CheckoutErrorCodeInsufficientStock, CheckoutErrorCodeInvalid, CheckoutErrorCodeInvalidShippingMethod, CheckoutErrorCodeNotFound, CheckoutErrorCodePaymentError, CheckoutErrorCodeQuantityGreaterThanLimit, CheckoutErrorCodeRequired, CheckoutErrorCodeShippingAddressNotSet, CheckoutErrorCodeShippingMethodNotApplicable, CheckoutErrorCodeDeliveryMethodNotApplicable, CheckoutErrorCodeShippingMethodNotSet, CheckoutErrorCodeShippingNotRequired, CheckoutErrorCodeTaxError, CheckoutErrorCodeUnique, CheckoutErrorCodeVoucherNotApplicable, CheckoutErrorCodeGiftCardNotApplicable, CheckoutErrorCodeZeroQuantity, CheckoutErrorCodeMissingChannelSlug, CheckoutErrorCodeChannelInactive, CheckoutErrorCodeUnavailableVariantInChannel:
		return true
	}
	return false
}

type CollectionErrorCode string

const (
	CollectionErrorCodeDuplicatedInputItem               CollectionErrorCode = "DUPLICATED_INPUT_ITEM"
	CollectionErrorCodeGraphqlError                      CollectionErrorCode = "GRAPHQL_ERROR"
	CollectionErrorCodeInvalid                           CollectionErrorCode = "INVALID"
	CollectionErrorCodeNotFound                          CollectionErrorCode = "NOT_FOUND"
	CollectionErrorCodeRequired                          CollectionErrorCode = "REQUIRED"
	CollectionErrorCodeUnique                            CollectionErrorCode = "UNIQUE"
	CollectionErrorCodeCannotManageProductWithoutVariant CollectionErrorCode = "CANNOT_MANAGE_PRODUCT_WITHOUT_VARIANT"
)

func (e CollectionErrorCode) IsValid() bool {
	switch e {
	case CollectionErrorCodeDuplicatedInputItem, CollectionErrorCodeGraphqlError, CollectionErrorCodeInvalid, CollectionErrorCodeNotFound, CollectionErrorCodeRequired, CollectionErrorCodeUnique, CollectionErrorCodeCannotManageProductWithoutVariant:
		return true
	}
	return false
}

type CollectionPublished string

const (
	CollectionPublishedPublished CollectionPublished = "PUBLISHED"
	CollectionPublishedHidden    CollectionPublished = "HIDDEN"
)

func (e CollectionPublished) IsValid() bool {
	switch e {
	case CollectionPublishedPublished, CollectionPublishedHidden:
		return true
	}
	return false
}

type CollectionSortField string

const (
	CollectionSortFieldName            CollectionSortField = "NAME"
	CollectionSortFieldAvailability    CollectionSortField = "AVAILABILITY"
	CollectionSortFieldProductCount    CollectionSortField = "PRODUCT_COUNT"
	CollectionSortFieldPublicationDate CollectionSortField = "PUBLICATION_DATE"
)

func (e CollectionSortField) IsValid() bool {
	switch e {
	case CollectionSortFieldName, CollectionSortFieldAvailability, CollectionSortFieldProductCount, CollectionSortFieldPublicationDate:
		return true
	}
	return false
}

type ConfigurationTypeFieldEnum string

const (
	ConfigurationTypeFieldEnumString          ConfigurationTypeFieldEnum = "STRING"
	ConfigurationTypeFieldEnumMultiline       ConfigurationTypeFieldEnum = "MULTILINE"
	ConfigurationTypeFieldEnumBoolean         ConfigurationTypeFieldEnum = "BOOLEAN"
	ConfigurationTypeFieldEnumSecret          ConfigurationTypeFieldEnum = "SECRET"
	ConfigurationTypeFieldEnumPassword        ConfigurationTypeFieldEnum = "PASSWORD"
	ConfigurationTypeFieldEnumSecretmultiline ConfigurationTypeFieldEnum = "SECRETMULTILINE"
	ConfigurationTypeFieldEnumOutput          ConfigurationTypeFieldEnum = "OUTPUT"
)

func (e ConfigurationTypeFieldEnum) IsValid() bool {
	switch e {
	case ConfigurationTypeFieldEnumString, ConfigurationTypeFieldEnumMultiline, ConfigurationTypeFieldEnumBoolean, ConfigurationTypeFieldEnumSecret, ConfigurationTypeFieldEnumPassword, ConfigurationTypeFieldEnumSecretmultiline, ConfigurationTypeFieldEnumOutput:
		return true
	}
	return false
}

type CountryCode = model.CountryCode

type CustomerEventsEnum string

const (
	CustomerEventsEnumAccountCreated        CustomerEventsEnum = "ACCOUNT_CREATED"
	CustomerEventsEnumPasswordResetLinkSent CustomerEventsEnum = "PASSWORD_RESET_LINK_SENT"
	CustomerEventsEnumPasswordReset         CustomerEventsEnum = "PASSWORD_RESET"
	CustomerEventsEnumEmailChangedRequest   CustomerEventsEnum = "EMAIL_CHANGED_REQUEST"
	CustomerEventsEnumPasswordChanged       CustomerEventsEnum = "PASSWORD_CHANGED"
	CustomerEventsEnumEmailChanged          CustomerEventsEnum = "EMAIL_CHANGED"
	CustomerEventsEnumPlacedOrder           CustomerEventsEnum = "PLACED_ORDER"
	CustomerEventsEnumNoteAddedToOrder      CustomerEventsEnum = "NOTE_ADDED_TO_ORDER"
	CustomerEventsEnumDigitalLinkDownloaded CustomerEventsEnum = "DIGITAL_LINK_DOWNLOADED"
	CustomerEventsEnumCustomerDeleted       CustomerEventsEnum = "CUSTOMER_DELETED"
	CustomerEventsEnumNameAssigned          CustomerEventsEnum = "NAME_ASSIGNED"
	CustomerEventsEnumEmailAssigned         CustomerEventsEnum = "EMAIL_ASSIGNED"
	CustomerEventsEnumNoteAdded             CustomerEventsEnum = "NOTE_ADDED"
)

func (e CustomerEventsEnum) IsValid() bool {
	switch e {
	case CustomerEventsEnumAccountCreated, CustomerEventsEnumPasswordResetLinkSent, CustomerEventsEnumPasswordReset, CustomerEventsEnumEmailChangedRequest, CustomerEventsEnumPasswordChanged, CustomerEventsEnumEmailChanged, CustomerEventsEnumPlacedOrder, CustomerEventsEnumNoteAddedToOrder, CustomerEventsEnumDigitalLinkDownloaded, CustomerEventsEnumCustomerDeleted, CustomerEventsEnumNameAssigned, CustomerEventsEnumEmailAssigned, CustomerEventsEnumNoteAdded:
		return true
	}
	return false
}

type DiscountErrorCode string

const (
	DiscountErrorCodeAlreadyExists                     DiscountErrorCode = "ALREADY_EXISTS"
	DiscountErrorCodeGraphqlError                      DiscountErrorCode = "GRAPHQL_ERROR"
	DiscountErrorCodeInvalid                           DiscountErrorCode = "INVALID"
	DiscountErrorCodeNotFound                          DiscountErrorCode = "NOT_FOUND"
	DiscountErrorCodeRequired                          DiscountErrorCode = "REQUIRED"
	DiscountErrorCodeUnique                            DiscountErrorCode = "UNIQUE"
	DiscountErrorCodeCannotManageProductWithoutVariant DiscountErrorCode = "CANNOT_MANAGE_PRODUCT_WITHOUT_VARIANT"
	DiscountErrorCodeDuplicatedInputItem               DiscountErrorCode = "DUPLICATED_INPUT_ITEM"
)

func (e DiscountErrorCode) IsValid() bool {
	switch e {
	case DiscountErrorCodeAlreadyExists, DiscountErrorCodeGraphqlError, DiscountErrorCodeInvalid, DiscountErrorCodeNotFound, DiscountErrorCodeRequired, DiscountErrorCodeUnique, DiscountErrorCodeCannotManageProductWithoutVariant, DiscountErrorCodeDuplicatedInputItem:
		return true
	}
	return false
}

type DiscountStatusEnum string

const (
	DiscountStatusEnumActive    DiscountStatusEnum = "ACTIVE"
	DiscountStatusEnumExpired   DiscountStatusEnum = "EXPIRED"
	DiscountStatusEnumScheduled DiscountStatusEnum = "SCHEDULED"
)

func (e DiscountStatusEnum) IsValid() bool {
	switch e {
	case DiscountStatusEnumActive, DiscountStatusEnumExpired, DiscountStatusEnumScheduled:
		return true
	}
	return false
}

type DiscountValueTypeEnum = model.DiscountType

type DistanceUnitsEnum = measurement.DistanceUnit

type ExportEventsEnum = model.ExportEventType

type ExportFileSortField string

const (
	ExportFileSortFieldStatus    ExportFileSortField = "STATUS"
	ExportFileSortFieldCreatedAt ExportFileSortField = "CREATED_AT"
	ExportFileSortFieldUpdatedAt ExportFileSortField = "UPDATED_AT"
)

func (e ExportFileSortField) IsValid() bool {
	switch e {
	case ExportFileSortFieldStatus, ExportFileSortFieldCreatedAt, ExportFileSortFieldUpdatedAt:
		return true
	}
	return false
}

type ExportScope string

const (
	ExportScopeAll    ExportScope = "ALL"
	ExportScopeIDS    ExportScope = "IDS"
	ExportScopeFilter ExportScope = "FILTER"
)

func (e ExportScope) IsValid() bool {
	switch e {
	case ExportScopeAll, ExportScopeIDS, ExportScopeFilter:
		return true
	}
	return false
}

type ExternalNotificationErrorCodes string

const (
	ExternalNotificationErrorCodesRequired         ExternalNotificationErrorCodes = "REQUIRED"
	ExternalNotificationErrorCodesInvalidModelType ExternalNotificationErrorCodes = "INVALID_MODEL_TYPE"
	ExternalNotificationErrorCodesNotFound         ExternalNotificationErrorCodes = "NOT_FOUND"
	ExternalNotificationErrorCodesChannelInactive  ExternalNotificationErrorCodes = "CHANNEL_INACTIVE"
)

func (e ExternalNotificationErrorCodes) IsValid() bool {
	switch e {
	case ExternalNotificationErrorCodesRequired, ExternalNotificationErrorCodesInvalidModelType, ExternalNotificationErrorCodesNotFound, ExternalNotificationErrorCodesChannelInactive:
		return true
	}
	return false
}

type FileTypesEnum string

const (
	FileTypesEnumCSV  FileTypesEnum = "csv"
	FileTypesEnumXlsx FileTypesEnum = "xlsx"
)

func (e FileTypesEnum) IsValid() bool {
	switch e {
	case FileTypesEnumCSV, FileTypesEnumXlsx:
		return true
	}
	return false
}

type FulfillmentStatus = model.FulfillmentStatus

type GiftCardErrorCode string

const (
	GiftCardErrorCodeAlreadyExists GiftCardErrorCode = "ALREADY_EXISTS"
	GiftCardErrorCodeGraphqlError  GiftCardErrorCode = "GRAPHQL_ERROR"
	GiftCardErrorCodeInvalid       GiftCardErrorCode = "INVALID"
	GiftCardErrorCodeNotFound      GiftCardErrorCode = "NOT_FOUND"
	GiftCardErrorCodeRequired      GiftCardErrorCode = "REQUIRED"
	GiftCardErrorCodeUnique        GiftCardErrorCode = "UNIQUE"
)

func (e GiftCardErrorCode) IsValid() bool {
	switch e {
	case GiftCardErrorCodeAlreadyExists, GiftCardErrorCodeGraphqlError, GiftCardErrorCodeInvalid, GiftCardErrorCodeNotFound, GiftCardErrorCodeRequired, GiftCardErrorCodeUnique:
		return true
	}
	return false
}

type GiftCardEventsEnum = model.GiftcardEventType

type GiftCardSettingsErrorCode string

const (
	GiftCardSettingsErrorCodeInvalid      GiftCardSettingsErrorCode = "INVALID"
	GiftCardSettingsErrorCodeRequired     GiftCardSettingsErrorCode = "REQUIRED"
	GiftCardSettingsErrorCodeGraphqlError GiftCardSettingsErrorCode = "GRAPHQL_ERROR"
)

func (e GiftCardSettingsErrorCode) IsValid() bool {
	switch e {
	case GiftCardSettingsErrorCodeInvalid, GiftCardSettingsErrorCodeRequired, GiftCardSettingsErrorCodeGraphqlError:
		return true
	}
	return false
}

type GiftCardSettingsExpiryTypeEnum = model.GiftCardSettingsExpiryType

type GiftCardSortField string

const (
	GiftCardSortFieldTag GiftCardSortField = "TAG"
	// GiftCardSortFieldProduct        GiftCardSortField = "PRODUCT"
	// GiftCardSortFieldUsedBy         GiftCardSortField = "USED_BY"
	GiftCardSortFieldCurrentBalance GiftCardSortField = "CURRENT_BALANCE"
)

func (e GiftCardSortField) IsValid() bool {
	switch e {
	case GiftCardSortFieldTag, GiftCardSortFieldCurrentBalance:
		return true
	}
	return false
}

type InvoiceErrorCode string

const (
	InvoiceErrorCodeRequired      InvoiceErrorCode = "REQUIRED"
	InvoiceErrorCodeNotReady      InvoiceErrorCode = "NOT_READY"
	InvoiceErrorCodeURLNotSet     InvoiceErrorCode = "URL_NOT_SET"
	InvoiceErrorCodeEmailNotSet   InvoiceErrorCode = "EMAIL_NOT_SET"
	InvoiceErrorCodeNumberNotSet  InvoiceErrorCode = "NUMBER_NOT_SET"
	InvoiceErrorCodeNotFound      InvoiceErrorCode = "NOT_FOUND"
	InvoiceErrorCodeInvalidStatus InvoiceErrorCode = "INVALID_STATUS"
)

func (e InvoiceErrorCode) IsValid() bool {
	switch e {
	case InvoiceErrorCodeRequired, InvoiceErrorCodeNotReady, InvoiceErrorCodeURLNotSet, InvoiceErrorCodeEmailNotSet, InvoiceErrorCodeNumberNotSet, InvoiceErrorCodeNotFound, InvoiceErrorCodeInvalidStatus:
		return true
	}
	return false
}

type JobStatusEnum string

const (
	JobStatusEnumPending JobStatusEnum = "PENDING"
	JobStatusEnumSuccess JobStatusEnum = "SUCCESS"
	JobStatusEnumFailed  JobStatusEnum = "FAILED"
	JobStatusEnumDeleted JobStatusEnum = "DELETED"
)

func (e JobStatusEnum) IsValid() bool {
	switch e {
	case JobStatusEnumPending, JobStatusEnumSuccess, JobStatusEnumFailed, JobStatusEnumDeleted:
		return true
	}
	return false
}

type LanguageCodeEnum = model.LanguageCodeEnum

type LoginErrorCode string

const (
	LoginErrorCodeGraphqlError LoginErrorCode = "GRAPHQL_ERROR"
	LoginErrorCodeInvalid      LoginErrorCode = "INVALID"
	LoginErrorCodeNotFound     LoginErrorCode = "NOT_FOUND"
	LoginErrorCodeRequired     LoginErrorCode = "REQUIRED"
	LoginErrorCodeUnique       LoginErrorCode = "UNIQUE"
)

func (e LoginErrorCode) IsValid() bool {
	switch e {
	case LoginErrorCodeGraphqlError, LoginErrorCodeInvalid, LoginErrorCodeNotFound, LoginErrorCodeRequired, LoginErrorCodeUnique:
		return true
	}
	return false
}

type MeasurementUnitsEnum string

const (
	MeasurementUnitsEnumCm              MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.CM)
	MeasurementUnitsEnumM               MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.M)
	MeasurementUnitsEnumKm              MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.KM)
	MeasurementUnitsEnumFt              MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.FT)
	MeasurementUnitsEnumYd              MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.YD)
	MeasurementUnitsEnumInch            MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.INCH)
	MeasurementUnitsEnumSqCm            MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.SQ_CM)
	MeasurementUnitsEnumSqM             MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.SQ_M)
	MeasurementUnitsEnumSqKm            MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.SQ_KM)
	MeasurementUnitsEnumSqFt            MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.SQ_FT)
	MeasurementUnitsEnumSqYd            MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.SQ_YD)
	MeasurementUnitsEnumSqInch          MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.SQ_INCH)
	MeasurementUnitsEnumCubicMillimeter MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.CUBIC_MILLIMETER)
	MeasurementUnitsEnumCubicCentimeter MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.CUBIC_CENTIMETER)
	MeasurementUnitsEnumCubicDecimeter  MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.CUBIC_DECIMETER)
	MeasurementUnitsEnumCubicMeter      MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.CUBIC_METER)
	MeasurementUnitsEnumLiter           MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.LITER)
	MeasurementUnitsEnumCubicFoot       MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.CUBIC_FOOT)
	MeasurementUnitsEnumCubicInch       MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.CUBIC_INCH)
	MeasurementUnitsEnumCubicYard       MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.CUBIC_YARD)
	MeasurementUnitsEnumQt              MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.QT)
	MeasurementUnitsEnumPint            MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.PINT)
	MeasurementUnitsEnumFlOz            MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.FL_OZ)
	MeasurementUnitsEnumAcreIn          MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.ACRE_IN)
	MeasurementUnitsEnumAcreFt          MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.ACRE_FT)
	MeasurementUnitsEnumG               MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.G)
	MeasurementUnitsEnumLb              MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.LB)
	MeasurementUnitsEnumOz              MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.OZ)
	MeasurementUnitsEnumKg              MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.KG)
	MeasurementUnitsEnumTonne           MeasurementUnitsEnum = MeasurementUnitsEnum(measurement.TONNE)
)

func (e MeasurementUnitsEnum) IsValid() bool {
	switch e {
	case MeasurementUnitsEnumCm, MeasurementUnitsEnumM, MeasurementUnitsEnumKm, MeasurementUnitsEnumFt, MeasurementUnitsEnumYd, MeasurementUnitsEnumInch, MeasurementUnitsEnumSqCm, MeasurementUnitsEnumSqM, MeasurementUnitsEnumSqKm, MeasurementUnitsEnumSqFt, MeasurementUnitsEnumSqYd, MeasurementUnitsEnumSqInch, MeasurementUnitsEnumCubicMillimeter, MeasurementUnitsEnumCubicCentimeter, MeasurementUnitsEnumCubicDecimeter, MeasurementUnitsEnumCubicMeter, MeasurementUnitsEnumLiter, MeasurementUnitsEnumCubicFoot, MeasurementUnitsEnumCubicInch, MeasurementUnitsEnumCubicYard, MeasurementUnitsEnumQt, MeasurementUnitsEnumPint, MeasurementUnitsEnumFlOz, MeasurementUnitsEnumAcreIn, MeasurementUnitsEnumAcreFt, MeasurementUnitsEnumG, MeasurementUnitsEnumLb, MeasurementUnitsEnumOz, MeasurementUnitsEnumKg, MeasurementUnitsEnumTonne:
		return true
	}
	return false
}

type MenuErrorCode string

const (
	MenuErrorCodeCannotAssignNode   MenuErrorCode = "CANNOT_ASSIGN_NODE"
	MenuErrorCodeGraphqlError       MenuErrorCode = "GRAPHQL_ERROR"
	MenuErrorCodeInvalid            MenuErrorCode = "INVALID"
	MenuErrorCodeInvalidMenuItem    MenuErrorCode = "INVALID_MENU_ITEM"
	MenuErrorCodeNoMenuItemProvided MenuErrorCode = "NO_MENU_ITEM_PROVIDED"
	MenuErrorCodeNotFound           MenuErrorCode = "NOT_FOUND"
	MenuErrorCodeRequired           MenuErrorCode = "REQUIRED"
	MenuErrorCodeTooManyMenuItems   MenuErrorCode = "TOO_MANY_MENU_ITEMS"
	MenuErrorCodeUnique             MenuErrorCode = "UNIQUE"
)

func (e MenuErrorCode) IsValid() bool {
	switch e {
	case MenuErrorCodeCannotAssignNode, MenuErrorCodeGraphqlError, MenuErrorCodeInvalid, MenuErrorCodeInvalidMenuItem, MenuErrorCodeNoMenuItemProvided, MenuErrorCodeNotFound, MenuErrorCodeRequired, MenuErrorCodeTooManyMenuItems, MenuErrorCodeUnique:
		return true
	}
	return false
}

type MenuItemsSortField string

const (
	MenuItemsSortFieldName MenuItemsSortField = "NAME"
)

func (e MenuItemsSortField) IsValid() bool {
	switch e {
	case MenuItemsSortFieldName:
		return true
	}
	return false
}

type MenuSortField string

const (
	MenuSortFieldName       MenuSortField = "NAME"
	MenuSortFieldItemsCount MenuSortField = "ITEMS_COUNT"
)

func (e MenuSortField) IsValid() bool {
	switch e {
	case MenuSortFieldName, MenuSortFieldItemsCount:
		return true
	}
	return false
}

type TransactionKind = model.TransactionKind

type MetadataErrorCode string

const (
	MetadataErrorCodeGraphqlError MetadataErrorCode = "GRAPHQL_ERROR"
	MetadataErrorCodeInvalid      MetadataErrorCode = "INVALID"
	MetadataErrorCodeNotFound     MetadataErrorCode = "NOT_FOUND"
	MetadataErrorCodeRequired     MetadataErrorCode = "REQUIRED"
)

func (e MetadataErrorCode) IsValid() bool {
	switch e {
	case MetadataErrorCodeGraphqlError, MetadataErrorCodeInvalid, MetadataErrorCodeNotFound, MetadataErrorCodeRequired:
		return true
	}
	return false
}

type NavigationType string

const (
	NavigationTypeMain      NavigationType = "MAIN"
	NavigationTypeSecondary NavigationType = "SECONDARY"
)

func (e NavigationType) IsValid() bool {
	switch e {
	case NavigationTypeMain, NavigationTypeSecondary:
		return true
	}
	return false
}

type OrderAction string

const (
	OrderActionCapture    OrderAction = "CAPTURE"
	OrderActionMarkAsPaid OrderAction = "MARK_AS_PAID"
	OrderActionRefund     OrderAction = "REFUND"
	OrderActionVoid       OrderAction = "VOID"
)

func (e OrderAction) IsValid() bool {
	switch e {
	case OrderActionCapture, OrderActionMarkAsPaid, OrderActionRefund, OrderActionVoid:
		return true
	}
	return false
}

func (o OrderAction) Description() string {
	switch o {
	case OrderActionCapture:
		return "Represents the capture action."
	case OrderActionMarkAsPaid:
		return "Represents a mark-as-paid action."
	case OrderActionRefund:
		return "Represents a refund action."
	case OrderActionVoid:
		return "Represents a void action."
	default:
		return "Unsupported enum value: " + string(o)
	}
}

type OrderDirection = model.OrderDirection

const (
	OrderDirectionAsc  OrderDirection = model.ASC
	OrderDirectionDesc OrderDirection = model.DESC
)

type OrderDiscountType = model.OrderDiscountType

type OrderErrorCode string

const (
	OrderErrorCodeBillingAddressNotSet          OrderErrorCode = "BILLING_ADDRESS_NOT_SET"
	OrderErrorCodeCannotCancelFulfillment       OrderErrorCode = "CANNOT_CANCEL_FULFILLMENT"
	OrderErrorCodeCannotCancelOrder             OrderErrorCode = "CANNOT_CANCEL_ORDER"
	OrderErrorCodeCannotDelete                  OrderErrorCode = "CANNOT_DELETE"
	OrderErrorCodeCannotDiscount                OrderErrorCode = "CANNOT_DISCOUNT"
	OrderErrorCodeCannotRefund                  OrderErrorCode = "CANNOT_REFUND"
	OrderErrorCodeCannotFulfillUnpaidOrder      OrderErrorCode = "CANNOT_FULFILL_UNPAID_ORDER"
	OrderErrorCodeCaptureInactivePayment        OrderErrorCode = "CAPTURE_INACTIVE_PAYMENT"
	OrderErrorCodeGiftCardLine                  OrderErrorCode = "GIFT_CARD_LINE"
	OrderErrorCodeNotEditable                   OrderErrorCode = "NOT_EDITABLE"
	OrderErrorCodeFulfillOrderLine              OrderErrorCode = "FULFILL_ORDER_LINE"
	OrderErrorCodeGraphqlError                  OrderErrorCode = "GRAPHQL_ERROR"
	OrderErrorCodeInvalid                       OrderErrorCode = "INVALID"
	OrderErrorCodeProductNotPublished           OrderErrorCode = "PRODUCT_NOT_PUBLISHED"
	OrderErrorCodeProductUnavailableForPurchase OrderErrorCode = "PRODUCT_UNAVAILABLE_FOR_PURCHASE"
	OrderErrorCodeNotFound                      OrderErrorCode = "NOT_FOUND"
	OrderErrorCodeOrderNoShippingAddress        OrderErrorCode = "ORDER_NO_SHIPPING_ADDRESS"
	OrderErrorCodePaymentError                  OrderErrorCode = "PAYMENT_ERROR"
	OrderErrorCodePaymentMissing                OrderErrorCode = "PAYMENT_MISSING"
	OrderErrorCodeRequired                      OrderErrorCode = "REQUIRED"
	OrderErrorCodeShippingMethodNotApplicable   OrderErrorCode = "SHIPPING_METHOD_NOT_APPLICABLE"
	OrderErrorCodeShippingMethodRequired        OrderErrorCode = "SHIPPING_METHOD_REQUIRED"
	OrderErrorCodeTaxError                      OrderErrorCode = "TAX_ERROR"
	OrderErrorCodeUnique                        OrderErrorCode = "UNIQUE"
	OrderErrorCodeVoidInactivePayment           OrderErrorCode = "VOID_INACTIVE_PAYMENT"
	OrderErrorCodeZeroQuantity                  OrderErrorCode = "ZERO_QUANTITY"
	OrderErrorCodeInvalidQuantity               OrderErrorCode = "INVALID_QUANTITY"
	OrderErrorCodeInsufficientStock             OrderErrorCode = "INSUFFICIENT_STOCK"
	OrderErrorCodeDuplicatedInputItem           OrderErrorCode = "DUPLICATED_INPUT_ITEM"
	OrderErrorCodeNotAvailableInChannel         OrderErrorCode = "NOT_AVAILABLE_IN_CHANNEL"
	OrderErrorCodeChannelInactive               OrderErrorCode = "CHANNEL_INACTIVE"
)

func (e OrderErrorCode) IsValid() bool {
	switch e {
	case OrderErrorCodeBillingAddressNotSet, OrderErrorCodeCannotCancelFulfillment, OrderErrorCodeCannotCancelOrder, OrderErrorCodeCannotDelete, OrderErrorCodeCannotDiscount, OrderErrorCodeCannotRefund, OrderErrorCodeCannotFulfillUnpaidOrder, OrderErrorCodeCaptureInactivePayment, OrderErrorCodeGiftCardLine, OrderErrorCodeNotEditable, OrderErrorCodeFulfillOrderLine, OrderErrorCodeGraphqlError, OrderErrorCodeInvalid, OrderErrorCodeProductNotPublished, OrderErrorCodeProductUnavailableForPurchase, OrderErrorCodeNotFound, OrderErrorCodeOrderNoShippingAddress, OrderErrorCodePaymentError, OrderErrorCodePaymentMissing, OrderErrorCodeRequired, OrderErrorCodeShippingMethodNotApplicable, OrderErrorCodeShippingMethodRequired, OrderErrorCodeTaxError, OrderErrorCodeUnique, OrderErrorCodeVoidInactivePayment, OrderErrorCodeZeroQuantity, OrderErrorCodeInvalidQuantity, OrderErrorCodeInsufficientStock, OrderErrorCodeDuplicatedInputItem, OrderErrorCodeNotAvailableInChannel, OrderErrorCodeChannelInactive:
		return true
	}
	return false
}

type OrderEventsEmailsEnum string

const (
	OrderEventsEmailsEnumPaymentConfirmation     OrderEventsEmailsEnum = "PAYMENT_CONFIRMATION"
	OrderEventsEmailsEnumConfirmed               OrderEventsEmailsEnum = "CONFIRMED"
	OrderEventsEmailsEnumShippingConfirmation    OrderEventsEmailsEnum = "SHIPPING_CONFIRMATION"
	OrderEventsEmailsEnumTrackingUpdated         OrderEventsEmailsEnum = "TRACKING_UPDATED"
	OrderEventsEmailsEnumOrderConfirmation       OrderEventsEmailsEnum = "ORDER_CONFIRMATION"
	OrderEventsEmailsEnumOrderCancel             OrderEventsEmailsEnum = "ORDER_CANCEL"
	OrderEventsEmailsEnumOrderRefund             OrderEventsEmailsEnum = "ORDER_REFUND"
	OrderEventsEmailsEnumFulfillmentConfirmation OrderEventsEmailsEnum = "FULFILLMENT_CONFIRMATION"
	OrderEventsEmailsEnumDigitalLinks            OrderEventsEmailsEnum = "DIGITAL_LINKS"
)

func (e OrderEventsEmailsEnum) IsValid() bool {
	switch e {
	case OrderEventsEmailsEnumPaymentConfirmation, OrderEventsEmailsEnumConfirmed, OrderEventsEmailsEnumShippingConfirmation, OrderEventsEmailsEnumTrackingUpdated, OrderEventsEmailsEnumOrderConfirmation, OrderEventsEmailsEnumOrderCancel, OrderEventsEmailsEnumOrderRefund, OrderEventsEmailsEnumFulfillmentConfirmation, OrderEventsEmailsEnumDigitalLinks:
		return true
	}
	return false
}

type OrderEventsEnum = model.OrderEventType

type OrderSettingsErrorCode string

const (
	OrderSettingsErrorCodeInvalid OrderSettingsErrorCode = "INVALID"
)

func (e OrderSettingsErrorCode) IsValid() bool {
	switch e {
	case OrderSettingsErrorCodeInvalid:
		return true
	}
	return false
}

type OrderStatus = model.OrderStatus

type OrderSortField string

const (
	OrderSortFieldNumber            OrderSortField = "NUMBER"
	OrderSortFieldCreationDate      OrderSortField = "CREATION_DATE"
	OrderSortFieldCustomer          OrderSortField = "CUSTOMER"
	OrderSortFieldPayment           OrderSortField = "PAYMENT"
	OrderSortFieldFulfillmentStatus OrderSortField = "FULFILLMENT_STATUS"
)

func (e OrderSortField) IsValid() bool {
	switch e {
	case OrderSortFieldNumber, OrderSortFieldCreationDate, OrderSortFieldCustomer, OrderSortFieldPayment, OrderSortFieldFulfillmentStatus:
		return true
	}
	return false
}

type OrderStatusFilter string

const (
	OrderStatusFilterReadyToFulfill     OrderStatusFilter = "READY_TO_FULFILL"
	OrderStatusFilterReadyToCapture     OrderStatusFilter = "READY_TO_CAPTURE"
	OrderStatusFilterUnfulfilled        OrderStatusFilter = "UNFULFILLED"
	OrderStatusFilterUnconfirmed        OrderStatusFilter = "UNCONFIRMED"
	OrderStatusFilterPartiallyFulfilled OrderStatusFilter = "PARTIALLY_FULFILLED"
	OrderStatusFilterFulfilled          OrderStatusFilter = "FULFILLED"
	OrderStatusFilterCanceled           OrderStatusFilter = "CANCELED"
)

func (e OrderStatusFilter) IsValid() bool {
	switch e {
	case OrderStatusFilterReadyToFulfill, OrderStatusFilterReadyToCapture, OrderStatusFilterUnfulfilled, OrderStatusFilterUnconfirmed, OrderStatusFilterPartiallyFulfilled, OrderStatusFilterFulfilled, OrderStatusFilterCanceled:
		return true
	}
	return false
}

type PageErrorCode string

const (
	PageErrorCodeGraphqlError             PageErrorCode = "GRAPHQL_ERROR"
	PageErrorCodeInvalid                  PageErrorCode = "INVALID"
	PageErrorCodeNotFound                 PageErrorCode = "NOT_FOUND"
	PageErrorCodeRequired                 PageErrorCode = "REQUIRED"
	PageErrorCodeUnique                   PageErrorCode = "UNIQUE"
	PageErrorCodeDuplicatedInputItem      PageErrorCode = "DUPLICATED_INPUT_ITEM"
	PageErrorCodeAttributeAlreadyAssigned PageErrorCode = "ATTRIBUTE_ALREADY_ASSIGNED"
)

func (e PageErrorCode) IsValid() bool {
	switch e {
	case PageErrorCodeGraphqlError, PageErrorCodeInvalid, PageErrorCodeNotFound, PageErrorCodeRequired, PageErrorCodeUnique, PageErrorCodeDuplicatedInputItem, PageErrorCodeAttributeAlreadyAssigned:
		return true
	}
	return false
}

type PageSortField string

const (
	PageSortFieldTitle           PageSortField = "TITLE"
	PageSortFieldSlug            PageSortField = "SLUG"
	PageSortFieldVisibility      PageSortField = "VISIBILITY"
	PageSortFieldCreationDate    PageSortField = "CREATION_DATE"
	PageSortFieldPublicationDate PageSortField = "PUBLICATION_DATE"
)

func (e PageSortField) IsValid() bool {
	switch e {
	case PageSortFieldTitle, PageSortFieldSlug, PageSortFieldVisibility, PageSortFieldCreationDate, PageSortFieldPublicationDate:
		return true
	}
	return false
}

type PageTypeSortField string

const (
	PageTypeSortFieldName PageTypeSortField = "NAME"
	PageTypeSortFieldSlug PageTypeSortField = "SLUG"
)

func (e PageTypeSortField) IsValid() bool {
	switch e {
	case PageTypeSortFieldName, PageTypeSortFieldSlug:
		return true
	}
	return false
}

type PaymentChargeStatusEnum = model.PaymentChargeStatus

type PaymentErrorCode string

const (
	PaymentErrorCodeBillingAddressNotSet     PaymentErrorCode = "BILLING_ADDRESS_NOT_SET"
	PaymentErrorCodeGraphqlError             PaymentErrorCode = "GRAPHQL_ERROR"
	PaymentErrorCodeInvalid                  PaymentErrorCode = "INVALID"
	PaymentErrorCodeNotFound                 PaymentErrorCode = "NOT_FOUND"
	PaymentErrorCodeRequired                 PaymentErrorCode = "REQUIRED"
	PaymentErrorCodeUnique                   PaymentErrorCode = "UNIQUE"
	PaymentErrorCodePartialPaymentNotAllowed PaymentErrorCode = "PARTIAL_PAYMENT_NOT_ALLOWED"
	PaymentErrorCodeShippingAddressNotSet    PaymentErrorCode = "SHIPPING_ADDRESS_NOT_SET"
	PaymentErrorCodeInvalidShippingMethod    PaymentErrorCode = "INVALID_SHIPPING_METHOD"
	PaymentErrorCodeShippingMethodNotSet     PaymentErrorCode = "SHIPPING_METHOD_NOT_SET"
	PaymentErrorCodePaymentError             PaymentErrorCode = "PAYMENT_ERROR"
	PaymentErrorCodeNotSupportedGateway      PaymentErrorCode = "NOT_SUPPORTED_GATEWAY"
	PaymentErrorCodeChannelInactive          PaymentErrorCode = "CHANNEL_INACTIVE"
)

func (e PaymentErrorCode) IsValid() bool {
	switch e {
	case PaymentErrorCodeBillingAddressNotSet, PaymentErrorCodeGraphqlError, PaymentErrorCodeInvalid, PaymentErrorCodeNotFound, PaymentErrorCodeRequired, PaymentErrorCodeUnique, PaymentErrorCodePartialPaymentNotAllowed, PaymentErrorCodeShippingAddressNotSet, PaymentErrorCodeInvalidShippingMethod, PaymentErrorCodeShippingMethodNotSet, PaymentErrorCodePaymentError, PaymentErrorCodeNotSupportedGateway, PaymentErrorCodeChannelInactive:
		return true
	}
	return false
}

type PermissionGroupErrorCode string

const (
	PermissionGroupErrorCodeAssignNonStaffMember        PermissionGroupErrorCode = "ASSIGN_NON_STAFF_MEMBER"
	PermissionGroupErrorCodeDuplicatedInputItem         PermissionGroupErrorCode = "DUPLICATED_INPUT_ITEM"
	PermissionGroupErrorCodeCannotRemoveFromLastGroup   PermissionGroupErrorCode = "CANNOT_REMOVE_FROM_LAST_GROUP"
	PermissionGroupErrorCodeLeftNotManageablePermission PermissionGroupErrorCode = "LEFT_NOT_MANAGEABLE_PERMISSION"
	PermissionGroupErrorCodeOutOfScopePermission        PermissionGroupErrorCode = "OUT_OF_SCOPE_PERMISSION"
	PermissionGroupErrorCodeOutOfScopeUser              PermissionGroupErrorCode = "OUT_OF_SCOPE_USER"
	PermissionGroupErrorCodeRequired                    PermissionGroupErrorCode = "REQUIRED"
	PermissionGroupErrorCodeUnique                      PermissionGroupErrorCode = "UNIQUE"
)

type PermissionGroupSortField string

const (
	PermissionGroupSortFieldName PermissionGroupSortField = "NAME"
)

type PluginConfigurationType string

const (
	PluginConfigurationTypePerChannel PluginConfigurationType = "PER_CHANNEL"
	PluginConfigurationTypeGlobal     PluginConfigurationType = "GLOBAL"
)

func (e PluginConfigurationType) IsValid() bool {
	switch e {
	case PluginConfigurationTypePerChannel, PluginConfigurationTypeGlobal:
		return true
	}
	return false
}

type PluginErrorCode string

const (
	PluginErrorCodeGraphqlError        PluginErrorCode = "GRAPHQL_ERROR"
	PluginErrorCodeInvalid             PluginErrorCode = "INVALID"
	PluginErrorCodePluginMisconfigured PluginErrorCode = "PLUGIN_MISCONFIGURED"
	PluginErrorCodeNotFound            PluginErrorCode = "NOT_FOUND"
	PluginErrorCodeRequired            PluginErrorCode = "REQUIRED"
	PluginErrorCodeUnique              PluginErrorCode = "UNIQUE"
)

func (e PluginErrorCode) IsValid() bool {
	switch e {
	case PluginErrorCodeGraphqlError, PluginErrorCodeInvalid, PluginErrorCodePluginMisconfigured, PluginErrorCodeNotFound, PluginErrorCodeRequired, PluginErrorCodeUnique:
		return true
	}
	return false
}

type PluginSortField string

const (
	PluginSortFieldName     PluginSortField = "NAME"
	PluginSortFieldIsActive PluginSortField = "IS_ACTIVE"
)

func (e PluginSortField) IsValid() bool {
	switch e {
	case PluginSortFieldName, PluginSortFieldIsActive:
		return true
	}
	return false
}

type PostalCodeRuleInclusionTypeEnum = model.InclusionType

type ProductAttributeType string

const (
	ProductAttributeTypeProduct ProductAttributeType = "PRODUCT"
	ProductAttributeTypeVariant ProductAttributeType = "VARIANT"
)

func (e ProductAttributeType) IsValid() bool {
	switch e {
	case ProductAttributeTypeProduct, ProductAttributeTypeVariant:
		return true
	}
	return false
}

type ProductErrorCode string

const (
	ProductErrorCodeAlreadyExists                     ProductErrorCode = "ALREADY_EXISTS"
	ProductErrorCodeAttributeAlreadyAssigned          ProductErrorCode = "ATTRIBUTE_ALREADY_ASSIGNED"
	ProductErrorCodeAttributeCannotBeAssigned         ProductErrorCode = "ATTRIBUTE_CANNOT_BE_ASSIGNED"
	ProductErrorCodeAttributeVariantsDisabled         ProductErrorCode = "ATTRIBUTE_VARIANTS_DISABLED"
	ProductErrorCodeDuplicatedInputItem               ProductErrorCode = "DUPLICATED_INPUT_ITEM"
	ProductErrorCodeGraphqlError                      ProductErrorCode = "GRAPHQL_ERROR"
	ProductErrorCodeInvalid                           ProductErrorCode = "INVALID"
	ProductErrorCodeProductWithoutCategory            ProductErrorCode = "PRODUCT_WITHOUT_CATEGORY"
	ProductErrorCodeNotProductsImage                  ProductErrorCode = "NOT_PRODUCTS_IMAGE"
	ProductErrorCodeNotProductsVariant                ProductErrorCode = "NOT_PRODUCTS_VARIANT"
	ProductErrorCodeNotFound                          ProductErrorCode = "NOT_FOUND"
	ProductErrorCodeRequired                          ProductErrorCode = "REQUIRED"
	ProductErrorCodeUnique                            ProductErrorCode = "UNIQUE"
	ProductErrorCodeVariantNoDigitalContent           ProductErrorCode = "VARIANT_NO_DIGITAL_CONTENT"
	ProductErrorCodeCannotManageProductWithoutVariant ProductErrorCode = "CANNOT_MANAGE_PRODUCT_WITHOUT_VARIANT"
	ProductErrorCodeProductNotAssignedToChannel       ProductErrorCode = "PRODUCT_NOT_ASSIGNED_TO_CHANNEL"
	ProductErrorCodeUnsupportedMediaProvider          ProductErrorCode = "UNSUPPORTED_MEDIA_PROVIDER"
)

func (e ProductErrorCode) IsValid() bool {
	switch e {
	case ProductErrorCodeAlreadyExists, ProductErrorCodeAttributeAlreadyAssigned, ProductErrorCodeAttributeCannotBeAssigned, ProductErrorCodeAttributeVariantsDisabled, ProductErrorCodeDuplicatedInputItem, ProductErrorCodeGraphqlError, ProductErrorCodeInvalid, ProductErrorCodeProductWithoutCategory, ProductErrorCodeNotProductsImage, ProductErrorCodeNotProductsVariant, ProductErrorCodeNotFound, ProductErrorCodeRequired, ProductErrorCodeUnique, ProductErrorCodeVariantNoDigitalContent, ProductErrorCodeCannotManageProductWithoutVariant, ProductErrorCodeProductNotAssignedToChannel, ProductErrorCodeUnsupportedMediaProvider:
		return true
	}
	return false
}

type ProductFieldEnum string

const (
	ProductFieldEnumName          ProductFieldEnum = "NAME"
	ProductFieldEnumDescription   ProductFieldEnum = "DESCRIPTION"
	ProductFieldEnumProductType   ProductFieldEnum = "PRODUCT_TYPE"
	ProductFieldEnumCategory      ProductFieldEnum = "CATEGORY"
	ProductFieldEnumProductWeight ProductFieldEnum = "PRODUCT_WEIGHT"
	ProductFieldEnumCollections   ProductFieldEnum = "COLLECTIONS"
	ProductFieldEnumChargeTaxes   ProductFieldEnum = "CHARGE_TAXES"
	ProductFieldEnumProductMedia  ProductFieldEnum = "PRODUCT_MEDIA"
	ProductFieldEnumVariantID     ProductFieldEnum = "VARIANT_ID"
	ProductFieldEnumVariantSku    ProductFieldEnum = "VARIANT_SKU"
	ProductFieldEnumVariantWeight ProductFieldEnum = "VARIANT_WEIGHT"
	ProductFieldEnumVariantMedia  ProductFieldEnum = "VARIANT_MEDIA"
)

func (e ProductFieldEnum) IsValid() bool {
	switch e {
	case ProductFieldEnumName, ProductFieldEnumDescription, ProductFieldEnumProductType, ProductFieldEnumCategory, ProductFieldEnumProductWeight, ProductFieldEnumCollections, ProductFieldEnumChargeTaxes, ProductFieldEnumProductMedia, ProductFieldEnumVariantID, ProductFieldEnumVariantSku, ProductFieldEnumVariantWeight, ProductFieldEnumVariantMedia:
		return true
	}
	return false
}

type ProductMediaType string

const (
	ProductMediaTypeImage ProductMediaType = model.IMAGE
	ProductMediaTypeVideo ProductMediaType = model.VIDEO
)

func (e ProductMediaType) IsValid() bool {
	switch e {
	case ProductMediaTypeImage, ProductMediaTypeVideo:
		return true
	}
	return false
}

type ProductOrderField = model.ProductOrderField

type ProductTypeConfigurable string

const (
	ProductTypeConfigurableConfigurable ProductTypeConfigurable = "CONFIGURABLE"
	ProductTypeConfigurableSimple       ProductTypeConfigurable = "SIMPLE"
)

type ProductTypeEnum string

const (
	ProductTypeEnumDigital   ProductTypeEnum = "DIGITAL"
	ProductTypeEnumShippable ProductTypeEnum = "SHIPPABLE"
)

type ProductTypeKindEnum = model.ProductTypeKind

type ProductTypeSortField string

const (
	ProductTypeSortFieldName             ProductTypeSortField = "NAME"
	ProductTypeSortFieldDigital          ProductTypeSortField = "DIGITAL"
	ProductTypeSortFieldShippingRequired ProductTypeSortField = "SHIPPING_REQUIRED"
)

type ReportingPeriod string

const (
	ReportingPeriodToday     ReportingPeriod = "TODAY"
	ReportingPeriodThisMonth ReportingPeriod = "THIS_MONTH"
)

type SaleSortField string

const (
	SaleSortFieldName      SaleSortField = "NAME"
	SaleSortFieldStartDate SaleSortField = "START_DATE"
	SaleSortFieldEndDate   SaleSortField = "END_DATE"
	SaleSortFieldValue     SaleSortField = "VALUE"
	SaleSortFieldType      SaleSortField = "TYPE"
)

type SaleType = model.DiscountType

type ShippingErrorCode string

const (
	ShippingErrorCodeAlreadyExists       ShippingErrorCode = "ALREADY_EXISTS"
	ShippingErrorCodeGraphqlError        ShippingErrorCode = "GRAPHQL_ERROR"
	ShippingErrorCodeInvalid             ShippingErrorCode = "INVALID"
	ShippingErrorCodeMaxLessThanMin      ShippingErrorCode = "MAX_LESS_THAN_MIN"
	ShippingErrorCodeNotFound            ShippingErrorCode = "NOT_FOUND"
	ShippingErrorCodeRequired            ShippingErrorCode = "REQUIRED"
	ShippingErrorCodeUnique              ShippingErrorCode = "UNIQUE"
	ShippingErrorCodeDuplicatedInputItem ShippingErrorCode = "DUPLICATED_INPUT_ITEM"
)

type ShippingMethodTypeEnum = model.ShippingMethodType

type ShopErrorCode string

const (
	ShopErrorCodeAlreadyExists       ShopErrorCode = "ALREADY_EXISTS"
	ShopErrorCodeCannotFetchTaxRates ShopErrorCode = "CANNOT_FETCH_TAX_RATES"
	ShopErrorCodeGraphqlError        ShopErrorCode = "GRAPHQL_ERROR"
	ShopErrorCodeInvalid             ShopErrorCode = "INVALID"
	ShopErrorCodeNotFound            ShopErrorCode = "NOT_FOUND"
	ShopErrorCodeRequired            ShopErrorCode = "REQUIRED"
	ShopErrorCodeUnique              ShopErrorCode = "UNIQUE"
)

type StaffMemberStatus string

const (
	StaffMemberStatusActive      StaffMemberStatus = "active"
	StaffMemberStatusDeactivated StaffMemberStatus = "deactivated"
)

func (s StaffMemberStatus) IsValid() bool {
	return s == StaffMemberStatusActive || s == StaffMemberStatusDeactivated
}

type StockAvailability = model.StockAvailability

type StockErrorCode string

const (
	StockErrorCodeAlreadyExists StockErrorCode = "ALREADY_EXISTS"
	StockErrorCodeGraphqlError  StockErrorCode = "GRAPHQL_ERROR"
	StockErrorCodeInvalid       StockErrorCode = "INVALID"
	StockErrorCodeNotFound      StockErrorCode = "NOT_FOUND"
	StockErrorCodeRequired      StockErrorCode = "REQUIRED"
	StockErrorCodeUnique        StockErrorCode = "UNIQUE"
)

type StorePaymentMethodEnum string

const (
	StorePaymentMethodEnumOnSession  StorePaymentMethodEnum = "ON_SESSION"
	StorePaymentMethodEnumOffSession StorePaymentMethodEnum = "OFF_SESSION"
	StorePaymentMethodEnumNone       StorePaymentMethodEnum = "NONE"
)

type TranslatableKinds string

const (
	TranslatableKindsAttribute      TranslatableKinds = "ATTRIBUTE"
	TranslatableKindsAttributeValue TranslatableKinds = "ATTRIBUTE_VALUE"
	TranslatableKindsCategory       TranslatableKinds = "CATEGORY"
	TranslatableKindsCollection     TranslatableKinds = "COLLECTION"
	TranslatableKindsMenuItem       TranslatableKinds = "MENU_ITEM"
	TranslatableKindsPage           TranslatableKinds = "PAGE"
	TranslatableKindsProduct        TranslatableKinds = "PRODUCT"
	TranslatableKindsSale           TranslatableKinds = "SALE"
	TranslatableKindsShippingMethod TranslatableKinds = "SHIPPING_METHOD"
	TranslatableKindsVariant        TranslatableKinds = "VARIANT"
	TranslatableKindsVoucher        TranslatableKinds = "VOUCHER"
)

type TranslationErrorCode string

const (
	TranslationErrorCodeGraphqlError TranslationErrorCode = "GRAPHQL_ERROR"
	TranslationErrorCodeNotFound     TranslationErrorCode = "NOT_FOUND"
	TranslationErrorCodeRequired     TranslationErrorCode = "REQUIRED"
)

type UploadErrorCode string

const (
	UploadErrorCodeGraphqlError UploadErrorCode = "GRAPHQL_ERROR"
)

type UserSortField string

const (
	UserSortFieldFirstName  UserSortField = "FIRST_NAME"
	UserSortFieldLastName   UserSortField = "LAST_NAME"
	UserSortFieldEmail      UserSortField = "EMAIL"
	UserSortFieldOrderCount UserSortField = "ORDER_COUNT"
)

type VariantAttributeScope string

const (
	VariantAttributeScopeAll                 VariantAttributeScope = "ALL"
	VariantAttributeScopeVariantSelection    VariantAttributeScope = "VARIANT_SELECTION"
	VariantAttributeScopeNotVariantSelection VariantAttributeScope = "NOT_VARIANT_SELECTION"
)

type VolumeUnitsEnum string

const (
	VolumeUnitsEnumCubicMillimeter VolumeUnitsEnum = measurement.CUBIC_MILLIMETER
	VolumeUnitsEnumCubicCentimeter VolumeUnitsEnum = measurement.CUBIC_CENTIMETER
	VolumeUnitsEnumCubicDecimeter  VolumeUnitsEnum = measurement.CUBIC_DECIMETER
	VolumeUnitsEnumCubicMeter      VolumeUnitsEnum = measurement.CUBIC_METER
	VolumeUnitsEnumLiter           VolumeUnitsEnum = measurement.LITER
	VolumeUnitsEnumCubicFoot       VolumeUnitsEnum = measurement.CUBIC_FOOT
	VolumeUnitsEnumCubicInch       VolumeUnitsEnum = measurement.CUBIC_INCH
	VolumeUnitsEnumCubicYard       VolumeUnitsEnum = measurement.CUBIC_YARD
	VolumeUnitsEnumQt              VolumeUnitsEnum = measurement.QT
	VolumeUnitsEnumPint            VolumeUnitsEnum = measurement.PINT
	VolumeUnitsEnumFlOz            VolumeUnitsEnum = measurement.FL_OZ
	VolumeUnitsEnumAcreIn          VolumeUnitsEnum = measurement.ACRE_IN
	VolumeUnitsEnumAcreFt          VolumeUnitsEnum = measurement.ACRE_FT
)

type VoucherDiscountType string

const (
	VoucherDiscountTypeFixed      VoucherDiscountType = "FIXED"
	VoucherDiscountTypePercentage VoucherDiscountType = "PERCENTAGE"
	VoucherDiscountTypeShipping   VoucherDiscountType = "SHIPPING"
)

type VoucherSortField string

const (
	VoucherSortFieldCode               VoucherSortField = "CODE"
	VoucherSortFieldStartDate          VoucherSortField = "START_DATE"
	VoucherSortFieldEndDate            VoucherSortField = "END_DATE"
	VoucherSortFieldValue              VoucherSortField = "VALUE"
	VoucherSortFieldType               VoucherSortField = "TYPE"
	VoucherSortFieldUsageLimit         VoucherSortField = "USAGE_LIMIT"
	VoucherSortFieldMinimumSpentAmount VoucherSortField = "MINIMUM_SPENT_AMOUNT"
)

type VoucherTypeEnum string

const (
	VoucherTypeEnumShipping        VoucherTypeEnum = model.SHIPPING
	VoucherTypeEnumEntireOrder     VoucherTypeEnum = model.ENTIRE_ORDER
	VoucherTypeEnumSpecificProduct VoucherTypeEnum = model.SPECIFIC_PRODUCT
)

type WarehouseErrorCode string

const (
	WarehouseErrorCodeAlreadyExists WarehouseErrorCode = "ALREADY_EXISTS"
	WarehouseErrorCodeGraphqlError  WarehouseErrorCode = "GRAPHQL_ERROR"
	WarehouseErrorCodeInvalid       WarehouseErrorCode = "INVALID"
	WarehouseErrorCodeNotFound      WarehouseErrorCode = "NOT_FOUND"
	WarehouseErrorCodeRequired      WarehouseErrorCode = "REQUIRED"
	WarehouseErrorCodeUnique        WarehouseErrorCode = "UNIQUE"
)

type WarehouseSortField string

const (
	WarehouseSortFieldName WarehouseSortField = "NAME"
)

type WebhookErrorCode string

const (
	WebhookErrorCodeGraphqlError WebhookErrorCode = "GRAPHQL_ERROR"
	WebhookErrorCodeInvalid      WebhookErrorCode = "INVALID"
	WebhookErrorCodeNotFound     WebhookErrorCode = "NOT_FOUND"
	WebhookErrorCodeRequired     WebhookErrorCode = "REQUIRED"
	WebhookErrorCodeUnique       WebhookErrorCode = "UNIQUE"
)

type WebhookEventTypeEnum string

const (
	WebhookEventTypeEnumAnyEvents                 WebhookEventTypeEnum = "ANY_EVENTS"
	WebhookEventTypeEnumOrderCreated              WebhookEventTypeEnum = "ORDER_CREATED"
	WebhookEventTypeEnumOrderConfirmed            WebhookEventTypeEnum = "ORDER_CONFIRMED"
	WebhookEventTypeEnumOrderFullyPaid            WebhookEventTypeEnum = "ORDER_FULLY_PAID"
	WebhookEventTypeEnumOrderUpdated              WebhookEventTypeEnum = "ORDER_UPDATED"
	WebhookEventTypeEnumOrderCancelled            WebhookEventTypeEnum = "ORDER_CANCELLED"
	WebhookEventTypeEnumOrderFulfilled            WebhookEventTypeEnum = "ORDER_FULFILLED"
	WebhookEventTypeEnumDraftOrderCreated         WebhookEventTypeEnum = "DRAFT_ORDER_CREATED"
	WebhookEventTypeEnumDraftOrderUpdated         WebhookEventTypeEnum = "DRAFT_ORDER_UPDATED"
	WebhookEventTypeEnumDraftOrderDeleted         WebhookEventTypeEnum = "DRAFT_ORDER_DELETED"
	WebhookEventTypeEnumSaleCreated               WebhookEventTypeEnum = "SALE_CREATED"
	WebhookEventTypeEnumSaleUpdated               WebhookEventTypeEnum = "SALE_UPDATED"
	WebhookEventTypeEnumSaleDeleted               WebhookEventTypeEnum = "SALE_DELETED"
	WebhookEventTypeEnumInvoiceRequested          WebhookEventTypeEnum = "INVOICE_REQUESTED"
	WebhookEventTypeEnumInvoiceDeleted            WebhookEventTypeEnum = "INVOICE_DELETED"
	WebhookEventTypeEnumInvoiceSent               WebhookEventTypeEnum = "INVOICE_SENT"
	WebhookEventTypeEnumCustomerCreated           WebhookEventTypeEnum = "CUSTOMER_CREATED"
	WebhookEventTypeEnumCustomerUpdated           WebhookEventTypeEnum = "CUSTOMER_UPDATED"
	WebhookEventTypeEnumProductCreated            WebhookEventTypeEnum = "PRODUCT_CREATED"
	WebhookEventTypeEnumProductUpdated            WebhookEventTypeEnum = "PRODUCT_UPDATED"
	WebhookEventTypeEnumProductDeleted            WebhookEventTypeEnum = "PRODUCT_DELETED"
	WebhookEventTypeEnumProductVariantCreated     WebhookEventTypeEnum = "PRODUCT_VARIANT_CREATED"
	WebhookEventTypeEnumProductVariantUpdated     WebhookEventTypeEnum = "PRODUCT_VARIANT_UPDATED"
	WebhookEventTypeEnumProductVariantDeleted     WebhookEventTypeEnum = "PRODUCT_VARIANT_DELETED"
	WebhookEventTypeEnumProductVariantOutOfStock  WebhookEventTypeEnum = "PRODUCT_VARIANT_OUT_OF_STOCK"
	WebhookEventTypeEnumProductVariantBackInStock WebhookEventTypeEnum = "PRODUCT_VARIANT_BACK_IN_STOCK"
	WebhookEventTypeEnumCheckoutCreated           WebhookEventTypeEnum = "CHECKOUT_CREATED"
	WebhookEventTypeEnumCheckoutUpdated           WebhookEventTypeEnum = "CHECKOUT_UPDATED"
	WebhookEventTypeEnumFulfillmentCreated        WebhookEventTypeEnum = "FULFILLMENT_CREATED"
	WebhookEventTypeEnumFulfillmentCanceled       WebhookEventTypeEnum = "FULFILLMENT_CANCELED"
	WebhookEventTypeEnumNotifyUser                WebhookEventTypeEnum = "NOTIFY_USER"
	WebhookEventTypeEnumPageCreated               WebhookEventTypeEnum = "PAGE_CREATED"
	WebhookEventTypeEnumPageUpdated               WebhookEventTypeEnum = "PAGE_UPDATED"
	WebhookEventTypeEnumPageDeleted               WebhookEventTypeEnum = "PAGE_DELETED"
	WebhookEventTypeEnumPaymentAuthorize          WebhookEventTypeEnum = "PAYMENT_AUTHORIZE"
	WebhookEventTypeEnumPaymentCapture            WebhookEventTypeEnum = "PAYMENT_CAPTURE"
	WebhookEventTypeEnumPaymentConfirm            WebhookEventTypeEnum = "PAYMENT_CONFIRM"
	WebhookEventTypeEnumPaymentListGateways       WebhookEventTypeEnum = "PAYMENT_LIST_GATEWAYS"
	WebhookEventTypeEnumPaymentProcess            WebhookEventTypeEnum = "PAYMENT_PROCESS"
	WebhookEventTypeEnumPaymentRefund             WebhookEventTypeEnum = "PAYMENT_REFUND"
	WebhookEventTypeEnumPaymentVoid               WebhookEventTypeEnum = "PAYMENT_VOID"
	WebhookEventTypeEnumTranslationCreated        WebhookEventTypeEnum = "TRANSLATION_CREATED"
	WebhookEventTypeEnumTranslationUpdated        WebhookEventTypeEnum = "TRANSLATION_UPDATED"
)

type WebhookSampleEventTypeEnum string

const (
	WebhookSampleEventTypeEnumOrderCreated              WebhookSampleEventTypeEnum = "ORDER_CREATED"
	WebhookSampleEventTypeEnumOrderConfirmed            WebhookSampleEventTypeEnum = "ORDER_CONFIRMED"
	WebhookSampleEventTypeEnumOrderFullyPaid            WebhookSampleEventTypeEnum = "ORDER_FULLY_PAID"
	WebhookSampleEventTypeEnumOrderUpdated              WebhookSampleEventTypeEnum = "ORDER_UPDATED"
	WebhookSampleEventTypeEnumOrderCancelled            WebhookSampleEventTypeEnum = "ORDER_CANCELLED"
	WebhookSampleEventTypeEnumOrderFulfilled            WebhookSampleEventTypeEnum = "ORDER_FULFILLED"
	WebhookSampleEventTypeEnumDraftOrderCreated         WebhookSampleEventTypeEnum = "DRAFT_ORDER_CREATED"
	WebhookSampleEventTypeEnumDraftOrderUpdated         WebhookSampleEventTypeEnum = "DRAFT_ORDER_UPDATED"
	WebhookSampleEventTypeEnumDraftOrderDeleted         WebhookSampleEventTypeEnum = "DRAFT_ORDER_DELETED"
	WebhookSampleEventTypeEnumSaleCreated               WebhookSampleEventTypeEnum = "SALE_CREATED"
	WebhookSampleEventTypeEnumSaleUpdated               WebhookSampleEventTypeEnum = "SALE_UPDATED"
	WebhookSampleEventTypeEnumSaleDeleted               WebhookSampleEventTypeEnum = "SALE_DELETED"
	WebhookSampleEventTypeEnumInvoiceRequested          WebhookSampleEventTypeEnum = "INVOICE_REQUESTED"
	WebhookSampleEventTypeEnumInvoiceDeleted            WebhookSampleEventTypeEnum = "INVOICE_DELETED"
	WebhookSampleEventTypeEnumInvoiceSent               WebhookSampleEventTypeEnum = "INVOICE_SENT"
	WebhookSampleEventTypeEnumCustomerCreated           WebhookSampleEventTypeEnum = "CUSTOMER_CREATED"
	WebhookSampleEventTypeEnumCustomerUpdated           WebhookSampleEventTypeEnum = "CUSTOMER_UPDATED"
	WebhookSampleEventTypeEnumProductCreated            WebhookSampleEventTypeEnum = "PRODUCT_CREATED"
	WebhookSampleEventTypeEnumProductUpdated            WebhookSampleEventTypeEnum = "PRODUCT_UPDATED"
	WebhookSampleEventTypeEnumProductDeleted            WebhookSampleEventTypeEnum = "PRODUCT_DELETED"
	WebhookSampleEventTypeEnumProductVariantCreated     WebhookSampleEventTypeEnum = "PRODUCT_VARIANT_CREATED"
	WebhookSampleEventTypeEnumProductVariantUpdated     WebhookSampleEventTypeEnum = "PRODUCT_VARIANT_UPDATED"
	WebhookSampleEventTypeEnumProductVariantDeleted     WebhookSampleEventTypeEnum = "PRODUCT_VARIANT_DELETED"
	WebhookSampleEventTypeEnumProductVariantOutOfStock  WebhookSampleEventTypeEnum = "PRODUCT_VARIANT_OUT_OF_STOCK"
	WebhookSampleEventTypeEnumProductVariantBackInStock WebhookSampleEventTypeEnum = "PRODUCT_VARIANT_BACK_IN_STOCK"
	WebhookSampleEventTypeEnumCheckoutCreated           WebhookSampleEventTypeEnum = "CHECKOUT_CREATED"
	WebhookSampleEventTypeEnumCheckoutUpdated           WebhookSampleEventTypeEnum = "CHECKOUT_UPDATED"
	WebhookSampleEventTypeEnumFulfillmentCreated        WebhookSampleEventTypeEnum = "FULFILLMENT_CREATED"
	WebhookSampleEventTypeEnumFulfillmentCanceled       WebhookSampleEventTypeEnum = "FULFILLMENT_CANCELED"
	WebhookSampleEventTypeEnumNotifyUser                WebhookSampleEventTypeEnum = "NOTIFY_USER"
	WebhookSampleEventTypeEnumPageCreated               WebhookSampleEventTypeEnum = "PAGE_CREATED"
	WebhookSampleEventTypeEnumPageUpdated               WebhookSampleEventTypeEnum = "PAGE_UPDATED"
	WebhookSampleEventTypeEnumPageDeleted               WebhookSampleEventTypeEnum = "PAGE_DELETED"
	WebhookSampleEventTypeEnumPaymentAuthorize          WebhookSampleEventTypeEnum = "PAYMENT_AUTHORIZE"
	WebhookSampleEventTypeEnumPaymentCapture            WebhookSampleEventTypeEnum = "PAYMENT_CAPTURE"
	WebhookSampleEventTypeEnumPaymentConfirm            WebhookSampleEventTypeEnum = "PAYMENT_CONFIRM"
	WebhookSampleEventTypeEnumPaymentListGateways       WebhookSampleEventTypeEnum = "PAYMENT_LIST_GATEWAYS"
	WebhookSampleEventTypeEnumPaymentProcess            WebhookSampleEventTypeEnum = "PAYMENT_PROCESS"
	WebhookSampleEventTypeEnumPaymentRefund             WebhookSampleEventTypeEnum = "PAYMENT_REFUND"
	WebhookSampleEventTypeEnumPaymentVoid               WebhookSampleEventTypeEnum = "PAYMENT_VOID"
	WebhookSampleEventTypeEnumTranslationCreated        WebhookSampleEventTypeEnum = "TRANSLATION_CREATED"
	WebhookSampleEventTypeEnumTranslationUpdated        WebhookSampleEventTypeEnum = "TRANSLATION_UPDATED"
)

type WeightUnitsEnum = measurement.WeightUnit
