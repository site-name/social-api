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
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web/shared"
)

const (
	graphqlApi            = "/api/graphql"
	graphqlPlayground     = "/playground"
	userUnauthenticatedId = "graphql.Authenticated.user_unauthenticated.app_error"
)

// InitGraphql registers graphql playground and graphql api endpoint routes
func (web *Web) InitGraphql() {
	config := generated.Config{
		Resolvers: &graphql.Resolver{
			AppIface: web.app,
		},
	}

	// Authenticated directive makes sure requesting user is authenticated before letting him/her proceed
	config.Directives.Authenticated = func(ctx context.Context, obj interface{}, next gqlgenGraphql.Resolver, _ bool) (res interface{}, err error) {
		embededContext := ctx.Value(shared.APIContextKey).(*shared.Context)
		if session := embededContext.AppContext.Session(); session == nil {
			return nil, model.NewAppError("Directives.Authenticated", userUnauthenticatedId, nil, "", http.StatusForbidden)
		}

		return next(ctx)
	}
	// HasPermissions directive makes sure requesting user is authenticated and has required permission(s) to proceed
	config.Directives.HasPermissions = func(ctx context.Context, obj interface{}, next gqlgenGraphql.Resolver, permissions []gqlmodel.PermissionEnum) (res interface{}, err error) {
		// 1) check if user is authenticated:
		embededContext := ctx.Value(shared.APIContextKey).(*shared.Context)
		session := embededContext.AppContext.Session()
		if session == nil {
			return nil, model.NewAppError("Directives.HasPermissions", userUnauthenticatedId, nil, "", http.StatusForbidden)
		}

		// 2) check if user has required permission(s) to proceed
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
