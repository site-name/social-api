package gqlmodel

import "github.com/sitename/sitename/model/menu"

type Menu struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Slug            string          `json:"slug"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	ItemIDs         []string        `json:"items"` // []*MenuItem
}

func (Menu) IsNode()               {}
func (Menu) IsObjectWithMetadata() {}

type MenuItem struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Level           int             `json:"level"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	URL             *string         `json:"url"`
	CollectionID    *string         `json:"collection"`  // *Collection
	MenuID          *string         `json:"menu"`        // *Menu
	ParentID        *string         `json:"parent"`      // *MenuItem
	CategoryID      *string         `json:"category"`    // *Category
	PageID          *string         `json:"page"`        // *Page
	ChildrenIDs     []string        `json:"children"`    // []*MenuItem
	TranslationID   *string         `json:"translation"` // *MenuItemTranslation
}

func (MenuItem) IsNode()               {}
func (MenuItem) IsObjectWithMetadata() {}

// DatabaseMenuToGraphqlMenu convert system menu into graphql menu
func DatabaseMenuToGraphqlMenu(m *menu.Menu) *Menu {
	return &Menu{
		ID:              m.Id,
		Name:            m.Name,
		Slug:            m.Slug,
		PrivateMetadata: MapToGraphqlMetaDataItems(m.PrivateMetadata),
		Metadata:        MapToGraphqlMetaDataItems(m.Metadata),
	}
}

// DatabaseMenuItemToGraphqlMenuItem converts menu item to graphql menu item
// func DatabaseMenuItemToGraphqlMenuItem(m *menu.MenuItem) *MenuItem {
// 	return &MenuItem{
// 		ID:   m.Id,
// 		Name: m.Name,
// 		// Level:           m.,
// 		PrivateMetadata: MapToGraphqlMetaDataItems(m.PrivateMetadata),
// 		Metadata:        MapToGraphqlMetaDataItems(m.Metadata),
// 		URL:             m.Url,
// 		CategoryID:      m.CategoryID,
// 	}
// }
