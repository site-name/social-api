package api

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

type Page struct {
	ID              string          `json:"id"`
	SeoTitle        *string         `json:"seoTitle"`
	SeoDescription  *string         `json:"seoDescription"`
	Title           string          `json:"title"`
	Content         JSONString      `json:"content"`
	PublicationDate *Date           `json:"publicationDate"`
	IsPublished     bool            `json:"isPublished"`
	Slug            string          `json:"slug"`
	Created         DateTime        `json:"created"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`

	// PageType        *PageType            `json:"pageType"`
	// Translation     *PageTranslation     `json:"translation"`
	// Attributes      []*SelectedAttribute `json:"attributes"`
}

func systemPageToGraphqlPage(p *model.Page) *Page {
	if p == nil {
		return nil
	}

	res := &Page{
		ID:              p.Id,
		SeoTitle:        &p.SeoTitle,
		SeoDescription:  &p.SeoDescription,
		Title:           p.Title,
		Slug:            p.Slug,
		Content:         make(map[string]any),
		IsPublished:     p.IsPublished,
		Metadata:        MetadataToSlice(p.Metadata),
		PrivateMetadata: MetadataToSlice(p.PrivateMetadata),
		Created:         DateTime{util.TimeFromMillis(p.CreateAt)},
	}
	if p.PublicationDate != nil {
		res.PublicationDate = &Date{DateTime{*p.PublicationDate}}
	}
	if p.Content != nil && len(p.Content) > 0 {
		res.Content = JSONString(p.Content)
	}
	return res
}

func (p *Page) PageType(ctx context.Context) (*PageType, error) {
	panic("not implemented")
}

func (p *Page) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*PageTranslation, error) {
	panic("not implemented")
}

func (p *Page) Attributes(ctx context.Context) ([]*SelectedAttribute, error) {
	panic("not implemented")
}

func pageByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Page] {
	res := make([]*dataloader.Result[*model.Page], len(ids))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	pages, appErr := embedCtx.App.Srv().PageService().FindPagesByOptions(&model.PageFilterOptions{
		Conditions: squirrel.Eq{model.PageTableName + ".Id": ids},
	})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.Page]{Error: appErr}
		}
		return res
	}

	pageMap := lo.SliceToMap(pages, func(p *model.Page) (string, *model.Page) { return p.Id, p })
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Page]{Data: pageMap[id]}
	}
	return res

}
