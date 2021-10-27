package web

import (
	"context"
	"net/http"

	gqlgenGraphql "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sitename/sitename/graphql"
	"github.com/sitename/sitename/graphql/dataloaders"
	"github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/web/shared"
)

const (
	graphqlApi        = "/api/graphql"
	graphqlPlayground = "/playground"
)

// InitGraphql registers graphql playground and graphql api endpoint routes
func (web *Web) InitGraphql() {
	config := generated.Config{
		Resolvers: &graphql.Resolver{
			AppIface: web.app,
		},
	}

	config.Directives.Authenticated = func(ctx context.Context, obj interface{}, next gqlgenGraphql.Resolver, _ bool) (res interface{}, err error) {
		_, appErr := graphql.CheckUserAuthenticated("Directives.Authenticated", ctx)
		if appErr != nil {
			return nil, appErr
		}

		return next(ctx)
	}
	config.Directives.HasPermissions = func(ctx context.Context, obj interface{}, next gqlgenGraphql.Resolver, permissions []gqlmodel.PermissionEnum) (res interface{}, err error) {
		session, appErr := graphql.CheckUserAuthenticated("Directives.Authenticated", ctx)
		if appErr != nil {
			return nil, appErr
		}

		if perms := gqlmodel.SaleorGraphqlPermissionsToSystemPermission(permissions...); !web.app.Srv().AccountService().SessionHasPermissionToAll(session, perms...) {
			return nil, web.app.Srv().AccountService().MakePermissionError(session, perms...)
		}

		return next(ctx)
	}

	graphqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(config))
	playgroundHandler := playground.Handler("Sitename", graphqlApi)

	web.MainRouter.Handle(graphqlPlayground, web.NewHandler(commonGraphHandler(playgroundHandler))).Methods(http.MethodGet)
	web.MainRouter.Handle(graphqlApi, web.NewHandler(commonGraphHandler(graphqlServer))).Methods(http.MethodPost)
}

// commonGraphHandler is used for both graphql playground/api
func commonGraphHandler(handler http.Handler) func(c *shared.Context, w http.ResponseWriter, r *http.Request) {
	return func(c *shared.Context, w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r.WithContext(
			context.WithValue(
				context.WithValue(r.Context(), shared.APIContextKey, c),
				dataloaders.DataloaderContextKey,
				dataloaders.NewLoaders(c.App),
			),
		))
	}
}
