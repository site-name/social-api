package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/uber/jaeger-client-go/utils"
)

func (r *Resolver) CustomerCreate(ctx context.Context, args struct{ Input UserCreateInput }) (*CustomerCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CustomerUpdate(ctx context.Context, args struct {
	Id    string
	Input CustomerInput
}) (*CustomerUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CustomerDelete(ctx context.Context, args struct{ Id string }) (*CustomerDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CustomerBulkDelete(ctx context.Context, args struct{ Ids []string }) (*CustomerBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: graphql directive(s) validated before this. Refer to ./schemas/customer.graphqls for details
func (r *Resolver) Customers(ctx context.Context, args struct {
	Filter *CustomerFilterInput
	SortBy *UserSortingInput
	GraphqlParams
}) (*UserCountableConnection, error) {
	// validate params
	if args.Filter != nil {
		appErr := args.Filter.validate("Customers")
		if appErr != nil {
			return nil, appErr
		}
	}
	paginValues, appErr := args.GraphqlParams.Parse("Customers")
	if appErr != nil {
		return nil, appErr
	}

	// parsing
	conditions := squirrel.And{}
	userFilterOpts := model.UserFilterOptions{}

	if filter := args.Filter; filter != nil {
		if dateJoin := args.Filter.DateJoined; dateJoin != nil {
			gte, lte := dateJoin.Gte, dateJoin.Lte

			if gte != nil {
				conditions = append(conditions, squirrel.Expr(model.UserTableName+".CreateAt >= ?", utils.TimeToMicrosecondsSinceEpochInt64(gte.Time)))
			}
			if lte != nil {
				conditions = append(conditions, squirrel.Expr(model.UserTableName+".CreateAt <= ?", utils.TimeToMicrosecondsSinceEpochInt64(lte.Time)))
			}
		}

		if numOfOrders := args.Filter.NumberOfOrders; numOfOrders != nil {
			userFilterOpts.AnnotateOrderCount = true

			if gte := numOfOrders.Gte; gte != nil {
				conditions = append(conditions, squirrel.Expr(model.UserTableName+".OrderCount >= ?", *gte))
			}
			if lte := numOfOrders.Lte; lte != nil {
				conditions = append(conditions, squirrel.Expr(model.UserTableName+".OrderCount <= ?", *lte))
			}
		}

		if placedOrders := args.Filter.PlacedOrders; placedOrders != nil {
			if gte := placedOrders.Gte; gte != nil {

			}
			if lte := placedOrders.Lte; lte != nil {

			}
		}
	}
}
