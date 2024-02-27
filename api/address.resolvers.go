package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"
	"regexp"

	"github.com/samber/lo"
	"github.com/site-name/i18naddress"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

// NOTE: Authenticated users can create address for themself only
// NOTE: Refer to ./schemas/address.graphqls for details on directive used
func (r *Resolver) AddressCreate(ctx context.Context, args struct {
	Input AddressInput
}) (*AddressCreate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	// validate input
	appErr := args.Input.validate("AddressCreate")
	if appErr != nil {
		return nil, appErr
	}

	// create new address:
	var address model.Address
	args.Input.PatchAddress(&address)

	savedAddress, appErr := embedCtx.App.AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	// get current user
	currentUser, appErr := embedCtx.App.AccountService().UserById(ctx, currentSession.UserID)
	if appErr != nil {
		return nil, appErr
	}

	pluginManager := embedCtx.App.PluginService().GetPluginManager()
	finalAddress, appErr := pluginManager.ChangeUserAddress(*savedAddress, nil, currentUser)
	if appErr != nil {
		return nil, appErr
	}

	// add user-address relation:
	appErr = embedCtx.App.Srv().Store.User().AddRelations(nil, currentSession.UserID, []*model.Address{{Id: finalAddress.Id}}, false)
	if appErr != nil {
		return nil, appErr
	}

	return &AddressCreate{
		User:    SystemUserToGraphqlUser(currentUser),
		Address: SystemAddressToGraphqlAddress(finalAddress),
	}, nil
}

// NOTE: Users can update their addresses only
// NOTE: Refer to ./schemas/address.graphqls for details on directive used
func (r *Resolver) AddressUpdate(ctx context.Context, args struct {
	Id    UUID
	Input AddressInput
}) (*AddressUpdate, error) {
	appErr := args.Input.validate("AddressUpdate")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	// check if requester really owns address:
	addresses, appErr := embedCtx.App.AccountService().AddressesByUserId(currentSession.UserID)
	if appErr != nil {
		return nil, appErr
	}
	address, found := lo.Find(addresses, func(addr *model.Address) bool { return addr != nil && addr.ID == args.Id })
	if !found || address == nil {
		return nil, MakeUnauthorizedError("AddressUpdate")
	}

	args.Input.PatchAddress(address)

	updatedAddress, appErr := embedCtx.App.AccountService().UpsertAddress(nil, *address)
	if appErr != nil {
		return nil, appErr
	}

	user, appErr := embedCtx.App.AccountService().UserById(ctx, currentSession.UserID)
	if appErr != nil {
		return nil, appErr
	}

	pluginManager := embedCtx.App.PluginService().GetPluginManager()
	_, appErr = pluginManager.CustomerUpdated(*user)
	if appErr != nil {
		return nil, appErr
	}

	finalAddress, appErr := pluginManager.ChangeUserAddress(*updatedAddress, nil, user)
	if appErr != nil {
		return nil, appErr
	}

	return &AddressUpdate{
		Address: SystemAddressToGraphqlAddress(finalAddress),
		User:    SystemUserToGraphqlUser(user),
	}, nil
}

// NOTE: Refer to ./schemas/address.graphqls for details on directive used
func (r *Resolver) AddressDelete(ctx context.Context, args struct{ Id UUID }) (*AddressDelete, error) {
	// TODO: investigate if deleting an address affects other parts like shipping/billing address of orders
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	// delete relation between user and address, address stil exists
	appErr := embedCtx.App.Srv().Store.User().RemoveRelations(nil, currentSession.UserID, []*model.Address{{Id: args.Id}}, false)
	if appErr != nil {
		return nil, appErr
	}

	// get current user
	user, appErr := embedCtx.App.AccountService().UserById(ctx, currentSession.UserID)
	if appErr != nil {
		return nil, appErr
	}

	address, appErr := embedCtx.App.AccountService().AddressById(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}

	pluginManager := embedCtx.App.PluginService().GetPluginManager()
	_, appErr = pluginManager.CustomerUpdated(*user)
	if appErr != nil {
		return nil, appErr
	}

	return &AddressDelete{
		Address: SystemAddressToGraphqlAddress(address),
		User:    SystemUserToGraphqlUser(user),
	}, nil
}

// NOTE: Refer to ./schemas/address.graphqls for details on directive used
func (r *Resolver) AddressSetDefault(ctx context.Context, args struct {
	AddressID string
	Type      AddressTypeEnum
}) (*AddressSetDefault, error) {
	// validate params
	if !model_helper.IsValidId(args.AddressID) {
		return nil, model_helper.NewAppError("AddressSetDefault", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "addressId"}, "please provide valid address id", http.StatusBadRequest)
	}
	if !args.Type.IsValid() {
		return nil, model_helper.NewAppError("AddressSetDefault", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "address type"}, "please provide valid address type", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	// check if requester own address
	updatedUser, appErr := embedCtx.App.AccountService().UserSetDefaultAddress(currentSession.UserID, args.AddressID, args.Type)
	if appErr != nil {
		return nil, appErr
	}

	pluginManager := embedCtx.App.PluginService().GetPluginManager()
	_, appErr = pluginManager.CustomerUpdated(*updatedUser)
	if appErr != nil {
		return nil, appErr
	}

	return &AddressSetDefault{
		User: SystemUserToGraphqlUser(updatedUser),
	}, nil
}

func choicesToChoiceValues(choices [][2]string) []*ChoiceValue {
	return lo.Map(choices, func(item [2]string, _ int) *ChoiceValue {
		return &ChoiceValue{
			Raw:     &item[0],
			Verbose: &item[1],
		}
	})
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
		return nil, model_helper.NewAppError("AddressValidationRules", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "args"}, err.Error(), http.StatusBadRequest)
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

// NOTE: Refer to ./schemas/address.graphqls for details on directive used
func (r *Resolver) Address(ctx context.Context, args struct{ Id UUID }) (*Address, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	address, appErr := embedCtx.App.AccountService().AddressById(args.Id.String())
	if appErr != nil {
		return nil, appErr
	}
	return SystemAddressToGraphqlAddress(address), nil
}
