package mutations

import (
	"context"
	"net/http"

	"github.com/sitename/sitename/graph/model"
	dbmodel "github.com/sitename/sitename/model"
)

func cleanManifestPermissions(requiredPermissions []string) []*dbmodel.Permission {
	missingPermissions := []string{}
	saleorPermissionEnumMap := make(map[string]string)
	for _, perm := range dbmodel.SaleorPermissionEnumList {
		saleorPermissionEnumMap[perm.Id] = perm.Name
	}

	// finds missions in requiredPermissions that are not supported by system
	for _, perm := range requiredPermissions {
		if saleorPermissionEnumMap[perm] == "" {
			missingPermissions = append(missingPermissions, perm)
		}
	}

	return
}

func AppFetchManifest(ctx context.Context, manifestURL string) (*model.AppFetchManifest, error) {

	// check if given url is valid
	if !dbmodel.IsValidHttpUrl(manifestURL) {
		return &model.AppFetchManifest{
			Errors: []model.AppError{
				{
					Message: dbmodel.NewString("Enter a valid URL."),
					Code:    model.AppErrorCodeInvalidURLFormat,
				},
			},
		}, nil
	}

	// try fetching given url
	res, err := http.DefaultClient.Get(manifestURL)
	if err != nil {
		return &model.AppFetchManifest{
			Errors: []model.AppError{
				{
					Message: dbmodel.NewString("Unable to fetch manifest data."),
					Code:    model.AppErrorCodeManifestURLCantConnect,
				},
			},
		}, nil
	}
	defer res.Body.Close() // remember to close io.Reader instance

	var manifestData dbmodel.StringInterface
	err = dbmodel.ModelFromJson(&manifestData, res.Body)
	if err != nil {
		return &model.AppFetchManifest{
			Errors: []model.AppError{
				{
					Message: dbmodel.NewString("Incorrect structure of manifest."),
					Code:    model.AppErrorCodeInvalidManifestFormat,
				},
			},
		}, nil
	}

	manifestData["permissions"] = cleanManifestPermissions(manifestData["permissions"].([]string))
}
