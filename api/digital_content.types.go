package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
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
	contentURLs, err := DigitalContentUrlsByDigitalContentIDLoader.Load(ctx, d.ID)()
	if err != nil {
		return nil, err
	}
	return DataloaderResultMap(contentURLs, systemDigitalContentURLToGraphqlDigitalContentURL), nil
}

func (d *DigitalContent) ProductVariant(ctx context.Context) (*ProductVariant, error) {
	variant, err := ProductVariantByIdLoader.Load(ctx, d.d.ProductVariantID)()
	if err != nil {
		return nil, err
	}
	return SystemProductVariantToGraphqlProductVariant(variant), nil
}

func digitalContentByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.DigitalContent] {
	var (
		res        = make([]*dataloader.Result[*model.DigitalContent], len(ids))
		contents   []*model.DigitalContent
		appErr     *model.AppError
		contentMap = map[string]*model.DigitalContent{} // keys are digital content ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	contents, appErr = embedCtx.App.Srv().ProductService().DigitalContentsbyOptions(&model.DigitalContentFilterOption{
		Id: squirrel.Eq{store.DigitalContentTableName + ".Id": ids},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}
	for _, content := range contents {
		contentMap[content.Id] = content
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.DigitalContent]{Data: contentMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.DigitalContent]{Error: err}
	}
	return res
}

func digitalContentUrlsByDigitalContentIDLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.DigitalContentUrl] {
	var (
		res                  = make([]*dataloader.Result[[]*model.DigitalContentUrl], len(ids))
		digitalContentURLs   []*model.DigitalContentUrl
		appErr               *model.AppError
		digitalContentURLMap = map[string][]*model.DigitalContentUrl{} // keys are digital content ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	digitalContentURLs, appErr = embedCtx.App.Srv().ProductService().
		DigitalContentURLSByOptions(&model.DigitalContentUrlFilterOptions{
			ContentID: squirrel.Eq{store.DigitalContentURLTableName + ".ContentID": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, url := range digitalContentURLs {
		digitalContentURLMap[url.ContentID] = append(digitalContentURLMap[url.ContentID], url)
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[[]*model.DigitalContentUrl]{Data: digitalContentURLMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[[]*model.DigitalContentUrl]{Error: err}
	}
	return res
}

type DigitalContentURL struct {
	Created     DateTime `json:"created"`
	DownloadNum int32    `json:"downloadNum"`
	ID          string   `json:"id"`
	Token       string   `json:"token"`
	u           *model.DigitalContentUrl

	// URL         *string  `json:"url"`
	// Content     *DigitalContent `json:"content"`
}

func systemDigitalContentURLToGraphqlDigitalContentURL(u *model.DigitalContentUrl) *DigitalContentURL {
	if u == nil {
		return nil
	}

	return &DigitalContentURL{
		ID:          u.Id,
		Token:       u.Token,
		DownloadNum: int32(u.DownloadNum),
		Created:     DateTime{util.TimeFromMillis(u.CreateAt)},
		u:           u,
	}
}

func (d *DigitalContentURL) Content(ctx context.Context) (*DigitalContent, error) {
	content, err := DigitalContentByIdLoader.Load(ctx, d.u.ContentID)()
	if err != nil {
		return nil, err
	}
	return systemDigitalContentToGraphqlDigitalContent(content), nil
}

func (d *DigitalContentURL) URL(ctx context.Context) (*string, error) {
	panic("not implemented")
}