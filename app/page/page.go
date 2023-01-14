/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package page

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
)

type ServicePage struct {
	srv *app.Server
}

func init() {
	app.RegisterPageService(func(s *app.Server) (sub_app_iface.PageService, error) {
		return &ServicePage{
			srv: s,
		}, nil
	})
}

func (s *ServicePage) FindPagesByOptions(options *model.PageFilterOptions) ([]*model.Page, *model.AppError) {
	pages, err := s.srv.Store.Page().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("FindPagesByOptions", "app.page.finding_pages_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return pages, nil
}
