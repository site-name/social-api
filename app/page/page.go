/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package page

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

type ServicePage struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Page = &ServicePage{s}
		return nil
	})
}

func (s *ServicePage) FindPagesByOptions(options model_helper.PageFilterOptions) (model.PageSlice, *model_helper.AppError) {
	pages, err := s.srv.Store.Page().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("FindPagesByOptions", "app.page.finding_pages_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return pages, nil
}
