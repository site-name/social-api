package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) AttributeCreate(ctx context.Context, args struct{ Input AttributeCreateInput }) (*AttributeCreate, error) {
	var permissionToCheck = model.PermissionManagePageTypesAndAttributes
	if args.Input.Type == AttributeTypeEnumProductType {
		permissionToCheck = model.PermissionManageProductTypesAndAttributes
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), permissionToCheck) {
		return nil, model.NewAppError("AttributeCreate", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
	}

	// clean input
	inputType := args.Input.InputType
	if inputType != nil &&
		*inputType == model.REFERENCE &&
		(args.Input.EntityType == nil || !args.Input.EntityType.IsValid()) {
		return nil, model.NewAppError("AttributeCreate", "api.attribute.entity_type_missing.app_error", nil, "entity type is required when REFERENCE input type is used", http.StatusBadRequest)
	}

	// no need to initialize Slug here, since it will be done in PreSave() method
	attribute := &model.Attribute{
		Name:       args.Input.Name,
		Type:       string(args.Input.Type),
		Unit:       (*string)(unsafe.Pointer(args.Input.Unit)),
		EntityType: (*string)(unsafe.Pointer(args.Input.EntityType)),
	}
	if inputType != nil {
		attribute.InputType = *inputType
	}
	if v := args.Input.ValueRequired; v != nil {
		attribute.ValueRequired = *v
	}
	if v := args.Input.IsVariantOnly; v != nil {
		attribute.IsVariantOnly = *v
	}
	if v := args.Input.VisibleInStorefront; v != nil {
		attribute.VisibleInStoreFront = *v
	}
	if v := args.Input.FilterableInStorefront; v != nil {
		attribute.FilterableInStorefront = *v
	}
	if v := args.Input.FilterableInDashboard; v != nil {
		attribute.FilterableInDashboard = *v
	}
	if v := args.Input.StorefrontSearchPosition; v != nil {
		attribute.StorefrontSearchPosition = int(*v)
	}
	if v := args.Input.AvailableInGrid; v != nil {
		attribute.AvailableInGrid = *v
	}

	// clean attribute
	appErr := cleanAttributeSettings(attribute, &args.Input)
	if appErr != nil {
		return nil, appErr
	}

	// clean values
	validatedAttributeValues, appErr := (&AttributeMixin[*AttributeCreateInput]{
		ATTRIBUTE_VALUES_FIELD: "values",
		srv:                    embedCtx.App.Srv(),
	}).cleanValues(&args.Input, attribute)
	if appErr != nil {
		return nil, appErr
	}

	// save
	savedAttr, appErr := embedCtx.App.Srv().AttributeService().UpsertAttribute(attribute)
	if appErr != nil {
		return nil, appErr
	}

	// save attribute values
	_, appErr = embedCtx.App.Srv().AttributeService().BulkUpsertAttributeValue(nil, validatedAttributeValues)
	if appErr != nil {
		return nil, appErr
	}

	return &AttributeCreate{
		Attribute: SystemAttributeToGraphqlAttribute(savedAttr),
	}, nil
}

func (r *Resolver) AttributeDelete(ctx context.Context, args struct{ Id string }) (*AttributeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeUpdate(ctx context.Context, args struct {
	Id    string
	Input AttributeUpdateInput
}) (*AttributeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeTranslate(ctx context.Context, args struct {
	Id           string
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeBulkDelete(ctx context.Context, args struct{ Ids []string }) (*AttributeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueBulkDelete(ctx context.Context, args struct{ Ids []string }) (*AttributeValueBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueCreate(ctx context.Context, args struct {
	AttributeID string
	Input       AttributeValueCreateInput
}) (*AttributeValueCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueDelete(ctx context.Context, args struct{ Id string }) (*AttributeValueDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueUpdate(ctx context.Context, args struct {
	Id    string
	Input AttributeValueUpdateInput
}) (*AttributeValueUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueTranslate(ctx context.Context, args struct {
	Id           string
	Input        AttributeValueTranslationInput
	LanguageCode LanguageCodeEnum
}) (*AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeReorderValues(ctx context.Context, args struct {
	AttributeID string
	Moves       []*ReorderInput
}) (*AttributeReorderValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Attributes(ctx context.Context, args struct {
	Filter     *AttributeFilterInput
	SortBy     *AttributeSortingInput
	ChanelSlug *string
	GraphqlParams
}) (*AttributeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Attribute(ctx context.Context, args struct {
	Id   *string
	Slug *string
}) (*Attribute, error) {
	panic(fmt.Errorf("not implemented"))
}
