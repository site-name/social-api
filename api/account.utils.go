package api

import (
	"strings"

	"github.com/sitename/sitename/modules/util"
)

func ToCamelCase(s string) string {
	var split = strings.Split(s, "_")
	if len(split) == 1 {
		return split[0]
	}

	var builder strings.Builder
	builder.WriteString(split[0])

	for _, item := range split[1:] {
		if item != "" {
			builder.WriteString(strings.Title(item))
			continue
		}
		builder.WriteByte('_')
	}

	return builder.String()
}

func ValidationFieldToCamelCase(name string) string {
	name = ToCamelCase(name)
	if name == "streetAddress" {
		return "streetAddress1"
	}
	return name
}

func GetRequiredFieldsCamelCase(requiredFields util.AnyArray[string]) util.AnyArray[string] {
	return requiredFields.
		Map(func(_ int, item string) string {
			return ValidationFieldToCamelCase(item)
		}).
		Dedup()
}

func GetUppserFieldsCamelCase(uppserFields util.AnyArray[string]) util.AnyArray[string] {
	return GetRequiredFieldsCamelCase(uppserFields)
}

func GetAllowedFieldsCamelCase(allowedFields util.AnyArray[string]) util.AnyArray[string] {
	fields := GetRequiredFieldsCamelCase(allowedFields)
	if fields.Contains("streetAddress1") {
		fields = append(fields, "streetAddress2")
	}
	return fields
}

// func userIsStaffOfShopLoader(ctx context.Context, userIDs []string) []*dataloader.Result[bool] {
// 	var (
// 		res      = make([]*dataloader.Result[bool], len(userIDs))
// 		staffMap map[string]struct{}
// 	)
// 	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

// 	staffs, appErr := embedCtx.App.Srv().ShopService().StaffsByOptions(&model.ShopStaffFilterOptions{})
// 	if appErr != nil {
// 		goto errorLabel
// 	}
// 	staffMap = lo.SliceToMap(staffs, func(st *model.ShopStaff) (string, struct{}) { return st.StaffID, struct{}{} })

// 	for idx, id := range userIDs {
// 		_, ok := staffMap[id]
// 		res[idx] = &dataloader.Result[bool]{Data: ok}
// 	}
// 	return res

// errorLabel:
// 	for idx := range userIDs {
// 		res[idx] = &dataloader.Result[bool]{Error: appErr}
// 	}
// 	return res
// }
