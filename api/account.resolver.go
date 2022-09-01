package api

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
	"github.com/sitename/sitename/model/account"
)

func (r *resolver) Address(ctx context.Context, args struct{ Id string }) (*gqlmodel.Address, error) {
	fmt.Println("------------id", args.Id)

	return gqlmodel.SystemAddressToGraphqlAddress(&account.Address{
		Id:        "dfdfdf",
		FirstName: "le minh son",
		LastName:  "Vu Thi Quyet",
	}), nil
}
