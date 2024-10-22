package api

import (
	"context"
	"encoding/json"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/web"
)

type graphQLInput struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName"`
	Variables     map[string]any `json:"variables"`
}

type Resolver struct{}

func (api *API) InitGraphql() error {
	schemaString, err := constructSchema()
	if err != nil {
		return err
	}

	opts := []graphql.SchemaOpt{
		graphql.UseFieldResolvers(),
		graphql.Logger(slog.NewGraphQLLogger(api.srv.Log)),
		graphql.MaxParallelism(200),
		graphql.MaxDepth(4),
		graphql.UseStringDescriptions(),
		graphql.Directives(
			&HasRolesDirective{},
			&HasRoleAnyDirective{},
			&HasPermissionsDirective{},
			&AuthenticatedDirective{},
			&HasPermissionAnyDirective{},
		),
		graphql.RestrictIntrospection(func(_ context.Context) bool {
			return model_helper.BuildNumber != "dev"
		}),
	}

	api.schema, err = graphql.ParseSchema(schemaString, &Resolver{}, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to parse graphql schema")
	}

	api.Router.Handle("/graphql", api.APIHandler(graphiQL)).Methods(http.MethodGet)
	api.Router.Handle("/graphql", api.APIHandler(api.graphql)).Methods(http.MethodPost)
	return nil
}

func (api *API) graphql(c *web.Context, w http.ResponseWriter, r *http.Request) {
	var response *graphql.Response
	defer func() {
		if response != nil {
			if err := json.NewEncoder(w).Encode(response); err != nil {
				c.Logger.Warn("Error while writing response", slog.Err(err))
			}
		}
	}()

	// Limit bodies to 100KiB.
	// We need to enforce a lower limit than the file upload size,
	// to prevent the library doing unnecessary parsing.
	r.Body = http.MaxBytesReader(w, r.Body, 102400)

	var params graphQLInput

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		err2 := gqlerrors.Errorf("invalid request body: %v", err)
		response = &graphql.Response{Errors: []*gqlerrors.QueryError{err2}}
		return
	}

	if params.OperationName == "" {
		err2 := gqlerrors.Errorf("operation name not passed")
		response = &graphql.Response{Errors: []*gqlerrors.QueryError{err2}}
		return
	}

	c.GraphQLOperationName = params.OperationName

	c.CurrentChannelID = r.URL.Query().Get("channel_id")
	if c.CurrentChannelID != "" && !model_helper.IsValidId(c.CurrentChannelID) {
		err2 := gqlerrors.Errorf("invalid request query param channel_id")
		response = &graphql.Response{Errors: []*gqlerrors.QueryError{err2}}
		return
	}

	reqCtx := r.Context()
	reqCtx = context.WithValue(reqCtx, WebCtx, c)

	response = api.schema.Exec(reqCtx, params.Query, params.OperationName, params.Variables)

	if len(response.Errors) > 0 {
		logFunc := slog.Error
		errorStrings := ", "

		for _, err := range response.Errors {
			if err != nil && err.Err != nil {
				if appErr, ok := err.Err.(*model_helper.AppError); ok && appErr.StatusCode < http.StatusInternalServerError {
					logFunc = slog.Debug
				}
				errorStrings += err.Error() + ", "
			}
		}

		logFunc("Error executing request", slog.String("operation", params.OperationName), slog.String("errors", errorStrings[1:]))
	}
}

func graphiQL(c *web.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(newGraphqlPlayground)
}

var newGraphqlPlayground = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<title>GraphiQL editor | Sitename</title>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css" integrity="sha256-gSgd+on4bTXigueyd/NSRNAy4cBY42RAVNaXnQDjOW8=" crossorigin="anonymous"/>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/es6-promise/4.1.1/es6-promise.auto.min.js" integrity="sha256-OI3N9zCKabDov2rZFzl8lJUXCcP7EmsGcGoP6DMXQCo=" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js" integrity="sha256-aB35laj7IZhLTx58xw/Gm1EKOoJJKZt6RY+bH1ReHxs=" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js" integrity="sha256-wouRkivKKXA3y6AuyFwcDcF50alCNV8LbghfYCH6Z98=" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js" integrity="sha256-9hrJxD4IQsWHdNpzLkJKYGiY/SEZFJJSUqyeZPNKd8g=" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js" integrity="sha256-oeWyQyKKUurcnbFRsfeSgrdOpXXiRYopnPjTVZ+6UmI=" crossorigin="anonymous"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/graphql", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
					headers: {
						'X-Requested-With': 'XMLHttpRequest'
					}
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
