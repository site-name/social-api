package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/samber/lo"
	"github.com/site-name/i18naddress"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
)

func (r *Resolver) AddressCreate(ctx context.Context, args struct {
	Input  AddressInput
	UserID string
}) (*AddressCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressUpdate(ctx context.Context, args struct {
	Id    string
	Input AddressInput
}) (*AddressUpdate, error) {
	// embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	// if err != nil {
	// 	return nil, err
	// }

	// // validate given id
	// if !model.IsValidId(args.Id) {
	// 	return nil, model.NewAppError("AddressUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid address id", http.StatusBadRequest)
	// }

	// // validate permission
	// if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageUsers) {
	// 	return nil, model.NewAppError("AddressUpdate", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
	// }

	// // find address with given id
	// address, appErr := embedCtx.App.Srv().AccountService().AddressById(args.Id)
	// if err != nil {
	// 	return nil, appErr
	// }

	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressDelete(ctx context.Context, args struct{ Id string }) (*AddressDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressSetDefault(ctx context.Context, args struct {
	AddressID string
	Type      AddressTypeEnum
	UserID    string
}) (*AddressSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AddressValidationRules(ctx context.Context, args struct {
	CountryCode CountryCode
	CountryArea *string
	City        *string
	CityArea    *string
}) (*AddressValidationData, error) {
	addressParam := &i18naddress.Params{
		CountryCode: args.CountryCode.String(),
	}
	if area := args.CountryArea; area != nil && *area != "" {
		addressParam.CountryArea = *area
	}
	if city := args.City; city != nil && *city != "" {
		addressParam.City = *city
	}
	if cArea := args.CityArea; cArea != nil && *cArea != "" {
		addressParam.CityArea = *cArea
	}
	validationRules, err := i18naddress.GetValidationRules(addressParam)

	if err != nil {
		return nil, model.NewAppError("AddressValidationRules", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "args"}, err.Error(), http.StatusBadRequest)
	}

	choicesToChoiceValues := func(choices [][2]string) []*ChoiceValue {
		return lo.Map(choices, func(item [2]string, _ int) *ChoiceValue {
			return &ChoiceValue{
				Raw:     &item[0],
				Verbose: &item[1],
			}
		})
	}

	return &AddressValidationData{
		CountryCode:        &validationRules.CountryCode,
		CountryName:        &validationRules.CountryName,
		AddressFormat:      &validationRules.AddressFormat,
		AddressLatinFormat: &validationRules.AddressLatinFormat,
		AllowedFields:      GetAllowedFieldsCamelCase(validationRules.AllowedFields),
		RequiredFields:     GetRequiredFieldsCamelCase(validationRules.RequiredFields),
		UpperFields:        GetUppserFieldsCamelCase(validationRules.UpperFields),
		CountryAreaType:    &validationRules.CountryAreaType,
		CityType:           &validationRules.CityType,
		CityAreaType:       &validationRules.CityAreaType,
		PostalCodeType:     &validationRules.PostalCodeType,
		PostalCodeExamples: validationRules.PostalCodeExamples,
		PostalCodePrefix:   &validationRules.PostalCodePrefix,
		CountryAreaChoices: choicesToChoiceValues(validationRules.CountryAreaChoices),
		CityChoices:        choicesToChoiceValues(validationRules.CityChoices),
		CityAreaChoices:    choicesToChoiceValues(validationRules.CityAreaChoices),
		PostalCodeMatchers: lo.Map(validationRules.PostalCodeMatchers, func(rg *regexp.Regexp, _ int) string { return rg.String() }),
	}, nil
}

func (r *Resolver) Address(ctx context.Context, args struct{ Id string }) (*Address, error) {
	panic(fmt.Errorf("not implemented"))
}
