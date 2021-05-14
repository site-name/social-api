package graph

import (
	"context"

	"github.com/sitename/sitename/graph/model"
	dbmodel "github.com/sitename/sitename/model"
)

// take permissions list required by app's manifest
// if some of them are not supported, return false
// else true
func cleanManifestPermissions(requiredPermissions []string) bool {
	missingPermissions := []string{}

	// finds missions in requiredPermissions that are not supported by system
	for _, perm := range requiredPermissions {
		if dbmodel.SaleorPermissionEnumMap[perm] == "" {
			missingPermissions = append(missingPermissions, perm)
		}
	}

	if len(missingPermissions) > 0 {
		return false
	}
	return true
}

func AppFetchManifest(r *Resolver, ctx context.Context, manifestURL string) (*model.AppFetchManifest, error) {

	// check if given url is valid
	// if !dbmodel.IsValidHttpUrl(manifestURL) {
	// 	return &model.AppFetchManifest{
	// 		Errors: []model.AppError{
	// 			{
	// 				Message: dbmodel.NewString("Enter a valid URL."),
	// 				Code:    model.AppErrorCodeInvalidURLFormat,
	// 			},
	// 		},
	// 	}, nil
	// }

	// // try fetching given url
	// client := &http.Client{
	// 	Timeout: 10 * time.Second,
	// }
	// res, err := client.Get(manifestURL)
	// if err != nil {
	// 	return &model.AppFetchManifest{
	// 		Errors: []model.AppError{
	// 			{
	// 				Message: dbmodel.NewString("Unable to fetch manifest data."),
	// 				Code:    model.AppErrorCodeManifestURLCantConnect,
	// 			},
	// 		},
	// 	}, nil
	// }
	// defer res.Body.Close() // remember to close io.Reader instance

	// // try doing json unmarshaling data returned by request
	// var manifestData dbmodel.StringInterface
	// err = dbmodel.ModelFromJson(&manifestData, res.Body)
	// if err != nil {
	// 	return &model.AppFetchManifest{
	// 		Errors: []model.AppError{
	// 			{
	// 				Message: dbmodel.NewString("Incorrect structure of manifest."),
	// 				Code:    model.AppErrorCodeInvalidManifestFormat,
	// 			},
	// 		},
	// 	}, nil
	// }

	// if valid := cleanManifestPermissions(manifestData["permissions"].([]string)); !valid {
	// 	return &model.AppFetchManifest{
	// 		Errors: []model.AppError{
	// 			{
	// 				Message: dbmodel.NewString("Given permissions don't exist"),
	// 				Code:    model.AppErrorCodeInvalidPermission,
	// 			},
	// 		},
	// 	}, nil
	// }

	panic("not implemented") // TODO: fixme
}

func AppInstall(r *Resolver, ctx context.Context, input model.AppInstallInput) (*model.AppInstall, error) {
	panic("not implemented") // TODO: fixme

}
