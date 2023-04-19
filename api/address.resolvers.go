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
	"github.com/sitename/sitename/web"
)

func (r *Resolver) AddressCreate(ctx context.Context, args struct {
	Input AddressInput
}) (*AddressCreate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}
	currentSession := embedCtx.AppContext.Session()

	// valudate input
	appErr := args.Input.Validate()
	if appErr != nil {
		return nil, appErr
	}

	// create new address:
	address := new(model.Address)
	args.Input.PatchAddress(address)

	savedAddress, appErr := r.srv.AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	// add user-address relation:
	_, appErr = r.srv.AccountService().AddUserAddress(&model.UserAddress{
		UserID:    currentSession.UserId,
		AddressID: savedAddress.Id,
	})
	if appErr != nil {
		return nil, appErr
	}

	panic("not implemented") // TODO : finish plugin feature

	return &AddressCreate{
		User:    &User{ID: currentSession.UserId},
		Address: SystemAddressToGraphqlAddress(savedAddress),
	}, nil
}

func (r *Resolver) AddressUpdate(ctx context.Context, args struct {
	Id    string
	Input AddressInput
}) (*AddressUpdate, error) {
	// requester can update address only if he is owner of that address
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// validate given id
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AddressUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid address id", http.StatusBadRequest)
	}

	appErr := args.Input.Validate()
	if appErr != nil {
		return nil, appErr
	}

	// check if requester really owns address:
	addresses, appErr := r.srv.AccountService().AddressesByUserId(embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	address, found := lo.Find(addresses, func(addr *model.Address) bool { return addr.Id == args.Id })
	if !found || address == nil {
		return nil, MakeUnauthorizedError("AddressUpdate")
	}

	args.Input.PatchAddress(address)

	updatedAddress, appErr := r.srv.AccountService().UpsertAddress(nil, address)
	if appErr != nil {
		return nil, appErr
	}

	panic(fmt.Errorf("not implemented")) // TODO: finish plugin feature

	return &AddressUpdate{
		Address: SystemAddressToGraphqlAddress(updatedAddress),
		User:    &User{ID: embedCtx.AppContext.Session().UserId},
	}, nil
}

func (r *Resolver) AddressDelete(ctx context.Context, args struct{ Id string }) (*AddressDelete, error) {
	// requester can delete address only if he is owner of given address
	// TODO: investigate if deleting an address affects other parts like shipping/billing address of orders

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// validate id input
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("AddressDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid address id", http.StatusBadRequest)
	}

	// delete relation between user and address, address stil exists
	appErr := embedCtx.App.Srv().AccountService().AddressDeleteForUser(embedCtx.AppContext.Session().UserId, args.Id)
	if appErr != nil {
		return nil, appErr
	}

	panic(fmt.Errorf("not implemented"))

	return &AddressDelete{
		Address: &Address{model.Address{Id: args.Id}},
		User:    &User{ID: embedCtx.AppContext.Session().UserId},
	}, nil
}

func (r *Resolver) AddressSetDefault(ctx context.Context, args struct {
	AddressID string
	Type      model.AddressTypeEnum
}) (*AddressSetDefault, error) {
	// requester can update default billing/shipping address of himself only
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// validate address id
	if !model.IsValidId(args.AddressID) {
		return nil, model.NewAppError("AddressSetDefault", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "addressId"}, "please provide valid address id", http.StatusBadRequest)
	}
	if !args.Type.IsValid() {
		return nil, model.NewAppError("AddressSetDefault", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "address type"}, "please provide valid address type", http.StatusBadRequest)
	}

	// check if requester own address
	updatedUser, appErr := embedCtx.App.Srv().AccountService().UserSetDefaultAddress(embedCtx.AppContext.Session().UserId, args.AddressID, args.Type)
	if appErr != nil {
		return nil, appErr
	}

	panic("not implemented") // TODO: finish plugin apis

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
		return nil, model.NewAppError("AddressValidationRules", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "args"}, err.Error(), http.StatusBadRequest)
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
	// +) requester can see address if he is owner of that address OR
	// +) requester is staff of currentshop, this address is shipping/billing address of an order belongs to current shop
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("Address", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid address id", http.StatusBadRequest)
	}

	address, appErr := r.srv.AccountService().AddressById(args.Id)
	if appErr != nil {
		return nil, appErr
	}
	return SystemAddressToGraphqlAddress(address), nil
}
